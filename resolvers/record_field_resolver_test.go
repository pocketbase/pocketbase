package resolvers_test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestRecordFieldResolverUpdateQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo4")
	if err != nil {
		t.Fatal(err)
	}

	requestData := map[string]any{
		"user": map[string]any{
			"id": "4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			"profile": map[string]any{
				"id":   "d13f60a4-5765-48c7-9e1d-3e782340f833",
				"name": "test",
			},
		},
	}

	scenarios := []struct {
		name        string
		fields      []string
		expectQuery string
	}{
		{
			"missing field",
			[]string{""},
			"SELECT `demo4`.* FROM `demo4`",
		},
		{
			"non relation field",
			[]string{"title"},
			"SELECT `demo4`.* FROM `demo4`",
		},
		{
			"incomplete rel",
			[]string{"onerel"},
			"SELECT `demo4`.* FROM `demo4`",
		},
		{
			"single rel",
			[]string{"onerel.title"},
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.onerel]]) THEN [[demo4.onerel]] ELSE json_array([[demo4.onerel]]) END) `demo4_onerel_je` LEFT JOIN `demo4` `demo4_onerel` ON [[demo4_onerel.id]] = [[demo4_onerel_je.value]]",
		},
		{
			"non-relation field + single rel",
			[]string{"title", "onerel.title"},
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.onerel]]) THEN [[demo4.onerel]] ELSE json_array([[demo4.onerel]]) END) `demo4_onerel_je` LEFT JOIN `demo4` `demo4_onerel` ON [[demo4_onerel.id]] = [[demo4_onerel_je.value]]",
		},
		{
			"nested incomplete rels",
			[]string{"manyrels.onerel"},
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.manyrels]]) THEN [[demo4.manyrels]] ELSE json_array([[demo4.manyrels]]) END) `demo4_manyrels_je` LEFT JOIN `demo4` `demo4_manyrels` ON [[demo4_manyrels.id]] = [[demo4_manyrels_je.value]]",
		},
		{
			"nested complete rels",
			[]string{"manyrels.onerel.title"},
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.manyrels]]) THEN [[demo4.manyrels]] ELSE json_array([[demo4.manyrels]]) END) `demo4_manyrels_je` LEFT JOIN `demo4` `demo4_manyrels` ON [[demo4_manyrels.id]] = [[demo4_manyrels_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_manyrels.onerel]]) THEN [[demo4_manyrels.onerel]] ELSE json_array([[demo4_manyrels.onerel]]) END) `demo4_manyrels_onerel_je` LEFT JOIN `demo4` `demo4_manyrels_onerel` ON [[demo4_manyrels_onerel.id]] = [[demo4_manyrels_onerel_je.value]]",
		},
		{
			"repeated nested rels",
			[]string{"manyrels.onerel.manyrels.onerel.title"},
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.manyrels]]) THEN [[demo4.manyrels]] ELSE json_array([[demo4.manyrels]]) END) `demo4_manyrels_je` LEFT JOIN `demo4` `demo4_manyrels` ON [[demo4_manyrels.id]] = [[demo4_manyrels_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_manyrels.onerel]]) THEN [[demo4_manyrels.onerel]] ELSE json_array([[demo4_manyrels.onerel]]) END) `demo4_manyrels_onerel_je` LEFT JOIN `demo4` `demo4_manyrels_onerel` ON [[demo4_manyrels_onerel.id]] = [[demo4_manyrels_onerel_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_manyrels_onerel.manyrels]]) THEN [[demo4_manyrels_onerel.manyrels]] ELSE json_array([[demo4_manyrels_onerel.manyrels]]) END) `demo4_manyrels_onerel_manyrels_je` LEFT JOIN `demo4` `demo4_manyrels_onerel_manyrels` ON [[demo4_manyrels_onerel_manyrels.id]] = [[demo4_manyrels_onerel_manyrels_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_manyrels_onerel_manyrels.onerel]]) THEN [[demo4_manyrels_onerel_manyrels.onerel]] ELSE json_array([[demo4_manyrels_onerel_manyrels.onerel]]) END) `demo4_manyrels_onerel_manyrels_onerel_je` LEFT JOIN `demo4` `demo4_manyrels_onerel_manyrels_onerel` ON [[demo4_manyrels_onerel_manyrels_onerel.id]] = [[demo4_manyrels_onerel_manyrels_onerel_je.value]]",
		},
		{
			"multiple rels",
			[]string{"manyrels.title", "onerel.onefile"},
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.manyrels]]) THEN [[demo4.manyrels]] ELSE json_array([[demo4.manyrels]]) END) `demo4_manyrels_je` LEFT JOIN `demo4` `demo4_manyrels` ON [[demo4_manyrels.id]] = [[demo4_manyrels_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4.onerel]]) THEN [[demo4.onerel]] ELSE json_array([[demo4.onerel]]) END) `demo4_onerel_je` LEFT JOIN `demo4` `demo4_onerel` ON [[demo4_onerel.id]] = [[demo4_onerel_je.value]]",
		},
		{
			"@collection join",
			[]string{"@collection.demo.title", "@collection.demo2.text", "@collection.demo.file"},
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo` `__collection_demo` LEFT JOIN `demo2` `__collection_demo2`",
		},
		{
			"static @request.user.profile fields",
			[]string{"@request.user.id", "@request.user.profile.id", "@request.data.demo"},
			"^" +
				regexp.QuoteMeta("SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `profiles` `__user_profiles` ON [[__user_profiles.id]] =") +
				" {:.*}$",
		},
		{
			"relational @request.user.profile fields",
			[]string{"@request.user.profile.rel.id", "@request.user.profile.rel.name"},
			"^" +
				regexp.QuoteMeta("SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `profiles` `__user_profiles` ON [[__user_profiles.id]] =") +
				" {:.*} " +
				regexp.QuoteMeta("LEFT JOIN json_each(CASE WHEN json_valid([[__user_profiles.rel]]) THEN [[__user_profiles.rel]] ELSE json_array([[__user_profiles.rel]]) END) `__user_profiles_rel_je` LEFT JOIN `profiles` `__user_profiles_rel` ON [[__user_profiles_rel.id]] = [[__user_profiles_rel_je.value]]") +
				"$",
		},
	}

	for _, s := range scenarios {
		query := app.Dao().RecordQuery(collection)

		r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestData)
		for _, field := range s.fields {
			r.Resolve(field)
		}

		if err := r.UpdateQuery(query); err != nil {
			t.Errorf("(%s) UpdateQuery failed with error %v", s.name, err)
			continue
		}

		rawQuery := query.Build().SQL()

		if !list.ExistInSliceWithRegex(rawQuery, []string{s.expectQuery}) {
			t.Errorf("(%s) Expected query\n %v \ngot:\n %v", s.name, s.expectQuery, rawQuery)
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

	requestData := map[string]any{
		"user": map[string]any{
			"id": "4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			"profile": map[string]any{
				"id": "d13f60a4-5765-48c7-9e1d-3e782340f833",
			},
		},
	}

	r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestData)

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
		// @request.user.profile relation join:
		{"@request.user.profile.rel", false, "[[__user_profiles.rel]]"},
		{"@request.user.profile.rel.name", false, "[[__user_profiles_rel.name]]"},
		// @collection fieds:
		{"@collect", true, ""},
		{"collection.demo4.title", true, ""},
		{"@collection", true, ""},
		{"@collection.unknown", true, ""},
		{"@collection.demo", true, ""},
		{"@collection.demo.", true, ""},
		{"@collection.demo.title", false, "[[__collection_demo.title]]"},
		{"@collection.demo4.title", false, "[[__collection_demo4.title]]"},
		{"@collection.demo4.id", false, "[[__collection_demo4.id]]"},
		{"@collection.demo4.created", false, "[[__collection_demo4.created]]"},
		{"@collection.demo4.updated", false, "[[__collection_demo4.updated]]"},
		{"@collection.demo4.manyrels.missing", true, ""},
		{"@collection.demo4.manyrels.onerel.manyrels.onerel.onefile", false, "[[__collection_demo4_manyrels_onerel_manyrels_onerel.onefile]]"},
	}

	for _, s := range scenarios {
		name, params, err := r.Resolve(s.fieldName)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%q) Expected hasErr %v, got %v (%v)", s.fieldName, s.expectError, hasErr, err)
			continue
		}

		if name != s.expectName {
			t.Errorf("(%q) Expected name %q, got %q", s.fieldName, s.expectName, name)
		}

		// params should be empty for non @request fields
		if len(params) != 0 {
			t.Errorf("(%q) Expected 0 params, got %v", s.fieldName, params)
		}
	}
}

func TestRecordFieldResolverResolveStaticRequestDataFields(t *testing.T) {
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
		"user": map[string]any{
			"id": "4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			"profile": map[string]any{
				"id":   "d13f60a4-5765-48c7-9e1d-3e782340f833",
				"name": "test",
			},
		},
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
		{"@request.user.id", false, `"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"`},
		{"@request.user.profile", false, `"{\"id\":\"d13f60a4-5765-48c7-9e1d-3e782340f833\",\"name\":\"test\"}"`},
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
