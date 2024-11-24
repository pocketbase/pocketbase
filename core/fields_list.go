package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
)

// NewFieldsList creates a new FieldsList instance with the provided fields.
func NewFieldsList(fields ...Field) FieldsList {
	l := make(FieldsList, 0, len(fields))

	for _, f := range fields {
		l.add(-1, f)
	}

	return l
}

// FieldsList defines a Collection slice of fields.
type FieldsList []Field

// Clone creates a deep clone of the current list.
func (l FieldsList) Clone() (FieldsList, error) {
	copyRaw, err := json.Marshal(l)
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
func (l FieldsList) FieldNames() []string {
	result := make([]string, len(l))

	for i, field := range l {
		result[i] = field.GetName()
	}

	return result
}

// AsMap returns a map with all registered list field.
// The returned map is indexed with each field name.
func (l FieldsList) AsMap() map[string]Field {
	result := make(map[string]Field, len(l))

	for _, field := range l {
		result[field.GetName()] = field
	}

	return result
}

// GetById returns a single field by its id.
func (l FieldsList) GetById(fieldId string) Field {
	for _, field := range l {
		if field.GetId() == fieldId {
			return field
		}
	}
	return nil
}

// GetByName returns a single field by its name.
func (l FieldsList) GetByName(fieldName string) Field {
	for _, field := range l {
		if field.GetName() == fieldName {
			return field
		}
	}
	return nil
}

// RemoveById removes a single field by its id.
//
// This method does nothing if field with the specified id doesn't exist.
func (l *FieldsList) RemoveById(fieldId string) {
	fields := *l
	for i, field := range fields {
		if field.GetId() == fieldId {
			*l = append(fields[:i], fields[i+1:]...)
			return
		}
	}
}

// RemoveByName removes a single field by its name.
//
// This method does nothing if field with the specified name doesn't exist.
func (l *FieldsList) RemoveByName(fieldName string) {
	fields := *l
	for i, field := range fields {
		if field.GetName() == fieldName {
			*l = append(fields[:i], fields[i+1:]...)
			return
		}
	}
}

// Add adds one or more fields to the current list.
//
// By default this method will try to REPLACE existing fields with
// the new ones by their id or by their name if the new field doesn't have an explicit id.
//
// If no matching existing field is found, it will APPEND the field to the end of the list.
//
// In all cases, if any of the new fields don't have an explicit id it will auto generate a default one for them
// (the id value doesn't really matter and it is mostly used as a stable identifier in case of a field rename).
func (l *FieldsList) Add(fields ...Field) {
	for _, f := range fields {
		l.add(-1, f)
	}
}

// AddAt is the same as Add but insert/move the fields at the specific position.
//
// If pos < 0, then this method acts the same as calling Add.
//
// If pos > FieldsList total items, then the specified fields are inserted/moved at the end of the list.
func (l *FieldsList) AddAt(pos int, fields ...Field) {
	total := len(*l)

	for i, f := range fields {
		if pos < 0 {
			l.add(-1, f)
		} else if pos > total {
			l.add(total+i, f)
		} else {
			l.add(pos+i, f)
		}
	}
}

// AddMarshaledJSON parses the provided raw json data and adds the
// found fields into the current list (following the same rule as the Add method).
//
// The rawJSON argument could be one of:
//   - serialized array of field objects
//   - single field object.
//
// Example:
//
//	l.AddMarshaledJSON([]byte{`{"type":"text", name: "test"}`})
//	l.AddMarshaledJSON([]byte{`[{"type":"text", name: "test1"}, {"type":"text", name: "test2"}]`})
func (l *FieldsList) AddMarshaledJSON(rawJSON []byte) error {
	extractedFields, err := marshaledJSONtoFieldsList(rawJSON)
	if err != nil {
		return err
	}

	l.Add(extractedFields...)

	return nil
}

// AddMarshaledJSONAt is the same as AddMarshaledJSON but insert/move the fields at the specific position.
//
// If pos < 0, then this method acts the same as calling AddMarshaledJSON.
//
// If pos > FieldsList total items, then the specified fields are inserted/moved at the end of the list.
func (l *FieldsList) AddMarshaledJSONAt(pos int, rawJSON []byte) error {
	extractedFields, err := marshaledJSONtoFieldsList(rawJSON)
	if err != nil {
		return err
	}

	l.AddAt(pos, extractedFields...)

	return nil
}

func marshaledJSONtoFieldsList(rawJSON []byte) (FieldsList, error) {
	extractedFields := FieldsList{}

	// nothing to add
	if len(rawJSON) == 0 {
		return extractedFields, nil
	}

	// try to unmarshal first into a new fieds list
	// (assuming that rawJSON is array of objects)
	err := json.Unmarshal(rawJSON, &extractedFields)
	if err != nil {
		// try again but wrap the rawJSON in []
		// (assuming that rawJSON is a single object)
		wrapped := make([]byte, 0, len(rawJSON)+2)
		wrapped = append(wrapped, '[')
		wrapped = append(wrapped, rawJSON...)
		wrapped = append(wrapped, ']')
		err = json.Unmarshal(wrapped, &extractedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal the provided JSON - expects array of objects or just single object: %w", err)
		}
	}

	return extractedFields, nil
}

func (l *FieldsList) add(pos int, newField Field) {
	fields := *l

	var replaceByName bool
	var replaceInPlace bool

	if pos < 0 {
		replaceInPlace = true
		pos = len(fields)
	} else if pos > len(fields) {
		pos = len(fields)
	}

	newFieldId := newField.GetId()

	// set default id
	if newFieldId == "" {
		replaceByName = true

		baseId := newField.Type() + crc32Checksum(newField.GetName())
		newFieldId = baseId
		for i := 2; i < 1000; i++ {
			if l.GetById(newFieldId) == nil {
				break // already unique
			}
			newFieldId = baseId + strconv.Itoa(i)
		}
		newField.SetId(newFieldId)
	}

	// try to replace existing
	for i, field := range fields {
		if replaceByName {
			if name := newField.GetName(); name != "" && field.GetName() == name {
				// reuse the original id
				newField.SetId(field.GetId())

				if replaceInPlace {
					(*l)[i] = newField
					return
				} else {
					// remove the current field and insert it later at the specific position
					*l = slices.Delete(*l, i, i+1)
					if total := len(*l); pos > total {
						pos = total
					}
					break
				}
			}
		} else {
			if field.GetId() == newFieldId {
				if replaceInPlace {
					(*l)[i] = newField
					return
				} else {
					// remove the current field and insert it later at the specific position
					*l = slices.Delete(*l, i, i+1)
					if total := len(*l); pos > total {
						pos = total
					}
					break
				}
			}
		}
	}

	// insert the new field
	*l = slices.Insert(*l, pos, newField)
}

// String returns the string representation of the current list.
func (l FieldsList) String() string {
	v, _ := json.Marshal(l)
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
func (l *FieldsList) UnmarshalJSON(data []byte) error {
	fwts := []fieldWithType{}

	if err := json.Unmarshal(data, &fwts); err != nil {
		return err
	}

	*l = []Field{} // reset

	for _, fwt := range fwts {
		l.add(-1, fwt.Field)
	}

	return nil
}

// MarshalJSON implements the [json.Marshaler] interface.
func (l FieldsList) MarshalJSON() ([]byte, error) {
	if l == nil {
		l = []Field{} // always init to ensure that it is serialized as empty array
	}

	wrapper := make([]map[string]any, 0, len(l))

	for _, f := range l {
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
func (l FieldsList) Value() (driver.Value, error) {
	data, err := json.Marshal(l)

	return string(data), err
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current FieldsList instance.
func (l *FieldsList) Scan(value any) error {
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

	return l.UnmarshalJSON(data)
}
