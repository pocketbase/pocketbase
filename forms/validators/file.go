package validators

import (
	"fmt"
	"image"
	"strconv"
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
			return validation.NewError("validation_file_size_limit", fmt.Sprintf("Maximum allowed file size is %v bytes.", maxBytes))
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

		if len(validTypes) == 0 {
			return validation.NewError("validation_invalid_mime_type", "Unsupported file type.")
		}

		f, err := v.Reader.Open()
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

// UploadedFileDimensions checks whether the validated `rest.UploadedFile`
// image dimensions (width x height) are within the provided dimensions rules.
//
// Example:
//
//	validDimensions := []string{"480x480", "800x1200"}
//
// In the above example, the upload image must either be 480x480 or 800x1200.
//
//	validation.Field(&form.File, validation.By(validators.UploadedFileDimensions(validDimensions)))
func UploadedFileDimensions(validDimensions []string) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(*filesystem.File)
		if v == nil {
			return nil // nothing to validate
		}

		if len(validDimensions) == 0 {
			return validation.NewError("validation_invalid_dimensions", "Unsupported file dimensions (width x height)")
		}

		f, err := v.Reader.Open()
		if err != nil {
			return validation.NewError("validation_invalid_dimensions", "Unsupported file dimensions (width x height)")
		}
		defer f.Close()

		imageConfig, _, err := image.DecodeConfig(f)
		if err != nil {
			return validation.NewError("validation_invalid_dimensions", "Unsupported file dimensions (width x height)")
		}

		// for each registered dimensions constraints, check if the image is valid
		for _, rule := range validDimensions {
			dimensions := strings.Split(rule, "x")
			if len(dimensions) < 2 {
				return validation.NewError("validation_invalid_dimensions", "Unsupported file dimensions (width x height)")
			}
			width, err := strconv.Atoi(dimensions[0])
			if err != nil {
				return validation.NewError("validation_invalid_dimensions", "Unsupported file dimensions (width x height)")
			}
			height, err := strconv.Atoi(dimensions[1])
			if err != nil {
				return validation.NewError("validation_invalid_dimensions", "Unsupported file dimensions (width x height)")
			}
			if width == imageConfig.Width && height == imageConfig.Height {
				return nil
			}
		}

		validDimensionsErr := strings.Join(validDimensions, ", ")
		return validation.NewError("validation_invalid_dimensions", fmt.Sprintf("The following dimensions are allowed: %s", validDimensionsErr))
	}
}
