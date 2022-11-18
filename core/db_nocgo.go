//go:build !cgo

package core

import (
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

func connectDB(dbPath string) (*dbx.DB, error) {
	// note: the busy_timeout pragma must be first because
	// the connection needs to be set to block on busy before WAL mode
	// is set in case it hasn't been already set by another connection
	pragmas := "_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=synchronous(NORMAL)&_pragma=journal_size_limit(100000000)"

	db, err := dbx.MustOpen("sqlite", fmt.Sprintf("%s?%s", dbPath, pragmas))
	if err != nil {
		return nil, err
	}

	// use a fixed connection pool to limit the SQLITE_BUSY errors and
	// reduce the open file descriptors
	// (the limits are arbitrary and may change in the future)
	//
	// @see https://gitlab.com/cznic/sqlite/-/issues/115
	db.DB().SetMaxOpenConns(1000)
	db.DB().SetMaxIdleConns(30)
	db.DB().SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}
