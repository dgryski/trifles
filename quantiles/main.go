package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/VividCortex/gohistogram"
	bquantile "github.com/bmizerany/perks/quantile"
	"github.com/dgryski/go-fastquantiles"
	squantile "github.com/streadway/quantile"
)

func main() {

	f := flag.String("f", "", "file to read")
	zwn := flag.Int("zwn", 0, "number of entries")

	rand.Seed(time.Now().UnixNano())

	flag.Parse()

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
	gk := fastquantiles.NewGK()
	sq := squantile.New(squantile.Unknown(0.1))
	vvh := gohistogram.NewHistogram(80)
	zeps := 1 / math.Exp(math.Log(float64(*zwn))-1)
	zw, _ := fastquantiles.New(zeps, *zwn)

	for v := range ch {
		stream = append(stream, v)
		bmq.Insert(float64(v))
		gk.Insert(float64(v))
		sq.Add(float64(v))
		vvh.Add(float64(v))
		zw.Update(float64(v))
	}

	fmt.Println("bq:")
	for i := 1; i <= 10; i++ {
		fmt.Printf(" % 4g", bmq.Query(float64(i)*0.1))
	}
	fmt.Println()

	fmt.Println("sq:")
	for i := 1; i <= 10; i++ {
		fmt.Printf(" % 4g", sq.Get(float64(i)*0.1))
	}
	fmt.Println()

	fmt.Println("gk:")
	for i := 1; i <= 10; i++ {
		fmt.Printf(" % 4g", gk.Query(float64(i)*0.1))
	}
	fmt.Println()

	fmt.Println("vv:")
	for i := 1; i <= 10; i++ {
		fmt.Printf(" % 4d", int(vvh.Quantile(float64(i)*0.1)))
	}
	fmt.Println()

	fmt.Println("zw: ")
	zw.Finish()
	for i := 1; i <= 10; i++ {
		fmt.Printf(" % 4g", zw.Query(float64(i)*0.1))
	}
	fmt.Println()

	sort.Ints(stream)
	fmt.Println("ex: ")
	for i := 1; i <= 10; i++ {
		fmt.Printf(" % 4d", stream[int(float64(len(stream))*float64(i)*0.1)-1])
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
