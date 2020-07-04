package main

import (
	"flag"
	"io/ioutil"
	"log"
	"unsafe"

	"github.com/dgryski/trifles/testlibfsm/fishasm"
	"github.com/dgryski/trifles/testlibfsm/fishfsm"
)

/*

int fishmain(const char *b, const char *e);

*/
import "C"

func main() {
	input := flag.String("-i", "../matcher/domains.txt", "input file")
	flag.Parse()

	data, err := ioutil.ReadFile(*input)
	if err != nil {
		log.Fatalf("Error reading %q: %v", *input, err)
	}

	var b []byte
	var caughtfsm int
	var caughtasm int
	var caughtvmc int
	var caughtrx int

	b = data
	for len(b) > 20 {
		if fishfsm.Match(b) != -1 {
			caughtfsm++
		}
		b = b[20:]
	}

	b = data
	for len(b) > 20 {
		if fishasm.Match(b) != -1 {
			caughtasm++
		}
		b = b[20:]
	}

	b = data
	for len(b) > 20 {
		if vmcMatch(b) != -1 {
			caughtvmc++
		}
		b = b[20:]
	}

	b = data
	for len(b) > 20 {
		if fishrx.Match(b) {
			caughtrx++
		}
		b = b[20:]
	}

	if caughtfsm != caughtasm || caughtasm != caughtrx || caughtfsm != caughtvmc {
		log.Fatalf("matcher count mismatch: fsm=%d asm=%d vmc=%d rx=%d", caughtfsm, caughtasm, caughtvmc, caughtrx)
	}
}

func vmcMatch(s []byte) int {
	return int(C.fishmain((*C.char)(unsafe.Pointer(&s[0])), (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(&s[0]))+uintptr(len(s))*unsafe.Sizeof(s[0])))))
}
