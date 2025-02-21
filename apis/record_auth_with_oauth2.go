package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"golang.org/x/oauth2"
)

func recordAuthWithOAuth2(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.OAuth2.Enabled {
		return e.ForbiddenError("The collection is not configured to allow OAuth2 authentication.", nil)
	}

	var fallbackAuthRecord *core.Record
	if e.Auth != nil && e.Auth.Collection().Id == collection.Id {
		fallbackAuthRecord = e.Auth
	}

	e.Set(core.RequestEventKeyInfoContext, core.RequestInfoContextOAuth2)

	form := new(recordOAuth2LoginForm)
	form.collection = collection
	if err = e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}

	if form.RedirectUrl != "" && form.RedirectURL == "" {
		e.App.Logger().Warn("[recordAuthWithOAuth2] redirectUrl body param is deprecated and will be removed in the future. Please replace it with redirectURL.")
		form.RedirectURL = form.RedirectUrl
	}

	if err = form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}

	// exchange token for OAuth2 user info and locate existing ExternalAuth rel
	// ---------------------------------------------------------------

	// load provider configuration
	providerConfig, ok := collection.OAuth2.GetProviderConfig(form.Provider)
	if !ok {
		return e.InternalServerError("Missing or invalid provider config.", nil)
	}

	provider, err := providerConfig.InitProvider()
	if err != nil {
		return firstApiError(err, e.InternalServerError("Failed to init provider "+form.Provider, err))
	}

	ctx, cancel := context.WithTimeout(e.Request.Context(), 30*time.Second)
	defer cancel()

	provider.SetContext(ctx)
	provider.SetRedirectURL(form.RedirectURL)

	var opts []oauth2.AuthCodeOption

	if provider.PKCE() {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", form.CodeVerifier))
	}

	// fetch token
	token, err := provider.FetchToken(form.Code, opts...)
	if err != nil {
		return firstApiError(err, e.BadRequestError("Failed to fetch OAuth2 token.", err))
	}

	// fetch external auth user
	authUser, err := provider.FetchAuthUser(token)
	if err != nil {
		return firstApiError(err, e.BadRequestError("Failed to fetch OAuth2 user.", err))
	}

	var authRecord *core.Record

	// check for existing relation with the auth collection
	externalAuthRel, err := e.App.FindFirstExternalAuthByExpr(dbx.HashExp{
		"collectionRef": form.collection.Id,
		"provider":      form.Provider,
		"providerId":    authUser.Id,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return e.InternalServerError("Failed OAuth2 relation check.", err)
	}

	switch {
	case err == nil && externalAuthRel != nil:
		authRecord, err = e.App.FindRecordById(form.collection, externalAuthRel.RecordRef())
		if err != nil {
			return err
		}
	case fallbackAuthRecord != nil && fallbackAuthRecord.Collection().Id == form.collection.Id:
		// fallback to the logged auth record (if any)
		authRecord = fallbackAuthRecord
	case authUser.Email != "":
		// look for an existing auth record by the external auth record's email
		authRecord, err = e.App.FindAuthRecordByEmail(form.collection.Id, authUser.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return e.InternalServerError("Failed OAuth2 auth record check.", err)
		}
	}

	// ---------------------------------------------------------------

	event := new(core.RecordAuthWithOAuth2RequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.ProviderName = form.Provider
	event.ProviderClient = provider
	event.OAuth2User = authUser
	event.CreateData = form.CreateData
	event.Record = authRecord
	event.IsNewRecord = authRecord == nil

	return e.App.OnRecordAuthWithOAuth2Request().Trigger(event, func(e *core.RecordAuthWithOAuth2RequestEvent) error {
		if err := oauth2Submit(e, externalAuthRel); err != nil {
			return firstApiError(err, e.BadRequestError("Failed to authenticate.", err))
		}

		// @todo revert back to struct after removing the custom auth.AuthUser marshalization
		meta := map[string]any{}
		rawOAuth2User, err := json.Marshal(e.OAuth2User)
		if err != nil {
			return err
		}
		err = json.Unmarshal(rawOAuth2User, &meta)
		if err != nil {
			return err
		}
		meta["isNew"] = e.IsNewRecord

		return RecordAuthResponse(e.RequestEvent, e.Record, core.MFAMethodOAuth2, meta)
	})
}

// -------------------------------------------------------------------

type recordOAuth2LoginForm struct {
	collection *core.Collection

	// Additional data that will be used for creating a new auth record
	// if an existing OAuth2 account doesn't exist.
	CreateData map[string]any `form:"createData" json:"createData"`

	// The name of the OAuth2 client provider (eg. "google")
	Provider string `form:"provider" json:"provider"`

	// The authorization code returned from the initial request.
	Code string `form:"code" json:"code"`

	// The optional PKCE code verifier as part of the code_challenge sent with the initial request.
	CodeVerifier string `form:"codeVerifier" json:"codeVerifier"`

	// The redirect url sent with the initial request.
	RedirectURL string `form:"redirectURL" json:"redirectURL"`

	// @todo
	// deprecated: use RedirectURL instead
	// RedirectUrl will be removed after dropping v0.22 support
	RedirectUrl string `form:"redirectUrl" json:"redirectUrl"`
}

func (form *recordOAuth2LoginForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Provider, validation.Required, validation.Length(0, 100), validation.By(form.checkProviderName)),
		validation.Field(&form.Code, validation.Required),
		validation.Field(&form.RedirectURL, validation.Required),
	)
}

func (form *recordOAuth2LoginForm) checkProviderName(value any) error {
	name, _ := value.(string)

	_, ok := form.collection.OAuth2.GetProviderConfig(name)
	if !ok {
		return validation.NewError("validation_invalid_provider", "Provider with name {{.name}} is missing or is not enabled.").
			SetParams(map[string]any{"name": name})
	}

	return nil
}

func oldCanAssignUsername(txApp core.App, collection *core.Collection, username string) bool {
	// ensure that username is unique
	index, hasUniqueue := dbutils.FindSingleColumnUniqueIndex(collection.Indexes, collection.OAuth2.MappedFields.Username)
	if hasUniqueue {
		var expr dbx.Expression
		if strings.EqualFold(index.Columns[0].Collate, "nocase") {
			// case-insensitive search
			expr = dbx.NewExp("username = {:username} COLLATE NOCASE", dbx.Params{"username": username})
		} else {
			expr = dbx.HashExp{"username": username}
		}

		var exists int
		_ = txApp.RecordQuery(collection).Select("(1)").AndWhere(expr).Limit(1).Row(&exists)
		if exists > 0 {
			return false
		}
	}

	// ensure that the value matches the pattern of the username field (if text)
	txtField, _ := collection.Fields.GetByName(collection.OAuth2.MappedFields.Username).(*core.TextField)

	return txtField != nil && txtField.ValidatePlainValue(username) == nil
}

func oauth2Submit(e *core.RecordAuthWithOAuth2RequestEvent, optExternalAuth *core.ExternalAuth) error {
	return e.App.RunInTransaction(func(txApp core.App) error {
		if e.Record == nil {
			// extra check to prevent creating a superuser record via
			// OAuth2 in case the method is used by another action
			if e.Collection.Name == core.CollectionNameSuperusers {
				return errors.New("superusers are not allowed to sign-up with OAuth2")
			}

			payload := maps.Clone(e.CreateData)
			if payload == nil {
				payload = map[string]any{}
			}

			// assign the OAuth2 user email only if the user hasn't submitted one
			// (ignore empty/invalid values for consistency with the OAuth2->existing user update flow)
			if v, _ := payload[core.FieldNameEmail].(string); v == "" {
				payload[core.FieldNameEmail] = e.OAuth2User.Email
			}

			// map known fields (unless the field was explicitly submitted as part of CreateData)
			if _, ok := payload[e.Collection.OAuth2.MappedFields.Id]; !ok && e.Collection.OAuth2.MappedFields.Id != "" {
				payload[e.Collection.OAuth2.MappedFields.Id] = e.OAuth2User.Id
			}
			if _, ok := payload[e.Collection.OAuth2.MappedFields.Name]; !ok && e.Collection.OAuth2.MappedFields.Name != "" {
				payload[e.Collection.OAuth2.MappedFields.Name] = e.OAuth2User.Name
			}
			if _, ok := payload[e.Collection.OAuth2.MappedFields.Username]; !ok &&
				// no explicit username payload value and existing OAuth2 mapping
				e.Collection.OAuth2.MappedFields.Username != "" &&
				// extra checks for backward compatibility with earlier versions
				oldCanAssignUsername(txApp, e.Collection, e.OAuth2User.Username) {
				payload[e.Collection.OAuth2.MappedFields.Username] = e.OAuth2User.Username
			}
			if _, ok := payload[e.Collection.OAuth2.MappedFields.AvatarURL]; !ok &&
				// no explicit avatar payload value and existing OAuth2 mapping
				e.Collection.OAuth2.MappedFields.AvatarURL != "" &&
				// non-empty OAuth2 avatar url
				e.OAuth2User.AvatarURL != "" {
				mappedField := e.Collection.Fields.GetByName(e.Collection.OAuth2.MappedFields.AvatarURL)
				if mappedField != nil && mappedField.Type() == core.FieldTypeFile {
					// download the avatar if the mapped field is a file
					avatarFile, err := func() (*filesystem.File, error) {
						ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
						defer cancel()
						return filesystem.NewFileFromURL(ctx, e.OAuth2User.AvatarURL)
					}()
					if err != nil {
						txApp.Logger().Warn("Failed to retrieve OAuth2 avatar", slog.String("error", err.Error()))
					} else {
						payload[e.Collection.OAuth2.MappedFields.AvatarURL] = avatarFile
					}
				} else {
					// otherwise - assign the url string
					payload[e.Collection.OAuth2.MappedFields.AvatarURL] = e.OAuth2User.AvatarURL
				}
			}

			createdRecord, err := sendOAuth2RecordCreateRequest(txApp, e, payload)
			if err != nil {
				return err
			}

			e.Record = createdRecord

			if e.Record.Email() == e.OAuth2User.Email && !e.Record.Verified() {
				// mark as verified as long as it matches the OAuth2 data (even if the email is empty)
				e.Record.SetVerified(true)
				if err := txApp.Save(e.Record); err != nil {
					return err
				}
			}
		} else {
			var needUpdate bool

			isLoggedAuthRecord := e.Auth != nil &&
				e.Auth.Id == e.Record.Id &&
				e.Auth.Collection().Id == e.Record.Collection().Id

			// set random password for users with unverified email
			// (this is in case a malicious actor has registered previously with the user email)
			if !isLoggedAuthRecord && e.Record.Email() != "" && !e.Record.Verified() {
				e.Record.SetRandomPassword()
				needUpdate = true
			}

			// update the existing auth record empty email if the data.OAuth2User has one
			// (this is in case previously the auth record was created
			// with an OAuth2 provider that didn't return an email address)
			if e.Record.Email() == "" && e.OAuth2User.Email != "" {
				e.Record.SetEmail(e.OAuth2User.Email)
				needUpdate = true
			}

			// update the existing auth record verified state
			// (only if the auth record doesn't have an email or the auth record email match with the one in data.OAuth2User)
			if !e.Record.Verified() && (e.Record.Email() == "" || e.Record.Email() == e.OAuth2User.Email) {
				e.Record.SetVerified(true)
				needUpdate = true
			}

			if needUpdate {
				if err := txApp.Save(e.Record); err != nil {
					return err
				}
			}
		}

		// create ExternalAuth relation if missing
		if optExternalAuth == nil {
			optExternalAuth = core.NewExternalAuth(txApp)
			optExternalAuth.SetCollectionRef(e.Record.Collection().Id)
			optExternalAuth.SetRecordRef(e.Record.Id)
			optExternalAuth.SetProvider(e.ProviderName)
			optExternalAuth.SetProviderId(e.OAuth2User.Id)

			if err := txApp.Save(optExternalAuth); err != nil {
				return fmt.Errorf("failed to save linked rel: %w", err)
			}
		}

		return nil
	})
}

func sendOAuth2RecordCreateRequest(txApp core.App, e *core.RecordAuthWithOAuth2RequestEvent, payload map[string]any) (*core.Record, error) {
	ir := &core.InternalRequest{
		Method: http.MethodPost,
		URL:    "/api/collections/" + e.Collection.Name + "/records",
		Body:   payload,
	}

	var createdRecord *core.Record
	response, err := processInternalRequest(txApp, e.RequestEvent, ir, core.RequestInfoContextOAuth2, func(data any) error {
		createdRecord, _ = data.(*core.Record)

		return nil
	})
	if err != nil {
		return nil, err
	}

	if response.Status != http.StatusOK || createdRecord == nil {
		return nil, errors.New("failed to create OAuth2 auth record")
	}

	return createdRecord, nil
}
