// most of the functions here are referred from; https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
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
	origin, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(origin)
	path := "/*catchall"

	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host

		wildcardIndex := strings.IndexAny(path, "*")
		proxyPath := singleJoiningSlash(origin.Path, req.URL.Path[wildcardIndex:])
		if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
			proxyPath = proxyPath[:len(proxyPath)-1]
		}
		req.URL.Path = proxyPath
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

// copied from https://golang.org/src/net/http/httputil/reverseproxy.go?s=3330:3391#L98
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
