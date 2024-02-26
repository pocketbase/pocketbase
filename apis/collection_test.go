package apis_test

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestCollectionsList(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/collections",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/collections",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/collections",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":11`,
				`"items":[{`,
				`"id":"_pb_users_auth_"`,
				`"id":"v851q4r790rhknl"`,
				`"id":"kpv709sk2lqbqk8"`,
				`"id":"wsmn24bux7wo113"`,
				`"id":"sz5l5z67tg7gku0"`,
				`"id":"wzlqyes4orhoygb"`,
				`"id":"4d1blo5cuycfaca"`,
				`"id":"9n89pl5vkct6330"`,
				`"id":"ib3m2700k5hlsjz"`,
				`"type":"auth"`,
				`"type":"base"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionsListRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + paging and sorting",
			Method: http.MethodGet,
			Url:    "/api/collections?page=2&perPage=2&sort=-created",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":2`,
				`"perPage":2`,
				`"totalItems":11`,
				`"items":[{`,
				`"id":"v9gwnfh02gjq1q0"`,
				`"id":"9n89pl5vkct6330"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionsListRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + invalid filter",
			Method: http.MethodGet,
			Url:    "/api/collections?filter=invalidfield~'demo2'",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + valid filter",
			Method: http.MethodGet,
			Url:    "/api/collections?filter=name~'demo'",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":5`,
				`"items":[{`,
				`"id":"wsmn24bux7wo113"`,
				`"id":"sz5l5z67tg7gku0"`,
				`"id":"wzlqyes4orhoygb"`,
				`"id":"4d1blo5cuycfaca"`,
				`"id":"9n89pl5vkct6330"`,
			},
			ExpectedEvents: map[string]int{
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
			Url:             "/api/collections/demo1",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting collection identifier",
			Method: http.MethodGet,
			Url:    "/api/collections/missing",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + using the collection name",
			Method: http.MethodGet,
			Url:    "/api/collections/demo1",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"wsmn24bux7wo113"`,
				`"name":"demo1"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionViewRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + using the collection id",
			Method: http.MethodGet,
			Url:    "/api/collections/wsmn24bux7wo113",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"wsmn24bux7wo113"`,
				`"name":"demo1"`,
			},
			ExpectedEvents: map[string]int{
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
			Url:             "/api/collections/demo1",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo1",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting collection identifier",
			Method: http.MethodDelete,
			Url:    "/api/collections/missing",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + using the collection name",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo5",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":             1,
				"OnModelAfterDelete":              1,
				"OnCollectionBeforeDeleteRequest": 1,
				"OnCollectionAfterDeleteRequest":  1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				ensureDeletedFiles(app, "9n89pl5vkct6330")
			},
		},
		{
			Name:   "authorized as admin + using the collection id",
			Method: http.MethodDelete,
			Url:    "/api/collections/9n89pl5vkct6330",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			Delay:          100 * time.Millisecond,
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":             1,
				"OnModelAfterDelete":              1,
				"OnCollectionBeforeDeleteRequest": 1,
				"OnCollectionAfterDeleteRequest":  1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				ensureDeletedFiles(app, "9n89pl5vkct6330")
			},
		},
		{
			Name:   "authorized as admin + trying to delete a system collection",
			Method: http.MethodDelete,
			Url:    "/api/collections/nologin",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnCollectionBeforeDeleteRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + trying to delete a referenced collection",
			Method: http.MethodDelete,
			Url:    "/api/collections/demo2",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnCollectionBeforeDeleteRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + deleting a view",
			Method: http.MethodDelete,
			Url:    "/api/collections/view2",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":             1,
				"OnModelAfterDelete":              1,
				"OnCollectionBeforeDeleteRequest": 1,
				"OnCollectionAfterDeleteRequest":  1,
			},
		},
		{
			Name:   "OnCollectionAfterDeleteRequest error response",
			Method: http.MethodDelete,
			Url:    "/api/collections/view2",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnCollectionAfterDeleteRequest().Add(func(e *core.CollectionDeleteEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":             1,
				"OnModelAfterDelete":              1,
				"OnCollectionBeforeDeleteRequest": 1,
				"OnCollectionAfterDeleteRequest":  1,
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
			Url:             "/api/collections",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPost,
			Url:    "/api/collections",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty data",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"name":{"code":"validation_required"`,
				`"schema":{"code":"validation_required"`,
			},
		},
		{
			Name:   "authorized as admin + invalid data (eg. existing name)",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body:   strings.NewReader(`{"name":"demo1","type":"base","schema":[{"type":"text","name":""}]}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"name":{"code":"validation_collection_name_exists"`,
				`"schema":{"0":{"name":{"code":"validation_required"`,
			},
		},
		{
			Name:   "authorized as admin + valid data",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body:   strings.NewReader(`{"name":"new","type":"base","schema":[{"type":"text","id":"12345789","name":"test"}]}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"name":"new"`,
				`"type":"base"`,
				`"system":false`,
				`"schema":[{"system":false,"id":"12345789","name":"test","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}}]`,
				`"options":{}`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":             1,
				"OnModelAfterCreate":              1,
				"OnCollectionBeforeCreateRequest": 1,
				"OnCollectionAfterCreateRequest":  1,
			},
		},
		{
			Name:   "creating auth collection without specified options",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body:   strings.NewReader(`{"name":"new","type":"auth","schema":[{"type":"text","id":"12345789","name":"test"}]}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"name":"new"`,
				`"type":"auth"`,
				`"system":false`,
				`"schema":[{"system":false,"id":"12345789","name":"test","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}}]`,
				`"options":{"allowEmailAuth":false,"allowOAuth2Auth":false,"allowUsernameAuth":false,"exceptEmailDomains":null,"manageRule":null,"minPasswordLength":0,"onlyEmailDomains":null,"onlyVerified":false,"requireEmail":false}`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":             1,
				"OnModelAfterCreate":              1,
				"OnCollectionBeforeCreateRequest": 1,
				"OnCollectionAfterCreateRequest":  1,
			},
		},
		{
			Name:   "trying to create auth collection with reserved auth fields",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"auth",
				"schema":[
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
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"schema":{`,
				`"0":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"1":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"2":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"3":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"4":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"5":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"6":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"7":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"8":{"name":{"code":"validation_reserved_auth_field_name"`,
			},
		},
		{
			Name:   "creating base collection with reserved auth fields",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"schema":[
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
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"new"`,
				`"type":"base"`,
				`"schema":[{`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":             1,
				"OnModelAfterCreate":              1,
				"OnCollectionBeforeCreateRequest": 1,
				"OnCollectionAfterCreateRequest":  1,
			},
		},
		{
			Name:   "trying to create base collection with reserved base fields",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"schema":[
					{"type":"text","name":"id"},
					{"type":"text","name":"created"},
					{"type":"text","name":"updated"},
					{"type":"text","name":"expand"},
					{"type":"text","name":"collectionId"},
					{"type":"text","name":"collectionName"}
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"schema":{`,
				`"0":{"name":{"code":"validation_not_in_invalid`,
				`"1":{"name":{"code":"validation_not_in_invalid`,
				`"2":{"name":{"code":"validation_not_in_invalid`,
				`"3":{"name":{"code":"validation_not_in_invalid`,
				`"4":{"name":{"code":"validation_not_in_invalid`,
				`"5":{"name":{"code":"validation_not_in_invalid`,
			},
		},
		{
			Name:   "trying to create auth collection with invalid options",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"auth",
				"schema":[{"type":"text","id":"12345789","name":"test"}],
				"options":{"allowUsernameAuth": true}
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"options":{"minPasswordLength":{"code":"validation_required"`,
			},
		},
		{
			Name:   "OnCollectionAfterCreateRequest error response",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body:   strings.NewReader(`{"name":"new","type":"base","schema":[{"type":"text","id":"12345789","name":"test"}]}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnCollectionAfterCreateRequest().Add(func(e *core.CollectionCreateEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":             1,
				"OnModelAfterCreate":              1,
				"OnCollectionBeforeCreateRequest": 1,
				"OnCollectionAfterCreateRequest":  1,
			},
		},

		// view
		// -----------------------------------------------------------
		{
			Name:   "trying to create view collection with invalid options",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"view",
				"schema":[{"type":"text","id":"12345789","name":"ignored!@#$"}],
				"options":{"query": "invalid"}
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"options":{"query":{"code":"validation_invalid_view_query`,
			},
		},
		{
			Name:   "creating view collection",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"view",
				"schema":[{"type":"text","id":"12345789","name":"ignored!@#$"}],
				"options": {
					"query": "select 1 as id from _admins"
				}
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"new"`,
				`"type":"view"`,
				`"schema":[]`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":             1,
				"OnModelAfterCreate":              1,
				"OnCollectionBeforeCreateRequest": 1,
				"OnCollectionAfterCreateRequest":  1,
			},
		},

		// indexes
		// -----------------------------------------------------------
		{
			Name:   "creating base collection with invalid indexes",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"schema":[
					{"type":"text","name":"test"}
				],
				"indexes": [
					"create index idx_test1 on new (test)",
					"create index idx_test2 on new (missing)"
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"indexes":{"1":{"code":"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionBeforeCreateRequest": 1,
				"OnModelBeforeCreate":             1,
			},
		},
		{
			Name:   "creating base collection with valid indexes (+ random table name)",
			Method: http.MethodPost,
			Url:    "/api/collections",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"schema":[
					{"type":"text","name":"test"}
				],
				"indexes": [
					"create index idx_test1 on new (test)",
					"create index idx_test2 on anything (id, test)"
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
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
				"OnModelBeforeCreate":             1,
				"OnModelAfterCreate":              1,
				"OnCollectionBeforeCreateRequest": 1,
				"OnCollectionAfterCreateRequest":  1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				indexes, err := app.Dao().TableIndexes("new")
				if err != nil {
					t.Fatal(err)
				}

				expected := []string{"idx_test1", "idx_test2"}
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
			Url:             "/api/collections/demo1",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo1",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + missing collection",
			Method: http.MethodPatch,
			Url:    "/api/collections/missing",
			Body:   strings.NewReader(`{}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty body",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo1",
			Body:   strings.NewReader(`{}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"wsmn24bux7wo113"`,
				`"name":"demo1"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionAfterUpdateRequest":  1,
				"OnCollectionBeforeUpdateRequest": 1,
				"OnModelAfterUpdate":              1,
				"OnModelBeforeUpdate":             1,
			},
		},
		{
			Name:   "OnCollectionAfterUpdateRequest error response",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo1",
			Body:   strings.NewReader(`{}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnCollectionAfterUpdateRequest().Add(func(e *core.CollectionUpdateEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnCollectionAfterUpdateRequest":  1,
				"OnCollectionBeforeUpdateRequest": 1,
				"OnModelAfterUpdate":              1,
				"OnModelBeforeUpdate":             1,
			},
		},
		{
			Name:   "authorized as admin + invalid data (eg. existing name)",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo1",
			Body: strings.NewReader(`{
				"name":"demo2",
				"type":"auth"
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"name":{"code":"validation_collection_name_exists"`,
				`"type":{"code":"validation_collection_type_change"`,
			},
		},
		{
			Name:   "authorized as admin + valid data",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo1",
			Body:   strings.NewReader(`{"name":"new"}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"name":"new"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":             1,
				"OnModelAfterUpdate":              1,
				"OnCollectionBeforeUpdateRequest": 1,
				"OnCollectionAfterUpdateRequest":  1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				// check if the record table was renamed
				if !app.Dao().HasTable("new") {
					t.Fatal("Couldn't find record table 'new'.")
				}
			},
		},
		{
			Name:   "trying to update auth collection with reserved auth fields",
			Method: http.MethodPatch,
			Url:    "/api/collections/users",
			Body: strings.NewReader(`{
				"schema":[
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
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"schema":{`,
				`"0":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"1":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"2":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"3":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"4":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"5":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"6":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"7":{"name":{"code":"validation_reserved_auth_field_name"`,
				`"8":{"name":{"code":"validation_reserved_auth_field_name"`,
			},
		},
		{
			Name:   "updating base collection with reserved auth fields",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo4",
			Body: strings.NewReader(`{
				"schema":[
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
				],
				"indexes": []
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"demo4"`,
				`"type":"base"`,
				`"schema":[{`,
				`"email"`,
				`"username"`,
				`"verified"`,
				`"emailVisibility"`,
				`"lastResetSentAt"`,
				`"lastVerificationSentAt"`,
				`"tokenKey"`,
				`"passwordHash"`,
				`"password"`,
				`"passwordConfirm"`,
				`"oldPassword"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":             1,
				"OnModelAfterUpdate":              1,
				"OnCollectionBeforeUpdateRequest": 1,
				"OnCollectionAfterUpdateRequest":  1,
			},
		},
		{
			Name:   "trying to update base collection with reserved base fields",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo1",
			Body: strings.NewReader(`{
				"name":"new",
				"type":"base",
				"schema":[
					{"type":"text","name":"id"},
					{"type":"text","name":"created"},
					{"type":"text","name":"updated"},
					{"type":"text","name":"expand"},
					{"type":"text","name":"collectionId"},
					{"type":"text","name":"collectionName"}
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{"schema":{`,
				`"0":{"name":{"code":"validation_not_in_invalid`,
				`"1":{"name":{"code":"validation_not_in_invalid`,
				`"2":{"name":{"code":"validation_not_in_invalid`,
				`"3":{"name":{"code":"validation_not_in_invalid`,
				`"4":{"name":{"code":"validation_not_in_invalid`,
				`"5":{"name":{"code":"validation_not_in_invalid`,
			},
		},
		{
			Name:   "trying to update auth collection with invalid options",
			Method: http.MethodPatch,
			Url:    "/api/collections/users",
			Body: strings.NewReader(`{
				"options":{"minPasswordLength": 4}
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"options":{"minPasswordLength":{"code":"validation_min_greater_equal_than_required"`,
			},
		},

		// view
		// -----------------------------------------------------------
		{
			Name:   "trying to update view collection with invalid options",
			Method: http.MethodPatch,
			Url:    "/api/collections/view1",
			Body: strings.NewReader(`{
				"schema":[{"type":"text","id":"12345789","name":"ignored!@#$"}],
				"options":{"query": "invalid"}
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"options":{"query":{"code":"validation_invalid_view_query`,
			},
		},
		{
			Name:   "updating view collection",
			Method: http.MethodPatch,
			Url:    "/api/collections/view2",
			Body: strings.NewReader(`{
				"name":"view2_update",
				"schema":[{"type":"text","id":"12345789","name":"ignored!@#$"}],
				"options": {
					"query": "select 2 as id, created, updated, email from _admins"
				}
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"view2_update"`,
				`"type":"view"`,
				`"schema":[{`,
				`"name":"email"`,
			},
			NotExpectedContent: []string{
				// base model fields are not part of the schema
				`"name":"id"`,
				`"name":"created"`,
				`"name":"updated"`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":             1,
				"OnModelAfterUpdate":              1,
				"OnCollectionBeforeUpdateRequest": 1,
				"OnCollectionAfterUpdateRequest":  1,
			},
		},

		// indexes
		// -----------------------------------------------------------
		{
			Name:   "updating base collection with invalid indexes",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo2",
			Body: strings.NewReader(`{
				"indexes": [
					"create unique idx_test1 on demo1 (text)",
					"create index idx_test2 on demo2 (id, title)"
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"indexes":{"0":{"code":"`,
			},
		},
		{
			Name:   "updating base collection with valid indexes (+ random table name)",
			Method: http.MethodPatch,
			Url:    "/api/collections/demo2",
			Body: strings.NewReader(`{
				"indexes": [
					"create unique index idx_test1 on demo2 (title)",
					"create index idx_test2 on anything (active)"
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"name":"demo2"`,
				`"indexes":[`,
				"CREATE UNIQUE INDEX `idx_test1`",
				"CREATE INDEX `idx_test2`",
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":             1,
				"OnModelAfterUpdate":              1,
				"OnCollectionBeforeUpdateRequest": 1,
				"OnCollectionAfterUpdateRequest":  1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				indexes, err := app.Dao().TableIndexes("new")
				if err != nil {
					t.Fatal(err)
				}

				expected := []string{"idx_test1", "idx_test2"}
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

func TestCollectionsImport(t *testing.T) {
	t.Parallel()

	totalCollections := 11

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPut,
			Url:             "/api/collections/import",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty collections",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body:   strings.NewReader(`{"collections":[]}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"collections":{"code":"validation_required"`,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				expected := totalCollections
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}
			},
		},
		{
			Name:   "authorized as admin + trying to delete system collections",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body:   strings.NewReader(`{"deleteMissing": true, "collections":[{"name": "test123"}]}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"collections":{"code":"collections_import_failure"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionsBeforeImportRequest": 1,
				"OnModelBeforeDelete":              1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				expected := totalCollections
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}
			},
		},
		{
			Name:   "authorized as admin + collections validator failure",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body: strings.NewReader(`{
				"collections":[
					{
						"name": "import1",
						"schema": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					},
					{"name": "import2"}
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"collections":{"code":"collections_import_validate_failure"`,
			},
			ExpectedEvents: map[string]int{
				"OnCollectionsBeforeImportRequest": 1,
				"OnModelBeforeCreate":              2,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				expected := totalCollections
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}
			},
		},
		{
			Name:   "authorized as admin + successful collections save",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body: strings.NewReader(`{
				"collections":[
					{
						"name": "import1",
						"schema": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					},
					{
						"name": "import2",
						"schema": [
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
						"name": "auth_without_schema",
						"type": "auth"
					}
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnCollectionsBeforeImportRequest": 1,
				"OnCollectionsAfterImportRequest":  1,
				"OnModelBeforeCreate":              3,
				"OnModelAfterCreate":               3,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}

				expected := totalCollections + 3
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}

				indexes, err := app.Dao().TableIndexes("import2")
				if err != nil || indexes["idx_test"] == "" {
					t.Fatalf("Missing index %s (%v)", "idx_test", err)
				}
			},
		},
		{
			Name:   "authorized as admin + successful collections save and old non-system collections deletion",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body: strings.NewReader(`{
				"deleteMissing": true,
				"collections":[
					{
						"name": "new_import",
						"schema": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					},
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
						"id":"wsmn24bux7wo113",
						"name":"demo1",
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
					}
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnCollectionsAfterImportRequest":  1,
				"OnCollectionsBeforeImportRequest": 1,
				"OnModelBeforeDelete":              9,
				"OnModelAfterDelete":               9,
				"OnModelBeforeUpdate":              2,
				"OnModelAfterUpdate":               2,
				"OnModelBeforeCreate":              1,
				"OnModelAfterCreate":               1,
			},
			AfterTestFunc: func(t *testing.T, app *tests.TestApp, res *http.Response) {
				collections := []*models.Collection{}
				if err := app.Dao().CollectionQuery().All(&collections); err != nil {
					t.Fatal(err)
				}
				expected := 3
				if len(collections) != expected {
					t.Fatalf("Expected %d collections, got %d", expected, len(collections))
				}
			},
		},
		{
			Name:   "authorized as admin + successful collections save",
			Method: http.MethodPut,
			Url:    "/api/collections/import",
			Body: strings.NewReader(`{
				"collections":[
					{
						"name": "import1",
						"schema": [
							{
							  "id": "koih1lqx",
							  "name": "test",
							  "type": "text"
							}
						]
					},
					{
						"name": "import2",
						"schema": [
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
						"name": "auth_without_schema",
						"type": "auth"
					}
				]
			}`),
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				app.OnCollectionsAfterImportRequest().Add(func(e *core.CollectionsImportEvent) error {
					return errors.New("error")
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnCollectionsBeforeImportRequest": 1,
				"OnCollectionsAfterImportRequest":  1,
				"OnModelBeforeCreate":              3,
				"OnModelAfterCreate":               3,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
