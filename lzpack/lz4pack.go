package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"io/ioutil"
	"log"
	"os"

	lz4 "github.com/bkaradzic/go-lz4"
	"github.com/dchest/siphash"
	binpack "github.com/dgryski/go-binpack"
)

func main() {

	verbose := flag.Bool("v", false, "verbose")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for scanner.Scan() {

		fname := scanner.Text()

		b, err := ioutil.ReadFile(fname)
		if err != nil {
			log.Println(err)
			continue
		}

		sip := siphash.Hash(0, 0, b)

		data, err := lz4.Encode(nil, b)
		if err != nil {
			log.Println("error packing: ", err)
			continue
		}

		const headerSize = 4 + 8 + 2 + 4 // + len(fname)

		h := struct {
			PacketSize       uint32
			SipHash          uint64
			Fname            []byte `binpack:"lenprefix=uint16"`
			UncompressedSize uint32
		}{
			PacketSize:       uint32(headerSize + len(fname) + len(data)),
			SipHash:          sip,
			Fname:            []byte(fname),
			UncompressedSize: uint32(len(b)),
		}

		if *verbose {
			log.Printf("%s: %d -> %d\n", fname, len(b), len(data))
		}

		err = binpack.Write(w, binary.LittleEndian, h)
		if err != nil {
			log.Printf("error writing header %+v\n", err)
			break
		}

		_, err = w.Write(data)
		if err != nil {
			log.Printf("error writing %+v\n", err)
			break
		}
	}
}
