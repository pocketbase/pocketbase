package migrations

import (
	"github.com/pocketbase/dbx"
)

// Fixes the unique _externalAuths constraint for old installations
// to allow a single OAuth2 provider to be registered for different auth collections.
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		_, createErr := db.NewQuery("CREATE UNIQUE INDEX IF NOT EXISTS _externalAuths_collection_provider_idx on {{_externalAuths}} ([[collectionId]], [[provider]], [[providerId]])").Execute()
		if createErr != nil {
			return createErr
		}

		_, dropErr := db.NewQuery("DROP INDEX IF EXISTS _externalAuths_provider_providerId_idx").Execute()
		if dropErr != nil {
			return dropErr
		}

		return nil
	}, nil)
}
