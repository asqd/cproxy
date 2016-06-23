package cproxy

import "time"

// Caching worker struct
type CacheWorker struct {
	Store Store
}

// Run worker for endless recache loop
func (w CacheWorker) Run() {
	for {
		w.Loop()
	}
}

// Method is trying looping thrue cache records
// and recaching them if needed
func (w CacheWorker) Loop() {
	c, err := w.Store.GetFirst()

	if err != nil {
		time.Sleep(100 * time.Millisecond)
		afterGetFirstFail()
		return
	}

	if c.ShouldRecache() {
		p := Prerender{Cache: c}
		p.Process()
		afterRecache()
		return
	}

	w.Store.Append(c)
}

var (
	afterGetFirstFail = func() {}
	afterRecache      = func() {}
)
