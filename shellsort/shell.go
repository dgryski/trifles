package shellsort

import "github.com/dgryski/go-pcgr"
import "time"

var gaps = []int{701, 301, 132, 57, 23, 10, 4, 1}

func Shell(a []int) {
	// Sort an array a[0...n-1].

	n := len(a)

	// Start with the largest gap and work down to a gap of 1
	for _, gap := range gaps {
		//  Do a gapped insertion sort for this gap size.
		// The first gap elements a[0..gap-1] are already in gapped order
		// keep adding one more element until the entire array is gap sorted
		for i := gap; i < n; i++ {
			// add a[i] to the elements that have been gap sorted
			// save a[i] in temp and make a hole at position i
			temp := a[i]
			// shift earlier gap-sorted elements up until the correct location for a[i] is found
			var j int
			for j = i; j >= gap && a[j-gap] > temp; j -= gap {
				a[j] = a[j-gap]
			}
			// put temp (the original a[i]) in its correct location
			a[j] = temp
		}
	}
}

func permuteRandom(array []int, r *pcgr.Rand) {
	for i := len(array) - 1; i >= 1; i-- {
		j := r.Bound(uint32(i + 1))
		array[i], array[j] = array[j], array[i]
	}
}

const _C = 16

// compare-exchange two regions of length offset each
func compareRegions(a []int, s, t, offset, round int, r *pcgr.Rand) {
	mate := make([]int, offset)
	for count := 0; count < _C; count++ {
		for i := 0; i < offset; i++ {
			mate[i] = i
		}
		permuteRandom(mate, r)
		for i := 0; i < offset; i++ {

			i2 := s + i
			j2 := t + mate[i]

			if i2 >= len(a) || j2 >= len(a) || r.Bound(uint32(round)) == 0 {
				continue
			}

			if ((i2 < j2) && (a[i2] > a[j2])) || ((i2 > j2) && (a[i2] < a[j2])) {
				a[i2], a[j2] = a[j2], a[i2]
			}
		}
	}
}

func RandShell(a []int) {

	r := &pcgr.Rand{}
	r.Seed(time.Now().UnixNano())

	n := len(a)
	round := 1

	for _, offset := range gaps {

		for i := 0; i < n-offset; i += offset { // compare-exchange up
			compareRegions(a, i, i+offset, offset, round, r)
		}
		for i := n - offset; i >= offset; i -= offset { // compare-exchange down
			compareRegions(a, i-offset, i, offset, round, r)
		}
		for i := 0; i < n-3*offset; i += offset { // compare 3 hops up
			compareRegions(a, i, i+3*offset, offset, round, r)
		}
		for i := 0; i < n-2*offset; i += offset { // compare 2 hops up
			compareRegions(a, i, i+2*offset, offset, round, r)
		}
		for i := 0; i < n; i += 2 * offset { // compare odd-even regions
			compareRegions(a, i, i+offset, offset, round, r)
		}
		for i := offset; i < n-offset; i += 2 * offset { // compare even-odd regions
			compareRegions(a, i, i+offset, offset, round, r)
		}

		round++
	}
}
