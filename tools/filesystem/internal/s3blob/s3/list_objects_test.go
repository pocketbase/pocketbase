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

func TestS3ListParamsEncode(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name     string
		params   s3.ListParams
		expected string
	}{
		{
			"blank",
			s3.ListParams{},
			"list-type=2",
		},
		{
			"filled",
			s3.ListParams{
				ContinuationToken: "test_ct",
				Delimiter:         "test_delimiter",
				Prefix:            "test_prefix",
				EncodingType:      "test_et",
				StartAfter:        "test_sa",
				MaxKeys:           1,
				FetchOwner:        true,
			},
			"continuation-token=test_ct&delimiter=test_delimiter&encoding-type=test_et&fetch-owner=true&list-type=2&max-keys=1&prefix=test_prefix&start-after=test_sa",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.params.Encode()
			if result != s.expected {
				t.Fatalf("Expected\n%s\ngot\n%s", s.expected, result)
			}
		})
	}
}

func TestS3ListObjects(t *testing.T) {
	t.Parallel()

	listParams := s3.ListParams{
		ContinuationToken: "test_ct",
		Delimiter:         "test_delimiter",
		Prefix:            "test_prefix",
		EncodingType:      "test_et",
		StartAfter:        "test_sa",
		MaxKeys:           10,
		FetchOwner:        true,
	}

	httpClient := tests.NewClient(
		&tests.RequestStub{
			Method: http.MethodGet,
			URL:    "http://test_bucket.example.com/?" + listParams.Encode(),
			Match: func(req *http.Request) bool {
				return tests.ExpectHeaders(req.Header, map[string]string{
					"test_header":   "test",
					"Authorization": "^.+Credential=123/.+$",
				})
			},
			Response: &http.Response{
				Body: io.NopCloser(strings.NewReader(`
					<?xml version="1.0" encoding="UTF-8"?>
					<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
						<Name>example</Name>
						<EncodingType>test_encoding</EncodingType>
						<Prefix>a/</Prefix>
						<Delimiter>/</Delimiter>
						<ContinuationToken>ct</ContinuationToken>
						<NextContinuationToken>nct</NextContinuationToken>
						<StartAfter>example0.txt</StartAfter>
						<KeyCount>1</KeyCount>
						<MaxKeys>3</MaxKeys>
						<IsTruncated>true</IsTruncated>
						<Contents>
							<Key>example1.txt</Key>
							<LastModified>2025-01-01T01:02:03.123Z</LastModified>
							<ChecksumAlgorithm>test_ca</ChecksumAlgorithm>
							<ETag>test_etag1</ETag>
							<Size>123</Size>
							<StorageClass>STANDARD</StorageClass>
							<Owner>
								<DisplayName>owner_dn</DisplayName>
								<ID>owner_id</ID>
							</Owner>
							<RestoreStatus>
								<RestoreExpiryDate>2025-01-02T01:02:03.123Z</RestoreExpiryDate>
								<IsRestoreInProgress>true</IsRestoreInProgress>
							</RestoreStatus>
						</Contents>
						<Contents>
							<Key>example2.txt</Key>
							<LastModified>2025-01-02T01:02:03.123Z</LastModified>
							<ETag>test_etag2</ETag>
							<Size>456</Size>
							<StorageClass>STANDARD</StorageClass>
						</Contents>
						<CommonPrefixes>
							<Prefix>a/b/</Prefix>
						</CommonPrefixes>
						<CommonPrefixes>
							<Prefix>a/c/</Prefix>
						</CommonPrefixes>
					</ListBucketResult>
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

	resp, err := s3Client.ListObjects(context.Background(), listParams, func(r *http.Request) {
		r.Header.Set("test_header", "test")
	})
	if err != nil {
		t.Fatal(err)
	}

	err = httpClient.AssertNoRemaining()
	if err != nil {
		t.Fatal(err)
	}

	raw, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)

	expected := `{"encodingType":"test_encoding","name":"example","prefix":"a/","delimiter":"/","continuationToken":"ct","nextContinuationToken":"nct","startAfter":"example0.txt","commonPrefixes":[{"prefix":"a/b/"},{"prefix":"a/c/"}],"contents":[{"owner":{"displayName":"owner_dn","id":"owner_id"},"checksumAlgorithm":"test_ca","etag":"test_etag1","key":"example1.txt","storageClass":"STANDARD","lastModified":"2025-01-01T01:02:03.123Z","restoreStatus":{"restoreExpiryDate":"2025-01-02T01:02:03.123Z","isRestoreInProgress":true},"size":123},{"owner":{"displayName":"","id":""},"checksumAlgorithm":"","etag":"test_etag2","key":"example2.txt","storageClass":"STANDARD","lastModified":"2025-01-02T01:02:03.123Z","restoreStatus":{"restoreExpiryDate":"0001-01-01T00:00:00Z","isRestoreInProgress":false},"size":456}],"keyCount":1,"maxKeys":3,"isTruncated":true}`

	if rawStr != expected {
		t.Fatalf("Expected response\n%s\ngot\n%s", expected, rawStr)
	}
}
