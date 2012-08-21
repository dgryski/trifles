package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"runtime/pprof"
)

func find_longest_match(buf []uint8, s []uint8) (length int, offs int) {

	max_length := 0
	max_offs := 0

	s0 := s[0]

	for i, bc := range buf {

		if bc != s0 {
			continue
		}

		j := 1
		for j < len(s) && i+j < len(buf) && buf[i+j] == s[j] {
			j++
		}

		if j > max_length {
			max_length = j
			max_offs = i
		}
	}

	return max_length, max_offs
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func compress(buf []byte, w io.Writer) {

	var ctrlbits uint8
	var ctrln int

	data := make([]byte, 4, 4096)

	n := len(buf)

	data[0] = uint8(n & 0xff)
	data[1] = uint8(n >> 8 & 0xff)
	data[2] = uint8(n >> 16 & 0xff)
	data[3] = uint8(n >> 24 & 0xff)

	w.Write(data)

	data = data[0:0]

	fileoffs := 0

	for fileoffs < n {
		minoffs := max(fileoffs-4096, 0)
		l, o := find_longest_match(buf[minoffs:fileoffs], buf[fileoffs:])

		ctrlbits <<= 1
		ctrln++

		if l > 2 {

			minoffs1 := max(1+fileoffs-4096, 0)
			l1, o1 := find_longest_match(buf[minoffs1:fileoffs+1], buf[fileoffs+1:])

			if l1 > l {
				// found a better match if we pretend this char didn't match
				// this code is basically the 'else' clause and the bottom/top of the loop
				data = append(data, buf[fileoffs])

				if ctrln == 8 {
					// control word is full -- flush
					w.Write([]byte{ctrlbits})
					w.Write(data)
					data = data[0:0]
					ctrln = 0
					ctrlbits = 0
				}
				fileoffs++

				ctrlbits <<= 1
				ctrln++
				l = l1
				o = o1
			}

			// worth compressing
			var lbits int
			if l >= 15 {
				lbits = 0x0f
			} else {
				lbits = l
			}
			b := uint8((lbits & 0x0f << 4) | ((o & 0xf00) >> 8))
			data = append(data, b)
			b = uint8(o & 0x00ff)
			data = append(data, b)
			if lbits == 0x0f {
				b = uint8(l - 0x0f)
				data = append(data, b)
			}
			ctrlbits |= 1
		} else {
			l = 1
			data = append(data, buf[fileoffs])
		}

		if ctrln == 8 {
			// control word is full -- flush
			w.Write([]byte{ctrlbits})
			w.Write(data)
			data = data[0:0]
			ctrln = 0
			ctrlbits = 0
		}

		fileoffs += l
	}

	// flush any remaining data
	if ctrln != 0 {
		ctrlbits <<= uint(8 - ctrln)
		w.Write([]byte{ctrlbits})
		w.Write(data)
	}
}

func decompress(buf []byte, w io.Writer) {

	data := make([]byte, 0, 4096)

	n := len(buf)

	decompressed_size := int(buf[0]) | int(buf[1])<<8 | int(buf[2])<<16 | int(buf[3])<<24
	fileoffs := 4

	var ctrlbits uint8

	for fileoffs < n {
		ctrlbits = buf[fileoffs]
		fileoffs++

		for ctrln := 0; ctrln < 8; ctrln++ {
			if ctrlbits&0x80 == 0x80 {
				l := (int(buf[fileoffs]&0xf0) >> 4) | 0
				o := (int(buf[fileoffs]&0x0f) << 8) | int(buf[fileoffs+1])
				if l == 0x0f {
					l += int(buf[fileoffs+2])
					fileoffs++
				}

				minoffs := max(len(data)-4096, 0)
				data = append(data, data[minoffs+o:minoffs+o+l]...)
				fileoffs += 2
			} else {
				data = append(data, buf[fileoffs])
				fileoffs++
			}
			ctrlbits <<= 1
			if len(data) >= decompressed_size {
				break
			}
		}
	}

	w.Write(data)
}

func main() {

	var optDecompress = flag.Bool("d", false, "decompress")
	var optCpuProfile = flag.String("cpuprofile", "", "profile")

	flag.Parse()

	if *optCpuProfile != "" {
		f, err := os.Create(*optCpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var bufin [16 * 1024]uint8
	var buf []uint8

	for {
		n, err := os.Stdin.Read(bufin[:])
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("error during read: ", err)
		}
		buf = append(buf, bufin[0:n]...)
	}

	if *optDecompress {
		decompress(buf, os.Stdout)
	} else {
		stdout := bufio.NewWriter(os.Stdout)
		compress(buf, stdout)
		stdout.Flush()
	}

}
