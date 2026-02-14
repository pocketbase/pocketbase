package apis_test

import (
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestPanicRecover(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "panic from route",
			Method: http.MethodGet,
			URL:    "/my/test",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					panic("123")
				})
			},
			ExpectedStatus:  500,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "panic from middleware",
			Method: http.MethodGet,
			URL:    "/my/test",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(http.StatusOK, "test")
				}).BindFunc(func(e *core.RequestEvent) error {
					panic(123)
				})
			},
			ExpectedStatus:  500,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireGuestOnly(t *testing.T) {
	t.Parallel()

	beforeTestFunc := func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
		e.Router.GET("/my/test", func(e *core.RequestEvent) error {
			return e.String(200, "test123")
		}).Bind(apis.RequireGuestOnly())
	}

	scenarios := []tests.ApiScenario{
		{
			Name:   "valid regular user token",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc:  beforeTestFunc,
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid superuser auth token",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc:  beforeTestFunc,
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoxNjQwOTkxNjYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.2D3tmqPn3vc5LoqqCz8V-iCDVXo9soYiH0d32G7FQT4",
			},
			BeforeTestFunc:  beforeTestFunc,
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:            "guest",
			Method:          http.MethodGet,
			URL:             "/my/test",
			BeforeTestFunc:  beforeTestFunc,
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			URL:    "/my/test",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireAuth())
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "expired token",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoxNjQwOTkxNjYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.2D3tmqPn3vc5LoqqCz8V-iCDVXo9soYiH0d32G7FQT4",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireAuth())
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "invalid token",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsImV4cCI6MjUyNDYwNDQ2MSwidHlwZSI6ImZpbGUiLCJjb2xsZWN0aW9uSWQiOiJwYmNfMzE0MjYzNTgyMyJ9.Lupz541xRvrktwkrl55p5pPCF77T69ZRsohsIcb2dxc",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireAuth())
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid record auth token with no collection restrictions",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				// regular user
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireAuth())
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid record static auth token",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				// regular user
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6ZmFsc2V9.4IsO6YMsR19crhwl_YWzvRH8pfq2Ri4Gv2dzGyneLak",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireAuth())
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid record auth token with collection not in the restricted list",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				// superuser
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireAuth("users", "demo1"))
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid record auth token with collection in the restricted list",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				// superuser
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireAuth("users", core.CollectionNameSuperusers))
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid record auth token with Bearer case-insensitive prefix",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				// regular user
				"Authorization": "BeArEr eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireAuth())
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireSuperuserAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			URL:    "/my/test",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserAuth())
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjE2NDA5OTE2NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.0pDcBPGDpL2Khh76ivlRi7ugiLBSYvasct3qpHV3rfs",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserAuth())
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid regular user auth token",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserAuth())
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid superuser auth token",
			Method: http.MethodGet,
			URL:    "/my/test",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserAuth())
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireSuperuserOrOwnerAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			URL:    "/my/test/4q1xlclmfloku33",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{id}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserOrOwnerAuth(""))
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			URL:    "/my/test/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjE2NDA5OTE2NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.0pDcBPGDpL2Khh76ivlRi7ugiLBSYvasct3qpHV3rfs",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{id}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserOrOwnerAuth(""))
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid record auth token (different user)",
			Method: http.MethodGet,
			URL:    "/my/test/oap640cot4yru2s",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{id}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserOrOwnerAuth(""))
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid record auth token (owner)",
			Method: http.MethodGet,
			URL:    "/my/test/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{id}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserOrOwnerAuth(""))
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid record auth token (owner + non-matching custom owner param)",
			Method: http.MethodGet,
			URL:    "/my/test/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{id}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserOrOwnerAuth("test"))
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid record auth token (owner + matching custom owner param)",
			Method: http.MethodGet,
			URL:    "/my/test/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{test}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserOrOwnerAuth("test"))
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid superuser auth token",
			Method: http.MethodGet,
			URL:    "/my/test/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{id}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserOrOwnerAuth(""))
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireSameCollectionContextAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			URL:    "/my/test/_pb_users_auth_",
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{collection}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSameCollectionContextAuth(""))
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			URL:    "/my/test/_pb_users_auth_",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoxNjQwOTkxNjYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.2D3tmqPn3vc5LoqqCz8V-iCDVXo9soYiH0d32G7FQT4",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{collection}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSameCollectionContextAuth(""))
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid record auth token (different collection)",
			Method: http.MethodGet,
			URL:    "/my/test/clients",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{collection}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSameCollectionContextAuth(""))
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid record auth token (same collection)",
			Method: http.MethodGet,
			URL:    "/my/test/_pb_users_auth_",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{collection}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSameCollectionContextAuth(""))
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid record auth token (non-matching/missing collection param)",
			Method: http.MethodGet,
			URL:    "/my/test/_pb_users_auth_",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{id}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserOrOwnerAuth(""))
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "valid record auth token (matching custom collection param)",
			Method: http.MethodGet,
			URL:    "/my/test/_pb_users_auth_",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{test}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSuperuserOrOwnerAuth("test"))
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "superuser no exception check",
			Method: http.MethodGet,
			URL:    "/my/test/_pb_users_auth_",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				e.Router.GET("/my/test/{collection}", func(e *core.RequestEvent) error {
					return e.String(200, "test123")
				}).Bind(apis.RequireSameCollectionContextAuth(""))
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
