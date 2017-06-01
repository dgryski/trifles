package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	bucketSize := flag.Int("n", 500, "bucket size")
	f := flag.String("f", "", "file to read")
	flag.Parse()

	var r io.Reader

	if *f == "" {
		r = os.Stdin
	} else {
		var err error
		r, err = os.Open(*f)
		if err != nil {
			log.Fatal(err)
		}
	}

	sc := bufio.NewScanner(r)

	counts := make(map[string]int)

	var bucket int

	for sc.Scan() {
		line := sc.Text()

		counts[line]++

		if bucket++; bucket == *bucketSize {
			for k, v := range counts {
				fmt.Printf("%s\t%d\n", k, v)
			}
			bucket = 0
			counts = map[string]int{}
		}
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	for k, v := range counts {
		fmt.Printf("%s\t%d\n", k, v)
	}

}
