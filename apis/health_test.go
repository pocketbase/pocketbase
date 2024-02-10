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
			Name:           "HEAD health status",
			Method:         http.MethodHead,
			Url:            "/api/health",
			ExpectedStatus: 200,
		},
		{
			Name:           "GET health status",
			Method:         http.MethodGet,
			Url:            "/api/health",
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"code":200`,
				`"data":{`,
				`"canBackup":true`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
