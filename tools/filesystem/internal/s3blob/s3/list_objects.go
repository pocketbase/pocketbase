package s3

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ListParams defines optional parameters for the ListObject request.
type ListParams struct {
	// ContinuationToken indicates that the list is being continued on this bucket with a token.
	// ContinuationToken is obfuscated and is not a real key.
	// You can use this ContinuationToken for pagination of the list results.
	ContinuationToken string `json:"continuationToken"`

	// Delimiter is a character that you use to group keys.
	//
	// For directory buckets, "/" is the only supported delimiter.
	Delimiter string `json:"delimiter"`

	// Prefix limits the response to keys that begin with the specified prefix.
	Prefix string `json:"prefix"`

	// Encoding type is used to encode the object keys in the response.
	// Responses are encoded only in UTF-8.
	// An object key can contain any Unicode character.
	// However, the XML 1.0 parser can't parse certain characters,
	// such as characters with an ASCII value from 0 to 10.
	// For characters that aren't supported in XML 1.0, you can add
	// this parameter to request that S3 encode the keys in the response.
	//
	// Valid Values: url
	EncodingType string `json:"encodingType"`

	// StartAfter is where you want S3 to start listing from.
	// S3 starts listing after this specified key.
	// StartAfter can be any key in the bucket.
	//
	// This functionality is not supported for directory buckets.
	StartAfter string `json:"startAfter"`

	// MaxKeys Sets the maximum number of keys returned in the response.
	// By default, the action returns up to 1,000 key names.
	// The response might contain fewer keys but will never contain more.
	MaxKeys int `json:"maxKeys"`

	// FetchOwner returns the owner field with each key in the result.
	FetchOwner bool `json:"fetchOwner"`
}

// Encode encodes the parameters in a properly formatted query string.
func (l *ListParams) Encode() string {
	query := url.Values{}

	query.Add("list-type", "2")

	if l.ContinuationToken != "" {
		query.Add("continuation-token", l.ContinuationToken)
	}

	if l.Delimiter != "" {
		query.Add("delimiter", l.Delimiter)
	}

	if l.Prefix != "" {
		query.Add("prefix", l.Prefix)
	}

	if l.EncodingType != "" {
		query.Add("encoding-type", l.EncodingType)
	}

	if l.FetchOwner {
		query.Add("fetch-owner", "true")
	}

	if l.MaxKeys > 0 {
		query.Add("max-keys", strconv.Itoa(l.MaxKeys))
	}

	if l.StartAfter != "" {
		query.Add("start-after", l.StartAfter)
	}

	return query.Encode()
}

// ListObjects retrieves paginated objects list.
//
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html
func (s3 *S3) ListObjects(ctx context.Context, params ListParams, optReqFuncs ...func(*http.Request)) (*ListObjectsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s3.URL("?"+params.Encode()), nil)
	if err != nil {
		return nil, err
	}

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

	result := &ListObjectsResponse{}

	err = xml.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html#API_ListObjectsV2_ResponseSyntax
type ListObjectsResponse struct {
	XMLName               xml.Name `json:"-" xml:"ListBucketResult"`
	EncodingType          string   `json:"encodingType" xml:"EncodingType"`
	Name                  string   `json:"name" xml:"Name"`
	Prefix                string   `json:"prefix" xml:"Prefix"`
	Delimiter             string   `json:"delimiter" xml:"Delimiter"`
	ContinuationToken     string   `json:"continuationToken" xml:"ContinuationToken"`
	NextContinuationToken string   `json:"nextContinuationToken" xml:"NextContinuationToken"`
	StartAfter            string   `json:"startAfter" xml:"StartAfter"`

	CommonPrefixes []*ListObjectCommonPrefix `json:"commonPrefixes" xml:"CommonPrefixes"`

	Contents []*ListObjectContent `json:"contents" xml:"Contents"`

	KeyCount    int  `json:"keyCount" xml:"KeyCount"`
	MaxKeys     int  `json:"maxKeys" xml:"MaxKeys"`
	IsTruncated bool `json:"isTruncated" xml:"IsTruncated"`
}

type ListObjectCommonPrefix struct {
	Prefix string `json:"prefix" xml:"Prefix"`
}

type ListObjectContent struct {
	Owner struct {
		DisplayName string `json:"displayName" xml:"DisplayName"`
		ID          string `json:"id" xml:"ID"`
	} `json:"owner" xml:"Owner"`

	ChecksumAlgorithm string    `json:"checksumAlgorithm" xml:"ChecksumAlgorithm"`
	ETag              string    `json:"etag" xml:"ETag"`
	Key               string    `json:"key" xml:"Key"`
	StorageClass      string    `json:"storageClass" xml:"StorageClass"`
	LastModified      time.Time `json:"lastModified" xml:"LastModified"`

	RestoreStatus struct {
		RestoreExpiryDate   time.Time `json:"restoreExpiryDate" xml:"RestoreExpiryDate"`
		IsRestoreInProgress bool      `json:"isRestoreInProgress" xml:"IsRestoreInProgress"`
	} `json:"restoreStatus" xml:"RestoreStatus"`

	Size int64 `json:"size" xml:"Size"`
}
