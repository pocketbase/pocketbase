package validators

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/rest"
)

// UploadedFileSize checks whether the validated `rest.UploadedFile`
// size is no more than the provided maxBytes.
//
// Example:
//	validation.Field(&form.File, validation.By(validators.UploadedFileSize(1000)))
func UploadedFileSize(maxBytes int) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(*rest.UploadedFile)
		if v == nil {
			return nil // nothing to validate
		}

		if binary.Size(v.Bytes()) > maxBytes {
			return validation.NewError("validation_file_size_limit", fmt.Sprintf("Maximum allowed file size is %v bytes.", maxBytes))
		}

		return nil
	}
}

// UploadedFileMimeType checks whether the validated `rest.UploadedFile`
// mimetype is within the provided allowed mime types.
//
// Example:
// 	validMimeTypes := []string{"test/plain","image/jpeg"}
//	validation.Field(&form.File, validation.By(validators.UploadedFileMimeType(validMimeTypes)))
func UploadedFileMimeType(validTypes []string) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(*rest.UploadedFile)
		if v == nil {
			return nil // nothing to validate
		}

		if len(validTypes) == 0 {
			return validation.NewError("validation_invalid_mime_type", "Unsupported file type.")
		}

		filetype := http.DetectContentType(v.Bytes())

		for _, t := range validTypes {
			if t == filetype {
				return nil // valid
			}
		}

		return validation.NewError("validation_invalid_mime_type", fmt.Sprintf(
			"The following mime types are only allowed: %s.",
			strings.Join(validTypes, ","),
		))
	}
}
