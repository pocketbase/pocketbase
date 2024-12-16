package main

import (
	"log"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

const (
	databaseURL = "" // replace with your actual database url
	token       = "" // replace with your actual token
)

var enableLibSQL = false

func init() {
	dbx.BuilderFuncMap["libsql"] = dbx.BuilderFuncMap["sqlite3"]
}

func main() {
	// config libsql
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DBConnect: func(dbPath string) (*dbx.DB, error) {
			if enableLibSQL {
				if strings.Contains(dbPath, "data.db") {
					return dbx.Open("libsql", databaseURL+"?authToken="+token)
				}
			}
			// optionally for the logs (aka. pb_data/auxiliary.db) use the default local filesystem driver
			return core.DefaultDBConnect(dbPath)
		},
		// DefaultDataDir: "pb_custom_dir_name",
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
