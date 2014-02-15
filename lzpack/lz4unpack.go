package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"os"

	lz4 "github.com/bkaradzic/go-lz4"
	"github.com/dchest/siphash"
	binpack "github.com/dgryski/go-binpack"
)

func main() {

	verbose := flag.Bool("v", false, "verbose")

	flag.Parse()

	r := bufio.NewReader(os.Stdin)

	for {

		var header struct {
			PacketSize       uint32
			SipHash          uint64
			Fname            []byte `binpack:"lenprefix=uint16"`
			UncompressedSize uint32
		}

		err := binpack.Read(r, binary.LittleEndian, &header)
		if err == io.EOF {
			break
		}

		compressedSize := header.PacketSize - 16 - 2 - uint32(len(header.Fname))

		if *verbose {
			log.Printf("%s: %d -> %d\n", string(header.Fname), compressedSize, header.UncompressedSize)
		}

		data := make([]byte, compressedSize)

		n, err := io.ReadFull(r, data)
		if n != len(data) || err != nil {
			log.Fatalf("bad read: n=%d len(data)=%d err=%s\n", n, len(data), err)
		}

		data, err = lz4.Decode(nil, data)

		h := siphash.Hash(0, 0, data)

		if err != nil || h != header.SipHash {
			log.Println("FAIL: err=", err)
		}
	}
}
