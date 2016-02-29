package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dchest/siphash"
)

func main() {

	keys := flag.String("k", "", "key0,key1")

	flag.Parse()

	var k0, k1 uint64

	if *keys != "" {
		var err error
		k0, k1, err = parseKeys(*keys)

		if err != nil {
			log.Fatal(err)
		}
	}

	if flag.NArg() == 0 {
		hashOne(k0, k1, os.Stdin, "stdin")
		return
	}

	for _, arg := range flag.Args() {
		f, err := os.Open(arg)
		if err != nil {
			log.Println(f, err)
			continue
		}
		hashOne(k0, k1, f, arg)
	}
}

func hashOne(k0, k1 uint64, r io.Reader, name string) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(name, err)
		return
	}
	h := siphash.Hash(k0, k1, buf)
	fmt.Printf("%016x %s\n", h, name)
}

func parseKeys(keys string) (uint64, uint64, error) {

	kstr := strings.SplitN(keys, ",", 2)

	var k64 []uint64

	for _, k := range kstr {
		k0, err := strconv.ParseUint(k, 16, 64)
		if err != nil {
			return 0, 0, err
		}
		k64 = append(k64, k0)
	}

	return k64[0], k64[1], nil
}
