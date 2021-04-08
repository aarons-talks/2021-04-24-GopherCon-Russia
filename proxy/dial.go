package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func DialContextWithRetry(coreDialer *net.Dialer, numRetries uint, retryPause time.Duration) DialContextFunc {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		var lastError error
		for i := uint(0); i < numRetries; i++ {
			conn, err := coreDialer.DialContext(ctx, network, addr)
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
