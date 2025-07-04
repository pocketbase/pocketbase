package migrations

import (
	"github.com/pocketbase/pocketbase/core"
)

func init() {
	Register(func(app core.App) error {
		// Find the lambdas collection
		collection, err := app.FindCollectionByNameOrId(core.CollectionNameLambdaFunctions)
		if err != nil {
			return err
		}

		// Add the description field if it doesn't exist
		field := collection.Fields.GetByName("description")
		if field == nil {
			collection.Fields.Add(&core.TextField{
				Name:   "description",
				System: true,
			})
			
			return app.Save(collection)
		}
		
		return nil
	}, func(app core.App) error {
		// Down migration - remove the description field
		collection, err := app.FindCollectionByNameOrId(core.CollectionNameLambdaFunctions)
		if err != nil {
			return nil // Collection doesn't exist, nothing to do
		}

		field := collection.Fields.GetByName("description")
		if field != nil {
			collection.Fields.RemoveById(field.GetId())
			return app.Save(collection)
		}
		
		return nil
	})
}