package blob

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
)

// Largely copied from gocloud.dev/blob.Reader to minimize breaking changes.
//
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
// -------------------------------------------------------------------

var _ io.ReadSeekCloser = (*Reader)(nil)

// Reader reads bytes from a blob.
// It implements io.ReadSeekCloser, and must be closed after reads are finished.
type Reader struct {
	ctx            context.Context // Used to recreate r after Seeks
	r              DriverReader
	drv            Driver
	key            string
	baseOffset     int64 // The base offset provided to NewRangeReader.
	baseLength     int64 // The length provided to NewRangeReader (may be negative).
	relativeOffset int64 // Current offset (relative to baseOffset).
	savedOffset    int64 // Last relativeOffset for r, saved after relativeOffset is changed in Seek, or -1 if no Seek.
	closed         bool
}

// Read implements io.Reader (https://golang.org/pkg/io/#Reader).
func (r *Reader) Read(p []byte) (int, error) {
	if r.savedOffset != -1 {
		// We've done one or more Seeks since the last read. We may have
		// to recreate the Reader.
		//
		// Note that remembering the savedOffset and lazily resetting the
		// reader like this allows the caller to Seek, then Seek again back,
		// to the original offset, without having to recreate the reader.
		// We only have to recreate the reader if we actually read after a Seek.
		// This is an important optimization because it's common to Seek
		// to (SeekEnd, 0) and use the return value to determine the size
		// of the data, then Seek back to (SeekStart, 0).
		saved := r.savedOffset
		if r.relativeOffset == saved {
			// Nope! We're at the same place we left off.
			r.savedOffset = -1
		} else {
			// Yep! We've changed the offset. Recreate the reader.
			length := r.baseLength
			if length >= 0 {
				length -= r.relativeOffset
				if length < 0 {
					// Shouldn't happen based on checks in Seek.
					return 0, fmt.Errorf("invalid Seek (base length %d, relative offset %d)", r.baseLength, r.relativeOffset)
				}
			}
			newR, err := r.drv.NewRangeReader(r.ctx, r.key, r.baseOffset+r.relativeOffset, length)
			if err != nil {
				return 0, wrapError(r.drv, err, r.key)
			}
			_ = r.r.Close()
			r.savedOffset = -1
			r.r = newR
		}
	}
	n, err := r.r.Read(p)
	r.relativeOffset += int64(n)
	return n, wrapError(r.drv, err, r.key)
}

// Seek implements io.Seeker (https://golang.org/pkg/io/#Seeker).
func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	if r.savedOffset == -1 {
		// Save the current offset for our reader. If the Seek changes the
		// offset, and then we try to read, we'll need to recreate the reader.
		// See comment above in Read for why we do it lazily.
		r.savedOffset = r.relativeOffset
	}
	// The maximum relative offset is the minimum of:
	// 1. The actual size of the blob, minus our initial baseOffset.
	// 2. The length provided to NewRangeReader (if it was non-negative).
	maxRelativeOffset := r.Size() - r.baseOffset
	if r.baseLength >= 0 && r.baseLength < maxRelativeOffset {
		maxRelativeOffset = r.baseLength
	}
	switch whence {
	case io.SeekStart:
		r.relativeOffset = offset
	case io.SeekCurrent:
		r.relativeOffset += offset
	case io.SeekEnd:
		r.relativeOffset = maxRelativeOffset + offset
	}
	if r.relativeOffset < 0 {
		// "Seeking to an offset before the start of the file is an error."
		invalidOffset := r.relativeOffset
		r.relativeOffset = 0
		return 0, fmt.Errorf("Seek resulted in invalid offset %d, using 0", invalidOffset)
	}
	if r.relativeOffset > maxRelativeOffset {
		// "Seeking to any positive offset is legal, but the behavior of subsequent
		// I/O operations on the underlying object is implementation-dependent."
		// We'll choose to set the offset to the EOF.
		log.Printf("blob.Reader.Seek set an offset after EOF (base offset/length from NewRangeReader %d, %d; actual blob size %d; relative offset %d -> absolute offset %d).", r.baseOffset, r.baseLength, r.Size(), r.relativeOffset, r.baseOffset+r.relativeOffset)
		r.relativeOffset = maxRelativeOffset
	}
	return r.relativeOffset, nil
}

// Close implements io.Closer (https://golang.org/pkg/io/#Closer).
func (r *Reader) Close() error {
	r.closed = true
	err := wrapError(r.drv, r.r.Close(), r.key)
	return err
}

// ContentType returns the MIME type of the blob.
func (r *Reader) ContentType() string {
	return r.r.Attributes().ContentType
}

// ModTime returns the time the blob was last modified.
func (r *Reader) ModTime() time.Time {
	return r.r.Attributes().ModTime
}

// Size returns the size of the blob content in bytes.
func (r *Reader) Size() int64 {
	return r.r.Attributes().Size
}

// WriteTo reads from r and writes to w until there's no more data or
// an error occurs.
// The return value is the number of bytes written to w.
//
// It implements the io.WriterTo interface.
func (r *Reader) WriteTo(w io.Writer) (int64, error) {
	// If the writer has a ReaderFrom method, use it to do the copy.
	// Don't do this for our own *Writer to avoid infinite recursion.
	// Avoids an allocation and a copy.
	switch w.(type) {
	case *Writer:
	default:
		if rf, ok := w.(io.ReaderFrom); ok {
			n, err := rf.ReadFrom(r)
			return n, err
		}
	}

	_, nw, err := readFromWriteTo(r, w)
	return nw, err
}

// readFromWriteTo is a helper for ReadFrom and WriteTo.
// It reads data from r and writes to w, until EOF or a read/write error.
// It returns the number of bytes read from r and the number of bytes
// written to w.
func readFromWriteTo(r io.Reader, w io.Writer) (int64, int64, error) {
	// Note: can't use io.Copy because it will try to use r.WriteTo
	// or w.WriteTo, which is recursive in this context.
	buf := make([]byte, 4096)
	var totalRead, totalWritten int64
	for {
		numRead, rerr := r.Read(buf)
		if numRead > 0 {
			totalRead += int64(numRead)
			numWritten, werr := w.Write(buf[0:numRead])
			totalWritten += int64(numWritten)
			if werr != nil {
				return totalRead, totalWritten, werr
			}
		}
		if rerr == io.EOF {
			// Done!
			return totalRead, totalWritten, nil
		}
		if rerr != nil {
			return totalRead, totalWritten, rerr
		}
	}
}
