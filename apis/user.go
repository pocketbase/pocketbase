package apis

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
	"golang.org/x/oauth2"
)

// BindUserApi registers the user api endpoints and the corresponding handlers.
func BindUserApi(app core.App, rg *echo.Group) {
	api := userApi{app: app}

	subGroup := rg.Group("/users", ActivityLogger(app))
	subGroup.GET("/auth-methods", api.authMethods)
	subGroup.POST("/auth-via-oauth2", api.oauth2Auth, RequireGuestOnly())
	subGroup.POST("/auth-via-email", api.emailAuth, RequireGuestOnly())
	subGroup.POST("/request-password-reset", api.requestPasswordReset)
	subGroup.POST("/confirm-password-reset", api.confirmPasswordReset)
	subGroup.POST("/request-verification", api.requestVerification)
	subGroup.POST("/confirm-verification", api.confirmVerification)
	subGroup.POST("/request-email-change", api.requestEmailChange, RequireUserAuth())
	subGroup.POST("/confirm-email-change", api.confirmEmailChange)
	subGroup.POST("/refresh", api.refresh, RequireUserAuth())
	// crud
	subGroup.GET("", api.list, RequireAdminAuth())
	subGroup.POST("", api.create)
	subGroup.GET("/:id", api.view, RequireAdminOrOwnerAuth("id"))
	subGroup.PATCH("/:id", api.update, RequireAdminAuth())
	subGroup.DELETE("/:id", api.delete, RequireAdminOrOwnerAuth("id"))
	subGroup.GET("/:id/external-auths", api.listExternalAuths, RequireAdminOrOwnerAuth("id"))
	subGroup.DELETE("/:id/external-auths/:provider", api.unlinkExternalAuth, RequireAdminOrOwnerAuth("id"))
}

type userApi struct {
	app core.App
}

func (api *userApi) authResponse(c echo.Context, user *models.User, meta any) error {
	token, tokenErr := tokens.NewUserAuthToken(api.app, user)
	if tokenErr != nil {
		return rest.NewBadRequestError("Failed to create auth token.", tokenErr)
	}

	event := &core.UserAuthEvent{
		HttpContext: c,
		User:        user,
		Token:       token,
		Meta:        meta,
	}

	return api.app.OnUserAuthRequest().Trigger(event, func(e *core.UserAuthEvent) error {
		result := map[string]any{
			"token": e.Token,
			"user":  e.User,
		}

		if e.Meta != nil {
			result["meta"] = e.Meta
		}

		return e.HttpContext.JSON(http.StatusOK, result)
	})
}

func (api *userApi) refresh(c echo.Context) error {
	user, _ := c.Get(ContextUserKey).(*models.User)
	if user == nil {
		return rest.NewNotFoundError("Missing auth user context.", nil)
	}

	return api.authResponse(c, user, nil)
}

type providerInfo struct {
	Name                string `json:"name"`
	State               string `json:"state"`
	CodeVerifier        string `json:"codeVerifier"`
	CodeChallenge       string `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
	AuthUrl             string `json:"authUrl"`
}

func (api *userApi) authMethods(c echo.Context) error {
	result := struct {
		EmailPassword bool           `json:"emailPassword"`
		AuthProviders []providerInfo `json:"authProviders"`
	}{
		EmailPassword: true,
		AuthProviders: []providerInfo{},
	}

	settings := api.app.Settings()

	result.EmailPassword = settings.EmailAuth.Enabled

	nameConfigMap := settings.NamedAuthProviderConfigs()

	for name, config := range nameConfigMap {
		if !config.Enabled {
			continue
		}

		provider, err := auth.NewProviderByName(name)
		if err != nil {
			if api.app.IsDebug() {
				log.Println(err)
			}

			// skip provider
			continue
		}

		if err := config.SetupProvider(provider); err != nil {
			if api.app.IsDebug() {
				log.Println(err)
			}

			// skip provider
			continue
		}

		state := security.RandomString(30)
		codeVerifier := security.RandomString(43)
		codeChallenge := security.S256Challenge(codeVerifier)
		codeChallengeMethod := "S256"
		result.AuthProviders = append(result.AuthProviders, providerInfo{
			Name:                name,
			State:               state,
			CodeVerifier:        codeVerifier,
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: codeChallengeMethod,
			AuthUrl: provider.BuildAuthUrl(
				state,
				oauth2.SetAuthURLParam("code_challenge", codeChallenge),
				oauth2.SetAuthURLParam("code_challenge_method", codeChallengeMethod),
			) + "&redirect_uri=", // empty redirect_uri so that users can append their url
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (api *userApi) oauth2Auth(c echo.Context) error {
	form := forms.NewUserOauth2Login(api.app)
	if readErr := c.Bind(form); readErr != nil {
		return rest.NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	user, authData, submitErr := form.Submit()
	if submitErr != nil {
		return rest.NewBadRequestError("Failed to authenticate.", submitErr)
	}

	return api.authResponse(c, user, authData)
}

func (api *userApi) emailAuth(c echo.Context) error {
	if !api.app.Settings().EmailAuth.Enabled {
		return rest.NewBadRequestError("Email/Password authentication is not enabled.", nil)
	}

	form := forms.NewUserEmailLogin(api.app)
	if readErr := c.Bind(form); readErr != nil {
		return rest.NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	user, submitErr := form.Submit()
	if submitErr != nil {
		return rest.NewBadRequestError("Failed to authenticate.", submitErr)
	}

	return api.authResponse(c, user, nil)
}

func (api *userApi) requestPasswordReset(c echo.Context) error {
	form := forms.NewUserPasswordResetRequest(api.app)
	if err := c.Bind(form); err != nil {
		return rest.NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	if err := form.Validate(); err != nil {
		return rest.NewBadRequestError("An error occurred while validating the form.", err)
	}

	// run in background because we don't need to show
	// the result to the user (prevents users enumeration)
	routine.FireAndForget(func() {
		if err := form.Submit(); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	})

	return c.NoContent(http.StatusNoContent)
}

func (api *userApi) confirmPasswordReset(c echo.Context) error {
	form := forms.NewUserPasswordResetConfirm(api.app)
	if readErr := c.Bind(form); readErr != nil {
		return rest.NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	user, submitErr := form.Submit()
	if submitErr != nil {
		return rest.NewBadRequestError("Failed to set new password.", submitErr)
	}

	return api.authResponse(c, user, nil)
}

func (api *userApi) requestEmailChange(c echo.Context) error {
	loggedUser, _ := c.Get(ContextUserKey).(*models.User)
	if loggedUser == nil {
		return rest.NewUnauthorizedError("The request requires valid authorized user.", nil)
	}

	form := forms.NewUserEmailChangeRequest(api.app, loggedUser)
	if err := c.Bind(form); err != nil {
		return rest.NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	if err := form.Submit(); err != nil {
		return rest.NewBadRequestError("Failed to request email change.", err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (api *userApi) confirmEmailChange(c echo.Context) error {
	form := forms.NewUserEmailChangeConfirm(api.app)
	if readErr := c.Bind(form); readErr != nil {
		return rest.NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	user, submitErr := form.Submit()
	if submitErr != nil {
		return rest.NewBadRequestError("Failed to confirm email change.", submitErr)
	}

	return api.authResponse(c, user, nil)
}

func (api *userApi) requestVerification(c echo.Context) error {
	form := forms.NewUserVerificationRequest(api.app)
	if err := c.Bind(form); err != nil {
		return rest.NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	if err := form.Validate(); err != nil {
		return rest.NewBadRequestError("An error occurred while validating the form.", err)
	}

	// run in background because we don't need to show
	// the result to the user (prevents users enumeration)
	routine.FireAndForget(func() {
		if err := form.Submit(); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	})

	return c.NoContent(http.StatusNoContent)
}

func (api *userApi) confirmVerification(c echo.Context) error {
	form := forms.NewUserVerificationConfirm(api.app)
	if readErr := c.Bind(form); readErr != nil {
		return rest.NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	user, submitErr := form.Submit()
	if submitErr != nil {
		return rest.NewBadRequestError("An error occurred while submitting the form.", submitErr)
	}

	return api.authResponse(c, user, nil)
}

// -------------------------------------------------------------------
// CRUD
// -------------------------------------------------------------------

func (api *userApi) list(c echo.Context) error {
	fieldResolver := search.NewSimpleFieldResolver(
		"id", "created", "updated", "email", "verified",
	)

	users := []*models.User{}

	result, searchErr := search.NewProvider(fieldResolver).
		Query(api.app.Dao().UserQuery()).
		ParseAndExec(c.QueryString(), &users)
	if searchErr != nil {
		return rest.NewBadRequestError("", searchErr)
	}

	// eager load user profiles (if any)
	if err := api.app.Dao().LoadProfiles(users); err != nil {
		return rest.NewBadRequestError("", err)
	}

	event := &core.UsersListEvent{
		HttpContext: c,
		Users:       users,
		Result:      result,
	}

	return api.app.OnUsersListRequest().Trigger(event, func(e *core.UsersListEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.Result)
	})
}

func (api *userApi) view(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return rest.NewNotFoundError("", nil)
	}

	user, err := api.app.Dao().FindUserById(id)
	if err != nil || user == nil {
		return rest.NewNotFoundError("", err)
	}

	event := &core.UserViewEvent{
		HttpContext: c,
		User:        user,
	}

	return api.app.OnUserViewRequest().Trigger(event, func(e *core.UserViewEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.User)
	})
}

func (api *userApi) create(c echo.Context) error {
	if !api.app.Settings().EmailAuth.Enabled {
		return rest.NewBadRequestError("Email/Password authentication is not enabled.", nil)
	}

	user := &models.User{}
	form := forms.NewUserUpsert(api.app, user)

	// load request
	if err := c.Bind(form); err != nil {
		return rest.NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.UserCreateEvent{
		HttpContext: c,
		User:        user,
	}

	// create the user
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnUserBeforeCreateRequest().Trigger(event, func(e *core.UserCreateEvent) error {
				if err := next(); err != nil {
					return rest.NewBadRequestError("Failed to create user.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.User)
			})
		}
	})

	if submitErr == nil {
		api.app.OnUserAfterCreateRequest().Trigger(event)
	}

	return submitErr
}

func (api *userApi) update(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return rest.NewNotFoundError("", nil)
	}

	user, err := api.app.Dao().FindUserById(id)
	if err != nil || user == nil {
		return rest.NewNotFoundError("", err)
	}

	form := forms.NewUserUpsert(api.app, user)

	// load request
	if err := c.Bind(form); err != nil {
		return rest.NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.UserUpdateEvent{
		HttpContext: c,
		User:        user,
	}

	// update the user
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnUserBeforeUpdateRequest().Trigger(event, func(e *core.UserUpdateEvent) error {
				if err := next(); err != nil {
					return rest.NewBadRequestError("Failed to update user.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.User)
			})
		}
	})

	if submitErr == nil {
		api.app.OnUserAfterUpdateRequest().Trigger(event)
	}

	return submitErr
}

func (api *userApi) delete(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return rest.NewNotFoundError("", nil)
	}

	user, err := api.app.Dao().FindUserById(id)
	if err != nil || user == nil {
		return rest.NewNotFoundError("", err)
	}

	event := &core.UserDeleteEvent{
		HttpContext: c,
		User:        user,
	}

	handlerErr := api.app.OnUserBeforeDeleteRequest().Trigger(event, func(e *core.UserDeleteEvent) error {
		// delete the user model
		if err := api.app.Dao().DeleteUser(e.User); err != nil {
			return rest.NewBadRequestError("Failed to delete user. Make sure that the user is not part of a required relation reference.", err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnUserAfterDeleteRequest().Trigger(event)
	}

	return handlerErr
}

func (api *userApi) listExternalAuths(c echo.Context) error {
	id := c.PathParam("id")
	if id == "" {
		return rest.NewNotFoundError("", nil)
	}

	user, err := api.app.Dao().FindUserById(id)
	if err != nil || user == nil {
		return rest.NewNotFoundError("", err)
	}

	externalAuths, err := api.app.Dao().FindAllExternalAuthsByUserId(user.Id)
	if err != nil {
		return rest.NewBadRequestError("Failed to fetch the external auths for the specified user.", err)
	}

	event := &core.UserListExternalAuthsEvent{
		HttpContext:   c,
		User:          user,
		ExternalAuths: externalAuths,
	}

	return api.app.OnUserListExternalAuths().Trigger(event, func(e *core.UserListExternalAuthsEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.ExternalAuths)
	})
}

func (api *userApi) unlinkExternalAuth(c echo.Context) error {
	id := c.PathParam("id")
	provider := c.PathParam("provider")
	if id == "" || provider == "" {
		return rest.NewNotFoundError("", nil)
	}

	user, err := api.app.Dao().FindUserById(id)
	if err != nil || user == nil {
		return rest.NewNotFoundError("", err)
	}

	externalAuth, err := api.app.Dao().FindExternalAuthByUserIdAndProvider(user.Id, provider)
	if err != nil {
		return rest.NewNotFoundError("Missing external auth provider relation.", err)
	}

	event := &core.UserUnlinkExternalAuthEvent{
		HttpContext:  c,
		User:         user,
		ExternalAuth: externalAuth,
	}

	handlerErr := api.app.OnUserBeforeUnlinkExternalAuthRequest().Trigger(event, func(e *core.UserUnlinkExternalAuthEvent) error {
		if err := api.app.Dao().DeleteExternalAuth(externalAuth); err != nil {
			return rest.NewBadRequestError("Cannot unlink the external auth provider. Make sure that the user has other linked auth providers OR has an email address.", err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnUserAfterUnlinkExternalAuthRequest().Trigger(event)
	}

	return handlerErr
}
