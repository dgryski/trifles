package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("p", 8080, "port to listen on")

	flag.Parse()

	log.Println("Serving files in the current directory on port", *port)

        fserver := http.FileServer(http.Dir("."))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            log.Println(r.URL)
            fserver.ServeHTTP(w, r)
        })

	portStr := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
