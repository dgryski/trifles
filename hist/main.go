package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

func main() {

	values := make([]int, 0, 10000)
	hist := make(map[int]int)

	br := bufio.NewReader(os.Stdin)

FOR:
	for i := 0; ; i++ {
		b, err := br.ReadBytes('\n')
		switch err {
		case io.EOF:
			break FOR
		case nil:
			n, _ := strconv.Atoi(string(b[:len(b)-1]))
			values = append(values, n)
			v := hist[n]
			hist[n] = v + 1

		default:
			fmt.Printf("Error reading input: %v", err)
			os.Exit(1)
		}
	}

	sum := float64(0)
	sumsq := float64(0)

	min := int(math.MaxInt32)
	max := 0

	for _, v := range values {
		f := float64(v)
		sum += f
		sumsq += f * f

		if min > v {
			min = v
		}

		if max < v {
			max = v
		}

	}

	n := len(values)

	sort.Ints(values)

	mean := sum / float64(n)
	stddev := math.Sqrt(sumsq/float64(n) - mean*mean)

	var median int

	if n&1 == 1 {
		median = values[(n-1)/2]
	} else {
		m := n / 2
		median = int(values[m]+values[m+1]) / 2.0
	}

	mode := 0
	maxcount := 0

	for k, v := range hist {
		if v > maxcount {
			mode = k
			maxcount = v
		}
	}

	printHistogram(hist, maxcount)

	fmt.Println("items: ", n)
	fmt.Println("min: ", min)
	fmt.Println("mean: ", mean)
	fmt.Println("median: ", median)
	fmt.Println("mode: ", mode)
	fmt.Println("count(mode): ", maxcount)
	fmt.Println("max: ", max)
	fmt.Println("stdev: ", stddev)
}

func dupRune(r rune, n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = r
	}
	return string(s)
}

func printHistogram(hist map[int]int, maxcount int) {

	const keycols = 8
	const starcols = 55
	const countcols = 8

	const headerFmt = "%8s  %-55s  %8s\n"
	const rowFmt = "%8d  %-55s  %8d\n"

	fmt.Printf(headerFmt, "Key", "", "Count")

	fmt.Println(dupRune('-', keycols+starcols+countcols+4))

	keys := make([]int, 0, len(hist))

	for k := range hist {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		stars := dupRune('*', int(float64(starcols*(float64(hist[k])/float64(maxcount)))))
		fmt.Printf(rowFmt, k, stars, hist[k])
	}
}
