package apis

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tools/router"
)

// bindSettingsApi registers the settings api endpoints.
func bindSettingsApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/settings").Bind(RequireSuperuserAuth())
	subGroup.GET("", settingsList)
	subGroup.PATCH("", settingsSet)
	subGroup.POST("/test/s3", settingsTestS3)
	subGroup.POST("/test/email", settingsTestEmail)
	subGroup.POST("/apple/generate-client-secret", settingsGenerateAppleClientSecret)
}

func settingsList(e *core.RequestEvent) error {
	clone, err := e.App.Settings().Clone()
	if err != nil {
		return e.InternalServerError("", err)
	}

	event := new(core.SettingsListRequestEvent)
	event.RequestEvent = e
	event.Settings = clone

	return e.App.OnSettingsListRequest().Trigger(event, func(e *core.SettingsListRequestEvent) error {
		return e.JSON(http.StatusOK, e.Settings)
	})
}

func settingsSet(e *core.RequestEvent) error {
	event := new(core.SettingsUpdateRequestEvent)
	event.RequestEvent = e

	if clone, err := e.App.Settings().Clone(); err == nil {
		event.OldSettings = clone
	} else {
		return e.BadRequestError("", err)
	}

	if clone, err := e.App.Settings().Clone(); err == nil {
		event.NewSettings = clone
	} else {
		return e.BadRequestError("", err)
	}

	if err := e.BindBody(&event.NewSettings); err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}

	return e.App.OnSettingsUpdateRequest().Trigger(event, func(e *core.SettingsUpdateRequestEvent) error {
		err := e.App.Save(e.NewSettings)
		if err != nil {
			return e.BadRequestError("An error occurred while saving the new settings.", err)
		}

		appSettings, err := e.App.Settings().Clone()
		if err != nil {
			return e.InternalServerError("Failed to clone app settings.", err)
		}

		return e.JSON(http.StatusOK, appSettings)
	})
}

func settingsTestS3(e *core.RequestEvent) error {
	form := forms.NewTestS3Filesystem(e.App)

	// load request
	if err := e.BindBody(form); err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}

	// send
	if err := form.Submit(); err != nil {
		// form error
		if fErr, ok := err.(validation.Errors); ok {
			return e.BadRequestError("Failed to test the S3 filesystem.", fErr)
		}

		// mailer error
		return e.BadRequestError("Failed to test the S3 filesystem. Raw error: \n"+err.Error(), nil)
	}

	return e.NoContent(http.StatusNoContent)
}

func settingsTestEmail(e *core.RequestEvent) error {
	form := forms.NewTestEmailSend(e.App)

	// load request
	if err := e.BindBody(form); err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}

	// send
	if err := form.Submit(); err != nil {
		// form error
		if fErr, ok := err.(validation.Errors); ok {
			return e.BadRequestError("Failed to send the test email.", fErr)
		}

		// mailer error
		return e.BadRequestError("Failed to send the test email. Raw error: \n"+err.Error(), nil)
	}

	return e.NoContent(http.StatusNoContent)
}

func settingsGenerateAppleClientSecret(e *core.RequestEvent) error {
	form := forms.NewAppleClientSecretCreate(e.App)

	// load request
	if err := e.BindBody(form); err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}

	// generate
	secret, err := form.Submit()
	if err != nil {
		// form error
		if fErr, ok := err.(validation.Errors); ok {
			return e.BadRequestError("Invalid client secret data.", fErr)
		}

		// secret generation error
		return e.BadRequestError("Failed to generate client secret. Raw error: \n"+err.Error(), nil)
	}

	return e.JSON(http.StatusOK, map[string]string{
		"secret": secret,
	})
}
