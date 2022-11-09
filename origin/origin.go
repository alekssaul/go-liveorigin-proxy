package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	// http handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HttpHandler(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HttpHandler(w http.ResponseWriter, r *http.Request) {
	reqHeadersBytes, err := json.Marshal(r.Header)
	if err != nil {
		log.Println("Could not Marshal Req Headers")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseString := "Request Headers: " + string(reqHeadersBytes)

	// Randomly add Expires Header or not
	if rand.Intn(2) == 0 {
		w.Header().Add("Expires", "foo")
	}

	// Randomly add Cache-Control header
	if rand.Intn(2) == 0 {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%v", 1440))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseString))

}
