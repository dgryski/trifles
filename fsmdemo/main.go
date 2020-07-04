package main

//go:generate sh -c "re -r pcre -k str -e fish  -pl go 'fi+sh' >fishfsm/fishfsm.go"
//go:generate sh -c "re -r pcre -k str -e fish  -pl amd64_go 'fi+sh' >fishasm/fishasm.s"
//go:generate sh -c "re -r pcre -k pair -e fish  -pl vmc 'fi+sh' >fishvmc.c"

import (
	"fmt"
	"unsafe"

	"github.com/dgryski/trifles/fsmdemo/fishasm"
	"github.com/dgryski/trifles/fsmdemo/fishfsm"
)

/*

#include <stdlib.h>

int fishmain(const char *b, const char *e);

*/
import "C"

func main() {
	goFishing("fsm", fishfsm.Match)
	goFishing("asm", fishasm.Match)
	goFishing("C", vmcMatch)
}

func goFishing(what string, match func(s string) int) {

	lake := []string{"fsh", "catfishes", "dogo", "fiiish"}

	for _, caught := range lake {
		if match(caught) != -1 {
			fmt.Println("caught fish with", what, ":", caught)
		}
	}
}

func vmcMatch(s string) int {
	cstr := C.CString(s)
	C.free(unsafe.Pointer(cstr))
	return int(C.fishmain(cstr, (*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cstr))+uintptr(len(s))))))
}
