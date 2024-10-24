package core

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/spf13/cast"
)

func init() {
	Fields[FieldTypeBool] = func() Field {
		return &BoolField{}
	}
}

const FieldTypeBool = "bool"

var _ Field = (*BoolField)(nil)

// BoolField defines "bool" type field to store a single true/false value.
//
// The respective zero record field value is false.
type BoolField struct {
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

	// Required will require the field value to be always "true".
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *BoolField) Type() string {
	return FieldTypeBool
}

// GetId implements [Field.GetId] interface method.
func (f *BoolField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *BoolField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *BoolField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *BoolField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *BoolField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *BoolField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *BoolField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *BoolField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *BoolField) ColumnType(app App) string {
	return "BOOLEAN DEFAULT FALSE NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *BoolField) PrepareValue(record *Record, raw any) (any, error) {
	return cast.ToBool(raw), nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *BoolField) ValidateValue(ctx context.Context, app App, record *Record) error {
	v, ok := record.GetRaw(f.Name).(bool)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if f.Required {
		return validation.Required.Validate(v)
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *BoolField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
	)
}
