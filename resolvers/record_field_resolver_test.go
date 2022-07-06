package resolvers_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordFieldResolverUpdateQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo4")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		fieldName        string
		expectQueryParts []string // we are matching parts of the query
		// since joins are added with map iteration and the order is not guaranteed
	}{
		// missing field
		{"", []string{
			"SELECT `demo4`.* FROM `demo4`",
		}},
		// non relation field
		{"title", []string{
			"SELECT `demo4`.* FROM `demo4`",
		}},
		// incomplete rel
		{"onerel", []string{
			"SELECT `demo4`.* FROM `demo4`",
		}},
		// single rel
		{"onerel.title", []string{
			"SELECT DISTINCT `demo4`.* FROM `demo4`",
			" LEFT JOIN `demo4` `demo4_onerel` ON [[demo4.onerel]] LIKE ('%' || [[demo4_onerel.id]] || '%')",
		}},
		// nested incomplete rels
		{"manyrels.onerel", []string{
			"SELECT DISTINCT `demo4`.* FROM `demo4`",
			" LEFT JOIN `demo4` `demo4_manyrels` ON [[demo4.manyrels]] LIKE ('%' || [[demo4_manyrels.id]] || '%')",
		}},
		// nested complete rels
		{"manyrels.onerel.title", []string{
			"SELECT DISTINCT `demo4`.* FROM `demo4`",
			" LEFT JOIN `demo4` `demo4_manyrels` ON [[demo4.manyrels]] LIKE ('%' || [[demo4_manyrels.id]] || '%')",
			" LEFT JOIN `demo4` `demo4_manyrels_onerel` ON [[demo4_manyrels.onerel]] LIKE ('%' || [[demo4_manyrels_onerel.id]] || '%')",
		}},
		// // repeated nested rels
		{"manyrels.onerel.manyrels.onerel.title", []string{
			"SELECT DISTINCT `demo4`.* FROM `demo4`",
			" LEFT JOIN `demo4` `demo4_manyrels` ON [[demo4.manyrels]] LIKE ('%' || [[demo4_manyrels.id]] || '%')",
			" LEFT JOIN `demo4` `demo4_manyrels_onerel` ON [[demo4_manyrels.onerel]] LIKE ('%' || [[demo4_manyrels_onerel.id]] || '%')",
			" LEFT JOIN `demo4` `demo4_manyrels_onerel_manyrels` ON [[demo4_manyrels_onerel.manyrels]] LIKE ('%' || [[demo4_manyrels_onerel_manyrels.id]] || '%')",
			" LEFT JOIN `demo4` `demo4_manyrels_onerel_manyrels_onerel` ON [[demo4_manyrels_onerel_manyrels.onerel]] LIKE ('%' || [[demo4_manyrels_onerel_manyrels_onerel.id]] || '%')",
		}},
	}

	for i, s := range scenarios {
		query := app.Dao().RecordQuery(collection)

		r := resolvers.NewRecordFieldResolver(app.Dao(), collection, nil)
		r.Resolve(s.fieldName)

		if err := r.UpdateQuery(query); err != nil {
			t.Errorf("(%d) UpdateQuery failed with error %v", i, err)
			continue
		}

		rawQuery := query.Build().SQL()

		partsLength := 0
		for _, part := range s.expectQueryParts {
			partsLength += len(part)
			if !strings.Contains(rawQuery, part) {
				t.Errorf("(%d) Part %v is missing from query \n%v", i, part, rawQuery)
			}
		}

		if partsLength != len(rawQuery) {
			t.Errorf("(%d) Expected %d characters, got %d in \n%v", i, partsLength, len(rawQuery), rawQuery)
		}
	}
}

func TestRecordFieldResolverResolveSchemaFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo4")
	if err != nil {
		t.Fatal(err)
	}

	r := resolvers.NewRecordFieldResolver(app.Dao(), collection, nil)

	scenarios := []struct {
		fieldName   string
		expectError bool
		expectName  string
	}{
		{"", true, ""},
		{" ", true, ""},
		{"unknown", true, ""},
		{"invalid format", true, ""},
		{"id", false, "[[demo4.id]]"},
		{"created", false, "[[demo4.created]]"},
		{"updated", false, "[[demo4.updated]]"},
		{"title", false, "[[demo4.title]]"},
		{"title.test", true, ""},
		{"manyrels", false, "[[demo4.manyrels]]"},
		{"manyrels.", true, ""},
		{"manyrels.unknown", true, ""},
		{"manyrels.title", false, "[[demo4_manyrels.title]]"},
		{"manyrels.onerel.manyrels.onefile", false, "[[demo4_manyrels_onerel_manyrels.onefile]]"},
		{"@collect", true, ""},
		{"collection.demo4.title", true, ""},
		{"@collection", true, ""},
		{"@collection.unknown", true, ""},
		{"@collection.demo", true, ""},
		{"@collection.demo.", true, ""},
		{"@collection.demo.title", false, "[[c_demo.title]]"},
		{"@collection.demo4.title", false, "[[c_demo4.title]]"},
		{"@collection.demo4.id", false, "[[c_demo4.id]]"},
		{"@collection.demo4.created", false, "[[c_demo4.created]]"},
		{"@collection.demo4.updated", false, "[[c_demo4.updated]]"},
		{"@collection.demo4.manyrels.missing", true, ""},
		{"@collection.demo4.manyrels.onerel.manyrels.onerel.onefile", false, "[[c_demo4_manyrels_onerel_manyrels_onerel.onefile]]"},
	}

	for i, s := range scenarios {
		name, params, err := r.Resolve(s.fieldName)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}

		if name != s.expectName {
			t.Errorf("(%d) Expected name %q, got %q", i, s.expectName, name)
		}

		// params should be empty for non @request fields
		if len(params) != 0 {
			t.Errorf("(%d) Expected 0 params, got %v", i, params)
		}
	}
}

func TestRecordFieldResolverResolveRequestDataFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo4")
	if err != nil {
		t.Fatal(err)
	}

	requestData := map[string]any{
		"method": "get",
		"query": map[string]any{
			"a": 123,
		},
		"data": map[string]any{
			"b": 456,
			"c": map[string]int{"sub": 1},
		},
		"user": nil,
	}

	r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestData)

	scenarios := []struct {
		fieldName        string
		expectError      bool
		expectParamValue string // encoded json
	}{
		{"@request", true, ""},
		{"@request.invalid format", true, ""},
		{"@request.invalid_format2!", true, ""},
		{"@request.missing", true, ""},
		{"@request.method", false, `"get"`},
		{"@request.query", true, ``},
		{"@request.query.a", false, `123`},
		{"@request.query.a.missing", false, ``},
		{"@request.data", true, ``},
		{"@request.data.b", false, `456`},
		{"@request.data.b.missing", false, ``},
		{"@request.data.c", false, `"{\"sub\":1}"`},
		{"@request.user", true, ""},
		{"@request.user.id", false, ""},
	}

	for i, s := range scenarios {
		name, params, err := r.Resolve(s.fieldName)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		// missing key
		// ---
		if len(params) == 0 {
			if name != "NULL" {
				t.Errorf("(%d) Expected 0 placeholder parameters, got %v", i, params)
			}
			continue
		}

		// existing key
		// ---
		if len(params) != 1 {
			t.Errorf("(%d) Expected 1 placeholder parameter, got %v", i, params)
			continue
		}

		var paramName string
		var paramValue any
		for k, v := range params {
			paramName = k
			paramValue = v
		}

		if name != ("{:" + paramName + "}") {
			t.Errorf("(%d) Expected parameter name %q, got %q", i, paramName, name)
		}

		encodedParamValue, _ := json.Marshal(paramValue)
		if string(encodedParamValue) != s.expectParamValue {
			t.Errorf("(%d) Expected params %v, got %v", i, s.expectParamValue, string(encodedParamValue))
		}
	}
}
