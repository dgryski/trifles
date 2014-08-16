package main

/*

sorted by misses

clockpro: 3.133172079s 10000000 total 2212461 misses
arc: 5.136880077s 10000000 total 2220016 misses
slru2080: 6.36564715s 10000000 total 2254569 misses
s4lru: 2.571416442s 10000000 total 2259629 misses
lfu: 5.194561433s 10000000 total 2326371 misses
slru5050: 6.266535668s 10000000 total 2360416 misses
clock: 2.497958238s 10000000 total 2587380 misses
lru: 5.020239605s 10000000 total 2643644 misses
random: 2.523026752s 10000000 total 2900727 misses

*/

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/calmh/lfucache"
	"github.com/davecheney/profile"
	"github.com/dgryski/go-arc"
	"github.com/dgryski/go-clockpro"
	"github.com/dgryski/go-s4lru"
	"github.com/dgryski/trifles/cachetest/clock"
	"github.com/dgryski/trifles/cachetest/random"
	"github.com/dgryski/trifles/cachetest/slru"
	"github.com/golang/groupcache/lru"
)

func main() {

	n := flag.Int("n", 1000, "cache size")
	alg := flag.String("alg", "", "algorithm")
	file := flag.String("f", "", "input file")
	cpuprofile := flag.Bool("cpuprofile", false, "cpuprofile")
	memprofile := flag.Bool("memprofile", false, "memprofile")

	flag.Parse()

	if *alg == "" {
		log.Fatalln("no algorithm provided (-alg)")
	}

	if *cpuprofile {
		defer profile.Start(profile.CPUProfile).Stop()
	}

	count := 0
	miss := 0

	t0 := time.Now()

	var f func(string) bool

	switch *alg {

	case "arc":

		cache := arc.New(*n)

		f = func(s string) bool {

			var miss bool

			cache.Get(s, func() interface{} {
				miss = true
				return s
			})

			return miss
		}

	case "random":

		cache := random.New(*n)

		f = func(s string) bool {
			if i := cache.Get(s); i == nil {
				cache.Set(s, s)
				return true
			}
			return false
		}

	case "lru":

		cache := lru.New(*n)

		f = func(s string) bool {
			if _, ok := cache.Get(s); !ok {
				cache.Add(s, s)
				return true
			}
			return false
		}

	case "lfu":

		cache := lfucache.New(*n)

		f = func(s string) bool {
			if _, ok := cache.Access(s); !ok {
				cache.Insert(s, s)
				return true
			}
			return false

		}

	case "clock":

		cache := clock.New(*n)

		f = func(s string) bool {
			if i := cache.Get(s); i == nil {
				cache.Set(s, s)
				return true
			}
			return false
		}

	case "clockpro":

		cache := clockpro.New(*n)

		f = func(s string) bool {
			if i := cache.Get(s); i == nil {
				cache.Set(s, s)
				return true
			}
			return false
		}

	case "slru":

		cache := slru.New(int(float64(*n)*0.2), int(float64(*n)*0.8))

		f = func(s string) bool {
			if i := cache.Get(s); i == nil {
				cache.Set(s, s)
				return true
			}
			return false
		}

	case "s4lru":

		cache := s4lru.New(*n)

		f = func(s string) bool {
			if _, ok := cache.Get(s); !ok {
				cache.Set(s, s)
				return true
			}
			return false

		}

	default:
		log.Fatalln("unknown algorithm")
	}

	var inputFile = os.Stdin
	if *file != "" {
		var err error
		inputFile, err = os.Open(*file)
		if err != nil {
			log.Fatalln(err)
		}
		defer inputFile.Close()
	}

	in := bufio.NewScanner(inputFile)
	for in.Scan() {
		if f(in.Text()) {
			miss++
		}
		count++
	}

	if *memprofile {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		e := json.NewEncoder(os.Stdout)
		e.Encode(m)
	}

	fmt.Printf("%s: %s %d total %d misses\n", *alg, time.Since(t0), count, miss)
}
