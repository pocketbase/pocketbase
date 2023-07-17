package migrations

import (
	"github.com/pocketbase/dbx"
)

// Renormalizes old single and multiple values of MultiValuer fields (file, select, relation)
// (see https://github.com/pocketbase/pocketbase/issues/2930).
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		return normalizeMultivaluerFields(db)
	}, func(db dbx.Builder) error {
		return nil
	})
}
