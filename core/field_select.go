package core

import (
	"context"
	"database/sql/driver"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	Fields[FieldTypeSelect] = func() Field {
		return &SelectField{}
	}
}

const FieldTypeSelect = "select"

var (
	_ Field        = (*SelectField)(nil)
	_ MultiValuer  = (*SelectField)(nil)
	_ DriverValuer = (*SelectField)(nil)
	_ SetterFinder = (*SelectField)(nil)
)

// SelectField defines "select" type field for storing single or
// multiple string values from a predefined list.
//
// Requires the Values option to be set.
//
// If MaxSelect is not set or <= 1, then the field value is expected to be a single Values element.
//
// If MaxSelect is > 1, then the field value is expected to be a subset of Values slice.
//
// The respective zero record field value is either empty string (single) or empty string slice (multiple).
//
// ---
//
// The following additional setter keys are available:
//
//   - "fieldName+" - append one or more values to the existing record one. For example:
//
//     record.Set("roles+", []string{"new1", "new2"}) // []string{"old1", "old2", "new1", "new2"}
//
//   - "+fieldName" - prepend one or more values to the existing record one. For example:
//
//     record.Set("+roles", []string{"new1", "new2"}) // []string{"new1", "new2", "old1", "old2"}
//
//   - "fieldName-" - subtract one or more values from the existing record one. For example:
//
//     record.Set("roles-", "old1") // []string{"old2"}
type SelectField struct {
	// Name (required) is the unique name of the field.
	Name string `form:"name" json:"name"`

	// Id is the unique stable field identifier.
	//
	// It is automatically generated from the name when adding to a collection FieldsList.
	Id string `form:"id" json:"id"`

	// System prevents the renaming and removal of the field.
	System bool `form:"system" json:"system"`

	// Hidden hides the field from the API response.
	Hidden bool `form:"hidden" json:"hidden"`

	// Presentable hints the Dashboard UI to use the underlying
	// field record value in the relation preview label.
	Presentable bool `form:"presentable" json:"presentable"`

	// ---

	// Values specifies the list of accepted values.
	Values []string `form:"values" json:"values"`

	// MaxSelect specifies the max allowed selected values.
	//
	// For multiple select the value must be > 1, otherwise fallbacks to single (default).
	MaxSelect int `form:"maxSelect" json:"maxSelect"`

	// Required will require the field value to be non-empty.
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *SelectField) Type() string {
	return FieldTypeSelect
}

// GetId implements [Field.GetId] interface method.
func (f *SelectField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *SelectField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *SelectField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *SelectField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *SelectField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *SelectField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *SelectField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *SelectField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// IsMultiple implements [MultiValuer] interface and checks whether the
// current field options support multiple values.
func (f *SelectField) IsMultiple() bool {
	return f.MaxSelect > 1
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *SelectField) ColumnType(app App) string {
	if f.IsMultiple() {
		return "JSON DEFAULT '[]' NOT NULL"
	}

	return "TEXT DEFAULT '' NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *SelectField) PrepareValue(record *Record, raw any) (any, error) {
	return f.normalizeValue(raw), nil
}

func (f *SelectField) normalizeValue(raw any) any {
	val := list.ToUniqueStringSlice(raw)

	if !f.IsMultiple() {
		if len(val) > 0 {
			return val[len(val)-1] // the last selected
		}
		return ""
	}

	return val
}

// DriverValue implements the [DriverValuer] interface.
func (f *SelectField) DriverValue(record *Record) (driver.Value, error) {
	val := list.ToUniqueStringSlice(record.GetRaw(f.Name))

	if !f.IsMultiple() {
		if len(val) > 0 {
			return val[len(val)-1], nil // the last selected
		}
		return "", nil
	}

	// serialize as json string array
	return append(types.JSONArray[string]{}, val...), nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *SelectField) ValidateValue(ctx context.Context, app App, record *Record) error {
	normalizedVal := list.ToUniqueStringSlice(record.GetRaw(f.Name))
	if len(normalizedVal) == 0 {
		if f.Required {
			return validation.ErrRequired
		}
		return nil // nothing to check
	}

	maxSelect := max(f.MaxSelect, 1)

	// check max selected items
	if len(normalizedVal) > maxSelect {
		return validation.NewError("validation_too_many_values", "Select no more than {{.maxSelect}}").
			SetParams(map[string]any{"maxSelect": maxSelect})
	}

	// check against the allowed values
	for _, val := range normalizedVal {
		if !slices.Contains(f.Values, val) {
			return validation.NewError("validation_invalid_value", "Invalid value {{.value}}").
				SetParams(map[string]any{"value": val})
		}
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *SelectField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	max := len(f.Values)
	if max == 0 {
		max = 1
	}

	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
		validation.Field(&f.Values, validation.Required),
		validation.Field(&f.MaxSelect, validation.Min(0), validation.Max(max)),
	)
}

// FindSetter implements the [SetterFinder] interface.
func (f *SelectField) FindSetter(key string) SetterFunc {
	switch key {
	case f.Name:
		return f.setValue
	case "+" + f.Name:
		return f.prependValue
	case f.Name + "+":
		return f.appendValue
	case f.Name + "-":
		return f.subtractValue
	default:
		return nil
	}
}

func (f *SelectField) setValue(record *Record, raw any) {
	record.SetRaw(f.Name, f.normalizeValue(raw))
}

func (f *SelectField) appendValue(record *Record, modifierValue any) {
	val := record.GetRaw(f.Name)

	val = append(
		list.ToUniqueStringSlice(val),
		list.ToUniqueStringSlice(modifierValue)...,
	)

	f.setValue(record, val)
}

func (f *SelectField) prependValue(record *Record, modifierValue any) {
	val := record.GetRaw(f.Name)

	val = append(
		list.ToUniqueStringSlice(modifierValue),
		list.ToUniqueStringSlice(val)...,
	)

	f.setValue(record, val)
}

func (f *SelectField) subtractValue(record *Record, modifierValue any) {
	val := record.GetRaw(f.Name)

	val = list.SubtractSlice(
		list.ToUniqueStringSlice(val),
		list.ToUniqueStringSlice(modifierValue),
	)

	f.setValue(record, val)
}
