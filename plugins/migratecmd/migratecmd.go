package migratecmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/migrate"
	"github.com/spf13/cobra"
)

type Options struct {
	Dir          string // the directory with user defined migrations
	AutoMigrate  bool
	TemplateLang string
}

type plugin struct {
	app     core.App
	options *Options
}

func MustRegister(app core.App, rootCmd *cobra.Command, options *Options) {
	if err := Register(app, rootCmd, options); err != nil {
		panic(err)
	}
}

func Register(app core.App, rootCmd *cobra.Command, options *Options) error {
	p := &plugin{app: app}

	if options != nil {
		p.options = options
	} else {
		p.options = &Options{}
	}

	if p.options.TemplateLang == "" {
		p.options.TemplateLang = TemplateLangGo
	}

	if p.options.Dir == "" {
		if p.options.TemplateLang == TemplateLangJS {
			p.options.Dir = filepath.Join(p.app.DataDir(), "../pb_migrations")
		} else {
			p.options.Dir = filepath.Join(p.app.DataDir(), "../migrations")
		}
	}

	// attach the migrate command
	if rootCmd != nil {
		rootCmd.AddCommand(p.createCommand())
	}

	// watch for collection changes
	if p.options.AutoMigrate {
		// @todo replace with AfterBootstrap
		p.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
			if err := p.tidyMigrationsTable(); err != nil && p.app.IsDebug() {
				log.Println("Failed to tidy the migrations table.")
			}

			return nil
		})

		p.app.OnModelAfterCreate().Add(p.onCollectionChange())
		p.app.OnModelAfterUpdate().Add(p.onCollectionChange())
		p.app.OnModelAfterDelete().Add(p.onCollectionChange())
	}

	return nil
}

func (p *plugin) createCommand() *cobra.Command {
	const cmdDesc = `Supported arguments are:
- up                   - runs all available migrations.
- down [number]        - reverts the last [number] applied migrations.
- create name [folder] - creates new blank migration template file.
- collections [folder] - creates new migration file with the latest local collections snapshot (similar to the automigrate but allows editing).
`

	command := &cobra.Command{
		Use:       "migrate",
		Short:     "Executes app DB migration scripts",
		ValidArgs: []string{"up", "down", "create", "collections"},
		Long:      cmdDesc,
		Run: func(command *cobra.Command, args []string) {
			cmd := ""
			if len(args) > 0 {
				cmd = args[0]
			}

			// additional commands
			// ---
			if cmd == "create" {
				if err := p.migrateCreateHandler("", args[1:]); err != nil {
					log.Fatal(err)
				}
				return
			}

			if cmd == "collections" {
				if err := p.migrateCollectionsHandler(args[1:]); err != nil {
					log.Fatal(err)
				}
				return
			}
			// ---

			runner, err := migrate.NewRunner(p.app.DB(), migrations.AppMigrations)
			if err != nil {
				log.Fatal(err)
			}

			if err := runner.Run(args...); err != nil {
				log.Fatal(err)
			}
		},
	}

	return command
}

func (p *plugin) migrateCreateHandler(template string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing migration file name")
	}

	name := args[0]

	var dir string
	if len(args) == 2 {
		dir = args[1]
	}
	if dir == "" {
		dir = p.options.Dir
	}

	resultFilePath := path.Join(
		dir,
		fmt.Sprintf("%d_%s.%s", time.Now().Unix(), inflector.Snakecase(name), p.options.TemplateLang),
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

	// get default create template
	if template == "" {
		var templateErr error
		if p.options.TemplateLang == TemplateLangJS {
			template, templateErr = p.jsCreateTemplate()
		} else {
			template, templateErr = p.goCreateTemplate()
		}
		if templateErr != nil {
			return fmt.Errorf("Failed to resolve create template: %v\n", templateErr)
		}
	}

	// ensure that the migrations dir exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// save the migration file
	if err := os.WriteFile(resultFilePath, []byte(template), 0644); err != nil {
		return fmt.Errorf("Failed to save migration file %q: %v\n", resultFilePath, err)
	}

	fmt.Printf("Successfully created file %q\n", resultFilePath)
	return nil
}

func (p *plugin) migrateCollectionsHandler(args []string) error {
	createArgs := []string{"collections_snapshot"}
	createArgs = append(createArgs, args...)

	collections := []*models.Collection{}
	if err := p.app.Dao().CollectionQuery().OrderBy("created ASC").All(&collections); err != nil {
		return fmt.Errorf("Failed to fetch migrations list: %v", err)
	}

	var template string
	var templateErr error
	if p.options.TemplateLang == TemplateLangJS {
		template, templateErr = p.jsSnapshotTemplate(collections)
	} else {
		template, templateErr = p.goSnapshotTemplate(collections)
	}
	if templateErr != nil {
		return fmt.Errorf("Failed to resolve template: %v", templateErr)
	}

	return p.migrateCreateHandler(template, createArgs)
}
