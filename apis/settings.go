package apis

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tools/security"
)

// bindSettingsApi registers the settings api endpoints.
func bindSettingsApi(app core.App, rg *echo.Group) {
	api := settingsApi{app: app}

	subGroup := rg.Group("/settings", ActivityLogger(app), RequireAdminAuth())
	subGroup.GET("", api.list)
	subGroup.PATCH("", api.set)
	subGroup.POST("/test/s3", api.testS3)
	subGroup.POST("/test/email", api.testEmail)
}

type settingsApi struct {
	app core.App
}

func (api *settingsApi) list(c echo.Context) error {
	settings, err := api.app.Settings().RedactClone()
	if err != nil {
		return NewBadRequestError("", err)
	}

	event := &core.SettingsListEvent{
		HttpContext:      c,
		RedactedSettings: settings,
	}

	return api.app.OnSettingsListRequest().Trigger(event, func(e *core.SettingsListEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.RedactedSettings)
	})
}

func (api *settingsApi) set(c echo.Context) error {
	form := forms.NewSettingsUpsert(api.app)

	// load request
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	event := &core.SettingsUpdateEvent{
		HttpContext: c,
		OldSettings: api.app.Settings(),
		NewSettings: form.Settings,
	}

	// update the settings
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnSettingsBeforeUpdateRequest().Trigger(event, func(e *core.SettingsUpdateEvent) error {
				if err := next(); err != nil {
					return NewBadRequestError("An error occurred while submitting the form.", err)
				}

				redactedSettings, err := api.app.Settings().RedactClone()
				if err != nil {
					return NewBadRequestError("", err)
				}

				return e.HttpContext.JSON(http.StatusOK, redactedSettings)
			})
		}
	})

	if submitErr == nil {
		api.app.OnSettingsAfterUpdateRequest().Trigger(event)
	}

	return submitErr
}

func (api *settingsApi) testS3(c echo.Context) error {
	if !api.app.Settings().S3.Enabled {
		return NewBadRequestError("S3 storage is not enabled.", nil)
	}

	fs, err := api.app.NewFilesystem()
	if err != nil {
		return NewBadRequestError("Failed to initialize the S3 storage. Raw error: \n"+err.Error(), nil)
	}
	defer fs.Close()

	testFileKey := "pb_test_" + security.PseudorandomString(5) + "/test.txt"

	if err := fs.Upload([]byte("test"), testFileKey); err != nil {
		return NewBadRequestError("Failed to upload a test file. Raw error: \n"+err.Error(), nil)
	}

	if err := fs.Delete(testFileKey); err != nil {
		return NewBadRequestError("Failed to delete a test file. Raw error: \n"+err.Error(), nil)
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *settingsApi) testEmail(c echo.Context) error {
	form := forms.NewTestEmailSend(api.app)

	// load request
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	// send
	if err := form.Submit(); err != nil {
		if fErr, ok := err.(validation.Errors); ok {
			// form error
			return NewBadRequestError("Failed to send the test email.", fErr)
		}

		// mailer error
		return NewBadRequestError("Failed to send the test email. Raw error: \n"+err.Error(), nil)
	}

	return c.NoContent(http.StatusNoContent)
}
