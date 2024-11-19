package apis_test

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestCollectionsList(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			URL:             "/api/collections",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodGet,
			URL:    "/api/collections",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser",
			Method: http.MethodGet,
			URL:    "/api/collections",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":16`,
				`"items":[{`,
				`"name":"` + core.CollectionNameSuperusers + `"`,
				`"name":"` + core.CollectionNameAuthOrigins + `"`,
				`"name":"` + core.CollectionNameExternalAuths + `"`,
				`"name":"` + core.CollectionNameMFAs + `"`,
				`"name":"` + core.CollectionNameOTPs + `"`,
				`"name":"users"`,
				`"name":"nologin"`,
				`"name":"clients"`,
				`"name":"demo1"`,
				`"name":"demo2"`,
				`"name":"demo3"`,
				`"name":"demo4"`,
				`"name":"demo5"`,
				`"name":"numeric_id_view"`,
				`"name":"view1"`,
				`"name":"view2"`,
				`"type":"auth"`,
				`"type":"base"`,
				`"type":"view"`,
			},
			ExpectedEvents: map[string]int{
				"*":                        0,
				"OnCollectionsListRequest": 1,
			},
		},
		{
			Name:   "authorized as superuser + paging and sorting",
			Method: http.MethodGet,
			URL:    "/api/collections?page=2&perPage=2&sort=-created",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":2`,
				`"perPage":2`,
				`"totalItems":16`,
				`"items":[{`,
				`"name":"` + core.CollectionNameMFAs + `"`,
			},
			ExpectedEvents: map[string]int{
				"*":                        0,
				"OnCollectionsListRequest": 1,
			},
		},
		{
			Name:   "authorized as superuser + invalid filter",
			Method: http.MethodGet,
			URL:    "/api/collections?filter=invalidfield~'demo2'",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + valid filter",
			Method: http.MethodGet,
			URL:    "/api/collections?filter=name~'demo'",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":5`,
				`"items":[{`,
				`"name":"demo1"`,
				`"name":"demo2"`,
				`"name":"demo3"`,
				`"name":"demo4"`,
				`"name":"demo5"`,
			},
			ExpectedEvents: map[string]int{
				"*":                        0,
				"OnCollectionsListRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionView(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			URL:             "/api/collections/demo1",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodGet,
			URL:    "/api/collections/demo1",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + nonexisting collection identifier",
			Method: http.MethodGet,
			URL:    "/api/collections/missing",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + using the collection name",
			Method: http.MethodGet,
			URL:    "/api/collections/demo1",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"wsmn24bux7wo113"`,
				`"name":"demo1"`,
			},
			ExpectedEvents: map[string]int{
				"*":                       0,
				"OnCollectionViewRequest": 1,
			},
		},
		{
			Name:   "authorized as superuser + using the collection id",
			Method: http.MethodGet,
			URL:    "/api/collections/wsmn24bux7wo113",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"wsmn24bux7wo113"`,
				`"name":"demo1"`,
			},
			ExpectedEvents: map[string]int{
				"*":                       0,
				"OnCollectionViewRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionDelete(t *testing.T) {
	t.Parallel()

	ensureDeletedFiles := func(app *tests.TestApp, collectionId string) {
		storageDir := filepath.Join(app.DataDir(), "storage", collectionId)

		entries, _ := os.ReadDir(storageDir)
		if len(entries) != 0 {
			t.Errorf("Expected empty/deleted dir, found %d", len(entries))
		}
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodDelete,
			URL:             "/api/collections/demo1",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodDelete,
			URL:    "/api/collections/demo1",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + nonexisting collection identifier",
			Method: http.MethodDelete,
			URL:    "/api/collections/missing",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + using the collection name",
			Method: http.MethodDelete,
			URL:    "/api/collections/demo5",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionDeleteRequest":      1,
				"OnCollectionDelete":             1,
				"OnCollectionDeleteExecute":      1,
				"OnCollectionAfterDeleteSuccess": 1,
				"OnModelDelete":                  1,
				"OnModelDeleteExecute":           1,
				"OnModelAfterDeleteSuccess":      1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				ensureDeletedFiles(app, "9n89pl5vkct6330")
			},
		},
		{
			Name:   "authorized as superuser + using the collection id",
			Method: http.MethodDelete,
			URL:    "/api/collections/9n89pl5vkct6330",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionDeleteRequest":      1,
				"OnCollectionDelete":             1,
				"OnCollectionDeleteExecute":      1,
				"OnCollectionAfterDeleteSuccess": 1,
				"OnModelDelete":                  1,
				"OnModelDeleteExecute":           1,
				"OnModelAfterDeleteSuccess":      1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				ensureDeletedFiles(app, "9n89pl5vkct6330")
			},
		},
		{
			Name:   "authorized as superuser + trying to delete a system collection",
			Method: http.MethodDelete,
			URL:    "/api/collections/" + core.CollectionNameMFAs,
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionDeleteRequest":    1,
				"OnCollectionDelete":           1,
				"OnCollectionDeleteExecute":    1,
				"OnCollectionAfterDeleteError": 1,
				"OnModelDelete":                1,
				"OnModelDeleteExecute":         1,
				"OnModelAfterDeleteError":      1,
			},
		},
		{
			Name:   "authorized as superuser + trying to delete a referenced collection",
			Method: http.MethodDelete,
			URL:    "/api/collections/demo2",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionDeleteRequest":    1,
				"OnCollectionDelete":           1,
				"OnCollectionDeleteExecute":    1,
				"OnCollectionAfterDeleteError": 1,
				"OnModelDelete":                1,
				"OnModelDeleteExecute":         1,
				"OnModelAfterDeleteError":      1,
			},
		},
		{
			Name:   "authorized as superuser + deleting a view",
			Method: http.MethodDelete,
			URL:    "/api/collections/view2",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionDeleteRequest":      1,
				"OnCollectionDelete":             1,
				"OnCollectionDeleteExecute":      1,
				"OnCollectionAfterDeleteSuccess": 1,
				"OnModelDelete":                  1,
				"OnModelDeleteExecute":           1,
				"OnModelAfterDeleteSuccess":      1,
			},
		},
		{
			Name:   "OnCollectionAfterDeleteSuccessRequest error response",
			Method: http.MethodDelete,
			URL:    "/api/collections/view2",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.OnCollectionDeleteRequest().BindFunc(func(e *core.CollectionRequestEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                         0,
				"OnCollectionDeleteRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionCreate(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			URL:             "/api/collections",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body:   strings.NewReader(`{"name":"new","type":"base","fields":[{"type":"text","name":"test"}]}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + empty data",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body:   strings.NewReader(``),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"name":{"code":"validation_required"`,
			},
			NotExpectedContent: []string{
				`"fields":{`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionCreateRequest":    1,
				"OnCollectionCreate":           1,
				"OnCollectionAfterCreateError": 1,
				"OnCollectionValidate":         1,
				"OnModelCreate":                1,
				"OnModelAfterCreateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "authorized as superuser + invalid data (eg. existing name)",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body:   strings.NewReader(`{"name":"demo1","type":"base","fields":[{"type":"text","name":""}]}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"fields":{`,
				`"name":{"code":"validation_collection_name_exists"`,
				`"name":{"code":"validation_required"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionCreateRequest":    1,
				"OnCollectionCreate":           1,
				"OnCollectionAfterCreateError": 1,
				"OnCollectionValidate":         1,
				"OnModelCreate":                1,
				"OnModelAfterCreateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "authorized as superuser + valid data",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body:   strings.NewReader(`{"name":"new","type":"base","fields":[{"type":"text","id":"12345789","name":"test"}]}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"name":"new"`,
				`"type":"base"`,
				`"system":false`,
				// ensures that id field was prepended
				`"fields":[{"autogeneratePattern":"[a-z0-9]{15}","hidden":false,"id":"text3208210256","max":15,"min":15,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"},{"autogeneratePattern":"","hidden":false,"id":"12345789","max":0,"min":0,"name":"test","pattern":"","presentable":false,"primaryKey":false,"required":false,"system":false,"type":"text"}]`,
			},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionCreateRequest":      1,
				"OnCollectionCreate":             1,
				"OnCollectionCreateExecute":      1,
				"OnCollectionAfterCreateSuccess": 1,
				"OnCollectionValidate":           1,
				"OnModelCreate":                  1,
				"OnModelCreateExecute":           1,
				"OnModelAfterCreateSuccess":      1,
				"OnModelValidate":                1,
			},
		},
		{
			Name:   "creating auth collection (default settings merge test)",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"auth",
				"emailChangeToken":{"duration":123},
				"fields":[
					{"type":"text","id":"12345789","name":"test"},
					{"type":"text","name":"tokenKey","system":true,"required":false,"min":10}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"name":"new"`,
				`"type":"auth"`,
				`"system":false`,
				`"passwordAuth":{"enabled":true,"identityFields":["email"]}`,
				`"authRule":""`,
				`"manageRule":null`,
				`"name":"test"`,
				`"name":"id"`,
				`"name":"tokenKey"`,
				`"name":"password"`,
				`"name":"email"`,
				`"name":"emailVisibility"`,
				`"name":"verified"`,
				`"duration":123`,
				// should overwrite the user required option but keep the min value
				`{"autogeneratePattern":"","hidden":true,"id":"text2504183744","max":0,"min":10,"name":"tokenKey","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":true,"type":"text"}`,
			},
			NotExpectedContent: []string{
				`"secret":"`,
			},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionCreateRequest":      1,
				"OnCollectionCreate":             1,
				"OnCollectionCreateExecute":      1,
				"OnCollectionAfterCreateSuccess": 1,
				"OnCollectionValidate":           1,
				"OnModelCreate":                  1,
				"OnModelCreateExecute":           1,
				"OnModelAfterCreateSuccess":      1,
				"OnModelValidate":                1,
			},
		},
		{
			Name:   "creating base collection with reserved auth fields",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"fields":[
					{"type":"text","name":"email"},
					{"type":"text","name":"username"},
					{"type":"text","name":"verified"},
					{"type":"text","name":"emailVisibility"},
					{"type":"text","name":"lastResetSentAt"},
					{"type":"text","name":"lastVerificationSentAt"},
					{"type":"text","name":"tokenKey"},
					{"type":"text","name":"passwordHash"},
					{"type":"text","name":"password"},
					{"type":"text","name":"passwordConfirm"},
					{"type":"text","name":"oldPassword"}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"new"`,
				`"type":"base"`,
				`"fields":[{`,
			},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionCreateRequest":      1,
				"OnCollectionCreate":             1,
				"OnCollectionCreateExecute":      1,
				"OnCollectionAfterCreateSuccess": 1,
				"OnCollectionValidate":           1,
				"OnModelCreate":                  1,
				"OnModelCreateExecute":           1,
				"OnModelAfterCreateSuccess":      1,
				"OnModelValidate":                1,
			},
		},
		{
			Name:   "trying to create base collection with reserved system fields",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"fields":[
					{"type":"text","name":"id"},
					{"type":"text","name":"expand"},
					{"type":"text","name":"collectionId"},
					{"type":"text","name":"collectionName"}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"fields":{`,
				`"1":{"name":{"code":"validation_not_in_invalid`,
				`"2":{"name":{"code":"validation_not_in_invalid`,
				`"3":{"name":{"code":"validation_not_in_invalid`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionCreateRequest":    1,
				"OnCollectionCreate":           1,
				"OnCollectionAfterCreateError": 1,
				"OnCollectionValidate":         1,
				"OnModelCreate":                1,
				"OnModelAfterCreateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "trying to create auth collection with reserved auth fields",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"auth",
				"fields":[
					{"type":"text","name":"oldPassword"},
					{"type":"text","name":"passwordConfirm"}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"fields":{`,
				`"1":{"name":{"code":"validation_reserved_field_name`,
				`"2":{"name":{"code":"validation_reserved_field_name`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionCreateRequest":    1,
				"OnCollectionCreate":           1,
				"OnCollectionAfterCreateError": 1,
				"OnCollectionValidate":         1,
				"OnModelCreate":                1,
				"OnModelAfterCreateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "OnCollectionCreateRequest error response",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body:   strings.NewReader(`{"name":"new","type":"base"}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.OnCollectionCreateRequest().BindFunc(func(e *core.CollectionRequestEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                         0,
				"OnCollectionCreateRequest": 1,
			},
		},

		// view
		// -----------------------------------------------------------
		{
			Name:   "trying to create view collection with invalid options",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"view",
				"fields":[{"type":"text","id":"12345789","name":"ignored!@#$"}],
				"viewQuery":"invalid"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"viewQuery":{"code":"validation_invalid_view_query`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionCreateRequest":    1,
				"OnCollectionCreate":           1,
				"OnCollectionAfterCreateError": 1,
				"OnCollectionValidate":         1,
				"OnModelCreate":                1,
				"OnModelAfterCreateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "creating view collection",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"view",
				"fields":[{"type":"text","id":"12345789","name":"ignored!@#$"}],
				"viewQuery": "select 1 as id from ` + core.CollectionNameSuperusers + `"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"new"`,
				`"type":"view"`,
				`"fields":[{"autogeneratePattern":"","hidden":false,"id":"text3208210256","max":0,"min":0,"name":"id","pattern":"^[a-z0-9]+$","presentable":false,"primaryKey":true,"required":true,"system":true,"type":"text"}]`,
			},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionCreateRequest":      1,
				"OnCollectionCreate":             1,
				"OnCollectionCreateExecute":      1,
				"OnCollectionAfterCreateSuccess": 1,
				"OnCollectionValidate":           1,
				"OnModelCreate":                  1,
				"OnModelCreateExecute":           1,
				"OnModelAfterCreateSuccess":      1,
				"OnModelValidate":                1,
			},
		},

		// indexes
		// -----------------------------------------------------------
		{
			Name:   "creating base collection with invalid indexes",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"fields":[
					{"type":"text","name":"test"}
				],
				"indexes": [
					"create index idx_test1 on new (test)",
					"create index idx_test2 on new (missing)"
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"indexes":{"1":{"code":"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionCreateRequest":    1,
				"OnCollectionCreate":           1,
				"OnCollectionCreateExecute":    1,
				"OnCollectionAfterCreateError": 1,
				"OnCollectionValidate":         1,
				"OnModelCreate":                1,
				"OnModelCreateExecute":         1,
				"OnModelAfterCreateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "creating base collection with index name from another collection",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"fields":[
					{"type":"text","name":"test"}
				],
				"indexes": [
					"create index exist_test on new (test)"
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				demo1, err := app.FindCollectionByNameOrId("demo1")
				if err != nil {
					t.Fatal(err)
				}
				demo1.AddIndex("exist_test", false, "updated", "")
				if err = app.Save(demo1); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"indexes":{`,
				`"0":{"code":"validation_existing_index_name"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionCreateRequest":    1,
				"OnCollectionCreate":           1,
				"OnCollectionAfterCreateError": 1,
				"OnCollectionValidate":         1,
				"OnModelCreate":                1,
				"OnModelAfterCreateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "creating base collection with 2 indexes using the same name",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"indexes": [
					"create index duplicate_idx on new (created)",
					"create index duplicate_idx on new (updated)"
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"indexes":{`,
				`"1":{"code":"validation_duplicated_index_name"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionCreateRequest":    1,
				"OnCollectionCreate":           1,
				"OnCollectionAfterCreateError": 1,
				"OnCollectionValidate":         1,
				"OnModelCreate":                1,
				"OnModelAfterCreateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "creating base collection with valid indexes (+ random table name)",
			Method: http.MethodPost,
			URL:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"fields":[
					{"type":"text","name":"test"}
				],
				"indexes": [
					"create index idx_test1 on new (test)",
					"create index idx_test2 on anything (id, test)"
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"new"`,
				`"type":"base"`,
				`"indexes":[`,
				`idx_test1`,
				`idx_test2`,
			},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionCreateRequest":      1,
				"OnCollectionCreate":             1,
				"OnCollectionCreateExecute":      1,
				"OnCollectionAfterCreateSuccess": 1,
				"OnCollectionValidate":           1,
				"OnModelCreate":                  1,
				"OnModelCreateExecute":           1,
				"OnModelAfterCreateSuccess":      1,
				"OnModelValidate":                1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				indexes, err := app.TableIndexes("new")
				if err != nil {
					t.Fatal(err)
				}

				expected := []string{"idx_test1", "idx_test2"}
				if len(indexes) != len(expected) {
					t.Fatalf("Expected %d indexes, got %d\n%v", len(expected), len(indexes), indexes)
				}
				for name := range indexes {
					if !list.ExistInSlice(name, expected) {
						t.Fatalf("Missing index %q", name)
					}
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionUpdate(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPatch,
			URL:             "/api/collections/demo1",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo1",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + missing collection",
			Method: http.MethodPatch,
			URL:    "/api/collections/missing",
			Body:   strings.NewReader(`{}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + empty body",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo1",
			Body:   strings.NewReader(`{}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"wsmn24bux7wo113"`,
				`"name":"demo1"`,
			},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionUpdateRequest":      1,
				"OnCollectionUpdate":             1,
				"OnCollectionUpdateExecute":      1,
				"OnCollectionAfterUpdateSuccess": 1,
				"OnCollectionValidate":           1,
				"OnModelUpdate":                  1,
				"OnModelUpdateExecute":           1,
				"OnModelAfterUpdateSuccess":      1,
				"OnModelValidate":                1,
			},
		},
		{
			Name:   "OnCollectionAfterUpdateSuccessRequest error response",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo1",
			Body:   strings.NewReader(`{}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				app.OnCollectionUpdateRequest().BindFunc(func(e *core.CollectionRequestEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                         0,
				"OnCollectionUpdateRequest": 1,
			},
		},
		{
			Name:   "authorized as superuser + invalid data (eg. existing name)",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo1",
			Body: strings.NewReader(`{
				"name":"demo2",
				"type":"auth"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"name":{"code":"validation_collection_name_exists"`,
				`"type":{"code":"validation_collection_type_change"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionUpdateRequest":    1,
				"OnCollectionUpdate":           1,
				"OnCollectionAfterUpdateError": 1,
				"OnCollectionValidate":         1,
				"OnModelUpdate":                1,
				"OnModelAfterUpdateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "authorized as superuser + valid data",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo1",
			Body:   strings.NewReader(`{"name":"new"}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"name":"new"`,
			},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionUpdateRequest":      1,
				"OnCollectionUpdate":             1,
				"OnCollectionUpdateExecute":      1,
				"OnCollectionAfterUpdateSuccess": 1,
				"OnCollectionValidate":           1,
				"OnModelUpdate":                  1,
				"OnModelUpdateExecute":           1,
				"OnModelAfterUpdateSuccess":      1,
				"OnModelValidate":                1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				// check if the record table was renamed
				if !app.HasTable("new") {
					t.Fatal("Couldn't find record table 'new'.")
				}
			},
		},
		{
			Name:   "trying to update collection with reserved fields",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo1",
			Body: strings.NewReader(`{
				"name":"new",
				"fields":[
					{"type":"text","name":"id","id":"_pbf_text_id_"},
					{"type":"text","name":"created"},
					{"type":"text","name":"updated"},
					{"type":"text","name":"expand"},
					{"type":"text","name":"collectionId"},
					{"type":"text","name":"collectionName"}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"fields":{`,
				`"3":{"name":{"code":"validation_not_in_invalid`,
				`"4":{"name":{"code":"validation_not_in_invalid`,
				`"5":{"name":{"code":"validation_not_in_invalid`,
			},
			NotExpectedContent: []string{
				`"0":`,
				`"1":`,
				`"2":`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionUpdateRequest":    1,
				"OnCollectionUpdate":           1,
				"OnCollectionAfterUpdateError": 1,
				"OnCollectionValidate":         1,
				"OnModelUpdate":                1,
				"OnModelAfterUpdateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "trying to update collection with changed/removed system fields",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo1",
			Body: strings.NewReader(`{
				"name":"new",
				"fields":[
					{"type":"text","name":"created"}
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"fields":{`,
				`"code":"validation_system_field_change"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionUpdateRequest":    1,
				"OnCollectionUpdate":           1,
				"OnCollectionAfterUpdateError": 1,
				"OnCollectionValidate":         1,
				"OnModelUpdate":                1,
				"OnModelAfterUpdateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "trying to update auth collection with invalid options",
			Method: http.MethodPatch,
			URL:    "/api/collections/users",
			Body: strings.NewReader(`{
				"passwordAuth":{"identityFields": ["missing"]}
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"passwordAuth":{"identityFields":{"code":"validation_missing_field"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionUpdateRequest":    1,
				"OnCollectionUpdate":           1,
				"OnCollectionAfterUpdateError": 1,
				"OnCollectionValidate":         1,
				"OnModelUpdate":                1,
				"OnModelAfterUpdateError":      1,
				"OnModelValidate":              1,
			},
		},

		// view
		// -----------------------------------------------------------
		{
			Name:   "trying to update view collection with invalid options",
			Method: http.MethodPatch,
			URL:    "/api/collections/view1",
			Body: strings.NewReader(`{
				"fields":[{"type":"text","id":"12345789","name":"ignored!@#$"}],
				"viewQuery":"invalid"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"viewQuery":{"code":"validation_invalid_view_query"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionUpdateRequest":    1,
				"OnCollectionUpdate":           1,
				"OnCollectionAfterUpdateError": 1,
				"OnCollectionValidate":         1,
				"OnModelUpdate":                1,
				"OnModelAfterUpdateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "updating view collection",
			Method: http.MethodPatch,
			URL:    "/api/collections/view2",
			Body: strings.NewReader(`{
				"name":"view2_update",
				"fields":[{"type":"text","id":"12345789","name":"ignored!@#$"}],
				"viewQuery": "select 2 as id, created, updated, email from ` + core.CollectionNameSuperusers + `"
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"view2_update"`,
				`"type":"view"`,
				`"fields":[{`,
				`"name":"email"`,
				`"name":"id"`,
				`"name":"created"`,
				`"name":"updated"`,
			},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionUpdateRequest":      1,
				"OnCollectionUpdate":             1,
				"OnCollectionUpdateExecute":      1,
				"OnCollectionAfterUpdateSuccess": 1,
				"OnCollectionValidate":           1,
				"OnModelUpdate":                  1,
				"OnModelUpdateExecute":           1,
				"OnModelAfterUpdateSuccess":      1,
				"OnModelValidate":                1,
			},
		},

		// indexes
		// -----------------------------------------------------------
		{
			Name:   "updating base collection with invalid indexes",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo2",
			Body: strings.NewReader(`{
				"indexes": [
					"create unique idx_test1 on demo1 (text)",
					"create index idx_test2 on demo2 (id, title)"
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"indexes":{"0":{"code":"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionUpdateRequest":    1,
				"OnCollectionUpdate":           1,
				"OnCollectionAfterUpdateError": 1,
				"OnCollectionValidate":         1,
				"OnModelUpdate":                1,
				"OnModelAfterUpdateError":      1,
				"OnModelValidate":              1,
			},
		},

		{
			Name:   "updating base collection with index name from another collection",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo2",
			Body: strings.NewReader(`{
				"indexes": [
					"create index exist_test on new (test)"
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				demo1, err := app.FindCollectionByNameOrId("demo1")
				if err != nil {
					t.Fatal(err)
				}
				demo1.AddIndex("exist_test", false, "updated", "")
				if err = app.Save(demo1); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"indexes":{`,
				`"0":{"code":"validation_existing_index_name"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionUpdateRequest":    1,
				"OnCollectionUpdate":           1,
				"OnCollectionAfterUpdateError": 1,
				"OnCollectionValidate":         1,
				"OnModelUpdate":                1,
				"OnModelAfterUpdateError":      1,
				"OnModelValidate":              1,
			},
		},
		{
			Name:   "updating base collection with 2 indexes using the same name",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo2",
			Body: strings.NewReader(`{
				"indexes": [
					"create index duplicate_idx on new (created)",
					"create index duplicate_idx on new (updated)"
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"indexes":{`,
				`"1":{"code":"validation_duplicated_index_name"`,
			},
			ExpectedEvents: map[string]int{
				"*":                            0,
				"OnCollectionUpdateRequest":    1,
				"OnCollectionUpdate":           1,
				"OnCollectionAfterUpdateError": 1,
				"OnCollectionValidate":         1,
				"OnModelUpdate":                1,
				"OnModelAfterUpdateError":      1,
				"OnModelValidate":              1,
			},
		},

		{
			Name:   "updating base collection with valid indexes (+ random table name)",
			Method: http.MethodPatch,
			URL:    "/api/collections/demo2",
			Body: strings.NewReader(`{
				"indexes": [
					"create unique index idx_test1 on demo2 (title)",
					"create index idx_test2 on anything (active)"
				]
			}`),
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"demo2"`,
				`"indexes":[`,
				`idx_test1`,
				`idx_test2`,
			},
			ExpectedEvents: map[string]int{
				"*":                              0,
				"OnCollectionUpdateRequest":      1,
				"OnCollectionUpdate":             1,
				"OnCollectionUpdateExecute":      1,
				"OnCollectionAfterUpdateSuccess": 1,
				"OnCollectionValidate":           1,
				"OnModelUpdate":                  1,
				"OnModelUpdateExecute":           1,
				"OnModelAfterUpdateSuccess":      1,
				"OnModelValidate":                1,
			},
			AfterTestFunc: func(t testing.TB, app *tests.TestApp, res *http.Response) {
				indexes, err := app.TableIndexes("demo2")
				if err != nil {
					t.Fatal(err)
				}

				expected := []string{"idx_test1", "idx_test2"}
				if len(indexes) != len(expected) {
					t.Fatalf("Expected %d indexes, got %d\n%v", len(expected), len(indexes), indexes)
				}
				for name := range indexes {
					if !list.ExistInSlice(name, expected) {
						t.Fatalf("Missing index %q", name)
					}
				}
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionScaffolds(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			URL:             "/api/collections/meta/scaffolds",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodGet,
			URL:    "/api/collections/meta/scaffolds",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser",
			Method: http.MethodGet,
			URL:    "/api/collections/meta/scaffolds",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":""`,
				`"name":""`,
				`"auth":{`,
				`"base":{`,
				`"view":{`,
				`"type":"auth"`,
				`"type":"base"`,
				`"type":"view"`,
				`"fields":[{`,
				`"fields":[{`,
				`"id":"text3208210256"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCollectionTruncate(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodDelete,
			URL:             "/api/collections/demo5/truncate",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodDelete,
			URL:    "/api/collections/demo5/truncate",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser",
			Method: http.MethodDelete,
			URL:    "/api/collections/demo5/truncate",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"*":                          0,
				"OnModelDelete":              2,
				"OnModelDeleteExecute":       2,
				"OnModelAfterDeleteSuccess":  2,
				"OnRecordDelete":             2,
				"OnRecordDeleteExecute":      2,
				"OnRecordAfterDeleteSuccess": 2,
			},
		},
		{
			Name:   "authorized as superuser but collection with required cascade delete references",
			Method: http.MethodDelete,
			URL:    "/api/collections/demo3/truncate",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"*":                        0,
				"OnModelDelete":            2,
				"OnModelDeleteExecute":     2,
				"OnModelAfterDeleteError":  2,
				"OnModelUpdate":            2,
				"OnModelUpdateExecute":     2,
				"OnModelAfterUpdateError":  2,
				"OnRecordDelete":           2,
				"OnRecordDeleteExecute":    2,
				"OnRecordAfterDeleteError": 2,
				"OnRecordUpdate":           2,
				"OnRecordUpdateExecute":    2,
				"OnRecordAfterUpdateError": 2,
			},
		},
		{
			Name:   "authorized as superuser trying to truncate view collection",
			Method: http.MethodDelete,
			URL:    "/api/collections/view2/truncate",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
