package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	Register(func(app core.App) error {
		// Create the lambdas collection
		collection := core.NewBaseCollection(core.CollectionNameLambdaFunctions)
		collection.System = true
		
		// Set rules to only allow superusers to manage functions
		superuserRule := "@request.auth.collectionName = '_superusers'"
		collection.ListRule = types.Pointer(superuserRule)
		collection.ViewRule = types.Pointer(superuserRule)
		collection.CreateRule = types.Pointer(superuserRule)
		collection.UpdateRule = types.Pointer(superuserRule)
		collection.DeleteRule = types.Pointer(superuserRule)

		// Add fields
		collection.Fields.Add(&core.TextField{
			Name:     "name",
			Required: true,
			System:   true,
			Pattern:  "^[a-zA-Z0-9][a-zA-Z0-9_-]*$",
			Min:      3,
			Max:      50,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "code",
			Required: true,
			System:   true,
		})

		collection.Fields.Add(&core.BoolField{
			Name:   "enabled",
			System: true,
		})

		collection.Fields.Add(&core.NumberField{
			Name:    "timeout",
			System:  true,
			OnlyInt: true,
			Min:     types.Pointer(float64(1000)),    // 1 second minimum
			Max:     types.Pointer(float64(300000)),  // 5 minutes maximum
		})

		collection.Fields.Add(&core.JSONField{
			Name:     "triggers",
			System:   true,
			Required: true,
		})

		collection.Fields.Add(&core.JSONField{
			Name:   "envVars",
			System: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:   "description",
			System: true,
		})

		collection.Fields.Add(&core.AutodateField{
			Name:     "created",
			System:   true,
			OnCreate: true,
		})

		collection.Fields.Add(&core.AutodateField{
			Name:      "updated",
			System:    true,
			OnCreate:  true,
			OnUpdate:  true,
		})

		// Add indexes
		collection.AddIndex("idx_lambdas_name", false, "name", "")
		collection.AddIndex("idx_lambdas_enabled", false, "enabled", "")

		// Create the collection
		return app.Save(collection)

	}, func(app core.App) error {
		// Down migration - Delete the collection
		collection, err := app.FindCollectionByNameOrId(core.CollectionNameLambdaFunctions)
		if err == nil {
			return app.Delete(collection)
		}
		return nil
	})
}