package main

import (
	"bufio"
	"log"
	"os"
	"sort"
)

var lookupBsearch []string

func initBsearch() {

	f, err := os.Open("domains.txt")
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lookupBsearch = append(lookupBsearch, scanner.Text())
	}
	sort.Strings(lookupBsearch)
}

func MatchBsearch(s string) bool {
	i := sort.SearchStrings(lookupBsearch, s)
	// the search fn gives the pos where it should be (if it's not there, where to insert it)
	// so we still need to verify
	return i < len(lookupBsearch) && lookupBsearch[i] == s
}
