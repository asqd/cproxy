package cproxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PreProxy struct {
	StatusCode int
}

func (p PreProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(p.StatusCode)
	fmt.Fprintln(w, "Write from fake prerender")
}

func mockPrerenderServer(p *PreProxy) *httptest.Server {
	return httptest.NewServer(p)
}

func resetCacheDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "")

	if err != nil {
		t.Errorf("%v", err)
	}

	cacheDir = &dir

	return dir
}

func getStore() RedisStore {
	conn := ConnectRedis()

	store := RedisStore{Conn: conn}

	return store
}
