package forms

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/security"
	"golang.org/x/oauth2"
)

// RecordOAuth2Login is an auth record OAuth2 login form.
type RecordOAuth2Login struct {
	app        core.App
	dao        *daos.Dao
	collection *models.Collection

	// Optional auth record that will be used if no external
	// auth relation is found (if it is from the same collection)
	loggedAuthRecord *models.Record

	// The name of the OAuth2 client provider (eg. "google")
	Provider string `form:"provider" json:"provider"`

	// The authorization code returned from the initial request.
	Code string `form:"code" json:"code"`

	// The code verifier sent with the initial request as part of the code_challenge.
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

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *RecordOAuth2Login) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Provider, validation.Required, validation.By(form.checkProviderName)),
		validation.Field(&form.Code, validation.Required),
		validation.Field(&form.CodeVerifier, validation.Required),
		validation.Field(&form.RedirectUrl, validation.Required, is.URL),
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
// You can intercept/modify the create form by setting the optional beforeCreateFuncs argument.
//
// On success returns the authorized record model and the fetched provider's data.
func (form *RecordOAuth2Login) Submit(
	beforeCreateFuncs ...func(createForm *RecordUpsert, authRecord *models.Record, authUser *auth.AuthUser) error,
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

	// load provider configuration
	providerConfig := form.app.Settings().NamedAuthProviderConfigs()[form.Provider]
	if err := providerConfig.SetupProvider(provider); err != nil {
		return nil, nil, err
	}

	provider.SetRedirectUrl(form.RedirectUrl)

	// fetch token
	token, err := provider.FetchToken(
		form.Code,
		oauth2.SetAuthURLParam("code_verifier", form.CodeVerifier),
	)
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
	rel, _ := form.dao.FindExternalAuthByProvider(form.Provider, authUser.Id)
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

	saveErr := form.dao.RunInTransaction(func(txDao *daos.Dao) error {
		if authRecord == nil {
			authRecord = models.NewRecord(form.collection)
			authRecord.RefreshId()
			authRecord.MarkAsNew()
			createForm := NewRecordUpsert(form.app, authRecord)
			createForm.SetFullManageAccess(true)
			createForm.SetDao(txDao)
			if authUser.Username != "" && usernameRegex.MatchString(authUser.Username) {
				createForm.Username = form.dao.SuggestUniqueAuthRecordUsername(form.collection.Id, authUser.Username)
			}

			// load custom data
			createForm.LoadData(form.CreateData)

			// load the OAuth2 profile data as fallback
			if createForm.Email == "" {
				createForm.Email = authUser.Email
			}
			createForm.Verified = false
			if createForm.Email == authUser.Email {
				// mark as verified as long as it matches the OAuth2 data (even if the email is empty)
				createForm.Verified = true
			}
			if createForm.Password == "" {
				createForm.Password = security.RandomString(30)
				createForm.PasswordConfirm = createForm.Password
			}

			for _, f := range beforeCreateFuncs {
				if f == nil {
					continue
				}
				if err := f(createForm, authRecord, authUser); err != nil {
					return err
				}
			}

			// create the new auth record
			if err := createForm.Submit(); err != nil {
				return err
			}
		} else {
			// update the existing auth record empty email if the authUser has one
			// (this is in case previously the auth record was created
			// with an OAuth2 provider that didn't return an email address)
			if authRecord.Email() == "" && authUser.Email != "" {
				authRecord.SetEmail(authUser.Email)
				if err := txDao.SaveRecord(authRecord); err != nil {
					return err
				}
			}

			// update the existing auth record verified state
			// (only if the auth record doesn't have an email or the auth record email match with the one in authUser)
			if !authRecord.Verified() && (authRecord.Email() == "" || authRecord.Email() == authUser.Email) {
				authRecord.SetVerified(true)
				if err := txDao.SaveRecord(authRecord); err != nil {
					return err
				}
			}
		}

		// create ExternalAuth relation if missing
		if rel == nil {
			rel = &models.ExternalAuth{
				CollectionId: authRecord.Collection().Id,
				RecordId:     authRecord.Id,
				Provider:     form.Provider,
				ProviderId:   authUser.Id,
			}
			if err := txDao.SaveExternalAuth(rel); err != nil {
				return err
			}
		}

		return nil
	})

	if saveErr != nil {
		return nil, authUser, saveErr
	}

	return authRecord, authUser, nil
}
