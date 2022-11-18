package filesystem_test

import (
	"bytes"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
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
		{"image.png", true},
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
		file              string
		expectError       bool
		expectContentType string
	}{
		{"sub1.txt", true, ""},
		{"test/sub1.txt", false, "application/octet-stream"},
		{"test/sub2.txt", false, "application/octet-stream"},
		{"image.png", false, "image/png"},
	}

	for i, scenario := range scenarios {
		attr, err := fs.Attributes(scenario.file)

		if err == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}

		if err != nil && !scenario.expectError {
			t.Errorf("(%d) Expected nil, got error, %v", i, err)
		}

		if err == nil && attr.ContentType != scenario.expectContentType {
			t.Errorf("(%d) Expected attr.ContentType to be %q, got %q", i, scenario.expectContentType, attr.ContentType)
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

	if err := fs.Delete("image.png"); err != nil {
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

func TestFileSystemUploadMultipart(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	// create multipart form file
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	w, err := mp.CreateFormFile("test", "test")
	if err != nil {
		t.Fatalf("Failed creating form file: %v", err)
	}
	w.Write([]byte("demo"))
	mp.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Add("Content-Type", mp.FormDataContentType())

	file, fh, err := req.FormFile("test")
	if err != nil {
		t.Fatalf("Failed to fetch form file: %v", err)
	}
	defer file.Close()
	// ---

	fs, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()

	uploadErr := fs.UploadMultipart(fh, "newdir/newkey.txt")
	if uploadErr != nil {
		t.Fatal(uploadErr)
	}

	if exists, _ := fs.Exists("newdir/newkey.txt"); !exists {
		t.Fatalf("Expected newdir/newkey.txt to exist")
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

	scenarios := []struct {
		path          string
		name          string
		expectError   bool
		expectHeaders map[string]string
	}{
		{
			// missing
			"missing.txt",
			"test_name.txt",
			true,
			nil,
		},
		{
			// existing regular file
			"test/sub1.txt",
			"test_name.txt",
			false,
			map[string]string{
				"Content-Disposition":     "attachment; filename=test_name.txt",
				"Content-Type":            "application/octet-stream",
				"Content-Length":          "0",
				"Content-Security-Policy": "default-src 'none'; media-src 'self'; style-src 'unsafe-inline'; sandbox",
			},
		},
		{
			// png inline
			"image.png",
			"test_name.png",
			false,
			map[string]string{
				"Content-Disposition":     "inline; filename=test_name.png",
				"Content-Type":            "image/png",
				"Content-Length":          "73",
				"Content-Security-Policy": "default-src 'none'; media-src 'self'; style-src 'unsafe-inline'; sandbox",
			},
		},
		{
			// svg exception
			"image.svg",
			"test_name.svg",
			false,
			map[string]string{
				"Content-Disposition":     "attachment; filename=test_name.svg",
				"Content-Type":            "image/svg+xml",
				"Content-Length":          "0",
				"Content-Security-Policy": "default-src 'none'; media-src 'self'; style-src 'unsafe-inline'; sandbox",
			},
		},
		{
			// css exception
			"style.css",
			"test_name.css",
			false,
			map[string]string{
				"Content-Disposition":     "attachment; filename=test_name.css",
				"Content-Type":            "text/css",
				"Content-Length":          "0",
				"Content-Security-Policy": "default-src 'none'; media-src 'self'; style-src 'unsafe-inline'; sandbox",
			},
		},
	}

	for _, scenario := range scenarios {
		r := httptest.NewRecorder()

		err := fs.Serve(r, scenario.path, scenario.name)
		hasErr := err != nil

		if hasErr != scenario.expectError {
			t.Errorf("(%s) Expected hasError %v, got %v (%v)", scenario.path, scenario.expectError, hasErr, err)
			continue
		}

		if scenario.expectError {
			continue
		}

		result := r.Result()

		for hName, hValue := range scenario.expectHeaders {
			v := result.Header.Get(hName)
			if v != hValue {
				t.Errorf("(%s) Expected value %q for header %q, got %q", scenario.path, hValue, hName, v)
			}
		}

		if v := result.Header.Get("X-Frame-Options"); v != "" {
			t.Errorf("(%s) Expected the X-Frame-Options header to be unset, got %v", scenario.path, v)
		}

		if v := result.Header.Get("Cache-Control"); v == "" {
			t.Errorf("(%s) Expected Cache-Control header to be set, got empty string", scenario.path)
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
		{"image.png", "thumb_file_center", true, false},
		// existing image file - crop top
		{"image.png", "thumb_file_top", false, false},
		// existing image file with existing thumb path = should fail
		{"image.png", "test", true, true},
	}

	for i, scenario := range scenarios {
		err := fs.CreateThumb(scenario.file, scenario.thumb, "100x100")

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

	file3, err := os.OpenFile(filepath.Join(dir, "image.png"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	// tiny 1x1 png
	imgRect := image.Rect(0, 0, 1, 1)
	png.Encode(file3, imgRect)
	file3.Close()
	err2 := os.WriteFile(filepath.Join(dir, "image.png.attrs"), []byte(`{"user.cache_control":"","user.content_disposition":"","user.content_encoding":"","user.content_language":"","user.content_type":"image/png","user.metadata":null}`), 0666)
	if err2 != nil {
		t.Fatal(err2)
	}

	file4, err := os.OpenFile(filepath.Join(dir, "image.svg"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	file4.Close()

	file5, err := os.OpenFile(filepath.Join(dir, "style.css"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	file5.Close()

	return dir
}
