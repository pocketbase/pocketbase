package apis

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

// note: for now allow superusers but it may change in the future to allow access
// also to users with "Manage API" rule access depending on the use cases that will arise
func recordAuthImpersonate(e *core.RequestEvent) error {
	if !e.HasSuperuserAuth() {
		return e.ForbiddenError("", nil)
	}

	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	record, err := e.App.FindRecordById(collection, e.Request.PathValue("id"))
	if err != nil {
		return e.NotFoundError("", err)
	}

	form := &impersonateForm{}
	if err = e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err = form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	token, err := record.NewStaticAuthToken(time.Duration(form.Duration) * time.Second)
	if err != nil {
		e.InternalServerError("Failed to generate static auth token", err)
	}

	return recordAuthResponse(e, record, token, "", nil)
}

// -------------------------------------------------------------------

type impersonateForm struct {
	// Duration is the optional custom token duration in seconds.
	Duration int64 `form:"duration" json:"duration"`
}

func (form *impersonateForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Duration, validation.Min(0)),
	)
}
