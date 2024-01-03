package apis_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tests"
)

func TestLogsList(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/logs",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as auth record",
			Method: http.MethodGet,
			Url:    "/api/logs",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/logs",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"873f2133-9f38-44fb-bf82-c8f53b310d91"`,
				`"id":"f2133873-44fb-9f38-bf82-c918f53b310d"`,
			},
		},
		{
			Name:   "authorized as admin + filter",
			Method: http.MethodGet,
			Url:    "/api/logs?filter=data.status>200",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":1`,
				`"items":[{`,
				`"id":"f2133873-44fb-9f38-bf82-c918f53b310d"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestLogView(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/logs/873f2133-9f38-44fb-bf82-c8f53b310d91",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as auth record",
			Method: http.MethodGet,
			Url:    "/api/logs/873f2133-9f38-44fb-bf82-c8f53b310d91",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (nonexisting request log)",
			Method: http.MethodGet,
			Url:    "/api/logs/missing1-9f38-44fb-bf82-c8f53b310d91",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (existing request log)",
			Method: http.MethodGet,
			Url:    "/api/logs/873f2133-9f38-44fb-bf82-c8f53b310d91",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"873f2133-9f38-44fb-bf82-c8f53b310d91"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestLogsStats(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/logs/stats",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as auth record",
			Method: http.MethodGet,
			Url:    "/api/logs/stats",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/logs/stats",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`[{"total":1,"date":"2022-05-01 10:00:00.000Z"},{"total":1,"date":"2022-05-02 10:00:00.000Z"}]`,
			},
		},
		{
			Name:   "authorized as admin + filter",
			Method: http.MethodGet,
			Url:    "/api/logs/stats?filter=data.status>200",
			RequestHeaders: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			},
			BeforeTestFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`[{"total":1,"date":"2022-05-02 10:00:00.000Z"}]`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
