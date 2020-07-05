package main

import "github.com/dgryski/trifles/fsmdemo/fishasm"

import "C"

//export fishasmMatch
func fishasmMatch(base *C.char, len C.int) C.int {
	return C.int(fishasm.Match(C.GoStringN(base, len)))
}

func main() {}
