package filesystem_test

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func TestFileSystemExists(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	scenarios := []struct {
		file   string
		exists bool
	}{
		{"sub1.txt", false},
		{"test/sub1.txt", true},
		{"test/sub2.txt", true},
		{"image.png", true},
	}

	for _, s := range scenarios {
		t.Run(s.file, func(t *testing.T) {
			exists, err := fsys.Exists(s.file)

			if err != nil {
				t.Fatal(err)
			}

			if exists != s.exists {
				t.Fatalf("Expected exists %v, got %v", s.exists, exists)
			}
		})
	}
}

func TestFileSystemAttributes(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

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

	for _, s := range scenarios {
		t.Run(s.file, func(t *testing.T) {
			attr, err := fsys.Attributes(s.file)

			hasErr := err != nil

			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectError, hasErr)
			}

			if hasErr && !errors.Is(err, filesystem.ErrNotFound) {
				t.Fatalf("Expected ErrNotFound err, got %q", err)
			}

			if !hasErr && attr.ContentType != s.expectContentType {
				t.Fatalf("Expected attr.ContentType to be %q, got %q", s.expectContentType, attr.ContentType)
			}
		})
	}
}

func TestFileSystemDelete(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	if err := fsys.Delete("missing.txt"); err == nil || !errors.Is(err, filesystem.ErrNotFound) {
		t.Fatalf("Expected ErrNotFound error, got %v", err)
	}

	if err := fsys.Delete("image.png"); err != nil {
		t.Fatalf("Expected nil, got error %v", err)
	}
}

func TestFileSystemDeletePrefixWithoutTrailingSlash(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	if errs := fsys.DeletePrefix(""); len(errs) == 0 {
		t.Fatal("Expected error, got nil", errs)
	}

	if errs := fsys.DeletePrefix("missing"); len(errs) != 0 {
		t.Fatalf("Not existing prefix shouldn't error, got %v", errs)
	}

	if errs := fsys.DeletePrefix("test"); len(errs) != 0 {
		t.Fatalf("Expected nil, got errors %v", errs)
	}

	// ensure that the test/* files are deleted
	if exists, _ := fsys.Exists("test/sub1.txt"); exists {
		t.Fatalf("Expected test/sub1.txt to be deleted")
	}
	if exists, _ := fsys.Exists("test/sub2.txt"); exists {
		t.Fatalf("Expected test/sub2.txt to be deleted")
	}

	// the test directory should remain since the prefix didn't have trailing slash
	if _, err := os.Stat(filepath.Join(dir, "test")); os.IsNotExist(err) {
		t.Fatal("Expected the prefix dir to remain")
	}
}

func TestFileSystemDeletePrefixWithTrailingSlash(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	if errs := fsys.DeletePrefix("missing/"); len(errs) != 0 {
		t.Fatalf("Not existing prefix shouldn't error, got %v", errs)
	}

	if errs := fsys.DeletePrefix("test/"); len(errs) != 0 {
		t.Fatalf("Expected nil, got errors %v", errs)
	}

	// ensure that the test/* files are deleted
	if exists, _ := fsys.Exists("test/sub1.txt"); exists {
		t.Fatalf("Expected test/sub1.txt to be deleted")
	}
	if exists, _ := fsys.Exists("test/sub2.txt"); exists {
		t.Fatalf("Expected test/sub2.txt to be deleted")
	}

	// the test directory should be also deleted since the prefix has trailing slash
	if _, err := os.Stat(filepath.Join(dir, "test")); !os.IsNotExist(err) {
		t.Fatal("Expected the prefix dir to be deleted")
	}
}

func TestFileSystemIsEmptyDir(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	scenarios := []struct {
		dir      string
		expected bool
	}{
		{"", false}, // special case that shouldn't be suffixed with delimiter to search for any files within the bucket
		{"/", true},
		{"missing", true},
		{"missing/", true},
		{"test", false},
		{"test/", false},
		{"empty", true},
		{"empty/", true},
	}

	for _, s := range scenarios {
		t.Run(s.dir, func(t *testing.T) {
			result := fsys.IsEmptyDir(s.dir)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
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

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	fileKey := "newdir/newkey.txt"

	uploadErr := fsys.UploadMultipart(fh, fileKey)
	if uploadErr != nil {
		t.Fatal(uploadErr)
	}

	if exists, _ := fsys.Exists(fileKey); !exists {
		t.Fatalf("Expected %q to exist", fileKey)
	}

	attrs, err := fsys.Attributes(fileKey)
	if err != nil {
		t.Fatalf("Failed to fetch file attributes: %v", err)
	}
	if name, ok := attrs.Metadata["original-filename"]; !ok || name != "test" {
		t.Fatalf("Expected original-filename to be %q, got %q", "test", name)
	}
}

func TestFileSystemUploadFile(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	fileKey := "newdir/newkey.txt"

	file, err := filesystem.NewFileFromPath(filepath.Join(dir, "image.svg"))
	if err != nil {
		t.Fatalf("Failed to load test file: %v", err)
	}

	file.OriginalName = "test.txt"

	uploadErr := fsys.UploadFile(file, fileKey)
	if uploadErr != nil {
		t.Fatal(uploadErr)
	}

	if exists, _ := fsys.Exists(fileKey); !exists {
		t.Fatalf("Expected %q to exist", fileKey)
	}

	attrs, err := fsys.Attributes(fileKey)
	if err != nil {
		t.Fatalf("Failed to fetch file attributes: %v", err)
	}
	if name, ok := attrs.Metadata["original-filename"]; !ok || name != file.OriginalName {
		t.Fatalf("Expected original-filename to be %q, got %q", file.OriginalName, name)
	}
}

func TestFileSystemUpload(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	fileKey := "newdir/newkey.txt"

	uploadErr := fsys.Upload([]byte("demo"), fileKey)
	if uploadErr != nil {
		t.Fatal(uploadErr)
	}

	if exists, _ := fsys.Exists(fileKey); !exists {
		t.Fatalf("Expected %s to exist", fileKey)
	}
}

func TestFileSystemServe(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	csp := "default-src 'none'; media-src 'self'; style-src 'unsafe-inline'; sandbox"
	cacheControl := "max-age=2592000, stale-while-revalidate=86400"

	scenarios := []struct {
		path          string
		name          string
		query         map[string]string
		headers       map[string]string
		expectError   bool
		expectHeaders map[string]string
	}{
		{
			// missing
			"missing.txt",
			"test_name.txt",
			nil,
			nil,
			true,
			nil,
		},
		{
			// existing regular file
			"test/sub1.txt",
			"test_name.txt",
			nil,
			nil,
			false,
			map[string]string{
				"Content-Disposition":     "attachment; filename=test_name.txt",
				"Content-Type":            "application/octet-stream",
				"Content-Length":          "4",
				"Content-Security-Policy": csp,
				"Cache-Control":           cacheControl,
			},
		},
		{
			// png inline
			"image.png",
			"test_name.png",
			nil,
			nil,
			false,
			map[string]string{
				"Content-Disposition":     "inline; filename=test_name.png",
				"Content-Type":            "image/png",
				"Content-Length":          "73",
				"Content-Security-Policy": csp,
				"Cache-Control":           cacheControl,
			},
		},
		{
			// png with forced attachment
			"image.png",
			"test_name_download.png",
			map[string]string{"download": "1"},
			nil,
			false,
			map[string]string{
				"Content-Disposition":     "attachment; filename=test_name_download.png",
				"Content-Type":            "image/png",
				"Content-Length":          "73",
				"Content-Security-Policy": csp,
				"Cache-Control":           cacheControl,
			},
		},
		{
			// svg exception
			"image.svg",
			"test_name.svg",
			nil,
			nil,
			false,
			map[string]string{
				"Content-Disposition":     "attachment; filename=test_name.svg",
				"Content-Type":            "image/svg+xml",
				"Content-Length":          "0",
				"Content-Security-Policy": csp,
				"Cache-Control":           cacheControl,
			},
		},
		{
			// css exception
			"style.css",
			"test_name.css",
			nil,
			nil,
			false,
			map[string]string{
				"Content-Disposition":     "attachment; filename=test_name.css",
				"Content-Type":            "text/css",
				"Content-Length":          "0",
				"Content-Security-Policy": csp,
				"Cache-Control":           cacheControl,
			},
		},
		{
			// js exception
			"main.js",
			"test_name.js",
			nil,
			nil,
			false,
			map[string]string{
				"Content-Disposition":     "attachment; filename=test_name.js",
				"Content-Type":            "text/javascript",
				"Content-Length":          "0",
				"Content-Security-Policy": csp,
				"Cache-Control":           cacheControl,
			},
		},
		{
			// mjs exception
			"main.mjs",
			"test_name.mjs",
			nil,
			nil,
			false,
			map[string]string{
				"Content-Disposition":     "attachment; filename=test_name.mjs",
				"Content-Type":            "text/javascript",
				"Content-Length":          "0",
				"Content-Security-Policy": csp,
				"Cache-Control":           cacheControl,
			},
		},
		{
			// custom header
			"test/sub2.txt",
			"test_name.txt",
			nil,
			map[string]string{
				"Content-Disposition":     "1",
				"Content-Type":            "2",
				"Content-Length":          "1",
				"Content-Security-Policy": "4",
				"Cache-Control":           "5",
				"X-Custom":                "6",
			},
			false,
			map[string]string{
				"Content-Disposition":     "1",
				"Content-Type":            "2",
				"Content-Length":          "4", // overwriten by http.ServeContent
				"Content-Security-Policy": "4",
				"Cache-Control":           "5",
				"X-Custom":                "6",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.path, func(t *testing.T) {
			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			query := req.URL.Query()
			for k, v := range s.query {
				query.Set(k, v)
			}
			req.URL.RawQuery = query.Encode()

			for k, v := range s.headers {
				res.Header().Set(k, v)
			}

			err := fsys.Serve(res, req, s.path, s.name)
			hasErr := err != nil

			if hasErr != s.expectError {
				t.Fatalf("Expected hasError %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if s.expectError {
				return
			}

			result := res.Result()
			defer result.Body.Close()

			for hName, hValue := range s.expectHeaders {
				v := result.Header.Get(hName)
				if v != hValue {
					t.Errorf("Expected value %q for header %q, got %q", hValue, hName, v)
				}
			}
		})
	}
}

func TestFileSystemGetFile(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	scenarios := []struct {
		file            string
		expectError     bool
		expectedContent string
	}{
		{"test/missing.txt", true, ""},
		{"test/sub1.txt", false, "sub1"},
	}

	for _, s := range scenarios {
		t.Run(s.file, func(t *testing.T) {
			f, err := fsys.GetFile(s.file)
			defer func() {
				if f != nil {
					f.Close()
				}
			}()

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectError, hasErr)
			}

			if hasErr {
				if !errors.Is(err, filesystem.ErrNotFound) {
					t.Fatalf("Expected ErrNotFound error, got %v", err)
				}
				return
			}

			raw, err := io.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			if str := string(raw); str != s.expectedContent {
				t.Fatalf("Expected content\n%s\ngot\n%s", s.expectedContent, str)
			}
		})
	}
}

func TestFileSystemCopy(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	src := "image.png"
	dst := "image.png_copy"

	// copy missing file
	if err := fsys.Copy(dst, src); err == nil {
		t.Fatalf("Expected to fail copying %q to %q, got nil", dst, src)
	}

	// copy existing file
	if err := fsys.Copy(src, dst); err != nil {
		t.Fatalf("Failed to copy %q to %q: %v", src, dst, err)
	}
	f, err := fsys.GetFile(dst)
	//nolint
	defer f.Close()
	if err != nil {
		t.Fatalf("Missing copied file %q: %v", dst, err)
	}
	if f.Size() != 73 {
		t.Fatalf("Expected file size %d, got %d", 73, f.Size())
	}
}

func TestFileSystemList(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	scenarios := []struct {
		prefix   string
		expected []string
	}{
		{
			"",
			[]string{
				"image.png",
				"image.jpg",
				"image.svg",
				"image.webp",
				"image_! noext",
				"style.css",
				"main.js",
				"main.mjs",
				"test/sub1.txt",
				"test/sub2.txt",
			},
		},
		{
			"test",
			[]string{
				"test/sub1.txt",
				"test/sub2.txt",
			},
		},
		{
			"missing",
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run("prefix_"+s.prefix, func(t *testing.T) {
			objs, err := fsys.List(s.prefix)
			if err != nil {
				t.Fatal(err)
			}

			if len(s.expected) != len(objs) {
				t.Fatalf("Expected %d files, got \n%v", len(s.expected), objs)
			}

			for _, obj := range objs {
				var exists bool
				for _, name := range s.expected {
					if name == obj.Key {
						exists = true
						break
					}
				}

				if !exists {
					t.Fatalf("Unexpected file %q", obj.Key)
				}
			}
		})
	}
}

func TestFileSystemServeSingleRange(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Range", "bytes=0-20")

	if err := fsys.Serve(res, req, "image.png", "image.png"); err != nil {
		t.Fatal(err)
	}

	result := res.Result()

	if result.StatusCode != http.StatusPartialContent {
		t.Fatalf("Expected StatusCode %d, got %d", http.StatusPartialContent, result.StatusCode)
	}

	expectedRange := "bytes 0-20/73"
	if cr := result.Header.Get("Content-Range"); cr != expectedRange {
		t.Fatalf("Expected Content-Range %q, got %q", expectedRange, cr)
	}

	if l := result.Header.Get("Content-Length"); l != "21" {
		t.Fatalf("Expected Content-Length %v, got %v", 21, l)
	}
}

func TestFileSystemServeMultiRange(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Range", "bytes=0-20, 25-30")

	if err := fsys.Serve(res, req, "image.png", "image.png"); err != nil {
		t.Fatal(err)
	}

	result := res.Result()

	if result.StatusCode != http.StatusPartialContent {
		t.Fatalf("Expected StatusCode %d, got %d", http.StatusPartialContent, result.StatusCode)
	}

	if ct := result.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/byteranges; boundary=") {
		t.Fatalf("Expected Content-Type to be multipart/byteranges, got %v", ct)
	}
}

func TestFileSystemCreateThumb(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	fsys, err := filesystem.NewLocal(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer fsys.Close()

	scenarios := []struct {
		file             string
		thumb            string
		size             string
		expectedMimeType string
	}{
		// missing
		{"missing.txt", "thumb_test_missing", "100x100", ""},
		// non-image existing file
		{"test/sub1.txt", "thumb_test_sub1", "100x100", ""},
		// existing image file with existing thumb path = should fail
		{"image.png", "test", "100x100", ""},
		// existing image file with invalid thumb size
		{"image.png", "thumb0", "invalid", ""},
		// existing image file with 0xH thumb size
		{"image.png", "thumb_0xH", "0x100", "image/png"},
		// existing image file with Wx0 thumb size
		{"image.png", "thumb_Wx0", "100x0", "image/png"},
		// existing image file with WxH thumb size
		{"image.png", "thumb_WxH", "100x100", "image/png"},
		// existing image file with WxHt thumb size
		{"image.png", "thumb_WxHt", "100x100t", "image/png"},
		// existing image file with WxHb thumb size
		{"image.png", "thumb_WxHb", "100x100b", "image/png"},
		// existing image file with WxHf thumb size
		{"image.png", "thumb_WxHf", "100x100f", "image/png"},
		// jpg
		{"image.jpg", "thumb.jpg", "100x100", "image/jpeg"},
		// webp (should produce png)
		{"image.webp", "thumb.webp", "100x100", "image/png"},
	}

	for _, s := range scenarios {
		t.Run(s.file+"_"+s.thumb+"_"+s.size, func(t *testing.T) {
			err := fsys.CreateThumb(s.file, s.thumb, s.size)

			expectErr := s.expectedMimeType == ""

			hasErr := err != nil
			if hasErr != expectErr {
				t.Fatalf("Expected hasErr to be %v, got %v (%v)", expectErr, hasErr, err)
			}

			if hasErr {
				return
			}

			f, err := fsys.GetFile(s.thumb)
			if err != nil {
				t.Fatalf("Missing expected thumb %s (%v)", s.thumb, err)
			}
			defer f.Close()

			mt, err := mimetype.DetectReader(f)
			if err != nil {
				t.Fatalf("Failed to detect thumb %s mimetype (%v)", s.thumb, err)
			}

			if mtStr := mt.String(); mtStr != s.expectedMimeType {
				t.Fatalf("Expected thumb %s MimeType %q, got %q", s.thumb, s.expectedMimeType, mtStr)
			}
		})
	}
}

// ---

func createTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp(os.TempDir(), "pb_test")
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(filepath.Join(dir, "empty"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(filepath.Join(dir, "test"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(dir, "test/sub1.txt"), []byte("sub1"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(dir, "test/sub2.txt"), []byte("sub2"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// png
	{

		file, err := os.OpenFile(filepath.Join(dir, "image.png"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		imgRect := image.Rect(0, 0, 1, 1) // tiny 1x1 png
		_ = png.Encode(file, imgRect)
		file.Close()
		err = os.WriteFile(filepath.Join(dir, "image.png.attrs"), []byte(`{"user.cache_control":"","user.content_disposition":"","user.content_encoding":"","user.content_language":"","user.content_type":"image/png","user.metadata":null}`), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// jpg
	{
		file, err := os.OpenFile(filepath.Join(dir, "image.jpg"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		imgRect := image.Rect(0, 0, 1, 1) // tiny 1x1 jpg
		_ = jpeg.Encode(file, imgRect, nil)
		file.Close()
		err = os.WriteFile(filepath.Join(dir, "image.jpg.attrs"), []byte(`{"user.cache_control":"","user.content_disposition":"","user.content_encoding":"","user.content_language":"","user.content_type":"image/jpeg","user.metadata":null}`), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// svg
	{
		file, err := os.OpenFile(filepath.Join(dir, "image.svg"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	// webp
	{
		err := os.WriteFile(filepath.Join(dir, "image.webp"), []byte{
			82, 73, 70, 70, 36, 0, 0, 0, 87, 69, 66, 80, 86, 80, 56, 32,
			24, 0, 0, 0, 48, 1, 0, 157, 1, 42, 1, 0, 1, 0, 2, 0, 52, 37,
			164, 0, 3, 112, 0, 254, 251, 253, 80, 0,
		}, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// no extension and invalid characters
	{
		file, err := os.OpenFile(filepath.Join(dir, "image_! noext"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		_ = png.Encode(file, image.Rect(0, 0, 1, 1)) // tiny 1x1 png
		file.Close()
	}

	// css
	{
		file, err := os.OpenFile(filepath.Join(dir, "style.css"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	// js
	{
		file, err := os.OpenFile(filepath.Join(dir, "main.js"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	// mjs
	{
		file, err := os.OpenFile(filepath.Join(dir, "main.mjs"), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	return dir
}
