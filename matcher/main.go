package main

//go:generate ragel -Z lexer.rl

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"
)

func main() {

	f := flag.String("f", "/dev/stdin", "input file")
	w := flag.String("which", "ragel", "which matcher")
	n := flag.Int("n", 2000000, "size")
	iter := flag.Int("i", 10, "iterations")

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

	switch *w {

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
