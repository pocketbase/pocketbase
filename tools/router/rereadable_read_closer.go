package router

import (
	"io"
	"os"
)

const defaultMaxMemory = 10 << 20 // 10 MB

var (
	_ io.ReadCloser = (*RereadableReadCloser)(nil)
	_ Rereader      = (*RereadableReadCloser)(nil)
)

// Rereader defines an interface for rewindable readers.
type Rereader interface {
	Reread()
}

// RereadableReadCloser defines a wrapper around a io.ReadCloser reader
// allowing to read the original reader multiple times.
//
// It caches the read data in memory up to MaxMemory and then switches
// to a temporary file on disk to prevent excessive memory allocations.
type RereadableReadCloser struct {
	io.ReadCloser

	MaxMemory int64

	file       *os.File
	memBuf     []byte
	readOffset int64
	writeSize  int64
	sourceEOF  bool
}

// Read implements the standard io.Reader interface.
//
// It reads up to len(b) bytes into b and at at the same time writes
// the read data into an internal bytes buffer or a temporary file.
//
// On EOF the r is "rewinded" to allow reading from r multiple times.
func (r *RereadableReadCloser) Read(b []byte) (int, error) {
	if r.MaxMemory == 0 {
		r.MaxMemory = defaultMaxMemory
	}

	if r.readOffset < r.writeSize {
		limit := r.writeSize - r.readOffset
		if int64(len(b)) > limit {
			b = b[:limit]
		}

		var n int
		var err error
		if r.file != nil {
			n, err = r.file.ReadAt(b, r.readOffset)
			if err == io.EOF {
				err = nil
			}
		} else {
			n = copy(b, r.memBuf[r.readOffset:])
		}

		r.readOffset += int64(n)

		if r.readOffset == r.writeSize && r.sourceEOF {
			r.Reread()
			return n, io.EOF
		}
		return n, err
	}

	if r.sourceEOF {
		r.Reread()
		return 0, io.EOF
	}

	n, err := r.ReadCloser.Read(b)
	if n > 0 {
		if storeErr := r.store(b[:n]); storeErr != nil {
			return n, storeErr
		}
		r.readOffset += int64(n)
	}

	if err == io.EOF {
		r.sourceEOF = true
		r.Reread()
	}

	return n, err
}

func (r *RereadableReadCloser) store(p []byte) error {
	r.writeSize += int64(len(p))

	if r.file != nil {
		_, err := r.file.Write(p)
		return err
	}

	if int64(len(r.memBuf))+int64(len(p)) > r.MaxMemory {
		f, err := os.CreateTemp("", "rereader-")
		if err != nil {
			return err
		}
		r.file = f

		if len(r.memBuf) > 0 {
			if _, err := r.file.Write(r.memBuf); err != nil {
				return err
			}
			r.memBuf = nil
		}

		_, err = r.file.Write(p)
		return err
	}

	r.memBuf = append(r.memBuf, p...)
	return nil
}

// Reread satisfies the [Rereader] interface and resets the r internal state to allow rereads.
//
// note: not named Reset to avoid conflicts with other reader interfaces.
func (r *RereadableReadCloser) Reread() {
	r.readOffset = 0
}

// Close closes the underlying reader and cleans up any temporary files.
func (r *RereadableReadCloser) Close() error {
	var err error
	if r.file != nil {
		name := r.file.Name()
		r.file.Close()
		err = os.Remove(name)
	}

	if r.ReadCloser != nil {
		closeErr := r.ReadCloser.Close()
		if err == nil {
			err = closeErr
		}
	}
	return err
}
