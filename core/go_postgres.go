package core

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pocketbase/dbx"
)

func connectDB(dbPath string) (*dbx.DB, error) {
	return dbx.MustOpen("postgres", fmt.Sprintf("%s", dbPath))
}
