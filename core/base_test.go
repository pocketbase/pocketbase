package core_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "unsafe"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/logger"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

func TestNewBaseApp(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "test_env",
		IsDev:         true,
	})

	if app.DataDir() != testDataDir {
		t.Fatalf("expected DataDir %q, got %q", testDataDir, app.DataDir())
	}

	if app.EncryptionEnv() != "test_env" {
		t.Fatalf("expected EncryptionEnv test_env, got %q", app.EncryptionEnv())
	}

	if !app.IsDev() {
		t.Fatalf("expected IsDev true, got %v", app.IsDev())
	}

	if app.Store() == nil {
		t.Fatal("expected Store to be set, got nil")
	}

	if app.Settings() == nil {
		t.Fatal("expected Settings to be set, got nil")
	}

	if app.SubscriptionsBroker() == nil {
		t.Fatal("expected SubscriptionsBroker to be set, got nil")
	}

	if app.Cron() == nil {
		t.Fatal("expected Cron to be set, got nil")
	}
}

func TestBaseAppBootstrap(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir: testDataDir,
	})
	defer app.ResetBootstrapState()

	if app.IsBootstrapped() {
		t.Fatal("Didn't expect the application to be bootstrapped.")
	}

	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	if !app.IsBootstrapped() {
		t.Fatal("Expected the application to be bootstrapped.")
	}

	if stat, err := os.Stat(testDataDir); err != nil || !stat.IsDir() {
		t.Fatal("Expected test data directory to be created.")
	}

	type nilCheck struct {
		name      string
		value     any
		expectNil bool
	}

	runNilChecks := func(checks []nilCheck) {
		for _, check := range checks {
			t.Run(check.name, func(t *testing.T) {
				isNil := check.value == nil
				if isNil != check.expectNil {
					t.Fatalf("Expected isNil %v, got %v", check.expectNil, isNil)
				}
			})
		}
	}

	nilChecksBeforeReset := []nilCheck{
		{"[before] concurrentDB", app.DB(), false},
		{"[before] nonconcurrentDB", app.NonconcurrentDB(), false},
		{"[before] auxConcurrentDB", app.AuxDB(), false},
		{"[before] auxNonconcurrentDB", app.AuxNonconcurrentDB(), false},
		{"[before] settings", app.Settings(), false},
		{"[before] logger", app.Logger(), false},
		{"[before] cached collections", app.Store().Get(core.StoreKeyCachedCollections), false},
	}

	runNilChecks(nilChecksBeforeReset)

	// reset
	if err := app.ResetBootstrapState(); err != nil {
		t.Fatal(err)
	}

	nilChecksAfterReset := []nilCheck{
		{"[after] concurrentDB", app.DB(), true},
		{"[after] nonconcurrentDB", app.NonconcurrentDB(), true},
		{"[after] auxConcurrentDB", app.AuxDB(), true},
		{"[after] auxNonconcurrentDB", app.AuxNonconcurrentDB(), true},
		{"[after] settings", app.Settings(), false},
		{"[after] logger", app.Logger(), false},
		{"[after] cached collections", app.Store().Get(core.StoreKeyCachedCollections), false},
	}

	runNilChecks(nilChecksAfterReset)
}

func TestNewBaseAppIsTransactional(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir: testDataDir,
	})
	defer app.ResetBootstrapState()

	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	if app.IsTransactional() {
		t.Fatalf("Didn't expect the app to be transactional")
	}

	app.RunInTransaction(func(txApp core.App) error {
		if !txApp.IsTransactional() {
			t.Fatalf("Expected the app to be transactional")
		}

		return nil
	})
}

func TestBaseAppNewMailClient(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "pb_test_env",
	})
	defer app.ResetBootstrapState()

	client1 := app.NewMailClient()
	m1, ok := client1.(*mailer.Sendmail)
	if !ok {
		t.Fatalf("Expected mailer.Sendmail instance, got %v", m1)
	}
	if m1.OnSend() == nil || m1.OnSend().Length() == 0 {
		t.Fatal("Expected OnSend hook to be registered")
	}

	app.Settings().SMTP.Enabled = true

	client2 := app.NewMailClient()
	m2, ok := client2.(*mailer.SMTPClient)
	if !ok {
		t.Fatalf("Expected mailer.SMTPClient instance, got %v", m2)
	}
	if m2.OnSend() == nil || m2.OnSend().Length() == 0 {
		t.Fatal("Expected OnSend hook to be registered")
	}
}

func TestBaseAppNewFilesystem(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir: testDataDir,
	})
	defer app.ResetBootstrapState()

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

	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir: testDataDir,
	})
	defer app.ResetBootstrapState()

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

func TestBaseAppLoggerWrites(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// reset
	if err := app.DeleteOldLogs(time.Now()); err != nil {
		t.Fatal(err)
	}

	const logsThreshold = 200

	totalLogs := func(app core.App, t *testing.T) int {
		var total int

		err := app.LogQuery().Select("count(*)").Row(&total)
		if err != nil {
			t.Fatalf("Failed to fetch total logs: %v", err)
		}

		return total
	}

	t.Run("disabled logs retention", func(t *testing.T) {
		app.Settings().Logs.MaxDays = 0

		for i := 0; i < logsThreshold+1; i++ {
			app.Logger().Error("test")
		}

		if total := totalLogs(app, t); total != 0 {
			t.Fatalf("Expected no logs, got %d", total)
		}
	})

	t.Run("test batch logs writes", func(t *testing.T) {
		app.Settings().Logs.MaxDays = 1

		for i := 0; i < logsThreshold-1; i++ {
			app.Logger().Error("test")
		}

		if total := totalLogs(app, t); total != 0 {
			t.Fatalf("Expected no logs, got %d", total)
		}

		// should trigger batch write
		app.Logger().Error("test")

		// should be added for the next batch write
		app.Logger().Error("test")

		if total := totalLogs(app, t); total != logsThreshold {
			t.Fatalf("Expected %d logs, got %d", logsThreshold, total)
		}

		// wait for ~3 secs to check the timer trigger
		time.Sleep(3200 * time.Millisecond)
		if total := totalLogs(app, t); total != logsThreshold+1 {
			t.Fatalf("Expected %d logs, got %d", logsThreshold+1, total)
		}
	})
}

func TestBaseAppRefreshSettingsLoggerMinLevelEnabled(t *testing.T) {
	scenarios := []struct {
		name  string
		isDev bool
		level int
		// level->enabled map
		expectations map[int]bool
	}{
		{
			"dev mode",
			true,
			4,
			map[int]bool{
				3: true,
				4: true,
				5: true,
			},
		},
		{
			"nondev mode",
			false,
			4,
			map[int]bool{
				3: false,
				4: true,
				5: true,
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			const testDataDir = "./pb_base_app_test_data_dir/"
			defer os.RemoveAll(testDataDir)

			app := core.NewBaseApp(core.BaseAppConfig{
				DataDir: testDataDir,
				IsDev:   s.isDev,
			})
			defer app.ResetBootstrapState()

			if err := app.Bootstrap(); err != nil {
				t.Fatal(err)
			}

			// silence query logs
			app.DB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {}
			app.DB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {}
			app.NonconcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {}
			app.NonconcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {}

			handler, ok := app.Logger().Handler().(*logger.BatchHandler)
			if !ok {
				t.Fatalf("Expected BatchHandler, got %v", app.Logger().Handler())
			}

			app.Settings().Logs.MinLevel = s.level

			if err := app.Save(app.Settings()); err != nil {
				t.Fatalf("Failed to save settings: %v", err)
			}

			for level, enabled := range s.expectations {
				if v := handler.Enabled(context.Background(), slog.Level(level)); v != enabled {
					t.Fatalf("Expected level %d Enabled() to be %v, got %v", level, enabled, v)
				}
			}
		})
	}
}
