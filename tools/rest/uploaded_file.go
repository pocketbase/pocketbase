package rest

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pocketbase/pocketbase/tools/inflector"
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

		originalExt := filepath.Ext(fh.Filename)
		sanitizedExt := extensionInvalidCharsRegex.ReplaceAllString(originalExt, "")

		originalName := strings.TrimSuffix(fh.Filename, originalExt)
		sanitizedName := inflector.Snakecase(originalName)

		if length := len(sanitizedName); length < 3 {
			// the name is too short so we concatenate an additional random part
			sanitizedName += ("_" + security.RandomString(10))
		} else if length > 100 {
			// keep only the first 100 characters (it is multibyte safe after Snakecase)
			sanitizedName = sanitizedName[:100]
		}

		uploadedFilename := fmt.Sprintf(
			"%s_%s%s",
			sanitizedName,
			security.RandomString(10), // ensure that there is always a random part
			sanitizedExt,
		)

		result[i] = &UploadedFile{
			name:   uploadedFilename,
			header: fh,
			bytes:  buf.Bytes(),
		}
	}

	return result, nil
}
