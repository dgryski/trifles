package main

import (
	"bufio"
	"log"
	"os"

	"github.com/cloudflare/ahocorasick"
)

var lookupAho *ahocorasick.Matcher

func initAho() {

	f, err := os.Open("domains.txt")
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	scanner := bufio.NewScanner(f)
	var arr [][]byte
	for scanner.Scan() {
		b := append([]byte(nil), scanner.Bytes()...)
		arr = append(arr, b)
	}

	lookupAho = ahocorasick.NewMatcher(arr)

}

func MatchAho(b []byte) bool {
	return lookupAho.Match(b) != nil
}
