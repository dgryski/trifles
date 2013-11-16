package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

func main() {
	port := flag.Int("p", 8080, "port to listen on")
	setNoCache := flag.Bool("no-cache", false, "set no-cache header on requests")
	shufflePhotos := flag.Bool("shuffle", false, "shuffle photo order")
	photosPerPage := flag.Int("n", 30, "photos per page")

	flag.Parse()

	dir1 := flag.Arg(0)
	dir2 := flag.Arg(1)

	log.Println("Comparing images from", dir1, "and", dir2, "on port", *port)

	d1files, _ := filepath.Glob(dir1 + "/*")
	d2files, _ := filepath.Glob(dir2 + "/*")

	if len(d1files) != len(d2files) {
		log.Println("file count mismatch: dir1=", len(d1files), "dir2=", len(d2files))
	}

	m := make(map[string]struct{})
	addFileNames(m, d1files)
	addFileNames(m, d2files)

	var photos []string
	for k := range m {
		photos = append(photos, k)
	}

	sort.Strings(photos)

	if *shufflePhotos {
		rand.Seed(time.Now().UnixNano())
		perm := rand.Perm(len(photos))
		dst := make([]string, len(photos))
		for i, v := range perm {
			dst[i] = photos[v]
		}

		photos = dst
	}

	// set up handler to serve our files
	http.Handle("/dir1/", http.StripPrefix("/dir1/", http.FileServer(http.Dir(dir1))))
	http.Handle("/dir2/", http.StripPrefix("/dir2/", http.FileServer(http.Dir(dir2))))

	http.HandleFunc("/compare", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if *setNoCache {
			w.Header().Set("Cache-Control", "no-cache")
		}

		pageStr := r.FormValue("page")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Println("invalid param for page, zeroing")
			page = 0
		}

		offset := page * *photosPerPage
		limit := offset + *photosPerPage

		if offset > len(photos) {
			offset = len(photos)
		}

		if limit > len(photos) {
			limit = len(photos)
		}

		prev := page - 1
		if prev < 0 {
			prev = 0
		}

		tmplParams := struct {
			Photos []string
			Next   int
			Prev   int
		}{
			Photos: photos[offset:limit],
			Next:   page + 1,
			Prev:   prev,
		}

		compareTMPL.Execute(w, tmplParams)
	})

	portStr := fmt.Sprintf(":%d", *port)
	log.Fatal("ListenAndServe:", http.ListenAndServe(portStr, nil))
}

func addFileNames(m map[string]struct{}, paths []string) {
	for _, fname := range paths {
		_, file := filepath.Split(fname)
		m[file] = struct{}{}
	}
}

var compareTMPL = template.Must(template.New("compare").Parse(
	`<html><head></head>
<body>
<table>
{{ range .Photos }}
<tr>
<td><img src="/dir1/{{ . }}"></td>
<td><img src="/dir2/{{ . }}"></td>
</tr>
{{ end }}
</table>

<a href="/compare?page={{ .Next }}">Next</a>
<a href="/compare?page={{ .Prev }}">Prev</a>
</body>`))
