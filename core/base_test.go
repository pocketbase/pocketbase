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

func TestBaseAppNewBackupsFilesystem(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := NewBaseApp(&BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "pb_test_env",
		IsDebug:       false,
	})

	// local
	local, localErr := app.NewBackupsFilesystem()
	if localErr != nil {
		t.Fatal(localErr)
	}
	if local == nil {
		t.Fatal("Expected local backups filesystem instance, got nil")
	}

	// misconfigured s3
	app.Settings().Backups.S3.Enabled = true
	s3, s3Err := app.NewBackupsFilesystem()
	if s3Err == nil {
		t.Fatal("Expected S3 error, got nil")
	}
	if s3 != nil {
		t.Fatalf("Expected nil s3 backups filesystem, got %v", s3)
	}
}
