package main

//go:generate ragel -Z -G2 lexer.rl
//go:generate go run gen-mafsa.go

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
	w := flag.String("which", "ragel", "which matcher (one of: bsearch, map, radix, ragel, aho, mafsa. defaults to ragel.)")
	n := flag.Int("n", 2000000, "size")
	iter := flag.Int("i", 10, "iterations")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to this file")
	bloom := flag.Bool("bloom", false, "use bloom filter")

	flag.Parse()

	inf, err := os.Open(*f)
	if err != nil {
		log.Fatalf("err=%v", err)
	}

	var arr [][]byte
	var strarr []string
	var rarr [][]rune

	scanner := bufio.NewScanner(inf)
	for scanner.Scan() {
		b := append([]byte(nil), scanner.Bytes()...)
		strarr = append(strarr, scanner.Text())
		arr = append(arr, b)
		rarr = append(rarr, []rune(scanner.Text()))
		if len(arr) > *n {
			break
		}
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

	if *bloom {
		initBloom()
	}

	switch *w {

	case "bsearch":
		initBsearch()

		t0 := time.Now()
		var found int
		for i := 0; i < *iter; i++ {
			for j, a := range arr[:*n] {
				if MatchBloom(a) && MatchBsearch(strarr[j]) {
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
			for j, a := range arr[:*n] {
				if MatchBloom(a) && MatchMap(strarr[j]) {
					found++
				}

			}
		}
		log.Printf("time.Since(t0)=%+v\n", time.Since(t0))
		log.Printf("found=%+v\n", found)

	case "mafsa":
		initMAFSA()

		t0 := time.Now()
		var found int
		for i := 0; i < *iter; i++ {
			for j, a := range arr[:*n] {
				if MatchBloom(a) && MatchMAFSA(rarr[j]) {
					found++
				}

			}
		}
		log.Printf("time.Since(t0)=%+v\n", time.Since(t0))
		log.Printf("found=%+v\n", found)

	case "mph":
		initMPH()

		t0 := time.Now()
		var found int
		for i := 0; i < *iter; i++ {
			for j, a := range arr[:*n] {
				if MatchBloom(a) && MatchMPH(strarr[j]) {
					found++
				}

			}
		}
		log.Printf("time.Since(t0)=%+v\n", time.Since(t0))
		log.Printf("found=%+v\n", found)

	case "dmph":
		initDMPH()

		t0 := time.Now()
		var found int
		for i := 0; i < *iter; i++ {
			for j, a := range arr[:*n] {
				if MatchBloom(a) && MatchDMPH(strarr[j]) {
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
			for j, a := range arr[:*n] {
				if MatchBloom(a) && MatchRadix(strarr[j]) {
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
				if MatchBloom(a) && MatchRagel(a) {
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
				if MatchBloom(a) && MatchAho(a) {
					found++
				}
			}
		}
		log.Printf("time.Since(t0)=%+v\n", time.Since(t0))
		log.Printf("found=%+v\n", found)
	default:
		log.Fatalf("unknown matcher %s", *w)
	}

	if *bloom {
		log.Printf("filtered=%+v\n", filtered)
	}
}
