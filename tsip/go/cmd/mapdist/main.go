package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/dchest/siphash"
	"github.com/dgryski/go-marvin32"
	"github.com/dgryski/go-sip13"
	tsip "github.com/dgryski/trifles/tsip/go"
)

func main() {

	bits := flag.Uint("s", 8, "table bits")
	f := flag.String("f", "tsip", "function to test")

	flag.Parse()

	var h func(uint64, uint64, []byte) uint64

	switch *f {
	case "tsip":
		h = tsip.Hash
	case "siphash":
		h = siphash.Hash
	case "sip13":
		h = sip13.Sum64
	case "marvin32":
		h = func(k0, k1 uint64, m []byte) uint64 { return uint64(marvin32.Sum32(0, m)) }
	}

	size := uint64(1<<*bits) - 1

	t := make([]uint32, size+1)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		b := scanner.Bytes()
		a := h(0, 0, b)
		t[a&size]++
	}

	for i, v := range t {
		fmt.Println(i, v)
	}
}
