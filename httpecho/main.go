package main

// TODO: reimplement useful features from http://httpbin.org/

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(r)

	//turn r.Body from a io.ReadCloser into the actual content
	body, _ := ioutil.ReadAll(r.Body)
	var jr map[string]interface{}
	json.Unmarshal(j, &jr)
	jr["Body"] = string(body)

	j, _ = json.Marshal(jr)
	w.Write(j)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
