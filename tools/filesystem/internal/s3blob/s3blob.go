// Package s3blob provides a blob.Bucket S3 driver implementation.
//
// NB! To minimize breaking changes with older PocketBase releases,
// the driver is based of the previously used gocloud.dev/blob/s3blob,
// hence many of the below doc comments, struct options and interface
// implementations are the same.
//
// The blob abstraction supports all UTF-8 strings; to make this work with services lacking
// full UTF-8 support, strings must be escaped (during writes) and unescaped
// (during reads). The following escapes are performed for s3blob:
//   - Blob keys: ASCII characters 0-31 are escaped to "__0x<hex>__".
//     Additionally, the "/" in "../" is escaped in the same way.
//   - Metadata keys: Escaped using URL encoding, then additionally "@:=" are
//     escaped using "__0x<hex>__". These characters were determined by
//     experimentation.
//   - Metadata values: Escaped using URL encoding.
//
// Example:
//
//	drv, _ := s3blob.New(&s3.S3{
//		Bucket:    "bucketName",
//		Region:    "region",
//		Endpoint:  "endpoint",
//		AccessKey: "accessKey",
//		SecretKey: "secretKey",
//	})
//	bucket := blob.NewBucket(drv)
package s3blob

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

	"github.com/pocketbase/pocketbase/tools/filesystem/blob"
	"github.com/pocketbase/pocketbase/tools/filesystem/internal/s3blob/s3"
)

const defaultPageSize = 1000

// New creates a new instance of the S3 driver backed by the internal S3 client.
func New(s3Client *s3.S3) (blob.Driver, error) {
	if s3Client.Bucket == "" {
		return nil, errors.New("s3blob.New: missing bucket name")
	}

	if s3Client.Endpoint == "" {
		return nil, errors.New("s3blob.New: missing endpoint")
	}

	if s3Client.Region == "" {
		return nil, errors.New("s3blob.New: missing region")
	}

	return &driver{s3: s3Client}, nil
}

type driver struct {
	s3 *s3.S3
}

// Close implements [blob/Driver.Close].
func (drv *driver) Close() error {
	return nil // nothing to close
}

// NormalizeError implements [blob/Driver.NormalizeError].
func (drv *driver) NormalizeError(err error) error {
	// already normalized
	if errors.Is(err, blob.ErrNotFound) {
		return err
	}

	// normalize base on its S3 error status or code
	var ae *s3.ResponseError
	if errors.As(err, &ae) {
		if ae.Status == 404 {
			return errors.Join(err, blob.ErrNotFound)
		}

		switch ae.Code {
		case "NoSuchBucket", "NoSuchKey", "NotFound":
			return errors.Join(err, blob.ErrNotFound)
		}
	}

	return err
}

// ListPaged implements [blob/Driver.ListPaged].
func (drv *driver) ListPaged(ctx context.Context, opts *blob.ListOptions) (*blob.ListPage, error) {
	pageSize := opts.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	listParams := s3.ListParams{
		MaxKeys: pageSize,
	}
	if len(opts.PageToken) > 0 {
		listParams.ContinuationToken = string(opts.PageToken)
	}
	if opts.Prefix != "" {
		listParams.Prefix = escapeKey(opts.Prefix)
	}
	if opts.Delimiter != "" {
		listParams.Delimiter = escapeKey(opts.Delimiter)
	}

	resp, err := drv.s3.ListObjects(ctx, listParams)
	if err != nil {
		return nil, err
	}

	page := blob.ListPage{}
	if resp.NextContinuationToken != "" {
		page.NextPageToken = []byte(resp.NextContinuationToken)
	}

	if n := len(resp.Contents) + len(resp.CommonPrefixes); n > 0 {
		page.Objects = make([]*blob.ListObject, n)
		for i, obj := range resp.Contents {
			page.Objects[i] = &blob.ListObject{
				Key:     unescapeKey(obj.Key),
				ModTime: obj.LastModified,
				Size:    obj.Size,
				MD5:     eTagToMD5(obj.ETag),
			}
		}

		for i, prefix := range resp.CommonPrefixes {
			page.Objects[i+len(resp.Contents)] = &blob.ListObject{
				Key:   unescapeKey(prefix.Prefix),
				IsDir: true,
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

// Attributes implements [blob/Driver.Attributes].
func (drv *driver) Attributes(ctx context.Context, key string) (*blob.Attributes, error) {
	key = escapeKey(key)

	resp, err := drv.s3.HeadObject(ctx, key)
	if err != nil {
		return nil, err
	}

	md := make(map[string]string, len(resp.Metadata))
	for k, v := range resp.Metadata {
		// See the package comments for more details on escaping of metadata keys & values.
		md[blob.HexUnescape(urlUnescape(k))] = urlUnescape(v)
	}

	return &blob.Attributes{
		CacheControl:       resp.CacheControl,
		ContentDisposition: resp.ContentDisposition,
		ContentEncoding:    resp.ContentEncoding,
		ContentLanguage:    resp.ContentLanguage,
		ContentType:        resp.ContentType,
		Metadata:           md,
		// CreateTime not supported; left as the zero time.
		ModTime: resp.LastModified,
		Size:    resp.ContentLength,
		MD5:     eTagToMD5(resp.ETag),
		ETag:    resp.ETag,
	}, nil
}

// NewRangeReader implements [blob/Driver.NewRangeReader].
func (drv *driver) NewRangeReader(ctx context.Context, key string, offset, length int64) (blob.DriverReader, error) {
	key = escapeKey(key)

	var byteRange string
	if offset > 0 && length < 0 {
		byteRange = fmt.Sprintf("bytes=%d-", offset)
	} else if length == 0 {
		// AWS doesn't support a zero-length read; we'll read 1 byte and then
		// ignore it in favor of http.NoBody below.
		byteRange = fmt.Sprintf("bytes=%d-%d", offset, offset)
	} else if length >= 0 {
		byteRange = fmt.Sprintf("bytes=%d-%d", offset, offset+length-1)
	}

	reqOpt := func(req *http.Request) {
		if byteRange != "" {
			req.Header.Set("Range", byteRange)
		}
	}

	resp, err := drv.s3.GetObject(ctx, key, reqOpt)
	if err != nil {
		return nil, err
	}

	body := resp.Body
	if length == 0 {
		body = http.NoBody
	}

	return &reader{
		body: body,
		attrs: &blob.ReaderAttributes{
			ContentType: resp.ContentType,
			ModTime:     resp.LastModified,
			Size:        getSize(resp.ContentLength, resp.ContentRange),
		},
	}, nil
}

// NewTypedWriter implements [blob/Driver.NewTypedWriter].
func (drv *driver) NewTypedWriter(ctx context.Context, key string, contentType string, opts *blob.WriterOptions) (blob.DriverWriter, error) {
	key = escapeKey(key)

	u := &s3.Uploader{
		S3:  drv.s3,
		Key: key,
	}

	if opts.BufferSize != 0 {
		u.MinPartSize = opts.BufferSize
	}

	if opts.MaxConcurrency != 0 {
		u.MaxConcurrency = opts.MaxConcurrency
	}

	md := make(map[string]string, len(opts.Metadata))
	for k, v := range opts.Metadata {
		// See the package comments for more details on escaping of metadata keys & values.
		k = blob.HexEscape(url.PathEscape(k), func(runes []rune, i int) bool {
			c := runes[i]
			return c == '@' || c == ':' || c == '='
		})
		md[k] = url.PathEscape(v)
	}
	u.Metadata = md

	var reqOptions []func(*http.Request)
	reqOptions = append(reqOptions, func(r *http.Request) {
		r.Header.Set("Content-Type", contentType)

		if opts.CacheControl != "" {
			r.Header.Set("Cache-Control", opts.CacheControl)
		}
		if opts.ContentDisposition != "" {
			r.Header.Set("Content-Disposition", opts.ContentDisposition)
		}
		if opts.ContentEncoding != "" {
			r.Header.Set("Content-Encoding", opts.ContentEncoding)
		}
		if opts.ContentLanguage != "" {
			r.Header.Set("Content-Language", opts.ContentLanguage)
		}
		if len(opts.ContentMD5) > 0 {
			r.Header.Set("Content-MD5", base64.StdEncoding.EncodeToString(opts.ContentMD5))
		}
	})

	return &writer{
		ctx:        ctx,
		uploader:   u,
		donec:      make(chan struct{}),
		reqOptions: reqOptions,
	}, nil
}

// Copy implements [blob/Driver.Copy].
func (drv *driver) Copy(ctx context.Context, dstKey, srcKey string) error {
	dstKey = escapeKey(dstKey)
	srcKey = escapeKey(srcKey)
	_, err := drv.s3.CopyObject(ctx, srcKey, dstKey)
	return err
}

// Delete implements [blob/Driver.Delete].
func (drv *driver) Delete(ctx context.Context, key string) error {
	key = escapeKey(key)
	return drv.s3.DeleteObject(ctx, key)
}

// -------------------------------------------------------------------

// reader reads an S3 object. It implements io.ReadCloser.
type reader struct {
	attrs *blob.ReaderAttributes
	body  io.ReadCloser
}

// Read implements [io/ReadCloser.Read].
func (r *reader) Read(p []byte) (int, error) {
	return r.body.Read(p)
}

// Close closes the reader itself. It must be called when done reading.
func (r *reader) Close() error {
	return r.body.Close()
}

// Attributes implements [blob/DriverReader.Attributes].
func (r *reader) Attributes() *blob.ReaderAttributes {
	return r.attrs
}

// -------------------------------------------------------------------

// writer writes an S3 object, it implements io.WriteCloser.
type writer struct {
	ctx      context.Context
	err      error // written before donec closes
	uploader *s3.Uploader

	// Ends of an io.Pipe, created when the first byte is written.
	pw *io.PipeWriter
	pr *io.PipeReader

	donec chan struct{} // closed when done writing

	reqOptions []func(*http.Request)
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

// r may be nil if we're Closing and no data was written.
// If closePipeOnError is true, w.pr will be closed if there's an
// error uploading to S3.
func (w *writer) open(r io.Reader, closePipeOnError bool) {
	// This goroutine will keep running until Close, unless there's an error.
	go func() {
		defer func() {
			close(w.donec)
		}()

		if r == nil {
			// AWS doesn't like a nil Body.
			r = http.NoBody
		}

		w.uploader.Payload = r

		err := w.uploader.Upload(w.ctx, w.reqOptions...)
		if err != nil {
			if closePipeOnError {
				w.pr.CloseWithError(err)
			}
			w.err = err
		}
	}()
}

// Close completes the writer and closes it. Any error occurring during write
// will be returned. If a writer is closed before any Write is called, Close
// will create an empty file at the given key.
func (w *writer) Close() error {
	if w.pr != nil {
		defer w.pr.Close()
	}

	if w.pw == nil {
		// We never got any bytes written. We'll write an http.NoBody.
		w.open(nil, false)
	} else if err := w.pw.Close(); err != nil {
		return err
	}

	<-w.donec

	return w.err
}

// -------------------------------------------------------------------

// etagToMD5 processes an ETag header and returns an MD5 hash if possible.
// S3's ETag header is sometimes a quoted hexstring of the MD5. Other times,
// notably when the object was uploaded in multiple parts, it is not.
// We do the best we can.
// Some links about ETag:
// https://docs.aws.amazon.com/AmazonS3/latest/API/RESTCommonResponseHeaders.html
// https://github.com/aws/aws-sdk-net/issues/815
// https://teppen.io/2018/06/23/aws_s3_etags/
func eTagToMD5(etag string) []byte {
	// No header at all.
	if etag == "" {
		return nil
	}

	// Strip the expected leading and trailing quotes.
	if len(etag) < 2 || etag[0] != '"' || etag[len(etag)-1] != '"' {
		return nil
	}
	unquoted := etag[1 : len(etag)-1]

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
	return blob.HexEscape(key, func(r []rune, i int) bool {
		c := r[i]

		// S3 doesn't handle these characters (determined via experimentation).
		if c < 32 {
			return true
		}

		// For "../", escape the trailing slash.
		if i > 1 && c == '/' && r[i-1] == '.' && r[i-2] == '.' {
			return true
		}

		return false
	})
}

// unescapeKey reverses escapeKey.
func unescapeKey(key string) string {
	return blob.HexUnescape(key)
}

// urlUnescape reverses URLEscape using url.PathUnescape. If the unescape
// returns an error, it returns s.
func urlUnescape(s string) string {
	if u, err := url.PathUnescape(s); err == nil {
		return u
	}

	return s
}
