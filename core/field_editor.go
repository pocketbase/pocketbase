package core

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/spf13/cast"
)

func init() {
	Fields[FieldTypeEditor] = func() Field {
		return &EditorField{}
	}
}

const FieldTypeEditor = "editor"

const DefaultEditorFieldMaxSize int64 = 5 << 20

var (
	_ Field                 = (*EditorField)(nil)
	_ MaxBodySizeCalculator = (*EditorField)(nil)
)

// EditorField defines "editor" type field to store HTML formatted text.
//
// The respective zero record field value is empty string.
type EditorField struct {
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

	// MaxSize specifies the maximum size of the allowed field value (in bytes and up to 2^53-1).
	//
	// If zero, a default limit of ~5MB is applied.
	MaxSize int64 `form:"maxSize" json:"maxSize"`

	// ConvertURLs is usually used to instruct the editor whether to
	// apply url conversion (eg. stripping the domain name in case the
	// urls are using the same domain as the one where the editor is loaded).
	//
	// (see also https://www.tiny.cloud/docs/tinymce/6/url-handling/#convert_urls)
	ConvertURLs bool `form:"convertURLs" json:"convertURLs"`

	// Required will require the field value to be non-empty string.
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *EditorField) Type() string {
	return FieldTypeEditor
}

// GetId implements [Field.GetId] interface method.
func (f *EditorField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *EditorField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *EditorField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *EditorField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *EditorField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *EditorField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *EditorField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *EditorField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *EditorField) ColumnType(app App) string {
	return "TEXT DEFAULT '' NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *EditorField) PrepareValue(record *Record, raw any) (any, error) {
	return cast.ToString(raw), nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *EditorField) ValidateValue(ctx context.Context, app App, record *Record) error {
	val, ok := record.GetRaw(f.Name).(string)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if f.Required {
		if err := validation.Required.Validate(val); err != nil {
			return err
		}
	}

	maxSize := f.CalculateMaxBodySize()

	if int64(len(val)) > maxSize {
		return validation.NewError(
			"validation_content_size_limit",
			"The maximum allowed content size is {{.maxSize}} bytes",
		).SetParams(map[string]any{"maxSize": maxSize})
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *EditorField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
		validation.Field(&f.MaxSize, validation.Min(0), validation.Max(maxSafeJSONInt)),
	)
}

// CalculateMaxBodySize implements the [MaxBodySizeCalculator] interface.
func (f *EditorField) CalculateMaxBodySize() int64 {
	if f.MaxSize <= 0 {
		return DefaultEditorFieldMaxSize
	}

	return f.MaxSize
}
