package main

import (
	"bufio"
	"log"
	"os"

	"github.com/dgryski/go-mph"
)

var lookupDMPH *mph.Table

var domains []string

func initDMPH() {

	f, err := os.Open("domains.txt")
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	lookupDMPH = mph.New(domains)

}

func MatchDMPH(s string) bool {
	i := lookupDMPH.Query(s)
	return domains[i] == s
}
