package logs

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/migrate"
)

var LogsMigrations migrate.MigrationsList

func init() {
	LogsMigrations.Register(func(db dbx.Builder) (err error) {
		_, err = db.NewQuery(`

			CREATE TABLE {{_requests}} (
				[[id]]        TEXT PRIMARY KEY,
				[[url]]       TEXT DEFAULT '' NOT NULL,
				[[method]]    TEXT DEFAULT 'get' NOT NULL,
				[[status]]    INTEGER DEFAULT 200 NOT NULL,
				[[auth]]      TEXT DEFAULT 'guest' NOT NULL,
				[[ip]]        TEXT DEFAULT '127.0.0.1' NOT NULL,
				[[referer]]   TEXT DEFAULT '' NOT NULL,
				[[userAgent]] TEXT DEFAULT '' NOT NULL,
				[[meta]]      JSON DEFAULT '{}' NOT NULL,
				[[created]]   TEXT DEFAULT '' NOT NULL,
				[[updated]]   TEXT DEFAULT '' NOT NULL
			);
		`).Execute()

		if err != nil {
			return err
		}

		_, err = db.NewQuery(`
			CREATE INDEX _request_status_idx on {{_requests}} ([[status]]);
			CREATE INDEX _request_auth_idx on {{_requests}} ([[auth]]);
			CREATE INDEX _request_ip_idx on {{_requests}} ([[ip]]);
		`).Execute()
		// TODO: Enable this again, by changing to to_date
		//CREATE INDEX _request_created_hour_idx on {{_requests}} (strftime('%Y-%m-%d %H:00:00', [[created]]));

		return err
	}, func(db dbx.Builder) error {
		_, err := db.DropTable("_requests").Execute()
		return err
	})
}
