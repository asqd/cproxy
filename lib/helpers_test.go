package cproxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type testProxy struct {
	Code int
	Body string
}

func (p testProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(p.Code)
	fmt.Fprintf(w, "%v", p.Body)
}

func mockPrerenderServer(p *testProxy) *httptest.Server {
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

func getStore() *RedisStore {
	conn := ConnectRedis()

	store := &RedisStore{conn, &sync.Mutex{}}

	return store
}
