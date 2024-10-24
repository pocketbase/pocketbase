package core

import (
	"context"
	"fmt"
	"math"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/spf13/cast"
)

func init() {
	Fields[FieldTypeNumber] = func() Field {
		return &NumberField{}
	}
}

const FieldTypeNumber = "number"

var (
	_ Field        = (*NumberField)(nil)
	_ SetterFinder = (*NumberField)(nil)
)

// NumberField defines "number" type field for storing numeric (float64) value.
//
// The respective zero record field value is 0.
//
// The following additional setter keys are available:
//
//   - "fieldName+" - appends to the existing record value. For example:
//     record.Set("total+", 5)
//   - "fieldName-" - subtracts from the existing record value. For example:
//     record.Set("total-", 5)
type NumberField struct {
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
	// Leave it nil to skip the validator.
	Min *float64 `form:"min" json:"min"`

	// Max specifies the max allowed field value.
	//
	// Leave it nil to skip the validator.
	Max *float64 `form:"max" json:"max"`

	// OnlyInt will require the field value to be integer.
	OnlyInt bool `form:"onlyInt" json:"onlyInt"`

	// Required will require the field value to be non-zero.
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *NumberField) Type() string {
	return FieldTypeNumber
}

// GetId implements [Field.GetId] interface method.
func (f *NumberField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *NumberField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *NumberField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *NumberField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *NumberField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *NumberField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *NumberField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *NumberField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *NumberField) ColumnType(app App) string {
	return "NUMERIC DEFAULT 0 NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *NumberField) PrepareValue(record *Record, raw any) (any, error) {
	return cast.ToFloat64(raw), nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *NumberField) ValidateValue(ctx context.Context, app App, record *Record) error {
	val, ok := record.GetRaw(f.Name).(float64)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if math.IsInf(val, 0) || math.IsNaN(val) {
		return validation.NewError("validation_not_a_number", "The submitted number is not properly formatted")
	}

	if val == 0 {
		if f.Required {
			if err := validation.Required.Validate(val); err != nil {
				return err
			}
		}
		return nil
	}

	if f.OnlyInt && val != float64(int64(val)) {
		return validation.NewError("validation_only_int_constraint", "Decimal numbers are not allowed")
	}

	if f.Min != nil && val < *f.Min {
		return validation.NewError("validation_min_number_constraint", fmt.Sprintf("Must be larger than %f", *f.Min))
	}

	if f.Max != nil && val > *f.Max {
		return validation.NewError("validation_max_number_constraint", fmt.Sprintf("Must be less than %f", *f.Max))
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *NumberField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	maxRules := []validation.Rule{
		validation.By(f.checkOnlyInt),
	}
	if f.Min != nil && f.Max != nil {
		maxRules = append(maxRules, validation.Min(*f.Min))
	}

	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
		validation.Field(&f.Min, validation.By(f.checkOnlyInt)),
		validation.Field(&f.Max, maxRules...),
	)
}

func (f *NumberField) checkOnlyInt(value any) error {
	v, _ := value.(*float64)
	if v == nil || !f.OnlyInt {
		return nil // nothing to check
	}

	if *v != float64(int64(*v)) {
		return validation.NewError("validation_only_int_constraint", "Decimal numbers are not allowed.")
	}

	return nil
}

// FindSetter implements the [SetterFinder] interface.
func (f *NumberField) FindSetter(key string) SetterFunc {
	switch key {
	case f.Name:
		return f.setValue
	case f.Name + "+":
		return f.addValue
	case f.Name + "-":
		return f.subtractValue
	default:
		return nil
	}
}

func (f *NumberField) setValue(record *Record, raw any) {
	record.SetRaw(f.Name, cast.ToFloat64(raw))
}

func (f *NumberField) addValue(record *Record, raw any) {
	val := cast.ToFloat64(record.GetRaw(f.Name))

	record.SetRaw(f.Name, val+cast.ToFloat64(raw))
}

func (f *NumberField) subtractValue(record *Record, raw any) {
	val := cast.ToFloat64(record.GetRaw(f.Name))

	record.SetRaw(f.Name, val-cast.ToFloat64(raw))
}
