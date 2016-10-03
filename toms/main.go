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
	dur := flag.Duration("d", 1*time.Millisecond, "output duration")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	if *field == 0 {
		for scanner.Scan() {

			d, err := time.ParseDuration(scanner.Text())
			if err != nil {
				log.Println("time.Parse():", err)
				continue
			}

			fmt.Println(int(float64(d.Nanoseconds())/float64(*dur) + 0.5))
		}

		return
	}

	for scanner.Scan() {

		var d time.Duration

		var duration string

		fields := strings.Fields(scanner.Text())
		duration = fields[*field]
		d, err := time.ParseDuration(duration)
		if err != nil {
			log.Println("time.Parse():", err)
			continue
		}

		fields[*field] = fmt.Sprint(int(float64(d.Nanoseconds())/float64(*dur) + 0.5))
		fmt.Printf("%s\n", strings.Join(fields, " "))
	}
}
