package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

// Copy the now deprecated RelationOptions.DisplayFields values from
// all relation fields and register its value as Presentable under
// the specific field in the related collection.
//
// If there is more than one relation to a single collection with explicitly
// set DisplayFields only one of the configuration will be copied.
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collections := []*models.Collection{}
		if err := dao.CollectionQuery().All(&collections); err != nil {
			return err
		}

		indexedCollections := make(map[string]*models.Collection, len(collections))
		for _, collection := range collections {
			indexedCollections[collection.Id] = collection
		}

		for _, collection := range indexedCollections {
			for _, f := range collection.Schema.Fields() {
				if f.Type != schema.FieldTypeRelation {
					continue
				}

				options, ok := f.Options.(*schema.RelationOptions)
				if !ok || len(options.DisplayFields) == 0 {
					continue
				}

				relCollection, ok := indexedCollections[options.CollectionId]
				if !ok {
					continue
				}

				for _, name := range options.DisplayFields {
					relField := relCollection.Schema.GetFieldByName(name)
					if relField != nil {
						relField.Presentable = true
					}
				}

				// only raw model save
				if err := dao.Save(relCollection); err != nil {
					return err
				}
			}
		}

		return nil
	}, nil)
}
