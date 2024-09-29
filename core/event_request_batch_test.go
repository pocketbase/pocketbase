package core_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestInternalRequestValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		request        core.InternalRequest
		expectedErrors []string
	}{
		{
			"empty struct",
			core.InternalRequest{},
			[]string{"method", "url"},
		},

		// method
		{
			"GET method",
			core.InternalRequest{URL: "test", Method: http.MethodGet},
			[]string{},
		},
		{
			"POST method",
			core.InternalRequest{URL: "test", Method: http.MethodPost},
			[]string{},
		},
		{
			"PUT method",
			core.InternalRequest{URL: "test", Method: http.MethodPut},
			[]string{},
		},
		{
			"PATCH method",
			core.InternalRequest{URL: "test", Method: http.MethodPatch},
			[]string{},
		},
		{
			"DELETE method",
			core.InternalRequest{URL: "test", Method: http.MethodDelete},
			[]string{},
		},
		{
			"unknown method",
			core.InternalRequest{URL: "test", Method: "unknown"},
			[]string{"method"},
		},

		// url
		{
			"url <= 2000",
			core.InternalRequest{URL: strings.Repeat("a", 2000), Method: http.MethodGet},
			[]string{},
		},
		{
			"url > 2000",
			core.InternalRequest{URL: strings.Repeat("a", 2001), Method: http.MethodGet},
			[]string{"url"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			tests.TestValidationErrors(t, s.request.Validate(), s.expectedErrors)
		})
	}
}
