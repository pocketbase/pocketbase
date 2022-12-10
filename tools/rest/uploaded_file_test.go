package rest_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tools/rest"
)

func TestFindUploadedFiles(t *testing.T) {
	scenarios := []struct {
		filename        string
		expectedPattern string
	}{
		{"ab.png", `^ab\w{10}_\w{10}\.png$`},
		{"test", `^test_\w{10}\.txt$`},
		{"a b c d!@$.j!@$pg", `^a_b_c_d_\w{10}\.jpg$`},
		{strings.Repeat("a", 150), `^a{100}_\w{10}\.txt$`},
	}

	for i, s := range scenarios {
		// create multipart form file body
		body := new(bytes.Buffer)
		mp := multipart.NewWriter(body)
		w, err := mp.CreateFormFile("test", s.filename)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte("test"))
		mp.Close()
		// ---

		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Add("Content-Type", mp.FormDataContentType())

		result, err := rest.FindUploadedFiles(req, "test")
		if err != nil {
			t.Fatal(err)
		}

		if len(result) != 1 {
			t.Errorf("[%d] Expected 1 file, got %d", i, len(result))
		}

		if result[0].Size != 4 {
			t.Errorf("[%d] Expected the file size to be 4 bytes, got %d", i, result[0].Size)
		}

		pattern, err := regexp.Compile(s.expectedPattern)
		if err != nil {
			t.Errorf("[%d] Invalid filename pattern %q: %v", i, s.expectedPattern, err)
		}
		if !pattern.MatchString(result[0].Name) {
			t.Fatalf("Expected filename to match %s, got filename %s", s.expectedPattern, result[0].Name)
		}
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
