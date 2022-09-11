package rest_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/rest"
)

func TestNormalizeUrl(t *testing.T) {
	scenarios := []struct {
		url         string
		expectError bool
		expectUrl   string
	}{
		{":/", true, ""},
		{"./", false, "./"},
		{"../../test////", false, "../../test/"},
		{"/a/b/c", false, "/a/b/c"},
		{"a/////b//c/", false, "a/b/c/"},
		{"/a/////b//c", false, "/a/b/c"},
		{"///a/b/c", false, "/a/b/c"},
		{"//a/b/c", false, "//a/b/c"}, // preserve "auto-schema"
		{"http://a/b/c", false, "http://a/b/c"},
		{"a//bc?test=1//dd", false, "a/bc?test=1//dd"},       // only the path is normalized
		{"a//bc?test=1#12///3", false, "a/bc?test=1#12///3"}, // only the path is normalized
	}

	for i, s := range scenarios {
		result, err := rest.NormalizeUrl(s.url)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v", i, s.expectError, hasErr)
		}

		if result != s.expectUrl {
			t.Errorf("(%d) Expected url %q, got %q", i, s.expectUrl, result)
		}
	}
}
