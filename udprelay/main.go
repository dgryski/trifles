package main

/*
TODO:
    actually write files on error
    better reconnect / error handling
*/

import (
	"encoding/binary"
	"expvar"
	"flag"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/couchbaselabs/go-slab"
)

type lockedArena struct {
	arena *slab.Arena
	sync.Mutex
}

func (a *lockedArena) Alloc(size int) []byte {
	a.Lock()
	b := a.arena.Alloc(size)
	a.Unlock()
	return b
}

func (a *lockedArena) AddRef(b []byte) {
	a.Lock()
	a.arena.AddRef(b)
	a.Unlock()
}

func (a *lockedArena) DecRef(b []byte) {
	a.Lock()
	a.arena.DecRef(b)
	a.Unlock()
}

func (a *lockedArena) Stats(m map[string]int64) map[string]int64 {
	a.Lock()
	m2 := a.arena.Stats(m)
	a.Unlock()
	return m2
}

var Arena lockedArena

// grouped expvars for /debug/vars and graphite
var Metrics = struct {
	PacketsIn  *expvar.Int
	PacketsOut *expvar.Int
	BytesIn    *expvar.Int
	BytesOut   *expvar.Int
	Failed     *expvar.Int
}{
	PacketsIn:  expvar.NewInt("packetsIn"),
	PacketsOut: expvar.NewInt("packetsOut"),
	BytesIn:    expvar.NewInt("bytesIn"),
	BytesOut:   expvar.NewInt("bytesOut"),
	Failed:     expvar.NewInt("failed"),
}

func main() {

	proto := flag.String("proto", "udp", "listen protocol")
	port := flag.Int("p", 12233, "listen port")
	debugPort := flag.Int("debugPort", 8080, "debug port")

	flag.Parse()

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(*debugPort), nil)
		if err != nil {
			log.Fatal("Error from ListenAndServe:", err)
		}
	}()

	quit := make(chan struct{})
	done := make(chan struct{})

	// start workers
	var workerchs []chan []byte
	for _, dst := range flag.Args() {
		w := make(chan []byte, 100000)
		workerchs = append(workerchs, w)
		go worker(w, dst, quit, done)
	}

	Arena.arena = slab.NewArena(8192, 32*1024*1024, 2, nil)

	expvar.Publish("arenastats", expvar.Func(func() interface{} {
		m := make(map[string]int64)
		Arena.Stats(m)
		return m
	}))

	switch *proto {
	case "udp":
		go udpHandler(*port, workerchs)
	case "tcp":
		go tcpHandler(*port, workerchs)
	default:
		log.Fatal("unknown protocol", *proto)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("Shutting down..")
	close(quit)
	for i := 0; i < len(workerchs); i++ {
		<-done
	}
}

func udpHandler(port int, workerchs []chan []byte) {
	pconn, e := net.ListenPacket("udp", ":"+strconv.Itoa(port))
	if e != nil {
		log.Fatal("udp listen error:", e)
	}

	log.Println("udp server starting on port", port)

	var b [math.MaxUint16]byte
	for {
		n, _, err := pconn.ReadFrom(b[:])
		if err != nil {
			log.Println(err)
			continue
		}

		Metrics.PacketsIn.Add(1)

		if n == 0 {
			// ignore 0 byte packets
			log.Println("ignoring 0-byte packet")
			continue
		}

		Metrics.BytesIn.Add(int64(n))

		pkt := Arena.Alloc(4 + n)
		copy(pkt[4:], b[:n])

		binary.LittleEndian.PutUint32(pkt[:], uint32(n))

		for _, ch := range workerchs {
			Arena.AddRef(pkt)
			select {
			case ch <- pkt:
				// success
			default:
				// channel is full :(
				Arena.DecRef(pkt)
			}
		}
		Arena.DecRef(pkt)
	}
}

func tcpHandler(port int, workerchs []chan []byte) {
	ln, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	if e != nil {
		log.Fatal("tcp listen error:", e)
	}

	log.Println("tcp server starting on port", port)

	for {

		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func(conn net.Conn) {

			defer conn.Close()

			var b [math.MaxUint16]byte
			var bsz [4]byte
			for {

				n, err := conn.Read(bsz[:])
				if n != 4 || err != nil {
					log.Println(err)
					return
				}

				n = int(binary.LittleEndian.Uint16(bsz[:]))

				if n == 0 {
					// ignore 0 byte packets
					log.Println("ignoring 0-byte packet")
					continue
				}

				if n > math.MaxUint16 {
					log.Println("packet > math.MaxUint16: ", n)
					return
				}

				_, err = io.ReadFull(conn, b[:n])
				if err != nil {
					log.Println("short read on packet")
					return
				}

				// same as for udp

				Metrics.PacketsIn.Add(1)
				Metrics.BytesIn.Add(int64(n) + 4)

				pkt := Arena.Alloc(n + 4)
				copy(pkt[4:], b[:n])
				binary.LittleEndian.PutUint32(pkt[:], uint32(n))

				for _, ch := range workerchs {
					Arena.AddRef(pkt)
					select {
					case ch <- pkt:
						// success
					default:
						// channel is full :(
						Arena.DecRef(pkt)
					}
				}
				Arena.DecRef(pkt)
			}
		}(conn)
	}
}

func worker(w chan []byte, dst string, quit, done chan struct{}) {
	conn, err := net.Dial("tcp", dst)
	if err != nil {
		log.Println("unable to connect: ", conn)
	}

	var failedPackets [][]byte

	tick := time.Tick(1 * time.Second)

	for {
		var pkt []byte
		select {
		case <-quit:
			// shut down immediately or try to drain the work queue?
			if len(failedPackets) != 0 {
				dumpToFile(dst, time.Now().Unix(), failedPackets)
			}
			done <- struct{}{}
			return

		case epoch := <-tick:
			if len(failedPackets) != 0 {
				dumpToFile(dst, epoch.Unix(), failedPackets)

				for _, b := range failedPackets {
					Arena.DecRef(b)
				}

				failedPackets = failedPackets[:0]

				// try reconnecting
				// TODO(dgryski): add our good friend 'exponential backoff' ?
				if conn != nil {
					conn.Close()
					conn = nil
				}
			}

			if conn == nil {
				conn, err = net.Dial("tcp", dst)
				if err != nil {
					log.Println("unable to re-connect: ", dst, err)
				}
			}
			continue

		case pkt = <-w:
		}

		if conn == nil {
			failedPackets = append(failedPackets, pkt)
			continue
		}

		n, err := conn.Write(pkt)

		if err != nil {
			failedPackets = append(failedPackets, pkt)
			continue
		}

		if n != len(pkt) {
			failedPackets = append(failedPackets, pkt)
			continue
		}
		Metrics.PacketsOut.Add(1)
		Metrics.BytesOut.Add(int64(n))
		Arena.DecRef(pkt)
	}
}

func dumpToFile(dst string, epoch int64, failedPackets [][]byte) {
	Metrics.Failed.Add(int64(len(failedPackets)))

	// just log, for now
	log.Println("dumping ", len(failedPackets), "for", dst, epoch)
}
