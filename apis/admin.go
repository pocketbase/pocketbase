package apis

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/search"
)

// bindAdminApi registers the admin api endpoints and the corresponding handlers.
func bindAdminApi(app core.App, rg *echo.Group) {
	api := adminApi{app: app}

	subGroup := rg.Group("/admins", ActivityLogger(app))
	subGroup.POST("/auth-with-password", api.authWithPassword)
	subGroup.POST("/request-password-reset", api.requestPasswordReset)
	subGroup.POST("/confirm-password-reset", api.confirmPasswordReset)
	subGroup.POST("/auth-refresh", api.authRefresh, RequireAdminAuth())
	subGroup.GET("", api.list, RequireAdminAuth())
	subGroup.POST("", api.create, RequireAdminAuthOnlyIfAny(app))
	subGroup.GET("/:id", api.view, RequireAdminAuth())
	subGroup.PATCH("/:id", api.update, RequireAdminAuth())
	subGroup.DELETE("/:id", api.delete, RequireAdminAuth())
}

type adminApi struct {
	app core.App
}

func (api *adminApi) authResponse(c echo.Context, admin *models.Admin) error {
	token, tokenErr := tokens.NewAdminAuthToken(api.app, admin)
	if tokenErr != nil {
		return NewBadRequestError("Failed to create auth token.", tokenErr)
	}

	event := &core.AdminAuthEvent{
		HttpContext: c,
		Admin:       admin,
		Token:       token,
	}

	return api.app.OnAdminAuthRequest().Trigger(event, func(e *core.AdminAuthEvent) error {
		return e.HttpContext.JSON(200, map[string]any{
			"token": e.Token,
			"admin": e.Admin,
		})
	})
}

func (api *adminApi) authRefresh(c echo.Context) error {
	admin, _ := c.Get(ContextAdminKey).(*models.Admin)
	if admin == nil {
		return NewNotFoundError("Missing auth admin context.", nil)
	}

	return api.authResponse(c, admin)
}

func (api *adminApi) authWithPassword(c echo.Context) error {
	form := forms.NewAdminLogin(api.app)
	if readErr := c.Bind(form); readErr != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	admin, submitErr := form.Submit()
	if submitErr != nil {
		return NewBadRequestError("Failed to authenticate.", submitErr)
	}

	return api.authResponse(c, admin)
}

func (api *adminApi) requestPasswordReset(c echo.Context) error {
	form := forms.NewAdminPasswordResetRequest(api.app)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	if err := form.Validate(); err != nil {
		return NewBadRequestError("An error occurred while validating the form.", err)
	}

	// run in background because we don't need to show the result
	// (prevents admins enumeration)
	routine.FireAndForget(func() {
		if err := form.Submit(); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	})

	return c.NoContent(http.StatusNoContent)
}

func (api *adminApi) confirmPasswordReset(c echo.Context) error {
	form := forms.NewAdminPasswordResetConfirm(api.app)
	if readErr := c.Bind(form); readErr != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	_, submitErr := form.Submit()
	if submitErr != nil {
		return NewBadRequestError("Failed to set new password.", submitErr)
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *adminApi) list(c echo.Context) error {
	fieldResolver := search.NewSimpleFieldResolver(
		"id", "created", "updated", "name", "email",
	)

	admins := []*models.Admin{}

	result, err := search.NewProvider(fieldResolver).
		Query(api.app.Dao().AdminQuery()).
		ParseAndExec(c.QueryParams().Encode(), &admins)

	if err != nil {
		return NewBadRequestError("", err)
	}

	event := &core.AdminsListEvent{
		HttpContext: c,
		Admins:      admins,
		Result:      result,
	}

	return api.app.OnAdminsListRequest().Trigger(event, func(e *core.AdminsListEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.Result)
	})
}

func (api *adminApi) view(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return NewNotFoundError("", nil)
	}

	admin, err := api.app.Dao().FindAdminById(id)
	if err != nil || admin == nil {
		return NewNotFoundError("", err)
	}

	event := &core.AdminViewEvent{
		HttpContext: c,
		Admin:       admin,
	}

	return api.app.OnAdminViewRequest().Trigger(event, func(e *core.AdminViewEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.Admin)
	})
}

func (api *adminApi) create(c echo.Context) error {
	admin := &models.Admin{}

	form := forms.NewAdminUpsert(api.app, admin)

	// load request
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.AdminCreateEvent{
		HttpContext: c,
		Admin:       admin,
	}

	// create the admin
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnAdminBeforeCreateRequest().Trigger(event, func(e *core.AdminCreateEvent) error {
				if err := next(); err != nil {
					return NewBadRequestError("Failed to create admin.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.Admin)
			})
		}
	})

	if submitErr == nil {
		api.app.OnAdminAfterCreateRequest().Trigger(event)
	}

	return submitErr
}

func (api *adminApi) update(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return NewNotFoundError("", nil)
	}

	admin, err := api.app.Dao().FindAdminById(id)
	if err != nil || admin == nil {
		return NewNotFoundError("", err)
	}

	form := forms.NewAdminUpsert(api.app, admin)

	// load request
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.AdminUpdateEvent{
		HttpContext: c,
		Admin:       admin,
	}

	// update the admin
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnAdminBeforeUpdateRequest().Trigger(event, func(e *core.AdminUpdateEvent) error {
				if err := next(); err != nil {
					return NewBadRequestError("Failed to update admin.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.Admin)
			})
		}
	})

	if submitErr == nil {
		api.app.OnAdminAfterUpdateRequest().Trigger(event)
	}

	return submitErr
}

func (api *adminApi) delete(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return NewNotFoundError("", nil)
	}

	admin, err := api.app.Dao().FindAdminById(id)
	if err != nil || admin == nil {
		return NewNotFoundError("", err)
	}

	event := &core.AdminDeleteEvent{
		HttpContext: c,
		Admin:       admin,
	}

	handlerErr := api.app.OnAdminBeforeDeleteRequest().Trigger(event, func(e *core.AdminDeleteEvent) error {
		if err := api.app.Dao().DeleteAdmin(e.Admin); err != nil {
			return NewBadRequestError("Failed to delete admin.", err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnAdminAfterDeleteRequest().Trigger(event)
	}

	return handlerErr
}
