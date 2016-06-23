package cproxy

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
)

var (
	googleHex = "687474703a2f2f676f6f676c652e636f6d"
	dirStub   = "/Users/roma/Repositories/prerender/"
)

func requestStub(ur string) (*http.Request, error) {
	url, err := url.Parse(ur)

	if err != nil {
		return &http.Request{}, err
	}

	r := &http.Request{URL: url, Host: "google.com"}

	return r, nil
}

func TestReqHelperFullUrl(t *testing.T) {
	dir := resetCacheDir(t)
	defer os.Remove(dir)

	expect := "http://google.com"

	r, err := requestStub(expect)

	if err != nil {
		t.Errorf("%v", err)
	}

	rh := ReqHelper{Request: r}

	got := rh.FullPath()

	if expect != got {
		t.Errorf("Expected %v got %v", expect, got)
	}
}

func TestReqHelperHextFromPath(t *testing.T) {
	dir := resetCacheDir(t)
	defer os.Remove(dir)

	expect := googleHex

	r, err := requestStub("http://google.com")

	if err != nil {
		t.Errorf("%v", err)
	}

	rh := ReqHelper{Request: r}

	got := rh.HexFromPath()

	if expect != got {
		t.Errorf("Expected %v got %v", expect, got)
	}
}

func TestReqHelperGotCache(t *testing.T) {
	dir := resetCacheDir(t)
	defer os.Remove(dir)

	r, err := requestStub("http://google.com")

	if err != nil {
		t.Errorf("%v", err)
	}

	rh := ReqHelper{Request: r}

	if rh.GotCache() == true {
		t.Errorf("When no file expected false got %v", rh.GotCache())
	}

	path := *cacheDir + googleHex
	b := BaseReader{}
	b.WriteBody(path, []byte{})

	defer func() {
		err := os.Remove(path)

		if err != nil {
			log.Println(err)
		}
	}()

	if rh.GotCache() == false {
		t.Errorf("When no file expected true got %v", rh.GotCache())
	}
}
