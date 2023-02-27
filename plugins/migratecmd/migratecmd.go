// Package migratecmd adds a new "migrate" command support to a PocketBase instance.
//
// It also comes with automigrations support and templates generation
// (both for JS and GO migration files).
//
// Example usage:
//
//	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
//		TemplateLang: migratecmd.TemplateLangJS, // default to migratecmd.TemplateLangGo
//		Automigrate:  true,
//		Dir:          "migrations_dir_path", // optional template migrations path; default to "pb_migrations" (for JS) and "migrations" (for Go)
//	})
//
//	Note: To allow running JS migrations you'll need to enable first
//	[jsvm.MustRegisterMigrations].
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

// Options defines optional struct to customize the default plugin behavior.
type Options struct {
	// Dir specifies the directory with the user defined migrations.
	//
	// If not set it fallbacks to a relative "pb_data/../pb_migrations" (for js)
	// or "pb_data/../migrations" (for go) directory.
	Dir string

	// Automigrate specifies whether to enable automigrations.
	Automigrate bool

	// TemplateLang specifies the template language to use when
	// generating migrations - js or go (default).
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
	if p.options.Automigrate {
		// refresh the cache right after app bootstap
		p.app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
			p.refreshCachedCollections()
			return nil
		})

		// refresh the cache to ensure that it constains the latest changes
		// when migrations are applied on server start
		p.app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
			p.refreshCachedCollections()

			cachedCollections, _ := p.getCachedCollections()
			// create a full initial snapshot, if there are no custom
			// migrations but there is already at least 1 collection created,
			// to ensure that the automigrate will work with up-to-date collections data
			if !p.hasCustomMigrations() && len(cachedCollections) > 1 {
				p.migrateCollectionsHandler(nil, false)
			}

			return nil
		})

		p.app.OnModelAfterCreate().Add(p.afterCollectionChange())
		p.app.OnModelAfterUpdate().Add(p.afterCollectionChange())
		p.app.OnModelAfterDelete().Add(p.afterCollectionChange())
	}

	return nil
}

func (p *plugin) createCommand() *cobra.Command {
	const cmdDesc = `Supported arguments are:
- up            - runs all available migrations
- down [number] - reverts the last [number] applied migrations
- create name   - creates new blank migration template file
- collections   - creates new migration file with snapshot of the local collections configuration
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

			switch cmd {
			case "create":
				if err := p.migrateCreateHandler("", args[1:], true); err != nil {
					log.Fatal(err)
				}
			case "collections":
				if err := p.migrateCollectionsHandler(args[1:], true); err != nil {
					log.Fatal(err)
				}
			default:
				runner, err := migrate.NewRunner(p.app.DB(), migrations.AppMigrations)
				if err != nil {
					log.Fatal(err)
				}

				if err := runner.Run(args...); err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	return command
}

func (p *plugin) migrateCreateHandler(template string, args []string, interactive bool) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing migration file name")
	}

	name := args[0]
	dir := p.options.Dir

	resultFilePath := path.Join(
		dir,
		fmt.Sprintf("%d_%s.%s", time.Now().Unix(), inflector.Snakecase(name), p.options.TemplateLang),
	)

	if interactive {
		confirm := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Do you really want to create migration %q?", resultFilePath),
		}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			fmt.Println("The command has been cancelled")
			return nil
		}
	}

	// get default create template
	if template == "" {
		var templateErr error
		if p.options.TemplateLang == TemplateLangJS {
			template, templateErr = p.jsBlankTemplate()
		} else {
			template, templateErr = p.goBlankTemplate()
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

	if interactive {
		fmt.Printf("Successfully created file %q\n", resultFilePath)
	}

	return nil
}

func (p *plugin) migrateCollectionsHandler(args []string, interactive bool) error {
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

	return p.migrateCreateHandler(template, createArgs, interactive)
}
