package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/dgryski/go-bitstream"
)

const SIZE = 1024

func vint_encode(enc []uint8, n uint32) []uint8 {

	for n >= 0x80 {
		b := byte(n) | 0x80
		enc = append(enc, b)
		n >>= 7
	}

	return append(enc, byte(n))
}

func vintEncodeArray(numbers []int32) []uint8 {

	enc := make([]uint8, 0, SIZE)

	for _, n := range numbers {
		enc = vint_encode(enc, uint32(n))
	}
	return enc
}

func vintDecodeArray(enc []uint8) []int32 {

	numbers := make([]int32, 0, SIZE)

	s := uint(0)
	var n int32
	for _, b := range enc {
		n |= int32(b&0x7f) << s
		s += 7

		if (b & 0x80) == 0 {
			numbers = append(numbers, n)
			s = 0
			n = 0
		}
	}
	return numbers
}

func deltaEncodeArray(numbers []int32) []uint8 {

	enc := make([]uint8, 0, SIZE)

	var p int32

	for _, n := range numbers {
		enc = vint_encode(enc, uint32(n-p))
		p = n
	}
	return enc
}

func deltaDecodeArray(enc []uint8) []int32 {

	numbers := make([]int32, 0, SIZE)

	var prev int32

	s := uint(0)
	var n int32

	for _, b := range enc {
		n |= int32(b&0x7f) << s
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

func grvintEncodeArray(numbers []int32) []uint8 {

	enc := make([]uint8, len(numbers)*4+len(numbers)/4)

	var prev int32

	var offs uint32

	for len(numbers) >= 4 {
		idx := offs

		var bits uint8
		var n, b uint32

		offs++

		n = uint32(numbers[0] - prev)
		binary.LittleEndian.PutUint32(enc[offs:], n)
		b = 3 - nlz(n|1)/8
		bits |= byte(b)
		offs += b + 1

		n = uint32(numbers[1] - numbers[0])
		binary.LittleEndian.PutUint32(enc[offs:], n)
		b = 3 - nlz(n|1)/8
		bits |= byte(b) << 2
		offs += b + 1

		n = uint32(numbers[2] - numbers[1])
		binary.LittleEndian.PutUint32(enc[offs:], n)
		b = 3 - nlz(n|1)/8
		bits |= byte(b) << 4
		offs += b + 1

		n = uint32(numbers[3] - numbers[2])
		binary.LittleEndian.PutUint32(enc[offs:], n)
		b = 3 - nlz(n|1)/8
		bits |= byte(b) << 6
		offs += b + 1

		enc[idx] = bits

		prev = numbers[3]
		numbers = numbers[4:]
	}

	enc = enc[:offs]

	for _, n := range numbers {
		enc = vint_encode(enc, uint32(n-prev))
		prev = n
	}

	for i := 0; i < 4-len(numbers); i++ {
		enc = append(enc, 0)
	}

	return enc
}

var grvintMask = [4]uint32{0xff, 0xffff, 0xffffff, 0xffffffff}

func grvintDecodeArray(enc []uint8, size int) []int32 {

	numbers := make([]int32, 0, SIZE)

	var prev int32

	s := uint(0)
	var n int32

	for size >= 4 {
		b := enc[0]
		enc = enc[1:]

		n := binary.LittleEndian.Uint32(enc) & grvintMask[b&3]
		enc = enc[1+(b&3):]
		b >>= 2
		n1 := prev + int32(n)

		n = binary.LittleEndian.Uint32(enc) & grvintMask[b&3]
		enc = enc[1+(b&3):]
		b >>= 2
		n2 := n1 + int32(n)

		n = binary.LittleEndian.Uint32(enc) & grvintMask[b&3]
		enc = enc[1+(b&3):]
		b >>= 2
		n3 := n2 + int32(n)

		n = binary.LittleEndian.Uint32(enc) & grvintMask[b&3]
		enc = enc[1+(b&3):]
		n4 := n3 + int32(n)

		numbers = append(numbers, n1, n2, n3, n4)
		prev = n4

		size -= 4
	}

	if size > 0 {

		// varint decode
		for _, b := range enc {
			n |= int32(b&0x7f) << s
			s += 7

			if (b & 0x80) == 0 {
				prev += n
				numbers = append(numbers, prev)
				s = 0
				n = 0
				size--
				if size == 0 {
					break
				}
			}
		}
	}
	return numbers
}

const M = 1

func riceEncodeArray(numbers []int32) []uint8 {

	// first entry is numbers[0], little-endian

	buf := bytes.NewBuffer(nil)

	binary.Write(buf, binary.LittleEndian, int32(numbers[0]))

	bw := bitstream.NewWriter(buf)

	l := len(numbers)

	for i := 1; i < l; i++ {
		var n = numbers[i] - numbers[i-1]

		q := (n - 1) / (1 << M)

		for i := int32(0); i < q; i++ {
			bw.WriteBit(bitstream.Zero)
		}
		bw.WriteBit(bitstream.One)

		r := (n - 1) & ((1 << M) - 1)

		for i := M; i >= 0; i-- {
			var b bitstream.Bit
			b = r&(1<<uint(i)) != 0
			bw.WriteBit(b)
		}
	}
	bw.Flush(bitstream.Zero)
	return buf.Bytes()
}

func riceDecodeArray(l int, enc []uint8) []int32 {

	buf := bytes.NewBuffer(enc)

	// first entry is numbers[0], little-endian int32
	var i32 int32
	binary.Read(buf, binary.LittleEndian, &i32)
	numbers := []int32{i32}

	br := bitstream.NewReader(buf)

	prev := numbers[0]

	for i := 1; i < l; i++ {
		var nr, q int32
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

func appendUint32(b []byte, i int32) []byte {
	b = append(b, byte(i))
	b = append(b, byte(i>>8))
	b = append(b, byte(i>>16))
	b = append(b, byte(i>>24))
	return b
}

func FoREncode(input []int32) []byte {

	if len(input) == 0 {
		return nil
	}

	var b []byte

	b = appendUint32(b, input[0])

	prev := input[0]
	for _, v := range input[1:] {
		diff := v - prev
		if -127 <= diff && diff <= 127 {
			b = append(b, byte(diff))
		} else {
			b = append(b, 0x80)
			b = appendUint32(b, v)
		}
		prev = v
	}

	return b
}

func FoRDecode(input []byte) []int32 {

	if len(input) == 0 {
		return nil
	}

	var output []int32

	n := int32(binary.LittleEndian.Uint32(input[:]))
	input = input[4:]

	output = append(output, n)

	prev := n
	for len(input) > 0 {
		diff := input[0]
		input = input[1:]
		if diff == 0x80 {
			n = int32(binary.LittleEndian.Uint32(input[:]))
			input = input[4:]
		} else {
			n = prev + int32(int8(diff))
		}
		prev = n
		output = append(output, n)
	}

	return output
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

func compare(dec, numbers []int32) {
	for i, d := range dec {
		if d != numbers[i] {
			fmt.Println("found mismatch: ", i, " => ", d, " and ", numbers[i])
		}
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	var numbers [SIZE]int32

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
			numbers[seq] = frame
			seq++
		}
		// */
	}

	sort.Sort(Int32s(numbers[:]))

	var enc []byte
	var dec []int32

	t0 := time.Now()

	fmt.Println("vint")
	t0 = time.Now()
	enc = vintEncodeArray(numbers[:])
	dec = vintDecodeArray(enc)
	fmt.Println("t = ", time.Since(t0))
	compare(dec, numbers[:])

	fmt.Println("size of original    : ", 4*SIZE)
	fmt.Println("size of encoded data: ", len(enc))
	fmt.Println()

	fmt.Println("delta")
	t0 = time.Now()
	enc = deltaEncodeArray(numbers[:])
	dec = deltaDecodeArray(enc)
	fmt.Println("t = ", time.Since(t0))
	compare(dec, numbers[:])

	fmt.Println("size of original    : ", 4*SIZE)
	fmt.Println("size of encoded data: ", len(enc))
	fmt.Println()

	fmt.Println("grvint delta")
	t0 = time.Now()
	enc = grvintEncodeArray(numbers[:])
	dec = grvintDecodeArray(enc[:], len(numbers))
	fmt.Println("t = ", time.Since(t0))
	compare(dec, numbers[:])

	fmt.Println("size of original    : ", 4*SIZE)
	fmt.Println("size of encoded data: ", len(enc))
	fmt.Println()

	fmt.Println("rice")
	t0 = time.Now()
	enc = riceEncodeArray(numbers[:])
	dec = riceDecodeArray(SIZE, enc)
	fmt.Println("t = ", time.Since(t0))
	compare(dec, numbers[:])

	fmt.Println("size of original    : ", 4*SIZE)
	fmt.Println("size of encoded data: ", len(enc))
	fmt.Println()

	fmt.Println("frame-of-reference")
	t0 = time.Now()
	enc = FoREncode(numbers[:])
	dec = FoRDecode(enc[:])
	fmt.Println("t = ", time.Since(t0))
	compare(dec, numbers[:])

	fmt.Println("size of original    : ", 4*SIZE)
	fmt.Println("size of encoded data: ", len(enc))

	num := []int32{10001, 10002, 10001, 10005, 10003, 10400, 10500, 10300, 10200}
	t0 = time.Now()
	enc = FoREncode(num[:])
	dec = FoRDecode(enc[:])
	fmt.Println("t = ", time.Since(t0))
	compare(dec, num[:])

}

type Int32s []int32

func (i32s Int32s) Len() int           { return len(i32s) }
func (i32s Int32s) Less(i, j int) bool { return i32s[i] < i32s[j] }
func (i32s Int32s) Swap(i, j int)      { i32s[i], i32s[j] = i32s[j], i32s[i] }
