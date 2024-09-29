package search_test

import (
	"fmt"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/search"
)

func TestSimpleFieldResolverUpdateQuery(t *testing.T) {
	r := search.NewSimpleFieldResolver("test")

	scenarios := []struct {
		fieldName   string
		expectQuery string
	}{
		// missing field (the query shouldn't change)
		{"", `SELECT "id" FROM "test"`},
		// unknown field (the query shouldn't change)
		{"unknown", `SELECT "id" FROM "test"`},
		// allowed field (the query shouldn't change)
		{"test", `SELECT "id" FROM "test"`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.fieldName), func(t *testing.T) {
			db := dbx.NewFromDB(nil, "")
			query := db.Select("id").From("test")

			r.Resolve(s.fieldName)

			if err := r.UpdateQuery(nil); err != nil {
				t.Fatalf("UpdateQuery failed with error %v", err)
			}

			rawQuery := query.Build().SQL()

			if rawQuery != s.expectQuery {
				t.Fatalf("Expected query %v, got \n%v", s.expectQuery, rawQuery)
			}
		})
	}
}

func TestSimpleFieldResolverResolve(t *testing.T) {
	r := search.NewSimpleFieldResolver("test", `^test_regex\d+$`, "Test columnify!", "data.test")

	scenarios := []struct {
		fieldName   string
		expectError bool
		expectName  string
	}{
		{"", true, ""},
		{" ", true, ""},
		{"unknown", true, ""},
		{"test", false, "[[test]]"},
		{"test.sub", true, ""},
		{"test_regex", true, ""},
		{"test_regex1", false, "[[test_regex1]]"},
		{"Test columnify!", false, "[[Testcolumnify]]"},
		{"data.test", false, "JSON_EXTRACT([[data]], '$.test')"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.fieldName), func(t *testing.T) {
			r, err := r.Resolve(s.fieldName)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if r.Identifier != s.expectName {
				t.Fatalf("Expected r.Identifier %q, got %q", s.expectName, r.Identifier)
			}

			if len(r.Params) != 0 {
				t.Fatalf("r.Params should be empty, got %v", r.Params)
			}
		})
	}
}
