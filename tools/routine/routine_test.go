package routine_test

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/pocketbase/pocketbase/tools/routine"
)

func TestFireAndForget(t *testing.T) {
	called := false

	fn := func() {
		called = true
		panic("test_recover")
	}

	wg := &sync.WaitGroup{}

	routine.FireAndForget(fn, wg)

	wg.Wait()

	if !called {
		t.Fatal("Expected fn to be called.")
	}
}

func TestSafeWrap(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		called := false

		fn := func() error {
			called = true
			panic("test_recover")
		}

		err := routine.SafeWrap(fn)()

		if !called {
			t.Fatal("Expected fn to be called.")
		}

		if err == nil {
			t.Fatal("Expected fn panic to be converted to error")
		}

		if !strings.Contains(err.Error(), "test_recover") {
			t.Fatal("Expected the returned error to contain the recovered panic value")
		}
	})

	t.Run("regular error", func(t *testing.T) {
		called := false

		fn := func() error {
			called = true
			return errors.New("test_error")
		}

		err := routine.SafeWrap(fn)()

		if !called {
			t.Fatal("Expected fn to be called.")
		}

		if err == nil {
			t.Fatal("Expected to return the wrapped fn error")
		}

		errStr := err.Error()
		expected := "test_error"
		if errStr != expected {
			t.Fatalf("Expected %s error, got %s", expected, errStr)
		}
	})
}
