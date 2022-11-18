package validators

import (
	"fmt"
	"strings"

	"github.com/gabriel-vasile/mimetype"
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

		if int(v.Header().Size) > maxBytes {
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

		f, err := v.Header().Open()
		if err != nil {
			return validation.NewError("validation_invalid_mime_type", "Unsupported file type.")
		}
		defer f.Close()

		filetype, err := mimetype.DetectReader(f)
		if err != nil {
			return validation.NewError("validation_invalid_mime_type", "Unsupported file type.")
		}

		for _, t := range validTypes {
			if filetype.Is(t) {
				return nil // valid
			}
		}

		return validation.NewError("validation_invalid_mime_type", fmt.Sprintf(
			"The following mime types are only allowed: %s.",
			strings.Join(validTypes, ","),
		))
	}
}
