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
	"strings"
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

type WalkResults struct {
	Packages []struct {
		ImportPath string `json:"import_path"`
		Synopsis   string `json:"synopsis"`
	} `json:"packages"`
}

func getWalk(q string) map[string]Hit {
	u := fmt.Sprintf("https://gowalker.org/api/v1/search?key=%s", url.QueryEscape(q))
	var results WalkResults
	get(u, &results)

	hits := make(map[string]Hit)
	for i, pkg := range results.Packages {
		//	log.Println("walk", i, pkg.Package)
		hits[pkg.ImportPath] = Hit{i, pkg.ImportPath, pkg.Synopsis}
	}
	return hits
}

type GithubResults struct {
	Repositories []struct {
		HtmlUrl     string `json:"html_url"`
		Description string `json:"description"`
	} `json:"items"`
}

func getGithub(q string) map[string]Hit {
	u := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+language:go", url.QueryEscape(q))
	var results GithubResults
	get(u, &results)

	hits := make(map[string]Hit)
	for i, pkg := range results.Repositories {
		path := strings.TrimPrefix(pkg.HtmlUrl, "https://")
		hits[path] = Hit{i, path, pkg.Description}
	}
	return hits
}

func main() {

	flag.Parse()

	ch := make(chan map[string]Hit)

	query := flag.Arg(0)

	searchers := []func(string) map[string]Hit{
		getGddo,
		getGcse,
		getWalk,
		getGithub,
	}

	for _, s := range searchers {
		go func(s func(string) map[string]Hit) { ch <- s(query) }(s)
	}

	var rmaps []map[string]Hit

	for i := 0; i < len(searchers); i++ {
		rmaps = append(rmaps, <-ch)
	}

	r := make(map[string]Hit)

	for _, m := range rmaps {
		for p, h := range m {
			if rh, ok := r[p]; ok {
				if rh.Synopsis == "" {
					rh.Synopsis = h.Synopsis
					r[p] = rh
				}
			} else {
				r[p] = h
			}
		}
	}

	// wipe the votes
	for p, h := range r {
		h.Rank = 0
		r[p] = h
	}

	for _, m := range rmaps {
		merge(r, m)
	}

	var rs ResultSet
	for _, v := range r {
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
		fmt.Fprintf(w, "%s \t%s\n", pkg.Path, pkg.Synopsis)
	}
	w.Flush()
}

func merge(r, r1 map[string]Hit) {

	// Add scores for all candidates in r1
	for p, h := range r {
		if r, ok := r1[h.Path]; ok {
			h.Rank += r.Rank
		} else {
			h.Rank += len(r1)
		}
		r[p] = h
	}
}
