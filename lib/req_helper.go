package cproxy

import (
	"encoding/hex"
	"flag"
	"net/http"
	"os"
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

	return u.String()
}

// Get hex form reqeust full url
func (r ReqHelper) HexFromPath() string {
	result := hex.EncodeToString([]byte(r.FullPath()))

	return result
}

// Writes header of response
func (r ReqHelper) WriteHeader(status int) {
	if r.Writer != nil {
		r.Writer.WriteHeader(status)
	}
}
