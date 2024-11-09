package apis

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

func recordAuthWithOTP(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.OTP.Enabled {
		return e.ForbiddenError("The collection is not configured to allow OTP authentication.", nil)
	}

	form := &authWithOTPForm{}
	if err = e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err = form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	event := new(core.RecordAuthWithOTPRequestEvent)
	event.RequestEvent = e
	event.Collection = collection

	// extra validations
	// (note: returns a generic 400 as a very basic OTPs enumeration protection)
	// ---
	event.OTP, err = e.App.FindOTPById(form.OTPId)
	if err != nil {
		return e.BadRequestError("Invalid or expired OTP", err)
	}

	if event.OTP.CollectionRef() != collection.Id {
		return e.BadRequestError("Invalid or expired OTP", errors.New("the OTP is for a different collection"))
	}

	if event.OTP.HasExpired(collection.OTP.DurationTime()) {
		return e.BadRequestError("Invalid or expired OTP", errors.New("the OTP is expired"))
	}

	event.Record, err = e.App.FindRecordById(event.OTP.CollectionRef(), event.OTP.RecordRef())
	if err != nil {
		return e.BadRequestError("Invalid or expired OTP", fmt.Errorf("missing auth record: %w", err))
	}

	// since otps are usually simple digit numbers we enforce an extra rate limit rule to prevent enumerations
	err = checkRateLimit(e, "@pb_otp_"+event.Record.Id, core.RateLimitRule{MaxRequests: 5, Duration: 180})
	if err != nil {
		return e.TooManyRequestsError("Too many attempts, please try again later with a new OTP.", nil)
	}

	if !event.OTP.ValidatePassword(form.Password) {
		return e.BadRequestError("Invalid or expired OTP", errors.New("incorrect password"))
	}
	// ---

	return e.App.OnRecordAuthWithOTPRequest().Trigger(event, func(e *core.RecordAuthWithOTPRequestEvent) error {
		err = RecordAuthResponse(e.RequestEvent, e.Record, core.MFAMethodOTP, nil)
		if err != nil {
			return err
		}

		// try to delete the used otp
		if e.OTP != nil {
			err = e.App.Delete(e.OTP)
			if err != nil {
				e.App.Logger().Error("Failed to delete used OTP", "error", err, "otpId", e.OTP.Id)
			}
		}

		// note: we don't update the user verified state the same way as in the password reset confirmation
		// at the moment because it is not clear whether the otp confirmation came from the user email
		// (e.g. it could be from an sms or some other channel)

		return nil
	})
}

// -------------------------------------------------------------------

type authWithOTPForm struct {
	OTPId    string `form:"otpId" json:"otpId"`
	Password string `form:"password" json:"password"`
}

func (form *authWithOTPForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.OTPId, validation.Required, validation.Length(1, 255)),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 71)),
	)
}
