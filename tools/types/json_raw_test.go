package types_test

import (
	"database/sql/driver"
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestParseJsonRaw(t *testing.T) {
	scenarios := []struct {
		value       any
		expectError bool
		expectJson  string
	}{
		{nil, false, `null`},
		{``, false, `null`},
		{[]byte{}, false, `null`},
		{types.JsonRaw{}, false, `null`},
		{`{}`, false, `{}`},
		{`[]`, false, `[]`},
		{123, false, `123`},
		{`""`, false, `""`},
		{`test`, false, `test`},
		{`{"invalid"`, false, `{"invalid"`}, // treated as a byte casted string
		{`{"test":1}`, false, `{"test":1}`},
		{[]byte(`[1,2,3]`), false, `[1,2,3]`},
		{[]int{1, 2, 3}, false, `[1,2,3]`},
		{map[string]int{"test": 1}, false, `{"test":1}`},
	}

	for i, s := range scenarios {
		raw, parseErr := types.ParseJsonRaw(s.value)
		hasErr := parseErr != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected %v, got %v (%v)", i, s.expectError, hasErr, parseErr)
			continue
		}

		result, _ := raw.MarshalJSON()

		if string(result) != s.expectJson {
			t.Errorf("(%d) Expected %s, got %v", i, s.expectJson, string(result))
		}
	}
}

func TestJsonRawString(t *testing.T) {
	scenarios := []struct {
		json     types.JsonRaw
		expected string
	}{
		{nil, ``},
		{types.JsonRaw{}, ``},
		{types.JsonRaw([]byte(`123`)), `123`},
		{types.JsonRaw(`{"demo":123}`), `{"demo":123}`},
	}

	for i, s := range scenarios {
		result := s.json.String()
		if result != s.expected {
			t.Errorf("(%d) Expected %q, got %q", i, s.expected, result)
		}
	}
}

func TestJsonRawMarshalJSON(t *testing.T) {
	scenarios := []struct {
		json     types.JsonRaw
		expected string
	}{
		{nil, `null`},
		{types.JsonRaw{}, `null`},
		{types.JsonRaw([]byte(`123`)), `123`},
		{types.JsonRaw(`{"demo":123}`), `{"demo":123}`},
	}

	for i, s := range scenarios {
		result, err := s.json.MarshalJSON()
		if err != nil {
			t.Errorf("(%d) %v", i, err)
			continue
		}

		if string(result) != s.expected {
			t.Errorf("(%d) Expected %q, got %q", i, s.expected, string(result))
		}
	}
}

func TestJsonRawUnmarshalJSON(t *testing.T) {
	scenarios := []struct {
		json         []byte
		expectString string
	}{
		{nil, ""},
		{[]byte{0, 1, 2}, "\x00\x01\x02"},
		{[]byte("123"), "123"},
		{[]byte("test"), "test"},
		{[]byte(`{"test":123}`), `{"test":123}`},
	}

	for i, s := range scenarios {
		raw := types.JsonRaw{}
		err := raw.UnmarshalJSON(s.json)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
			continue
		}

		if raw.String() != s.expectString {
			t.Errorf("(%d) Expected %q, got %q", i, s.expectString, raw.String())
		}
	}
}

func TestJsonRawValue(t *testing.T) {
	scenarios := []struct {
		json     types.JsonRaw
		expected driver.Value
	}{
		{nil, nil},
		{types.JsonRaw{}, nil},
		{types.JsonRaw(``), nil},
		{types.JsonRaw(`test`), `test`},
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

func TestJsonRawScan(t *testing.T) {
	scenarios := []struct {
		value       any
		expectError bool
		expectJson  string
	}{
		{nil, false, `null`},
		{``, false, `null`},
		{[]byte{}, false, `null`},
		{types.JsonRaw{}, false, `null`},
		{types.JsonRaw(`test`), false, `test`},
		{`{}`, false, `{}`},
		{`[]`, false, `[]`},
		{123, false, `123`},
		{`""`, false, `""`},
		{`test`, false, `test`},
		{`{"invalid"`, false, `{"invalid"`}, // treated as a byte casted string
		{`{"test":1}`, false, `{"test":1}`},
		{[]byte(`[1,2,3]`), false, `[1,2,3]`},
		{[]int{1, 2, 3}, false, `[1,2,3]`},
		{map[string]int{"test": 1}, false, `{"test":1}`},
	}

	for i, s := range scenarios {
		raw := types.JsonRaw{}
		scanErr := raw.Scan(s.value)
		hasErr := scanErr != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected %v, got %v (%v)", i, s.expectError, hasErr, scanErr)
			continue
		}

		result, _ := raw.MarshalJSON()

		if string(result) != s.expectJson {
			t.Errorf("(%d) Expected %s, got %v", i, s.expectJson, string(result))
		}
	}
}
