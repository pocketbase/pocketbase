package core

import (
	"fmt"

	"github.com/pocketbase/dbx"
)

// LambdaFunctionQuery returns a new LambdaFunction select query.
func (app *BaseApp) LambdaFunctionQuery() *dbx.SelectQuery {
	return app.ModelQuery(&LambdaFunction{})
}

// FindLambdaFunctionById finds the first LambdaFunction by its id.
func (app *BaseApp) FindLambdaFunctionById(id string) (*LambdaFunction, error) {
	model := &LambdaFunction{}

	err := app.LambdaFunctionQuery().
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

// FindLambdaFunctionByName finds the first LambdaFunction by its name (case insensitive).
func (app *BaseApp) FindLambdaFunctionByName(name string) (*LambdaFunction, error) {
	model := &LambdaFunction{}

	err := app.LambdaFunctionQuery().
		AndWhere(dbx.NewExp("LOWER(name)={:name}", dbx.Params{"name": name})).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

// FindAllLambdaFunctions finds all LambdaFunction models.
func (app *BaseApp) FindAllLambdaFunctions() ([]*LambdaFunction, error) {
	models := []*LambdaFunction{}

	err := app.LambdaFunctionQuery().
		OrderBy("name ASC").
		All(&models)

	if err != nil {
		return nil, err
	}

	return models, nil
}

// FindEnabledLambdaFunctions finds all enabled LambdaFunction models.
func (app *BaseApp) FindEnabledLambdaFunctions() ([]*LambdaFunction, error) {
	models := []*LambdaFunction{}

	err := app.LambdaFunctionQuery().
		AndWhere(dbx.HashExp{"enabled": true}).
		OrderBy("name ASC").
		All(&models)

	if err != nil {
		return nil, err
	}

	return models, nil
}

// FindLambdaFunctionsByTriggerType finds all LambdaFunction models with a specific trigger type.
func (app *BaseApp) FindLambdaFunctionsByTriggerType(triggerType string) ([]*LambdaFunction, error) {
	models := []*LambdaFunction{}

	err := app.LambdaFunctionQuery().
		AndWhere(dbx.NewExp("json_extract(triggers, '$[*].type') LIKE {:type}", dbx.Params{
			"type": fmt.Sprintf("%%%s%%", triggerType),
		})).
		OrderBy("name ASC").
		All(&models)

	if err != nil {
		return nil, err
	}

	// Filter in Go to ensure exact match (SQLite JSON functions can be imprecise)
	var filtered []*LambdaFunction
	for _, model := range models {
		for _, trigger := range model.Triggers {
			if trigger.Type == triggerType {
				filtered = append(filtered, model)
				break
			}
		}
	}

	return filtered, nil
}

// FindLambdaFunctionsByCollection finds all LambdaFunction models that have database triggers for a specific collection.
func (app *BaseApp) FindLambdaFunctionsByCollection(collectionName string) ([]*LambdaFunction, error) {
	models := []*LambdaFunction{}

	err := app.LambdaFunctionQuery().
		AndWhere(dbx.NewExp("json_extract(triggers, '$[*].config.collection') LIKE {:collection}", dbx.Params{
			"collection": fmt.Sprintf("%%%s%%", collectionName),
		})).
		OrderBy("name ASC").
		All(&models)

	if err != nil {
		return nil, err
	}

	// Filter in Go to ensure exact match
	var filtered []*LambdaFunction
	for _, model := range models {
		configs, _ := model.GetDatabaseTriggers()
		for _, config := range configs {
			if config.Collection == collectionName {
				filtered = append(filtered, model)
				break
			}
		}
	}

	return filtered, nil
}

// IsLambdaFunctionNameUnique checks if a lambda function name is unique.
// For new functions, oldNames should be empty.
// For existing functions, oldNames should contain the current name.
func (app *BaseApp) IsLambdaFunctionNameUnique(name string, oldNames ...string) bool {
	if name == "" {
		return false
	}

	query := app.LambdaFunctionQuery().
		Select("COUNT(*)").
		AndWhere(dbx.NewExp("LOWER(name)={:name}", dbx.Params{"name": name})).
		Limit(1)

	if len(oldNames) > 0 {
		excludeIds := make([]any, len(oldNames))
		for i, oldName := range oldNames {
			excludeIds[i] = oldName
		}
		query.AndWhere(dbx.NotIn("LOWER(name)", excludeIds...))
	}

	var exists bool
	query.Row(&exists)

	return !exists
}