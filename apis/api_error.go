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
	return &ApiError{
		rawData: data,
		Data:    safeErrorsData(data),
		Code:    status,
		Message: strings.TrimSpace(inflector.Sentenize(message)),
	}
}

func safeErrorsData(data any) map[string]any {
	switch v := data.(type) {
	case validation.Errors:
		return resolveSafeErrorsData[error](v)
	case map[string]validation.Error:
		return resolveSafeErrorsData[validation.Error](v)
	case map[string]error:
		return resolveSafeErrorsData[error](v)
	case map[string]any:
		return resolveSafeErrorsData[any](v)
	default:
		return map[string]any{} // not nil to ensure that is json serialized as object
	}
}

func resolveSafeErrorsData[T any](data map[string]T) map[string]any {
	result := map[string]any{}

	for name, err := range data {
		if isNestedError(err) {
			result[name] = safeErrorsData(err)
			continue
		}
		result[name] = resolveSafeErrorItem(err)
	}

	return result
}

func isNestedError(err any) bool {
	switch err.(type) {
	case validation.Errors, map[string]validation.Error, map[string]error, map[string]any:
		return true
	}

	return false
}

// resolveSafeErrorItem extracts from each validation error its
// public safe error code and message.
func resolveSafeErrorItem(err any) map[string]string {
	// default public safe error values
	code := "validation_invalid_value"
	msg := "Invalid value."

	// only validation errors are public safe
	if obj, ok := err.(validation.Error); ok {
		code = obj.Code()
		msg = inflector.Sentenize(obj.Error())
	}

	return map[string]string{
		"code":    code,
		"message": msg,
	}
}
