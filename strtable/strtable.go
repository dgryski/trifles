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

		if b.next == 0 {
			b.next = len(t.overflow) + 1
			t.overflow = append(t.overflow, bucket{})
		}
		b = &t.overflow[b.next-1]
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
						lb.next = len(overflow) + 1
						overflow = append(overflow, bucket{})
						lb = &overflow[lb.next-1]
					}

					lb.hash[lidx] = bh
					lb.keys[lidx] = b.keys[i]
					lb.vals[lidx] = b.vals[i]
					lidx++
				} else {
					if hidx == bucketSize {
						hidx = 0
						hb.next = len(overflow) + 1
						overflow = append(overflow, bucket{})
						hb = &overflow[hb.next-1]
					}

					hb.hash[hidx] = bh
					hb.keys[hidx] = b.keys[i]
					hb.vals[hidx] = b.vals[i]
					hidx++
				}
			}
			if b.next == 0 {
				break
			}
			b = &t.overflow[b.next-1]
		}
	}

	t.b = newTable
	t.overflow = overflow
}

// leapfrog: preshing.com/20160314/leapfrog-probing/

type leapfrog struct {
	hash  uint32
	key   []byte
	val   uint32
	delta [2]byte
}

type FrogTable struct {
	t []leapfrog
}

func NewFrog(keys int) *FrogTable {
	t := make([]leapfrog, 2*keys)
	return &FrogTable{t: t}
}

func (t *FrogTable) Insert(k []byte, val uint32) (uint32, bool) {
	h := leveldbHash(k)
	if h == 0 {
		h++
	}

	for {
		if v, b, ok := insertWithHash(t.t, k, val, h); ok {
			return v, b
		}

		t.double()
	}
}

func insertWithHash(t []leapfrog, k []byte, val uint32, h uint32) (uint32, bool, bool) {
	mask := uint32(len(t) - 1)
	slot := h & mask

	if (t)[slot].hash == h && bytes.Equal((t)[slot].key, k) {
		return (t)[slot].val, true, true
	}

	didx := 0
	delta := uint32((t)[slot].delta[didx])

	for delta != 0 {
		slot = (slot + delta) & mask
		if (t)[slot].hash == h && bytes.Equal((t)[slot].key, k) {
			return (t)[slot].val, true, true
		}
		didx = 1
		delta = uint32((t)[slot].delta[didx])
	}

	// reached the end of chain for this bucket; key not found
	// linear scan to find the next slot
	old := slot
	for i := uint32(0); i < 256; i++ {
		slot = (old + i) & mask
		if (t)[slot].key == nil {
			(t)[old].delta[didx] = byte(i)
			(t)[slot].key = k
			(t)[slot].hash = h
			(t)[slot].val = val
			return val, false, true
		}
	}

	return 0, false, false
}

func (t *FrogTable) double() {

	newSize := 2 * len(t.t)

retry:
	for {
		newTable := make([]leapfrog, newSize)
		for i := range t.t {
			if t.t[i].key != nil {
				_, _, ok := insertWithHash(newTable, t.t[i].key, t.t[i].val, t.t[i].hash)
				if !ok {
					newSize *= 2
					continue retry
				}
			}
		}
		t.t = newTable
		break
	}
}

type rhelement struct {
	hash uint32
	val  uint32
	dist uint32
	key  []byte
}

type RHTable struct {
	t    []rhelement
	keys int
}

func NewRH(keys int) *RHTable {
	t := RHTable{t: make([]rhelement, 2*keys)}
	return &t
}

func (t *RHTable) Insert(k []byte, val uint32) (uint32, bool) {

	mask := uint32(len(t.t) - 1)

	h := leveldbHash(k)
	//  you can hack your runtime to expose this..
	// h := uint32(runtime.BytesHash(k, 0)) & mask
	slot := h & mask

	var dist uint32
	var originalKey = true
	var retval uint32 = val

	for {
		if (t.t)[slot].key == nil {
			// doesn't exist -- add it
			(t.t)[slot].key = k
			(t.t)[slot].hash = h
			(t.t)[slot].val = val
			(t.t)[slot].dist = dist

			t.keys++
			if t.keys > 3*len(t.t)/4 {
				t.double()
			}

			return retval, false
		}

		// slot key is not nil
		if originalKey && (t.t)[slot].hash == h && bytes.Equal((t.t)[slot].key, k) {
			return (t.t)[slot].val, true
		}

		if (t.t)[slot].dist < dist {
			originalKey = false
			dist, (t.t)[slot].dist = (t.t)[slot].dist, dist
			k, (t.t)[slot].key = (t.t)[slot].key, k
			h, (t.t)[slot].hash = (t.t)[slot].hash, h
			val, (t.t)[slot].val = (t.t)[slot].val, val
		}

		slot = (slot + 1) & mask
		dist++
	}
}

func (t *RHTable) double() {
	newTable := make([]rhelement, 2*len(t.t))
	mask := uint32(len(newTable) - 1)

	for _, elt := range t.t {

		if elt.key == nil {
			continue
		}

		slot := elt.hash & mask

		var dist uint32
		for {
			if newTable[slot].key == nil {
				// doesn't exist -- add it
				newTable[slot] = elt
				newTable[slot].dist = dist
				break
			}

			if newTable[slot].dist < dist {
				elt, newTable[slot] = newTable[slot], elt
				newTable[slot].dist = dist
				dist = elt.dist
			}

			slot = (slot + 1) & mask
			dist++
		}

	}
	t.t = newTable
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
