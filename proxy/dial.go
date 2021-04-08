package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func newDialContextFuncWithRetry(coreDialer *net.Dialer, numRetries uint, retryPause time.Duration) DialContextFunc {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		var lastError error
		for i := uint(0); i < numRetries; i++ {
			log.Printf("dialing try %d", i)
			conn, err := coreDialer.DialContext(ctx, network, addr)
			log.Printf("dialed, connection = %#v, error = %#v", conn, err)
			if err == nil {
				return conn, nil
			}
			lastError = err
			t := time.NewTimer(retryPause)
			select {
			case <-ctx.Done():
				t.Stop()
				return nil, fmt.Errorf("context timed out: %s", ctx.Err())
			case <-t.C:
				t.Stop()
			}
		}
		return nil, lastError
	}
}
