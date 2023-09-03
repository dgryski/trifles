package main

import (
	_ "embed"
	"encoding/json"
	"flag"
)

//go:embed repos.json
var response []byte

func main() {
	n := flag.Int("n", 1000, "number of iterations")
	flag.Parse()

	var total int

	var repos []*Repository
	for i := 0; i < *n; i++ {
		if err := json.Unmarshal(response, &repos); err != nil {
			println("error unmarshaling", err)
			return
		}
		total += len(repos)
	}
	println(total)
}
