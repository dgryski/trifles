package main

import (
	"bufio"
	"log"
	"os"

	"github.com/cespare/mph"
)

var lookupMPH *mph.Table

func initMPH() {

	f, err := os.Open("domains.txt")
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	scanner := bufio.NewScanner(f)
	var arr []string
	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}

	lookupMPH = mph.Build(arr)

}

func MatchMPH(s string) bool {
	_, ok := lookupMPH.Lookup(s)
	return ok
}
