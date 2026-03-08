package router

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestRereadableReadCloser(t *testing.T) {
	content := "test"

	rereadable := &RereadableReadCloser{
		ReadCloser: io.NopCloser(strings.NewReader(content)),
	}

	// read multiple times
	for i := 0; i < 3; i++ {
		result, err := io.ReadAll(rereadable)
		if err != nil {
			t.Fatalf("[read:%d] %v", i, err)
		}
		if str := string(result); str != content {
			t.Fatalf("[read:%d] Expected %q, got %q", i, content, result)
		}
	}

	// verify that no file was created for small payloads
	if rereadable.file != nil {
		t.Fatal("Expected no temp file to be created")
	}
}

func TestRereadableReadCloserFileFallback(t *testing.T) {
	content := "this content is strictly longer than the defined memory limit"

	rereadable := &RereadableReadCloser{
		ReadCloser: io.NopCloser(strings.NewReader(content)),
		MaxMemory:  10,
	}

	// read multiple times
	for i := 0; i < 3; i++ {
		result, err := io.ReadAll(rereadable)
		if err != nil {
			t.Fatalf("[read:%d] %v", i, err)
		}
		if str := string(result); str != content {
			t.Fatalf("[read:%d] Expected %q, got %q", i, content, result)
		}
	}

	// verify that a temp file was created
	if rereadable.file == nil {
		t.Fatal("Expected a temp file to be created")
	}

	fileName := rereadable.file.Name()
	if _, err := os.Stat(fileName); err != nil {
		t.Fatalf("Expected temp file %q to exist on disk, got error: %v", fileName, err)
	}

	// ensure resources and temp files are cleaned up
	if err := rereadable.Close(); err != nil {
		t.Fatalf("Failed to close reader: %v", err)
	}

	// verify that the temp file was deleted
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		t.Fatalf("Expected temp file %q to be deleted, got error: %v", fileName, err)
	}
}
