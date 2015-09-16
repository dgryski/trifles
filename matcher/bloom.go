package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"os"
)

var lookupBloom bitvector

const bloomK = 2
const bloomM = 1 << 14

func bkeys(b []byte) (uint32, uint32) {
	h1 := binary.LittleEndian.Uint32(b)
	i := 4
	if i > len(b)-4 {
		i = len(b) - 4
	}
	h2 := binary.LittleEndian.Uint32(b[i:])
	return Xorshift32(h1), Xorshift32(h2)
}

func initBloom() {

	lookupBloom = newbv(bloomM)

	f, err := os.Open("domains.txt")
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b := scanner.Bytes()

		h1, h2 := bkeys(b)

		for i := uint32(0); i < bloomK; i++ {
			lookupBloom.set((h1 + (i * h2)) & (bloomM - 1))
		}
	}
}

var filtered int

// Lookup checks the bloom filter for the byte array b
func MatchBloom(b []byte) bool {

	if lookupBloom == nil {
		return true
	}

	if len(b) < 4 {
		return true
	}

	h1, h2 := bkeys(b)

	for i := uint32(0); i < bloomK; i++ {
		if lookupBloom.get((h1+(i*h2))&(bloomM-1)) == 0 {
			filtered++
			return false
		}
	}

	return true
}

// Internal routines for the bit vector
type bitvector []uint64

func newbv(size uint32) bitvector {
	return make([]uint64, uint(size+63)/64)
}

// get bit 'bit' in the bitvector d
func (b bitvector) get(bit uint32) uint {

	shift := bit % 64
	bb := b[bit/64]
	bb &= (1 << shift)

	return uint(bb >> shift)
}

// set bit 'bit' in the bitvector d
func (b bitvector) set(bit uint32) {
	b[bit/64] |= (1 << (bit % 64))
}

// Xorshift32 is an xorshift RNG
func Xorshift32(y uint32) uint32 {

	// http://www.jstatsoft.org/v08/i14/paper
	// Marasaglia's "favourite"

	y ^= (y << 13)
	y ^= (y >> 17)
	y ^= (y << 5)
	return y
}
