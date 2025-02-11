package store_test

import (
	"bytes"
	"encoding/json"
	"slices"
	"strconv"
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

	for _, scenario := range scenarios {
		t.Run(scenario.key, func(t *testing.T) {
			val := s.Get(scenario.key)
			if val != scenario.expect {
				t.Fatalf("Expected %v, got %v", scenario.expect, val)
			}
		})
	}
}

func TestGetOk(t *testing.T) {
	s := store.New(map[string]int{"test1": 0, "test2": 1})

	scenarios := []struct {
		key         string
		expectValue int
		expectOk    bool
	}{
		{"test1", 0, true},
		{"test2", 1, true},
		{"missing", 0, false}, // should auto fallback to the zero value
	}

	for _, scenario := range scenarios {
		t.Run(scenario.key, func(t *testing.T) {
			val, ok := s.GetOk(scenario.key)

			if ok != scenario.expectOk {
				t.Fatalf("Expected ok %v, got %v", scenario.expectOk, ok)
			}

			if val != scenario.expectValue {
				t.Fatalf("Expected %v, got %v", scenario.expectValue, val)
			}
		})
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

func TestValues(t *testing.T) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}

	values := store.New(data).Values()

	expected := []int{1, 2}

	if len(values) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(values))
	}

	for _, v := range expected {
		if !slices.Contains(values, v) {
			t.Fatalf("Missing value %v in\n%v", v, values)
		}
	}
}

func TestSet(t *testing.T) {
	s := store.Store[string, int]{}

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

func TestSetFunc(t *testing.T) {
	s := store.Store[string, int]{}

	// non existing value
	s.SetFunc("test", func(old int) int {
		if old != 0 {
			t.Fatalf("Expected old value %d, got %d", 0, old)
		}
		return old + 2
	})
	if v := s.Get("test"); v != 2 {
		t.Fatalf("Expected the stored value to be %d, got %d", 2, v)
	}

	// increment existing value
	s.SetFunc("test", func(old int) int {
		if old != 2 {
			t.Fatalf("Expected old value %d, got %d", 2, old)
		}
		return old + 1
	})
	if v := s.Get("test"); v != 3 {
		t.Fatalf("Expected the stored value to be %d, got %d", 3, v)
	}
}

func TestGetOrSet(t *testing.T) {
	s := store.New(map[string]int{
		"test1": 0,
		"test2": 1,
		"test3": 3,
	})

	scenarios := []struct {
		key      string
		value    int
		expected int
	}{
		{"test2", 20, 1},
		{"test3", 2, 3},
		{"test_new", 20, 20},
		{"test_new", 50, 20}, // should return the previously inserted value
	}

	for _, scenario := range scenarios {
		t.Run(scenario.key, func(t *testing.T) {
			result := s.GetOrSet(scenario.key, func() int {
				return scenario.value
			})

			if result != scenario.expected {
				t.Fatalf("Expected %v, got %v", scenario.expected, result)
			}
		})
	}
}

func TestSetIfLessThanLimit(t *testing.T) {
	s := store.Store[string, int]{}

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

func TestUnmarshalJSON(t *testing.T) {
	s := store.Store[string, string]{}
	s.Set("b", "old")   // should be overwritten
	s.Set("c", "test3") // ensures that the old values are not removed

	raw := []byte(`{"a":"test1", "b":"test2"}`)
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatal(err)
	}

	if v := s.Get("a"); v != "test1" {
		t.Fatalf("Expected store.a to be %q, got %q", "test1", v)
	}

	if v := s.Get("b"); v != "test2" {
		t.Fatalf("Expected store.b to be %q, got %q", "test2", v)
	}

	if v := s.Get("c"); v != "test3" {
		t.Fatalf("Expected store.c to be %q, got %q", "test3", v)
	}
}

func TestMarshalJSON(t *testing.T) {
	s := &store.Store[string, string]{}
	s.Set("a", "test1")
	s.Set("b", "test2")

	expected := []byte(`{"a":"test1", "b":"test2"}`)

	result, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(result, expected) {
		t.Fatalf("Expected\n%s\ngot\n%s", expected, result)
	}
}

func TestShrink(t *testing.T) {
	s := &store.Store[string, int]{}

	total := 1000

	for i := 0; i < total; i++ {
		s.Set(strconv.Itoa(i), i)
	}

	if s.Length() != total {
		t.Fatalf("Expected %d items, got %d", total, s.Length())
	}

	// trigger map "shrink"
	for i := 0; i < store.ShrinkThreshold; i++ {
		s.Remove(strconv.Itoa(i))
	}

	// ensure that after the deletion, the new map was copied properly
	if s.Length() != total-store.ShrinkThreshold {
		t.Fatalf("Expected %d items, got %d", total-store.ShrinkThreshold, s.Length())
	}

	for k := range s.GetAll() {
		kInt, err := strconv.Atoi(k)
		if err != nil {
			t.Fatalf("failed to convert %s into int: %v", k, err)
		}
		if kInt < store.ShrinkThreshold {
			t.Fatalf("Key %q should have been deleted", k)
		}
	}
}
