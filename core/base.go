package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/logger"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
	"golang.org/x/sync/semaphore"
)

const (
	DefaultDataMaxOpenConns int           = 120
	DefaultDataMaxIdleConns int           = 15
	DefaultAuxMaxOpenConns  int           = 20
	DefaultAuxMaxIdleConns  int           = 3
	DefaultQueryTimeout     time.Duration = 30 * time.Second

	LocalStorageDirName       string = "storage"
	LocalBackupsDirName       string = "backups"
	LocalTempDirName          string = ".pb_temp_to_delete" // temp pb_data sub directory that will be deleted on each app.Bootstrap()
	LocalAutocertCacheDirName string = ".autocert_cache"

	// @todo consider removing after backups refactoring
	lostFoundDirName string = "lost+found"
)

// FilesManager defines an interface with common methods that files manager models should implement.
type FilesManager interface {
	// BaseFilesPath returns the storage dir path used by the interface instance.
	BaseFilesPath() string
}

// DBConnectFunc defines a database connection initialization function.
type DBConnectFunc func(dbPath string) (*dbx.DB, error)

// BaseAppConfig defines a BaseApp configuration option
type BaseAppConfig struct {
	DBConnect        DBConnectFunc
	DataDir          string
	EncryptionEnv    string
	QueryTimeout     time.Duration
	DataMaxOpenConns int
	DataMaxIdleConns int
	AuxMaxOpenConns  int
	AuxMaxIdleConns  int
	IsDev            bool
}

// ensures that the BaseApp implements the App interface.
var _ App = (*BaseApp)(nil)

// BaseApp implements core.App and defines the base PocketBase app structure.
type BaseApp struct {
	config              *BaseAppConfig
	txInfo              *TxAppInfo
	store               *store.Store[string, any]
	cron                *cron.Cron
	settings            *Settings
	subscriptionsBroker *subscriptions.Broker
	logger              *slog.Logger
	concurrentDB        dbx.Builder
	nonconcurrentDB     dbx.Builder
	auxConcurrentDB     dbx.Builder
	auxNonconcurrentDB  dbx.Builder

	// app event hooks
	onBootstrap     *hook.Hook[*BootstrapEvent]
	onServe         *hook.Hook[*ServeEvent]
	onTerminate     *hook.Hook[*TerminateEvent]
	onBackupCreate  *hook.Hook[*BackupEvent]
	onBackupRestore *hook.Hook[*BackupEvent]

	// db model hooks
	onModelValidate           *hook.Hook[*ModelEvent]
	onModelCreate             *hook.Hook[*ModelEvent]
	onModelCreateExecute      *hook.Hook[*ModelEvent]
	onModelAfterCreateSuccess *hook.Hook[*ModelEvent]
	onModelAfterCreateError   *hook.Hook[*ModelErrorEvent]
	onModelUpdate             *hook.Hook[*ModelEvent]
	onModelUpdateWrite        *hook.Hook[*ModelEvent]
	onModelAfterUpdateSuccess *hook.Hook[*ModelEvent]
	onModelAfterUpdateError   *hook.Hook[*ModelErrorEvent]
	onModelDelete             *hook.Hook[*ModelEvent]
	onModelDeleteExecute      *hook.Hook[*ModelEvent]
	onModelAfterDeleteSuccess *hook.Hook[*ModelEvent]
	onModelAfterDeleteError   *hook.Hook[*ModelErrorEvent]

	// db record hooks
	onRecordEnrich             *hook.Hook[*RecordEnrichEvent]
	onRecordValidate           *hook.Hook[*RecordEvent]
	onRecordCreate             *hook.Hook[*RecordEvent]
	onRecordCreateExecute      *hook.Hook[*RecordEvent]
	onRecordAfterCreateSuccess *hook.Hook[*RecordEvent]
	onRecordAfterCreateError   *hook.Hook[*RecordErrorEvent]
	onRecordUpdate             *hook.Hook[*RecordEvent]
	onRecordUpdateExecute      *hook.Hook[*RecordEvent]
	onRecordAfterUpdateSuccess *hook.Hook[*RecordEvent]
	onRecordAfterUpdateError   *hook.Hook[*RecordErrorEvent]
	onRecordDelete             *hook.Hook[*RecordEvent]
	onRecordDeleteExecute      *hook.Hook[*RecordEvent]
	onRecordAfterDeleteSuccess *hook.Hook[*RecordEvent]
	onRecordAfterDeleteError   *hook.Hook[*RecordErrorEvent]

	// db collection hooks
	onCollectionValidate           *hook.Hook[*CollectionEvent]
	onCollectionCreate             *hook.Hook[*CollectionEvent]
	onCollectionCreateExecute      *hook.Hook[*CollectionEvent]
	onCollectionAfterCreateSuccess *hook.Hook[*CollectionEvent]
	onCollectionAfterCreateError   *hook.Hook[*CollectionErrorEvent]
	onCollectionUpdate             *hook.Hook[*CollectionEvent]
	onCollectionUpdateExecute      *hook.Hook[*CollectionEvent]
	onCollectionAfterUpdateSuccess *hook.Hook[*CollectionEvent]
	onCollectionAfterUpdateError   *hook.Hook[*CollectionErrorEvent]
	onCollectionDelete             *hook.Hook[*CollectionEvent]
	onCollectionDeleteExecute      *hook.Hook[*CollectionEvent]
	onCollectionAfterDeleteSuccess *hook.Hook[*CollectionEvent]
	onCollectionAfterDeleteError   *hook.Hook[*CollectionErrorEvent]

	// mailer event hooks
	onMailerSend                    *hook.Hook[*MailerEvent]
	onMailerRecordPasswordResetSend *hook.Hook[*MailerRecordEvent]
	onMailerRecordVerificationSend  *hook.Hook[*MailerRecordEvent]
	onMailerRecordEmailChangeSend   *hook.Hook[*MailerRecordEvent]
	onMailerRecordOTPSend           *hook.Hook[*MailerRecordEvent]
	onMailerRecordAuthAlertSend     *hook.Hook[*MailerRecordEvent]

	// realtime api event hooks
	onRealtimeConnectRequest   *hook.Hook[*RealtimeConnectRequestEvent]
	onRealtimeMessageSend      *hook.Hook[*RealtimeMessageEvent]
	onRealtimeSubscribeRequest *hook.Hook[*RealtimeSubscribeRequestEvent]

	// settings event hooks
	onSettingsListRequest   *hook.Hook[*SettingsListRequestEvent]
	onSettingsUpdateRequest *hook.Hook[*SettingsUpdateRequestEvent]
	onSettingsReload        *hook.Hook[*SettingsReloadEvent]

	// file api event hooks
	onFileDownloadRequest *hook.Hook[*FileDownloadRequestEvent]
	onFileTokenRequest    *hook.Hook[*FileTokenRequestEvent]

	// record auth API event hooks
	onRecordAuthRequest                 *hook.Hook[*RecordAuthRequestEvent]
	onRecordAuthWithPasswordRequest     *hook.Hook[*RecordAuthWithPasswordRequestEvent]
	onRecordAuthWithOAuth2Request       *hook.Hook[*RecordAuthWithOAuth2RequestEvent]
	onRecordAuthRefreshRequest          *hook.Hook[*RecordAuthRefreshRequestEvent]
	onRecordRequestPasswordResetRequest *hook.Hook[*RecordRequestPasswordResetRequestEvent]
	onRecordConfirmPasswordResetRequest *hook.Hook[*RecordConfirmPasswordResetRequestEvent]
	onRecordRequestVerificationRequest  *hook.Hook[*RecordRequestVerificationRequestEvent]
	onRecordConfirmVerificationRequest  *hook.Hook[*RecordConfirmVerificationRequestEvent]
	onRecordRequestEmailChangeRequest   *hook.Hook[*RecordRequestEmailChangeRequestEvent]
	onRecordConfirmEmailChangeRequest   *hook.Hook[*RecordConfirmEmailChangeRequestEvent]
	onRecordRequestOTPRequest           *hook.Hook[*RecordCreateOTPRequestEvent]
	onRecordAuthWithOTPRequest          *hook.Hook[*RecordAuthWithOTPRequestEvent]

	// record crud API event hooks
	onRecordsListRequest  *hook.Hook[*RecordsListRequestEvent]
	onRecordViewRequest   *hook.Hook[*RecordRequestEvent]
	onRecordCreateRequest *hook.Hook[*RecordRequestEvent]
	onRecordUpdateRequest *hook.Hook[*RecordRequestEvent]
	onRecordDeleteRequest *hook.Hook[*RecordRequestEvent]

	// collection API event hooks
	onCollectionsListRequest   *hook.Hook[*CollectionsListRequestEvent]
	onCollectionViewRequest    *hook.Hook[*CollectionRequestEvent]
	onCollectionCreateRequest  *hook.Hook[*CollectionRequestEvent]
	onCollectionUpdateRequest  *hook.Hook[*CollectionRequestEvent]
	onCollectionDeleteRequest  *hook.Hook[*CollectionRequestEvent]
	onCollectionsImportRequest *hook.Hook[*CollectionsImportRequestEvent]

	onBatchRequest *hook.Hook[*BatchRequestEvent]
}

// NewBaseApp creates and returns a new BaseApp instance
// configured with the provided arguments.
//
// To initialize the app, you need to call `app.Bootstrap()`.
func NewBaseApp(config BaseAppConfig) *BaseApp {
	app := &BaseApp{
		settings:            newDefaultSettings(),
		store:               store.New[string, any](nil),
		cron:                cron.New(),
		subscriptionsBroker: subscriptions.NewBroker(),
		config:              &config,
	}

	// apply config defaults
	if app.config.DBConnect == nil {
		app.config.DBConnect = DefaultDBConnect
	}
	if app.config.DataMaxOpenConns <= 0 {
		app.config.DataMaxOpenConns = DefaultDataMaxOpenConns
	}
	if app.config.DataMaxIdleConns <= 0 {
		app.config.DataMaxIdleConns = DefaultDataMaxIdleConns
	}
	if app.config.AuxMaxOpenConns <= 0 {
		app.config.AuxMaxOpenConns = DefaultAuxMaxOpenConns
	}
	if app.config.AuxMaxIdleConns <= 0 {
		app.config.AuxMaxIdleConns = DefaultAuxMaxIdleConns
	}
	if app.config.QueryTimeout <= 0 {
		app.config.QueryTimeout = DefaultQueryTimeout
	}

	app.initHooks()
	app.registerBaseHooks()

	return app
}

// initHooks initializes all app hook handlers.
func (app *BaseApp) initHooks() {
	// app event hooks
	app.onBootstrap = &hook.Hook[*BootstrapEvent]{}
	app.onServe = &hook.Hook[*ServeEvent]{}
	app.onTerminate = &hook.Hook[*TerminateEvent]{}
	app.onBackupCreate = &hook.Hook[*BackupEvent]{}
	app.onBackupRestore = &hook.Hook[*BackupEvent]{}

	// db model hooks
	app.onModelValidate = &hook.Hook[*ModelEvent]{}
	app.onModelCreate = &hook.Hook[*ModelEvent]{}
	app.onModelCreateExecute = &hook.Hook[*ModelEvent]{}
	app.onModelAfterCreateSuccess = &hook.Hook[*ModelEvent]{}
	app.onModelAfterCreateError = &hook.Hook[*ModelErrorEvent]{}
	app.onModelUpdate = &hook.Hook[*ModelEvent]{}
	app.onModelUpdateWrite = &hook.Hook[*ModelEvent]{}
	app.onModelAfterUpdateSuccess = &hook.Hook[*ModelEvent]{}
	app.onModelAfterUpdateError = &hook.Hook[*ModelErrorEvent]{}
	app.onModelDelete = &hook.Hook[*ModelEvent]{}
	app.onModelDeleteExecute = &hook.Hook[*ModelEvent]{}
	app.onModelAfterDeleteSuccess = &hook.Hook[*ModelEvent]{}
	app.onModelAfterDeleteError = &hook.Hook[*ModelErrorEvent]{}

	// db record hooks
	app.onRecordEnrich = &hook.Hook[*RecordEnrichEvent]{}
	app.onRecordValidate = &hook.Hook[*RecordEvent]{}
	app.onRecordCreate = &hook.Hook[*RecordEvent]{}
	app.onRecordCreateExecute = &hook.Hook[*RecordEvent]{}
	app.onRecordAfterCreateSuccess = &hook.Hook[*RecordEvent]{}
	app.onRecordAfterCreateError = &hook.Hook[*RecordErrorEvent]{}
	app.onRecordUpdate = &hook.Hook[*RecordEvent]{}
	app.onRecordUpdateExecute = &hook.Hook[*RecordEvent]{}
	app.onRecordAfterUpdateSuccess = &hook.Hook[*RecordEvent]{}
	app.onRecordAfterUpdateError = &hook.Hook[*RecordErrorEvent]{}
	app.onRecordDelete = &hook.Hook[*RecordEvent]{}
	app.onRecordDeleteExecute = &hook.Hook[*RecordEvent]{}
	app.onRecordAfterDeleteSuccess = &hook.Hook[*RecordEvent]{}
	app.onRecordAfterDeleteError = &hook.Hook[*RecordErrorEvent]{}

	// db collection hooks
	app.onCollectionValidate = &hook.Hook[*CollectionEvent]{}
	app.onCollectionCreate = &hook.Hook[*CollectionEvent]{}
	app.onCollectionCreateExecute = &hook.Hook[*CollectionEvent]{}
	app.onCollectionAfterCreateSuccess = &hook.Hook[*CollectionEvent]{}
	app.onCollectionAfterCreateError = &hook.Hook[*CollectionErrorEvent]{}
	app.onCollectionUpdate = &hook.Hook[*CollectionEvent]{}
	app.onCollectionUpdateExecute = &hook.Hook[*CollectionEvent]{}
	app.onCollectionAfterUpdateSuccess = &hook.Hook[*CollectionEvent]{}
	app.onCollectionAfterUpdateError = &hook.Hook[*CollectionErrorEvent]{}
	app.onCollectionDelete = &hook.Hook[*CollectionEvent]{}
	app.onCollectionAfterDeleteSuccess = &hook.Hook[*CollectionEvent]{}
	app.onCollectionDeleteExecute = &hook.Hook[*CollectionEvent]{}
	app.onCollectionAfterDeleteError = &hook.Hook[*CollectionErrorEvent]{}

	// mailer event hooks
	app.onMailerSend = &hook.Hook[*MailerEvent]{}
	app.onMailerRecordPasswordResetSend = &hook.Hook[*MailerRecordEvent]{}
	app.onMailerRecordVerificationSend = &hook.Hook[*MailerRecordEvent]{}
	app.onMailerRecordEmailChangeSend = &hook.Hook[*MailerRecordEvent]{}
	app.onMailerRecordOTPSend = &hook.Hook[*MailerRecordEvent]{}
	app.onMailerRecordAuthAlertSend = &hook.Hook[*MailerRecordEvent]{}

	// realtime API event hooks
	app.onRealtimeConnectRequest = &hook.Hook[*RealtimeConnectRequestEvent]{}
	app.onRealtimeMessageSend = &hook.Hook[*RealtimeMessageEvent]{}
	app.onRealtimeSubscribeRequest = &hook.Hook[*RealtimeSubscribeRequestEvent]{}

	// settings event hooks
	app.onSettingsListRequest = &hook.Hook[*SettingsListRequestEvent]{}
	app.onSettingsUpdateRequest = &hook.Hook[*SettingsUpdateRequestEvent]{}
	app.onSettingsReload = &hook.Hook[*SettingsReloadEvent]{}

	// file API event hooks
	app.onFileDownloadRequest = &hook.Hook[*FileDownloadRequestEvent]{}
	app.onFileTokenRequest = &hook.Hook[*FileTokenRequestEvent]{}

	// record auth API event hooks
	app.onRecordAuthRequest = &hook.Hook[*RecordAuthRequestEvent]{}
	app.onRecordAuthWithPasswordRequest = &hook.Hook[*RecordAuthWithPasswordRequestEvent]{}
	app.onRecordAuthWithOAuth2Request = &hook.Hook[*RecordAuthWithOAuth2RequestEvent]{}
	app.onRecordAuthRefreshRequest = &hook.Hook[*RecordAuthRefreshRequestEvent]{}
	app.onRecordRequestPasswordResetRequest = &hook.Hook[*RecordRequestPasswordResetRequestEvent]{}
	app.onRecordConfirmPasswordResetRequest = &hook.Hook[*RecordConfirmPasswordResetRequestEvent]{}
	app.onRecordRequestVerificationRequest = &hook.Hook[*RecordRequestVerificationRequestEvent]{}
	app.onRecordConfirmVerificationRequest = &hook.Hook[*RecordConfirmVerificationRequestEvent]{}
	app.onRecordRequestEmailChangeRequest = &hook.Hook[*RecordRequestEmailChangeRequestEvent]{}
	app.onRecordConfirmEmailChangeRequest = &hook.Hook[*RecordConfirmEmailChangeRequestEvent]{}
	app.onRecordRequestOTPRequest = &hook.Hook[*RecordCreateOTPRequestEvent]{}
	app.onRecordAuthWithOTPRequest = &hook.Hook[*RecordAuthWithOTPRequestEvent]{}

	// record crud API event hooks
	app.onRecordsListRequest = &hook.Hook[*RecordsListRequestEvent]{}
	app.onRecordViewRequest = &hook.Hook[*RecordRequestEvent]{}
	app.onRecordCreateRequest = &hook.Hook[*RecordRequestEvent]{}
	app.onRecordUpdateRequest = &hook.Hook[*RecordRequestEvent]{}
	app.onRecordDeleteRequest = &hook.Hook[*RecordRequestEvent]{}

	// collection API event hooks
	app.onCollectionsListRequest = &hook.Hook[*CollectionsListRequestEvent]{}
	app.onCollectionViewRequest = &hook.Hook[*CollectionRequestEvent]{}
	app.onCollectionCreateRequest = &hook.Hook[*CollectionRequestEvent]{}
	app.onCollectionUpdateRequest = &hook.Hook[*CollectionRequestEvent]{}
	app.onCollectionDeleteRequest = &hook.Hook[*CollectionRequestEvent]{}
	app.onCollectionsImportRequest = &hook.Hook[*CollectionsImportRequestEvent]{}

	app.onBatchRequest = &hook.Hook[*BatchRequestEvent]{}
}

// UnsafeWithoutHooks returns a shallow copy of the current app WITHOUT any registered hooks.
//
// NB! Note that using the returned app instance may cause data integrity errors
// since the Record validations and data normalizations (including files uploads)
// rely on the app hooks to work.
func (app *BaseApp) UnsafeWithoutHooks() App {
	clone := *app

	// reset all hook handlers
	clone.initHooks()

	return &clone
}

// Logger returns the default app logger.
//
// If the application is not bootstrapped yet, fallbacks to slog.Default().
func (app *BaseApp) Logger() *slog.Logger {
	if app.logger == nil {
		return slog.Default()
	}

	return app.logger
}

// TxInfo returns the transaction associated with the current app instance (if any).
//
// Could be used if you want to execute indirectly a function after
// the related app transaction completes using `app.TxInfo().OnAfterFunc(callback)`.
func (app *BaseApp) TxInfo() *TxAppInfo {
	return app.txInfo
}

// IsTransactional checks if the current app instance is part of a transaction.
func (app *BaseApp) IsTransactional() bool {
	return app.TxInfo() != nil
}

// IsBootstrapped checks if the application was initialized
// (aka. whether Bootstrap() was called).
func (app *BaseApp) IsBootstrapped() bool {
	return app.concurrentDB != nil && app.auxConcurrentDB != nil
}

// Bootstrap initializes the application
// (aka. create data dir, open db connections, load settings, etc.).
//
// It will call ResetBootstrapState() if the application was already bootstrapped.
func (app *BaseApp) Bootstrap() error {
	event := &BootstrapEvent{}
	event.App = app

	err := app.OnBootstrap().Trigger(event, func(e *BootstrapEvent) error {
		// clear resources of previous core state (if any)
		if err := app.ResetBootstrapState(); err != nil {
			return err
		}

		// ensure that data dir exist
		if err := os.MkdirAll(app.DataDir(), os.ModePerm); err != nil {
			return err
		}

		if err := app.initDataDB(); err != nil {
			return err
		}

		if err := app.initAuxDB(); err != nil {
			return err
		}

		if err := app.initLogger(); err != nil {
			return err
		}

		if err := app.RunSystemMigrations(); err != nil {
			return err
		}

		if err := app.ReloadCachedCollections(); err != nil {
			return err
		}

		if err := app.ReloadSettings(); err != nil {
			return err
		}

		// try to cleanup the pb_data temp directory (if any)
		_ = os.RemoveAll(filepath.Join(app.DataDir(), LocalTempDirName))

		return nil
	})

	// add a more user friendly message in case users forgot to call
	// e.Next() as part of their bootstrap hook
	if err == nil && !app.IsBootstrapped() {
		app.Logger().Warn("OnBootstrap hook didn't fail but the app is still not bootstrapped - maybe missing e.Next()?")
	}

	return err
}

type closer interface {
	Close() error
}

// ResetBootstrapState releases the initialized core app resources
// (closing db connections, stopping cron ticker, etc.).
func (app *BaseApp) ResetBootstrapState() error {
	app.Cron().Stop()

	var errs []error

	dbs := []*dbx.Builder{
		&app.concurrentDB,
		&app.nonconcurrentDB,
		&app.auxConcurrentDB,
		&app.auxNonconcurrentDB,
	}

	for _, db := range dbs {
		if db == nil {
			continue
		}
		if v, ok := (*db).(closer); ok {
			if err := v.Close(); err != nil {
				errs = append(errs, err)
			}
		}
		*db = nil
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// DB returns the default app data.db builder instance.
//
// To minimize SQLITE_BUSY errors, it automatically routes the
// SELECT queries to the underlying concurrent db pool and everything
// else to the nonconcurrent one.
//
// For more finer control over the used connections pools you can
// call directly ConcurrentDB() or NonconcurrentDB().
func (app *BaseApp) DB() dbx.Builder {
	// transactional or both are nil
	if app.concurrentDB == app.nonconcurrentDB {
		return app.concurrentDB
	}

	return &dualDBBuilder{
		concurrentDB:    app.concurrentDB,
		nonconcurrentDB: app.nonconcurrentDB,
	}
}

// ConcurrentDB returns the concurrent app data.db builder instance.
//
// This method is used mainly internally for executing db read
// operations in a concurrent/non-blocking manner.
//
// Most users should use simply DB() as it will automatically
// route the query execution to ConcurrentDB() or NonconcurrentDB().
//
// In a transaction the ConcurrentDB() and NonconcurrentDB() refer to the same *dbx.TX instance.
func (app *BaseApp) ConcurrentDB() dbx.Builder {
	return app.concurrentDB
}

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
func (app *BaseApp) NonconcurrentDB() dbx.Builder {
	return app.nonconcurrentDB
}

// AuxDB returns the app auxiliary.db builder instance.
//
// To minimize SQLITE_BUSY errors, it automatically routes the
// SELECT queries to the underlying concurrent db pool and everything
// else to the nonconcurrent one.
//
// For more finer control over the used connections pools you can
// call directly AuxConcurrentDB() or AuxNonconcurrentDB().
func (app *BaseApp) AuxDB() dbx.Builder {
	// transactional or both are nil
	if app.auxConcurrentDB == app.auxNonconcurrentDB {
		return app.auxConcurrentDB
	}

	return &dualDBBuilder{
		concurrentDB:    app.auxConcurrentDB,
		nonconcurrentDB: app.auxNonconcurrentDB,
	}
}

// AuxConcurrentDB returns the concurrent app auxiliary.db builder instance.
//
// This method is used mainly internally for executing db read
// operations in a concurrent/non-blocking manner.
//
// Most users should use simply AuxDB() as it will automatically
// route the query execution to AuxConcurrentDB() or AuxNonconcurrentDB().
//
// In a transaction the AuxConcurrentDB() and AuxNonconcurrentDB() refer to the same *dbx.TX instance.
func (app *BaseApp) AuxConcurrentDB() dbx.Builder {
	return app.auxConcurrentDB
}

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
func (app *BaseApp) AuxNonconcurrentDB() dbx.Builder {
	return app.auxNonconcurrentDB
}

// DataDir returns the app data directory path.
func (app *BaseApp) DataDir() string {
	return app.config.DataDir
}

// EncryptionEnv returns the name of the app secret env key
// (currently used primarily for optional settings encryption but this may change in the future).
func (app *BaseApp) EncryptionEnv() string {
	return app.config.EncryptionEnv
}

// IsDev returns whether the app is in dev mode.
//
// When enabled logs, executed sql statements, etc. are printed to the stderr.
func (app *BaseApp) IsDev() bool {
	return app.config.IsDev
}

// Settings returns the loaded app settings.
func (app *BaseApp) Settings() *Settings {
	return app.settings
}

// Store returns the app runtime store.
func (app *BaseApp) Store() *store.Store[string, any] {
	return app.store
}

// Cron returns the app cron instance.
func (app *BaseApp) Cron() *cron.Cron {
	return app.cron
}

// SubscriptionsBroker returns the app realtime subscriptions broker instance.
func (app *BaseApp) SubscriptionsBroker() *subscriptions.Broker {
	return app.subscriptionsBroker
}

// NewMailClient creates and returns a new SMTP or Sendmail client
// based on the current app settings.
func (app *BaseApp) NewMailClient() mailer.Mailer {
	var client mailer.Mailer

	// init mailer client
	if app.Settings().SMTP.Enabled {
		client = &mailer.SMTPClient{
			Host:       app.Settings().SMTP.Host,
			Port:       app.Settings().SMTP.Port,
			Username:   app.Settings().SMTP.Username,
			Password:   app.Settings().SMTP.Password,
			TLS:        app.Settings().SMTP.TLS,
			AuthMethod: app.Settings().SMTP.AuthMethod,
			LocalName:  app.Settings().SMTP.LocalName,
		}
	} else {
		client = &mailer.Sendmail{}
	}

	// register the app level hook
	if h, ok := client.(mailer.SendInterceptor); ok {
		h.OnSend().Bind(&hook.Handler[*mailer.SendEvent]{
			Id: "__pbMailerOnSend__",
			Func: func(e *mailer.SendEvent) error {
				appEvent := new(MailerEvent)
				appEvent.App = app
				appEvent.Mailer = client
				appEvent.Message = e.Message

				return app.OnMailerSend().Trigger(appEvent, func(ae *MailerEvent) error {
					e.Message = ae.Message

					// print the mail in the console to assist with the debugging
					if app.IsDev() {
						logDate := new(strings.Builder)
						log.New(logDate, "", log.LstdFlags).Print()

						mailLog := new(strings.Builder)
						mailLog.WriteString(strings.TrimSpace(logDate.String()))
						mailLog.WriteString(" Mail sent\n")
						fmt.Fprintf(mailLog, "├─ From: %v\n", ae.Message.From)
						fmt.Fprintf(mailLog, "├─ To: %v\n", ae.Message.To)
						fmt.Fprintf(mailLog, "├─ Cc: %v\n", ae.Message.Cc)
						fmt.Fprintf(mailLog, "├─ Bcc: %v\n", ae.Message.Bcc)
						fmt.Fprintf(mailLog, "├─ Subject: %v\n", ae.Message.Subject)

						if len(ae.Message.Attachments) > 0 {
							attachmentKeys := make([]string, 0, len(ae.Message.Attachments))
							for k := range ae.Message.Attachments {
								attachmentKeys = append(attachmentKeys, k)
							}
							fmt.Fprintf(mailLog, "├─ Attachments: %v\n", attachmentKeys)
						}

						if len(ae.Message.InlineAttachments) > 0 {
							attachmentKeys := make([]string, 0, len(ae.Message.InlineAttachments))
							for k := range ae.Message.InlineAttachments {
								attachmentKeys = append(attachmentKeys, k)
							}
							fmt.Fprintf(mailLog, "├─ InlineAttachments: %v\n", attachmentKeys)
						}

						const indentation = "        "
						if ae.Message.Text != "" {
							textParts := strings.Split(strings.TrimSpace(ae.Message.Text), "\n")
							textIndented := indentation + strings.Join(textParts, "\n"+indentation)
							fmt.Fprintf(mailLog, "└─ Text:\n%s", textIndented)
						} else {
							htmlParts := strings.Split(strings.TrimSpace(ae.Message.HTML), "\n")
							htmlIndented := indentation + strings.Join(htmlParts, "\n"+indentation)
							fmt.Fprintf(mailLog, "└─ HTML:\n%s", htmlIndented)
						}

						color.HiBlack("%s", mailLog.String())
					}

					// send the email with the new mailer in case it was replaced
					if client != ae.Mailer {
						return ae.Mailer.Send(e.Message)
					}

					return e.Next()
				})
			},
		})
	}

	return client
}

// NewFilesystem creates a new local or S3 filesystem instance
// for managing regular app files (ex. record uploads)
// based on the current app settings.
//
// NB! Make sure to call Close() on the returned result
// after you are done working with it.
func (app *BaseApp) NewFilesystem() (*filesystem.System, error) {
	if app.settings != nil && app.settings.S3.Enabled {
		return filesystem.NewS3(
			app.settings.S3.Bucket,
			app.settings.S3.Region,
			app.settings.S3.Endpoint,
			app.settings.S3.AccessKey,
			app.settings.S3.Secret,
			app.settings.S3.ForcePathStyle,
		)
	}

	// fallback to local filesystem
	return filesystem.NewLocal(filepath.Join(app.DataDir(), LocalStorageDirName))
}

// NewBackupsFilesystem creates a new local or S3 filesystem instance
// for managing app backups based on the current app settings.
//
// NB! Make sure to call Close() on the returned result
// after you are done working with it.
func (app *BaseApp) NewBackupsFilesystem() (*filesystem.System, error) {
	if app.settings != nil && app.settings.Backups.S3.Enabled {
		return filesystem.NewS3(
			app.settings.Backups.S3.Bucket,
			app.settings.Backups.S3.Region,
			app.settings.Backups.S3.Endpoint,
			app.settings.Backups.S3.AccessKey,
			app.settings.Backups.S3.Secret,
			app.settings.Backups.S3.ForcePathStyle,
		)
	}

	// fallback to local filesystem
	return filesystem.NewLocal(filepath.Join(app.DataDir(), LocalBackupsDirName))
}

// Restart restarts (aka. replaces) the current running application process.
//
// NB! It relies on execve which is supported only on UNIX based systems.
func (app *BaseApp) Restart() error {
	if runtime.GOOS == "windows" {
		return errors.New("restart is not supported on windows")
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	event := &TerminateEvent{}
	event.App = app
	event.IsRestart = true

	return app.OnTerminate().Trigger(event, func(e *TerminateEvent) error {
		_ = e.App.ResetBootstrapState()

		// attempt to restart the bootstrap process in case execve returns an error for some reason
		defer func() {
			if err := e.App.Bootstrap(); err != nil {
				app.Logger().Error("Failed to rebootstrap the application after failed app.Restart()", "error", err)
			}
		}()

		return execve(execPath, os.Args, os.Environ())
	})
}

// RunSystemMigrations applies all new migrations registered in the [core.SystemMigrations] list.
func (app *BaseApp) RunSystemMigrations() error {
	_, err := NewMigrationsRunner(app, SystemMigrations).Up()
	return err
}

// RunAppMigrations applies all new migrations registered in the [core.AppMigrations] list.
func (app *BaseApp) RunAppMigrations() error {
	_, err := NewMigrationsRunner(app, AppMigrations).Up()
	return err
}

// RunAllMigrations applies all system and app migrations
// (aka. from both [core.SystemMigrations] and [core.AppMigrations]).
func (app *BaseApp) RunAllMigrations() error {
	list := MigrationsList{}
	list.Copy(SystemMigrations)
	list.Copy(AppMigrations)
	_, err := NewMigrationsRunner(app, list).Up()
	return err
}

// -------------------------------------------------------------------
// App event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnBootstrap() *hook.Hook[*BootstrapEvent] {
	return app.onBootstrap
}

func (app *BaseApp) OnServe() *hook.Hook[*ServeEvent] {
	return app.onServe
}

func (app *BaseApp) OnTerminate() *hook.Hook[*TerminateEvent] {
	return app.onTerminate
}

func (app *BaseApp) OnBackupCreate() *hook.Hook[*BackupEvent] {
	return app.onBackupCreate
}

func (app *BaseApp) OnBackupRestore() *hook.Hook[*BackupEvent] {
	return app.onBackupRestore
}

// ---------------------------------------------------------------

func (app *BaseApp) OnModelCreate(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelCreate, tags...)
}

func (app *BaseApp) OnModelCreateExecute(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelCreateExecute, tags...)
}

func (app *BaseApp) OnModelAfterCreateSuccess(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelAfterCreateSuccess, tags...)
}

func (app *BaseApp) OnModelAfterCreateError(tags ...string) *hook.TaggedHook[*ModelErrorEvent] {
	return hook.NewTaggedHook(app.onModelAfterCreateError, tags...)
}

func (app *BaseApp) OnModelUpdate(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelUpdate, tags...)
}

func (app *BaseApp) OnModelUpdateExecute(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelUpdateWrite, tags...)
}

func (app *BaseApp) OnModelAfterUpdateSuccess(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelAfterUpdateSuccess, tags...)
}

func (app *BaseApp) OnModelAfterUpdateError(tags ...string) *hook.TaggedHook[*ModelErrorEvent] {
	return hook.NewTaggedHook(app.onModelAfterUpdateError, tags...)
}

func (app *BaseApp) OnModelValidate(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelValidate, tags...)
}

func (app *BaseApp) OnModelDelete(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelDelete, tags...)
}

func (app *BaseApp) OnModelDeleteExecute(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelDeleteExecute, tags...)
}

func (app *BaseApp) OnModelAfterDeleteSuccess(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelAfterDeleteSuccess, tags...)
}

func (app *BaseApp) OnModelAfterDeleteError(tags ...string) *hook.TaggedHook[*ModelErrorEvent] {
	return hook.NewTaggedHook(app.onModelAfterDeleteError, tags...)
}

func (app *BaseApp) OnRecordEnrich(tags ...string) *hook.TaggedHook[*RecordEnrichEvent] {
	return hook.NewTaggedHook(app.onRecordEnrich, tags...)
}

func (app *BaseApp) OnRecordValidate(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordValidate, tags...)
}

func (app *BaseApp) OnRecordCreate(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordCreate, tags...)
}

func (app *BaseApp) OnRecordCreateExecute(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordCreateExecute, tags...)
}

func (app *BaseApp) OnRecordAfterCreateSuccess(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordAfterCreateSuccess, tags...)
}

func (app *BaseApp) OnRecordAfterCreateError(tags ...string) *hook.TaggedHook[*RecordErrorEvent] {
	return hook.NewTaggedHook(app.onRecordAfterCreateError, tags...)
}

func (app *BaseApp) OnRecordUpdate(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordUpdate, tags...)
}

func (app *BaseApp) OnRecordUpdateExecute(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordUpdateExecute, tags...)
}

func (app *BaseApp) OnRecordAfterUpdateSuccess(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordAfterUpdateSuccess, tags...)
}

func (app *BaseApp) OnRecordAfterUpdateError(tags ...string) *hook.TaggedHook[*RecordErrorEvent] {
	return hook.NewTaggedHook(app.onRecordAfterUpdateError, tags...)
}

func (app *BaseApp) OnRecordDelete(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordDelete, tags...)
}

func (app *BaseApp) OnRecordDeleteExecute(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordDeleteExecute, tags...)
}

func (app *BaseApp) OnRecordAfterDeleteSuccess(tags ...string) *hook.TaggedHook[*RecordEvent] {
	return hook.NewTaggedHook(app.onRecordAfterDeleteSuccess, tags...)
}

func (app *BaseApp) OnRecordAfterDeleteError(tags ...string) *hook.TaggedHook[*RecordErrorEvent] {
	return hook.NewTaggedHook(app.onRecordAfterDeleteError, tags...)
}

func (app *BaseApp) OnCollectionValidate(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionValidate, tags...)
}

func (app *BaseApp) OnCollectionCreate(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionCreate, tags...)
}

func (app *BaseApp) OnCollectionCreateExecute(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionCreateExecute, tags...)
}

func (app *BaseApp) OnCollectionAfterCreateSuccess(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionAfterCreateSuccess, tags...)
}

func (app *BaseApp) OnCollectionAfterCreateError(tags ...string) *hook.TaggedHook[*CollectionErrorEvent] {
	return hook.NewTaggedHook(app.onCollectionAfterCreateError, tags...)
}

func (app *BaseApp) OnCollectionUpdate(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionUpdate, tags...)
}

func (app *BaseApp) OnCollectionUpdateExecute(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionUpdateExecute, tags...)
}

func (app *BaseApp) OnCollectionAfterUpdateSuccess(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionAfterUpdateSuccess, tags...)
}

func (app *BaseApp) OnCollectionAfterUpdateError(tags ...string) *hook.TaggedHook[*CollectionErrorEvent] {
	return hook.NewTaggedHook(app.onCollectionAfterUpdateError, tags...)
}

func (app *BaseApp) OnCollectionDelete(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionDelete, tags...)
}

func (app *BaseApp) OnCollectionDeleteExecute(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionDeleteExecute, tags...)
}

func (app *BaseApp) OnCollectionAfterDeleteSuccess(tags ...string) *hook.TaggedHook[*CollectionEvent] {
	return hook.NewTaggedHook(app.onCollectionAfterDeleteSuccess, tags...)
}

func (app *BaseApp) OnCollectionAfterDeleteError(tags ...string) *hook.TaggedHook[*CollectionErrorEvent] {
	return hook.NewTaggedHook(app.onCollectionAfterDeleteError, tags...)
}

// -------------------------------------------------------------------
// Mailer event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnMailerSend() *hook.Hook[*MailerEvent] {
	return app.onMailerSend
}

func (app *BaseApp) OnMailerRecordPasswordResetSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerRecordPasswordResetSend, tags...)
}

func (app *BaseApp) OnMailerRecordVerificationSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerRecordVerificationSend, tags...)
}

func (app *BaseApp) OnMailerRecordEmailChangeSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerRecordEmailChangeSend, tags...)
}

func (app *BaseApp) OnMailerRecordOTPSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerRecordOTPSend, tags...)
}

func (app *BaseApp) OnMailerRecordAuthAlertSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerRecordAuthAlertSend, tags...)
}

// -------------------------------------------------------------------
// Realtime API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRealtimeConnectRequest() *hook.Hook[*RealtimeConnectRequestEvent] {
	return app.onRealtimeConnectRequest
}

func (app *BaseApp) OnRealtimeMessageSend() *hook.Hook[*RealtimeMessageEvent] {
	return app.onRealtimeMessageSend
}

func (app *BaseApp) OnRealtimeSubscribeRequest() *hook.Hook[*RealtimeSubscribeRequestEvent] {
	return app.onRealtimeSubscribeRequest
}

// -------------------------------------------------------------------
// Settings API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnSettingsListRequest() *hook.Hook[*SettingsListRequestEvent] {
	return app.onSettingsListRequest
}

func (app *BaseApp) OnSettingsUpdateRequest() *hook.Hook[*SettingsUpdateRequestEvent] {
	return app.onSettingsUpdateRequest
}

func (app *BaseApp) OnSettingsReload() *hook.Hook[*SettingsReloadEvent] {
	return app.onSettingsReload
}

// -------------------------------------------------------------------
// File API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnFileDownloadRequest(tags ...string) *hook.TaggedHook[*FileDownloadRequestEvent] {
	return hook.NewTaggedHook(app.onFileDownloadRequest, tags...)
}

func (app *BaseApp) OnFileTokenRequest(tags ...string) *hook.TaggedHook[*FileTokenRequestEvent] {
	return hook.NewTaggedHook(app.onFileTokenRequest, tags...)
}

// -------------------------------------------------------------------
// Record auth API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRecordAuthRequest(tags ...string) *hook.TaggedHook[*RecordAuthRequestEvent] {
	return hook.NewTaggedHook(app.onRecordAuthRequest, tags...)
}

func (app *BaseApp) OnRecordAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*RecordAuthWithPasswordRequestEvent] {
	return hook.NewTaggedHook(app.onRecordAuthWithPasswordRequest, tags...)
}

func (app *BaseApp) OnRecordAuthWithOAuth2Request(tags ...string) *hook.TaggedHook[*RecordAuthWithOAuth2RequestEvent] {
	return hook.NewTaggedHook(app.onRecordAuthWithOAuth2Request, tags...)
}

func (app *BaseApp) OnRecordAuthRefreshRequest(tags ...string) *hook.TaggedHook[*RecordAuthRefreshRequestEvent] {
	return hook.NewTaggedHook(app.onRecordAuthRefreshRequest, tags...)
}

func (app *BaseApp) OnRecordRequestPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordRequestPasswordResetRequestEvent] {
	return hook.NewTaggedHook(app.onRecordRequestPasswordResetRequest, tags...)
}

func (app *BaseApp) OnRecordConfirmPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordConfirmPasswordResetRequestEvent] {
	return hook.NewTaggedHook(app.onRecordConfirmPasswordResetRequest, tags...)
}

func (app *BaseApp) OnRecordRequestVerificationRequest(tags ...string) *hook.TaggedHook[*RecordRequestVerificationRequestEvent] {
	return hook.NewTaggedHook(app.onRecordRequestVerificationRequest, tags...)
}

func (app *BaseApp) OnRecordConfirmVerificationRequest(tags ...string) *hook.TaggedHook[*RecordConfirmVerificationRequestEvent] {
	return hook.NewTaggedHook(app.onRecordConfirmVerificationRequest, tags...)
}

func (app *BaseApp) OnRecordRequestEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordRequestEmailChangeRequestEvent] {
	return hook.NewTaggedHook(app.onRecordRequestEmailChangeRequest, tags...)
}

func (app *BaseApp) OnRecordConfirmEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordConfirmEmailChangeRequestEvent] {
	return hook.NewTaggedHook(app.onRecordConfirmEmailChangeRequest, tags...)
}

func (app *BaseApp) OnRecordRequestOTPRequest(tags ...string) *hook.TaggedHook[*RecordCreateOTPRequestEvent] {
	return hook.NewTaggedHook(app.onRecordRequestOTPRequest, tags...)
}

func (app *BaseApp) OnRecordAuthWithOTPRequest(tags ...string) *hook.TaggedHook[*RecordAuthWithOTPRequestEvent] {
	return hook.NewTaggedHook(app.onRecordAuthWithOTPRequest, tags...)
}

// -------------------------------------------------------------------
// Record CRUD API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRecordsListRequest(tags ...string) *hook.TaggedHook[*RecordsListRequestEvent] {
	return hook.NewTaggedHook(app.onRecordsListRequest, tags...)
}

func (app *BaseApp) OnRecordViewRequest(tags ...string) *hook.TaggedHook[*RecordRequestEvent] {
	return hook.NewTaggedHook(app.onRecordViewRequest, tags...)
}

func (app *BaseApp) OnRecordCreateRequest(tags ...string) *hook.TaggedHook[*RecordRequestEvent] {
	return hook.NewTaggedHook(app.onRecordCreateRequest, tags...)
}

func (app *BaseApp) OnRecordUpdateRequest(tags ...string) *hook.TaggedHook[*RecordRequestEvent] {
	return hook.NewTaggedHook(app.onRecordUpdateRequest, tags...)
}

func (app *BaseApp) OnRecordDeleteRequest(tags ...string) *hook.TaggedHook[*RecordRequestEvent] {
	return hook.NewTaggedHook(app.onRecordDeleteRequest, tags...)
}

// -------------------------------------------------------------------
// Collection API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnCollectionsListRequest() *hook.Hook[*CollectionsListRequestEvent] {
	return app.onCollectionsListRequest
}

func (app *BaseApp) OnCollectionViewRequest() *hook.Hook[*CollectionRequestEvent] {
	return app.onCollectionViewRequest
}

func (app *BaseApp) OnCollectionCreateRequest() *hook.Hook[*CollectionRequestEvent] {
	return app.onCollectionCreateRequest
}

func (app *BaseApp) OnCollectionUpdateRequest() *hook.Hook[*CollectionRequestEvent] {
	return app.onCollectionUpdateRequest
}

func (app *BaseApp) OnCollectionDeleteRequest() *hook.Hook[*CollectionRequestEvent] {
	return app.onCollectionDeleteRequest
}

func (app *BaseApp) OnCollectionsImportRequest() *hook.Hook[*CollectionsImportRequestEvent] {
	return app.onCollectionsImportRequest
}

func (app *BaseApp) OnBatchRequest() *hook.Hook[*BatchRequestEvent] {
	return app.onBatchRequest
}

// -------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------

func (app *BaseApp) initDataDB() error {
	dbPath := filepath.Join(app.DataDir(), "data.db")

	concurrentDB, err := app.config.DBConnect(dbPath)
	if err != nil {
		return err
	}
	concurrentDB.DB().SetMaxOpenConns(app.config.DataMaxOpenConns)
	concurrentDB.DB().SetMaxIdleConns(app.config.DataMaxIdleConns)
	concurrentDB.DB().SetConnMaxIdleTime(3 * time.Minute)

	nonconcurrentDB, err := app.config.DBConnect(dbPath)
	if err != nil {
		return err
	}
	nonconcurrentDB.DB().SetMaxOpenConns(1)
	nonconcurrentDB.DB().SetMaxIdleConns(1)
	nonconcurrentDB.DB().SetConnMaxIdleTime(3 * time.Minute)

	if app.IsDev() {
		nonconcurrentDB.QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
			color.HiBlack("[%.2fms] %v\n", float64(t.Milliseconds()), normalizeSQLLog(sql))
		}
		nonconcurrentDB.ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
			color.HiBlack("[%.2fms] %v\n", float64(t.Milliseconds()), normalizeSQLLog(sql))
		}
		concurrentDB.QueryLogFunc = nonconcurrentDB.QueryLogFunc
		concurrentDB.ExecLogFunc = nonconcurrentDB.ExecLogFunc
	}

	app.concurrentDB = concurrentDB
	app.nonconcurrentDB = nonconcurrentDB

	return nil
}

var sqlLogReplacements = []struct {
	pattern     *regexp.Regexp
	replacement string
}{
	{regexp.MustCompile(`\[\[([^\[\]\{\}\.]+)\.([^\[\]\{\}\.]+)\]\]`), "`$1`.`$2`"},
	{regexp.MustCompile(`\{\{([^\[\]\{\}\.]+)\.([^\[\]\{\}\.]+)\}\}`), "`$1`.`$2`"},
	{regexp.MustCompile(`([^'"])?\{\{`), "$1`"},
	{regexp.MustCompile(`\}\}([^'"])?`), "`$1"},
	{regexp.MustCompile(`([^'"])?\[\[`), "$1`"},
	{regexp.MustCompile(`\]\]([^'"])?`), "`$1"},
	{regexp.MustCompile(`<nil>`), "NULL"},
}

// normalizeSQLLog replaces common query builder charactes with their plain SQL version for easier debugging.
// The query is still not suitable for execution and should be used only for log and debug purposes
// (the normalization is done here to avoid breaking changes in dbx).
func normalizeSQLLog(sql string) string {
	for _, item := range sqlLogReplacements {
		sql = item.pattern.ReplaceAllString(sql, item.replacement)
	}

	return sql
}

func (app *BaseApp) initAuxDB() error {
	// note: renamed to "auxiliary" because "aux" is a reserved Windows filename
	// (see https://github.com/pocketbase/pocketbase/issues/5607)
	dbPath := filepath.Join(app.DataDir(), "auxiliary.db")

	concurrentDB, err := app.config.DBConnect(dbPath)
	if err != nil {
		return err
	}
	concurrentDB.DB().SetMaxOpenConns(app.config.AuxMaxOpenConns)
	concurrentDB.DB().SetMaxIdleConns(app.config.AuxMaxIdleConns)
	concurrentDB.DB().SetConnMaxIdleTime(3 * time.Minute)

	nonconcurrentDB, err := app.config.DBConnect(dbPath)
	if err != nil {
		return err
	}
	nonconcurrentDB.DB().SetMaxOpenConns(1)
	nonconcurrentDB.DB().SetMaxIdleConns(1)
	nonconcurrentDB.DB().SetConnMaxIdleTime(3 * time.Minute)

	app.auxConcurrentDB = concurrentDB
	app.auxNonconcurrentDB = nonconcurrentDB

	return nil
}

// @todo remove after refactoring the FilesManager interface
func supportFiles(m Model) bool {
	var collection *Collection
	switch v := m.(type) {
	case *Collection:
		collection = v
	case *Record:
		collection = v.Collection()
	case RecordProxy:
		if v.ProxyRecord() != nil {
			collection = v.ProxyRecord().Collection()
		}
	}

	if collection == nil {
		return true
	}

	for _, f := range collection.Fields {
		if f.Type() == FieldTypeFile {
			return true
		}
	}

	return false
}

func (app *BaseApp) registerBaseHooks() {
	deletePrefix := func(prefix string) error {
		fs, err := app.NewFilesystem()
		if err != nil {
			return err
		}
		defer fs.Close()

		failed := fs.DeletePrefix(prefix)
		if len(failed) > 0 {
			return errors.New("failed to delete the files at " + prefix)
		}

		return nil
	}

	maxFilesDeleteWorkers := cast.ToInt64(os.Getenv("PB_FILES_DELETE_MAX_WORKERS"))
	if maxFilesDeleteWorkers <= 0 {
		maxFilesDeleteWorkers = 2000 // the value is arbitrary chosen and may change in the future
	}

	deleteSem := semaphore.NewWeighted(maxFilesDeleteWorkers)

	// try to delete the storage files from deleted Collection, Records, etc. model
	app.OnModelAfterDeleteSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: "__pbFilesManagerDelete__",
		Func: func(e *ModelEvent) error {
			if m, ok := e.Model.(FilesManager); ok && m.BaseFilesPath() != "" && supportFiles(e.Model) {
				// ensure that there is a trailing slash so that the list iterator could start walking from the prefix dir
				// (https://github.com/pocketbase/pocketbase/discussions/5246#discussioncomment-10128955)
				prefix := strings.TrimRight(m.BaseFilesPath(), "/") + "/"

				// note: for now assume no context cancellation
				err := deleteSem.Acquire(context.Background(), 1)
				if err != nil {
					app.Logger().Error(
						"Failed to delete storage prefix (couldn't acquire a worker)",
						slog.String("prefix", prefix),
						slog.String("error", err.Error()),
					)
				} else {
					// run in the background for "optimistic" delete to avoid blocking the delete transaction
					routine.FireAndForget(func() {
						defer deleteSem.Release(1)

						if err := deletePrefix(prefix); err != nil {
							app.Logger().Error(
								"Failed to delete storage prefix (non critical error; usually could happen because of S3 api limits)",
								slog.String("prefix", prefix),
								slog.String("error", err.Error()),
							)
						}
					})
				}
			}

			return e.Next()
		},
		Priority: -99,
	})

	app.OnServe().Bind(&hook.Handler[*ServeEvent]{
		Id: "__pbCronStart__",
		Func: func(e *ServeEvent) error {
			app.Cron().Start()

			return e.Next()
		},
		Priority: 999,
	})

	app.Cron().Add("__pbDBOptimize__", "0 0 * * *", func() {
		_, execErr := app.NonconcurrentDB().NewQuery("PRAGMA wal_checkpoint(TRUNCATE)").Execute()
		if execErr != nil {
			app.Logger().Warn("Failed to run periodic PRAGMA wal_checkpoint for the main DB", slog.String("error", execErr.Error()))
		}

		_, execErr = app.AuxNonconcurrentDB().NewQuery("PRAGMA wal_checkpoint(TRUNCATE)").Execute()
		if execErr != nil {
			app.Logger().Warn("Failed to run periodic PRAGMA wal_checkpoint for the auxiliary DB", slog.String("error", execErr.Error()))
		}

		_, execErr = app.NonconcurrentDB().NewQuery("PRAGMA optimize").Execute()
		if execErr != nil {
			app.Logger().Warn("Failed to run periodic PRAGMA optimize", slog.String("error", execErr.Error()))
		}
	})

	app.registerSettingsHooks()
	app.registerAutobackupHooks()
	app.registerCollectionHooks()
	app.registerRecordHooks()
	app.registerSuperuserHooks()
	app.registerExternalAuthHooks()
	app.registerMFAHooks()
	app.registerOTPHooks()
	app.registerAuthOriginHooks()
}

// getLoggerMinLevel returns the logger min level based on the
// app configurations (dev mode, settings, etc.).
//
// If not in dev mode - returns the level from the app settings.
//
// If the app is in dev mode it returns -9999 level allowing to print
// practically all logs to the terminal.
// In this case DB logs are still filtered but the checks for the min level are done
// in the BatchOptions.BeforeAddFunc instead of the slog.Handler.Enabled() method.
func getLoggerMinLevel(app App) slog.Level {
	var minLevel slog.Level

	if app.IsDev() {
		minLevel = -99999
	} else if app.Settings() != nil {
		minLevel = slog.Level(app.Settings().Logs.MinLevel)
	}

	return minLevel
}

func (app *BaseApp) initLogger() error {
	duration := 3 * time.Second
	ticker := time.NewTicker(duration)
	done := make(chan bool, 1)

	handler := logger.NewBatchHandler(logger.BatchOptions{
		Level:     getLoggerMinLevel(app),
		BatchSize: 200,
		BeforeAddFunc: func(ctx context.Context, log *logger.Log) bool {
			if app.IsDev() {
				printLog(log)

				// manually check the log level and skip if necessary
				if log.Level < slog.Level(app.Settings().Logs.MinLevel) {
					return false
				}
			}

			ticker.Reset(duration)

			return app.Settings().Logs.MaxDays > 0
		},
		WriteFunc: func(ctx context.Context, logs []*logger.Log) error {
			if !app.IsBootstrapped() || app.Settings().Logs.MaxDays == 0 {
				return nil
			}

			// write the accumulated logs
			// (note: based on several local tests there is no significant performance difference between small number of separate write queries vs 1 big INSERT)
			app.AuxRunInTransaction(func(txApp App) error {
				model := &Log{}
				for _, l := range logs {
					model.MarkAsNew()
					model.Id = GenerateDefaultRandomId()
					model.Level = int(l.Level)
					model.Message = l.Message
					model.Data = l.Data
					model.Created, _ = types.ParseDateTime(l.Time)

					if err := txApp.AuxSave(model); err != nil {
						log.Println("Failed to write log", model, err)
					}
				}

				return nil
			})

			return nil
		},
	})

	go func() {
		ctx := context.Background()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				handler.WriteAll(ctx)
			}
		}
	}()

	app.logger = slog.New(handler)

	// write all remaining logs before ticker.Stop to avoid races with ResetBootstrap user calls
	app.OnTerminate().Bind(&hook.Handler[*TerminateEvent]{
		Id: "__pbAppLoggerOnTerminate__",
		Func: func(e *TerminateEvent) error {
			handler.WriteAll(context.Background())

			ticker.Stop()

			// don't block in case OnTerminate is triggered more than once
			select {
			case done <- true:
			default:
			}

			return e.Next()
		},
		Priority: -999,
	})

	// reload log handler level (if initialized)
	app.OnSettingsReload().Bind(&hook.Handler[*SettingsReloadEvent]{
		Id: "__pbAppLoggerOnSettingsReload__",
		Func: func(e *SettingsReloadEvent) error {
			err := e.Next()
			if err != nil {
				return err
			}

			if e.App.Logger() != nil {
				if h, ok := e.App.Logger().Handler().(*logger.BatchHandler); ok {
					h.SetLevel(getLoggerMinLevel(e.App))
				}
			}

			// try to clear old logs not matching the new settings
			createdBefore := types.NowDateTime().AddDate(0, 0, -1*e.App.Settings().Logs.MaxDays)
			expr := dbx.NewExp("[[created]] <= {:date} OR [[level]] < {:level}", dbx.Params{
				"date":  createdBefore.String(),
				"level": e.App.Settings().Logs.MinLevel,
			})
			_, err = e.App.AuxNonconcurrentDB().Delete((&Log{}).TableName(), expr).Execute()
			if err != nil {
				e.App.Logger().Debug("Failed to cleanup old logs", "error", err)
			}

			// no logs are allowed -> try to reclaim preserved disk space after the previous delete operation
			if e.App.Settings().Logs.MaxDays == 0 {
				err = e.App.AuxVacuum()
				if err != nil {
					e.App.Logger().Debug("Failed to VACUUM aux database", "error", err)
				}
			}

			return nil
		},
		Priority: -999,
	})

	// cleanup old logs
	app.Cron().Add("__pbLogsCleanup__", "0 */6 * * *", func() {
		deleteErr := app.DeleteOldLogs(time.Now().AddDate(0, 0, -1*app.Settings().Logs.MaxDays))
		if deleteErr != nil {
			app.Logger().Warn("Failed to delete old logs", "error", deleteErr)
		}
	})

	return nil
}
