package core

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/logger"
)

func TestBaseAppLoggerLevelDevPrint(t *testing.T) {
	testLogLevel := 4

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
			const testDataDir = "./pb_base_app_test_data_dir/"
			defer os.RemoveAll(testDataDir)

			app := NewBaseApp(BaseAppConfig{
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

			app.Settings().Logs.MinLevel = testLogLevel
			if err := app.Save(app.Settings()); err != nil {
				t.Fatal(err)
			}

			var printedLevels []int
			var persistedLevels []int

			ctx := context.Background()

			// track printed logs
			originalPrintLog := printLog
			defer func() {
				printLog = originalPrintLog
			}()
			printLog = func(log *logger.Log) {
				printedLevels = append(printedLevels, int(log.Level))
			}

			// track persisted logs
			app.OnModelAfterCreateSuccess("_logs").BindFunc(func(e *ModelEvent) error {
				l, ok := e.Model.(*Log)
				if ok {
					persistedLevels = append(persistedLevels, l.Level)
				}
				return e.Next()
			})

			// write and persist logs
			for _, l := range s.levels {
				app.Logger().Log(ctx, slog.Level(l), "test")
			}
			handler, ok := app.Logger().Handler().(*logger.BatchHandler)
			if !ok {
				t.Fatalf("Expected BatchHandler, got %v", app.Logger().Handler())
			}
			if err := handler.WriteAll(ctx); err != nil {
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
