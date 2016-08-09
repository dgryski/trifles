package strtable

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

var strData [][]byte

type Inserter interface {
	Insert([]byte, uint32) (uint32, bool)
}

func Xorshift32(y uint32) uint32 {
	y ^= (y << 13)
	y ^= (y >> 17)
	y ^= (y << 5)
	return y
}

func testInsert(t *testing.T, which string, creator func(int) Inserter) {

	n := NewNative(8)
	m := creator(8)

	x := uint32(0xdeadbeef)

	for i := 0; i < 1e6; i++ {

		k := []byte(fmt.Sprintf("%d", Xorshift32(uint32(i))%8192))
		x := Xorshift32(x)

		vn, okn := n.Insert(k, x)
		vm, okm := m.Insert(k, x)

		if vn != vm || okn != okm {
			t.Fatalf("%v.%d: Insert(%q,%x)=%v,%t, want %v,%t\n", which, i, k, x, vm, okm, vn, okn)
		}
	}
}

func TestNative(t *testing.T) { testInsert(t, "native", nativeInserter) }
func TestTable(t *testing.T)  { testInsert(t, "table", tableInserter) }
func TestBTable(t *testing.T) { testInsert(t, "btable", btableInserter) }
func TestFrog(t *testing.T)   { testInsert(t, "frog", frogInserter) }
func TestRH(t *testing.T)     { testInsert(t, "rh", rhInserter) }

func benchmarkInsertSome(b *testing.B, size uint32, creator func(int) Inserter) {

	if strData == nil {
		loadStringData("/home/dgryski/strings.out")
	}

	b.ResetTimer()

	var total uint32
	var hits int
	for i := 0; i < b.N; i++ {

		ht := creator(int(size / 2))

		for j := uint32(0); j < size; j++ {
			v, ok := ht.Insert(strData[j], j)
			total += v
			if ok {
				hits += 1
			}
		}
	}

	b.StopTimer()

	//	b.Log("N=", b.N, "size=", size, "total=", total, "hits=", hits)
}

func benchmarkInsertAll(b *testing.B, size uint32, creator func(int) Inserter) {

	if strData == nil {
		loadStringData("/home/dgryski/strings.out")
	}

	size *= 2

	if size > uint32(len(strData)) {
		size = uint32(len(strData))
	}

	b.ResetTimer()

	var total uint32
	var hits int
	for i := 0; i < b.N; i++ {

		ht := creator(int(size))

		for j, s := range strData[:size] {
			v, ok := ht.Insert(s, uint32(j))
			total += v
			if ok {
				hits += 1
			}
		}
	}

	b.StopTimer()

	// b.Log("N=", b.N, "size=", size, "total=", total, "hits=", hits)
}

func nativeInserter(size int) Inserter { return NewNative(size) }
func tableInserter(size int) Inserter  { return New(size) }
func btableInserter(size int) Inserter { return NewBucket(size) }
func frogInserter(size int) Inserter   { return NewFrog(size) }
func rhInserter(size int) Inserter     { return NewRH(size) }

func BenchmarkNative1024(b *testing.B) { benchmarkInsertAll(b, 1024, nativeInserter) }
func BenchmarkTable1024(b *testing.B)  { benchmarkInsertAll(b, 1024, tableInserter) }
func BenchmarkBTable1024(b *testing.B) { benchmarkInsertAll(b, 1024, btableInserter) }

func BenchmarkNative64k(b *testing.B) { benchmarkInsertAll(b, 1<<16, nativeInserter) }
func BenchmarkTable64k(b *testing.B)  { benchmarkInsertAll(b, 1<<16, tableInserter) }
func BenchmarkBTable64k(b *testing.B) { benchmarkInsertAll(b, 1<<16, btableInserter) }

func BenchmarkNative256k(b *testing.B) { benchmarkInsertAll(b, 1<<18, nativeInserter) }
func BenchmarkTable256k(b *testing.B)  { benchmarkInsertAll(b, 1<<18, tableInserter) }
func BenchmarkBTable256k(b *testing.B) { benchmarkInsertAll(b, 1<<18, btableInserter) }

func BenchmarkNative1M(b *testing.B) { benchmarkInsertAll(b, 1<<20, nativeInserter) }
func BenchmarkTable1M(b *testing.B)  { benchmarkInsertAll(b, 1<<20, tableInserter) }
func BenchmarkBTable1M(b *testing.B) { benchmarkInsertAll(b, 1<<20, btableInserter) }

func BenchmarkNative2M(b *testing.B) { benchmarkInsertAll(b, 1<<21, nativeInserter) }
func BenchmarkTable2M(b *testing.B)  { benchmarkInsertAll(b, 1<<21, tableInserter) }
func BenchmarkBTable2M(b *testing.B) { benchmarkInsertAll(b, 1<<21, btableInserter) }

func BenchmarkSomeNative1024(b *testing.B) { benchmarkInsertSome(b, 1024, nativeInserter) }
func BenchmarkSomeTable1024(b *testing.B)  { benchmarkInsertSome(b, 1024, tableInserter) }
func BenchmarkSomeBTable1024(b *testing.B) { benchmarkInsertSome(b, 1024, btableInserter) }
func BenchmarkSomeFrog1024(b *testing.B)   { benchmarkInsertSome(b, 1024, frogInserter) }
func BenchmarkSomeRH1024(b *testing.B)     { benchmarkInsertSome(b, 1024, rhInserter) }

func BenchmarkSomeNative64k(b *testing.B) { benchmarkInsertSome(b, 1<<16, nativeInserter) }
func BenchmarkSomeTable64k(b *testing.B)  { benchmarkInsertSome(b, 1<<16, tableInserter) }
func BenchmarkSomeBTable64k(b *testing.B) { benchmarkInsertSome(b, 1<<16, btableInserter) }
func BenchmarkSomeFrog64k(b *testing.B)   { benchmarkInsertSome(b, 1<<16, frogInserter) }
func BenchmarkSomeRH64k(b *testing.B)     { benchmarkInsertSome(b, 1<<16, rhInserter) }

func BenchmarkSomeNative256k(b *testing.B) { benchmarkInsertSome(b, 1<<18, nativeInserter) }
func BenchmarkSomeTable256k(b *testing.B)  { benchmarkInsertSome(b, 1<<18, tableInserter) }
func BenchmarkSomeBTable256k(b *testing.B) { benchmarkInsertSome(b, 1<<18, btableInserter) }
func BenchmarkSomeFrog256k(b *testing.B)   { benchmarkInsertSome(b, 1<<18, frogInserter) }
func BenchmarkSomeRH256k(b *testing.B)     { benchmarkInsertSome(b, 1<<18, rhInserter) }

func BenchmarkSomeNative1M(b *testing.B) { benchmarkInsertSome(b, 1<<20, nativeInserter) }
func BenchmarkSomeTable1M(b *testing.B)  { benchmarkInsertSome(b, 1<<20, tableInserter) }
func BenchmarkSomeBTable1M(b *testing.B) { benchmarkInsertSome(b, 1<<20, btableInserter) }
func BenchmarkSomeFrog1M(b *testing.B)   { benchmarkInsertSome(b, 1<<20, frogInserter) }
func BenchmarkSomeRH1M(b *testing.B)     { benchmarkInsertSome(b, 1<<20, rhInserter) }

/*
func BenchmarkSomeNative2M(b *testing.B) { benchmarkInsertSome(b, 1<<21, nativeInserter) }
func BenchmarkSomeTable2M(b *testing.B)  { benchmarkInsertSome(b, 1<<21, tableInserter) }
func BenchmarkSomeBTable2M(b *testing.B) { benchmarkInsertSome(b, 1<<21, btableInserter) }
*/

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
