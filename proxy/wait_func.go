package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type forwardWaitFunc func(ctx context.Context) error

func wait(ctx context.Context, dur time.Duration) error {
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

func newScalerForwardWaitFunc(scalerURL *url.URL, errTolerance uint, errWait time.Duration) forwardWaitFunc {
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
				if err := wait(ctx, errWait); err != nil {
					return err
				}
				continue
			}
			defer res.Body.Close()
			resBytes, err := io.ReadAll(res.Body)
			if toleranceMet(err) {
				return err
			} else if err != nil {
				if err := wait(ctx, errWait); err != nil {
					return err
				}
				continue
			}
			numScaled, err := strconv.Atoi(string(resBytes))
			if toleranceMet(err) {
				return err
			} else if err != nil {
				if err := wait(ctx, errWait); err != nil {
					return err
				}
				continue
			}
			// if the scaler reports it has scaled up to 1 or more replicas, we're done!
			if numScaled > 0 {
				return nil
			} else {
				log.Printf("Scaler reported 0 replicas, waiting")
				if err := wait(ctx, errWait); err != nil {
					return err
				}
			}
			// return nil
		}
	}
}
