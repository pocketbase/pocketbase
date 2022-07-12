package apis_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestAdminAuth(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/admins/auth-via-email",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."},"password":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/admins/auth-via-email",
			Body:            strings.NewReader(`{`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "wrong email/password",
			Method:          http.MethodPost,
			Url:             "/api/admins/auth-via-email",
			Body:            strings.NewReader(`{"email":"missing@example.com","password":"wrong_pass"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "valid email/password (already authorized)",
			Method: http.MethodPost,
			Url:    "/api/admins/auth-via-email",
			Body:   strings.NewReader(`{"email":"test@example.com","password":"1234567890"}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message":"The request can be accessed only by guests.","data":{}`},
		},
		{
			Name:           "valid email/password (guest)",
			Method:         http.MethodPost,
			Url:            "/api/admins/auth-via-email",
			Body:           strings.NewReader(`{"email":"test@example.com","password":"1234567890"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"admin":{"id":"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c"`,
				`"token":`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminAuthRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminRequestPasswordReset(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/admins/request-password-reset",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/admins/request-password-reset",
			Body:            strings.NewReader(`{"email`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:           "missing admin",
			Method:         http.MethodPost,
			Url:            "/api/admins/request-password-reset",
			Body:           strings.NewReader(`{"email":"missing@example.com"}`),
			ExpectedStatus: 204,
		},
		{
			Name:           "existing admin",
			Method:         http.MethodPost,
			Url:            "/api/admins/request-password-reset",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			ExpectedStatus: 204,
			// usually this events are fired but since the submit is
			// executed in a separate go routine they are fired async
			// ExpectedEvents: map[string]int{
			// 	"OnModelBeforeUpdate": 1,
			// 	"OnModelAfterUpdate":  1,
			// 	"OnMailerBeforeUserResetPasswordSend:1":  1,
			// 	"OnMailerAfterUserResetPasswordSend:1":  1,
			// },
		},
		{
			Name:           "existing admin (after already sent)",
			Method:         http.MethodPost,
			Url:            "/api/admins/request-password-reset",
			Body:           strings.NewReader(`{"email":"test@example.com"}`),
			ExpectedStatus: 204,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminConfirmPasswordReset(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodPost,
			Url:             "/api/admins/confirm-password-reset",
			Body:            strings.NewReader(``),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"password":{"code":"validation_required","message":"Cannot be blank."},"passwordConfirm":{"code":"validation_required","message":"Cannot be blank."},"token":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:            "invalid data",
			Method:          http.MethodPost,
			Url:             "/api/admins/confirm-password-reset",
			Body:            strings.NewReader(`{"password`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "expired token",
			Method:          http.MethodPost,
			Url:             "/api/admins/confirm-password-reset",
			Body:            strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMzIwMH0.Gp_1b5WVhqjj2o3nJhNUlJmpdiwFLXN72LbMP-26gjA","password":"1234567890","passwordConfirm":"1234567890"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"token":{"code":"validation_invalid_token","message":"Invalid or expired token."}}}`},
		},
		{
			Name:           "valid token",
			Method:         http.MethodPost,
			Url:            "/api/admins/confirm-password-reset",
			Body:           strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg5MzQ3NDAwMH0.72IhlL_5CpNGE0ZKM7sV9aAKa3wxQaMZdDiHBo0orpw","password":"1234567890","passwordConfirm":"1234567890"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"admin":{"id":"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c"`,
				`"token":`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate": 1,
				"OnModelAfterUpdate":  1,
				"OnAdminAuthRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminRefresh(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			Url:             "/api/admins/refresh",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPost,
			Url:    "/api/admins/refresh",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodPost,
			Url:    "/api/admins/refresh",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"admin":{"id":"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c"`,
				`"token":`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminAuthRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminsList(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/admins",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/admins",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/admins",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c"`,
				`"id":"3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8"`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminsListRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + paging and sorting",
			Method: http.MethodGet,
			Url:    "/api/admins?page=2&perPage=1&sort=-created",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":2`,
				`"perPage":1`,
				`"totalItems":2`,
				`"items":[{`,
				`"id":"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c"`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminsListRequest": 1,
			},
		},
		{
			Name:   "authorized as admin + invalid filter",
			Method: http.MethodGet,
			Url:    "/api/admins?filter=invalidfield~'test2'",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + valid filter",
			Method: http.MethodGet,
			Url:    "/api/admins?filter=email~'test2'",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"page":1`,
				`"perPage":30`,
				`"totalItems":1`,
				`"items":[{`,
				`"id":"3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8"`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminsListRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminView(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + invalid admin id",
			Method: http.MethodGet,
			Url:    "/api/admins/invalid",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting admin id",
			Method: http.MethodGet,
			Url:    "/api/admins/b97ccf83-34a2-4d01-a26b-3d77bc842d3c",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + existing admin id",
			Method: http.MethodGet,
			Url:    "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8"`,
			},
			ExpectedEvents: map[string]int{
				"OnAdminViewRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminDelete(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodDelete,
			Url:             "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodDelete,
			Url:    "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + invalid admin id",
			Method: http.MethodDelete,
			Url:    "/api/admins/invalid",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting admin id",
			Method: http.MethodDelete,
			Url:    "/api/admins/b97ccf83-34a2-4d01-a26b-3d77bc842d3c",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + existing admin id",
			Method: http.MethodDelete,
			Url:    "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 204,
			ExpectedEvents: map[string]int{
				"OnModelBeforeDelete":        1,
				"OnModelAfterDelete":         1,
				"OnAdminBeforeDeleteRequest": 1,
				"OnAdminAfterDeleteRequest":  1,
			},
		},
		{
			Name:   "authorized as admin - try to delete the only remaining admin",
			Method: http.MethodDelete,
			Url:    "/api/admins/2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// delete all admins except the authorized one
				adminModel := &models.Admin{}
				_, err := app.Dao().DB().Delete(adminModel.TableName(), dbx.Not(dbx.HashExp{
					"id": "2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
				})).Execute()
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents: map[string]int{
				"OnAdminBeforeDeleteRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminCreate(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized (while having at least 1 existing admin)",
			Method:          http.MethodPost,
			Url:             "/api/admins",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "unauthorized (while having 0 existing admins)",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body:   strings.NewReader(`{"email":"testnew@example.com","password":"1234567890","passwordConfirm":"1234567890","avatar":3}`),
			BeforeFunc: func(t *testing.T, app *tests.TestApp, e *echo.Echo) {
				// delete all admins
				_, err := app.Dao().DB().NewQuery("DELETE FROM {{_admins}}").Execute()
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"email":"testnew@example.com"`,
				`"avatar":3`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":        1,
				"OnModelAfterCreate":         1,
				"OnAdminBeforeCreateRequest": 1,
				"OnAdminAfterCreateRequest":  1,
			},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPost,
			Url:    "/api/admins",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty data",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"email":{"code":"validation_required","message":"Cannot be blank."},"password":{"code":"validation_required","message":"Cannot be blank."}}`},
		},
		{
			Name:   "authorized as admin + invalid data format",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body:   strings.NewReader(`{`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + invalid data",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body:   strings.NewReader(`{"email":"test@example.com","password":"1234","passwordConfirm":"4321","avatar":99}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"avatar":{"code":"validation_max_less_equal_than_required","message":"Must be no greater than 9."},"email":{"code":"validation_admin_email_exists","message":"Admin email already exists."},"password":{"code":"validation_length_out_of_range","message":"The length must be between 10 and 100."},"passwordConfirm":{"code":"validation_values_mismatch","message":"Values don't match."}}`},
		},
		{
			Name:   "authorized as admin + valid data",
			Method: http.MethodPost,
			Url:    "/api/admins",
			Body:   strings.NewReader(`{"email":"testnew@example.com","password":"1234567890","passwordConfirm":"1234567890","avatar":3}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":`,
				`"email":"testnew@example.com"`,
				`"avatar":3`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeCreate":        1,
				"OnModelAfterCreate":         1,
				"OnAdminBeforeCreateRequest": 1,
				"OnAdminAfterCreateRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAdminUpdate(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPatch,
			Url:             "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPatch,
			Url:    "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + invalid admin id",
			Method: http.MethodPatch,
			Url:    "/api/admins/invalid",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + nonexisting admin id",
			Method: http.MethodPatch,
			Url:    "/api/admins/b97ccf83-34a2-4d01-a26b-3d77bc842d3c",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + empty data",
			Method: http.MethodPatch,
			Url:    "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8"`,
				`"email":"test2@example.com"`,
				`"avatar":2`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":        1,
				"OnModelAfterUpdate":         1,
				"OnAdminBeforeUpdateRequest": 1,
				"OnAdminAfterUpdateRequest":  1,
			},
		},
		{
			Name:   "authorized as admin + invalid formatted data",
			Method: http.MethodPatch,
			Url:    "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			Body:   strings.NewReader(`{`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin + invalid data",
			Method: http.MethodPatch,
			Url:    "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			Body:   strings.NewReader(`{"email":"test@example.com","password":"1234","passwordConfirm":"4321","avatar":99}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"data":{"avatar":{"code":"validation_max_less_equal_than_required","message":"Must be no greater than 9."},"email":{"code":"validation_admin_email_exists","message":"Admin email already exists."},"password":{"code":"validation_length_out_of_range","message":"The length must be between 10 and 100."},"passwordConfirm":{"code":"validation_values_mismatch","message":"Values don't match."}}`},
		},
		{
			Method: http.MethodPatch,
			Url:    "/api/admins/3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8",
			Body:   strings.NewReader(`{"email":"testnew@example.com","password":"1234567890","passwordConfirm":"1234567890","avatar":5}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"id":"3f8397cc-2b4a-a26b-4d01-42d3c3d77bc8"`,
				`"email":"testnew@example.com"`,
				`"avatar":5`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":        1,
				"OnModelAfterUpdate":         1,
				"OnAdminBeforeUpdateRequest": 1,
				"OnAdminAfterUpdateRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
