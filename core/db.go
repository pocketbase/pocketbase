package core

import (
	"github.com/pocketbase/dbx"
)

func initPragmas(db *dbx.DB) error {
	// note: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection
	_, err := db.NewQuery(`
		PRAGMA busy_timeout       = 10000;
		PRAGMA journal_mode       = WAL;
		PRAGMA journal_size_limit = 200000000;
		PRAGMA synchronous        = NORMAL;
		PRAGMA foreign_keys       = TRUE;
	`).Execute()

	return err
}
