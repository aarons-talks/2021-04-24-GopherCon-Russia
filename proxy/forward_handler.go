package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"gcruaaron.dev/pkg/proxy"
	"golang.org/x/sync/errgroup"
)

func newForwardingHandler(
	fwdSvcURL *url.URL,
	dialCtxFunc proxy.DialContextFunc,
	waitFunc proxy.ForwardWaitFunc,
	waitTimeout time.Duration,
	respTimeout time.Duration,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, done := context.WithTimeout(r.Context(), waitTimeout)
		defer done()
		grp, _ := errgroup.WithContext(ctx)
		grp.Go(func() error {
			return waitFunc(ctx)
		})
		waitErr := grp.Wait()
		if waitErr != nil {
			log.Printf("Error, not forwarding request")
			w.WriteHeader(502)
			w.Write([]byte(fmt.Sprintf("error on backend (%s)", waitErr)))
			return
		}
		log.Printf("forwarding request to %#v", *fwdSvcURL)
		forwardRequest(w, r, dialCtxFunc, respTimeout, fwdSvcURL)
	})
}

func forwardRequest(
	w http.ResponseWriter,
	r *http.Request,
	dialCtxFunc proxy.DialContextFunc,
	respHeaderTimeout time.Duration,
	fwdSvcURL *url.URL,
) {
	roundTripper := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialCtxFunc,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: respHeaderTimeout,
	}
	proxy := httputil.NewSingleHostReverseProxy(fwdSvcURL)
	proxy.Transport = roundTripper
	proxy.Director = func(req *http.Request) {
		log.Printf("forwarding request %#v", *req)
		req.URL = fwdSvcURL
		req.Host = fwdSvcURL.Host
		req.URL.Path = r.URL.Path
		req.URL.RawQuery = r.URL.RawQuery
		// delete the incoming X-Forwarded-For header so the proxy
		// puts its own in. This is also important to prevent IP spoofing
		req.Header.Del("X-Forwarded-For ")
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(502)
		errMsg := fmt.Errorf("error on backend (%w)", err).Error()
		w.Write([]byte(errMsg))
	}

	proxy.ServeHTTP(w, r)
}
