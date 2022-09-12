package logs

import (
	"github.com/pocketbase/dbx"
	"log"
)

func init() {
	LogsMigrations.Register(func(db dbx.Builder) error {
		_, err := db.NewQuery(`
			DROP INDEX IF EXISTS _request_ip_idx;
			ALTER TABLE {{_requests}} RENAME COLUMN [[ip]] TO [[remoteIp]];
			ALTER TABLE {{_requests}} ADD COLUMN [[userIp]] TEXT DEFAULT '127.0.0.1' NOT NULL;
			CREATE INDEX _request_remote_ip_idx ON {{_requests}} ([[remoteIp]]);
			CREATE INDEX _request_user_ip_idx ON {{_requests}} ([[userIp]]);
		`).Execute()

		if err != nil {
			log.Println("up", err)
			return err
		}

		return nil
	}, func(db dbx.Builder) error {
		_, err := db.NewQuery(`
			DROP INDEX _request_remote_ip_idx;
			DROP INDEX _request_user_ip_idx;
			CREATE INDEX _request_ip_idx ON {{_requests}} ([[ip]]);

			ALTER TABLE {{_requests}} DROP COLUMN [[userIp]];
			ALTER TABLE {{_requests}} RENAME COLUMN [[remoteIp]] TO [[ip]];
			
		`).Execute()
		if err != nil {
			log.Println("drop", err)
			return err
		}

		return nil
	})
}
