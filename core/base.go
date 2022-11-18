package core

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
)

var _ App = (*BaseApp)(nil)

// BaseApp implements core.App and defines the base PocketBase app structure.
type BaseApp struct {
	// configurable parameters
	isDebug       bool
	dataDir       string
	encryptionEnv string

	// internals
	cache               *store.Store[any]
	settings            *Settings
	db                  *dbx.DB
	dao                 *daos.Dao
	logsDB              *dbx.DB
	logsDao             *daos.Dao
	subscriptionsBroker *subscriptions.Broker

	// serve event hooks
	onBeforeServe *hook.Hook[*ServeEvent]

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
	onRealtimeBeforeSubscribeRequest *hook.Hook[*RealtimeSubscribeEvent]
	onRealtimeAfterSubscribeRequest  *hook.Hook[*RealtimeSubscribeEvent]

	// settings api event hooks
	onSettingsListRequest         *hook.Hook[*SettingsListEvent]
	onSettingsBeforeUpdateRequest *hook.Hook[*SettingsUpdateEvent]
	onSettingsAfterUpdateRequest  *hook.Hook[*SettingsUpdateEvent]

	// file api event hooks
	onFileDownloadRequest *hook.Hook[*FileDownloadEvent]

	// admin api event hooks
	onAdminsListRequest        *hook.Hook[*AdminsListEvent]
	onAdminViewRequest         *hook.Hook[*AdminViewEvent]
	onAdminBeforeCreateRequest *hook.Hook[*AdminCreateEvent]
	onAdminAfterCreateRequest  *hook.Hook[*AdminCreateEvent]
	onAdminBeforeUpdateRequest *hook.Hook[*AdminUpdateEvent]
	onAdminAfterUpdateRequest  *hook.Hook[*AdminUpdateEvent]
	onAdminBeforeDeleteRequest *hook.Hook[*AdminDeleteEvent]
	onAdminAfterDeleteRequest  *hook.Hook[*AdminDeleteEvent]
	onAdminAuthRequest         *hook.Hook[*AdminAuthEvent]

	// user api event hooks
	onRecordAuthRequest                     *hook.Hook[*RecordAuthEvent]
	onRecordListExternalAuthsRequest        *hook.Hook[*RecordListExternalAuthsEvent]
	onRecordBeforeUnlinkExternalAuthRequest *hook.Hook[*RecordUnlinkExternalAuthEvent]
	onRecordAfterUnlinkExternalAuthRequest  *hook.Hook[*RecordUnlinkExternalAuthEvent]

	// record api event hooks
	onRecordsListRequest        *hook.Hook[*RecordsListEvent]
	onRecordViewRequest         *hook.Hook[*RecordViewEvent]
	onRecordBeforeCreateRequest *hook.Hook[*RecordCreateEvent]
	onRecordAfterCreateRequest  *hook.Hook[*RecordCreateEvent]
	onRecordBeforeUpdateRequest *hook.Hook[*RecordUpdateEvent]
	onRecordAfterUpdateRequest  *hook.Hook[*RecordUpdateEvent]
	onRecordBeforeDeleteRequest *hook.Hook[*RecordDeleteEvent]
	onRecordAfterDeleteRequest  *hook.Hook[*RecordDeleteEvent]

	// collection api event hooks
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

	// View api event hooks
	onRecordsFromViewListRequest *hook.Hook[*RecordsFromViewListEvent]
	onViewListRequest            *hook.Hook[*ViewListEvent]
	onViewViewRequest            *hook.Hook[*ViewViewEvent]
	onViewBeforeCreateRequest    *hook.Hook[*ViewCreateEvent]
	onViewAfterCreateRequest     *hook.Hook[*ViewCreateEvent]
	onViewBeforeUpdateRequest    *hook.Hook[*ViewUpdateEvent]
	onViewAfterUpdateRequest     *hook.Hook[*ViewUpdateEvent]
	onViewBeforeDeleteRequest    *hook.Hook[*ViewDeleteEvent]
	onViewAfterDeleteRequest     *hook.Hook[*ViewDeleteEvent]
}

// NewBaseApp creates and returns a new BaseApp instance
// configured with the provided arguments.
//
// To initialize the app, you need to call `app.Bootstrap()`.
func NewBaseApp(dataDir string, encryptionEnv string, isDebug bool) *BaseApp {
	app := &BaseApp{
		dataDir:             dataDir,
		isDebug:             isDebug,
		encryptionEnv:       encryptionEnv,
		cache:               store.New[any](nil),
		settings:            NewSettings(),
		subscriptionsBroker: subscriptions.NewBroker(),

		// serve event hooks
		onBeforeServe: &hook.Hook[*ServeEvent]{},

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
		onRealtimeBeforeSubscribeRequest: &hook.Hook[*RealtimeSubscribeEvent]{},
		onRealtimeAfterSubscribeRequest:  &hook.Hook[*RealtimeSubscribeEvent]{},

		// settings API event hooks
		onSettingsListRequest:         &hook.Hook[*SettingsListEvent]{},
		onSettingsBeforeUpdateRequest: &hook.Hook[*SettingsUpdateEvent]{},
		onSettingsAfterUpdateRequest:  &hook.Hook[*SettingsUpdateEvent]{},

		// file API event hooks
		onFileDownloadRequest: &hook.Hook[*FileDownloadEvent]{},

		// admin API event hooks
		onAdminsListRequest:        &hook.Hook[*AdminsListEvent]{},
		onAdminViewRequest:         &hook.Hook[*AdminViewEvent]{},
		onAdminBeforeCreateRequest: &hook.Hook[*AdminCreateEvent]{},
		onAdminAfterCreateRequest:  &hook.Hook[*AdminCreateEvent]{},
		onAdminBeforeUpdateRequest: &hook.Hook[*AdminUpdateEvent]{},
		onAdminAfterUpdateRequest:  &hook.Hook[*AdminUpdateEvent]{},
		onAdminBeforeDeleteRequest: &hook.Hook[*AdminDeleteEvent]{},
		onAdminAfterDeleteRequest:  &hook.Hook[*AdminDeleteEvent]{},
		onAdminAuthRequest:         &hook.Hook[*AdminAuthEvent]{},

		// user API event hooks
		onRecordAuthRequest:                     &hook.Hook[*RecordAuthEvent]{},
		onRecordListExternalAuthsRequest:        &hook.Hook[*RecordListExternalAuthsEvent]{},
		onRecordBeforeUnlinkExternalAuthRequest: &hook.Hook[*RecordUnlinkExternalAuthEvent]{},
		onRecordAfterUnlinkExternalAuthRequest:  &hook.Hook[*RecordUnlinkExternalAuthEvent]{},

		// record API event hooks
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

		// View api event hooks
		onRecordsFromViewListRequest: &hook.Hook[*RecordsFromViewListEvent]{},
		onViewListRequest:            &hook.Hook[*ViewListEvent]{},
		onViewViewRequest:            &hook.Hook[*ViewViewEvent]{},
		onViewBeforeCreateRequest:    &hook.Hook[*ViewCreateEvent]{},
		onViewAfterCreateRequest:     &hook.Hook[*ViewCreateEvent]{},
		onViewBeforeUpdateRequest:    &hook.Hook[*ViewUpdateEvent]{},
		onViewAfterUpdateRequest:     &hook.Hook[*ViewUpdateEvent]{},
		onViewBeforeDeleteRequest:    &hook.Hook[*ViewDeleteEvent]{},
		onViewAfterDeleteRequest:     &hook.Hook[*ViewDeleteEvent]{},
	}

	app.registerDefaultHooks()

	return app
}

// Bootstrap initializes the application
// (aka. create data dir, open db connections, load settings, etc.)
func (app *BaseApp) Bootstrap() error {
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

	// we don't check for an error because the db migrations may
	// have not been executed yet.
	app.RefreshSettings()

	return nil
}

// ResetBootstrapState takes care for releasing initialized app resources
// (eg. closing db connections).
func (app *BaseApp) ResetBootstrapState() error {
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			return err
		}
	}

	if app.logsDB != nil {
		if err := app.logsDB.Close(); err != nil {
			return err
		}
	}

	app.dao = nil
	app.logsDao = nil
	app.settings = nil

	return nil
}

// DB returns the default app database instance.
func (app *BaseApp) DB() *dbx.DB {
	return app.db
}

// Dao returns the default app Dao instance.
func (app *BaseApp) Dao() *daos.Dao {
	return app.dao
}

// LogsDB returns the app logs database instance.
func (app *BaseApp) LogsDB() *dbx.DB {
	return app.logsDB
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

// IsDebug returns whether the app is in debug mode
// (showing more detailed error logs, executed sql statements, etc.).
func (app *BaseApp) IsDebug() bool {
	return app.isDebug
}

// Settings returns the loaded app settings.
func (app *BaseApp) Settings() *Settings {
	return app.settings
}

// Cache returns the app internal cache store.
func (app *BaseApp) Cache() *store.Store[any] {
	return app.cache
}

// SubscriptionsBroker returns the app realtime subscriptions broker instance.
func (app *BaseApp) SubscriptionsBroker() *subscriptions.Broker {
	return app.subscriptionsBroker
}

// NewMailClient creates and returns a new SMTP or Sendmail client
// based on the current app settings.
func (app *BaseApp) NewMailClient() mailer.Mailer {
	if app.Settings().Smtp.Enabled {
		return mailer.NewSmtpClient(
			app.Settings().Smtp.Host,
			app.Settings().Smtp.Port,
			app.Settings().Smtp.Username,
			app.Settings().Smtp.Password,
			app.Settings().Smtp.Tls,
		)
	}

	return &mailer.Sendmail{}
}

// NewFilesystem creates a new local or S3 filesystem instance
// based on the current app settings.
//
// NB! Make sure to call `Close()` on the returned result
// after you are done working with it.
func (app *BaseApp) NewFilesystem() (*filesystem.System, error) {
	if app.settings.S3.Enabled {
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
	return filesystem.NewLocal(filepath.Join(app.DataDir(), "storage"))
}

// RefreshSettings reinitializes and reloads the stored application settings.
func (app *BaseApp) RefreshSettings() error {
	if app.settings == nil {
		app.settings = NewSettings()
	}

	encryptionKey := os.Getenv(app.EncryptionEnv())

	param, err := app.Dao().FindParamByKey(models.ParamAppSettings)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// no settings were previously stored
	if param == nil {
		return app.Dao().SaveParam(models.ParamAppSettings, app.settings, encryptionKey)
	}

	// load the settings from the stored param into the app ones
	// ---
	newSettings := NewSettings()

	// try first without decryption
	plainDecodeErr := json.Unmarshal(param.Value, newSettings)

	// failed, try to decrypt
	if plainDecodeErr != nil {
		// load without decrypt has failed and there is no encryption key to use for decrypt
		if encryptionKey == "" {
			return errors.New("Failed to load the stored app settings (missing or invalid encryption key).")
		}

		// decrypt
		decrypted, decryptErr := security.Decrypt(string(param.Value), encryptionKey)
		if decryptErr != nil {
			return decryptErr
		}

		// decode again
		decryptedDecodeErr := json.Unmarshal(decrypted, newSettings)
		if decryptedDecodeErr != nil {
			return decryptedDecodeErr
		}
	}

	if err := app.settings.Merge(newSettings); err != nil {
		return err
	}

	afterMergeRaw, err := json.Marshal(app.settings)
	if err != nil {
		return err
	}

	if
	// save because previously the settings weren't stored encrypted
	(plainDecodeErr == nil && encryptionKey != "") ||
		// or save because there are new fields after the merge
		!bytes.Equal(param.Value, afterMergeRaw) {
		saveErr := app.Dao().SaveParam(models.ParamAppSettings, app.settings, encryptionKey)
		if saveErr != nil {
			return saveErr
		}
	}

	return nil
}

// -------------------------------------------------------------------
// Serve event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnBeforeServe() *hook.Hook[*ServeEvent] {
	return app.onBeforeServe
}

// -------------------------------------------------------------------
// Dao event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnModelBeforeCreate() *hook.Hook[*ModelEvent] {
	return app.onModelBeforeCreate
}

func (app *BaseApp) OnModelAfterCreate() *hook.Hook[*ModelEvent] {
	return app.onModelAfterCreate
}

func (app *BaseApp) OnModelBeforeUpdate() *hook.Hook[*ModelEvent] {
	return app.onModelBeforeUpdate
}

func (app *BaseApp) OnModelAfterUpdate() *hook.Hook[*ModelEvent] {
	return app.onModelAfterUpdate
}

func (app *BaseApp) OnModelBeforeDelete() *hook.Hook[*ModelEvent] {
	return app.onModelBeforeDelete
}

func (app *BaseApp) OnModelAfterDelete() *hook.Hook[*ModelEvent] {
	return app.onModelAfterDelete
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

func (app *BaseApp) OnMailerBeforeRecordResetPasswordSend() *hook.Hook[*MailerRecordEvent] {
	return app.onMailerBeforeRecordResetPasswordSend
}

func (app *BaseApp) OnMailerAfterRecordResetPasswordSend() *hook.Hook[*MailerRecordEvent] {
	return app.onMailerAfterRecordResetPasswordSend
}

func (app *BaseApp) OnMailerBeforeRecordVerificationSend() *hook.Hook[*MailerRecordEvent] {
	return app.onMailerBeforeRecordVerificationSend
}

func (app *BaseApp) OnMailerAfterRecordVerificationSend() *hook.Hook[*MailerRecordEvent] {
	return app.onMailerAfterRecordVerificationSend
}

func (app *BaseApp) OnMailerBeforeRecordChangeEmailSend() *hook.Hook[*MailerRecordEvent] {
	return app.onMailerBeforeRecordChangeEmailSend
}

func (app *BaseApp) OnMailerAfterRecordChangeEmailSend() *hook.Hook[*MailerRecordEvent] {
	return app.onMailerAfterRecordChangeEmailSend
}

// -------------------------------------------------------------------
// Realtime API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRealtimeConnectRequest() *hook.Hook[*RealtimeConnectEvent] {
	return app.onRealtimeConnectRequest
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

func (app *BaseApp) OnFileDownloadRequest() *hook.Hook[*FileDownloadEvent] {
	return app.onFileDownloadRequest
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

// -------------------------------------------------------------------
// Auth Record API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRecordAuthRequest() *hook.Hook[*RecordAuthEvent] {
	return app.onRecordAuthRequest
}

func (app *BaseApp) OnRecordListExternalAuthsRequest() *hook.Hook[*RecordListExternalAuthsEvent] {
	return app.onRecordListExternalAuthsRequest
}

func (app *BaseApp) OnRecordBeforeUnlinkExternalAuthRequest() *hook.Hook[*RecordUnlinkExternalAuthEvent] {
	return app.onRecordBeforeUnlinkExternalAuthRequest
}

func (app *BaseApp) OnRecordAfterUnlinkExternalAuthRequest() *hook.Hook[*RecordUnlinkExternalAuthEvent] {
	return app.onRecordAfterUnlinkExternalAuthRequest
}

// -------------------------------------------------------------------
// Record API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRecordsListRequest() *hook.Hook[*RecordsListEvent] {
	return app.onRecordsListRequest
}

func (app *BaseApp) OnRecordViewRequest() *hook.Hook[*RecordViewEvent] {
	return app.onRecordViewRequest
}

func (app *BaseApp) OnRecordBeforeCreateRequest() *hook.Hook[*RecordCreateEvent] {
	return app.onRecordBeforeCreateRequest
}

func (app *BaseApp) OnRecordAfterCreateRequest() *hook.Hook[*RecordCreateEvent] {
	return app.onRecordAfterCreateRequest
}

func (app *BaseApp) OnRecordBeforeUpdateRequest() *hook.Hook[*RecordUpdateEvent] {
	return app.onRecordBeforeUpdateRequest
}

func (app *BaseApp) OnRecordAfterUpdateRequest() *hook.Hook[*RecordUpdateEvent] {
	return app.onRecordAfterUpdateRequest
}

func (app *BaseApp) OnRecordBeforeDeleteRequest() *hook.Hook[*RecordDeleteEvent] {
	return app.onRecordBeforeDeleteRequest
}

func (app *BaseApp) OnRecordAfterDeleteRequest() *hook.Hook[*RecordDeleteEvent] {
	return app.onRecordAfterDeleteRequest
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
// View API event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnRecordsFromViewListRequest() *hook.Hook[*RecordsFromViewListEvent] {
	return app.onRecordsFromViewListRequest
}

func (app *BaseApp) OnViewListRequest() *hook.Hook[*ViewListEvent] {
	return app.onViewListRequest
}

func (app *BaseApp) OnViewViewRequest() *hook.Hook[*ViewViewEvent] {
	return app.onViewViewRequest
}

func (app *BaseApp) OnViewBeforeCreateRequest() *hook.Hook[*ViewCreateEvent] {
	return app.onViewBeforeCreateRequest
}

func (app *BaseApp) OnViewAfterCreateRequest() *hook.Hook[*ViewCreateEvent] {
	return app.onViewAfterCreateRequest
}

func (app *BaseApp) OnViewBeforeUpdateRequest() *hook.Hook[*ViewUpdateEvent] {
	return app.onViewBeforeUpdateRequest
}

func (app *BaseApp) OnViewAfterUpdateRequest() *hook.Hook[*ViewUpdateEvent] {
	return app.onViewAfterUpdateRequest
}

func (app *BaseApp) OnViewBeforeDeleteRequest() *hook.Hook[*ViewDeleteEvent] {
	return app.onViewBeforeDeleteRequest
}

func (app *BaseApp) OnViewAfterDeleteRequest() *hook.Hook[*ViewDeleteEvent] {
	return app.onViewAfterDeleteRequest
}

// -------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------

func (app *BaseApp) initLogsDB() error {
	var connectErr error
	app.logsDB, connectErr = connectDB(filepath.Join(app.DataDir(), "logs.db"))
	if connectErr != nil {
		return connectErr
	}

	app.logsDao = daos.New(app.logsDB)

	return nil
}

func (app *BaseApp) initDataDB() error {
	var connectErr error
	app.db, connectErr = connectDB(filepath.Join(app.DataDir(), "data.db"))
	if connectErr != nil {
		return connectErr
	}

	app.db.QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if app.IsDebug() {
			color.HiBlack("[%.2fms] %v\n", float64(t.Milliseconds()), sql)
		}
	}

	app.db.ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if app.IsDebug() {
			color.HiBlack("[%.2fms] %v\n", float64(t.Milliseconds()), sql)
		}
	}

	app.dao = app.createDaoWithHooks(app.db)

	return nil
}

func (app *BaseApp) createDaoWithHooks(db dbx.Builder) *daos.Dao {
	dao := daos.New(db)

	dao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		return app.OnModelBeforeCreate().Trigger(&ModelEvent{eventDao, m})
	}

	dao.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) {
		app.OnModelAfterCreate().Trigger(&ModelEvent{eventDao, m})
	}

	dao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		return app.OnModelBeforeUpdate().Trigger(&ModelEvent{eventDao, m})
	}

	dao.AfterUpdateFunc = func(eventDao *daos.Dao, m models.Model) {
		app.OnModelAfterUpdate().Trigger(&ModelEvent{eventDao, m})
	}

	dao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		return app.OnModelBeforeDelete().Trigger(&ModelEvent{eventDao, m})
	}

	dao.AfterDeleteFunc = func(eventDao *daos.Dao, m models.Model) {
		app.OnModelAfterDelete().Trigger(&ModelEvent{eventDao, m})
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
			return errors.New("Failed to delete the files at " + prefix)
		}

		return nil
	}

	// delete storage files from deleted Collection, Records, etc.
	app.OnModelAfterDelete().Add(func(e *ModelEvent) error {
		if m, ok := e.Model.(models.FilesManager); ok && m.BaseFilesPath() != "" {
			if err := deletePrefix(m.BaseFilesPath()); err != nil && app.IsDebug() {
				// non critical error - only log for debug
				// (usually could happen because of S3 api limits)
				log.Println(err)
			}
		}

		return nil
	})
}
