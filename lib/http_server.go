package cproxy

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var port = flag.String("port", ":80", "Port where server will listen")

// Starts server on specified port by -port flag, e.g. -port=:80
func RunServer(store Store) {
	server := Server{store}
	http.Handle("/", server)
	http.ListenAndServe(*port, nil)
}

// http.Handler
type Server struct {
	Store Store
}

// Handler for all requests
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	helper := ReqHelper{w, r}

	beg := time.Now()

	cache := &Cache{time.Now(), helper.FullPath(), s.Store}

	if cache.Exists() {
		body, err := cache.Read()

		if err != nil {
			helper.WriteHeader(http.StatusInternalServerError)
		}

		helper.Write(body)

		spend := time.Since(beg)
		log.Printf("Request: %v - %.2f sec\n", helper.FullPath(), spend.Seconds())
		return
	}

	p := Prerender{helper, cache}
	p.Process()

	spend := time.Since(beg)
	log.Printf("Request: %v - %.2f sec\n", helper.FullPath(), spend.Seconds())
}
