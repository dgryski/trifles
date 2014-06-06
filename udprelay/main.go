package main

import (
	"encoding/binary"
	"flag"
	"log"
	"math"
	"net"
	"net/http"
	_ "net/http/pprof"
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

	// start workers
	var workerchs []chan []byte
	for _, dst := range flag.Args() {
		w := make(chan []byte, 100000)
		workerchs = append(workerchs, w)
		go worker(w, dst)
	}

	pconn, e := net.ListenPacket("udp", ":"+strconv.Itoa(*port))
	if e != nil {
		log.Fatal("udp listen error:", e)
	}

	log.Println("udp server starting on port", *port)

	Arena.arena = slab.NewArena(4+math.MaxUint16, 16*1024*1024, 2, nil)

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

	for {
		b := Arena.Alloc(4 + math.MaxUint16)
		n, _, err := pconn.ReadFrom(b[4:])
		if err != nil {
			log.Println(err)
			Arena.DecRef(b)
			continue
		}
		binary.LittleEndian.PutUint32(b[:], uint32(n))
		pkt := b[:4+n] // trim

		go func(pkt []byte) {
			for _, ch := range workerchs {
				Arena.AddRef(pkt)
				ch <- pkt
			}
			Arena.DecRef(b)
		}(pkt)
	}
}

func worker(w chan []byte, dst string) {
	conn, err := net.Dial("tcp", dst)
	if err != nil {
		log.Println("unable to connect: ", conn)
	}

	var failedPackets [][]byte

	tick := time.Tick(1 * time.Second)

	for {
		var pkt []byte
		select {
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
				}
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
