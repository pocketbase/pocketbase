package router

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestRereadableReadCloser(t *testing.T) {
	content := "test"

	rereadable := &RereadableReadCloser{
		ReadCloser: io.NopCloser(strings.NewReader(content)),
		MaxMemory:  2, // should store the rest 2 bytes in temp file
	}

	totalRereads := 5

	tempFilenames := make([]string, 0, totalRereads)

	// reread multiple times
	for i := 0; i < totalRereads; i++ {
		t.Run("run_"+strconv.Itoa(i), func(t *testing.T) {
			if i > 3 {
				// test allso with manual Reread calls to ensure that
				// r.copy is reseted and written to only when there are n>0 bytes
				rereadable.Reread()
			}

			result, err := io.ReadAll(rereadable)
			if err != nil {
				t.Fatalf("[read:%d] %v", i, err)
			}
			if str := string(result); str != content {
				t.Fatalf("[read:%d] Expected %q, got %q", i, content, result)
			}

			b, ok := rereadable.ReadCloser.(*bufferWithFile)
			if !ok {
				t.Fatalf("Expected bufferWithFile replacement, got %v", b)
			}

			if b.file != nil {
				tempFilenames = append(tempFilenames, b.file.Name())
			}
		})
	}

	if v := len(tempFilenames); v != totalRereads {
		t.Fatalf("Expected %d temp files to have been created during the previous rereads, got %d", totalRereads, v)
	}

	err := rereadable.Close()
	if err != nil {
		t.Fatalf("Expected no close errors, got %v", err)
	}

	// ensure that no lingering temp files are left after close
	for _, name := range tempFilenames {
		info, err := os.Stat(name)
		if err == nil || !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("Expected file name %q to be deleted, got %v (%v)", name, info, err)
		}
	}
}
