package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/dchest/siphash"
	"github.com/dgryski/go-onlinestats"
	"github.com/dgryski/go-shardedkv"
	"github.com/dgryski/go-shardedkv/choosers/chash"
	"github.com/dgryski/go-shardedkv/choosers/jump"
	"github.com/dgryski/go-shardedkv/choosers/ketama"
)

func main() {

	chooserType := flag.String("chooser", "jump", "shardedkv chooser")
	expectedShard := flag.String("shard", "", "expected shard name")
	numShards := flag.Int("n", 16, "number of shards")
	fname := flag.String("f", "", "input file name")

	flag.Parse()

	var buckets []string
	for i := 1; i <= *numShards; i++ {
		buckets = append(buckets, fmt.Sprintf("shard%d", i))
	}

	var chooser shardedkv.Chooser

	switch *chooserType {
	case "jump":
		chooser = jump.New(func(b []byte) uint64 { return siphash.Hash(0, 0, b) })
	case "ketama":
		chooser = ketama.New()
	case "chash":
		chooser = chash.New()
	default:
		log.Fatalf("unknown chooser: %s; known jump,ketama,chash", *chooserType)
	}

	chooser.SetBuckets(buckets)

	results := make(map[string]int)

	var f io.Reader

	f, err := os.Open(*fname)
	if err != nil {
		log.Fatal("error loading input file", fname, ":", err)
	}

	if strings.HasSuffix(*fname, ".gz") {
		f, _ = gzip.NewReader(f)
	}

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		shard := chooser.Choose(scan.Text())
		if *expectedShard != "" && shard != *expectedShard {
			fmt.Printf("wrong shard %s %+v\n", shard, scan.Text())
		}
		results[shard]++
	}

	stats := onlinestats.NewRunning()

	fmt.Println("chooser", *chooserType)
	for _, b := range buckets {
		v := results[b]
		fmt.Printf("%-10s %d\n", b, v)
		stats.Push(float64(v))
	}

	if *expectedShard == "" {
		fmt.Println("mean", stats.Mean())
		fmt.Println("stdev", stats.Stddev(), stats.Stddev()/stats.Mean())
	}
}
