// +build ignore

package main

import (
	"bufio"
	"log"
	"os"
	"sort"

	"github.com/smartystreets/mafsa"
)

func main() {

	f, err := os.Open("domains.txt")
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	scanner := bufio.NewScanner(f)
	var arr []string
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}

	sort.Strings(arr)

	bt := mafsa.New()

	for _, s := range arr {
		if err := bt.Insert(s); err != nil {
			log.Fatalf("error inserting %q: %v", s, err)
		}
	}

	bt.Finish()

	if err := bt.Save("mafsa.out"); err != nil {
		log.Fatalf("error saving mafsa: %v", err)
	}
}
