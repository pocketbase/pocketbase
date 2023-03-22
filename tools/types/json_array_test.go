package types_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestJsonArrayMarshalJSON(t *testing.T) {
	scenarios := []struct {
		json     json.Marshaler
		expected string
	}{
		{new(types.JsonArray[any]), "[]"},
		{types.JsonArray[any]{}, `[]`},
		{types.JsonArray[int]{1, 2, 3}, `[1,2,3]`},
		{types.JsonArray[string]{"test1", "test2", "test3"}, `["test1","test2","test3"]`},
		{types.JsonArray[any]{1, "test"}, `[1,"test"]`},
	}

	for i, s := range scenarios {
		result, err := s.json.MarshalJSON()
		if err != nil {
			t.Errorf("(%d) %v", i, err)
			continue
		}
		if string(result) != s.expected {
			t.Errorf("(%d) Expected %s, got %s", i, s.expected, string(result))
		}
	}
}

func TestJsonArrayValue(t *testing.T) {
	scenarios := []struct {
		json     driver.Valuer
		expected driver.Value
	}{
		{new(types.JsonArray[any]), `[]`},
		{types.JsonArray[any]{}, `[]`},
		{types.JsonArray[int]{1, 2, 3}, `[1,2,3]`},
		{types.JsonArray[string]{"test1", "test2", "test3"}, `["test1","test2","test3"]`},
		{types.JsonArray[any]{1, "test"}, `[1,"test"]`},
	}

	for i, s := range scenarios {
		result, err := s.json.Value()
		if err != nil {
			t.Errorf("(%d) %v", i, err)
			continue
		}
		if result != s.expected {
			t.Errorf("(%d) Expected %s, got %v", i, s.expected, result)
		}
	}
}

func TestJsonArrayScan(t *testing.T) {
	scenarios := []struct {
		value       any
		expectError bool
		expectJson  string
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
		arr := types.JsonArray[any]{}
		scanErr := arr.Scan(s.value)

		hasErr := scanErr != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected %v, got %v (%v)", i, s.expectError, hasErr, scanErr)
			continue
		}

		result, _ := arr.MarshalJSON()

		if string(result) != s.expectJson {
			t.Errorf("(%d) Expected %s, got %v", i, s.expectJson, string(result))
		}
	}
}
