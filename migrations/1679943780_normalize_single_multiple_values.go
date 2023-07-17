package migrations

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

// Normalizes old single and multiple values of MultiValuer fields (file, select, relation).
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		return normalizeMultivaluerFields(db)
	}, func(db dbx.Builder) error {
		return nil
	})
}

func normalizeMultivaluerFields(db dbx.Builder) error {
	dao := daos.New(db)

	collections := []*models.Collection{}
	if err := dao.CollectionQuery().All(&collections); err != nil {
		return err
	}

	for _, c := range collections {
		if c.IsView() {
			// skip view collections
			continue
		}

		for _, f := range c.Schema.Fields() {
			opt, ok := f.Options.(schema.MultiValuer)
			if !ok {
				continue
			}

			var updateQuery *dbx.Query

			if opt.IsMultiple() {
				updateQuery = dao.DB().NewQuery(fmt.Sprintf(
					`UPDATE {{%s}} set [[%s]] = (
						CASE
							WHEN COALESCE([[%s]], '') = ''
							THEN '[]'
							ELSE (
								CASE
									WHEN json_valid([[%s]]) AND json_type([[%s]]) == 'array'
									THEN [[%s]]
									ELSE json_array([[%s]])
								END
							)
						END
					)`,
					c.Name,
					f.Name,
					f.Name,
					f.Name,
					f.Name,
					f.Name,
					f.Name,
				))
			} else {
				updateQuery = dao.DB().NewQuery(fmt.Sprintf(
					`UPDATE {{%s}} set [[%s]] = (
						CASE
							WHEN COALESCE([[%s]], '[]') = '[]'
							THEN ''
							ELSE (
								CASE
									WHEN json_valid([[%s]]) AND json_type([[%s]]) == 'array'
									THEN COALESCE(json_extract([[%s]], '$[#-1]'), '')
									ELSE [[%s]]
								END
							)
						END
					)`,
					c.Name,
					f.Name,
					f.Name,
					f.Name,
					f.Name,
					f.Name,
					f.Name,
				))
			}

			if _, err := updateQuery.Execute(); err != nil {
				return err
			}
		}
	}

	// trigger view query update after the records normalization
	// (ignore save error in case of invalid query to allow users to change it from the UI)
	for _, c := range collections {
		if !c.IsView() {
			continue
		}

		dao.SaveCollection(c)
	}

	return nil
}
