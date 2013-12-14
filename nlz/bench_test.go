package main

import (
	"testing"
)

func BenchmarkBSR(b *testing.B) {

	//	var x uint32 = 0xFFFFFFFF

	N := uint32(b.N)
	for i := uint32(0); i < N; i++ {
		/* x += x
		if int(x) < 0 {
			x ^= 0x88888EEF
		} */
		nlz(uint32(i))
	}

}

func BenchmarkGo(b *testing.B) {

	//	var x uint32 = 0xFFFFFFFF

	N := uint32(b.N)
	for i := uint32(0); i < N; i++ {
		/* x += x
		if int(x) < 0 {
			x ^= 0x88888EEF
		} */
		nlz2(uint32(i))
	}
}
