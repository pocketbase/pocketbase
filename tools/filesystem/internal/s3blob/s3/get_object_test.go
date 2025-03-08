package s3_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3/tests"
)

func TestS3GetObject(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method: http.MethodGet,
			URL:    "http://test_bucket.example.com/test_key",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Header: http.Header{
					"Last-Modified":       []string{"Mon, 01 Feb 2025 03:04:05 GMT"},
					"Cache-Control":       []string{"test_cache"},
					"Content-Disposition": []string{"test_disposition"},
					"Content-Encoding":    []string{"test_encoding"},
					"Content-Language":    []string{"test_language"},
					"Content-Type":        []string{"test_type"},
					"Content-Range":       []string{"test_range"},
					"Etag":                []string{"test_etag"},
					"Content-Length":      []string{"100"},
					"x-amz-meta-AbC":      []string{"test_meta_a"},
					"x-amz-meta-Def":      []string{"test_meta_b"},
				},
				Body: io.NopCloser(strings.NewReader("test")),
			},
		},
	)

	s3Client := &s3.S3{
		Client:    httpClient,
		Region:    "test_region",
		Bucket:    "test_bucket",
		Endpoint:  "http://example.com",
		AccessKey: "123",
		SecretKey: "abc",
	}

	resp, err := s3Client.GetObject(context.Background(), "test_key", func(r *http.Request) {
		r.Header.Set("test_header", "test")
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}

	// check body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	bodyStr := string(body)

	if bodyStr != "test" {
		t.Fatalf("Expected body\n%q\ngot\n%q", "test", bodyStr)
	}

	// check serialized attributes
	raw, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)

	expected := `{"metadata":{"abc":"test_meta_a","def":"test_meta_b"},"lastModified":"2025-02-01T03:04:05Z","cacheControl":"test_cache","contentDisposition":"test_disposition","contentEncoding":"test_encoding","contentLanguage":"test_language","contentType":"test_type","contentRange":"test_range","etag":"test_etag","contentLength":100}`

	if rawStr != expected {
		t.Fatalf("Expected attributes\n%s\ngot\n%s", expected, rawStr)
	}
}
