package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const scalerPort = 8082
const originPort = 8081

func main() {
	scalerURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", scalerPort))
	if err != nil {
		log.Fatalf("Invalid scaler URL: %s", err)
	}
	originURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", originPort))
	if err != nil {
		log.Fatalf("Invalid forwarding URL: %s", err)
	}
	port := 8080
	proxyMux := http.NewServeMux()
	coreDialer := &net.Dialer{
		Timeout:   500 * time.Millisecond,
		KeepAlive: 1 * time.Second,
	}

	dialContextFunc := newDialContextFuncWithRetry(coreDialer, 100, 1*time.Second)
	waitFunc := newScalerForwardWaitFunc(scalerURL, 100, 1*time.Second)
	proxyHdl := newForwardingHandler(
		originURL,
		dialContextFunc,
		waitFunc,
		10*time.Second,
		10*time.Second,
	)
	proxyMux.Handle("/", proxyHdl)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("proxy server starting on %s", addr)
	log.Printf("Using scalerURL = %s, originURL = %s", scalerURL, originURL)
	http.ListenAndServe(addr, proxyMux)
}
