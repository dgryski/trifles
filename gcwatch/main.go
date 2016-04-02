package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

func getMemStats(url string) (*runtime.MemStats, error) {

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode >= 400 {
		return nil, errors.New(r.Status)
	}

	var expvars struct {
		Memstats runtime.MemStats `json:"memstats"`
	}

	json.NewDecoder(r.Body).Decode(&expvars)

	return &expvars.Memstats, nil

}

func main() {

	delay := flag.Duration("d", 10*time.Second, "time between fetches")
	units := flag.Duration("u", time.Millisecond, "gc time units")

	flag.Parse()

	host := "http://localhost:8080/debug/vars"

	if flag.NArg() >= 1 {
		host = flag.Arg(0)
	}

	log.Printf("fetching %s every %s", host, *delay)

	m, err := getMemStats(host)
	if err != nil {
		log.Fatal(err)
	}

	ngcs := m.NumGC

	for range time.Tick(*delay) {
		m, err := getMemStats(host)
		if err != nil {
			log.Println(err)
			continue
		}

		var times []uint64
		for i := ngcs + 1; i <= m.NumGC; i++ {
			times = append(times, m.PauseNs[(i+255)%256]/uint64(*units))
		}

		fmt.Println(time.Now().Format(time.Stamp), ": ", m.NumGC-ngcs, times)
		ngcs = m.NumGC
	}

}
