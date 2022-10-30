package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JsonArray defines a slice that is safe for json and db read/write.
type JsonArray []any

// MarshalJSON implements the [json.Marshaler] interface.
func (m JsonArray) MarshalJSON() ([]byte, error) {
	type alias JsonArray // prevent recursion

	// initialize an empty map to ensure that `[]` is returned as json
	if m == nil {
		m = JsonArray{}
	}

	return json.Marshal(alias(m))
}

// Value implements the [driver.Valuer] interface.
func (m JsonArray) Value() (driver.Value, error) {
	data, err := json.Marshal(m)

	return string(data), err
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current `JsonArray` instance.
func (m *JsonArray) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case nil:
		// no cast needed
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("Failed to unmarshal JsonArray value: %q.", value)
	}

	if len(data) == 0 {
		data = []byte("[]")
	}

	return json.Unmarshal(data, m)
}
