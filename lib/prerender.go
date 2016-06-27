package cproxy

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var prerenderUrl = flag.String("prerender-url", "http://localhost:3000", "Full path to prerender service")

// Type for getting content directly from prerender service
// and cache it
type Prerender struct {
	ReqHelper
	*Cache
}

// Process directives for recieving content from prerender
// service and write cache and write to http.ResponseWriter
func (p Prerender) Process() {
	body, code := p.GetContent()

	if code == http.StatusOK {
		err := p.Update(body)

		if err != nil {
			p.ReqHelper.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	p.ReqHelper.WriteHeader(code)
	p.ReqHelper.Write(body)
}

// Return url for prerender service request e.g.:
// http://localhost:3000/http://google.com
func (p Prerender) UrlWithPrerender() string {
	return *prerenderUrl + "/" + p.Cache.URL
}

// Requests to prerender service and return
// body and status code from it, if there was
// errors it will return empty body and code 500
func (p Prerender) GetContent() ([]byte, int) {
	response, err := http.Get(p.UrlWithPrerender())

	if err != nil {
		return p.GetContent()
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println(err)

		return p.GetContent()
	}

	return body, response.StatusCode
}
