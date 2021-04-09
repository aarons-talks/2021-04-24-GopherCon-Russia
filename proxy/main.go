package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"gcruaaron.dev/pkg/proxy"
)

func main() {
	scalerURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", proxy.ScalerPort))
	if err != nil {
		log.Fatalf("Invalid scaler URL: %s", err)
	}
	originURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", proxy.OriginPort))
	if err != nil {
		log.Fatalf("Invalid forwarding URL: %s", err)
	}
	proxyMux := http.NewServeMux()
	coreDialer := &net.Dialer{
		Timeout:   500 * time.Millisecond,
		KeepAlive: 1 * time.Second,
	}

	dialContextFunc := newDialContextFuncWithRetry(coreDialer, 100, 1*time.Second)
	waitFunc := newScalerForwardWaitFunc(scalerURL, 100, 1*time.Second)
	proxyHandler := newForwardingHandler(
		originURL,
		dialContextFunc,
		waitFunc,
		10*time.Second,
		10*time.Second,
	)
	proxyMux.Handle("/", proxyHandler)

	addr := fmt.Sprintf("0.0.0.0:%d", proxy.ProxyPort)
	log.Printf("proxy server starting on %s", addr)
	log.Printf("Using scalerURL = %s, originURL = %s", scalerURL, originURL)
	http.ListenAndServe(addr, proxyMux)
}
