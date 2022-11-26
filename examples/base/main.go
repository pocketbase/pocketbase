package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/plugins/publicdir"
)

func main() {
	app := pocketbase.New()

	// load js pb_migrations
	jsvm.MustRegisterMigrationsLoader(app, nil)

	// migrate command (with js templates)
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		TemplateLang: migratecmd.TemplateLangJS,
		AutoMigrate:  true,
	})

	// pb_public dir
	publicdir.MustRegister(app, &publicdir.Options{
		FlagsCmd:      app.RootCmd,
		IndexFallback: true,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
