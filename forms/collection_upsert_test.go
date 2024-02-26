package forms_test

import (
	"encoding/json"
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

func TestNewCollectionUpsert(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := &models.Collection{}
	collection.Name = "test_name"
	collection.Type = "test_type"
	collection.System = true
	listRule := "test_list"
	collection.ListRule = &listRule
	viewRule := "test_view"
	collection.ViewRule = &viewRule
	createRule := "test_create"
	collection.CreateRule = &createRule
	updateRule := "test_update"
	collection.UpdateRule = &updateRule
	deleteRule := "test_delete"
	collection.DeleteRule = &deleteRule
	collection.Schema = schema.NewSchema(&schema.SchemaField{
		Name: "test",
		Type: schema.FieldTypeText,
	})

	form := forms.NewCollectionUpsert(app, collection)

	if form.Name != collection.Name {
		t.Errorf("Expected Name %q, got %q", collection.Name, form.Name)
	}

	if form.Type != collection.Type {
		t.Errorf("Expected Type %q, got %q", collection.Type, form.Type)
	}

	if form.System != collection.System {
		t.Errorf("Expected System %v, got %v", collection.System, form.System)
	}

	if form.ListRule != collection.ListRule {
		t.Errorf("Expected ListRule %v, got %v", collection.ListRule, form.ListRule)
	}

	if form.ViewRule != collection.ViewRule {
		t.Errorf("Expected ViewRule %v, got %v", collection.ViewRule, form.ViewRule)
	}

	if form.CreateRule != collection.CreateRule {
		t.Errorf("Expected CreateRule %v, got %v", collection.CreateRule, form.CreateRule)
	}

	if form.UpdateRule != collection.UpdateRule {
		t.Errorf("Expected UpdateRule %v, got %v", collection.UpdateRule, form.UpdateRule)
	}

	if form.DeleteRule != collection.DeleteRule {
		t.Errorf("Expected DeleteRule %v, got %v", collection.DeleteRule, form.DeleteRule)
	}

	// store previous state and modify the collection schema to verify
	// that the form.Schema is a deep clone
	loadedSchema, _ := collection.Schema.MarshalJSON()
	collection.Schema.AddField(&schema.SchemaField{
		Name: "new_field",
		Type: schema.FieldTypeBool,
	})

	formSchema, _ := form.Schema.MarshalJSON()

	if string(formSchema) != string(loadedSchema) {
		t.Errorf("Expected Schema %v, got %v", string(loadedSchema), string(formSchema))
	}
}

func TestCollectionUpsertValidateAndSubmit(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		testName       string
		existingName   string
		jsonData       string
		expectedErrors []string
	}{
		{"empty create (base)", "", "{}", []string{"name", "schema"}},
		{"empty create (auth)", "", `{"type":"auth"}`, []string{"name"}},
		{"empty create (view)", "", `{"type":"view"}`, []string{"name", "options"}},
		{"empty update", "demo2", "{}", []string{}},
		{
			"collection and field with _via_ names",
			"",
			`{
				"name": "a_via_b",
				"schema": [
					{"name":"c_via_d","type":"text"}
				]
			}`,
			[]string{"name", "schema"},
		},
		{
			"create failure",
			"",
			`{
				"name": "test ?!@#$",
				"type": "invalid",
				"system": true,
				"schema": [
					{"name":"","type":"text"}
				],
				"listRule": "missing = '123'",
				"viewRule": "missing = '123'",
				"createRule": "missing = '123'",
				"updateRule": "missing = '123'",
				"deleteRule": "missing = '123'",
				"indexes": ["create index '' on '' ()"]
			}`,
			[]string{"name", "type", "schema", "listRule", "viewRule", "createRule", "updateRule", "deleteRule", "indexes"},
		},
		{
			"create failure - existing name",
			"",
			`{
				"name": "demo1",
				"system": true,
				"schema": [
					{"name":"test","type":"text"}
				],
				"listRule": "test='123'",
				"viewRule": "test='123'",
				"createRule": "test='123'",
				"updateRule": "test='123'",
				"deleteRule": "test='123'"
			}`,
			[]string{"name"},
		},
		{
			"create failure - existing internal table",
			"",
			`{
				"name": "_admins",
				"schema": [
					{"name":"test","type":"text"}
				]
			}`,
			[]string{"name"},
		},
		{
			"create failure - name starting with underscore",
			"",
			`{
				"name": "_test_new",
				"schema": [
					{"name":"test","type":"text"}
				]
			}`,
			[]string{"name"},
		},
		{
			"create failure - duplicated field names (case insensitive)",
			"",
			`{
				"name": "test_new",
				"schema": [
					{"name":"test","type":"text"},
					{"name":"tESt","type":"text"}
				]
			}`,
			[]string{"schema"},
		},
		{
			"create failure - check auth options validators",
			"",
			`{
				"name": "test_new",
				"type": "auth",
				"schema": [
					{"name":"test","type":"text"}
				],
				"options": { "minPasswordLength": 3 }
			}`,
			[]string{"options"},
		},
		{
			"create failure - check view options validators",
			"",
			`{
				"name": "test_new",
				"type": "view",
				"options": { "query": "invalid query" }
			}`,
			[]string{"options"},
		},
		{
			"create success",
			"",
			`{
				"name": "test_new",
				"type": "auth",
				"system": true,
				"schema": [
					{"id":"a123456","name":"test1","type":"text"},
					{"id":"b123456","name":"test2","type":"email"},
					{
						"name":"test3",
						"type":"relation",
						"options":{
							"collectionId":"v851q4r790rhknl",
							"displayFields":["name","id","created","updated","username","email","emailVisibility","verified"]
						}
					}
				],
				"listRule": "test1='123' && verified = true",
				"viewRule": "test1='123' && emailVisibility = true",
				"createRule": "test1='123' && email != ''",
				"updateRule": "test1='123' && username != ''",
				"deleteRule": "test1='123' && id != ''",
				"indexes": ["create index idx_test_new on anything (test1)"]
			}`,
			[]string{},
		},
		{
			"update failure - changing field type",
			"test_new",
			`{
				"schema": [
					{"id":"a123456","name":"test1","type":"url"},
					{"id":"b123456","name":"test2","type":"bool"}
				],
				"indexes": ["create index idx_test_new on test_new (test1)", "invalid"]
			}`,
			[]string{"schema", "indexes"},
		},
		{
			"update success - rename fields to existing field names (aka. reusing field names)",
			"test_new",
			`{
				"schema": [
					{"id":"a123456","name":"test2","type":"text"},
					{"id":"b123456","name":"test1","type":"email"}
				]
			}`,
			[]string{},
		},
		{
			"update failure - existing name",
			"demo2",
			`{"name": "demo3"}`,
			[]string{"name"},
		},
		{
			"update failure - changing system collection",
			"nologin",
			`{
				"name": "update",
				"system": false,
				"schema": [
					{"id":"koih1lqx","name":"abc","type":"text"}
				],
				"listRule": "abc = '123'",
				"viewRule": "abc = '123'",
				"createRule": "abc = '123'",
				"updateRule": "abc = '123'",
				"deleteRule": "abc = '123'"
			}`,
			[]string{"name", "system"},
		},
		{
			"update failure - changing collection type",
			"demo3",
			`{
				"type": "auth"
			}`,
			[]string{"type"},
		},
		{
			"update failure - changing relation collection",
			"users",
			`{
				"schema": [
					{
						"id": "lkeigvv3",
						"name": "rel",
						"type": "relation",
						"options": {
							"collectionId": "wzlqyes4orhoygb",
							"cascadeDelete": false,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				]
			}`,
			[]string{"schema"},
		},
		{
			"update failure - all fields",
			"demo2",
			`{
				"name": "test ?!@#$",
				"type": "invalid",
				"system": true,
				"schema": [
					{"name":"","type":"text"}
				],
				"listRule": "missing = '123'",
				"viewRule": "missing = '123'",
				"createRule": "missing = '123'",
				"updateRule": "missing = '123'",
				"deleteRule": "missing = '123'",
				"options": {"test": 123},
				"indexes": ["create index '' from demo2 on (id)"]
			}`,
			[]string{"name", "type", "system", "schema", "listRule", "viewRule", "createRule", "updateRule", "deleteRule", "indexes"},
		},
		{
			"update success - update all fields",
			"clients",
			`{
				"name": "demo_update",
				"type": "auth",
				"schema": [
					{"id":"_2hlxbmp","name":"test","type":"text"}
				],
				"listRule": "test='123' && verified = true",
				"viewRule": "test='123' && emailVisibility = true",
				"createRule": "test='123' && email != ''",
				"updateRule": "test='123' && username != ''",
				"deleteRule": "test='123' && id != ''",
				"options": {"minPasswordLength": 10},
				"indexes": [
					"create index idx_clients_test1 on anything (id, email, test)",
					"create unique index idx_clients_test2 on clients (id, username, email)"
				]
			}`,
			[]string{},
		},
		// (fail due to filters old field references)
		{
			"update failure - rename the schema field of the last updated collection",
			"demo_update",
			`{
				"schema": [
					{"id":"_2hlxbmp","name":"test_renamed","type":"text"}
				]
			}`,
			[]string{"listRule", "viewRule", "createRule", "updateRule", "deleteRule"},
		},
		// (cleared filter references)
		{
			"update success - rename the schema field of the last updated collection",
			"demo_update",
			`{
				"schema": [
					{"id":"_2hlxbmp","name":"test_renamed","type":"text"}
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"indexes": []
			}`,
			[]string{},
		},
		{
			"update success - system collection",
			"nologin",
			`{
				"listRule": "name='123'",
				"viewRule": "name='123'",
				"createRule": "name='123'",
				"updateRule": "name='123'",
				"deleteRule": "name='123'"
			}`,
			[]string{},
		},

		// view tests
		// -----------------------------------------------------------
		{
			"base->view relation",
			"",
			`{
				"name": "test_view_relation",
				"type": "base",
				"schema": [
					{
						"name": "test",
						"type": "relation",
						"options":{
							"collectionId": "v9gwnfh02gjq1q0"
						}
					}
				]
			}`,
			[]string{"schema"}, // not allowed
		},
		{
			"auth->view relation",
			"",
			`{
				"name": "test_view_relation",
				"type": "auth",
				"schema": [
					{
						"name": "test",
						"type": "relation",
						"options": {
							"collectionId": "v9gwnfh02gjq1q0"
						}
					}
				]
			}`,
			[]string{"schema"}, // not allowed
		},
		{
			"view->view relation",
			"",
			`{
				"name": "test_view_relation",
				"type": "view",
				"options": {
					"query": "select view1.id, view1.id as rel from view1"
				}
			}`,
			[]string{}, // allowed
		},
		{
			"view create failure",
			"",
			`{
				"name": "upsert_view",
				"type": "view",
				"listRule": "id='123' && verified = true",
				"viewRule": "id='123' && emailVisibility = true",
				"schema": [
					{"id":"abc123","name":"some invalid field name that will be overwritten !@#$","type":"bool"}
				],
				"options": {
					"query": "select id, email from users; drop table _admins;"
				},
				"indexes": ["create index idx_test_view on upsert_view (id)"]
			}`,
			[]string{
				"listRule",
				"viewRule",
				"options",
				"indexes", // views don't have indexes
			},
		},
		{
			"view create success",
			"",
			`{
				"name": "upsert_view",
				"type": "view",
				"listRule": "id='123' && verified = true",
				"viewRule": "id='123' && emailVisibility = true",
				"schema": [
					{"id":"abc123","name":"some invalid field name that will be overwritten !@#$","type":"bool"}
				],
				"options": {
					"query": "select id, emailVisibility, verified from users"
				}
			}`,
			[]string{
				// "schema", should be overwritten by an autogenerated from the query
			},
		},
		{
			"view update failure (schema autogeneration and rule fields check)",
			"upsert_view",
			`{
				"name": "upsert_view_2",
				"listRule": "id='456' && verified = true",
				"viewRule": "id='456'",
				"createRule": "id='123'",
				"updateRule": "id='123'",
				"deleteRule": "id='123'",
				"schema": [
					{"id":"abc123","name":"verified","type":"bool"}
				],
				"options": {
					"query": "select 1 as id"
				}
			}`,
			[]string{
				"listRule",   // missing field (ignoring the old or explicit schema)
				"createRule", // not allowed
				"updateRule", // not allowed
				"deleteRule", // not allowed
			},
		},
		{
			"view update failure (check query identifiers format)",
			"upsert_view",
			`{
				"listRule": null,
				"viewRule": null,
				"options": {
					"query": "select 1 as id, 2 as [invalid!@#]"
				}
			}`,
			[]string{
				"schema", // should fail due to invalid field name
			},
		},
		{
			"view update success",
			"upsert_view",
			`{
				"listRule": null,
				"viewRule": null,
				"options": {
					"query": "select 1 as id, 2 as valid"
				}
			}`,
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			collection := &models.Collection{}
			if s.existingName != "" {
				var err error
				collection, err = app.Dao().FindCollectionByNameOrId(s.existingName)
				if err != nil {
					t.Fatal(err)
				}
			}

			form := forms.NewCollectionUpsert(app, collection)

			// load data
			loadErr := json.Unmarshal([]byte(s.jsonData), form)
			if loadErr != nil {
				t.Fatalf("Failed to load form data: %v", loadErr)
			}

			interceptorCalls := 0
			interceptor := func(next forms.InterceptorNextFunc[*models.Collection]) forms.InterceptorNextFunc[*models.Collection] {
				return func(c *models.Collection) error {
					interceptorCalls++
					return next(c)
				}
			}

			// parse errors
			result := form.Submit(interceptor)
			errs, ok := result.(validation.Errors)
			if !ok && result != nil {
				t.Fatalf("Failed to parse errors %v", result)
			}

			// check interceptor calls
			expectInterceptorCalls := 1
			if len(s.expectedErrors) > 0 {
				expectInterceptorCalls = 0
			}
			if interceptorCalls != expectInterceptorCalls {
				t.Fatalf("Expected interceptor to be called %d, got %d", expectInterceptorCalls, interceptorCalls)
			}

			// check errors
			if len(errs) > len(s.expectedErrors) {
				t.Fatalf("Expected error keys %v, got %v", s.expectedErrors, errs)
			}
			for _, k := range s.expectedErrors {
				if _, ok := errs[k]; !ok {
					t.Fatalf("Missing expected error key %q in %v", k, errs)
				}
			}

			if len(s.expectedErrors) > 0 {
				return
			}

			collection, _ = app.Dao().FindCollectionByNameOrId(form.Name)
			if collection == nil {
				t.Fatalf("Expected to find collection %q, got nil", form.Name)
			}

			if form.Name != collection.Name {
				t.Fatalf("Expected Name %q, got %q", collection.Name, form.Name)
			}

			if form.Type != collection.Type {
				t.Fatalf("Expected Type %q, got %q", collection.Type, form.Type)
			}

			if form.System != collection.System {
				t.Fatalf("Expected System %v, got %v", collection.System, form.System)
			}

			if cast.ToString(form.ListRule) != cast.ToString(collection.ListRule) {
				t.Fatalf("Expected ListRule %v, got %v", collection.ListRule, form.ListRule)
			}

			if cast.ToString(form.ViewRule) != cast.ToString(collection.ViewRule) {
				t.Fatalf("Expected ViewRule %v, got %v", collection.ViewRule, form.ViewRule)
			}

			if cast.ToString(form.CreateRule) != cast.ToString(collection.CreateRule) {
				t.Fatalf("Expected CreateRule %v, got %v", collection.CreateRule, form.CreateRule)
			}

			if cast.ToString(form.UpdateRule) != cast.ToString(collection.UpdateRule) {
				t.Fatalf("Expected UpdateRule %v, got %v", collection.UpdateRule, form.UpdateRule)
			}

			if cast.ToString(form.DeleteRule) != cast.ToString(collection.DeleteRule) {
				t.Fatalf("Expected DeleteRule %v, got %v", collection.DeleteRule, form.DeleteRule)
			}

			rawFormSchema, _ := form.Schema.MarshalJSON()
			rawCollectionSchema, _ := collection.Schema.MarshalJSON()

			if len(form.Schema.Fields()) != len(collection.Schema.Fields()) {
				t.Fatalf("Expected Schema \n%v, \ngot \n%v", string(rawCollectionSchema), string(rawFormSchema))
			}

			for _, f := range form.Schema.Fields() {
				if collection.Schema.GetFieldByName(f.Name) == nil {
					t.Fatalf("Missing field %s \nin \n%v", f.Name, string(rawFormSchema))
				}
			}

			// check indexes (if any)
			allIndexes, _ := app.Dao().TableIndexes(form.Name)
			for _, formIdx := range form.Indexes {
				parsed := dbutils.ParseIndex(formIdx)
				parsed.TableName = form.Name
				normalizedIdx := parsed.Build()

				var exists bool
				for _, idx := range allIndexes {
					if dbutils.ParseIndex(idx).Build() == normalizedIdx {
						exists = true
						continue
					}
				}

				if !exists {
					t.Fatalf("Missing index %s \nin \n%v", normalizedIdx, allIndexes)
				}
			}
		})
	}
}

func TestCollectionUpsertSubmitInterceptors(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}

	form := forms.NewCollectionUpsert(app, collection)
	form.Name = "test_new"

	testErr := errors.New("test_error")
	interceptorCollectionName := ""

	interceptor1Called := false
	interceptor1 := func(next forms.InterceptorNextFunc[*models.Collection]) forms.InterceptorNextFunc[*models.Collection] {
		return func(c *models.Collection) error {
			interceptor1Called = true
			return next(c)
		}
	}

	interceptor2Called := false
	interceptor2 := func(next forms.InterceptorNextFunc[*models.Collection]) forms.InterceptorNextFunc[*models.Collection] {
		return func(c *models.Collection) error {
			interceptorCollectionName = collection.Name // to check if the record was filled
			interceptor2Called = true
			return testErr
		}
	}

	submitErr := form.Submit(interceptor1, interceptor2)
	if submitErr != testErr {
		t.Fatalf("Expected submitError %v, got %v", testErr, submitErr)
	}

	if !interceptor1Called {
		t.Fatalf("Expected interceptor1 to be called")
	}

	if !interceptor2Called {
		t.Fatalf("Expected interceptor2 to be called")
	}

	if interceptorCollectionName != form.Name {
		t.Fatalf("Expected the form model to be filled before calling the interceptors")
	}
}

func TestCollectionUpsertWithCustomId(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	existingCollection, err := app.Dao().FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}

	newCollection := func() *models.Collection {
		return &models.Collection{
			Name:   "c_" + security.PseudorandomString(4),
			Schema: existingCollection.Schema,
		}
	}

	scenarios := []struct {
		name        string
		jsonData    string
		collection  *models.Collection
		expectError bool
	}{
		{
			"empty data",
			"{}",
			newCollection(),
			false,
		},
		{
			"empty id",
			`{"id":""}`,
			newCollection(),
			false,
		},
		{
			"id < 15 chars",
			`{"id":"a23"}`,
			newCollection(),
			true,
		},
		{
			"id > 15 chars",
			`{"id":"a234567890123456"}`,
			newCollection(),
			true,
		},
		{
			"id = 15 chars (invalid chars)",
			`{"id":"a@3456789012345"}`,
			newCollection(),
			true,
		},
		{
			"id = 15 chars (valid chars)",
			`{"id":"a23456789012345"}`,
			newCollection(),
			false,
		},
		{
			"changing the id of an existing item",
			`{"id":"b23456789012345"}`,
			existingCollection,
			true,
		},
		{
			"using the same existing item id",
			`{"id":"` + existingCollection.Id + `"}`,
			existingCollection,
			false,
		},
		{
			"skipping the id for existing item",
			`{}`,
			existingCollection,
			false,
		},
	}

	for _, s := range scenarios {
		form := forms.NewCollectionUpsert(app, s.collection)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("[%s] Failed to load form data: %v", s.name, loadErr)
			continue
		}

		submitErr := form.Submit()
		hasErr := submitErr != nil

		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr to be %v, got %v (%v)", s.name, s.expectError, hasErr, submitErr)
		}

		if !hasErr && form.Id != "" {
			_, err := app.Dao().FindCollectionByNameOrId(form.Id)
			if err != nil {
				t.Errorf("[%s] Expected to find record with id %s, got %v", s.name, form.Id, err)
			}
		}
	}
}
