package main

import (
	"bufio"
	"log"
	"os"

	"github.com/armon/go-radix"
)

var lookupTree = radix.New()

func initRadix() {

	f, err := os.Open("domains.txt")
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lookupTree.Insert(scanner.Text(), struct{}{})
	}
}

func MatchRadix(s string) bool {
	_, ok := lookupTree.Get(s)
	return ok
}
