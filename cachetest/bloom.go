package main

import (
	"github.com/dgryski/go-metro"
	"github.com/dgryski/go-bloomf"
)

type bloomkeeper struct {
	bf bloomf.BF
}

func newBloomkeeper(size int) *bloomkeeper {
	return &bloomkeeper{
		bf: *bloomf.New(size*9, 0.01, func(b []byte) uint64 { return metro.Hash64(b, 0) }),
	}
}

func (b *bloomkeeper) allow(s string) bool {
	alreadyPresent := b.bf.Insert([]byte(s))
	if b.bf.Len() > b.bf.Cap() {
		b.bf.Reset()
	}
	return alreadyPresent
}
