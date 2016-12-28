package main

import (
	"log"

	"github.com/smartystreets/mafsa"
)

var lookupMAFSA *mafsa.MinTree

func initMAFSA() {

	var err error

	lookupMAFSA, err = mafsa.Load("mafsa.out")
	if err != nil {
		log.Fatalf("error loading mafsa.out: %v", err)
	}

}

func MatchMAFSA(r []rune) bool {
	result := lookupMAFSA.Traverse(r)
	return result != nil && result.Final
}
