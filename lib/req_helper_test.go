package cproxy

import (
	"net/http"
	"net/url"
	"os"
	"testing"
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

	r, err := requestStub(expect + "?_escaped_fragment_=")

	if err != nil {
		t.Errorf("%v", err)
	}

	rh := ReqHelper{Request: r}

	got := rh.FullPath()

	if expect != got {
		t.Errorf("Expected %v got %v", expect, got)
	}
}
