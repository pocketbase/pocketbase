package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

func init() {
	Fields[FieldTypeText] = func() Field {
		return &TextField{}
	}
}

const FieldTypeText = "text"

const autogenerateModifier = ":autogenerate"

var (
	_ Field             = (*TextField)(nil)
	_ SetterFinder      = (*TextField)(nil)
	_ RecordInterceptor = (*TextField)(nil)
)

var forbiddenPKCharacters = []string{
	".", "/", `\`, "|", `"`, "'", "`",
	"<", ">", ":", "?", "*", "%", "$",
	"\000", "\t", "\n", "\r", " ",
}

// (see largestReservedPKLength)
var caseInsensitiveReservedPKs = []string{
	// reserved Windows files names
	"CON", "PRN", "AUX", "NUL",
	"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
	"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
}

const largestReservedPKLength = 4

// TextField defines "text" type field for storing any string value.
//
// The respective zero record field value is empty string.
//
// The following additional setter keys are available:
//
// - "fieldName:autogenerate" - autogenerate field value if AutogeneratePattern is set. For example:
//
//	record.Set("slug:autogenerate", "") // [random value]
//	record.Set("slug:autogenerate", "abc-") // abc-[random value]
type TextField struct {
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

	// Min specifies the minimum required string characters.
	//
	// if zero value, no min limit is applied.
	Min int `form:"min" json:"min"`

	// Max specifies the maximum allowed string characters.
	//
	// If zero, a default limit of 5000 is applied.
	Max int `form:"max" json:"max"`

	// Pattern specifies an optional regex pattern to match against the field value.
	//
	// Leave it empty to skip the pattern check.
	Pattern string `form:"pattern" json:"pattern"`

	// AutogeneratePattern specifies an optional regex pattern that could
	// be used to generate random string from it and set it automatically
	// on record create if no explicit value is set or when the `:autogenerate` modifier is used.
	//
	// Note: the generated value still needs to satisfy min, max, pattern (if set)
	AutogeneratePattern string `form:"autogeneratePattern" json:"autogeneratePattern"`

	// Required will require the field value to be non-empty string.
	Required bool `form:"required" json:"required"`

	// PrimaryKey will mark the field as primary key.
	//
	// A single collection can have only 1 field marked as primary key.
	PrimaryKey bool `form:"primaryKey" json:"primaryKey"`
}

// Type implements [Field.Type] interface method.
func (f *TextField) Type() string {
	return FieldTypeText
}

// GetId implements [Field.GetId] interface method.
func (f *TextField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *TextField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *TextField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *TextField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *TextField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *TextField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *TextField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *TextField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *TextField) ColumnType(app App) string {
	if f.PrimaryKey {
		// note: the default is just a last resort fallback to avoid empty
		// string values in case the record was inserted with raw sql and
		// it is not actually used when operating with the db abstraction
		return "TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL"
	}

	return "TEXT DEFAULT '' NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *TextField) PrepareValue(record *Record, raw any) (any, error) {
	return cast.ToString(raw), nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *TextField) ValidateValue(ctx context.Context, app App, record *Record) error {
	newVal, ok := record.GetRaw(f.Name).(string)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if f.PrimaryKey {
		// disallow PK change
		if !record.IsNew() {
			oldVal := record.LastSavedPK()
			if oldVal != newVal {
				return validation.NewError("validation_pk_change", "The record primary key cannot be changed.")
			}
			if oldVal != "" {
				// no need to further validate because the id can't be updated
				// and because the id could have been inserted manually by migration from another system
				// that may not comply with the user defined PocketBase validations
				return nil
			}
		} else {
			// this technically shouldn't be necessarily but again to
			// minimize misuse of the Pattern validator that could cause
			// side-effects on some platforms check for duplicates in a case-insensitive manner
			//
			// (@todo eventually may get replaced in the future with a system unique constraint to avoid races or wrapping the request in a transaction)
			if f.Pattern != defaultLowercaseRecordIdPattern {
				var exists int
				err := app.ConcurrentDB().
					Select("(1)").
					From(record.TableName()).
					Where(dbx.NewExp("id = {:id} COLLATE NOCASE", dbx.Params{"id": newVal})).
					Limit(1).
					Row(&exists)
				if exists > 0 || (err != nil && !errors.Is(err, sql.ErrNoRows)) {
					return validation.NewError("validation_pk_invalid", "The record primary key is invalid or already exists.")
				}
			}
		}
	}

	return f.ValidatePlainValue(newVal)
}

// ValidatePlainValue validates the provided string against the field options.
func (f *TextField) ValidatePlainValue(value string) error {
	if f.Required || f.PrimaryKey {
		if err := validation.Required.Validate(value); err != nil {
			return err
		}
	}

	if value == "" {
		return nil // nothing to check
	}

	// note: casted to []rune to count multi-byte chars as one
	length := len([]rune(value))

	if f.Min > 0 && length < f.Min {
		return validation.NewError("validation_min_text_constraint", "Must be at least {{.min}} character(s).").
			SetParams(map[string]any{"min": f.Min})
	}

	max := f.Max
	if max == 0 {
		max = 5000
	}

	if max > 0 && length > max {
		return validation.NewError("validation_max_text_constraint", "Must be no more than {{.max}} character(s).").
			SetParams(map[string]any{"max": max})
	}

	if f.Pattern != "" {
		match, _ := regexp.MatchString(f.Pattern, value)
		if !match {
			return validation.NewError("validation_invalid_format", "Invalid value format.")
		}
	}

	// additional primary key checks to minimize eventual filesystem compatibility issues
	// because the primary key is often used as a file/directory name
	if f.PrimaryKey && f.Pattern != defaultLowercaseRecordIdPattern {
		for _, ch := range forbiddenPKCharacters {
			if strings.Contains(value, ch) {
				return validation.NewError("validation_forbidden_pk_character", "'{{.ch}}' is not a valid primary key character.").
					SetParams(map[string]any{"ch": ch})
			}
		}

		if largestReservedPKLength >= length {
			for _, reserved := range caseInsensitiveReservedPKs {
				if strings.EqualFold(value, reserved) {
					return validation.NewError("validation_reserved_pk", "The primary key '{{.reserved}}' is reserved and cannot be used.").
						SetParams(map[string]any{"reserved": reserved})
				}
			}
		}
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *TextField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name,
			validation.By(DefaultFieldNameValidationRule),
			validation.When(f.PrimaryKey, validation.In(idColumn).Error(`The primary key must be named "id".`)),
		),
		validation.Field(&f.PrimaryKey, validation.By(f.checkOtherFieldsForPK(collection))),
		validation.Field(&f.Min, validation.Min(0), validation.Max(maxSafeJSONInt)),
		validation.Field(&f.Max, validation.Min(f.Min), validation.Max(maxSafeJSONInt)),
		validation.Field(&f.Pattern, validation.When(f.PrimaryKey, validation.Required), validation.By(validators.IsRegex)),
		validation.Field(&f.Hidden, validation.When(f.PrimaryKey, validation.Empty)),
		validation.Field(&f.Required, validation.When(f.PrimaryKey, validation.Required)),
		validation.Field(&f.AutogeneratePattern, validation.By(validators.IsRegex), validation.By(f.checkAutogeneratePattern)),
	)
}

func (f *TextField) checkOtherFieldsForPK(collection *Collection) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(bool)
		if !v {
			return nil // not a pk
		}

		totalPrimaryKeys := 0
		for _, field := range collection.Fields {
			if text, ok := field.(*TextField); ok && text.PrimaryKey {
				totalPrimaryKeys++
			}

			if totalPrimaryKeys > 1 {
				return validation.NewError("validation_unsupported_composite_pk", "Composite PKs are not supported and the collection must have only 1 PK.")
			}
		}

		return nil
	}
}

func (f *TextField) checkAutogeneratePattern(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	// run 10 tests to check for conflicts with the other field validators
	for i := 0; i < 10; i++ {
		generated, err := security.RandomStringByRegex(v)
		if err != nil {
			return validation.NewError("validation_invalid_autogenerate_pattern", err.Error())
		}

		// (loosely) check whether the generated pattern satisfies the current field settings
		if err := f.ValidatePlainValue(generated); err != nil {
			return validation.NewError(
				"validation_invalid_autogenerate_pattern_value",
				fmt.Sprintf("The provided autogenerate pattern could produce invalid field values, ex.: %q", generated),
			)
		}
	}

	return nil
}

// Intercept implements the [RecordInterceptor] interface.
func (f *TextField) Intercept(
	ctx context.Context,
	app App,
	record *Record,
	actionName string,
	actionFunc func() error,
) error {
	// set autogenerated value if missing for new records
	switch actionName {
	case InterceptorActionValidate, InterceptorActionCreate:
		if f.AutogeneratePattern != "" && f.hasZeroValue(record) && record.IsNew() {
			v, err := security.RandomStringByRegex(f.AutogeneratePattern)
			if err != nil {
				return fmt.Errorf("failed to autogenerate %q value: %w", f.Name, err)
			}
			record.SetRaw(f.Name, v)
		}
	}

	return actionFunc()
}

func (f *TextField) hasZeroValue(record *Record) bool {
	v, _ := record.GetRaw(f.Name).(string)
	return v == ""
}

// FindSetter implements the [SetterFinder] interface.
func (f *TextField) FindSetter(key string) SetterFunc {
	switch key {
	case f.Name:
		return func(record *Record, raw any) {
			record.SetRaw(f.Name, cast.ToString(raw))
		}
	case f.Name + autogenerateModifier:
		return func(record *Record, raw any) {
			v := cast.ToString(raw)

			if f.AutogeneratePattern != "" {
				generated, _ := security.RandomStringByRegex(f.AutogeneratePattern)
				v += generated
			}

			record.SetRaw(f.Name, v)
		}
	default:
		return nil
	}
}
