package core

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	Fields[FieldTypePoint] = func() Field {
		return &PointField{}
	}
}

const FieldTypePoint = "point"

var _ Field = (*PointField)(nil)

// PointField defines "point" type field to store a single [types.Point] value.
//
// The respective zero record field value is the zero [types.Point].
type PointField struct {
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

	// Required will require the field value to be non-empty coordinate pair.
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *PointField) Type() string {
	return FieldTypeDate
}

// GetId implements [Field.GetId] interface method.
func (f *PointField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *PointField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *PointField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *PointField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *PointField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *PointField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *PointField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *PointField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *PointField) ColumnType(app App) string {
	return "TEXT DEFAULT '' NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *PointField) PrepareValue(record *Record, raw any) (any, error) {
	// ignore scan errors since the format may change between versions
	// and to allow running db adjusting migrations
	val, _ := types.ParseDateTime(raw)
	return val, nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *PointField) ValidateValue(ctx context.Context, app App, record *Record) error {
	val, ok := record.GetRaw(f.Name).(types.Point)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if val.IsEmpty() {
		if f.Required {
			if err := validation.Required.Validate(val); err != nil {
				return err
			}
		}
		return nil
	}

	lat := val.Lat()
	if lat > 90 {
		return fmt.Errorf("latitude cannot be >90, got %f", val.Lat())
	}
	if lat < -90 {
		return fmt.Errorf("latitude cannot be <-90, got %f", val.Lat())
	}

	long := val.Long()
	if long > 180 {
		return fmt.Errorf("longitude cannot be >180, got %f", val.Lat())
	}
	if lat < -180 {
		return fmt.Errorf("longitude cannot be <-180, got %f", val.Lat())
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *PointField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
	)
}
