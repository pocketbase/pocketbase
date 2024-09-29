package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pocketbase/pocketbase/tools/security"
)

// NewFieldsList creates a new FieldsList instance with the provided fields.
func NewFieldsList(fields ...Field) FieldsList {
	s := make(FieldsList, 0, len(fields))

	for _, f := range fields {
		s.Add(f)
	}

	return s
}

// FieldsList defines a Collection slice of fields.
type FieldsList []Field

// Clone creates a deep clone of the current list.
func (s FieldsList) Clone() (FieldsList, error) {
	copyRaw, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	result := FieldsList{}
	if err := json.Unmarshal(copyRaw, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// FieldNames returns a slice with the name of all list fields.
func (s FieldsList) FieldNames() []string {
	result := make([]string, len(s))

	for i, field := range s {
		result[i] = field.GetName()
	}

	return result
}

// AsMap returns a map with all registered list field.
// The returned map is indexed with each field name.
func (s FieldsList) AsMap() map[string]Field {
	result := make(map[string]Field, len(s))

	for _, field := range s {
		result[field.GetName()] = field
	}

	return result
}

// GetById returns a single field by its id.
func (s FieldsList) GetById(fieldId string) Field {
	for _, field := range s {
		if field.GetId() == fieldId {
			return field
		}
	}
	return nil
}

// GetByName returns a single field by its name.
func (s FieldsList) GetByName(fieldName string) Field {
	for _, field := range s {
		if field.GetName() == fieldName {
			return field
		}
	}
	return nil
}

// RemoveById removes a single field by its id.
//
// This method does nothing if field with the specified id doesn't exist.
func (s *FieldsList) RemoveById(fieldId string) {
	fields := *s
	for i, field := range fields {
		if field.GetId() == fieldId {
			*s = append(fields[:i], fields[i+1:]...)
			return
		}
	}
}

// RemoveByName removes a single field by its name.
//
// This method does nothing if field with the specified name doesn't exist.
func (s *FieldsList) RemoveByName(fieldName string) {
	fields := *s
	for i, field := range fields {
		if field.GetName() == fieldName {
			*s = append(fields[:i], fields[i+1:]...)
			return
		}
	}
}

// Add adds one or more fields to the current list.
//
// If any of the new fields doesn't have an id it will try to set a
// default one based on its type and name.
//
// If the list already has a field with the same id,
// then the existing field is replaced with the new one.
//
// Otherwise the new field is appended after the other list fields.
func (s *FieldsList) Add(fields ...Field) {
	for _, f := range fields {
		s.add(f)
	}
}

func (s *FieldsList) add(newField Field) {
	newFieldId := newField.GetId()

	// set default id
	if newFieldId == "" {
		if newField.GetName() != "" {
			newFieldId = newField.Type() + crc32Checksum(newField.GetName())
		} else {
			newFieldId = newField.Type() + security.RandomString(5)
		}
		newField.SetId(newFieldId)
	}

	fields := *s

	for i, field := range fields {
		// replace existing
		if newFieldId != "" && field.GetId() == newFieldId {
			(*s)[i] = newField
			return
		}
	}

	// add new field
	*s = append(fields, newField)
}

// String returns the string representation of the current list.
func (s FieldsList) String() string {
	v, _ := json.Marshal(s)
	return string(v)
}

type onlyFieldType struct {
	Type string `json:"type"`
}

type fieldWithType struct {
	Field
	Type string `json:"type"`
}

func (fwt *fieldWithType) UnmarshalJSON(data []byte) error {
	// extract the field type to init a blank factory
	t := &onlyFieldType{}
	if err := json.Unmarshal(data, t); err != nil {
		return fmt.Errorf("failed to unmarshal field type: %w", err)
	}

	factory, ok := Fields[t.Type]
	if !ok {
		return fmt.Errorf("missing or unknown field type in %s", data)
	}

	fwt.Type = t.Type
	fwt.Field = factory()

	// unmarshal the rest of the data into the created field
	if err := json.Unmarshal(data, fwt.Field); err != nil {
		return fmt.Errorf("failed to unmarshal field: %w", err)
	}

	return nil
}

// UnmarshalJSON implements [json.Unmarshaler] and
// loads the provided json data into the current FieldsList.
func (s *FieldsList) UnmarshalJSON(data []byte) error {
	fwts := []fieldWithType{}

	if err := json.Unmarshal(data, &fwts); err != nil {
		return err
	}

	*s = []Field{} // reset

	for _, fwt := range fwts {
		s.Add(fwt.Field)
	}

	return nil
}

// MarshalJSON implements the [json.Marshaler] interface.
func (s FieldsList) MarshalJSON() ([]byte, error) {
	if s == nil {
		s = []Field{} // always init to ensure that it is serialized as empty array
	}

	wrapper := make([]map[string]any, 0, len(s))

	for _, f := range s {
		// precompute the json into a map so that we can append the type to a flatten object
		raw, err := json.Marshal(f)
		if err != nil {
			return nil, err
		}

		data := map[string]any{}
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, err
		}
		data["type"] = f.Type()

		wrapper = append(wrapper, data)
	}

	return json.Marshal(wrapper)
}

// Value implements the [driver.Valuer] interface.
func (s FieldsList) Value() (driver.Value, error) {
	data, err := json.Marshal(s)

	return string(data), err
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current FieldsList instance.
func (s *FieldsList) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case nil:
		// no cast needed
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal FieldsList value %q", value)
	}

	if len(data) == 0 {
		data = []byte("[]")
	}

	return s.UnmarshalJSON(data)
}
