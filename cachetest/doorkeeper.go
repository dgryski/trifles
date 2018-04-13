package main

import (
	"github.com/dgryski/go-metro"
	"github.com/dgryski/go-bloomf"
)

type doorkeeper struct {
	bf *bloomf.BF
}

func newDoorkeeper(size int) *doorkeeper {
	return &doorkeeper{
		bf: bloomf.New(size*9, 0.01, func(b []byte) uint64 { return metro.Hash64(b, 0) }),
	}
}

func (d *doorkeeper) allow(s string) bool {
	if d == nil {
		return true
	}

	alreadyPresent := d.bf.Insert([]byte(s))
	if d.bf.Len() > d.bf.Cap() {
		d.bf.Reset()
	}
	return alreadyPresent
}
