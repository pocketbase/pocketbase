package types_test

import (
	"database/sql/driver"
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestJsonMapMarshalJSON(t *testing.T) {
	scenarios := []struct {
		json     types.JsonMap
		expected string
	}{
		{nil, "{}"},
		{types.JsonMap{}, `{}`},
		{types.JsonMap{"test1": 123, "test2": "lorem"}, `{"test1":123,"test2":"lorem"}`},
		{types.JsonMap{"test": []int{1, 2, 3}}, `{"test":[1,2,3]}`},
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

func TestJsonMapGet(t *testing.T) {
	scenarios := []struct {
		json     types.JsonMap
		key      string
		expected any
	}{
		{nil, "test", nil},
		{types.JsonMap{"test": 123}, "test", 123},
		{types.JsonMap{"test": 123}, "missing", nil},
	}

	for i, s := range scenarios {
		result := s.json.Get(s.key)
		if result != s.expected {
			t.Errorf("(%d) Expected %s, got %v", i, s.expected, result)
		}
	}
}

func TestJsonMapSet(t *testing.T) {
	scenarios := []struct {
		key   string
		value any
	}{
		{"a", nil},
		{"a", 123},
		{"b", "test"},
	}

	for i, s := range scenarios {
		j := types.JsonMap{}

		j.Set(s.key, s.value)

		if v := j[s.key]; v != s.value {
			t.Errorf("(%d) Expected %s, got %v", i, s.value, v)
		}
	}
}

func TestJsonMapValue(t *testing.T) {
	scenarios := []struct {
		json     types.JsonMap
		expected driver.Value
	}{
		{nil, `{}`},
		{types.JsonMap{}, `{}`},
		{types.JsonMap{"test1": 123, "test2": "lorem"}, `{"test1":123,"test2":"lorem"}`},
		{types.JsonMap{"test": []int{1, 2, 3}}, `{"test":[1,2,3]}`},
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

func TestJsonArrayMapScan(t *testing.T) {
	scenarios := []struct {
		value       any
		expectError bool
		expectJson  string
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
		arr := types.JsonMap{}
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
