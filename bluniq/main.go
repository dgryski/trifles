package main

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"hash"
	"log"
	"os"

	"github.com/dchest/siphash"
	"github.com/dgryski/dgobloom"
)

type h32 struct {
	hash.Hash64
}

func (h h32) Sum32() uint32 {
	return uint32(h.Sum64())
}

func main() {

	f := flag.String("f", "", "input file")
	n := flag.Int("n", 100e6, "cardinality estimate of set")
	fprate := flag.Float64("fp", 0.00000001, "false positive rate")

	flag.Parse()

	if *f == "" {
		flag.Usage()
	}

	file, err := os.Open(*f)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	lines := 0

	var k [16]byte

	var salts []uint32

	for i, needed := uint(0), dgobloom.SaltsRequired(uint32(*n), *fprate); i < needed; i++ {
		var u [4]byte
		rand.Read(u[:])
		salts = append(salts, binary.LittleEndian.Uint32(u[:]))
	}

	b := dgobloom.NewBloomFilter(uint32(*n), *fprate, h32{siphash.New(k[:])}, salts)

	for scanner.Scan() {
		lines++
		if ok := b.Exists(scanner.Bytes()); !ok {
			os.Stdout.Write(scanner.Bytes())
			os.Stdout.Write([]byte("\n"))
		}
		b.Insert(scanner.Bytes())
		if lines%(1<<20) == 0 {
			log.Println(lines)
		}
	}
}
