package search_test

import (
	"regexp"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/search"
)

func TestFilterDataBuildExpr(t *testing.T) {
	resolver := search.NewSimpleFieldResolver("test1", "test2", "test3", "test4.sub")

	scenarios := []struct {
		filterData    search.FilterData
		expectError   bool
		expectPattern string
	}{
		// empty
		{"", true, ""},
		// invalid format
		{"(test1 > 1", true, ""},
		// invalid operator
		{"test1 + 123", true, ""},
		// unknown field
		{"test1 = 'example' && unknown > 1", true, ""},
		// simple expression
		{"test1 > 1", false,
			"^" +
				regexp.QuoteMeta("[[test1]] > {:") +
				".+" +
				regexp.QuoteMeta("}") +
				"$",
		},
		// like with 2 columns
		{"test1 ~ test2", false,
			"^" +
				regexp.QuoteMeta("[[test1]] LIKE ('%' || [[test2]] || '%')") +
				"$",
		},
		// like with right column operand
		{"'lorem' ~ test1", false,
			"^" +
				regexp.QuoteMeta("{:") +
				".+" +
				regexp.QuoteMeta("} LIKE ('%' || [[test1]] || '%')") +
				"$",
		},
		// like with left column operand and text as right operand
		{"test1 ~ 'lorem'", false,
			"^" +
				regexp.QuoteMeta("[[test1]] LIKE {:") +
				".+" +
				regexp.QuoteMeta("}") +
				"$",
		},
		// not like with 2 columns
		{"test1 !~ test2", false,
			"^" +
				regexp.QuoteMeta("[[test1]] NOT LIKE ('%' || [[test2]] || '%')") +
				"$",
		},
		// not like with right column operand
		{"'lorem' !~ test1", false,
			"^" +
				regexp.QuoteMeta("{:") +
				".+" +
				regexp.QuoteMeta("} NOT LIKE ('%' || [[test1]] || '%')") +
				"$",
		},
		// like with left column operand and text as right operand
		{"test1 !~ 'lorem'", false,
			"^" +
				regexp.QuoteMeta("[[test1]] NOT LIKE {:") +
				".+" +
				regexp.QuoteMeta("}") +
				"$",
		},
		// current datetime constant
		{"test1 > @now", false,
			"^" +
				regexp.QuoteMeta("[[test1]] > {:") +
				".+" +
				regexp.QuoteMeta("}") +
				"$",
		},
		// complex expression
		{
			"((test1 > 1) || (test2 != 2)) && test3 ~ '%%example' && test4.sub = null",
			false,
			"^" +
				regexp.QuoteMeta("((([[test1]] > {:") +
				".+" +
				regexp.QuoteMeta("}) OR (COALESCE([[test2]], '') != COALESCE({:") +
				".+" +
				regexp.QuoteMeta("}, ''))) AND ([[test3]] LIKE {:") +
				".+" +
				regexp.QuoteMeta("})) AND (COALESCE([[test4.sub]], '') = COALESCE(NULL, ''))") +
				"$",
		},
		// combination of special literals (null, true, false)
		{
			"test1=true && test2 != false && test3 = null || test4.sub != null",
			false,
			"^" + regexp.QuoteMeta("(((COALESCE([[test1]], '') = COALESCE(1, '')) AND (COALESCE([[test2]], '') != COALESCE(0, ''))) AND (COALESCE([[test3]], '') = COALESCE(NULL, ''))) OR (COALESCE([[test4.sub]], '') != COALESCE(NULL, ''))") + "$",
		},
		// all operators
		{
			"(test1 = test2 || test2 != test3) && (test2 ~ 'example' || test2 !~ '%%abc') && 'switch1%%' ~ test1 && 'switch2' !~ test2 && test3 > 1 && test3 >= 0 && test3 <= 4 && 2 < 5",
			false,
			"^" +
				regexp.QuoteMeta("((((((((COALESCE([[test1]], '') = COALESCE([[test2]], '')) OR (COALESCE([[test2]], '') != COALESCE([[test3]], ''))) AND (([[test2]] LIKE {:") +
				".+" +
				regexp.QuoteMeta("}) OR ([[test2]] NOT LIKE {:") +
				".+" +
				regexp.QuoteMeta("}))) AND ({:") +
				".+" +
				regexp.QuoteMeta("} LIKE ('%' || [[test1]] || '%'))) AND ({:") +
				".+" +
				regexp.QuoteMeta("} NOT LIKE ('%' || [[test2]] || '%'))) AND ([[test3]] > {:") +
				".+" +
				regexp.QuoteMeta("})) AND ([[test3]] >= {:") +
				".+" +
				regexp.QuoteMeta("})) AND ([[test3]] <= {:") +
				".+" +
				regexp.QuoteMeta("})) AND ({:") +
				".+" +
				regexp.QuoteMeta("} < {:") +
				".+" +
				regexp.QuoteMeta("})") +
				"$",
		},
	}

	for i, s := range scenarios {
		expr, err := s.filterData.BuildExpr(resolver)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		dummyDB := &dbx.DB{}
		rawSql := expr.Build(dummyDB, map[string]any{})

		pattern := regexp.MustCompile(s.expectPattern)
		if !pattern.MatchString(rawSql) {
			t.Errorf("(%d) Pattern %v don't match with expression: \n%v", i, s.expectPattern, rawSql)
		}
	}
}
