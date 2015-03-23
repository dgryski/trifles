package interp

import (
	"math/rand"
	"sort"
	"testing"
)

const maxLimit = 3e7

var Ints []int

func benchmarkSearch(b *testing.B, limit int, search func([]int, int) int) {

	if Ints == nil {

		rand.Seed(0)

		Ints = make([]int, maxLimit)

		for i := 0; i < limit; i++ {
			Ints[i] = int(int32(rand.Int()))
		}

		sort.Ints(Ints)
	}

	ints := Ints[:limit]

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		elt := ints[rand.Intn(limit)]
		search(ints, elt)
	}
}

func BenchmarkInterp100(b *testing.B)   { benchmarkSearch(b, 100, Search) }
func BenchmarkInterp1000(b *testing.B)  { benchmarkSearch(b, 1000, Search) }
func BenchmarkInterp10000(b *testing.B) { benchmarkSearch(b, 10000, Search) }
func BenchmarkInterp1e5(b *testing.B)   { benchmarkSearch(b, 1e5, Search) }
func BenchmarkInterp1e6(b *testing.B)   { benchmarkSearch(b, 1e6, Search) }
func BenchmarkInterp1e7(b *testing.B)   { benchmarkSearch(b, 1e7, Search) }
func BenchmarkInterp3e7(b *testing.B)   { benchmarkSearch(b, 3e7, Search) }

func BenchmarkBin100(b *testing.B)   { benchmarkSearch(b, 100, binsearch) }
func BenchmarkBin1000(b *testing.B)  { benchmarkSearch(b, 1000, binsearch) }
func BenchmarkBin10000(b *testing.B) { benchmarkSearch(b, 10000, binsearch) }
func BenchmarkBin1e5(b *testing.B)   { benchmarkSearch(b, 1e5, binsearch) }
func BenchmarkBin1e6(b *testing.B)   { benchmarkSearch(b, 1e6, binsearch) }
func BenchmarkBin1e7(b *testing.B)   { benchmarkSearch(b, 1e7, binsearch) }
func BenchmarkBin3e7(b *testing.B)   { benchmarkSearch(b, 3e7, binsearch) }

func BenchmarkStd100(b *testing.B)   { benchmarkSearch(b, 100, sort.SearchInts) }
func BenchmarkStd1000(b *testing.B)  { benchmarkSearch(b, 1000, sort.SearchInts) }
func BenchmarkStd10000(b *testing.B) { benchmarkSearch(b, 10000, sort.SearchInts) }
func BenchmarkStd1e5(b *testing.B)   { benchmarkSearch(b, 1e5, sort.SearchInts) }
func BenchmarkStd1e6(b *testing.B)   { benchmarkSearch(b, 1e6, sort.SearchInts) }
func BenchmarkStd1e7(b *testing.B)   { benchmarkSearch(b, 1e7, sort.SearchInts) }
func BenchmarkStd3e7(b *testing.B)   { benchmarkSearch(b, 3e7, sort.SearchInts) }

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

func binsearch(sortedArray []int, toFind int) int {

	low, high := 0, len(sortedArray)

	for low < high {
		mid := low + (high-low)/2
		if sortedArray[mid] > toFind {
			high = mid
		} else {
			low = mid + 1
		}
	}

	if low >= len(sortedArray) {
		low = -1
	}
	return low
}
