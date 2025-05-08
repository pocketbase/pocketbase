package core

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	Fields[FieldTypeGeoPoint] = func() Field {
		return &GeoPointField{}
	}
}

const FieldTypeGeoPoint = "geoPoint"

var (
	_ Field = (*GeoPointField)(nil)
)

// GeoPointField defines "geoPoint" type field for storing latitude and longitude GPS coordinates.
//
// You can set the record field value as [types.GeoPoint], map or serialized json object with lat-lon props.
// The stored value is always converted to [types.GeoPoint].
// Nil, empty map, empty bytes slice, etc. results in zero [types.GeoPoint].
//
// Examples of updating a record's GeoPointField value programmatically:
//
//	record.Set("location", types.GeoPoint{Lat: 123, Lon: 456})
//	record.Set("location", map[string]any{"lat":123, "lon":456})
//	record.Set("location", []byte(`{"lat":123, "lon":456}`)
type GeoPointField struct {
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

	// Required will require the field coordinates to be non-zero (aka. not "Null Island").
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *GeoPointField) Type() string {
	return FieldTypeGeoPoint
}

// GetId implements [Field.GetId] interface method.
func (f *GeoPointField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *GeoPointField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *GeoPointField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *GeoPointField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *GeoPointField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *GeoPointField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *GeoPointField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *GeoPointField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *GeoPointField) ColumnType(app App) string {
	return `JSONB DEFAULT '{"lon":0,"lat":0}' NOT NULL`
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *GeoPointField) PrepareValue(record *Record, raw any) (any, error) {
	point := types.GeoPoint{}
	err := point.Scan(raw)
	return point, err
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *GeoPointField) ValidateValue(ctx context.Context, app App, record *Record) error {
	val, ok := record.GetRaw(f.Name).(types.GeoPoint)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	// zero value
	if val.Lat == 0 && val.Lon == 0 {
		if f.Required {
			return validation.ErrRequired
		}
		return nil
	}

	if val.Lat < -90 || val.Lat > 90 {
		return validation.NewError("validation_invalid_latitude", "Latitude must be between -90 and 90 degrees.")
	}

	if val.Lon < -180 || val.Lon > 180 {
		return validation.NewError("validation_invalid_longitude", "Longitude must be between -180 and 180 degrees.")
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *GeoPointField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
	)
}
