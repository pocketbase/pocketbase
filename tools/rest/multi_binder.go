package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

// BindBody binds request body content to i.
//
// This is similar to `echo.BindBody()`, but for JSON requests uses
// custom json reader that **copies** the request body, allowing multiple reads.
func BindBody(c echo.Context, i interface{}) error {
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
	default:
		// fallback to the default binder
		return echo.BindBody(c, i)
	}
}

// CopyJsonBody reads the request body into i by
// creating a copy of `r.Body` to allow multiple reads.
func CopyJsonBody(r *http.Request, i interface{}) error {
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
