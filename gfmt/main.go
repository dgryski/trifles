package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dgryski/go-linebreak"
)

func main() {

	var wrap func(s string, w, g int) string

	width := flag.Int("w", 75, "line width")
	goal := flag.Int("g", 72, "line goal")
	fast := flag.Bool("f", false, "fast")

	flag.Parse()

	wrap = linebreak.Wrap
	if *fast {
		wrap = linebreak.Greedy
	}

	scanner := bufio.NewScanner(os.Stdin)

	var line []string
	for scanner.Scan() {
		t := scanner.Text()
		if t == "" {
			text := strings.Join(line, " ")
			fmt.Println(wrap(text, *goal, *width))
			fmt.Println()
			line = line[:0]
		}
		line = append(line, t)
	}

	text := strings.Join(line, " ")
	fmt.Println(wrap(text, *goal, *width))
}
