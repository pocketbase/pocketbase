package rest_test

import (
	"encoding/json"
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
		expectBody  string
		expectError bool
	}{
		{
			strings.NewReader(""),
			echo.MIMEApplicationJSON,
			`{}`,
			false,
		},
		{
			strings.NewReader(`{"test":"invalid`),
			echo.MIMEApplicationJSON,
			`{}`,
			true,
		},
		{
			strings.NewReader(`{"test":123}`),
			echo.MIMEApplicationJSON,
			`{"test":123}`,
			false,
		},
		{
			strings.NewReader(
				url.Values{
					"string":  []string{"str"},
					"stings":  []string{"str1", "str2", ""},
					"number":  []string{"-123"},
					"numbers": []string{"123", "456.789"},
					"bool":    []string{"true"},
					"bools":   []string{"true", "false"},
				}.Encode(),
			),
			echo.MIMEApplicationForm,
			`{"bool":true,"bools":[true,false],"number":-123,"numbers":[123,456.789],"stings":["str1","str2",""],"string":"str"}`,
			false,
		},
	}

	for i, scenario := range scenarios {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", scenario.body)
		req.Header.Set(echo.HeaderContentType, scenario.contentType)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		data := map[string]any{}
		err := rest.BindBody(c, &data)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("[%d] Expected hasErr %v, got %v", i, scenario.expectError, hasErr)
		}

		rawBody, err := json.Marshal(data)
		if err != nil {
			t.Errorf("[%d] Failed to marshal binded body: %v", i, err)

		}

		if scenario.expectBody != string(rawBody) {
			t.Errorf("[%d] Expected body \n%s, \ngot \n%s", i, scenario.expectBody, rawBody)
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
