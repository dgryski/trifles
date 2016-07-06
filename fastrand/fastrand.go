// Package fastrand provides a fast random number generator based on xorshift
/*
This sort of thing is frequently needed during benchmarks, and using math/rand
will make create a single-threaded bottleneck.  Using rand.New() is expensive
so always creating multiple sources isn't always an answer either.
*/
package fastrand

type Rng uint64

// Next returns the next random number
func (r *Rng) Next() uint64 {
	*r ^= *r >> 12 // a
	*r ^= *r << 25 // b
	*r ^= *r >> 27 // c
	return uint64(*r * 2685821657736338717)
}

// Intn returns a uniform random integer [0,bound)
func (r *Rng) Intn(bound uint64) uint64 {
	threshold := -bound % bound
	n := r.Next()
	for n < threshold {
		n = r.Next()
	}
	return n % bound
}

func (r *Rng) Int32n(n uint32) uint32 {

	bits := r.Next()

	random32bit := uint32(bits)
	bits >>= 32
	full := false

	multiresult := uint64(random32bit) * uint64(n)
	leftover := uint32(multiresult)
	if leftover < n {
		threshold := uint32(-n) % n
		for leftover < threshold {
			random32bit = uint32(bits)
			bits >>= 32
			if !full {
				bits = r.Next()
			}
			full = !full
			multiresult = uint64(random32bit) * uint64(n)
			leftover = uint32(multiresult)
		}
	}
	return uint32(multiresult >> 32) // [0, n)
}
