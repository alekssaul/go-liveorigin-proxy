package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	origin := os.Getenv("ORIGIN")

	// initialize a reverse proxy and pass the actual backend server url here
	proxy, err := NewProxy(origin)
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
	}

	proxy.ModifyResponse = modifyResponse()

	return proxy, nil
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		// Randomly add Expires Header or not
		if rand.Intn(2) == 0 {
			max := 86400 // max
			min := 6     // min 6 seconds
			delay := rand.Intn(max-min+1) + min

			timein := time.Now().UTC().Add(time.Duration(delay))
			resp.Header.Set("Expires", timein.Format(http.TimeFormat))
		}

		// Randomly add Cache-Control header
		if rand.Intn(2) == 0 {
			resp.Header.Set("Cache-Control", fmt.Sprintf("max-age=%v", 86400))
		}

		return nil
	}
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
