package migrate

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/list"
	_ "modernc.org/sqlite"
)

func TestNewRunner(t *testing.T) {
	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	l := MigrationsList{}
	l.Register(nil, nil, "1_test.go")
	l.Register(nil, nil, "2_test.go")
	l.Register(nil, nil, "3_test.go")

	r, err := NewRunner(testDB.DB, l)
	if err != nil {
		t.Fatal(err)
	}

	if len(r.migrationsList.Items()) != len(l.Items()) {
		t.Fatalf("Expected the same migrations list to be assigned, got \n%#v", r.migrationsList)
	}

	expectedQueries := []string{
		"CREATE TABLE IF NOT EXISTS `_migrations` (file VARCHAR(255) PRIMARY KEY NOT NULL, applied INTEGER NOT NULL)",
	}
	if len(expectedQueries) != len(testDB.CalledQueries) {
		t.Fatalf("Expected %d queries, got %d: \n%v", len(expectedQueries), len(testDB.CalledQueries), testDB.CalledQueries)
	}
	for _, q := range expectedQueries {
		if !list.ExistInSlice(q, testDB.CalledQueries) {
			t.Fatalf("Query %s was not found in \n%v", q, testDB.CalledQueries)
		}
	}
}

func TestRunnerUpAndDown(t *testing.T) {
	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	var test1UpCalled bool
	var test1DownCalled bool
	var test2UpCalled bool
	var test2DownCalled bool

	l := MigrationsList{}
	l.Register(func(db dbx.Builder) error {
		test1UpCalled = true
		return nil
	}, func(db dbx.Builder) error {
		test1DownCalled = true
		return nil
	}, "1_test")
	l.Register(func(db dbx.Builder) error {
		test2UpCalled = true
		return nil
	}, func(db dbx.Builder) error {
		test2DownCalled = true
		return nil
	}, "2_test")

	r, err := NewRunner(testDB.DB, l)
	if err != nil {
		t.Fatal(err)
	}

	// simulate partially run migration
	r.saveAppliedMigration(testDB, r.migrationsList.Item(0).File)

	// Up()
	// ---
	if _, err := r.Up(); err != nil {
		t.Fatal(err)
	}

	if test1UpCalled {
		t.Fatalf("Didn't expect 1_test to be called")
	}

	if !test2UpCalled {
		t.Fatalf("Expected 2_test to be called")
	}

	// simulate unrun migration
	var test3DownCalled bool
	r.migrationsList.Register(nil, func(db dbx.Builder) error {
		test3DownCalled = true
		return nil
	}, "3_test")

	// Down()
	// ---
	// revert one migration
	if _, err := r.Down(1); err != nil {
		t.Fatal(err)
	}

	if test3DownCalled {
		t.Fatal("Didn't expect 3_test to be reverted.")
	}

	if !test2DownCalled {
		t.Fatal("Expected 2_test to be reverted.")
	}

	if test1DownCalled {
		t.Fatal("Didn't expect 1_test to be reverted.")
	}
}

// -------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------

type testDB struct {
	*dbx.DB
	CalledQueries []string
}

// NB! Don't forget to call `db.Close()` at the end of the test.
func createTestDB() (*testDB, error) {
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	db := testDB{DB: dbx.NewFromDB(sqlDB, "sqlite")}
	db.QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		db.CalledQueries = append(db.CalledQueries, sql)
	}
	db.ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		db.CalledQueries = append(db.CalledQueries, sql)
	}

	return &db, nil
}
