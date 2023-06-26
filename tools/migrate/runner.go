package migrate

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/pocketbase/dbx"
	"github.com/spf13/cast"
)

const DefaultMigrationsTable = "_migrations"

// Runner defines a simple struct for managing the execution of db migrations.
type Runner struct {
	db             *dbx.DB
	migrationsList MigrationsList
	tableName      string
}

// NewRunner creates and initializes a new db migrations Runner instance.
func NewRunner(db *dbx.DB, migrationsList MigrationsList) (*Runner, error) {
	runner := &Runner{
		db:             db,
		migrationsList: migrationsList,
		tableName:      DefaultMigrationsTable,
	}

	if err := runner.createMigrationsTable(); err != nil {
		return nil, err
	}

	return runner, nil
}

// Run interactively executes the current runner with the provided args.
//
// The following commands are supported:
// - up       - applies all migrations
// - down [n] - reverts the last n applied migrations
func (r *Runner) Run(args ...string) error {
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

		confirm := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf(
				"\n%v\nDo you really want to revert the last %d applied migration(s)?",
				strings.Join(names, "\n"),
				toRevertCount,
			),
		}
		survey.AskOne(prompt, &confirm)
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
		if err := r.removeMissingAppliedMigrations(); err != nil {
			return err
		}

		color.Green("The %s table was synced with the available migrations.", r.tableName)
		return nil
	default:
		return fmt.Errorf("Unsupported command: %q\n", cmd)
	}
}

// Up executes all unapplied migrations for the provided runner.
//
// On success returns list with the applied migrations file names.
func (r *Runner) Up() ([]string, error) {
	applied := []string{}

	err := r.db.Transactional(func(tx *dbx.Tx) error {
		for _, m := range r.migrationsList.Items() {
			// skip applied
			if r.isMigrationApplied(tx, m.File) {
				continue
			}

			// ignore empty Up action
			if m.Up != nil {
				if err := m.Up(tx); err != nil {
					return fmt.Errorf("Failed to apply migration %s: %w", m.File, err)
				}
			}

			if err := r.saveAppliedMigration(tx, m.File); err != nil {
				return fmt.Errorf("Failed to save applied migration info for %s: %w", m.File, err)
			}

			applied = append(applied, m.File)
		}

		return nil
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
func (r *Runner) Down(toRevertCount int) ([]string, error) {
	reverted := make([]string, 0, toRevertCount)

	names, appliedErr := r.lastAppliedMigrations(toRevertCount)
	if appliedErr != nil {
		return nil, appliedErr
	}

	err := r.db.Transactional(func(tx *dbx.Tx) error {
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
					if err := m.Down(tx); err != nil {
						return fmt.Errorf("Failed to revert migration %s: %w", m.File, err)
					}
				}

				if err := r.saveRevertedMigration(tx, m.File); err != nil {
					return fmt.Errorf("Failed to save reverted migration info for %s: %w", m.File, err)
				}

				reverted = append(reverted, m.File)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return reverted, nil
}

func (r *Runner) createMigrationsTable() error {
	rawQuery := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %v (file VARCHAR(255) PRIMARY KEY NOT NULL, applied INTEGER NOT NULL)",
		r.db.QuoteTableName(r.tableName),
	)

	_, err := r.db.NewQuery(rawQuery).Execute()

	return err
}

func (r *Runner) isMigrationApplied(tx dbx.Builder, file string) bool {
	var exists bool

	err := tx.Select("count(*)").
		From(r.tableName).
		Where(dbx.HashExp{"file": file}).
		Limit(1).
		Row(&exists)

	return err == nil && exists
}

func (r *Runner) saveAppliedMigration(tx dbx.Builder, file string) error {
	_, err := tx.Insert(r.tableName, dbx.Params{
		"file":    file,
		"applied": time.Now().UnixMicro(),
	}).Execute()

	return err
}

func (r *Runner) saveRevertedMigration(tx dbx.Builder, file string) error {
	_, err := tx.Delete(r.tableName, dbx.HashExp{"file": file}).Execute()

	return err
}

func (r *Runner) lastAppliedMigrations(limit int) ([]string, error) {
	var files = make([]string, 0, limit)

	err := r.db.Select("file").
		From(r.tableName).
		Where(dbx.Not(dbx.HashExp{"applied": nil})).
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

func (r *Runner) removeMissingAppliedMigrations() error {
	loadedMigrations := r.migrationsList.Items()

	names := make([]any, len(loadedMigrations))
	for i, migration := range loadedMigrations {
		names[i] = migration.File
	}

	_, err := r.db.Delete(r.tableName, dbx.Not(dbx.HashExp{
		"file": names,
	})).Execute()

	return err
}
