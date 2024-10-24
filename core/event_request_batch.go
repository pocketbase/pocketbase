package core

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/hook"
)

type BatchRequestEvent struct {
	hook.Event
	*RequestEvent

	Batch []*InternalRequest
}

type InternalRequest struct {
	// note: for uploading files the value must be either *filesystem.File or []*filesystem.File
	Body map[string]any `form:"body" json:"body"`

	Headers map[string]string `form:"headers" json:"headers"`

	Method string `form:"method" json:"method"`

	URL string `form:"url" json:"url"`
}

func (br InternalRequest) Validate() error {
	return validation.ValidateStruct(&br,
		validation.Field(&br.Method, validation.Required, validation.In(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete)),
		validation.Field(&br.URL, validation.Required, validation.Length(0, 2000)),
	)
}
