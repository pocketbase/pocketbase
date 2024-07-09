package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JsonMap defines a map that is safe for json and db read/write.
type JsonMap map[string]any

// MarshalJSON implements the [json.Marshaler] interface.
func (m JsonMap) MarshalJSON() ([]byte, error) {
	type alias JsonMap // prevent recursion

	// initialize an empty map to ensure that `{}` is returned as json
	if m == nil {
		m = JsonMap{}
	}

	return json.Marshal(alias(m))
}

// Get retrieves a single value from the current JsonMap.
//
// This helper was added primarily to assist the goja integration since custom map types
// don't have direct access to the map keys (https://pkg.go.dev/github.com/dop251/goja#hdr-Maps_with_methods).
func (m JsonMap) Get(key string) any {
	return m[key]
}

// Set sets a single value in the current JsonMap.
//
// This helper was added primarily to assist the goja integration since custom map types
// don't have direct access to the map keys (https://pkg.go.dev/github.com/dop251/goja#hdr-Maps_with_methods).
func (m JsonMap) Set(key string, value any) {
	m[key] = value
}

// Value implements the [driver.Valuer] interface.
func (m JsonMap) Value() (driver.Value, error) {
	data, err := json.Marshal(m)

	return string(data), err
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current `JsonMap` instance.
func (m *JsonMap) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case nil:
		// no cast needed
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JsonMap value: %q", value)
	}

	if len(data) == 0 {
		data = []byte("{}")
	}

	return json.Unmarshal(data, m)
}
