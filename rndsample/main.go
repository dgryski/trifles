package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	sampleSize := flag.Int("n", 1000, "sample size")

	flag.Parse()

	if *sampleSize < 0 {
		log.Fatal("sample size must be >0: got ", *sampleSize)
	}

	rand.Seed(time.Now().UnixNano())

	r := newReservoir(*sampleSize)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		r.push(scanner.Text())
	}

	for _, s := range r.data {
		fmt.Println(s)
	}
}

type reservoir struct {
	data []string
	n    int
}

func newReservoir(capacity int) *reservoir {
	return &reservoir{
		data: make([]string, 0, capacity),
	}
}

func (r *reservoir) push(s string) {

	index := r.n

	r.n++

	// not enough samples yet -- add it
	if index < cap(r.data) {
		r.data = append(r.data, s)
		return
	}

	ridx := rand.Intn(r.n) // == index+1, so we're 0..index inclusive

	if ridx >= len(r.data) {
		// ignore this one
		return
	}

	// add to our sample
	r.data[ridx] = s
}
