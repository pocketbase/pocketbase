package routine_test

import (
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
}
