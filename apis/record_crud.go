package apis

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/search"
)

const expandQueryParam = "expand"

// bindRecordCrudApi registers the record crud api endpoints and
// the corresponding handlers.
func bindRecordCrudApi(app core.App, rg *echo.Group) {
	api := recordApi{app: app}

	subGroup := rg.Group(
		"/collections/:collection",
		ActivityLogger(app),
		LoadCollectionContext(app),
	)

	subGroup.GET("/records", api.list)
	subGroup.POST("/records", api.create)
	subGroup.GET("/records/:id", api.view)
	subGroup.PATCH("/records/:id", api.update)
	subGroup.DELETE("/records/:id", api.delete)
}

type recordApi struct {
	app core.App
}

func (api *recordApi) list(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("", "Missing collection context.")
	}

	// forbid users and guests to query special filter/sort fields
	if err := api.checkForForbiddenQueryFields(c); err != nil {
		return err
	}

	requestData := RequestData(c)

	if requestData.Admin == nil && collection.ListRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	fieldsResolver := resolvers.NewRecordFieldResolver(
		api.app.Dao(),
		collection,
		requestData,
		// hidden fields are searchable only by admins
		requestData.Admin != nil,
	)

	searchProvider := search.NewProvider(fieldsResolver).
		Query(api.app.Dao().RecordQuery(collection))

	if requestData.Admin == nil && collection.ListRule != nil {
		searchProvider.AddFilter(search.FilterData(*collection.ListRule))
	}

	var rawRecords = []dbx.NullStringMap{}
	result, err := searchProvider.ParseAndExec(c.QueryString(), &rawRecords)
	if err != nil {
		return NewBadRequestError("Invalid filter parameters.", err)
	}

	records := models.NewRecordsFromNullStringMaps(collection, rawRecords)

	result.Items = records

	event := &core.RecordsListEvent{
		HttpContext: c,
		Collection:  collection,
		Records:     records,
		Result:      result,
	}

	return api.app.OnRecordsListRequest().Trigger(event, func(e *core.RecordsListEvent) error {
		if err := EnrichRecords(e.HttpContext, api.app.Dao(), e.Records); err != nil && api.app.IsDebug() {
			log.Println(err)
		}

		return e.HttpContext.JSON(http.StatusOK, e.Result)
	})
}

func (api *recordApi) view(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("", "Missing collection context.")
	}

	recordId := c.PathParam("id")
	if recordId == "" {
		return NewNotFoundError("", nil)
	}

	requestData := RequestData(c)

	if requestData.Admin == nil && collection.ViewRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		if requestData.Admin == nil && collection.ViewRule != nil && *collection.ViewRule != "" {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestData, true)
			expr, err := search.FilterData(*collection.ViewRule).BuildExpr(resolver)
			if err != nil {
				return err
			}
			resolver.UpdateQuery(q)
			q.AndWhere(expr)
		}
		return nil
	}

	record, fetchErr := api.app.Dao().FindRecordById(collection.Id, recordId, ruleFunc)
	if fetchErr != nil || record == nil {
		return NewNotFoundError("", fetchErr)
	}

	event := &core.RecordViewEvent{
		HttpContext: c,
		Record:      record,
	}

	return api.app.OnRecordViewRequest().Trigger(event, func(e *core.RecordViewEvent) error {
		if err := EnrichRecord(e.HttpContext, api.app.Dao(), e.Record); err != nil && api.app.IsDebug() {
			log.Println(err)
		}

		return e.HttpContext.JSON(http.StatusOK, e.Record)
	})
}

func (api *recordApi) create(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("", "Missing collection context.")
	}

	requestData := RequestData(c)

	if requestData.Admin == nil && collection.CreateRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	hasFullManageAccess := requestData.Admin != nil

	// temporary save the record and check it against the create rule
	if requestData.Admin == nil && collection.CreateRule != nil {
		createRuleFunc := func(q *dbx.SelectQuery) error {
			if *collection.CreateRule == "" {
				return nil // no create rule to resolve
			}

			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestData, true)
			expr, err := search.FilterData(*collection.CreateRule).BuildExpr(resolver)
			if err != nil {
				return err
			}
			resolver.UpdateQuery(q)
			q.AndWhere(expr)
			return nil
		}

		testRecord := models.NewRecord(collection)
		testForm := forms.NewRecordUpsert(api.app, testRecord)
		testForm.SetFullManageAccess(true)
		if err := testForm.LoadRequest(c.Request(), ""); err != nil {
			return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
		}

		testErr := testForm.DrySubmit(func(txDao *daos.Dao) error {
			foundRecord, err := txDao.FindRecordById(collection.Id, testRecord.Id, createRuleFunc)
			if err != nil {
				return fmt.Errorf("DrySubmit create rule failure: %v", err)
			}
			hasFullManageAccess = hasAuthManageAccess(txDao, foundRecord, requestData)
			return nil
		})

		if testErr != nil {
			return NewBadRequestError("Failed to create record.", testErr)
		}
	}

	record := models.NewRecord(collection)
	form := forms.NewRecordUpsert(api.app, record)
	form.SetFullManageAccess(hasFullManageAccess)

	// load request
	if err := form.LoadRequest(c.Request(), ""); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.RecordCreateEvent{
		HttpContext: c,
		Record:      record,
	}

	// create the record
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnRecordBeforeCreateRequest().Trigger(event, func(e *core.RecordCreateEvent) error {
				if err := next(); err != nil {
					return NewBadRequestError("Failed to create record.", err)
				}

				if err := EnrichRecord(e.HttpContext, api.app.Dao(), e.Record); err != nil && api.app.IsDebug() {
					log.Println(err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.Record)
			})
		}
	})

	if submitErr == nil {
		api.app.OnRecordAfterCreateRequest().Trigger(event)
	}

	return submitErr
}

func (api *recordApi) update(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("", "Missing collection context.")
	}

	recordId := c.PathParam("id")
	if recordId == "" {
		return NewNotFoundError("", nil)
	}

	requestData := RequestData(c)

	if requestData.Admin == nil && collection.UpdateRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		if requestData.Admin == nil && collection.UpdateRule != nil && *collection.UpdateRule != "" {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestData, true)
			expr, err := search.FilterData(*collection.UpdateRule).BuildExpr(resolver)
			if err != nil {
				return err
			}
			resolver.UpdateQuery(q)
			q.AndWhere(expr)
		}
		return nil
	}

	// fetch record
	record, fetchErr := api.app.Dao().FindRecordById(collection.Id, recordId, ruleFunc)
	if fetchErr != nil || record == nil {
		return NewNotFoundError("", fetchErr)
	}

	form := forms.NewRecordUpsert(api.app, record)
	form.SetFullManageAccess(requestData.Admin != nil || hasAuthManageAccess(api.app.Dao(), record, requestData))

	// load request
	if err := form.LoadRequest(c.Request(), ""); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.RecordUpdateEvent{
		HttpContext: c,
		Record:      record,
	}

	// update the record
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnRecordBeforeUpdateRequest().Trigger(event, func(e *core.RecordUpdateEvent) error {
				if err := next(); err != nil {
					return NewBadRequestError("Failed to update record.", err)
				}

				if err := EnrichRecord(e.HttpContext, api.app.Dao(), e.Record); err != nil && api.app.IsDebug() {
					log.Println(err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.Record)
			})
		}
	})

	if submitErr == nil {
		api.app.OnRecordAfterUpdateRequest().Trigger(event)
	}

	return submitErr
}

func (api *recordApi) delete(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("", "Missing collection context.")
	}

	recordId := c.PathParam("id")
	if recordId == "" {
		return NewNotFoundError("", nil)
	}

	requestData := RequestData(c)

	if requestData.Admin == nil && collection.DeleteRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		if requestData.Admin == nil && collection.DeleteRule != nil && *collection.DeleteRule != "" {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestData, true)
			expr, err := search.FilterData(*collection.DeleteRule).BuildExpr(resolver)
			if err != nil {
				return err
			}
			resolver.UpdateQuery(q)
			q.AndWhere(expr)
		}
		return nil
	}

	record, fetchErr := api.app.Dao().FindRecordById(collection.Id, recordId, ruleFunc)
	if fetchErr != nil || record == nil {
		return NewNotFoundError("", fetchErr)
	}

	event := &core.RecordDeleteEvent{
		HttpContext: c,
		Record:      record,
	}

	handlerErr := api.app.OnRecordBeforeDeleteRequest().Trigger(event, func(e *core.RecordDeleteEvent) error {
		// delete the record
		if err := api.app.Dao().DeleteRecord(e.Record); err != nil {
			return NewBadRequestError("Failed to delete record. Make sure that the record is not part of a required relation reference.", err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnRecordAfterDeleteRequest().Trigger(event)
	}

	return handlerErr
}

func (api *recordApi) checkForForbiddenQueryFields(c echo.Context) error {
	admin, _ := c.Get(ContextAdminKey).(*models.Admin)
	if admin != nil {
		return nil // admins are allowed to query everything
	}

	decodedQuery := c.QueryParam(search.FilterQueryParam) + c.QueryParam(search.SortQueryParam)
	forbiddenFields := []string{"@collection.", "@request."}

	for _, field := range forbiddenFields {
		if strings.Contains(decodedQuery, field) {
			return NewForbiddenError("Only admins can filter by @collection and @request query params", nil)
		}
	}

	return nil
}
