package hook

import (
	"errors"
	"testing"
)

func TestHookAddAndPreAdd(t *testing.T) {
	h := Hook[int]{}

	if total := len(h.handlers); total != 0 {
		t.Fatalf("Expected no handlers, found %d", total)
	}

	triggerSequence := ""

	f1 := func(data int) error { triggerSequence += "f1"; return nil }
	f2 := func(data int) error { triggerSequence += "f2"; return nil }
	f3 := func(data int) error { triggerSequence += "f3"; return nil }
	f4 := func(data int) error { triggerSequence += "f4"; return nil }

	h.Add(f1)
	h.Add(f2)
	h.PreAdd(f3)
	h.PreAdd(f4)
	h.Trigger(1)

	if total := len(h.handlers); total != 4 {
		t.Fatalf("Expected %d handlers, found %d", 4, total)
	}

	expectedTriggerSequence := "f4f3f1f2"

	if triggerSequence != expectedTriggerSequence {
		t.Fatalf("Expected trigger sequence %s, got %s", expectedTriggerSequence, triggerSequence)
	}
}

func TestHookRemove(t *testing.T) {
	h := Hook[int]{}

	h1Called := false
	h2Called := false

	id1 := h.Add(func(data int) error { h1Called = true; return nil })
	h.Add(func(data int) error { h2Called = true; return nil })

	h.Remove("missing") // should do nothing and not panic

	if total := len(h.handlers); total != 2 {
		t.Fatalf("Expected %d handlers, got %d", 2, total)
	}

	h.Remove(id1)

	if total := len(h.handlers); total != 1 {
		t.Fatalf("Expected %d handlers, got %d", 1, total)
	}

	if err := h.Trigger(1); err != nil {
		t.Fatal(err)
	}

	if h1Called {
		t.Fatalf("Expected hook 1 to be removed and not called")
	}

	if !h2Called {
		t.Fatalf("Expected hook 2 to be called")
	}
}

func TestHookRemoveAll(t *testing.T) {
	h := Hook[int]{}

	h.RemoveAll() // should do nothing and not panic

	h.Add(func(data int) error { return nil })
	h.Add(func(data int) error { return nil })

	if total := len(h.handlers); total != 2 {
		t.Fatalf("Expected 2 handlers before RemoveAll, found %d", total)
	}

	h.RemoveAll()

	if total := len(h.handlers); total != 0 {
		t.Fatalf("Expected no handlers after RemoveAll, found %d", total)
	}
}

func TestHookTrigger(t *testing.T) {
	err1 := errors.New("demo")
	err2 := errors.New("demo")

	scenarios := []struct {
		handlers      []Handler[int]
		expectedError error
	}{
		{
			[]Handler[int]{
				func(data int) error { return nil },
				func(data int) error { return nil },
			},
			nil,
		},
		{
			[]Handler[int]{
				func(data int) error { return nil },
				func(data int) error { return err1 },
				func(data int) error { return err2 },
			},
			err1,
		},
	}

	for i, scenario := range scenarios {
		h := Hook[int]{}
		for _, handler := range scenario.handlers {
			h.Add(handler)
		}
		result := h.Trigger(1)
		if result != scenario.expectedError {
			t.Fatalf("(%d) Expected %v, got %v", i, scenario.expectedError, result)
		}
	}
}

func TestHookTriggerStopPropagation(t *testing.T) {
	called1 := false
	f1 := func(data int) error { called1 = true; return nil }

	called2 := false
	f2 := func(data int) error { called2 = true; return nil }

	called3 := false
	f3 := func(data int) error { called3 = true; return nil }

	called4 := false
	f4 := func(data int) error { called4 = true; return StopPropagation }

	called5 := false
	f5 := func(data int) error { called5 = true; return nil }

	called6 := false
	f6 := func(data int) error { called6 = true; return nil }

	h := Hook[int]{}
	h.Add(f1)
	h.Add(f2)

	result := h.Trigger(123, f3, f4, f5, f6)

	if result != nil {
		t.Fatalf("Expected nil after StopPropagation, got %v", result)
	}

	// ensure that the trigger handler were not persisted
	if total := len(h.handlers); total != 2 {
		t.Fatalf("Expected 2 handlers, found %d", total)
	}

	scenarios := []struct {
		called   bool
		expected bool
	}{
		{called1, true},
		{called2, true},
		{called3, true},
		{called4, true}, // StopPropagation
		{called5, false},
		{called6, false},
	}
	for i, scenario := range scenarios {
		if scenario.called != scenario.expected {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expected, scenario.called)
		}
	}
}
