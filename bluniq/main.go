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
	n := flag.Int("n", 10e6, "cardinality estimate of set")
	q := flag.Bool("q", false, "quiet")
	fprate := flag.Float64("fp", 0.00000001, "false positive rate")

	flag.Parse()

	file, err := os.Open(*f)
	if err != nil {
		log.Fatal(err)
	}

	var lines int
	b := bloomf.New(*n, *fprate, func(b []byte) uint64 { return siphash.Hash(0, 0, b) })
	scanner := bufio.NewScanner(file)

	var total int

	for scanner.Scan() {
		lines++
		l := scanner.Bytes()
		if !b.Lookup(l) {
			total++
			if !*q {
				os.Stdout.Write(l)
				os.Stdout.Write([]byte("\n"))
			}
		}
		b.Insert(l)
		if lines%(1<<20) == 0 {
			log.Println(lines)
		}
	}
	log.Printf("unique %d\n", total)
}
