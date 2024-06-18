package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/mails"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/oauth2"
)

// bindRecordAuthApi registers the auth record api endpoints and
// the corresponding handlers.
func bindRecordAuthApi(app core.App, rg *echo.Group) {
	api := recordAuthApi{app: app}

	// global oauth2 subscription redirect handler
	rg.GET("/oauth2-redirect", api.oauth2SubscriptionRedirect)
	rg.POST("/oauth2-redirect", api.oauth2SubscriptionRedirect) // needed in case of response_mode=form_post

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

	return api.app.OnRecordBeforeAuthRefreshRequest().Trigger(event, func(e *core.RecordAuthRefreshEvent) error {
		return api.app.OnRecordAfterAuthRefreshRequest().Trigger(event, func(e *core.RecordAuthRefreshEvent) error {
			return RecordAuthResponse(api.app, e.HttpContext, e.Record, nil)
		})
	})
}

type providerInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	State       string `json:"state"`
	AuthUrl     string `json:"authUrl"`
	// technically could be omitted if the provider doesn't support PKCE,
	// but to avoid breaking existing typed clients we'll return them as empty string
	CodeVerifier        string `json:"codeVerifier"`
	CodeChallenge       string `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
}

func (api *recordAuthApi) authMethods(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("Missing collection context.", nil)
	}

	authOptions := collection.AuthOptions()

	result := struct {
		AuthProviders    []providerInfo `json:"authProviders"`
		UsernamePassword bool           `json:"usernamePassword"`
		EmailPassword    bool           `json:"emailPassword"`
		OnlyVerified     bool           `json:"onlyVerified"`
	}{
		UsernamePassword: authOptions.AllowUsernameAuth,
		EmailPassword:    authOptions.AllowEmailAuth,
		OnlyVerified:     authOptions.OnlyVerified,
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
			api.app.Logger().Debug("Missing or invalid provider name", slog.String("name", name))
			continue // skip provider
		}

		if err := config.SetupProvider(provider); err != nil {
			api.app.Logger().Debug(
				"Failed to setup provider",
				slog.String("name", name),
				slog.String("error", err.Error()),
			)
			continue // skip provider
		}

		info := providerInfo{
			Name:        name,
			DisplayName: provider.DisplayName(),
			State:       security.RandomString(30),
		}

		if info.DisplayName == "" {
			info.DisplayName = name
		}

		urlOpts := []oauth2.AuthCodeOption{}

		// custom providers url options
		switch name {
		case auth.NameApple:
			// see https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_js/incorporating_sign_in_with_apple_into_other_platforms#3332113
			urlOpts = append(urlOpts, oauth2.SetAuthURLParam("response_mode", "form_post"))
		}

		if provider.PKCE() {
			info.CodeVerifier = security.RandomString(43)
			info.CodeChallenge = security.S256Challenge(info.CodeVerifier)
			info.CodeChallengeMethod = "S256"
			urlOpts = append(urlOpts,
				oauth2.SetAuthURLParam("code_challenge", info.CodeChallenge),
				oauth2.SetAuthURLParam("code_challenge_method", info.CodeChallengeMethod),
			)
		}

		info.AuthUrl = provider.BuildAuthUrl(
			info.State,
			urlOpts...,
		) + "&redirect_uri=" // empty redirect_uri so that users can append their redirect url

		result.AuthProviders = append(result.AuthProviders, info)
	}

	// sort providers
	sort.SliceStable(result.AuthProviders, func(i, j int) bool {
		return result.AuthProviders[i].Name < result.AuthProviders[j].Name
	})

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

	form.SetBeforeNewRecordCreateFunc(func(createForm *forms.RecordUpsert, authRecord *models.Record, authUser *auth.AuthUser) error {
		return createForm.DrySubmit(func(txDao *daos.Dao) error {
			event.IsNewRecord = true

			// clone the current request data and assign the form create data as its body data
			requestInfo := *RequestInfo(c)
			requestInfo.Context = models.RequestInfoContextOAuth2
			requestInfo.Data = form.CreateData

			createRuleFunc := func(q *dbx.SelectQuery) error {
				admin, _ := c.Get(ContextAdminKey).(*models.Admin)
				if admin != nil {
					return nil // either admin or the rule is empty
				}

				if collection.CreateRule == nil {
					return errors.New("Only admins can create new accounts with OAuth2")
				}

				if *collection.CreateRule != "" {
					resolver := resolvers.NewRecordFieldResolver(txDao, collection, &requestInfo, true)
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
			event.IsNewRecord = data.Record == nil

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

				return api.app.OnRecordAfterAuthWithOAuth2Request().Trigger(event, func(e *core.RecordAuthWithOAuth2Event) error {
					// clear the lastLoginAlertSentAt field so that we can enforce password auth notifications
					if !e.Record.LastLoginAlertSentAt().IsZero() {
						e.Record.Set(schema.FieldNameLastLoginAlertSentAt, "")
						if err := api.app.Dao().SaveRecord(e.Record); err != nil {
							api.app.Logger().Warn("Failed to reset lastLoginAlertSentAt", "error", err, "recordId", e.Record.Id)
						}
					}

					return RecordAuthResponse(api.app, e.HttpContext, e.Record, meta)
				})
			})
		}
	})

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

				// @todo remove after the refactoring
				if collection.AuthOptions().AllowOAuth2Auth && e.Record.Email() != "" {
					externalAuths, err := api.app.Dao().FindAllExternalAuthsByRecord(e.Record)
					if err != nil {
						return NewBadRequestError("Failed to authenticate.", err)
					}
					if len(externalAuths) > 0 {
						lastLoginAlert := e.Record.LastLoginAlertSentAt().Time()

						// send an email alert if the password auth is after OAuth2 auth (lastLoginAlert will be empty)
						// or if it has been ~7 days since the last alert
						if lastLoginAlert.IsZero() || time.Now().UTC().Sub(lastLoginAlert).Hours() > 168 {
							providerNames := make([]string, len(externalAuths))
							for i, ea := range externalAuths {
								var name string
								if provider, err := auth.NewProviderByName(ea.Provider); err == nil {
									name = provider.DisplayName()
								}
								if name == "" {
									name = ea.Provider
								}
								providerNames[i] = name
							}

							if err := mails.SendRecordPasswordLoginAlert(api.app, e.Record, providerNames...); err != nil {
								return NewBadRequestError("Failed to authenticate.", err)
							}

							e.Record.SetLastLoginAlertSentAt(types.NowDateTime())
							if err := api.app.Dao().SaveRecord(e.Record); err != nil {
								api.app.Logger().Warn("Failed to update lastLoginAlertSentAt", "error", err, "recordId", e.Record.Id)
							}
						}
					}
				}

				return api.app.OnRecordAfterAuthWithPasswordRequest().Trigger(event, func(e *core.RecordAuthWithPasswordEvent) error {
					return RecordAuthResponse(api.app, e.HttpContext, e.Record, nil)
				})
			})
		}
	})

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
					if err := next(e.Record); err != nil {
						api.app.Logger().Debug(
							"Failed to send password reset email",
							slog.String("error", err.Error()),
						)
					}
				})

				return api.app.OnRecordAfterRequestPasswordResetRequest().Trigger(event, func(e *core.RecordRequestPasswordResetEvent) error {
					if e.HttpContext.Response().Committed {
						return nil
					}

					return e.HttpContext.NoContent(http.StatusNoContent)
				})
			})
		}
	})

	// eagerly write 204 response and skip submit errors
	// as a measure against emails enumeration
	if !c.Response().Committed {
		c.NoContent(http.StatusNoContent)
	}

	return submitErr
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

				return api.app.OnRecordAfterConfirmPasswordResetRequest().Trigger(event, func(e *core.RecordConfirmPasswordResetEvent) error {
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
					if err := next(e.Record); err != nil {
						api.app.Logger().Debug(
							"Failed to send verification email",
							slog.String("error", err.Error()),
						)
					}
				})

				return api.app.OnRecordAfterRequestVerificationRequest().Trigger(event, func(e *core.RecordRequestVerificationEvent) error {
					if e.HttpContext.Response().Committed {
						return nil
					}

					return e.HttpContext.NoContent(http.StatusNoContent)
				})
			})
		}
	})

	// eagerly write 204 response and skip submit errors
	// as a measure against users enumeration
	if !c.Response().Committed {
		c.NoContent(http.StatusNoContent)
	}

	return submitErr
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

				return api.app.OnRecordAfterConfirmVerificationRequest().Trigger(event, func(e *core.RecordConfirmVerificationEvent) error {
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

	return form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(record *models.Record) error {
			return api.app.OnRecordBeforeRequestEmailChangeRequest().Trigger(event, func(e *core.RecordRequestEmailChangeEvent) error {
				if err := next(e.Record); err != nil {
					return NewBadRequestError("Failed to request email change.", err)
				}

				return api.app.OnRecordAfterRequestEmailChangeRequest().Trigger(event, func(e *core.RecordRequestEmailChangeEvent) error {
					if e.HttpContext.Response().Committed {
						return nil
					}

					return e.HttpContext.NoContent(http.StatusNoContent)
				})
			})
		}
	})
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

				return api.app.OnRecordAfterConfirmEmailChangeRequest().Trigger(event, func(e *core.RecordConfirmEmailChangeEvent) error {
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

	return api.app.OnRecordBeforeUnlinkExternalAuthRequest().Trigger(event, func(e *core.RecordUnlinkExternalAuthEvent) error {
		if err := api.app.Dao().DeleteExternalAuth(externalAuth); err != nil {
			return NewBadRequestError("Cannot unlink the external auth provider.", err)
		}

		return api.app.OnRecordAfterUnlinkExternalAuthRequest().Trigger(event, func(e *core.RecordUnlinkExternalAuthEvent) error {
			if e.HttpContext.Response().Committed {
				return nil
			}

			return e.HttpContext.NoContent(http.StatusNoContent)
		})
	})
}

// -------------------------------------------------------------------

const (
	oauth2SubscriptionTopic   string = "@oauth2"
	oauth2RedirectFailurePath string = "../_/#/auth/oauth2-redirect-failure"
	oauth2RedirectSuccessPath string = "../_/#/auth/oauth2-redirect-success"
)

type oauth2RedirectData struct {
	State string `form:"state" query:"state" json:"state"`
	Code  string `form:"code" query:"code" json:"code"`
	Error string `form:"error" query:"error" json:"error,omitempty"`
}

func (api *recordAuthApi) oauth2SubscriptionRedirect(c echo.Context) error {
	redirectStatusCode := http.StatusTemporaryRedirect
	if c.Request().Method != http.MethodGet {
		redirectStatusCode = http.StatusSeeOther
	}

	data := oauth2RedirectData{}
	if err := c.Bind(&data); err != nil {
		api.app.Logger().Debug("Failed to read OAuth2 redirect data", "error", err)
		return c.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
	}

	if data.State == "" {
		api.app.Logger().Debug("Missing OAuth2 state parameter")
		return c.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
	}

	client, err := api.app.SubscriptionsBroker().ClientById(data.State)
	if err != nil || client.IsDiscarded() || !client.HasSubscription(oauth2SubscriptionTopic) {
		api.app.Logger().Debug("Missing or invalid OAuth2 subscription client", "error", err, "clientId", data.State)
		return c.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
	}
	defer client.Unsubscribe(oauth2SubscriptionTopic)

	encodedData, err := json.Marshal(data)
	if err != nil {
		api.app.Logger().Debug("Failed to marshalize OAuth2 redirect data", "error", err)
		return c.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
	}

	msg := subscriptions.Message{
		Name: oauth2SubscriptionTopic,
		Data: encodedData,
	}

	client.Send(msg)

	if data.Error != "" || data.Code == "" {
		api.app.Logger().Debug("Failed OAuth2 redirect due to an error or missing code parameter", "error", data.Error, "clientId", data.State)
		return c.Redirect(redirectStatusCode, oauth2RedirectFailurePath)
	}

	return c.Redirect(redirectStatusCode, oauth2RedirectSuccessPath)
}
