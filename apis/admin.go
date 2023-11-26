package apis

import (
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

func (api *adminApi) authResponse(c echo.Context, admin *models.Admin, finalizers ...func(token string) error) error {
	token, tokenErr := tokens.NewAdminAuthToken(api.app, admin)
	if tokenErr != nil {
		return NewBadRequestError("Failed to create auth token.", tokenErr)
	}

	for _, f := range finalizers {
		if err := f(token); err != nil {
			return err
		}
	}

	event := new(core.AdminAuthEvent)
	event.HttpContext = c
	event.Admin = admin
	event.Token = token

	return api.app.OnAdminAuthRequest().Trigger(event, func(e *core.AdminAuthEvent) error {
		if e.HttpContext.Response().Committed {
			return nil
		}

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

	return api.app.OnAdminBeforeAuthRefreshRequest().Trigger(event, func(e *core.AdminAuthRefreshEvent) error {
		return api.app.OnAdminAfterAuthRefreshRequest().Trigger(event, func(e *core.AdminAuthRefreshEvent) error {
			return api.authResponse(e.HttpContext, e.Admin)
		})
	})
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

				return api.app.OnAdminAfterAuthWithPasswordRequest().Trigger(event, func(e *core.AdminAuthWithPasswordEvent) error {
					return api.authResponse(e.HttpContext, e.Admin)
				})
			})
		}
	})

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
					if err := next(e.Admin); err != nil {
						api.app.Logger().Error("Failed to send admin password reset request.", "error", err)
					}
				})

				return api.app.OnAdminAfterRequestPasswordResetRequest().Trigger(event, func(e *core.AdminRequestPasswordResetEvent) error {
					if e.HttpContext.Response().Committed {
						return nil
					}

					return e.HttpContext.NoContent(http.StatusNoContent)
				})
			})
		}
	})

	// eagerly write 204 response and skip submit errors
	// as a measure against admins enumeration
	if !c.Response().Committed {
		c.NoContent(http.StatusNoContent)
	}

	return submitErr
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

				return api.app.OnAdminAfterConfirmPasswordResetRequest().Trigger(event, func(e *core.AdminConfirmPasswordResetEvent) error {
					if e.HttpContext.Response().Committed {
						return nil
					}

					return e.HttpContext.NoContent(http.StatusNoContent)
				})
			})
		}
	})

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
		if e.HttpContext.Response().Committed {
			return nil
		}

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
		if e.HttpContext.Response().Committed {
			return nil
		}

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

				return api.app.OnAdminAfterCreateRequest().Trigger(event, func(e *core.AdminCreateEvent) error {
					if e.HttpContext.Response().Committed {
						return nil
					}

					return e.HttpContext.JSON(http.StatusOK, e.Admin)
				})
			})
		}
	})

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

				return api.app.OnAdminAfterUpdateRequest().Trigger(event, func(e *core.AdminUpdateEvent) error {
					if e.HttpContext.Response().Committed {
						return nil
					}

					return e.HttpContext.JSON(http.StatusOK, e.Admin)
				})
			})
		}
	})

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

	return api.app.OnAdminBeforeDeleteRequest().Trigger(event, func(e *core.AdminDeleteEvent) error {
		if err := api.app.Dao().DeleteAdmin(e.Admin); err != nil {
			return NewBadRequestError("Failed to delete admin.", err)
		}

		return api.app.OnAdminAfterDeleteRequest().Trigger(event, func(e *core.AdminDeleteEvent) error {
			if e.HttpContext.Response().Committed {
				return nil
			}

			return e.HttpContext.NoContent(http.StatusNoContent)
		})
	})
}
