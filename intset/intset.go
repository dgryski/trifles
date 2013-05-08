package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dgryski/dbitstream"
	"math/rand"
	"sort"
	"time"
)

const SIZE = 1024

func vint_encode(enc []uint8, n uint) []uint8 {

	for n >= 0x80 {
		b := byte(n) | 0x80
		enc = append(enc, b)
		n >>= 7
	}

	return append(enc, byte(n))
}

func vintEncodeArray(numbers []int) []uint8 {

	enc := make([]uint8, 0, SIZE)

	for _, n := range numbers {
		enc = vint_encode(enc, uint(n))
	}
	return enc
}

func vintDecodeArray(enc []uint8) []int {

	numbers := make([]int, 0, SIZE)

	s := uint(0)
	n := 0
	for _, b := range enc {
		n |= int(b&0x7f) << s
		s += 7

		if (b & 0x80) == 0 {
			numbers = append(numbers, n)
			s = 0
			n = 0
		}
	}
	return numbers
}

func deltaEncodeArray(numbers []int) []uint8 {

	enc := make([]uint8, 0, SIZE)

	p := 0

	for _, n := range numbers {
		enc = vint_encode(enc, uint(n-p))
		p = n
	}
	return enc
}

func deltaDecodeArray(enc []uint8) []int {

	numbers := make([]int, 0, SIZE)

	prev := 0

	s := uint(0)
	n := 0
	for _, b := range enc {
		n |= int(b&0x7f) << s
		s += 7

		if (b & 0x80) == 0 {
			prev += n
			numbers = append(numbers, prev)
			s = 0
			n = 0
		}
	}
	return numbers
}

const M = 1

func riceEncodeArray(numbers []int) []uint8 {

	// first entry is numbers[0], little-endian

	buf := bytes.NewBuffer(nil)

	binary.Write(buf, binary.LittleEndian, int32(numbers[0]))

	bw := dbitstream.NewWriter(buf)

	l := len(numbers)

	for i := 1; i < l; i++ {
		var n = numbers[i] - numbers[i-1]

		q := (n - 1) / (1 << M)

		for i := 0; i < q; i++ {
			bw.WriteBit(dbitstream.Zero)
		}
		bw.WriteBit(dbitstream.One)

		r := (n - 1) & ((1 << M) - 1)

		for i := M; i >= 0; i-- {
			var b dbitstream.Bit
			b = r&(1<<uint(i)) != 0
			bw.WriteBit(b)
		}
	}
	bw.Flush(dbitstream.Zero)
	return buf.Bytes()
}

func riceDecodeArray(l int, enc []uint8) []int {

	buf := bytes.NewBuffer(enc)

	// first entry is numbers[0], little-endian int32
	var i32 int32
	binary.Read(buf, binary.LittleEndian, &i32)
	numbers := []int{int(i32)}

	br := dbitstream.NewReader(buf)

	prev := numbers[0]

	for i := 1; i < l; i++ {
		var nr, q int
		for {
			b, _ := br.ReadBit()
			if b {
				break
			}
			q++
		}

		for j := M; j >= 0; j-- {
			if b, _ := br.ReadBit(); b {
				nr += 1 << uint(j)
			}
		}

		nr += 1 + q*(1<<M)

		prev += nr

		numbers = append(numbers, prev)

	}
	return numbers
}

/*
func encodeFrameList(numbers [SIZE]int) []uint8 {
   // first, turn numbers into a bitset

   var bitset []uint8;
   n = 0;

   state = enum { ZEROS, ONES };

   while n < len(bitset) {

     if (bitset[n] == 0xff) {
        state = ONES;
        count = 8;

        n++;
        while(bitset[n] == 0xff) {
           count += 8;
        }
     }
   }
}
*/

func compare(dec, numbers []int) {

	for i := 0; i < SIZE; i++ {
		if dec[i] != numbers[i] {
			fmt.Println("found mismatch: ", i, " => ", dec[i], " and ", numbers[i])
		}
	}

}

func main() {

	rand.Seed(time.Now().UnixNano())

	var numbers [SIZE]int

	quality := float32(0.4)

	start := rand.Int31n(900000) + 100000

	seq := 0
	for frame := start; seq < SIZE; frame++ {
		/*
		   numbers[seq] = int(rand.Int31n(10000))
		   seq++;
		*/
		// /*
		if rand.Float32() < quality {
			numbers[seq] = int(frame)
			seq++
		}
		// */
	}

	sort.Ints(numbers[:])

	var enc []byte
	var dec []int

	t0 := time.Now()

	t0 = time.Now()
	enc = vintEncodeArray(numbers[:])
	dec = vintDecodeArray(enc)
	fmt.Println("t = ", time.Since(t0))
	compare(dec, numbers[:])

	fmt.Println("vint")
	fmt.Println("size of original    : ", 4*SIZE)
	fmt.Println("size of encoded data: ", len(enc))
	fmt.Println()

	t0 = time.Now()
	enc = deltaEncodeArray(numbers[:])
	dec = deltaDecodeArray(enc)
	fmt.Println("t = ", time.Since(t0))
	compare(dec, numbers[:])

	fmt.Println("delta")
	fmt.Println("size of original    : ", 4*SIZE)
	fmt.Println("size of encoded data: ", len(enc))
	fmt.Println()

	t0 = time.Now()
	enc = riceEncodeArray(numbers[:])
	dec = riceDecodeArray(SIZE, enc)
	fmt.Println("t = ", time.Since(t0))
	compare(dec, numbers[:])

	fmt.Println("rice")
	fmt.Println("size of original    : ", 4*SIZE)
	fmt.Println("size of encoded data: ", len(enc))
}
