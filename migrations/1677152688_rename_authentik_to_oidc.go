package migrations

import (
	"github.com/pocketbase/dbx"
)

// This migration replaces the "authentikAuth" setting with "oidc".
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		_, err := db.NewQuery(`
			UPDATE {{_params}}
			SET [[value]] = replace([[value]]::text, '"authentikAuth":', '"oidcAuth":')::json
			WHERE [[key]] = 'settings'
		`).Execute()

		return err
	}, func(db dbx.Builder) error {
		_, err := db.NewQuery(`
			UPDATE {{_params}}
			SET [[value]] = replace([[value]]::text, '"oidcAuth":', '"authentikAuth":')::json
			WHERE [[key]] = 'settings'
		`).Execute()

		return err
	})
}
