package cmd

import (
	"log"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/migrations/logs"
	"github.com/pocketbase/pocketbase/tools/migrate"
	"github.com/spf13/cobra"
)

// NewMigrateCommand creates and returns new command for handling DB migrations.
func NewMigrateCommand(app core.App) *cobra.Command {
	desc := `
Supported arguments are:
- up                 - runs all available migrations.
- down [number]      - reverts the last [number] applied migrations.
- create folder name - creates new migration template file.
`
	var databaseFlag string

	command := &cobra.Command{
		Use:       "migrate",
		Short:     "Executes DB migration scripts",
		ValidArgs: []string{"up", "down", "create"},
		Long:      desc,
		Run: func(command *cobra.Command, args []string) {
			// normalize
			if databaseFlag != "logs" {
				databaseFlag = "db"
			}

			connections := migrationsConnectionsMap(app)

			runner, err := migrate.NewRunner(
				connections[databaseFlag].DB,
				connections[databaseFlag].MigrationsList,
			)
			if err != nil {
				log.Fatal(err)
			}

			if err := runner.Run(args...); err != nil {
				log.Fatal(err)
			}
		},
	}

	command.PersistentFlags().StringVar(
		&databaseFlag,
		"database",
		"db",
		"specify the database connection to use (db or logs)",
	)

	return command
}

type migrationsConnection struct {
	DB             *dbx.DB
	MigrationsList migrate.MigrationsList
}

func migrationsConnectionsMap(app core.App) map[string]migrationsConnection {
	return map[string]migrationsConnection{
		"db": {
			DB:             app.DB(),
			MigrationsList: migrations.AppMigrations,
		},
		"logs": {
			DB:             app.LogsDB(),
			MigrationsList: logs.LogsMigrations,
		},
	}
}
