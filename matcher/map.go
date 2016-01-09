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

func MatchMap(s string) bool {
	_, ok := lookupMap[s]
	return ok
}
