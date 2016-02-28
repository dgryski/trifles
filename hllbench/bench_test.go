package main

import (
	"testing"

	eclesh "github.com/eclesh/hyperloglog"
	lytics "github.com/lytics/hll"
)

func BenchmarkLyticsPP(b *testing.B) {

	h := lytics.NewHll(14, 20)

	for i := 0; i < b.N; i++ {
		h.Add(uint64(i))
	}
}

func BenchmarkEclsh(b *testing.B) {

	h, _ := eclesh.New(1 << 14)

	for i := 0; i < b.N; i++ {
		h.Add(uint32(i))
	}
}
