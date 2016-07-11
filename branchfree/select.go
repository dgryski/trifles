package main

import (
	"log"
)

func inlineMin(x, y uint64) uint64 {
	return (uint64((int64(x)-int64(y))>>63) & (x ^ y)) ^ y
}

func fastMin(x, y uint64) uint64 {
	diff := int64(x) - int64(y)

	// sets bit63 to 0xFFFFFFFF if x<y, else 0x00000000
	bit63 := uint64(diff >> 63)

	// return ifXisSmaller if x is smaller than y, else y
	return (bit63 & (x ^ y)) ^ y
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}

	return y
}

func main() {
	log.Printf("min(1,20)=%+v\n", fastMin(1, 20))
}
