package main

import (
	"log"
	"net/http"
)

const scaledTo = "1"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("scaler reporting %s replicas", scaledTo)
		w.Write([]byte(scaledTo))
	})
	log.Printf("Serving the Scaler on port 8082 with %s replicas", scaledTo)
	log.Fatal(http.ListenAndServe("0.0.0.0:8082", mux))
}
