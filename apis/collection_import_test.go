package apis_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestCollectionsImport(t *testing.T) {
	t.Parallel()

	totalCollections := 16

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPut,
			URL:             "/api/collections/import",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodPut,
			URL:    "/api/collections/import",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + empty collections",
			Method: http.MethodPut,
			URL:    "/api/collections/import",
			Body:   strings.NewReader(`{"collections":[]}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"collections":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				collections := []*core.Collection{}
				if err := app.CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				expected := totalCollections
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}
			},
		},
		{
			Name:   "authorized as superuser + collections validator failure",
			Method: http.MethodPut,
			URL:    "/api/collections/import",
			Body: strings.NewReader(`{
				"collections":[
					{"name": "import1"},
					{
						"name": "import2",
						"fields": [
							{
							  "id": "koih1lqx",
							  "name": "expand",
							  "type": "text"
							}
						]
					}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"collections":{"code":"validation_collections_import_failure"`,
				`import2`,
				`fields`,
			},
			NotExpectedContent: []string{"Raw error:"},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionsImportRequest":   1,
				"OnCollectionCreate":           2,
				"OnCollectionCreateExecute":    2,
				"OnCollectionAfterCreateError": 2,
				"OnModelCreate":                2,
				"OnModelCreateExecute":         2,
				"OnModelAfterCreateError":      2,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				collections := []*core.Collection{}
				if err := app.CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				expected := totalCollections
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}
			},
		},
		{
			Name:   "authorized as superuser + non-validator failure",
			Method: http.MethodPut,
			URL:    "/api/collections/import",
			Body: strings.NewReader(`{
				"collections":[
					{
						"name": "import1",
						"fields": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					},
					{
						"name": "import2",
						"fields": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						],
						"indexes": [
							"create index idx_test on import2 (test)"
						]
					}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"collections":{"code":"validation_collections_import_failure"`,
				`Raw error:`,
				`custom_error`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionsImportRequest":   1,
				"OnCollectionCreate":           1,
				"OnCollectionAfterCreateError": 1,
				"OnModelCreate":                1,
				"OnModelAfterCreateError":      1,
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.OnCollectionCreate().BindFunc(func(e *core.CollectionEvent) error {
					return errors.New("custom_error")
				})
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				collections := []*core.Collection{}
				if err := app.CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				expected := totalCollections
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}
			},
		},
		{
			Name:   "authorized as superuser + successful collections create",
			Method: http.MethodPut,
			URL:    "/api/collections/import",
			Body: strings.NewReader(`{
				"collections":[
					{
						"name": "import1",
						"fields": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					},
					{
						"name": "import2",
						"fields": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						],
						"indexes": [
							"create index idx_test on import2 (test)"
						]
					},
					{
						"name": "auth_without_fields",
						"type": "auth"
					}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionsImportRequest":     1,
				"OnCollectionCreate":             3,
				"OnCollectionCreateExecute":      3,
				"OnCollectionAfterCreateSuccess": 3,
				"OnModelCreate":                  3,
				"OnModelCreateExecute":           3,
				"OnModelAfterCreateSuccess":      3,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				collections := []*core.Collection{}
				if err := app.CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}

				expected := totalCollections + 3
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}

				indexes, err := app.TableIndexes("import2")
				if err != nil || indexes["idx_test"] == "" {
					t.Fatalf("Missing index %s (%v)", "idx_test", err)
				}
			},
		},
		{
			Name:   "authorized as superuser + create/update/delete",
			Method: http.MethodPut,
			URL:    "/api/collections/import",
			Body: strings.NewReader(`{
				"deleteMissing": true,
				"collections":[
					{"name": "test123"},
					{
						"id":"wsmn24bux7wo113",
						"name":"demo1",
						"fields":[
							{
								"id":"_2hlxbmp",
								"name":"title",
								"type":"text",
								"required":true
							}
						],
						"indexes": []
					}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                          0,
				"OnCollectionsImportRequest": 1,
				// ---
				"OnModelCreate":                  1,
				"OnModelCreateExecute":           1,
				"OnModelAfterCreateSuccess":      1,
				"OnCollectionCreate":             1,
				"OnCollectionCreateExecute":      1,
				"OnCollectionAfterCreateSuccess": 1,
				// ---
				"OnModelUpdate":                  1,
				"OnModelUpdateExecute":           1,
				"OnModelAfterUpdateSuccess":      1,
				"OnCollectionUpdate":             1,
				"OnCollectionUpdateExecute":      1,
				"OnCollectionAfterUpdateSuccess": 1,
				// ---
				"OnModelDelete":                  14,
				"OnModelAfterDeleteSuccess":      14,
				"OnModelDeleteExecute":           14,
				"OnCollectionDelete":             9,
				"OnCollectionDeleteExecute":      9,
				"OnCollectionAfterDeleteSuccess": 9,
				"OnRecordAfterDeleteSuccess":     5,
				"OnRecordDelete":                 5,
				"OnRecordDeleteExecute":          5,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				collections := []*core.Collection{}
				if err := app.CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}

				systemCollections := 0
				for _, c := range collections {
					if c.System {
						systemCollections++
					}
				}

				expected := systemCollections + 2
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
