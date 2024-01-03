package apis_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRequireGuestOnly(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "valid record token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireGuestOnly(),
					},
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid admin token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireGuestOnly(),
					},
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoxNjQwOTkxNjYxfQ.HqvpCpM0RAk3Qu9PfCMuZsk_DKh9UYuzFLwXBMTZd1w",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireGuestOnly(),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireGuestOnly(),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireRecordAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireRecordAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoxNjQwOTkxNjYxfQ.HqvpCpM0RAk3Qu9PfCMuZsk_DKh9UYuzFLwXBMTZd1w",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireRecordAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid admin token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireRecordAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireRecordAuth(),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid record token with collection not in the restricted list",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireRecordAuth("demo1", "demo2"),
					},
				})
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token with collection in the restricted list",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireRecordAuth("demo1", "demo2", "users"),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireSameContextRecordAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/users/test",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireSameContextRecordAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			Url:    "/my/users/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoxNjQwOTkxNjYxfQ.HqvpCpM0RAk3Qu9PfCMuZsk_DKh9UYuzFLwXBMTZd1w",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireSameContextRecordAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid admin token",
			Method: http.MethodGet,
			Url:    "/my/users/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireSameContextRecordAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token but from different collection",
			Method: http.MethodGet,
			Url:    "/my/users/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImdrMzkwcWVnczR5NDd3biIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoidjg1MXE0cjc5MHJoa25sIiwiZXhwIjoyMjA4OTg1MjYxfQ.q34IWXrRWsjLvbbVNRfAs_J4SoTHloNBfdGEiLmy-D8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireSameContextRecordAuth(),
					},
				})
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireRecordAuth(),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireAdminAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMzIwMH0.Gp_1b5WVhqjj2o3nJhNUlJmpdiwFLXN72LbMP-26gjA",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid admin token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminAuth(),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireAdminAuthOnlyIfAny(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest (while having at least 1 existing admin)",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminAuthOnlyIfAny(app),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "guest (while having 0 existing admins)",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// delete all admins
				_, err := app.Dao().DB().NewQuery("DELETE FROM {{_admins}}").Execute()
				if err != nil {
					t.Fatal(err)
				}

				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminAuthOnlyIfAny(app),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMzIwMH0.Gp_1b5WVhqjj2o3nJhNUlJmpdiwFLXN72LbMP-26gjA",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminAuthOnlyIfAny(app),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminAuthOnlyIfAny(app),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid admin token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminAuthOnlyIfAny(app),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireAdminOrRecordAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrRecordAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMzIwMH0.Gp_1b5WVhqjj2o3nJhNUlJmpdiwFLXN72LbMP-26gjA",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrRecordAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrRecordAuth(),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid record token with collection not in the restricted list",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrRecordAuth("demo1", "demo2", "clients"),
					},
				})
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token with collection in the restricted list",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrRecordAuth("demo1", "demo2", "users"),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid admin token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrRecordAuth(),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid admin token + restricted collections list (should be ignored)",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrRecordAuth("demo1", "demo2"),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRequireAdminOrOwnerAuth(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/test/4q1xlclmfloku33",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test/:id",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrOwnerAuth(""),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "expired/invalid token",
			Method: http.MethodGet,
			Url:    "/my/test/4q1xlclmfloku33",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoxNjQwOTkxNjYxfQ.HqvpCpM0RAk3Qu9PfCMuZsk_DKh9UYuzFLwXBMTZd1w",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test/:id",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrOwnerAuth(""),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token (different user)",
			Method: http.MethodGet,
			Url:    "/my/test/4q1xlclmfloku33",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImJnczgyMG4zNjF2ajFxZCIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.tW4NZWZ0mHBgvSZsQ0OOQhWajpUNFPCvNrOF9aCZLZs",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test/:id",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrOwnerAuth(""),
					},
				})
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token (different collection)",
			Method: http.MethodGet,
			Url:    "/my/test/4q1xlclmfloku33",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6ImdrMzkwcWVnczR5NDd3biIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoidjg1MXE0cjc5MHJoa25sIiwiZXhwIjoyMjA4OTg1MjYxfQ.q34IWXrRWsjLvbbVNRfAs_J4SoTHloNBfdGEiLmy-D8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test/:id",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrOwnerAuth(""),
					},
				})
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid record token (owner)",
			Method: http.MethodGet,
			Url:    "/my/test/4q1xlclmfloku33",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test/:id",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrOwnerAuth(""),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid admin token",
			Method: http.MethodGet,
			Url:    "/my/test/4q1xlclmfloku33",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test/:custom",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrOwnerAuth("custom"),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestLoadCollectionContext(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "missing collection",
			Method: http.MethodGet,
			Url:    "/my/missing",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.LoadCollectionContext(app),
					},
				})
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/demo1",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.LoadCollectionContext(app),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid record token",
			Method: http.MethodGet,
			Url:    "/my/demo1",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.LoadCollectionContext(app),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "valid admin token",
			Method: http.MethodGet,
			Url:    "/my/demo1",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.LoadCollectionContext(app),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "mismatched type",
			Method: http.MethodGet,
			Url:    "/my/demo1",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.LoadCollectionContext(app, "auth"),
					},
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "matched type",
			Method: http.MethodGet,
			Url:    "/my/users",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/:collection",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.LoadCollectionContext(app, "auth"),
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
