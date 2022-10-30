package apis

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/spf13/cast"
)

// exportRequestData exports a map with common request fields.
//
// @todo consider changing the map to a typed struct after v0.8 and the
// IN operator support.
func exportRequestData(c echo.Context) map[string]any {
	result := map[string]any{}
	queryParams := map[string]any{}
	bodyData := map[string]any{}
	method := c.Request().Method

	echo.BindQueryParams(c, &queryParams)

	rest.BindBody(c, &bodyData)

	result["method"] = method
	result["query"] = queryParams
	result["data"] = bodyData
	result["auth"] = nil

	auth, _ := c.Get(ContextAuthRecordKey).(*models.Record)
	if auth != nil {
		result["auth"] = auth.PublicExport()
	}

	return result
}

// expandFetch is the records fetch function that is used to expand related records.
func expandFetch(
	dao *daos.Dao,
	isAdmin bool,
	requestData map[string]any,
) daos.ExpandFetchFunc {
	return func(relCollection *models.Collection, relIds []string) ([]*models.Record, error) {
		records, err := dao.FindRecordsByIds(relCollection.Id, relIds, func(q *dbx.SelectQuery) error {
			if isAdmin {
				return nil // admins can access everything
			}

			if relCollection.ViewRule == nil {
				return fmt.Errorf("Only admins can view collection %q records", relCollection.Name)
			}

			if *relCollection.ViewRule != "" {
				resolver := resolvers.NewRecordFieldResolver(dao, relCollection, requestData, true)
				expr, err := search.FilterData(*(relCollection.ViewRule)).BuildExpr(resolver)
				if err != nil {
					return err
				}
				resolver.UpdateQuery(q)
				q.AndWhere(expr)
			}

			return nil
		})

		if err == nil && len(records) > 0 {
			autoIgnoreAuthRecordsEmailVisibility(dao, records, isAdmin, requestData)
		}

		return records, err
	}
}

// autoIgnoreAuthRecordsEmailVisibility ignores the email visibility check for
// the provided record if the current auth model is admin, owner or a "manager".
//
// Note: Expects all records to be from the same auth collection!
func autoIgnoreAuthRecordsEmailVisibility(
	dao *daos.Dao,
	records []*models.Record,
	isAdmin bool,
	requestData map[string]any,
) error {
	if len(records) == 0 || !records[0].Collection().IsAuth() {
		return nil // nothing to check
	}

	if isAdmin {
		for _, rec := range records {
			rec.IgnoreEmailVisibility(true)
		}
		return nil
	}

	collection := records[0].Collection()

	mappedRecords := make(map[string]*models.Record, len(records))
	recordIds := make([]any, 0, len(records))
	for _, rec := range records {
		mappedRecords[rec.Id] = rec
		recordIds = append(recordIds, rec.Id)
	}

	if auth, ok := requestData["auth"].(map[string]any); ok && mappedRecords[cast.ToString(auth["id"])] != nil {
		mappedRecords[cast.ToString(auth["id"])].IgnoreEmailVisibility(true)
	}

	authOptions := collection.AuthOptions()
	if authOptions.ManageRule == nil || *authOptions.ManageRule == "" {
		return nil // no manage rule to check
	}

	// fetch the ids of the managed records
	// ---
	managedIds := []string{}

	query := dao.RecordQuery(collection).
		Select(dao.DB().QuoteSimpleColumnName(collection.Name) + ".id").
		AndWhere(dbx.In(dao.DB().QuoteSimpleColumnName(collection.Name)+".id", recordIds...))

	resolver := resolvers.NewRecordFieldResolver(dao, collection, requestData, true)
	expr, err := search.FilterData(*authOptions.ManageRule).BuildExpr(resolver)
	if err != nil {
		return err
	}
	resolver.UpdateQuery(query)
	query.AndWhere(expr)

	if err := query.Column(&managedIds); err != nil {
		return err
	}
	// ---

	// ignore the email visibility check for the managed records
	for _, id := range managedIds {
		if rec, ok := mappedRecords[id]; ok {
			rec.IgnoreEmailVisibility(true)
		}
	}

	return nil
}

// hasAuthManageAccess checks whether the client is allowed to have full
// [forms.RecordUpsert] auth management permissions
// (aka. allowing to change system auth fields without oldPassword).
func hasAuthManageAccess(
	dao *daos.Dao,
	record *models.Record,
	requestData map[string]any,
) bool {
	if !record.Collection().IsAuth() {
		return false
	}

	manageRule := record.Collection().AuthOptions().ManageRule

	if manageRule == nil || *manageRule == "" {
		return false // only for admins (manageRule can't be empty)
	}

	if auth, ok := requestData["auth"].(map[string]any); !ok || cast.ToString(auth["id"]) == "" {
		return false // no auth record
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		resolver := resolvers.NewRecordFieldResolver(dao, record.Collection(), requestData, true)
		expr, err := search.FilterData(*manageRule).BuildExpr(resolver)
		if err != nil {
			return err
		}
		resolver.UpdateQuery(q)
		q.AndWhere(expr)
		return nil
	}

	_, findErr := dao.FindRecordById(record.Collection().Id, record.Id, ruleFunc)

	return findErr == nil
}
