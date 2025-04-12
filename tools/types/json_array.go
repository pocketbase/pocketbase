package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONArray defines a slice that is safe for json and db read/write.
type JSONArray[T any] []T

// internal alias to prevent recursion during marshalization.
type jsonArrayAlias[T any] JSONArray[T]

// MarshalJSON implements the [json.Marshaler] interface.
func (m JSONArray[T]) MarshalJSON() ([]byte, error) {
	// initialize an empty map to ensure that `[]` is returned as json
	if m == nil {
		m = JSONArray[T]{}
	}

	return json.Marshal(jsonArrayAlias[T](m))
}

// String returns the string representation of the current json array.
func (m JSONArray[T]) String() string {
	v, _ := m.MarshalJSON()
	return string(v)
}

// Value implements the [driver.Valuer] interface.
func (m JSONArray[T]) Value() (driver.Value, error) {
	data, err := json.Marshal(m)

	return string(data), err
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current JSONArray[T] instance.
func (m *JSONArray[T]) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case nil:
		// no cast needed
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONArray value: %q", value)
	}

	if len(data) == 0 {
		data = []byte("[]")
	}

	return json.Unmarshal(data, m)
}
