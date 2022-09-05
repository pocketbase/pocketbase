package list

import (
	"encoding/json"
	"github.com/pocketbase/pocketbase/tools/store"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

var cachedPatterns = store.New(make(map[string]*regexp.Regexp, 50))

// ExistInSlice checks whether a comparable element exists in a slice of the same type.
func ExistInSlice[T comparable](item T, list []T) bool {
	if len(list) == 0 {
		return false
	}

	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}

// ExistInSliceWithRegex checks whether a string exists in a slice
// either by direct match, or by a regular expression (eg. `^\w+$`).
//
// _Note: Only list items starting with '^' and ending with '$' are treated as regular expressions!_
func ExistInSliceWithRegex(str string, list []string) bool {
	for _, field := range list {
		isRegex := strings.HasPrefix(field, "^") && strings.HasSuffix(field, "$")

		if !isRegex {
			// check for direct match
			if str == field {
				return true
			}
			continue
		}

		// check for regex match
		pattern, ok := cachedPatterns.MightGet(field)
		if !ok {
			var err error
			pattern, err = regexp.Compile(field)
			if err != nil {
				continue
			}
			// "cache" the pattern to avoid compiling it every time
			cachedPatterns.Set(field, pattern)
		}

		if pattern != nil && pattern.MatchString(str) {
			return true
		}
	}

	return false
}

// ToInterfaceSlice converts a generic slice to slice of interfaces.
func ToInterfaceSlice[T any](list []T) []any {
	result := make([]any, len(list))

	for i := range list {
		result[i] = list[i]
	}

	return result
}

// NonzeroUniques returns only the nonzero unique values from a slice.
func NonzeroUniques[T comparable](list []T) []T {
	result := make([]T, 0, len(list))
	existMap := make(map[T]struct{}, len(list))

	var zeroVal T

	for _, val := range list {
		if _, ok := existMap[val]; ok || val == zeroVal {
			continue
		}
		existMap[val] = struct{}{}
		result = append(result, val)
	}

	return result
}

// ToUniqueStringSlice casts `value` to a slice of non-zero unique strings.
func ToUniqueStringSlice(value any) (result []string) {
	switch val := value.(type) {
	case nil:
		// nothing to cast
	case []string:
		result = val
	case string:
		if val == "" {
			break
		}

		// check if it is a json encoded array of strings
		if err := json.Unmarshal([]byte(val), &result); err != nil {
			// not a json array, just add the string as single array element
			result = append(result, val)
		}
	case json.Marshaler: // eg. JsonArray
		raw, _ := val.MarshalJSON()
		_ = json.Unmarshal(raw, &result)
	default:
		result = cast.ToStringSlice(value)
	}

	return NonzeroUniques(result)
}
