// Package schema implements custom Schema and SchemaField datatypes
// for handling the Collection schema definitions.
package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
)

// NewSchema creates a new Schema instance with the provided fields.
func NewSchema(fields ...*SchemaField) Schema {
	s := Schema{}

	for _, f := range fields {
		s.AddField(f)
	}

	return s
}

// Schema defines a dynamic db schema as a slice of `SchemaField`s.
type Schema struct {
	fields []*SchemaField
}

// Fields returns the registered schema fields.
func (s *Schema) Fields() []*SchemaField {
	return s.fields
}

// InitFieldsOptions calls `InitOptions()` for all schema fields.
func (s *Schema) InitFieldsOptions() error {
	for _, field := range s.Fields() {
		if err := field.InitOptions(); err != nil {
			return err
		}
	}
	return nil
}

// Clone creates a deep clone of the current schema.
func (s *Schema) Clone() (*Schema, error) {
	copyRaw, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	result := &Schema{}
	if err := json.Unmarshal(copyRaw, result); err != nil {
		return nil, err
	}

	return result, nil
}

// AsMap returns a map with all registered schema field.
// The returned map is indexed with each field name.
func (s *Schema) AsMap() map[string]*SchemaField {
	result := map[string]*SchemaField{}

	for _, field := range s.fields {
		result[field.Name] = field
	}

	return result
}

// GetFieldById returns a single field by its id.
func (s *Schema) GetFieldById(id string) *SchemaField {
	for _, field := range s.fields {
		if field.Id == id {
			return field
		}
	}
	return nil
}

// GetFieldByName returns a single field by its name.
func (s *Schema) GetFieldByName(name string) *SchemaField {
	for _, field := range s.fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

// RemoveField removes a single schema field by its id.
//
// This method does nothing if field with `id` doesn't exist.
func (s *Schema) RemoveField(id string) {
	for i, field := range s.fields {
		if field.Id == id {
			s.fields = append(s.fields[:i], s.fields[i+1:]...)
			return
		}
	}
}

// AddField registers the provided newField to the current schema.
//
// If field with `newField.Id` already exist, the existing field is
// replaced with the new one.
//
// Otherwise the new field is appended to the other schema fields.
func (s *Schema) AddField(newField *SchemaField) {
	if newField.Id == "" {
		// set default id
		newField.Id = strings.ToLower(security.PseudorandomString(8))
	}

	for i, field := range s.fields {
		// replace existing
		if field.Id == newField.Id {
			s.fields[i] = newField
			return
		}
	}

	// add new field
	s.fields = append(s.fields, newField)
}

// Validate makes Schema validatable by implementing [validation.Validatable] interface.
//
// Internally calls each individual field's validator and additionally
// checks for invalid renamed fields and field name duplications.
func (s Schema) Validate() error {
	return validation.Validate(&s.fields, validation.By(func(value any) error {
		fields := s.fields // use directly the schema value to avoid unnecessary interface casting

		ids := []string{}
		names := []string{}
		for i, field := range fields {
			if list.ExistInSlice(field.Id, ids) {
				return validation.Errors{
					strconv.Itoa(i): validation.Errors{
						"id": validation.NewError(
							"validation_duplicated_field_id",
							"Duplicated or invalid schema field id",
						),
					},
				}
			}

			// field names are used as db columns and should be case insensitive
			nameLower := strings.ToLower(field.Name)

			if list.ExistInSlice(nameLower, names) {
				return validation.Errors{
					strconv.Itoa(i): validation.Errors{
						"name": validation.NewError(
							"validation_duplicated_field_name",
							"Duplicated or invalid schema field name",
						),
					},
				}
			}

			ids = append(ids, field.Id)
			names = append(names, nameLower)
		}

		return nil
	}))
}

// MarshalJSON implements the [json.Marshaler] interface.
func (s Schema) MarshalJSON() ([]byte, error) {
	if s.fields == nil {
		s.fields = []*SchemaField{}
	}
	return json.Marshal(s.fields)
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
//
// On success, all schema field options are auto initialized.
func (s *Schema) UnmarshalJSON(data []byte) error {
	fields := []*SchemaField{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	s.fields = []*SchemaField{}

	for _, f := range fields {
		s.AddField(f)
	}

	for _, field := range s.fields {
		if err := field.InitOptions(); err != nil {
			// ignore the error and remove the invalid field
			s.RemoveField(field.Id)
		}
	}

	return nil
}

// Value implements the [driver.Valuer] interface.
func (s Schema) Value() (driver.Value, error) {
	if s.fields == nil {
		// initialize an empty slice to ensure that `[]` is returned
		s.fields = []*SchemaField{}
	}

	data, err := json.Marshal(s.fields)

	return string(data), err
}

// Scan implements [sql.Scanner] interface to scan the provided value
// into the current Schema instance.
func (s *Schema) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case nil:
		// no cast needed
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("Failed to unmarshal Schema value %q.", value)
	}

	if len(data) == 0 {
		data = []byte("[]")
	}

	return s.UnmarshalJSON(data)
}
