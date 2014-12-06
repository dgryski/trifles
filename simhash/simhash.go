/*
Package simhash implements the simhash document similarity hashing function.

By itself, this isn't that useful.  If I ever decide to create  more fully
featured SimHash implementation, it will get moved to its own repository.

http://www.cs.princeton.edu/courses/archive/spr04/cos598B/bib/CharikarEstim.pdf
http://infolab.stanford.edu/~manku/papers/07www-duplicates.pdf
http://irl.cse.tamu.edu/people/sadhan/papers/cikm2011.pdf

Left here for legacy reasons, but you probably want
http://github.com/dgryski/go-simstore now.
*/
package simhash

import (
	"github.com/dchest/siphash"
)

// SimHash returns a simhash value for the document returned by the scanner
func SimHash(scanner FeatureScanner) uint64 {
	var signs [64]int64

	for scanner.Scan() {
		h := siphash.Hash(0, 0, scanner.Bytes())

		for i := 0; i < 64; i++ {
			negate := int(h) & 1
			// if negate is 1, we will negate '-1', below
			r := (-1 ^ -negate) + negate
			signs[i] += int64(r)
			h >>= 1
		}
	}

	var shash uint64

	// TODO: can probably be done with SSE?
	for i := 63; i >= 0; i-- {
		shash <<= 1
		shash |= uint64(signs[i]>>63) & 1
	}

	return shash
}

func hammingDistance(v1 uint64, v2 uint64) int {
	// TODO: replace with magic multiplication or lookup?

	v := v1 ^ v2
	var c int
	for c = 0; v != 0; c++ {
		v &= v - 1
	}

	return c
}
