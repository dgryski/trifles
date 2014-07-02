package strtable

import (
	"bufio"
	"os"
	"testing"
)

var strData [][]byte

type Inserter interface {
	Insert([]byte, uint32) (uint32, bool)
}

func benchmarkCommon(b *testing.B, size uint32, inserter Inserter) {

	if strData == nil {
		loadStringData("/home/dgryski/shuf.out")
	}

	b.ResetTimer()

	var total uint32
	var hits int
	limit := uint32(b.N)
	for i := uint32(0); i < limit; i++ {
		v, ok := inserter.Insert(strData[i&(size-1)], i)
		total += v
		if ok {
			hits += 1
		}
	}

	b.StopTimer()

	//b.Log("size=", size, "total=", total, "hits=", hits)
}

func BenchmarkNative1024(b *testing.B) { benchmarkCommon(b, 1024, NewNative()) }
func BenchmarkTable1024(b *testing.B)  { benchmarkCommon(b, 1024, New(1024)) }

func BenchmarkNative65k(b *testing.B) { benchmarkCommon(b, 1<<16, NewNative()) }
func BenchmarkTable65k(b *testing.B)  { benchmarkCommon(b, 1<<16, New(1<<16)) }

func BenchmarkNative1M(b *testing.B) { benchmarkCommon(b, 1<<20, NewNative()) }
func BenchmarkTable1M(b *testing.B)  { benchmarkCommon(b, 1<<20, New(1<<20)) }

func BenchmarkNative2M(b *testing.B) { benchmarkCommon(b, 1<<21, NewNative()) }
func BenchmarkTable2M(b *testing.B)  { benchmarkCommon(b, 1<<21, New(1<<21)) }

func loadStringData(path string) {

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		b := make([]byte, len(scanner.Bytes()))
		copy(b, scanner.Bytes())
		strData = append(strData, b)
	}
}
