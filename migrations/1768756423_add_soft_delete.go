package migrations

import (
	"github.com/pocketbase/pocketbase/core"
)

func init() {
	core.SystemMigrations.Register(func(txApp core.App) error {
		// Check if column already exists (for new installations where init migration already has it)
		tableInfo, err := txApp.TableInfo("_collections")
		if err != nil {
			return err
		}

		for _, col := range tableInfo {
			if col.Name == "softDelete" {
				return nil // Column already exists, skip migration
			}
		}

		_, err = txApp.DB().NewQuery(`
			ALTER TABLE {{_collections}} ADD COLUMN [[softDelete]] BOOLEAN DEFAULT FALSE NOT NULL;
		`).Execute()
		if err != nil {
			return err
		}

		return nil
	}, func(txApp core.App) error {
		// SQLite doesn't support DROP COLUMN, so we have to recreate the table
		// For the down migration, we'll just leave the column as is
		// since it doesn't affect functionality
		return nil
	})
}
