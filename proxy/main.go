package main

import "time"
import "net"
import "net/http"

func main() {
	port := 8080
	proxyMux := http.NewServeMux()
	dialer := 	&net.Dialer{
		Timeout:   500 * time.Millisecond,
		KeepAlive: 1 * time.Second,
	}

	dialContextFunc := kedanet.DialContextWithRetry(dialer, timeouts.DefaultBackoff())
	proxyHdl := newForwardingHandler(
		svcURL,
		dialContextFunc,
		waitFunc,
		timeouts.DeploymentReplicas,
		timeouts.ResponseHeader,
	)
	proxyMux.Handle("/*", countMiddleware(q, proxyHdl))

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("proxy server starting on %s", addr)
	nethttp.ListenAndServe(addr, proxyMux)
}

}
