package daos

import (
	"errors"
	"testing"
)

func TestGetDefaultRetryInterval(t *testing.T) {
	t.Parallel()

	if i := getDefaultRetryInterval(-1); i.Milliseconds() != 1000 {
		t.Fatalf("Expected 1000ms, got %v", i)
	}

	if i := getDefaultRetryInterval(999); i.Milliseconds() != 1000 {
		t.Fatalf("Expected 1000ms, got %v", i)
	}

	if i := getDefaultRetryInterval(3); i.Milliseconds() != 500 {
		t.Fatalf("Expected 500ms, got %v", i)
	}
}

func TestBaseLockRetry(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		err              error
		failUntilAttempt int
		expectedAttempts int
	}{
		{nil, 3, 1},
		{errors.New("test"), 3, 1},
		{errors.New("database is locked"), 3, 3},
	}

	for i, s := range scenarios {
		lastAttempt := 0

		err := baseLockRetry(func(attempt int) error {
			lastAttempt = attempt

			if attempt < s.failUntilAttempt {
				return s.err
			}

			return nil
		}, s.failUntilAttempt+2)

		if lastAttempt != s.expectedAttempts {
			t.Errorf("[%d] Expected lastAttempt to be %d, got %d", i, s.expectedAttempts, lastAttempt)
		}

		if s.failUntilAttempt == s.expectedAttempts && err != nil {
			t.Errorf("[%d] Expected nil, got err %v", i, err)
			continue
		}

		if s.failUntilAttempt != s.expectedAttempts && s.err != nil && err == nil {
			t.Errorf("[%d] Expected error %q, got nil", i, s.err)
			continue
		}
	}
}
