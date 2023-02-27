package validators

import (
	"fmt"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// UploadedFileSize checks whether the validated `rest.UploadedFile`
// size is no more than the provided maxBytes.
//
// Example:
//
//	validation.Field(&form.File, validation.By(validators.UploadedFileSize(1000)))
func UploadedFileSize(maxBytes int) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(*filesystem.File)
		if v == nil {
			return nil // nothing to validate
		}

		if int(v.Size) > maxBytes {
			return validation.NewError(
				"validation_file_size_limit",
				fmt.Sprintf("Failed to upload %q - the maximum allowed file size is %v bytes.", v.OriginalName, maxBytes),
			)
		}

		return nil
	}
}

// UploadedFileMimeType checks whether the validated `rest.UploadedFile`
// mimetype is within the provided allowed mime types.
//
// Example:
//
//	validMimeTypes := []string{"test/plain","image/jpeg"}
//	validation.Field(&form.File, validation.By(validators.UploadedFileMimeType(validMimeTypes)))
func UploadedFileMimeType(validTypes []string) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(*filesystem.File)
		if v == nil {
			return nil // nothing to validate
		}

		baseErr := validation.NewError(
			"validation_invalid_mime_type",
			fmt.Sprintf("Failed to upload %q due to unsupported file type.", v.OriginalName),
		)

		if len(validTypes) == 0 {
			return baseErr
		}

		f, err := v.Reader.Open()
		if err != nil {
			return baseErr
		}
		defer f.Close()

		filetype, err := mimetype.DetectReader(f)
		if err != nil {
			return baseErr
		}

		for _, t := range validTypes {
			if filetype.Is(t) {
				return nil // valid
			}
		}

		return validation.NewError(
			"validation_invalid_mime_type",
			fmt.Sprintf(
				"%q mime type must be one of: %s.",
				v.Name,
				strings.Join(validTypes, ", "),
			),
		)
	}
}
