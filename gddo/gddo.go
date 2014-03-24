// Package main is a command-line client search for godoc.org
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"
)

type Results struct {
	Results []struct {
		Path     string
		Synopsis string
	}
}

func main() {

	flag.Parse()

	u := fmt.Sprintf("http://api.godoc.org/search?q=%s", url.QueryEscape(flag.Arg(0)))

	r, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var results Results

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &results)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for _, pkg := range results.Results {
		fmt.Fprintf(w, "%s\t%s\n", pkg.Path, pkg.Synopsis)
	}
	w.Flush()
}
