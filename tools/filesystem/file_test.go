package filesystem_test

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func TestNewFileFromPath(t *testing.T) {
	testDir := createTestDir(t)
	defer os.RemoveAll(testDir)

	// missing file
	_, err := filesystem.NewFileFromPath("missing")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// existing file
	originalName := "image_! noext"
	normalizedNamePattern := regexp.QuoteMeta("image_noext_") + `\w{10}` + regexp.QuoteMeta(".png")
	f, err := filesystem.NewFileFromPath(filepath.Join(testDir, originalName))
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if f.OriginalName != originalName {
		t.Fatalf("Expected OriginalName %q, got %q", originalName, f.OriginalName)
	}
	if match, _ := regexp.Match(normalizedNamePattern, []byte(f.Name)); !match {
		t.Fatalf("Expected Name to match %v, got %q (%v)", normalizedNamePattern, f.Name, err)
	}
	if f.Size != 73 {
		t.Fatalf("Expected Size %v, got %v", 73, f.Size)
	}
	if _, ok := f.Reader.(*filesystem.PathReader); !ok {
		t.Fatalf("Expected Reader to be PathReader, got %v", f.Reader)
	}
}

func TestNewFileFromBytes(t *testing.T) {
	// nil bytes
	if _, err := filesystem.NewFileFromBytes(nil, "photo.jpg"); err == nil {
		t.Fatal("Expected error, got nil")
	}

	// zero bytes
	if _, err := filesystem.NewFileFromBytes([]byte{}, "photo.jpg"); err == nil {
		t.Fatal("Expected error, got nil")
	}

	originalName := "image_! noext"
	normalizedNamePattern := regexp.QuoteMeta("image_noext_") + `\w{10}` + regexp.QuoteMeta(".txt")
	f, err := filesystem.NewFileFromBytes([]byte("text\n"), originalName)
	if err != nil {
		t.Fatal(err)
	}
	if f.Size != 5 {
		t.Fatalf("Expected Size %v, got %v", 5, f.Size)
	}
	if f.OriginalName != originalName {
		t.Fatalf("Expected OriginalName %q, got %q", originalName, f.OriginalName)
	}
	if match, _ := regexp.Match(normalizedNamePattern, []byte(f.Name)); !match {
		t.Fatalf("Expected Name to match %v, got %q (%v)", normalizedNamePattern, f.Name, err)
	}
}

func TestNewFileFromMultipart(t *testing.T) {
	formData, mp, err := tests.MockMultipartData(nil, "test")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("", "/", formData)
	req.Header.Set(echo.HeaderContentType, mp.FormDataContentType())
	req.ParseMultipartForm(32 << 20)

	_, mh, err := req.FormFile("test")
	if err != nil {
		t.Fatal(err)
	}

	f, err := filesystem.NewFileFromMultipart(mh)
	if err != nil {
		t.Fatal(err)
	}

	originalNamePattern := regexp.QuoteMeta("tmpfile-") + `\w+` + regexp.QuoteMeta(".txt")
	if match, _ := regexp.Match(originalNamePattern, []byte(f.OriginalName)); !match {
		t.Fatalf("Expected OriginalName to match %v, got %q (%v)", originalNamePattern, f.OriginalName, err)
	}

	normalizedNamePattern := regexp.QuoteMeta("tmpfile_") + `\w+\_\w{10}` + regexp.QuoteMeta(".txt")
	if match, _ := regexp.Match(normalizedNamePattern, []byte(f.Name)); !match {
		t.Fatalf("Expected Name to match %v, got %q (%v)", normalizedNamePattern, f.Name, err)
	}

	if f.Size != 4 {
		t.Fatalf("Expected Size %v, got %v", 4, f.Size)
	}

	if _, ok := f.Reader.(*filesystem.MultipartReader); !ok {
		t.Fatalf("Expected Reader to be MultipartReader, got %v", f.Reader)
	}
}
