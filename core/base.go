package core

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/logger"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

const (
	DefaultDataMaxOpenConns int = 120
	DefaultDataMaxIdleConns int = 20
	DefaultLogsMaxOpenConns int = 10
	DefaultLogsMaxIdleConns int = 2

	LocalStorageDirName string = "storage"
	LocalBackupsDirName string = "backups"
	LocalTempDirName    string = ".pb_temp_to_delete" // temp pb_data sub directory that will be deleted on each app.Bootstrap()
)

var _ App = (*BaseApp)(nil)

// BaseApp implements core.App and defines the base PocketBase app structure.
type BaseApp struct {
	// @todo consider introducing a mutex to allow safe concurrent config changes during runtime

	// configurable parameters
	isDev            bool
	dataDir          string
	encryptionEnv    string
	dataMaxOpenConns int
	dataMaxIdleConns int
	logsMaxOpenConns int
	logsMaxIdleConns int

	// internals
	store               *store.Store[any]
	settings            *settings.Settings
	dao                 *daos.Dao
	logsDao             *daos.Dao
	subscriptionsBroker *subscriptions.Broker
	logger              *slog.Logger

	// app event hooks
	onBeforeBootstrap *hook.Hook[*BootstrapEvent]
	onAfterBootstrap  *hook.Hook[*BootstrapEvent]
	onBeforeServe     *hook.Hook[*ServeEvent]
	onBeforeApiError  *hook.Hook[*ApiErrorEvent]
	onAfterApiError   *hook.Hook[*ApiErrorEvent]
	onTerminate       *hook.Hook[*TerminateEvent]

	// dao event hooks
	onModelBeforeCreate *hook.Hook[*ModelEvent]
	onModelAfterCreate  *hook.Hook[*ModelEvent]
	onModelBeforeUpdate *hook.Hook[*ModelEvent]
	onModelAfterUpdate  *hook.Hook[*ModelEvent]
	onModelBeforeDelete *hook.Hook[*ModelEvent]
	onModelAfterDelete  *hook.Hook[*ModelEvent]

	// mailer event hooks
	onMailerBeforeAdminResetPasswordSend  *hook.Hook[*MailerAdminEvent]
	onMailerAfterAdminResetPasswordSend   *hook.Hook[*MailerAdminEvent]
	onMailerBeforeRecordResetPasswordSend *hook.Hook[*MailerRecordEvent]
	onMailerAfterRecordResetPasswordSend  *hook.Hook[*MailerRecordEvent]
	onMailerBeforeRecordVerificationSend  *hook.Hook[*MailerRecordEvent]
	onMailerAfterRecordVerificationSend   *hook.Hook[*MailerRecordEvent]
	onMailerBeforeRecordChangeEmailSend   *hook.Hook[*MailerRecordEvent]
	onMailerAfterRecordChangeEmailSend    *hook.Hook[*MailerRecordEvent]

	// realtime api event hooks
	onRealtimeConnectRequest         *hook.Hook[*RealtimeConnectEvent]
	onRealtimeDisconnectRequest      *hook.Hook[*RealtimeDisconnectEvent]
	onRealtimeBeforeMessageSend      *hook.Hook[*RealtimeMessageEvent]
	onRealtimeAfterMessageSend       *hook.Hook[*RealtimeMessageEvent]
	onRealtimeBeforeSubscribeRequest *hook.Hook[*RealtimeSubscribeEvent]
	onRealtimeAfterSubscribeRequest  *hook.Hook[*RealtimeSubscribeEvent]

	// settings api event hooks
	onSettingsListRequest         *hook.Hook[*SettingsListEvent]
	onSettingsBeforeUpdateRequest *hook.Hook[*SettingsUpdateEvent]
	onSettingsAfterUpdateRequest  *hook.Hook[*SettingsUpdateEvent]

	// file api event hooks
	onFileDownloadRequest    *hook.Hook[*FileDownloadEvent]
	onFileBeforeTokenRequest *hook.Hook[*FileTokenEvent]
	onFileAfterTokenRequest  *hook.Hook[*FileTokenEvent]

	// admin api event hooks
	onAdminsListRequest                      *hook.Hook[*AdminsListEvent]
	onAdminViewRequest                       *hook.Hook[*AdminViewEvent]
	onAdminBeforeCreateRequest               *hook.Hook[*AdminCreateEvent]
	onAdminAfterCreateRequest                *hook.Hook[*AdminCreateEvent]
	onAdminBeforeUpdateRequest               *hook.Hook[*AdminUpdateEvent]
	onAdminAfterUpdateRequest                *hook.Hook[*AdminUpdateEvent]
	onAdminBeforeDeleteRequest               *hook.Hook[*AdminDeleteEvent]
	onAdminAfterDeleteRequest                *hook.Hook[*AdminDeleteEvent]
	onAdminAuthRequest                       *hook.Hook[*AdminAuthEvent]
	onAdminBeforeAuthWithPasswordRequest     *hook.Hook[*AdminAuthWithPasswordEvent]
	onAdminAfterAuthWithPasswordRequest      *hook.Hook[*AdminAuthWithPasswordEvent]
	onAdminBeforeAuthRefreshRequest          *hook.Hook[*AdminAuthRefreshEvent]
	onAdminAfterAuthRefreshRequest           *hook.Hook[*AdminAuthRefreshEvent]
	onAdminBeforeRequestPasswordResetRequest *hook.Hook[*AdminRequestPasswordResetEvent]
	onAdminAfterRequestPasswordResetRequest  *hook.Hook[*AdminRequestPasswordResetEvent]
	onAdminBeforeConfirmPasswordResetRequest *hook.Hook[*AdminConfirmPasswordResetEvent]
	onAdminAfterConfirmPasswordResetRequest  *hook.Hook[*AdminConfirmPasswordResetEvent]

	// record auth API event hooks
	onRecordAuthRequest                       *hook.Hook[*RecordAuthEvent]
	onRecordBeforeAuthWithPasswordRequest     *hook.Hook[*RecordAuthWithPasswordEvent]
	onRecordAfterAuthWithPasswordRequest      *hook.Hook[*RecordAuthWithPasswordEvent]
	onRecordBeforeAuthWithOAuth2Request       *hook.Hook[*RecordAuthWithOAuth2Event]
	onRecordAfterAuthWithOAuth2Request        *hook.Hook[*RecordAuthWithOAuth2Event]
	onRecordBeforeAuthRefreshRequest          *hook.Hook[*RecordAuthRefreshEvent]
	onRecordAfterAuthRefreshRequest           *hook.Hook[*RecordAuthRefreshEvent]
	onRecordBeforeRequestPasswordResetRequest *hook.Hook[*RecordRequestPasswordResetEvent]
	onRecordAfterRequestPasswordResetRequest  *hook.Hook[*RecordRequestPasswordResetEvent]
	onRecordBeforeConfirmPasswordResetRequest *hook.Hook[*RecordConfirmPasswordResetEvent]
	onRecordAfterConfirmPasswordResetRequest  *hook.Hook[*RecordConfirmPasswordResetEvent]
	onRecordBeforeRequestVerificationRequest  *hook.Hook[*RecordRequestVerificationEvent]
	onRecordAfterRequestVerificationRequest   *hook.Hook[*RecordRequestVerificationEvent]
	onRecordBeforeConfirmVerificationRequest  *hook.Hook[*RecordConfirmVerificationEvent]
	onRecordAfterConfirmVerificationRequest   *hook.Hook[*RecordConfirmVerificationEvent]
	onRecordBeforeRequestEmailChangeRequest   *hook.Hook[*RecordRequestEmailChangeEvent]
	onRecordAfterRequestEmailChangeRequest    *hook.Hook[*RecordRequestEmailChangeEvent]
	onRecordBeforeConfirmEmailChangeRequest   *hook.Hook[*RecordConfirmEmailChangeEvent]
	onRecordAfterConfirmEmailChangeRequest    *hook.Hook[*RecordConfirmEmailChangeEvent]
	onRecordListExternalAuthsRequest          *hook.Hook[*RecordListExternalAuthsEvent]
	onRecordBeforeUnlinkExternalAuthRequest   *hook.Hook[*RecordUnlinkExternalAuthEvent]
	onRecordAfterUnlinkExternalAuthRequest    *hook.Hook[*RecordUnlinkExternalAuthEvent]

	// record crud API event hooks
	onRecordsListRequest        *hook.Hook[*RecordsListEvent]
	onRecordViewRequest         *hook.Hook[*RecordViewEvent]
	onRecordBeforeCreateRequest *hook.Hook[*RecordCreateEvent]
	onRecordAfterCreateRequest  *hook.Hook[*RecordCreateEvent]
	onRecordBeforeUpdateRequest *hook.Hook[*RecordUpdateEvent]
	onRecordAfterUpdateRequest  *hook.Hook[*RecordUpdateEvent]
	onRecordBeforeDeleteRequest *hook.Hook[*RecordDeleteEvent]
	onRecordAfterDeleteRequest  *hook.Hook[*RecordDeleteEvent]

	// collection API event hooks
	onCollectionsListRequest         *hook.Hook[*CollectionsListEvent]
	onCollectionViewRequest          *hook.Hook[*CollectionViewEvent]
	onCollectionBeforeCreateRequest  *hook.Hook[*CollectionCreateEvent]
	onCollectionAfterCreateRequest   *hook.Hook[*CollectionCreateEvent]
	onCollectionBeforeUpdateRequest  *hook.Hook[*CollectionUpdateEvent]
	onCollectionAfterUpdateRequest   *hook.Hook[*CollectionUpdateEvent]
	onCollectionBeforeDeleteRequest  *hook.Hook[*CollectionDeleteEvent]
	onCollectionAfterDeleteRequest   *hook.Hook[*CollectionDeleteEvent]
	onCollectionsBeforeImportRequest *hook.Hook[*CollectionsImportEvent]
	onCollectionsAfterImportRequest  *hook.Hook[*CollectionsImportEvent]
}

// BaseAppConfig defines a BaseApp configuration option
type BaseAppConfig struct {
	IsDev            bool
	DataDir          string
	EncryptionEnv    string
	DataMaxOpenConns int // default to 500
	DataMaxIdleConns int // default 20
	LogsMaxOpenConns int // default to 100
	LogsMaxIdleConns int // default to 5
}

// NewBaseApp creates and returns a new BaseApp instance
// configured with the provided arguments.
//
// To initialize the app, you need to call `app.Bootstrap()`.
func NewBaseApp(config BaseAppConfig) *BaseApp {
	app := &BaseApp{
		isDev:               config.IsDev,
		dataDir:             config.DataDir,
		encryptionEnv:       config.EncryptionEnv,
		dataMaxOpenConns:    config.DataMaxOpenConns,
		dataMaxIdleConns:    config.DataMaxIdleConns,
		logsMaxOpenConns:    config.LogsMaxOpenConns,
		logsMaxIdleConns:    config.LogsMaxIdleConns,
		store:               store.New[any](nil),
		settings:            settings.New(),
		subscriptionsBroker: subscriptions.NewBroker(),

		// app event hooks
		onBeforeBootstrap: &hook.Hook[*BootstrapEvent]{},
		onAfterBootstrap:  &hook.Hook[*BootstrapEvent]{},
		onBeforeServe:     &hook.Hook[*ServeEvent]{},
		onBeforeApiError:  &hook.Hook[*ApiErrorEvent]{},
		onAfterApiError:   &hook.Hook[*ApiErrorEvent]{},
		onTerminate:       &hook.Hook[*TerminateEvent]{},

		// dao event hooks
		onModelBeforeCreate: &hook.Hook[*ModelEvent]{},
		onModelAfterCreate:  &hook.Hook[*ModelEvent]{},
		onModelBeforeUpdate: &hook.Hook[*ModelEvent]{},
		onModelAfterUpdate:  &hook.Hook[*ModelEvent]{},
		onModelBeforeDelete: &hook.Hook[*ModelEvent]{},
		onModelAfterDelete:  &hook.Hook[*ModelEvent]{},

		// mailer event hooks
		onMailerBeforeAdminResetPasswordSend:  &hook.Hook[*MailerAdminEvent]{},
		onMailerAfterAdminResetPasswordSend:   &hook.Hook[*MailerAdminEvent]{},
		onMailerBeforeRecordResetPasswordSend: &hook.Hook[*MailerRecordEvent]{},
		onMailerAfterRecordResetPasswordSend:  &hook.Hook[*MailerRecordEvent]{},
		onMailerBeforeRecordVerificationSend:  &hook.Hook[*MailerRecordEvent]{},
		onMailerAfterRecordVerificationSend:   &hook.Hook[*MailerRecordEvent]{},
		onMailerBeforeRecordChangeEmailSend:   &hook.Hook[*MailerRecordEvent]{},
		onMailerAfterRecordChangeEmailSend:    &hook.Hook[*MailerRecordEvent]{},

		// realtime API event hooks
		onRealtimeConnectRequest:         &hook.Hook[*RealtimeConnectEvent]{},
		onRealtimeDisconnectRequest:      &hook.Hook[*RealtimeDisconnectEvent]{},
		onRealtimeBeforeMessageSend:      &hook.Hook[*RealtimeMessageEvent]{},
		onRealtimeAfterMessageSend:       &hook.Hook[*RealtimeMessageEvent]{},
		onRealtimeBeforeSubscribeRequest: &hook.Hook[*RealtimeSubscribeEvent]{},
		onRealtimeAfterSubscribeRequest:  &hook.Hook[*RealtimeSubscribeEvent]{},

		// settings API event hooks
		onSettingsListRequest:         &hook.Hook[*SettingsListEvent]{},
		onSettingsBeforeUpdateRequest: &hook.Hook[*SettingsUpdateEvent]{},
		onSettingsAfterUpdateRequest:  &hook.Hook[*SettingsUpdateEvent]{},

		// file API event hooks
		onFileDownloadRequest:    &hook.Hook[*FileDownloadEvent]{},
		onFileBeforeTokenRequest: &hook.Hook[*FileTokenEvent]{},
		onFileAfterTokenRequest:  &hook.Hook[*FileTokenEvent]{},

		// admin API event hooks
		onAdminsListRequest:                      &hook.Hook[*AdminsListEvent]{},
		onAdminViewRequest:                       &hook.Hook[*AdminViewEvent]{},
		onAdminBeforeCreateRequest:               &hook.Hook[*AdminCreateEvent]{},
		onAdminAfterCreateRequest:                &hook.Hook[*AdminCreateEvent]{},
		onAdminBeforeUpdateRequest:               &hook.Hook[*AdminUpdateEvent]{},
		onAdminAfterUpdateRequest:                &hook.Hook[*AdminUpdateEvent]{},
		onAdminBeforeDeleteRequest:               &hook.Hook[*AdminDeleteEvent]{},
		onAdminAfterDeleteRequest:                &hook.Hook[*AdminDeleteEvent]{},
		onAdminAuthRequest:                       &hook.Hook[*AdminAuthEvent]{},
		onAdminBeforeAuthWithPasswordRequest:     &hook.Hook[*AdminAuthWithPasswordEvent]{},
		onAdminAfterAuthWithPasswordRequest:      &hook.Hook[*AdminAuthWithPasswordEvent]{},
		onAdminBeforeAuthRefreshRequest:          &hook.Hook[*AdminAuthRefreshEvent]{},
		onAdminAfterAuthRefreshRequest:           &hook.Hook[*AdminAuthRefreshEvent]{},
		onAdminBeforeRequestPasswordResetRequest: &hook.Hook[*AdminRequestPasswordResetEvent]{},
		onAdminAfterRequestPasswordResetRequest:  &hook.Hook[*AdminRequestPasswordResetEvent]{},
		onAdminBeforeConfirmPasswordResetRequest: &hook.Hook[*AdminConfirmPasswordResetEvent]{},
		onAdminAfterConfirmPasswordResetRequest:  &hook.Hook[*AdminConfirmPasswordResetEvent]{},

		// record auth API event hooks
		onRecordAuthRequest:                       &hook.Hook[*RecordAuthEvent]{},
		onRecordBeforeAuthWithPasswordRequest:     &hook.Hook[*RecordAuthWithPasswordEvent]{},
		onRecordAfterAuthWithPasswordRequest:      &hook.Hook[*RecordAuthWithPasswordEvent]{},
		onRecordBeforeAuthWithOAuth2Request:       &hook.Hook[*RecordAuthWithOAuth2Event]{},
		onRecordAfterAuthWithOAuth2Request:        &hook.Hook[*RecordAuthWithOAuth2Event]{},
		onRecordBeforeAuthRefreshRequest:          &hook.Hook[*RecordAuthRefreshEvent]{},
		onRecordAfterAuthRefreshRequest:           &hook.Hook[*RecordAuthRefreshEvent]{},
		onRecordBeforeRequestPasswordResetRequest: &hook.Hook[*RecordRequestPasswordResetEvent]{},
		onRecordAfterRequestPasswordResetRequest:  &hook.Hook[*RecordRequestPasswordResetEvent]{},
		onRecordBeforeConfirmPasswordResetRequest: &hook.Hook[*RecordConfirmPasswordResetEvent]{},
		onRecordAfterConfirmPasswordResetRequest:  &hook.Hook[*RecordConfirmPasswordResetEvent]{},
		onRecordBeforeRequestVerificationRequest:  &hook.Hook[*RecordRequestVerificationEvent]{},
		onRecordAfterRequestVerificationRequest:   &hook.Hook[*RecordRequestVerificationEvent]{},
		onRecordBeforeConfirmVerificationRequest:  &hook.Hook[*RecordConfirmVerificationEvent]{},
		onRecordAfterConfirmVerificationRequest:   &hook.Hook[*RecordConfirmVerificationEvent]{},
		onRecordBeforeRequestEmailChangeRequest:   &hook.Hook[*RecordRequestEmailChangeEvent]{},
		onRecordAfterRequestEmailChangeRequest:    &hook.Hook[*RecordRequestEmailChangeEvent]{},
		onRecordBeforeConfirmEmailChangeRequest:   &hook.Hook[*RecordConfirmEmailChangeEvent]{},
		onRecordAfterConfirmEmailChangeRequest:    &hook.Hook[*RecordConfirmEmailChangeEvent]{},
		onRecordListExternalAuthsRequest:          &hook.Hook[*RecordListExternalAuthsEvent]{},
		onRecordBeforeUnlinkExternalAuthRequest:   &hook.Hook[*RecordUnlinkExternalAuthEvent]{},
		onRecordAfterUnlinkExternalAuthRequest:    &hook.Hook[*RecordUnlinkExternalAuthEvent]{},

		// record crud API event hooks
		onRecordsListRequest:        &hook.Hook[*RecordsListEvent]{},
		onRecordViewRequest:         &hook.Hook[*RecordViewEvent]{},
		onRecordBeforeCreateRequest: &hook.Hook[*RecordCreateEvent]{},
		onRecordAfterCreateRequest:  &hook.Hook[*RecordCreateEvent]{},
		onRecordBeforeUpdateRequest: &hook.Hook[*RecordUpdateEvent]{},
		onRecordAfterUpdateRequest:  &hook.Hook[*RecordUpdateEvent]{},
		onRecordBeforeDeleteRequest: &hook.Hook[*RecordDeleteEvent]{},
		onRecordAfterDeleteRequest:  &hook.Hook[*RecordDeleteEvent]{},

		// collection API event hooks
		onCollectionsListRequest:         &hook.Hook[*CollectionsListEvent]{},
		onCollectionViewRequest:          &hook.Hook[*CollectionViewEvent]{},
		onCollectionBeforeCreateRequest:  &hook.Hook[*CollectionCreateEvent]{},
		onCollectionAfterCreateRequest:   &hook.Hook[*CollectionCreateEvent]{},
		onCollectionBeforeUpdateRequest:  &hook.Hook[*CollectionUpdateEvent]{},
		onCollectionAfterUpdateRequest:   &hook.Hook[*CollectionUpdateEvent]{},
		onCollectionBeforeDeleteRequest:  &hook.Hook[*CollectionDeleteEvent]{},
		onCollectionAfterDeleteRequest:   &hook.Hook[*CollectionDeleteEvent]{},
		onCollectionsBeforeImportRequest: &hook.Hook[*CollectionsImportEvent]{},
		onCollectionsAfterImportRequest:  &hook.Hook[*CollectionsImportEvent]{},
	}

	app.registerDefaultHooks()

	return app
}

// IsBootstrapped checks if the application was initialized
// (aka. whether Bootstrap() was called).
func (app *BaseApp) IsBootstrapped() bool {
	return app.dao != nil && app.logsDao != nil && app.settings != nil
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

// Bootstrap initializes the application
// (aka. create data dir, open db connections, load settings, etc.).
//
// It will call ResetBootstrapState() if the application was already bootstrapped.
func (app *BaseApp) Bootstrap() error {
	event := &BootstrapEvent{app}

	if err := app.OnBeforeBootstrap().Trigger(event); err != nil {
		return err
	}

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

	if err := app.initLogsDB(); err != nil {
		return err
	}

	if err := app.initLogger(); err != nil {
		return err
	}

	// we don't check for an error because the db migrations may have not been executed yet
	app.RefreshSettings()

	// cleanup the pb_data temp directory (if any)
	os.RemoveAll(filepath.Join(app.DataDir(), LocalTempDirName))

	return app.OnAfterBootstrap().Trigger(event)
}

// ResetBootstrapState takes care for releasing initialized app resources
// (eg. closing db connections).
func (app *BaseApp) ResetBootstrapState() error {
	if app.Dao() != nil {
		if err := app.Dao().ConcurrentDB().(*dbx.DB).Close(); err != nil {
			return err
		}
		if err := app.Dao().NonconcurrentDB().(*dbx.DB).Close(); err != nil {
			return err
		}
	}

	if app.LogsDao() != nil {
		if err := app.LogsDao().ConcurrentDB().(*dbx.DB).Close(); err != nil {
			return err
		}
		if err := app.LogsDao().NonconcurrentDB().(*dbx.DB).Close(); err != nil {
			return err
		}
	}

	app.dao = nil
	app.logsDao = nil

	return nil
}

// Deprecated:
// This method may get removed in the near future.
// It is recommended to access the db instance from app.Dao().DB() or
// if you want more flexibility - app.Dao().ConcurrentDB() and app.Dao().NonconcurrentDB().
//
// DB returns the default app database instance.
func (app *BaseApp) DB() *dbx.DB {
	if app.Dao() == nil {
		return nil
	}

	db, ok := app.Dao().DB().(*dbx.DB)
	if !ok {
		return nil
	}

	return db
}

// Dao returns the default app Dao instance.
func (app *BaseApp) Dao() *daos.Dao {
	return app.dao
}

// Deprecated:
// This method may get removed in the near future.
// It is recommended to access the logs db instance from app.LogsDao().DB() or
// if you want more flexibility - app.LogsDao().ConcurrentDB() and app.LogsDao().NonconcurrentDB().
//
// LogsDB returns the app logs database instance.
func (app *BaseApp) LogsDB() *dbx.DB {
	if app.LogsDao() == nil {
		return nil
	}

	db, ok := app.LogsDao().DB().(*dbx.DB)
	if !ok {
		return nil
	}

	return db
}

// LogsDao returns the app logs Dao instance.
func (app *BaseApp) LogsDao() *daos.Dao {
	return app.logsDao
}

// DataDir returns the app data directory path.
func (app *BaseApp) DataDir() string {
	return app.dataDir
}

// EncryptionEnv returns the name of the app secret env key
// (used for settings encryption).
func (app *BaseApp) EncryptionEnv() string {
	return app.encryptionEnv
}

// IsDev returns whether the app is in dev mode.
//
// When enabled logs, executed sql statements, etc. are printed to the stderr.
func (app *BaseApp) IsDev() bool {
	return app.isDev
}

// Settings returns the loaded app settings.
func (app *BaseApp) Settings() *settings.Settings {
	return app.settings
}

// Deprecated: Use app.Store() instead.
func (app *BaseApp) Cache() *store.Store[any] {
	color.Yellow("app.Store() is soft-deprecated. Please replace it with app.Store().")
	return app.Store()
}

// Store returns the app internal runtime store.
func (app *BaseApp) Store() *store.Store[any] {
	return app.store
}

// SubscriptionsBroker returns the app realtime subscriptions broker instance.
func (app *BaseApp) SubscriptionsBroker() *subscriptions.Broker {
	return app.subscriptionsBroker
}

// NewMailClient creates and returns a new SMTP or Sendmail client
// based on the current app settings.
func (app *BaseApp) NewMailClient() mailer.Mailer {
	if app.Settings().Smtp.Enabled {
		return &mailer.SmtpClient{
			Host:       app.Settings().Smtp.Host,
			Port:       app.Settings().Smtp.Port,
			Username:   app.Settings().Smtp.Username,
			Password:   app.Settings().Smtp.Password,
			Tls:        app.Settings().Smtp.Tls,
			AuthMethod: app.Settings().Smtp.AuthMethod,
			LocalName:  app.Settings().Smtp.LocalName,
		}
	}

	return &mailer.Sendmail{}
}

// NewFilesystem creates a new local or S3 filesystem instance
// for managing regular app files (eg. collection uploads)
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

// NewFilesystem creates a new local or S3 filesystem instance
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

	return app.OnTerminate().Trigger(&TerminateEvent{
		App:       app,
		IsRestart: true,
	}, func(e *TerminateEvent) error {
		e.App.ResetBootstrapState()

		// attempt to restart the bootstrap process in case execve returns an error for some reason
		defer e.App.Bootstrap()

		return syscall.Exec(execPath, os.Args, os.Environ())
	})
}

// RefreshSettings reinitializes and reloads the stored application settings.
func (app *BaseApp) RefreshSettings() error {
	if app.settings == nil {
		app.settings = settings.New()
	}

	encryptionKey := os.Getenv(app.EncryptionEnv())

	storedSettings, err := app.Dao().FindSettings(encryptionKey)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// no settings were previously stored
	if storedSettings == nil {
		return app.Dao().SaveSettings(app.settings, encryptionKey)
	}

	// load the settings from the stored param into the app ones
	if err := app.settings.Merge(storedSettings); err != nil {
		return err
	}

	// reload handler level (if initialized)
	if app.Logger() != nil {
		if h, ok := app.Logger().Handler().(*logger.BatchHandler); ok {
			h.SetLevel(app.getLoggerMinLevel())
		}
	}

	return nil
}

// -------------------------------------------------------------------
// App event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnBeforeBootstrap() *hook.Hook[*BootstrapEvent] {
	return app.onBeforeBootstrap
}

func (app *BaseApp) OnAfterBootstrap() *hook.Hook[*BootstrapEvent] {
	return app.onAfterBootstrap
}

func (app *BaseApp) OnBeforeServe() *hook.Hook[*ServeEvent] {
	return app.onBeforeServe
}

func (app *BaseApp) OnBeforeApiError() *hook.Hook[*ApiErrorEvent] {
	return app.onBeforeApiError
}

func (app *BaseApp) OnAfterApiError() *hook.Hook[*ApiErrorEvent] {
	return app.onAfterApiError
}

func (app *BaseApp) OnTerminate() *hook.Hook[*TerminateEvent] {
	return app.onTerminate
}

// -------------------------------------------------------------------
// Dao event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnModelBeforeCreate(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelBeforeCreate, tags...)
}

func (app *BaseApp) OnModelAfterCreate(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelAfterCreate, tags...)
}

func (app *BaseApp) OnModelBeforeUpdate(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelBeforeUpdate, tags...)
}

func (app *BaseApp) OnModelAfterUpdate(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelAfterUpdate, tags...)
}

func (app *BaseApp) OnModelBeforeDelete(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelBeforeDelete, tags...)
}

func (app *BaseApp) OnModelAfterDelete(tags ...string) *hook.TaggedHook[*ModelEvent] {
	return hook.NewTaggedHook(app.onModelAfterDelete, tags...)
}

// -------------------------------------------------------------------
// Mailer event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnMailerBeforeAdminResetPasswordSend() *hook.Hook[*MailerAdminEvent] {
	return app.onMailerBeforeAdminResetPasswordSend
}

func (app *BaseApp) OnMailerAfterAdminResetPasswordSend() *hook.Hook[*MailerAdminEvent] {
	return app.onMailerAfterAdminResetPasswordSend
}

func (app *BaseApp) OnMailerBeforeRecordResetPasswordSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerBeforeRecordResetPasswordSend, tags...)
}

func (app *BaseApp) OnMailerAfterRecordResetPasswordSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerAfterRecordResetPasswordSend, tags...)
}

func (app *BaseApp) OnMailerBeforeRecordVerificationSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerBeforeRecordVerificationSend, tags...)
}

func (app *BaseApp) OnMailerAfterRecordVerificationSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerAfterRecordVerificationSend, tags...)
}

func (app *BaseApp) OnMailerBeforeRecordChangeEmailSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerBeforeRecordChangeEmailSend, tags...)
}

func (app *BaseApp) OnMailerAfterRecordChangeEmailSend(tags ...string) *hook.TaggedHook[*MailerRecordEvent] {
	return hook.NewTaggedHook(app.onMailerAfterRecordChangeEmailSend, tags...)
}

// -------------------------------------------------------------------
// Realtime API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRealtimeConnectRequest() *hook.Hook[*RealtimeConnectEvent] {
	return app.onRealtimeConnectRequest
}

func (app *BaseApp) OnRealtimeDisconnectRequest() *hook.Hook[*RealtimeDisconnectEvent] {
	return app.onRealtimeDisconnectRequest
}

func (app *BaseApp) OnRealtimeBeforeMessageSend() *hook.Hook[*RealtimeMessageEvent] {
	return app.onRealtimeBeforeMessageSend
}

func (app *BaseApp) OnRealtimeAfterMessageSend() *hook.Hook[*RealtimeMessageEvent] {
	return app.onRealtimeAfterMessageSend
}

func (app *BaseApp) OnRealtimeBeforeSubscribeRequest() *hook.Hook[*RealtimeSubscribeEvent] {
	return app.onRealtimeBeforeSubscribeRequest
}

func (app *BaseApp) OnRealtimeAfterSubscribeRequest() *hook.Hook[*RealtimeSubscribeEvent] {
	return app.onRealtimeAfterSubscribeRequest
}

// -------------------------------------------------------------------
// Settings API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnSettingsListRequest() *hook.Hook[*SettingsListEvent] {
	return app.onSettingsListRequest
}

func (app *BaseApp) OnSettingsBeforeUpdateRequest() *hook.Hook[*SettingsUpdateEvent] {
	return app.onSettingsBeforeUpdateRequest
}

func (app *BaseApp) OnSettingsAfterUpdateRequest() *hook.Hook[*SettingsUpdateEvent] {
	return app.onSettingsAfterUpdateRequest
}

// -------------------------------------------------------------------
// File API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnFileDownloadRequest(tags ...string) *hook.TaggedHook[*FileDownloadEvent] {
	return hook.NewTaggedHook(app.onFileDownloadRequest, tags...)
}

func (app *BaseApp) OnFileBeforeTokenRequest(tags ...string) *hook.TaggedHook[*FileTokenEvent] {
	return hook.NewTaggedHook(app.onFileBeforeTokenRequest, tags...)
}

func (app *BaseApp) OnFileAfterTokenRequest(tags ...string) *hook.TaggedHook[*FileTokenEvent] {
	return hook.NewTaggedHook(app.onFileAfterTokenRequest, tags...)
}

// -------------------------------------------------------------------
// Admin API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnAdminsListRequest() *hook.Hook[*AdminsListEvent] {
	return app.onAdminsListRequest
}

func (app *BaseApp) OnAdminViewRequest() *hook.Hook[*AdminViewEvent] {
	return app.onAdminViewRequest
}

func (app *BaseApp) OnAdminBeforeCreateRequest() *hook.Hook[*AdminCreateEvent] {
	return app.onAdminBeforeCreateRequest
}

func (app *BaseApp) OnAdminAfterCreateRequest() *hook.Hook[*AdminCreateEvent] {
	return app.onAdminAfterCreateRequest
}

func (app *BaseApp) OnAdminBeforeUpdateRequest() *hook.Hook[*AdminUpdateEvent] {
	return app.onAdminBeforeUpdateRequest
}

func (app *BaseApp) OnAdminAfterUpdateRequest() *hook.Hook[*AdminUpdateEvent] {
	return app.onAdminAfterUpdateRequest
}

func (app *BaseApp) OnAdminBeforeDeleteRequest() *hook.Hook[*AdminDeleteEvent] {
	return app.onAdminBeforeDeleteRequest
}

func (app *BaseApp) OnAdminAfterDeleteRequest() *hook.Hook[*AdminDeleteEvent] {
	return app.onAdminAfterDeleteRequest
}

func (app *BaseApp) OnAdminAuthRequest() *hook.Hook[*AdminAuthEvent] {
	return app.onAdminAuthRequest
}

func (app *BaseApp) OnAdminBeforeAuthWithPasswordRequest() *hook.Hook[*AdminAuthWithPasswordEvent] {
	return app.onAdminBeforeAuthWithPasswordRequest
}

func (app *BaseApp) OnAdminAfterAuthWithPasswordRequest() *hook.Hook[*AdminAuthWithPasswordEvent] {
	return app.onAdminAfterAuthWithPasswordRequest
}

func (app *BaseApp) OnAdminBeforeAuthRefreshRequest() *hook.Hook[*AdminAuthRefreshEvent] {
	return app.onAdminBeforeAuthRefreshRequest
}

func (app *BaseApp) OnAdminAfterAuthRefreshRequest() *hook.Hook[*AdminAuthRefreshEvent] {
	return app.onAdminAfterAuthRefreshRequest
}

func (app *BaseApp) OnAdminBeforeRequestPasswordResetRequest() *hook.Hook[*AdminRequestPasswordResetEvent] {
	return app.onAdminBeforeRequestPasswordResetRequest
}

func (app *BaseApp) OnAdminAfterRequestPasswordResetRequest() *hook.Hook[*AdminRequestPasswordResetEvent] {
	return app.onAdminAfterRequestPasswordResetRequest
}

func (app *BaseApp) OnAdminBeforeConfirmPasswordResetRequest() *hook.Hook[*AdminConfirmPasswordResetEvent] {
	return app.onAdminBeforeConfirmPasswordResetRequest
}

func (app *BaseApp) OnAdminAfterConfirmPasswordResetRequest() *hook.Hook[*AdminConfirmPasswordResetEvent] {
	return app.onAdminAfterConfirmPasswordResetRequest
}

// -------------------------------------------------------------------
// Record auth API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRecordAuthRequest(tags ...string) *hook.TaggedHook[*RecordAuthEvent] {
	return hook.NewTaggedHook(app.onRecordAuthRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*RecordAuthWithPasswordEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeAuthWithPasswordRequest, tags...)
}

func (app *BaseApp) OnRecordAfterAuthWithPasswordRequest(tags ...string) *hook.TaggedHook[*RecordAuthWithPasswordEvent] {
	return hook.NewTaggedHook(app.onRecordAfterAuthWithPasswordRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeAuthWithOAuth2Request(tags ...string) *hook.TaggedHook[*RecordAuthWithOAuth2Event] {
	return hook.NewTaggedHook(app.onRecordBeforeAuthWithOAuth2Request, tags...)
}

func (app *BaseApp) OnRecordAfterAuthWithOAuth2Request(tags ...string) *hook.TaggedHook[*RecordAuthWithOAuth2Event] {
	return hook.NewTaggedHook(app.onRecordAfterAuthWithOAuth2Request, tags...)
}

func (app *BaseApp) OnRecordBeforeAuthRefreshRequest(tags ...string) *hook.TaggedHook[*RecordAuthRefreshEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeAuthRefreshRequest, tags...)
}

func (app *BaseApp) OnRecordAfterAuthRefreshRequest(tags ...string) *hook.TaggedHook[*RecordAuthRefreshEvent] {
	return hook.NewTaggedHook(app.onRecordAfterAuthRefreshRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeRequestPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordRequestPasswordResetEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeRequestPasswordResetRequest, tags...)
}

func (app *BaseApp) OnRecordAfterRequestPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordRequestPasswordResetEvent] {
	return hook.NewTaggedHook(app.onRecordAfterRequestPasswordResetRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeConfirmPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordConfirmPasswordResetEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeConfirmPasswordResetRequest, tags...)
}

func (app *BaseApp) OnRecordAfterConfirmPasswordResetRequest(tags ...string) *hook.TaggedHook[*RecordConfirmPasswordResetEvent] {
	return hook.NewTaggedHook(app.onRecordAfterConfirmPasswordResetRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeRequestVerificationRequest(tags ...string) *hook.TaggedHook[*RecordRequestVerificationEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeRequestVerificationRequest, tags...)
}

func (app *BaseApp) OnRecordAfterRequestVerificationRequest(tags ...string) *hook.TaggedHook[*RecordRequestVerificationEvent] {
	return hook.NewTaggedHook(app.onRecordAfterRequestVerificationRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeConfirmVerificationRequest(tags ...string) *hook.TaggedHook[*RecordConfirmVerificationEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeConfirmVerificationRequest, tags...)
}

func (app *BaseApp) OnRecordAfterConfirmVerificationRequest(tags ...string) *hook.TaggedHook[*RecordConfirmVerificationEvent] {
	return hook.NewTaggedHook(app.onRecordAfterConfirmVerificationRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeRequestEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordRequestEmailChangeEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeRequestEmailChangeRequest, tags...)
}

func (app *BaseApp) OnRecordAfterRequestEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordRequestEmailChangeEvent] {
	return hook.NewTaggedHook(app.onRecordAfterRequestEmailChangeRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeConfirmEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordConfirmEmailChangeEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeConfirmEmailChangeRequest, tags...)
}

func (app *BaseApp) OnRecordAfterConfirmEmailChangeRequest(tags ...string) *hook.TaggedHook[*RecordConfirmEmailChangeEvent] {
	return hook.NewTaggedHook(app.onRecordAfterConfirmEmailChangeRequest, tags...)
}

func (app *BaseApp) OnRecordListExternalAuthsRequest(tags ...string) *hook.TaggedHook[*RecordListExternalAuthsEvent] {
	return hook.NewTaggedHook(app.onRecordListExternalAuthsRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeUnlinkExternalAuthRequest(tags ...string) *hook.TaggedHook[*RecordUnlinkExternalAuthEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeUnlinkExternalAuthRequest, tags...)
}

func (app *BaseApp) OnRecordAfterUnlinkExternalAuthRequest(tags ...string) *hook.TaggedHook[*RecordUnlinkExternalAuthEvent] {
	return hook.NewTaggedHook(app.onRecordAfterUnlinkExternalAuthRequest, tags...)
}

// -------------------------------------------------------------------
// Record CRUD API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRecordsListRequest(tags ...string) *hook.TaggedHook[*RecordsListEvent] {
	return hook.NewTaggedHook(app.onRecordsListRequest, tags...)
}

func (app *BaseApp) OnRecordViewRequest(tags ...string) *hook.TaggedHook[*RecordViewEvent] {
	return hook.NewTaggedHook(app.onRecordViewRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeCreateRequest(tags ...string) *hook.TaggedHook[*RecordCreateEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeCreateRequest, tags...)
}

func (app *BaseApp) OnRecordAfterCreateRequest(tags ...string) *hook.TaggedHook[*RecordCreateEvent] {
	return hook.NewTaggedHook(app.onRecordAfterCreateRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeUpdateRequest(tags ...string) *hook.TaggedHook[*RecordUpdateEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeUpdateRequest, tags...)
}

func (app *BaseApp) OnRecordAfterUpdateRequest(tags ...string) *hook.TaggedHook[*RecordUpdateEvent] {
	return hook.NewTaggedHook(app.onRecordAfterUpdateRequest, tags...)
}

func (app *BaseApp) OnRecordBeforeDeleteRequest(tags ...string) *hook.TaggedHook[*RecordDeleteEvent] {
	return hook.NewTaggedHook(app.onRecordBeforeDeleteRequest, tags...)
}

func (app *BaseApp) OnRecordAfterDeleteRequest(tags ...string) *hook.TaggedHook[*RecordDeleteEvent] {
	return hook.NewTaggedHook(app.onRecordAfterDeleteRequest, tags...)
}

// -------------------------------------------------------------------
// Collection API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnCollectionsListRequest() *hook.Hook[*CollectionsListEvent] {
	return app.onCollectionsListRequest
}

func (app *BaseApp) OnCollectionViewRequest() *hook.Hook[*CollectionViewEvent] {
	return app.onCollectionViewRequest
}

func (app *BaseApp) OnCollectionBeforeCreateRequest() *hook.Hook[*CollectionCreateEvent] {
	return app.onCollectionBeforeCreateRequest
}

func (app *BaseApp) OnCollectionAfterCreateRequest() *hook.Hook[*CollectionCreateEvent] {
	return app.onCollectionAfterCreateRequest
}

func (app *BaseApp) OnCollectionBeforeUpdateRequest() *hook.Hook[*CollectionUpdateEvent] {
	return app.onCollectionBeforeUpdateRequest
}

func (app *BaseApp) OnCollectionAfterUpdateRequest() *hook.Hook[*CollectionUpdateEvent] {
	return app.onCollectionAfterUpdateRequest
}

func (app *BaseApp) OnCollectionBeforeDeleteRequest() *hook.Hook[*CollectionDeleteEvent] {
	return app.onCollectionBeforeDeleteRequest
}

func (app *BaseApp) OnCollectionAfterDeleteRequest() *hook.Hook[*CollectionDeleteEvent] {
	return app.onCollectionAfterDeleteRequest
}

func (app *BaseApp) OnCollectionsBeforeImportRequest() *hook.Hook[*CollectionsImportEvent] {
	return app.onCollectionsBeforeImportRequest
}

func (app *BaseApp) OnCollectionsAfterImportRequest() *hook.Hook[*CollectionsImportEvent] {
	return app.onCollectionsAfterImportRequest
}

// -------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------

func (app *BaseApp) initLogsDB() error {
	maxOpenConns := DefaultLogsMaxOpenConns
	maxIdleConns := DefaultLogsMaxIdleConns
	if app.logsMaxOpenConns > 0 {
		maxOpenConns = app.logsMaxOpenConns
	}
	if app.logsMaxIdleConns > 0 {
		maxIdleConns = app.logsMaxIdleConns
	}

	concurrentDB, err := connectDB(filepath.Join(app.DataDir(), "logs.db"))
	if err != nil {
		return err
	}
	concurrentDB.DB().SetMaxOpenConns(maxOpenConns)
	concurrentDB.DB().SetMaxIdleConns(maxIdleConns)
	concurrentDB.DB().SetConnMaxIdleTime(3 * time.Minute)

	nonconcurrentDB, err := connectDB(filepath.Join(app.DataDir(), "logs.db"))
	if err != nil {
		return err
	}
	nonconcurrentDB.DB().SetMaxOpenConns(1)
	nonconcurrentDB.DB().SetMaxIdleConns(1)
	nonconcurrentDB.DB().SetConnMaxIdleTime(3 * time.Minute)

	app.logsDao = daos.NewMultiDB(concurrentDB, nonconcurrentDB)

	return nil
}

func (app *BaseApp) initDataDB() error {
	maxOpenConns := DefaultDataMaxOpenConns
	maxIdleConns := DefaultDataMaxIdleConns
	if app.dataMaxOpenConns > 0 {
		maxOpenConns = app.dataMaxOpenConns
	}
	if app.dataMaxIdleConns > 0 {
		maxIdleConns = app.dataMaxIdleConns
	}

	concurrentDB, err := connectDB(filepath.Join(app.DataDir(), "data.db"))
	if err != nil {
		return err
	}
	concurrentDB.DB().SetMaxOpenConns(maxOpenConns)
	concurrentDB.DB().SetMaxIdleConns(maxIdleConns)
	concurrentDB.DB().SetConnMaxIdleTime(3 * time.Minute)

	nonconcurrentDB, err := connectDB(filepath.Join(app.DataDir(), "data.db"))
	if err != nil {
		return err
	}
	nonconcurrentDB.DB().SetMaxOpenConns(1)
	nonconcurrentDB.DB().SetMaxIdleConns(1)
	nonconcurrentDB.DB().SetConnMaxIdleTime(3 * time.Minute)

	if app.IsDev() {
		nonconcurrentDB.QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
			color.HiBlack("[%.2fms] %v\n", float64(t.Milliseconds()), sql)
		}
		nonconcurrentDB.ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
			color.HiBlack("[%.2fms] %v\n", float64(t.Milliseconds()), sql)
		}
		concurrentDB.QueryLogFunc = nonconcurrentDB.QueryLogFunc
		concurrentDB.ExecLogFunc = nonconcurrentDB.ExecLogFunc
	}

	app.dao = app.createDaoWithHooks(concurrentDB, nonconcurrentDB)

	return nil
}

func (app *BaseApp) createDaoWithHooks(concurrentDB, nonconcurrentDB dbx.Builder) *daos.Dao {
	dao := daos.NewMultiDB(concurrentDB, nonconcurrentDB)

	dao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		e := new(ModelEvent)
		e.Dao = eventDao
		e.Model = m

		return app.OnModelBeforeCreate().Trigger(e, func(e *ModelEvent) error {
			return action()
		})
	}

	dao.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		e := new(ModelEvent)
		e.Dao = eventDao
		e.Model = m

		return app.OnModelAfterCreate().Trigger(e)
	}

	dao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		e := new(ModelEvent)
		e.Dao = eventDao
		e.Model = m

		return app.OnModelBeforeUpdate().Trigger(e, func(e *ModelEvent) error {
			return action()
		})
	}

	dao.AfterUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		e := new(ModelEvent)
		e.Dao = eventDao
		e.Model = m

		return app.OnModelAfterUpdate().Trigger(e)
	}

	dao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		e := new(ModelEvent)
		e.Dao = eventDao
		e.Model = m

		return app.OnModelBeforeDelete().Trigger(e, func(e *ModelEvent) error {
			return action()
		})
	}

	dao.AfterDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		e := new(ModelEvent)
		e.Dao = eventDao
		e.Model = m

		return app.OnModelAfterDelete().Trigger(e)
	}

	return dao
}

func (app *BaseApp) registerDefaultHooks() {
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

	// try to delete the storage files from deleted Collection, Records, etc. model
	app.OnModelAfterDelete().Add(func(e *ModelEvent) error {
		if m, ok := e.Model.(models.FilesManager); ok && m.BaseFilesPath() != "" {
			// ensure that there is a trailing slash so that the list iterator could start walking from the prefix
			// (https://github.com/pocketbase/pocketbase/discussions/5246#discussioncomment-10128955)
			prefix := strings.TrimRight(m.BaseFilesPath(), "/") + "/"

			// run in the background for "optimistic" delete to avoid
			// blocking the delete transaction
			routine.FireAndForget(func() {
				if err := deletePrefix(prefix); err != nil {
					app.Logger().Error(
						"Failed to delete storage prefix (non critical error; usually could happen because of S3 api limits)",
						slog.String("prefix", prefix),
						slog.String("error", err.Error()),
					)
				}
			})
		}

		return nil
	})

	if err := app.initAutobackupHooks(); err != nil {
		app.Logger().Error("Failed to init auto backup hooks", slog.String("error", err.Error()))
	}

	registerCachedCollectionsAppHooks(app)
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
func (app *BaseApp) getLoggerMinLevel() slog.Level {
	var minLevel slog.Level

	if app.IsDev() {
		minLevel = -9999
	} else if app.Settings() != nil {
		minLevel = slog.Level(app.Settings().Logs.MinLevel)
	}

	return minLevel
}

func (app *BaseApp) initLogger() error {
	duration := 3 * time.Second
	ticker := time.NewTicker(duration)
	done := make(chan bool)

	handler := logger.NewBatchHandler(logger.BatchOptions{
		Level:     app.getLoggerMinLevel(),
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
			app.LogsDao().RunInTransaction(func(txDao *daos.Dao) error {
				model := &models.Log{}
				for _, l := range logs {
					model.MarkAsNew()
					// note: using pseudorandom for a slightly better performance
					model.Id = security.PseudorandomStringWithAlphabet(models.DefaultIdLength, models.DefaultIdAlphabet)
					model.Level = int(l.Level)
					model.Message = l.Message
					model.Data = l.Data
					model.Created, _ = types.ParseDateTime(l.Time)
					model.Updated = model.Created

					if err := txDao.SaveLog(model); err != nil {
						log.Println("Failed to write log", model, err)
					}
				}

				return nil
			})

			// @todo replace with cron so that it doesn't rely on the logs write
			//
			// delete old logs
			// ---
			now := time.Now()
			lastLogsDeletedAt := cast.ToTime(app.Store().Get("lastLogsDeletedAt"))
			if now.Sub(lastLogsDeletedAt).Hours() >= 6 {
				deleteErr := app.LogsDao().DeleteOldLogs(now.AddDate(0, 0, -1*app.Settings().Logs.MaxDays))
				if deleteErr == nil {
					app.Store().Set("lastLogsDeletedAt", now)
				} else {
					log.Println("Logs delete failed", deleteErr)
				}
			}

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

	app.OnTerminate().PreAdd(func(e *TerminateEvent) error {
		// write all remaining logs before ticker.Stop to avoid races with ResetBootstrap user calls
		handler.WriteAll(context.Background())

		ticker.Stop()

		done <- true

		return nil
	})

	return nil
}
