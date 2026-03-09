package router

import (
	"errors"
	"io"
)

var (
	_ Rereader      = (*RereadableReadCloser)(nil)
	_ io.ReadCloser = (*RereadableReadCloser)(nil)
)

// Rereader defines an interface for rewindable readers.
type Rereader interface {
	Reread()
}

// RereadableReadCloser defines a wrapper around a [io.ReadCloser] reader
// allowing to read the original reader multiple times.
//
// NB! Make sure to call Close after done working with the reader.
type RereadableReadCloser struct {
	io.ReadCloser

	copy        io.ReadWriteCloser
	closeErrors []error

	// MaxMemory specifies the max size of the in memory copy buffer
	// before switching to read/write from temp disk file.
	//
	// If negative or zero, defaults to [DefaultMaxMemory].
	MaxMemory int64
}

// Read implements the standard [io.Reader] interface.
//
// It reads up to len(p) bytes into p and and at the same time copies
// the read data into an internal buffer (memory + temp file).
//
// On EOF r is "rewinded" to allow reading multiple times.
func (r *RereadableReadCloser) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)

	// copy the read bytes into the internal buffer
	if n > 0 {
		if r.copy == nil {
			r.copy = newBufferWithFile(r.MaxMemory)
		}

		if n, err := r.copy.Write(p[:n]); err != nil {
			return n, err
		}
	}

	// end reached -> reset for the next read
	if err == io.EOF {
		r.Reread()
	}

	return n, err
}

// Reread satisfies the [Rereader] interface and resets the r internal state to allow rereads.
//
// note: not named Reset to avoid conflicts with other reader interfaces.
func (r *RereadableReadCloser) Reread() {
	if r.copy == nil {
		return // nothing to reset
	}

	// eagerly close the old reader to prevent accumulating too much memory or temp files
	if err := r.ReadCloser.Close(); err != nil {
		r.closeErrors = append(r.closeErrors, err)
	}

	r.ReadCloser = r.copy
	r.copy = nil
}

// Close implements the standard [io.Closer] interface by cleaning up related resources.
//
// It is safe to call Close multiple times.
// Once Close is invoked the reader no longer can be used and should be discarded.
func (r *RereadableReadCloser) Close() error {
	if r.copy != nil {
		if err := r.copy.Close(); err != nil {
			r.closeErrors = append(r.closeErrors, err)
		}
	}

	if err := r.ReadCloser.Close(); err != nil {
		r.closeErrors = append(r.closeErrors, err)
	}

	return errors.Join(r.closeErrors...)
}
