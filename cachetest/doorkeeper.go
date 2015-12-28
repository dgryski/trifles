package main

import (
	"github.com/dchest/siphash"
	"github.com/dgryski/go-bloomf"
)

type doorkeeper struct {
	bf *bloomf.BF
}

func newDoorkeeper(size int) *doorkeeper {
	return &doorkeeper{
		bf: bloomf.New(size*9, 0.01, func(b []byte) uint64 { return siphash.Hash(0, 0, b) }),
	}
}

func (d *doorkeeper) allow(s string) bool {
	if d == nil {
		return true
	}

	ok := d.bf.InsertLookup([]byte(s))
	if d.bf.Len() > d.bf.Cap() {
		d.bf.Reset()
	}
	return ok
}
