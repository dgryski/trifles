// Package main is a command-line client search for godoc.org
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"text/tabwriter"
)

type Hit struct {
	Rank     float64
	Path     string
	Synopsis string
}

type ResultSet []Hit

func (r ResultSet) Len() int           { return len(r) }
func (r ResultSet) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ResultSet) Less(i, j int) bool { return r[i].Rank < r[j].Rank }

func get(url string, jsr interface{}) {
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&jsr)
}

type GddoResults struct {
	Results []struct {
		Path     string
		Synopsis string
	}
}

func getGddo(q string) []Hit {
	u := fmt.Sprintf("http://api.godoc.org/search?q=%s", url.QueryEscape(q))
	var results GddoResults
	get(u, &results)

	var hits []Hit
	for i, pkg := range results.Results {
		hits = append(hits, Hit{float64(i), pkg.Path, pkg.Synopsis})
	}
	return hits
}

type GcseResults struct {
	Hits []struct {
		Package  string
		Synopsis string
	}
}

func getGcse(q string) []Hit {
	u := fmt.Sprintf("http://go-search.org/api?action=search&q=%s", url.QueryEscape(q))
	var results GcseResults
	get(u, &results)

	var hits []Hit
	for i, pkg := range results.Hits {
		hits = append(hits, Hit{float64(i), pkg.Package, pkg.Synopsis})
	}
	return hits
}

func main() {

	flag.Parse()

	ch := make(chan []Hit)

	query := flag.Arg(0)

	go func() { ch <- getGddo(query) }()
	go func() { ch <- getGcse(query) }()

	hits := <-ch
	hits = append(hits, <-ch...)

	results := make(map[string]Hit)

	for _, h := range hits {
		if r, ok := results[h.Path]; ok {
			r.Rank = (-r.Rank + h.Rank) / 2
			results[h.Path] = r
		} else {
			results[h.Path] = Hit{-h.Rank, h.Path, h.Synopsis}
		}
	}

	var rs ResultSet
	for _, v := range results {
		if v.Rank >= 0 {
			rs = append(rs, v)
		}
	}

	sort.Sort(rs)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	end := 20
	if len(rs) < end {
		end = len(rs)
	}

	for _, pkg := range rs[:end] {
		fmt.Fprintf(w, "%s\t%s\n", pkg.Path, pkg.Synopsis)
	}
	w.Flush()
}
