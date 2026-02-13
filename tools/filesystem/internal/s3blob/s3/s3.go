// Package s3 implements a lightweight client for interacting with the
// REST APIs of any S3 compatible service.
//
// It implements only the minimal functionality required by PocketBase
// such as objects list, get, copy, delete and upload.
//
// For more details why we don't use the official aws-sdk-go-v2, you could check
// https://github.com/pocketbase/pocketbase/discussions/6562.
//
// Example:
//
//	client := &s3.S3{
//		Endpoint:     "example.com",
//		Region:       "us-east-1",
//		Bucket:       "test",
//		AccessKey:    "...",
//		SecretKey:    "...",
//		UsePathStyle: true,
//	}
//	resp, err := client.GetObject(context.Background(), "abc.txt")
package s3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

const (
	awsS3ServiceCode     = "s3"
	awsSignAlgorithm     = "AWS4-HMAC-SHA256"
	awsTerminationString = "aws4_request"
	metadataPrefix       = "x-amz-meta-"
	dateTimeFormat       = "20060102T150405Z"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type S3 struct {
	// Client specifies a custom HTTP client to send the request with.
	//
	// If not explicitly set, fallbacks to http.DefaultClient.
	Client HTTPClient

	Bucket       string
	Region       string
	Endpoint     string // can be with or without the schema
	AccessKey    string
	SecretKey    string
	UsePathStyle bool
}

// URL constructs an S3 request URL based on the current configuration.
//
// Note that the path will be URL escaped based on the AWS [UriEncode rules]
// for broader compatibility with some providers that expect the same
// path format as the one in the canonical signed header
// (see also https://github.com/pocketbase/pocketbase/issues/7153).
//
// [UriEncode rules]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html
func (s3 *S3) URL(path string) string {
	scheme := "https"
	endpoint := strings.TrimRight(s3.Endpoint, "/")
	if after, ok := strings.CutPrefix(endpoint, "https://"); ok {
		endpoint = after
	} else if after, ok := strings.CutPrefix(endpoint, "http://"); ok {
		endpoint = after
		scheme = "http"
	}

	// to prevent double escaping we first parse/unescape it
	parsed, err := url.Parse(path)
	if err != nil {
		// truly rare case, keep the path as it is
	} else {
		path = escapePath(parsed.Path)

		// the rest is usually not expected to be part of the S3 path but it is kept to avoid surprises
		// (it will be further escaped if necessery by the Go HTTP client)
		if parsed.RawQuery != "" {
			path += "?" + parsed.RawQuery
		}
		if parsed.RawFragment != "" {
			path += "#" + parsed.RawFragment
		}
	}

	path = strings.TrimLeft(path, "/")

	if s3.UsePathStyle {
		return fmt.Sprintf("%s://%s/%s/%s", scheme, endpoint, s3.Bucket, path)
	}

	return fmt.Sprintf("%s://%s.%s/%s", scheme, s3.Bucket, endpoint, path)
}

// SignAndSend signs the provided request per AWS Signature v4 and sends it.
//
// It automatically normalizes all 40x/50x responses to ResponseError.
//
// Note: Don't forget to call resp.Body.Close() after done with the result.
func (s3 *S3) SignAndSend(req *http.Request) (*http.Response, error) {
	s3.sign(req)

	client := s3.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()

		respErr := &ResponseError{
			Status: resp.StatusCode,
		}

		respErr.Raw, err = io.ReadAll(resp.Body)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, errors.Join(err, respErr)
		}

		if len(respErr.Raw) > 0 {
			err = xml.Unmarshal(respErr.Raw, respErr)
			if err != nil {
				return nil, errors.Join(err, respErr)
			}
		}

		return nil, respErr
	}

	return resp, nil
}

// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html#create-signed-request-steps
func (s3 *S3) sign(req *http.Request) {
	// explicitly set Accept-Encoding to avoid transparent decompression
	// and Content-Length zeroing (https://github.com/pocketbase/pocketbase/issues/7523)
	if req.Header.Get("Accept-Encoding") == "" {
		req.Header.Set("Accept-Encoding", "identity")
	}

	// fallback to the Unsigned payload option
	// (data integrity checks could be still applied via the content-md5 or x-amz-checksum-* headers)
	if req.Header.Get("x-amz-content-sha256") == "" {
		req.Header.Set("x-amz-content-sha256", "UNSIGNED-PAYLOAD")
	}

	reqDateTime, _ := time.Parse(dateTimeFormat, req.Header.Get("x-amz-date"))
	if reqDateTime.IsZero() {
		reqDateTime = time.Now().UTC()
		req.Header.Set("x-amz-date", reqDateTime.Format(dateTimeFormat))
	}

	req.Header.Set("host", req.URL.Host)

	date := reqDateTime.Format("20060102")

	dateTime := reqDateTime.Format(dateTimeFormat)

	// 1. Create canonical request
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html#create-canonical-request
	// ---------------------------------------------------------------
	canonicalHeaders, signedHeaders := canonicalAndSignedHeaders(req)

	canonicalParts := []string{
		req.Method,
		escapePath(req.URL.Path),
		escapeQuery(req.URL.Query()),
		canonicalHeaders,
		signedHeaders,
		req.Header.Get("x-amz-content-sha256"),
	}

	// 2. Create a hash of the canonical request
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html#create-canonical-request-hash
	// ---------------------------------------------------------------
	hashedCanonicalRequest := sha256Hex([]byte(strings.Join(canonicalParts, "\n")))

	// 3. Create a string to sign
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html#create-string-to-sign
	// ---------------------------------------------------------------
	scope := strings.Join([]string{
		date,
		s3.Region,
		awsS3ServiceCode,
		awsTerminationString,
	}, "/")

	stringToSign := strings.Join([]string{
		awsSignAlgorithm,
		dateTime,
		scope,
		hashedCanonicalRequest,
	}, "\n")

	// 4. Derive a signing key for SigV4
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html#derive-signing-key
	// ---------------------------------------------------------------
	dateKey := hmacSHA256([]byte("AWS4"+s3.SecretKey), date)
	dateRegionKey := hmacSHA256(dateKey, s3.Region)
	dateRegionServiceKey := hmacSHA256(dateRegionKey, awsS3ServiceCode)
	signingKey := hmacSHA256(dateRegionServiceKey, awsTerminationString)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	// 5. Add the signature to the request
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html#add-signature-to-request
	authorization := fmt.Sprintf(
		"%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		awsSignAlgorithm,
		s3.AccessKey,
		scope,
		signedHeaders,
		signature,
	)

	req.Header.Set("authorization", authorization)
}

func sha256Hex(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}

func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func canonicalAndSignedHeaders(req *http.Request) (string, string) {
	signed := []string{}
	canonical := map[string]string{}

	for key, values := range req.Header {
		normalizedKey := strings.ToLower(key)

		if normalizedKey != "host" &&
			normalizedKey != "content-type" &&
			!strings.HasPrefix(normalizedKey, "x-amz-") {
			continue
		}

		signed = append(signed, normalizedKey)

		// for each value:
		// trim any leading or trailing spaces
		// convert sequential spaces to a single space
		normalizedValues := make([]string, len(values))
		for i, v := range values {
			normalizedValues[i] = strings.ReplaceAll(strings.TrimSpace(v), "  ", " ")
		}

		canonical[normalizedKey] = strings.Join(normalizedValues, ",")
	}

	slices.Sort(signed)

	var sortedCanonical strings.Builder
	for _, key := range signed {
		sortedCanonical.WriteString(key)
		sortedCanonical.WriteString(":")
		sortedCanonical.WriteString(canonical[key])
		sortedCanonical.WriteString("\n")
	}

	return sortedCanonical.String(), strings.Join(signed, ";")
}

// extractMetadata parses and extracts and the metadata from the specified request headers.
//
// The metadata keys are all lowercased and without the "x-amz-meta-" prefix.
func extractMetadata(headers http.Header) map[string]string {
	result := map[string]string{}

	for k, v := range headers {
		if len(v) == 0 {
			continue
		}

		metadataKey, ok := strings.CutPrefix(strings.ToLower(k), metadataPrefix)
		if !ok {
			continue
		}

		result[metadataKey] = v[0]
	}

	return result
}

// escapeQuery returns the URI encoded request query parameters according to the AWS S3 spec requirements
// (it is similar to url.Values.Encode but instead of url.QueryEscape uses our own escape method).
func escapeQuery(values url.Values) string {
	if len(values) == 0 {
		return ""
	}

	var buf strings.Builder

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		vs := values[k]
		keyEscaped := escape(k)
		for _, values := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(escape(values))
		}
	}

	return buf.String()
}

// escapePath returns the URI encoded request path according to the AWS S3 spec requirements.
func escapePath(path string) string {
	parts := strings.Split(path, "/")

	for i, part := range parts {
		parts[i] = escape(part)
	}

	return strings.Join(parts, "/")
}

const upperhex = "0123456789ABCDEF"

// escape is similar to the std url.escape but implements the AWS [UriEncode requirements]:
//   - URI encode every byte except the unreserved characters: 'A'-'Z', 'a'-'z', '0'-'9', '-', '.', '_', and '~'.
//   - The space character is a reserved character and must be encoded as "%20" (and not as "+").
//   - Each URI encoded byte is formed by a '%' and the two-digit hexadecimal value of the byte.
//   - Letters in the hexadecimal value must be uppercase, for example "%1A".
//
// [UriEncode requirements]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html
func escape(s string) string {
	hexCount := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	result := make([]byte, len(s)+2*hexCount)

	j := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			result[j] = '%'
			result[j+1] = upperhex[c>>4]
			result[j+2] = upperhex[c&15]
			j += 3
		} else {
			result[j] = c
			j++
		}
	}

	return string(result)
}

// > "URI encode every byte except the unreserved characters: 'A'-'Z', 'a'-'z', '0'-'9', '-', '.', '_', and '~'."
func shouldEscape(c byte) bool {
	isUnreserved := (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '.' || c == '_' || c == '~'

	return !isUnreserved
}
