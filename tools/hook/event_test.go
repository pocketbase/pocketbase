package hook

import "testing"

func TestEventNext(t *testing.T) {
	calls := 0

	e := Event{}

	if e.nextFunc() != nil {
		t.Fatalf("Expected nextFunc to be nil")
	}

	e.setNextFunc(func() error {
		calls++
		return nil
	})

	if e.nextFunc() == nil {
		t.Fatalf("Expected nextFunc to be non-nil")
	}

	e.Next()
	e.Next()

	if calls != 2 {
		t.Fatalf("Expected %d calls, got %d", 2, calls)
	}
}
