package main

//go:generate ragel -Z lexer.rl

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func main() {

	f := flag.String("f", "/dev/stdin", "input file")
	w := flag.String("which", "ragel", "which matcher (one of: bsearch, map, radix, ragel, aho. defaults to ragel.)")
	n := flag.Int("n", 2000000, "size")
	iter := flag.Int("i", 10, "iterations")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to this file")

	flag.Parse()

	inf, err := os.Open(*f)
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	var arr [][]byte

	scanner := bufio.NewScanner(inf)
	for scanner.Scan() {
		b := append([]byte(nil), scanner.Bytes()...)
		arr = append(arr, b)
	}

	log.Printf("using matcher=%+v\n", *w)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		defer pprof.WriteHeapProfile(f)
	}

	switch *w {

	case "bsearch":
		initBsearch()

		t0 := time.Now()
		var found int
		for i := 0; i < *iter; i++ {
			for _, a := range arr[:*n] {
				if MatchBsearch(a) {
					found++
				}

			}
		}
		log.Printf("time.Since(t0)=%+v\n", time.Since(t0))
		log.Printf("found=%+v\n", found)

	case "map":
		initMap()

		t0 := time.Now()
		var found int
		for i := 0; i < *iter; i++ {
			for _, a := range arr[:*n] {
				if MatchMap(a) {
					found++
				}

			}
		}
		log.Printf("time.Since(t0)=%+v\n", time.Since(t0))
		log.Printf("found=%+v\n", found)

	case "radix":
		initRadix()

		t0 := time.Now()
		var found int
		for i := 0; i < *iter; i++ {
			for _, a := range arr[:*n] {
				if MatchRadix(string(a)) {
					found++
				}

			}
		}
		log.Printf("time.Since(t0)=%+v\n", time.Since(t0))
		log.Printf("found=%+v\n", found)

	case "ragel":
		t0 := time.Now()
		var found int
		for i := 0; i < *iter; i++ {
			for _, a := range arr[:*n] {
				if MatchRagel(a) {
					found++
				}
			}
		}
		log.Printf("time.Since(t0)=%+v\n", time.Since(t0))
		log.Printf("found=%+v\n", found)

	case "aho":
		initAho()

		t0 := time.Now()
		var found int
		for i := 0; i < *iter; i++ {
			for _, a := range arr[:*n] {
				if MatchAho(a) {
					found++
				}
			}
		}
		log.Printf("time.Since(t0)=%+v\n", time.Since(t0))
		log.Printf("found=%+v\n", found)
	default:
		log.Fatalf("unknown matcher %s", *w)
	}
}
