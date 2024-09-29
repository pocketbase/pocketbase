package router

import (
	"bytes"
	"io"
)

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
type RereadableReadCloser struct {
	io.ReadCloser

	copy   *bytes.Buffer
	active io.Reader
}

// Read implements the standard io.Reader interface.
//
// It reads up to len(b) bytes into b and at at the same time writes
// the read data into an internal bytes buffer.
//
// On EOF the r is "rewinded" to allow reading from r multiple times.
func (r *RereadableReadCloser) Read(b []byte) (int, error) {
	if r.active == nil {
		if r.copy == nil {
			r.copy = &bytes.Buffer{}
		}
		r.active = io.TeeReader(r.ReadCloser, r.copy)
	}

	n, err := r.active.Read(b)
	if err == io.EOF {
		r.Reread()
	}

	return n, err
}

// Reread satisfies the [Rereader] interface and resets the r internal state to allow rereads.
//
// note: not named Reset to avoid conflicts with other reader interfaces.
func (r *RereadableReadCloser) Reread() {
	if r.copy == nil || r.copy.Len() == 0 {
		return // nothing to reset or it has been already reset
	}

	oldCopy := r.copy
	r.copy = &bytes.Buffer{}
	r.active = io.TeeReader(oldCopy, r.copy)
}
