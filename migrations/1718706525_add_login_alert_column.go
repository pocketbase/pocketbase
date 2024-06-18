package migrations

import (
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/security"
)

// adds a "lastLoginAlertSentAt" column to all auth collection tables (if not already)
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collections := []*models.Collection{}
		err := dao.CollectionQuery().AndWhere(dbx.HashExp{"type": models.CollectionTypeAuth}).All(&collections)
		if err != nil {
			return err
		}

		var needToResetTokens bool

		for _, c := range collections {
			columns, err := dao.TableColumns(c.Name)
			if err != nil {
				return err
			}
			if slices.Contains(columns, schema.FieldNameLastLoginAlertSentAt) {
				continue // already inserted
			}

			_, err = db.AddColumn(c.Name, schema.FieldNameLastLoginAlertSentAt, "TEXT DEFAULT '' NOT NULL").Execute()
			if err != nil {
				return err
			}

			opts := c.AuthOptions()
			if opts.AllowOAuth2Auth && (opts.AllowEmailAuth || opts.AllowUsernameAuth) {
				needToResetTokens = true
			}
		}

		settings, _ := dao.FindSettings()
		if needToResetTokens && settings != nil {
			settings.RecordAuthToken.Secret = security.RandomString(50)
			if err := dao.SaveSettings(settings); err != nil {
				return err
			}
		}

		return nil
	}, nil)
}
