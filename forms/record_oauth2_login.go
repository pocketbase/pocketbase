package forms

import (
	"context"
	"errors"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/security"
	"golang.org/x/oauth2"
)

// RecordOAuth2LoginData defines the OA
type RecordOAuth2LoginData struct {
	ExternalAuth   *models.ExternalAuth
	Record         *models.Record
	OAuth2User     *auth.AuthUser
	ProviderClient auth.Provider
}

// BeforeOAuth2RecordCreateFunc defines a callback function that will
// be called before OAuth2 new Record creation.
type BeforeOAuth2RecordCreateFunc func(createForm *RecordUpsert, authRecord *models.Record, authUser *auth.AuthUser) error

// RecordOAuth2Login is an auth record OAuth2 login form.
type RecordOAuth2Login struct {
	app        core.App
	dao        *daos.Dao
	collection *models.Collection

	beforeOAuth2RecordCreateFunc BeforeOAuth2RecordCreateFunc

	// Optional auth record that will be used if no external
	// auth relation is found (if it is from the same collection)
	loggedAuthRecord *models.Record

	// The name of the OAuth2 client provider (eg. "google")
	Provider string `form:"provider" json:"provider"`

	// The authorization code returned from the initial request.
	Code string `form:"code" json:"code"`

	// The optional PKCE code verifier as part of the code_challenge sent with the initial request.
	CodeVerifier string `form:"codeVerifier" json:"codeVerifier"`

	// The redirect url sent with the initial request.
	RedirectUrl string `form:"redirectUrl" json:"redirectUrl"`

	// Additional data that will be used for creating a new auth record
	// if an existing OAuth2 account doesn't exist.
	CreateData map[string]any `form:"createData" json:"createData"`
}

// NewRecordOAuth2Login creates a new [RecordOAuth2Login] form with
// initialized with from the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewRecordOAuth2Login(app core.App, collection *models.Collection, optAuthRecord *models.Record) *RecordOAuth2Login {
	form := &RecordOAuth2Login{
		app:              app,
		dao:              app.Dao(),
		collection:       collection,
		loggedAuthRecord: optAuthRecord,
	}

	return form
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *RecordOAuth2Login) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// SetBeforeNewRecordCreateFunc sets a before OAuth2 record create callback handler.
func (form *RecordOAuth2Login) SetBeforeNewRecordCreateFunc(f BeforeOAuth2RecordCreateFunc) {
	form.beforeOAuth2RecordCreateFunc = f
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordOAuth2Login) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Provider, validation.Required, validation.By(form.checkProviderName)),
		validation.Field(&form.Code, validation.Required),
		validation.Field(&form.RedirectUrl, validation.Required),
	)
}

func (form *RecordOAuth2Login) checkProviderName(value any) error {
	name, _ := value.(string)

	config, ok := form.app.Settings().NamedAuthProviderConfigs()[name]
	if !ok || !config.Enabled {
		return validation.NewError("validation_invalid_provider", fmt.Sprintf("%q is missing or is not enabled.", name))
	}

	return nil
}

// Submit validates and submits the form.
//
// If an auth record doesn't exist, it will make an attempt to create it
// based on the fetched OAuth2 profile data via a local [RecordUpsert] form.
// You can intercept/modify the Record create form with [form.SetBeforeNewRecordCreateFunc()].
//
// You can also optionally provide a list of InterceptorFunc to
// further modify the form behavior before persisting it.
//
// On success returns the authorized record model and the fetched provider's data.
func (form *RecordOAuth2Login) Submit(
	interceptors ...InterceptorFunc[*RecordOAuth2LoginData],
) (*models.Record, *auth.AuthUser, error) {
	if err := form.Validate(); err != nil {
		return nil, nil, err
	}

	if !form.collection.AuthOptions().AllowOAuth2Auth {
		return nil, nil, errors.New("OAuth2 authentication is not allowed for the auth collection.")
	}

	provider, err := auth.NewProviderByName(form.Provider)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	provider.SetContext(ctx)

	// load provider configuration
	providerConfig := form.app.Settings().NamedAuthProviderConfigs()[form.Provider]
	if err := providerConfig.SetupProvider(provider); err != nil {
		return nil, nil, err
	}

	provider.SetRedirectUrl(form.RedirectUrl)

	var opts []oauth2.AuthCodeOption

	if provider.PKCE() {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", form.CodeVerifier))
	}

	// fetch token
	token, err := provider.FetchToken(form.Code, opts...)
	if err != nil {
		return nil, nil, err
	}

	// fetch external auth user
	authUser, err := provider.FetchAuthUser(token)
	if err != nil {
		return nil, nil, err
	}

	var authRecord *models.Record

	// check for existing relation with the auth record
	rel, _ := form.dao.FindFirstExternalAuthByExpr(dbx.HashExp{
		"collectionId": form.collection.Id,
		"provider":     form.Provider,
		"providerId":   authUser.Id,
	})
	switch {
	case rel != nil:
		authRecord, err = form.dao.FindRecordById(form.collection.Id, rel.RecordId)
		if err != nil {
			return nil, authUser, err
		}
	case form.loggedAuthRecord != nil && form.loggedAuthRecord.Collection().Id == form.collection.Id:
		// fallback to the logged auth record (if any)
		authRecord = form.loggedAuthRecord
	case authUser.Email != "":
		// look for an existing auth record by the external auth record's email
		authRecord, _ = form.dao.FindAuthRecordByEmail(form.collection.Id, authUser.Email)
	}

	interceptorData := &RecordOAuth2LoginData{
		ExternalAuth:   rel,
		Record:         authRecord,
		OAuth2User:     authUser,
		ProviderClient: provider,
	}

	interceptorsErr := runInterceptors(interceptorData, func(newData *RecordOAuth2LoginData) error {
		return form.submit(newData)
	}, interceptors...)

	if interceptorsErr != nil {
		return nil, interceptorData.OAuth2User, interceptorsErr
	}

	return interceptorData.Record, interceptorData.OAuth2User, nil
}

func (form *RecordOAuth2Login) submit(data *RecordOAuth2LoginData) error {
	return form.dao.RunInTransaction(func(txDao *daos.Dao) error {
		if data.Record == nil {
			data.Record = models.NewRecord(form.collection)
			data.Record.RefreshId()
			data.Record.MarkAsNew()
			createForm := NewRecordUpsert(form.app, data.Record)
			createForm.SetFullManageAccess(true)
			createForm.SetDao(txDao)
			if data.OAuth2User.Username != "" &&
				len(data.OAuth2User.Username) >= 3 &&
				len(data.OAuth2User.Username) <= 150 &&
				usernameRegex.MatchString(data.OAuth2User.Username) {
				createForm.Username = form.dao.SuggestUniqueAuthRecordUsername(
					form.collection.Id,
					data.OAuth2User.Username,
				)
			}

			// load custom data
			createForm.LoadData(form.CreateData)

			// load the OAuth2 user data
			createForm.Email = data.OAuth2User.Email
			createForm.Verified = true // mark as verified as long as it matches the OAuth2 data (even if the email is empty)

			// generate a random password if not explicitly set
			if createForm.Password == "" {
				createForm.Password = security.RandomString(30)
				createForm.PasswordConfirm = createForm.Password
			}

			if form.beforeOAuth2RecordCreateFunc != nil {
				if err := form.beforeOAuth2RecordCreateFunc(createForm, data.Record, data.OAuth2User); err != nil {
					return err
				}
			}

			// create the new auth record
			if err := createForm.Submit(); err != nil {
				return err
			}
		} else {
			isLoggedAuthRecord := form.loggedAuthRecord != nil &&
				form.loggedAuthRecord.Id == data.Record.Id &&
				form.loggedAuthRecord.Collection().Id == data.Record.Collection().Id

			// set random password for users with unverified email
			// (this is in case a malicious actor has registered via password using the user email)
			if !isLoggedAuthRecord && data.Record.Email() != "" && !data.Record.Verified() {
				data.Record.SetPassword(security.RandomString(30))
				if err := txDao.SaveRecord(data.Record); err != nil {
					return err
				}
			}

			// update the existing auth record empty email if the data.OAuth2User has one
			// (this is in case previously the auth record was created
			// with an OAuth2 provider that didn't return an email address)
			if data.Record.Email() == "" && data.OAuth2User.Email != "" {
				data.Record.SetEmail(data.OAuth2User.Email)
				if err := txDao.SaveRecord(data.Record); err != nil {
					return err
				}
			}

			// update the existing auth record verified state
			// (only if the auth record doesn't have an email or the auth record email match with the one in data.OAuth2User)
			if !data.Record.Verified() && (data.Record.Email() == "" || data.Record.Email() == data.OAuth2User.Email) {
				data.Record.SetVerified(true)
				if err := txDao.SaveRecord(data.Record); err != nil {
					return err
				}
			}
		}

		// create ExternalAuth relation if missing
		if data.ExternalAuth == nil {
			data.ExternalAuth = &models.ExternalAuth{
				CollectionId: data.Record.Collection().Id,
				RecordId:     data.Record.Id,
				Provider:     form.Provider,
				ProviderId:   data.OAuth2User.Id,
			}
			if err := txDao.SaveExternalAuth(data.ExternalAuth); err != nil {
				return err
			}
		}

		return nil
	})
}
