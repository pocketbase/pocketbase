package s3

import (
	"context"
	"io"
	"net/http"
)

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html#API_GetObject_ResponseElements
type GetObjectResponse struct {
	Body io.ReadCloser `json:"-" xml:"-"`

	HeadObjectResponse
}

// GetObject retrieves a single object by its key.
//
// NB! Make sure to call GetObjectResponse.Body.Close() after done working with the result.
//
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (s3 *S3) GetObject(ctx context.Context, key string, optFuncs ...func(*http.Request)) (*GetObjectResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s3.URL(key), nil)
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

	result := &GetObjectResponse{Body: resp.Body}
	result.load(resp.Header)

	return result, nil
}
