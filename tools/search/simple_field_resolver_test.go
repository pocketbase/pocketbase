package search_test

import (
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
		db := dbx.NewFromDB(nil, "")
		query := db.Select("id").From("test")

		r.Resolve(s.fieldName)

		if err := r.UpdateQuery(nil); err != nil {
			t.Errorf("(%d) UpdateQuery failed with error %v", i, err)
			continue
		}

		rawQuery := query.Build().SQL()
		// rawQuery := s.expectQuery

		if rawQuery != s.expectQuery {
			t.Errorf("(%d) Expected query %v, got \n%v", i, s.expectQuery, rawQuery)
		}
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
		r, err := r.Resolve(s.fieldName)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		if r.Identifier != s.expectName {
			t.Errorf("(%d) Expected r.Identifier %q, got %q", i, s.expectName, r.Identifier)
		}

		// params should be empty
		if len(r.Params) != 0 {
			t.Errorf("(%d) Expected 0 r.Params, got %v", i, r.Params)
		}
	}
}
