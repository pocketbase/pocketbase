package s3_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3"
)

func TestResponseErrorSerialization(t *testing.T) {
	raw := `
		<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		  <Code>test_code</Code>
		  <Message>test_message</Message>
		  <RequestId>test_request_id</RequestId>
		  <Resource>test_resource</Resource>
		</Error>
	`

	respErr := &s3.ResponseError{
		Status: 123,
		Raw:    []byte("test"),
	}

	err := xml.Unmarshal([]byte(raw), &respErr)
	if err != nil {
		t.Fatal(err)
	}

	jsonRaw, err := json.Marshal(respErr)
	if err != nil {
		t.Fatal(err)
	}
	jsonStr := string(jsonRaw)

	expected := `{"code":"test_code","message":"test_message","requestId":"test_request_id","resource":"test_resource","status":123}`

	if expected != jsonStr {
		t.Fatalf("Expected JSON\n%s\ngot\n%s", expected, jsonStr)
	}
}

func TestResponseErrorErrorInterface(t *testing.T) {
	scenarios := []struct {
		name     string
		err      *s3.ResponseError
		expected string
	}{
		{
			"empty",
			&s3.ResponseError{},
			"0 S3ResponseError",
		},
		{
			"with code and message (nil raw)",
			&s3.ResponseError{
				Status:  123,
				Code:    "test_code",
				Message: "test_message",
			},
			"123 test_code: test_message",
		},
		{
			"with code and message (non-nil raw)",
			&s3.ResponseError{
				Status:  123,
				Code:    "test_code",
				Message: "test_message",
				Raw:     []byte("test_raw"),
			},
			"123 test_code: test_message\n(RAW: test_raw)",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.err.Error()

			if result != s.expected {
				t.Fatalf("Expected\n%s\ngot\n%s", s.expected, result)
			}
		})
	}
}
