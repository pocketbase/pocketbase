package router

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"
)

func TestNewBufferWithFile(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name     string
		limit    int64
		expected int64
	}{
		{"negative limit", -1, DefaultMaxMemory},
		{"zero limit", 0, DefaultMaxMemory},
		{"> 0", 1, 1},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			b := newBufferWithFile(s.limit)

			if b.file != nil {
				t.Fatalf("Expected no file descriptor to be open, got %v", b.file)
			}

			if b.buf == nil {
				t.Fatal("Expected buf to be initialized, got nil")
			}

			if b.memoryLimit != s.expected {
				t.Fatalf("Expected %d limit, got %d", 10, b.memoryLimit)
			}
		})
	}
}

func TestBufferWithFile_WriteReadClose(t *testing.T) {
	t.Parallel()

	b := newBufferWithFile(4)

	t.Run("write under limit", func(t *testing.T) {
		n, err := b.Write([]byte("ab"))
		if err != nil {
			t.Fatal(err)
		}

		if n != 2 {
			t.Fatalf("Expected %d bytes to be written, got %v", 2, n)
		}

		if l := b.buf.Len(); l != 2 {
			t.Fatalf("Expected memory buf lenth %d, got %d", 2, l)
		}

		if b.file != nil {
			t.Fatalf("Expected temp file to remain nil, got %v", b.file)
		}
	})

	t.Run("write under limit (again)", func(t *testing.T) {
		n, err := b.Write([]byte("c"))
		if err != nil {
			t.Fatal(err)
		}

		if n != 1 {
			t.Fatalf("Expected %d bytes to be written, got %v", 1, n)
		}

		if l := b.buf.Len(); l != 3 {
			t.Fatalf("Expected memory buf lenth %d, got %d", 3, l)
		}

		if b.file != nil {
			t.Fatalf("Expected temp file to remain nil, got %v", b.file)
		}
	})

	t.Run("write beyound limit (aka. skip memory buf and write into file)", func(t *testing.T) {
		n, err := b.Write([]byte("de"))
		if err != nil {
			t.Fatal(err)
		}

		if n != 2 {
			t.Fatalf("Expected %d bytes to be written, got %v", 2, n)
		}

		if l := b.buf.Len(); l != 3 {
			t.Fatalf("Expected memory buf lenth to be unchanged (%d), got %d", 3, l)
		}

		if b.file == nil {
			t.Fatal("Expected temp file to be initialized")
		}
	})

	t.Run("read 0 bytes fromm non-empty buffer", func(t *testing.T) {
		temp := []byte{}

		n, err := b.Read(temp)
		if err != nil { // should return nil for consistency with bytes.Buffer
			t.Fatalf("Expected nil, got %v", err)
		}

		if n != 0 {
			t.Fatalf("Expected 0 bytes to be read, got %d (%q)", n, temp)
		}
	})

	t.Run("read under limit", func(t *testing.T) {
		expected := "ab"
		temp := make([]byte, 2)

		n, err := b.Read(temp)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}

		if n != len(temp) {
			t.Fatalf("Expected %d bytes to be read, got %d (%q)", len(temp), n, temp)
		}

		if str := string(temp); str != expected {
			t.Fatalf("Expected to read %q, got %q", expected, str)
		}
	})

	t.Run("read beyound limit", func(t *testing.T) {
		expected := "cde"
		temp := make([]byte, 3)

		n, err := b.Read(temp)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}

		if n != len(temp) {
			t.Fatalf("Expected %d bytes to be read, got %d (%q)", len(temp), n, temp)
		}

		if str := string(temp); str != expected {
			t.Fatalf("Expected to read %q, got %q", expected, str)
		}
	})

	t.Run("read from empty buffers", func(t *testing.T) {
		temp := make([]byte, 3)

		n, err := b.Read(temp)
		if err != io.EOF {
			t.Fatalf("Expected EOF, got %v", err)
		}

		if n != 0 {
			t.Fatalf("Expected 0 bytes to be read, got %d (%q)", n, temp)
		}
	})

	t.Run("close cleanup", func(t *testing.T) {
		if b.file == nil {
			t.Fatal("Expected temp file to be initialized, got nil")
		}

		filename := b.file.Name()

		_, err := os.Stat(filename)
		if err != nil || errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("Expected the temp file to exist and be accessible, got %v", err)
		}

		err = b.Close()
		if err != nil {
			t.Fatal(err)
		}

		info, err := os.Stat(filename)
		if err == nil {
			t.Fatalf("Expected the temp file to be deleted after close, got %v", info)
		}

		if b.buf != nil || b.file != nil {
			t.Fatal("Expected the internal buffers to be nil after close")
		}
	})
}
