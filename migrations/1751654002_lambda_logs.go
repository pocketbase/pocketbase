package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	Register(func(app core.App) error {
		// Create the lambda_logs collection
		collection := core.NewBaseCollection("lambda_logs")
		collection.System = true
		
		// Set rules to only allow superusers to view logs
		superuserRule := "@request.auth.collectionName = '_superusers'"
		collection.ListRule = types.Pointer(superuserRule)
		collection.ViewRule = types.Pointer(superuserRule)
		// No create/update/delete rules - only the system can create logs

		// Add fields
		collection.Fields.Add(&core.TextField{
			Name:     "function_id",
			Required: true,
			System:   true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "function_name",
			Required: true,
			System:   true,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "trigger_type",
			Required: true,
			System:   true,
		})

		collection.Fields.Add(&core.BoolField{
			Name:   "success",
			System: true,
		})

		collection.Fields.Add(&core.JSONField{
			Name:   "output",
			System: true,
		})

		collection.Fields.Add(&core.TextField{
			Name:   "error",
			System: true,
		})

		collection.Fields.Add(&core.NumberField{
			Name:    "duration_ms",
			System:  true,
			OnlyInt: true,
			Min:     types.Pointer(float64(0)),
		})

		collection.Fields.Add(&core.JSONField{
			Name:   "context",
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
		collection.AddIndex("idx_lambda_logs_function_id", false, "function_id", "")
		collection.AddIndex("idx_lambda_logs_created", false, "created", "")

		// Create the collection
		return app.Save(collection)

	}, func(app core.App) error {
		// Down migration - Delete the collection
		collection, err := app.FindCollectionByNameOrId("lambda_logs")
		if err == nil {
			return app.Delete(collection)
		}
		return nil
	})
}