package core

import (
	"os"
	"testing"

	"github.com/pocketbase/pocketbase/tools/mailer"
)

func TestNewBaseApp(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := NewBaseApp(&BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "test_env",
		IsDebug:       true,
	})

	if app.dataDir != testDataDir {
		t.Fatalf("expected dataDir %q, got %q", testDataDir, app.dataDir)
	}

	if app.encryptionEnv != "test_env" {
		t.Fatalf("expected encryptionEnv test_env, got %q", app.dataDir)
	}

	if !app.isDebug {
		t.Fatalf("expected isDebug true, got %v", app.isDebug)
	}

	if app.cache == nil {
		t.Fatal("expected cache to be set, got nil")
	}

	if app.settings == nil {
		t.Fatal("expected settings to be set, got nil")
	}

	if app.subscriptionsBroker == nil {
		t.Fatal("expected subscriptionsBroker to be set, got nil")
	}
}

func TestBaseAppBootstrap(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := NewBaseApp(&BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "pb_test_env",
		IsDebug:       false,
	})
	defer app.ResetBootstrapState()

	if app.IsBootstrapped() {
		t.Fatal("Didn't expect the application to be bootstrapped.")
	}

	// bootstrap
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	if !app.IsBootstrapped() {
		t.Fatal("Expected the application to be bootstrapped.")
	}

	if stat, err := os.Stat(testDataDir); err != nil || !stat.IsDir() {
		t.Fatal("Expected test data directory to be created.")
	}

	if app.dao == nil {
		t.Fatal("Expected app.dao to be initialized, got nil.")
	}

	if app.dao.BeforeCreateFunc == nil {
		t.Fatal("Expected app.dao.BeforeCreateFunc to be set, got nil.")
	}

	if app.dao.AfterCreateFunc == nil {
		t.Fatal("Expected app.dao.AfterCreateFunc to be set, got nil.")
	}

	if app.dao.BeforeUpdateFunc == nil {
		t.Fatal("Expected app.dao.BeforeUpdateFunc to be set, got nil.")
	}

	if app.dao.AfterUpdateFunc == nil {
		t.Fatal("Expected app.dao.AfterUpdateFunc to be set, got nil.")
	}

	if app.dao.BeforeDeleteFunc == nil {
		t.Fatal("Expected app.dao.BeforeDeleteFunc to be set, got nil.")
	}

	if app.dao.AfterDeleteFunc == nil {
		t.Fatal("Expected app.dao.AfterDeleteFunc to be set, got nil.")
	}

	if app.logsDao == nil {
		t.Fatal("Expected app.logsDao to be initialized, got nil.")
	}

	if app.settings == nil {
		t.Fatal("Expected app.settings to be initialized, got nil.")
	}

	// reset
	if err := app.ResetBootstrapState(); err != nil {
		t.Fatal(err)
	}

	if app.dao != nil {
		t.Fatalf("Expected app.dao to be nil, got %v.", app.dao)
	}

	if app.logsDao != nil {
		t.Fatalf("Expected app.logsDao to be nil, got %v.", app.logsDao)
	}

	if app.settings != nil {
		t.Fatalf("Expected app.settings to be nil, got %v.", app.settings)
	}
}

func TestBaseAppGetters(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := NewBaseApp(&BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "pb_test_env",
		IsDebug:       false,
	})
	defer app.ResetBootstrapState()

	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	if app.dao != app.Dao() {
		t.Fatalf("Expected app.Dao %v, got %v", app.Dao(), app.dao)
	}

	if app.dao.ConcurrentDB() != app.DB() {
		t.Fatalf("Expected app.DB %v, got %v", app.DB(), app.dao.ConcurrentDB())
	}

	if app.logsDao != app.LogsDao() {
		t.Fatalf("Expected app.LogsDao %v, got %v", app.LogsDao(), app.logsDao)
	}

	if app.logsDao.ConcurrentDB() != app.LogsDB() {
		t.Fatalf("Expected app.LogsDB %v, got %v", app.LogsDB(), app.logsDao.ConcurrentDB())
	}

	if app.dataDir != app.DataDir() {
		t.Fatalf("Expected app.DataDir %v, got %v", app.DataDir(), app.dataDir)
	}

	if app.encryptionEnv != app.EncryptionEnv() {
		t.Fatalf("Expected app.EncryptionEnv %v, got %v", app.EncryptionEnv(), app.encryptionEnv)
	}

	if app.isDebug != app.IsDebug() {
		t.Fatalf("Expected app.IsDebug %v, got %v", app.IsDebug(), app.isDebug)
	}

	if app.settings != app.Settings() {
		t.Fatalf("Expected app.Settings %v, got %v", app.Settings(), app.settings)
	}

	if app.cache != app.Cache() {
		t.Fatalf("Expected app.Cache %v, got %v", app.Cache(), app.cache)
	}

	if app.subscriptionsBroker != app.SubscriptionsBroker() {
		t.Fatalf("Expected app.SubscriptionsBroker %v, got %v", app.SubscriptionsBroker(), app.subscriptionsBroker)
	}

	if app.onBeforeServe != app.OnBeforeServe() || app.OnBeforeServe() == nil {
		t.Fatalf("Getter app.OnBeforeServe does not match or nil (%v vs %v)", app.OnBeforeServe(), app.onBeforeServe)
	}

	if app.onModelBeforeCreate != app.OnModelBeforeCreate() || app.OnModelBeforeCreate() == nil {
		t.Fatalf("Getter app.OnModelBeforeCreate does not match or nil (%v vs %v)", app.OnModelBeforeCreate(), app.onModelBeforeCreate)
	}

	if app.onModelAfterCreate != app.OnModelAfterCreate() || app.OnModelAfterCreate() == nil {
		t.Fatalf("Getter app.OnModelAfterCreate does not match or nil (%v vs %v)", app.OnModelAfterCreate(), app.onModelAfterCreate)
	}

	if app.onModelBeforeUpdate != app.OnModelBeforeUpdate() || app.OnModelBeforeUpdate() == nil {
		t.Fatalf("Getter app.OnModelBeforeUpdate does not match or nil (%v vs %v)", app.OnModelBeforeUpdate(), app.onModelBeforeUpdate)
	}

	if app.onModelAfterUpdate != app.OnModelAfterUpdate() || app.OnModelAfterUpdate() == nil {
		t.Fatalf("Getter app.OnModelAfterUpdate does not match or nil (%v vs %v)", app.OnModelAfterUpdate(), app.onModelAfterUpdate)
	}

	if app.onModelBeforeDelete != app.OnModelBeforeDelete() || app.OnModelBeforeDelete() == nil {
		t.Fatalf("Getter app.OnModelBeforeDelete does not match or nil (%v vs %v)", app.OnModelBeforeDelete(), app.onModelBeforeDelete)
	}

	if app.onModelAfterDelete != app.OnModelAfterDelete() || app.OnModelAfterDelete() == nil {
		t.Fatalf("Getter app.OnModelAfterDelete does not match or nil (%v vs %v)", app.OnModelAfterDelete(), app.onModelAfterDelete)
	}

	if app.onMailerBeforeAdminResetPasswordSend != app.OnMailerBeforeAdminResetPasswordSend() || app.OnMailerBeforeAdminResetPasswordSend() == nil {
		t.Fatalf("Getter app.OnMailerBeforeAdminResetPasswordSend does not match or nil (%v vs %v)", app.OnMailerBeforeAdminResetPasswordSend(), app.onMailerBeforeAdminResetPasswordSend)
	}

	if app.onMailerAfterAdminResetPasswordSend != app.OnMailerAfterAdminResetPasswordSend() || app.OnMailerAfterAdminResetPasswordSend() == nil {
		t.Fatalf("Getter app.OnMailerAfterAdminResetPasswordSend does not match or nil (%v vs %v)", app.OnMailerAfterAdminResetPasswordSend(), app.onMailerAfterAdminResetPasswordSend)
	}

	if app.onMailerBeforeRecordResetPasswordSend != app.OnMailerBeforeRecordResetPasswordSend() || app.OnMailerBeforeRecordResetPasswordSend() == nil {
		t.Fatalf("Getter app.OnMailerBeforeRecordResetPasswordSend does not match or nil (%v vs %v)", app.OnMailerBeforeRecordResetPasswordSend(), app.onMailerBeforeRecordResetPasswordSend)
	}

	if app.onMailerAfterRecordResetPasswordSend != app.OnMailerAfterRecordResetPasswordSend() || app.OnMailerAfterRecordResetPasswordSend() == nil {
		t.Fatalf("Getter app.OnMailerAfterRecordResetPasswordSend does not match or nil (%v vs %v)", app.OnMailerAfterRecordResetPasswordSend(), app.onMailerAfterRecordResetPasswordSend)
	}

	if app.onMailerBeforeRecordVerificationSend != app.OnMailerBeforeRecordVerificationSend() || app.OnMailerBeforeRecordVerificationSend() == nil {
		t.Fatalf("Getter app.OnMailerBeforeRecordVerificationSend does not match or nil (%v vs %v)", app.OnMailerBeforeRecordVerificationSend(), app.onMailerBeforeRecordVerificationSend)
	}

	if app.onMailerAfterRecordVerificationSend != app.OnMailerAfterRecordVerificationSend() || app.OnMailerAfterRecordVerificationSend() == nil {
		t.Fatalf("Getter app.OnMailerAfterRecordVerificationSend does not match or nil (%v vs %v)", app.OnMailerAfterRecordVerificationSend(), app.onMailerAfterRecordVerificationSend)
	}

	if app.onMailerBeforeRecordChangeEmailSend != app.OnMailerBeforeRecordChangeEmailSend() || app.OnMailerBeforeRecordChangeEmailSend() == nil {
		t.Fatalf("Getter app.OnMailerBeforeRecordChangeEmailSend does not match or nil (%v vs %v)", app.OnMailerBeforeRecordChangeEmailSend(), app.onMailerBeforeRecordChangeEmailSend)
	}

	if app.onMailerAfterRecordChangeEmailSend != app.OnMailerAfterRecordChangeEmailSend() || app.OnMailerAfterRecordChangeEmailSend() == nil {
		t.Fatalf("Getter app.OnMailerAfterRecordChangeEmailSend does not match or nil (%v vs %v)", app.OnMailerAfterRecordChangeEmailSend(), app.onMailerAfterRecordChangeEmailSend)
	}

	if app.onRealtimeConnectRequest != app.OnRealtimeConnectRequest() || app.OnRealtimeConnectRequest() == nil {
		t.Fatalf("Getter app.OnRealtimeConnectRequest does not match or nil (%v vs %v)", app.OnRealtimeConnectRequest(), app.onRealtimeConnectRequest)
	}

	if app.onRealtimeBeforeSubscribeRequest != app.OnRealtimeBeforeSubscribeRequest() || app.OnRealtimeBeforeSubscribeRequest() == nil {
		t.Fatalf("Getter app.OnRealtimeBeforeSubscribeRequest does not match or nil (%v vs %v)", app.OnRealtimeBeforeSubscribeRequest(), app.onRealtimeBeforeSubscribeRequest)
	}

	if app.onRealtimeAfterSubscribeRequest != app.OnRealtimeAfterSubscribeRequest() || app.OnRealtimeAfterSubscribeRequest() == nil {
		t.Fatalf("Getter app.OnRealtimeAfterSubscribeRequest does not match or nil (%v vs %v)", app.OnRealtimeAfterSubscribeRequest(), app.onRealtimeAfterSubscribeRequest)
	}

	if app.onSettingsListRequest != app.OnSettingsListRequest() || app.OnSettingsListRequest() == nil {
		t.Fatalf("Getter app.OnSettingsListRequest does not match or nil (%v vs %v)", app.OnSettingsListRequest(), app.onSettingsListRequest)
	}

	if app.onSettingsBeforeUpdateRequest != app.OnSettingsBeforeUpdateRequest() || app.OnSettingsBeforeUpdateRequest() == nil {
		t.Fatalf("Getter app.OnSettingsBeforeUpdateRequest does not match or nil (%v vs %v)", app.OnSettingsBeforeUpdateRequest(), app.onSettingsBeforeUpdateRequest)
	}

	if app.onSettingsAfterUpdateRequest != app.OnSettingsAfterUpdateRequest() || app.OnSettingsAfterUpdateRequest() == nil {
		t.Fatalf("Getter app.OnSettingsAfterUpdateRequest does not match or nil (%v vs %v)", app.OnSettingsAfterUpdateRequest(), app.onSettingsAfterUpdateRequest)
	}

	if app.onFileDownloadRequest != app.OnFileDownloadRequest() || app.OnFileDownloadRequest() == nil {
		t.Fatalf("Getter app.OnFileDownloadRequest does not match or nil (%v vs %v)", app.OnFileDownloadRequest(), app.onFileDownloadRequest)
	}

	if app.onAdminsListRequest != app.OnAdminsListRequest() || app.OnAdminsListRequest() == nil {
		t.Fatalf("Getter app.OnAdminsListRequest does not match or nil (%v vs %v)", app.OnAdminsListRequest(), app.onAdminsListRequest)
	}

	if app.onAdminViewRequest != app.OnAdminViewRequest() || app.OnAdminViewRequest() == nil {
		t.Fatalf("Getter app.OnAdminViewRequest does not match or nil (%v vs %v)", app.OnAdminViewRequest(), app.onAdminViewRequest)
	}

	if app.onAdminBeforeCreateRequest != app.OnAdminBeforeCreateRequest() || app.OnAdminBeforeCreateRequest() == nil {
		t.Fatalf("Getter app.OnAdminBeforeCreateRequest does not match or nil (%v vs %v)", app.OnAdminBeforeCreateRequest(), app.onAdminBeforeCreateRequest)
	}

	if app.onAdminAfterCreateRequest != app.OnAdminAfterCreateRequest() || app.OnAdminAfterCreateRequest() == nil {
		t.Fatalf("Getter app.OnAdminAfterCreateRequest does not match or nil (%v vs %v)", app.OnAdminAfterCreateRequest(), app.onAdminAfterCreateRequest)
	}

	if app.onAdminBeforeUpdateRequest != app.OnAdminBeforeUpdateRequest() || app.OnAdminBeforeUpdateRequest() == nil {
		t.Fatalf("Getter app.OnAdminBeforeUpdateRequest does not match or nil (%v vs %v)", app.OnAdminBeforeUpdateRequest(), app.onAdminBeforeUpdateRequest)
	}

	if app.onAdminAfterUpdateRequest != app.OnAdminAfterUpdateRequest() || app.OnAdminAfterUpdateRequest() == nil {
		t.Fatalf("Getter app.OnAdminAfterUpdateRequest does not match or nil (%v vs %v)", app.OnAdminAfterUpdateRequest(), app.onAdminAfterUpdateRequest)
	}

	if app.onAdminBeforeDeleteRequest != app.OnAdminBeforeDeleteRequest() || app.OnAdminBeforeDeleteRequest() == nil {
		t.Fatalf("Getter app.OnAdminBeforeDeleteRequest does not match or nil (%v vs %v)", app.OnAdminBeforeDeleteRequest(), app.onAdminBeforeDeleteRequest)
	}

	if app.onAdminAfterDeleteRequest != app.OnAdminAfterDeleteRequest() || app.OnAdminAfterDeleteRequest() == nil {
		t.Fatalf("Getter app.OnAdminAfterDeleteRequest does not match or nil (%v vs %v)", app.OnAdminAfterDeleteRequest(), app.onAdminAfterDeleteRequest)
	}

	if app.onAdminAuthRequest != app.OnAdminAuthRequest() || app.OnAdminAuthRequest() == nil {
		t.Fatalf("Getter app.OnAdminAuthRequest does not match or nil (%v vs %v)", app.OnAdminAuthRequest(), app.onAdminAuthRequest)
	}

	if app.onRecordsListRequest != app.OnRecordsListRequest() || app.OnRecordsListRequest() == nil {
		t.Fatalf("Getter app.OnRecordsListRequest does not match or nil (%v vs %v)", app.OnRecordsListRequest(), app.onRecordsListRequest)
	}

	if app.onRecordViewRequest != app.OnRecordViewRequest() || app.OnRecordViewRequest() == nil {
		t.Fatalf("Getter app.OnRecordViewRequest does not match or nil (%v vs %v)", app.OnRecordViewRequest(), app.onRecordViewRequest)
	}

	if app.onRecordBeforeCreateRequest != app.OnRecordBeforeCreateRequest() || app.OnRecordBeforeCreateRequest() == nil {
		t.Fatalf("Getter app.OnRecordBeforeCreateRequest does not match or nil (%v vs %v)", app.OnRecordBeforeCreateRequest(), app.onRecordBeforeCreateRequest)
	}

	if app.onRecordAfterCreateRequest != app.OnRecordAfterCreateRequest() || app.OnRecordAfterCreateRequest() == nil {
		t.Fatalf("Getter app.OnRecordAfterCreateRequest does not match or nil (%v vs %v)", app.OnRecordAfterCreateRequest(), app.onRecordAfterCreateRequest)
	}

	if app.onRecordBeforeUpdateRequest != app.OnRecordBeforeUpdateRequest() || app.OnRecordBeforeUpdateRequest() == nil {
		t.Fatalf("Getter app.OnRecordBeforeUpdateRequest does not match or nil (%v vs %v)", app.OnRecordBeforeUpdateRequest(), app.onRecordBeforeUpdateRequest)
	}

	if app.onRecordAfterUpdateRequest != app.OnRecordAfterUpdateRequest() || app.OnRecordAfterUpdateRequest() == nil {
		t.Fatalf("Getter app.OnRecordAfterUpdateRequest does not match or nil (%v vs %v)", app.OnRecordAfterUpdateRequest(), app.onRecordAfterUpdateRequest)
	}

	if app.onRecordBeforeDeleteRequest != app.OnRecordBeforeDeleteRequest() || app.OnRecordBeforeDeleteRequest() == nil {
		t.Fatalf("Getter app.OnRecordBeforeDeleteRequest does not match or nil (%v vs %v)", app.OnRecordBeforeDeleteRequest(), app.onRecordBeforeDeleteRequest)
	}

	if app.onRecordAfterDeleteRequest != app.OnRecordAfterDeleteRequest() || app.OnRecordAfterDeleteRequest() == nil {
		t.Fatalf("Getter app.OnRecordAfterDeleteRequest does not match or nil (%v vs %v)", app.OnRecordAfterDeleteRequest(), app.onRecordAfterDeleteRequest)
	}

	if app.onRecordAuthRequest != app.OnRecordAuthRequest() || app.OnRecordAuthRequest() == nil {
		t.Fatalf("Getter app.OnRecordAuthRequest does not match or nil (%v vs %v)", app.OnRecordAuthRequest(), app.onRecordAuthRequest)
	}

	if app.onRecordListExternalAuthsRequest != app.OnRecordListExternalAuthsRequest() || app.OnRecordListExternalAuthsRequest() == nil {
		t.Fatalf("Getter app.OnRecordListExternalAuthsRequest does not match or nil (%v vs %v)", app.OnRecordListExternalAuthsRequest(), app.onRecordListExternalAuthsRequest)
	}

	if app.onRecordBeforeUnlinkExternalAuthRequest != app.OnRecordBeforeUnlinkExternalAuthRequest() || app.OnRecordBeforeUnlinkExternalAuthRequest() == nil {
		t.Fatalf("Getter app.OnRecordBeforeUnlinkExternalAuthRequest does not match or nil (%v vs %v)", app.OnRecordBeforeUnlinkExternalAuthRequest(), app.onRecordBeforeUnlinkExternalAuthRequest)
	}

	if app.onRecordAfterUnlinkExternalAuthRequest != app.OnRecordAfterUnlinkExternalAuthRequest() || app.OnRecordAfterUnlinkExternalAuthRequest() == nil {
		t.Fatalf("Getter app.OnRecordAfterUnlinkExternalAuthRequest does not match or nil (%v vs %v)", app.OnRecordAfterUnlinkExternalAuthRequest(), app.onRecordAfterUnlinkExternalAuthRequest)
	}

	if app.onRecordsListRequest != app.OnRecordsListRequest() || app.OnRecordsListRequest() == nil {
		t.Fatalf("Getter app.OnRecordsListRequest does not match or nil (%v vs %v)", app.OnRecordsListRequest(), app.onRecordsListRequest)
	}

	if app.onRecordViewRequest != app.OnRecordViewRequest() || app.OnRecordViewRequest() == nil {
		t.Fatalf("Getter app.OnRecordViewRequest does not match or nil (%v vs %v)", app.OnRecordViewRequest(), app.onRecordViewRequest)
	}

	if app.onRecordBeforeCreateRequest != app.OnRecordBeforeCreateRequest() || app.OnRecordBeforeCreateRequest() == nil {
		t.Fatalf("Getter app.OnRecordBeforeCreateRequest does not match or nil (%v vs %v)", app.OnRecordBeforeCreateRequest(), app.onRecordBeforeCreateRequest)
	}

	if app.onRecordAfterCreateRequest != app.OnRecordAfterCreateRequest() || app.OnRecordAfterCreateRequest() == nil {
		t.Fatalf("Getter app.OnRecordAfterCreateRequest does not match or nil (%v vs %v)", app.OnRecordAfterCreateRequest(), app.onRecordAfterCreateRequest)
	}

	if app.onRecordBeforeUpdateRequest != app.OnRecordBeforeUpdateRequest() || app.OnRecordBeforeUpdateRequest() == nil {
		t.Fatalf("Getter app.OnRecordBeforeUpdateRequest does not match or nil (%v vs %v)", app.OnRecordBeforeUpdateRequest(), app.onRecordBeforeUpdateRequest)
	}

	if app.onRecordAfterUpdateRequest != app.OnRecordAfterUpdateRequest() || app.OnRecordAfterUpdateRequest() == nil {
		t.Fatalf("Getter app.OnRecordAfterUpdateRequest does not match or nil (%v vs %v)", app.OnRecordAfterUpdateRequest(), app.onRecordAfterUpdateRequest)
	}

	if app.onRecordBeforeDeleteRequest != app.OnRecordBeforeDeleteRequest() || app.OnRecordBeforeDeleteRequest() == nil {
		t.Fatalf("Getter app.OnRecordBeforeDeleteRequest does not match or nil (%v vs %v)", app.OnRecordBeforeDeleteRequest(), app.onRecordBeforeDeleteRequest)
	}

	if app.onRecordAfterDeleteRequest != app.OnRecordAfterDeleteRequest() || app.OnRecordAfterDeleteRequest() == nil {
		t.Fatalf("Getter app.OnRecordAfterDeleteRequest does not match or nil (%v vs %v)", app.OnRecordAfterDeleteRequest(), app.onRecordAfterDeleteRequest)
	}

	if app.onCollectionsListRequest != app.OnCollectionsListRequest() || app.OnCollectionsListRequest() == nil {
		t.Fatalf("Getter app.OnCollectionsListRequest does not match or nil (%v vs %v)", app.OnCollectionsListRequest(), app.onCollectionsListRequest)
	}

	if app.onCollectionViewRequest != app.OnCollectionViewRequest() || app.OnCollectionViewRequest() == nil {
		t.Fatalf("Getter app.OnCollectionViewRequest does not match or nil (%v vs %v)", app.OnCollectionViewRequest(), app.onCollectionViewRequest)
	}

	if app.onCollectionBeforeCreateRequest != app.OnCollectionBeforeCreateRequest() || app.OnCollectionBeforeCreateRequest() == nil {
		t.Fatalf("Getter app.OnCollectionBeforeCreateRequest does not match or nil (%v vs %v)", app.OnCollectionBeforeCreateRequest(), app.onCollectionBeforeCreateRequest)
	}

	if app.onCollectionAfterCreateRequest != app.OnCollectionAfterCreateRequest() || app.OnCollectionAfterCreateRequest() == nil {
		t.Fatalf("Getter app.OnCollectionAfterCreateRequest does not match or nil (%v vs %v)", app.OnCollectionAfterCreateRequest(), app.onCollectionAfterCreateRequest)
	}

	if app.onCollectionBeforeUpdateRequest != app.OnCollectionBeforeUpdateRequest() || app.OnCollectionBeforeUpdateRequest() == nil {
		t.Fatalf("Getter app.OnCollectionBeforeUpdateRequest does not match or nil (%v vs %v)", app.OnCollectionBeforeUpdateRequest(), app.onCollectionBeforeUpdateRequest)
	}

	if app.onCollectionAfterUpdateRequest != app.OnCollectionAfterUpdateRequest() || app.OnCollectionAfterUpdateRequest() == nil {
		t.Fatalf("Getter app.OnCollectionAfterUpdateRequest does not match or nil (%v vs %v)", app.OnCollectionAfterUpdateRequest(), app.onCollectionAfterUpdateRequest)
	}

	if app.onCollectionBeforeDeleteRequest != app.OnCollectionBeforeDeleteRequest() || app.OnCollectionBeforeDeleteRequest() == nil {
		t.Fatalf("Getter app.OnCollectionBeforeDeleteRequest does not match or nil (%v vs %v)", app.OnCollectionBeforeDeleteRequest(), app.onCollectionBeforeDeleteRequest)
	}

	if app.onCollectionAfterDeleteRequest != app.OnCollectionAfterDeleteRequest() || app.OnCollectionAfterDeleteRequest() == nil {
		t.Fatalf("Getter app.OnCollectionAfterDeleteRequest does not match or nil (%v vs %v)", app.OnCollectionAfterDeleteRequest(), app.onCollectionAfterDeleteRequest)
	}
}

func TestBaseAppNewMailClient(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := NewBaseApp(&BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "pb_test_env",
		IsDebug:       false,
	})

	client1 := app.NewMailClient()
	if val, ok := client1.(*mailer.Sendmail); !ok {
		t.Fatalf("Expected mailer.Sendmail instance, got %v", val)
	}

	app.Settings().Smtp.Enabled = true

	client2 := app.NewMailClient()
	if val, ok := client2.(*mailer.SmtpClient); !ok {
		t.Fatalf("Expected mailer.SmtpClient instance, got %v", val)
	}
}

func TestBaseAppNewFilesystem(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := NewBaseApp(&BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "pb_test_env",
		IsDebug:       false,
	})

	// local
	local, localErr := app.NewFilesystem()
	if localErr != nil {
		t.Fatal(localErr)
	}
	if local == nil {
		t.Fatal("Expected local filesystem instance, got nil")
	}

	// misconfigured s3
	app.Settings().S3.Enabled = true
	s3, s3Err := app.NewFilesystem()
	if s3Err == nil {
		t.Fatal("Expected S3 error, got nil")
	}
	if s3 != nil {
		t.Fatalf("Expected nil s3 filesystem, got %v", s3)
	}
}
