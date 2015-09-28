package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
)

func StablePartition(d sort.Interface, f, l int, p func(i int) bool) int {

	n := l - f
	if n == 0 {
		return f
	}

	if n == 1 {
		r := f
		if p(f) {
			r++
		}
		return r
	}

	m := f + n/2

	return Rotate(d, StablePartition(d, f, m, p), m, StablePartition(d, m, l, p))
}

func Rotate(d sort.Interface, f, k, l int) int {
	Reverse(d, f, k)
	Reverse(d, k, l)
	Reverse(d, f, l)
	return f + l - k
}

func Reverse(d sort.Interface, f, l int) {
	lend := (l - f)
	for i, j := 0, lend-1; i < lend/2; i, j = i+1, j-1 {
		d.Swap(f+i, f+j)
	}
}

type element struct {
	val int
	idx int
}

type elements []element

func (e elements) Len() int           { return len(e) }
func (e elements) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e elements) Less(i, j int) bool { return e[i].val < e[j].val }

func verify(pre, elts elements, f, l int, which string, p func(i int) bool) {

	// check the predicate
	for i := f; i < l; i++ {
		if !p(f) {
			fmt.Printf("%s predicate failed %v\n", which, elts[0])
			fmt.Printf("pre =%+v\n", pre)
			fmt.Printf("elts=%+v\n", elts)
			log.Fatalln("verify failed")
		}

	}

	// check the order
	for i := f + 1; i < l; i++ {
		if elts[i-1].idx > elts[i].idx {
			fmt.Printf("%s partition mismatch: %v %v\n", which, elts[i-1], elts[i])
			fmt.Printf("pre =%+v\n", pre)
			fmt.Printf("elts=%+v\n", elts)
			log.Fatalln("verify failed")
		}
	}
}

func main() {

	const numTests = 100000

	for j := 0; j < numTests; j++ {

		// generate random slice
		sz := rand.Intn(100)
		elts := make(elements, sz)
		for i := range elts {
			elts[i].val = rand.Intn(1000)
			elts[i].idx = i
		}

		// save a copy
		pre := make(elements, sz)
		copy(pre, elts)

		predicate := func(i int) bool { return elts[i].val&1 == 0 }

		// partition
		p := StablePartition(elts, 0, elts.Len(), predicate)

		// verify it's correct
		verify(pre, elts, 0, p, "pre", predicate)
		verify(pre, elts, p, sz, "post", func(i int) bool { return !predicate(i) })
	}

	fmt.Printf("%d tests succeeded\n", numTests)
}
