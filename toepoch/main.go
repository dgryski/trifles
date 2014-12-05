package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	field := flag.Int("f", 0, "field to transform")
	format := flag.String("d", "", "date format")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	var d time.Time
	var prev string

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		t := fields[*field]
		if t != prev {
			var err error
			d, err = time.Parse(*format, t)
			if err != nil {
				log.Println("time.Parse():", err)
				continue
			}
		}

		fields[*field] = fmt.Sprint(d.Unix())
		fmt.Printf("%s\n", strings.Join(fields, " "))
		prev = t
	}
}
