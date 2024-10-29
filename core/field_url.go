package core

import (
	"context"
	"net/url"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/spf13/cast"
)

func init() {
	Fields[FieldTypeURL] = func() Field {
		return &URLField{}
	}
}

const FieldTypeURL = "url"

var _ Field = (*URLField)(nil)

// URLField defines "url" type field for storing a single URL string value.
//
// The respective zero record field value is empty string.
type URLField struct {
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

	// ExceptDomains will require the URL domain to NOT be included in the listed ones.
	//
	// This validator can be set only if OnlyDomains is empty.
	ExceptDomains []string `form:"exceptDomains" json:"exceptDomains"`

	// OnlyDomains will require the URL domain to be included in the listed ones.
	//
	// This validator can be set only if ExceptDomains is empty.
	OnlyDomains []string `form:"onlyDomains" json:"onlyDomains"`

	// Required will require the field value to be non-empty URL string.
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *URLField) Type() string {
	return FieldTypeURL
}

// GetId implements [Field.GetId] interface method.
func (f *URLField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *URLField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *URLField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *URLField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *URLField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *URLField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *URLField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *URLField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *URLField) ColumnType(app App) string {
	return "TEXT DEFAULT '' NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *URLField) PrepareValue(record *Record, raw any) (any, error) {
	return cast.ToString(raw), nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *URLField) ValidateValue(ctx context.Context, app App, record *Record) error {
	val, ok := record.GetRaw(f.Name).(string)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if f.Required {
		if err := validation.Required.Validate(val); err != nil {
			return err
		}
	}

	if val == "" {
		return nil // nothing to check
	}

	if is.URL.Validate(val) != nil {
		return validation.NewError("validation_invalid_url", "Must be a valid url")
	}

	// extract host/domain
	u, _ := url.Parse(val)

	// only domains check
	if len(f.OnlyDomains) > 0 && !slices.Contains(f.OnlyDomains, u.Host) {
		return validation.NewError("validation_url_domain_not_allowed", "Url domain is not allowed")
	}

	// except domains check
	if len(f.ExceptDomains) > 0 && slices.Contains(f.ExceptDomains, u.Host) {
		return validation.NewError("validation_url_domain_not_allowed", "Url domain is not allowed")
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *URLField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
		validation.Field(
			&f.ExceptDomains,
			validation.When(len(f.OnlyDomains) > 0, validation.Empty).Else(validation.Each(is.Domain)),
		),
		validation.Field(
			&f.OnlyDomains,
			validation.When(len(f.ExceptDomains) > 0, validation.Empty).Else(validation.Each(is.Domain)),
		),
	)
}
