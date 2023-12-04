package daos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
)

// default retries intervals (in ms)
var defaultRetryIntervals = []int{100, 250, 350, 500, 700, 1000}

func execLockRetry(timeout time.Duration, maxRetries int) dbx.ExecHookFunc {
	return func(q *dbx.Query, op func() error) error {
		if q.Context() == nil {
			cancelCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer func() {
				cancel()
				//nolint:staticcheck
				q.WithContext(nil) // reset
			}()
			q.WithContext(cancelCtx)
		}

		execErr := baseLockRetry(func(attempt int) error {
			return op()
		}, maxRetries)
		if execErr != nil && !errors.Is(execErr, sql.ErrNoRows) {
			execErr = fmt.Errorf("%w; failed query: %s", execErr, q.SQL())
		}

		return execErr
	}
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
