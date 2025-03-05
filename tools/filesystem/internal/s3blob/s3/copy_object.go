package s3

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_CopyObject.html#API_CopyObject_ResponseSyntax
type CopyObjectResponse struct {
	CopyObjectResult  xml.Name  `json:"copyObjectResult" xml:"CopyObjectResult"`
	ETag              string    `json:"etag" xml:"ETag"`
	LastModified      time.Time `json:"lastModified" xml:"LastModified"`
	ChecksumType      string    `json:"checksumType" xml:"ChecksumType"`
	ChecksumCRC32     string    `json:"checksumCRC32" xml:"ChecksumCRC32"`
	ChecksumCRC32C    string    `json:"checksumCRC32C" xml:"ChecksumCRC32C"`
	ChecksumCRC64NVME string    `json:"checksumCRC64NVME" xml:"ChecksumCRC64NVME"`
	ChecksumSHA1      string    `json:"checksumSHA1" xml:"ChecksumSHA1"`
	ChecksumSHA256    string    `json:"checksumSHA256" xml:"ChecksumSHA256"`
}

// CopyObject copies a single object from srcKey to dstKey destination.
// (both keys are expected to be operating within the same bucket).
//
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_CopyObject.html
func (s3 *S3) CopyObject(ctx context.Context, srcKey string, dstKey string, optReqFuncs ...func(*http.Request)) (*CopyObjectResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s3.URL(dstKey), nil)
	if err != nil {
		return nil, err
	}

	// per the doc the header value must be URL-encoded
	req.Header.Set("x-amz-copy-source", url.PathEscape(s3.Bucket+"/"+strings.TrimLeft(srcKey, "/")))

	// apply optional request funcs
	for _, fn := range optReqFuncs {
		if fn != nil {
			fn(req)
		}
	}

	resp, err := s3.SignAndSend(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &CopyObjectResponse{}

	err = xml.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
