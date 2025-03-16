// Package blob defines a lightweight abstration for interacting with
// various storage services (local filesystem, S3, etc.).
//
// NB!
// For compatibility with earlier PocketBase versions and to prevent
// unnecessary breaking changes, this package is based and implemented
// as a minimal, stripped down version of the previously used gocloud.dev/blob.
// While there is no promise that it won't diverge in the future to accommodate
// better some PocketBase specific use cases, currently it copies and
// tries to follow as close as possible the same implementations,
// conventions and rules for the key escaping/unescaping, blob read/write
// interfaces and struct options as gocloud.dev/blob, therefore the
// credits goes to the original Go Cloud Development Kit Authors.
package blob

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrClosed   = errors.New("bucket or blob is closed")
)

// Bucket provides an easy and portable way to interact with blobs
// within a "bucket", including read, write, and list operations.
// To create a Bucket, use constructors found in driver subpackages.
type Bucket struct {
	drv Driver

	// mu protects the closed variable.
	// Read locks are kept to allow holding a read lock for long-running calls,
	// and thereby prevent closing until a call finishes.
	mu     sync.RWMutex
	closed bool
}

// NewBucket creates a new *Bucket based on a specific driver implementation.
func NewBucket(drv Driver) *Bucket {
	return &Bucket{drv: drv}
}

// ListOptions sets options for listing blobs via Bucket.List.
type ListOptions struct {
	// Prefix indicates that only blobs with a key starting with this prefix
	// should be returned.
	Prefix string

	// Delimiter sets the delimiter used to define a hierarchical namespace,
	// like a filesystem with "directories". It is highly recommended that you
	// use "" or "/" as the Delimiter. Other values should work through this API,
	// but service UIs generally assume "/".
	//
	// An empty delimiter means that the bucket is treated as a single flat
	// namespace.
	//
	// A non-empty delimiter means that any result with the delimiter in its key
	// after Prefix is stripped will be returned with ListObject.IsDir = true,
	// ListObject.Key truncated after the delimiter, and zero values for other
	// ListObject fields. These results represent "directories". Multiple results
	// in a "directory" are returned as a single result.
	Delimiter string

	// PageSize sets the maximum number of objects to be returned.
	// 0 means no maximum; driver implementations should choose a reasonable
	// max. It is guaranteed to be >= 0.
	PageSize int

	// PageToken may be filled in with the NextPageToken from a previous
	// ListPaged call.
	PageToken []byte
}

// ListPage represents a page of results return from ListPaged.
type ListPage struct {
	// Objects is the slice of objects found. If ListOptions.PageSize > 0,
	// it should have at most ListOptions.PageSize entries.
	//
	// Objects should be returned in lexicographical order of UTF-8 encoded keys,
	// including across pages. I.e., all objects returned from a ListPage request
	// made using a PageToken from a previous ListPage request's NextPageToken
	// should have Key >= the Key for all objects from the previous request.
	Objects []*ListObject `json:"objects"`

	// NextPageToken should be left empty unless there are more objects
	// to return. The value may be returned as ListOptions.PageToken on a
	// subsequent ListPaged call, to fetch the next page of results.
	// It can be an arbitrary []byte; it need not be a valid key.
	NextPageToken []byte `json:"nextPageToken"`
}

// ListIterator iterates over List results.
type ListIterator struct {
	b       *Bucket
	opts    *ListOptions
	page    *ListPage
	nextIdx int
}

// Next returns a *ListObject for the next blob.
// It returns (nil, io.EOF) if there are no more.
func (i *ListIterator) Next(ctx context.Context) (*ListObject, error) {
	if i.page != nil {
		// We've already got a page of results.
		if i.nextIdx < len(i.page.Objects) {
			// Next object is in the page; return it.
			dobj := i.page.Objects[i.nextIdx]
			i.nextIdx++
			return &ListObject{
				Key:     dobj.Key,
				ModTime: dobj.ModTime,
				Size:    dobj.Size,
				MD5:     dobj.MD5,
				IsDir:   dobj.IsDir,
			}, nil
		}

		if len(i.page.NextPageToken) == 0 {
			// Done with current page, and there are no more; return io.EOF.
			return nil, io.EOF
		}

		// We need to load the next page.
		i.opts.PageToken = i.page.NextPageToken
	}

	i.b.mu.RLock()
	defer i.b.mu.RUnlock()

	if i.b.closed {
		return nil, ErrClosed
	}

	// Loading a new page.
	p, err := i.b.drv.ListPaged(ctx, i.opts)
	if err != nil {
		return nil, wrapError(i.b.drv, err, "")
	}

	i.page = p
	i.nextIdx = 0

	return i.Next(ctx)
}

// ListObject represents a single blob returned from List.
type ListObject struct {
	// Key is the key for this blob.
	Key string `json:"key"`

	// ModTime is the time the blob was last modified.
	ModTime time.Time `json:"modTime"`

	// Size is the size of the blob's content in bytes.
	Size int64 `json:"size"`

	// MD5 is an MD5 hash of the blob contents or nil if not available.
	MD5 []byte `json:"md5"`

	// IsDir indicates that this result represents a "directory" in the
	// hierarchical namespace, ending in ListOptions.Delimiter. Key can be
	// passed as ListOptions.Prefix to list items in the "directory".
	// Fields other than Key and IsDir will not be set if IsDir is true.
	IsDir bool `json:"isDir"`
}

// List returns a ListIterator that can be used to iterate over blobs in a
// bucket, in lexicographical order of UTF-8 encoded keys. The underlying
// implementation fetches results in pages.
//
// A nil ListOptions is treated the same as the zero value.
//
// List is not guaranteed to include all recently-written blobs;
// some services are only eventually consistent.
func (b *Bucket) List(opts *ListOptions) *ListIterator {
	if opts == nil {
		opts = &ListOptions{}
	}

	dopts := &ListOptions{
		Prefix:    opts.Prefix,
		Delimiter: opts.Delimiter,
	}

	return &ListIterator{b: b, opts: dopts}
}

// FirstPageToken is the pageToken to pass to ListPage to retrieve the first page of results.
var FirstPageToken = []byte("first page")

// ListPage returns a page of ListObject results for blobs in a bucket, in lexicographical
// order of UTF-8 encoded keys.
//
// To fetch the first page, pass FirstPageToken as the pageToken. For subsequent pages, pass
// the pageToken returned from a previous call to ListPage.
// It is not possible to "skip ahead" pages.
//
// Each call will return pageSize results, unless there are not enough blobs to fill the
// page, in which case it will return fewer results (possibly 0).
//
// If there are no more blobs available, ListPage will return an empty pageToken. Note that
// this may happen regardless of the number of returned results -- the last page might have
// 0 results (i.e., if the last item was deleted), pageSize results, or anything in between.
//
// Calling ListPage with an empty pageToken will immediately return io.EOF. When looping
// over pages, callers can either check for an empty pageToken, or they can make one more
// call and check for io.EOF.
//
// The underlying implementation fetches results in pages, but one call to ListPage may
// require multiple page fetches (and therefore, multiple calls to the BeforeList callback).
//
// A nil ListOptions is treated the same as the zero value.
//
// ListPage is not guaranteed to include all recently-written blobs;
// some services are only eventually consistent.
func (b *Bucket) ListPage(ctx context.Context, pageToken []byte, pageSize int, opts *ListOptions) (retval []*ListObject, nextPageToken []byte, err error) {
	if opts == nil {
		opts = &ListOptions{}
	}
	if pageSize <= 0 {
		return nil, nil, fmt.Errorf("pageSize must be > 0 (%d)", pageSize)
	}

	// Nil pageToken means no more results.
	if len(pageToken) == 0 {
		return nil, nil, io.EOF
	}

	// FirstPageToken fetches the first page. Drivers use nil.
	// The public API doesn't use nil for the first page because it would be too easy to
	// keep fetching forever (since the last page return nil for the next pageToken).
	if bytes.Equal(pageToken, FirstPageToken) {
		pageToken = nil
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return nil, nil, ErrClosed
	}

	dopts := &ListOptions{
		Prefix:    opts.Prefix,
		Delimiter: opts.Delimiter,
		PageToken: pageToken,
		PageSize:  pageSize,
	}

	retval = make([]*ListObject, 0, pageSize)
	for len(retval) < pageSize {
		p, err := b.drv.ListPaged(ctx, dopts)
		if err != nil {
			return nil, nil, wrapError(b.drv, err, "")
		}

		for _, dobj := range p.Objects {
			retval = append(retval, &ListObject{
				Key:     dobj.Key,
				ModTime: dobj.ModTime,
				Size:    dobj.Size,
				MD5:     dobj.MD5,
				IsDir:   dobj.IsDir,
			})
		}

		// ListPaged may return fewer results than pageSize. If there are more results
		// available, signalled by non-empty p.NextPageToken, try to fetch the remainder
		// of the page.
		// It does not work to ask for more results than we need, because then we'd have
		// a NextPageToken on a non-page boundary.
		dopts.PageSize = pageSize - len(retval)
		dopts.PageToken = p.NextPageToken
		if len(dopts.PageToken) == 0 {
			dopts.PageToken = nil
			break
		}
	}

	return retval, dopts.PageToken, nil
}

// Attributes contains attributes about a blob.
type Attributes struct {
	// CacheControl specifies caching attributes that services may use
	// when serving the blob.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
	CacheControl string `json:"cacheControl"`

	// ContentDisposition specifies whether the blob content is expected to be
	// displayed inline or as an attachment.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
	ContentDisposition string `json:"contentDisposition"`

	// ContentEncoding specifies the encoding used for the blob's content, if any.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
	ContentEncoding string `json:"contentEncoding"`

	// ContentLanguage specifies the language used in the blob's content, if any.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Language
	ContentLanguage string `json:"contentLanguage"`

	// ContentType is the MIME type of the blob. It will not be empty.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	ContentType string `json:"contentType"`

	// Metadata holds key/value pairs associated with the blob.
	// Keys are guaranteed to be in lowercase, even if the backend service
	// has case-sensitive keys (although note that Metadata written via
	// this package will always be lowercased). If there are duplicate
	// case-insensitive keys (e.g., "foo" and "FOO"), only one value
	// will be kept, and it is undefined which one.
	Metadata map[string]string `json:"metadata"`

	// CreateTime is the time the blob was created, if available. If not available,
	// CreateTime will be the zero time.
	CreateTime time.Time `json:"createTime"`

	// ModTime is the time the blob was last modified.
	ModTime time.Time `json:"modTime"`

	// Size is the size of the blob's content in bytes.
	Size int64 `json:"size"`

	// MD5 is an MD5 hash of the blob contents or nil if not available.
	MD5 []byte `json:"md5"`

	// ETag for the blob; see https://en.wikipedia.org/wiki/HTTP_ETag.
	ETag string `json:"etag"`
}

// Attributes returns attributes for the blob stored at key.
//
// If the blob does not exist, Attributes returns an error for which
// gcerrors.Code will return gcerrors.NotFound.
func (b *Bucket) Attributes(ctx context.Context, key string) (_ *Attributes, err error) {
	if !utf8.ValidString(key) {
		return nil, fmt.Errorf("Attributes key must be a valid UTF-8 string: %q", key)
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return nil, ErrClosed
	}

	a, err := b.drv.Attributes(ctx, key)
	if err != nil {
		return nil, wrapError(b.drv, err, key)
	}

	var md map[string]string
	if len(a.Metadata) > 0 {
		// Services are inconsistent, but at least some treat keys
		// as case-insensitive. To make the behavior consistent, we
		// force-lowercase them when writing and reading.
		md = make(map[string]string, len(a.Metadata))
		for k, v := range a.Metadata {
			md[strings.ToLower(k)] = v
		}
	}

	return &Attributes{
		CacheControl:       a.CacheControl,
		ContentDisposition: a.ContentDisposition,
		ContentEncoding:    a.ContentEncoding,
		ContentLanguage:    a.ContentLanguage,
		ContentType:        a.ContentType,
		Metadata:           md,
		CreateTime:         a.CreateTime,
		ModTime:            a.ModTime,
		Size:               a.Size,
		MD5:                a.MD5,
		ETag:               a.ETag,
	}, nil
}

// Exists returns true if a blob exists at key, false if it does not exist, or
// an error.
//
// It is a shortcut for calling Attributes and checking if it returns an error
// with code ErrNotFound.
func (b *Bucket) Exists(ctx context.Context, key string) (bool, error) {
	_, err := b.Attributes(ctx, key)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, ErrNotFound) {
		return false, nil
	}

	return false, err
}

// NewReader is a shortcut for NewRangeReader with offset=0 and length=-1.
func (b *Bucket) NewReader(ctx context.Context, key string) (*Reader, error) {
	return b.newRangeReader(ctx, key, 0, -1)
}

// NewRangeReader returns a Reader to read content from the blob stored at key.
// It reads at most length bytes starting at offset (>= 0).
// If length is negative, it will read till the end of the blob.
//
// For the purposes of Seek, the returned Reader will start at offset and
// end at the minimum of the actual end of the blob or (if length > 0) offset + length.
//
// Note that ctx is used for all reads performed during the lifetime of the reader.
//
// If the blob does not exist, NewRangeReader returns an error for which
// gcerrors.Code will return gcerrors.NotFound. Exists is a lighter-weight way
// to check for existence.
//
// A nil ReaderOptions is treated the same as the zero value.
//
// The caller must call Close on the returned Reader when done reading.
func (b *Bucket) NewRangeReader(ctx context.Context, key string, offset, length int64) (_ *Reader, err error) {
	return b.newRangeReader(ctx, key, offset, length)
}

func (b *Bucket) newRangeReader(ctx context.Context, key string, offset, length int64) (_ *Reader, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return nil, ErrClosed
	}

	if offset < 0 {
		return nil, fmt.Errorf("NewRangeReader offset must be non-negative (%d)", offset)
	}

	if !utf8.ValidString(key) {
		return nil, fmt.Errorf("NewRangeReader key must be a valid UTF-8 string: %q", key)
	}

	var dr DriverReader
	dr, err = b.drv.NewRangeReader(ctx, key, offset, length)
	if err != nil {
		return nil, wrapError(b.drv, err, key)
	}

	r := &Reader{
		drv:         b.drv,
		r:           dr,
		key:         key,
		ctx:         ctx,
		baseOffset:  offset,
		baseLength:  length,
		savedOffset: -1,
	}

	_, file, lineno, ok := runtime.Caller(2)
	runtime.SetFinalizer(r, func(r *Reader) {
		if !r.closed {
			var caller string
			if ok {
				caller = fmt.Sprintf(" (%s:%d)", file, lineno)
			}
			log.Printf("A blob.Reader reading from %q was never closed%s", key, caller)
		}
	})

	return r, nil
}

// WriterOptions sets options for NewWriter.
type WriterOptions struct {
	// BufferSize changes the default size in bytes of the chunks that
	// Writer will upload in a single request; larger blobs will be split into
	// multiple requests.
	//
	// This option may be ignored by some drivers.
	//
	// If 0, the driver will choose a reasonable default.
	//
	// If the Writer is used to do many small writes concurrently, using a
	// smaller BufferSize may reduce memory usage.
	BufferSize int

	// MaxConcurrency changes the default concurrency for parts of an upload.
	//
	// This option may be ignored by some drivers.
	//
	// If 0, the driver will choose a reasonable default.
	MaxConcurrency int

	// CacheControl specifies caching attributes that services may use
	// when serving the blob.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
	CacheControl string

	// ContentDisposition specifies whether the blob content is expected to be
	// displayed inline or as an attachment.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
	ContentDisposition string

	// ContentEncoding specifies the encoding used for the blob's content, if any.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
	ContentEncoding string

	// ContentLanguage specifies the language used in the blob's content, if any.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Language
	ContentLanguage string

	// ContentType specifies the MIME type of the blob being written. If not set,
	// it will be inferred from the content using the algorithm described at
	// http://mimesniff.spec.whatwg.org/.
	// Set DisableContentTypeDetection to true to disable the above and force
	// the ContentType to stay empty.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	ContentType string

	// When true, if ContentType is the empty string, it will stay the empty
	// string rather than being inferred from the content.
	// Note that while the blob will be written with an empty string ContentType,
	// most providers will fill one in during reads, so don't expect an empty
	// ContentType if you read the blob back.
	DisableContentTypeDetection bool

	// ContentMD5 is used as a message integrity check.
	// If len(ContentMD5) > 0, the MD5 hash of the bytes written must match
	// ContentMD5, or Close will return an error without completing the write.
	// https://tools.ietf.org/html/rfc1864
	ContentMD5 []byte

	// Metadata holds key/value strings to be associated with the blob, or nil.
	// Keys may not be empty, and are lowercased before being written.
	// Duplicate case-insensitive keys (e.g., "foo" and "FOO") will result in
	// an error.
	Metadata map[string]string
}

// NewWriter returns a Writer that writes to the blob stored at key.
// A nil WriterOptions is treated the same as the zero value.
//
// If a blob with this key already exists, it will be replaced.
// The blob being written is not guaranteed to be readable until Close
// has been called; until then, any previous blob will still be readable.
// Even after Close is called, newly written blobs are not guaranteed to be
// returned from List; some services are only eventually consistent.
//
// The returned Writer will store ctx for later use in Write and/or Close.
// To abort a write, cancel ctx; otherwise, it must remain open until
// Close is called.
//
// The caller must call Close on the returned Writer, even if the write is
// aborted.
func (b *Bucket) NewWriter(ctx context.Context, key string, opts *WriterOptions) (_ *Writer, err error) {
	if !utf8.ValidString(key) {
		return nil, fmt.Errorf("NewWriter key must be a valid UTF-8 string: %q", key)
	}
	if opts == nil {
		opts = &WriterOptions{}
	}
	dopts := &WriterOptions{
		CacheControl:                opts.CacheControl,
		ContentDisposition:          opts.ContentDisposition,
		ContentEncoding:             opts.ContentEncoding,
		ContentLanguage:             opts.ContentLanguage,
		ContentMD5:                  opts.ContentMD5,
		BufferSize:                  opts.BufferSize,
		MaxConcurrency:              opts.MaxConcurrency,
		DisableContentTypeDetection: opts.DisableContentTypeDetection,
	}

	if len(opts.Metadata) > 0 {
		// Services are inconsistent, but at least some treat keys
		// as case-insensitive. To make the behavior consistent, we
		// force-lowercase them when writing and reading.
		md := make(map[string]string, len(opts.Metadata))
		for k, v := range opts.Metadata {
			if k == "" {
				return nil, errors.New("WriterOptions.Metadata keys may not be empty strings")
			}
			if !utf8.ValidString(k) {
				return nil, fmt.Errorf("WriterOptions.Metadata keys must be valid UTF-8 strings: %q", k)
			}
			if !utf8.ValidString(v) {
				return nil, fmt.Errorf("WriterOptions.Metadata values must be valid UTF-8 strings: %q", v)
			}
			lowerK := strings.ToLower(k)
			if _, found := md[lowerK]; found {
				return nil, fmt.Errorf("WriterOptions.Metadata has a duplicate case-insensitive metadata key: %q", lowerK)
			}
			md[lowerK] = v
		}
		dopts.Metadata = md
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return nil, ErrClosed
	}

	ctx, cancel := context.WithCancel(ctx)

	w := &Writer{
		drv:        b.drv,
		cancel:     cancel,
		key:        key,
		contentMD5: opts.ContentMD5,
		md5hash:    md5.New(),
	}

	if opts.ContentType != "" || opts.DisableContentTypeDetection {
		var ct string
		if opts.ContentType != "" {
			t, p, err := mime.ParseMediaType(opts.ContentType)
			if err != nil {
				cancel()
				return nil, err
			}
			ct = mime.FormatMediaType(t, p)
		}
		dw, err := b.drv.NewTypedWriter(ctx, key, ct, dopts)
		if err != nil {
			cancel()
			return nil, wrapError(b.drv, err, key)
		}
		w.w = dw
	} else {
		// Save the fields needed to called NewTypedWriter later, once we've gotten
		// sniffLen bytes; see the comment on Writer.
		w.ctx = ctx
		w.opts = dopts
		w.buf = bytes.NewBuffer([]byte{})
	}

	_, file, lineno, ok := runtime.Caller(1)
	runtime.SetFinalizer(w, func(w *Writer) {
		if !w.closed {
			var caller string
			if ok {
				caller = fmt.Sprintf(" (%s:%d)", file, lineno)
			}
			log.Printf("A blob.Writer writing to %q was never closed%s", key, caller)
		}
	})

	return w, nil
}

// Copy the blob stored at srcKey to dstKey.
// A nil CopyOptions is treated the same as the zero value.
//
// If the source blob does not exist, Copy returns an error for which
// gcerrors.Code will return gcerrors.NotFound.
//
// If the destination blob already exists, it is overwritten.
func (b *Bucket) Copy(ctx context.Context, dstKey, srcKey string) (err error) {
	if !utf8.ValidString(srcKey) {
		return fmt.Errorf("Copy srcKey must be a valid UTF-8 string: %q", srcKey)
	}

	if !utf8.ValidString(dstKey) {
		return fmt.Errorf("Copy dstKey must be a valid UTF-8 string: %q", dstKey)
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return ErrClosed
	}

	return wrapError(b.drv, b.drv.Copy(ctx, dstKey, srcKey), fmt.Sprintf("%s -> %s", srcKey, dstKey))
}

// Delete deletes the blob stored at key.
//
// If the blob does not exist, Delete returns an error for which
// gcerrors.Code will return gcerrors.NotFound.
func (b *Bucket) Delete(ctx context.Context, key string) (err error) {
	if !utf8.ValidString(key) {
		return fmt.Errorf("Delete key must be a valid UTF-8 string: %q", key)
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return ErrClosed
	}

	return wrapError(b.drv, b.drv.Delete(ctx, key), key)
}

// Close releases any resources used for the bucket.
func (b *Bucket) Close() error {
	b.mu.Lock()
	prev := b.closed
	b.closed = true
	b.mu.Unlock()

	if prev {
		return ErrClosed
	}

	return wrapError(b.drv, b.drv.Close(), "")
}

func wrapError(b Driver, err error, key string) error {
	if err == nil {
		return nil
	}

	// don't wrap or normalize EOF errors since there are many places
	// in the standard library (e.g. io.ReadAll) that rely on checks
	// such as "err == io.EOF" and they will fail
	if errors.Is(err, io.EOF) {
		return err
	}

	err = b.NormalizeError(err)

	if key != "" {
		err = fmt.Errorf("[key: %s] %w", key, err)
	}

	return err
}
