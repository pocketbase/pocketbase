// This is a trimmed version of the original go-cloud/s3blob driver
// to avoid loading both aws-sdk-go and aws-sdk-go-v2 dependencies.
// It helps reducing the final binary with ~11MB.
//
// In the future this may get replaced entirely with an even slimmer
// version without relying on aws-sdk-go-v2.
//
//--------------------------------------------------------------------
//
// Copyright 2018 The Go Cloud Development Kit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package s3blob provides a blob implementation that uses S3. Use OpenBucket
// to construct a *blob.Bucket.
//
// # URLs
//
// For blob.OpenBucket, s3blob registers for the scheme "s3".
// The default URL opener will use an AWS session with the default credentials
// and configuration; see https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
// for more details.
// Use "awssdk=v1" or "awssdk=v2" to force a specific AWS SDK version.
// To customize the URL opener, or for more details on the URL format,
// see URLOpener.
// See https://gocloud.dev/concepts/urls/ for background information.
//
// # Escaping
//
// Go CDK supports all UTF-8 strings; to make this work with services lacking
// full UTF-8 support, strings must be escaped (during writes) and unescaped
// (during reads). The following escapes are performed for s3blob:
//   - Blob keys: ASCII characters 0-31 are escaped to "__0x<hex>__".
//     Additionally, the "/" in "../" is escaped in the same way.
//   - Metadata keys: Escaped using URL encoding, then additionally "@:=" are
//     escaped using "__0x<hex>__". These characters were determined by
//     experimentation.
//   - Metadata values: Escaped using URL encoding.
//
// # As
//
// s3blob exposes the following types for As:
//   - Bucket: (V1) *s3.S3; (V2) *s3v2.Client
//   - Error: (V1) awserr.Error; (V2) any error type returned by the service, notably smithy.APIError
//   - ListObject: (V1) s3.Object for objects, s3.CommonPrefix for "directories"; (V2) typesv2.Object for objects, typesv2.CommonPrefix for "directories"
//   - ListOptions.BeforeList: (V1) *s3.ListObjectsV2Input or *s3.ListObjectsInput
//     when Options.UseLegacyList == true; (V2) *s3v2.ListObjectsV2Input or *[]func(*s3v2.Options), or *s3v2.ListObjectsInput
//     when Options.UseLegacyList == true
//   - Reader: (V1) s3.GetObjectOutput; (V2) s3v2.GetObjectInput
//   - ReaderOptions.BeforeRead: (V1) *s3.GetObjectInput; (V2) *s3v2.GetObjectInput or *[]func(*s3v2.Options)
//   - Attributes: (V1) s3.HeadObjectOutput; (V2)s3v2.HeadObjectOutput
//   - CopyOptions.BeforeCopy: *(V1) s3.CopyObjectInput; (V2) s3v2.CopyObjectInput
//   - WriterOptions.BeforeWrite: (V1) *s3manager.UploadInput, *s3manager.Uploader; (V2) *s3v2.PutObjectInput, *s3v2manager.Uploader
//   - SignedURLOptions.BeforeSign:
//     (V1) *s3.GetObjectInput; (V2) *s3v2.GetObjectInput, when Options.Method == http.MethodGet, or
//     (V1) *s3.PutObjectInput; (V2) *s3v2.PutObjectInput, when Options.Method == http.MethodPut, or
//     (V1) *s3.DeleteObjectInput; (V2) [not supported] when Options.Method == http.MethodDelete

package filesystem

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	awsv2cfg "github.com/aws/aws-sdk-go-v2/config"
	s3managerv2 "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	typesv2 "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/aws/smithy-go"
	"gocloud.dev/blob"
	"gocloud.dev/blob/driver"
	"gocloud.dev/gcerrors"
)

// -------------------------------------------------------------------

// Copyright 2019 The Go Cloud Development Kit Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// HexEscape returns s, with all runes for which shouldEscape returns true
// escaped to "__0xXXX__", where XXX is the hex representation of the rune
// value. For example, " " would escape to "__0x20__".
//
// Non-UTF-8 strings will have their non-UTF-8 characters escaped to
// unicode.ReplacementChar; the original value is lost. Please file an
// issue if you need non-UTF8 support.
//
// Note: shouldEscape takes the whole string as a slice of runes and an
// index. Passing it a single byte or a single rune doesn't provide
// enough context for some escape decisions; for example, the caller might
// want to escape the second "/" in "//" but not the first one.
// We pass a slice of runes instead of the string or a slice of bytes
// because some decisions will be made on a rune basis (e.g., encode
// all non-ASCII runes).
func HexEscape(s string, shouldEscape func(s []rune, i int) bool) string {
	// Do a first pass to see which runes (if any) need escaping.
	runes := []rune(s)
	var toEscape []int
	for i := range runes {
		if shouldEscape(runes, i) {
			toEscape = append(toEscape, i)
		}
	}
	if len(toEscape) == 0 {
		return s
	}
	// Each escaped rune turns into at most 14 runes ("__0x7fffffff__"),
	// so allocate an extra 13 for each. We'll reslice at the end
	// if we didn't end up using them.
	escaped := make([]rune, len(runes)+13*len(toEscape))
	n := 0 // current index into toEscape
	j := 0 // current index into escaped
	for i, r := range runes {
		if n < len(toEscape) && i == toEscape[n] {
			// We were asked to escape this rune.
			for _, x := range fmt.Sprintf("__%#x__", r) {
				escaped[j] = x
				j++
			}
			n++
		} else {
			escaped[j] = r
			j++
		}
	}
	return string(escaped[0:j])
}

// unescape tries to unescape starting at r[i].
// It returns a boolean indicating whether the unescaping was successful,
// and (if true) the unescaped rune and the last index of r that was used
// during unescaping.
func unescape(r []rune, i int) (bool, rune, int) {
	// Look for "__0x".
	if r[i] != '_' {
		return false, 0, 0
	}
	i++
	if i >= len(r) || r[i] != '_' {
		return false, 0, 0
	}
	i++
	if i >= len(r) || r[i] != '0' {
		return false, 0, 0
	}
	i++
	if i >= len(r) || r[i] != 'x' {
		return false, 0, 0
	}
	i++
	// Capture the digits until the next "_" (if any).
	var hexdigits []rune
	for ; i < len(r) && r[i] != '_'; i++ {
		hexdigits = append(hexdigits, r[i])
	}
	// Look for the trailing "__".
	if i >= len(r) || r[i] != '_' {
		return false, 0, 0
	}
	i++
	if i >= len(r) || r[i] != '_' {
		return false, 0, 0
	}
	// Parse the hex digits into an int32.
	retval, err := strconv.ParseInt(string(hexdigits), 16, 32)
	if err != nil {
		return false, 0, 0
	}
	return true, rune(retval), i
}

// HexUnescape reverses HexEscape.
func HexUnescape(s string) string {
	var unescaped []rune
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if ok, newR, newI := unescape(runes, i); ok {
			// We unescaped some runes starting at i, resulting in the
			// unescaped rune newR. The last rune used was newI.
			if unescaped == nil {
				// This is the first rune we've encountered that
				// needed unescaping. Allocate a buffer and copy any
				// previous runes.
				unescaped = make([]rune, i)
				copy(unescaped, runes)
			}
			unescaped = append(unescaped, newR)
			i = newI
		} else if unescaped != nil {
			unescaped = append(unescaped, runes[i])
		}
	}
	if unescaped == nil {
		return s
	}
	return string(unescaped)
}

// URLEscape uses url.PathEscape to escape s.
func URLEscape(s string) string {
	return url.PathEscape(s)
}

// URLUnescape reverses URLEscape using url.PathUnescape. If the unescape
// returns an error, it returns s.
func URLUnescape(s string) string {
	if u, err := url.PathUnescape(s); err == nil {
		return u
	}
	return s
}

// -------------------------------------------------------------------

// UseV2 returns true iff the URL parameters indicate that the provider
// should use the AWS SDK v2.
//
// "awssdk=v1" will force V1.
// "awssdk=v2" will force V2.
// No "awssdk" parameter (or any other value) will return the default, currently V1.
// Note that the default may change in the future.
func UseV2(q url.Values) bool {
	if values, ok := q["awssdk"]; ok {
		if values[0] == "v2" || values[0] == "V2" {
			return true
		}
	}
	return false
}

// NewDefaultV2Config returns a aws.Config for AWS SDK v2, using the default options.
func NewDefaultV2Config(ctx context.Context) (awsv2.Config, error) {
	return awsv2cfg.LoadDefaultConfig(ctx)
}

// V2ConfigFromURLParams returns an aws.Config for AWS SDK v2 initialized based on the URL
// parameters in q. It is intended to be used by URLOpeners for AWS services if
// UseV2 returns true.
//
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#Config
//
// It returns an error if q contains any unknown query parameters; callers
// should remove any query parameters they know about from q before calling
// V2ConfigFromURLParams.
//
// The following query options are supported:
//   - region: The AWS region for requests; sets WithRegion.
//   - profile: The shared config profile to use; sets SharedConfigProfile.
//   - endpoint: The AWS service endpoint to send HTTP request.
func V2ConfigFromURLParams(ctx context.Context, q url.Values) (awsv2.Config, error) {
	var opts []func(*awsv2cfg.LoadOptions) error
	for param, values := range q {
		value := values[0]
		switch param {
		case "region":
			opts = append(opts, awsv2cfg.WithRegion(value))
		case "endpoint":
			customResolver := awsv2.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (awsv2.Endpoint, error) {
					return awsv2.Endpoint{
						PartitionID:   "aws",
						URL:           value,
						SigningRegion: region,
					}, nil
				})
			opts = append(opts, awsv2cfg.WithEndpointResolverWithOptions(customResolver))
		case "profile":
			opts = append(opts, awsv2cfg.WithSharedConfigProfile(value))
		case "awssdk":
			// ignore, should be handled before this
		default:
			return awsv2.Config{}, fmt.Errorf("unknown query parameter %q", param)
		}
	}
	return awsv2cfg.LoadDefaultConfig(ctx, opts...)
}

// -------------------------------------------------------------------

const defaultPageSize = 1000

func init() {
	blob.DefaultURLMux().RegisterBucket(Scheme, new(urlSessionOpener))
}

type urlSessionOpener struct{}

func (o *urlSessionOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {
	opener := &URLOpener{UseV2: true}
	return opener.OpenBucketURL(ctx, u)
}

// Scheme is the URL scheme s3blob registers its URLOpener under on
// blob.DefaultMux.
const Scheme = "s3"

// URLOpener opens S3 URLs like "s3://mybucket".
//
// The URL host is used as the bucket name.
//
// Use "awssdk=v1" to force using AWS SDK v1, "awssdk=v2" to force using AWS SDK v2,
// or anything else to accept the default.
//
// For V1, see gocloud.dev/aws/ConfigFromURLParams for supported query parameters
// for overriding the aws.Session from the URL.
// For V2, see gocloud.dev/aws/V2ConfigFromURLParams.
type URLOpener struct {
	// UseV2 indicates whether the AWS SDK V2 should be used.
	UseV2 bool

	// Options specifies the options to pass to OpenBucket.
	Options Options
}

// OpenBucketURL opens a blob.Bucket based on u.
func (o *URLOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {
	cfg, err := V2ConfigFromURLParams(ctx, u.Query())
	if err != nil {
		return nil, fmt.Errorf("open bucket %v: %v", u, err)
	}
	clientV2 := s3v2.NewFromConfig(cfg)
	return OpenBucketV2(ctx, clientV2, u.Host, &o.Options)
}

// Options sets options for constructing a *blob.Bucket backed by fileblob.
type Options struct {
	// UseLegacyList forces the use of ListObjects instead of ListObjectsV2.
	// Some S3-compatible services (like CEPH) do not currently support
	// ListObjectsV2.
	UseLegacyList bool
}

// openBucket returns an S3 Bucket.
func openBucket(ctx context.Context, useV2 bool, clientV2 *s3v2.Client, bucketName string, opts *Options) (*bucket, error) {
	if bucketName == "" {
		return nil, errors.New("s3blob.OpenBucket: bucketName is required")
	}
	if opts == nil {
		opts = &Options{}
	}
	if clientV2 == nil {
		return nil, errors.New("s3blob.OpenBucketV2: client is required")
	}
	return &bucket{
		useV2:         useV2,
		name:          bucketName,
		clientV2:      clientV2,
		useLegacyList: opts.UseLegacyList,
	}, nil
}

// OpenBucketV2 returns a *blob.Bucket backed by S3, using AWS SDK v2.
func OpenBucketV2(ctx context.Context, client *s3v2.Client, bucketName string, opts *Options) (*blob.Bucket, error) {
	drv, err := openBucket(ctx, true, client, bucketName, opts)
	if err != nil {
		return nil, err
	}
	return blob.NewBucket(drv), nil
}

// reader reads an S3 object. It implements io.ReadCloser.
type reader struct {
	useV2 bool
	body  io.ReadCloser
	attrs driver.ReaderAttributes
	rawV2 *s3v2.GetObjectOutput
}

func (r *reader) Read(p []byte) (int, error) {
	return r.body.Read(p)
}

// Close closes the reader itself. It must be called when done reading.
func (r *reader) Close() error {
	return r.body.Close()
}

func (r *reader) As(i interface{}) bool {
	p, ok := i.(*s3v2.GetObjectOutput)
	if !ok {
		return false
	}
	*p = *r.rawV2
	return true
}

func (r *reader) Attributes() *driver.ReaderAttributes {
	return &r.attrs
}

// writer writes an S3 object, it implements io.WriteCloser.
type writer struct {
	// Ends of an io.Pipe, created when the first byte is written.
	pw *io.PipeWriter
	pr *io.PipeReader

	// Alternatively, upload is set to true when Upload was
	// used to upload data.
	upload bool

	ctx   context.Context
	useV2 bool

	// v2
	uploaderV2 *s3managerv2.Uploader
	reqV2      *s3v2.PutObjectInput

	donec chan struct{} // closed when done writing
	// The following fields will be written before donec closes:
	err error
}

// Write appends p to w.pw. User must call Close to close the w after done writing.
func (w *writer) Write(p []byte) (int, error) {
	// Avoid opening the pipe for a zero-length write;
	// the concrete can do these for empty blobs.
	if len(p) == 0 {
		return 0, nil
	}
	if w.pw == nil {
		// We'll write into pw and use pr as an io.Reader for the
		// Upload call to S3.
		w.pr, w.pw = io.Pipe()
		w.open(w.pr, true)
	}
	return w.pw.Write(p)
}

// Upload reads from r. Per the driver, it is guaranteed to be the only
// write call for this writer.
func (w *writer) Upload(r io.Reader) error {
	w.upload = true
	w.open(r, false)
	return nil
}

// r may be nil if we're Closing and no data was written.
// If closePipeOnError is true, w.pr will be closed if there's an
// error uploading to S3.
func (w *writer) open(r io.Reader, closePipeOnError bool) {
	// This goroutine will keep running until Close, unless there's an error.
	go func() {
		defer close(w.donec)

		if r == nil {
			// AWS doesn't like a nil Body.
			r = http.NoBody
		}
		var err error
		w.reqV2.Body = r
		_, err = w.uploaderV2.Upload(w.ctx, w.reqV2)
		if err != nil {
			if closePipeOnError {
				w.pr.CloseWithError(err)
				w.pr = nil
			}
			w.err = err
		}
	}()
}

// Close completes the writer and closes it. Any error occurring during write
// will be returned. If a writer is closed before any Write is called, Close
// will create an empty file at the given key.
func (w *writer) Close() error {
	if !w.upload {
		if w.pr != nil {
			defer w.pr.Close()
		}
		if w.pw == nil {
			// We never got any bytes written. We'll write an http.NoBody.
			w.open(nil, false)
		} else if err := w.pw.Close(); err != nil {
			return err
		}
	}
	<-w.donec
	return w.err
}

// bucket represents an S3 bucket and handles read, write and delete operations.
type bucket struct {
	name          string
	useV2         bool
	clientV2      *s3v2.Client
	useLegacyList bool
}

func (b *bucket) Close() error {
	return nil
}

func (b *bucket) ErrorCode(err error) gcerrors.ErrorCode {
	var code string
	var ae smithy.APIError
	var oe *smithy.OperationError
	if errors.As(err, &oe) && strings.Contains(oe.Error(), "301") {
		// V2 returns an OperationError with a missing redirect for invalid buckets.
		code = "NoSuchBucket"
	} else if errors.As(err, &ae) {
		code = ae.ErrorCode()
	} else {
		return gcerrors.Unknown
	}
	switch {
	case code == "NoSuchBucket" || code == "NoSuchKey" || code == "NotFound":
		return gcerrors.NotFound
	default:
		return gcerrors.Unknown
	}
}

// ListPaged implements driver.ListPaged.
func (b *bucket) ListPaged(ctx context.Context, opts *driver.ListOptions) (*driver.ListPage, error) {
	pageSize := opts.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	in := &s3v2.ListObjectsV2Input{
		Bucket:  awsv2.String(b.name),
		MaxKeys: awsv2.Int32(int32(pageSize)),
	}
	if len(opts.PageToken) > 0 {
		in.ContinuationToken = awsv2.String(string(opts.PageToken))
	}
	if opts.Prefix != "" {
		in.Prefix = awsv2.String(escapeKey(opts.Prefix))
	}
	if opts.Delimiter != "" {
		in.Delimiter = awsv2.String(escapeKey(opts.Delimiter))
	}
	resp, err := b.listObjectsV2(ctx, in, opts)
	if err != nil {
		return nil, err
	}
	page := driver.ListPage{}
	if resp.NextContinuationToken != nil {
		page.NextPageToken = []byte(*resp.NextContinuationToken)
	}
	if n := len(resp.Contents) + len(resp.CommonPrefixes); n > 0 {
		page.Objects = make([]*driver.ListObject, n)
		for i, obj := range resp.Contents {
			obj := obj
			page.Objects[i] = &driver.ListObject{
				Key:     unescapeKey(awsv2.ToString(obj.Key)),
				ModTime: *obj.LastModified,
				Size:    awsv2.ToInt64(obj.Size),
				MD5:     eTagToMD5(obj.ETag),
				AsFunc: func(i interface{}) bool {
					p, ok := i.(*typesv2.Object)
					if !ok {
						return false
					}
					*p = obj
					return true
				},
			}
		}
		for i, prefix := range resp.CommonPrefixes {
			prefix := prefix
			page.Objects[i+len(resp.Contents)] = &driver.ListObject{
				Key:   unescapeKey(awsv2.ToString(prefix.Prefix)),
				IsDir: true,
				AsFunc: func(i interface{}) bool {
					p, ok := i.(*typesv2.CommonPrefix)
					if !ok {
						return false
					}
					*p = prefix
					return true
				},
			}
		}
		if len(resp.Contents) > 0 && len(resp.CommonPrefixes) > 0 {
			// S3 gives us blobs and "directories" in separate lists; sort them.
			sort.Slice(page.Objects, func(i, j int) bool {
				return page.Objects[i].Key < page.Objects[j].Key
			})
		}
	}
	return &page, nil
}

func (b *bucket) listObjectsV2(ctx context.Context, in *s3v2.ListObjectsV2Input, opts *driver.ListOptions) (*s3v2.ListObjectsV2Output, error) {
	if !b.useLegacyList {
		var varopt []func(*s3v2.Options)
		if opts.BeforeList != nil {
			asFunc := func(i interface{}) bool {
				if p, ok := i.(**s3v2.ListObjectsV2Input); ok {
					*p = in
					return true
				}
				if p, ok := i.(**[]func(*s3v2.Options)); ok {
					*p = &varopt
					return true
				}
				return false
			}
			if err := opts.BeforeList(asFunc); err != nil {
				return nil, err
			}
		}
		return b.clientV2.ListObjectsV2(ctx, in, varopt...)
	}

	// Use the legacy ListObjects request.
	legacyIn := &s3v2.ListObjectsInput{
		Bucket:       in.Bucket,
		Delimiter:    in.Delimiter,
		EncodingType: in.EncodingType,
		Marker:       in.ContinuationToken,
		MaxKeys:      in.MaxKeys,
		Prefix:       in.Prefix,
		RequestPayer: in.RequestPayer,
	}
	if opts.BeforeList != nil {
		asFunc := func(i interface{}) bool {
			p, ok := i.(**s3v2.ListObjectsInput)
			if !ok {
				return false
			}
			*p = legacyIn
			return true
		}
		if err := opts.BeforeList(asFunc); err != nil {
			return nil, err
		}
	}
	legacyResp, err := b.clientV2.ListObjects(ctx, legacyIn)
	if err != nil {
		return nil, err
	}

	var nextContinuationToken *string
	if legacyResp.NextMarker != nil {
		nextContinuationToken = legacyResp.NextMarker
	} else if awsv2.ToBool(legacyResp.IsTruncated) {
		nextContinuationToken = awsv2.String(awsv2.ToString(legacyResp.Contents[len(legacyResp.Contents)-1].Key))
	}
	return &s3v2.ListObjectsV2Output{
		CommonPrefixes:        legacyResp.CommonPrefixes,
		Contents:              legacyResp.Contents,
		NextContinuationToken: nextContinuationToken,
	}, nil
}

// func (b *bucket) listObjects(ctx context.Context, in *s3.ListObjectsV2Input, opts *driver.ListOptions) (*s3.ListObjectsV2Output, error) {
//  if !b.useLegacyList {
//      if opts.BeforeList != nil {
//          asFunc := func(i interface{}) bool {
//              if p, ok := i.(**s3.ListObjectsV2Input); ok {
//                  *p = in
//                  return true
//              }
//              return false
//          }
//          if err := opts.BeforeList(asFunc); err != nil {
//              return nil, err
//          }
//      }
//      return b.client.ListObjectsV2WithContext(ctx, in)
//  }

//  // Use the legacy ListObjects request.
//  legacyIn := &s3.ListObjectsInput{
//      Bucket:       in.Bucket,
//      Delimiter:    in.Delimiter,
//      EncodingType: in.EncodingType,
//      Marker:       in.ContinuationToken,
//      MaxKeys:      in.MaxKeys,
//      Prefix:       in.Prefix,
//      RequestPayer: in.RequestPayer,
//  }
//  if opts.BeforeList != nil {
//      asFunc := func(i interface{}) bool {
//          p, ok := i.(**s3.ListObjectsInput)
//          if !ok {
//              return false
//          }
//          *p = legacyIn
//          return true
//      }
//      if err := opts.BeforeList(asFunc); err != nil {
//          return nil, err
//      }
//  }
//  legacyResp, err := b.client.ListObjectsWithContext(ctx, legacyIn)
//  if err != nil {
//      return nil, err
//  }

//  var nextContinuationToken *string
//  if legacyResp.NextMarker != nil {
//      nextContinuationToken = legacyResp.NextMarker
//  } else if awsv2.ToBool(legacyResp.IsTruncated) {
//      nextContinuationToken = awsv2.String(awsv2.ToString(legacyResp.Contents[len(legacyResp.Contents)-1].Key))
//  }
//  return &s3.ListObjectsV2Output{
//      CommonPrefixes:        legacyResp.CommonPrefixes,
//      Contents:              legacyResp.Contents,
//      NextContinuationToken: nextContinuationToken,
//  }, nil
// }

// As implements driver.As.
func (b *bucket) As(i interface{}) bool {
	p, ok := i.(**s3v2.Client)
	if !ok {
		return false
	}
	*p = b.clientV2
	return true
}

// As implements driver.ErrorAs.
func (b *bucket) ErrorAs(err error, i interface{}) bool {
	return errors.As(err, i)
}

// Attributes implements driver.Attributes.
func (b *bucket) Attributes(ctx context.Context, key string) (*driver.Attributes, error) {
	key = escapeKey(key)
	in := &s3v2.HeadObjectInput{
		Bucket: awsv2.String(b.name),
		Key:    awsv2.String(key),
	}
	resp, err := b.clientV2.HeadObject(ctx, in)
	if err != nil {
		return nil, err
	}

	md := make(map[string]string, len(resp.Metadata))
	for k, v := range resp.Metadata {
		// See the package comments for more details on escaping of metadata
		// keys & values.
		md[HexUnescape(URLUnescape(k))] = URLUnescape(v)
	}
	return &driver.Attributes{
		CacheControl:       awsv2.ToString(resp.CacheControl),
		ContentDisposition: awsv2.ToString(resp.ContentDisposition),
		ContentEncoding:    awsv2.ToString(resp.ContentEncoding),
		ContentLanguage:    awsv2.ToString(resp.ContentLanguage),
		ContentType:        awsv2.ToString(resp.ContentType),
		Metadata:           md,
		// CreateTime not supported; left as the zero time.
		ModTime: awsv2.ToTime(resp.LastModified),
		Size:    awsv2.ToInt64(resp.ContentLength),
		MD5:     eTagToMD5(resp.ETag),
		ETag:    awsv2.ToString(resp.ETag),
		AsFunc: func(i interface{}) bool {
			p, ok := i.(*s3v2.HeadObjectOutput)
			if !ok {
				return false
			}
			*p = *resp
			return true
		},
	}, nil
}

// NewRangeReader implements driver.NewRangeReader.
func (b *bucket) NewRangeReader(ctx context.Context, key string, offset, length int64, opts *driver.ReaderOptions) (driver.Reader, error) {
	key = escapeKey(key)
	var byteRange *string
	if offset > 0 && length < 0 {
		byteRange = awsv2.String(fmt.Sprintf("bytes=%d-", offset))
	} else if length == 0 {
		// AWS doesn't support a zero-length read; we'll read 1 byte and then
		// ignore it in favor of http.NoBody below.
		byteRange = awsv2.String(fmt.Sprintf("bytes=%d-%d", offset, offset))
	} else if length >= 0 {
		byteRange = awsv2.String(fmt.Sprintf("bytes=%d-%d", offset, offset+length-1))
	}
	in := &s3v2.GetObjectInput{
		Bucket: awsv2.String(b.name),
		Key:    awsv2.String(key),
		Range:  byteRange,
	}
	var varopt []func(*s3v2.Options)
	if opts.BeforeRead != nil {
		asFunc := func(i interface{}) bool {
			if p, ok := i.(**s3v2.GetObjectInput); ok {
				*p = in
				return true
			}
			if p, ok := i.(**[]func(*s3v2.Options)); ok {
				*p = &varopt
				return true
			}
			return false
		}
		if err := opts.BeforeRead(asFunc); err != nil {
			return nil, err
		}
	}
	resp, err := b.clientV2.GetObject(ctx, in, varopt...)
	if err != nil {
		return nil, err
	}
	body := resp.Body
	if length == 0 {
		body = http.NoBody
	}
	return &reader{
		useV2: true,
		body:  body,
		attrs: driver.ReaderAttributes{
			ContentType: awsv2.ToString(resp.ContentType),
			ModTime:     awsv2.ToTime(resp.LastModified),
			Size:        getSize(awsv2.ToInt64(resp.ContentLength), awsv2.ToString(resp.ContentRange)),
		},
		rawV2: resp,
	}, nil
}

// etagToMD5 processes an ETag header and returns an MD5 hash if possible.
// S3's ETag header is sometimes a quoted hexstring of the MD5. Other times,
// notably when the object was uploaded in multiple parts, it is not.
// We do the best we can.
// Some links about ETag:
// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTCommonResponseHeaders.html
// https://github.com/aws/aws-sdk-net/issues/815
// https://teppen.io/2018/06/23/aws_s3_etags/
func eTagToMD5(etag *string) []byte {
	if etag == nil {
		// No header at all.
		return nil
	}
	// Strip the expected leading and trailing quotes.
	quoted := *etag
	if len(quoted) < 2 || quoted[0] != '"' || quoted[len(quoted)-1] != '"' {
		return nil
	}
	unquoted := quoted[1 : len(quoted)-1]
	// Un-hex; we return nil on error. In particular, we'll get an error here
	// for multi-part uploaded blobs, whose ETag will contain a "-" and so will
	// never be a legal hex encoding.
	md5, err := hex.DecodeString(unquoted)
	if err != nil {
		return nil
	}
	return md5
}

func getSize(contentLength int64, contentRange string) int64 {
	// Default size to ContentLength, but that's incorrect for partial-length reads,
	// where ContentLength refers to the size of the returned Body, not the entire
	// size of the blob. ContentRange has the full size.
	size := contentLength
	if contentRange != "" {
		// Sample: bytes 10-14/27 (where 27 is the full size).
		parts := strings.Split(contentRange, "/")
		if len(parts) == 2 {
			if i, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
				size = i
			}
		}
	}
	return size
}

// escapeKey does all required escaping for UTF-8 strings to work with S3.
func escapeKey(key string) string {
	return HexEscape(key, func(r []rune, i int) bool {
		c := r[i]
		switch {
		// S3 doesn't handle these characters (determined via experimentation).
		case c < 32:
			return true
		// For "../", escape the trailing slash.
		case i > 1 && c == '/' && r[i-1] == '.' && r[i-2] == '.':
			return true
		}
		return false
	})
}

// unescapeKey reverses escapeKey.
func unescapeKey(key string) string {
	return HexUnescape(key)
}

// NewTypedWriter implements driver.NewTypedWriter.
func (b *bucket) NewTypedWriter(ctx context.Context, key string, contentType string, opts *driver.WriterOptions) (driver.Writer, error) {
	key = escapeKey(key)
	uploaderV2 := s3managerv2.NewUploader(b.clientV2, func(u *s3managerv2.Uploader) {
		if opts.BufferSize != 0 {
			u.PartSize = int64(opts.BufferSize)
		}
		if opts.MaxConcurrency != 0 {
			u.Concurrency = opts.MaxConcurrency
		}
	})
	md := make(map[string]string, len(opts.Metadata))
	for k, v := range opts.Metadata {
		// See the package comments for more details on escaping of metadata
		// keys & values.
		k = HexEscape(url.PathEscape(k), func(runes []rune, i int) bool {
			c := runes[i]
			return c == '@' || c == ':' || c == '='
		})
		md[k] = url.PathEscape(v)
	}
	reqV2 := &s3v2.PutObjectInput{
		Bucket:      awsv2.String(b.name),
		ContentType: awsv2.String(contentType),
		Key:         awsv2.String(key),
		Metadata:    md,
	}
	if opts.CacheControl != "" {
		reqV2.CacheControl = awsv2.String(opts.CacheControl)
	}
	if opts.ContentDisposition != "" {
		reqV2.ContentDisposition = awsv2.String(opts.ContentDisposition)
	}
	if opts.ContentEncoding != "" {
		reqV2.ContentEncoding = awsv2.String(opts.ContentEncoding)
	}
	if opts.ContentLanguage != "" {
		reqV2.ContentLanguage = awsv2.String(opts.ContentLanguage)
	}
	if len(opts.ContentMD5) > 0 {
		reqV2.ContentMD5 = awsv2.String(base64.StdEncoding.EncodeToString(opts.ContentMD5))
	}
	if opts.BeforeWrite != nil {
		asFunc := func(i interface{}) bool {
			// Note that since the Go CDK Blob
			// abstraction does not expose AWS's
			// Uploader concept, there does not
			// appear to be any utility in
			// exposing the options list to the v2
			// Uploader's Upload() method.
			// Instead, applications can
			// manipulate the exposed *Uploader
			// directly, including by setting
			// ClientOptions if needed.
			if p, ok := i.(**s3managerv2.Uploader); ok {
				*p = uploaderV2
				return true
			}
			if p, ok := i.(**s3v2.PutObjectInput); ok {
				*p = reqV2
				return true
			}
			return false
		}
		if err := opts.BeforeWrite(asFunc); err != nil {
			return nil, err
		}
	}
	return &writer{
		ctx:        ctx,
		useV2:      true,
		uploaderV2: uploaderV2,
		reqV2:      reqV2,
		donec:      make(chan struct{}),
	}, nil
}

// Copy implements driver.Copy.
func (b *bucket) Copy(ctx context.Context, dstKey, srcKey string, opts *driver.CopyOptions) error {
	dstKey = escapeKey(dstKey)
	srcKey = escapeKey(srcKey)
	input := &s3v2.CopyObjectInput{
		Bucket:     awsv2.String(b.name),
		CopySource: awsv2.String(b.name + "/" + srcKey),
		Key:        awsv2.String(dstKey),
	}
	if opts.BeforeCopy != nil {
		asFunc := func(i interface{}) bool {
			switch v := i.(type) {
			case **s3v2.CopyObjectInput:
				*v = input
				return true
			}
			return false
		}
		if err := opts.BeforeCopy(asFunc); err != nil {
			return err
		}
	}
	_, err := b.clientV2.CopyObject(ctx, input)
	return err
}

// Delete implements driver.Delete.
func (b *bucket) Delete(ctx context.Context, key string) error {
	if _, err := b.Attributes(ctx, key); err != nil {
		return err
	}
	key = escapeKey(key)
	input := &s3v2.DeleteObjectInput{
		Bucket: awsv2.String(b.name),
		Key:    awsv2.String(key),
	}
	_, err := b.clientV2.DeleteObject(ctx, input)
	return err
}

func (b *bucket) SignedURL(ctx context.Context, key string, opts *driver.SignedURLOptions) (string, error) {
	key = escapeKey(key)
	switch opts.Method {
	case http.MethodGet:
		in := &s3v2.GetObjectInput{
			Bucket: awsv2.String(b.name),
			Key:    awsv2.String(key),
		}
		if opts.BeforeSign != nil {
			asFunc := func(i interface{}) bool {
				v, ok := i.(**s3v2.GetObjectInput)
				if ok {
					*v = in
				}
				return ok
			}
			if err := opts.BeforeSign(asFunc); err != nil {
				return "", err
			}
		}
		p, err := s3v2.NewPresignClient(b.clientV2, s3v2.WithPresignExpires(opts.Expiry)).PresignGetObject(ctx, in)
		if err != nil {
			return "", err
		}
		return p.URL, nil
	case http.MethodPut:
		in := &s3v2.PutObjectInput{
			Bucket: awsv2.String(b.name),
			Key:    awsv2.String(key),
		}
		if opts.EnforceAbsentContentType || opts.ContentType != "" {
			// https://github.com/aws/aws-sdk-go-v2/issues/1475
			return "", errors.New("s3blob: AWS SDK v2 does not supported enforcing ContentType in SignedURLs for PUT")
		}
		if opts.BeforeSign != nil {
			asFunc := func(i interface{}) bool {
				v, ok := i.(**s3v2.PutObjectInput)
				if ok {
					*v = in
				}
				return ok
			}
			if err := opts.BeforeSign(asFunc); err != nil {
				return "", err
			}
		}
		p, err := s3v2.NewPresignClient(b.clientV2, s3v2.WithPresignExpires(opts.Expiry)).PresignPutObject(ctx, in)
		if err != nil {
			return "", err
		}
		return p.URL, nil
	}
	return "", fmt.Errorf("unsupported Method %q", opts.Method)
}
