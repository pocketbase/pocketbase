package apis

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tools/rest"
)

// BindSettingsApi registers the settings api endpoints.
func BindSettingsApi(app core.App, rg *echo.Group) {
	api := settingsApi{app: app}

	subGroup := rg.Group("/settings", ActivityLogger(app), RequireAdminAuth())
	subGroup.GET("", api.list)
	subGroup.PATCH("", api.set)
}

type settingsApi struct {
	app core.App
}

func (api *settingsApi) list(c echo.Context) error {
	settings, err := api.app.Settings().RedactClone()
	if err != nil {
		return rest.NewBadRequestError("", err)
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
	if err := c.Bind(form); err != nil {
		return rest.NewBadRequestError("An error occured while reading the submitted data.", err)
	}

	event := &core.SettingsUpdateEvent{
		HttpContext: c,
		OldSettings: api.app.Settings(),
		NewSettings: form.Settings,
	}

	handlerErr := api.app.OnSettingsBeforeUpdateRequest().Trigger(event, func(e *core.SettingsUpdateEvent) error {
		if err := form.Submit(); err != nil {
			return rest.NewBadRequestError("An error occured while submitting the form.", err)
		}

		redactedSettings, err := api.app.Settings().RedactClone()
		if err != nil {
			return rest.NewBadRequestError("", err)
		}

		return e.HttpContext.JSON(http.StatusOK, redactedSettings)
	})

	if handlerErr == nil {
		api.app.OnSettingsAfterUpdateRequest().Trigger(event)
	}

	return handlerErr
}
