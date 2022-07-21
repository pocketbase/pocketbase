package rest

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/pocketbase/pocketbase/tools/security"
)

// DefaultMaxMemory defines the default max memory bytes that
// will be used when parsing a form request body.
const DefaultMaxMemory = 32 << 20 // 32mb

var extensionInvalidCharsRegex = regexp.MustCompile(`[^\w\.\*\-\+\=\#]+`)

// UploadedFile defines a single multipart uploaded file instance.
type UploadedFile struct {
	name   string
	header *multipart.FileHeader
	bytes  []byte
}

// Name returns an assigned unique name to the uploaded file.
func (f *UploadedFile) Name() string {
	return f.name
}

// Header returns the file header that comes with the multipart request.
func (f *UploadedFile) Header() *multipart.FileHeader {
	return f.header
}

// Bytes returns a slice with the file content.
func (f *UploadedFile) Bytes() []byte {
	return f.bytes
}

// FindUploadedFiles extracts all form files of `key` from a http request
// and returns a slice with `UploadedFile` instances (if any).
func FindUploadedFiles(r *http.Request, key string) ([]*UploadedFile, error) {
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(DefaultMaxMemory)
		if err != nil {
			return nil, err
		}
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil || len(r.MultipartForm.File[key]) == 0 {
		return nil, http.ErrMissingFile
	}

	result := make([]*UploadedFile, len(r.MultipartForm.File[key]))

	for i, fh := range r.MultipartForm.File[key] {
		file, err := fh.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return nil, err
		}

		ext := extensionInvalidCharsRegex.ReplaceAllString(filepath.Ext(fh.Filename), "")

		result[i] = &UploadedFile{
			name:   fmt.Sprintf("%s%s", security.RandomString(32), ext),
			header: fh,
			bytes:  buf.Bytes(),
		}
	}

	return result, nil
}
