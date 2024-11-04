package apis_test

import (
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestLogsList(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			URL:             "/api/logs",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodGet,
			URL:    "/api/logs",
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
			URL:    "/api/logs",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				if err := tests.StubLogsData(app); err != nil {
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
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser + filter",
			Method: http.MethodGet,
			URL:    "/api/logs?filter=data.status>200",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				if err := tests.StubLogsData(app); err != nil {
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
			ExpectedEvents: map[string]int{"*": 0},
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
			URL:             "/api/logs/873f2133-9f38-44fb-bf82-c8f53b310d91",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodGet,
			URL:    "/api/logs/873f2133-9f38-44fb-bf82-c8f53b310d91",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (nonexisting request log)",
			Method: http.MethodGet,
			URL:    "/api/logs/missing1-9f38-44fb-bf82-c8f53b310d91",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				if err := tests.StubLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser (existing request log)",
			Method: http.MethodGet,
			URL:    "/api/logs/873f2133-9f38-44fb-bf82-c8f53b310d91",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				if err := tests.StubLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"873f2133-9f38-44fb-bf82-c8f53b310d91"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
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
			URL:             "/api/logs/stats",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodGet,
			URL:    "/api/logs/stats",
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
			URL:    "/api/logs/stats",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				if err := tests.StubLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`[{"date":"2022-05-01 10:00:00.000Z","total":1},{"date":"2022-05-02 10:00:00.000Z","total":1}]`,
			},
		},
		{
			Name:   "authorized as superuser + filter",
			Method: http.MethodGet,
			URL:    "/api/logs/stats?filter=data.status>200",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				if err := tests.StubLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`[{"date":"2022-05-02 10:00:00.000Z","total":1}]`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
