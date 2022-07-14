package store_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/store"
)

func TestNew(t *testing.T) {
	s := store.New(map[string]int{"test": 1})

	if s.Get("test") != 1 {
		t.Error("Expected the initizialized store map to be loaded")
	}
}

func TestRemoveAll(t *testing.T) {
	s := store.New(map[string]bool{"test1": true, "test2": true})

	keys := []string{"test1", "test2"}

	s.RemoveAll()

	for i, key := range keys {
		if s.Has(key) {
			t.Errorf("(%d) Expected %q to be removed", i, key)
		}
	}
}

func TestRemove(t *testing.T) {
	s := store.New(map[string]bool{"test": true})

	keys := []string{"test", "missing"}

	for i, key := range keys {
		s.Remove(key)
		if s.Has(key) {
			t.Errorf("(%d) Expected %q to be removed", i, key)
		}
	}
}

func TestHas(t *testing.T) {
	s := store.New(map[string]int{"test1": 0, "test2": 1})

	scenarios := []struct {
		key   string
		exist bool
	}{
		{"test1", true},
		{"test2", true},
		{"missing", false},
	}

	for i, scenario := range scenarios {
		exist := s.Has(scenario.key)
		if exist != scenario.exist {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.exist, exist)
		}
	}
}

func TestGet(t *testing.T) {
	s := store.New(map[string]int{"test1": 0, "test2": 1})

	scenarios := []struct {
		key    string
		expect int
	}{
		{"test1", 0},
		{"test2", 1},
		{"missing", 0}, // should auto fallback to the zero value
	}

	for i, scenario := range scenarios {
		val := s.Get(scenario.key)
		if val != scenario.expect {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expect, val)
		}
	}
}

func TestSet(t *testing.T) {
	s := store.New[int](nil)

	data := map[string]int{"test1": 0, "test2": 1, "test3": 3}

	// set values
	for k, v := range data {
		s.Set(k, v)
	}

	// verify that the values are set
	for k, v := range data {
		if !s.Has(k) {
			t.Errorf("Expected key %q", k)
		}

		val := s.Get(k)
		if val != v {
			t.Errorf("Expected %v, got %v for key %q", v, val, k)
		}
	}
}

func TestSetIfLessThanLimit(t *testing.T) {
	s := store.New[int](nil)

	limit := 2

	// set values
	scenarios := []struct {
		key      string
		value    int
		expected bool
	}{
		{"test1", 1, true},
		{"test2", 2, true},
		{"test3", 3, false},
		{"test2", 4, true}, // overwrite
	}

	for i, scenario := range scenarios {
		result := s.SetIfLessThanLimit(scenario.key, scenario.value, limit)

		if result != scenario.expected {
			t.Errorf("(%d) Expected result %v, got %v", i, scenario.expected, result)
		}

		if !scenario.expected && s.Has(scenario.key) {
			t.Errorf("(%d) Expected key %q to not be set", i, scenario.key)
		}

		val := s.Get(scenario.key)
		if scenario.expected && val != scenario.value {
			t.Errorf("(%d) Expected value %v, got %v", i, scenario.value, val)
		}
	}
}
