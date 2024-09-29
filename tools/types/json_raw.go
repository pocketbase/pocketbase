package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONRaw defines a json value type that is safe for db read/write.
type JSONRaw []byte

// ParseJSONRaw creates a new JSONRaw instance from the provided value
// (could be JSONRaw, int, float, string, []byte, etc.).
func ParseJSONRaw(value any) (JSONRaw, error) {
	result := JSONRaw{}
	err := result.Scan(value)
	return result, err
}

// String returns the current JSONRaw instance as a json encoded string.
func (j JSONRaw) String() string {
	raw, _ := j.MarshalJSON()
	return string(raw)
}

// MarshalJSON implements the [json.Marshaler] interface.
func (j JSONRaw) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}

	return j, nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (j *JSONRaw) UnmarshalJSON(b []byte) error {
	if j == nil {
		return errors.New("JSONRaw: UnmarshalJSON on nil pointer")
	}

	*j = append((*j)[0:0], b...)

	return nil
}

// Value implements the [driver.Valuer] interface.
func (j JSONRaw) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}

	return j.String(), nil
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current JSONRaw instance.
func (j *JSONRaw) Scan(value any) error {
	var data []byte

	switch v := value.(type) {
	case nil:
		// no cast is needed
	case []byte:
		if len(v) != 0 {
			data = v
		}
	case string:
		if v != "" {
			data = []byte(v)
		}
	case JSONRaw:
		if len(v) != 0 {
			data = []byte(v)
		}
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		data = bytes
	}

	return j.UnmarshalJSON(data)
}
