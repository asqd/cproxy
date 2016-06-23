package cproxy

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestCacheWorkerLoop(t *testing.T) {
	store := getStore()
	s := mockPrerenderServer(&testProxy{Code: http.StatusOK, Body: "Prerender"})

	url := s.URL
	def := *prerenderUrl

	// Set new prerenderUrl and defer default
	prerenderUrl = &url
	defer func() {
		prerenderUrl = &def
	}()

	defer s.Close()

	dir := resetCacheDir(t)
	defer os.Remove(dir)

	cw := &CacheWorker{store}

	afterGetBool := true
	afterGetFirstFail = func() {
		afterGetBool = false
	}

	afterRecacheBool := true
	afterRecache = func() {
		afterRecacheBool = false
	}

	defer func() {
		afterGetFirstFail = func() {}
		afterRecache = func() {}
	}()

	cw.Loop()

	if afterGetBool {
		t.Errorf("Loop() with empty %v should be triggred", *table)
	}

	store.Append(&Cache{time.Now().Add(-(*timeout)), "https://google.com", store})

	cw.Loop()

	if afterRecacheBool {
		t.Errorf("Loop() with cache in %v what should be recached should be trigerred", *table)
	}

	afterRecacheBool = true

	cw.Loop()

	if !afterRecacheBool {
		t.Errorf("Loop() with cache in %v what shouldn't be recached shouldn't be trigerred", *table)
	}
}
