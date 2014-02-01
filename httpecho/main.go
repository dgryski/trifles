package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type response struct {
	Args map[string]string `json:"args"`
	Data string            `json:"data"`
	// File    map[string]string `json:"file"`
	Form    map[string]string `json:"form"`
	Headers map[string]string `json:"headers"`
	URL     string            `json:"url"`
	Json    interface{}       `json:"json"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp response

	r.ParseForm()

	resp.Args = make(map[string]string)
	for k, v := range r.URL.Query() {
		resp.Args[k] = v[0]
	}

	body, _ := ioutil.ReadAll(r.Body)
	resp.Data = string(body)

	resp.Form = make(map[string]string)
	for k, v := range r.PostForm {
		resp.Form[k] = v[0]
	}

	resp.Headers = make(map[string]string)
	for k, v := range r.Header {
		resp.Headers[k] = v[0]
	}

	resp.URL = r.URL.String()

	if r.Header.Get("Content-Type") == "application/json" {
		var j interface{}

		err := json.Unmarshal(body, &j)
		if err == nil {
			resp.Json = j
		} else {
			// invalid JSON passed in body
			resp.Json = "<invalid json>"
		}
	}

	j, _ := json.Marshal(resp)
	w.Write(j)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
