package schema

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

var schemaFieldNameRegex = regexp.MustCompile(`^\w+$`)

// field value modifiers
const (
	FieldValueModifierAdd      string = "+"
	FieldValueModifierSubtract string = "-"
)

// FieldValueModifiers returns a list with all available field modifier tokens.
func FieldValueModifiers() []string {
	return []string{
		FieldValueModifierAdd,
		FieldValueModifierSubtract,
	}
}

// commonly used field names
const (
	FieldNameId                     string = "id"
	FieldNameCreated                string = "created"
	FieldNameUpdated                string = "updated"
	FieldNameCollectionId           string = "collectionId"
	FieldNameCollectionName         string = "collectionName"
	FieldNameExpand                 string = "expand"
	FieldNameUsername               string = "username"
	FieldNameEmail                  string = "email"
	FieldNameEmailVisibility        string = "emailVisibility"
	FieldNameVerified               string = "verified"
	FieldNameTokenKey               string = "tokenKey"
	FieldNamePasswordHash           string = "passwordHash"
	FieldNameLastResetSentAt        string = "lastResetSentAt"
	FieldNameLastVerificationSentAt string = "lastVerificationSentAt"
	FieldNameLastLoginAlertSentAt   string = "lastLoginAlertSentAt"
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
		FieldNameLastLoginAlertSentAt,
	}
}

// All valid field types
const (
	FieldTypeText     string = "text"
	FieldTypeNumber   string = "number"
	FieldTypeBool     string = "bool"
	FieldTypeEmail    string = "email"
	FieldTypeUrl      string = "url"
	FieldTypeEditor   string = "editor"
	FieldTypeDate     string = "date"
	FieldTypeSelect   string = "select"
	FieldTypeJson     string = "json"
	FieldTypeFile     string = "file"
	FieldTypeRelation string = "relation"

	// Deprecated: Will be removed in v0.9+
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
		FieldTypeEditor,
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

	// Presentable indicates whether the field is suitable for
	// visualization purposes (eg. in the Admin UI relation views).
	Presentable bool `form:"presentable" json:"presentable"`

	// Deprecated: This field is no-op and will be removed in future versions.
	// Please use the collection.Indexes field to define a unique constraint.
	Unique bool `form:"unique" json:"unique"`

	Options any `form:"options" json:"options"`
}

// ColDefinition returns the field db column type definition as string.
func (f *SchemaField) ColDefinition() string {
	switch f.Type {
	case FieldTypeNumber:
		return "NUMERIC DEFAULT 0 NOT NULL"
	case FieldTypeBool:
		return "BOOLEAN DEFAULT FALSE NOT NULL"
	case FieldTypeJson:
		return "JSON DEFAULT NULL"
	default:
		if opt, ok := f.Options.(MultiValuer); ok && opt.IsMultiple() {
			return "JSON DEFAULT '[]' NOT NULL"
		}

		return "TEXT DEFAULT '' NOT NULL"
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
	// exclude special filter literals
	excludeNames = append(excludeNames, "null", "true", "false", "_rowid_")
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
			validation.By(f.checkForVia),
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

// @todo merge with the collections during the refactoring
func (f *SchemaField) checkForVia(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil
	}

	if strings.Contains(strings.ToLower(v), "_via_") {
		return validation.NewError("validation_invalid_name", "The name of the field cannot contain '_via_'.")
	}

	return nil
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
	case FieldTypeEditor:
		options = &EditorOptions{}
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

	// Deprecated: Will be removed in v0.9+
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
	case FieldTypeText, FieldTypeEmail, FieldTypeUrl, FieldTypeEditor:
		return cast.ToString(value)
	case FieldTypeJson:
		val := value

		if str, ok := val.(string); ok {
			// in order to support seamlessly both json and multipart/form-data requests,
			// the following normalization rules are applied for plain string values:
			// - "true" is converted to the json `true`
			// - "false" is converted to the json `false`
			// - "null" is converted to the json `null`
			// - "[1,2,3]" is converted to the json `[1,2,3]`
			// - "{\"a\":1,\"b\":2}" is converted to the json `{"a":1,"b":2}`
			// - numeric strings are converted to json number
			// - double quoted strings are left as they are (aka. without normalizations)
			// - any other string (empty string too) is double quoted
			if str == "" {
				val = strconv.Quote(str)
			} else if str == "null" || str == "true" || str == "false" {
				val = str
			} else if ((str[0] >= '0' && str[0] <= '9') ||
				str[0] == '-' ||
				str[0] == '"' ||
				str[0] == '[' ||
				str[0] == '{') &&
				is.JSON.Validate(str) == nil {
				val = str
			} else {
				val = strconv.Quote(str)
			}
		}

		val, _ = types.ParseJsonRaw(val)
		return val
	case FieldTypeNumber:
		return cast.ToFloat64(value)
	case FieldTypeBool:
		return cast.ToBool(value)
	case FieldTypeDate:
		val, _ := types.ParseDateTime(value)
		return val
	case FieldTypeSelect:
		val := list.ToUniqueStringSlice(value)

		options, _ := f.Options.(*SelectOptions)
		if !options.IsMultiple() {
			if len(val) > 0 {
				return val[len(val)-1] // the last selected
			}
			return ""
		}

		return val
	case FieldTypeFile:
		val := list.ToUniqueStringSlice(value)

		options, _ := f.Options.(*FileOptions)
		if !options.IsMultiple() {
			if len(val) > 0 {
				return val[len(val)-1] // the last selected
			}
			return ""
		}

		return val
	case FieldTypeRelation:
		ids := list.ToUniqueStringSlice(value)

		options, _ := f.Options.(*RelationOptions)
		if !options.IsMultiple() {
			if len(ids) > 0 {
				return ids[len(ids)-1] // the last selected
			}
			return ""
		}

		return ids
	default:
		return value // unmodified
	}
}

// PrepareValueWithModifier returns normalized and properly formatted field value
// by "merging" baseValue with the modifierValue based on the specified modifier (+ or -).
func (f *SchemaField) PrepareValueWithModifier(baseValue any, modifier string, modifierValue any) any {
	resolvedValue := baseValue

	switch f.Type {
	case FieldTypeNumber:
		switch modifier {
		case FieldValueModifierAdd:
			resolvedValue = cast.ToFloat64(baseValue) + cast.ToFloat64(modifierValue)
		case FieldValueModifierSubtract:
			resolvedValue = cast.ToFloat64(baseValue) - cast.ToFloat64(modifierValue)
		}
	case FieldTypeSelect, FieldTypeRelation:
		switch modifier {
		case FieldValueModifierAdd:
			resolvedValue = append(
				list.ToUniqueStringSlice(baseValue),
				list.ToUniqueStringSlice(modifierValue)...,
			)
		case FieldValueModifierSubtract:
			resolvedValue = list.SubtractSlice(
				list.ToUniqueStringSlice(baseValue),
				list.ToUniqueStringSlice(modifierValue),
			)
		}
	case FieldTypeFile:
		// note: file for now supports only the subtract modifier
		if modifier == FieldValueModifierSubtract {
			resolvedValue = list.SubtractSlice(
				list.ToUniqueStringSlice(baseValue),
				list.ToUniqueStringSlice(modifierValue),
			)
		}
	}

	return f.PrepareValue(resolvedValue)
}

// -------------------------------------------------------------------

// MultiValuer defines common interface methods that every multi-valued (eg. with MaxSelect) field option struct has.
type MultiValuer interface {
	IsMultiple() bool
}

// FieldOptions defines common interface methods that every field option struct has.
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
	Min       *float64 `form:"min" json:"min"`
	Max       *float64 `form:"max" json:"max"`
	NoDecimal bool     `form:"noDecimal" json:"noDecimal"`
}

func (o NumberOptions) Validate() error {
	var maxRules []validation.Rule
	if o.Min != nil && o.Max != nil {
		maxRules = append(maxRules, validation.Min(*o.Min), validation.By(o.checkNoDecimal))
	}

	return validation.ValidateStruct(&o,
		validation.Field(&o.Min, validation.By(o.checkNoDecimal)),
		validation.Field(&o.Max, maxRules...),
	)
}

func (o *NumberOptions) checkNoDecimal(value any) error {
	v, _ := value.(*float64)
	if v == nil || !o.NoDecimal {
		return nil // nothing to check
	}

	if *v != float64(int64(*v)) {
		return validation.NewError("validation_no_decimal_constraint", "Decimal numbers are not allowed.")
	}

	return nil
}

// -------------------------------------------------------------------

type BoolOptions struct {
}

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

type EditorOptions struct {
	// ConvertUrls is usually used to instruct the editor whether to
	// apply url conversion (eg. stripping the domain name in case the
	// urls are using the same domain as the one where the editor is loaded).
	//
	// (see also https://www.tiny.cloud/docs/tinymce/6/url-handling/#convert_urls)
	ConvertUrls bool `form:"convertUrls" json:"convertUrls"`
}

func (o EditorOptions) Validate() error {
	return nil
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

// IsMultiple implements MultiValuer interface and checks whether the
// current field options support multiple values.
func (o SelectOptions) IsMultiple() bool {
	return o.MaxSelect > 1
}

// -------------------------------------------------------------------

type JsonOptions struct {
	MaxSize int `form:"maxSize" json:"maxSize"`
}

func (o JsonOptions) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.MaxSize, validation.Required, validation.Min(1)),
	)
}

// -------------------------------------------------------------------

var _ MultiValuer = (*FileOptions)(nil)

type FileOptions struct {
	MimeTypes []string `form:"mimeTypes" json:"mimeTypes"`
	Thumbs    []string `form:"thumbs" json:"thumbs"`
	MaxSelect int      `form:"maxSelect" json:"maxSelect"`
	MaxSize   int      `form:"maxSize" json:"maxSize"`
	Protected bool     `form:"protected" json:"protected"`
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

// IsMultiple implements MultiValuer interface and checks whether the
// current field options support multiple values.
func (o FileOptions) IsMultiple() bool {
	return o.MaxSelect > 1
}

// -------------------------------------------------------------------

var _ MultiValuer = (*RelationOptions)(nil)

type RelationOptions struct {
	// CollectionId is the id of the related collection.
	CollectionId string `form:"collectionId" json:"collectionId"`

	// CascadeDelete indicates whether the root model should be deleted
	// in case of delete of all linked relations.
	CascadeDelete bool `form:"cascadeDelete" json:"cascadeDelete"`

	// MinSelect indicates the min number of allowed relation records
	// that could be linked to the main model.
	//
	// If nil no limits are applied.
	MinSelect *int `form:"minSelect" json:"minSelect"`

	// MaxSelect indicates the max number of allowed relation records
	// that could be linked to the main model.
	//
	// If nil no limits are applied.
	MaxSelect *int `form:"maxSelect" json:"maxSelect"`

	// Deprecated: This field is no-op and will be removed in future versions.
	// Instead use the individula SchemaField.Presentable option for each field in the relation collection.
	DisplayFields []string `form:"displayFields" json:"displayFields"`
}

func (o RelationOptions) Validate() error {
	minVal := 0
	if o.MinSelect != nil {
		minVal = *o.MinSelect
	}

	return validation.ValidateStruct(&o,
		validation.Field(&o.CollectionId, validation.Required),
		validation.Field(&o.MinSelect, validation.Min(0)),
		validation.Field(&o.MaxSelect, validation.NilOrNotEmpty, validation.Min(minVal)),
	)
}

// IsMultiple implements MultiValuer interface and checks whether the
// current field options support multiple values.
func (o RelationOptions) IsMultiple() bool {
	return o.MaxSelect == nil || *o.MaxSelect > 1
}

// -------------------------------------------------------------------

// Deprecated: Will be removed in v0.9+
type UserOptions struct {
	MaxSelect     int  `form:"maxSelect" json:"maxSelect"`
	CascadeDelete bool `form:"cascadeDelete" json:"cascadeDelete"`
}

// Deprecated: Will be removed in v0.9+
func (o UserOptions) Validate() error {
	return nil
}
