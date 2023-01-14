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
	OnModelBeforeCreate() *hook.Hook[*ModelEvent]

	// OnModelAfterCreate hook is triggered after successfully
	// inserting a new entry in the DB.
	OnModelAfterCreate() *hook.Hook[*ModelEvent]

	// OnModelBeforeUpdate hook is triggered before updating existing
	// entry in the DB, allowing you to modify or validate the stored data.
	OnModelBeforeUpdate() *hook.Hook[*ModelEvent]

	// OnModelAfterUpdate hook is triggered after successfully updating
	// existing entry in the DB.
	OnModelAfterUpdate() *hook.Hook[*ModelEvent]

	// OnModelBeforeDelete hook is triggered before deleting an
	// existing entry from the DB.
	OnModelBeforeDelete() *hook.Hook[*ModelEvent]

	// OnModelAfterDelete is triggered after successfully deleting an
	// existing entry from the DB.
	OnModelAfterDelete() *hook.Hook[*ModelEvent]

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
	OnMailerBeforeRecordResetPasswordSend() *hook.Hook[*MailerRecordEvent]

	// OnMailerAfterRecordResetPasswordSend hook is triggered after
	// an auth record password reset email was successfully sent.
	OnMailerAfterRecordResetPasswordSend() *hook.Hook[*MailerRecordEvent]

	// OnMailerBeforeRecordVerificationSend hook is triggered right before
	// sending a verification email to an auth record.
	//
	// Could be used to send your own custom email template if
	// [hook.StopPropagation] is returned in one of its listeners.
	OnMailerBeforeRecordVerificationSend() *hook.Hook[*MailerRecordEvent]

	// OnMailerAfterRecordVerificationSend hook is triggered after a
	// verification email was successfully sent to an auth record.
	OnMailerAfterRecordVerificationSend() *hook.Hook[*MailerRecordEvent]

	// OnMailerBeforeRecordChangeEmailSend hook is triggered right before
	// sending a confirmation new address email to an auth record.
	//
	// Could be used to send your own custom email template if
	// [hook.StopPropagation] is returned in one of its listeners.
	OnMailerBeforeRecordChangeEmailSend() *hook.Hook[*MailerRecordEvent]

	// OnMailerAfterRecordChangeEmailSend hook is triggered after a
	// verification email was successfully sent to an auth record.
	OnMailerAfterRecordChangeEmailSend() *hook.Hook[*MailerRecordEvent]

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
	OnFileDownloadRequest() *hook.Hook[*FileDownloadEvent]

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

	// ---------------------------------------------------------------
	// Record Auth API event hooks
	// ---------------------------------------------------------------

	// OnRecordAuthRequest hook is triggered on each successful API
	// record authentication request (sign-in, token refresh, etc.).
	//
	// Could be used to additionally validate or modify the authenticated
	// record data and token.
	OnRecordAuthRequest() *hook.Hook[*RecordAuthEvent]

	// OnRecordBeforeRequestPasswordResetRequest hook is triggered before each Record
	// request password reset API request (after request data load and before sending the reset email).
	//
	// Could be used to additionally validate the request data or implement
	// completely different password reset behavior (returning [hook.StopPropagation]).
	OnRecordBeforeRequestPasswordResetRequest() *hook.Hook[*RecordRequestPasswordResetEvent]

	// OnRecordAfterRequestPasswordResetRequest hook is triggered after each
	// successful request password reset API request.
	OnRecordAfterRequestPasswordResetRequest() *hook.Hook[*RecordRequestPasswordResetEvent]

	// OnRecordBeforeConfirmPasswordResetRequest hook is triggered before each Record
	// confirm password reset API request (after request data load and before persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnRecordBeforeConfirmPasswordResetRequest() *hook.Hook[*RecordConfirmPasswordResetEvent]

	// OnRecordAfterConfirmPasswordResetRequest hook is triggered after each
	// successful confirm password reset API request.
	OnRecordAfterConfirmPasswordResetRequest() *hook.Hook[*RecordConfirmPasswordResetEvent]

	// OnRecordBeforeRequestVerificationRequest hook is triggered before each Record
	// request verification API request (after request data load and before sending the verification email).
	//
	// Could be used to additionally validate the loaded request data or implement
	// completely different verification behavior (returning [hook.StopPropagation]).
	OnRecordBeforeRequestVerificationRequest() *hook.Hook[*RecordRequestVerificationEvent]

	// OnRecordAfterRequestVerificationRequest hook is triggered after each
	// successful request verification API request.
	OnRecordAfterRequestVerificationRequest() *hook.Hook[*RecordRequestVerificationEvent]

	// OnRecordBeforeConfirmVerificationRequest hook is triggered before each Record
	// confirm verification API request (after request data load and before persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnRecordBeforeConfirmVerificationRequest() *hook.Hook[*RecordConfirmVerificationEvent]

	// OnRecordAfterConfirmVerificationRequest hook is triggered after each
	// successful confirm verification API request.
	OnRecordAfterConfirmVerificationRequest() *hook.Hook[*RecordConfirmVerificationEvent]

	// OnRecordBeforeRequestEmailChangeRequest hook is triggered before each Record request email change API request
	// (after request data load and before sending the email link to confirm the change).
	//
	// Could be used to additionally validate the request data or implement
	// completely different request email change behavior (returning [hook.StopPropagation]).
	OnRecordBeforeRequestEmailChangeRequest() *hook.Hook[*RecordRequestEmailChangeEvent]

	// OnRecordAfterRequestEmailChangeRequest hook is triggered after each
	// successful request email change API request.
	OnRecordAfterRequestEmailChangeRequest() *hook.Hook[*RecordRequestEmailChangeEvent]

	// OnRecordBeforeConfirmEmailChangeRequest hook is triggered before each Record
	// confirm email change API request (after request data load and before persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnRecordBeforeConfirmEmailChangeRequest() *hook.Hook[*RecordConfirmEmailChangeEvent]

	// OnRecordAfterConfirmEmailChangeRequest hook is triggered after each
	// successful confirm email change API request.
	OnRecordAfterConfirmEmailChangeRequest() *hook.Hook[*RecordConfirmEmailChangeEvent]

	// OnRecordListExternalAuthsRequest hook is triggered on each API record external auths list request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnRecordListExternalAuthsRequest() *hook.Hook[*RecordListExternalAuthsEvent]

	// OnRecordBeforeUnlinkExternalAuthRequest hook is triggered before each API record
	// external auth unlink request (after models load and before the actual relation deletion).
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior (returning [hook.StopPropagation]).
	OnRecordBeforeUnlinkExternalAuthRequest() *hook.Hook[*RecordUnlinkExternalAuthEvent]

	// OnRecordAfterUnlinkExternalAuthRequest hook is triggered after each
	// successful API record external auth unlink request.
	OnRecordAfterUnlinkExternalAuthRequest() *hook.Hook[*RecordUnlinkExternalAuthEvent]

	// ---------------------------------------------------------------
	// Record CRUD API event hooks
	// ---------------------------------------------------------------

	// OnRecordsListRequest hook is triggered on each API Records list request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnRecordsListRequest() *hook.Hook[*RecordsListEvent]

	// OnRecordViewRequest hook is triggered on each API Record view request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnRecordViewRequest() *hook.Hook[*RecordViewEvent]

	// OnRecordBeforeCreateRequest hook is triggered before each API Record
	// create request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnRecordBeforeCreateRequest() *hook.Hook[*RecordCreateEvent]

	// OnRecordAfterCreateRequest hook is triggered after each
	// successful API Record create request.
	OnRecordAfterCreateRequest() *hook.Hook[*RecordCreateEvent]

	// OnRecordBeforeUpdateRequest hook is triggered before each API Record
	// update request (after request data load and before model persistence).
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior (returning [hook.StopPropagation]).
	OnRecordBeforeUpdateRequest() *hook.Hook[*RecordUpdateEvent]

	// OnRecordAfterUpdateRequest hook is triggered after each
	// successful API Record update request.
	OnRecordAfterUpdateRequest() *hook.Hook[*RecordUpdateEvent]

	// OnRecordBeforeDeleteRequest hook is triggered before each API Record
	// delete request (after model load and before actual deletion).
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior (returning [hook.StopPropagation]).
	OnRecordBeforeDeleteRequest() *hook.Hook[*RecordDeleteEvent]

	// OnRecordAfterDeleteRequest hook is triggered after each
	// successful API Record delete request.
	OnRecordAfterDeleteRequest() *hook.Hook[*RecordDeleteEvent]

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
