package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONMap defines a map that is safe for json and db read/write.
type JSONMap[T any] map[string]T

// MarshalJSON implements the [json.Marshaler] interface.
func (m JSONMap[T]) MarshalJSON() ([]byte, error) {
	type alias JSONMap[T] // prevent recursion

	// initialize an empty map to ensure that `{}` is returned as json
	if m == nil {
		m = JSONMap[T]{}
	}

	return json.Marshal(alias(m))
}

// String returns the string representation of the current json map.
func (m JSONMap[T]) String() string {
	v, _ := m.MarshalJSON()
	return string(v)
}

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

// Value implements the [driver.Valuer] interface.
func (m JSONMap[T]) Value() (driver.Value, error) {
	data, err := json.Marshal(m)

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

	return json.Unmarshal(data, m)
}
