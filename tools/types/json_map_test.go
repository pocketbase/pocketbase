package types_test

import (
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestJSONMapMarshalJSON(t *testing.T) {
	scenarios := []struct {
		json     types.JSONMap[any]
		expected string
	}{
		{nil, "{}"},
		{types.JSONMap[any]{}, `{}`},
		{types.JSONMap[any]{"test1": 123, "test2": "lorem"}, `{"test1":123,"test2":"lorem"}`},
		{types.JSONMap[any]{"test": []int{1, 2, 3}}, `{"test":[1,2,3]}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.expected), func(t *testing.T) {
			result, err := s.json.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if string(result) != s.expected {
				t.Fatalf("Expected\n%s\ngot\n%s", s.expected, result)
			}
		})
	}
}

func TestJSONMapMarshalString(t *testing.T) {
	scenarios := []struct {
		json     types.JSONMap[any]
		expected string
	}{
		{nil, "{}"},
		{types.JSONMap[any]{}, `{}`},
		{types.JSONMap[any]{"test1": 123, "test2": "lorem"}, `{"test1":123,"test2":"lorem"}`},
		{types.JSONMap[any]{"test": []int{1, 2, 3}}, `{"test":[1,2,3]}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.expected), func(t *testing.T) {
			result := s.json.String()

			if result != s.expected {
				t.Fatalf("Expected\n%s\ngot\n%s", s.expected, result)
			}
		})
	}
}

func TestJSONMapGet(t *testing.T) {
	scenarios := []struct {
		json     types.JSONMap[any]
		key      string
		expected any
	}{
		{nil, "test", nil},
		{types.JSONMap[any]{"test": 123}, "test", 123},
		{types.JSONMap[any]{"test": 123}, "missing", nil},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.key), func(t *testing.T) {
			result := s.json.Get(s.key)
			if result != s.expected {
				t.Fatalf("Expected %s, got %#v", s.expected, result)
			}
		})
	}
}

func TestJSONMapSet(t *testing.T) {
	scenarios := []struct {
		key   string
		value any
	}{
		{"a", nil},
		{"a", 123},
		{"b", "test"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.key), func(t *testing.T) {
			j := types.JSONMap[any]{}

			j.Set(s.key, s.value)

			if v := j[s.key]; v != s.value {
				t.Fatalf("Expected %s, got %#v", s.value, v)
			}
		})
	}
}

func TestJSONMapValue(t *testing.T) {
	scenarios := []struct {
		json     types.JSONMap[any]
		expected driver.Value
	}{
		{nil, `{}`},
		{types.JSONMap[any]{}, `{}`},
		{types.JSONMap[any]{"test1": 123, "test2": "lorem"}, `{"test1":123,"test2":"lorem"}`},
		{types.JSONMap[any]{"test": []int{1, 2, 3}}, `{"test":[1,2,3]}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.expected), func(t *testing.T) {
			result, err := s.json.Value()
			if err != nil {
				t.Fatal(err)
			}

			if result != s.expected {
				t.Fatalf("Expected %s, got %#v", s.expected, result)
			}
		})
	}
}

func TestJSONArrayMapScan(t *testing.T) {
	scenarios := []struct {
		value       any
		expectError bool
		expectJSON  string
	}{
		{``, false, `{}`},
		{nil, false, `{}`},
		{[]byte{}, false, `{}`},
		{`{}`, false, `{}`},
		{123, true, `{}`},
		{`""`, true, `{}`},
		{`invalid_json`, true, `{}`},
		{`"test"`, true, `{}`},
		{`1,2,3`, true, `{}`},
		{`{"test": 1`, true, `{}`},
		{`{"test": 1}`, false, `{"test":1}`},
		{[]byte(`{"test": 1}`), false, `{"test":1}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			arr := types.JSONMap[any]{}
			scanErr := arr.Scan(s.value)

			hasErr := scanErr != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected %v, got %v (%v)", s.expectError, hasErr, scanErr)
			}

			result, _ := arr.MarshalJSON()

			if string(result) != s.expectJSON {
				t.Fatalf("Expected %s, got %s", s.expectJSON, result)
			}
		})
	}
}
