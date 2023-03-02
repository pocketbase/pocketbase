package logs

import (
	"github.com/pocketbase/dbx"
)

// This migration normalizes the request logs method to UPPERCASE (eg. "get" => "GET").
func init() {
	LogsMigrations.Register(func(db dbx.Builder) error {
		_, err := db.NewQuery("UPDATE {{_requests}} SET method=UPPER(method)").Execute()

		return err
	}, func(db dbx.Builder) error {
		_, err := db.NewQuery("UPDATE {{_requests}} SET method=LOWER(method)").Execute()

		return err
	})
}
