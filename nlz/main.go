package main

import "fmt"

func nlz2(x uint64) uint64 {
	var n uint64

	n = 1

	if (x >> 32) == 0 {
		n = n + 32
		x = x << 32
	}
	if (x >> (32 + 16)) == 0 {
		n = n + 16
		x = x << 16
	}

	if (x >> (32 + 16 + 8)) == 0 {
		n = n + 8
		x = x << 8
	}

	if (x >> (32 + 16 + 8 + 4)) == 0 {
		n = n + 4
		x = x << 4
	}

	if (x >> (32 + 16 + 8 + 4 + 2)) == 0 {
		n = n + 2
		x = x << 2
	}

	n = n - (x >> 63)
	return uint64(n)
}

func main() {

	// From http://www.hackersdelight.org/hdcodetxt/nlz.c.txt
	tests := []uint64{0, 32, 1, 31, 2, 30, 3, 30, 4, 29, 5, 29, 6, 29,
		7, 29, 8, 28, 9, 28, 16, 27, 32, 26, 64, 25, 128, 24, 255, 24, 256, 23,
		512, 22, 1024, 21, 2048, 20, 4096, 19, 8192, 18, 16384, 17, 32768, 16,
		65536, 15, 0x20000, 14, 0x40000, 13, 0x80000, 12, 0x100000, 11,
		0x200000, 10, 0x400000, 9, 0x800000, 8, 0x1000000, 7, 0x2000000, 6,
		0x4000000, 5, 0x8000000, 4, 0x0FFFFFFF, 4, 0x10000000, 3,
		0x3000FFFF, 2, 0x50003333, 1, 0x7FFFFFFF, 1, 0x80000000, 0,
		0xFFFFFFFF, 0}

	for i := 0; i < len(tests); i += 2 {
		if lz := nlz(tests[i]); lz != tests[i+1]+32 {
			fmt.Printf("failed for test %d -> nlz(%b)=%d, got %d\n", i/2, tests[i], tests[i+1]+32, lz)
		}
	}
}
