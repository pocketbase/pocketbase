package apis_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRequestsList(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/logs/requests",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/logs/requests",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/logs/requests",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockRequestLogsData(app); err != nil {
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
			Url:    "/api/logs/requests?filter=status>200",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockRequestLogsData(app); err != nil {
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

func TestRequestView(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/logs/requests/873f2133-9f38-44fb-bf82-c8f53b310d91",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/logs/requests/873f2133-9f38-44fb-bf82-c8f53b310d91",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (nonexisting request log)",
			Method: http.MethodGet,
			Url:    "/api/logs/requests/missing1-9f38-44fb-bf82-c8f53b310d91",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockRequestLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin (existing request log)",
			Method: http.MethodGet,
			Url:    "/api/logs/requests/873f2133-9f38-44fb-bf82-c8f53b310d91",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockRequestLogsData(app); err != nil {
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

func TestRequestsStats(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/logs/requests/stats",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/logs/requests/stats",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/logs/requests/stats",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockRequestLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`[{"total":1,"date":"2022-05-01 10:00:00.000"},{"total":1,"date":"2022-05-02 10:00:00.000"}]`,
			},
		},
		{
			Name:   "authorized as admin + filter",
			Method: http.MethodGet,
			Url:    "/api/logs/requests/stats?filter=status>200",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				if err := tests.MockRequestLogsData(app); err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`[{"total":1,"date":"2022-05-02 10:00:00.000"}]`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
