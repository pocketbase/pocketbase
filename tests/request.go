package tests

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type TestFile struct {
	Field   string
	IsImage bool
}

// DefaultTestFile - default file used for many tests
var DefaultTestFile = TestFile{
	Field: "test",
}

// MockMultipartData creates a mocked multipart/form-data payload.
//
// Example
//
//	data, mp, err := tests.MockMultipartData(
//		map[string]string{"title": "new"},
//		"file1",
//		"file2",
//		...
//	)
func MockMultipartData(data map[string]string, fileFields ...TestFile) (*bytes.Buffer, *multipart.Writer, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()

	// write data fields
	for k, v := range data {
		mp.WriteField(k, v)
	}

	// write TestFile fields
	for _, fileField := range fileFields {
		// create a test temporary TestFile
		err := func() error {
			var tmpFile *os.File
			var err error

			if fileField.IsImage {
				tmpFile, err = os.CreateTemp(os.TempDir(), "tmpfile-*.jpg")
				if err != nil {
					return err
				}

				// read content of the 100x100.jpg TestFile
				_, currentFile, _, _ := runtime.Caller(0)
				testFilesDir := filepath.Join(path.Dir(currentFile), "files")

				source, err := os.Open(testFilesDir + "/100x100.jpg")
				if err != nil {
					return err
				}
				defer source.Close()

				if _, err := io.Copy(tmpFile, source); err != nil {
					return err
				}
			} else {
				tmpFile, err = os.CreateTemp(os.TempDir(), "tmpfile-*.txt")
				if err != nil {
					return err
				}

				if _, err := tmpFile.Write([]byte("test")); err != nil {
					return err
				}
			}

			tmpFile.Seek(0, 0)
			defer tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			// stub uploaded TestFile
			w, err := mp.CreateFormFile(fileField.Field, tmpFile.Name())
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, tmpFile); err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return nil, nil, err
		}
	}

	return body, mp, nil
}
