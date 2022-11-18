package apis

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/inflector"
)

// ApiError defines the struct for a basic api error response.
type ApiError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`

	// stores unformatted error data (could be an internal error, text, etc.)
	rawData any
}

// Error makes it compatible with the `error` interface.
func (e *ApiError) Error() string {
	return e.Message
}

// RawData returns the unformatted error data (could be an internal error, text, etc.)
func (e *ApiError) RawData() any {
	return e.rawData
}

// NewNotFoundError creates and returns 404 `ApiError`.
func NewNotFoundError(message string, data any) *ApiError {
	if message == "" {
		message = "The requested resource wasn't found."
	}

	return NewApiError(http.StatusNotFound, message, data)
}

// NewBadRequestError creates and returns 400 `ApiError`.
func NewBadRequestError(message string, data any) *ApiError {
	if message == "" {
		message = "Something went wrong while processing your request."
	}

	return NewApiError(http.StatusBadRequest, message, data)
}

// NewForbiddenError creates and returns 403 `ApiError`.
func NewForbiddenError(message string, data any) *ApiError {
	if message == "" {
		message = "You are not allowed to perform this request."
	}

	return NewApiError(http.StatusForbidden, message, data)
}

// NewUnauthorizedError creates and returns 401 `ApiError`.
func NewUnauthorizedError(message string, data any) *ApiError {
	if message == "" {
		message = "Missing or invalid authentication token."
	}

	return NewApiError(http.StatusUnauthorized, message, data)
}

// NewApiError creates and returns new normalized `ApiError` instance.
func NewApiError(status int, message string, data any) *ApiError {
	message = inflector.Sentenize(message)

	formattedData := map[string]any{}

	if v, ok := data.(validation.Errors); ok {
		formattedData = resolveValidationErrors(v)
	}

	return &ApiError{
		rawData: data,
		Data:    formattedData,
		Code:    status,
		Message: strings.TrimSpace(message),
	}
}

func resolveValidationErrors(validationErrors validation.Errors) map[string]any {
	result := map[string]any{}

	// extract from each validation error its error code and message.
	for name, err := range validationErrors {
		// check for nested errors
		if nestedErrs, ok := err.(validation.Errors); ok {
			result[name] = resolveValidationErrors(nestedErrs)
			continue
		}

		errCode := "validation_invalid_value" // default
		if errObj, ok := err.(validation.ErrorObject); ok {
			errCode = errObj.Code()
		}

		result[name] = map[string]string{
			"code":    errCode,
			"message": inflector.Sentenize(err.Error()),
		}
	}

	return result
}
