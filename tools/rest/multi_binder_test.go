package rest_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tools/rest"
)

func TestBindBody(t *testing.T) {
	scenarios := []struct {
		body        io.Reader
		contentType string
		result      map[string]string
		expectError bool
	}{
		{
			strings.NewReader(""),
			echo.MIMEApplicationJSON,
			map[string]string{},
			false,
		},
		{
			strings.NewReader(`{"test":"invalid`),
			echo.MIMEApplicationJSON,
			map[string]string{},
			true,
		},
		{
			strings.NewReader(`{"test":"test123"}`),
			echo.MIMEApplicationJSON,
			map[string]string{"test": "test123"},
			false,
		},
		{
			strings.NewReader(url.Values{"test": []string{"test123"}}.Encode()),
			echo.MIMEApplicationForm,
			map[string]string{"test": "test123"},
			false,
		},
	}

	for i, scenario := range scenarios {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", scenario.body)
		req.Header.Set(echo.HeaderContentType, scenario.contentType)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		result := map[string]string{}
		err := rest.BindBody(c, &result)

		if err == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}

		if err != nil && !scenario.expectError {
			t.Errorf("(%d) Expected nil, got error %v", i, err)
		}

		if len(result) != len(scenario.result) {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.result, result)
		}

		for k, v := range result {
			if sv, ok := scenario.result[k]; !ok || v != sv {
				t.Errorf("(%d) Expected value %v for key %s, got %v", i, sv, k, v)
			}
		}
	}
}

func TestCopyJsonBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"test":"test123"}`))

	// simulate multiple reads from the same request
	result1 := map[string]string{}
	rest.CopyJsonBody(req, &result1)
	result2 := map[string]string{}
	rest.CopyJsonBody(req, &result2)

	if len(result1) == 0 {
		t.Error("Expected result1 to be filled")
	}

	if len(result2) == 0 {
		t.Error("Expected result2 to be filled")
	}

	if v, ok := result1["test"]; !ok || v != "test123" {
		t.Errorf("Expected result1.test to be %q, got %q", "test123", v)
	}

	if v, ok := result2["test"]; !ok || v != "test123" {
		t.Errorf("Expected result2.test to be %q, got %q", "test123", v)
	}
}
