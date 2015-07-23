package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/dgryski/go-linlog"
	"github.com/dgryski/go-onlinestats"
)

func main() {

	width := flag.Int("w", 100, "display width")
	llog := flag.Bool("l", false, "linear-log bucketing")
	percentiles := flag.String("p", "0.5,0.75,0.95,0.99,0.999", "percentile values to print")
	flag.Parse()

	hist := make(map[int]int)

	stats := onlinestats.NewRunning()

	min := int(math.MaxInt32)
	max := 0

	maxcount := 0

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		v, err := strconv.Atoi(sc.Text())
		if err != nil {
			log.Fatal(err)
		}

		if *llog {
			r, _ := linlog.BinOf(uint64(v), 4, 2)
			v = int(r)
		}

		if v < min {
			min = v
		}

		if v > max {
			max = v
		}
		c := hist[v] + 1
		if c > maxcount {
			maxcount = c
		}

		hist[v] = c

		stats.Push(float64(v))

	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	printHistogram(hist, *width, maxcount)

	fmt.Println("items: ", stats.Len())
	fmt.Println("min: ", min)

	keys := make([]int, 0, len(hist))

	for k := range hist {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var quantiles []float64

	if *percentiles == "all" {
		for i := 0; i < 100; i++ {
			quantiles = append(quantiles, float64(i)/100.0)
		}
	} else {
		pfloats := strings.Split(*percentiles, ",")
		for _, p := range pfloats {
			f, err := strconv.ParseFloat(p, 64)
			if err != nil {
				log.Printf("unable to parse float %s: %s\n", p, err)
				continue
			}
			quantiles = append(quantiles, f)
		}
	}

	pidx := 0
	total := 0
KEYS:
	for _, k := range keys {
		total += hist[k]
		p := float64(total) / float64(stats.Len())
		for quantiles[pidx] <= p {
			fmt.Printf("% 4.1f: %d\n", quantiles[pidx]*100, k)
			pidx++
			if pidx == len(quantiles) {
				break KEYS
			}
		}
	}

	fmt.Println("max: ", max)
	fmt.Println("stdev: ", stats.Stddev())
}

func dupRune(r rune, n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = r
	}
	return string(s)
}

func printHistogram(hist map[int]int, width, maxcount int) {

	const keycols = 8
	var starcols = width
	const countcols = 8

	var headerFmt = "%8s  %-" + strconv.Itoa(width) + "s %8s\n"
	var rowFmt = "%8d  %-" + strconv.Itoa(width) + "s %8d\n"

	fmt.Printf(headerFmt, "Key", "", "Count")

	fmt.Println(dupRune('-', keycols+starcols+countcols+4))

	keys := make([]int, 0, len(hist))

	for k := range hist {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		stars := dupRune('*', int(float64(float64(starcols)*(float64(hist[k])/float64(maxcount)))))
		fmt.Printf(rowFmt, k, stars, hist[k])
	}
}
