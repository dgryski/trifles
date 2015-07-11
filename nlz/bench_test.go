package main

import (
	"testing"
)

func BenchmarkBSR(b *testing.B) {

	//	var x uint64 = 0xFFFFFFFF

	N := uint64(b.N)
	for i := uint64(0); i < N; i++ {
		/* x += x
		if int(x) < 0 {
			x ^= 0x88888EEF
		} */
		nlz(uint64(i))
	}

}

func BenchmarkGo(b *testing.B) {

	//	var x uint64 = 0xFFFFFFFF

	N := uint64(b.N)
	for i := uint64(0); i < N; i++ {
		/* x += x
		if int(x) < 0 {
			x ^= 0x88888EEF
		} */
		nlz2(uint64(i))
	}
}

func BenchmarkPopcount(b *testing.B) {

	N := uint64(b.N)
	for i := uint64(0); i < N; i++ {
		popcount(uint64(i))
	}

}

func BenchmarkPopcountGo(b *testing.B) {

	N := uint64(b.N)
	for i := uint64(0); i < N; i++ {
		popcountGo(uint64(i))
	}
}
