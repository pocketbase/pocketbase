package forms_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestCollectionsImportPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewCollectionsImport(nil)
}

func TestCollectionsImportValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	form := forms.NewCollectionsImport(app)

	scenarios := []struct {
		collections []*models.Collection
		expectError bool
	}{
		{nil, true},
		{[]*models.Collection{}, true},
		{[]*models.Collection{{}}, false},
	}

	for i, s := range scenarios {
		form.Collections = s.collections

		err := form.Validate()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}
	}
}

func TestCollectionsImportSubmit(t *testing.T) {
	scenarios := []struct {
		name                   string
		jsonData               string
		expectError            bool
		expectCollectionsCount int
		expectEvents           map[string]int
	}{
		{
			name: "empty collections",
			jsonData: `{
				"deleteMissing": true,
				"collections": []
			}`,
			expectError:            true,
			expectCollectionsCount: 5,
			expectEvents:           nil,
		},
		{
			name: "one of the collections has invalid data",
			jsonData: `{
				"collections": [
					{
						"name": "import1",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					},
					{
						"name": "import 2",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					}
				]
			}`,
			expectError:            true,
			expectCollectionsCount: 5,
			expectEvents: map[string]int{
				"OnModelBeforeCreate": 2,
			},
		},
		{
			name: "all imported collections has valid data",
			jsonData: `{
				"collections": [
					{
						"name": "import1",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					},
					{
						"name": "import2",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					}
				]
			}`,
			expectError:            false,
			expectCollectionsCount: 7,
			expectEvents: map[string]int{
				"OnModelBeforeCreate": 2,
				"OnModelAfterCreate":  2,
			},
		},
		{
			name: "new collection with existing name",
			jsonData: `{
				"collections": [
					{
						"name": "demo2",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					}
				]
			}`,
			expectError:            true,
			expectCollectionsCount: 5,
			expectEvents: map[string]int{
				"OnModelBeforeCreate": 1,
			},
		},
		{
			name: "delete system + modified + new collection",
			jsonData: `{
				"deleteMissing": true,
				"collections": [
					{
						"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc",
						"name":"demo",
						"schema":[
							{
								"id":"_2hlxbmp",
								"name":"title",
								"type":"text",
								"system":false,
								"required":true,
								"unique":false,
								"options":{
									"min":3,
									"max":null,
									"pattern":""
								}
							}
						]
					},
					{
						"name": "import1",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					}
				]
			}`,
			expectError:            true,
			expectCollectionsCount: 5,
		},
		{
			name: "modified + new collection",
			jsonData: `{
				"collections": [
					{
						"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc",
						"name":"demo",
						"schema":[
							{
								"id":"_2hlxbmp",
								"name":"title",
								"type":"text",
								"system":false,
								"required":true,
								"unique":false,
								"options":{
									"min":3,
									"max":null,
									"pattern":""
								}
							}
						]
					},
					{
						"name": "import1",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					},
					{
						"name": "import2",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					}
				]
			}`,
			expectError:            false,
			expectCollectionsCount: 7,
			expectEvents: map[string]int{
				"OnModelBeforeUpdate": 1,
				"OnModelAfterUpdate":  1,
				"OnModelBeforeCreate": 2,
				"OnModelAfterCreate":  2,
			},
		},
		{
			name: "delete non-system + modified + new collection",
			jsonData: `{
				"deleteMissing": true,
				"collections": [
					{
						"id":"abe78266-fd4d-4aea-962d-8c0138ac522b",
						"name":"profiles",
						"system":true,
						"listRule":"userId = @request.user.id",
						"viewRule":"created > 'test_change'",
						"createRule":"userId = @request.user.id",
						"updateRule":"userId = @request.user.id",
						"deleteRule":"userId = @request.user.id",
						"schema":[
							{
								"id":"koih1lqx",
								"name":"userId",
								"type":"user",
								"system":true,
								"required":true,
								"unique":true,
								"options":{
									"maxSelect":1,
									"cascadeDelete":true
								}
							},
							{
								"id":"69ycbg3q",
								"name":"rel",
								"type":"relation",
								"system":false,
								"required":false,
								"unique":false,
								"options":{
									"maxSelect":2,
									"collectionId":"abe78266-fd4d-4aea-962d-8c0138ac522b",
									"cascadeDelete":false
								}
							}
						]
					},
					{
						"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc",
						"name":"demo",
						"schema":[
							{
								"id":"_2hlxbmp",
								"name":"title",
								"type":"text",
								"system":false,
								"required":true,
								"unique":false,
								"options":{
									"min":3,
									"max":null,
									"pattern":""
								}
							}
						]
					},
					{
						"id": "test_deleted_collection_name_reuse",
						"name": "demo2",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					}
				]
			}`,
			expectError:            false,
			expectCollectionsCount: 3,
			expectEvents: map[string]int{
				"OnModelBeforeUpdate": 2,
				"OnModelAfterUpdate":  2,
				"OnModelBeforeCreate": 1,
				"OnModelAfterCreate":  1,
				"OnModelBeforeDelete": 3,
				"OnModelAfterDelete":  3,
			},
		},
	}

	for _, s := range scenarios {
		testApp, _ := tests.NewTestApp()
		defer testApp.Cleanup()

		form := forms.NewCollectionsImport(testApp)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("[%s] Failed to load form data: %v", s.name, loadErr)
			continue
		}

		err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr to be %v, got %v (%v)", s.name, s.expectError, hasErr, err)
		}

		// check collections count
		collections := []*models.Collection{}
		if err := testApp.Dao().CollectionQuery().All(&collections); err != nil {
			t.Fatal(err)
		}
		if len(collections) != s.expectCollectionsCount {
			t.Errorf("[%s] Expected %d collections, got %d", s.name, s.expectCollectionsCount, len(collections))
		}

		// check events
		if len(testApp.EventCalls) > len(s.expectEvents) {
			t.Errorf("[%s] Expected events %v, got %v", s.name, s.expectEvents, testApp.EventCalls)
		}
		for event, expectedCalls := range s.expectEvents {
			actualCalls := testApp.EventCalls[event]
			if actualCalls != expectedCalls {
				t.Errorf("[%s] Expected event %s to be called %d, got %d", s.name, event, expectedCalls, actualCalls)
			}
		}
	}
}

func TestCollectionsImportSubmitInterceptors(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collections := []*models.Collection{}
	if err := app.Dao().CollectionQuery().All(&collections); err != nil {
		t.Fatal(err)
	}

	form := forms.NewCollectionsImport(app)
	form.Collections = collections

	testErr := errors.New("test_error")

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
}
