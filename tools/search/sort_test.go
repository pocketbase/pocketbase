package search_test

import (
	"encoding/json"
	"testing"

	"github.com/pocketbase/pocketbase/tools/search"
)

func TestSortFieldBuildExpr(t *testing.T) {
	resolver := search.NewSimpleFieldResolver("test1", "test2", "test3", "test4.sub")

	scenarios := []struct {
		sortField        search.SortField
		expectError      bool
		expectExpression string
	}{
		// empty
		{search.SortField{"", search.SortDesc}, true, ""},
		// unknown field
		{search.SortField{"unknown", search.SortAsc}, true, ""},
		// placeholder field
		{search.SortField{"'test'", search.SortAsc}, true, ""},
		// null field
		{search.SortField{"null", search.SortAsc}, true, ""},
		// allowed field - asc
		{search.SortField{"test1", search.SortAsc}, false, "[[test1]] ASC"},
		// allowed field - desc
		{search.SortField{"test1", search.SortDesc}, false, "[[test1]] DESC"},
		// special @random field (ignore direction)
		{search.SortField{"@random", search.SortDesc}, false, "RANDOM()"},
	}

	for i, s := range scenarios {
		result, err := s.sortField.BuildExpr(resolver)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}

		if result != s.expectExpression {
			t.Errorf("(%d) Expected expression %v, got %v", i, s.expectExpression, result)
		}
	}
}

func TestParseSortFromString(t *testing.T) {
	scenarios := []struct {
		value        string
		expectedJson string
	}{
		{"", `[{"name":"","direction":"ASC"}]`},
		{"test", `[{"name":"test","direction":"ASC"}]`},
		{"+test", `[{"name":"test","direction":"ASC"}]`},
		{"-test", `[{"name":"test","direction":"DESC"}]`},
		{"test1,-test2,+test3", `[{"name":"test1","direction":"ASC"},{"name":"test2","direction":"DESC"},{"name":"test3","direction":"ASC"}]`},
		{"@random,-test", `[{"name":"@random","direction":"ASC"},{"name":"test","direction":"DESC"}]`},
	}

	for i, s := range scenarios {
		result := search.ParseSortFromString(s.value)
		encoded, _ := json.Marshal(result)

		if string(encoded) != s.expectedJson {
			t.Errorf("(%d) Expected expression %v, got %v", i, s.expectedJson, string(encoded))
		}
	}
}
