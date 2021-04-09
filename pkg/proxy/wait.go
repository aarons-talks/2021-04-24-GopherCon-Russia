package proxy

import (
	"context"
	"log"
	"time"
)

// ForwardWaitFunc is a function that takes a context and waits for a condition to be true.
// It returns a nil error if the condition is met successfully. Otherwise, it returns a descriptive
// error if the context was cancelled or another failure condition happened (i.e. another timeout)
//
// It's highly recommended that you the incoming request context -- or one derived from the incoming
// request context -- to this function.
type ForwardWaitFunc func(ctx context.Context) error

// Wait waits for dur or until ctx is done, whichever comes first. Returns ctx.Err()
// if the context is done before dur, nil otherwise
func Wait(ctx context.Context, dur time.Duration) error {
	timer := time.NewTimer(dur)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		log.Printf("scaler context done: %s", ctx.Err())
		return ctx.Err()
	}
}
