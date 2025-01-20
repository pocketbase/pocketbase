package core_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestMigrationsRunnerUpAndDown(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	callsOrder := []string{}

	l := core.MigrationsList{}
	l.Register(func(app core.App) error {
		callsOrder = append(callsOrder, "up2")
		return nil
	}, func(app core.App) error {
		callsOrder = append(callsOrder, "down2")
		return nil
	}, "2_test")
	l.Register(func(app core.App) error {
		callsOrder = append(callsOrder, "up3")
		return nil
	}, func(app core.App) error {
		callsOrder = append(callsOrder, "down3")
		return nil
	}, "3_test")
	l.Register(func(app core.App) error {
		callsOrder = append(callsOrder, "up1")
		return nil
	}, func(app core.App) error {
		callsOrder = append(callsOrder, "down1")
		return nil
	}, "1_test")
	l.Register(func(app core.App) error {
		callsOrder = append(callsOrder, "up4")
		return nil
	}, func(app core.App) error {
		callsOrder = append(callsOrder, "down4")
		return nil
	}, "4_test")
	l.Add(&core.Migration{
		Up: func(app core.App) error {
			callsOrder = append(callsOrder, "up5")
			return nil
		},
		Down: func(app core.App) error {
			callsOrder = append(callsOrder, "down5")
			return nil
		},
		File: "5_test",
		ReapplyCondition: func(txApp core.App, runner *core.MigrationsRunner, fileName string) (bool, error) {
			return true, nil
		},
	})

	runner := core.NewMigrationsRunner(app, l)

	// ---------------------------------------------------------------
	// simulate partially out-of-order applied migration
	// ---------------------------------------------------------------

	_, err := app.DB().Insert(core.DefaultMigrationsTable, dbx.Params{
		"file":    "4_test",
		"applied": time.Now().UnixMicro() - 2,
	}).Execute()
	if err != nil {
		t.Fatalf("Failed to insert 5_test migration: %v", err)
	}

	_, err = app.DB().Insert(core.DefaultMigrationsTable, dbx.Params{
		"file":    "5_test",
		"applied": time.Now().UnixMicro() - 1,
	}).Execute()
	if err != nil {
		t.Fatalf("Failed to insert 5_test migration: %v", err)
	}

	_, err = app.DB().Insert(core.DefaultMigrationsTable, dbx.Params{
		"file":    "2_test",
		"applied": time.Now().UnixMicro(),
	}).Execute()
	if err != nil {
		t.Fatalf("Failed to insert 2_test migration: %v", err)
	}

	// ---------------------------------------------------------------
	// Up()
	// ---------------------------------------------------------------

	if _, err := runner.Up(); err != nil {
		t.Fatal(err)
	}

	expectedUpCallsOrder := `["up1","up3","up5"]` // skip up2 and up4 since they were applied already (up5 has extra reapply condition)

	upCallsOrder, err := json.Marshal(callsOrder)
	if err != nil {
		t.Fatal(err)
	}

	if v := string(upCallsOrder); v != expectedUpCallsOrder {
		t.Fatalf("Expected Up() calls order %s, got %s", expectedUpCallsOrder, upCallsOrder)
	}

	// ---------------------------------------------------------------

	// reset callsOrder
	callsOrder = []string{}

	// simulate unrun migration
	l.Register(nil, func(app core.App) error {
		callsOrder = append(callsOrder, "down6")
		return nil
	}, "6_test")

	// simulate applied migrations from different migrations list
	_, err = app.DB().Insert(core.DefaultMigrationsTable, dbx.Params{
		"file":    "from_different_list",
		"applied": time.Now().UnixMicro(),
	}).Execute()
	if err != nil {
		t.Fatalf("Failed to insert from_different_list migration: %v", err)
	}

	// ---------------------------------------------------------------

	// ---------------------------------------------------------------
	// Down()
	// ---------------------------------------------------------------

	if _, err := runner.Down(2); err != nil {
		t.Fatal(err)
	}

	expectedDownCallsOrder := `["down5","down3"]` // revert in the applied order

	downCallsOrder, err := json.Marshal(callsOrder)
	if err != nil {
		t.Fatal(err)
	}

	if v := string(downCallsOrder); v != expectedDownCallsOrder {
		t.Fatalf("Expected Down() calls order %s, got %s", expectedDownCallsOrder, downCallsOrder)
	}
}

func TestMigrationsRunnerRemoveMissingAppliedMigrations(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// mock migrations history
	for i := 1; i <= 3; i++ {
		_, err := app.DB().Insert(core.DefaultMigrationsTable, dbx.Params{
			"file":    fmt.Sprintf("%d_test", i),
			"applied": time.Now().UnixMicro(),
		}).Execute()
		if err != nil {
			t.Fatal(err)
		}
	}

	if !isMigrationApplied(app, "2_test") {
		t.Fatalf("Expected 2_test migration to be applied")
	}

	// create a runner without 2_test to mock deleted migration
	l := core.MigrationsList{}
	l.Register(func(app core.App) error {
		return nil
	}, func(app core.App) error {
		return nil
	}, "1_test")
	l.Register(func(app core.App) error {
		return nil
	}, func(app core.App) error {
		return nil
	}, "3_test")

	r := core.NewMigrationsRunner(app, l)

	if err := r.RemoveMissingAppliedMigrations(); err != nil {
		t.Fatalf("Failed to remove missing applied migrations: %v", err)
	}

	if isMigrationApplied(app, "2_test") {
		t.Fatalf("Expected 2_test migration to NOT be applied")
	}
}

func isMigrationApplied(app core.App, file string) bool {
	var exists int

	err := app.DB().Select("(1)").
		From(core.DefaultMigrationsTable).
		Where(dbx.HashExp{"file": file}).
		Limit(1).
		Row(&exists)

	return err == nil && exists > 0
}
