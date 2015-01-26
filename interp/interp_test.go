package interp

import (
	"math/rand"
	"sort"
	"testing"
)

func benchmarkInterp(b *testing.B, limit int) {

	rand.Seed(0)

	ints := make([]int, limit)

	for i := 0; i < limit; i++ {
		ints[i] = int(int32(rand.Int()))
	}

	sort.Ints(ints)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		elt := ints[rand.Intn(limit)]
		Search(ints, elt)
	}
}

func benchmarkBin(b *testing.B, limit int) {

	rand.Seed(0)

	ints := make([]int, limit)

	for i := 0; i < limit; i++ {
		ints[i] = int(int32(rand.Int()))
	}

	sort.Ints(ints)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		elt := ints[rand.Intn(limit)]
		sort.SearchInts(ints, elt)
	}
}

func BenchmarkInterp100(b *testing.B)   { benchmarkInterp(b, 100) }
func BenchmarkInterp1000(b *testing.B)  { benchmarkInterp(b, 1000) }
func BenchmarkInterp10000(b *testing.B) { benchmarkInterp(b, 10000) }
func BenchmarkInterp1e5(b *testing.B)   { benchmarkInterp(b, 1e5) }
func BenchmarkInterp1e6(b *testing.B)   { benchmarkInterp(b, 1e6) }

func BenchmarkBin100(b *testing.B)   { benchmarkBin(b, 100) }
func BenchmarkBin1000(b *testing.B)  { benchmarkBin(b, 1000) }
func BenchmarkBin10000(b *testing.B) { benchmarkBin(b, 10000) }
func BenchmarkBin1e5(b *testing.B)   { benchmarkBin(b, 1e5) }
func BenchmarkBin1e6(b *testing.B)   { benchmarkBin(b, 1e6) }

func TestSearchSmall(t *testing.T) {

	const limit = 100

	t.Log("generating")

	var ints []int

	for i := 0; i < limit; i++ {
		ints = append(ints, i*2)
	}

	t.Log("searching")

	queries := []int{0, 300, 51, 50, 18, 19, 20, 21, 22, 23, 24, 198, 199, 200, 201, 202}

	for _, q := range queries {

		var want int
		if q >= 200 || q&1 == 1 {
			want = -1
		} else {
			want = q / 2
		}

		if idx := Search(ints, q); idx != want {
			t.Errorf("Search(ints, %v)=%v, want %v", q, idx, want)
		}
	}
}
