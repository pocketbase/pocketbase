// Package models implements various services used for request data
// validation and applying changes to existing DB models through the app Dao.
package forms

import (
	"regexp"
)

// base ID value regex pattern
var idRegex = regexp.MustCompile(`^[^\@\#\$\&\|\.\,\'\"\\\/\s]+$`)

// InterceptorNextFunc is a interceptor handler function.
// Usually used in combination with InterceptorFunc.
type InterceptorNextFunc[T any] func(t T) error

// InterceptorFunc defines a single interceptor function that
// will execute the provided next func handler.
type InterceptorFunc[T any] func(next InterceptorNextFunc[T]) InterceptorNextFunc[T]

// runInterceptors executes the provided list of interceptors.
func runInterceptors[T any](
	data T,
	next InterceptorNextFunc[T],
	interceptors ...InterceptorFunc[T],
) error {
	for i := len(interceptors) - 1; i >= 0; i-- {
		next = interceptors[i](next)
	}

	return next(data)
}
