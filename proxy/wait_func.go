package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gcruaaron.dev/pkg/proxy"
)

// newScalerForwardWaitFunc creates a ForwardWaitFunc that watches scaler metrics from scalerURL
// and returns after it reports either >0 replicas, the given error tolerance errTolerance is reached
// (when trying to make RPCs to the scaler), the errWait timeout is reached, or the given context
// is cancelled (or times out).
//
// In the first case (>0 replicas), the returned ForwardWaitFunc will return nil and in all other cases,
// it will return a non-nil error.
func newScalerForwardWaitFunc(
	scalerURL *url.URL,
	errTolerance uint,
	errWait time.Duration,
) proxy.ForwardWaitFunc {
	return func(ctx context.Context) error {
		numErrs := uint(0)
		toleranceMet := func(err error) bool {
			if err != nil && numErrs >= errTolerance {
				log.Printf("scaler wait func tolerange met, bailing")
				return true
			} else if err != nil {
				log.Printf("scaler error tolerance not met, waiting")
				numErrs++
			}
			return false
		}

		for {
			res, err := http.Get(scalerURL.String())
			if toleranceMet(err) {
				return err
			} else if err != nil {
				if err := proxy.Wait(ctx, errWait); err != nil {
					return err
				}
				continue
			}
			defer res.Body.Close()
			resBytes, err := io.ReadAll(res.Body)
			if toleranceMet(err) {
				return err
			} else if err != nil {
				if err := proxy.Wait(ctx, errWait); err != nil {
					return err
				}
				continue
			}
			numScaled, err := strconv.Atoi(string(resBytes))
			if toleranceMet(err) {
				return err
			} else if err != nil {
				if err := proxy.Wait(ctx, errWait); err != nil {
					return err
				}
				continue
			}
			// if the scaler reports it has scaled up to 1 or more replicas, we're done!
			if numScaled > 0 {
				return nil
			} else {
				log.Printf("Scaler reported 0 replicas, waiting")
				if err := proxy.Wait(ctx, errWait); err != nil {
					return err
				}
			}
		}
	}
}
