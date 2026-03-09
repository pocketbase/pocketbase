package router

import (
	"bytes"
	"errors"
	"io"
	"os"
)

var _ io.ReadWriteCloser = (*bufferWithFile)(nil)

// newBufferWithFile initializes and returns a new bufferWithFile with the specified memoryLimit.
//
// If memoryLimit is negative or zero, defaults to [DefaultMaxMemory].
func newBufferWithFile(memoryLimit int64) *bufferWithFile {
	if memoryLimit <= 0 {
		memoryLimit = DefaultMaxMemory
	}

	return &bufferWithFile{
		buf:         new(bytes.Buffer),
		memoryLimit: memoryLimit,
	}
}

// bufferWithFile is similar to [bytes.Buffer] but after the limit it
// fallbacks to a temporary file to minimize excessive memory usage.
type bufferWithFile struct {
	buf            *bytes.Buffer
	file           *os.File
	memoryLimit    int64
	fileReadOffset int64
}

// Read implements the standard [io.Reader] interface by reading
// up to len(p) bytes into p.
func (b *bufferWithFile) Read(p []byte) (n int, err error) {
	if b.buf == nil {
		return 0, errors.New("[bufferWithFile.Read] not initialized or already closed")
	}

	// eagerly get length because bytes.Buffer may resize and change it
	maxToRead := len(p)

	// read first from the memory buffer
	if b.buf.Len() > 0 {
		n, err = b.buf.Read(p)
		if err != nil && err != io.EOF {
			return n, err
		}
	}

	// continue reading from the file to fill the remaining bytes
	if n < maxToRead && b.file != nil {
		fileN, fileErr := b.file.ReadAt(p[n:maxToRead], b.fileReadOffset)
		b.fileReadOffset += int64(fileN)
		n += fileN
		err = fileErr
	}

	// return EOF if the buffers are empty and nothing has been read
	// (to minimize potential breaking changes and for consistency with the bytes.Buffer rules)
	if n == 0 && maxToRead > 0 && err == nil {
		return 0, io.EOF
	}

	return n, err
}

// Write implements the standard [io.Writer] interface by writing the
// content of p into the buffer.
//
// If the current memory buffer doesn't have enough space to hold len(p),
// it write p into a temp disk file.
func (b *bufferWithFile) Write(p []byte) (int, error) {
	if b.buf == nil {
		return 0, errors.New("[bufferWithFile.Write] not initialized or already closed")
	}

	// already above the limit -> continue with the file
	if b.file != nil {
		return b.file.Write(p)
	}

	// above limit -> create and write to file
	if int64(b.buf.Len()+len(p)) > b.memoryLimit {
		if b.file == nil {
			var err error
			b.file, err = os.CreateTemp("", "pb_buffer_file_*")
			if err != nil {
				return 0, err
			}
		}

		return b.file.Write(p)
	}

	// write in memory
	return b.buf.Write(p)
}

// Close implements the standard [io.Closer] interface.
//
// It unsets the memory buffer and will cleanup after the fallback
// temporary file (if exists).
//
// It is safe to call Close multiple times.
// Once Close is invoked the buffer no longer can be used and should be discarded.
func (b *bufferWithFile) Close() error {
	if b.file != nil {
		err := errors.Join(
			b.file.Close(),
			os.Remove(b.file.Name()),
		)
		if err != nil {
			return err
		}

		b.file = nil
	}

	b.buf = nil

	return nil
}
