package cproxy

import (
	"flag"
	"time"
	"net/http"
)

var port = flag.String("port", ":80", "Port where server will listen")

// Starts server on specified port by -port flag, e.g. -port=:80
func RunServer(store Store) {
	server := Server{store}
	http.Handle("/", server)
	http.ListenAndServe(*port, nil)
}

// ServerHTTP struct
type Server struct {
	Store Store
}

// Handler for all requests
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	helper := ReqHelper{w, r}

	cache := &Cache{time.Now(), helper.FullPath(), s.Store}

	if cache.Exists() {
		body, err := cache.Read()

		if err != nil {
			helper.WriteHeader(http.StatusInternalServerError)
		}

		helper.Write(body)
		return
	}

	p := Prerender{helper, cache}
	p.Process()
}
