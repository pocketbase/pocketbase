package daos

import (
	"context"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
)

const defaultQueryTimeout time.Duration = 2 * time.Minute

const defaultMaxRetries int = 10

var defaultRetryIntervals = []int{100, 250, 350, 500, 700, 1000, 1200, 1500}

func onLockErrorRetry(s *dbx.SelectQuery, op func() error) error {
	return baseLockRetry(func(attempt int) error {
		// load a default timeout context if not set explicitly
		if s.Context() == nil {
			ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
			defer func() {
				cancel()
				s.WithContext(nil) // reset
			}()
			s.WithContext(ctx)
		}

		return op()
	}, defaultMaxRetries)
}

func baseLockRetry(op func(attempt int) error, maxRetries int) error {
	attempt := 1

Retry:
	err := op(attempt)

	if err != nil &&
		attempt <= maxRetries &&
		// we are checking the err message to handle both the cgo and noncgo errors
		strings.Contains(err.Error(), "database is locked") {
		// wait and retry
		time.Sleep(getDefaultRetryInterval(attempt))
		attempt++
		goto Retry
	}

	return err
}

func getDefaultRetryInterval(attempt int) time.Duration {
	if attempt < 0 || attempt > len(defaultRetryIntervals)-1 {
		return time.Duration(defaultRetryIntervals[len(defaultRetryIntervals)-1]) * time.Millisecond
	}

	return time.Duration(defaultRetryIntervals[attempt]) * time.Millisecond
}
