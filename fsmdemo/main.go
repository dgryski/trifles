package main

//go:generate sh -c "re -r pcre -k str -e fish  -pl go 'fi+sh' >fishfsm/fishfsm.go"

import (
	"fmt"

	"github.com/dgryski/trifles/fsmdemo/fishfsm"
)

func main() {
	net := []string{"fsh", "catfish", "dogo", "fiiish"}
	for _, caught := range net {
		if fishfsm.Match(caught) != -1 {
			fmt.Println("caught a fish: ", caught)
		}
	}
}
