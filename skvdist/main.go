package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dchest/siphash"
	"github.com/dgryski/go-onlinestats"
	"github.com/dgryski/go-shardedkv"
	"github.com/dgryski/go-shardedkv/choosers/chash"
	"github.com/dgryski/go-shardedkv/choosers/jump"
	"github.com/dgryski/go-shardedkv/choosers/ketama"
	"github.com/dgryski/go-shardedkv/choosers/mpc"
)

func main() {

	numBuckets := flag.Int("b", 8, "buckets")
	chooserType := flag.String("chooser", "jump", "shardedkv chooser")
	verbose := flag.Bool("v", false, "verbose")
	replicas := flag.Int("r", 0, "replicas")

	flag.Parse()

	var buckets []string
	for i := 0; i < *numBuckets; i++ {
		//  	buckets = append(buckets, fmt.Sprintf("%016x", siphash.Hash(0x0ddc0ffeebadf00d, 0xcafebabedeadbeef, []byte(fmt.Sprintf("192.168.0.%d", 10+i)))))
		// buckets = append(buckets, fmt.Sprintf("%016x", siphash.Hash(0, 0, []byte(fmt.Sprintf("192.168.0.%d", 10+i)))))
		buckets = append(buckets, fmt.Sprintf("192.168.0.%d", 10+i))
	}

	var chooser shardedkv.Chooser

	switch *chooserType {
	case "jump":
		chooser = jump.New(func(b []byte) uint64 { return siphash.Hash(0, 0, b) })
	case "ketama":
		chooser = ketama.New()
	case "mpc":
		chooser = mpc.New(func(b []byte, seed uint64) uint64 { return siphash.Hash(seed, 0, b) }, [2]uint64{0x0ddc0ffeebadf00d, 0xcafebabedeadbeef}, 21)
	case "chash":
		// 5120*8 + 5120*(8+16)*2 bytes
		chooser = chash.New()
	default:
		log.Fatalf("unknown chooser: %s; known jump,ketama,chash,mpc", *chooserType)
	}

	chooser.SetBuckets(buckets)

	results := make(map[string]int)

	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		key := scan.Text()
		shard := chooser.Choose(key)
		if *verbose {
			fmt.Printf("%s %s\n", shard, key)
		}
		results[shard]++
	}

	stats := onlinestats.NewRunning()

	min := int(^uint(0) / 2)
	max := 0
	fmt.Println("chooser", *chooserType)
	for _, b := range buckets {
		v := results[b]
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		fmt.Printf("%-10s %d\n", b, v)
		stats.Push(float64(v))
	}

	fmt.Printf("min=%d max=%d\n", min, max)
	fmt.Printf("mean %.2f\n", stats.Mean())
	fmt.Printf("stdev %.2f %.2f%%\n", stats.Stddev(), 100*stats.Stddev()/stats.Mean())
	fmt.Printf("peak/mean %.2f\n", float64(max)/stats.Mean())
}
