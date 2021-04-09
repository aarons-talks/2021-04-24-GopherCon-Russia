package main

import (
	"fmt"
	"log"
	"net/http"

	pkgnet "gcruaaron.dev/pkg/net"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("origin got request for %s", r.URL.Path)
		w.Write([]byte("hello from the origin!"))
	})
	addr := fmt.Sprintf("0.0.0.0:%d", pkgnet.OriginPort)
	log.Printf("Serving the origin on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
