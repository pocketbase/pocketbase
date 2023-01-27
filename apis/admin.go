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

	event := new(core.AdminAuthEvent)
	event.HttpContext = c
	event.Admin = admin
	event.Token = token

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

	event := new(core.AdminAuthRefreshEvent)
	event.HttpContext = c
	event.Admin = admin

	handlerErr := api.app.OnAdminBeforeAuthRefreshRequest().Trigger(event, func(e *core.AdminAuthRefreshEvent) error {
		return api.authResponse(e.HttpContext, e.Admin)
	})

	if handlerErr == nil {
		if err := api.app.OnAdminAfterAuthRefreshRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return handlerErr
}

func (api *adminApi) authWithPassword(c echo.Context) error {
	form := forms.NewAdminLogin(api.app)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	event := new(core.AdminAuthWithPasswordEvent)
	event.HttpContext = c
	event.Password = form.Password
	event.Identity = form.Identity

	_, submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
		return func(admin *models.Admin) error {
			event.Admin = admin

			return api.app.OnAdminBeforeAuthWithPasswordRequest().Trigger(event, func(e *core.AdminAuthWithPasswordEvent) error {
				if err := next(e.Admin); err != nil {
					return NewBadRequestError("Failed to authenticate.", err)
				}

				return api.authResponse(e.HttpContext, e.Admin)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnAdminAfterAuthWithPasswordRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return submitErr
}

func (api *adminApi) requestPasswordReset(c echo.Context) error {
	form := forms.NewAdminPasswordResetRequest(api.app)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	if err := form.Validate(); err != nil {
		return NewBadRequestError("An error occurred while validating the form.", err)
	}

	event := new(core.AdminRequestPasswordResetEvent)
	event.HttpContext = c

	submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
		return func(Admin *models.Admin) error {
			event.Admin = Admin

			return api.app.OnAdminBeforeRequestPasswordResetRequest().Trigger(event, func(e *core.AdminRequestPasswordResetEvent) error {
				// run in background because we don't need to show the result to the client
				routine.FireAndForget(func() {
					if err := next(e.Admin); err != nil && api.app.IsDebug() {
						log.Println(err)
					}
				})

				return e.HttpContext.NoContent(http.StatusNoContent)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnAdminAfterRequestPasswordResetRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	} else if api.app.IsDebug() {
		log.Println(submitErr)
	}

	// don't return the response error to prevent emails enumeration
	if !c.Response().Committed {
		c.NoContent(http.StatusNoContent)
	}

	return nil
}

func (api *adminApi) confirmPasswordReset(c echo.Context) error {
	form := forms.NewAdminPasswordResetConfirm(api.app)
	if readErr := c.Bind(form); readErr != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	event := new(core.AdminConfirmPasswordResetEvent)
	event.HttpContext = c

	_, submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
		return func(admin *models.Admin) error {
			event.Admin = admin

			return api.app.OnAdminBeforeConfirmPasswordResetRequest().Trigger(event, func(e *core.AdminConfirmPasswordResetEvent) error {
				if err := next(e.Admin); err != nil {
					return NewBadRequestError("Failed to set new password.", err)
				}

				return e.HttpContext.NoContent(http.StatusNoContent)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnAdminAfterConfirmPasswordResetRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return submitErr
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

	event := new(core.AdminsListEvent)
	event.HttpContext = c
	event.Admins = admins
	event.Result = result

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

	event := new(core.AdminViewEvent)
	event.HttpContext = c
	event.Admin = admin

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

	event := new(core.AdminCreateEvent)
	event.HttpContext = c
	event.Admin = admin

	// create the admin
	submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
		return func(m *models.Admin) error {
			event.Admin = m

			return api.app.OnAdminBeforeCreateRequest().Trigger(event, func(e *core.AdminCreateEvent) error {
				if err := next(e.Admin); err != nil {
					return NewBadRequestError("Failed to create admin.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.Admin)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnAdminAfterCreateRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
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

	event := new(core.AdminUpdateEvent)
	event.HttpContext = c
	event.Admin = admin

	// update the admin
	submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
		return func(m *models.Admin) error {
			event.Admin = m

			return api.app.OnAdminBeforeUpdateRequest().Trigger(event, func(e *core.AdminUpdateEvent) error {
				if err := next(e.Admin); err != nil {
					return NewBadRequestError("Failed to update admin.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.Admin)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnAdminAfterUpdateRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
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

	event := new(core.AdminDeleteEvent)
	event.HttpContext = c
	event.Admin = admin

	handlerErr := api.app.OnAdminBeforeDeleteRequest().Trigger(event, func(e *core.AdminDeleteEvent) error {
		if err := api.app.Dao().DeleteAdmin(e.Admin); err != nil {
			return NewBadRequestError("Failed to delete admin.", err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		if err := api.app.OnAdminAfterDeleteRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return handlerErr
}
