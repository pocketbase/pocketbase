package router

import "github.com/pocketbase/pocketbase/tools/hook"

type Route[T hook.Resolver] struct {
	excludedMiddlewares map[string]struct{}

	Action      func(e T) error
	Method      string
	Path        string
	Middlewares []*hook.Handler[T]
}

// BindFunc registers one or multiple middleware functions to the current route.
//
// The registered middleware functions are "anonymous" and with default priority,
// aka. executes in the order they were registered.
//
// If you need to specify a named middleware (ex. so that it can be removed)
// or middleware with custom exec prirority, use the [Route.Bind] method.
func (route *Route[T]) BindFunc(middlewareFuncs ...func(e T) error) *Route[T] {
	for _, m := range middlewareFuncs {
		route.Middlewares = append(route.Middlewares, &hook.Handler[T]{Func: m})
	}

	return route
}

// Bind registers one or multiple middleware handlers to the current route.
func (route *Route[T]) Bind(middlewares ...*hook.Handler[T]) *Route[T] {
	route.Middlewares = append(route.Middlewares, middlewares...)

	// unmark the newly added middlewares in case they were previously "excluded"
	if route.excludedMiddlewares != nil {
		for _, m := range middlewares {
			if m.Id != "" {
				delete(route.excludedMiddlewares, m.Id)
			}
		}
	}

	return route
}

// Unbind removes one or more middlewares with the specified id(s) from the current route.
//
// It also adds the removed middleware ids to an exclude list so that they could be skipped from
// the execution chain in case the middleware is registered in a parent group.
//
// Anonymous middlewares are considered non-removable, aka. this method
// does nothing if the middleware id is an empty string.
func (route *Route[T]) Unbind(middlewareIds ...string) *Route[T] {
	for _, middlewareId := range middlewareIds {
		if middlewareId == "" {
			continue
		}

		// remove from the route's middlewares
		for i := len(route.Middlewares) - 1; i >= 0; i-- {
			if route.Middlewares[i].Id == middlewareId {
				route.Middlewares = append(route.Middlewares[:i], route.Middlewares[i+1:]...)
			}
		}

		// add to the exclude list
		if route.excludedMiddlewares == nil {
			route.excludedMiddlewares = map[string]struct{}{}
		}
		route.excludedMiddlewares[middlewareId] = struct{}{}
	}

	return route
}
