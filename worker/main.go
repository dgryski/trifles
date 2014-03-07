package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
)

func worker(urlch chan string, donech chan struct{}) {
	for u := range urlch {
		log.Println("processing u=", u)
	}

	donech <- struct{}{}
}

func main() {

	infname := flag.String("f", "", "input file name")
	workers := flag.Int("w", 4, "workers")

	urlch := make(chan string)
	donech := make(chan struct{})

	log.Println("using", *workers, "workers")

	for i := 0; i < *workers; i++ {
		go worker(urlch, donech)
	}

	var f io.Reader
	if *infname == "" {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(*infname)

		if err != nil {
			log.Fatal("error loading input file", *infname, ":", err)
		}
	}

	scan := bufio.NewScanner(f)
	var items []string
	for scan.Scan() {
		items = append(items, scan.Text())
	}

	log.Println("Total items:", len(items))

	processedItems := 0

	// send off the work
	for _, item := range items {
		urlch <- item
		processedItems++
	}

	close(urlch)

	log.Println("waiting for workers to complete")
	for i := 0; i < *workers; i++ {
		<-donech
	}
}
