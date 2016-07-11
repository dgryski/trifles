package main

import (
	"testing"
	"testing/quick"
)

func TestFastMin(t *testing.T) {

	f := func(x, y uint64) bool {
		// must be positive
		x &^= 1 << 63
		y &^= 1 << 63
		m := fastMin(x, y)
		n := min(x, y)
		o := inlineMin(x, y)
		return m == n && n == o
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

var total uint64

func BenchmarkMin(b *testing.B) {
	x := uint64(1)
	y := XorshiftMult64(x)

	for i := 0; i < b.N; i++ {

		xp := x &^ 1 << 63
		yp := y &^ 1 << 63

		total += min(xp, yp)
		x, y = y, XorshiftMult64(y)
	}
}

func BenchmarkFastMin(b *testing.B) {
	x := uint64(1)
	y := XorshiftMult64(x)

	total = 0

	for i := 0; i < b.N; i++ {

		xp := x &^ 1 << 63
		yp := y &^ 1 << 63

		total += fastMin(xp, yp)
		x, y = y, XorshiftMult64(y)
	}
}

func BenchmarkInlineMin(b *testing.B) {
	x := uint64(1)
	y := XorshiftMult64(x)

	total = 0

	for i := 0; i < b.N; i++ {

		xp := x &^ 1 << 63
		yp := y &^ 1 << 63

		total += inlineMin(xp, yp)
		x, y = y, XorshiftMult64(y)
	}
}

// 64-bit xorshift multiply rng from http://vigna.di.unimi.it/ftp/papers/xorshift.pdf
func XorshiftMult64(x uint64) uint64 {
	x ^= x >> 12 // a "te
	x ^= x << 25 // b
	x ^= x >> 27 // c
	return x * 2685821657736338717
}
