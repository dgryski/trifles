package main

import (
	"bufio"
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/dgryski/go-onlinestats"
	"github.com/dgryski/go-shardedkv"
	"github.com/dgryski/go-shardedkv/choosers/chash"
	"github.com/dgryski/go-shardedkv/choosers/jump"
	"github.com/dgryski/go-shardedkv/choosers/ketama"
)

type replichooser interface {
	shardedkv.Chooser
	ChooseReplicas(key string, n int) []string
}

func main() {

	chooser := flag.String("chooser", "jump", "chooser to use")
	numBuckets := flag.Int("b", 8, "buckets")
	replicas := flag.Int("r", 1, "replicas")
	cpuprofile := flag.String("cpuprofile", "", "cpu profile")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var buckets []string
	for i := 0; i < *numBuckets; i++ {
		buckets = append(buckets, fmt.Sprintf("192.168.%d.%d", i, 10+2*i))
	}
	results := make(map[string]int)

	scan := bufio.NewScanner(os.Stdin)

	hasher := func(b []byte) uint64 { h := fnv.New64a(); h.Write(b); return h.Sum64() }

	var j replichooser

	switch *chooser {
	case "ketama":
		j = ketama.New()
	case "chash":
		j = chash.New()
	case "jump":
		j = jump.New(hasher)
	default:
		log.Fatal("unknown chooser", *chooser)
	}

	j.SetBuckets(buckets)

	var totalMetrics int
	t0 := time.Now()
	for scan.Scan() {
		key := scan.Text()
		totalMetrics++

		r := j.ChooseReplicas(key, *replicas)

		for _, rr := range r {
			results[rr]++
		}
	}
	t1 := time.Since(t0)

	stats := onlinestats.NewRunning()

	fmt.Printf("total metrics processed: %d, time spent ~%ds (%d/s)\n", totalMetrics, int(t1.Seconds()), int(float64(totalMetrics)/float64(t1.Seconds())))
	fmt.Printf("replication count: %d, server count: %d\n", *replicas, *numBuckets)
	fmt.Println("server distribution:")

	min := int(^uint(0) / 2)
	max := 0
	for _, b := range buckets {
		v := results[b]
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		fmt.Printf("- server %15s: %6d (%.2f%%)\n", b, v, 100*float64(v)/float64(totalMetrics))
		stats.Push(float64(v))
	}

	fmt.Printf("band: %d - %d (diff %d, %.1f%%), mean: %.2f, stddev: %.2f\n", min, max, (max - min), 100*float64(max-min)/float64(totalMetrics), stats.Mean(), stats.Stddev())
}
