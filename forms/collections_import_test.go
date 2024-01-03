package forms_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestCollectionsImportValidate(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	totalCollections := 11

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
			expectCollectionsCount: totalCollections,
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
			expectCollectionsCount: totalCollections,
			expectEvents: map[string]int{
				"OnModelBeforeCreate": 2,
			},
		},
		{
			name: "test empty base collection schema",
			jsonData: `{
				"collections": [
					{
						"name": "import1"
					},
					{
						"name": "import2",
						"type": "auth"
					}
				]
			}`,
			expectError:            true,
			expectCollectionsCount: totalCollections,
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
					},
					{
						"name": "import3",
						"type": "auth"
					}
				]
			}`,
			expectError:            false,
			expectCollectionsCount: totalCollections + 3,
			expectEvents: map[string]int{
				"OnModelBeforeCreate": 3,
				"OnModelAfterCreate":  3,
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
			expectCollectionsCount: totalCollections,
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
						"id":"sz5l5z67tg7gku0",
						"name":"demo2",
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
			expectCollectionsCount: totalCollections,
			expectEvents: map[string]int{
				"OnModelBeforeDelete": 1,
			},
		},
		{
			name: "modified + new collection",
			jsonData: `{
				"collections": [
					{
						"id":"sz5l5z67tg7gku0",
						"name":"demo2_rename",
						"schema":[
							{
								"id":"_2hlxbmp",
								"name":"title_new",
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
			expectCollectionsCount: totalCollections + 2,
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
						"id": "kpv709sk2lqbqk8",
						"system": true,
						"name": "nologin",
						"type": "auth",
						"options": {
							"allowEmailAuth": false,
							"allowOAuth2Auth": false,
							"allowUsernameAuth": false,
							"exceptEmailDomains": [],
							"manageRule": "@request.auth.collectionName = 'users'",
							"minPasswordLength": 8,
							"onlyEmailDomains": [],
							"requireEmail": true
						},
						"listRule": "",
						"viewRule": "",
						"createRule": "",
						"updateRule": "",
						"deleteRule": "",
						"schema": [
							{
								"id": "x8zzktwe",
								"name": "name",
								"type": "text",
								"system": false,
								"required": false,
								"unique": false,
								"options": {
									"min": null,
									"max": null,
									"pattern": ""
								}
							}
						]
					},
					{
						"id":"sz5l5z67tg7gku0",
						"name":"demo2",
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
						"name": "demo1",
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
				"OnModelBeforeDelete": totalCollections - 2,
				"OnModelAfterDelete":  totalCollections - 2,
			},
		},
		{
			name: "lazy system table name error",
			jsonData: `{
				"collections": [
					{
						"name": "_admins",
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
			expectCollectionsCount: totalCollections,
			expectEvents: map[string]int{
				"OnModelBeforeCreate": 1,
			},
		},
		{
			name: "lazy view evaluation",
			jsonData: `{
				"collections": [
					{
						"name": "view_before",
						"type": "view",
						"options": {
							"query": "select id, active from base_test"
						}
					},
					{
						"name": "base_test",
						"schema": [
							{
								"id":"fz6iql2m",
								"name":"active",
								"type":"bool"
							}
						]
					},
					{
						"name": "view_after_new",
						"type": "view",
						"options": {
							"query": "select id, active from base_test"
						}
					},
					{
						"name": "view_after_old",
						"type": "view",
						"options": {
							"query": "select id from demo1"
						}
					}
				]
			}`,
			expectError:            false,
			expectCollectionsCount: totalCollections + 4,
			expectEvents: map[string]int{
				"OnModelBeforeUpdate": 3,
				"OnModelAfterUpdate":  3,
				"OnModelBeforeCreate": 4,
				"OnModelAfterCreate":  4,
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testApp, _ := tests.NewTestApp()
			defer testApp.Cleanup()

			form := forms.NewCollectionsImport(testApp)

			// load data
			loadErr := json.Unmarshal([]byte(s.jsonData), form)
			if loadErr != nil {
				t.Fatalf("Failed to load form data: %v", loadErr)
			}

			err := form.Submit()

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr to be %v, got %v (%v)", s.expectError, hasErr, err)
			}

			// check collections count
			collections := []*models.Collection{}
			if err := testApp.Dao().CollectionQuery().All(&collections); err != nil {
				t.Fatal(err)
			}
			if len(collections) != s.expectCollectionsCount {
				t.Fatalf("Expected %d collections, got %d", s.expectCollectionsCount, len(collections))
			}

			// check events
			if len(testApp.EventCalls) > len(s.expectEvents) {
				t.Fatalf("Expected events %v, got %v", s.expectEvents, testApp.EventCalls)
			}
			for event, expectedCalls := range s.expectEvents {
				actualCalls := testApp.EventCalls[event]
				if actualCalls != expectedCalls {
					t.Fatalf("Expected event %s to be called %d, got %d", event, expectedCalls, actualCalls)
				}
			}
		})
	}
}

func TestCollectionsImportSubmitInterceptors(t *testing.T) {
	t.Parallel()

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
	interceptor1 := func(next forms.InterceptorNextFunc[[]*models.Collection]) forms.InterceptorNextFunc[[]*models.Collection] {
		return func(imports []*models.Collection) error {
			interceptor1Called = true
			return next(imports)
		}
	}

	interceptor2Called := false
	interceptor2 := func(next forms.InterceptorNextFunc[[]*models.Collection]) forms.InterceptorNextFunc[[]*models.Collection] {
		return func(imports []*models.Collection) error {
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
