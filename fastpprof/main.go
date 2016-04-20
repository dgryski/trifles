package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func registerPProf(h func(string, func(http.ResponseWriter, *http.Request))) {
	h("/debug/pprof/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		pprof.Index(w, r)
	})
	h("/debug/pprof/cmdline", pprof.Cmdline)
	h("/debug/pprof/profile", pprof.Profile)
	h("/debug/pprof/symbol", pprof.Symbol)
	h("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	h("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	h("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	h("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}

func main() {

	port := flag.Int("p", 8080, "http port")
	flag.Parse()

	pprofMux := http.NewServeMux()
	registerPProf(pprofMux.HandleFunc)
	fastpprof := fasthttpadaptor.NewFastHTTPHandler(pprofMux)

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch {
		case bytes.Equal(ctx.Path(), []byte("/debug/pprof")):
			ctx.Redirect("/debug/pprof/", 301)
		case bytes.HasPrefix(ctx.Path(), []byte("/debug/pprof/")):
			fastpprof(ctx)
		default:
			ctx.NotFound()
		}
	}

	log.Println("listening on port", *port)

	err := fasthttp.ListenAndServe(":"+strconv.Itoa(*port), requestHandler)
	if err != nil {
		log.Println("error during fasthttp.ListenAndServe(): ", err)
	}
}
