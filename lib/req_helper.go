package cproxy

import (
	"net/http"
	"strings"
)

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

	return strings.Replace(u.String(), "?_escaped_fragment_=", "", -1)

}

// Writes header of response
func (r ReqHelper) WriteHeader(status int) {
	if r.Writer != nil {
		r.Writer.WriteHeader(status)
	}
}

// Writes body of response
func (r ReqHelper) Write(body []byte) {
	if r.Writer != nil {
		r.Writer.Write(body)
	}
}
