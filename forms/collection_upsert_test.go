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
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

func TestNewCollectionUpsert(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := &models.Collection{}
	collection.Name = "test_name"
	collection.Type = "test_type"
	collection.System = true
	listRule := "testview"
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
		{"empty update", "demo2", "{}", []string{}},
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
				"deleteRule": "missing = '123'"
			}`,
			[]string{"name", "type", "schema", "listRule", "viewRule", "createRule", "updateRule", "deleteRule"},
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
			"create failure - check type options validators",
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
			"create success",
			"",
			`{
				"name": "test_new",
				"type": "auth",
				"system": true,
				"schema": [
					{"id":"a123456","name":"test1","type":"text"},
					{"id":"b123456","name":"test2","type":"email"}
				],
				"listRule": "test1='123' && verified = true",
				"viewRule": "test1='123' && emailVisibility = true",
				"createRule": "test1='123' && email != ''",
				"updateRule": "test1='123' && username != ''",
				"deleteRule": "test1='123' && id != ''"
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
				]
			}`,
			[]string{"schema"},
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
				"options": {"test": 123}
			}`,
			[]string{"name", "type", "system", "schema", "listRule", "viewRule", "createRule", "updateRule", "deleteRule"},
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
				"options": {"minPasswordLength": 10}
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
				"deleteRule": null
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
	}

	for _, s := range scenarios {
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
			t.Errorf("[%s] Failed to load form data: %v", s.testName, loadErr)
			continue
		}

		interceptorCalls := 0
		interceptor := func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
			return func() error {
				interceptorCalls++
				return next()
			}
		}

		// parse errors
		result := form.Submit(interceptor)
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("[%s] Failed to parse errors %v", s.testName, result)
			continue
		}

		// check interceptor calls
		expectInterceptorCalls := 1
		if len(s.expectedErrors) > 0 {
			expectInterceptorCalls = 0
		}
		if interceptorCalls != expectInterceptorCalls {
			t.Errorf("[%s] Expected interceptor to be called %d, got %d", s.testName, expectInterceptorCalls, interceptorCalls)
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("[%s] Expected error keys %v, got %v", s.testName, s.expectedErrors, errs)
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("[%s] Missing expected error key %q in %v", s.testName, k, errs)
			}
		}

		if len(s.expectedErrors) > 0 {
			continue
		}

		collection, _ = app.Dao().FindCollectionByNameOrId(form.Name)
		if collection == nil {
			t.Errorf("[%s] Expected to find collection %q, got nil", s.testName, form.Name)
			continue
		}

		if form.Name != collection.Name {
			t.Errorf("[%s] Expected Name %q, got %q", s.testName, collection.Name, form.Name)
		}

		if form.Type != collection.Type {
			t.Errorf("[%s] Expected Type %q, got %q", s.testName, collection.Type, form.Type)
		}

		if form.System != collection.System {
			t.Errorf("[%s] Expected System %v, got %v", s.testName, collection.System, form.System)
		}

		if cast.ToString(form.ListRule) != cast.ToString(collection.ListRule) {
			t.Errorf("[%s] Expected ListRule %v, got %v", s.testName, collection.ListRule, form.ListRule)
		}

		if cast.ToString(form.ViewRule) != cast.ToString(collection.ViewRule) {
			t.Errorf("[%s] Expected ViewRule %v, got %v", s.testName, collection.ViewRule, form.ViewRule)
		}

		if cast.ToString(form.CreateRule) != cast.ToString(collection.CreateRule) {
			t.Errorf("[%s] Expected CreateRule %v, got %v", s.testName, collection.CreateRule, form.CreateRule)
		}

		if cast.ToString(form.UpdateRule) != cast.ToString(collection.UpdateRule) {
			t.Errorf("[%s] Expected UpdateRule %v, got %v", s.testName, collection.UpdateRule, form.UpdateRule)
		}

		if cast.ToString(form.DeleteRule) != cast.ToString(collection.DeleteRule) {
			t.Errorf("[%s] Expected DeleteRule %v, got %v", s.testName, collection.DeleteRule, form.DeleteRule)
		}

		formSchema, _ := form.Schema.MarshalJSON()
		collectionSchema, _ := collection.Schema.MarshalJSON()
		if string(formSchema) != string(collectionSchema) {
			t.Errorf("[%s] Expected Schema %v, got %v", s.testName, string(collectionSchema), string(formSchema))
		}
	}
}

func TestCollectionUpsertSubmitInterceptors(t *testing.T) {
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
	interceptor1 := func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			interceptor1Called = true
			return next()
		}
	}

	interceptor2Called := false
	interceptor2 := func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
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
