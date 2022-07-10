package apis_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRequireGuestOnly(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:   "valid user token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxNjQwOTkxNjYxfQ.HkAldxpbn0EybkMfFGQKEJUIYKE5UJA0AjcsrV7Q6Io",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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

func TestRequireUserAuth(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireUserAuth(),
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
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxNjQwOTkxNjYxfQ.HkAldxpbn0EybkMfFGQKEJUIYKE5UJA0AjcsrV7Q6Io",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireUserAuth(),
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireUserAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid user token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireUserAuth(),
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
	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMzIwMH0.Gp_1b5WVhqjj2o3nJhNUlJmpdiwFLXN72LbMP-26gjA",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
			Name:   "valid user token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
	scenarios := []tests.ApiScenario{
		{
			Name:   "guest (while having at least 1 existing admin)",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMzIwMH0.Gp_1b5WVhqjj2o3nJhNUlJmpdiwFLXN72LbMP-26gjA",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
			Name:   "valid user token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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

func TestRequireAdminOrUserAuth(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/test",
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrUserAuth(),
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMzIwMH0.Gp_1b5WVhqjj2o3nJhNUlJmpdiwFLXN72LbMP-26gjA",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrUserAuth(),
					},
				})
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid user token",
			Method: http.MethodGet,
			Url:    "/my/test",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrUserAuth(),
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
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/my/test",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
					Middlewares: []echo.MiddlewareFunc{
						apis.RequireAdminOrUserAuth(),
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
	scenarios := []tests.ApiScenario{
		{
			Name:   "guest",
			Method: http.MethodGet,
			Url:    "/my/test/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
			Url:    "/my/test/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxNjQwOTkxNjYxfQ.HkAldxpbn0EybkMfFGQKEJUIYKE5UJA0AjcsrV7Q6Io",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
			Name:   "valid user token (different user)",
			Method: http.MethodGet,
			Url:    "/my/test/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			RequestHeaders: map[string]string{
				// test3@example.com
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoidXNlciIsImVtYWlsIjoidGVzdDNAZXhhbXBsZS5jb20iLCJpZCI6Ijk3Y2MzZDNkLTZiYTItMzgzZi1iNDJhLTdiYzg0ZDI3NDEwYyIsImV4cCI6MTg5MzUxNTU3Nn0.Q965uvlTxxOsZbACXSgJQNXykYK0TKZ87nyPzemvN4E",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
			Name:   "valid user token (owner)",
			Method: http.MethodGet,
			Url:    "/my/test/4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
			Url:    "/my/test/2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
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
