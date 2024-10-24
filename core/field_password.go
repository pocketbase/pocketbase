package core

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	Fields[FieldTypePassword] = func() Field {
		return &PasswordField{}
	}
}

const FieldTypePassword = "password"

var (
	_ Field             = (*PasswordField)(nil)
	_ GetterFinder      = (*PasswordField)(nil)
	_ SetterFinder      = (*PasswordField)(nil)
	_ DriverValuer      = (*PasswordField)(nil)
	_ RecordInterceptor = (*PasswordField)(nil)
)

// PasswordField defines "password" type field for storing bcrypt hashed strings
// (usually used only internally for the "password" auth collection system field).
//
// If you want to set a direct bcrypt hash as record field value you can use the SetRaw method, for example:
//
//	// generates a bcrypt hash of "123456" and set it as field value
//	// (record.GetString("password") returns the plain password until persisted, otherwise empty string)
//	record.Set("password", "123456")
//
//	// set directly a bcrypt hash of "123456" as field value
//	// (record.GetString("password") returns empty string)
//	record.SetRaw("password", "$2a$10$.5Elh8fgxypNUWhpUUr/xOa2sZm0VIaE0qWuGGl9otUfobb46T1Pq")
//
// The following additional getter keys are available:
//
//   - "fieldName:hash" - returns the bcrypt hash string of the record field value (if any). For example:
//     record.GetString("password:hash")
type PasswordField struct {
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

	// Pattern specifies an optional regex pattern to match against the field value.
	//
	// Leave it empty to skip the pattern check.
	Pattern string `form:"pattern" json:"pattern"`

	// Min specifies an optional required field string length.
	Min int `form:"min" json:"min"`

	// Max specifies an optional required field string length.
	//
	// If zero, fallback to max 71 bytes.
	Max int `form:"max" json:"max"`

	// Cost specifies the cost/weight/iteration/etc. bcrypt factor.
	//
	// If zero, fallback to [bcrypt.DefaultCost].
	//
	// If explicitly set, must be between [bcrypt.MinCost] and [bcrypt.MaxCost].
	Cost int `form:"cost" json:"cost"`

	// Required will require the field value to be non-empty string.
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *PasswordField) Type() string {
	return FieldTypePassword
}

// GetId implements [Field.GetId] interface method.
func (f *PasswordField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *PasswordField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *PasswordField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *PasswordField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *PasswordField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *PasswordField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *PasswordField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *PasswordField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *PasswordField) ColumnType(app App) string {
	return "TEXT DEFAULT '' NOT NULL"
}

// DriverValue implements the [DriverValuer] interface.
func (f *PasswordField) DriverValue(record *Record) (driver.Value, error) {
	fp := f.getPasswordValue(record)
	return fp.Hash, fp.LastError
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *PasswordField) PrepareValue(record *Record, raw any) (any, error) {
	return &PasswordFieldValue{
		Hash: cast.ToString(raw),
	}, nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *PasswordField) ValidateValue(ctx context.Context, app App, record *Record) error {
	fp, ok := record.GetRaw(f.Name).(*PasswordFieldValue)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if fp.LastError != nil {
		return fp.LastError
	}

	if f.Required {
		if err := validation.Required.Validate(fp.Hash); err != nil {
			return err
		}
	}

	if fp.Plain == "" {
		return nil // nothing to check
	}

	// note: casted to []rune to count multi-byte chars as one for the
	// sake of more intuitive UX and clearer user error messages
	//
	// note2: technically multi-byte strings could produce bigger length than the bcrypt limit
	// but it should be fine as it will be just truncated (even if it cuts a byte sequence in the middle)
	length := len([]rune(fp.Plain))

	if length < f.Min {
		return validation.NewError("validation_min_text_constraint", fmt.Sprintf("Must be at least %d character(s)", f.Min))
	}

	maxLength := f.Max
	if maxLength <= 0 {
		maxLength = 71
	}
	if length > maxLength {
		return validation.NewError("validation_max_text_constraint", fmt.Sprintf("Must be less than %d character(s)", maxLength))
	}

	if f.Pattern != "" {
		match, _ := regexp.MatchString(f.Pattern, fp.Plain)
		if !match {
			return validation.NewError("validation_invalid_format", "Invalid value format")
		}
	}

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *PasswordField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
		validation.Field(&f.Min, validation.Min(1), validation.Max(71)),
		validation.Field(&f.Max, validation.Min(f.Min), validation.Max(71)),
		validation.Field(&f.Cost, validation.Min(bcrypt.MinCost), validation.Max(bcrypt.MaxCost)),
		validation.Field(&f.Pattern, validation.By(validators.IsRegex)),
	)
}

func (f *PasswordField) getPasswordValue(record *Record) *PasswordFieldValue {
	raw := record.GetRaw(f.Name)

	switch v := raw.(type) {
	case *PasswordFieldValue:
		return v
	case string:
		// we assume that any raw string starting with $2 is bcrypt hash
		if strings.HasPrefix(v, "$2") {
			return &PasswordFieldValue{Hash: v}
		}
	}

	return &PasswordFieldValue{}
}

// Intercept implements the [RecordInterceptor] interface.
func (f *PasswordField) Intercept(
	ctx context.Context,
	app App,
	record *Record,
	actionName string,
	actionFunc func() error,
) error {
	switch actionName {
	case InterceptorActionAfterCreate, InterceptorActionAfterUpdate:
		// unset the plain field value after successful create/update
		fp := f.getPasswordValue(record)
		fp.Plain = ""
	}

	return actionFunc()
}

// FindGetter implements the [GetterFinder] interface.
func (f *PasswordField) FindGetter(key string) GetterFunc {
	switch key {
	case f.Name:
		return func(record *Record) any {
			return f.getPasswordValue(record).Plain
		}
	case f.Name + ":hash":
		return func(record *Record) any {
			return f.getPasswordValue(record).Hash
		}
	default:
		return nil
	}
}

// FindSetter implements the [SetterFinder] interface.
func (f *PasswordField) FindSetter(key string) SetterFunc {
	switch key {
	case f.Name:
		return f.setValue
	default:
		return nil
	}
}

func (f *PasswordField) setValue(record *Record, raw any) {
	fv := &PasswordFieldValue{
		Plain: cast.ToString(raw),
	}

	// hash the password
	if fv.Plain != "" {
		cost := f.Cost
		if cost <= 0 {
			cost = bcrypt.DefaultCost
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(fv.Plain), cost)
		if err != nil {
			fv.LastError = err
		}

		fv.Hash = string(hash)
	}

	record.SetRaw(f.Name, fv)
}

// -------------------------------------------------------------------

type PasswordFieldValue struct {
	LastError error
	Hash      string
	Plain     string
}

func (pv PasswordFieldValue) Validate(pass string) bool {
	if pv.Hash == "" || pv.LastError != nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(pv.Hash), []byte(pass))

	return err == nil
}
