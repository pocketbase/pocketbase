// Package core is the backbone of PocketBase.
//
// It defines the main PocketBase App interface and its base implementation.
package core

import (
	"context"
	"log/slog"

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

	// Logger returns the active app logger.
	Logger() *slog.Logger

	// DataDir returns the app data directory path.
	DataDir() string

	// EncryptionEnv returns the name of the app secret env key
	// (used for settings encryption).
	EncryptionEnv() string

	// IsDev returns whether the app is in dev mode.
	IsDev() bool

	// Settings returns the loaded app settings.
	Settings() *settings.Settings

	// Deprecated: Use app.Store() instead.
	Cache() *store.Store[any]

	// Store returns the app runtime store.
	Store() *store.Store[any]

	// SubscriptionsBroker returns the app realtime subscriptions broker instance.
	SubscriptionsBroker() *subscriptions.Broker

	// NewMailClient creates and returns a configured app mail client.
	NewMailClient() mailer.Mailer

	// NewFilesystem creates and returns a configured filesystem.System instance
	// for managing regular app files (eg. collection uploads).
	//
	// NB! Make sure to call Close() on the returned result
	// after you are done working with it.
	NewFilesystem() (*filesystem.System, error)

	// NewBackupsFilesystem creates and returns a configured filesystem.System instance
	// for managing app backups.
	//
	// NB! Make sure to call Close() on the returned result
	// after you are done working with it.
	NewBackupsFilesystem() (*filesystem.System, error)

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

	// CreateBackup creates a new backup of the current app pb_data directory.
	//
	// Backups can be stored on S3 if it is configured in app.Settings().Backups.
	//
	// Please refer to the godoc of the specific core.App implementation
	// for details on the backup procedures.
	CreateBackup(ctx context.Context, name string) error

	// RestoreBackup restores the backup with the specified name and restarts
	// the current running application process.
	//
	// The safely perform the restore it is recommended to have free disk space
	// for at least 2x the size of the restored pb_data backup.
	//
	// Please refer to the godoc of the specific core.App implementation
	// for details on the restore procedures.
	//
	// NB! This feature is experimental and currently is expected to work only on UNIX based systems.
	RestoreBackup(ctx context.Context, name string) error

	// Restart restarts the current running application process.
	//
	// Currently it is relying on execve so it is supported only on UNIX based systems.
	Restart() error

	// ---------------------------------------------------------------
	// App event hooks
	// ---------------------------------------------------------------

	// OnBeforeBootstrap hook is triggered before initializing the main
	// application resources (eg. before db open and initial settings load).
	OnBeforeBootstrap() *hook.Hook[*BootstrapEvent]

	// OnAfterBootstrap hook is triggered after initializing the main
	// application resources (eg. after db open and initial settings load).
	OnAfterBootstrap() *hook.Hook[*BootstrapEvent]

	// OnBeforeServe hook is triggered before serving the internal router (echo),
	// allowing you to adjust its options and attach new routes or middlewares.
	OnBeforeServe() *hook.Hook[*ServeEvent]

	// OnBeforeApiError hook is triggered right before sending an error API
	// response to the client, allowing you to further modify the error data
	// or to return a completely different API response.
	OnBeforeApiError() *hook.Hook[*ApiErrorEvent]

	// OnAfterApiError hook is triggered right after sending an error API
	// response to the client.
	// It could be used to log the final API error in external services.
	OnAfterApiError() *hook.Hook[*ApiErrorEvent]

	// OnTerminate hook is triggered when the app is in the process
	// of being terminated (eg. on SIGTERM signal).
	OnTerminate() *hook.Hook[*TerminateEvent]

	// ---------------------------------------------------------------
	// Dao event hooks
	// ---------------------------------------------------------------

	// OnModelBeforeCreate hook is triggered before inserting a new
	// model in the DB, allowing you to modify or validate the stored data.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelBeforeCreate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterCreate hook is triggered after successfully
	// inserting a new model in the DB.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelAfterCreate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelBeforeUpdate hook is triggered before updating existing
	// model in the DB, allowing you to modify or validate the stored data.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelBeforeUpdate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterUpdate hook is triggered after successfully updating
	// existing model in the DB.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelAfterUpdate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelBeforeDelete hook is triggered before deleting an
	// existing model from the DB.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelBeforeDelete(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterDelete hook is triggered after successfully deleting an
	// existing model from the DB.
	//
	// If the optional "tags" list (table names and/or the Collection id for Record models)
	// is specified, then all event handlers registered via the created hook
	// will be triggered and called only if their event data origin matches the tags.
	OnModelAfterDelete(tags ...string) *hook.TaggedHook[*ModelEvent]

	// ---------------------------------------------------------------
	// Mailer event hooks
	// ---------------------------------------------------------------

	// OnMailerBeforeAdminResetPasswordSend hook is triggered right
	// before sending a password reset email to an admin, allowing you
	// to inspect and customize the email message that is being sent.
	OnMailerBeforeAdminResetPasswordSend() *hook.Hook[*MailerAdminEvent]

	// OnMailerAfterAdminResetPasswordSend hook is triggered after
	// admin password reset email was successfully sent.
	OnMailerAfterAdminResetPasswordSend() *hook.Hook[*MailerAdminEvent]

	// OnMailerBeforeRecordResetPasswordSend hook is triggered right
	// before sending a password reset email to an auth record, allowing
	// you to inspect and customize the email message that is being sent.
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

	// OnMailerBeforeRecordVerificationSend hook is triggered right
	// before sending a verification email to an auth record, allowing
	// you to inspect and customize the email message that is being sent.
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
	// sending a confirmation new address email to an auth record, allowing
	// you to inspect and customize the email message that is being sent.
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

	// OnRealtimeBeforeMessageSend hook is triggered right before sending
	// an SSE message to a client.
	//
	// Returning [hook.StopPropagation] will prevent sending the message.
	// Returning any other non-nil error will close the realtime connection.
	OnRealtimeBeforeMessageSend() *hook.Hook[*RealtimeMessageEvent]

	// OnRealtimeAfterMessageSend hook is triggered right after sending
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
	// implement completely different persistence behavior.
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

	// OnFileBeforeTokenRequest hook is triggered before each file
	// token API request.
	//
	// If no token or model was submitted, e.Model and e.Token will be empty,
	// allowing you to implement your own custom model file auth implementation.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnFileBeforeTokenRequest(tags ...string) *hook.TaggedHook[*FileTokenEvent]

	// OnFileAfterTokenRequest hook is triggered after each
	// successful file token API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnFileAfterTokenRequest(tags ...string) *hook.TaggedHook[*FileTokenEvent]

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
	// completely different persistence behavior.
	OnAdminBeforeCreateRequest() *hook.Hook[*AdminCreateEvent]

	// OnAdminAfterCreateRequest hook is triggered after each
	// successful API Admin create request.
	OnAdminAfterCreateRequest() *hook.Hook[*AdminCreateEvent]

	// OnAdminBeforeUpdateRequest hook is triggered before each API
	// Admin update request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
	OnAdminBeforeUpdateRequest() *hook.Hook[*AdminUpdateEvent]

	// OnAdminAfterUpdateRequest hook is triggered after each
	// successful API Admin update request.
	OnAdminAfterUpdateRequest() *hook.Hook[*AdminUpdateEvent]

	// OnAdminBeforeDeleteRequest hook is triggered before each API
	// Admin delete request (after model load and before actual deletion).
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior.
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
	// completely different auth refresh behavior.
	OnAdminBeforeAuthRefreshRequest() *hook.Hook[*AdminAuthRefreshEvent]

	// OnAdminAfterAuthRefreshRequest hook is triggered after each
	// successful auth refresh API request (right after generating a new auth token).
	OnAdminAfterAuthRefreshRequest() *hook.Hook[*AdminAuthRefreshEvent]

	// OnAdminBeforeRequestPasswordResetRequest hook is triggered before each Admin
	// request password reset API request (after request data load and before sending the reset email).
	//
	// Could be used to additionally validate the request data or implement
	// completely different password reset behavior.
	OnAdminBeforeRequestPasswordResetRequest() *hook.Hook[*AdminRequestPasswordResetEvent]

	// OnAdminAfterRequestPasswordResetRequest hook is triggered after each
	// successful request password reset API request.
	OnAdminAfterRequestPasswordResetRequest() *hook.Hook[*AdminRequestPasswordResetEvent]

	// OnAdminBeforeConfirmPasswordResetRequest hook is triggered before each Admin
	// confirm password reset API request (after request data load and before persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
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
	// or to locate a different Record model (by reassigning [RecordAuthWithPasswordEvent.Record]).
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
	// If the [RecordAuthWithOAuth2Event.Record] is not set, then the OAuth2
	// request will try to create a new auth Record.
	//
	// To assign or link a different existing record model you can
	// change the [RecordAuthWithOAuth2Event.Record] field.
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
	// completely different auth refresh behavior.
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
	// completely different delete behavior.
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

	// OnRecordBeforeRequestPasswordResetRequest hook is triggered before each Record
	// request password reset API request (after request data load and before sending the reset email).
	//
	// Could be used to additionally validate the request data or implement
	// completely different password reset behavior.
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
	// completely different persistence behavior.
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
	// completely different verification behavior.
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
	// completely different persistence behavior.
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
	// completely different request email change behavior.
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
	// completely different persistence behavior.
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
	// completely different persistence behavior.
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
	// completely different persistence behavior.
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
	// completely different delete behavior.
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
	// completely different persistence behavior.
	OnCollectionBeforeCreateRequest() *hook.Hook[*CollectionCreateEvent]

	// OnCollectionAfterCreateRequest hook is triggered after each
	// successful API Collection create request.
	OnCollectionAfterCreateRequest() *hook.Hook[*CollectionCreateEvent]

	// OnCollectionBeforeUpdateRequest hook is triggered before each API Collection
	// update request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
	OnCollectionBeforeUpdateRequest() *hook.Hook[*CollectionUpdateEvent]

	// OnCollectionAfterUpdateRequest hook is triggered after each
	// successful API Collection update request.
	OnCollectionAfterUpdateRequest() *hook.Hook[*CollectionUpdateEvent]

	// OnCollectionBeforeDeleteRequest hook is triggered before each API
	// Collection delete request (after model load and before actual deletion).
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior.
	OnCollectionBeforeDeleteRequest() *hook.Hook[*CollectionDeleteEvent]

	// OnCollectionAfterDeleteRequest hook is triggered after each
	// successful API Collection delete request.
	OnCollectionAfterDeleteRequest() *hook.Hook[*CollectionDeleteEvent]

	// OnCollectionsBeforeImportRequest hook is triggered before each API
	// collections import request (after request data load and before the actual import).
	//
	// Could be used to additionally validate the imported collections or
	// to implement completely different import behavior.
	OnCollectionsBeforeImportRequest() *hook.Hook[*CollectionsImportEvent]

	// OnCollectionsAfterImportRequest hook is triggered after each
	// successful API collections import request.
	OnCollectionsAfterImportRequest() *hook.Hook[*CollectionsImportEvent]
}
