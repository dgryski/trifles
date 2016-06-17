package main

import (
	"hash/fnv"
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

func BenchmarkSpooky8(b *testing.B)  { benchmarkHashn(b, 8, hspooky) }
func BenchmarkSpooky16(b *testing.B) { benchmarkHashn(b, 16, hspooky) }
func BenchmarkSpooky40(b *testing.B) { benchmarkHashn(b, 40, hspooky) }
func BenchmarkSpooky64(b *testing.B) { benchmarkHashn(b, 64, hspooky) }
func BenchmarkSpooky1K(b *testing.B) { benchmarkHashn(b, 1024, hspooky) }
func BenchmarkSpooky8K(b *testing.B) { benchmarkHashn(b, 8192, hspooky) }

var hsiphash = func(k []byte) uint64 { return siphash.Hash(0, 0, k) }

func BenchmarkSiphash8(b *testing.B)  { benchmarkHashn(b, 8, hsiphash) }
func BenchmarkSiphash16(b *testing.B) { benchmarkHashn(b, 16, hsiphash) }
func BenchmarkSiphash40(b *testing.B) { benchmarkHashn(b, 40, hsiphash) }
func BenchmarkSiphash64(b *testing.B) { benchmarkHashn(b, 64, hsiphash) }
func BenchmarkSiphash1K(b *testing.B) { benchmarkHashn(b, 1024, hsiphash) }
func BenchmarkSiphash8K(b *testing.B) { benchmarkHashn(b, 8192, hsiphash) }

var hfarm = func(k []byte) uint64 { return farm.Hash64(k) }

func BenchmarkFarm8(b *testing.B)  { benchmarkHashn(b, 8, hfarm) }
func BenchmarkFarm16(b *testing.B) { benchmarkHashn(b, 16, hfarm) }
func BenchmarkFarm40(b *testing.B) { benchmarkHashn(b, 40, hfarm) }
func BenchmarkFarm64(b *testing.B) { benchmarkHashn(b, 64, hfarm) }
func BenchmarkFarm1K(b *testing.B) { benchmarkHashn(b, 1024, hfarm) }
func BenchmarkFarm8K(b *testing.B) { benchmarkHashn(b, 8192, hfarm) }

var hcfarmhash = func(k []byte) uint64 { return cfarmhash.Hash64(k) }

func BenchmarkCFarmhash8(b *testing.B)  { benchmarkHashn(b, 8, hcfarmhash) }
func BenchmarkCFarmhash16(b *testing.B) { benchmarkHashn(b, 16, hcfarmhash) }
func BenchmarkCFarmhash40(b *testing.B) { benchmarkHashn(b, 40, hcfarmhash) }
func BenchmarkCFarmhash64(b *testing.B) { benchmarkHashn(b, 64, hcfarmhash) }
func BenchmarkCFarmhash1K(b *testing.B) { benchmarkHashn(b, 1024, hcfarmhash) }
func BenchmarkCFarmhash8K(b *testing.B) { benchmarkHashn(b, 8192, hcfarmhash) }

var hcity = func(k []byte) uint64 { return cityhash.CityHash64(k, uint32(len(k))) }

func BenchmarkCity8(b *testing.B)  { benchmarkHashn(b, 8, hcity) }
func BenchmarkCity16(b *testing.B) { benchmarkHashn(b, 16, hcity) }
func BenchmarkCity40(b *testing.B) { benchmarkHashn(b, 40, hcity) }
func BenchmarkCity64(b *testing.B) { benchmarkHashn(b, 64, hcity) }
func BenchmarkCity1K(b *testing.B) { benchmarkHashn(b, 1024, hcity) }
func BenchmarkCity8K(b *testing.B) { benchmarkHashn(b, 8192, hcity) }

var hmetro = func(k []byte) uint64 { return metro.Hash64(k, 0) }

func BenchmarkMetro8(b *testing.B)  { benchmarkHashn(b, 8, hmetro) }
func BenchmarkMetro16(b *testing.B) { benchmarkHashn(b, 16, hmetro) }
func BenchmarkMetro40(b *testing.B) { benchmarkHashn(b, 40, hmetro) }
func BenchmarkMetro64(b *testing.B) { benchmarkHashn(b, 64, hmetro) }
func BenchmarkMetro1K(b *testing.B) { benchmarkHashn(b, 1024, hmetro) }
func BenchmarkMetro8K(b *testing.B) { benchmarkHashn(b, 8192, hmetro) }

var hxxhash = func(k []byte) uint64 { return xxhash.Checksum64(k) }

func BenchmarkXXHash8(b *testing.B)  { benchmarkHashn(b, 8, hxxhash) }
func BenchmarkXXHash16(b *testing.B) { benchmarkHashn(b, 16, hxxhash) }
func BenchmarkXXHash40(b *testing.B) { benchmarkHashn(b, 40, hxxhash) }
func BenchmarkXXHash64(b *testing.B) { benchmarkHashn(b, 64, hxxhash) }
func BenchmarkXXHash1K(b *testing.B) { benchmarkHashn(b, 1024, hxxhash) }
func BenchmarkXXHash8K(b *testing.B) { benchmarkHashn(b, 8192, hxxhash) }

var fsthash = func(k []byte) uint64 { return fasthash.Hash64(0, k) }

func BenchmarkFasthash8(b *testing.B)  { benchmarkHashn(b, 8, fsthash) }
func BenchmarkFasthash16(b *testing.B) { benchmarkHashn(b, 16, fsthash) }
func BenchmarkFasthash40(b *testing.B) { benchmarkHashn(b, 40, fsthash) }
func BenchmarkFasthash64(b *testing.B) { benchmarkHashn(b, 64, fsthash) }
func BenchmarkFasthash1K(b *testing.B) { benchmarkHashn(b, 1024, fsthash) }
func BenchmarkFasthash8K(b *testing.B) { benchmarkHashn(b, 8192, fsthash) }

var high = func(k []byte) uint64 { return highway.Hash(highway.Lanes{}, k) }

func BenchmarkHighway8(b *testing.B)  { benchmarkHashn(b, 8, high) }
func BenchmarkHighway16(b *testing.B) { benchmarkHashn(b, 16, high) }
func BenchmarkHighway40(b *testing.B) { benchmarkHashn(b, 40, high) }
func BenchmarkHighway64(b *testing.B) { benchmarkHashn(b, 64, high) }
func BenchmarkHighway1K(b *testing.B) { benchmarkHashn(b, 1024, high) }
func BenchmarkHighway8K(b *testing.B) { benchmarkHashn(b, 8192, high) }

var fnv64 = func(k []byte) uint64 {
	h := fnv.New64a()
	h.Write(k)
	return h.Sum64()
}

func BenchmarkFNV8(b *testing.B)  { benchmarkHashn(b, 8, fnv64) }
func BenchmarkFNV16(b *testing.B) { benchmarkHashn(b, 16, fnv64) }
func BenchmarkFNV40(b *testing.B) { benchmarkHashn(b, 40, fnv64) }
func BenchmarkFNV64(b *testing.B) { benchmarkHashn(b, 64, fnv64) }
func BenchmarkFNV1K(b *testing.B) { benchmarkHashn(b, 1024, fnv64) }
func BenchmarkFNV8K(b *testing.B) { benchmarkHashn(b, 8192, fnv64) }

var total uint64

func benchmarkHashn(b *testing.B, size int64, h func([]byte) uint64) {
	b.SetBytes(size)
	for i := 0; i < b.N; i++ {
		total += h(buf[:size])
	}
}
