package apis

// -------------------------------------------------------------------
// This middleware is ported from echo/middleware to minimize the breaking
// changes and differences in the API behavior from earlier PocketBase versions
// (https://github.com/labstack/echo/blob/ec5b858dab6105ab4c3ed2627d1ebdfb6ae1ecb8/middleware/compress.go).
// -------------------------------------------------------------------

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/router"
)

const (
	gzipScheme = "gzip"
)

const (
	DefaultGzipMiddlewareId = "pbGzip"
)

// GzipConfig defines the config for Gzip middleware.
type GzipConfig struct {
	// Gzip compression level.
	// Optional. Default value -1.
	Level int

	// Length threshold before gzip compression is applied.
	// Optional. Default value 0.
	//
	// Most of the time you will not need to change the default. Compressing
	// a short response might increase the transmitted data because of the
	// gzip format overhead. Compressing the response will also consume CPU
	// and time on the server and the client (for decompressing). Depending on
	// your use case such a threshold might be useful.
	//
	// See also:
	// https://webmasters.stackexchange.com/questions/31750/what-is-recommended-minimum-object-size-for-gzip-performance-benefits
	MinLength int
}

// Gzip returns a middleware which compresses HTTP response using Gzip compression scheme.
func Gzip() *hook.Handler[*core.RequestEvent] {
	return GzipWithConfig(GzipConfig{})
}

// GzipWithConfig returns a middleware which compresses HTTP response using gzip compression scheme.
func GzipWithConfig(config GzipConfig) *hook.Handler[*core.RequestEvent] {
	if config.Level < -2 || config.Level > 9 { // these are consts: gzip.HuffmanOnly and gzip.BestCompression
		panic(errors.New("invalid gzip level"))
	}
	if config.Level == 0 {
		config.Level = -1
	}
	if config.MinLength < 0 {
		config.MinLength = 0
	}

	pool := sync.Pool{
		New: func() interface{} {
			w, err := gzip.NewWriterLevel(io.Discard, config.Level)
			if err != nil {
				return err
			}
			return w
		},
	}

	bpool := sync.Pool{
		New: func() interface{} {
			b := &bytes.Buffer{}
			return b
		},
	}

	return &hook.Handler[*core.RequestEvent]{
		Id: DefaultGzipMiddlewareId,
		Func: func(e *core.RequestEvent) error {
			e.Response.Header().Add("Vary", "Accept-Encoding")
			if strings.Contains(e.Request.Header.Get("Accept-Encoding"), gzipScheme) {
				w, ok := pool.Get().(*gzip.Writer)
				if !ok {
					return e.InternalServerError("", errors.New("failed to get gzip.Writer"))
				}

				rw := e.Response
				w.Reset(rw)

				buf := bpool.Get().(*bytes.Buffer)
				buf.Reset()

				grw := &gzipResponseWriter{Writer: w, ResponseWriter: rw, minLength: config.MinLength, buffer: buf}
				defer func() {
					// There are different reasons for cases when we have not yet written response to the client and now need to do so.
					// a) handler response had only response code and no response body (ala 404 or redirects etc). Response code need to be written now.
					// b) body is shorter than our minimum length threshold and being buffered currently and needs to be written
					if !grw.wroteBody {
						if rw.Header().Get("Content-Encoding") == gzipScheme {
							rw.Header().Del("Content-Encoding")
						}
						if grw.wroteHeader {
							rw.WriteHeader(grw.code)
						}
						// We have to reset response to it's pristine state when
						// nothing is written to body or error is returned.
						// See issue echo#424, echo#407.
						e.Response = rw
						w.Reset(io.Discard)
					} else if !grw.minLengthExceeded {
						// Write uncompressed response
						e.Response = rw
						if grw.wroteHeader {
							rw.WriteHeader(grw.code)
						}
						grw.buffer.WriteTo(rw)
						w.Reset(io.Discard)
					}
					w.Close()
					bpool.Put(buf)
					pool.Put(w)
				}()
				e.Response = grw
			}

			return e.Next()
		},
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
	buffer            *bytes.Buffer
	minLength         int
	code              int
	wroteHeader       bool
	wroteBody         bool
	minLengthExceeded bool
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	w.Header().Del("Content-Length") // Issue echo#444

	w.wroteHeader = true

	// Delay writing of the header until we know if we'll actually compress the response
	w.code = code
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	w.wroteBody = true

	if !w.minLengthExceeded {
		n, err := w.buffer.Write(b)

		if w.buffer.Len() >= w.minLength {
			w.minLengthExceeded = true

			// The minimum length is exceeded, add Content-Encoding header and write the header
			w.Header().Set("Content-Encoding", gzipScheme)
			if w.wroteHeader {
				w.ResponseWriter.WriteHeader(w.code)
			}

			return w.Writer.Write(w.buffer.Bytes())
		}

		return n, err
	}

	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	if !w.minLengthExceeded {
		// Enforce compression because we will not know how much more data will come
		w.minLengthExceeded = true
		w.Header().Set("Content-Encoding", gzipScheme)
		if w.wroteHeader {
			w.ResponseWriter.WriteHeader(w.code)
		}

		_, _ = w.Writer.Write(w.buffer.Bytes())
	}

	_ = w.Writer.(*gzip.Writer).Flush()

	_ = http.NewResponseController(w.ResponseWriter).Flush()
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(w.ResponseWriter).Hijack()
}

func (w *gzipResponseWriter) Push(target string, opts *http.PushOptions) error {
	rw := w.ResponseWriter
	for {
		switch p := rw.(type) {
		case http.Pusher:
			return p.Push(target, opts)
		case router.RWUnwrapper:
			rw = p.Unwrap()
		default:
			return http.ErrNotSupported
		}
	}
}

// Note: Disable the implementation for now because in case the platform
// supports the sendfile fast-path it won't run gzipResponseWriter.Write,
// preventing compression on the fly.
//
// func (w *gzipResponseWriter) ReadFrom(r io.Reader) (n int64, err error) {
// 	if w.wroteHeader {
// 		w.ResponseWriter.WriteHeader(w.code)
// 	}
// 	rw := w.ResponseWriter
// 	for {
// 		switch rf := rw.(type) {
// 		case io.ReaderFrom:
// 			return rf.ReadFrom(r)
// 		case router.RWUnwrapper:
// 			rw = rf.Unwrap()
// 		default:
// 			return io.Copy(w.ResponseWriter, r)
// 		}
// 	}
// }

func (w *gzipResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
