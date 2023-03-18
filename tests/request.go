package tests

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

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
func MockMultipartData(data map[string]string, fileFields ...string) (*bytes.Buffer, *multipart.Writer, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()

	// write data fields
	for k, v := range data {
		mp.WriteField(k, v)
	}

	// write file fields
	for _, fileField := range fileFields {
		// create a test temporary file
		err := func() error {
			tmpFile, err := os.CreateTemp(os.TempDir(), "tmpfile-*.txt")
			if err != nil {
				return err
			}

			if _, err := tmpFile.Write([]byte("test")); err != nil {
				return err
			}
			tmpFile.Seek(0, 0)
			defer tmpFile.Close()
			defer os.Remove(tmpFile.Name())

			// stub uploaded file
			w, err := mp.CreateFormFile(fileField, tmpFile.Name())
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
