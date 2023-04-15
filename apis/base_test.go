package apis_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/spf13/cast"
)

func Test404(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Method:          http.MethodGet,
			Url:             "/api/missing",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Method:          http.MethodPost,
			Url:             "/api/missing",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Method:          http.MethodPatch,
			Url:             "/api/missing",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Method:          http.MethodDelete,
			Url:             "/api/missing",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Method:         http.MethodHead,
			Url:            "/api/missing",
			ExpectedStatus: 404,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestCustomRoutesAndErrorsHandling(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:   "custom route",
			Method: http.MethodGet,
			Url:    "/custom",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/custom",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "custom route with url encoded parameter",
			Method: http.MethodGet,
			Url:    "/a%2Bb%2Bc",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/:param",
					Handler: func(c echo.Context) error {
						return c.String(200, c.PathParam("param"))
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"a+b+c"},
		},
		{
			Name:   "route with HTTPError",
			Method: http.MethodGet,
			Url:    "/http-error",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/http-error",
					Handler: func(c echo.Context) error {
						return echo.ErrBadRequest
					},
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`{"code":400,"message":"Bad Request.","data":{}}`},
		},
		{
			Name:   "route with api error",
			Method: http.MethodGet,
			Url:    "/api-error",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/api-error",
					Handler: func(c echo.Context) error {
						return apis.NewApiError(500, "test message", errors.New("internal_test"))
					},
				})
			},
			ExpectedStatus:  500,
			ExpectedContent: []string{`{"code":500,"message":"Test message.","data":{}}`},
		},
		{
			Name:   "route with plain error",
			Method: http.MethodGet,
			Url:    "/plain-error",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/plain-error",
					Handler: func(c echo.Context) error {
						return errors.New("Test error")
					},
				})
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`{"code":400,"message":"Something went wrong while processing your request.","data":{}}`},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRemoveTrailingSlashMiddleware(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:   "non /api/* route (exact match)",
			Method: http.MethodGet,
			Url:    "/custom",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/custom",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "non /api/* route (with trailing slash)",
			Method: http.MethodGet,
			Url:    "/custom/",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/custom",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
				})
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "/api/* route (exact match)",
			Method: http.MethodGet,
			Url:    "/api/custom",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/api/custom",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
		{
			Name:   "/api/* route (with trailing slash)",
			Method: http.MethodGet,
			Url:    "/api/custom/",
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: http.MethodGet,
					Path:   "/api/custom",
					Handler: func(c echo.Context) error {
						return c.String(200, "test123")
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

func TestEagerRequestDataCache(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:   "[UNKNOWN] unsupported eager cached request method",
			Method: "UNKNOWN",
			Url:    "/custom",
			Body:   strings.NewReader(`{"name":"test123"}`),
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				e.AddRoute(echo.Route{
					Method: "UNKNOWN",
					Path:   "/custom",
					Handler: func(c echo.Context) error {
						data := &struct {
							Name string `json:"name"`
						}{}

						if err := c.Bind(data); err != nil {
							return err
						}

						// since the unknown method is not eager cache support
						// it should fail reading the json body twice
						r := apis.RequestData(c)
						if v := cast.ToString(r.Data["name"]); v != "" {
							t.Fatalf("Expected empty request data body, got, %v", r.Data)
						}

						return c.String(200, data.Name)
					},
				})
			},
			ExpectedStatus:  200,
			ExpectedContent: []string{"test123"},
		},
	}

	// supported eager cache request methods
	supportedMethods := []string{"POST", "PUT", "PATCH", "DELETE"}
	for _, m := range supportedMethods {
		scenarios = append(
			scenarios,
			tests.ApiScenario{
				Name:   fmt.Sprintf("[%s] valid cached json body request", m),
				Method: http.MethodPost,
				Url:    "/custom",
				Body:   strings.NewReader(`{"name":"test123"}`),
				BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
					e.AddRoute(echo.Route{
						Method: http.MethodPost,
						Path:   "/custom",
						Handler: func(c echo.Context) error {
							data := &struct {
								Name string `json:"name"`
							}{}

							if err := c.Bind(data); err != nil {
								return err
							}

							// try to read the body again
							r := apis.RequestData(c)
							if v := cast.ToString(r.Data["name"]); v != "test123" {
								t.Fatalf("Expected request data with name %q, got, %q", "test123", v)
							}

							return c.String(200, data.Name)
						},
					})
				},
				ExpectedStatus:  200,
				ExpectedContent: []string{"test123"},
			},
		)
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
