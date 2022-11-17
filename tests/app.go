// Package tests provides common helpers and mocks used in PocketBase application tests.
package tests

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

// TestApp is a wrapper app instance used for testing.
type TestApp struct {
	*core.BaseApp

	// EventCalls defines a map to inspect which app events
	// (and how many times) were triggered.
	EventCalls map[string]int
	TestMailer *TestMailer
}

// Cleanup resets the test application state and removes the test
// app's dataDir from the filesystem.
//
// After this call, the app instance shouldn't be used anymore.
func (t *TestApp) Cleanup() {
	t.ResetEventCalls()
	t.ResetBootstrapState()

	if t.DataDir() != "" {
		os.RemoveAll(t.DataDir())
	}
}

func (t *TestApp) NewMailClient() mailer.Mailer {
	t.TestMailer.Reset()
	return t.TestMailer
}

// ResetEventCalls resets the EventCalls counter.
func (t *TestApp) ResetEventCalls() {
	t.EventCalls = make(map[string]int)
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

	app := core.NewBaseApp(tempDir, "pb_test_env", false)

	// load data dir and db connections
	if err := app.Bootstrap(); err != nil {
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

	// no need to count since this is executed always
	// t.OnBeforeServe().Add(func(e *core.ServeEvent) error {
	// 	t.EventCalls["OnBeforeServe"]++
	// 	return nil
	// })

	t.OnModelBeforeCreate().Add(func(e *core.ModelEvent) error {
		t.EventCalls["OnModelBeforeCreate"]++
		return nil
	})

	t.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		t.EventCalls["OnModelAfterCreate"]++
		return nil
	})

	t.OnModelBeforeUpdate().Add(func(e *core.ModelEvent) error {
		t.EventCalls["OnModelBeforeUpdate"]++
		return nil
	})

	t.OnModelAfterUpdate().Add(func(e *core.ModelEvent) error {
		t.EventCalls["OnModelAfterUpdate"]++
		return nil
	})

	t.OnModelBeforeDelete().Add(func(e *core.ModelEvent) error {
		t.EventCalls["OnModelBeforeDelete"]++
		return nil
	})

	t.OnModelAfterDelete().Add(func(e *core.ModelEvent) error {
		t.EventCalls["OnModelAfterDelete"]++
		return nil
	})

	t.OnRecordsListRequest().Add(func(e *core.RecordsListEvent) error {
		t.EventCalls["OnRecordsListRequest"]++
		return nil
	})

	t.OnRecordViewRequest().Add(func(e *core.RecordViewEvent) error {
		t.EventCalls["OnRecordViewRequest"]++
		return nil
	})

	t.OnRecordBeforeCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		t.EventCalls["OnRecordBeforeCreateRequest"]++
		return nil
	})

	t.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		t.EventCalls["OnRecordAfterCreateRequest"]++
		return nil
	})

	t.OnRecordBeforeUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		t.EventCalls["OnRecordBeforeUpdateRequest"]++
		return nil
	})

	t.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		t.EventCalls["OnRecordAfterUpdateRequest"]++
		return nil
	})

	t.OnRecordBeforeDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		t.EventCalls["OnRecordBeforeDeleteRequest"]++
		return nil
	})

	t.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		t.EventCalls["OnRecordAfterDeleteRequest"]++
		return nil
	})

	t.OnRecordAuthRequest().Add(func(e *core.RecordAuthEvent) error {
		t.EventCalls["OnRecordAuthRequest"]++
		return nil
	})

	t.OnRecordListExternalAuthsRequest().Add(func(e *core.RecordListExternalAuthsEvent) error {
		t.EventCalls["OnRecordListExternalAuthsRequest"]++
		return nil
	})

	t.OnRecordBeforeUnlinkExternalAuthRequest().Add(func(e *core.RecordUnlinkExternalAuthEvent) error {
		t.EventCalls["OnRecordBeforeUnlinkExternalAuthRequest"]++
		return nil
	})

	t.OnRecordAfterUnlinkExternalAuthRequest().Add(func(e *core.RecordUnlinkExternalAuthEvent) error {
		t.EventCalls["OnRecordAfterUnlinkExternalAuthRequest"]++
		return nil
	})

	t.OnMailerBeforeAdminResetPasswordSend().Add(func(e *core.MailerAdminEvent) error {
		t.EventCalls["OnMailerBeforeAdminResetPasswordSend"]++
		return nil
	})

	t.OnMailerAfterAdminResetPasswordSend().Add(func(e *core.MailerAdminEvent) error {
		t.EventCalls["OnMailerAfterAdminResetPasswordSend"]++
		return nil
	})

	t.OnMailerBeforeRecordResetPasswordSend().Add(func(e *core.MailerRecordEvent) error {
		t.EventCalls["OnMailerBeforeRecordResetPasswordSend"]++
		return nil
	})

	t.OnMailerAfterRecordResetPasswordSend().Add(func(e *core.MailerRecordEvent) error {
		t.EventCalls["OnMailerAfterRecordResetPasswordSend"]++
		return nil
	})

	t.OnMailerBeforeRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		t.EventCalls["OnMailerBeforeRecordVerificationSend"]++
		return nil
	})

	t.OnMailerAfterRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		t.EventCalls["OnMailerAfterRecordVerificationSend"]++
		return nil
	})

	t.OnMailerBeforeRecordChangeEmailSend().Add(func(e *core.MailerRecordEvent) error {
		t.EventCalls["OnMailerBeforeRecordChangeEmailSend"]++
		return nil
	})

	t.OnMailerAfterRecordChangeEmailSend().Add(func(e *core.MailerRecordEvent) error {
		t.EventCalls["OnMailerAfterRecordChangeEmailSend"]++
		return nil
	})

	t.OnRealtimeConnectRequest().Add(func(e *core.RealtimeConnectEvent) error {
		t.EventCalls["OnRealtimeConnectRequest"]++
		return nil
	})

	t.OnRealtimeBeforeSubscribeRequest().Add(func(e *core.RealtimeSubscribeEvent) error {
		t.EventCalls["OnRealtimeBeforeSubscribeRequest"]++
		return nil
	})

	t.OnRealtimeAfterSubscribeRequest().Add(func(e *core.RealtimeSubscribeEvent) error {
		t.EventCalls["OnRealtimeAfterSubscribeRequest"]++
		return nil
	})

	t.OnSettingsListRequest().Add(func(e *core.SettingsListEvent) error {
		t.EventCalls["OnSettingsListRequest"]++
		return nil
	})

	t.OnSettingsBeforeUpdateRequest().Add(func(e *core.SettingsUpdateEvent) error {
		t.EventCalls["OnSettingsBeforeUpdateRequest"]++
		return nil
	})

	t.OnSettingsAfterUpdateRequest().Add(func(e *core.SettingsUpdateEvent) error {
		t.EventCalls["OnSettingsAfterUpdateRequest"]++
		return nil
	})

	t.OnCollectionsListRequest().Add(func(e *core.CollectionsListEvent) error {
		t.EventCalls["OnCollectionsListRequest"]++
		return nil
	})

	t.OnCollectionViewRequest().Add(func(e *core.CollectionViewEvent) error {
		t.EventCalls["OnCollectionViewRequest"]++
		return nil
	})

	t.OnCollectionBeforeCreateRequest().Add(func(e *core.CollectionCreateEvent) error {
		t.EventCalls["OnCollectionBeforeCreateRequest"]++
		return nil
	})

	t.OnCollectionAfterCreateRequest().Add(func(e *core.CollectionCreateEvent) error {
		t.EventCalls["OnCollectionAfterCreateRequest"]++
		return nil
	})

	t.OnCollectionBeforeUpdateRequest().Add(func(e *core.CollectionUpdateEvent) error {
		t.EventCalls["OnCollectionBeforeUpdateRequest"]++
		return nil
	})

	t.OnCollectionAfterUpdateRequest().Add(func(e *core.CollectionUpdateEvent) error {
		t.EventCalls["OnCollectionAfterUpdateRequest"]++
		return nil
	})

	t.OnCollectionBeforeDeleteRequest().Add(func(e *core.CollectionDeleteEvent) error {
		t.EventCalls["OnCollectionBeforeDeleteRequest"]++
		return nil
	})

	t.OnCollectionAfterDeleteRequest().Add(func(e *core.CollectionDeleteEvent) error {
		t.EventCalls["OnCollectionAfterDeleteRequest"]++
		return nil
	})

	t.OnCollectionsBeforeImportRequest().Add(func(e *core.CollectionsImportEvent) error {
		t.EventCalls["OnCollectionsBeforeImportRequest"]++
		return nil
	})

	t.OnCollectionsAfterImportRequest().Add(func(e *core.CollectionsImportEvent) error {
		t.EventCalls["OnCollectionsAfterImportRequest"]++
		return nil
	})

	t.OnAdminsListRequest().Add(func(e *core.AdminsListEvent) error {
		t.EventCalls["OnAdminsListRequest"]++
		return nil
	})

	t.OnAdminViewRequest().Add(func(e *core.AdminViewEvent) error {
		t.EventCalls["OnAdminViewRequest"]++
		return nil
	})

	t.OnAdminBeforeCreateRequest().Add(func(e *core.AdminCreateEvent) error {
		t.EventCalls["OnAdminBeforeCreateRequest"]++
		return nil
	})

	t.OnAdminAfterCreateRequest().Add(func(e *core.AdminCreateEvent) error {
		t.EventCalls["OnAdminAfterCreateRequest"]++
		return nil
	})

	t.OnAdminBeforeUpdateRequest().Add(func(e *core.AdminUpdateEvent) error {
		t.EventCalls["OnAdminBeforeUpdateRequest"]++
		return nil
	})

	t.OnAdminAfterUpdateRequest().Add(func(e *core.AdminUpdateEvent) error {
		t.EventCalls["OnAdminAfterUpdateRequest"]++
		return nil
	})

	t.OnAdminBeforeDeleteRequest().Add(func(e *core.AdminDeleteEvent) error {
		t.EventCalls["OnAdminBeforeDeleteRequest"]++
		return nil
	})

	t.OnAdminAfterDeleteRequest().Add(func(e *core.AdminDeleteEvent) error {
		t.EventCalls["OnAdminAfterDeleteRequest"]++
		return nil
	})

	t.OnAdminAuthRequest().Add(func(e *core.AdminAuthEvent) error {
		t.EventCalls["OnAdminAuthRequest"]++
		return nil
	})

	t.OnFileDownloadRequest().Add(func(e *core.FileDownloadEvent) error {
		t.EventCalls["OnFileDownloadRequest"]++
		return nil
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
