//go:build !cgo

package core

import (
	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

func connectDB(dbPath string) (*dbx.DB, error) {
	db, err := dbx.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := initPragmas(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
