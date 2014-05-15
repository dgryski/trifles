// Package main is a reducer for calculating entropy for values for various categories
/*
Input format:
       epoch   value   category

The file must be sorted by epoch.

*/
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/dgryski/go-entropy"
)

func main() {

	ex := make(map[string]entropy.Exact)

	scanner := bufio.NewScanner(os.Stdin)

	var epoch int

	for scanner.Scan() {
		fields := bytes.Fields(scanner.Bytes())

		e, _ := strconv.Atoi(string(fields[0]))
		if e != epoch {
			for k, v := range ex {
				fmt.Println(epoch, k, v.Entropy())
			}

			ex = make(map[string]entropy.Exact)
			epoch = e
		}

		m, ok := ex[string(fields[2])]
		if !ok {
			m = entropy.NewExact()
			ex[string(fields[2])] = m
		}

		m.Push(fields[1], 1)
	}

	for k, v := range ex {
		fmt.Println(epoch, k, v.Entropy())
	}

}
