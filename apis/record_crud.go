package apis

import (
	cryptoRand "crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/search"
)

// bindRecordCrudApi registers the record crud api endpoints and
// the corresponding handlers.
func bindRecordCrudApi(app core.App, rg *echo.Group) {
	api := recordApi{app: app}

	subGroup := rg.Group(
		"/collections/:collection",
		ActivityLogger(app),
	)

	subGroup.GET("/records", api.list, LoadCollectionContext(app))
	subGroup.GET("/records/:id", api.view, LoadCollectionContext(app))
	subGroup.POST("/records", api.create, LoadCollectionContext(app, models.CollectionTypeBase, models.CollectionTypeAuth))
	subGroup.PATCH("/records/:id", api.update, LoadCollectionContext(app, models.CollectionTypeBase, models.CollectionTypeAuth))
	subGroup.DELETE("/records/:id", api.delete, LoadCollectionContext(app, models.CollectionTypeBase, models.CollectionTypeAuth))
}

type recordApi struct {
	app core.App
}

func (api *recordApi) list(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("", "Missing collection context.")
	}

	requestInfo := RequestInfo(c)

	// forbid users and guests to query special filter/sort fields
	if err := checkForAdminOnlyRuleFields(requestInfo); err != nil {
		return err
	}

	if requestInfo.Admin == nil && collection.ListRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	fieldsResolver := resolvers.NewRecordFieldResolver(
		api.app.Dao(),
		collection,
		requestInfo,
		// hidden fields are searchable only by admins
		requestInfo.Admin != nil,
	)

	searchProvider := search.NewProvider(fieldsResolver).
		Query(api.app.Dao().RecordQuery(collection))

	if requestInfo.Admin == nil && collection.ListRule != nil {
		searchProvider.AddFilter(search.FilterData(*collection.ListRule))
	}

	records := []*models.Record{}

	// note: in v0.23.0 this has been migrated as option check in the search.Provider
	queryStr := c.QueryParams().Encode()
	if len(queryStr) > 2048 {
		return NewBadRequestError("query string is too large", nil)
	}

	result, err := searchProvider.ParseAndExec(queryStr, &records)
	if err != nil {
		return NewBadRequestError("", err)
	}

	event := new(core.RecordsListEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Records = records
	event.Result = result

	return api.app.OnRecordsListRequest().Trigger(event, func(e *core.RecordsListEvent) error {
		if e.HttpContext.Response().Committed {
			return nil
		}

		if err := EnrichRecords(e.HttpContext, api.app.Dao(), e.Records); err != nil {
			api.app.Logger().Debug("Failed to enrich list records", slog.String("error", err.Error()))
		}

		// note: in v0.23.0 this is combined with extra check for repeated attempts
		//
		// Add a randomized throttle in case of empty search filter attempts.
		//
		// This is just for extra precaution since security researches raised concern regarding the possibity of eventual
		// timing attacks because the List API rule acts also as filter and executes in a single run with the client-side filters.
		// This is by design and it is an accepted tradeoff between performance, usability and correctness.
		//
		// While technically the below doesn't fully guarantee protection against filter timing attacks, in practice combined with the network latency it makes them even less feasible.
		// A properly configured rate limiter or individual fields Hidden checks are better suited if you are really concerned about eventual information disclosure by side-channel attacks.
		//
		// In all cases it doesn't really matter that much because it doesn't affect the builtin PocketBase security sensitive fields (e.g. password and tokenKey) since they
		// are not client-side filterable and in the few places where they need to be compared against an external value, a constant time check is used.
		if requestInfo.Admin == nil &&
			(collection.ListRule != nil && *collection.ListRule != "") &&
			(requestInfo.Query["filter"] != "") &&
			len(e.Records) == 0 {
			api.app.Logger().Debug("Randomized throttle because of failed filter search", "collectionId", collection.Id)
			randomizedThrottle(100)
		}

		return e.HttpContext.JSON(http.StatusOK, e.Result)
	})
}

func randomizedThrottle(softMax int64) {
	var timeout int64
	randRange, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(softMax))
	if err == nil {
		timeout = randRange.Int64()
	} else {
		timeout = softMax
	}

	time.Sleep(time.Duration(timeout) * time.Millisecond)
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

	requestInfo := RequestInfo(c)

	if requestInfo.Admin == nil && collection.ViewRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		if requestInfo.Admin == nil && collection.ViewRule != nil && *collection.ViewRule != "" {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestInfo, true)
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

	event := new(core.RecordViewEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Record = record

	return api.app.OnRecordViewRequest().Trigger(event, func(e *core.RecordViewEvent) error {
		if e.HttpContext.Response().Committed {
			return nil
		}

		if err := EnrichRecord(e.HttpContext, api.app.Dao(), e.Record); err != nil {
			api.app.Logger().Debug(
				"Failed to enrich view record",
				slog.String("id", e.Record.Id),
				slog.String("collectionName", e.Record.Collection().Name),
				slog.String("error", err.Error()),
			)
		}

		return e.HttpContext.JSON(http.StatusOK, e.Record)
	})
}

func (api *recordApi) create(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("", "Missing collection context.")
	}

	requestInfo := RequestInfo(c)

	if requestInfo.Admin == nil && collection.CreateRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	hasFullManageAccess := requestInfo.Admin != nil

	// temporary save the record and check it against the create rule
	if requestInfo.Admin == nil && collection.CreateRule != nil {
		testRecord := models.NewRecord(collection)

		// replace modifiers fields so that the resolved value is always
		// available when accessing requestInfo.Data using just the field name
		if requestInfo.HasModifierDataKeys() {
			requestInfo.Data = testRecord.ReplaceModifers(requestInfo.Data)
		}

		testForm := forms.NewRecordUpsert(api.app, testRecord)
		testForm.SetFullManageAccess(true)
		if err := testForm.LoadRequest(c.Request(), ""); err != nil {
			return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
		}

		// force unset the verified state to prevent ManageRule misuse
		if !hasFullManageAccess {
			testForm.Verified = false
		}

		createRuleFunc := func(q *dbx.SelectQuery) error {
			if *collection.CreateRule == "" {
				return nil // no create rule to resolve
			}

			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestInfo, true)
			expr, err := search.FilterData(*collection.CreateRule).BuildExpr(resolver)
			if err != nil {
				return err
			}
			resolver.UpdateQuery(q)
			q.AndWhere(expr)
			return nil
		}

		testErr := testForm.DrySubmit(func(txDao *daos.Dao) error {
			foundRecord, err := txDao.FindRecordById(collection.Id, testRecord.Id, createRuleFunc)
			if err != nil {
				return fmt.Errorf("DrySubmit create rule failure: %w", err)
			}
			hasFullManageAccess = hasAuthManageAccess(txDao, foundRecord, requestInfo)
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

	event := new(core.RecordCreateEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Record = record
	event.UploadedFiles = form.FilesToUpload()

	// create the record
	return form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(m *models.Record) error {
			event.Record = m

			return api.app.OnRecordBeforeCreateRequest().Trigger(event, func(e *core.RecordCreateEvent) error {
				if err := next(e.Record); err != nil {
					return NewBadRequestError("Failed to create record.", err)
				}

				if err := EnrichRecord(e.HttpContext, api.app.Dao(), e.Record); err != nil {
					api.app.Logger().Debug(
						"Failed to enrich create record",
						slog.String("id", e.Record.Id),
						slog.String("collectionName", e.Record.Collection().Name),
						slog.String("error", err.Error()),
					)
				}

				return api.app.OnRecordAfterCreateRequest().Trigger(event, func(e *core.RecordCreateEvent) error {
					if e.HttpContext.Response().Committed {
						return nil
					}

					return e.HttpContext.JSON(http.StatusOK, e.Record)
				})
			})
		}
	})
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

	requestInfo := RequestInfo(c)

	if requestInfo.Admin == nil && collection.UpdateRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	// eager fetch the record so that the modifier field values are replaced
	// and available when accessing requestInfo.Data using just the field name
	if requestInfo.HasModifierDataKeys() {
		record, err := api.app.Dao().FindRecordById(collection.Id, recordId)
		if err != nil || record == nil {
			return NewNotFoundError("", err)
		}
		requestInfo.Data = record.ReplaceModifers(requestInfo.Data)
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		if requestInfo.Admin == nil && collection.UpdateRule != nil && *collection.UpdateRule != "" {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestInfo, true)
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
	form.SetFullManageAccess(requestInfo.Admin != nil || hasAuthManageAccess(api.app.Dao(), record, requestInfo))

	// load request
	if err := form.LoadRequest(c.Request(), ""); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := new(core.RecordUpdateEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Record = record
	event.UploadedFiles = form.FilesToUpload()

	// update the record
	return form.Submit(func(next forms.InterceptorNextFunc[*models.Record]) forms.InterceptorNextFunc[*models.Record] {
		return func(m *models.Record) error {
			event.Record = m

			return api.app.OnRecordBeforeUpdateRequest().Trigger(event, func(e *core.RecordUpdateEvent) error {
				if err := next(e.Record); err != nil {
					return NewBadRequestError("Failed to update record.", err)
				}

				if err := EnrichRecord(e.HttpContext, api.app.Dao(), e.Record); err != nil {
					api.app.Logger().Debug(
						"Failed to enrich update record",
						slog.String("id", e.Record.Id),
						slog.String("collectionName", e.Record.Collection().Name),
						slog.String("error", err.Error()),
					)
				}

				return api.app.OnRecordAfterUpdateRequest().Trigger(event, func(e *core.RecordUpdateEvent) error {
					if e.HttpContext.Response().Committed {
						return nil
					}

					return e.HttpContext.JSON(http.StatusOK, e.Record)
				})
			})
		}
	})
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

	requestInfo := RequestInfo(c)

	if requestInfo.Admin == nil && collection.DeleteRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		if requestInfo.Admin == nil && collection.DeleteRule != nil && *collection.DeleteRule != "" {
			resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), collection, requestInfo, true)
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

	event := new(core.RecordDeleteEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Record = record

	return api.app.OnRecordBeforeDeleteRequest().Trigger(event, func(e *core.RecordDeleteEvent) error {
		// delete the record
		if err := api.app.Dao().DeleteRecord(e.Record); err != nil {
			return NewBadRequestError("Failed to delete record. Make sure that the record is not part of a required relation reference.", err)
		}

		return api.app.OnRecordAfterDeleteRequest().Trigger(event, func(e *core.RecordDeleteEvent) error {
			if e.HttpContext.Response().Committed {
				return nil
			}

			return e.HttpContext.NoContent(http.StatusNoContent)
		})
	})
}
