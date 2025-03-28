package s3_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3/tests"
)

func TestS3URL(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name     string
		s3Client *s3.S3
		expected string
	}{
		{
			"no schema",
			&s3.S3{
				Region:    "test_region",
				Bucket:    "test_bucket",
				Endpoint:  "example.com/",
				AccessKey: "123",
				SecretKey: "abc",
			},
			"https://test_bucket.example.com/test_key/a/b/c?q=1",
		},
		{
			"with https schema",
			&s3.S3{
				Region:    "test_region",
				Bucket:    "test_bucket",
				Endpoint:  "https://example.com/",
				AccessKey: "123",
				SecretKey: "abc",
			},
			"https://test_bucket.example.com/test_key/a/b/c?q=1",
		},
		{
			"with http schema",
			&s3.S3{
				Region:    "test_region",
				Bucket:    "test_bucket",
				Endpoint:  "http://example.com/",
				AccessKey: "123",
				SecretKey: "abc",
			},
			"http://test_bucket.example.com/test_key/a/b/c?q=1",
		},
		{
			"path style addressing (non-explicit schema)",
			&s3.S3{
				Region:       "test_region",
				Bucket:       "test_bucket",
				Endpoint:     "example.com/",
				AccessKey:    "123",
				SecretKey:    "abc",
				UsePathStyle: true,
			},
			"https://example.com/test_bucket/test_key/a/b/c?q=1",
		},
		{
			"path style addressing (explicit schema)",
			&s3.S3{
				Region:       "test_region",
				Bucket:       "test_bucket",
				Endpoint:     "http://example.com/",
				AccessKey:    "123",
				SecretKey:    "abc",
				UsePathStyle: true,
			},
			"http://example.com/test_bucket/test_key/a/b/c?q=1",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.s3Client.URL("/test_key/a/b/c?q=1")
			if result != s.expected {
				t.Fatalf("Expected URL\n%s\ngot\n%s", s.expected, result)
			}
		})
	}
}

func TestS3SignAndSend(t *testing.T) {
	t.Parallel()

	testResponse := func() *http.Response {
		return &http.Response{
			Body: io.NopCloser(strings.NewReader("test_response")),
		}
	}

	scenarios := []struct {
		name     string
		path     string
		reqFunc  func(req *http.Request)
		s3Client *s3.S3
	}{
		{
			"minimal",
			"/test",
			func(req *http.Request) {
				req.Header.Set("x-amz-date", "20250102T150405Z")
			},
			&s3.S3{
				Region:    "test_region",
				Bucket:    "test_bucket",
				Endpoint:  "https://example.com/",
				AccessKey: "123",
				SecretKey: "abc",
				Client: tests.NewClient(&tests.RequestStub{
					Method:   http.MethodGet,
					URL:      "https://test_bucket.example.com/test",
					Response: testResponse(),
					Match: func(req *http.Request) bool {
						return tests.ExpectHeaders(req.Header, map[string]string{
							"Authorization":        "AWS4-HMAC-SHA256 Credential=123/20250102/test_region/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=ea093662bc1deef08dfb4ac35453dfaad5ea89edf102e9dd3b7156c9a27e4c1f",
							"Host":                 "test_bucket.example.com",
							"X-Amz-Content-Sha256": "UNSIGNED-PAYLOAD",
							"X-Amz-Date":           "20250102T150405Z",
						})
					},
				}),
			},
		},
		{
			"minimal with different access and secret keys",
			"/test",
			func(req *http.Request) {
				req.Header.Set("x-amz-date", "20250102T150405Z")
			},
			&s3.S3{
				Region:    "test_region",
				Bucket:    "test_bucket",
				Endpoint:  "https://example.com/",
				AccessKey: "456",
				SecretKey: "def",
				Client: tests.NewClient(&tests.RequestStub{
					Method:   http.MethodGet,
					URL:      "https://test_bucket.example.com/test",
					Response: testResponse(),
					Match: func(req *http.Request) bool {
						return tests.ExpectHeaders(req.Header, map[string]string{
							"Authorization":        "AWS4-HMAC-SHA256 Credential=456/20250102/test_region/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=17510fa1f724403dd0a563b61c9b31d1d718f877fcbd75455620d17a8afce5fb",
							"Host":                 "test_bucket.example.com",
							"X-Amz-Content-Sha256": "UNSIGNED-PAYLOAD",
							"X-Amz-Date":           "20250102T150405Z",
						})
					},
				}),
			},
		},
		{
			"minimal with special characters",
			"/ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~!@#$^&*()=/@sub?a=1&@b=@2",
			func(req *http.Request) {
				req.Header.Set("x-amz-date", "20250102T150405Z")
			},
			&s3.S3{
				Region:    "test_region",
				Bucket:    "test_bucket",
				Endpoint:  "https://example.com/",
				AccessKey: "456",
				SecretKey: "def",
				Client: tests.NewClient(&tests.RequestStub{
					Method:   http.MethodGet,
					URL:      "https://test_bucket.example.com/ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.~!@#$%5E&*()=/@sub?a=1&@b=@2",
					Response: testResponse(),
					Match: func(req *http.Request) bool {
						return tests.ExpectHeaders(req.Header, map[string]string{
							"Authorization":        "AWS4-HMAC-SHA256 Credential=456/20250102/test_region/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=e0001982deef1652704f74503203e77d83d4c88369421f9fca644d96f2a62a3c",
							"Host":                 "test_bucket.example.com",
							"X-Amz-Content-Sha256": "UNSIGNED-PAYLOAD",
							"X-Amz-Date":           "20250102T150405Z",
						})
					},
				}),
			},
		},
		{
			"with extra headers",
			"/test",
			func(req *http.Request) {
				req.Header.Set("x-amz-date", "20250102T150405Z")
				req.Header.Set("x-amz-content-sha256", "test_sha256")
				req.Header.Set("x-amz-example", "123")
				req.Header.Set("x-amz-meta-a", "456")
				req.Header.Set("content-type", "image/png")
				req.Header.Set("x-test", "789") // shouldn't be included in the signing headers
			},
			&s3.S3{
				Region:    "test_region",
				Bucket:    "test_bucket",
				Endpoint:  "https://example.com/",
				AccessKey: "123",
				SecretKey: "abc",
				Client: tests.NewClient(&tests.RequestStub{
					Method:   http.MethodGet,
					URL:      "https://test_bucket.example.com/test",
					Response: testResponse(),
					Match: func(req *http.Request) bool {
						return tests.ExpectHeaders(req.Header, map[string]string{
							"authorization":        "AWS4-HMAC-SHA256 Credential=123/20250102/test_region/s3/aws4_request, SignedHeaders=content-type;host;x-amz-content-sha256;x-amz-date;x-amz-example;x-amz-meta-a, Signature=86dccbcd012c33073dc99e9d0a9e0b717a4d8c11c37848cfa9a4a02716bc0db3",
							"host":                 "test_bucket.example.com",
							"x-amz-date":           "20250102T150405Z",
							"x-amz-content-sha256": "test_sha256",
							"x-amz-example":        "123",
							"x-amz-meta-a":         "456",
							"x-test":               "789",
						})
					},
				}),
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, s.s3Client.URL(s.path), strings.NewReader("test_request"))
			if err != nil {
				t.Fatal(err)
			}

			if s.reqFunc != nil {
				s.reqFunc(req)
			}

			resp, err := s.s3Client.SignAndSend(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			err = s.s3Client.Client.(*tests.Client).AssertNoRemaining()
			if err != nil {
				t.Fatal(err)
			}

			expectedBody := "test_response"

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if str := string(body); str != expectedBody {
				t.Fatalf("Expected body %q, got %q", expectedBody, str)
			}
		})
	}
}
