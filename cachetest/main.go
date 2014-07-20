package main

/*

slru2080: 6.36564715s 10000000 total 2254569 misses
s4lru: 3.176211561s 10000000 total 2259629 misses
lfu: 5.194561433s 10000000 total 2326371 misses
slru5050: 6.266535668s 10000000 total 2360416 misses
clock: 2.497958238s 10000000 total 2587380 misses
lru: 5.020239605s 10000000 total 2643644 misses
random: 2.523026752s 10000000 total 2900727 misses

*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/calmh/lfucache"
	"github.com/dgryski/go-s4lru"
	"github.com/dgryski/trifles/cachetest/clock"
	"github.com/dgryski/trifles/cachetest/random"
	"github.com/dgryski/trifles/cachetest/slru"
	"github.com/golang/groupcache/lru"
)

func main() {

	n := flag.Int("n", 1000, "cache size")
	alg := flag.String("alg", "", "algorithm")

	flag.Parse()

	count := 0
	miss := 0

	t0 := time.Now()

	var f func(string) bool

	switch *alg {

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

	}

	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		if f(in.Text()) {
			miss += 1
		}
		count += 1
	}

	fmt.Printf("%s: %s %d total %d misses\n", *alg, time.Since(t0), count, miss)
}
