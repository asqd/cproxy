package cproxy

import "net/http"

// Helper struct for working with response
// and request
type ReqHelper struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

// Get full url form request
func (r ReqHelper) FullPath() string {
	u := r.Request.URL
	u.Scheme = "http"
	u.Host = r.Request.Host

	return u.String()
}

// Writes header of response
func (r ReqHelper) WriteHeader(status int) {
	if r.Writer != nil {
		r.Writer.WriteHeader(status)
	}
}
