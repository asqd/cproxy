package cproxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func TestServeHTTP(t *testing.T) {
	pr := &testProxy{Code: 500, Body: ""}
	ps := mockPrerenderServer(pr)

	def := *prerenderUrl
	prerenderUrl = &ps.URL

	defer func() {
		prerenderUrl = &def
	}()

	defer ps.Close()

	store := getStore()
	ser := Server{store}

	s := httptest.NewServer(ser)

	_, err := store.Conn.Do("DEL", *table)

	if err != nil {
		t.Error(err)
	}

	defer s.Close()

	dir := resetCacheDir(t)
	defer os.Remove(dir)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", s.URL, nil)

	ser.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected recorder code to be %v got %v", http.StatusInternalServerError, w.Code)
	}

	pr.Code = 200
	pr.Body = "Test"

	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", s.URL, nil)

	ser.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected recorder code to be %v got %v", http.StatusOK, w.Code)
	}

	count, err := redis.Int(store.Conn.Do("LLEN", *table))

	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Errorf("Expected request to create record in redis")
	}

	_, err = store.Conn.Do("DEL", *table)

	if err != nil {
		t.Error(err)
	}

	c := &Cache{time.Now(), s.URL, store}
	err = c.Write([]byte(pr.Body))

	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	r, err = http.NewRequest("GET", s.URL, nil)

	ser.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected recorder code to be %v got %v", http.StatusOK, w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)

	if err != nil {
		t.Error(err)
	}

	if string(body) != pr.Body {
		t.Errorf("Expected content to be %v got %v", pr.Body, string(body))
	}
}
