package core_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/store"
	"golang.org/x/sync/semaphore"
)

func TestNotifyWatcher_SettingsUpdate(t *testing.T) {
	t.Parallel()

	testEvents := store.New[core.App, int](nil)

	tmpDir, err := os.MkdirTemp("", "pb_notify_test*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	app1 := core.NewBaseApp(core.BaseAppConfig{
		DataDir: tmpDir,
	})
	if err := app1.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	app2 := core.NewBaseApp(core.BaseAppConfig{
		DataDir: tmpDir,
	})
	if err := app2.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	timeout := time.After(3 * time.Second)
	done := make(chan struct{})

	app1.OnSettingsReload().BindFunc(func(e *core.SettingsReloadEvent) error {
		testEvents.SetFunc(app1, func(old int) int {
			return old + 1
		})
		return e.Next()
	})

	app2.OnSettingsReload().BindFunc(func(e *core.SettingsReloadEvent) error {
		testEvents.SetFunc(app2, func(old int) int {
			defer func() {
				done <- struct{}{}
			}()

			return old + 1
		})
		return e.Next()
	})

	// updating app1 settings should trigger a reload in app2
	app1.Settings().SuperuserIPs = []string{"127.0.0.1"}
	if err := app1.Save(app1.Settings()); err != nil {
		t.Fatal(err)
	}

	// wait for the event
	select {
	case <-timeout:
		t.Fatal("app2 reload event timeout")
	case <-done:
		// ready
	}

	if app1Total := testEvents.Get(app1); app1Total != 1 {
		t.Fatalf("Expected 1 app1 event, got %d", app1Total)
	}

	if app2Total := testEvents.Get(app2); app2Total != 1 {
		t.Fatalf("Expected 1 app2 event, got %d", app2Total)
	}

	app2SuperuserIPs := app2.Settings().SuperuserIPs
	if len(app2SuperuserIPs) != 1 || app2SuperuserIPs[0] != "127.0.0.1" {
		t.Fatalf("Expected exactly 127.0.0.1 superuser IP in app2 settings event, got %v", app2SuperuserIPs)
	}
}

func TestNotifyWatcher_CollectionsUpdate(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "pb_notify_test*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	app1 := core.NewBaseApp(core.BaseAppConfig{
		DataDir: tmpDir,
	})
	if err := app1.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	app2 := core.NewBaseApp(core.BaseAppConfig{
		DataDir: tmpDir,
	})
	if err := app2.Bootstrap(); err != nil {
		t.Fatal(err)
	}

	testQueries := store.New[string, []string](nil)
	app2.ConcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		testQueries.SetFunc("concurrent", func(old []string) []string {
			return append(old, sql)
		})
	}
	app2.ConcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		testQueries.SetFunc("concurrent", func(old []string) []string {
			return append(old, sql)
		})
	}
	app2.NonconcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		testQueries.SetFunc("nonconcurrent", func(old []string) []string {
			return append(old, sql)
		})
	}
	app2.NonconcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		testQueries.SetFunc("nonconcurrent", func(old []string) []string {
			return append(old, sql)
		})
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelCtx()

	sem := semaphore.NewWeighted(1)
	sem.Acquire(ctx, 1)

	// currently there is no hook for the collections cache reload so we pool instead
	done := make(chan bool, 1)
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				if len(testQueries.Get("concurrent")) == 1 {
					sem.Release(1)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// create/update/delete app1 collections should trigger a reload in app2
	dummyCollection := core.NewBaseCollection("test")
	if err := app1.Save(dummyCollection); err != nil {
		t.Fatal(err)
	}
	dummyCollection.Fields.Add(&core.TextField{Name: "test"})
	if err := app1.Save(dummyCollection); err != nil {
		t.Fatal(err)
	}
	if err := app1.Delete(dummyCollection); err != nil {
		t.Fatal(err)
	}

	// block until released or timeouted
	sem.Acquire(ctx, 1)
	ticker.Stop()
	done <- true

	nonconcurrentQueries := testQueries.Get("nonconcurrent")
	concurrentQueries := testQueries.Get("concurrent")

	if len(nonconcurrentQueries) != 0 {
		t.Fatalf("Expected 0 concurrent queries, got %d (%v)", len(nonconcurrentQueries), nonconcurrentQueries)
	}
	if len(concurrentQueries) != 1 {
		t.Fatalf("Expected 1 concurrent query, got %d (%v)", len(concurrentQueries), concurrentQueries)
	}

	expectedQuery := "SELECT {{_collections}}.* FROM `_collections` ORDER BY `rowid` ASC"
	if concurrentQueries[0] != expectedQuery {
		t.Fatalf("Expected query\n%s\ngot\n%s", expectedQuery, concurrentQueries[0])
	}
}
