package main

import (
	"math/rand"
	"testing"
)

func benchmarkInts(b *testing.B, n int) {

	rand.Seed(0)

	ints := make([]int, n)
	for i := range ints {
		ints[i] = rand.Intn(1000)
	}

	orig := make([]int, n)
	copy(orig, ints)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stablePartitionInts(ints, 0, len(ints), func(i int) bool { return i&1 == 0 })
		copy(ints, orig)
	}
}

func benchmarkMove(b *testing.B, n int) {

	rand.Seed(0)

	ints := make([]int, n)
	for i := range ints {
		ints[i] = rand.Intn(1000)
	}

	orig := make([]int, n)
	copy(orig, ints)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		stablePartitionMove(ints, len(ints), func(i int) bool { return i&1 == 0 })
		copy(ints, orig)
	}
}

func BenchmarkInts10(b *testing.B)    { benchmarkInts(b, 10) }
func BenchmarkInts20(b *testing.B)    { benchmarkInts(b, 20) }
func BenchmarkInts50(b *testing.B)    { benchmarkInts(b, 50) }
func BenchmarkInts100(b *testing.B)   { benchmarkInts(b, 100) }
func BenchmarkInts200(b *testing.B)   { benchmarkInts(b, 200) }
func BenchmarkInts500(b *testing.B)   { benchmarkInts(b, 500) }
func BenchmarkInts1000(b *testing.B)  { benchmarkInts(b, 1000) }
func BenchmarkInts2000(b *testing.B)  { benchmarkInts(b, 2000) }
func BenchmarkInts5000(b *testing.B)  { benchmarkInts(b, 5000) }
func BenchmarkInts10000(b *testing.B) { benchmarkInts(b, 10000) }

func BenchmarkMove10(b *testing.B)    { benchmarkMove(b, 10) }
func BenchmarkMove20(b *testing.B)    { benchmarkMove(b, 20) }
func BenchmarkMove50(b *testing.B)    { benchmarkMove(b, 50) }
func BenchmarkMove100(b *testing.B)   { benchmarkMove(b, 100) }
func BenchmarkMove200(b *testing.B)   { benchmarkMove(b, 200) }
func BenchmarkMove500(b *testing.B)   { benchmarkMove(b, 500) }
func BenchmarkMove1000(b *testing.B)  { benchmarkMove(b, 1000) }
func BenchmarkMove2000(b *testing.B)  { benchmarkMove(b, 2000) }
func BenchmarkMove5000(b *testing.B)  { benchmarkMove(b, 5000) }
func BenchmarkMove10000(b *testing.B) { benchmarkMove(b, 10000) }
