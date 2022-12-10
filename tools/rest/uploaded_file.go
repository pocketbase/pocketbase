package rest

import (
	"net/http"

	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// DefaultMaxMemory defines the default max memory bytes that
// will be used when parsing a form request body.
const DefaultMaxMemory = 32 << 20 // 32mb

// FindUploadedFiles extracts all form files of "key" from a http request
// and returns a slice with filesystem.File instances (if any).
func FindUploadedFiles(r *http.Request, key string) ([]*filesystem.File, error) {
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(DefaultMaxMemory)
		if err != nil {
			return nil, err
		}
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil || len(r.MultipartForm.File[key]) == 0 {
		return nil, http.ErrMissingFile
	}

	result := make([]*filesystem.File, 0, len(r.MultipartForm.File[key]))

	for _, fh := range r.MultipartForm.File[key] {
		file, err := filesystem.NewFileFromMultipart(fh)
		if err != nil {
			return nil, err
		}

		result = append(result, file)
	}

	return result, nil
}
