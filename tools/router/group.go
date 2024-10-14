package router

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/pocketbase/pocketbase/tools/hook"
)

// (note: the struct is named RouterGroup instead of Group so that it can
// be embedded in the Router without conflicting with the Group method)

// RouterGroup represents a collection of routes and other sub groups
// that share common pattern prefix and middlewares.
type RouterGroup[T hook.Resolver] struct {
	excludedMiddlewares map[string]struct{}
	children            []any // Route or RouterGroup

	Prefix      string
	Middlewares []*hook.Handler[T]
}

// Group creates and register a new child Group into the current one
// with the specified prefix.
//
// The prefix follows the standard Go net/http ServeMux pattern format ("[HOST]/[PATH]")
// and will be concatenated recursively into the final route path, meaning that
// only the root level group could have HOST as part of the prefix.
//
// Returns the newly created group to allow chaining and registering
// sub-routes and group specific middlewares.
func (group *RouterGroup[T]) Group(prefix string) *RouterGroup[T] {
	newGroup := &RouterGroup[T]{}
	newGroup.Prefix = prefix

	group.children = append(group.children, newGroup)

	return newGroup
}

// BindFunc registers one or multiple middleware functions to the current group.
//
// The registered middleware functions are "anonymous" and with default priority,
// aka. executes in the order they were registered.
//
// If you need to specify a named middleware (ex. so that it can be removed)
// or middleware with custom exec prirority, use [RouterGroup.Bind] method.
func (group *RouterGroup[T]) BindFunc(middlewareFuncs ...func(e T) error) *RouterGroup[T] {
	for _, m := range middlewareFuncs {
		group.Middlewares = append(group.Middlewares, &hook.Handler[T]{Func: m})
	}

	return group
}

// Bind registers one or multiple middleware handlers to the current group.
func (group *RouterGroup[T]) Bind(middlewares ...*hook.Handler[T]) *RouterGroup[T] {
	group.Middlewares = append(group.Middlewares, middlewares...)

	// unmark the newly added middlewares in case they were previously "excluded"
	if group.excludedMiddlewares != nil {
		for _, m := range middlewares {
			if m.Id != "" {
				delete(group.excludedMiddlewares, m.Id)
			}
		}
	}

	return group
}

// Unbind removes one or more middlewares with the specified id(s)
// from the current group and its children (if any).
//
// Anonymous middlewares are not removable, aka. this method does nothing
// if the middleware id is an empty string.
func (group *RouterGroup[T]) Unbind(middlewareIds ...string) *RouterGroup[T] {
	for _, middlewareId := range middlewareIds {
		if middlewareId == "" {
			continue
		}

		// remove from the group middlwares
		for i := len(group.Middlewares) - 1; i >= 0; i-- {
			if group.Middlewares[i].Id == middlewareId {
				group.Middlewares = append(group.Middlewares[:i], group.Middlewares[i+1:]...)
			}
		}

		// remove from the group children
		for i := len(group.children) - 1; i >= 0; i-- {
			switch v := group.children[i].(type) {
			case *RouterGroup[T]:
				v.Unbind(middlewareId)
			case *Route[T]:
				v.Unbind(middlewareId)
			}
		}

		// add to the exclude list
		if group.excludedMiddlewares == nil {
			group.excludedMiddlewares = map[string]struct{}{}
		}
		group.excludedMiddlewares[middlewareId] = struct{}{}
	}

	return group
}

// Route registers a single route into the current group.
//
// Note that the final route path will be the concatenation of all parent groups prefixes + the route path.
// The path follows the standard Go net/http ServeMux format ("[HOST]/[PATH]"),
// meaning that only a top level group route could have HOST as part of the prefix.
//
// Returns the newly created route to allow attaching route-only middlewares.
func (group *RouterGroup[T]) Route(method string, path string, action func(e T) error) *Route[T] {
	route := &Route[T]{
		Method: method,
		Path:   path,
		Action: action,
	}

	group.children = append(group.children, route)

	return route
}

// Any is a shorthand for [RouterGroup.AddRoute] with "" as route method (aka. matches any method).
func (group *RouterGroup[T]) Any(path string, action func(e T) error) *Route[T] {
	return group.Route("", path, action)
}

// GET is a shorthand for [RouterGroup.AddRoute] with GET as route method.
func (group *RouterGroup[T]) GET(path string, action func(e T) error) *Route[T] {
	return group.Route(http.MethodGet, path, action)
}

// SEARCH is a shorthand for [RouterGroup.AddRoute] with SEARCH as route method.
func (group *RouterGroup[T]) SEARCH(path string, action func(e T) error) *Route[T] {
	return group.Route("SEARCH", path, action)
}

// POST is a shorthand for [RouterGroup.AddRoute] with POST as route method.
func (group *RouterGroup[T]) POST(path string, action func(e T) error) *Route[T] {
	return group.Route(http.MethodPost, path, action)
}

// DELETE is a shorthand for [RouterGroup.AddRoute] with DELETE as route method.
func (group *RouterGroup[T]) DELETE(path string, action func(e T) error) *Route[T] {
	return group.Route(http.MethodDelete, path, action)
}

// PATCH is a shorthand for [RouterGroup.AddRoute] with PATCH as route method.
func (group *RouterGroup[T]) PATCH(path string, action func(e T) error) *Route[T] {
	return group.Route(http.MethodPatch, path, action)
}

// PUT is a shorthand for [RouterGroup.AddRoute] with PUT as route method.
func (group *RouterGroup[T]) PUT(path string, action func(e T) error) *Route[T] {
	return group.Route(http.MethodPut, path, action)
}

// HEAD is a shorthand for [RouterGroup.AddRoute] with HEAD as route method.
func (group *RouterGroup[T]) HEAD(path string, action func(e T) error) *Route[T] {
	return group.Route(http.MethodHead, path, action)
}

// OPTIONS is a shorthand for [RouterGroup.AddRoute] with OPTIONS as route method.
func (group *RouterGroup[T]) OPTIONS(path string, action func(e T) error) *Route[T] {
	return group.Route(http.MethodOptions, path, action)
}

// HasRoute checks whether the specified route pattern (method + path)
// is registered in the current group or its children.
//
// This could be useful to conditionally register and checks for routes
// in order prevent panic on duplicated routes.
//
// Note that routes with anonymous and named wildcard placeholder are treated as equal,
// aka. "GET /abc/" is considered the same as "GET /abc/{something...}".
func (group *RouterGroup[T]) HasRoute(method string, path string) bool {
	pattern := path
	if method != "" {
		pattern = strings.ToUpper(method) + " " + pattern
	}

	return group.hasRoute(pattern, nil)
}

func (group *RouterGroup[T]) hasRoute(pattern string, parents []*RouterGroup[T]) bool {
	for _, child := range group.children {
		switch v := child.(type) {
		case *RouterGroup[T]:
			if v.hasRoute(pattern, append(parents, group)) {
				return true
			}
		case *Route[T]:
			var result string

			if v.Method != "" {
				result += v.Method + " "
			}

			// add parent groups prefixes
			for _, p := range parents {
				result += p.Prefix
			}

			// add current group prefix
			result += group.Prefix

			// add current route path
			result += v.Path

			if result == pattern || // direct match
				// compares without the named wildcard, aka. /abc/{test...} is equal to /abc/
				stripWildcard(result) == stripWildcard(pattern) {
				return true
			}
		}
	}
	return false
}

var wildcardPlaceholderRegex = regexp.MustCompile(`/{.+\.\.\.}$`)

func stripWildcard(pattern string) string {
	return wildcardPlaceholderRegex.ReplaceAllString(pattern, "/")
}
