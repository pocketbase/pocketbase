package migrations

import (
	"github.com/pocketbase/dbx"
)

func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		_, createErr := db.NewQuery(`
			CREATE TABLE {{_views}} (
				[[id]]         TEXT PRIMARY KEY,
				[[name]]     TEXT NOT NULL,
                [[sql]]   TEXT NOT NULL,
				[[created]]    TEXT DEFAULT "" NOT NULL,
				[[schema]]     JSON DEFAULT "[]" NOT NULL,
				[[updated]]    TEXT DEFAULT "" NOT NULL,
				[[listRule]]   TEXT DEFAULT NULL
			);
		`).Execute()
		return createErr
	}, func(db dbx.Builder) error {
		_, err := db.DropTable("_views").Execute()
		return err
	})
}
