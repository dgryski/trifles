package main

import (
	"crypto/sha1"
	"encoding/binary"
	"hash/fnv"
	"strconv"
	"testing"

	"blainsmith.com/go/seahash"
	xxhash "github.com/OneOfOne/xxhash"
	xxhashfast "github.com/cespare/xxhash"
	dchestsip "github.com/dchest/siphash"
	"github.com/dgryski/go-farm"
	cfarmhash "github.com/dgryski/go-farmhash"
	"github.com/dgryski/go-highway"
	"github.com/dgryski/go-marvin32"
	"github.com/dgryski/go-metro"
	"github.com/dgryski/go-sip13"
	"github.com/dgryski/go-spooky"
	"github.com/dgryski/go-stadtx"
	"github.com/dgryski/go-t1ha"
	"github.com/dgryski/go-wyhash"
	"github.com/dgryski/go-xxh3"
	"github.com/mmcloughlin/meow"
	"github.com/rbastic/go-zaphod64"
	"github.com/surge/cityhash"
	zxxh3 "github.com/zeebo/xxh3"
	"golang.org/x/crypto/blake2b"

	tsip "github.com/dgryski/trifles/tsip/go"
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

func BenchmarkSpooky(b *testing.B) { benchmarkHash(b, "Spooky", hspooky) }

var hsiphash = func(k []byte) uint64 { return dchestsip.Hash(0, 0, k) }

func BenchmarkSipHash(b *testing.B) { benchmarkHash(b, "SipHash", hsiphash) }

var hfarm = func(k []byte) uint64 { return farm.Fingerprint64(k) }

func BenchmarkFarm(b *testing.B) { benchmarkHash(b, "Farm", hfarm) }

var hcfarmhash = func(k []byte) uint64 { return cfarmhash.Hash64(k) }

func BenchmarkCFarm(b *testing.B) { benchmarkHash(b, "CFarm", hcfarmhash) }

var hcity = func(k []byte) uint64 { return cityhash.CityHash64(k, uint32(len(k))) }

func BenchmarkCity(b *testing.B) { benchmarkHash(b, "City", hcity) }

var hmetro = func(k []byte) uint64 { return metro.Hash64(k, 0) }

func BenchmarkMetro(b *testing.B) { benchmarkHash(b, "Metro", hmetro) }

var hxxhash = func(k []byte) uint64 { return xxhash.Checksum64(k) }

func BenchmarkXXHash(b *testing.B) { benchmarkHash(b, "XXHash", hxxhash) }

var hxxhashfast = func(k []byte) uint64 { return xxhashfast.Sum64(k) }

func BenchmarkXXFast(b *testing.B) { benchmarkHash(b, "XXFast", hxxhashfast) }

var high = func(k []byte) uint64 { return highway.Hash(highway.Lanes{}, k) }

func BenchmarkHighway(b *testing.B) { benchmarkHash(b, "Highway", high) }

var marvin = func(k []byte) uint64 { return uint64(marvin32.Sum32(0, k)) }

func BenchmarkMarvin32(b *testing.B) { benchmarkHash(b, "Marvin32", marvin) }

var sip13hash = func(k []byte) uint64 { return sip13.Sum64(0, 0, k) }

func BenchmarkSip13Hash(b *testing.B) { benchmarkHash(b, "Sip13", sip13hash) }

var fnv64 = func(k []byte) uint64 {
	h := fnv.New64a()
	h.Write(k)
	return h.Sum64()
}

func BenchmarkFNV1(b *testing.B) { benchmarkHash(b, "fnv1a", fnv64) }

var ht1ha = func(k []byte) uint64 { return t1ha.Sum64(k, 0) }

func BenchmarkT1ha(b *testing.B) { benchmarkHash(b, "T1ha", ht1ha) }

var zaphodSeed zaphod64.State

var hzaphod64 = func(k []byte) uint64 { return uint64(zaphod64.HashWithState(&zaphodSeed, k, uint64(len(k)))) }

func BenchmarkZaphod64(b *testing.B) { benchmarkHash(b, "Zaphod64", hzaphod64) }

var stadtxState stadtx.State

var hstadtx = func(k []byte) uint64 { return stadtx.Hash(&stadtxState, k) }

func BenchmarkStadtx(b *testing.B) { benchmarkHash(b, "Stadtx", hstadtx) }

var htsip = func(k []byte) uint64 { return tsip.HashASM(0, 0, k) }

func BenchmarkTsip(b *testing.B) { benchmarkHash(b, "Tsip", htsip) }

var hwyhash = func(k []byte) uint64 { return wyhash.Hash(k, 0) }

func BenchmarkWyhash(b *testing.B) { benchmarkHash(b, "Wyhash", hwyhash) }

var hxxh3 = func(k []byte) uint64 { return xxh3.Hash(k, 0) }

func BenchmarkXXH3(b *testing.B) { benchmarkHash(b, "XXH3", hxxh3) }

var hzxxh3 = func(k []byte) uint64 { return zxxh3.Hash(k) }

func BenchmarkZeeboXXH3(b *testing.B) { benchmarkHash(b, "ZeeboXXH3", hzxxh3) }

var hseahash = func(k []byte) uint64 { return seahash.Sum64(k) }

func BenchmarkSeahash(b *testing.B) { benchmarkHash(b, "Seahash", hseahash) }

var hsha1 = func(k []byte) uint64 {
	b := sha1.Sum(k)
	return binary.LittleEndian.Uint64(b[:8])
}

func BenchmarkSHA1(b *testing.B) { benchmarkHash(b, "SHA1", hsha1) }

var hblake = func(k []byte) uint64 {
	b := blake2b.Sum256(k)
	return binary.LittleEndian.Uint64(b[:8])
}

func BenchmarkBlake2b(b *testing.B) { benchmarkHash(b, "Blake", hblake) }

var hmeow = func(k []byte) uint64 {
	b := meow.Checksum(0, k)
	return binary.LittleEndian.Uint64(b[:8])
}

func BenchmarkMeow(b *testing.B) { benchmarkHash(b, "Meow", hmeow) }

func benchmarkHash(b *testing.B, str string, h func([]byte) uint64) {
	var sizes = []int{1, 2, 3, 4, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 1024, 8192}
	for _, n := range sizes {
		b.Run(strconv.Itoa(n), func(b *testing.B) { benchmarkHashn(b, int64(n), h) })
	}
}

var total uint64

func benchmarkHashn(b *testing.B, size int64, h func([]byte) uint64) {
	b.SetBytes(size)
	for i := 0; i < b.N; i++ {
		total += h(buf[:size])
	}
}
