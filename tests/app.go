// Package tests provides common helpers and mocks used in PocketBase application tests.
package tests

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/migrations/logs"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/migrate"
)

// TestApp is a wrapper app instance used for testing.
type TestApp struct {
	*core.BaseApp

	mux sync.Mutex

	// EventCalls defines a map to inspect which app events
	// (and how many times) were triggered.
	//
	// The following events are not counted because they execute always:
	// - OnBeforeBootstrap
	// - OnAfterBootstrap
	// - OnBeforeServe
	EventCalls map[string]int

	TestMailer *TestMailer
}

// Cleanup resets the test application state and removes the test
// app's dataDir from the filesystem.
//
// After this call, the app instance shouldn't be used anymore.
func (t *TestApp) Cleanup() {
	t.OnTerminate().Trigger(&core.TerminateEvent{App: t}, func(e *core.TerminateEvent) error {
		t.TestMailer.Reset()
		t.ResetEventCalls()
		t.ResetBootstrapState()

		return nil
	})

	if t.DataDir() != "" {
		os.RemoveAll(t.DataDir())
	}
}

// NewMailClient initializes (if not already) a test app mail client.
func (t *TestApp) NewMailClient() mailer.Mailer {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.TestMailer == nil {
		t.TestMailer = &TestMailer{}
	}

	return t.TestMailer
}

// ResetEventCalls resets the EventCalls counter.
func (t *TestApp) ResetEventCalls() {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.EventCalls = make(map[string]int)
}

func (t *TestApp) registerEventCall(name string) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.EventCalls == nil {
		t.EventCalls = make(map[string]int)
	}

	t.EventCalls[name]++

	return nil
}

// NewTestApp creates and initializes a test application instance.
//
// It is the caller's responsibility to call `app.Cleanup()`
// when the app is no longer needed.
func NewTestApp(optTestDataDir ...string) (*TestApp, error) {
	var testDataDir string
	if len(optTestDataDir) == 0 || optTestDataDir[0] == "" {
		// fallback to the default test data directory
		_, currentFile, _, _ := runtime.Caller(0)
		testDataDir = filepath.Join(path.Dir(currentFile), "data")
	} else {
		testDataDir = optTestDataDir[0]
	}

	tempDir, err := TempDirClone(testDataDir)
	if err != nil {
		return nil, err
	}

	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir:       tempDir,
		EncryptionEnv: "pb_test_env",
	})

	// load data dir and db connections
	if err := app.Bootstrap(); err != nil {
		return nil, err
	}

	// ensure that the Dao and DB configurations are properly loaded
	if _, err := app.Dao().DB().NewQuery("Select 1").Execute(); err != nil {
		return nil, err
	}
	if _, err := app.LogsDao().DB().NewQuery("Select 1").Execute(); err != nil {
		return nil, err
	}

	// apply any missing migrations
	if err := runMigrations(app); err != nil {
		return nil, err
	}

	// force disable request logs because the logs db call execute in a separate
	// go routine and it is possible to panic due to earlier api test completion.
	app.Settings().Logs.MaxDays = 0

	t := &TestApp{
		BaseApp:    app,
		EventCalls: make(map[string]int),
		TestMailer: &TestMailer{},
	}

	t.OnBeforeApiError().Add(func(e *core.ApiErrorEvent) error {
		return t.registerEventCall("OnBeforeApiError")
	})

	t.OnAfterApiError().Add(func(e *core.ApiErrorEvent) error {
		return t.registerEventCall("OnAfterApiError")
	})

	t.OnModelBeforeCreate().Add(func(e *core.ModelEvent) error {
		return t.registerEventCall("OnModelBeforeCreate")
	})

	t.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		return t.registerEventCall("OnModelAfterCreate")
	})

	t.OnModelBeforeUpdate().Add(func(e *core.ModelEvent) error {
		return t.registerEventCall("OnModelBeforeUpdate")
	})

	t.OnModelAfterUpdate().Add(func(e *core.ModelEvent) error {
		return t.registerEventCall("OnModelAfterUpdate")
	})

	t.OnModelBeforeDelete().Add(func(e *core.ModelEvent) error {
		return t.registerEventCall("OnModelBeforeDelete")
	})

	t.OnModelAfterDelete().Add(func(e *core.ModelEvent) error {
		return t.registerEventCall("OnModelAfterDelete")
	})

	t.OnRecordsListRequest().Add(func(e *core.RecordsListEvent) error {
		return t.registerEventCall("OnRecordsListRequest")
	})

	t.OnRecordViewRequest().Add(func(e *core.RecordViewEvent) error {
		return t.registerEventCall("OnRecordViewRequest")
	})

	t.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return t.registerEventCall("OnRecordBeforeCreateRequest")
	})

	t.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		return t.registerEventCall("OnRecordAfterCreateRequest")
	})

	t.OnRecordBeforeUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		return t.registerEventCall("OnRecordBeforeUpdateRequest")
	})

	t.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		return t.registerEventCall("OnRecordAfterUpdateRequest")
	})

	t.OnRecordBeforeDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		return t.registerEventCall("OnRecordBeforeDeleteRequest")
	})

	t.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		return t.registerEventCall("OnRecordAfterDeleteRequest")
	})

	t.OnRecordAuthRequest().Add(func(e *core.RecordAuthEvent) error {
		return t.registerEventCall("OnRecordAuthRequest")
	})

	t.OnRecordBeforeAuthWithPasswordRequest().Add(func(e *core.RecordAuthWithPasswordEvent) error {
		return t.registerEventCall("OnRecordBeforeAuthWithPasswordRequest")
	})

	t.OnRecordAfterAuthWithPasswordRequest().Add(func(e *core.RecordAuthWithPasswordEvent) error {
		return t.registerEventCall("OnRecordAfterAuthWithPasswordRequest")
	})

	t.OnRecordBeforeAuthWithOAuth2Request().Add(func(e *core.RecordAuthWithOAuth2Event) error {
		return t.registerEventCall("OnRecordBeforeAuthWithOAuth2Request")
	})

	t.OnRecordAfterAuthWithOAuth2Request().Add(func(e *core.RecordAuthWithOAuth2Event) error {
		return t.registerEventCall("OnRecordAfterAuthWithOAuth2Request")
	})

	t.OnRecordBeforeAuthRefreshRequest().Add(func(e *core.RecordAuthRefreshEvent) error {
		return t.registerEventCall("OnRecordBeforeAuthRefreshRequest")
	})

	t.OnRecordAfterAuthRefreshRequest().Add(func(e *core.RecordAuthRefreshEvent) error {
		return t.registerEventCall("OnRecordAfterAuthRefreshRequest")
	})

	t.OnRecordBeforeRequestPasswordResetRequest().Add(func(e *core.RecordRequestPasswordResetEvent) error {
		return t.registerEventCall("OnRecordBeforeRequestPasswordResetRequest")
	})

	t.OnRecordAfterRequestPasswordResetRequest().Add(func(e *core.RecordRequestPasswordResetEvent) error {
		return t.registerEventCall("OnRecordAfterRequestPasswordResetRequest")
	})

	t.OnRecordBeforeConfirmPasswordResetRequest().Add(func(e *core.RecordConfirmPasswordResetEvent) error {
		return t.registerEventCall("OnRecordBeforeConfirmPasswordResetRequest")
	})

	t.OnRecordAfterConfirmPasswordResetRequest().Add(func(e *core.RecordConfirmPasswordResetEvent) error {
		return t.registerEventCall("OnRecordAfterConfirmPasswordResetRequest")
	})

	t.OnRecordBeforeRequestVerificationRequest().Add(func(e *core.RecordRequestVerificationEvent) error {
		return t.registerEventCall("OnRecordBeforeRequestVerificationRequest")
	})

	t.OnRecordAfterRequestVerificationRequest().Add(func(e *core.RecordRequestVerificationEvent) error {
		return t.registerEventCall("OnRecordAfterRequestVerificationRequest")
	})

	t.OnRecordBeforeConfirmVerificationRequest().Add(func(e *core.RecordConfirmVerificationEvent) error {
		return t.registerEventCall("OnRecordBeforeConfirmVerificationRequest")
	})

	t.OnRecordAfterConfirmVerificationRequest().Add(func(e *core.RecordConfirmVerificationEvent) error {
		return t.registerEventCall("OnRecordAfterConfirmVerificationRequest")
	})

	t.OnRecordBeforeRequestEmailChangeRequest().Add(func(e *core.RecordRequestEmailChangeEvent) error {
		return t.registerEventCall("OnRecordBeforeRequestEmailChangeRequest")
	})

	t.OnRecordAfterRequestEmailChangeRequest().Add(func(e *core.RecordRequestEmailChangeEvent) error {
		return t.registerEventCall("OnRecordAfterRequestEmailChangeRequest")
	})

	t.OnRecordBeforeConfirmEmailChangeRequest().Add(func(e *core.RecordConfirmEmailChangeEvent) error {
		return t.registerEventCall("OnRecordBeforeConfirmEmailChangeRequest")
	})

	t.OnRecordAfterConfirmEmailChangeRequest().Add(func(e *core.RecordConfirmEmailChangeEvent) error {
		return t.registerEventCall("OnRecordAfterConfirmEmailChangeRequest")
	})

	t.OnRecordListExternalAuthsRequest().Add(func(e *core.RecordListExternalAuthsEvent) error {
		return t.registerEventCall("OnRecordListExternalAuthsRequest")
	})

	t.OnRecordBeforeUnlinkExternalAuthRequest().Add(func(e *core.RecordUnlinkExternalAuthEvent) error {
		return t.registerEventCall("OnRecordBeforeUnlinkExternalAuthRequest")
	})

	t.OnRecordAfterUnlinkExternalAuthRequest().Add(func(e *core.RecordUnlinkExternalAuthEvent) error {
		return t.registerEventCall("OnRecordAfterUnlinkExternalAuthRequest")
	})

	t.OnMailerBeforeAdminResetPasswordSend().Add(func(e *core.MailerAdminEvent) error {
		return t.registerEventCall("OnMailerBeforeAdminResetPasswordSend")
	})

	t.OnMailerAfterAdminResetPasswordSend().Add(func(e *core.MailerAdminEvent) error {
		return t.registerEventCall("OnMailerAfterAdminResetPasswordSend")
	})

	t.OnMailerBeforeRecordResetPasswordSend().Add(func(e *core.MailerRecordEvent) error {
		return t.registerEventCall("OnMailerBeforeRecordResetPasswordSend")
	})

	t.OnMailerAfterRecordResetPasswordSend().Add(func(e *core.MailerRecordEvent) error {
		return t.registerEventCall("OnMailerAfterRecordResetPasswordSend")
	})

	t.OnMailerBeforeRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		return t.registerEventCall("OnMailerBeforeRecordVerificationSend")
	})

	t.OnMailerAfterRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		return t.registerEventCall("OnMailerAfterRecordVerificationSend")
	})

	t.OnMailerBeforeRecordChangeEmailSend().Add(func(e *core.MailerRecordEvent) error {
		return t.registerEventCall("OnMailerBeforeRecordChangeEmailSend")
	})

	t.OnMailerAfterRecordChangeEmailSend().Add(func(e *core.MailerRecordEvent) error {
		return t.registerEventCall("OnMailerAfterRecordChangeEmailSend")
	})

	t.OnRealtimeConnectRequest().Add(func(e *core.RealtimeConnectEvent) error {
		return t.registerEventCall("OnRealtimeConnectRequest")
	})

	t.OnRealtimeDisconnectRequest().Add(func(e *core.RealtimeDisconnectEvent) error {
		return t.registerEventCall("OnRealtimeDisconnectRequest")
	})

	t.OnRealtimeBeforeMessageSend().Add(func(e *core.RealtimeMessageEvent) error {
		return t.registerEventCall("OnRealtimeBeforeMessageSend")
	})

	t.OnRealtimeAfterMessageSend().Add(func(e *core.RealtimeMessageEvent) error {
		return t.registerEventCall("OnRealtimeAfterMessageSend")
	})

	t.OnRealtimeBeforeSubscribeRequest().Add(func(e *core.RealtimeSubscribeEvent) error {
		return t.registerEventCall("OnRealtimeBeforeSubscribeRequest")
	})

	t.OnRealtimeAfterSubscribeRequest().Add(func(e *core.RealtimeSubscribeEvent) error {
		return t.registerEventCall("OnRealtimeAfterSubscribeRequest")
	})

	t.OnSettingsListRequest().Add(func(e *core.SettingsListEvent) error {
		return t.registerEventCall("OnSettingsListRequest")
	})

	t.OnSettingsBeforeUpdateRequest().Add(func(e *core.SettingsUpdateEvent) error {
		return t.registerEventCall("OnSettingsBeforeUpdateRequest")
	})

	t.OnSettingsAfterUpdateRequest().Add(func(e *core.SettingsUpdateEvent) error {
		return t.registerEventCall("OnSettingsAfterUpdateRequest")
	})

	t.OnCollectionsListRequest().Add(func(e *core.CollectionsListEvent) error {
		return t.registerEventCall("OnCollectionsListRequest")
	})

	t.OnCollectionViewRequest().Add(func(e *core.CollectionViewEvent) error {
		return t.registerEventCall("OnCollectionViewRequest")
	})

	t.OnCollectionBeforeCreateRequest().Add(func(e *core.CollectionCreateEvent) error {
		return t.registerEventCall("OnCollectionBeforeCreateRequest")
	})

	t.OnCollectionAfterCreateRequest().Add(func(e *core.CollectionCreateEvent) error {
		return t.registerEventCall("OnCollectionAfterCreateRequest")
	})

	t.OnCollectionBeforeUpdateRequest().Add(func(e *core.CollectionUpdateEvent) error {
		return t.registerEventCall("OnCollectionBeforeUpdateRequest")
	})

	t.OnCollectionAfterUpdateRequest().Add(func(e *core.CollectionUpdateEvent) error {
		return t.registerEventCall("OnCollectionAfterUpdateRequest")
	})

	t.OnCollectionBeforeDeleteRequest().Add(func(e *core.CollectionDeleteEvent) error {
		return t.registerEventCall("OnCollectionBeforeDeleteRequest")
	})

	t.OnCollectionAfterDeleteRequest().Add(func(e *core.CollectionDeleteEvent) error {
		return t.registerEventCall("OnCollectionAfterDeleteRequest")
	})

	t.OnCollectionsBeforeImportRequest().Add(func(e *core.CollectionsImportEvent) error {
		return t.registerEventCall("OnCollectionsBeforeImportRequest")
	})

	t.OnCollectionsAfterImportRequest().Add(func(e *core.CollectionsImportEvent) error {
		return t.registerEventCall("OnCollectionsAfterImportRequest")
	})

	t.OnAdminsListRequest().Add(func(e *core.AdminsListEvent) error {
		return t.registerEventCall("OnAdminsListRequest")
	})

	t.OnAdminViewRequest().Add(func(e *core.AdminViewEvent) error {
		return t.registerEventCall("OnAdminViewRequest")
	})

	t.OnAdminBeforeCreateRequest().Add(func(e *core.AdminCreateEvent) error {
		return t.registerEventCall("OnAdminBeforeCreateRequest")
	})

	t.OnAdminAfterCreateRequest().Add(func(e *core.AdminCreateEvent) error {
		return t.registerEventCall("OnAdminAfterCreateRequest")
	})

	t.OnAdminBeforeUpdateRequest().Add(func(e *core.AdminUpdateEvent) error {
		return t.registerEventCall("OnAdminBeforeUpdateRequest")
	})

	t.OnAdminAfterUpdateRequest().Add(func(e *core.AdminUpdateEvent) error {
		return t.registerEventCall("OnAdminAfterUpdateRequest")
	})

	t.OnAdminBeforeDeleteRequest().Add(func(e *core.AdminDeleteEvent) error {
		return t.registerEventCall("OnAdminBeforeDeleteRequest")
	})

	t.OnAdminAfterDeleteRequest().Add(func(e *core.AdminDeleteEvent) error {
		return t.registerEventCall("OnAdminAfterDeleteRequest")
	})

	t.OnAdminAuthRequest().Add(func(e *core.AdminAuthEvent) error {
		return t.registerEventCall("OnAdminAuthRequest")
	})

	t.OnAdminBeforeAuthWithPasswordRequest().Add(func(e *core.AdminAuthWithPasswordEvent) error {
		return t.registerEventCall("OnAdminBeforeAuthWithPasswordRequest")
	})

	t.OnAdminAfterAuthWithPasswordRequest().Add(func(e *core.AdminAuthWithPasswordEvent) error {
		return t.registerEventCall("OnAdminAfterAuthWithPasswordRequest")
	})

	t.OnAdminBeforeAuthRefreshRequest().Add(func(e *core.AdminAuthRefreshEvent) error {
		return t.registerEventCall("OnAdminBeforeAuthRefreshRequest")
	})

	t.OnAdminAfterAuthRefreshRequest().Add(func(e *core.AdminAuthRefreshEvent) error {
		return t.registerEventCall("OnAdminAfterAuthRefreshRequest")
	})

	t.OnAdminBeforeRequestPasswordResetRequest().Add(func(e *core.AdminRequestPasswordResetEvent) error {
		return t.registerEventCall("OnAdminBeforeRequestPasswordResetRequest")
	})

	t.OnAdminAfterRequestPasswordResetRequest().Add(func(e *core.AdminRequestPasswordResetEvent) error {
		return t.registerEventCall("OnAdminAfterRequestPasswordResetRequest")
	})

	t.OnAdminBeforeConfirmPasswordResetRequest().Add(func(e *core.AdminConfirmPasswordResetEvent) error {
		return t.registerEventCall("OnAdminBeforeConfirmPasswordResetRequest")
	})

	t.OnAdminAfterConfirmPasswordResetRequest().Add(func(e *core.AdminConfirmPasswordResetEvent) error {
		return t.registerEventCall("OnAdminAfterConfirmPasswordResetRequest")
	})

	t.OnFileDownloadRequest().Add(func(e *core.FileDownloadEvent) error {
		return t.registerEventCall("OnFileDownloadRequest")
	})

	t.OnFileBeforeTokenRequest().Add(func(e *core.FileTokenEvent) error {
		return t.registerEventCall("OnFileBeforeTokenRequest")
	})

	t.OnFileAfterTokenRequest().Add(func(e *core.FileTokenEvent) error {
		return t.registerEventCall("OnFileAfterTokenRequest")
	})

	return t, nil
}

// TempDirClone creates a new temporary directory copy from the
// provided directory path.
//
// It is the caller's responsibility to call `os.RemoveAll(tempDir)`
// when the directory is no longer needed!
func TempDirClone(dirToClone string) (string, error) {
	tempDir, err := os.MkdirTemp("", "pb_test_*")
	if err != nil {
		return "", err
	}

	// copy everything from testDataDir to tempDir
	if err := copyDir(dirToClone, tempDir); err != nil {
		return "", err
	}

	return tempDir, nil
}

// -------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------

func copyDir(src string, dest string) error {
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	sourceDir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceDir.Close()

	items, err := sourceDir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, item := range items {
		fullSrcPath := filepath.Join(src, item.Name())
		fullDestPath := filepath.Join(dest, item.Name())

		var copyErr error
		if item.IsDir() {
			copyErr = copyDir(fullSrcPath, fullDestPath)
		} else {
			copyErr = copyFile(fullSrcPath, fullDestPath)
		}

		if copyErr != nil {
			return copyErr
		}
	}

	return nil
}

func copyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}

// @todo replace with app.RunMigrations on merge with the refactoring.
func runMigrations(app core.App) error {
	connections := []struct {
		db             *dbx.DB
		migrationsList migrate.MigrationsList
	}{
		{
			db:             app.DB(),
			migrationsList: migrations.AppMigrations,
		},
		{
			db:             app.LogsDB(),
			migrationsList: logs.LogsMigrations,
		},
	}

	for _, c := range connections {
		runner, err := migrate.NewRunner(c.db, c.migrationsList)
		if err != nil {
			return err
		}

		if _, err := runner.Up(); err != nil {
			return err
		}
	}

	return nil
}
