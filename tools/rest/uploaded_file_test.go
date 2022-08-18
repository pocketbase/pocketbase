package rest_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/rest"
)

func TestFindUploadedFiles(t *testing.T) {
	// create a test temporary file (with very large prefix to test if it will be truncated)
	tmpFile, err := os.CreateTemp(os.TempDir(), strings.Repeat("a", 150)+"tmpfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpFile.Write([]byte("test")); err != nil {
		t.Fatal(err)
	}
	tmpFile.Seek(0, 0)
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	// ---

	// stub multipart form file body
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	w, err := mp.CreateFormFile("test", tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(w, tmpFile); err != nil {
		t.Fatal(err)
	}
	mp.Close()
	// ---

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Add("Content-Type", mp.FormDataContentType())

	result, err := rest.FindUploadedFiles(req, "test")
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(result))
	}

	if result[0].Header().Size != 4 {
		t.Fatalf("Expected the file size to be 4 bytes, got %d", result[0].Header().Size)
	}

	if !strings.HasSuffix(result[0].Name(), ".txt") {
		t.Fatalf("Expected the file name to have suffix .txt, got %v", result[0].Name())
	}

	if length := len(result[0].Name()); length != 115 { // truncated + random part + ext
		t.Fatalf("Expected the file name to have length of 115, got %d\n%q", length, result[0].Name())
	}

	if string(result[0].Bytes()) != "test" {
		t.Fatalf("Expected the file content to be %q, got %q", "test", string(result[0].Bytes()))
	}
}

func TestFindUploadedFilesMissing(t *testing.T) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	mp.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Add("Content-Type", mp.FormDataContentType())

	result, err := rest.FindUploadedFiles(req, "test")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got %v", result)
	}
}
