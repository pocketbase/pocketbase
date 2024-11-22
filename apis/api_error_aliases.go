package apis

import "github.com/pocketbase/pocketbase/tools/router"

// ApiError aliases to minimize the breaking changes with earlier versions
// and for consistency with the JSVM binds.
// -------------------------------------------------------------------

// ToApiError wraps err into ApiError instance (if not already).
func ToApiError(err error) *router.ApiError {
	return router.ToApiError(err)
}

// NewApiError is an alias for [router.NewApiError].
func NewApiError(status int, message string, errData any) *router.ApiError {
	return router.NewApiError(status, message, errData)
}

// NewBadRequestError is an alias for [router.NewBadRequestError].
func NewBadRequestError(message string, errData any) *router.ApiError {
	return router.NewBadRequestError(message, errData)
}

// NewNotFoundError is an alias for [router.NewNotFoundError].
func NewNotFoundError(message string, errData any) *router.ApiError {
	return router.NewNotFoundError(message, errData)
}

// NewForbiddenError is an alias for [router.NewForbiddenError].
func NewForbiddenError(message string, errData any) *router.ApiError {
	return router.NewForbiddenError(message, errData)
}

// NewUnauthorizedError is an alias for [router.NewUnauthorizedError].
func NewUnauthorizedError(message string, errData any) *router.ApiError {
	return router.NewUnauthorizedError(message, errData)
}

// NewTooManyRequestsError is an alias for [router.NewTooManyRequestsError].
func NewTooManyRequestsError(message string, errData any) *router.ApiError {
	return router.NewTooManyRequestsError(message, errData)
}

// NewInternalServerError is an alias for [router.NewInternalServerError].
func NewInternalServerError(message string, errData any) *router.ApiError {
	return router.NewInternalServerError(message, errData)
}
