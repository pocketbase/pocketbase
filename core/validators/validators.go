// Package validators implements some common custom PocketBase validators.
package validators

import (
	"errors"
	"maps"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var ErrUnsupportedValueType = validation.NewError("validation_unsupported_value_type", "Invalid or unsupported value type.")

// JoinValidationErrors attempts to join the provided [validation.Errors] arguments.
//
// If only one of the arguments is [validation.Errors], it returns the first non-empty [validation.Errors].
//
// If both arguments are not [validation.Errors] then it returns a combined [errors.Join] error.
func JoinValidationErrors(errA, errB error) error {
	vErrA, okA := errA.(validation.Errors)
	vErrB, okB := errB.(validation.Errors)

	// merge
	if okA && okB {
		result := maps.Clone(vErrA)
		maps.Copy(result, vErrB)
		if len(result) > 0 {
			return result
		}
	}

	if okA && len(vErrA) > 0 {
		return vErrA
	}

	if okB && len(vErrB) > 0 {
		return vErrB
	}

	return errors.Join(errA, errB)
}
