package logs

import (
	"github.com/pocketbase/dbx"
)

func init() {
	LogsMigrations.Register(func(db dbx.Builder) error {
		if _, err := db.DropTable("_requests").Execute(); err != nil {
			return err
		}

		_, err := db.NewQuery(`
			CREATE TABLE {{_logs}} (
				[[id]]      TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL,
				[[level]]   INTEGER DEFAULT 0 NOT NULL,
				[[message]] TEXT DEFAULT "" NOT NULL,
				[[data]]    JSON DEFAULT "{}" NOT NULL,
				[[created]] TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
				[[updated]] TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
			);

			CREATE INDEX _logs_level_idx on {{_logs}} ([[level]]);
			CREATE INDEX _logs_message_idx on {{_logs}} ([[message]]);
			CREATE INDEX _logs_created_hour_idx on {{_logs}} (strftime('%Y-%m-%d %H:00:00', [[created]]));
		`).Execute()

		return err
	}, func(db dbx.Builder) error {
		if _, err := db.DropTable("_logs").Execute(); err != nil {
			return err
		}

		_, err := db.NewQuery(`
			CREATE TABLE {{_requests}} (
				[[id]]        TEXT PRIMARY KEY NOT NULL,
				[[url]]       TEXT DEFAULT "" NOT NULL,
				[[method]]    TEXT DEFAULT "get" NOT NULL,
				[[status]]    INTEGER DEFAULT 200 NOT NULL,
				[[auth]]      TEXT DEFAULT "guest" NOT NULL,
				[[ip]]        TEXT DEFAULT "127.0.0.1" NOT NULL,
				[[referer]]   TEXT DEFAULT "" NOT NULL,
				[[userAgent]] TEXT DEFAULT "" NOT NULL,
				[[meta]]      JSON DEFAULT "{}" NOT NULL,
				[[created]]   TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL,
				[[updated]]   TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL
			);

			CREATE INDEX _request_status_idx on {{_requests}} ([[status]]);
			CREATE INDEX _request_auth_idx on {{_requests}} ([[auth]]);
			CREATE INDEX _request_ip_idx on {{_requests}} ([[ip]]);
			CREATE INDEX _request_created_hour_idx on {{_requests}} (strftime('%Y-%m-%d %H:00:00', [[created]]));
		`).Execute()

		return err
	})
}
