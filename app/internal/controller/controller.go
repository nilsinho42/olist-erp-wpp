package controller

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

var OlistERPURL = os.Getenv("API_OLIST_URL")

type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, threshold int) Circuit {
	var failures int
	var last = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock()
		d := failures - threshold
		fmt.Printf("Failures: %d\n", failures)
		if d >= 0 {
			backoff := calculateBackOff(d)
			shouldRetryAt := last.Add(backoff)

			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", fmt.Errorf("circuit breaker open, retry after %s", shouldRetryAt)
			}
		}
		m.RUnlock()

		response, err := circuit(ctx)

		m.Lock()
		defer m.Unlock()
		last = time.Now()
		if err != nil {
			failures++
			return response, err
		}
		failures = 0
		return response, nil
	}
}

func DebounceContext(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time
	var m sync.Mutex
	var lastCtx context.Context
	var lastCancel context.CancelFunc

	return func(ctx context.Context) (string, error) {
		m.Lock()

		if time.Now().Before(threshold) {
			lastCancel()
		}

		lastCtx, lastCancel = context.WithCancel(ctx)
		threshold = time.Now().Add(d)

		m.Unlock()

		result, err := circuit(lastCtx)

		return result, err

	}
}

func calculateBackOff(d int) time.Duration {
	jitter := rand.Int63n(int64(d * 3))
	if backoff := (2 << d) * (time.Second + time.Duration(jitter)); backoff < 60*time.Second {
		return backoff
	}

	return 60*time.Second + time.Duration(jitter)
}
