package types_test

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestJSONArrayMarshalJSON(t *testing.T) {
	scenarios := []struct {
		json     json.Marshaler
		expected string
	}{
		{new(types.JSONArray[any]), "[]"},
		{types.JSONArray[any]{}, `[]`},
		{types.JSONArray[int]{1, 2, 3}, `[1,2,3]`},
		{types.JSONArray[string]{"test1", "test2", "test3"}, `["test1","test2","test3"]`},
		{types.JSONArray[any]{1, "test"}, `[1,"test"]`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.expected), func(t *testing.T) {
			result, err := s.json.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if string(result) != s.expected {
				t.Fatalf("Expected %s, got %s", s.expected, result)
			}
		})
	}
}

func TestJSONArrayString(t *testing.T) {
	scenarios := []struct {
		json     fmt.Stringer
		expected string
	}{
		{new(types.JSONArray[any]), "[]"},
		{types.JSONArray[any]{}, `[]`},
		{types.JSONArray[int]{1, 2, 3}, `[1,2,3]`},
		{types.JSONArray[string]{"test1", "test2", "test3"}, `["test1","test2","test3"]`},
		{types.JSONArray[any]{1, "test"}, `[1,"test"]`},
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

func TestJSONArrayValue(t *testing.T) {
	scenarios := []struct {
		json     driver.Valuer
		expected driver.Value
	}{
		{new(types.JSONArray[any]), `[]`},
		{types.JSONArray[any]{}, `[]`},
		{types.JSONArray[int]{1, 2, 3}, `[1,2,3]`},
		{types.JSONArray[string]{"test1", "test2", "test3"}, `["test1","test2","test3"]`},
		{types.JSONArray[any]{1, "test"}, `[1,"test"]`},
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

func TestJSONArrayScan(t *testing.T) {
	scenarios := []struct {
		value       any
		expectError bool
		expectJSON  string
	}{
		{``, false, `[]`},
		{[]byte{}, false, `[]`},
		{nil, false, `[]`},
		{123, true, `[]`},
		{`""`, true, `[]`},
		{`invalid_json`, true, `[]`},
		{`"test"`, true, `[]`},
		{`1,2,3`, true, `[]`},
		{`[1, 2, 3`, true, `[]`},
		{`[1, 2, 3]`, false, `[1,2,3]`},
		{[]byte(`[1, 2, 3]`), false, `[1,2,3]`},
		{`[1, "test"]`, false, `[1,"test"]`},
		{`[]`, false, `[]`},
	}

	for i, s := range scenarios {
		arr := types.JSONArray[any]{}
		scanErr := arr.Scan(s.value)

		hasErr := scanErr != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected %v, got %v (%v)", i, s.expectError, hasErr, scanErr)
			continue
		}

		result, _ := arr.MarshalJSON()

		if string(result) != s.expectJSON {
			t.Errorf("(%d) Expected %s, got %v", i, s.expectJSON, string(result))
		}
	}
}
