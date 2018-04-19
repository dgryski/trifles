package main

import (
	"math/rand"

	"github.com/dgryski/go-pcgr"
)

type probkeeper struct {
	r pcgr.Rand
	f float32
}

func newProbkeeper(f float32) *probkeeper {
	return &probkeeper{
		f:f,
		r:pcgr.New(rand.Int63(), rand.Int63()),
	}
}

func (p *probkeeper) allow(s string) bool {
	return p.r.Float32() < p.f
}
