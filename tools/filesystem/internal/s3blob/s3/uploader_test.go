package s3_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3/tests"
)

func TestUploaderRequiredFields(t *testing.T) {
	t.Parallel()

	s3Client := &s3.S3{
		Client:    tests.NewClient(&tests.RequestStub{Method: "PUT", URL: `^.+$`}), // match every upload
		Region:    "test_region",
		Bucket:    "test_bucket",
		Endpoint:  "http://example.com",
		AccessKey: "123",
		SecretKey: "abc",
	}

	payload := strings.NewReader("test")

	scenarios := []struct {
		name          string
		uploader      *s3.Uploader
		expectedError bool
	}{
		{
			"blank",
			&s3.Uploader{},
			true,
		},
		{
			"no Key",
			&s3.Uploader{S3: s3Client, Payload: payload},
			true,
		},
		{
			"no S3",
			&s3.Uploader{Key: "abc", Payload: payload},
			true,
		},
		{
			"no Payload",
			&s3.Uploader{S3: s3Client, Key: "abc"},
			true,
		},
		{
			"with S3, Key and Payload",
			&s3.Uploader{S3: s3Client, Key: "abc", Payload: payload},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			err := s.uploader.Upload(context.Background())

			hasErr := err != nil
			if hasErr != s.expectedError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectedError, hasErr)
			}
		})
	}
}

func TestUploaderSingleUpload(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "http://test_bucket.example.com/test_key",
			Match: func(req *http.Request) bool {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}

				return string(body) == "abcdefg" && tests.ExpectHeaders(req.Header, map[string]string{
					"Content-Length": "7",
					"x-amz-meta-a":   "123",
					"x-amz-meta-b":   "456",
					"test_header":    "test",
					"Authorization":  "^.+Credential=123/.+$",
				})
			},
		},
	)

	uploader := &s3.Uploader{
		S3: &s3.S3{
			Client:    httpClient,
			Region:    "test_region",
			Bucket:    "test_bucket",
			Endpoint:  "http://example.com",
			AccessKey: "123",
			SecretKey: "abc",
		},
		Key:         "test_key",
		Payload:     strings.NewReader("abcdefg"),
		Metadata:    map[string]string{"a": "123", "b": "456"},
		MinPartSize: 8,
	}

	err := uploader.Upload(context.Background(), func(r *http.Request) {
		r.Header.Set("test_header", "test")
	})
	if err != nil {
		t.Fatal(err)
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploaderMultipartUploadSuccess(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method: http.MethodPost,
			URL:    "http://test_bucket.example.com/test_key?uploads",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"x-amz-meta-a":  "123",
					"x-amz-meta-b":  "456",
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Body: io.NopCloser(strings.NewReader(`
					<?xml version="1.0" encoding="UTF-8"?>
					<InitiateMultipartUploadResult>
					   <Bucket>test_bucket</Bucket>
					   <Key>test_key</Key>
					   <UploadId>test_id</UploadId>
					</InitiateMultipartUploadResult>
				`)),
			},
		},
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "http://test_bucket.example.com/test_key?partNumber=1&uploadId=test_id",
			Match: func(req *http.Request) bool {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}

				return string(body) == "abc" && tests.ExpectHeaders(req.Header, map[string]string{
					"Content-Length": "3",
					"test_header":    "test",
					"Authorization":  "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Header: http.Header{"Etag": []string{"etag1"}},
			},
		},
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "http://test_bucket.example.com/test_key?partNumber=2&uploadId=test_id",
			Match: func(req *http.Request) bool {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}

				return string(body) == "def" && tests.ExpectHeaders(req.Header, map[string]string{
					"Content-Length": "3",
					"test_header":    "test",
					"Authorization":  "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Header: http.Header{"Etag": []string{"etag2"}},
			},
		},
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "http://test_bucket.example.com/test_key?partNumber=3&uploadId=test_id",
			Match: func(req *http.Request) bool {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}
				return string(body) == "g" && tests.ExpectHeaders(req.Header, map[string]string{
					"Content-Length": "1",
					"test_header":    "test",
					"Authorization":  "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Header: http.Header{"Etag": []string{"etag3"}},
			},
		},
		&tests.RequestStub{
			Method: http.MethodPost,
			URL:    "http://test_bucket.example.com/test_key?uploadId=test_id",
			Match: func(req *http.Request) bool {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}

				expected := `<CompleteMultipartUpload><Part><ETag>etag1</ETag><PartNumber>1</PartNumber></Part><Part><ETag>etag2</ETag><PartNumber>2</PartNumber></Part><Part><ETag>etag3</ETag><PartNumber>3</PartNumber></Part></CompleteMultipartUpload>`

				return strings.Contains(string(body), expected) && tests.ExpectHeaders(req.Header, map[string]string{
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
		},
	)

	uploader := &s3.Uploader{
		S3: &s3.S3{
			Client:    httpClient,
			Region:    "test_region",
			Bucket:    "test_bucket",
			Endpoint:  "http://example.com",
			AccessKey: "123",
			SecretKey: "abc",
		},
		Key:         "test_key",
		Payload:     strings.NewReader("abcdefg"),
		Metadata:    map[string]string{"a": "123", "b": "456"},
		MinPartSize: 3,
	}

	err := uploader.Upload(context.Background(), func(r *http.Request) {
		r.Header.Set("test_header", "test")
	})
	if err != nil {
		t.Fatal(err)
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploaderMultipartUploadPartFailure(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method: http.MethodPost,
			URL:    "http://test_bucket.example.com/test_key?uploads",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"x-amz-meta-a":  "123",
					"x-amz-meta-b":  "456",
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Body: io.NopCloser(strings.NewReader(`
					<?xml version="1.0" encoding="UTF-8"?>
					<InitiateMultipartUploadResult>
					   <Bucket>test_bucket</Bucket>
					   <Key>test_key</Key>
					   <UploadId>test_id</UploadId>
					</InitiateMultipartUploadResult>
				`)),
			},
		},
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "http://test_bucket.example.com/test_key?partNumber=1&uploadId=test_id",
			Match: func(req *http.Request) bool {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}
				return string(body) == "abc" && tests.ExpectHeaders(req.Header, map[string]string{
					"Content-Length": "3",
					"test_header":    "test",
					"Authorization":  "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Header: http.Header{"Etag": []string{"etag1"}},
			},
		},
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "http://test_bucket.example.com/test_key?partNumber=2&uploadId=test_id",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				StatusCode: 400,
			},
		},
		&tests.RequestStub{
			Method: http.MethodDelete,
			URL:    "http://test_bucket.example.com/test_key?uploadId=test_id",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
		},
	)

	uploader := &s3.Uploader{
		S3: &s3.S3{
			Client:    httpClient,
			Region:    "test_region",
			Bucket:    "test_bucket",
			Endpoint:  "http://example.com",
			AccessKey: "123",
			SecretKey: "abc",
		},
		Key:         "test_key",
		Payload:     strings.NewReader("abcdefg"),
		Metadata:    map[string]string{"a": "123", "b": "456"},
		MinPartSize: 3,
	}

	err := uploader.Upload(context.Background(), func(r *http.Request) {
		r.Header.Set("test_header", "test")
	})
	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploaderMultipartUploadCompleteFailure(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method: http.MethodPost,
			URL:    "http://test_bucket.example.com/test_key?uploads",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"x-amz-meta-a":  "123",
					"x-amz-meta-b":  "456",
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Body: io.NopCloser(strings.NewReader(`
					<?xml version="1.0" encoding="UTF-8"?>
					<InitiateMultipartUploadResult>
					   <Bucket>test_bucket</Bucket>
					   <Key>test_key</Key>
					   <UploadId>test_id</UploadId>
					</InitiateMultipartUploadResult>
				`)),
			},
		},
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "http://test_bucket.example.com/test_key?partNumber=1&uploadId=test_id",
			Match: func(req *http.Request) bool {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}
				return string(body) == "abc" && tests.ExpectHeaders(req.Header, map[string]string{
					"Content-Length": "3",
					"test_header":    "test",
					"Authorization":  "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Header: http.Header{"Etag": []string{"etag1"}},
			},
		},
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "http://test_bucket.example.com/test_key?partNumber=2&uploadId=test_id",
			Match: func(req *http.Request) bool {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}
				return string(body) == "def" && tests.ExpectHeaders(req.Header, map[string]string{
					"Content-Length": "3",
					"test_header":    "test",
					"Authorization":  "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Header: http.Header{"Etag": []string{"etag2"}},
			},
		},
		&tests.RequestStub{
			Method: http.MethodPost,
			URL:    "http://test_bucket.example.com/test_key?uploadId=test_id",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				StatusCode: 400,
			},
		},
		&tests.RequestStub{
			Method: http.MethodDelete,
			URL:    "http://test_bucket.example.com/test_key?uploadId=test_id",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
		},
	)

	uploader := &s3.Uploader{
		S3: &s3.S3{
			Client:    httpClient,
			Region:    "test_region",
			Bucket:    "test_bucket",
			Endpoint:  "http://example.com",
			AccessKey: "123",
			SecretKey: "abc",
		},
		Key:         "test_key",
		Payload:     strings.NewReader("abcdef"),
		Metadata:    map[string]string{"a": "123", "b": "456"},
		MinPartSize: 3,
	}

	err := uploader.Upload(context.Background(), func(r *http.Request) {
		r.Header.Set("test_header", "test")
	})
	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}
}
