package u64cmp

import "testing"

var dataA = make([]uint64, 500*1024)
var dataB = make([]uint64, 500*1024)

var Sink bool

func BenchmarkNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sink = Sink || naive(dataA, dataB)
	}
}

func BenchmarkUnrollXor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sink = Sink || unrollxor(dataA, dataB)
	}
}
