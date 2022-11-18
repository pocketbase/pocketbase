package resolvers_test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestRecordFieldResolverUpdateQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	authRecord, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestData := &models.RequestData{
		AuthRecord: authRecord,
	}

	scenarios := []struct {
		name               string
		collectionIdOrName string
		fields             []string
		allowHiddenFields  bool
		expectQuery        string
	}{
		{
			"missing field",
			"demo4",
			[]string{""},
			false,
			"SELECT `demo4`.* FROM `demo4`",
		},
		{
			"non relation field",
			"demo4",
			[]string{"title"},
			false,
			"SELECT `demo4`.* FROM `demo4`",
		},
		{
			"incomplete rel",
			"demo4",
			[]string{"self_rel_one"},
			false,
			"SELECT `demo4`.* FROM `demo4`",
		},
		{
			"single rel (self rel)",
			"demo4",
			[]string{"self_rel_one.title"},
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_one]]) THEN [[demo4.self_rel_one]] ELSE json_array([[demo4.self_rel_one]]) END) `demo4_self_rel_one_je` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4_self_rel_one_je.value]]",
		},
		{
			"single rel (other collection)",
			"demo4",
			[]string{"rel_one_cascade.title"},
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.rel_one_cascade]]) THEN [[demo4.rel_one_cascade]] ELSE json_array([[demo4.rel_one_cascade]]) END) `demo4_rel_one_cascade_je` LEFT JOIN `demo3` `demo4_rel_one_cascade` ON [[demo4_rel_one_cascade.id]] = [[demo4_rel_one_cascade_je.value]]",
		},
		{
			"non-relation field + single rel",
			"demo4",
			[]string{"title", "self_rel_one.title"},
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_one]]) THEN [[demo4.self_rel_one]] ELSE json_array([[demo4.self_rel_one]]) END) `demo4_self_rel_one_je` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4_self_rel_one_je.value]]",
		},
		{
			"nested incomplete rels",
			"demo4",
			[]string{"self_rel_many.self_rel_one"},
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]]",
		},
		{
			"nested complete rels",
			"demo4",
			[]string{"self_rel_many.self_rel_one.title"},
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_many.self_rel_one]]) THEN [[demo4_self_rel_many.self_rel_one]] ELSE json_array([[demo4_self_rel_many.self_rel_one]]) END) `demo4_self_rel_many_self_rel_one_je` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many_self_rel_one_je.value]]",
		},
		{
			"repeated nested rels",
			"demo4",
			[]string{"self_rel_many.self_rel_one.self_rel_many.self_rel_one.title"},
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_many.self_rel_one]]) THEN [[demo4_self_rel_many.self_rel_one]] ELSE json_array([[demo4_self_rel_many.self_rel_one]]) END) `demo4_self_rel_many_self_rel_one_je` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many_self_rel_one_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_many_self_rel_one.self_rel_many]]) THEN [[demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `demo4_self_rel_many_self_rel_one_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many` ON [[demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]]) THEN [[demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] ELSE json_array([[demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]]) END) `demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one_je` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one_je.value]]",
		},
		{
			"multiple rels",
			"demo4",
			[]string{"self_rel_many.title", "self_rel_one.onefile"},
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_one]]) THEN [[demo4.self_rel_one]] ELSE json_array([[demo4.self_rel_one]]) END) `demo4_self_rel_one_je` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4_self_rel_one_je.value]]",
		},
		{
			"@collection join",
			"demo4",
			[]string{"@collection.demo1.text", "@collection.demo2.active", "@collection.demo1.file_one"},
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo1` `__collection_demo1` LEFT JOIN `demo2` `__collection_demo2`",
		},
		{
			"@request.auth fields",
			"demo4",
			[]string{"@request.auth.id", "@request.auth.username", "@request.auth.rel.title", "@request.data.demo"},
			false,
			"^" +
				regexp.QuoteMeta("SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `users` `__auth_users` ON [[__auth_users.id]] =") +
				" {:.*} " +
				regexp.QuoteMeta("LEFT JOIN json_each(CASE WHEN json_valid([[__auth_users.rel]]) THEN [[__auth_users.rel]] ELSE json_array([[__auth_users.rel]]) END) `__auth_users_rel_je` LEFT JOIN `demo2` `__auth_users_rel` ON [[__auth_users_rel.id]] = [[__auth_users_rel_je.value]]") +
				"$",
		},
		{
			"hidden field with system filters (ignore emailVisibility)",
			"demo4",
			[]string{"@collection.users.email", "@request.auth.email"},
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `users` `__collection_users`",
		},
		{
			"hidden field (add emailVisibility)",
			"users",
			[]string{"email"},
			false,
			"SELECT `users`.* FROM `users` WHERE [[users.emailVisibility]] = TRUE",
		},
		{
			"hidden field (force ignore emailVisibility)",
			"users",
			[]string{"email"},
			true,
			"SELECT `users`.* FROM `users`",
		},
	}

	for _, s := range scenarios {
		collection, err := app.Dao().FindCollectionByNameOrId(s.collectionIdOrName)
		if err != nil {
			t.Errorf("[%s] Failed to load collection %s: %v", s.name, s.collectionIdOrName, err)
		}

		query := app.Dao().RecordQuery(collection)

		r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestData, s.allowHiddenFields)
		for _, field := range s.fields {
			r.Resolve(field)
		}

		if err := r.UpdateQuery(query); err != nil {
			t.Errorf("[%s] UpdateQuery failed with error %v", s.name, err)
			continue
		}

		rawQuery := query.Build().SQL()

		if !list.ExistInSliceWithRegex(rawQuery, []string{s.expectQuery}) {
			t.Errorf("[%s] Expected query\n %v \ngot:\n %v", s.name, s.expectQuery, rawQuery)
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

	authRecord, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestData := &models.RequestData{
		AuthRecord: authRecord,
	}

	r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestData, true)

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
		{"self_rel_many", false, "[[demo4.self_rel_many]]"},
		{"self_rel_many.", true, ""},
		{"self_rel_many.unknown", true, ""},
		{"self_rel_many.title", false, "[[demo4_self_rel_many.title]]"},
		{"self_rel_many.self_rel_one.self_rel_many.title", false, "[[demo4_self_rel_many_self_rel_one_self_rel_many.title]]"},
		// json_extract
		{"json_array.0", false, "JSON_EXTRACT([[demo4.json_array]], '$[0]')"},
		{"json_object.a.b.c", false, "JSON_EXTRACT([[demo4.json_object]], '$.a.b.c')"},
		// @request.auth relation join:
		{"@request.auth.rel", false, "[[__auth_users.rel]]"},
		{"@request.auth.rel.title", false, "[[__auth_users_rel.title]]"},
		// @collection fieds:
		{"@collect", true, ""},
		{"collection.demo4.title", true, ""},
		{"@collection", true, ""},
		{"@collection.unknown", true, ""},
		{"@collection.demo2", true, ""},
		{"@collection.demo2.", true, ""},
		{"@collection.demo2.title", false, "[[__collection_demo2.title]]"},
		{"@collection.demo4.title", false, "[[__collection_demo4.title]]"},
		{"@collection.demo4.id", false, "[[__collection_demo4.id]]"},
		{"@collection.demo4.created", false, "[[__collection_demo4.created]]"},
		{"@collection.demo4.updated", false, "[[__collection_demo4.updated]]"},
		{"@collection.demo4.self_rel_many.missing", true, ""},
		{"@collection.demo4.self_rel_many.self_rel_one.self_rel_many.self_rel_one.title", false, "[[__collection_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]]"},
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

	authRecord, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestData := &models.RequestData{
		Method: "get",
		Query: map[string]any{
			"a": 123,
		},
		Data: map[string]any{
			"b": 456,
			"c": map[string]int{"sub": 1},
		},
		AuthRecord: authRecord,
	}

	r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestData, true)

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
		{"@request.auth", true, ""},
		{"@request.auth.id", false, `"4q1xlclmfloku33"`},
		{"@request.auth.email", false, `"test@example.com"`},
		{"@request.auth.username", false, `"users75657"`},
		{"@request.auth.verified", false, `false`},
		{"@request.auth.emailVisibility", false, `false`},
		{"@request.auth.missing", false, `NULL`},
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
				t.Errorf("(%d) Expected 0 placeholder parameters for %v, got %v", i, name, params)
			}
			continue
		}

		// existing key
		// ---
		if len(params) != 1 {
			t.Errorf("(%d) Expected 1 placeholder parameter for %v, got %v", i, name, params)
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
			t.Errorf("(%d) Expected params %v for %v, got %v", i, s.expectParamValue, name, string(encodedParamValue))
		}
	}
}
