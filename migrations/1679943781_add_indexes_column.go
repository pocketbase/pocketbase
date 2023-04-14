package migrations

import (
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/list"
)

// Adds _collections indexes column (if not already).
//
// Note: This migration will be deleted once schema.SchemaField.Unuique is removed.
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		// cleanup failed remaining/"dangling" temp views to prevent
		// errors during the indexes upsert
		// ---
		tempViews := []string{}
		viewsErr := db.Select("name").
			From("sqlite_schema").
			AndWhere(dbx.HashExp{"type": "view"}).
			AndWhere(dbx.NewExp(`[[name]] LIKE '\_temp\_%' ESCAPE '\'`)).
			Column(&tempViews)
		if viewsErr != nil {
			return viewsErr
		}
		for _, name := range tempViews {
			if err := dao.DeleteView(name); err != nil {
				return err
			}
		}
		// ---

		cols, err := dao.TableColumns("_collections")
		if err != nil {
			return err
		}

		var hasIndexesColumn bool
		for _, col := range cols {
			if col == "indexes" {
				// already existing (probably via the init migration)
				hasIndexesColumn = true
				break
			}
		}

		if !hasIndexesColumn {
			if _, err := db.AddColumn("_collections", "indexes", `JSON DEFAULT "[]" NOT NULL`).Execute(); err != nil {
				return err
			}
		}

		collections := []*models.Collection{}
		if err := dao.CollectionQuery().AndWhere(dbx.NewExp("type != 'view'")).All(&collections); err != nil {
			return err
		}

		type indexInfo struct {
			Sql       string `db:"sql"`
			IndexName string `db:"name"`
			TableName string `db:"tbl_name"`
		}

		indexesQuery := db.NewQuery(`SELECT * FROM sqlite_master WHERE type = "index" and sql is not null`)
		rawIndexes := []indexInfo{}
		if err := indexesQuery.All(&rawIndexes); err != nil {
			return err
		}

		indexesByTableName := map[string][]indexInfo{}
		for _, idx := range rawIndexes {
			indexesByTableName[idx.TableName] = append(indexesByTableName[idx.TableName], idx)
		}

		for _, c := range collections {
			c.Indexes = nil // reset

			excludeIndexes := []string{
				"_" + c.Id + "_email_idx",
				"_" + c.Id + "_username_idx",
				"_" + c.Id + "_tokenKey_idx",
			}

			// convert custom indexes into the related collections
			for _, idx := range indexesByTableName[c.Name] {
				if strings.Contains(idx.IndexName, "sqlite_autoindex_") ||
					list.ExistInSlice(idx.IndexName, excludeIndexes) {
					continue
				}

				// drop old index (it will be recreated with the collection)
				if _, err := db.DropIndex(idx.TableName, idx.IndexName).Execute(); err != nil {
					return err
				}

				c.Indexes = append(c.Indexes, idx.Sql)
			}

			// convert unique fields to indexes
		FieldsLoop:
			for _, f := range c.Schema.Fields() {
				if !f.Unique {
					continue
				}

				for _, idx := range indexesByTableName[c.Name] {
					parsed := dbutils.ParseIndex(idx.Sql)
					if parsed.Unique && len(parsed.Columns) == 1 && strings.EqualFold(parsed.Columns[0].Name, f.Name) {
						continue FieldsLoop // already added
					}
				}

				c.Indexes = append(c.Indexes, fmt.Sprintf(
					`CREATE UNIQUE INDEX "idx_unique_%s" on "%s" ("%s")`,
					f.Id,
					c.Name,
					f.Name,
				))
			}

			if len(c.Indexes) > 0 {
				if err := dao.SaveCollection(c); err != nil {
					return err
				}
			}
		}

		return nil
	}, func(db dbx.Builder) error {
		_, err := db.DropColumn("_collections", "indexes").Execute()

		return err
	})
}
