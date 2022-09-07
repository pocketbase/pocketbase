package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/migrations/logs"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/migrate"
	"github.com/spf13/cobra"
)

// NewMigrateCommand creates and returns new command for handling DB migrations.
func NewMigrateCommand(app core.App) *cobra.Command {
	desc := `
Supported arguments are:
- up                   - runs all available migrations.
- down [number]        - reverts the last [number] applied migrations.
- create name [folder] - creates new migration template file.
- collections [folder] - (Experimental) creates new migration file with the most recent local collections configuration.
`
	var databaseFlag string

	command := &cobra.Command{
		Use:       "migrate",
		Short:     "Executes DB migration scripts",
		ValidArgs: []string{"up", "down", "create", "collections"},
		Long:      desc,
		Run: func(command *cobra.Command, args []string) {
			cmd := ""
			if len(args) > 0 {
				cmd = args[0]
			}

			// additional commands
			// ---
			if cmd == "create" {
				if err := migrateCreateHandler(defaultMigrateCreateTemplate, args[1:]); err != nil {
					log.Fatal(err)
				}
				return
			}
			if cmd == "collections" {
				if err := migrateCollectionsHandler(app, args[1:]); err != nil {
					log.Fatal(err)
				}
				return
			}
			// ---

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

// -------------------------------------------------------------------
// migrate create
// -------------------------------------------------------------------

const defaultMigrateCreateTemplate = `package migrations

import (
	"github.com/pocketbase/dbx"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		// add up queries...

		return nil
	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
`

func migrateCreateHandler(template string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing migration file name")
	}

	name := args[0]

	var dir string
	if len(args) == 2 {
		dir = args[1]
	}
	if dir == "" {
		// If not specified, auto point to the default migrations folder.
		//
		// NB!
		// Since the create command makes sense only during development,
		// it is expected the user to be in the app working directory
		// and to be using `go run`
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = path.Join(wd, "migrations")
	}

	resultFilePath := path.Join(
		dir,
		fmt.Sprintf("%d_%s.go", time.Now().Unix(), inflector.Snakecase(name)),
	)

	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Do you really want to create migration %q?", resultFilePath),
	}
	survey.AskOne(prompt, &confirm)
	if !confirm {
		fmt.Println("The command has been cancelled")
		return nil
	}

	// ensure that migrations dir exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(resultFilePath, []byte(template), 0644); err != nil {
		return fmt.Errorf("Failed to save migration file %q\n", resultFilePath)
	}

	fmt.Printf("Successfully created file %q\n", resultFilePath)
	return nil
}

// -------------------------------------------------------------------
// migrate collections
// -------------------------------------------------------------------

const collectionsMigrateCreateTemplate = `package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

// Auto generated migration with the most recent collections configuration.
func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := ` + "`" + `%s` + "`" + `

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		// no revert since the configuration on the environment, on which
		// the migration was executed, could have changed via the UI/API
		return nil
	})
}
`

func migrateCollectionsHandler(app core.App, args []string) error {
	createArgs := []string{"collections_snapshot"}
	createArgs = append(createArgs, args...)

	dao := daos.New(app.DB())

	collections := []*models.Collection{}
	if err := dao.CollectionQuery().OrderBy("created ASC").All(&collections); err != nil {
		return fmt.Errorf("Failed to fetch migrations list: %v", err)
	}

	serialized, err := json.MarshalIndent(collections, "\t\t", "\t")
	if err != nil {
		return fmt.Errorf("Failed to serialize collections list: %v", err)
	}

	return migrateCreateHandler(
		fmt.Sprintf(collectionsMigrateCreateTemplate, string(serialized)),
		createArgs,
	)
}
