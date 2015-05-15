// Package strtable is a fast hash table for string->uint32
package strtable

import "bytes"

type element struct {
	// NOTE(dgryski): don't reorder these.  Cacheline magic discovered by benchmarking
	hash uint32
	val  uint32
	key  []byte
}

type Table struct {
	t    []element
	keys int
}

func New(keys int) *Table {
	t := Table{t: make([]element, 2*keys)}
	return &t
}

func (t *Table) Insert(k []byte, val uint32) (uint32, bool) {

	mask := uint32(len(t.t) - 1)

	h := leveldbHash(k)
	//  you can hack your runtime to expose this..
	// h := uint32(runtime.BytesHash(k, 0)) & mask
	slot := h & mask

	for {
		if (t.t)[slot].key == nil {
			// doesn't exist -- add it
			(t.t)[slot].key = k
			(t.t)[slot].hash = h
			(t.t)[slot].val = val

			t.keys++
			if t.keys > len(t.t)/2 {
				t.double()
			}

			return val, false
		}

		// slot key is not nil
		if (t.t)[slot].hash == h && bytes.Equal((t.t)[slot].key, k) {
			return (t.t)[slot].val, true
		}

		slot = (slot + 1) & mask
	}
}

func (t *Table) double() {
	newTable := make([]element, 2*len(t.t))
	mask := uint32(len(newTable) - 1)

	for _, elt := range t.t {

		if elt.key == nil {
			continue
		}

		slot := elt.hash & mask

		var steps int

		for {
			if newTable[slot].key == nil {
				// found a spot
				newTable[slot] = elt
				break
			}

			slot = (slot + 1) & mask
			steps++
		}
	}

	t.t = newTable

}

// leveldb's bloom filter hash, a murmur-lite
func leveldbHash(b []byte) uint32 {

	const (
		seed = 0xbc9f1d34
		m    = 0xc6a4a793
	)

	h := uint32(seed) ^ uint32(len(b)*m)

	for ; len(b) >= 4; b = b[4:] {

		h += uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
		h *= m
		h ^= h >> 16
	}
	switch len(b) {
	case 3:
		h += uint32(b[2]) << 16
		fallthrough
	case 2:
		h += uint32(b[1]) << 8
		fallthrough
	case 1:
		h += uint32(b[0])
		h *= m
		h ^= h >> 24
	}

	return h
}

const bucketSize = 8

type bucket struct {
	hash [bucketSize]uint32
	keys [bucketSize][]byte
	vals [bucketSize]uint32
	next int
}

type BTable struct {
	b        []bucket
	overflow []bucket
	keys     int
}

func NewBucket(keys int) *BTable {
	t := BTable{b: make([]bucket, keys/2)}
	return &t
}

func (t *BTable) Insert(k []byte, val uint32) (uint32, bool) {

	mask := uint32(len(t.b) - 1)

	h := leveldbHash(k)
	slot := h & mask

	if h == 0 {
		h++
	}

	b := &t.b[slot]

	// Search the chain of buckets for this slot.
	for {
		for i := range b.hash {
			bh := b.hash[i]
			if bh == 0 {
				b.keys[i] = k
				b.hash[i] = h
				b.vals[i] = val

				t.keys++
				if t.keys > 6*len(t.b) {
					t.double()
				}

				return val, false
			}
			if bh == h && bytes.Equal(b.keys[i], k) {
				return b.vals[i], true
			}
		}

		b.next = len(t.overflow)
		t.overflow = append(t.overflow, bucket{})
		b = &t.overflow[b.next]
	}
}

func (t *BTable) double() {

	newTable := make([]bucket, 2*len(t.b))
	var overflow []bucket
	mask := uint32(len(newTable) - 1)

	for i := range t.b {
		slot := uint32(i)

		b := &t.b[slot]

		var lb = &newTable[slot]
		var lidx int

		var hb = &newTable[i+len(t.b)]
		var hidx int

	split:
		for b != nil {

			for i := range b.hash {
				bh := b.hash[i]
				if bh == 0 {
					break split
				}
				if bh&mask == slot {
					if lidx == bucketSize {
						lidx = 0
						lb.next = len(overflow)
						overflow = append(overflow, bucket{})
						lb = &overflow[b.next]
					}

					lb.hash[lidx] = bh
					lb.keys[lidx] = b.keys[i]
					lb.vals[lidx] = b.vals[i]
					lidx++
				} else {
					if hidx == bucketSize {
						hidx = 0
						hb.next = len(overflow)
						overflow = append(overflow, bucket{})
						hb = &overflow[b.next]
					}

					hb.hash[hidx] = bh
					hb.keys[hidx] = b.keys[i]
					hb.vals[hidx] = b.vals[i]
					hidx++
				}
			}
			b = &t.overflow[b.next]
		}
	}

	t.b = newTable
	t.overflow = overflow
}

type Native map[string]uint32

func NewNative(size int) Native {
	return Native(make(map[string]uint32, size))
}

func (n Native) Insert(x []byte, val uint32) (uint32, bool) {

	if val, ok := n[string(x)]; ok {
		return val, true
	}

	n[string(x)] = val
	return val, false
}
