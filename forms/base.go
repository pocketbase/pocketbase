// Package models implements various services used for request data
// validation and applying changes to existing DB models through the app Dao.
package forms

import "regexp"

// base ID value regex pattern
var idRegex = regexp.MustCompile(`^[^\@\#\$\&\|\.\,\'\"\\\/\s]+$`)

// InterceptorNextFunc is a interceptor handler function.
// Usually used in combination with InterceptorFunc.
type InterceptorNextFunc = func() error

// InterceptorFunc defines a single interceptor function that will execute the provided next func handler.
type InterceptorFunc func(next InterceptorNextFunc) InterceptorNextFunc

// runInterceptors executes the provided list of interceptors.
func runInterceptors(next InterceptorNextFunc, interceptors ...InterceptorFunc) error {
	for i := len(interceptors) - 1; i >= 0; i-- {
		next = interceptors[i](next)
	}

	return next()
}
