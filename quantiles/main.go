package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"time"

	"github.com/VividCortex/gohistogram"
	bquantile "github.com/bmizerany/perks/quantile"
	"github.com/caio/go-tdigest"
	gogk "github.com/dgryski/go-gk"
	"github.com/dgryski/go-kll"
	squantile "github.com/streadway/quantile"
)

func main() {

	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

	f := flag.String("f", "", "file to read")

	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ch := make(chan int)

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

	go sendInts(r, ch)

	var stream []int

	bmq := bquantile.NewBiased()
	gk := gogk.New(0.0001)
	sq := squantile.New(squantile.Unknown(0.0001))
	vvh := gohistogram.NewHistogram(80)
	td := tdigest.New(100)
	kl := kll.New(1024)

	for v := range ch {
		stream = append(stream, v)
		bmq.Insert(float64(v))
		gk.Insert(float64(v))
		sq.Add(float64(v))
		vvh.Add(float64(v))
		td.Add(float64(v), 1)
		kl.Update(float64(v))
	}

	qq := []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.95, 0.99, 0.999, 0.9999, 0.99999}

	fmt.Println("bq:")
	for _, q := range qq {
		fmt.Printf(" % 4g", bmq.Query(q))
	}
	fmt.Println()

	fmt.Println("sq:")
	for _, q := range qq {
		fmt.Printf(" % 4g", sq.Get(q))
	}
	fmt.Println()

	fmt.Println("gk:")
	for _, q := range qq {
		fmt.Printf(" % 4g", gk.Query(q))
	}
	fmt.Println()

	fmt.Println("vv:")
	for _, q := range qq {
		fmt.Printf(" % 4d", int(vvh.Quantile(q)))
	}
	fmt.Println()

	fmt.Println("td: ")
	for _, q := range qq {
		fmt.Printf(" % 4d", int(td.Quantile(q)))
	}
	fmt.Println()

	klcdf := kl.CDF()
	fmt.Println("kl: ")
	for _, q := range qq {
		fmt.Printf(" % 4d", int(klcdf.Query(q)))
	}
	fmt.Println()

	sort.Ints(stream)
	fmt.Println("ex: ")
	for _, q := range qq {
		fmt.Printf(" % 4d", stream[int(float64(len(stream))*q)])
	}
	fmt.Println()
}

func sendInts(r io.Reader, ch chan<- int) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		b := sc.Bytes()
		v, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		ch <- int(v)
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
	close(ch)
}
