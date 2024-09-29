package types_test

import (
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/tools/types"
)

func TestParseJSONRaw(t *testing.T) {
	scenarios := []struct {
		value       any
		expectError bool
		expectJSON  string
	}{
		{nil, false, `null`},
		{``, false, `null`},
		{[]byte{}, false, `null`},
		{types.JSONRaw{}, false, `null`},
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
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			raw, parseErr := types.ParseJSONRaw(s.value)

			hasErr := parseErr != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected %v, got %v (%v)", s.expectError, hasErr, parseErr)
			}

			result, _ := raw.MarshalJSON()

			if string(result) != s.expectJSON {
				t.Fatalf("Expected %s, got %s", s.expectJSON, string(result))
			}
		})
	}
}

func TestJSONRawString(t *testing.T) {
	scenarios := []struct {
		json     types.JSONRaw
		expected string
	}{
		{nil, `null`},
		{types.JSONRaw{}, `null`},
		{types.JSONRaw([]byte(`123`)), `123`},
		{types.JSONRaw(`{"demo":123}`), `{"demo":123}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.expected), func(t *testing.T) {
			result := s.json.String()
			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}

func TestJSONRawMarshalJSON(t *testing.T) {
	scenarios := []struct {
		json     types.JSONRaw
		expected string
	}{
		{nil, `null`},
		{types.JSONRaw{}, `null`},
		{types.JSONRaw([]byte(`123`)), `123`},
		{types.JSONRaw(`{"demo":123}`), `{"demo":123}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.expected), func(t *testing.T) {
			result, err := s.json.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if string(result) != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, string(result))
			}
		})
	}
}

func TestJSONRawUnmarshalJSON(t *testing.T) {
	scenarios := []struct {
		json         []byte
		expectString string
	}{
		{nil, `null`},
		{[]byte{0, 1, 2}, "\x00\x01\x02"},
		{[]byte("123"), "123"},
		{[]byte("test"), "test"},
		{[]byte(`{"test":123}`), `{"test":123}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.expectString), func(t *testing.T) {
			raw := types.JSONRaw{}

			err := raw.UnmarshalJSON(s.json)
			if err != nil {
				t.Fatal(err)
			}

			if raw.String() != s.expectString {
				t.Fatalf("Expected %q, got %q", s.expectString, raw.String())
			}
		})
	}
}

func TestJSONRawValue(t *testing.T) {
	scenarios := []struct {
		json     types.JSONRaw
		expected driver.Value
	}{
		{nil, nil},
		{types.JSONRaw{}, nil},
		{types.JSONRaw(``), nil},
		{types.JSONRaw(`test`), `test`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.json), func(t *testing.T) {
			result, err := s.json.Value()
			if err != nil {
				t.Fatal(err)
			}

			if result != s.expected {
				t.Fatalf("Expected %s, got %v", s.expected, result)
			}
		})
	}
}

func TestJSONRawScan(t *testing.T) {
	scenarios := []struct {
		value       any
		expectError bool
		expectJSON  string
	}{
		{nil, false, `null`},
		{``, false, `null`},
		{[]byte{}, false, `null`},
		{types.JSONRaw{}, false, `null`},
		{types.JSONRaw(`test`), false, `test`},
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
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			raw := types.JSONRaw{}
			scanErr := raw.Scan(s.value)
			hasErr := scanErr != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected %v, got %v (%v)", s.expectError, hasErr, scanErr)
			}

			result, _ := raw.MarshalJSON()

			if string(result) != s.expectJSON {
				t.Fatalf("Expected %s, got %v", s.expectJSON, string(result))
			}
		})
	}
}
