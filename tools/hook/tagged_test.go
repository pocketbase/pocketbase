package hook

import "testing"

type mockTagsData struct {
	tags []string
}

func (m mockTagsData) Tags() []string {
	return m.tags
}

func TestTaggedHook(t *testing.T) {
	triggerSequence := ""

	base := &Hook[mockTagsData]{}
	base.Add(func(data mockTagsData) error { triggerSequence += "f0"; return nil })

	hA := NewTaggedHook(base)
	hA.Add(func(data mockTagsData) error { triggerSequence += "a1"; return nil })
	hA.PreAdd(func(data mockTagsData) error { triggerSequence += "a2"; return nil })

	hB := NewTaggedHook(base, "b1", "b2")
	hB.Add(func(data mockTagsData) error { triggerSequence += "b1"; return nil })
	hB.PreAdd(func(data mockTagsData) error { triggerSequence += "b2"; return nil })

	hC := NewTaggedHook(base, "c1", "c2")
	hC.Add(func(data mockTagsData) error { triggerSequence += "c1"; return nil })
	hC.PreAdd(func(data mockTagsData) error { triggerSequence += "c2"; return nil })

	scenarios := []struct {
		data             mockTagsData
		expectedSequence string
	}{
		{
			mockTagsData{},
			"a2f0a1",
		},
		{
			mockTagsData{[]string{"missing"}},
			"a2f0a1",
		},
		{
			mockTagsData{[]string{"b2"}},
			"b2a2f0a1b1",
		},
		{
			mockTagsData{[]string{"c1"}},
			"c2a2f0a1c1",
		},
		{
			mockTagsData{[]string{"b1", "c2"}},
			"c2b2a2f0a1b1c1",
		},
	}

	for i, s := range scenarios {
		triggerSequence = "" // reset

		err := hA.Trigger(s.data)
		if err != nil {
			t.Fatalf("[%d] Unexpected trigger error: %v", i, err)
		}

		if triggerSequence != s.expectedSequence {
			t.Fatalf("[%d] Expected trigger sequence %s, got %s", i, s.expectedSequence, triggerSequence)
		}
	}
}
