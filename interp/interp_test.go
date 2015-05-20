package interp

import (
	"math/rand"
	"sort"
	"testing"
)

const maxLimit = 4e7

var Ints []int

func fillInts() {
	rand.Seed(0)

	Ints = make([]int, maxLimit)

	for i := 0; i < maxLimit; i++ {
		Ints[i] = rand.Int()
	}

	sort.Ints(Ints)
}

func benchmarkSearch(b *testing.B, limit int, search func([]int, int) int) {

	if Ints == nil {
		fillInts()
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
func BenchmarkInterp4e7(b *testing.B)   { benchmarkSearch(b, 4e7, Search) }

func BenchmarkBin100(b *testing.B)   { benchmarkSearch(b, 100, binsearch) }
func BenchmarkBin1000(b *testing.B)  { benchmarkSearch(b, 1000, binsearch) }
func BenchmarkBin10000(b *testing.B) { benchmarkSearch(b, 10000, binsearch) }
func BenchmarkBin1e5(b *testing.B)   { benchmarkSearch(b, 1e5, binsearch) }
func BenchmarkBin1e6(b *testing.B)   { benchmarkSearch(b, 1e6, binsearch) }
func BenchmarkBin1e7(b *testing.B)   { benchmarkSearch(b, 1e7, binsearch) }
func BenchmarkBin4e7(b *testing.B)   { benchmarkSearch(b, 4e7, binsearch) }

func BenchmarkStd100(b *testing.B)   { benchmarkSearch(b, 100, sort.SearchInts) }
func BenchmarkStd1000(b *testing.B)  { benchmarkSearch(b, 1000, sort.SearchInts) }
func BenchmarkStd10000(b *testing.B) { benchmarkSearch(b, 10000, sort.SearchInts) }
func BenchmarkStd1e5(b *testing.B)   { benchmarkSearch(b, 1e5, sort.SearchInts) }
func BenchmarkStd1e6(b *testing.B)   { benchmarkSearch(b, 1e6, sort.SearchInts) }
func BenchmarkStd1e7(b *testing.B)   { benchmarkSearch(b, 1e7, sort.SearchInts) }
func BenchmarkStd4e7(b *testing.B)   { benchmarkSearch(b, 4e7, sort.SearchInts) }

func TestSearchSmall(t *testing.T) {

	rand.Seed(0)

	const limit = 10000

	var ints []int

	for i := 0; i < limit; i++ {
		ints = append(ints, rand.Int())
	}

	sort.Ints(ints)

	for want, q := range ints {
		if idx := Search(ints, q); idx != want {
			t.Errorf("Search(ints, %v)=%v, want %v", q, idx, want)
		}
	}

	for i := 0; i < 10000; i++ {

		q := rand.Int()

		want := sort.SearchInts(ints, q)

		if idx := Search(ints, q); idx != want {
			t.Errorf("Search(ints, %v)=%v, want %v", q, idx, want)
		}
	}
}

func TestSearchBig(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping big test search")
	}

	rand.Seed(0)

	const maxLimit = 100e6

	ints := make([]int, maxLimit)

	for i := 0; i < maxLimit; i++ {
		ints[i] = rand.Int()
	}

	sort.Ints(ints)

	for i := 0; i < 100000000; i++ {
		elt := rand.Int()
		vidx := Search(ints, elt)
		idx := sort.SearchInts(ints, elt)
		if vidx != idx {
			t.Fatalf("Search failed for elt=%d: got %d want %d\n", elt, vidx, idx)
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
