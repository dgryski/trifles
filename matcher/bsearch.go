package main

import (
	"bufio"
	"log"
	"os"
	"sort"
)

type slice []string

var data slice

func (a slice) Len() int { return len(a) }

func (a slice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a slice) Less(i, j int) bool { return a[i] < a[j] }

func initBsearch() {

	f, err := os.Open("domains.txt")
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	sort.Sort(data)
}

func MatchBsearch(k []byte) bool {
	in := string(k)
	i := sort.SearchStrings(data, in)
	// the search fn gives the pos where it should be (if it's not there, where to insert it)
	// so we still need to verify
	if i < len(data) && data[i] == in {
		return true
	}
	return false
}
