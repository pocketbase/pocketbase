package hook

import (
	"strings"
	"testing"
)

type mockTagsEvent struct {
	Event
	tags []string
}

func (m mockTagsEvent) Tags() []string {
	return m.tags
}

func TestTaggedHook(t *testing.T) {
	calls := ""

	base := &Hook[*mockTagsEvent]{}
	base.BindFunc(func(e *mockTagsEvent) error { calls += "f0"; return e.Next() })

	hA := NewTaggedHook(base)
	hA.BindFunc(func(e *mockTagsEvent) error { calls += "a1"; return e.Next() })
	hA.Bind(&Handler[*mockTagsEvent]{
		Func:     func(e *mockTagsEvent) error { calls += "a2"; return e.Next() },
		Priority: -1,
	})

	hB := NewTaggedHook(base, "b1", "b2")
	hB.BindFunc(func(e *mockTagsEvent) error { calls += "b1"; return e.Next() })
	hB.Bind(&Handler[*mockTagsEvent]{
		Func:     func(e *mockTagsEvent) error { calls += "b2"; return e.Next() },
		Priority: -2,
	})

	hC := NewTaggedHook(base, "c1", "c2")
	hC.BindFunc(func(e *mockTagsEvent) error { calls += "c1"; return e.Next() })
	hC.Bind(&Handler[*mockTagsEvent]{
		Func:     func(e *mockTagsEvent) error { calls += "c2"; return e.Next() },
		Priority: -3,
	})

	scenarios := []struct {
		event         *mockTagsEvent
		expectedCalls string
	}{
		{
			&mockTagsEvent{},
			"a2f0a1",
		},
		{
			&mockTagsEvent{tags: []string{"missing"}},
			"a2f0a1",
		},
		{
			&mockTagsEvent{tags: []string{"b2"}},
			"b2a2f0a1b1",
		},
		{
			&mockTagsEvent{tags: []string{"c1"}},
			"c2a2f0a1c1",
		},
		{
			&mockTagsEvent{tags: []string{"b1", "c2"}},
			"c2b2a2f0a1b1c1",
		},
	}

	for _, s := range scenarios {
		t.Run(strings.Join(s.event.tags, "_"), func(t *testing.T) {
			calls = "" // reset

			err := base.Trigger(s.event)
			if err != nil {
				t.Fatalf("Unexpected trigger error: %v", err)
			}

			if calls != s.expectedCalls {
				t.Fatalf("Expected calls sequence %q, got %q", s.expectedCalls, calls)
			}
		})
	}
}
