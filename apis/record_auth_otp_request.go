package apis

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/security"
)

func recordRequestOTP(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.OTP.Enabled {
		return e.ForbiddenError("The collection is not configured to allow OTP authentication.", nil)
	}

	form := &createOTPForm{}
	if err = e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err = form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	record, err := e.App.FindAuthRecordByEmail(collection, form.Email)

	// ignore not found errors to allow custom record find implementations
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return e.InternalServerError("", err)
	}

	event := new(core.RecordCreateOTPRequestEvent)
	event.RequestEvent = e
	event.Password = security.RandomStringWithAlphabet(collection.OTP.Length, "1234567890")
	event.Collection = collection
	event.Record = record

	originalApp := e.App

	return e.App.OnRecordRequestOTPRequest().Trigger(event, func(e *core.RecordCreateOTPRequestEvent) error {
		if e.Record == nil {
			// write a dummy 200 response as a very rudimentary emails enumeration "protection"
			e.JSON(http.StatusOK, map[string]string{
				"otpId": core.GenerateDefaultRandomId(),
			})

			return fmt.Errorf("missing or invalid %s OTP auth record with email %s", collection.Name, form.Email)
		}

		var otp *core.OTP

		// limit the new OTP creations for a single user
		if !e.App.IsDev() {
			otps, err := e.App.FindAllOTPsByRecord(e.Record)
			if err != nil {
				return firstApiError(err, e.InternalServerError("Failed to fetch previous record OTPs.", err))
			}

			totalRecent := 0
			for _, existingOTP := range otps {
				if !existingOTP.HasExpired(collection.OTP.DurationTime()) {
					totalRecent++
				}
				// use the last issued one
				if totalRecent > 9 {
					otp = otps[0] // otps are DESC sorted
					e.App.Logger().Warn(
						"Too many OTP requests - reusing the last issued",
						"email", form.Email,
						"recordId", e.Record.Id,
						"otpId", existingOTP.Id,
					)
					break
				}
			}
		}

		if otp == nil {
			// create new OTP
			// ---
			otp = core.NewOTP(e.App)
			otp.SetCollectionRef(e.Record.Collection().Id)
			otp.SetRecordRef(e.Record.Id)
			otp.SetPassword(e.Password)
			err = e.App.Save(otp)
			if err != nil {
				return err
			}

			// send OTP email
			// (in the background as a very basic timing attacks and emails enumeration protection)
			// ---
			routine.FireAndForget(func() {
				err = mails.SendRecordOTP(originalApp, e.Record, otp.Id, e.Password)
				if err != nil {
					originalApp.Logger().Error("Failed to send OTP email", "error", errors.Join(err, originalApp.Delete(otp)))
				}
			})
		}

		return e.JSON(http.StatusOK, map[string]string{
			"otpId": otp.Id,
		})
	})
}

// -------------------------------------------------------------------

type createOTPForm struct {
	Email string `form:"email" json:"email"`
}

func (form createOTPForm) validate() error {
	return validation.ValidateStruct(&form,
		validation.Field(&form.Email, validation.Required, validation.Length(1, 255), is.EmailFormat),
	)
}
