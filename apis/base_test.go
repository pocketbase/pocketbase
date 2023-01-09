package apis_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tests"
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
