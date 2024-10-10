package validators

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// IsRegex checks whether the validated value is a valid regular expression pattern.
//
// Example:
//
//	validation.Field(&form.Pattern, validation.By(validators.IsRegex))
func IsRegex(value any) error {
	v, ok := value.(string)
	if !ok {
		return ErrUnsupportedValueType
	}

	if v == "" {
		return nil // nothing to check
	}

	if _, err := regexp.Compile(v); err != nil {
		return validation.NewError("validation_invalid_regex", err.Error())
	}

	return nil
}
