package migrations

import "github.com/pocketbase/dbx"

func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		_, createErr := db.NewQuery(`
			CREATE TABLE {{_externalAuths}} (
				[[id]]         TEXT PRIMARY KEY,
				[[userId]]     TEXT NOT NULL,
				[[provider]]   TEXT NOT NULL,
				[[providerId]] TEXT NOT NULL,
				[[created]]    TEXT DEFAULT "" NOT NULL,
				[[updated]]    TEXT DEFAULT "" NOT NULL,
				---
				FOREIGN KEY ([[userId]]) REFERENCES {{_users}} ([[id]]) ON UPDATE CASCADE ON DELETE CASCADE
			);

			CREATE UNIQUE INDEX _externalAuths_userId_provider_idx on {{_externalAuths}} ([[userId]], [[provider]]);
			CREATE UNIQUE INDEX _externalAuths_provider_providerId_idx on {{_externalAuths}} ([[provider]], [[providerId]]);
		`).Execute()
		if createErr != nil {
			return createErr
		}

		// remove the unique email index from the _users table and
		// replace it with partial index
		_, alterErr := db.NewQuery(`
			-- crate new users table
			CREATE TABLE {{_newUsers}} (
				[[id]]                     TEXT PRIMARY KEY,
				[[verified]]               BOOLEAN DEFAULT FALSE NOT NULL,
				[[email]]                  TEXT DEFAULT "" NOT NULL,
				[[tokenKey]]               TEXT NOT NULL,
				[[passwordHash]]           TEXT NOT NULL,
				[[lastResetSentAt]]        TEXT DEFAULT "" NOT NULL,
				[[lastVerificationSentAt]] TEXT DEFAULT "" NOT NULL,
				[[created]]                TEXT DEFAULT "" NOT NULL,
				[[updated]]                TEXT DEFAULT "" NOT NULL
			);

			-- copy all data from the old users table to the new one
			INSERT INTO {{_newUsers}} SELECT * FROM {{_users}};

			-- drop old table
			DROP TABLE {{_users}};

			-- rename new table
			ALTER TABLE {{_newUsers}} RENAME TO {{_users}};

			-- create named indexes
			CREATE UNIQUE INDEX _users_email_idx ON {{_users}} ([[email]]) WHERE [[email]] != "";
			CREATE UNIQUE INDEX _users_tokenKey_idx ON {{_users}} ([[tokenKey]]);
		`).Execute()
		if alterErr != nil {
			return alterErr
		}

		return nil
	}, func(db dbx.Builder) error {
		if _, err := db.DropTable("_externalAuths").Execute(); err != nil {
			return err
		}

		// drop the partial email unique index and replace it with normal unique index
		_, indexErr := db.NewQuery(`
			DROP INDEX IF EXISTS _users_email_idx;
			CREATE UNIQUE INDEX _users_email_idx on {{_users}} ([[email]]);
		`).Execute()
		if indexErr != nil {
			return indexErr
		}

		return nil
	})
}
