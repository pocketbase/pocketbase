package validators

import (
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Equal checks whether the validated value matches another one from the same type.
//
// It expects the compared values to be from the same type and works
// with booleans, numbers, strings and their pointer variants.
//
// If one of the value is pointer, the comparison is based on its
// underlying value (when possible to determine).
//
// Note that empty/zero values are also compared (this differ from other validation.RuleFunc).
//
// Example:
//
//	validation.Field(&form.PasswordConfirm, validation.By(validators.Equal(form.Password)))
func Equal[T comparable](valueToCompare T) validation.RuleFunc {
	return func(value any) error {
		if compareValues(value, valueToCompare) {
			return nil
		}

		return validation.NewError("validation_values_mismatch", "Values don't match.")
	}
}

func compareValues(a, b any) bool {
	if a == b {
		return true
	}

	if checkIsNil(a) && checkIsNil(b) {
		return true
	}

	var result bool

	defer func() {
		if err := recover(); err != nil {
			result = false
		}
	}()

	reflectA := reflect.ValueOf(a)
	reflectB := reflect.ValueOf(b)

	dereferencedA := dereference(reflectA)
	dereferencedB := dereference(reflectB)
	if dereferencedA.CanInterface() && dereferencedB.CanInterface() {
		result = dereferencedA.Interface() == dereferencedB.Interface()
	}

	return result
}

// note https://github.com/golang/go/issues/51649
func checkIsNil(value any) bool {
	if value == nil {
		return true
	}

	var result bool

	defer func() {
		if err := recover(); err != nil {
			result = false
		}
	}()

	result = reflect.ValueOf(value).IsNil()

	return result
}

func dereference(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	return v
}
