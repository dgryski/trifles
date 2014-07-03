// Package strtable is a fast hash table for string->uint32
package strtable

import "bytes"

type element struct {
	// NOTE(dgryski): don't reorder these.  Cacheline magic discovered by benchmarking
	hash uint32
	val  uint32
	key  []byte
}

type Table []element

func New(keys int) *Table {
	t := Table(make([]element, 2*keys))
	return &t
}

func (t *Table) Insert(k []byte, val uint32) (uint32, bool) {

	mask := uint32(len(*t) - 1)

	h := leveldbHash(k) & mask
	//  you can hack your runtime to expose this..
	// h := uint32(runtime.BytesHash(k, 0)) & mask

	slot := h
	for {
		if (*t)[slot].key == nil {
			// doesn't exist -- add it
			(*t)[slot].key = k
			(*t)[slot].hash = h
			(*t)[slot].val = val
			return val, false
		}

		// slot key is not nil
		if (*t)[slot].hash == h && bytes.Equal((*t)[slot].key, k) {
			return (*t)[slot].val, true
		}

		slot = (slot + 1) & mask
	}
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
