package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

// Update all collections with json fields to have a default MaxSize json field option.
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		// note: update even the view collections to prevent
		// unnecessary change detections during the automigrate
		collections := []*models.Collection{}
		if err := dao.CollectionQuery().All(&collections); err != nil {
			return err
		}

		for _, collection := range collections {
			var needSave bool

			for _, f := range collection.Schema.Fields() {
				if f.Type != schema.FieldTypeJson {
					continue
				}

				options, _ := f.Options.(*schema.JsonOptions)
				if options == nil {
					options = &schema.JsonOptions{}
				}
				options.MaxSize = 2000000 // 2mb
				f.Options = options
				needSave = true
			}

			if !needSave {
				continue
			}

			// save only the collection model without updating its records table
			if err := dao.Save(collection); err != nil {
				return err
			}
		}

		return nil
	}, nil)
}
