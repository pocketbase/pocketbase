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
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/search"
)

const expandQueryParam = "expand"

// BindRecordApi registers the record api endpoints and the corresponding handlers.
func BindRecordApi(app core.App, rg *echo.Group) {
	api := recordApi{app: app}

	subGroup := rg.Group(
		"/collections/:collection/records",
		ActivityLogger(app),
		LoadCollectionContext(app),
	)

	subGroup.GET("", api.list)
	subGroup.POST("", api.create)
	subGroup.GET("/:id", api.view)
	subGroup.PATCH("/:id", api.update)
	subGroup.DELETE("/:id", api.delete)
}

type recordApi struct {
	app core.App
}

func (api *recordApi) list(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return rest.NewNotFoundError("", "Missing collection context.")
	}

	admin, _ := c.Get(ContextAdminKey).(*models.Admin)
	if admin == nil && collection.ListRule == nil {
		// only admins can access if the rule is nil
		return rest.NewForbiddenError("Only admins can perform this action.", nil)
	}

	// forbid users and guests to query special filter/sort fields
	if err := api.checkForForbiddenQueryFields(c); err != nil {
		return err
	}

	requestData := api.exportRequestData(c)

	fieldsResolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestData)

	searchProvider := search.NewProvider(fieldsResolver).
		Query(api.app.Dao().RecordQuery(collection)).
		CountColumn(fmt.Sprintf("%s.id", api.app.Dao().DB().QuoteSimpleColumnName(collection.Name)))

	if admin == nil && collection.ListRule != nil {
		searchProvider.AddFilter(search.FilterData(*collection.ListRule))
	}

	var rawRecords = []dbx.NullStringMap{}
	result, err := searchProvider.ParseAndExec(c.QueryString(), &rawRecords)
	if err != nil {
		return rest.NewBadRequestError("Invalid filter parameters.", err)
	}

	records := models.NewRecordsFromNullStringMaps(collection, rawRecords)

	// expand records relations
	expands := strings.Split(c.QueryParam(expandQueryParam), ",")
	if len(expands) > 0 {
		failed := api.app.Dao().ExpandRecords(
			records,
			expands,
			api.expandFunc(c, requestData),
		)
		if len(failed) > 0 && api.app.IsDebug() {
			log.Println("Failed to expand relations: ", failed)
		}
	}

	result.Items = records

	event := &core.RecordsListEvent{
		HttpContext: c,
		Collection:  collection,
		Records:     records,
		Result:      result,
	}

	return api.app.OnRecordsListRequest().Trigger(event, func(e *core.RecordsListEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.Result)
	})
}

func (api *recordApi) view(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return rest.NewNotFoundError("", "Missing collection context.")
	}

	admin, _ := c.Get(ContextAdminKey).(*models.Admin)
	if admin == nil && collection.ViewRule == nil {
		// only admins can access if the rule is nil
		return rest.NewForbiddenError("Only admins can perform this action.", nil)
	}

	recordId := c.PathParam("id")
	if recordId == "" {
		return rest.NewNotFoundError("", nil)
	}

	requestData := api.exportRequestData(c)

	ruleFunc := func(q *dbx.SelectQuery) error {
		if admin == nil && collection.ViewRule != nil && *collection.ViewRule != "" {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestData)
			expr, err := search.FilterData(*collection.ViewRule).BuildExpr(resolver)
			if err != nil {
				return err
			}
			resolver.UpdateQuery(q)
			q.AndWhere(expr)
		}
		return nil
	}

	record, fetchErr := api.app.Dao().FindRecordById(collection, recordId, ruleFunc)
	if fetchErr != nil || record == nil {
		return rest.NewNotFoundError("", fetchErr)
	}

	// expand record relations
	failed := api.app.Dao().ExpandRecord(
		record,
		strings.Split(c.QueryParam(expandQueryParam), ","),
		api.expandFunc(c, requestData),
	)
	if len(failed) > 0 && api.app.IsDebug() {
		log.Println("Failed to expand relations: ", failed)
	}

	event := &core.RecordViewEvent{
		HttpContext: c,
		Record:      record,
	}

	return api.app.OnRecordViewRequest().Trigger(event, func(e *core.RecordViewEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.Record)
	})
}

func (api *recordApi) create(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return rest.NewNotFoundError("", "Missing collection context.")
	}

	admin, _ := c.Get(ContextAdminKey).(*models.Admin)
	if admin == nil && collection.CreateRule == nil {
		// only admins can access if the rule is nil
		return rest.NewForbiddenError("Only admins can perform this action.", nil)
	}

	requestData := api.exportRequestData(c)

	// temporary save the record and check it against the create rule
	if admin == nil && collection.CreateRule != nil && *collection.CreateRule != "" {
		ruleFunc := func(q *dbx.SelectQuery) error {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestData)
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
		if err := testForm.LoadData(c.Request()); err != nil {
			return rest.NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
		}

		testErr := testForm.DrySubmit(func(txDao *daos.Dao) error {
			_, fetchErr := txDao.FindRecordById(collection, testRecord.Id, ruleFunc)
			return fetchErr
		})
		if testErr != nil {
			return rest.NewBadRequestError("Failed to create record.", testErr)
		}
	}

	record := models.NewRecord(collection)
	form := forms.NewRecordUpsert(api.app, record)

	// load request
	if err := form.LoadData(c.Request()); err != nil {
		return rest.NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
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
					return rest.NewBadRequestError("Failed to create record.", err)
				}

				// expand record relations
				failed := api.app.Dao().ExpandRecord(
					e.Record,
					strings.Split(e.HttpContext.QueryParam(expandQueryParam), ","),
					api.expandFunc(e.HttpContext, requestData),
				)
				if len(failed) > 0 && api.app.IsDebug() {
					log.Println("Failed to expand relations: ", failed)
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
		return rest.NewNotFoundError("", "Missing collection context.")
	}

	admin, _ := c.Get(ContextAdminKey).(*models.Admin)
	if admin == nil && collection.UpdateRule == nil {
		// only admins can access if the rule is nil
		return rest.NewForbiddenError("Only admins can perform this action.", nil)
	}

	recordId := c.PathParam("id")
	if recordId == "" {
		return rest.NewNotFoundError("", nil)
	}

	requestData := api.exportRequestData(c)

	ruleFunc := func(q *dbx.SelectQuery) error {
		if admin == nil && collection.UpdateRule != nil && *collection.UpdateRule != "" {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestData)
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
	record, fetchErr := api.app.Dao().FindRecordById(collection, recordId, ruleFunc)
	if fetchErr != nil || record == nil {
		return rest.NewNotFoundError("", fetchErr)
	}

	form := forms.NewRecordUpsert(api.app, record)

	// load request
	if err := form.LoadData(c.Request()); err != nil {
		return rest.NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
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
					return rest.NewBadRequestError("Failed to update record.", err)
				}

				// expand record relations
				failed := api.app.Dao().ExpandRecord(
					e.Record,
					strings.Split(e.HttpContext.QueryParam(expandQueryParam), ","),
					api.expandFunc(e.HttpContext, requestData),
				)
				if len(failed) > 0 && api.app.IsDebug() {
					log.Println("Failed to expand relations: ", failed)
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
		return rest.NewNotFoundError("", "Missing collection context.")
	}

	admin, _ := c.Get(ContextAdminKey).(*models.Admin)
	if admin == nil && collection.DeleteRule == nil {
		// only admins can access if the rule is nil
		return rest.NewForbiddenError("Only admins can perform this action.", nil)
	}

	recordId := c.PathParam("id")
	if recordId == "" {
		return rest.NewNotFoundError("", nil)
	}

	requestData := api.exportRequestData(c)

	ruleFunc := func(q *dbx.SelectQuery) error {
		if admin == nil && collection.DeleteRule != nil && *collection.DeleteRule != "" {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestData)
			expr, err := search.FilterData(*collection.DeleteRule).BuildExpr(resolver)
			if err != nil {
				return err
			}
			resolver.UpdateQuery(q)
			q.AndWhere(expr)
		}
		return nil
	}

	record, fetchErr := api.app.Dao().FindRecordById(collection, recordId, ruleFunc)
	if fetchErr != nil || record == nil {
		return rest.NewNotFoundError("", fetchErr)
	}

	event := &core.RecordDeleteEvent{
		HttpContext: c,
		Record:      record,
	}

	handlerErr := api.app.OnRecordBeforeDeleteRequest().Trigger(event, func(e *core.RecordDeleteEvent) error {
		// delete the record
		if err := api.app.Dao().DeleteRecord(e.Record); err != nil {
			return rest.NewBadRequestError("Failed to delete record. Make sure that the record is not part of a required relation reference.", err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnRecordAfterDeleteRequest().Trigger(event)
	}

	return handlerErr
}

func (api *recordApi) exportRequestData(c echo.Context) map[string]any {
	result := map[string]any{}
	queryParams := map[string]any{}
	bodyData := map[string]any{}
	method := c.Request().Method

	echo.BindQueryParams(c, &queryParams)

	rest.BindBody(c, &bodyData)

	result["method"] = method
	result["query"] = queryParams
	result["data"] = bodyData
	result["user"] = nil

	loggedUser, _ := c.Get(ContextUserKey).(*models.User)
	if loggedUser != nil {
		result["user"], _ = loggedUser.AsMap()
	}

	return result
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
			return rest.NewForbiddenError("Only admins can filter by @collection and @request query params", nil)
		}
	}

	return nil
}

func (api *recordApi) expandFunc(c echo.Context, requestData map[string]any) daos.ExpandFetchFunc {
	admin, _ := c.Get(ContextAdminKey).(*models.Admin)

	return func(relCollection *models.Collection, relIds []string) ([]*models.Record, error) {
		return api.app.Dao().FindRecordsByIds(relCollection, relIds, func(q *dbx.SelectQuery) error {
			if admin != nil {
				return nil // admin can access everything
			}

			if relCollection.ViewRule == nil {
				return fmt.Errorf("Only admins can view collection %q records", relCollection.Name)
			}

			if *relCollection.ViewRule != "" {
				resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), relCollection, requestData)
				expr, err := search.FilterData(*(relCollection.ViewRule)).BuildExpr(resolver)
				if err != nil {
					return err
				}
				resolver.UpdateQuery(q)
				q.AndWhere(expr)
			}

			return nil
		})
	}
}
