package filesystem_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/thinkonmay/pocketbase/tests"
	"github.com/thinkonmay/pocketbase/tools/filesystem"
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
	if match, err := regexp.Match(normalizedNamePattern, []byte(f.Name)); !match {
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
	if match, err := regexp.Match(normalizedNamePattern, []byte(f.Name)); !match {
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
	if match, err := regexp.Match(originalNamePattern, []byte(f.OriginalName)); !match {
		t.Fatalf("Expected OriginalName to match %v, got %q (%v)", originalNamePattern, f.OriginalName, err)
	}

	normalizedNamePattern := regexp.QuoteMeta("tmpfile_") + `\w+\_\w{10}` + regexp.QuoteMeta(".txt")
	if match, err := regexp.Match(normalizedNamePattern, []byte(f.Name)); !match {
		t.Fatalf("Expected Name to match %v, got %q (%v)", normalizedNamePattern, f.Name, err)
	}

	if f.Size != 4 {
		t.Fatalf("Expected Size %v, got %v", 4, f.Size)
	}

	if _, ok := f.Reader.(*filesystem.MultipartReader); !ok {
		t.Fatalf("Expected Reader to be MultipartReader, got %v", f.Reader)
	}
}

func TestNewFileFromUrlTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "test")
	}))
	defer srv.Close()

	// cancelled context
	{
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		f, err := filesystem.NewFileFromUrl(ctx, srv.URL+"/cancel")
		if err == nil {
			t.Fatal("[ctx_cancel] Expected error, got nil")
		}
		if f != nil {
			t.Fatalf("[ctx_cancel] Expected file to be nil, got %v", f)
		}
	}

	// error response
	{
		f, err := filesystem.NewFileFromUrl(context.Background(), srv.URL+"/error")
		if err == nil {
			t.Fatal("[error_status] Expected error, got nil")
		}
		if f != nil {
			t.Fatalf("[error_status] Expected file to be nil, got %v", f)
		}
	}

	// valid response
	{
		originalName := "image_! noext"
		normalizedNamePattern := regexp.QuoteMeta("image_noext_") + `\w{10}` + regexp.QuoteMeta(".txt")

		f, err := filesystem.NewFileFromUrl(context.Background(), srv.URL+"/"+originalName)
		if err != nil {
			t.Fatalf("[valid] Unexpected error %v", err)
		}
		if f == nil {
			t.Fatal("[valid] Expected non-nil file")
		}

		// check the created file fields
		if f.OriginalName != originalName {
			t.Fatalf("Expected OriginalName %q, got %q", originalName, f.OriginalName)
		}
		if match, err := regexp.Match(normalizedNamePattern, []byte(f.Name)); !match {
			t.Fatalf("Expected Name to match %v, got %q (%v)", normalizedNamePattern, f.Name, err)
		}
		if f.Size != 4 {
			t.Fatalf("Expected Size %v, got %v", 4, f.Size)
		}
		if _, ok := f.Reader.(*filesystem.BytesReader); !ok {
			t.Fatalf("Expected Reader to be BytesReader, got %v", f.Reader)
		}
	}
}

func TestFileNameNormalizations(t *testing.T) {
	scenarios := []struct {
		name    string
		pattern string
	}{
		{"", `^\w{10}_\w{10}\.txt$`},
		{".png", `^\w{10}_\w{10}\.png$`},
		{".tar.gz", `^\w{10}_\w{10}\.tar\.gz$`},
		{"a.tar.gz", `^a\w{10}_\w{10}\.tar\.gz$`},
		{"a.b.c.d.tar.gz", `^a_b_c_d_\w{10}\.tar\.gz$`},
		{"abcd", `^abcd_\w{10}\.txt$`},
		{"a  b! c d  . 456", `^a_b_c_d_\w{10}\.456$`},                                        // normalize spaces
		{strings.Repeat("a", 101) + "." + strings.Repeat("b", 21), `^a{100}_\w{10}\.b{20}$`}, // name and extension length trim
	}

	for i, s := range scenarios {
		t.Run(strconv.Itoa(i)+"_"+s.name, func(t *testing.T) {
			f, err := filesystem.NewFileFromBytes([]byte("abc"), s.name)
			if err != nil {
				t.Fatal(err)
			}
			if match, err := regexp.Match(s.pattern, []byte(f.Name)); !match {
				t.Fatalf("Expected Name to match %v, got %q (%v)", s.pattern, f.Name, err)
			}
		})
	}
}
