package rangebench

import (
	"strconv"
	"testing"
)

var buf = make([]int, 1024*1024)

var sink int

func benchmarkRangeIndex(b *testing.B, n int) {

	var sum = 0

	for i := 0; i < b.N; i++ {
		for idx := range buf[:n] {
			sum += buf[idx]
		}
	}

	sink += sum

}

func benchmarkRangeIndexValue(b *testing.B, n int) {

	var sum = 0

	for i := 0; i < b.N; i++ {
		for _, val := range buf[:n] {
			sum += val
		}
	}

	sink += sum
}

var benchSizes = []int{5, 10, 50, 100, 1024, 8192, 64 * 1024, 512 * 1024, 1024 * 1024}

func BenchmarkRangeIndexValue(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(strconv.Itoa(n), func(b *testing.B) { benchmarkRangeIndexValue(b, n) })
	}
}

func BenchmarkRangeIndex(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(strconv.Itoa(n), func(b *testing.B) { benchmarkRangeIndex(b, n) })
	}
}
