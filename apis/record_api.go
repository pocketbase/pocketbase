package apis

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/search"
)

// Create - create single record in collection with provided fields/values (fieldsMap)
// check rights/rules
// don't trigger hooks/logs event
// return error or record struct
func RecordCreate(app core.App, collection *models.Collection, admin *models.Admin, fieldsMap map[string]any) (*models.Record, error) {
	hasFullManageAccess, errCreate := createTest(app, collection, admin, fieldsMap)
	if errCreate != nil {
		return nil, errCreate
	}

	record := models.NewRecord(collection)
	form := forms.NewRecordUpsert(app, record)
	form.SetFullManageAccess(hasFullManageAccess)

	if err := form.LoadData(fieldsMap); err != nil {
		return nil, err
	}
	// create the record
	if err := form.Submit(); err != nil {
		return nil, err
	}
	return record, nil
}

// createTest check access and create rules before create if admin is nil
func createTest(app core.App, collection *models.Collection, admin *models.Admin, fieldsMap map[string]any) (bool, error) {
	hasFullManageAccess := admin != nil

	if collection == nil {
		return hasFullManageAccess, NewNotFoundError("", "Missing collection context.")
	}
	if admin == nil && collection.CreateRule == nil {
		// only admins can access if the rule is nil
		return hasFullManageAccess, NewForbiddenError("Only admins can perform this action.", nil)
	}

	// temporary save the record and check it against the create rule
	if admin == nil && collection.CreateRule != nil {
		createRuleFunc := func(q *dbx.SelectQuery) error {
			if *collection.CreateRule == "" {
				return nil // no create rule to resolve
			}

			resolver := resolvers.NewRecordFieldResolver(app.Dao(), collection, fieldsMap, true)
			expr, err := search.FilterData(*collection.CreateRule).BuildExpr(resolver)
			if err != nil {
				return err
			}
			resolver.UpdateQuery(q)
			q.AndWhere(expr)
			return nil
		}

		testRecord := models.NewRecord(collection)
		testForm := forms.NewRecordUpsert(app, testRecord)
		testForm.SetFullManageAccess(true)
		if err := testForm.LoadData(fieldsMap); err != nil {
			return hasFullManageAccess, NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
		}

		testErr := testForm.DrySubmit(func(txDao *daos.Dao) error {
			foundRecord, err := txDao.FindRecordById(collection.Id, testRecord.Id, createRuleFunc)
			if err != nil {
				return fmt.Errorf("DrySubmit create rule failure: %v", err)
			}
			hasFullManageAccess = hasAuthManageAccess(txDao, foundRecord, fieldsMap)
			return nil
		})

		if testErr != nil {
			return hasFullManageAccess, NewBadRequestError("Failed to create record.", testErr)
		}
	}
	return hasFullManageAccess, nil
}
