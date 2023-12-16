package core

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/migrations/logs"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/logger"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/migrate"
)

func TestNewBaseApp(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := NewBaseApp(BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "test_env",
		IsDev:         true,
	})

	if app.dataDir != testDataDir {
		t.Fatalf("expected dataDir %q, got %q", testDataDir, app.dataDir)
	}

	if app.encryptionEnv != "test_env" {
		t.Fatalf("expected encryptionEnv test_env, got %q", app.dataDir)
	}

	if !app.isDev {
		t.Fatalf("expected isDev true, got %v", app.isDev)
	}

	if app.store == nil {
		t.Fatal("expected store to be set, got nil")
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

	app := NewBaseApp(BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "pb_test_env",
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

	if app.logger == nil {
		t.Fatal("Expected app.logger to be initialized, got nil.")
	}

	if _, ok := app.logger.Handler().(*logger.BatchHandler); !ok {
		t.Fatal("Expected app.logger handler to be initialized.")
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
}

func TestBaseAppGetters(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := NewBaseApp(BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "pb_test_env",
		IsDev:         true,
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

	if app.isDev != app.IsDev() {
		t.Fatalf("Expected app.IsDev %v, got %v", app.IsDev(), app.isDev)
	}

	if app.settings != app.Settings() {
		t.Fatalf("Expected app.Settings %v, got %v", app.Settings(), app.settings)
	}

	if app.store != app.Store() {
		t.Fatalf("Expected app.Store %v, got %v", app.Store(), app.store)
	}

	if app.logger != app.Logger() {
		t.Fatalf("Expected app.Logger %v, got %v", app.Logger(), app.logger)
	}

	if app.subscriptionsBroker != app.SubscriptionsBroker() {
		t.Fatalf("Expected app.SubscriptionsBroker %v, got %v", app.SubscriptionsBroker(), app.subscriptionsBroker)
	}

	if app.onBeforeServe != app.OnBeforeServe() || app.OnBeforeServe() == nil {
		t.Fatalf("Getter app.OnBeforeServe does not match or nil (%v vs %v)", app.OnBeforeServe(), app.onBeforeServe)
	}
}

func TestBaseAppNewMailClient(t *testing.T) {
	app, cleanup, err := initTestBaseApp()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

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
	app, cleanup, err := initTestBaseApp()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

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
	app, cleanup, err := initTestBaseApp()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

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
	app, cleanup, err := initTestBaseApp()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	threshold := 200

	totalLogs := func(app App, t *testing.T) int {
		var total int

		err := app.LogsDao().LogQuery().Select("count(*)").Row(&total)
		if err != nil {
			t.Fatalf("Failed to fetch total logs: %v", err)
		}

		return total
	}

	// disabled logs retention
	{
		app.Settings().Logs.MaxDays = 0

		for i := 0; i < threshold+1; i++ {
			app.Logger().Error("test")
		}

		if total := totalLogs(app, t); total != 0 {
			t.Fatalf("Expected no logs, got %d", total)
		}
	}

	// test batch logs writes
	{
		app.Settings().Logs.MaxDays = 1

		for i := 0; i < threshold-1; i++ {
			app.Logger().Error("test")
		}

		if total := totalLogs(app, t); total != 0 {
			t.Fatalf("Expected no logs, got %d", total)
		}

		// should trigger batch write
		app.Logger().Error("test")

		// should be added for the next batch write
		app.Logger().Error("test")

		if total := totalLogs(app, t); total != threshold {
			t.Fatalf("Expected %d logs, got %d", threshold, total)
		}

		// wait for ~3 secs to check the timer trigger
		time.Sleep(3200 * time.Millisecond)
		if total := totalLogs(app, t); total != threshold+1 {
			t.Fatalf("Expected %d logs, got %d", threshold+1, total)
		}
	}
}

func TestBaseAppRefreshSettingsLoggerMinLevelEnabled(t *testing.T) {
	app, cleanup, err := initTestBaseApp()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	handler, ok := app.Logger().Handler().(*logger.BatchHandler)
	if !ok {
		t.Fatalf("Expected BatchHandler, got %v", app.Logger().Handler())
	}

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
			app.isDev = s.isDev

			app.Settings().Logs.MinLevel = s.level

			if err := app.Dao().SaveSettings(app.Settings()); err != nil {
				t.Fatalf("Failed to save settings: %v", err)
			}

			if err := app.RefreshSettings(); err != nil {
				t.Fatalf("Failed to refresh app settings: %v", err)
			}

			for level, enabled := range s.expectations {
				if v := handler.Enabled(nil, slog.Level(level)); v != enabled {
					t.Fatalf("Expected level %d Enabled() to be %v, got %v", level, enabled, v)
				}
			}
		})
	}
}

func TestBaseAppLoggerLevelDevPrint(t *testing.T) {
	app, cleanup, err := initTestBaseApp()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	testLogLevel := 4

	app.Settings().Logs.MinLevel = testLogLevel
	if err := app.Dao().SaveSettings(app.Settings()); err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name            string
		isDev           bool
		levels          []int
		printedLevels   []int
		persistedLevels []int
	}{
		{
			"dev mode",
			true,
			[]int{testLogLevel - 1, testLogLevel, testLogLevel + 1},
			[]int{testLogLevel - 1, testLogLevel, testLogLevel + 1},
			[]int{testLogLevel, testLogLevel + 1},
		},
		{
			"nondev mode",
			false,
			[]int{testLogLevel - 1, testLogLevel, testLogLevel + 1},
			[]int{},
			[]int{testLogLevel, testLogLevel + 1},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			var printedLevels []int
			var persistedLevels []int

			app.isDev = s.isDev

			// trigger slog handler min level refresh
			if err := app.RefreshSettings(); err != nil {
				t.Fatal(err)
			}

			// track printed logs
			originalPrintLog := printLog
			defer func() {
				printLog = originalPrintLog
			}()
			printLog = func(log *logger.Log) {
				printedLevels = append(printedLevels, int(log.Level))
			}

			// track persisted logs
			app.LogsDao().AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
				l, ok := m.(*models.Log)
				if ok {
					persistedLevels = append(persistedLevels, l.Level)
				}
				return nil
			}

			// write and persist logs
			for _, l := range s.levels {
				app.Logger().Log(nil, slog.Level(l), "test")
			}
			handler, ok := app.Logger().Handler().(*logger.BatchHandler)
			if !ok {
				t.Fatalf("Expected BatchHandler, got %v", app.Logger().Handler())
			}
			if err := handler.WriteAll(nil); err != nil {
				t.Fatalf("Failed to write all logs: %v", err)
			}

			// check persisted log levels
			if len(s.persistedLevels) != len(persistedLevels) {
				t.Fatalf("Expected persisted levels \n%v\ngot\n%v", s.persistedLevels, persistedLevels)
			}
			for _, l := range persistedLevels {
				if !list.ExistInSlice(l, s.persistedLevels) {
					t.Fatalf("Missing expected persisted level %v in %v", l, persistedLevels)
				}
			}

			// check printed log levels
			if len(s.printedLevels) != len(printedLevels) {
				t.Fatalf("Expected printed levels \n%v\ngot\n%v", s.printedLevels, printedLevels)
			}
			for _, l := range printedLevels {
				if !list.ExistInSlice(l, s.printedLevels) {
					t.Fatalf("Missing expected printed level %v in %v", l, printedLevels)
				}
			}
		})
	}
}

// -------------------------------------------------------------------

// note: make sure to call `defer cleanup()` when the app is no longer needed.
func initTestBaseApp() (app *BaseApp, cleanup func(), err error) {
	testDataDir, err := os.MkdirTemp("", "test_base_app")
	if err != nil {
		return nil, nil, err
	}

	cleanup = func() {
		os.RemoveAll(testDataDir)
	}

	app = NewBaseApp(BaseAppConfig{
		DataDir: testDataDir,
	})

	initErr := func() error {
		if err := app.Bootstrap(); err != nil {
			return fmt.Errorf("bootstrap error: %w", err)
		}

		logsRunner, err := migrate.NewRunner(app.LogsDB(), logs.LogsMigrations)
		if err != nil {
			return fmt.Errorf("logsRunner error: %w", err)
		}
		if _, err := logsRunner.Up(); err != nil {
			return fmt.Errorf("logsRunner migrations execution error: %w", err)
		}

		dataRunner, err := migrate.NewRunner(app.DB(), migrations.AppMigrations)
		if err != nil {
			return fmt.Errorf("logsRunner error: %w", err)
		}
		if _, err := dataRunner.Up(); err != nil {
			return fmt.Errorf("dataRunner migrations execution error: %w", err)
		}

		return nil
	}()
	if initErr != nil {
		cleanup()
		return nil, nil, initErr
	}

	return app, cleanup, nil
}
