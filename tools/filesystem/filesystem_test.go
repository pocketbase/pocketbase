package filesystem_test

import (
	"image"
	"image/png"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func TestFileSystemExists(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fs, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	scenarios := []struct {
		file   string
		exists bool
	}{
		{"sub1.txt", false},
		{"test/sub1.txt", true},
		{"test/sub2.txt", true},
		{"file.png", true},
	}

	for i, scenario := range scenarios {
		exists, _ := fs.Exists(scenario.file)

		if exists != scenario.exists {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.exists, exists)
		}
	}
}

func TestFileSystemAttributes(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fs, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	scenarios := []struct {
		file        string
		expectError bool
	}{
		{"sub1.txt", true},
		{"test/sub1.txt", false},
		{"test/sub2.txt", false},
		{"file.png", false},
	}

	for i, scenario := range scenarios {
		attr, err := fs.Attributes(scenario.file)

		if err == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}

		if err != nil && !scenario.expectError {
			t.Errorf("(%d) Expected nil, got error, %v", i, err)
		}

		if err == nil && attr.ContentType != "application/octet-stream" {
			t.Errorf("(%d) Expected attr.ContentType to be %q, got %q", i, "application/octet-stream", attr.ContentType)
		}
	}
}

func TestFileSystemDelete(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fs, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	if err := fs.Delete("missing.txt"); err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err := fs.Delete("file.png"); err != nil {
		t.Fatalf("Expected nil, got error %v", err)
	}
}

func TestFileSystemDeletePrefix(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fs, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	if errs := fs.DeletePrefix(""); len(errs) == 0 {
		t.Fatal("Expected error, got nil", errs)
	}

	if errs := fs.DeletePrefix("missing/"); len(errs) != 0 {
		t.Fatalf("Not existing prefix shouldn't error, got %v", errs)
	}

	if errs := fs.DeletePrefix("test"); len(errs) != 0 {
		t.Fatalf("Expected nil, got errors %v", errs)
	}

	// ensure that the test/ files are deleted
	if exists, _ := fs.Exists("test/sub1.txt"); exists {
		t.Fatalf("Expected test/sub1.txt to be deleted")
	}
	if exists, _ := fs.Exists("test/sub2.txt"); exists {
		t.Fatalf("Expected test/sub2.txt to be deleted")
	}
}

func TestFileSystemUpload(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fs, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	uploadErr := fs.Upload([]byte("demo"), "newdir/newkey.txt")
	if uploadErr != nil {
		t.Fatal(uploadErr)
	}

	if exists, _ := fs.Exists("newdir/newkey.txt"); !exists {
		t.Fatalf("Expected newdir/newkey.txt to exist")
	}
}

func TestFileSystemServe(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fs, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	r := httptest.NewRecorder()

	// serve missing file
	if err := fs.Serve(r, "missing.txt", "download.txt"); err == nil {
		t.Fatal("Expected error, got nil")
	}

	// serve existing file
	if err := fs.Serve(r, "test/sub1.txt", "download.txt"); err != nil {
		t.Fatal("Expected nil, got error")
	}

	result := r.Result()

	// check headers
	scenarios := []struct {
		header   string
		expected string
	}{
		{"Content-Disposition", "attachment; filename=download.txt"},
		{"Content-Type", "application/octet-stream"},
		{"Content-Length", "0"},
	}
	for i, scenario := range scenarios {
		v := result.Header.Get(scenario.header)
		if v != scenario.expected {
			t.Errorf("(%d) Expected value %q for header %q, got %q", i, scenario.expected, scenario.header, v)
		}
	}
}

func TestFileSystemCreateThumb(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fs, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	scenarios := []struct {
		file        string
		thumb       string
		cropCenter  bool
		expectError bool
	}{
		// missing
		{"missing.txt", "thumb_test_missing", true, true},
		// non-image existing file
		{"test/sub1.txt", "thumb_test_sub1", true, true},
		// existing image file - crop center
		{"file.png", "thumb_file_center", true, false},
		// existing image file - crop top
		{"file.png", "thumb_file_top", false, false},
		// existing image file with existing thumb path = should fail
		{"file.png", "test", true, true},
	}

	for i, scenario := range scenarios {
		err := fs.CreateThumb(scenario.file, scenario.thumb, "100x100", scenario.cropCenter)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if scenario.expectError {
			continue
		}

		if exists, _ := fs.Exists(scenario.thumb); !exists {
			t.Errorf("(%d) Couldn't find %q thumb", i, scenario.thumb)
		}
	}
}

// ---

func createTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp(os.TempDir(), "pb_test")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	file1, err := os.OpenFile(filepath.Join(dir, "test/sub1.txt"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	file1.Close()

	file2, err := os.OpenFile(filepath.Join(dir, "test/sub2.txt"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	file2.Close()

	file3, err := os.OpenFile(filepath.Join(dir, "file.png"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	// tiny 1x1 png
	imgRect := image.Rect(0, 0, 1, 1)
	png.Encode(file3, imgRect)
	file3.Close()

	return dir
}
