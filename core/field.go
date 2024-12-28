package core

import (
	"context"
	"database/sql/driver"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/list"
)

var fieldNameRegex = regexp.MustCompile(`^\w+$`)

const maxSafeJSONInt int64 = 1<<53 - 1

// Commonly used field names.
const (
	FieldNameId              = "id"
	FieldNameCollectionId    = "collectionId"
	FieldNameCollectionName  = "collectionName"
	FieldNameExpand          = "expand"
	FieldNameEmail           = "email"
	FieldNameEmailVisibility = "emailVisibility"
	FieldNameVerified        = "verified"
	FieldNameTokenKey        = "tokenKey"
	FieldNamePassword        = "password"
)

// SystemFields returns special internal field names that are usually readonly.
var SystemDynamicFieldNames = []string{
	FieldNameCollectionId,
	FieldNameCollectionName,
	FieldNameExpand,
}

// Common RecordInterceptor action names.
const (
	InterceptorActionValidate         = "validate"
	InterceptorActionDelete           = "delete"
	InterceptorActionDeleteExecute    = "deleteExecute"
	InterceptorActionAfterDelete      = "afterDelete"
	InterceptorActionAfterDeleteError = "afterDeleteError"
	InterceptorActionCreate           = "create"
	InterceptorActionCreateExecute    = "createExecute"
	InterceptorActionAfterCreate      = "afterCreate"
	InterceptorActionAfterCreateError = "afterCreateFailure"
	InterceptorActionUpdate           = "update"
	InterceptorActionUpdateExecute    = "updateExecute"
	InterceptorActionAfterUpdate      = "afterUpdate"
	InterceptorActionAfterUpdateError = "afterUpdateError"
)

// Common field errors.
var (
	ErrUnknownField          = validation.NewError("validation_unknown_field", "Unknown or invalid field.")
	ErrInvalidFieldValue     = validation.NewError("validation_invalid_field_value", "Invalid field value.")
	ErrMustBeSystemAndHidden = validation.NewError("validation_must_be_system_and_hidden", `The field must be marked as "System" and "Hidden".`)
	ErrMustBeSystem          = validation.NewError("validation_must_be_system", `The field must be marked as "System".`)
)

// FieldFactoryFunc defines a simple function to construct a specific Field instance.
type FieldFactoryFunc func() Field

// Fields holds all available collection fields.
var Fields = map[string]FieldFactoryFunc{}

// Field defines a common interface that all Collection fields should implement.
type Field interface {
	// note: the getters has an explicit "Get" prefix to avoid conflicts with their related field members

	// GetId returns the field id.
	GetId() string

	// SetId changes the field id.
	SetId(id string)

	// GetName returns the field name.
	GetName() string

	// SetName changes the field name.
	SetName(name string)

	// GetSystem returns the field system flag state.
	GetSystem() bool

	// SetSystem changes the field system flag state.
	SetSystem(system bool)

	// GetHidden returns the field hidden flag state.
	GetHidden() bool

	// SetHidden changes the field hidden flag state.
	SetHidden(hidden bool)

	// Type returns the unique type of the field.
	Type() string

	// ColumnType returns the DB column definition of the field.
	ColumnType(app App) string

	// PrepareValue returns a properly formatted field value based on the provided raw one.
	//
	// This method is also called on record construction to initialize its default field value.
	PrepareValue(record *Record, raw any) (any, error)

	// ValidateSettings validates the current field value associated with the provided record.
	ValidateValue(ctx context.Context, app App, record *Record) error

	// ValidateSettings validates the current field settings.
	ValidateSettings(ctx context.Context, app App, collection *Collection) error
}

// MaxBodySizeCalculator defines an optional field interface for
// specifying the max size of a field value.
type MaxBodySizeCalculator interface {
	// CalculateMaxBodySize returns the approximate max body size of a field value.
	CalculateMaxBodySize() int64
}

type (
	SetterFunc func(record *Record, raw any)

	// SetterFinder defines a field interface for registering custom field value setters.
	SetterFinder interface {
		// FindSetter returns a single field value setter function
		// by performing pattern-like field matching using the specified key.
		//
		// The key is usually just the field name but it could also
		// contains "modifier" characters based on which you can perform custom set operations
		// (ex. "users+" could be mapped to a function that will append new user to the existing field value).
		//
		// Return nil if you want to fallback to the default field value setter.
		FindSetter(key string) SetterFunc
	}
)

type (
	GetterFunc func(record *Record) any

	// GetterFinder defines a field interface for registering custom field value getters.
	GetterFinder interface {
		// FindGetter returns a single field value getter function
		// by performing pattern-like field matching using the specified key.
		//
		// The key is usually just the field name but it could also
		// contains "modifier" characters based on which you can perform custom get operations
		// (ex. "description:excerpt" could be mapped to a function that will return an excerpt of the current field value).
		//
		// Return nil if you want to fallback to the default field value setter.
		FindGetter(key string) GetterFunc
	}
)

// DriverValuer defines a Field interface for exporting and formatting
// a field value for the database.
type DriverValuer interface {
	// DriverValue exports a single field value for persistence in the database.
	DriverValue(record *Record) (driver.Value, error)
}

// MultiValuer defines a field interface that every multi-valued (eg. with MaxSelect) field has.
type MultiValuer interface {
	// IsMultiple checks whether the field is configured to support multiple or single values.
	IsMultiple() bool
}

// RecordInterceptor defines a field interface for reacting to various
// Record related operations (create, delete, validate, etc.).
type RecordInterceptor interface {
	// Interceptor is invoked when a specific record action occurs
	// allowing you to perform extra validations and normalization
	// (ex. uploading or deleting files).
	//
	// Note that users must call actionFunc() manually if they want to
	// execute the specific record action.
	Intercept(
		ctx context.Context,
		app App,
		record *Record,
		actionName string,
		actionFunc func() error,
	) error
}

// DefaultFieldIdValidationRule performs base validation on a field id value.
func DefaultFieldIdValidationRule(value any) error {
	v, ok := value.(string)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	rules := []validation.Rule{
		validation.Required,
		validation.Length(1, 100),
	}

	for _, r := range rules {
		if err := r.Validate(v); err != nil {
			return err
		}
	}

	return nil
}

// exclude special filter and system literals
var excludeNames = append([]any{
	"null", "true", "false", "_rowid_",
}, list.ToInterfaceSlice(SystemDynamicFieldNames)...)

// DefaultFieldIdValidationRule performs base validation on a field name value.
func DefaultFieldNameValidationRule(value any) error {
	v, ok := value.(string)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	rules := []validation.Rule{
		validation.Required,
		validation.Length(1, 100),
		validation.Match(fieldNameRegex),
		validation.NotIn(excludeNames...),
		validation.By(checkForVia),
	}

	for _, r := range rules {
		if err := r.Validate(v); err != nil {
			return err
		}
	}

	return nil
}

func checkForVia(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil
	}

	if strings.Contains(strings.ToLower(v), "_via_") {
		return validation.NewError("validation_found_via", `The value cannot contain "_via_".`)
	}

	return nil
}

func noopSetter(record *Record, raw any) {
	// do nothing
}
