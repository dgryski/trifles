package main

import (
	"testing"

	xxhash "github.com/OneOfOne/xxhash/native"
	"github.com/dchest/siphash"
	"github.com/dgryski/go-farm"
	cfarmhash "github.com/dgryski/go-farmhash"
	"github.com/dgryski/go-highway"
	"github.com/dgryski/go-metro"
	"github.com/dgryski/go-spooky"
	"github.com/opennota/fasthash"
	"github.com/surge/cityhash"
)

// Written in 2012 by Dmitry Chestnykh.
// Expanded to include multiple hash functions by Damian Gryski, 2015.
//
// To the extent possible under law, the author have dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
// http://creativecommons.org/publicdomain/zero/1.0/

var (
	key0, key1 uint64
	buf        = make([]byte, 8<<10)
)

var hspooky = func(k []byte) uint64 { return spooky.Hash64(k) }

func BenchmarkSpooky8(b *testing.B)  { benchmarkHash8(b, hspooky) }
func BenchmarkSpooky16(b *testing.B) { benchmarkHash16(b, hspooky) }
func BenchmarkSpooky40(b *testing.B) { benchmarkHash40(b, hspooky) }
func BenchmarkSpooky64(b *testing.B) { benchmarkHash64(b, hspooky) }
func BenchmarkSpooky1K(b *testing.B) { benchmarkHash1K(b, hspooky) }
func BenchmarkSpooky8K(b *testing.B) { benchmarkHash8K(b, hspooky) }

var hsiphash = func(k []byte) uint64 { return siphash.Hash(0, 0, k) }

func BenchmarkSiphash8(b *testing.B)  { benchmarkHash8(b, hsiphash) }
func BenchmarkSiphash16(b *testing.B) { benchmarkHash16(b, hsiphash) }
func BenchmarkSiphash40(b *testing.B) { benchmarkHash40(b, hsiphash) }
func BenchmarkSiphash64(b *testing.B) { benchmarkHash64(b, hsiphash) }
func BenchmarkSiphash1K(b *testing.B) { benchmarkHash1K(b, hsiphash) }
func BenchmarkSiphash8K(b *testing.B) { benchmarkHash8K(b, hsiphash) }

var hfarm = func(k []byte) uint64 { return farm.Hash64(k) }

func BenchmarkFarm8(b *testing.B)  { benchmarkHash8(b, hfarm) }
func BenchmarkFarm16(b *testing.B) { benchmarkHash16(b, hfarm) }
func BenchmarkFarm40(b *testing.B) { benchmarkHash40(b, hfarm) }
func BenchmarkFarm64(b *testing.B) { benchmarkHash64(b, hfarm) }
func BenchmarkFarm1K(b *testing.B) { benchmarkHash1K(b, hfarm) }
func BenchmarkFarm8K(b *testing.B) { benchmarkHash8K(b, hfarm) }

var hcfarmhash = func(k []byte) uint64 { return cfarmhash.Hash64(k) }

func BenchmarkCFarmhash8(b *testing.B)  { benchmarkHash8(b, hcfarmhash) }
func BenchmarkCFarmhash16(b *testing.B) { benchmarkHash16(b, hcfarmhash) }
func BenchmarkCFarmhash40(b *testing.B) { benchmarkHash40(b, hcfarmhash) }
func BenchmarkCFarmhash64(b *testing.B) { benchmarkHash64(b, hcfarmhash) }
func BenchmarkCFarmhash1K(b *testing.B) { benchmarkHash1K(b, hcfarmhash) }
func BenchmarkCFarmhash8K(b *testing.B) { benchmarkHash8K(b, hcfarmhash) }

var hcity = func(k []byte) uint64 { return cityhash.CityHash64(k, uint32(len(k))) }

func BenchmarkCity8(b *testing.B)  { benchmarkHash8(b, hcity) }
func BenchmarkCity16(b *testing.B) { benchmarkHash16(b, hcity) }
func BenchmarkCity40(b *testing.B) { benchmarkHash40(b, hcity) }
func BenchmarkCity64(b *testing.B) { benchmarkHash64(b, hcity) }
func BenchmarkCity1K(b *testing.B) { benchmarkHash1K(b, hcity) }
func BenchmarkCity8K(b *testing.B) { benchmarkHash8K(b, hcity) }

var hmetro = func(k []byte) uint64 { return metro.Hash64_1(k, 0) }

func BenchmarkMetro8(b *testing.B)  { benchmarkHash8(b, hmetro) }
func BenchmarkMetro16(b *testing.B) { benchmarkHash16(b, hmetro) }
func BenchmarkMetro40(b *testing.B) { benchmarkHash40(b, hmetro) }
func BenchmarkMetro64(b *testing.B) { benchmarkHash64(b, hmetro) }
func BenchmarkMetro1K(b *testing.B) { benchmarkHash1K(b, hmetro) }
func BenchmarkMetro8K(b *testing.B) { benchmarkHash8K(b, hmetro) }

var hxxhash = func(k []byte) uint64 { return xxhash.Checksum64(k) }

func BenchmarkXXHash8(b *testing.B)  { benchmarkHash8(b, hxxhash) }
func BenchmarkXXHash16(b *testing.B) { benchmarkHash16(b, hxxhash) }
func BenchmarkXXHash40(b *testing.B) { benchmarkHash40(b, hxxhash) }
func BenchmarkXXHash64(b *testing.B) { benchmarkHash64(b, hxxhash) }
func BenchmarkXXHash1K(b *testing.B) { benchmarkHash1K(b, hxxhash) }
func BenchmarkXXHash8K(b *testing.B) { benchmarkHash8K(b, hxxhash) }

var fsthash = func(k []byte) uint64 { return fasthash.Hash64(0, k) }

func BenchmarkFasthash8(b *testing.B)  { benchmarkHash8(b, fsthash) }
func BenchmarkFasthash16(b *testing.B) { benchmarkHash16(b, fsthash) }
func BenchmarkFasthash40(b *testing.B) { benchmarkHash40(b, fsthash) }
func BenchmarkFasthash64(b *testing.B) { benchmarkHash64(b, fsthash) }
func BenchmarkFasthash1K(b *testing.B) { benchmarkHash1K(b, fsthash) }
func BenchmarkFasthash8K(b *testing.B) { benchmarkHash8K(b, fsthash) }

var high = func(k []byte) uint64 { return highway.Hash(highway.Lanes{}, k) }

func BenchmarkHighway8(b *testing.B)  { benchmarkHash8(b, high) }
func BenchmarkHighway16(b *testing.B) { benchmarkHash16(b, high) }
func BenchmarkHighway40(b *testing.B) { benchmarkHash40(b, high) }
func BenchmarkHighway64(b *testing.B) { benchmarkHash64(b, high) }
func BenchmarkHighway1K(b *testing.B) { benchmarkHash1K(b, high) }
func BenchmarkHighway8K(b *testing.B) { benchmarkHash8K(b, high) }

func benchmarkHash8(b *testing.B, h func([]byte) uint64) {
	b.SetBytes(8)
	for i := 0; i < b.N; i++ {
		h(buf[:8])
	}
}

func benchmarkHash16(b *testing.B, h func([]byte) uint64) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		h(buf[:16])
	}
}

func benchmarkHash40(b *testing.B, h func([]byte) uint64) {
	b.SetBytes(40)
	for i := 0; i < b.N; i++ {
		h(buf[:40])
	}
}

func benchmarkHash64(b *testing.B, h func([]byte) uint64) {
	b.SetBytes(64)
	for i := 0; i < b.N; i++ {
		h(buf[:64])
	}
}

func benchmarkHash128(b *testing.B, h func([]byte) uint64) {
	b.SetBytes(128)
	for i := 0; i < b.N; i++ {
		h(buf[:128])
	}
}

func benchmarkHash1K(b *testing.B, h func([]byte) uint64) {
	b.SetBytes(1024)
	for i := 0; i < b.N; i++ {
		h(buf[:1024])
	}
}

func benchmarkHash8K(b *testing.B, h func([]byte) uint64) {
	b.SetBytes(int64(len(buf)))
	for i := 0; i < b.N; i++ {
		h(buf)
	}
}
