package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/dchest/siphash"
	"github.com/dgryski/go-bloomf"
)

func main() {

	f := flag.String("f", "/dev/stdin", "input file")
	n := flag.Int("n", 100e6, "cardinality estimate of set")
	fprate := flag.Float64("fp", 0.00000001, "false positive rate")

	flag.Parse()

	file, err := os.Open(*f)
	if err != nil {
		log.Fatal(err)
	}

	var lines int
	b := bloomf.New(*n, *fprate, func(b []byte) uint64 { return siphash.Hash(0, 0, b) })
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines++
		if !b.Lookup(scanner.Bytes()) {
			os.Stdout.Write(scanner.Bytes())
			os.Stdout.Write([]byte("\n"))
		}
		b.Insert(scanner.Bytes())
		if lines%(1<<20) == 0 {
			log.Println(lines)
		}
	}
}
