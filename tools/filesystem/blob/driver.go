package blob

import (
	"context"
	"io"
	"time"
)

// ReaderAttributes contains a subset of attributes about a blob that are
// accessible from Reader.
type ReaderAttributes struct {
	// ContentType is the MIME type of the blob object. It must not be empty.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	ContentType string `json:"contentType"`

	// ModTime is the time the blob object was last modified.
	ModTime time.Time `json:"modTime"`

	// Size is the size of the object in bytes.
	Size int64 `json:"size"`
}

// DriverReader reads an object from the blob.
type DriverReader interface {
	io.ReadCloser

	// Attributes returns a subset of attributes about the blob.
	// The portable type will not modify the returned ReaderAttributes.
	Attributes() *ReaderAttributes
}

// DriverWriter writes an object to the blob.
type DriverWriter interface {
	io.WriteCloser
}

// Driver provides read, write and delete operations on objects within it on the
// blob service.
type Driver interface {
	NormalizeError(err error) error

	// Attributes returns attributes for the blob. If the specified object does
	// not exist, Attributes must return an error for which ErrorCode returns
	// gcerrors.NotFound.
	// The portable type will not modify the returned Attributes.
	Attributes(ctx context.Context, key string) (*Attributes, error)

	// ListPaged lists objects in the bucket, in lexicographical order by
	// UTF-8-encoded key, returning pages of objects at a time.
	// Services are only required to be eventually consistent with respect
	// to recently written or deleted objects. That is to say, there is no
	// guarantee that an object that's been written will immediately be returned
	// from ListPaged.
	// opts is guaranteed to be non-nil.
	ListPaged(ctx context.Context, opts *ListOptions) (*ListPage, error)

	// NewRangeReader returns a Reader that reads part of an object, reading at
	// most length bytes starting at the given offset. If length is negative, it
	// will read until the end of the object. If the specified object does not
	// exist, NewRangeReader must return an error for which ErrorCode returns
	// gcerrors.NotFound.
	// opts is guaranteed to be non-nil.
	//
	// The returned Reader *may* also implement Downloader if the underlying
	// implementation can take advantage of that. The Download call is guaranteed
	// to be the only call to the Reader. For such readers, offset will always
	// be 0 and length will always be -1.
	NewRangeReader(ctx context.Context, key string, offset, length int64) (DriverReader, error)

	// NewTypedWriter returns Writer that writes to an object associated with key.
	//
	// A new object will be created unless an object with this key already exists.
	// Otherwise any previous object with the same key will be replaced.
	// The object may not be available (and any previous object will remain)
	// until Close has been called.
	//
	// contentType sets the MIME type of the object to be written.
	// opts is guaranteed to be non-nil.
	//
	// The caller must call Close on the returned Writer when done writing.
	//
	// Implementations should abort an ongoing write if ctx is later canceled,
	// and do any necessary cleanup in Close. Close should then return ctx.Err().
	//
	// The returned Writer *may* also implement Uploader if the underlying
	// implementation can take advantage of that. The Upload call is guaranteed
	// to be the only non-Close call to the Writer..
	NewTypedWriter(ctx context.Context, key, contentType string, opts *WriterOptions) (DriverWriter, error)

	// Copy copies the object associated with srcKey to dstKey.
	//
	// If the source object does not exist, Copy must return an error for which
	// ErrorCode returns gcerrors.NotFound.
	//
	// If the destination object already exists, it should be overwritten.
	//
	// opts is guaranteed to be non-nil.
	Copy(ctx context.Context, dstKey, srcKey string) error

	// Delete deletes the object associated with key. If the specified object does
	// not exist, Delete must return an error for which ErrorCode returns
	// gcerrors.NotFound.
	Delete(ctx context.Context, key string) error

	// Close cleans up any resources used by the Bucket. Once Close is called,
	// there will be no method calls to the Bucket other than As, ErrorAs, and
	// ErrorCode. There may be open readers or writers that will receive calls.
	// It is up to the driver as to how these will be handled.
	Close() error
}
