package main

import (
	"fmt"
	"sort"
)

func StablePartition(d sort.Interface, f, l int, p func(i int) bool) int {

	n := l - f
	if n == 0 {
		return f
	}

	if n == 1 {
		if p(f) {
			return f + 1
		} else {
			return f
		}
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

func main() {
	d := []int{1, 4, 2, 5, 6, 3, 4, 1, 7, 5, 4, 1}
	fmt.Println(d)
	o := StablePartition(sort.IntSlice(d), 0, len(d), func(i int) bool { return d[i]&1 == 0 })
	fmt.Println(d, o)
}
