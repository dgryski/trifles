package main

//go:generate ragel -Z -G2 parse.rl

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/dchest/siphash"
	"github.com/lytics/hll"
)

const (
	lastMinute = 60
	last5m     = 5 * 60
	last10m    = 10 * 60
	last15m    = 15 * 60
	last30m    = 30 * 60
	last60m    = 60 * 60
	last2h     = 2 * 60 * 60
	last4h     = 4 * 60 * 60
	last8h     = 8 * 60 * 60
	last12h    = 12 * 60 * 60
	last24h    = 24 * 60 * 60
	last2d     = 2 * 24 * 60 * 60
	last5d     = 5 * 24 * 60 * 60
	last7d     = 7 * 24 * 60 * 60
	last14d    = 14 * 24 * 60 * 60
	last1month = 30 * 24 * 60 * 60
	last6month = 180 * 24 * 60 * 60
	last1y     = 365 * 24 * 60 * 60
	allTime    = 5 * 365 * 24 * 60 * 60
)

var seen = []struct {
	epoch int64
	count uint64
}{
	{lastMinute, 0},
	{last5m, 0},
	{last10m, 0},
	{last15m, 0},
	{last30m, 0},
	{last60m, 0},
	{last2h, 0},
	{last4h, 0},
	{last8h, 0},
	{last12h, 0},
	{last24h, 0},
	{last2d, 0},
	{last5d, 0},
	{last7d, 0},
	{last14d, 0},
	{last1month, 0},
	{last6month, 0},
	{last1y, 0},
	{allTime, 0},
}

const logFormat = "2006/01/02 15:04:05"

func main() {

	cpuprofile := flag.String("cpuprofile", "", "pprof output file")

	tzoffset := flag.Int64("tz", 0, "time zone offset between log and query")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("%v", err)
		}
		defer pprof.StopCPUProfile()
	}

	scanner := bufio.NewScanner(os.Stdin)

	h := hll.NewHll(14, 20)

	var t0 int64
	var prevt []byte

	for scanner.Scan() {
		t, m, e1, _, err := parseLine(scanner.Bytes())
		if err != nil {
			continue
		}

		h.Add(siphash.Hash(0, 0, m))

		if !bytes.Equal(prevt, t) {
			prevt = append(prevt[:0], t...)
			tm, _ := time.Parse(logFormat, string(prevt))
			t0 = tm.Unix()
		}
		diff := t0 - e1 - *tzoffset

		var i int
		for i < len(seen) && seen[i].epoch < diff {
			i++
		}
		if i == len(seen) {
			i = len(seen) - 1
		}
		seen[i].count++
	}

	fmt.Printf("unique metrics %8d\n", h.Cardinality())

	var maxcount uint64
	var total uint64
	for i := range seen {
		if seen[i].count > maxcount {
			maxcount = seen[i].count
		}
		total += seen[i].count
	}

	const starcols = 70
	stars := dupRune('*', starcols)
	var cumulative uint64
	for i := range seen {
		d := time.Duration(seen[i].epoch) * time.Second
		s := stars[:int(float64(float64(starcols)*(float64(seen[i].count)/float64(maxcount))))]
		cumulative += seen[i].count
		fmt.Printf("%12s %10d (%0.2f, cum %0.2f) %s\n", d.String(), seen[i].count, float64(seen[i].count)/float64(total), float64(cumulative)/float64(total), s)
	}
}

func dupRune(r rune, n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = r
	}
	return string(s)
}
