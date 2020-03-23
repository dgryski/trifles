// Package inthash is integer hashing functions
/*
Taken from
    https://web.archive.org/web/20120720045250/http://www.cris.com/~Ttwang/tech/inthash.htm
    http://burtleburtle.net/bob/hash/integer.html

*/
package inthash

import "math/bits"

// Hash64 hashes a uint64
func Hash64(key uint64) uint64 {
	key = ^key + (key << 21) // key = (key << 21) - key - 1;
	key = key ^ (key >> 24)
	key = (key + (key << 3)) + (key << 8) // key * 265
	key = key ^ (key >> 14)
	key = (key + (key << 2)) + (key << 4) // key * 21
	key = key ^ (key >> 28)
	key = key + (key << 31)
	return key
}

// Hash64Inv is the inverse of Hash64.  Hash64Inv(Hash64(x)) == x.
func Hash64Inv(key uint64) uint64 {

	// From http://naml.us/blog/2012/03/inverse-of-a-hash-function

	tmp := uint64(0)

	// Invert key = key + (key << 31)
	tmp = key - (key << 31)
	key = key - (tmp << 31)

	// Invert key = key ^ (key >> 28)
	tmp = key ^ key>>28
	key = key ^ tmp>>28

	// Invert key *= 21
	key *= 14933078535860113213

	// Invert key = key ^ (key >> 14)
	tmp = key ^ key>>14
	tmp = key ^ tmp>>14
	tmp = key ^ tmp>>14
	key = key ^ tmp>>14

	// Invert key *= 265
	key *= 15244667743933553977

	// Invert key = key ^ (key >> 24)
	tmp = key ^ key>>24
	key = key ^ tmp>>24

	// Invert key = (~key) + (key << 21)
	tmp = ^key
	tmp = ^(key - (tmp << 21))
	tmp = ^(key - (tmp << 21))
	key = ^(key - (tmp << 21))

	return key
}

// Hash32 hashes a uint64
func Hash32(key uint32) uint32 {
	key = ^key + (key << 15) // key = (key << 15) - key - 1;
	key = key ^ (key >> 12)
	key = key + (key << 2)
	key = key ^ (key >> 4)
	key = key * 2057 // key = (key + (key << 3)) + (key << 11);
	key = key ^ (key >> 16)
	return key
}

// Jenkins32 is Robert Jenkins' 32-bit integer hash function
func Jenkins32(a uint32) uint32 {
	a = (a + 0x7ed55d16) + (a << 12)
	a = (a ^ 0xc761c23c) ^ (a >> 19)
	a = (a + 0x165667b1) + (a << 5)
	a = (a + 0xd3a2646c) ^ (a << 9)
	a = (a + 0xfd7046c5) + (a << 3)
	a = (a ^ 0xb55a4f09) ^ (a >> 16)
	return a
}

func JenkinsShift32(a uint32) uint32 {
	a -= (a << 6)
	a ^= (a >> 17)
	a -= (a << 9)
	a ^= (a << 4)
	a -= (a << 3)
	a ^= (a << 10)
	a ^= (a >> 15)
	return a
}

// Jenkins96 is Robert Jenkins' 96-bit integer hash function, mixing three uint32s into a single uint32 value.
func Jenkins96(a, b, c uint32) uint32 {

	a = a - b
	a = a - c
	a = a ^ (c >> 13)

	b = b - c
	b = b - a
	b = b ^ (a << 8)

	c = c - a
	c = c - b
	c = c ^ (b >> 13)

	a = a - b
	a = a - c
	a = a ^ (c >> 12)

	b = b - c
	b = b - a
	b = b ^ (a << 16)

	c = c - a
	c = c - b
	c = c ^ (b >> 5)

	a = a - b
	a = a - c
	a = a ^ (c >> 3)

	b = b - c
	b = b - a
	b = b ^ (a << 10)

	c = c - a
	c = c - b
	c = c ^ (b >> 15)

	return c
}

// https://en.wikipedia.org/wiki/Xorshift

// Xorshift32 is an xorshift RNG
func Xorshift32(y uint32) uint32 {

	// http://www.jstatsoft.org/v08/i14/paper
	// Marasaglia's "favourite"

	y ^= (y << 13)
	y ^= (y >> 17)
	y ^= (y << 5)
	return y
}

func LecuyerPanneton32(y uint32) uint32 {

	// From the paper "On the xorshift random number generators"
	// http://www.iro.umontreal.ca/~lecuyer/myftp/papers/xorshift.pdf
	// This sequence had the best equidistribution property

	y ^= (y << 7)
	y ^= (y >> 1)
	y ^= (y << 9)
	return y
}

// 64-bit xorshift multiply rng from http://vigna.di.unimi.it/ftp/papers/xorshift.pdf
func XorshiftMult64(x uint64) uint64 {
	x ^= x >> 12 // a
	x ^= x << 25 // b
	x ^= x >> 27 // c
	return x * 2685821657736338717
}

func Murmur3_64_Finalizer(key uint64) uint64 {
	key ^= key >> 33
	key *= 0xff51afd7ed558ccd
	key ^= key >> 33
	key *= 0xc4ceb9fe1a85ec53
	key ^= key >> 33
	return key
}

// https://nullprogram.com/blog/2018/07/31/

// exact bias: 0.17353355999581582
func Lowbias32(x uint32) uint32 {
	x ^= x >> 16
	x *= 0x7feb352d
	x ^= x >> 15
	x *= 0x846ca68b
	x ^= x >> 16
	return x
}

// inverse
func Lowbias32_r(x uint32) uint32 {
	x ^= x >> 16
	x *= 0x43021123
	x ^= x>>15 ^ x>>30
	x *= 0x1d69e2a5
	x ^= x >> 16
	return x
}

// exact bias: 0.020888578919738908
func Triple32(x uint32) uint32 {
	x ^= x >> 17
	x *= 0xed5ad4bb
	x ^= x >> 11
	x *= 0xac4c1b51
	x ^= x >> 15
	x *= 0x31848bab
	x ^= x >> 14
	return x
}

// http://mostlymangling.blogspot.com/2019/01/better-stronger-mixer-and-test-procedure.html
func Rrxmrrxmsx_0(v uint64) uint64 {
	v ^= bits.RotateLeft64(v, -25) ^ bits.RotateLeft64(v, -50)
	v *= 0xA24BAED4963EE407
	v ^= bits.RotateLeft64(v, -24) ^ bits.RotateLeft64(v, -49)
	v *= 0x9FB21C651E98DF25
	return v ^ v>>28
}

// https://chromium.googlesource.com/chromium/blink/+/master/Source/wtf/HashFunctions.h#78
func Chrome64(key uint64) uint64 {
	key += ^(key << 32)
	key ^= (key >> 22)
	key += ^(key << 13)
	key ^= (key >> 8)
	key += (key << 3)
	key ^= (key >> 15)
	key += ^(key << 27)
	key ^= (key >> 31)
	return key
}
