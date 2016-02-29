package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dchest/siphash"
)

func main() {

	scan := bufio.NewScanner(os.Stdin)

	for scan.Scan() {
		u := siphash.Hash(0, 0, scan.Bytes())
		fmt.Printf("%016x\n", u)
	}
}
