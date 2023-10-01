package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/spf13/cast"
)

// BindBody binds request body content to i.
//
// This is similar to `echo.BindBody()`, but for JSON requests uses
// custom json reader that **copies** the request body, allowing multiple reads.
func BindBody(c echo.Context, i any) error {
	req := c.Request()
	if req.ContentLength == 0 {
		return nil
	}

	ctype := req.Header.Get(echo.HeaderContentType)
	switch {
	case strings.HasPrefix(ctype, echo.MIMEApplicationJSON):
		err := CopyJsonBody(c.Request(), i)
		if err != nil {
			return echo.NewHTTPErrorWithInternal(http.StatusBadRequest, err, err.Error())
		}
		return nil
	case strings.HasPrefix(ctype, echo.MIMEApplicationForm), strings.HasPrefix(ctype, echo.MIMEMultipartForm):
		return bindFormData(c, i)
	}

	// fallback to the default binder
	return echo.BindBody(c, i)
}

// CopyJsonBody reads the request body into i by
// creating a copy of `r.Body` to allow multiple reads.
func CopyJsonBody(r *http.Request, i any) error {
	body := r.Body

	// this usually shouldn't be needed because the Server calls close for us
	// but we are changing the request body with a new reader
	defer body.Close()

	limitReader := io.LimitReader(body, DefaultMaxMemory)

	bodyBytes, readErr := io.ReadAll(limitReader)
	if readErr != nil {
		return readErr
	}

	err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(i)

	// set new body reader
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return err
}

// This is temp hotfix for properly binding multipart/form-data array values
// when a map destination is used.
//
// It should be replaced with echo.BindBody(c, i) once the issue is fixed in echo.
func bindFormData(c echo.Context, i any) error {
	if i == nil {
		return nil
	}

	values, err := c.FormValues()
	if err != nil {
		return echo.NewHTTPErrorWithInternal(http.StatusBadRequest, err, err.Error())
	}

	if len(values) == 0 {
		return nil
	}

	rt := reflect.TypeOf(i).Elem()

	// map
	if rt.Kind() == reflect.Map {
		rv := reflect.ValueOf(i).Elem()

		for k, v := range values {
			if total := len(v); total == 1 {
				rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(normalizeMultipartValue(v[0])))
			} else {
				normalized := make([]any, total)
				for i, vItem := range v {
					normalized[i] = normalizeMultipartValue(vItem)
				}
				rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(normalized))
			}
		}

		return nil
	}

	// anything else
	return echo.BindBody(c, i)
}

// In order to support more seamlessly both json and multipart/form-data requests,
// the following normalization rules are applied for plain multipart string values:
// - "true" is converted to the json `true`
// - "false" is converted to the json `false`
// - numeric (non-scientific) strings are converted to json number
// - any other string (empty string too) is left as it is
func normalizeMultipartValue(raw string) any {
	switch raw {
	case "":
		return raw
	case "true":
		return true
	case "false":
		return false
	default:
		if raw[0] == '-' || (raw[0] >= '0' && raw[0] <= '9') {
			if v, err := cast.ToFloat64E(raw); err == nil {
				return v
			}
		}

		return raw
	}
}
