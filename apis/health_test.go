package apis_test

import (
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func TestHealthAPI(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:           "GET health status (guest)",
			Method:         http.MethodGet, // automatically matches also HEAD as a side-effect of the Go std mux
			URL:            "/api/health",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"code":200`,
				`"data":{}`,
			},
			NotExpectedContent: []string{
				"canBackup",
				"realIP",
				"possibleProxyHeader",
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "GET health status (regular user)",
			Method: http.MethodGet,
			URL:    "/api/health",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"code":200`,
				`"data":{}`,
			},
			NotExpectedContent: []string{
				"canBackup",
				"realIP",
				"possibleProxyHeader",
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "GET health status (superuser)",
			Method: http.MethodGet,
			URL:    "/api/health",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"code":200`,
				`"data":{`,
				`"canBackup":true`,
				`"realIP"`,
				`"possibleProxyHeader"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
