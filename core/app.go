// Package core is the backbone of PocketBase.
//
// It defines the main PocketBase App interface and its base implementation.
package core

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

// App defines the main PocketBase app interface.
type App interface {
	// Deprecated:
	// This method may get removed in the near future.
	// It is recommended to access the app db instance from app.Dao().DB() or
	// if you want more flexibility - app.Dao().ConcurrentDB() and app.Dao().NonconcurrentDB().
	//
	// DB returns the default app database instance.
	DB() *dbx.DB

	// Dao returns the default app Dao instance.
	//
	// This Dao could operate only on the tables and models
	// associated with the default app database. For example,
	// trying to access the request logs table will result in error.
	Dao() *daos.Dao

	// Deprecated:
	// This method may get removed in the near future.
	// It is recommended to access the logs db instance from app.LogsDao().DB() or
	// if you want more flexibility - app.LogsDao().ConcurrentDB() and app.LogsDao().NonconcurrentDB().
	//
	// LogsDB returns the app logs database instance.
	LogsDB() *dbx.DB

	// LogsDao returns the app logs Dao instance.
	//
	// This Dao could operate only on the tables and models
	// associated with the logs database. For example, trying to access
	// the users table from LogsDao will result in error.
	LogsDao() *daos.Dao

	// DataDir returns the app data directory path.
	DataDir() string

	// EncryptionEnv returns the name of the app secret env key
	// (used for settings encryption).
	EncryptionEnv() string

	// IsDebug returns whether the app is in debug mode
	// (showing more detailed error logs, executed sql statements, etc.).
	IsDebug() bool

	// Settings returns the loaded app settings.
	Settings() *settings.Settings

	// Cache returns the app internal cache store.
	Cache() *store.Store[any]

	// SubscriptionsBroker returns the app realtime subscriptions broker instance.
	SubscriptionsBroker() *subscriptions.Broker

	// NewMailClient creates and returns a configured app mail client.
	NewMailClient() mailer.Mailer

	// NewFilesystem creates and returns a configured filesystem.System instance.
	//
	// NB! Make sure to call `Close()` on the returned result
	// after you are done working with it.
	NewFilesystem() (*filesystem.System, error)

	// RefreshSettings reinitializes and reloads the stored application settings.
	RefreshSettings() error

	// IsBootstrapped checks if the application was initialized
	// (aka. whether Bootstrap() was called).
	IsBootstrapped() bool

	// Bootstrap takes care for initializing the application
	// (open db connections, load settings, etc.).
	//
	// It will call ResetBootstrapState() if the application was already bootstrapped.
	Bootstrap() error

	// ResetBootstrapState takes care for releasing initialized app resources
	// (eg. closing db connections).
	ResetBootstrapState() error

	// ---------------------------------------------------------------
	// App event hooks
	// ---------------------------------------------------------------

	// OnBeforeBootstrap hook is triggered before initializing the base
	// application resources (eg. before db open and initial settings load).
	OnBeforeBootstrap() *hook.Hook[*BootstrapEvent]

	// OnAfterBootstrap hook is triggered after initializing the base
	// application resources (eg. after db open and initial settings load).
	OnAfterBootstrap() *hook.Hook[*BootstrapEvent]

	// OnBeforeServe hook is triggered before serving the internal router (echo),
	// allowing you to adjust its options and attach new routes.
	OnBeforeServe() *hook.Hook[*ServeEvent]

	// OnBeforeApiError hook is triggered right before sending an error API
	// response to the client, allowing you to further modify the error data
	// or to return a completely different API response (using [hook.StopPropagation]).
	OnBeforeApiError() *hook.Hook[*ApiErrorEvent]

	// OnAfterApiError hook is triggered right after sending an error API
	// response to the client.
	// It could be used to log the final API error in external services.
	OnAfterApiError() *hook.Hook[*ApiErrorEvent]

	// ---------------------------------------------------------------
	// Dao event hooks
	// ---------------------------------------------------------------

	// OnModelBeforeCreate hook is triggered before inserting a new
	// entry in the DB, allowing you to modify or validate the stored data.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelBeforeCreate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterCreate hook is triggered after successfully
	// inserting a new entry in the DB.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelAfterCreate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelBeforeUpdate hook is triggered before updating existing
	// entry in the DB, allowing you to modify or validate the stored data.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelBeforeUpdate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterUpdate hook is triggered after successfully updating
	// existing entry in the DB.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelAfterUpdate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelBeforeDelete hook is triggered before deleting an
	// existing entry from the DB.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelBeforeDelete(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterDelete is triggered after successfully deleting an
	// existing entry from the DB.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelAfterDelete(tags ...string) *hook.TaggedHook[*ModelEvent]

	// ---------------------------------------------------------------
	// Mailer event hooks
	// ---------------------------------------------------------------

	// OnMailerBeforeAdminResetPasswordSend hook is triggered right before
	// sending a password reset email to an admin.
	//
	// Could be used to send your own custom email template if
	// [hook.StopPropagation] is returned in one of its listeners.
	OnMailerBeforeAdminResetPasswordSend() *hook.Hook[*MailerAdminEvent]

	// OnMailerAfterAdminResetPasswordSend hook is triggered after
	// admin password reset email was successfully sent.
	OnMailerAfterAdminResetPasswordSend() *hook.Hook[*MailerAdminEvent]

	// OnMailerBeforeRecordResetPasswordSend hook is triggered right before
	// sending a password reset email to an auth record.
	//
	// Could be used to send your own custom email template if
	// [hook.StopPropagation] is returned in one of its listeners.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerBeforeRecordResetPasswordSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// OnMailerAfterRecordResetPasswordSend hook is triggered after
	// an auth record password reset email was successfully sent.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerAfterRecordResetPasswordSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// OnMailerBeforeRecordVerificationSend hook is triggered right before
	// sending a verification email to an auth record.
	//
	// Could be used to send your own custom email template if
	// [hook.StopPropagation] is returned in one of its listeners.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerBeforeRecordVerificationSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// OnMailerAfterRecordVerificationSend hook is triggered after a
	// verification email was successfully sent to an auth record.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerAfterRecordVerificationSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// OnMailerBeforeRecordChangeEmailSend hook is triggered right before
	// sending a confirmation new address email to an auth record.
	//
	// Could be used to send your own custom email template if
	// [hook.StopPropagation] is returned in one of its listeners.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerBeforeRecordChangeEmailSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// OnMailerAfterRecordChangeEmailSend hook is triggered after a
	// verification email was successfully sent to an auth record.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerAfterRecordChangeEmailSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// ---------------------------------------------------------------
	// Realtime API event hooks
	// ---------------------------------------------------------------

	// OnRealtimeConnectRequest hook is triggered right before establishing
	// the SSE client connection.
	OnRealtimeConnectRequest() *hook.Hook[*RealtimeConnectEvent]

	// OnRealtimeDisconnectRequest hook is triggered on disconnected/interrupted
	// SSE client connection.
	OnRealtimeDisconnectRequest() *hook.Hook[*RealtimeDisconnectEvent]

	// OnRealtimeBeforeMessage hook is triggered right before sending
	// an SSE message to a client.
	//
	// Returning [hook.StopPropagation] will prevent sending the message.
	// Returning any other non-nil error will close the realtime connection.
	OnRealtimeBeforeMessageSend() *hook.Hook[*RealtimeMessageEvent]

	// OnRealtimeBeforeMessage hook is triggered right after sending
	// an SSE message to a client.
	OnRealtimeAfterMessageSend() *hook.Hook[*RealtimeMessageEvent]

	// OnRealtimeBeforeSubscribeRequest hook is triggered before changing
	// the client subscriptions, allowing you to further validate and
	// modify the submitted change.
	OnRealtimeBeforeSubscribeRequest() *hook.Hook[*RealtimeSubscribeEvent]

	// OnRealtimeAfterSubscribeRequest hook is triggered after the client
	// subscriptions were successfully changed.
	OnRealtimeAfterSubscribeRequest() *hook.Hook[*RealtimeSubscribeEvent]

	// ---------------------------------------------------------------
	// Settings API event hooks
	// ---------------------------------------------------------------

	// OnSettingsListRequest hook is triggered on each successful
	// API Settings list request.
	//
	// Could be used to validate or modify the response before
	// returning it to the client.
	OnSettingsListRequest() *hook.Hook[*SettingsListEvent]

	// OnSettingsBeforeUpdateRequest hook is triggered before each API
	// Settings update request (after request data load and before settings persistence).
	//
	// Could be used to additionally validate the request data or
	// implement completely different persistence behavior
	// (returning [hook.StopPropagation]).
	OnSettingsBeforeUpdateRequest() *hook.Hook[*SettingsUpdateEvent]

	// OnSettingsAfterUpdateRequest hook is triggered after each
	// successful API Settings update request.
	OnSettingsAfterUpdateRequest() *hook.Hook[*SettingsUpdateEvent]

	// ---------------------------------------------------------------
	// File API event hooks
	// ---------------------------------------------------------------

	// OnFileDownloadRequest hook is triggered before each API File download request.
	//
	// Could be used to validate or modify the file response before
	// returning it to the client.
	OnFileDownloadRequest(tags ...string) *hook.TaggedHook[*FileDownloadEvent]

	// ---------------------------------------------------------------
	// Admin API event hooks
	// ---------------------------------------------------------------

	// OnAdminsListRequest hook is triggered on each API Admins list request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnAdminsListRequest() *hook.Hook[*AdminsListEvent]

	// OnAdminViewRequest hook is triggered on each API Admin view request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnAdminViewRequest() *hook.Hook[*AdminViewEvent]

	// OnAdminBeforeCreateRequest hook is triggered before each API
	// Admin create request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnAdminBeforeCreateRequest() *hook.Hook[*AdminCreateEvent]

	// OnAdminAfterCreateRequest hook is triggered after each
	// successful API Admin create request.
	OnAdminAfterCreateRequest() *hook.Hook[*AdminCreateEvent]

	// OnAdminBeforeUpdateRequest hook is triggered before each API
	// Admin update request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnAdminBeforeUpdateRequest() *hook.Hook[*AdminUpdateEvent]

	// OnAdminAfterUpdateRequest hook is triggered after each
	// successful API Admin update request.
	OnAdminAfterUpdateRequest() *hook.Hook[*AdminUpdateEvent]

	// OnAdminBeforeDeleteRequest hook is triggered before each API
	// Admin delete request (after model load and before actual deletion).
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior (returning [hook.StopPropagation]).
	OnAdminBeforeDeleteRequest() *hook.Hook[*AdminDeleteEvent]

	// OnAdminAfterDeleteRequest hook is triggered after each
	// successful API Admin delete request.
	OnAdminAfterDeleteRequest() *hook.Hook[*AdminDeleteEvent]

	// OnAdminAuthRequest hook is triggered on each successful API Admin
	// authentication request (sign-in, token refresh, etc.).
	//
	// Could be used to additionally validate or modify the
	// authenticated admin data and token.
	OnAdminAuthRequest() *hook.Hook[*AdminAuthEvent]

	// OnAdminBeforeAuthWithPasswordRequest hook is triggered before each Admin
	// auth with password API request (after request data load and before password validation).
	//
	// Could be used to implement for example a custom password validation
	// or to locate a different Admin identity (by assigning [AdminAuthWithPasswordEvent.Admin]).
	OnAdminBeforeAuthWithPasswordRequest() *hook.Hook[*AdminAuthWithPasswordEvent]

	// OnAdminAfterAuthWithPasswordRequest hook is triggered after each
	// successful Admin auth with password API request.
	OnAdminAfterAuthWithPasswordRequest() *hook.Hook[*AdminAuthWithPasswordEvent]

	// OnAdminBeforeAuthRefreshRequest hook is triggered before each Admin
	// auth refresh API request (right before generating a new auth token).
	//
	// Could be used to additionally validate the request data or implement
	// completely different auth refresh behavior (returning [hook.StopPropagation]).
	OnAdminBeforeAuthRefreshRequest() *hook.Hook[*AdminAuthRefreshEvent]

	// OnAdminAfterAuthRefreshRequest hook is triggered after each
	// successful auth refresh API request (right after generating a new auth token).
	OnAdminAfterAuthRefreshRequest() *hook.Hook[*AdminAuthRefreshEvent]

	// OnAdminBeforeRequestPasswordResetRequest hook is triggered before each Admin
	// request password reset API request (after request data load and before sending the reset email).
	//
	// Could be used to additionally validate the request data or implement
	// completely different password reset behavior (returning [hook.StopPropagation]).
	OnAdminBeforeRequestPasswordResetRequest() *hook.Hook[*AdminRequestPasswordResetEvent]

	// OnAdminAfterRequestPasswordResetRequest hook is triggered after each
	// successful request password reset API request.
	OnAdminAfterRequestPasswordResetRequest() *hook.Hook[*AdminRequestPasswordResetEvent]

	// OnAdminBeforeConfirmPasswordResetRequest hook is triggered before each Admin
	// confirm password reset API request (after request data load and before persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnAdminBeforeConfirmPasswordResetRequest() *hook.Hook[*AdminConfirmPasswordResetEvent]

	// OnAdminAfterConfirmPasswordResetRequest hook is triggered after each
	// successful confirm password reset API request.
	OnAdminAfterConfirmPasswordResetRequest() *hook.Hook[*AdminConfirmPasswordResetEvent]

	// ---------------------------------------------------------------
	// Record Auth API event hooks
	// ---------------------------------------------------------------

	// OnRecordAuthRequest hook is triggered on each successful API
	// record authentication request (sign-in, token refresh, etc.).
	//
	// Could be used to additionally validate or modify the authenticated
	// record data and token.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAuthRequest(tags ...string) *hook.TaggedHook[*RecordAuthEvent]

	// OnRecordBeforeAuthWithPasswordRequest hook is triggered before each Record
	// auth with password API request (after request data load and before password validation).
	//
	// Could be used to implement for example a custom password validation
	// or to locate a different Record identity (by assigning [RecordAuthWithPasswordEvent.Record]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*RecordAuthWithPasswordEvent]

	// OnRecordAfterAuthWithPasswordRequest hook is triggered after each
	// successful Record auth with password API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*RecordAuthWithPasswordEvent]

	// OnRecordBeforeAuthWithOAuth2Request hook is triggered before each Record
	// OAuth2 sign-in/sign-up API request (after token exchange and before external provider linking).
	//
	// If the [RecordAuthWithOAuth2Event.Record] is nil, then the OAuth2
	// request will try to create a new auth Record.
	//
	// To assign or link a different existing record model you can
	// overwrite/modify the [RecordAuthWithOAuth2Event.Record] field.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeAuthWithOAuth2Request(tags ...string) *hook.TaggedHook[*RecordAuthWithOAuth2Event]

	// OnRecordAfterAuthWithOAuth2Request hook is triggered after each
	// successful Record OAuth2 API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterAuthWithOAuth2Request(tags ...string) *hook.TaggedHook[*RecordAuthWithOAuth2Event]

	// OnRecordBeforeAuthRefreshRequest hook is triggered before each Record
	// auth refresh API request (right before generating a new auth token).
	//
	// Could be used to additionally validate the request data or implement
	// completely different auth refresh behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeAuthRefreshRequest(tags ...string) *hook.TaggedHook[*RecordAuthRefreshEvent]

	// OnRecordAfterAuthRefreshRequest hook is triggered after each
	// successful auth refresh API request (right after generating a new auth token).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterAuthRefreshRequest(tags ...string) *hook.TaggedHook[*RecordAuthRefreshEvent]

	// OnRecordBeforeRequestPasswordResetRequest hook is triggered before each Record
	// request password reset API request (after request data load and before sending the reset email).
	//
	// Could be used to additionally validate the request data or implement
	// completely different password reset behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeRequestPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordRequestPasswordResetEvent]

	// OnRecordAfterRequestPasswordResetRequest hook is triggered after each
	// successful request password reset API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterRequestPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordRequestPasswordResetEvent]

	// OnRecordBeforeConfirmPasswordResetRequest hook is triggered before each Record
	// confirm password reset API request (after request data load and before persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeConfirmPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordConfirmPasswordResetEvent]

	// OnRecordAfterConfirmPasswordResetRequest hook is triggered after each
	// successful confirm password reset API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterConfirmPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordConfirmPasswordResetEvent]

	// OnRecordBeforeRequestVerificationRequest hook is triggered before each Record
	// request verification API request (after request data load and before sending the verification email).
	//
	// Could be used to additionally validate the loaded request data or implement
	// completely different verification behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeRequestVerificationRequest(tags ...string) *hook.TaggedHook[*RecordRequestVerificationEvent]

	// OnRecordAfterRequestVerificationRequest hook is triggered after each
	// successful request verification API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterRequestVerificationRequest(tags ...string) *hook.TaggedHook[*RecordRequestVerificationEvent]

	// OnRecordBeforeConfirmVerificationRequest hook is triggered before each Record
	// confirm verification API request (after request data load and before persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeConfirmVerificationRequest(tags ...string) *hook.TaggedHook[*RecordConfirmVerificationEvent]

	// OnRecordAfterConfirmVerificationRequest hook is triggered after each
	// successful confirm verification API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterConfirmVerificationRequest(tags ...string) *hook.TaggedHook[*RecordConfirmVerificationEvent]

	// OnRecordBeforeRequestEmailChangeRequest hook is triggered before each Record request email change API request
	// (after request data load and before sending the email link to confirm the change).
	//
	// Could be used to additionally validate the request data or implement
	// completely different request email change behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeRequestEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordRequestEmailChangeEvent]

	// OnRecordAfterRequestEmailChangeRequest hook is triggered after each
	// successful request email change API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterRequestEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordRequestEmailChangeEvent]

	// OnRecordBeforeConfirmEmailChangeRequest hook is triggered before each Record
	// confirm email change API request (after request data load and before persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeConfirmEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordConfirmEmailChangeEvent]

	// OnRecordAfterConfirmEmailChangeRequest hook is triggered after each
	// successful confirm email change API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterConfirmEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordConfirmEmailChangeEvent]

	// OnRecordListExternalAuthsRequest hook is triggered on each API record external auths list request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordListExternalAuthsRequest(tags ...string) *hook.TaggedHook[*RecordListExternalAuthsEvent]

	// OnRecordBeforeUnlinkExternalAuthRequest hook is triggered before each API record
	// external auth unlink request (after models load and before the actual relation deletion).
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeUnlinkExternalAuthRequest(tags ...string) *hook.TaggedHook[*RecordUnlinkExternalAuthEvent]

	// OnRecordAfterUnlinkExternalAuthRequest hook is triggered after each
	// successful API record external auth unlink request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterUnlinkExternalAuthRequest(tags ...string) *hook.TaggedHook[*RecordUnlinkExternalAuthEvent]

	// ---------------------------------------------------------------
	// Record CRUD API event hooks
	// ---------------------------------------------------------------

	// OnRecordsListRequest hook is triggered on each API Records list request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordsListRequest(tags ...string) *hook.TaggedHook[*RecordsListEvent]

	// OnRecordViewRequest hook is triggered on each API Record view request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordViewRequest(tags ...string) *hook.TaggedHook[*RecordViewEvent]

	// OnRecordBeforeCreateRequest hook is triggered before each API Record
	// create request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeCreateRequest(tags ...string) *hook.TaggedHook[*RecordCreateEvent]

	// OnRecordAfterCreateRequest hook is triggered after each
	// successful API Record create request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterCreateRequest(tags ...string) *hook.TaggedHook[*RecordCreateEvent]

	// OnRecordBeforeUpdateRequest hook is triggered before each API Record
	// update request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeUpdateRequest(tags ...string) *hook.TaggedHook[*RecordUpdateEvent]

	// OnRecordAfterUpdateRequest hook is triggered after each
	// successful API Record update request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterUpdateRequest(tags ...string) *hook.TaggedHook[*RecordUpdateEvent]

	// OnRecordBeforeDeleteRequest hook is triggered before each API Record
	// delete request (after model load and before actual deletion).
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior (returning [hook.StopPropagation]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordBeforeDeleteRequest(tags ...string) *hook.TaggedHook[*RecordDeleteEvent]

	// OnRecordAfterDeleteRequest hook is triggered after each
	// successful API Record delete request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterDeleteRequest(tags ...string) *hook.TaggedHook[*RecordDeleteEvent]

	// ---------------------------------------------------------------
	// Collection API event hooks
	// ---------------------------------------------------------------

	// OnCollectionsListRequest hook is triggered on each API Collections list request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnCollectionsListRequest() *hook.Hook[*CollectionsListEvent]

	// OnCollectionViewRequest hook is triggered on each API Collection view request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnCollectionViewRequest() *hook.Hook[*CollectionViewEvent]

	// OnCollectionBeforeCreateRequest hook is triggered before each API Collection
	// create request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnCollectionBeforeCreateRequest() *hook.Hook[*CollectionCreateEvent]

	// OnCollectionAfterCreateRequest hook is triggered after each
	// successful API Collection create request.
	OnCollectionAfterCreateRequest() *hook.Hook[*CollectionCreateEvent]

	// OnCollectionBeforeUpdateRequest hook is triggered before each API Collection
	// update request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnCollectionBeforeUpdateRequest() *hook.Hook[*CollectionUpdateEvent]

	// OnCollectionAfterUpdateRequest hook is triggered after each
	// successful API Collection update request.
	OnCollectionAfterUpdateRequest() *hook.Hook[*CollectionUpdateEvent]

	// OnCollectionBeforeDeleteRequest hook is triggered before each API
	// Collection delete request (after model load and before actual deletion).
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior (returning [hook.StopPropagation]).
	OnCollectionBeforeDeleteRequest() *hook.Hook[*CollectionDeleteEvent]

	// OnCollectionAfterDeleteRequest hook is triggered after each
	// successful API Collection delete request.
	OnCollectionAfterDeleteRequest() *hook.Hook[*CollectionDeleteEvent]

	// OnCollectionsBeforeImportRequest hook is triggered before each API
	// collections import request (after request data load and before the actual import).
	//
	// Could be used to additionally validate the imported collections or
	// to implement completely different import behavior (returning [hook.StopPropagation]).
	OnCollectionsBeforeImportRequest() *hook.Hook[*CollectionsImportEvent]

	// OnCollectionsAfterImportRequest hook is triggered after each
	// successful API collections import request.
	OnCollectionsAfterImportRequest() *hook.Hook[*CollectionsImportEvent]
}
