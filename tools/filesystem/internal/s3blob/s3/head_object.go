package s3

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadObject.html#API_HeadObject_ResponseElements
type HeadObjectResponse struct {
	// Metadata is the extra data that is stored with the S3 object (aka. the "x-amz-meta-*" header values).
	//
	// The map keys are normalized to lower-case.
	Metadata map[string]string `json:"metadata"`

	// LastModified date and time when the object was last modified.
	LastModified time.Time `json:"lastModified"`

	// CacheControl specifies caching behavior along the request/reply chain.
	CacheControl string `json:"cacheControl"`

	// ContentDisposition specifies presentational information for the object.
	ContentDisposition string `json:"contentDisposition"`

	// ContentEncoding indicates what content encodings have been applied to the object
	// and thus what decoding mechanisms must be applied to obtain the
	// media-type referenced by the Content-Type header field.
	ContentEncoding string `json:"contentEncoding"`

	// ContentLanguage indicates the language the content is in.
	ContentLanguage string `json:"contentLanguage"`

	// ContentType is a standard MIME type describing the format of the object data.
	ContentType string `json:"contentType"`

	// ContentRange is the portion of the object usually returned in the response for a GET request.
	ContentRange string `json:"contentRange"`

	// ETag is an opaque identifier assigned by a web
	// server to a specific version of a resource found at a URL.
	ETag string `json:"etag"`

	// ContentLength is size of the body in bytes.
	ContentLength int64 `json:"contentLength"`
}

// load parses and load the header values into the current HeadObjectResponse fields.
func (o *HeadObjectResponse) load(headers http.Header) {
	o.LastModified, _ = time.Parse(time.RFC1123, headers.Get("Last-Modified"))
	o.CacheControl = headers.Get("Cache-Control")
	o.ContentDisposition = headers.Get("Content-Disposition")
	o.ContentEncoding = headers.Get("Content-Encoding")
	o.ContentLanguage = headers.Get("Content-Language")
	o.ContentType = headers.Get("Content-Type")
	o.ContentRange = headers.Get("Content-Range")
	o.ETag = headers.Get("ETag")
	o.ContentLength, _ = strconv.ParseInt(headers.Get("Content-Length"), 10, 0)
	o.Metadata = extractMetadata(headers)
}

// HeadObject sends a HEAD request for a single object to check its
// existence and to retrieve its metadata.
//
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadObject.html
func (s3 *S3) HeadObject(ctx context.Context, key string, optFuncs ...func(*http.Request)) (*HeadObjectResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, s3.URL(key), nil)
	if err != nil {
		return nil, err
	}

	// apply optional request funcs
	for _, fn := range optFuncs {
		if fn != nil {
			fn(req)
		}
	}

	resp, err := s3.SignAndSend(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &HeadObjectResponse{}
	result.load(resp.Header)

	return result, nil
}
