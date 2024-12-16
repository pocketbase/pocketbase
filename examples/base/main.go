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
    token = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhIjoicnciLCJpYXQiOjE3MzQyODE1NzUsImlkIjoiZGUzYTVkNGYtZDNlOS00MTc5LWE5ZjYtZWM4MzE5YTljNTk3In0.hoKahQLjmaEJmoVANfUQ3e9uouft_GsMHDKz2pnYJMcPXCq2heDc2gFWbLHGXcPD0oCfXN3RXm8QHhv0zEYQCQ"// replace with your actual token
    databaseURL = "libsql://pb-1st-turso-db-naol-bm.turso.io" // replace with your actual database url
)


func init() {
    dbx.BuilderFuncMap["libsql"] = dbx.BuilderFuncMap["sqlite3"]
}

var enableLibSQL = false

func main() {
	// config libsql
   app := pocketbase.NewWithConfig(pocketbase.Config{
        DBConnect: func(dbPath string) (*dbx.DB, error) {
        	if(enableLibSQL){
            	if strings.Contains(dbPath, "data.db") {
                	return dbx.Open("libsql", databaseURL+"?authToken="+token)
            	}
        	}
            // optionally for the logs (aka. pb_data/auxiliary.db) use the default local filesystem driver
            return core.DefaultDBConnect(dbPath)
        },
        DefaultDataDir: "pb_custom_dir_name",
    })

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
