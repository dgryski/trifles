package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	k := flag.Int("n", 500, "k")
	goldenName := flag.String("golden", "", "golden file to read")
	gotName := flag.String("f", "", "file name to test")

	flag.Parse()

	golden := loadCounts(*goldenName, *k*2)
	got := loadCounts(*gotName, *k)

	var total int

	for k, v := range got {
		e := golden[k] - v
		total += e * e
		delete(got, k)
		delete(golden, k)
	}

	// add the missing
	if false {
		for _, v := range golden {
			total += v * v
		}
	}

	log.Printf("float64(total)/float64(*k) = %+v\n", float64(total)/float64(*k))
}

func loadCounts(fname string, k int) map[string]int {

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	counts := make(map[string]int)

	sc := bufio.NewScanner(f)

	var lines int

	for sc.Scan() {
		line := sc.Text()
		tab := strings.IndexByte(line, ' ')
		count, _ := strconv.Atoi(line[tab+1:])
		item := line[:tab]
		counts[item] = count
		lines++
		if lines > k {
			break
		}
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	return counts
}
