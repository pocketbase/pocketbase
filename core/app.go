// Package core is the backbone of PocketBase.
//
// It defines the main PocketBase App interface and its base implementation.
package core

import (
	"context"
	"log/slog"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

// App defines the main PocketBase app interface.
//
// Note that the interface is not intended to be implemented manually by users
// and instead they should use core.BaseApp (either directly or as embedded field in a custom struct).
//
// This interface exists to make testing easier and to allow users to
// create common and pluggable helpers and methods that doesn't rely
// on a specific wrapped app struct (hence the large interface size).
type App interface {
	// UnsafeWithoutHooks returns a shallow copy of the current app WITHOUT any registered hooks.
	//
	// NB! Note that using the returned app instance may cause data integrity errors
	// since the Record validations and data normalizations (including files uploads)
	// rely on the app hooks to work.
	UnsafeWithoutHooks() App

	// Logger returns the default app logger.
	//
	// If the application is not bootstrapped yet, fallbacks to slog.Default().
	Logger() *slog.Logger

	// IsBootstrapped checks if the application was initialized
	// (aka. whether Bootstrap() was called).
	IsBootstrapped() bool

	// IsTransactional checks if the current app instance is part of a transaction.
	IsTransactional() bool

	// TxInfo returns the transaction associated with the current app instance (if any).
	//
	// Could be used if you want to execute indirectly a function after
	// the related app transaction completes using `app.TxInfo().OnAfterFunc(callback)`.
	TxInfo() *TxAppInfo

	// Bootstrap initializes the application
	// (aka. create data dir, open db connections, load settings, etc.).
	//
	// It will call ResetBootstrapState() if the application was already bootstrapped.
	Bootstrap() error

	// ResetBootstrapState releases the initialized core app resources
	// (closing db connections, stopping cron ticker, etc.).
	ResetBootstrapState() error

	// DataDir returns the app data directory path.
	DataDir() string

	// EncryptionEnv returns the name of the app secret env key
	// (currently used primarily for optional settings encryption but this may change in the future).
	EncryptionEnv() string

	// IsDev returns whether the app is in dev mode.
	//
	// When enabled logs, executed sql statements, etc. are printed to the stderr.
	IsDev() bool

	// Settings returns the loaded app settings.
	Settings() *Settings

	// Store returns the app runtime store.
	Store() *store.Store[string, any]

	// Cron returns the app cron instance.
	Cron() *cron.Cron

	// SubscriptionsBroker returns the app realtime subscriptions broker instance.
	SubscriptionsBroker() *subscriptions.Broker

	// NewMailClient creates and returns a new SMTP or Sendmail client
	// based on the current app settings.
	NewMailClient() mailer.Mailer

	// NewFilesystem creates a new local or S3 filesystem instance
	// for managing regular app files (ex. record uploads)
	// based on the current app settings.
	//
	// NB! Make sure to call Close() on the returned result
	// after you are done working with it.
	NewFilesystem() (*filesystem.System, error)

	// NewBackupsFilesystem creates a new local or S3 filesystem instance
	// for managing app backups based on the current app settings.
	//
	// NB! Make sure to call Close() on the returned result
	// after you are done working with it.
	NewBackupsFilesystem() (*filesystem.System, error)

	// ReloadSettings reinitializes and reloads the stored application settings.
	ReloadSettings() error

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

	// Restart restarts (aka. replaces) the current running application process.
	//
	// NB! It relies on execve which is supported only on UNIX based systems.
	Restart() error

	// RunSystemMigrations applies all new migrations registered in the [core.SystemMigrations] list.
	RunSystemMigrations() error

	// RunAppMigrations applies all new migrations registered in the [core.AppMigrations] list.
	RunAppMigrations() error

	// RunAllMigrations applies all system and app migrations
	// (aka. from both [core.SystemMigrations] and [core.AppMigrations]).
	RunAllMigrations() error

	// ---------------------------------------------------------------
	// DB methods
	// ---------------------------------------------------------------

	// DB returns the default app data.db builder instance.
	//
	// To minimize SQLITE_BUSY errors, it automatically routes the
	// SELECT queries to the underlying concurrent db pool and everything else
	// to the nonconcurrent one.
	//
	// For more finer control over the used connections pools you can
	// call directly ConcurrentDB() or NonconcurrentDB().
	DB() dbx.Builder

	// ConcurrentDB returns the concurrent app data.db builder instance.
	//
	// This method is used mainly internally for executing db read
	// operations in a concurrent/non-blocking manner.
	//
	// Most users should use simply DB() as it will automatically
	// route the query execution to ConcurrentDB() or NonconcurrentDB().
	//
	// In a transaction the ConcurrentDB() and NonconcurrentDB() refer to the same *dbx.TX instance.
	ConcurrentDB() dbx.Builder

	// NonconcurrentDB returns the nonconcurrent app data.db builder instance.
	//
	// The returned db instance is limited only to a single open connection,
	// meaning that it can process only 1 db operation at a time (other queries queue up).
	//
	// This method is used mainly internally and in the tests to execute write
	// (save/delete) db operations as it helps with minimizing the SQLITE_BUSY errors.
	//
	// Most users should use simply DB() as it will automatically
	// route the query execution to ConcurrentDB() or NonconcurrentDB().
	//
	// In a transaction the ConcurrentDB() and NonconcurrentDB() refer to the same *dbx.TX instance.
	NonconcurrentDB() dbx.Builder

	// AuxDB returns the app auxiliary.db builder instance.
	//
	// To minimize SQLITE_BUSY errors, it automatically routes the
	// SELECT queries to the underlying concurrent db pool and everything else
	// to the nonconcurrent one.
	//
	// For more finer control over the used connections pools you can
	// call directly AuxConcurrentDB() or AuxNonconcurrentDB().
	AuxDB() dbx.Builder

	// AuxConcurrentDB returns the concurrent app auxiliary.db builder instance.
	//
	// This method is used mainly internally for executing db read
	// operations in a concurrent/non-blocking manner.
	//
	// Most users should use simply AuxDB() as it will automatically
	// route the query execution to AuxConcurrentDB() or AuxNonconcurrentDB().
	//
	// In a transaction the AuxConcurrentDB() and AuxNonconcurrentDB() refer to the same *dbx.TX instance.
	AuxConcurrentDB() dbx.Builder

	// AuxNonconcurrentDB returns the nonconcurrent app auxiliary.db builder instance.
	//
	// The returned db instance is limited only to a single open connection,
	// meaning that it can process only 1 db operation at a time (other queries queue up).
	//
	// This method is used mainly internally and in the tests to execute write
	// (save/delete) db operations as it helps with minimizing the SQLITE_BUSY errors.
	//
	// Most users should use simply AuxDB() as it will automatically
	// route the query execution to AuxConcurrentDB() or AuxNonconcurrentDB().
	//
	// In a transaction the AuxConcurrentDB() and AuxNonconcurrentDB() refer to the same *dbx.TX instance.
	AuxNonconcurrentDB() dbx.Builder

	// HasTable checks if a table (or view) with the provided name exists (case insensitive).
	// in the data.db.
	HasTable(tableName string) bool

	// AuxHasTable checks if a table (or view) with the provided name exists (case insensitive)
	// in the auxiliary.db.
	AuxHasTable(tableName string) bool

	// TableColumns returns all column names of a single table by its name.
	TableColumns(tableName string) ([]string, error)

	// TableInfo returns the "table_info" pragma result for the specified table.
	TableInfo(tableName string) ([]*TableInfoRow, error)

	// TableIndexes returns a name grouped map with all non empty index of the specified table.
	//
	// Note: This method doesn't return an error on nonexisting table.
	TableIndexes(tableName string) (map[string]string, error)

	// DeleteTable drops the specified table.
	//
	// This method is a no-op if a table with the provided name doesn't exist.
	//
	// NB! Be aware that this method is vulnerable to SQL injection and the
	// "tableName" argument must come only from trusted input!
	DeleteTable(tableName string) error

	// DeleteView drops the specified view name.
	//
	// This method is a no-op if a view with the provided name doesn't exist.
	//
	// NB! Be aware that this method is vulnerable to SQL injection and the
	// "name" argument must come only from trusted input!
	DeleteView(name string) error

	// SaveView creates (or updates already existing) persistent SQL view.
	//
	// NB! Be aware that this method is vulnerable to SQL injection and the
	// "selectQuery" argument must come only from trusted input!
	SaveView(name string, selectQuery string) error

	// CreateViewFields creates a new FieldsList from the provided select query.
	//
	// There are some caveats:
	// - The select query must have an "id" column.
	// - Wildcard ("*") columns are not supported to avoid accidentally leaking sensitive data.
	CreateViewFields(selectQuery string) (FieldsList, error)

	// FindRecordByViewFile returns the original Record of the provided view collection file.
	FindRecordByViewFile(viewCollectionModelOrIdentifier any, fileFieldName string, filename string) (*Record, error)

	// Vacuum executes VACUUM on the data.db in order to reclaim unused data db disk space.
	Vacuum() error

	// AuxVacuum executes VACUUM on the auxiliary.db in order to reclaim unused auxiliary db disk space.
	AuxVacuum() error

	// ---------------------------------------------------------------

	// ModelQuery creates a new preconfigured select data.db query with preset
	// SELECT, FROM and other common fields based on the provided model.
	ModelQuery(model Model) *dbx.SelectQuery

	// AuxModelQuery creates a new preconfigured select auxiliary.db query with preset
	// SELECT, FROM and other common fields based on the provided model.
	AuxModelQuery(model Model) *dbx.SelectQuery

	// Delete deletes the specified model from the regular app database.
	Delete(model Model) error

	// Delete deletes the specified model from the regular app database
	// (the context could be used to limit the query execution).
	DeleteWithContext(ctx context.Context, model Model) error

	// AuxDelete deletes the specified model from the auxiliary database.
	AuxDelete(model Model) error

	// AuxDeleteWithContext deletes the specified model from the auxiliary database
	// (the context could be used to limit the query execution).
	AuxDeleteWithContext(ctx context.Context, model Model) error

	// Save validates and saves the specified model into the regular app database.
	//
	// If you don't want to run validations, use [App.SaveNoValidate()].
	Save(model Model) error

	// SaveWithContext is the same as [App.Save()] but allows specifying a context to limit the db execution.
	//
	// If you don't want to run validations, use [App.SaveNoValidateWithContext()].
	SaveWithContext(ctx context.Context, model Model) error

	// SaveNoValidate saves the specified model into the regular app database without performing validations.
	//
	// If you want to also run validations before persisting, use [App.Save()].
	SaveNoValidate(model Model) error

	// SaveNoValidateWithContext is the same as [App.SaveNoValidate()]
	// but allows specifying a context to limit the db execution.
	//
	// If you want to also run validations before persisting, use [App.SaveWithContext()].
	SaveNoValidateWithContext(ctx context.Context, model Model) error

	// AuxSave validates and saves the specified model into the auxiliary app database.
	//
	// If you don't want to run validations, use [App.AuxSaveNoValidate()].
	AuxSave(model Model) error

	// AuxSaveWithContext is the same as [App.AuxSave()] but allows specifying a context to limit the db execution.
	//
	// If you don't want to run validations, use [App.AuxSaveNoValidateWithContext()].
	AuxSaveWithContext(ctx context.Context, model Model) error

	// AuxSaveNoValidate saves the specified model into the auxiliary app database without performing validations.
	//
	// If you want to also run validations before persisting, use [App.AuxSave()].
	AuxSaveNoValidate(model Model) error

	// AuxSaveNoValidateWithContext is the same as [App.AuxSaveNoValidate()]
	// but allows specifying a context to limit the db execution.
	//
	// If you want to also run validations before persisting, use [App.AuxSaveWithContext()].
	AuxSaveNoValidateWithContext(ctx context.Context, model Model) error

	// Validate triggers the OnModelValidate hook for the specified model.
	Validate(model Model) error

	// ValidateWithContext is the same as Validate but allows specifying the ModelEvent context.
	ValidateWithContext(ctx context.Context, model Model) error

	// RunInTransaction wraps fn into a transaction for the regular app database.
	//
	// It is safe to nest RunInTransaction calls as long as you use the callback's txApp.
	RunInTransaction(fn func(txApp App) error) error

	// AuxRunInTransaction wraps fn into a transaction for the auxiliary app database.
	//
	// It is safe to nest RunInTransaction calls as long as you use the callback's txApp.
	AuxRunInTransaction(fn func(txApp App) error) error

	// ---------------------------------------------------------------

	// LogQuery returns a new Log select query.
	LogQuery() *dbx.SelectQuery

	// FindLogById finds a single Log entry by its id.
	FindLogById(id string) (*Log, error)

	// LogsStatsItem returns hourly grouped logs statistics.
	LogsStats(expr dbx.Expression) ([]*LogsStatsItem, error)

	// DeleteOldLogs delete all logs that are created before createdBefore.
	DeleteOldLogs(createdBefore time.Time) error

	// ---------------------------------------------------------------

	// CollectionQuery returns a new Collection select query.
	CollectionQuery() *dbx.SelectQuery

	// FindCollections finds all collections by the given type(s).
	//
	// If collectionTypes is not set, it returns all collections.
	//
	// Example:
	//
	//	app.FindAllCollections() // all collections
	//	app.FindAllCollections("auth", "view") // only auth and view collections
	FindAllCollections(collectionTypes ...string) ([]*Collection, error)

	// ReloadCachedCollections fetches all collections and caches them into the app store.
	ReloadCachedCollections() error

	// FindCollectionByNameOrId finds a single collection by its name (case insensitive) or id.s
	FindCollectionByNameOrId(nameOrId string) (*Collection, error)

	// FindCachedCollectionByNameOrId is similar to [App.FindCollectionByNameOrId]
	// but retrieves the Collection from the app cache instead of making a db call.
	//
	// NB! This method is suitable for read-only Collection operations.
	//
	// Returns [sql.ErrNoRows] if no Collection is found for consistency
	// with the [App.FindCollectionByNameOrId] method.
	//
	// If you plan making changes to the returned Collection model,
	// use [App.FindCollectionByNameOrId] instead.
	//
	// Caveats:
	//
	//   - The returned Collection should be used only for read-only operations.
	//     Avoid directly modifying the returned cached Collection as it will affect
	//     the global cached value even if you don't persist the changes in the database!
	//   - If you are updating a Collection in a transaction and then call this method before commit,
	//     it'll return the cached Collection state and not the one from the uncommitted transaction.
	//   - The cache is automatically updated on collections db change (create/update/delete).
	//     To manually reload the cache you can call [App.ReloadCachedCollections]
	FindCachedCollectionByNameOrId(nameOrId string) (*Collection, error)

	// FindCollectionReferences returns information for all relation
	// fields referencing the provided collection.
	//
	// If the provided collection has reference to itself then it will be
	// also included in the result. To exclude it, pass the collection id
	// as the excludeIds argument.
	FindCollectionReferences(collection *Collection, excludeIds ...string) (map[*Collection][]Field, error)

	// FindCachedCollectionReferences is similar to [App.FindCollectionReferences]
	// but retrieves the Collection from the app cache instead of making a db call.
	//
	// NB! This method is suitable for read-only Collection operations.
	//
	// If you plan making changes to the returned Collection model,
	// use [App.FindCollectionReferences] instead.
	//
	// Caveats:
	//
	//   - The returned Collection should be used only for read-only operations.
	//     Avoid directly modifying the returned cached Collection as it will affect
	//     the global cached value even if you don't persist the changes in the database!
	//   - If you are updating a Collection in a transaction and then call this method before commit,
	//     it'll return the cached Collection state and not the one from the uncommitted transaction.
	//   - The cache is automatically updated on collections db change (create/update/delete).
	//     To manually reload the cache you can call [App.ReloadCachedCollections].
	FindCachedCollectionReferences(collection *Collection, excludeIds ...string) (map[*Collection][]Field, error)

	// IsCollectionNameUnique checks that there is no existing collection
	// with the provided name (case insensitive!).
	//
	// Note: case insensitive check because the name is used also as
	// table name for the records.
	IsCollectionNameUnique(name string, excludeIds ...string) bool

	// TruncateCollection deletes all records associated with the provided collection.
	//
	// The truncate operation is executed in a single transaction,
	// aka. either everything is deleted or none.
	//
	// Note that this method will also trigger the records related
	// cascade and file delete actions.
	TruncateCollection(collection *Collection) error

	// ImportCollections imports the provided collections data in a single transaction.
	//
	// For existing matching collections, the imported data is unmarshaled on top of the existing model.
	//
	// NB! If deleteMissing is true, ALL NON-SYSTEM COLLECTIONS AND SCHEMA FIELDS,
	// that are not present in the imported configuration, WILL BE DELETED
	// (this includes their related records data).
	ImportCollections(toImport []map[string]any, deleteMissing bool) error

	// ImportCollectionsByMarshaledJSON is the same as [ImportCollections]
	// but accept marshaled json array as import data (usually used for the autogenerated snapshots).
	ImportCollectionsByMarshaledJSON(rawSliceOfMaps []byte, deleteMissing bool) error

	// SyncRecordTableSchema compares the two provided collections
	// and applies the necessary related record table changes.
	//
	// If oldCollection is null, then only newCollection is used to create the record table.
	//
	// This method is automatically invoked as part of a collection create/update/delete operation.
	SyncRecordTableSchema(newCollection *Collection, oldCollection *Collection) error

	// ---------------------------------------------------------------

	// FindAllExternalAuthsByRecord returns all ExternalAuth models
	// linked to the provided auth record.
	FindAllExternalAuthsByRecord(authRecord *Record) ([]*ExternalAuth, error)

	// FindAllExternalAuthsByCollection returns all ExternalAuth models
	// linked to the provided auth collection.
	FindAllExternalAuthsByCollection(collection *Collection) ([]*ExternalAuth, error)

	// FindFirstExternalAuthByExpr returns the first available (the most recent created)
	// ExternalAuth model that satisfies the non-nil expression.
	FindFirstExternalAuthByExpr(expr dbx.Expression) (*ExternalAuth, error)

	// ---------------------------------------------------------------

	// FindAllMFAsByRecord returns all MFA models linked to the provided auth record.
	FindAllMFAsByRecord(authRecord *Record) ([]*MFA, error)

	// FindAllMFAsByCollection returns all MFA models linked to the provided collection.
	FindAllMFAsByCollection(collection *Collection) ([]*MFA, error)

	// FindMFAById returns a single MFA model by its id.
	FindMFAById(id string) (*MFA, error)

	// DeleteAllMFAsByRecord deletes all MFA models associated with the provided record.
	//
	// Returns a combined error with the failed deletes.
	DeleteAllMFAsByRecord(authRecord *Record) error

	// DeleteExpiredMFAs deletes the expired MFAs for all auth collections.
	DeleteExpiredMFAs() error

	// ---------------------------------------------------------------

	// FindAllOTPsByRecord returns all OTP models linked to the provided auth record.
	FindAllOTPsByRecord(authRecord *Record) ([]*OTP, error)

	// FindAllOTPsByCollection returns all OTP models linked to the provided collection.
	FindAllOTPsByCollection(collection *Collection) ([]*OTP, error)

	// FindOTPById returns a single OTP model by its id.
	FindOTPById(id string) (*OTP, error)

	// DeleteAllOTPsByRecord deletes all OTP models associated with the provided record.
	//
	// Returns a combined error with the failed deletes.
	DeleteAllOTPsByRecord(authRecord *Record) error

	// DeleteExpiredOTPs deletes the expired OTPs for all auth collections.
	DeleteExpiredOTPs() error

	// ---------------------------------------------------------------

	// FindAllAuthOriginsByRecord returns all AuthOrigin models linked to the provided auth record (in DESC order).
	FindAllAuthOriginsByRecord(authRecord *Record) ([]*AuthOrigin, error)

	// FindAllAuthOriginsByCollection returns all AuthOrigin models linked to the provided collection (in DESC order).
	FindAllAuthOriginsByCollection(collection *Collection) ([]*AuthOrigin, error)

	// FindAuthOriginById returns a single AuthOrigin model by its id.
	FindAuthOriginById(id string) (*AuthOrigin, error)

	// FindAuthOriginByRecordAndFingerprint returns a single AuthOrigin model
	// by its authRecord relation and fingerprint.
	FindAuthOriginByRecordAndFingerprint(authRecord *Record, fingerprint string) (*AuthOrigin, error)

	// DeleteAllAuthOriginsByRecord deletes all AuthOrigin models associated with the provided record.
	//
	// Returns a combined error with the failed deletes.
	DeleteAllAuthOriginsByRecord(authRecord *Record) error

	// ---------------------------------------------------------------

	// RecordQuery returns a new Record select query from a collection model, id or name.
	//
	// In case a collection id or name is provided and that collection doesn't
	// actually exists, the generated query will be created with a cancelled context
	// and will fail once an executor (Row(), One(), All(), etc.) is called.
	RecordQuery(collectionModelOrIdentifier any) *dbx.SelectQuery

	// FindRecordById finds the Record model by its id.
	FindRecordById(collectionModelOrIdentifier any, recordId string, optFilters ...func(q *dbx.SelectQuery) error) (*Record, error)

	// FindRecordsByIds finds all records by the specified ids.
	// If no records are found, returns an empty slice.
	FindRecordsByIds(collectionModelOrIdentifier any, recordIds []string, optFilters ...func(q *dbx.SelectQuery) error) ([]*Record, error)

	// FindAllRecords finds all records matching specified db expressions.
	//
	// Returns all collection records if no expression is provided.
	//
	// Returns an empty slice if no records are found.
	//
	// Example:
	//
	//	// no extra expressions
	//	app.FindAllRecords("example")
	//
	//	// with extra expressions
	//	expr1 := dbx.HashExp{"email": "test@example.com"}
	//	expr2 := dbx.NewExp("LOWER(username) = {:username}", dbx.Params{"username": "test"})
	//	app.FindAllRecords("example", expr1, expr2)
	FindAllRecords(collectionModelOrIdentifier any, exprs ...dbx.Expression) ([]*Record, error)

	// FindFirstRecordByData returns the first found record matching
	// the provided key-value pair.
	FindFirstRecordByData(collectionModelOrIdentifier any, key string, value any) (*Record, error)

	// FindRecordsByFilter returns limit number of records matching the
	// provided string filter.
	//
	// NB! Use the last "params" argument to bind untrusted user variables!
	//
	// The filter argument is optional and can be empty string to target
	// all available records.
	//
	// The sort argument is optional and can be empty string OR the same format
	// used in the web APIs, ex. "-created,title".
	//
	// If the limit argument is <= 0, no limit is applied to the query and
	// all matching records are returned.
	//
	// Returns an empty slice if no records are found.
	//
	// Example:
	//
	//	app.FindRecordsByFilter(
	//		"posts",
	//		"title ~ {:title} && visible = {:visible}",
	//		"-created",
	//		10,
	//		0,
	//		dbx.Params{"title": "lorem ipsum", "visible": true}
	//	)
	FindRecordsByFilter(
		collectionModelOrIdentifier any,
		filter string,
		sort string,
		limit int,
		offset int,
		params ...dbx.Params,
	) ([]*Record, error)

	// FindFirstRecordByFilter returns the first available record matching the provided filter (if any).
	//
	// NB! Use the last params argument to bind untrusted user variables!
	//
	// Returns sql.ErrNoRows if no record is found.
	//
	// Example:
	//
	//	app.FindFirstRecordByFilter("posts", "")
	//	app.FindFirstRecordByFilter("posts", "slug={:slug} && status='public'", dbx.Params{"slug": "test"})
	FindFirstRecordByFilter(
		collectionModelOrIdentifier any,
		filter string,
		params ...dbx.Params,
	) (*Record, error)

	// CountRecords returns the total number of records in a collection.
	CountRecords(collectionModelOrIdentifier any, exprs ...dbx.Expression) (int64, error)

	// FindAuthRecordByToken finds the auth record associated with the provided JWT
	// (auth, file, verifyEmail, changeEmail, passwordReset types).
	//
	// Optionally specify a list of validTypes to check tokens only from those types.
	//
	// Returns an error if the JWT is invalid, expired or not associated to an auth collection record.
	FindAuthRecordByToken(token string, validTypes ...string) (*Record, error)

	// FindAuthRecordByEmail finds the auth record associated with the provided email.
	//
	// Returns an error if it is not an auth collection or the record is not found.
	FindAuthRecordByEmail(collectionModelOrIdentifier any, email string) (*Record, error)

	// CanAccessRecord checks if a record is allowed to be accessed by the
	// specified requestInfo and accessRule.
	//
	// Rule and db checks are ignored in case requestInfo.Auth is a superuser.
	//
	// The returned error indicate that something unexpected happened during
	// the check (eg. invalid rule or db query error).
	//
	// The method always return false on invalid rule or db query error.
	//
	// Example:
	//
	//	requestInfo, _ := e.RequestInfo()
	//	record, _ := app.FindRecordById("example", "RECORD_ID")
	//	rule := types.Pointer("@request.auth.id != '' || status = 'public'")
	//	// ... or use one of the record collection's rule, eg. record.Collection().ViewRule
	//
	//	if ok, _ := app.CanAccessRecord(record, requestInfo, rule); ok { ... }
	CanAccessRecord(record *Record, requestInfo *RequestInfo, accessRule *string) (bool, error)

	// ExpandRecord expands the relations of a single Record model.
	//
	// If optFetchFunc is not set, then a default function will be used
	// that returns all relation records.
	//
	// Returns a map with the failed expand parameters and their errors.
	ExpandRecord(record *Record, expands []string, optFetchFunc ExpandFetchFunc) map[string]error

	// ExpandRecords expands the relations of the provided Record models list.
	//
	// If optFetchFunc is not set, then a default function will be used
	// that returns all relation records.
	//
	// Returns a map with the failed expand parameters and their errors.
	ExpandRecords(records []*Record, expands []string, optFetchFunc ExpandFetchFunc) map[string]error

	// ---------------------------------------------------------------
	// App event hooks
	// ---------------------------------------------------------------

	// OnBootstrap hook is triggered when initializing the main application
	// resources (db, app settings, etc).
	OnBootstrap() *hook.Hook[*BootstrapEvent]

	// OnServe hook is triggered when the app web server is started
	// (after starting the TCP listener but before initializing the blocking serve task),
	// allowing you to adjust its options and attach new routes or middlewares.
	OnServe() *hook.Hook[*ServeEvent]

	// OnTerminate hook is triggered when the app is in the process
	// of being terminated (ex. on SIGTERM signal).
	//
	// Note that the app could be terminated abruptly without awaiting the hook completion.
	OnTerminate() *hook.Hook[*TerminateEvent]

	// OnBackupCreate hook is triggered on each [App.CreateBackup] call.
	OnBackupCreate() *hook.Hook[*BackupEvent]

	// OnBackupRestore hook is triggered before app backup restore (aka. [App.RestoreBackup] call).
	//
	// Note that by default on success the application is restarted and the after state of the hook is ignored.
	OnBackupRestore() *hook.Hook[*BackupEvent]

	// ---------------------------------------------------------------
	// DB models event hooks
	// ---------------------------------------------------------------

	// OnModelValidate is triggered every time when a model is being validated
	// (e.g. triggered by App.Validate() or App.Save()).
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelValidate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// ---------------------------------------------------------------

	// OnModelCreate is triggered every time when a new model is being created
	// (e.g. triggered by App.Save()).
	//
	// Operations BEFORE the e.Next() execute before the model validation
	// and the INSERT DB statement.
	//
	// Operations AFTER the e.Next() execute after the model validation
	// and the INSERT DB statement.
	//
	// Note that successful execution doesn't guarantee that the model
	// is persisted in the database since its wrapping transaction may
	// not have been committed yet.
	// If you want to listen to only the actual persisted events, you can
	// bind to [OnModelAfterCreateSuccess] or [OnModelAfterCreateError] hooks.
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelCreate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelCreateExecute is triggered after successful Model validation
	// and right before the model INSERT DB statement execution.
	//
	// Usually it is triggered as part of the App.Save() in the following firing order:
	// OnModelCreate {
	//    -> OnModelValidate (skipped with App.SaveNoValidate())
	//    -> OnModelCreateExecute
	// }
	//
	// Note that successful execution doesn't guarantee that the model
	// is persisted in the database since its wrapping transaction may have been
	// committed yet.
	// If you want to listen to only the actual persisted events,
	// you can bind to [OnModelAfterCreateSuccess] or [OnModelAfterCreateError] hooks.
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelCreateExecute(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterCreateSuccess is triggered after each successful
	// Model DB create persistence.
	//
	// Note that when a Model is persisted as part of a transaction,
	// this hook is delayed and executed only AFTER the transaction has been committed.
	// This hook is NOT triggered in case the transaction rollbacks
	// (aka. when the model wasn't persisted).
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelAfterCreateSuccess(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterCreateError is triggered after each failed
	// Model DB create persistence.
	//
	// Note that the execution of this hook is either immediate or delayed
	// depending on the error:
	//   - "immediate" on App.Save() failure
	//   - "delayed" on transaction rollback
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelAfterCreateError(tags ...string) *hook.TaggedHook[*ModelErrorEvent]

	// ---------------------------------------------------------------

	// OnModelUpdate is triggered every time when a new model is being updated
	// (e.g. triggered by App.Save()).
	//
	// Operations BEFORE the e.Next() execute before the model validation
	// and the UPDATE DB statement.
	//
	// Operations AFTER the e.Next() execute after the model validation
	// and the UPDATE DB statement.
	//
	// Note that successful execution doesn't guarantee that the model
	// is persisted in the database since its wrapping transaction may
	// not have been committed yet.
	// If you want to listen to only the actual persisted events, you can
	// bind to [OnModelAfterUpdateSuccess] or [OnModelAfterUpdateError] hooks.
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelUpdate(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelUpdateExecute is triggered after successful Model validation
	// and right before the model UPDATE DB statement execution.
	//
	// Usually it is triggered as part of the App.Save() in the following firing order:
	// OnModelUpdate {
	//    -> OnModelValidate (skipped with App.SaveNoValidate())
	//    -> OnModelUpdateExecute
	// }
	//
	// Note that successful execution doesn't guarantee that the model
	// is persisted in the database since its wrapping transaction may have been
	// committed yet.
	// If you want to listen to only the actual persisted events,
	// you can bind to [OnModelAfterUpdateSuccess] or [OnModelAfterUpdateError] hooks.
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelUpdateExecute(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterUpdateSuccess is triggered after each successful
	// Model DB update persistence.
	//
	// Note that when a Model is persisted as part of a transaction,
	// this hook is delayed and executed only AFTER the transaction has been committed.
	// This hook is NOT triggered in case the transaction rollbacks
	// (aka. when the model changes weren't persisted).
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelAfterUpdateSuccess(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterUpdateError is triggered after each failed
	// Model DB update persistence.
	//
	// Note that the execution of this hook is either immediate or delayed
	// depending on the error:
	//   - "immediate" on App.Save() failure
	//   - "delayed" on transaction rollback
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelAfterUpdateError(tags ...string) *hook.TaggedHook[*ModelErrorEvent]

	// ---------------------------------------------------------------

	// OnModelDelete is triggered every time when a new model is being deleted
	// (e.g. triggered by App.Delete()).
	//
	// Note that successful execution doesn't guarantee that the model
	// is deleted from the database since its wrapping transaction may
	// not have been committed yet.
	// If you want to listen to only the actual persisted deleted events, you can
	// bind to [OnModelAfterDeleteSuccess] or [OnModelAfterDeleteError] hooks.
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelDelete(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelUpdateExecute is triggered right before the model
	// DELETE DB statement execution.
	//
	// Usually it is triggered as part of the App.Delete() in the following firing order:
	// OnModelDelete {
	//    -> (internal delete checks)
	//    -> OnModelDeleteExecute
	// }
	//
	// Note that successful execution doesn't guarantee that the model
	// is deleted from the database since its wrapping transaction may
	// not have been committed yet.
	// If you want to listen to only the actual persisted deleted events, you can
	// bind to [OnModelAfterDeleteSuccess] or [OnModelAfterDeleteError] hooks.
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelDeleteExecute(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterDeleteSuccess is triggered after each successful
	// Model DB delete persistence.
	//
	// Note that when a Model is deleted as part of a transaction,
	// this hook is delayed and executed only AFTER the transaction has been committed.
	// This hook is NOT triggered in case the transaction rollbacks
	// (aka. when the model delete wasn't persisted).
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelAfterDeleteSuccess(tags ...string) *hook.TaggedHook[*ModelEvent]

	// OnModelAfterDeleteError is triggered after each failed
	// Model DB delete persistence.
	//
	// Note that the execution of this hook is either immediate or delayed
	// depending on the error:
	//   - "immediate" on App.Delete() failure
	//   - "delayed" on transaction rollback
	//
	// For convenience, if you want to listen to only the Record models
	// events without doing manual type assertion, you can attach to the OnRecord* proxy hooks.
	//
	// If the optional "tags" list (Collection id/name, Model table name, etc.) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnModelAfterDeleteError(tags ...string) *hook.TaggedHook[*ModelErrorEvent]

	// ---------------------------------------------------------------
	// Record models event hooks
	// ---------------------------------------------------------------

	// OnRecordEnrich is triggered every time when a record is enriched
	// (as part of the builtin Record responses, during realtime message seriazation, or when [apis.EnrichRecord] is invoked).
	//
	// It could be used for example to redact/hide or add computed temporary
	// Record model props only for the specific request info. For example:
	//
	//  app.OnRecordEnrich("posts").BindFunc(func(e core.*RecordEnrichEvent) {
	//      // hide one or more fields
	//      e.Record.Hide("role")
	//
	//      // add new custom field for registered users
	//      if e.RequestInfo.Auth != nil && e.RequestInfo.Auth.Collection().Name == "users" {
	//          e.Record.WithCustomData(true) // for security requires explicitly allowing it
	//          e.Record.Set("computedScore", e.Record.GetInt("score") * e.RequestInfo.Auth.GetInt("baseScore"))
	//      }
	//
	//      return e.Next()
	//  })
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordEnrich(tags ...string) *hook.TaggedHook[*RecordEnrichEvent]

	// OnRecordValidate is a Record proxy model hook of [OnModelValidate].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordValidate(tags ...string) *hook.TaggedHook[*RecordEvent]

	// ---------------------------------------------------------------

	// OnRecordCreate is a Record proxy model hook of [OnModelCreate].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordCreate(tags ...string) *hook.TaggedHook[*RecordEvent]

	// OnRecordCreateExecute is a Record proxy model hook of [OnModelCreateExecute].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordCreateExecute(tags ...string) *hook.TaggedHook[*RecordEvent]

	// OnRecordAfterCreateSuccess is a Record proxy model hook of [OnModelAfterCreateSuccess].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterCreateSuccess(tags ...string) *hook.TaggedHook[*RecordEvent]

	// OnRecordAfterCreateError is a Record proxy model hook of [OnModelAfterCreateError].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterCreateError(tags ...string) *hook.TaggedHook[*RecordErrorEvent]

	// ---------------------------------------------------------------

	// OnRecordUpdate is a Record proxy model hook of [OnModelUpdate].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordUpdate(tags ...string) *hook.TaggedHook[*RecordEvent]

	// OnRecordUpdateExecute is a Record proxy model hook of [OnModelUpdateExecute].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordUpdateExecute(tags ...string) *hook.TaggedHook[*RecordEvent]

	// OnRecordAfterUpdateSuccess is a Record proxy model hook of [OnModelAfterUpdateSuccess].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterUpdateSuccess(tags ...string) *hook.TaggedHook[*RecordEvent]

	// OnRecordAfterUpdateError is a Record proxy model hook of [OnModelAfterUpdateError].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterUpdateError(tags ...string) *hook.TaggedHook[*RecordErrorEvent]

	// ---------------------------------------------------------------

	// OnRecordDelete is a Record proxy model hook of [OnModelDelete].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordDelete(tags ...string) *hook.TaggedHook[*RecordEvent]

	// OnRecordDeleteExecute is a Record proxy model hook of [OnModelDeleteExecute].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordDeleteExecute(tags ...string) *hook.TaggedHook[*RecordEvent]

	// OnRecordAfterDeleteSuccess is a Record proxy model hook of [OnModelAfterDeleteSuccess].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterDeleteSuccess(tags ...string) *hook.TaggedHook[*RecordEvent]

	// OnRecordAfterDeleteError is a Record proxy model hook of [OnModelAfterDeleteError].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAfterDeleteError(tags ...string) *hook.TaggedHook[*RecordErrorEvent]

	// ---------------------------------------------------------------
	// Collection models event hooks
	// ---------------------------------------------------------------

	// OnCollectionValidate is a Collection proxy model hook of [OnModelValidate].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionValidate(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// ---------------------------------------------------------------

	// OnCollectionCreate is a Collection proxy model hook of [OnModelCreate].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionCreate(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// OnCollectionCreateExecute is a Collection proxy model hook of [OnModelCreateExecute].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionCreateExecute(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// OnCollectionAfterCreateSuccess is a Collection proxy model hook of [OnModelAfterCreateSuccess].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionAfterCreateSuccess(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// OnCollectionAfterCreateError is a Collection proxy model hook of [OnModelAfterCreateError].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionAfterCreateError(tags ...string) *hook.TaggedHook[*CollectionErrorEvent]

	// ---------------------------------------------------------------

	// OnCollectionUpdate is a Collection proxy model hook of [OnModelUpdate].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionUpdate(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// OnCollectionUpdateExecute is a Collection proxy model hook of [OnModelUpdateExecute].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionUpdateExecute(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// OnCollectionAfterUpdateSuccess is a Collection proxy model hook of [OnModelAfterUpdateSuccess].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionAfterUpdateSuccess(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// OnCollectionAfterUpdateError is a Collection proxy model hook of [OnModelAfterUpdateError].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionAfterUpdateError(tags ...string) *hook.TaggedHook[*CollectionErrorEvent]

	// ---------------------------------------------------------------

	// OnCollectionDelete is a Collection proxy model hook of [OnModelDelete].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionDelete(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// OnCollectionDeleteExecute is a Collection proxy model hook of [OnModelDeleteExecute].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionDeleteExecute(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// OnCollectionAfterDeleteSuccess is a Collection proxy model hook of [OnModelAfterDeleteSuccess].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionAfterDeleteSuccess(tags ...string) *hook.TaggedHook[*CollectionEvent]

	// OnCollectionAfterDeleteError is a Collection proxy model hook of [OnModelAfterDeleteError].
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnCollectionAfterDeleteError(tags ...string) *hook.TaggedHook[*CollectionErrorEvent]

	// ---------------------------------------------------------------
	// Mailer event hooks
	// ---------------------------------------------------------------

	// OnMailerSend hook is triggered every time when a new email is
	// being sent using the [App.NewMailClient()] instance.
	//
	// It allows intercepting the email message or to use a custom mailer client.
	OnMailerSend() *hook.Hook[*MailerEvent]

	// OnMailerRecordAuthAlertSend hook is triggered when
	// sending a new device login auth alert email, allowing you to
	// intercept and customize the email message that is being sent.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerRecordAuthAlertSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// OnMailerBeforeRecordResetPasswordSend hook is triggered when
	// sending a password reset email to an auth record, allowing
	// you to intercept and customize the email message that is being sent.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerRecordPasswordResetSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// OnMailerBeforeRecordVerificationSend hook is triggered when
	// sending a verification email to an auth record, allowing
	// you to intercept and customize the email message that is being sent.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerRecordVerificationSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// OnMailerRecordEmailChangeSend hook is triggered when sending a
	// confirmation new address email to an auth record, allowing
	// you to intercept and customize the email message that is being sent.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerRecordEmailChangeSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// OnMailerRecordOTPSend hook is triggered when sending an OTP email
	// to an auth record, allowing you to intercept and customize the
	// email message that is being sent.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnMailerRecordOTPSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent]

	// ---------------------------------------------------------------
	// Realtime API event hooks
	// ---------------------------------------------------------------

	// OnRealtimeConnectRequest hook is triggered when establishing the SSE client connection.
	//
	// Any execution after e.Next() of a hook handler happens after the client disconnects.
	OnRealtimeConnectRequest() *hook.Hook[*RealtimeConnectRequestEvent]

	// OnRealtimeMessageSend hook is triggered when sending an SSE message to a client.
	OnRealtimeMessageSend() *hook.Hook[*RealtimeMessageEvent]

	// OnRealtimeSubscribeRequest hook is triggered when updating the
	// client subscriptions, allowing you to further validate and
	// modify the submitted change.
	OnRealtimeSubscribeRequest() *hook.Hook[*RealtimeSubscribeRequestEvent]

	// ---------------------------------------------------------------
	// Settings API event hooks
	// ---------------------------------------------------------------

	// OnSettingsListRequest hook is triggered on each API Settings list request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnSettingsListRequest() *hook.Hook[*SettingsListRequestEvent]

	// OnSettingsUpdateRequest hook is triggered on each API Settings update request.
	//
	// Could be used to additionally validate the request data or
	// implement completely different persistence behavior.
	OnSettingsUpdateRequest() *hook.Hook[*SettingsUpdateRequestEvent]

	// OnSettingsReload hook is triggered every time when the App.Settings()
	// is being replaced with a new state.
	//
	// Calling App.Settings() after e.Next() returns the new state.
	OnSettingsReload() *hook.Hook[*SettingsReloadEvent]

	// ---------------------------------------------------------------
	// File API event hooks
	// ---------------------------------------------------------------

	// OnFileDownloadRequest hook is triggered before each API File download request.
	//
	// Could be used to validate or modify the file response before
	// returning it to the client.
	OnFileDownloadRequest(tags ...string) *hook.TaggedHook[*FileDownloadRequestEvent]

	// OnFileBeforeTokenRequest hook is triggered on each auth file token API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnFileTokenRequest(tags ...string) *hook.TaggedHook[*FileTokenRequestEvent]

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
	OnRecordAuthRequest(tags ...string) *hook.TaggedHook[*RecordAuthRequestEvent]

	// OnRecordAuthWithPasswordRequest hook is triggered on each
	// Record auth with password API request.
	//
	// [RecordAuthWithPasswordRequestEvent.Record] could be nil if no matching identity is found, allowing
	// you to manually locate a different Record model (by reassigning [RecordAuthWithPasswordRequestEvent.Record]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*RecordAuthWithPasswordRequestEvent]

	// OnRecordAuthWithOAuth2Request hook is triggered on each Record
	// OAuth2 sign-in/sign-up API request (after token exchange and before external provider linking).
	//
	// If [RecordAuthWithOAuth2RequestEvent.Record] is not set, then the OAuth2
	// request will try to create a new auth Record.
	//
	// To assign or link a different existing record model you can
	// change the [RecordAuthWithOAuth2RequestEvent.Record] field.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAuthWithOAuth2Request(tags ...string) *hook.TaggedHook[*RecordAuthWithOAuth2RequestEvent]

	// OnRecordAuthRefreshRequest hook is triggered on each Record
	// auth refresh API request (right before generating a new auth token).
	//
	// Could be used to additionally validate the request data or implement
	// completely different auth refresh behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAuthRefreshRequest(tags ...string) *hook.TaggedHook[*RecordAuthRefreshRequestEvent]

	// OnRecordRequestPasswordResetRequest hook is triggered on
	// each Record request password reset API request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different password reset behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordRequestPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordRequestPasswordResetRequestEvent]

	// OnRecordConfirmPasswordResetRequest hook is triggered on
	// each Record confirm password reset API request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordConfirmPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordConfirmPasswordResetRequestEvent]

	// OnRecordRequestVerificationRequest hook is triggered on
	// each Record request verification API request.
	//
	// Could be used to additionally validate the loaded request data or implement
	// completely different verification behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordRequestVerificationRequest(tags ...string) *hook.TaggedHook[*RecordRequestVerificationRequestEvent]

	// OnRecordConfirmVerificationRequest hook is triggered on each
	// Record confirm verification API request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordConfirmVerificationRequest(tags ...string) *hook.TaggedHook[*RecordConfirmVerificationRequestEvent]

	// OnRecordRequestEmailChangeRequest hook is triggered on each
	// Record request email change API request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different request email change behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordRequestEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordRequestEmailChangeRequestEvent]

	// OnRecordConfirmEmailChangeRequest hook is triggered on each
	// Record confirm email change API request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordConfirmEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordConfirmEmailChangeRequestEvent]

	// OnRecordRequestOTPRequest hook is triggered on each Record
	// request OTP API request.
	//
	// [RecordCreateOTPRequestEvent.Record] could be nil if no matching identity is found, allowing
	// you to manually create or locate a different Record model (by reassigning [RecordCreateOTPRequestEvent.Record]).
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordRequestOTPRequest(tags ...string) *hook.TaggedHook[*RecordCreateOTPRequestEvent]

	// OnRecordAuthWithOTPRequest hook is triggered on each Record
	// auth with OTP API request.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordAuthWithOTPRequest(tags ...string) *hook.TaggedHook[*RecordAuthWithOTPRequestEvent]

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
	OnRecordsListRequest(tags ...string) *hook.TaggedHook[*RecordsListRequestEvent]

	// OnRecordViewRequest hook is triggered on each API Record view request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordViewRequest(tags ...string) *hook.TaggedHook[*RecordRequestEvent]

	// OnRecordCreateRequest hook is triggered on each API Record create request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordCreateRequest(tags ...string) *hook.TaggedHook[*RecordRequestEvent]

	// OnRecordUpdateRequest hook is triggered on each API Record update request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordUpdateRequest(tags ...string) *hook.TaggedHook[*RecordRequestEvent]

	// OnRecordDeleteRequest hook is triggered on each API Record delete request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior.
	//
	// If the optional "tags" list (Collection ids or names) is specified,
	// then all event handlers registered via the created hook will be
	// triggered and called only if their event data origin matches the tags.
	OnRecordDeleteRequest(tags ...string) *hook.TaggedHook[*RecordRequestEvent]

	// ---------------------------------------------------------------
	// Collection API event hooks
	// ---------------------------------------------------------------

	// OnCollectionsListRequest hook is triggered on each API Collections list request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnCollectionsListRequest() *hook.Hook[*CollectionsListRequestEvent]

	// OnCollectionViewRequest hook is triggered on each API Collection view request.
	//
	// Could be used to validate or modify the response before returning it to the client.
	OnCollectionViewRequest() *hook.Hook[*CollectionRequestEvent]

	// OnCollectionCreateRequest hook is triggered on each API Collection create request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
	OnCollectionCreateRequest() *hook.Hook[*CollectionRequestEvent]

	// OnCollectionUpdateRequest hook is triggered on each API Collection update request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different persistence behavior.
	OnCollectionUpdateRequest() *hook.Hook[*CollectionRequestEvent]

	// OnCollectionDeleteRequest hook is triggered on each API Collection delete request.
	//
	// Could be used to additionally validate the request data or implement
	// completely different delete behavior.
	OnCollectionDeleteRequest() *hook.Hook[*CollectionRequestEvent]

	// OnCollectionsBeforeImportRequest hook is triggered on each API
	// collections import request.
	//
	// Could be used to additionally validate the imported collections or
	// to implement completely different import behavior.
	OnCollectionsImportRequest() *hook.Hook[*CollectionsImportRequestEvent]

	// ---------------------------------------------------------------
	// Batch API event hooks
	// ---------------------------------------------------------------

	// OnBatchRequest hook is triggered on each API batch request.
	//
	// Could be used to additionally validate or modify the submitted batch requests.
	OnBatchRequest() *hook.Hook[*BatchRequestEvent]
}
