package main

import (
	"bufio"
	"log"
	"os"
)

var lookupMap = make(map[string]struct{})

func initMap() {

	f, err := os.Open("domains.txt")
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lookupMap[scanner.Text()] = struct{}{}
	}
}

func MatchMap(k []byte) bool {
	_, ok := lookupMap[string(k)]
	return ok
}
