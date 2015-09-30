package main

import "log"

func stablePartitionInts(d []int, f, l int, p func(i int) bool) int {

	n := l - f
	if n == 0 {
		return f
	}

	if n == 1 {
		r := f
		if p(d[f]) {
			r++
		}
		return r
	}

	m := f + n/2

	return rotateInts(d, stablePartitionInts(d, f, m, p), m, stablePartitionInts(d, m, l, p))
}

func rotateInts(d []int, f, k, l int) int {
	reverseInts(d, f, k)
	reverseInts(d, k, l)
	reverseInts(d, f, l)
	return f + l - k
}

func reverseInts(d []int, f, l int) {
	lend := (l - f)
	for i, j := 0, lend-1; i < lend/2; i, j = i+1, j-1 {
		d[f+i], d[f+j] = d[f+j], d[f+i]
	}
}

func stablePartitionMove(array []int, sz int, p func(i int) bool) {

	var shift_index int
	for shift_index < sz && p(array[shift_index]) {
		shift_index++
	}

	index := shift_index + 1

	for index < sz {
		// find the next number that needs to be shifted left
		if !p(array[index]) {
			index++
		} else {
			value := array[index]
			num_elements_to_move := index - shift_index
			copy(array[shift_index+1:shift_index+1+num_elements_to_move], array[shift_index:shift_index+num_elements_to_move])
			array[shift_index] = value // move the element to the hole left at the end
			shift_index++
			index++
		}
	}
}

func main() {
	array := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 8, 9, 7, 9}

	ints := append([]int(nil), array...)
	stablePartitionMove(ints, len(ints), func(i int) bool { return i%2 == 0 })
	log.Printf("ints=%+v\n", ints)

	copy(ints, array)
	stablePartitionMove(ints, len(ints), func(i int) bool { return i%2 == 0 })
	log.Printf("ints=%+v\n", ints)
}
