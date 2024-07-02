package apis

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/search"
)

const ContextRequestInfoKey = "requestInfo"

const expandQueryParam = "expand"
const fieldsQueryParam = "fields"

// Deprecated: Use RequestInfo instead.
func RequestData(c echo.Context) *models.RequestInfo {
	log.Println("RequestData(c) is deprecated and will be removed in the future! You can replace it with RequestInfo(c).")
	return RequestInfo(c)
}

// RequestInfo exports cached common request data fields
// (query, body, logged auth state, etc.) from the provided context.
func RequestInfo(c echo.Context) *models.RequestInfo {
	// return cached to avoid copying the body multiple times
	if v := c.Get(ContextRequestInfoKey); v != nil {
		if data, ok := v.(*models.RequestInfo); ok {
			// refresh auth state
			data.AuthRecord, _ = c.Get(ContextAuthRecordKey).(*models.Record)
			data.Admin, _ = c.Get(ContextAdminKey).(*models.Admin)
			return data
		}
	}

	result := &models.RequestInfo{
		Context: models.RequestInfoContextDefault,
		Method:  c.Request().Method,
		Query:   map[string]any{},
		Data:    map[string]any{},
		Headers: map[string]any{},
	}

	// extract the first value of all headers and normalizes the keys
	// ("X-Token" is converted to "x_token")
	for k, v := range c.Request().Header {
		if len(v) > 0 {
			result.Headers[inflector.Snakecase(k)] = v[0]
		}
	}

	result.AuthRecord, _ = c.Get(ContextAuthRecordKey).(*models.Record)
	result.Admin, _ = c.Get(ContextAdminKey).(*models.Admin)
	echo.BindQueryParams(c, &result.Query)
	rest.BindBody(c, &result.Data)

	c.Set(ContextRequestInfoKey, result)

	return result
}

// RecordAuthResponse writes standardised json record auth response
// into the specified request context.
func RecordAuthResponse(
	app core.App,
	c echo.Context,
	authRecord *models.Record,
	meta any,
	finalizers ...func(token string) error,
) error {
	if !authRecord.Verified() && authRecord.Collection().AuthOptions().OnlyVerified {
		return NewForbiddenError("Please verify your account first.", nil)
	}

	token, tokenErr := tokens.NewRecordAuthToken(app, authRecord)
	if tokenErr != nil {
		return NewBadRequestError("Failed to create auth token.", tokenErr)
	}

	event := new(core.RecordAuthEvent)
	event.HttpContext = c
	event.Collection = authRecord.Collection()
	event.Record = authRecord
	event.Token = token
	event.Meta = meta

	return app.OnRecordAuthRequest().Trigger(event, func(e *core.RecordAuthEvent) error {
		if e.HttpContext.Response().Committed {
			return nil
		}

		// allow always returning the email address of the authenticated account
		e.Record.IgnoreEmailVisibility(true)

		// expand record relations
		expands := strings.Split(c.QueryParam(expandQueryParam), ",")
		if len(expands) > 0 {
			// create a copy of the cached request data and adjust it to the current auth record
			requestInfo := *RequestInfo(e.HttpContext)
			requestInfo.Admin = nil
			requestInfo.AuthRecord = e.Record
			failed := app.Dao().ExpandRecord(
				e.Record,
				expands,
				expandFetch(app.Dao(), &requestInfo),
			)
			if len(failed) > 0 {
				app.Logger().Debug("[RecordAuthResponse] Failed to expand relations", slog.Any("errors", failed))
			}
		}

		result := map[string]any{
			"token":  e.Token,
			"record": e.Record,
		}

		if e.Meta != nil {
			result["meta"] = e.Meta
		}

		for _, f := range finalizers {
			if err := f(e.Token); err != nil {
				return err
			}
		}

		return e.HttpContext.JSON(http.StatusOK, result)
	})
}

// EnrichRecord parses the request context and enrich the provided record:
//   - expands relations (if defaultExpands and/or ?expand query param is set)
//   - ensures that the emails of the auth record and its expanded auth relations
//     are visible only for the current logged admin, record owner or record with manage access
func EnrichRecord(c echo.Context, dao *daos.Dao, record *models.Record, defaultExpands ...string) error {
	return EnrichRecords(c, dao, []*models.Record{record}, defaultExpands...)
}

// EnrichRecords parses the request context and enriches the provided records:
//   - expands relations (if defaultExpands and/or ?expand query param is set)
//   - ensures that the emails of the auth records and their expanded auth relations
//     are visible only for the current logged admin, record owner or record with manage access
func EnrichRecords(c echo.Context, dao *daos.Dao, records []*models.Record, defaultExpands ...string) error {
	requestInfo := RequestInfo(c)

	if err := autoIgnoreAuthRecordsEmailVisibility(dao, records, requestInfo); err != nil {
		return fmt.Errorf("failed to resolve email visibility: %w", err)
	}

	expands := defaultExpands
	if param := c.QueryParam(expandQueryParam); param != "" {
		expands = append(expands, strings.Split(param, ",")...)
	}
	if len(expands) == 0 {
		return nil // nothing to expand
	}

	errs := dao.ExpandRecords(records, expands, expandFetch(dao, requestInfo))
	if len(errs) > 0 {
		return fmt.Errorf("failed to expand: %v", errs)
	}

	return nil
}

// expandFetch is the records fetch function that is used to expand related records.
func expandFetch(
	dao *daos.Dao,
	requestInfo *models.RequestInfo,
) daos.ExpandFetchFunc {
	return func(relCollection *models.Collection, relIds []string) ([]*models.Record, error) {
		records, err := dao.FindRecordsByIds(relCollection.Id, relIds, func(q *dbx.SelectQuery) error {
			if requestInfo.Admin != nil {
				return nil // admins can access everything
			}

			if relCollection.ViewRule == nil {
				return fmt.Errorf("only admins can view collection %q records", relCollection.Name)
			}

			if *relCollection.ViewRule != "" {
				resolver := resolvers.NewRecordFieldResolver(dao, relCollection, requestInfo, true)
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
			autoIgnoreAuthRecordsEmailVisibility(dao, records, requestInfo)
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
	requestInfo *models.RequestInfo,
) error {
	if len(records) == 0 || !records[0].Collection().IsAuth() {
		return nil // nothing to check
	}

	if requestInfo.Admin != nil {
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

	if requestInfo != nil && requestInfo.AuthRecord != nil && mappedRecords[requestInfo.AuthRecord.Id] != nil {
		mappedRecords[requestInfo.AuthRecord.Id].IgnoreEmailVisibility(true)
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

	resolver := resolvers.NewRecordFieldResolver(dao, collection, requestInfo, true)
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
	requestInfo *models.RequestInfo,
) bool {
	if !record.Collection().IsAuth() {
		return false
	}

	manageRule := record.Collection().AuthOptions().ManageRule

	if manageRule == nil || *manageRule == "" {
		return false // only for admins (manageRule can't be empty)
	}

	if requestInfo == nil || requestInfo.AuthRecord == nil {
		return false // no auth record
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		resolver := resolvers.NewRecordFieldResolver(dao, record.Collection(), requestInfo, true)
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

var ruleQueryParams = []string{search.FilterQueryParam, search.SortQueryParam}
var adminOnlyRuleFields = []string{"@collection.", "@request."}

// @todo consider moving the rules check to the RecordFieldResolver.
//
// checkForAdminOnlyRuleFields loosely checks and returns an error if
// the provided RequestInfo contains rule fields that only the admin can use.
func checkForAdminOnlyRuleFields(requestInfo *models.RequestInfo) error {
	if requestInfo.Admin != nil || len(requestInfo.Query) == 0 {
		return nil // admin or nothing to check
	}

	for _, param := range ruleQueryParams {
		v, _ := requestInfo.Query[param].(string)
		if v == "" {
			continue
		}

		for _, field := range adminOnlyRuleFields {
			if strings.Contains(v, field) {
				return NewForbiddenError("Only admins can filter by "+field, nil)
			}
		}
	}

	return nil
}
