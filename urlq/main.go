package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

func main() {

	qparam := flag.String("q", "q", "query parameter to extract")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		s := scanner.Text()
		u, err := url.Parse(s)
		if err != nil {
			log.Printf("%v: err %+v\n", s, err)
			continue
		}
		query := u.Query()
		if query == nil {
			continue
		}
		for _, qq := range query[*qparam] {
			fmt.Println(qq)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error during scan: %+v\n", err)
	}
}
