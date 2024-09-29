package router

import (
	"database/sql"
	"errors"
	"io/fs"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/inflector"
)

// SafeErrorItem defines a common error interface for a printable public safe error.
type SafeErrorItem interface {
	// Code represents a fixed unique identifier of the error (usually used as translation key).
	Code() string

	// Error is the default English human readable error message that will be returned.
	Error() string
}

// SafeErrorParamsResolver defines an optional interface for specifying dynamic error parameters.
type SafeErrorParamsResolver interface {
	// Params defines a map with dynamic parameters to return as part of the public safe error view.
	Params() map[string]any
}

// SafeErrorResolver defines an error interface for resolving the public safe error fields.
type SafeErrorResolver interface {
	// Resolve allows modifying and returning a new public safe error data map.
	Resolve(errData map[string]any) any
}

// ApiError defines the struct for a basic api error response.
type ApiError struct {
	rawData any

	Data    map[string]any `json:"data"`
	Message string         `json:"message"`
	Status  int            `json:"status"`
}

// Error makes it compatible with the `error` interface.
func (e *ApiError) Error() string {
	return e.Message
}

// RawData returns the unformatted error data (could be an internal error, text, etc.)
func (e *ApiError) RawData() any {
	return e.rawData
}

// Is reports whether the current ApiError wraps the target.
func (e *ApiError) Is(target error) bool {
	err, ok := e.rawData.(error)
	if ok {
		return errors.Is(err, target)
	}

	apiErr, ok := target.(*ApiError)

	return ok && e == apiErr
}

// NewNotFoundError creates and returns 404 ApiError.
func NewNotFoundError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "The requested resource wasn't found."
	}

	return NewApiError(http.StatusNotFound, message, rawErrData)
}

// NewBadRequestError creates and returns 400 ApiError.
func NewBadRequestError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "Something went wrong while processing your request."
	}

	return NewApiError(http.StatusBadRequest, message, rawErrData)
}

// NewForbiddenError creates and returns 403 ApiError.
func NewForbiddenError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "You are not allowed to perform this request."
	}

	return NewApiError(http.StatusForbidden, message, rawErrData)
}

// NewUnauthorizedError creates and returns 401 ApiError.
func NewUnauthorizedError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "Missing or invalid authentication."
	}

	return NewApiError(http.StatusUnauthorized, message, rawErrData)
}

// NewInternalServerError creates and returns 500 ApiError.
func NewInternalServerError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "Something went wrong while processing your request."
	}

	return NewApiError(http.StatusInternalServerError, message, rawErrData)
}

func NewTooManyRequestsError(message string, rawErrData any) *ApiError {
	if message == "" {
		message = "Too Many Requests."
	}

	return NewApiError(http.StatusTooManyRequests, message, rawErrData)
}

// NewApiError creates and returns new normalized ApiError instance.
func NewApiError(status int, message string, rawErrData any) *ApiError {
	if message == "" {
		message = http.StatusText(status)
	}

	return &ApiError{
		rawData: rawErrData,
		Data:    safeErrorsData(rawErrData),
		Status:  status,
		Message: strings.TrimSpace(inflector.Sentenize(message)),
	}
}

// ToApiError wraps err into ApiError instance (if not already).
func ToApiError(err error) *ApiError {
	var apiErr *ApiError

	if !errors.As(err, &apiErr) {
		// no ApiError found -> assign a generic one
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, fs.ErrNotExist) {
			apiErr = NewNotFoundError("", err)
		} else {
			apiErr = NewBadRequestError("", err)
		}
	}

	return apiErr
}

// -------------------------------------------------------------------

func safeErrorsData(data any) map[string]any {
	switch v := data.(type) {
	case validation.Errors:
		return resolveSafeErrorsData(v)
	case error:
		validationErrors := validation.Errors{}
		if errors.As(v, &validationErrors) {
			return resolveSafeErrorsData(validationErrors)
		}
		return map[string]any{} // not nil to ensure that is json serialized as object
	case map[string]validation.Error:
		return resolveSafeErrorsData(v)
	case map[string]SafeErrorItem:
		return resolveSafeErrorsData(v)
	case map[string]error:
		return resolveSafeErrorsData(v)
	case map[string]string:
		return resolveSafeErrorsData(v)
	case map[string]any:
		return resolveSafeErrorsData(v)
	default:
		return map[string]any{} // not nil to ensure that is json serialized as object
	}
}

func resolveSafeErrorsData[T any](data map[string]T) map[string]any {
	result := map[string]any{}

	for name, err := range data {
		if isNestedError(err) {
			result[name] = safeErrorsData(err)
		} else {
			result[name] = resolveSafeErrorItem(err)
		}
	}

	return result
}

func isNestedError(err any) bool {
	switch err.(type) {
	case validation.Errors,
		map[string]validation.Error,
		map[string]SafeErrorItem,
		map[string]error,
		map[string]string,
		map[string]any:
		return true
	}

	return false
}

// resolveSafeErrorItem extracts from each validation error its
// public safe error code and message.
func resolveSafeErrorItem(err any) any {
	data := map[string]any{}

	if obj, ok := err.(SafeErrorItem); ok {
		// extract the specific error code and message
		data["code"] = obj.Code()
		data["message"] = inflector.Sentenize(obj.Error())
	} else {
		// fallback to the default public safe values
		data["code"] = "validation_invalid_value"
		data["message"] = "Invalid value."
	}

	if s, ok := err.(SafeErrorParamsResolver); ok {
		params := s.Params()
		if len(params) > 0 {
			data["params"] = params
		}
	}

	if s, ok := err.(SafeErrorResolver); ok {
		return s.Resolve(data)
	}

	return data
}
