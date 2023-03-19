package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

// Adds _collections indexes column (if not already).
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		cols, err := dao.GetTableColumns("_collections")
		if err != nil {
			return err
		}

		for _, col := range cols {
			if col == "indexes" {
				return nil // already existing (probably via the init migration)
			}
		}

		_, err = db.AddColumn("_collections", "indexes", `JSON DEFAULT "[]" NOT NULL`).Execute()

		// @todo populate existing indexes...

		return err
	}, func(db dbx.Builder) error {
		_, err := db.DropColumn("_collections", "indexes").Execute()

		return err
	})
}
