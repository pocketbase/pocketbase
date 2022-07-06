//go:build !cgo

package core

import (
	"fmt"

	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

func connectDB(dbPath string) (*dbx.DB, error) {
	pragmas := "_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(8000)&_pragma=journal_size_limit(100000000)"

	return dbx.MustOpen("sqlite", fmt.Sprintf("%s?%s", dbPath, pragmas))
}
