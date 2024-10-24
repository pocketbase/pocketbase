package core

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	Fields[FieldTypeDate] = func() Field {
		return &DateField{}
	}
}

const FieldTypeDate = "date"

var _ Field = (*DateField)(nil)

// DateField defines "date" type field to store a single [types.DateTime] value.
//
// The respective zero record field value is the zero [types.DateTime].
type DateField struct {
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

	// Min specifies the min allowed field value.
	//
	// Leave it empty to skip the validator.
	Min types.DateTime `form:"min" json:"min"`

	// Max specifies the max allowed field value.
	//
	// Leave it empty to skip the validator.
	Max types.DateTime `form:"max" json:"max"`

	// Required will require the field value to be non-zero [types.DateTime].
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *DateField) Type() string {
	return FieldTypeDate
}

// GetId implements [Field.GetId] interface method.
func (f *DateField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *DateField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *DateField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *DateField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *DateField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *DateField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *DateField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *DateField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *DateField) ColumnType(app App) string {
	return "TEXT DEFAULT '' NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *DateField) PrepareValue(record *Record, raw any) (any, error) {
	// ignore scan errors since the format may change between versions
	// and to allow running db adjusting migrations
	val, _ := types.ParseDateTime(raw)
	return val, nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *DateField) ValidateValue(ctx context.Context, app App, record *Record) error {
	val, ok := record.GetRaw(f.Name).(types.DateTime)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if val.IsZero() {
		if f.Required {
			return validation.ErrRequired
		}
		return nil // nothing to check
	}

	if !f.Min.IsZero() {
		if err := validation.Min(f.Min.Time()).Validate(val.Time()); err != nil {
			return err
		}
	}

	if !f.Max.IsZero() {
		if err := validation.Max(f.Max.Time()).Validate(val.Time()); err != nil {
			return err
		}
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *DateField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
		validation.Field(&f.Max, validation.By(f.checkRange(f.Min, f.Max))),
	)
}

func (f *DateField) checkRange(min types.DateTime, max types.DateTime) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(types.DateTime)
		if v.IsZero() {
			return nil // nothing to check
		}

		dr := validation.Date(types.DefaultDateLayout)

		if !min.IsZero() {
			dr.Min(min.Time())
		}

		if !max.IsZero() {
			dr.Max(max.Time())
		}

		return dr.Validate(v.String())
	}
}
