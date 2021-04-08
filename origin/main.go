package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("origin got request for %s", r.URL.Path)
		w.Write([]byte("hello from the origin!"))
	})
	log.Printf("Serving the origin on port 8081")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", mux))
}
