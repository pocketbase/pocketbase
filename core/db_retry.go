package core

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
var defaultRetryIntervals = []int{50, 100, 150, 200, 300, 400, 500, 700, 1000}

// default max retry attempts
const defaultMaxLockRetries = 12

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

	if err != nil && attempt <= maxRetries {
		errStr := err.Error()
		// we are checking the error against the plain error texts since the codes could vary between drivers
		if strings.Contains(errStr, "database is locked") ||
			strings.Contains(errStr, "table is locked") {
			// wait and retry
			time.Sleep(getDefaultRetryInterval(attempt))
			attempt++
			goto Retry
		}
	}

	return err
}

func getDefaultRetryInterval(attempt int) time.Duration {
	if attempt < 0 || attempt > len(defaultRetryIntervals)-1 {
		return time.Duration(defaultRetryIntervals[len(defaultRetryIntervals)-1]) * time.Millisecond
	}

	return time.Duration(defaultRetryIntervals[attempt]) * time.Millisecond
}
