package jsvm

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/dop251/goja"
)

var (
	_ goja.FieldNameMapper = (*FieldMapper)(nil)
)

// FieldMapper provides custom mapping between Go and JavaScript property names.
//
// It is similar to the builtin "uncapFieldNameMapper" but also converts
// all uppercase identifiers to their lowercase equivalent (eg. "GET" -> "get").
type FieldMapper struct {
}

// FieldName implements the [FieldNameMapper.FieldName] interface method.
func (u FieldMapper) FieldName(_ reflect.Type, f reflect.StructField) string {
	return convertGoToJSName(f.Name)
}

// MethodName implements the [FieldNameMapper.MethodName] interface method.
func (u FieldMapper) MethodName(_ reflect.Type, m reflect.Method) string {
	return convertGoToJSName(m.Name)
}

func convertGoToJSName(name string) string {
	allUppercase := true
	for _, c := range name {
		if c != '_' && !unicode.IsUpper(c) {
			allUppercase = false
			break
		}
	}

	// eg. "JSON" -> "json"
	if allUppercase {
		return strings.ToLower(name)
	}

	// eg. "GetField" -> "getField"
	return strings.ToLower(name[0:1]) + name[1:]
}
