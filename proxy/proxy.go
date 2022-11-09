package main

import (
	"log"
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
		max := 86400 // max
		min := 6     // min 6 seconds
		delay := rand.Intn(max-min+1) + min

		timein := time.Now().UTC().Add(time.Duration(delay))
		w.Header().Add("Expires", timein.Format(http.TimeFormat))
	}

	// Randomly add Cache-Control header
	if rand.Intn(2) == 0 {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%v", 86400))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseString))

}
