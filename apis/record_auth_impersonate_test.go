package apis_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func TestRecordAuthImpersonate(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			URL:             "/api/collections/users/impersonate/4q1xlclmfloku33",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as different user",
			Method: http.MethodPost,
			URL:    "/api/collections/users/impersonate/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6Im9hcDY0MGNvdDR5cnUycyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.GfJo6EHIobgas_AXt-M-tj5IoQendPnrkMSe9ExuSEY",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as the same user",
			Method: http.MethodPost,
			URL:    "/api/collections/users/impersonate/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser",
			Method: http.MethodPost,
			URL:    "/api/collections/users/impersonate/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":"`,
				`"id":"4q1xlclmfloku33"`,
				`"record":{`,
			},
			NotExpectedContent: []string{
				// hidden fields should remain hidden even though we are authenticated as superuser
				`"tokenKey"`,
				`"password"`,
			},
			ExpectedEvents: map[string]int{
				"*":                   0,
				"OnRecordAuthRequest": 1,
				"OnRecordEnrich":      1,
			},
		},
		{
			Name:   "authorized as superuser with custom invalid duration",
			Method: http.MethodPost,
			URL:    "/api/collections/users/impersonate/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body:           strings.NewReader(`{"duration":-1}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"duration":{`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser with custom valid duration",
			Method: http.MethodPost,
			URL:    "/api/collections/users/impersonate/4q1xlclmfloku33",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body:           strings.NewReader(`{"duration":100}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"token":"`,
				`"id":"4q1xlclmfloku33"`,
				`"record":{`,
			},
			ExpectedEvents: map[string]int{
				"*":                   0,
				"OnRecordAuthRequest": 1,
				"OnRecordEnrich":      1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
