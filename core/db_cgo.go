//go:build cgo

package core

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/dbx"
)

func connectDB(dbPath string) (*dbx.DB, error) {
	pragmas := "_foreign_keys=1&_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=8000"

	db, openErr := dbx.MustOpen("sqlite3", fmt.Sprintf("%s?%s", dbPath, pragmas))
	if openErr != nil {
		return nil, openErr
	}

	// additional pragmas not supported through the dsn string
	_, err := db.NewQuery(`
		pragma journal_size_limit = 100000000;
	`).Execute()

	return db, err
}
