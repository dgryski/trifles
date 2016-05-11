package main

import (
	"github.com/dgryski/trifles/ewmaest"
	"log"
	"math/rand"
	"time"
)

func main() {

	const (
		totalItem    = 1000
		maxBlockSize = 20
	)

	est := ewmaest.New(totalItem)

	est.StartBlock()
	blockSize := rand.Intn(maxBlockSize) + 1
	for i := 0; i < totalItem; i += blockSize {
		time.Sleep(time.Duration(blockSize) * (50*time.Millisecond + time.Duration(rand.Intn(5*int(time.Millisecond)))))
		eta := est.CompleteItems(blockSize)
		log.Println(est.Progress(eta))
		est.StartBlock()
		blockSize = rand.Intn(maxBlockSize) + 1
	}
}
