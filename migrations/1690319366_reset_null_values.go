package migrations

import (
	"fmt"

	"github.com/lilysnc/pocketbasepg/daos"
	"github.com/lilysnc/pocketbasepg/models"
	"github.com/lilysnc/pocketbasepg/models/schema"
	"github.com/pocketbase/dbx"
)

// Reset all previously inserted NULL values to the fields zero-default.
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collections := []*models.Collection{}
		if err := dao.CollectionQuery().All(&collections); err != nil {
			return err
		}

		for _, collection := range collections {
			if collection.IsView() {
				continue
			}

			for _, f := range collection.Schema.Fields() {
				defaultVal := "''"

				switch f.Type {
				case schema.FieldTypeJson:
					continue
				case schema.FieldTypeBool:
					defaultVal = "FALSE"
				case schema.FieldTypeNumber:
					defaultVal = "0"
				default:
					if opt, ok := f.Options.(schema.MultiValuer); ok && opt.IsMultiple() {
						defaultVal = "'[]'"
					}
				}

				_, err := db.NewQuery(fmt.Sprintf(
					"UPDATE {{%s}} SET [[%s]] = %s WHERE [[%s]] IS NULL",
					collection.Name,
					f.Name,
					defaultVal,
					f.Name,
				)).Execute()
				if err != nil {
					return err
				}
			}
		}

		return nil
	}, nil)
}
