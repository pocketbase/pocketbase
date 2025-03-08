package s3_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3/tests"
)

func TestS3DeleteObject(t *testing.T) {
	t.Parallel()

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method: http.MethodDelete,
			URL:    "http://test_bucket.example.com/test_key",
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
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

	err := s3Client.DeleteObject(context.Background(), "test_key", func(r *http.Request) {
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
