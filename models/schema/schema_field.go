package schema

import (
	"encoding/json"
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

var schemaFieldNameRegex = regexp.MustCompile(`^\w+$`)

// commonly used field names
const (
	FieldNameId                     = "id"
	FieldNameCreated                = "created"
	FieldNameUpdated                = "updated"
	FieldNameCollectionId           = "collectionId"
	FieldNameCollectionName         = "collectionName"
	FieldNameExpand                 = "expand"
	FieldNameUsername               = "username"
	FieldNameEmail                  = "email"
	FieldNameEmailVisibility        = "emailVisibility"
	FieldNameVerified               = "verified"
	FieldNameTokenKey               = "tokenKey"
	FieldNamePasswordHash           = "passwordHash"
	FieldNameLastResetSentAt        = "lastResetSentAt"
	FieldNameLastVerificationSentAt = "lastVerificationSentAt"
)

// BaseModelFieldNames returns the field names that all models have (id, created, updated).
func BaseModelFieldNames() []string {
	return []string{
		FieldNameId,
		FieldNameCreated,
		FieldNameUpdated,
	}
}

// SystemFields returns special internal field names that are usually readonly.
func SystemFieldNames() []string {
	return []string{
		FieldNameCollectionId,
		FieldNameCollectionName,
		FieldNameExpand,
	}
}

// AuthFieldNames returns the reserved "auth" collection auth field names.
func AuthFieldNames() []string {
	return []string{
		FieldNameUsername,
		FieldNameEmail,
		FieldNameEmailVisibility,
		FieldNameVerified,
		FieldNameTokenKey,
		FieldNamePasswordHash,
		FieldNameLastResetSentAt,
		FieldNameLastVerificationSentAt,
	}
}

// All valid field types
const (
	FieldTypeText     string = "text"
	FieldTypeNumber   string = "number"
	FieldTypeBool     string = "bool"
	FieldTypeEmail    string = "email"
	FieldTypeUrl      string = "url"
	FieldTypeDate     string = "date"
	FieldTypeSelect   string = "select"
	FieldTypeJson     string = "json"
	FieldTypeFile     string = "file"
	FieldTypeRelation string = "relation"

	// Deprecated: Will be removed in v0.9!
	FieldTypeUser string = "user"
)

// FieldTypes returns slice with all supported field types.
func FieldTypes() []string {
	return []string{
		FieldTypeText,
		FieldTypeNumber,
		FieldTypeBool,
		FieldTypeEmail,
		FieldTypeUrl,
		FieldTypeDate,
		FieldTypeSelect,
		FieldTypeJson,
		FieldTypeFile,
		FieldTypeRelation,
	}
}

// ArraybleFieldTypes returns slice with all array value supported field types.
func ArraybleFieldTypes() []string {
	return []string{
		FieldTypeSelect,
		FieldTypeFile,
		FieldTypeRelation,
	}
}

// SchemaField defines a single schema field structure.
type SchemaField struct {
	System   bool   `form:"system" json:"system"`
	Id       string `form:"id" json:"id"`
	Name     string `form:"name" json:"name"`
	Type     string `form:"type" json:"type"`
	Required bool   `form:"required" json:"required"`
	Unique   bool   `form:"unique" json:"unique"`
	Options  any    `form:"options" json:"options"`
}

// ColDefinition returns the field db column type definition as string.
func (f *SchemaField) ColDefinition() string {
	switch f.Type {
	case FieldTypeNumber:
		return "REAL"
	case FieldTypeBool:
		return "BOOLEAN"
	case FieldTypeJson:
		return "JSON"
	default:
		return "TEXT"
	}
}

// String serializes and returns the current field as string.
func (f SchemaField) String() string {
	data, _ := f.MarshalJSON()
	return string(data)
}

// MarshalJSON implements the [json.Marshaler] interface.
func (f SchemaField) MarshalJSON() ([]byte, error) {
	type alias SchemaField // alias to prevent recursion

	f.InitOptions()

	return json.Marshal(alias(f))
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
//
// The schema field options are auto initialized on success.
func (f *SchemaField) UnmarshalJSON(data []byte) error {
	type alias *SchemaField // alias to prevent recursion

	a := alias(f)

	if err := json.Unmarshal(data, a); err != nil {
		return err
	}

	return f.InitOptions()
}

// Validate makes `SchemaField` validatable by implementing [validation.Validatable] interface.
func (f SchemaField) Validate() error {
	// init field options (if not already)
	f.InitOptions()

	excludeNames := BaseModelFieldNames()
	// exclude filter literals
	excludeNames = append(excludeNames, "null", "true", "false")
	// exclude system literals
	excludeNames = append(excludeNames, SystemFieldNames()...)

	return validation.ValidateStruct(&f,
		validation.Field(&f.Options, validation.Required, validation.By(f.checkOptions)),
		validation.Field(&f.Id, validation.Required, validation.Length(5, 255)),
		validation.Field(
			&f.Name,
			validation.Required,
			validation.Length(1, 255),
			validation.Match(schemaFieldNameRegex),
			validation.NotIn(list.ToInterfaceSlice(excludeNames)...),
		),
		validation.Field(&f.Type, validation.Required, validation.In(list.ToInterfaceSlice(FieldTypes())...)),
		// currently file fields cannot be unique because a proper
		// hash/content check could cause performance issues
		validation.Field(&f.Unique, validation.When(f.Type == FieldTypeFile, validation.Empty)),
	)
}

func (f *SchemaField) checkOptions(value any) error {
	v, ok := value.(FieldOptions)
	if !ok {
		return validation.NewError("validation_invalid_options", "Failed to initialize field options")
	}

	return v.Validate()
}

// InitOptions initializes the current field options based on its type.
//
// Returns error on unknown field type.
func (f *SchemaField) InitOptions() error {
	if _, ok := f.Options.(FieldOptions); ok {
		return nil // already inited
	}

	serialized, err := json.Marshal(f.Options)
	if err != nil {
		return err
	}

	var options any
	switch f.Type {
	case FieldTypeText:
		options = &TextOptions{}
	case FieldTypeNumber:
		options = &NumberOptions{}
	case FieldTypeBool:
		options = &BoolOptions{}
	case FieldTypeEmail:
		options = &EmailOptions{}
	case FieldTypeUrl:
		options = &UrlOptions{}
	case FieldTypeDate:
		options = &DateOptions{}
	case FieldTypeSelect:
		options = &SelectOptions{}
	case FieldTypeJson:
		options = &JsonOptions{}
	case FieldTypeFile:
		options = &FileOptions{}
	case FieldTypeRelation:
		options = &RelationOptions{}

	// Deprecated: Will be removed in v0.9!
	case FieldTypeUser:
		options = &UserOptions{}

	default:
		return errors.New("Missing or unknown field field type.")
	}

	if err := json.Unmarshal(serialized, options); err != nil {
		return err
	}

	f.Options = options

	return nil
}

// PrepareValue returns normalized and properly formatted field value.
func (f *SchemaField) PrepareValue(value any) any {
	// init field options (if not already)
	f.InitOptions()

	switch f.Type {
	case FieldTypeText, FieldTypeEmail, FieldTypeUrl:
		return value
	case FieldTypeJson:
		val, _ := types.ParseJsonRaw(value)
		return val
	case FieldTypeNumber:
		if value != nil {
			return cast.ToFloat64(value)
		}
	case FieldTypeBool:
		if value != nil {
			return cast.ToBool(value)
		}
	case FieldTypeDate:
		val, _ := types.ParseDateTime(value)
		return val
	case FieldTypeSelect:
		val := list.ToUniqueStringSlice(value)

		options, _ := f.Options.(*SelectOptions)
		if options.MaxSelect <= 1 {
			if len(val) > 0 {
				return val[0]
			}
			return ""
		}

		return val
	case FieldTypeFile:
		val := list.ToUniqueStringSlice(value)

		options, _ := f.Options.(*FileOptions)
		if options.MaxSelect <= 1 {
			if len(val) > 0 {
				return val[0]
			}
			return ""
		}

		return val
	case FieldTypeRelation:
		ids := list.ToUniqueStringSlice(value)

		options, _ := f.Options.(*RelationOptions)
		if options.MaxSelect != nil && *options.MaxSelect <= 1 {
			if len(ids) > 0 {
				return ids[0]
			}
			return ""
		}

		return ids
	}
	return value
}

// -------------------------------------------------------------------

// FieldOptions interfaces that defines common methods that every field options struct has.
type FieldOptions interface {
	Validate() error
}

type TextOptions struct {
	Min     *int   `form:"min" json:"min"`
	Max     *int   `form:"max" json:"max"`
	Pattern string `form:"pattern" json:"pattern"`
}

func (o TextOptions) Validate() error {
	minVal := 0
	if o.Min != nil {
		minVal = *o.Min
	}

	return validation.ValidateStruct(&o,
		validation.Field(&o.Min, validation.Min(0)),
		validation.Field(&o.Max, validation.Min(minVal)),
		validation.Field(&o.Pattern, validation.By(o.checkRegex)),
	)
}

func (o *TextOptions) checkRegex(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	if _, err := regexp.Compile(v); err != nil {
		return validation.NewError("validation_invalid_regex", err.Error())
	}

	return nil
}

// -------------------------------------------------------------------

type NumberOptions struct {
	Min *float64 `form:"min" json:"min"`
	Max *float64 `form:"max" json:"max"`
}

func (o NumberOptions) Validate() error {
	var maxRules []validation.Rule
	if o.Min != nil && o.Max != nil {
		maxRules = append(maxRules, validation.Min(*o.Min))
	}

	return validation.ValidateStruct(&o,
		validation.Field(&o.Max, maxRules...),
	)
}

// -------------------------------------------------------------------

type BoolOptions struct{}

func (o BoolOptions) Validate() error {
	return nil
}

// -------------------------------------------------------------------

type EmailOptions struct {
	ExceptDomains []string `form:"exceptDomains" json:"exceptDomains"`
	OnlyDomains   []string `form:"onlyDomains" json:"onlyDomains"`
}

func (o EmailOptions) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(
			&o.ExceptDomains,
			validation.When(len(o.OnlyDomains) > 0, validation.Empty).Else(validation.Each(is.Domain)),
		),
		validation.Field(
			&o.OnlyDomains,
			validation.When(len(o.ExceptDomains) > 0, validation.Empty).Else(validation.Each(is.Domain)),
		),
	)
}

// -------------------------------------------------------------------

type UrlOptions struct {
	ExceptDomains []string `form:"exceptDomains" json:"exceptDomains"`
	OnlyDomains   []string `form:"onlyDomains" json:"onlyDomains"`
}

func (o UrlOptions) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(
			&o.ExceptDomains,
			validation.When(len(o.OnlyDomains) > 0, validation.Empty).Else(validation.Each(is.Domain)),
		),
		validation.Field(
			&o.OnlyDomains,
			validation.When(len(o.ExceptDomains) > 0, validation.Empty).Else(validation.Each(is.Domain)),
		),
	)
}

// -------------------------------------------------------------------

type DateOptions struct {
	Min types.DateTime `form:"min" json:"min"`
	Max types.DateTime `form:"max" json:"max"`
}

func (o DateOptions) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Max, validation.By(o.checkRange(o.Min, o.Max))),
	)
}

func (o *DateOptions) checkRange(min types.DateTime, max types.DateTime) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(types.DateTime)

		if v.IsZero() || min.IsZero() || max.IsZero() {
			return nil // nothing to check
		}

		return validation.Date(types.DefaultDateLayout).
			Min(min.Time()).
			Max(max.Time()).
			Validate(v.String())
	}
}

// -------------------------------------------------------------------

type SelectOptions struct {
	MaxSelect int      `form:"maxSelect" json:"maxSelect"`
	Values    []string `form:"values" json:"values"`
}

func (o SelectOptions) Validate() error {
	max := len(o.Values)
	if max == 0 {
		max = 1
	}

	return validation.ValidateStruct(&o,
		validation.Field(&o.Values, validation.Required),
		validation.Field(
			&o.MaxSelect,
			validation.Required,
			validation.Min(1),
			validation.Max(max),
		),
	)
}

// -------------------------------------------------------------------

type JsonOptions struct{}

func (o JsonOptions) Validate() error {
	return nil
}

// -------------------------------------------------------------------

type FileOptions struct {
	MaxSelect int      `form:"maxSelect" json:"maxSelect"`
	MaxSize   int      `form:"maxSize" json:"maxSize"` // in bytes
	MimeTypes []string `form:"mimeTypes" json:"mimeTypes"`
	Thumbs    []string `form:"thumbs" json:"thumbs"`
}

func (o FileOptions) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.MaxSelect, validation.Required, validation.Min(1)),
		validation.Field(&o.MaxSize, validation.Required, validation.Min(1)),
		validation.Field(&o.Thumbs, validation.Each(
			validation.NotIn("0x0", "0x0t", "0x0b", "0x0f"),
			validation.Match(filesystem.ThumbSizeRegex),
		)),
	)
}

// -------------------------------------------------------------------

type RelationOptions struct {
	MaxSelect     *int   `form:"maxSelect" json:"maxSelect"`
	CollectionId  string `form:"collectionId" json:"collectionId"`
	CascadeDelete bool   `form:"cascadeDelete" json:"cascadeDelete"`
}

func (o RelationOptions) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.CollectionId, validation.Required),
		validation.Field(&o.MaxSelect, validation.NilOrNotEmpty, validation.Min(1)),
	)
}

// -------------------------------------------------------------------

// Deprecated: Will be removed in v0.9!
type UserOptions struct {
	MaxSelect     int  `form:"maxSelect" json:"maxSelect"`
	CascadeDelete bool `form:"cascadeDelete" json:"cascadeDelete"`
}

// Deprecated: Will be removed in v0.9!
func (o UserOptions) Validate() error {
	return nil
}
