package router

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/pocketbase/pocketbase/tools/hook"
)

type EventCleanupFunc func()

// EventFactoryFunc defines the function responsible for creating a Route specific event
// based on the provided request handler ServeHTTP data.
//
// Optionally return a clean up function that will be invoked right after the route execution.
type EventFactoryFunc[T hook.Resolver] func(w http.ResponseWriter, r *http.Request) (T, EventCleanupFunc)

// Router defines a thin wrapper around the standard Go [http.ServeMux] by
// adding support for routing sub-groups, middlewares and other common utils.
//
// Example:
//
//	r := NewRouter[*MyEvent](eventFactory)
//
//	// middlewares
//	r.BindFunc(m1, m2)
//
//	// routes
//	r.GET("/test", handler1)
//
//	// sub-routers/groups
//	api := r.Group("/api")
//	api.GET("/admins", handler2)
//
//	// generate a http.ServeMux instance based on the router configurations
//	mux, _ := r.BuildMux()
//
//	http.ListenAndServe("localhost:8090", mux)
type Router[T hook.Resolver] struct {
	// @todo consider renaming the type to just Group and replace the embed type
	// with an alias after Go 1.24 adds support for generic type aliases
	*RouterGroup[T]

	eventFactory EventFactoryFunc[T]
}

// NewRouter creates a new Router instance with the provided event factory function.
func NewRouter[T hook.Resolver](eventFactory EventFactoryFunc[T]) *Router[T] {
	return &Router[T]{
		RouterGroup:  &RouterGroup[T]{},
		eventFactory: eventFactory,
	}
}

// BuildMux constructs a new mux [http.Handler] instance from the current router configurations.
func (r *Router[T]) BuildMux() (http.Handler, error) {
	// Note that some of the default std Go handlers like the [http.NotFoundHandler]
	// cannot be currently extended and requires defining a custom "catch-all" route
	// so that the group middlewares could be executed.
	//
	// https://github.com/golang/go/issues/65648
	if !r.HasRoute("", "/") {
		r.Route("", "/", func(e T) error {
			return NewNotFoundError("", nil)
		})
	}

	mux := http.NewServeMux()

	if err := r.loadMux(mux, r.RouterGroup, nil); err != nil {
		return nil, err
	}

	return mux, nil
}

func (r *Router[T]) loadMux(mux *http.ServeMux, group *RouterGroup[T], parents []*RouterGroup[T]) error {
	for _, child := range group.children {
		switch v := child.(type) {
		case *RouterGroup[T]:
			if err := r.loadMux(mux, v, append(parents, group)); err != nil {
				return err
			}
		case *Route[T]:
			routeHook := &hook.Hook[T]{}

			var pattern string

			if v.Method != "" {
				pattern = v.Method + " "
			}

			// add parent groups middlewares
			for _, p := range parents {
				pattern += p.Prefix
				for _, h := range p.Middlewares {
					if _, ok := p.excludedMiddlewares[h.Id]; !ok {
						if _, ok = group.excludedMiddlewares[h.Id]; !ok {
							if _, ok = v.excludedMiddlewares[h.Id]; !ok {
								routeHook.Bind(h)
							}
						}
					}
				}
			}

			// add current groups middlewares
			pattern += group.Prefix
			for _, h := range group.Middlewares {
				if _, ok := group.excludedMiddlewares[h.Id]; !ok {
					if _, ok = v.excludedMiddlewares[h.Id]; !ok {
						routeHook.Bind(h)
					}
				}
			}

			// add current route middlewares
			pattern += v.Path
			for _, h := range v.Middlewares {
				if _, ok := v.excludedMiddlewares[h.Id]; !ok {
					routeHook.Bind(h)
				}
			}

			mux.HandleFunc(pattern, func(resp http.ResponseWriter, req *http.Request) {
				// wrap the response to add write and status tracking
				resp = &ResponseWriter{ResponseWriter: resp}

				// wrap the request body to allow multiple reads
				req.Body = &RereadableReadCloser{ReadCloser: req.Body}

				event, cleanupFunc := r.eventFactory(resp, req)

				// trigger the handler hook chain
				err := routeHook.Trigger(event, v.Action)
				if err != nil {
					ErrorHandler(resp, req, err)
				}

				if cleanupFunc != nil {
					cleanupFunc()
				}
			})
		default:
			return errors.New("invalid Group item type")
		}
	}

	return nil
}

func ErrorHandler(resp http.ResponseWriter, req *http.Request, err error) {
	if err == nil {
		return
	}

	if ok, _ := getWritten(resp); ok {
		return // a response was already written (aka. already handled)
	}

	header := resp.Header()
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", "application/json")
	}

	apiErr := ToApiError(err)

	resp.WriteHeader(apiErr.Status)

	if req.Method != http.MethodHead {
		if jsonErr := json.NewEncoder(resp).Encode(apiErr); jsonErr != nil {
			log.Println(jsonErr) // truly rare case, log to stderr only for dev purposes
		}
	}
}

// -------------------------------------------------------------------

type WriteTracker interface {
	// Written reports whether a write operation has occurred.
	Written() bool
}

type StatusTracker interface {
	// Status reports the written response status code.
	Status() int
}

type flushErrorer interface {
	FlushError() error
}

var (
	_ WriteTracker  = (*ResponseWriter)(nil)
	_ StatusTracker = (*ResponseWriter)(nil)
	_ http.Flusher  = (*ResponseWriter)(nil)
	_ http.Hijacker = (*ResponseWriter)(nil)
	_ http.Pusher   = (*ResponseWriter)(nil)
	_ io.ReaderFrom = (*ResponseWriter)(nil)
	_ flushErrorer  = (*ResponseWriter)(nil)
)

// ResponseWriter wraps a http.ResponseWriter to track its write state.
type ResponseWriter struct {
	http.ResponseWriter

	written bool
	status  int
}

func (rw *ResponseWriter) WriteHeader(status int) {
	if rw.written {
		return
	}

	rw.written = true
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(b)
}

// Written implements [WriteTracker] and returns whether the current response body has been already written.
func (rw *ResponseWriter) Written() bool {
	return rw.written
}

// Written implements [StatusTracker] and returns the written status code of the current response.
func (rw *ResponseWriter) Status() int {
	return rw.status
}

// Flush implements [http.Flusher] and allows an HTTP handler to flush buffered data to the client.
// This method is no-op if the wrapped writer doesn't support it.
func (rw *ResponseWriter) Flush() {
	_ = rw.FlushError()
}

// FlushError is similar to [Flush] but returns [http.ErrNotSupported]
// if the wrapped writer doesn't support it.
func (rw *ResponseWriter) FlushError() error {
	err := http.NewResponseController(rw.ResponseWriter).Flush()
	if err == nil || !errors.Is(err, http.ErrNotSupported) {
		rw.written = true
	}
	return err
}

// Hijack implements [http.Hijacker] and allows an HTTP handler to take over the current connection.
func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(rw.ResponseWriter).Hijack()
}

// Pusher implements [http.Pusher] to indicate HTTP/2 server push support.
func (rw *ResponseWriter) Push(target string, opts *http.PushOptions) error {
	w := rw.ResponseWriter
	for {
		switch p := w.(type) {
		case http.Pusher:
			return p.Push(target, opts)
		case RWUnwrapper:
			w = p.Unwrap()
		default:
			return http.ErrNotSupported
		}
	}
}

// ReaderFrom implements [io.ReaderFrom] by checking if the underlying writer supports it.
// Otherwise calls [io.Copy].
func (rw *ResponseWriter) ReadFrom(r io.Reader) (n int64, err error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}

	w := rw.ResponseWriter
	for {
		switch rf := w.(type) {
		case io.ReaderFrom:
			return rf.ReadFrom(r)
		case RWUnwrapper:
			w = rf.Unwrap()
		default:
			return io.Copy(rw.ResponseWriter, r)
		}
	}
}

// Unwrap returns the underlying ResponseWritter instance (usually used by [http.ResponseController]).
func (rw *ResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func getWritten(rw http.ResponseWriter) (bool, error) {
	for {
		switch w := rw.(type) {
		case WriteTracker:
			return w.Written(), nil
		case RWUnwrapper:
			rw = w.Unwrap()
		default:
			return false, http.ErrNotSupported
		}
	}
}

func getStatus(rw http.ResponseWriter) (int, error) {
	for {
		switch w := rw.(type) {
		case StatusTracker:
			return w.Status(), nil
		case RWUnwrapper:
			rw = w.Unwrap()
		default:
			return 0, http.ErrNotSupported
		}
	}
}
