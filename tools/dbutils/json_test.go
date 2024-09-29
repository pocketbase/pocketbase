package dbutils_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/dbutils"
)

func TestJSONEach(t *testing.T) {
	result := dbutils.JSONEach("a.b")

	expected := "json_each(CASE WHEN json_valid([[a.b]]) THEN [[a.b]] ELSE json_array([[a.b]]) END)"

	if result != expected {
		t.Fatalf("Expected\n%v\ngot\n%v", expected, result)
	}
}

func TestJSONArrayLength(t *testing.T) {
	result := dbutils.JSONArrayLength("a.b")

	expected := "json_array_length(CASE WHEN json_valid([[a.b]]) THEN [[a.b]] ELSE (CASE WHEN [[a.b]] = '' OR [[a.b]] IS NULL THEN json_array() ELSE json_array([[a.b]]) END) END)"

	if result != expected {
		t.Fatalf("Expected\n%v\ngot\n%v", expected, result)
	}
}

func TestJSONExtract(t *testing.T) {
	scenarios := []struct {
		name     string
		column   string
		path     string
		expected string
	}{
		{
			"empty path",
			"a.b",
			"",
			"(CASE WHEN json_valid([[a.b]]) THEN JSON_EXTRACT([[a.b]], '$') ELSE JSON_EXTRACT(json_object('pb', [[a.b]]), '$.pb') END)",
		},
		{
			"starting with array index",
			"a.b",
			"[1].a[2]",
			"(CASE WHEN json_valid([[a.b]]) THEN JSON_EXTRACT([[a.b]], '$[1].a[2]') ELSE JSON_EXTRACT(json_object('pb', [[a.b]]), '$.pb[1].a[2]') END)",
		},
		{
			"starting with key",
			"a.b",
			"a.b[2].c",
			"(CASE WHEN json_valid([[a.b]]) THEN JSON_EXTRACT([[a.b]], '$.a.b[2].c') ELSE JSON_EXTRACT(json_object('pb', [[a.b]]), '$.pb.a.b[2].c') END)",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := dbutils.JSONExtract(s.column, s.path)

			if result != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, result)
			}
		})
	}
}
