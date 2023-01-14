package apis

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/search"
)

const ContextRequestDataKey = "requestData"

// RequestData exports cached common request data fields
// (query, body, logged auth state, etc.) from the provided context.
func RequestData(c echo.Context) *models.RequestData {
	// return cached to avoid copying the body multiple times
	if v := c.Get(ContextRequestDataKey); v != nil {
		if data, ok := v.(*models.RequestData); ok {
			return data
		}
	}

	result := &models.RequestData{
		Method: c.Request().Method,
		Query:  map[string]any{},
		Data:   map[string]any{},
	}

	result.AuthRecord, _ = c.Get(ContextAuthRecordKey).(*models.Record)
	result.Admin, _ = c.Get(ContextAdminKey).(*models.Admin)
	echo.BindQueryParams(c, &result.Query)
	rest.BindBody(c, &result.Data)

	c.Set(ContextRequestDataKey, result)

	return result
}

// EnrichRecord parses the request context and enrich the provided record:
// - expands relations (if defaultExpands and/or ?expand query param is set)
// - ensures that the emails of the auth record and its expanded auth relations
//   are visibe only for the current logged admin, record owner or record with manage access
func EnrichRecord(c echo.Context, dao *daos.Dao, record *models.Record, defaultExpands ...string) error {
	return EnrichRecords(c, dao, []*models.Record{record}, defaultExpands...)
}

// EnrichRecords parses the request context and enriches the provided records:
// - expands relations (if defaultExpands and/or ?expand query param is set)
// - ensures that the emails of the auth records and their expanded auth relations
//   are visibe only for the current logged admin, record owner or record with manage access
func EnrichRecords(c echo.Context, dao *daos.Dao, records []*models.Record, defaultExpands ...string) error {
	requestData := RequestData(c)

	if err := autoIgnoreAuthRecordsEmailVisibility(dao, records, requestData); err != nil {
		return fmt.Errorf("Failed to resolve email visibility: %w", err)
	}

	expands := defaultExpands
	expands = append(expands, strings.Split(c.QueryParam(expandQueryParam), ",")...)
	if len(expands) == 0 {
		return nil // nothing to expand
	}

	errs := dao.ExpandRecords(records, expands, expandFetch(dao, requestData))
	if len(errs) > 0 {
		return fmt.Errorf("Failed to expand: %v", errs)
	}

	return nil
}

// expandFetch is the records fetch function that is used to expand related records.
func expandFetch(
	dao *daos.Dao,
	requestData *models.RequestData,
) daos.ExpandFetchFunc {
	return func(relCollection *models.Collection, relIds []string) ([]*models.Record, error) {
		records, err := dao.FindRecordsByIds(relCollection.Id, relIds, func(q *dbx.SelectQuery) error {
			if requestData.Admin != nil {
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
			autoIgnoreAuthRecordsEmailVisibility(dao, records, requestData)
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
	requestData *models.RequestData,
) error {
	if len(records) == 0 || !records[0].Collection().IsAuth() {
		return nil // nothing to check
	}

	if requestData.Admin != nil {
		for _, rec := range records {
			rec.IgnoreEmailVisibility(true)
		}
		return nil
	}

	collection := records[0].Collection()

	mappedRecords := make(map[string]*models.Record, len(records))
	recordIds := make([]any, len(records))
	for i, rec := range records {
		mappedRecords[rec.Id] = rec
		recordIds[i] = rec.Id
	}

	if requestData != nil && requestData.AuthRecord != nil && mappedRecords[requestData.AuthRecord.Id] != nil {
		mappedRecords[requestData.AuthRecord.Id].IgnoreEmailVisibility(true)
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
	requestData *models.RequestData,
) bool {
	if !record.Collection().IsAuth() {
		return false
	}

	manageRule := record.Collection().AuthOptions().ManageRule

	if manageRule == nil || *manageRule == "" {
		return false // only for admins (manageRule can't be empty)
	}

	if requestData == nil || requestData.AuthRecord == nil {
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
