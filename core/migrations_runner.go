package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/osutils"
	"github.com/spf13/cast"
)

var AppMigrations MigrationsList
var SystemMigrations MigrationsList

const DefaultMigrationsTable = "_migrations"

// MigrationsRunner defines a simple struct for managing the execution of db migrations.
type MigrationsRunner struct {
	app            App
	tableName      string
	migrationsList MigrationsList
	inited         bool
}

// NewMigrationsRunner creates and initializes a new db migrations MigrationsRunner instance.
func NewMigrationsRunner(app App, migrationsList MigrationsList) *MigrationsRunner {
	return &MigrationsRunner{
		app:            app,
		migrationsList: migrationsList,
		tableName:      DefaultMigrationsTable,
	}
}

// Run interactively executes the current runner with the provided args.
//
// The following commands are supported:
// - up           - applies all migrations
// - down [n]     - reverts the last n (default 1) applied migrations
// - history-sync - syncs the migrations table with the runner's migrations list
func (r *MigrationsRunner) Run(args ...string) error {
	if err := r.initMigrationsTable(); err != nil {
		return err
	}

	cmd := "up"
	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case "up":
		applied, err := r.Up()
		if err != nil {
			return err
		}

		if len(applied) == 0 {
			color.Green("No new migrations to apply.")
		} else {
			for _, file := range applied {
				color.Green("Applied %s", file)
			}
		}

		return nil
	case "down":
		toRevertCount := 1
		if len(args) > 1 {
			toRevertCount = cast.ToInt(args[1])
			if toRevertCount < 0 {
				// revert all applied migrations
				toRevertCount = len(r.migrationsList.Items())
			}
		}

		names, err := r.lastAppliedMigrations(toRevertCount)
		if err != nil {
			return err
		}

		confirm := osutils.YesNoPrompt(fmt.Sprintf(
			"\n%v\nDo you really want to revert the last %d applied migration(s)?",
			strings.Join(names, "\n"),
			toRevertCount,
		), false)
		if !confirm {
			fmt.Println("The command has been cancelled")
			return nil
		}

		reverted, err := r.Down(toRevertCount)
		if err != nil {
			return err
		}

		if len(reverted) == 0 {
			color.Green("No migrations to revert.")
		} else {
			for _, file := range reverted {
				color.Green("Reverted %s", file)
			}
		}

		return nil
	case "history-sync":
		if err := r.RemoveMissingAppliedMigrations(); err != nil {
			return err
		}

		color.Green("The %s table was synced with the available migrations.", r.tableName)
		return nil
	default:
		return fmt.Errorf("unsupported command: %q", cmd)
	}
}

// Up executes all unapplied migrations for the provided runner.
//
// On success returns list with the applied migrations file names.
func (r *MigrationsRunner) Up() ([]string, error) {
	if err := r.initMigrationsTable(); err != nil {
		return nil, err
	}

	applied := []string{}

	err := r.app.AuxRunInTransaction(func(txApp App) error {
		return txApp.RunInTransaction(func(txApp App) error {
			for _, m := range r.migrationsList.Items() {
				// applied migrations check
				if r.isMigrationApplied(txApp, m.File) {
					if m.ReapplyCondition == nil {
						continue // no need to reapply
					}

					shouldReapply, err := m.ReapplyCondition(txApp, r, m.File)
					if err != nil {
						return err
					}
					if !shouldReapply {
						continue
					}

					// clear previous history stored entry
					// (it will be recreated after successful execution)
					r.saveRevertedMigration(txApp, m.File)
				}

				// ignore empty Up action
				if m.Up != nil {
					if err := m.Up(txApp); err != nil {
						return fmt.Errorf("failed to apply migration %s: %w", m.File, err)
					}
				}

				if err := r.saveAppliedMigration(txApp, m.File); err != nil {
					return fmt.Errorf("failed to save applied migration info for %s: %w", m.File, err)
				}

				applied = append(applied, m.File)
			}

			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return applied, nil
}

// Down reverts the last `toRevertCount` applied migrations
// (in the order they were applied).
//
// On success returns list with the reverted migrations file names.
func (r *MigrationsRunner) Down(toRevertCount int) ([]string, error) {
	if err := r.initMigrationsTable(); err != nil {
		return nil, err
	}

	reverted := make([]string, 0, toRevertCount)

	names, appliedErr := r.lastAppliedMigrations(toRevertCount)
	if appliedErr != nil {
		return nil, appliedErr
	}

	err := r.app.AuxRunInTransaction(func(txApp App) error {
		return txApp.RunInTransaction(func(txApp App) error {
			for _, name := range names {
				for _, m := range r.migrationsList.Items() {
					if m.File != name {
						continue
					}

					// revert limit reached
					if toRevertCount-len(reverted) <= 0 {
						return nil
					}

					// ignore empty Down action
					if m.Down != nil {
						if err := m.Down(txApp); err != nil {
							return fmt.Errorf("failed to revert migration %s: %w", m.File, err)
						}
					}

					if err := r.saveRevertedMigration(txApp, m.File); err != nil {
						return fmt.Errorf("failed to save reverted migration info for %s: %w", m.File, err)
					}

					reverted = append(reverted, m.File)
				}
			}
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return reverted, nil
}

// RemoveMissingAppliedMigrations removes the db entries of all applied migrations
// that are not listed in the runner's migrations list.
func (r *MigrationsRunner) RemoveMissingAppliedMigrations() error {
	loadedMigrations := r.migrationsList.Items()

	names := make([]any, len(loadedMigrations))
	for i, migration := range loadedMigrations {
		names[i] = migration.File
	}

	_, err := r.app.DB().Delete(r.tableName, dbx.Not(dbx.HashExp{
		"file": names,
	})).Execute()

	return err
}

func (r *MigrationsRunner) initMigrationsTable() error {
	if r.inited {
		return nil // already inited
	}

	rawQuery := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS {{%s}} (file VARCHAR(255) PRIMARY KEY NOT NULL, applied INTEGER NOT NULL)",
		r.tableName,
	)

	_, err := r.app.DB().NewQuery(rawQuery).Execute()

	if err == nil {
		r.inited = true
	}

	return err
}

func (r *MigrationsRunner) isMigrationApplied(txApp App, file string) bool {
	var exists int

	err := txApp.DB().Select("(1)").
		From(r.tableName).
		Where(dbx.HashExp{"file": file}).
		Limit(1).
		Row(&exists)

	return err == nil && exists > 0
}

func (r *MigrationsRunner) saveAppliedMigration(txApp App, file string) error {
	_, err := txApp.DB().Insert(r.tableName, dbx.Params{
		"file":    file,
		"applied": time.Now().UnixMicro(),
	}).Execute()

	return err
}

func (r *MigrationsRunner) saveRevertedMigration(txApp App, file string) error {
	_, err := txApp.DB().Delete(r.tableName, dbx.HashExp{"file": file}).Execute()

	return err
}

func (r *MigrationsRunner) lastAppliedMigrations(limit int) ([]string, error) {
	var files = make([]string, 0, limit)

	loadedMigrations := r.migrationsList.Items()

	names := make([]any, len(loadedMigrations))
	for i, migration := range loadedMigrations {
		names[i] = migration.File
	}

	err := r.app.DB().Select("file").
		From(r.tableName).
		Where(dbx.Not(dbx.HashExp{"applied": nil})).
		AndWhere(dbx.HashExp{"file": names}).
		// unify microseconds and seconds applied time for backward compatibility
		OrderBy("substr(applied||'0000000000000000', 0, 17) DESC").
		AndOrderBy("file DESC").
		Limit(int64(limit)).
		Column(&files)

	if err != nil {
		return nil, err
	}

	return files, nil
}
