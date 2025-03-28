package s3blob_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/filesystem/blob"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3/tests"
)

func TestNew(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name        string
		s3Client    *s3.S3
		expectError bool
	}{
		{
			"blank",
			&s3.S3{},
			true,
		},
		{
			"no bucket",
			&s3.S3{Region: "b", Endpoint: "c"},
			true,
		},
		{
			"no endpoint",
			&s3.S3{Bucket: "a", Region: "b"},
			true,
		},
		{
			"no region",
			&s3.S3{Bucket: "a", Endpoint: "c"},
			true,
		},
		{
			"with bucket, endpoint and region",
			&s3.S3{Bucket: "a", Region: "b", Endpoint: "c"},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			drv, err := s3blob.New(s.s3Client)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if err == nil && drv == nil {
				t.Fatal("Expected non-nil driver instance")
			}
		})
	}
}

func TestDriverClose(t *testing.T) {
	t.Parallel()

	drv, err := s3blob.New(&s3.S3{Bucket: "a", Region: "b", Endpoint: "c"})
	if err != nil {
		t.Fatal(err)
	}

	err = drv.Close()
	if err != nil {
		t.Fatalf("Expected nil, got error %v", err)
	}
}

func TestDriverNormilizeError(t *testing.T) {
	t.Parallel()

	drv, err := s3blob.New(&s3.S3{Bucket: "a", Region: "b", Endpoint: "c"})
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name              string
		err               error
		expectErrNotFound bool
	}{
		{
			"plain error",
			errors.New("test"),
			false,
		},
		{
			"response error with only status (non-404)",
			&s3.ResponseError{Status: 123},
			false,
		},
		{
			"response error with only status (404)",
			&s3.ResponseError{Status: 404},
			true,
		},
		{
			"response error with custom code",
			&s3.ResponseError{Code: "test"},
			false,
		},
		{
			"response error with NoSuchBucket code",
			&s3.ResponseError{Code: "NoSuchBucket"},
			true,
		},
		{
			"response error with NoSuchKey code",
			&s3.ResponseError{Code: "NoSuchKey"},
			true,
		},
		{
			"response error with NotFound code",
			&s3.ResponseError{Code: "NotFound"},
			true,
		},
		{
			"wrapped response error with NotFound code", // ensures that the entire error's tree is checked
			fmt.Errorf("test: %w", &s3.ResponseError{Code: "NotFound"}),
			true,
		},
		{
			"already normalized error",
			fmt.Errorf("test: %w", blob.ErrNotFound),
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			err := drv.NormalizeError(s.err)
			if err == nil {
				t.Fatal("Expected non-nil error")
			}

			isErrNotFound := errors.Is(err, blob.ErrNotFound)
			if isErrNotFound != s.expectErrNotFound {
				t.Fatalf("Expected isErrNotFound %v, got %v (%v)", s.expectErrNotFound, isErrNotFound, err)
			}
		})
	}
}

func TestDriverDeleteEscaping(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(&tests.RequestStub{
		Method: http.MethodDelete,
		URL:    "https://test_bucket.example.com/..__0x2f__abc/test/",
	})

	drv, err := s3blob.New(&s3.S3{
		Bucket:   "test_bucket",
		Region:   "test_region",
		Endpoint: "https://example.com",
		Client:   httpClient,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = drv.Delete(context.Background(), "../abc/test/")
	if err != nil {
		t.Fatal(err)
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDriverCopyEscaping(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(&tests.RequestStub{
		Method: http.MethodPut,
		URL:    "https://test_bucket.example.com/..__0x2f__a/",
		Match: func(req *http.Request) bool {
			return tests.ExpectHeaders(req.Header, map[string]string{
				"x-amz-copy-source": "test_bucket%2F..__0x2f__b%2F",
			})
		},
		Response: &http.Response{
			Body: io.NopCloser(strings.NewReader(`<CopyObjectResult></CopyObjectResult>`)),
		},
	})

	drv, err := s3blob.New(&s3.S3{
		Bucket:   "test_bucket",
		Region:   "test_region",
		Endpoint: "https://example.com",
		Client:   httpClient,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = drv.Copy(context.Background(), "../a/", "../b/")
	if err != nil {
		t.Fatal(err)
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDriverAttributes(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(&tests.RequestStub{
		Method: http.MethodHead,
		URL:    "https://test_bucket.example.com/..__0x2f__a/",
		Response: &http.Response{
			Header: http.Header{
				"Last-Modified":       []string{"Mon, 01 Feb 2025 03:04:05 GMT"},
				"Cache-Control":       []string{"test_cache"},
				"Content-Disposition": []string{"test_disposition"},
				"Content-Encoding":    []string{"test_encoding"},
				"Content-Language":    []string{"test_language"},
				"Content-Type":        []string{"test_type"},
				"Content-Range":       []string{"test_range"},
				"Etag":                []string{`"ce5be8b6f53645c596306c4572ece521"`},
				"Content-Length":      []string{"100"},
				"x-amz-meta-AbC%40":   []string{"%40test_meta_a"},
				"x-amz-meta-Def":      []string{"test_meta_b"},
			},
			Body: http.NoBody,
		},
	})

	drv, err := s3blob.New(&s3.S3{
		Bucket:   "test_bucket",
		Region:   "test_region",
		Endpoint: "https://example.com",
		Client:   httpClient,
	})
	if err != nil {
		t.Fatal(err)
	}

	attrs, err := drv.Attributes(context.Background(), "../a/")
	if err != nil {
		t.Fatal(err)
	}

	raw, err := json.Marshal(attrs)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"cacheControl":"test_cache","contentDisposition":"test_disposition","contentEncoding":"test_encoding","contentLanguage":"test_language","contentType":"test_type","metadata":{"abc@":"@test_meta_a","def":"test_meta_b"},"createTime":"0001-01-01T00:00:00Z","modTime":"2025-02-01T03:04:05Z","size":100,"md5":"zlvotvU2RcWWMGxFcuzlIQ==","etag":"\"ce5be8b6f53645c596306c4572ece521\""}`
	if str := string(raw); str != expected {
		t.Fatalf("Expected attributes\n%s\ngot\n%s", expected, str)
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDriverListPaged(t *testing.T) {
	t.Parallel()

	listResponse := func() *http.Response {
		return &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				<?xml version="1.0" encoding="UTF-8"?>
				<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
					<Name>example</Name>
					<ContinuationToken>ct</ContinuationToken>
					<NextContinuationToken>test_next</NextContinuationToken>
					<StartAfter>example0.txt</StartAfter>
					<KeyCount>1</KeyCount>
					<MaxKeys>3</MaxKeys>
					<Contents>
						<Key>..__0x2f__prefixB/test/example.txt</Key>
						<LastModified>2025-01-01T01:02:03.123Z</LastModified>
						<ETag>"ce5be8b6f53645c596306c4572ece521"</ETag>
						<Size>123</Size>
					</Contents>
					<Contents>
						<Key>prefixA/..__0x2f__escape.txt</Key>
						<LastModified>2025-01-02T01:02:03.123Z</LastModified>
						<Size>456</Size>
					</Contents>
					<CommonPrefixes>
						<Prefix>prefixA</Prefix>
					</CommonPrefixes>
					<CommonPrefixes>
						<Prefix>..__0x2f__prefixB</Prefix>
					</CommonPrefixes>
				</ListBucketResult>
			`)),
		}
	}

	expectedPage := `{"objects":[{"key":"../prefixB","modTime":"0001-01-01T00:00:00Z","size":0,"md5":null,"isDir":true},{"key":"../prefixB/test/example.txt","modTime":"2025-01-01T01:02:03.123Z","size":123,"md5":"zlvotvU2RcWWMGxFcuzlIQ==","isDir":false},{"key":"prefixA","modTime":"0001-01-01T00:00:00Z","size":0,"md5":null,"isDir":true},{"key":"prefixA/../escape.txt","modTime":"2025-01-02T01:02:03.123Z","size":456,"md5":null,"isDir":false}],"nextPageToken":"dGVzdF9uZXh0"}`

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method:   http.MethodGet,
			URL:      "https://test_bucket.example.com/?list-type=2&max-keys=1000",
			Response: listResponse(),
		},
		&tests.RequestStub{
			Method:   http.MethodGet,
			URL:      "https://test_bucket.example.com/?continuation-token=test_token&delimiter=test_delimiter&list-type=2&max-keys=123&prefix=test_prefix",
			Response: listResponse(),
		},
	)

	drv, err := s3blob.New(&s3.S3{
		Bucket:   "test_bucket",
		Region:   "test_region",
		Endpoint: "https://example.com",
		Client:   httpClient,
	})
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name     string
		opts     *blob.ListOptions
		expected string
	}{
		{
			"empty options",
			&blob.ListOptions{},
			expectedPage,
		},
		{
			"filled options",
			&blob.ListOptions{Prefix: "test_prefix", Delimiter: "test_delimiter", PageSize: 123, PageToken: []byte("test_token")},
			expectedPage,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			page, err := drv.ListPaged(context.Background(), s.opts)
			if err != nil {
				t.Fatal(err)
			}

			raw, err := json.Marshal(page)
			if err != nil {
				t.Fatal(err)
			}

			if str := string(raw); s.expected != str {
				t.Fatalf("Expected page result\n%s\ngot\n%s", s.expected, str)
			}
		})
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDriverNewRangeReader(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		offset        int64
		length        int64
		httpClient    *tests.Client
		expectedAttrs string
	}{
		{
			0,
			0,
			tests.NewClient(&tests.RequestStub{
				Method: http.MethodGet,
				URL:    "https://test_bucket.example.com/..__0x2f__abc/test.txt",
				Match: func(req *http.Request) bool {
					return tests.ExpectHeaders(req.Header, map[string]string{
						"Range": "bytes=0-0",
					})
				},
				Response: &http.Response{
					Header: http.Header{
						"Last-Modified":  []string{"Mon, 01 Feb 2025 03:04:05 GMT"},
						"Content-Type":   []string{"test_ct"},
						"Content-Length": []string{"123"},
					},
					Body: io.NopCloser(strings.NewReader("test")),
				},
			}),
			`{"contentType":"test_ct","modTime":"2025-02-01T03:04:05Z","size":123}`,
		},
		{
			10,
			-1,
			tests.NewClient(&tests.RequestStub{
				Method: http.MethodGet,
				URL:    "https://test_bucket.example.com/..__0x2f__abc/test.txt",
				Match: func(req *http.Request) bool {
					return tests.ExpectHeaders(req.Header, map[string]string{
						"Range": "bytes=10-",
					})
				},
				Response: &http.Response{
					Header: http.Header{
						"Last-Modified":  []string{"Mon, 01 Feb 2025 03:04:05 GMT"},
						"Content-Type":   []string{"test_ct"},
						"Content-Range":  []string{"bytes 1-1/456"}, // should take precedence over content-length
						"Content-Length": []string{"123"},
					},
					Body: io.NopCloser(strings.NewReader("test")),
				},
			}),
			`{"contentType":"test_ct","modTime":"2025-02-01T03:04:05Z","size":456}`,
		},
		{
			10,
			0,
			tests.NewClient(&tests.RequestStub{
				Method: http.MethodGet,
				URL:    "https://test_bucket.example.com/..__0x2f__abc/test.txt",
				Match: func(req *http.Request) bool {
					return tests.ExpectHeaders(req.Header, map[string]string{
						"Range": "bytes=10-10",
					})
				},
				Response: &http.Response{
					Header: http.Header{
						"Last-Modified": []string{"Mon, 01 Feb 2025 03:04:05 GMT"},
						"Content-Type":  []string{"test_ct"},
						// no range and length headers
						// "Content-Range":  []string{"bytes 1-1/456"},
						// "Content-Length": []string{"123"},
					},
					Body: io.NopCloser(strings.NewReader("test")),
				},
			}),
			`{"contentType":"test_ct","modTime":"2025-02-01T03:04:05Z","size":0}`,
		},
		{
			10,
			20,
			tests.NewClient(&tests.RequestStub{
				Method: http.MethodGet,
				URL:    "https://test_bucket.example.com/..__0x2f__abc/test.txt",
				Match: func(req *http.Request) bool {
					return tests.ExpectHeaders(req.Header, map[string]string{
						"Range": "bytes=10-29",
					})
				},
				Response: &http.Response{
					Header: http.Header{
						"Last-Modified": []string{"Mon, 01 Feb 2025 03:04:05 GMT"},
						"Content-Type":  []string{"test_ct"},
						// with range header but invalid format -> content-length takes precedence
						"Content-Range":  []string{"bytes invalid-456"},
						"Content-Length": []string{"123"},
					},
					Body: io.NopCloser(strings.NewReader("test")),
				},
			}),
			`{"contentType":"test_ct","modTime":"2025-02-01T03:04:05Z","size":123}`,
		},
	}

	for _, s := range scenarios {
		t.Run(fmt.Sprintf("offset_%d_length_%d", s.offset, s.length), func(t *testing.T) {
			drv, err := s3blob.New(&s3.S3{
				Bucket:   "test_bucket",
				Region:   "tesst_region",
				Endpoint: "https://example.com",
				Client:   s.httpClient,
			})
			if err != nil {
				t.Fatal(err)
			}

			r, err := drv.NewRangeReader(context.Background(), "../abc/test.txt", s.offset, s.length)
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			// the response body should be always replaced with http.NoBody
			if s.length == 0 {
				body := make([]byte, 1)
				n, err := r.Read(body)
				if n != 0 || !errors.Is(err, io.EOF) {
					t.Fatalf("Expected body to be http.NoBody, got %v (%v)", body, err)
				}
			}

			rawAttrs, err := json.Marshal(r.Attributes())
			if err != nil {
				t.Fatal(err)
			}

			if str := string(rawAttrs); str != s.expectedAttrs {
				t.Fatalf("Expected attributes\n%s\ngot\n%s", s.expectedAttrs, str)
			}

			err = s.httpClient.AssertNoRemaining()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDriverNewTypedWriter(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "https://test_bucket.example.com/..__0x2f__abc/test/",
			Match: func(req *http.Request) bool {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}

				return string(body) == "test" && tests.ExpectHeaders(req.Header, map[string]string{
					"cache-control":       "test_cache_control",
					"content-disposition": "test_content_disposition",
					"content-encoding":    "test_content_encoding",
					"content-language":    "test_content_language",
					"content-type":        "test_ct",
					"content-md5":         "dGVzdA==",
				})
			},
		},
	)

	drv, err := s3blob.New(&s3.S3{
		Bucket:   "test_bucket",
		Region:   "test_region",
		Endpoint: "https://example.com",
		Client:   httpClient,
	})
	if err != nil {
		t.Fatal(err)
	}

	options := &blob.WriterOptions{
		CacheControl:       "test_cache_control",
		ContentDisposition: "test_content_disposition",
		ContentEncoding:    "test_content_encoding",
		ContentLanguage:    "test_content_language",
		ContentType:        "test_content_type", // should be ignored
		ContentMD5:         []byte("test"),
		Metadata:           map[string]string{"@test_meta_a": "@test"},
	}

	w, err := drv.NewTypedWriter(context.Background(), "../abc/test/", "test_ct", options)
	if err != nil {
		t.Fatal(err)
	}

	n, err := w.Write(nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("Expected nil write to result in %d written bytes, got %d", 0, n)
	}

	n, err = w.Write([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 4 {
		t.Fatalf("Expected nil write to result in %d written bytes, got %d", 4, n)
	}

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}
}
