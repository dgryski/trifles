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
	Rank     int
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

func getGddo(q string) map[string]Hit {
	u := fmt.Sprintf("http://api.godoc.org/search?q=%s", url.QueryEscape(q))
	var results GddoResults
	get(u, &results)

	hits := make(map[string]Hit)
	for i, pkg := range results.Results {
		//	log.Println("gddo", i, pkg.Path)
		hits[pkg.Path] = Hit{i, pkg.Path, pkg.Synopsis}
	}
	return hits
}

type GcseResults struct {
	Hits []struct {
		Package  string
		Synopsis string
	}
}

func getGcse(q string) map[string]Hit {
	u := fmt.Sprintf("http://go-search.org/api?action=search&q=%s", url.QueryEscape(q))
	var results GcseResults
	get(u, &results)

	hits := make(map[string]Hit)
	for i, pkg := range results.Hits {
		//	log.Println("gcse", i, pkg.Package)
		hits[pkg.Package] = Hit{i, pkg.Package, pkg.Synopsis}
	}
	return hits
}

func main() {

	flag.Parse()

	ch := make(chan map[string]Hit)

	query := flag.Arg(0)

	go func() { ch <- getGddo(query) }()
	go func() { ch <- getGcse(query) }()

	r1 := <-ch
	r2 := <-ch

	for p, h := range r1 {
		if r, ok := r2[h.Path]; ok {
			h.Rank += r.Rank
		} else {
			h.Rank += len(r2)
		}
		r1[p] = h
	}

	for p, h := range r2 {
		if _, ok := r1[h.Path]; !ok {
			h.Rank += len(r1)
			r1[p] = h
		}
	}

	var rs ResultSet
	for _, v := range r1 {
		rs = append(rs, v)
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
