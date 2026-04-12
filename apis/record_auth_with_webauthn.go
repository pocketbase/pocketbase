package apis

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

const (
	webauthnSessionPrefix = "webauthn:session:"
	webauthnSessionTTL    = 2 * time.Minute
)

// webauthnSessionEntry holds a WebAuthn challenge session along with its expiry.
type webauthnSessionEntry struct {
	Session   *webauthn.SessionData
	ExpiresAt time.Time
	RecordId  string
}

// webauthnUserAdapter wraps a PocketBase auth record to implement the
// webauthn.User interface required by the go-webauthn library.
type webauthnUserAdapter struct {
	record      *core.Record
	credentials []webauthn.Credential
}

func (u *webauthnUserAdapter) WebAuthnID() []byte {
	return []byte(u.record.Id)
}

func (u *webauthnUserAdapter) WebAuthnName() string {
	name := u.record.GetString("username")
	if name == "" {
		name = u.record.GetString("email")
	}
	if name == "" {
		name = u.record.Id
	}
	return name
}

func (u *webauthnUserAdapter) WebAuthnDisplayName() string {
	name := u.record.GetString("name")
	if name == "" {
		name = u.WebAuthnName()
	}
	return name
}

func (u *webauthnUserAdapter) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

// initWebAuthn creates a new webauthn.WebAuthn instance configured from app settings.
// The relying party ID and origin are derived from the application URL.
func initWebAuthn(app core.App) (*webauthn.WebAuthn, error) {
	appURL := app.Settings().Meta.AppURL
	if appURL == "" {
		return nil, errors.New("application URL is not configured")
	}

	parsed, err := url.Parse(appURL)
	if err != nil {
		return nil, fmt.Errorf("invalid application URL: %w", err)
	}

	rpID := parsed.Hostname()
	origin := strings.TrimRight(appURL, "/")

	config := &webauthn.Config{
		RPDisplayName: app.Settings().Meta.AppName,
		RPID:          rpID,
		RPOrigins:     []string{origin},
	}

	return webauthn.New(config)
}

// loadUserCredentials loads all WebAuthn credentials for a given record from the database.
func loadUserCredentials(app core.App, record *core.Record) ([]webauthn.Credential, error) {
	records, err := app.FindAllRecords(
		core.CollectionNameWebAuthnCredentials,
		dbx.HashExp{
			"collectionRef": record.Collection().Id,
			"recordRef":     record.Id,
		},
	)
	if err != nil {
		return nil, err
	}

	credentials := make([]webauthn.Credential, 0, len(records))
	for _, r := range records {
		proxy := &core.WebAuthnCredential{Record: r}
		cred, err := proxy.ToWebAuthnCredential()
		if err != nil {
			continue
		}
		credentials = append(credentials, cred)
	}
	return credentials, nil
}

// storeWebAuthnSession saves a webauthn session to the app store with
// a cryptographically random key and short TTL.
func storeWebAuthnSession(app core.App, session *webauthn.SessionData, recordId string) string {
	token := security.RandomString(32)
	app.Store().Set(webauthnSessionPrefix+token, &webauthnSessionEntry{
		Session:   session,
		ExpiresAt: time.Now().Add(webauthnSessionTTL),
		RecordId:  recordId,
	})
	return token
}

// retrieveWebAuthnSession retrieves and validates a session from the app store.
func retrieveWebAuthnSession(app core.App, token string) (*webauthnSessionEntry, error) {
	raw := app.Store().Get(webauthnSessionPrefix + token)
	if raw == nil {
		return nil, errors.New("missing or expired webauthn session")
	}

	entry, ok := raw.(*webauthnSessionEntry)
	if !ok || entry == nil {
		return nil, errors.New("invalid webauthn session data")
	}

	if time.Now().After(entry.ExpiresAt) {
		app.Store().Remove(webauthnSessionPrefix + token)
		return nil, errors.New("webauthn session has expired")
	}

	return entry, nil
}

// deleteWebAuthnSession removes a session from the app store.
func deleteWebAuthnSession(app core.App, token string) {
	app.Store().Remove(webauthnSessionPrefix + token)
}

// -------------------------------------------------------------------
// Registration flow
// -------------------------------------------------------------------

// recordWebAuthnRegisterBegin initiates a WebAuthn credential registration.
//
// Requires an authenticated user. Returns PublicKeyCredentialCreationOptions
// and a session token for the finish step.
func recordWebAuthnRegisterBegin(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.WebAuthn.Enabled {
		return e.ForbiddenError("The collection is not configured to allow WebAuthn authentication.", nil)
	}

	record := e.Auth
	if record == nil {
		return e.UnauthorizedError("Authentication is required.", nil)
	}

	wa, err := initWebAuthn(e.App)
	if err != nil {
		return e.InternalServerError("Failed to initialize WebAuthn.", err)
	}

	existingCreds, err := loadUserCredentials(e.App, record)
	if err != nil {
		return e.InternalServerError("Failed to load existing credentials.", err)
	}

	user := &webauthnUserAdapter{
		record:      record,
		credentials: existingCreds,
	}

	options, session, err := wa.BeginRegistration(user)
	if err != nil {
		return e.InternalServerError("Failed to begin WebAuthn registration.", err)
	}

	token := storeWebAuthnSession(e.App, session, record.Id)

	return e.JSON(http.StatusOK, map[string]any{
		"options":      options,
		"sessionToken": token,
	})
}

// recordWebAuthnRegisterFinish completes WebAuthn credential registration.
//
// Requires the session token from register-begin and the attestation
// response from the authenticator. Returns 201 with the new credential ID.
func recordWebAuthnRegisterFinish(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.WebAuthn.Enabled {
		return e.ForbiddenError("The collection is not configured to allow WebAuthn authentication.", nil)
	}

	record := e.Auth
	if record == nil {
		return e.UnauthorizedError("Authentication is required.", nil)
	}

	// Parse the session token and optional credential name from the form
	form := &webauthnRegisterFinishForm{}
	if err := e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err := form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	entry, err := retrieveWebAuthnSession(e.App, form.SessionToken)
	if err != nil {
		return e.BadRequestError("Invalid or expired registration session.", err)
	}
	defer deleteWebAuthnSession(e.App, form.SessionToken)

	if entry.RecordId != record.Id {
		return e.BadRequestError("Session does not match the authenticated user.", nil)
	}

	wa, err := initWebAuthn(e.App)
	if err != nil {
		return e.InternalServerError("Failed to initialize WebAuthn.", err)
	}

	existingCreds, err := loadUserCredentials(e.App, record)
	if err != nil {
		return e.InternalServerError("Failed to load existing credentials.", err)
	}

	user := &webauthnUserAdapter{
		record:      record,
		credentials: existingCreds,
	}

	cred, err := wa.FinishRegistration(user, *entry.Session, e.Request)
	if err != nil {
		return e.BadRequestError("Failed to verify WebAuthn registration.", err)
	}

	// Check for duplicate credential registration
	credIdEncoded := base64.RawURLEncoding.EncodeToString(cred.ID)
	_, findErr := e.App.FindFirstRecordByFilter(
		core.CollectionNameWebAuthnCredentials,
		"credentialId = {:credId}",
		dbx.Params{"credId": credIdEncoded},
	)
	if findErr == nil {
		return e.BadRequestError("This passkey is already registered.", nil)
	}

	// Store credential
	wcRecord := core.NewWebAuthnCredential(e.App)
	wcRecord.SetCollectionRef(collection.Id)
	wcRecord.SetRecordRef(record.Id)
	wcRecord.FromWebAuthnCredential(*cred)
	if form.Name != "" {
		wcRecord.SetName(form.Name)
	}

	if err := e.App.Save(wcRecord); err != nil {
		return e.InternalServerError("Failed to save WebAuthn credential.", err)
	}

	return e.JSON(http.StatusCreated, map[string]any{
		"id":   wcRecord.Id,
		"name": wcRecord.Name(),
	})
}

// -------------------------------------------------------------------
// Login flow
// -------------------------------------------------------------------

// recordWebAuthnLoginBegin initiates a WebAuthn authentication.
//
// Accepts an identity (email or username) and returns
// PublicKeyCredentialRequestOptions and a session token.
func recordWebAuthnLoginBegin(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.WebAuthn.Enabled {
		return e.ForbiddenError("The collection is not configured to allow WebAuthn authentication.", nil)
	}

	form := &webauthnLoginBeginForm{}
	if err := e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err := form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	e.Set(core.RequestEventKeyInfoContext, core.RequestInfoContextWebAuthn)

	// find the user record (generic error to prevent enumeration)
	record, err := e.App.FindAuthRecordByEmail(collection, form.Identity)
	if err != nil {
		// try by username if email lookup failed
		record, err = e.App.FindFirstRecordByFilter(
			collection,
			"username = {:identity}",
			dbx.Params{"identity": form.Identity},
		)
		if err != nil {
			return e.BadRequestError("Failed to authenticate.", errors.New("invalid identity"))
		}
	}

	wa, err := initWebAuthn(e.App)
	if err != nil {
		return e.InternalServerError("Failed to initialize WebAuthn.", err)
	}

	existingCreds, err := loadUserCredentials(e.App, record)
	if err != nil || len(existingCreds) == 0 {
		return e.BadRequestError("Failed to authenticate.", errors.New("no registered passkeys"))
	}

	user := &webauthnUserAdapter{
		record:      record,
		credentials: existingCreds,
	}

	options, session, err := wa.BeginLogin(user)
	if err != nil {
		return e.InternalServerError("Failed to begin WebAuthn login.", err)
	}

	token := storeWebAuthnSession(e.App, session, record.Id)

	return e.JSON(http.StatusOK, map[string]any{
		"options":      options,
		"sessionToken": token,
	})
}

// recordWebAuthnLoginFinish completes WebAuthn authentication.
//
// Validates the assertion response against the stored credential
// and returns an auth token on success.
func recordWebAuthnLoginFinish(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.WebAuthn.Enabled {
		return e.ForbiddenError("The collection is not configured to allow WebAuthn authentication.", nil)
	}

	form := &webauthnLoginFinishForm{}
	if err := e.BindBody(form); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while loading the submitted data.", err))
	}
	if err := form.validate(); err != nil {
		return firstApiError(err, e.BadRequestError("An error occurred while validating the submitted data.", err))
	}

	e.Set(core.RequestEventKeyInfoContext, core.RequestInfoContextWebAuthn)

	// find the user record
	record, err := e.App.FindAuthRecordByEmail(collection, form.Identity)
	if err != nil {
		record, err = e.App.FindFirstRecordByFilter(
			collection,
			"username = {:identity}",
			dbx.Params{"identity": form.Identity},
		)
		if err != nil {
			return e.BadRequestError("Failed to authenticate.", errors.New("invalid identity"))
		}
	}

	// rate limit per-record (same threshold as OTP: 5 attempts per 180 seconds)
	err = checkRateLimit(e, "@pb_webauthn_"+record.Id, core.RateLimitRule{MaxRequests: 5, Duration: 180})
	if err != nil {
		return e.TooManyRequestsError("Too many attempts, please try again later.", nil)
	}

	entry, err := retrieveWebAuthnSession(e.App, form.SessionToken)
	if err != nil {
		return e.BadRequestError("Invalid or expired login session.", err)
	}
	defer deleteWebAuthnSession(e.App, form.SessionToken)

	if entry.RecordId != record.Id {
		return e.BadRequestError("Failed to authenticate.", errors.New("session mismatch"))
	}

	wa, err := initWebAuthn(e.App)
	if err != nil {
		return e.InternalServerError("Failed to initialize WebAuthn.", err)
	}

	existingCreds, err := loadUserCredentials(e.App, record)
	if err != nil {
		return e.InternalServerError("Failed to load credentials.", err)
	}

	user := &webauthnUserAdapter{
		record:      record,
		credentials: existingCreds,
	}

	cred, err := wa.FinishLogin(user, *entry.Session, e.Request)
	if err != nil {
		return e.BadRequestError("Failed to authenticate.", err)
	}

	// Update sign count on the stored credential (single query, reused for event)
	var storedProxy *core.WebAuthnCredential
	if cred != nil {
		credIdEncoded := base64.RawURLEncoding.EncodeToString(cred.ID)
		storedRecord, findErr := e.App.FindFirstRecordByFilter(
			core.CollectionNameWebAuthnCredentials,
			"credentialId = {:credId} && recordRef = {:recordRef}",
			dbx.Params{"credId": credIdEncoded, "recordRef": record.Id},
		)
		if findErr == nil {
			storedProxy = &core.WebAuthnCredential{Record: storedRecord}
			storedProxy.SetSignCount(cred.Authenticator.SignCount)
			storedProxy.SetFlags(cred.Flags)
			if saveErr := e.App.Save(storedProxy); saveErr != nil {
				e.App.Logger().Error("Failed to update WebAuthn credential sign count",
					"error", saveErr,
					"credentialId", credIdEncoded,
				)
			}
		}
	}

	// Trigger WebAuthn auth event
	event := new(core.RecordAuthWithWebAuthnRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Record = record
	event.Credential = storedProxy

	return e.App.OnRecordAuthWithWebAuthnRequest().Trigger(event, func(e *core.RecordAuthWithWebAuthnRequestEvent) error {
		return RecordAuthResponse(e.RequestEvent, e.Record, core.MFAMethodWebAuthn, nil)
	})
}

// -------------------------------------------------------------------
// Forms
// -------------------------------------------------------------------

// webauthnRegisterFinishForm captures the session token and optional credential
// name for the WebAuthn registration completion step.
type webauthnRegisterFinishForm struct {
	SessionToken string `form:"sessionToken" json:"sessionToken"`
	Name         string `form:"name" json:"name"`
}

func (form *webauthnRegisterFinishForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.SessionToken, validation.Required, validation.Length(1, 255)),
		validation.Field(&form.Name, validation.Length(0, 255)),
	)
}

// webauthnLoginBeginForm captures the identity (email or username) for
// WebAuthn login initiation.
type webauthnLoginBeginForm struct {
	Identity string `form:"identity" json:"identity"`
}

func (form *webauthnLoginBeginForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Identity, validation.Required, validation.Length(1, 255)),
	)
}

// webauthnLoginFinishForm captures the identity and session token for
// WebAuthn login completion.
type webauthnLoginFinishForm struct {
	Identity     string `form:"identity" json:"identity"`
	SessionToken string `form:"sessionToken" json:"sessionToken"`
}

func (form *webauthnLoginFinishForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Identity, validation.Required, validation.Length(1, 255)),
		validation.Field(&form.SessionToken, validation.Required, validation.Length(1, 255)),
	)
}

// -------------------------------------------------------------------
// Credential management (recovery support)
// -------------------------------------------------------------------

// recordWebAuthnListCredentials lists the authenticated user's WebAuthn credentials.
//
// Returns a summary of each credential (id, name, created, signCount)
// without exposing sensitive fields like public keys.
func recordWebAuthnListCredentials(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.WebAuthn.Enabled {
		return e.ForbiddenError("The collection is not configured to allow WebAuthn authentication.", nil)
	}

	record := e.Auth
	if record == nil {
		return e.UnauthorizedError("Authentication is required.", nil)
	}

	records, err := e.App.FindAllRecords(
		core.CollectionNameWebAuthnCredentials,
		dbx.HashExp{
			"collectionRef": collection.Id,
			"recordRef":     record.Id,
		},
	)
	if err != nil {
		return e.InternalServerError("Failed to load credentials.", err)
	}

	result := make([]map[string]any, 0, len(records))
	for _, r := range records {
		proxy := &core.WebAuthnCredential{Record: r}
		result = append(result, map[string]any{
			"id":        r.Id,
			"name":      proxy.Name(),
			"created":   r.GetDateTime("created"),
			"signCount": proxy.SignCount(),
		})
	}

	return e.JSON(http.StatusOK, result)
}

// recordWebAuthnDeleteCredential lets an authenticated user delete
// one of their own WebAuthn credentials by ID.
func recordWebAuthnDeleteCredential(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	if !collection.WebAuthn.Enabled {
		return e.ForbiddenError("The collection is not configured to allow WebAuthn authentication.", nil)
	}

	record := e.Auth
	if record == nil {
		return e.UnauthorizedError("Authentication is required.", nil)
	}

	credentialId := e.Request.PathValue("credentialId")
	if credentialId == "" {
		return e.BadRequestError("Missing credential ID.", nil)
	}

	credRecord, err := e.App.FindRecordById(core.CollectionNameWebAuthnCredentials, credentialId)
	if err != nil {
		return e.NotFoundError("Credential not found.", err)
	}

	proxy := &core.WebAuthnCredential{Record: credRecord}

	// Verify the credential belongs to the authenticated user
	if proxy.RecordRef() != record.Id || proxy.CollectionRef() != collection.Id {
		return e.NotFoundError("Credential not found.", nil)
	}

	if err := e.App.Delete(credRecord); err != nil {
		return e.InternalServerError("Failed to delete credential.", err)
	}

	return e.NoContent(http.StatusNoContent)
}

// recordWebAuthnAdminClearCredentials allows a superuser to delete
// all WebAuthn credentials for a specific auth record.
//
// This is the recovery mechanism for users who have lost their authenticator
// device and cannot log in with their passkey.
func recordWebAuthnAdminClearCredentials(e *core.RequestEvent) error {
	collection, err := findAuthCollection(e)
	if err != nil {
		return err
	}

	recordId := e.Request.PathValue("id")
	if recordId == "" {
		return e.BadRequestError("Missing record ID.", nil)
	}

	// Verify the target record exists
	_, err = e.App.FindRecordById(collection, recordId)
	if err != nil {
		return e.NotFoundError("Auth record not found.", err)
	}

	// Find and delete all WebAuthn credentials for this user
	records, err := e.App.FindAllRecords(
		core.CollectionNameWebAuthnCredentials,
		dbx.HashExp{
			"collectionRef": collection.Id,
			"recordRef":     recordId,
		},
	)
	if err != nil {
		return e.InternalServerError("Failed to load credentials.", err)
	}

	for _, r := range records {
		if err := e.App.Delete(r); err != nil {
			return e.InternalServerError("Failed to delete credential.", err)
		}
	}

	return e.JSON(http.StatusOK, map[string]any{
		"deleted": len(records),
	})
}
