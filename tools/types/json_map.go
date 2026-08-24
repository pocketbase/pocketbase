package types

import (
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

// JSONMap defines a map that is safe for json and db read/write.
type JSONMap[T any] map[string]T

// Get retrieves a single value from the current JSONMap[T].
//
// This helper was added primarily to assist the goja integration since custom map types
// don't have direct access to the map keys (https://pkg.go.dev/github.com/dop251/goja#hdr-Maps_with_methods).
func (m JSONMap[T]) Get(key string) T {
	return m[key]
}

// Set sets a single value in the current JSONMap[T].
//
// This helper was added primarily to assist the goja integration since custom map types
// don't have direct access to the map keys (https://pkg.go.dev/github.com/dop251/goja#hdr-Maps_with_methods).
func (m JSONMap[T]) Set(key string, value T) {
	m[key] = value
}

// MarshalJSON implements the [json.Marshaler] interface.
func (m JSONMap[T]) MarshalJSON() ([]byte, error) {
	type alias JSONMap[T] // prevent recursion

	// note: forces the Deterministic and AllowInvalidUTF8 options to
	// ensure consistent output in mixed json v1 and v2 configurations
	return json.Marshal(
		alias(m),
		json.Deterministic(true),
		jsontext.AllowInvalidUTF8(true),
	)
}

// String returns the string representation of the current json map.
func (m JSONMap[T]) String() string {
	v, _ := m.MarshalJSON()
	return string(v)
}

// Value implements the [driver.Valuer] interface.
func (m JSONMap[T]) Value() (driver.Value, error) {
	data, err := m.MarshalJSON()
	return string(data), err
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current JSONMap[T] instance.
func (m *JSONMap[T]) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case nil:
		// no cast needed
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONMap[T] value: %q", value)
	}

	if len(data) == 0 {
		data = []byte("{}")
	}

	err := json.Unmarshal(data, m)
	if err != nil {
		// reset because jsonv2 performs streaming decoding and mutates the dst even on error
		*m = JSONMap[T]{}
		return err
	}

	return nil
}
