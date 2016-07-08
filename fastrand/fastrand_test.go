package fastrand

import (
	"testing"
)

func BenchmarkIntn(b *testing.B) {

	r := Rng(0xdeadbeef)

	for i := 0; i < b.N; i++ {
		r.Intn(1000)
	}
}

func BenchmarkInt32n(b *testing.B) {
	r := Rng(0xdeadbeef)
	for i := 0; i < b.N; i++ {
		r.Int32n(1000)
	}
}

func TestIntn(t *testing.T) {

	var counts [1000]int

	r := Rng(0xdeadbeef)
	for i := 0; i < 100e6; i++ {
		counts[r.Intn(1000)]++
	}

	t.Logf("counts = %+v\n", counts)
}

func TestInt32n(t *testing.T) {

	var counts [1000]int

	r := Rng(0xdeadbeef)
	for i := 0; i < 100e6; i++ {
		counts[r.Int32n(1000)]++
	}

	t.Logf("counts = %+v\n", counts)
}
