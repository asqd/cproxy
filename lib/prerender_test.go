package cproxy

import (
	"net/http"
	"testing"
	"time"
)

func createPrerender() Prerender {
	r := ReqHelper{}
	c := &Cache{time.Now(), "http://google.com", nil}

	p := Prerender{r, c}

	return p
}

func TestPrerenderUrlWithPrerender(t *testing.T) {
	p := createPrerender()

	expect := "http://localhost:3000/http://google.com"

	got := p.UrlWithPrerender()

	if expect != got {
		t.Errorf("Expected %v got %v", expect, got)
	}
}

func TestPrerenderGetContent(t *testing.T) {
	proxy := &testProxy{Code: http.StatusOK, Body: "Prerender"}
	s := mockPrerenderServer(proxy)

	p := createPrerender()

	url := s.URL
	def := *prerenderUrl

	// Set new prerenderUrl and defer default
	prerenderUrl = &url
	defer func() {
		prerenderUrl = &def
	}()

	defer s.Close()

	body, code := p.GetContent()

	if string(body) != proxy.Body && code != proxy.Code {
		t.Errorf("Expected code %v got %, body %v got %v", proxy.Code, code, proxy.Body, string(body))
	}
}

func TestPrerenderProcess(t *testing.T) {
	proxy := &testProxy{Code: http.StatusOK, Body: "Prerender"}
	s := mockPrerenderServer(proxy)
	store := getStore()

	p := createPrerender()
	p.Cache.Store = store
	p.Cache.CreatedAt = time.Now().Add(-(*timeout))

	url := s.URL
	def := *prerenderUrl

	// Set new prerenderUrl and defer default
	prerenderUrl = &url
	defer func() {
		prerenderUrl = &def
	}()

	defer s.Close()

	p.Process()

	got, err := p.Cache.Read()

	if err != nil {
		t.Error(err)
	}

	if proxy.Body != string(got) {
		t.Errorf("After rewriting cache expected content to be %v got %v", proxy.Body, string(got))
	}

}
