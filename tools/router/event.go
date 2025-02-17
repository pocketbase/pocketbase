package router

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/netip"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/picker"
	"github.com/pocketbase/pocketbase/tools/store"
)

var ErrUnsupportedContentType = NewBadRequestError("Unsupported Content-Type", nil)
var ErrInvalidRedirectStatusCode = NewInternalServerError("Invalid redirect status code", nil)
var ErrFileNotFound = NewNotFoundError("File not found", nil)

const IndexPage = "index.html"

// Event specifies based Route handler event that is usually intended
// to be embedded as part of a custom event struct.
//
// NB! It is expected that the Response and Request fields are always set.
type Event struct {
	Response http.ResponseWriter
	Request  *http.Request

	hook.Event

	data store.Store[string, any]
}

// RWUnwrapper specifies that an http.ResponseWriter could be "unwrapped"
// (usually used with [http.ResponseController]).
type RWUnwrapper interface {
	Unwrap() http.ResponseWriter
}

// Written reports whether the current response has already been written.
//
// This method always returns false if e.ResponseWritter doesn't implement the WriteTracker interface
// (all router package handlers receives a ResponseWritter that implements it unless explicitly replaced with a custom one).
func (e *Event) Written() bool {
	written, _ := getWritten(e.Response)
	return written
}

// Status reports the status code of the current response.
//
// This method always returns 0 if e.Response doesn't implement the StatusTracker interface
// (all router package handlers receives a ResponseWritter that implements it unless explicitly replaced with a custom one).
func (e *Event) Status() int {
	status, _ := getStatus(e.Response)
	return status
}

// Flush flushes buffered data to the current response.
//
// Returns [http.ErrNotSupported] if e.Response doesn't implement the [http.Flusher] interface
// (all router package handlers receives a ResponseWritter that implements it unless explicitly replaced with a custom one).
func (e *Event) Flush() error {
	return http.NewResponseController(e.Response).Flush()
}

// IsTLS reports whether the connection on which the request was received is TLS.
func (e *Event) IsTLS() bool {
	return e.Request.TLS != nil
}

// SetCookie is an alias for [http.SetCookie].
//
// SetCookie adds a Set-Cookie header to the current response's headers.
// The provided cookie must have a valid Name.
// Invalid cookies may be silently dropped.
func (e *Event) SetCookie(cookie *http.Cookie) {
	http.SetCookie(e.Response, cookie)
}

// RemoteIP returns the IP address of the client that sent the request.
//
// IPv6 addresses are returned expanded.
// For example, "2001:db8::1" becomes "2001:0db8:0000:0000:0000:0000:0000:0001".
//
// Note that if you are behind reverse proxy(ies), this method returns
// the IP of the last connecting proxy.
func (e *Event) RemoteIP() string {
	ip, _, _ := net.SplitHostPort(e.Request.RemoteAddr)
	parsed, _ := netip.ParseAddr(ip)
	return parsed.StringExpanded()
}

// FindUploadedFiles extracts all form files of "key" from a http request
// and returns a slice with filesystem.File instances (if any).
func (e *Event) FindUploadedFiles(key string) ([]*filesystem.File, error) {
	if e.Request.MultipartForm == nil {
		err := e.Request.ParseMultipartForm(DefaultMaxMemory)
		if err != nil {
			return nil, err
		}
	}

	if e.Request.MultipartForm == nil || e.Request.MultipartForm.File == nil || len(e.Request.MultipartForm.File[key]) == 0 {
		return nil, http.ErrMissingFile
	}

	result := make([]*filesystem.File, 0, len(e.Request.MultipartForm.File[key]))

	for _, fh := range e.Request.MultipartForm.File[key] {
		file, err := filesystem.NewFileFromMultipart(fh)
		if err != nil {
			return nil, err
		}

		result = append(result, file)
	}

	return result, nil
}

// Store
// -------------------------------------------------------------------

// Get retrieves single value from the current event data store.
func (e *Event) Get(key string) any {
	return e.data.Get(key)
}

// GetAll returns a copy of the current event data store.
func (e *Event) GetAll() map[string]any {
	return e.data.GetAll()
}

// Set saves single value into the current event data store.
func (e *Event) Set(key string, value any) {
	e.data.Set(key, value)
}

// SetAll saves all items from m into the current event data store.
func (e *Event) SetAll(m map[string]any) {
	for k, v := range m {
		e.Set(k, v)
	}
}

// Response writers
// -------------------------------------------------------------------

const headerContentType = "Content-Type"

func (e *Event) setResponseHeaderIfEmpty(key, value string) {
	header := e.Response.Header()
	if header.Get(key) == "" {
		header.Set(key, value)
	}
}

// String writes a plain string response.
func (e *Event) String(status int, data string) error {
	e.setResponseHeaderIfEmpty(headerContentType, "text/plain; charset=utf-8")
	e.Response.WriteHeader(status)
	_, err := e.Response.Write([]byte(data))
	return err
}

// HTML writes an HTML response.
func (e *Event) HTML(status int, data string) error {
	e.setResponseHeaderIfEmpty(headerContentType, "text/html; charset=utf-8")
	e.Response.WriteHeader(status)
	_, err := e.Response.Write([]byte(data))
	return err
}

const jsonFieldsParam = "fields"

// JSON writes a JSON response.
//
// It also provides a generic response data fields picker if the "fields" query parameter is set.
// For example, if you are requesting `?fields=a,b` for `e.JSON(200, map[string]int{ "a":1, "b":2, "c":3 })`,
// it should result in a JSON response like: `{"a":1, "b": 2}`.
func (e *Event) JSON(status int, data any) error {
	e.setResponseHeaderIfEmpty(headerContentType, "application/json")
	e.Response.WriteHeader(status)

	rawFields := e.Request.URL.Query().Get(jsonFieldsParam)

	// error response or no fields to pick
	if rawFields == "" || status < 200 || status > 299 {
		return json.NewEncoder(e.Response).Encode(data)
	}

	// pick only the requested fields
	modified, err := picker.Pick(data, rawFields)
	if err != nil {
		return err
	}

	return json.NewEncoder(e.Response).Encode(modified)
}

// XML writes an XML response.
// It automatically prepends the generic [xml.Header] string to the response.
func (e *Event) XML(status int, data any) error {
	e.setResponseHeaderIfEmpty(headerContentType, "application/xml; charset=utf-8")
	e.Response.WriteHeader(status)
	if _, err := e.Response.Write([]byte(xml.Header)); err != nil {
		return err
	}
	return xml.NewEncoder(e.Response).Encode(data)
}

// Stream streams the specified reader into the response.
func (e *Event) Stream(status int, contentType string, reader io.Reader) error {
	e.Response.Header().Set(headerContentType, contentType)
	e.Response.WriteHeader(status)
	_, err := io.Copy(e.Response, reader)
	return err
}

// Blob writes a blob (bytes slice) response.
func (e *Event) Blob(status int, contentType string, b []byte) error {
	e.setResponseHeaderIfEmpty(headerContentType, contentType)
	e.Response.WriteHeader(status)
	_, err := e.Response.Write(b)
	return err
}

// FileFS serves the specified filename from fsys.
//
// It is similar to [echo.FileFS] for consistency with earlier versions.
func (e *Event) FileFS(fsys fs.FS, filename string) error {
	f, err := fsys.Open(filename)
	if err != nil {
		return ErrFileNotFound
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	// if it is a directory try to open its index.html file
	if fi.IsDir() {
		filename = filepath.ToSlash(filepath.Join(filename, IndexPage))
		f, err = fsys.Open(filename)
		if err != nil {
			return ErrFileNotFound
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil {
			return err
		}
	}

	ff, ok := f.(io.ReadSeeker)
	if !ok {
		return errors.New("[FileFS] file does not implement io.ReadSeeker")
	}

	http.ServeContent(e.Response, e.Request, fi.Name(), fi.ModTime(), ff)

	return nil
}

// NoContent writes a response with no body (ex. 204).
func (e *Event) NoContent(status int) error {
	e.Response.WriteHeader(status)
	return nil
}

// Redirect writes a redirect response to the specified url.
// The status code must be in between 300 – 399 range.
func (e *Event) Redirect(status int, url string) error {
	if status < 300 || status > 399 {
		return ErrInvalidRedirectStatusCode
	}
	e.Response.Header().Set("Location", url)
	e.Response.WriteHeader(status)
	return nil
}

// ApiError helpers
// -------------------------------------------------------------------

func (e *Event) Error(status int, message string, errData any) *ApiError {
	return NewApiError(status, message, errData)
}

func (e *Event) BadRequestError(message string, errData any) *ApiError {
	return NewBadRequestError(message, errData)
}

func (e *Event) NotFoundError(message string, errData any) *ApiError {
	return NewNotFoundError(message, errData)
}

func (e *Event) ForbiddenError(message string, errData any) *ApiError {
	return NewForbiddenError(message, errData)
}

func (e *Event) UnauthorizedError(message string, errData any) *ApiError {
	return NewUnauthorizedError(message, errData)
}

func (e *Event) TooManyRequestsError(message string, errData any) *ApiError {
	return NewTooManyRequestsError(message, errData)
}

func (e *Event) InternalServerError(message string, errData any) *ApiError {
	return NewInternalServerError(message, errData)
}

// Binders
// -------------------------------------------------------------------

const DefaultMaxMemory = 32 << 20 // 32mb

// BindBody unmarshal the request body into the provided dst.
//
// dst must be either a struct pointer or map[string]any.
//
// The rules how the body will be scanned depends on the request Content-Type.
//
// Currently the following Content-Types are supported:
//   - application/json
//   - text/xml, application/xml
//   - multipart/form-data, application/x-www-form-urlencoded
//
// Respectively the following struct tags are supported (again, which one will be used depends on the Content-Type):
//   - "json" (json body)- uses the builtin Go json package for unmarshaling.
//   - "xml" (xml body) - uses the builtin Go xml package for unmarshaling.
//   - "form" (form data) - utilizes the custom [router.UnmarshalRequestData] method.
//
// NB! When dst is a struct make sure that it doesn't have public fields
// that shouldn't be bindable and it is advisible such fields to be unexported
// or have a separate struct just for the binding. For example:
//
//	data := struct{
//	   somethingPrivate string
//
//	   Title string `json:"title" form:"title"`
//	   Total int    `json:"total" form:"total"`
//	}
//	err := e.BindBody(&data)
func (e *Event) BindBody(dst any) error {
	if e.Request.ContentLength == 0 {
		return nil
	}

	contentType := e.Request.Header.Get(headerContentType)

	if strings.HasPrefix(contentType, "application/json") {
		dec := json.NewDecoder(e.Request.Body)
		err := dec.Decode(dst)
		if err == nil {
			// manually call Reread because single call of json.Decoder.Decode()
			// doesn't ensure that the entire body is a valid json string
			// and it is not guaranteed that it will reach EOF to trigger the reread reset
			// (ex. in case of trailing spaces or invalid trailing parts like: `{"test":1},something`)
			if body, ok := e.Request.Body.(Rereader); ok {
				body.Reread()
			}
		}
		return err
	}

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := e.Request.ParseMultipartForm(DefaultMaxMemory); err != nil {
			return err
		}

		return UnmarshalRequestData(e.Request.Form, dst, "", "")
	}

	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		if err := e.Request.ParseForm(); err != nil {
			return err
		}

		return UnmarshalRequestData(e.Request.Form, dst, "", "")
	}

	if strings.HasPrefix(contentType, "text/xml") ||
		strings.HasPrefix(contentType, "application/xml") {
		return xml.NewDecoder(e.Request.Body).Decode(dst)
	}

	return ErrUnsupportedContentType
}
