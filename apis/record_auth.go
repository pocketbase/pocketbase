package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"golang.org/x/oauth2"
)

// bindRecordAuthApi registers the auth record api endpoints and
// the corresponding handlers.
func bindRecordAuthApi(app core.App, rg *echo.Group) {
	api := recordAuthApi{app: app}

	// global oauth2 subscription redirect handler
	rg.GET("/oauth2-redirect", api.oauth2SubscriptionRedirect)

	// common collection record related routes
	subGroup := rg.Group(
		"/collections/:collection",
		ActivityLogger(app),
		LoadCollectionContext(app, models.CollectionTypeAuth),
	)
	subGroup.GET("/auth-methods", api.authMethods)
	subGroup.POST("/auth-refresh", api.authRefresh, RequireSameContextRecordAuth())
	subGroup.POST("/auth-with-oauth2", api.authWithOAuth2)
	subGroup.POST("/auth-with-password", api.authWithPassword)
	subGroup.POST("/request-password-reset", api.requestPasswordReset)
	subGroup.POST("/confirm-password-reset", api.confirmPasswordReset)
	subGroup.POST("/request-verification", api.requestVerification)
	subGroup.POST("/confirm-verification", api.confirmVerification)
	subGroup.POST("/request-email-change", api.requestEmailChange, RequireSameContextRecordAuth())
	subGroup.POST("/confirm-email-change", api.confirmEmailChange)
	subGroup.GET("/records/:id/external-auths", api.listExternalAuths, RequireAdminOrOwnerAuth("id"))
	subGroup.DELETE("/records/:id/external-auths/:provider", api.unlinkExternalAuth, RequireAdminOrOwnerAuth("id"))
}

type recordAuthApi struct {
	app core.App
}

func (api *recordAuthApi) authRefresh(c echo.Context) error {
	record, _ := c.Get(ContextAuthRecordKey).(*models.Record)
	if record == nil {
		return NewNotFoundError("Missing auth record context.", nil)
	}

	event := new(core.RecordAuthRefreshEvent)
	event.HttpContext = c
	event.Collection = record.Collection()
	event.Record = record

	handlerErr := api.app.OnRecordBeforeAuthRefreshRequest().Trigger(event, func(e *core.RecordAuthRefreshEvent) error {
		return RecordAuthResponse(api.app, e.HttpContext, e.Record, nil)
	})

	if handlerErr == nil {
		if err := api.app.OnRecordAfterAuthRefreshRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return handlerErr
}

type providerInfo struct {
	Name                string `json:"name"`
	State               string `json:"state"`
	CodeVerifier        string `json:"codeVerifier"`
	CodeChallenge       string `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
	AuthUrl             string `json:"authUrl"`
}

func (api *recordAuthApi) authMethods(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	authOptions := collection.AuthOptions()

	result := struct {
		UsernamePassword bool           `json:"usernamePassword"`
		EmailPassword    bool           `json:"emailPassword"`
		AuthProviders    []providerInfo `json:"authProviders"`
	}{
		UsernamePassword: authOptions.AllowUsernameAuth,
		EmailPassword:    authOptions.AllowEmailAuth,
		AuthProviders:    []providerInfo{},
	}

	if !authOptions.AllowOAuth2Auth {
		return c.JSON(http.StatusOK, result)
	}

	nameConfigMap := api.app.Settings().NamedAuthProviderConfigs()
	for name, config := range nameConfigMap {
		if !config.Enabled {
			continue
		}

		provider, err := auth.NewProviderByName(name)
		if err != nil {
			if api.app.IsDebug() {
				log.Println(err)
			}
			continue // skip provider
		}

		if err := config.SetupProvider(provider); err != nil {
			if api.app.IsDebug() {
				log.Println(err)
			}
			continue // skip provider
		}

		state := security.RandomString(30)
		codeVerifier := security.RandomString(43)
		codeChallenge := security.S256Challenge(codeVerifier)
		codeChallengeMethod := "S256"
		urlOpts := []oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", codeChallengeMethod),
		}

		if name == auth.NameApple {
			// see https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_js/incorporating_sign_in_with_apple_into_other_platforms#3332113
			urlOpts = append(urlOpts, oauth2.SetAuthURLParam("response_mode", "query"))
		}

		result.AuthProviders = append(result.AuthProviders, providerInfo{
			Name:                name,
			State:               state,
			CodeVerifier:        codeVerifier,
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: codeChallengeMethod,
			AuthUrl: provider.BuildAuthUrl(
				state,
				urlOpts...,
			) + "&redirect_uri=", // empty redirect_uri so that users can append their redirect url
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (api *recordAuthApi) authWithOAuth2(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	if !collection.AuthOptions().AllowOAuth2Auth {
		return NewBadRequestError("The collection is not configured to allow OAuth2 authentication.", nil)
	}

	var fallbackAuthRecord *models.Record

	loggedAuthRecord, _ := c.Get(ContextAuthRecordKey).(*models.Record)
	if loggedAuthRecord != nil && loggedAuthRecord.Collection().Id == collection.Id {
		fallbackAuthRecord = loggedAuthRecord
	}

	form := forms.NewRecordOAuth2Login(api.app, collection, fallbackAuthRecord)
	if readErr := c.Bind(form); readErr != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	event := new(core.RecordAuthWithOAuth2Event)
	event.HttpContext = c
	event.Collection = collection
	event.ProviderName = form.Provider
	event.IsNewRecord = false

	form.SetBeforeNewRecordCreateFunc(func(createForm *forms.RecordUpsert, authRecord *models.Record, authUser *auth.AuthUser) error {
		return createForm.DrySubmit(func(txDao *daos.Dao) error {
			event.IsNewRecord = true
			// clone the current request data and assign the form create data as its body data
			requestData := *RequestData(c)
			requestData.Data = form.CreateData

			createRuleFunc := func(q *dbx.SelectQuery) error {
				admin, _ := c.Get(ContextAdminKey).(*models.Admin)
				if admin != nil {
					return nil // either admin or the rule is empty
				}

				if collection.CreateRule == nil {
					return errors.New("Only admins can create new accounts with OAuth2")
				}

				if *collection.CreateRule != "" {
					resolver := resolvers.NewRecordFieldResolver(txDao, collection, &requestData, true)
					expr, err := search.FilterData(*collection.CreateRule).BuildExpr(resolver)
					if err != nil {
						return err
					}
					resolver.UpdateQuery(q)
					q.AndWhere(expr)
				}

				return nil
			}

			if _, err := txDao.FindRecordById(collection.Id, createForm.Id, createRuleFunc); err != nil {
				return fmt.Errorf("Failed create rule constraint: %w", err)
			}

			return nil
		})
	})

	_, _, submitErr := form.Submit(func(next forms.InterceptorNextFunc[*forms.RecordOAuth2LoginData]) forms.InterceptorNextFunc[*forms.RecordOAuth2LoginData] {
		return func(data *forms.RecordOAuth2LoginData) error {
			event.Record = data.Record
			event.OAuth2User = data.OAuth2User
			event.ProviderClient = data.ProviderClient

			return api.app.OnRecordBeforeAuthWithOAuth2Request().Trigger(event, func(e *core.RecordAuthWithOAuth2Event) error {
				data.Record = e.Record
				data.OAuth2User = e.OAuth2User

				if err := next(data); err != nil {
					return NewBadRequestError("Failed to authenticate.", err)
				}

				e.Record = data.Record
				e.OAuth2User = data.OAuth2User

				meta := struct {
					*auth.AuthUser
					IsNew bool `json:"isNew"`
				}{
					AuthUser: e.OAuth2User,
					IsNew:    event.IsNewRecord,
				}

				return RecordAuthResponse(api.app, e.HttpContext, e.Record, meta)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnRecordAfterAuthWithOAuth2Request().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return submitErr
}

func (api *recordAuthApi) authWithPassword(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	form := forms.NewRecordPasswordLogin(api.app, collection)
	if readErr := c.Bind(form); readErr != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	event := new(core.RecordAuthWithPasswordEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Password = form.Password
	event.Identity = form.Identity

	_, submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			event.Record = record

			return api.app.OnRecordBeforeAuthWithPasswordRequest().Trigger(event, func(e *core.RecordAuthWithPasswordEvent) error {
				if err := next(e.Record); err != nil {
					return NewBadRequestError("Failed to authenticate.", err)
				}

				return RecordAuthResponse(api.app, e.HttpContext, e.Record, nil)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnRecordAfterAuthWithPasswordRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return submitErr
}

func (api *recordAuthApi) requestPasswordReset(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	authOptions := collection.AuthOptions()
	if !authOptions.AllowUsernameAuth && !authOptions.AllowEmailAuth {
		return NewBadRequestError("The collection is not configured to allow password authentication.", nil)
	}

	form := forms.NewRecordPasswordResetRequest(api.app, collection)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	if err := form.Validate(); err != nil {
		return NewBadRequestError("An error occurred while validating the form.", err)
	}

	event := new(core.RecordRequestPasswordResetEvent)
	event.HttpContext = c
	event.Collection = collection

	submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			event.Record = record

			return api.app.OnRecordBeforeRequestPasswordResetRequest().Trigger(event, func(e *core.RecordRequestPasswordResetEvent) error {
				// run in background because we don't need to show the result to the client
				routine.FireAndForget(func() {
					if err := next(e.Record); err != nil && api.app.IsDebug() {
						log.Println(err)
					}
				})

				return e.HttpContext.NoContent(http.StatusNoContent)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnRecordAfterRequestPasswordResetRequest().Trigger(event); err != nil && api.app.IsDebug() {
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

func (api *recordAuthApi) confirmPasswordReset(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	form := forms.NewRecordPasswordResetConfirm(api.app, collection)
	if readErr := c.Bind(form); readErr != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	event := new(core.RecordConfirmPasswordResetEvent)
	event.HttpContext = c
	event.Collection = collection

	_, submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			event.Record = record

			return api.app.OnRecordBeforeConfirmPasswordResetRequest().Trigger(event, func(e *core.RecordConfirmPasswordResetEvent) error {
				if err := next(e.Record); err != nil {
					return NewBadRequestError("Failed to set new password.", err)
				}

				return e.HttpContext.NoContent(http.StatusNoContent)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnRecordAfterConfirmPasswordResetRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return submitErr
}

func (api *recordAuthApi) requestVerification(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	form := forms.NewRecordVerificationRequest(api.app, collection)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	if err := form.Validate(); err != nil {
		return NewBadRequestError("An error occurred while validating the form.", err)
	}

	event := new(core.RecordRequestVerificationEvent)
	event.HttpContext = c
	event.Collection = collection

	submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			event.Record = record

			return api.app.OnRecordBeforeRequestVerificationRequest().Trigger(event, func(e *core.RecordRequestVerificationEvent) error {
				// run in background because we don't need to show the result to the client
				routine.FireAndForget(func() {
					if err := next(e.Record); err != nil && api.app.IsDebug() {
						log.Println(err)
					}
				})

				return e.HttpContext.NoContent(http.StatusNoContent)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnRecordAfterRequestVerificationRequest().Trigger(event); err != nil && api.app.IsDebug() {
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

func (api *recordAuthApi) confirmVerification(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	form := forms.NewRecordVerificationConfirm(api.app, collection)
	if readErr := c.Bind(form); readErr != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	event := new(core.RecordConfirmVerificationEvent)
	event.HttpContext = c
	event.Collection = collection

	_, submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			event.Record = record

			return api.app.OnRecordBeforeConfirmVerificationRequest().Trigger(event, func(e *core.RecordConfirmVerificationEvent) error {
				if err := next(e.Record); err != nil {
					return NewBadRequestError("An error occurred while submitting the form.", err)
				}

				return e.HttpContext.NoContent(http.StatusNoContent)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnRecordAfterConfirmVerificationRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return submitErr
}

func (api *recordAuthApi) requestEmailChange(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	record, _ := c.Get(ContextAuthRecordKey).(*models.Record)
	if record == nil {
		return NewUnauthorizedError("The request requires valid auth record.", nil)
	}

	form := forms.NewRecordEmailChangeRequest(api.app, record)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	event := new(core.RecordRequestEmailChangeEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Record = record

	submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			return api.app.OnRecordBeforeRequestEmailChangeRequest().Trigger(event, func(e *core.RecordRequestEmailChangeEvent) error {
				if err := next(e.Record); err != nil {
					return NewBadRequestError("Failed to request email change.", err)
				}

				return e.HttpContext.NoContent(http.StatusNoContent)
			})
		}
	})

	if submitErr == nil {
		api.app.OnRecordAfterRequestEmailChangeRequest().Trigger(event)
	}

	return submitErr
}

func (api *recordAuthApi) confirmEmailChange(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	form := forms.NewRecordEmailChangeConfirm(api.app, collection)
	if readErr := c.Bind(form); readErr != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", readErr)
	}

	event := new(core.RecordConfirmEmailChangeEvent)
	event.HttpContext = c
	event.Collection = collection

	_, submitErr := form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			event.Record = record

			return api.app.OnRecordBeforeConfirmEmailChangeRequest().Trigger(event, func(e *core.RecordConfirmEmailChangeEvent) error {
				if err := next(e.Record); err != nil {
					return NewBadRequestError("Failed to confirm email change.", err)
				}

				return e.HttpContext.NoContent(http.StatusNoContent)
			})
		}
	})

	if submitErr == nil {
		if err := api.app.OnRecordAfterConfirmEmailChangeRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return submitErr
}

func (api *recordAuthApi) listExternalAuths(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	id := c.PathParam("id")
	if id == "" {
		return NewNotFoundError("", nil)
	}

	record, err := api.app.Dao().FindRecordById(collection.Id, id)
	if err != nil || record == nil {
		return NewNotFoundError("", err)
	}

	externalAuths, err := api.app.Dao().FindAllExternalAuthsByRecord(record)
	if err != nil {
		return NewBadRequestError("Failed to fetch the external auths for the specified auth record.", err)
	}

	event := new(core.RecordListExternalAuthsEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Record = record
	event.ExternalAuths = externalAuths

	return api.app.OnRecordListExternalAuthsRequest().Trigger(event, func(e *core.RecordListExternalAuthsEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.ExternalAuths)
	})
}

func (api *recordAuthApi) unlinkExternalAuth(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	id := c.PathParam("id")
	provider := c.PathParam("provider")
	if id == "" || provider == "" {
		return NewNotFoundError("", nil)
	}

	record, err := api.app.Dao().FindRecordById(collection.Id, id)
	if err != nil || record == nil {
		return NewNotFoundError("", err)
	}

	externalAuth, err := api.app.Dao().FindExternalAuthByRecordAndProvider(record, provider)
	if err != nil {
		return NewNotFoundError("Missing external auth provider relation.", err)
	}

	event := new(core.RecordUnlinkExternalAuthEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Record = record
	event.ExternalAuth = externalAuth

	handlerErr := api.app.OnRecordBeforeUnlinkExternalAuthRequest().Trigger(event, func(e *core.RecordUnlinkExternalAuthEvent) error {
		if err := api.app.Dao().DeleteExternalAuth(externalAuth); err != nil {
			return NewBadRequestError("Cannot unlink the external auth provider.", err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnRecordAfterUnlinkExternalAuthRequest().Trigger(event)
	}

	return handlerErr
}

// -------------------------------------------------------------------

const oauth2SubscriptionTopic = "@oauth2"

func (api *recordAuthApi) oauth2SubscriptionRedirect(c echo.Context) error {
	state := c.QueryParam("state")
	code := c.QueryParam("code")

	if code == "" || state == "" {
		return NewBadRequestError("Invalid OAuth2 redirect parameters.", nil)
	}

	client, err := api.app.SubscriptionsBroker().ClientById(state)
	if err != nil || client.IsDiscarded() || !client.HasSubscription(oauth2SubscriptionTopic) {
		return NewNotFoundError("Missing or invalid OAuth2 subscription client.", err)
	}

	data := map[string]string{
		"state": state,
		"code":  code,
	}

	encodedData, err := json.Marshal(data)
	if err != nil {
		return NewBadRequestError("Failed to marshalize OAuth2 redirect data.", err)
	}

	msg := subscriptions.Message{
		Name: oauth2SubscriptionTopic,
		Data: string(encodedData),
	}

	client.Channel() <- msg

	return c.Redirect(http.StatusTemporaryRedirect, "../_/#/auth/oauth2-redirect")
}
