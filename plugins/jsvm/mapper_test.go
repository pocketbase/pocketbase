package jsvm_test

import (
	"reflect"
	"testing"

	"github.com/pocketbase/pocketbase/plugins/jsvm"
)

func TestFieldMapper(t *testing.T) {
	mapper := jsvm.FieldMapper{}

	scenarios := []struct {
		name     string
		expected string
	}{
		{"", ""},
		{"test", "test"},
		{"Test", "test"},
		{"miXeD", "miXeD"},
		{"MiXeD", "miXeD"},
		{"ResolveRequestAsJSON", "resolveRequestAsJSON"},
		{"Variable_with_underscore", "variable_with_underscore"},
		{"ALLCAPS", "allcaps"},
		{"ALL_CAPS_WITH_UNDERSCORE", "all_caps_with_underscore"},
		{"OIDCMap", "oidcMap"},
		{"MD5", "md5"},
		{"OAuth2", "oauth2"},
	}

	for i, s := range scenarios {
		field := reflect.StructField{Name: s.name}
		if v := mapper.FieldName(nil, field); v != s.expected {
			t.Fatalf("[%d] Expected FieldName %q, got %q", i, s.expected, v)
		}

		method := reflect.Method{Name: s.name}
		if v := mapper.MethodName(nil, method); v != s.expected {
			t.Fatalf("[%d] Expected MethodName %q, got %q", i, s.expected, v)
		}
	}
}
