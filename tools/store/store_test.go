package store_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pocketbase/pocketbase/tools/store"
)

func TestNew(t *testing.T) {
	data := map[string]int{"test1": 1, "test2": 2}
	originalRawData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	s := store.New(data)
	s.Set("test3", 3) // add 1 item
	s.Remove("test1") // remove 1 item

	// check if data was shallow copied
	rawData, _ := json.Marshal(data)
	if !bytes.Equal(originalRawData, rawData) {
		t.Fatalf("Expected data \n%s, \ngot \n%s", originalRawData, rawData)
	}

	if s.Has("test1") {
		t.Fatalf("Expected test1 to be deleted, got %v", s.Get("test1"))
	}

	if v := s.Get("test2"); v != 2 {
		t.Fatalf("Expected test2 to be %v, got %v", 2, v)
	}

	if v := s.Get("test3"); v != 3 {
		t.Fatalf("Expected test3 to be %v, got %v", 3, v)
	}
}

func TestReset(t *testing.T) {
	s := store.New(map[string]int{"test1": 1})

	data := map[string]int{"test2": 2}
	originalRawData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	s.Reset(data)
	s.Set("test3", 3)

	// check if data was shallow copied
	rawData, _ := json.Marshal(data)
	if !bytes.Equal(originalRawData, rawData) {
		t.Fatalf("Expected data \n%s, \ngot \n%s", originalRawData, rawData)
	}

	if s.Has("test1") {
		t.Fatalf("Expected test1 to be deleted, got %v", s.Get("test1"))
	}

	if v := s.Get("test2"); v != 2 {
		t.Fatalf("Expected test2 to be %v, got %v", 2, v)
	}

	if v := s.Get("test3"); v != 3 {
		t.Fatalf("Expected test3 to be %v, got %v", 3, v)
	}
}

func TestLength(t *testing.T) {
	s := store.New(map[string]int{"test1": 1})
	s.Set("test2", 2)

	if v := s.Length(); v != 2 {
		t.Fatalf("Expected length %d, got %d", 2, v)
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

func TestGetAll(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}

	s := store.New(data)

	// fetch and delete each key to make sure that it was shallow copied
	result := s.GetAll()
	for k := range result {
		delete(result, k)
	}

	// refetch again
	result = s.GetAll()

	if len(result) != len(data) {
		t.Fatalf("Expected %d, got %d items", len(data), len(result))
	}

	for k := range result {
		if result[k] != data[k] {
			t.Fatalf("Expected %s to be %v, got %v", k, data[k], result[k])
		}
	}
}

func TestSet(t *testing.T) {
	s := store.Store[int]{}

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
	s := store.Store[int]{}

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
