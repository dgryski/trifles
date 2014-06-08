package main

/*
TODO:
    actually write files on error
    better reconnect / error handling
    tcp listener
    expvar support
*/

import (
	"encoding/binary"
	"flag"
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

func main() {

	port := flag.Int("p", 12233, "udp listen port")
	debugPort := flag.Int("debugPort", 0, "debug port")
	arenaStats := flag.Bool("memstats", false, "dump arena stats")

	flag.Parse()

	if *debugPort != 0 {
		go http.ListenAndServe(":"+strconv.Itoa(*debugPort), nil)
	}

	quit := make(chan struct{})
	done := make(chan struct{})

	// start workers
	var workerchs []chan []byte
	for _, dst := range flag.Args() {
		w := make(chan []byte, 100000)
		workerchs = append(workerchs, w)
		go worker(w, dst, quit, done)
	}

	pconn, e := net.ListenPacket("udp", ":"+strconv.Itoa(*port))
	if e != nil {
		log.Fatal("udp listen error:", e)
	}

	log.Println("udp server starting on port", *port)

	Arena.arena = slab.NewArena(8192, 32*1024*1024, 2, nil)

	// TODO(dgryski): one we have expvar, export these too
	if *arenaStats {
		go func() {
			arenaStats := make(map[string]int64)

			for {
				time.Sleep(10 * time.Second)
				log.Println(Arena.Stats(arenaStats))
			}
		}()
	}

	var b [math.MaxUint16]byte

	go func() {
		for {
			n, _, err := pconn.ReadFrom(b[:])
			if err != nil {
				log.Println(err)
				continue
			}

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
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("Shutting down..")
	close(quit)
	for i := 0; i < len(workerchs); i++ {
		<-done
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
		Arena.DecRef(pkt)
	}
}

func dumpToFile(dst string, epoch int64, failedPackets [][]byte) {
	// nothing, for now
	log.Println("dumping ", len(failedPackets), "for", dst, epoch)
}
