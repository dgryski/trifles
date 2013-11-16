package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	dir := flag.String("d", ".", "directory to serve")
	port := flag.Int("p", 8080, "port to listen on")
	setNoCache := flag.Bool("no-cache", false, "set no-cache header on requests")

	flag.Parse()

	log.Println("Serving files from", *dir, "on port", *port)

	fserver := http.FileServer(http.Dir(*dir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.URL)
		if *setNoCache {
			w.Header().Set("Cache-Control", "no-cache")
		}
		fserver.ServeHTTP(w, r)
	})

	portStr := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
