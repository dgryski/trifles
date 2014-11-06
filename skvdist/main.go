package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dchest/siphash"
	"github.com/dgryski/go-onlinestats"
	"github.com/dgryski/go-shardedkv"
	"github.com/dgryski/go-shardedkv/choosers/chash"
	"github.com/dgryski/go-shardedkv/choosers/jump"
	"github.com/dgryski/go-shardedkv/choosers/ketama"
)

func worker(skv shardedkv.Chooser, filech chan string, donech chan map[string]int, verbose bool) {

	results := make(map[string]int)

	for file := range filech {

		var f io.Reader

		f, err := os.Open(file)
		if err != nil {
			log.Fatal("error loading input file", file, ":", err)
		}

		if strings.HasSuffix(file, ".gz") {
			f, _ = gzip.NewReader(f)
		}

		scan := bufio.NewScanner(f)
		for scan.Scan() {
			key := scan.Text()
			shard := skv.Choose(key)
			if verbose {
				fmt.Printf("%s %s\n", shard, key)
			}
			results[shard]++
		}
	}

	donech <- results
}

func main() {

	workers := flag.Int("w", 4, "workers")
	chooserType := flag.String("chooser", "jump", "shardedkv chooser")
	verbose := flag.Bool("v", false, "verbose")

	flag.Parse()

	filech := make(chan string)
	donech := make(chan map[string]int)

	var buckets []string
	for i := 1; i <= 32; i++ {
		buckets = append(buckets, fmt.Sprintf("shard%d", i))
	}

	var chooser shardedkv.Chooser

	switch *chooserType {
	case "jump":
		chooser = jump.New(func(b []byte) uint64 { return siphash.Hash(0, 0, b) })
	case "ketama":
		chooser = ketama.New()
	case "chash":
		// 5120*8 + 5120*(8+16)*2 bytes
		chooser = chash.New()
	default:
		log.Fatalf("unknown chooser: %s; known jump,ketama,chash", *chooserType)
	}

	chooser.SetBuckets(buckets)

	for i := 0; i < *workers; i++ {
		go worker(chooser, filech, donech, *verbose)
	}

	files, _ := filepath.Glob("data/*")

	for _, f := range files {
		filech <- f
	}

	close(filech)

	results := make(map[string]int)

	for i := 0; i < *workers; i++ {
		shards := <-donech
		for k, v := range shards {
			results[k] += v
		}
	}

	stats := onlinestats.NewRunning()

	fmt.Println("chooser", *chooserType)
	for _, b := range buckets {
		v := results[b]
		fmt.Printf("%-10s %d\n", b, v)
		stats.Push(float64(v))
	}

	fmt.Println("mean", stats.Mean())
	fmt.Println("stdev", stats.Stddev(), stats.Stddev()/stats.Mean())
}
