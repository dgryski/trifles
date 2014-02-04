package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

type response struct {
	Args    map[string]string `json:"args"`
	Data    string            `json:"data"`
	Files   map[string]string `json:"files"`
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

	resp.Form = make(map[string]string)
	for k, v := range r.PostForm {
		resp.Form[k] = v[0]
	}

	resp.Headers = make(map[string]string)
	for k, v := range r.Header {
		resp.Headers[k] = v[0]
	}

	resp.URL = r.URL.String()

	switch {

	case strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data"):

		m, err := r.MultipartReader()
		if err != nil {
			resp.Files["<error>"] = err.Error()
			break
		}

		resp.Files = make(map[string]string)
		// we can parse multipart form data
		for {
			p, err := m.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				resp.Files["<error>"] = err.Error()
				break
			}

			b, _ := ioutil.ReadAll(p)

			var s string
			if utf8.Valid(b) {
				s = string(b)
			} else {
				ctype := p.Header.Get("Content-Type")
				if ctype == "" {
					ctype = "application/octet-stream"
				}
				s = fmt.Sprintf("data:%s;base64,%s", ctype, base64.StdEncoding.EncodeToString(b))
			}

			if p.FileName() != "" {
				resp.Files[p.FormName()] = s
			} else {
				resp.Form[p.FormName()] = s
			}
		}

	case strings.HasPrefix(r.Header.Get("Content-Type"), "application/json"):
		body, _ := ioutil.ReadAll(r.Body)
		resp.Data = string(body)

		var j interface{}

		err := json.Unmarshal(body, &j)
		if err == nil {
			resp.Json = j
		} else {
			// invalid JSON passed in body
			resp.Json = "<invalid json>"
		}

	default:
		body, _ := ioutil.ReadAll(r.Body)
		resp.Data = string(body)
	}

	j, _ := json.Marshal(resp)
	w.Write(j)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
