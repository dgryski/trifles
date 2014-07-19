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
	"github.com/dgryski/trifles/cachetest/clock"
	"github.com/dgryski/trifles/cachetest/random"
	"github.com/dgryski/trifles/cachetest/s4lru"
	"github.com/dgryski/trifles/cachetest/slru"
	"github.com/golang/groupcache/lru"
)

func main() {

	alg := flag.String("alg", "", "algorithm")

	flag.Parse()

	count := 0
	miss := 0

	t0 := time.Now()

	switch *alg {

	case "random":

		cache := random.New(1000)

		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			domain := in.Text()
			if i := cache.Get(domain); i == nil {
				cache.Set(domain, domain)
				miss += 1
			}

			count += 1
		}

	case "lru":

		cache := lru.New(1000)

		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			domain := in.Text()
			if _, ok := cache.Get(domain); !ok {
				cache.Add(domain, domain)
				miss += 1
			}

			count += 1
		}

	case "lfu":

		cache := lfucache.New(1000)

		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			domain := in.Text()
			if _, ok := cache.Access(domain); !ok {
				cache.Insert(domain, domain)
				miss += 1
			}

			count += 1
		}

	case "clock":

		cache := clock.New(1000)

		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			domain := in.Text()
			if i := cache.Get(domain); i == nil {
				cache.Set(domain, domain)
				i = domain
				miss += 1
			}

			count += 1
		}

	case "slru":

		cache := slru.New(200, 800)

		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			domain := in.Text()
			if i := cache.Get(domain); i == nil {
				cache.Put(domain, domain)
				miss += 1
			}
			count += 1
		}

	case "s4lru":

		cache := s4lru.New(1000)

		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			domain := in.Text()
			if i := cache.Get(domain); i == nil {
				cache.Put(domain, domain)
				miss += 1
			}
			count += 1
		}
	}

	fmt.Printf("%s: %s %d total %d misses\n", *alg, time.Since(t0), count, miss)
}
