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

func TestS3CopyObject(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method: http.MethodPut,
			URL:    "http://test_bucket.example.com/@dst_test",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"test_header":       "test",
					"x-amz-copy-source": "test_bucket%2F@src_test",
					"Authorization":     "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Body: io.NopCloser(strings.NewReader(`
					<CopyObjectResult>
						<LastModified>2025-01-01T01:02:03.456Z</LastModified>
						<ETag>test_etag</ETag>
					</CopyObjectResult>
				`)),
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

	copyResp, err := s3Client.CopyObject(context.Background(), "@src_test", "@dst_test", func(r *http.Request) {
		r.Header.Set("test_header", "test")
	})
	if err != nil {
		t.Fatal(err)
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}

	if copyResp.ETag != "test_etag" {
		t.Fatalf("Expected ETag %q, got %q", "test_etag", copyResp.ETag)
	}

	if date := copyResp.LastModified.Format("2006-01-02T15:04:05.000Z"); date != "2025-01-01T01:02:03.456Z" {
		t.Fatalf("Expected LastModified %q, got %q", "2025-01-01T01:02:03.456Z", date)
	}
}
