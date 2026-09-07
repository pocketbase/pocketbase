package core_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"slices"
	"testing"
	"testing/synctest"
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
	defer app.ClearBootstrap()

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
		{"[before] db", app.DB(), false},
		{"[before] concurrentDB", app.ConcurrentDB(), false},
		{"[before] nonconcurrentDB", app.NonconcurrentDB(), false},
		{"[before] auxDB", app.AuxDB(), false},
		{"[before] auxConcurrentDB", app.AuxConcurrentDB(), false},
		{"[before] auxNonconcurrentDB", app.AuxNonconcurrentDB(), false},
		{"[before] settings", app.Settings(), false},
		{"[before] logger", app.Logger(), false},
		{"[before] cached collections", app.Store().Get(core.StoreKeyCachedCollections), false},
	}

	runNilChecks(nilChecksBeforeReset)

	// reset
	if err := app.ClearBootstrap(); err != nil {
		t.Fatal(err)
	}

	nilChecksAfterReset := []nilCheck{
		{"[after] db", app.DB(), true},
		{"[after] concurrentDB", app.ConcurrentDB(), true},
		{"[after] nonconcurrentDB", app.NonconcurrentDB(), true},
		{"[after] auxDB", app.AuxDB(), true},
		{"[after] auxConcurrentDB", app.AuxConcurrentDB(), true},
		{"[after] auxNonconcurrentDB", app.AuxNonconcurrentDB(), true},
		{"[after] settings", app.Settings(), false},
		{"[after] logger", app.Logger(), false},
		{"[after] cached collections", app.Store().Get(core.StoreKeyCachedCollections), false},
	}

	runNilChecks(nilChecksAfterReset)
}

func TestNewBaseAppTx(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir: testDataDir,
	})
	defer app.ClearBootstrap()

	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	mustNotHaveTx := func(app core.App) {
		if app.IsTransactional() {
			t.Fatalf("Didn't expect the app to be transactional")
		}

		if app.TxInfo() != nil {
			t.Fatalf("Didn't expect the app.txInfo to be loaded")
		}
	}

	mustHaveTx := func(app core.App) {
		if !app.IsTransactional() {
			t.Fatalf("Expected the app to be transactional")
		}

		if app.TxInfo() == nil {
			t.Fatalf("Expected the app.txInfo to be loaded")
		}
	}

	mustNotHaveTx(app)

	app.RunInTransaction(func(txApp core.App) error {
		mustHaveTx(txApp)
		return nil
	})

	mustNotHaveTx(app)
}

func TestBaseAppNewMailClient(t *testing.T) {
	const testDataDir = "./pb_base_app_test_data_dir/"
	defer os.RemoveAll(testDataDir)

	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir:       testDataDir,
		EncryptionEnv: "pb_test_env",
	})
	defer app.ClearBootstrap()

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
	defer app.ClearBootstrap()

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
	defer app.ClearBootstrap()

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

const logsThreshold = 200

func assertLogsCount(t *testing.T, app core.App, expected int) {
	var total int

	err := app.LogQuery().Select("count(*)").Row(&total)
	if err != nil {
		t.Fatalf("Failed to fetch total logs: %v", err)
	}

	if total != expected {
		t.Fatalf("Expected %d log(s), got %d", expected, total)
	}
}

func TestBaseAppLoggerWrites(t *testing.T) {
	t.Parallel()

	// note: outside of synctest because the bootstrap tickers could deadlock
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// clear old logs
	err := app.DeleteOldLogs(time.Now())
	if err != nil {
		t.Fatal(err)
	}

	t.Run("disabled logs retention", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			app.Settings().Logs.MaxDays = 0

			for i := 0; i < logsThreshold+1; i++ {
				app.Logger().Error("test")
			}

			// short delay for the non-blocking write goroutine
			synctest.Sleep(time.Nanosecond)

			assertLogsCount(t, app, 0)
		})
	})

	t.Run("test batch logs writes", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			app.Settings().Logs.MaxDays = 2

			for i := 0; i < logsThreshold-1; i++ {
				app.Logger().Error("test")
			}

			// short delay for the non-blocking write goroutine
			synctest.Sleep(time.Nanosecond)

			// below threshold
			assertLogsCount(t, app, 0)

			// threshold reached -> should trigger batch write
			app.Logger().Error("test")

			// should be skipped from this batch and added for the next
			app.Logger().Error("test")

			// short delay for the non-blocking write goroutine
			synctest.Sleep(time.Nanosecond)

			assertLogsCount(t, app, logsThreshold)

			// note: we can't test the flush timer here because the ticker
			// was started out of the synctest buble to avoid deadlocks
			// (see TestBaseAppLoggerWritesAwaited for a flaky but real timer test)
		})
	})
}

func TestBaseAppLoggerWritesAwaited(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// clear old logs
	err := app.DeleteOldLogs(time.Now())
	if err != nil {
		t.Fatal(err)
	}

	// enable logs persistence
	app.Settings().Logs.MaxDays = 1
	err = app.Save(app.Settings())
	if err != nil {
		t.Fatal(err)
	}

	t.Run("flush on timer tick", func(t *testing.T) {
		timeout := time.After(5 * time.Second)
		done := make(chan struct{})

		logsHook := app.OnModelAfterCreateSuccess("_logs")
		hookId := logsHook.BindFunc(func(e *core.ModelEvent) error {
			done <- struct{}{}
			return e.Next()
		})
		defer logsHook.Unbind(hookId)

		app.Logger().Error("test")

		// short wait to ensure that there is no non-blocking write
		time.Sleep(500 * time.Millisecond)

		assertLogsCount(t, app, 0)

		// wait for the ticker to write the db record
		select {
		case <-timeout:
			t.Fatal("ticker wait timeout")
		case <-done:
		}

		assertLogsCount(t, app, 1)
	})

	t.Run("before ClearBootstrap flush", func(t *testing.T) {
		app.Logger().Error("test")

		app.Bootstrap()

		assertLogsCount(t, app, 2)
	})

	t.Run("batch flush inside aux transaction shouldn't hang", func(t *testing.T) {
		timeout := time.After(1 * time.Second)
		done := make(chan struct{})
		totalCreated := 0

		logsHook := app.OnModelAfterCreateSuccess("_logs")
		hookId := logsHook.BindFunc(func(e *core.ModelEvent) error {
			totalCreated++
			if totalCreated == 200 {
				done <- struct{}{}
			}
			return e.Next()
		})
		defer logsHook.Unbind(hookId)

		app.AuxRunInTransaction(func(txApp core.App) error {
			for range logsThreshold {
				txApp.Logger().Error("test")
			}

			return nil
		})

		// wait for the non-blocking write
		select {
		case <-timeout:
			t.Fatal("non-blocking write timeout")
		case <-done:
		}

		assertLogsCount(t, app, 202)

		// force clear to ensure that there are no other logs
		app.Bootstrap()

		assertLogsCount(t, app, 202)
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
			defer app.ClearBootstrap()

			if err := app.Bootstrap(); err != nil {
				t.Fatal(err)
			}

			// silence query logs
			app.ConcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {}
			app.ConcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {}
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

func TestBaseAppDBDualBuilder(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	concurrentQueries := []string{}
	nonconcurrentQueries := []string{}
	app.ConcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		concurrentQueries = append(concurrentQueries, sql)
	}
	app.ConcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		concurrentQueries = append(concurrentQueries, sql)
	}
	app.NonconcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		nonconcurrentQueries = append(nonconcurrentQueries, sql)
	}
	app.NonconcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		nonconcurrentQueries = append(nonconcurrentQueries, sql)
	}

	type testQuery struct {
		query        string
		isConcurrent bool
	}

	regularTests := []testQuery{
		{"  \n  sEleCt 1", true},
		{"With abc(x) AS (select 2) SELECT x FROM abc", true},
		{"create table t1(x int)", false},
		{"insert into t1(x) values(1)", false},
		{"update t1 set x = 2", false},
		{"delete from t1", false},
	}

	txTests := []testQuery{
		{"select 3", false},
		{" \n WITH abc(x) AS (select 4) SELECT x FROM abc", false},
		{"create table t2(x int)", false},
		{"insert into t2(x) values(1)", false},
		{"update t2 set x = 2", false},
		{"delete from t2", false},
	}

	for _, item := range regularTests {
		_, err := app.DB().NewQuery(item.query).Execute()
		if err != nil {
			t.Fatalf("Failed to execute query %q error: %v", item.query, err)
		}
	}

	app.RunInTransaction(func(txApp core.App) error {
		for _, item := range txTests {
			_, err := txApp.DB().NewQuery(item.query).Execute()
			if err != nil {
				t.Fatalf("Failed to execute query %q error: %v", item.query, err)
			}
		}

		return nil
	})

	allTests := append(regularTests, txTests...)
	for _, item := range allTests {
		if item.isConcurrent {
			if !slices.Contains(concurrentQueries, item.query) {
				t.Fatalf("Expected concurrent query\n%q\ngot\nconcurrent:%v\nnonconcurrent:%v", item.query, concurrentQueries, nonconcurrentQueries)
			}
		} else {
			if !slices.Contains(nonconcurrentQueries, item.query) {
				t.Fatalf("Expected nonconcurrent query\n%q\ngot\nconcurrent:%v\nnonconcurrent:%v", item.query, concurrentQueries, nonconcurrentQueries)
			}
		}
	}
}

func TestBaseAppAuxDBDualBuilder(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	concurrentQueries := []string{}
	nonconcurrentQueries := []string{}
	app.AuxConcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		concurrentQueries = append(concurrentQueries, sql)
	}
	app.AuxConcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		concurrentQueries = append(concurrentQueries, sql)
	}
	app.AuxNonconcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		nonconcurrentQueries = append(nonconcurrentQueries, sql)
	}
	app.AuxNonconcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		nonconcurrentQueries = append(nonconcurrentQueries, sql)
	}

	type testQuery struct {
		query        string
		isConcurrent bool
	}

	regularTests := []testQuery{
		{"  \n  sEleCt 1", true},
		{"With abc(x) AS (select 2) SELECT x FROM abc", true},
		{"create table t1(x int)", false},
		{"insert into t1(x) values(1)", false},
		{"update t1 set x = 2", false},
		{"delete from t1", false},
	}

	txTests := []testQuery{
		{"select 3", false},
		{" \n WITH abc(x) AS (select 4) SELECT x FROM abc", false},
		{"create table t2(x int)", false},
		{"insert into t2(x) values(1)", false},
		{"update t2 set x = 2", false},
		{"delete from t2", false},
	}

	for _, item := range regularTests {
		_, err := app.AuxDB().NewQuery(item.query).Execute()
		if err != nil {
			t.Fatalf("Failed to execute query %q error: %v", item.query, err)
		}
	}

	app.AuxRunInTransaction(func(txApp core.App) error {
		for _, item := range txTests {
			_, err := txApp.AuxDB().NewQuery(item.query).Execute()
			if err != nil {
				t.Fatalf("Failed to execute query %q error: %v", item.query, err)
			}
		}

		return nil
	})

	allTests := append(regularTests, txTests...)
	for _, item := range allTests {
		if item.isConcurrent {
			if !slices.Contains(concurrentQueries, item.query) {
				t.Fatalf("Expected concurrent query\n%q\ngot\nconcurrent:%v\nnonconcurrent:%v", item.query, concurrentQueries, nonconcurrentQueries)
			}
		} else {
			if !slices.Contains(nonconcurrentQueries, item.query) {
				t.Fatalf("Expected nonconcurrent query\n%q\ngot\nconcurrent:%v\nnonconcurrent:%v", item.query, concurrentQueries, nonconcurrentQueries)
			}
		}
	}
}

func TestBaseAppTriggerOnTerminate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	event := new(core.TerminateEvent)
	event.App = app

	// trigger OnTerminate multiple times to ensure that it doesn't deadlock
	// https://github.com/pocketbase/pocketbase/pull/7305
	app.OnTerminate().Trigger(event)
	app.OnTerminate().Trigger(event)
	app.OnTerminate().Trigger(event)
}
