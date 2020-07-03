package main

//go:generate sh -c "re -r pcre -k str -e fish  -pl go 'fi+sh' >fishfsm/fishfsm.go"
//go:generate sh -c "re -r pcre -k str -e fish  -pl amd64_go 'fi+sh' >fishasm/fishasm.s"

import (
	"fmt"

	"github.com/dgryski/trifles/fsmdemo/fishasm"
	"github.com/dgryski/trifles/fsmdemo/fishfsm"
)

func main() {
	net := []string{"fsh", "catfish", "dogo", "fiiish"}
	for _, caught := range net {
		if fishfsm.Match(caught) != -1 {
			fmt.Println("caught a fish: ", caught)
		}
		if fishasm.Match(caught) != -1 {
			fmt.Println("caught an asm fish: ", caught)
		}
	}
}
