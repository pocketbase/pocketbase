package apis_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func TestSettingsList(t *testing.T) {
	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodGet,
			Url:             "/api/settings",
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodGet,
			Url:    "/api/settings",
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin",
			Method: http.MethodGet,
			Url:    "/api/settings",
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"userAuthToken":{`,
				`"userPasswordResetToken":{`,
				`"userEmailChangeToken":{`,
				`"userVerificationToken":{`,
				`"emailAuth":{`,
				`"googleAuth":{`,
				`"facebookAuth":{`,
				`"githubAuth":{`,
				`"gitlabAuth":{`,
				`"secret":"******"`,
				`"clientSecret":"******"`,
			},
			ExpectedEvents: map[string]int{
				"OnSettingsListRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestSettingsSet(t *testing.T) {
	validData := `{"meta":{"appName":"update_test"},"emailAuth":{"minPasswordLength": 12}}`

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPatch,
			Url:             "/api/settings",
			Body:            strings.NewReader(validData),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as user",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(validData),
			RequestHeaders: map[string]string{
				"Authorization": "User eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			},
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:   "authorized as admin submitting empty data",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(``),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"userAuthToken":{`,
				`"userPasswordResetToken":{`,
				`"userEmailChangeToken":{`,
				`"userVerificationToken":{`,
				`"emailAuth":{`,
				`"googleAuth":{`,
				`"facebookAuth":{`,
				`"githubAuth":{`,
				`"gitlabAuth":{`,
				`"secret":"******"`,
				`"clientSecret":"******"`,
				`"appName":"Acme"`,
				`"minPasswordLength":8`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":           1,
				"OnModelAfterUpdate":            1,
				"OnSettingsBeforeUpdateRequest": 1,
				"OnSettingsAfterUpdateRequest":  1,
			},
		},
		{
			Name:   "authorized as admin submitting invalid data",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(`{"meta":{"appName":""},"emailAuth":{"minPasswordLength": 3}}`),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"data":{`,
				`"emailAuth":{"minPasswordLength":{"code":"validation_min_greater_equal_than_required","message":"Must be no less than 5."}}`,
				`"meta":{"appName":{"code":"validation_required","message":"Cannot be blank."}}`,
			},
		},
		{
			Name:   "authorized as admin submitting valid data",
			Method: http.MethodPatch,
			Url:    "/api/settings",
			Body:   strings.NewReader(validData),
			RequestHeaders: map[string]string{
				"Authorization": "Admin eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjJiNGE5N2NjLTNmODMtNGQwMS1hMjZiLTNkNzdiYzg0MmQzYyIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTg3MzQ2Mjc5Mn0.AtRtXR6FHBrCUGkj5OffhmxLbSZaQ4L_Qgw4gfoHyfo",
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"meta":{`,
				`"logs":{`,
				`"smtp":{`,
				`"s3":{`,
				`"adminAuthToken":{`,
				`"adminPasswordResetToken":{`,
				`"userAuthToken":{`,
				`"userPasswordResetToken":{`,
				`"userEmailChangeToken":{`,
				`"userVerificationToken":{`,
				`"emailAuth":{`,
				`"googleAuth":{`,
				`"facebookAuth":{`,
				`"githubAuth":{`,
				`"gitlabAuth":{`,
				`"secret":"******"`,
				`"clientSecret":"******"`,
				`"appName":"update_test"`,
				`"minPasswordLength":12`,
			},
			ExpectedEvents: map[string]int{
				"OnModelBeforeUpdate":           1,
				"OnModelAfterUpdate":            1,
				"OnSettingsBeforeUpdateRequest": 1,
				"OnSettingsAfterUpdateRequest":  1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
