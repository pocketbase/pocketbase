package router

import (
	"encoding"
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strconv"
)

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

// JSONPayloadKey is the key for the special UnmarshalRequestData case
// used for reading serialized json payload without normalization.
const JSONPayloadKey string = "@jsonPayload"

// UnmarshalRequestData unmarshals url.Values type of data (query, multipart/form-data, etc.) into dst.
//
// dst must be a pointer to a map[string]any or struct.
//
// If dst is a map[string]any, each data value will be inferred and
// converted to its bool, numeric, or string equivalent value
// (refer to inferValue() for the exact rules).
//
// If dst is a struct, the following field types are supported:
//   - bool
//   - string
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - serialized json string if submitted under the special "@jsonPayload" key
//   - encoding.TextUnmarshaler
//   - pointer and slice variations of the above primitives (ex. *string, []string, *[]string []*string, etc.)
//   - named/anonymous struct fields
//     Dot-notation is used to target nested fields, ex. "nestedStructField.title".
//   - embedded struct fields
//     The embedded struct fields are treated by default as if they were defined in their parent struct.
//     If the embedded struct has a tag matching structTagKey then to set its fields the data keys must be prefixed with that tag
//     similar to the regular nested struct fields.
//
// structTagKey and structPrefix are used only when dst is a struct.
//
// structTagKey represents the tag to use to match a data entry with a struct field (defaults to "form").
// If the struct field doesn't have the structTagKey tag, then the exported struct field name will be used as it is.
//
// structPrefix could be provided if all of the data keys are prefixed with a common string
// and you want the struct field to match only the value without the structPrefix
// (ex. for "user.name", "user.email" data keys and structPrefix "user", it will match "name" and "email" struct fields).
//
// Note that while the method was inspired by binders from echo, gorrila/schema, ozzo-routing
// and other similar common routing packages, it is not intended to be a drop-in replacement.
//
// @todo Consider adding support for dot-notation keys, in addition to the prefix, (ex. parent.child.title) to express nested object keys.
func UnmarshalRequestData(data map[string][]string, dst any, structTagKey string, structPrefix string) error {
	if len(data) == 0 {
		return nil // nothing to unmarshal
	}

	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Pointer {
		return errors.New("dst must be a pointer")
	}

	dstValue = dereference(dstValue)

	dstType := dstValue.Type()

	switch dstType.Kind() {
	case reflect.Map: // map[string]any
		if dstType.Elem().Kind() != reflect.Interface {
			return errors.New("dst map value type must be any/interface{}")
		}

		for k, v := range data {
			if k == JSONPayloadKey {
				continue // unmarshaled separately
			}

			total := len(v)

			if total == 1 {
				dstValue.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(inferValue(v[0])))
			} else {
				normalized := make([]any, total)
				for i, vItem := range v {
					normalized[i] = inferValue(vItem)
				}
				dstValue.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(normalized))
			}
		}
	case reflect.Struct:
		// set a default tag key
		if structTagKey == "" {
			structTagKey = "form"
		}

		err := unmarshalInStructValue(data, dstValue, structTagKey, structPrefix)
		if err != nil {
			return err
		}
	default:
		return errors.New("dst must be a map[string]any or struct")
	}

	// @jsonPayload
	//
	// Special case to scan serialized json string without
	// normalization alongside the other data values
	// ---------------------------------------------------------------
	jsonPayloadValues := data[JSONPayloadKey]
	for _, payload := range jsonPayloadValues {
		if err := json.Unmarshal([]byte(payload), dst); err != nil {
			return err
		}
	}

	return nil
}

// unmarshalInStructValue unmarshals data into the provided struct reflect.Value fields.
func unmarshalInStructValue(
	data map[string][]string,
	dstStructValue reflect.Value,
	structTagKey string,
	structPrefix string,
) error {
	dstStructType := dstStructValue.Type()

	for i := 0; i < dstStructValue.NumField(); i++ {
		fieldType := dstStructType.Field(i)

		tag := fieldType.Tag.Get(structTagKey)

		if tag == "-" || (!fieldType.Anonymous && !fieldType.IsExported()) {
			continue // disabled or unexported non-anonymous struct field
		}

		fieldValue := dereference(dstStructValue.Field(i))

		ft := fieldType.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		isSlice := ft.Kind() == reflect.Slice
		if isSlice {
			ft = ft.Elem()
		}

		name := tag
		if name == "" && !fieldType.Anonymous {
			name = fieldType.Name
		}
		if name != "" && structPrefix != "" {
			name = structPrefix + "." + name
		}

		// (*)encoding.TextUnmarshaler field
		// ---
		if ft.Implements(textUnmarshalerType) || reflect.PointerTo(ft).Implements(textUnmarshalerType) {
			values, ok := data[name]
			if !ok || len(values) == 0 || !fieldValue.CanSet() {
				continue // no value to load or the field cannot be set
			}

			if isSlice {
				n := len(values)
				slice := reflect.MakeSlice(fieldValue.Type(), n, n)
				for i, v := range values {
					unmarshaler, ok := dereference(slice.Index(i)).Addr().Interface().(encoding.TextUnmarshaler)
					if ok {
						if err := unmarshaler.UnmarshalText([]byte(v)); err != nil {
							return err
						}
					}
				}
				fieldValue.Set(slice)
			} else {
				unmarshaler, ok := fieldValue.Addr().Interface().(encoding.TextUnmarshaler)
				if ok {
					if err := unmarshaler.UnmarshalText([]byte(values[0])); err != nil {
						return err
					}
				}
			}
			continue
		}

		// "regular" field
		// ---
		if ft.Kind() != reflect.Struct {
			values, ok := data[name]
			if !ok || len(values) == 0 || !fieldValue.CanSet() {
				continue // no value to load
			}

			if isSlice {
				n := len(values)
				slice := reflect.MakeSlice(fieldValue.Type(), n, n)
				for i, v := range values {
					if err := setRegularReflectedValue(dereference(slice.Index(i)), v); err != nil {
						return err
					}
				}
				fieldValue.Set(slice)
			} else {
				if err := setRegularReflectedValue(fieldValue, values[0]); err != nil {
					return err
				}
			}
			continue
		}

		// structs (embedded or nested)
		// ---
		// slice of structs
		if isSlice {
			// populating slice of structs is not supported at the moment
			// because the filling rules are ambiguous
			continue
		}

		if tag != "" {
			structPrefix = tag
		} else {
			structPrefix = name // name is empty for anonymous structs -> no prefix
		}

		if err := unmarshalInStructValue(data, fieldValue, structTagKey, structPrefix); err != nil {
			return err
		}
	}

	return nil
}

// dereference returns the underlying value v points to.
func dereference(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// initialize with a new value and continue searching
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

// setRegularReflectedValue sets and casts value into rv.
func setRegularReflectedValue(rv reflect.Value, value string) error {
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(value)
	case reflect.Bool:
		if value == "" {
			value = "f"
		}

		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		rv.SetBool(v)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if value == "" {
			value = "0"
		}

		v, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return err
		}

		rv.SetInt(v)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if value == "" {
			value = "0"
		}

		v, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return err
		}

		rv.SetUint(v)
	case reflect.Float32, reflect.Float64:
		if value == "" {
			value = "0"
		}

		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		rv.SetFloat(v)
	default:
		return errors.New("unknown value type " + rv.Kind().String())
	}

	return nil
}

var inferNumberCharsRegex = regexp.MustCompile(`^[\-\.\d]+$`)

// In order to support more seamlessly both json and multipart/form-data requests,
// the following normalization rules are applied for plain multipart string values:
//   - "true" is converted to the json "true"
//   - "false" is converted to the json "false"
//   - numeric strings are converted to json number ONLY if the resulted
//     minimal number string representation is the same as the provided raw string
//     (aka. scientific notations, "Infinity", "0.0", "0001", etc. are kept as string)
//   - any other string (empty string too) is left as it is
func inferValue(raw string) any {
	switch raw {
	case "":
		return raw
	case "true":
		return true
	case "false":
		return false
	default:
		// try to convert to number
		//
		// note: expects the provided raw string to match exactly with the minimal string representation of the parsed float
		if (raw[0] == '-' || (raw[0] >= '0' && raw[0] <= '9')) &&
			inferNumberCharsRegex.Match([]byte(raw)) {
			v, err := strconv.ParseFloat(raw, 64)
			if err == nil && strconv.FormatFloat(v, 'f', -1, 64) == raw {
				return v
			}
		}

		return raw
	}
}
