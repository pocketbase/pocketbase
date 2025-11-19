package apis

import (
	cryptoRand "crypto/rand"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
)

// bindRecordCrudApi registers the record crud api endpoints and
// the corresponding handlers.
//
// note: the rate limiter is "inlined" because some of the crud actions are also used in the batch APIs
func bindRecordCrudApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/collections/{collection}/records").Unbind(DefaultRateLimitMiddlewareId)
	subGroup.GET("", recordsList)
	subGroup.GET("/{id}", recordView)
	subGroup.POST("", recordCreate(true, nil)).Bind(dynamicCollectionBodyLimit(""))
	subGroup.PATCH("/{id}", recordUpdate(true, nil)).Bind(dynamicCollectionBodyLimit(""))
	subGroup.DELETE("/{id}", recordDelete(true, nil))
	// CSV导出和导入API
	subGroup.GET("/export", recordsExport)
	subGroup.POST("/import", recordsImport)
}

func recordsList(e *core.RequestEvent) error {
	collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
	if err != nil || collection == nil {
		return e.NotFoundError("Missing collection context.", err)
	}

	err = checkCollectionRateLimit(e, collection, "list")
	if err != nil {
		return err
	}

	requestInfo, err := e.RequestInfo()
	if err != nil {
		return firstApiError(err, e.BadRequestError("", err))
	}

	if collection.ListRule == nil && !requestInfo.HasSuperuserAuth() {
		return e.ForbiddenError("Only superusers can perform this action.", nil)
	}

	// forbid users and guests to query special filter/sort fields
	err = checkForSuperuserOnlyRuleFields(requestInfo)
	if err != nil {
		return err
	}

	query := e.App.RecordQuery(collection)

	fieldsResolver := core.NewRecordFieldResolver(e.App, collection, requestInfo, true)

	if !requestInfo.HasSuperuserAuth() && collection.ListRule != nil && *collection.ListRule != "" {
		expr, err := search.FilterData(*collection.ListRule).BuildExpr(fieldsResolver)
		if err != nil {
			return err
		}
		query.AndWhere(expr)

		// will be applied by the search provider right before executing the query
		// fieldsResolver.UpdateQuery(query)
	}

	// hidden fields are searchable only by superusers
	fieldsResolver.SetAllowHiddenFields(requestInfo.HasSuperuserAuth())

	searchProvider := search.NewProvider(fieldsResolver).Query(query)

	// use rowid when available to minimize the need of a covering index with the "id" field
	if !collection.IsView() {
		searchProvider.CountCol("_rowid_")
	}

	records := []*core.Record{}
	result, err := searchProvider.ParseAndExec(e.Request.URL.Query().Encode(), &records)

	if err != nil {
		return firstApiError(err, e.BadRequestError("", err))
	}

	event := new(core.RecordsListRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Records = records
	event.Result = result

	return e.App.OnRecordsListRequest().Trigger(event, func(e *core.RecordsListRequestEvent) error {
		if err := EnrichRecords(e.RequestEvent, e.Records); err != nil {
			return firstApiError(err, e.InternalServerError("Failed to enrich records", err))
		}

		// Add a randomized throttle in case of too many empty search filter attempts.
		//
		// This is just for extra precaution since security researches raised concern regarding the possibility of eventual
		// timing attacks because the List API rule acts also as filter and executes in a single run with the client-side filters.
		// This is by design and it is an accepted trade off between performance, usability and correctness.
		//
		// While technically the below doesn't fully guarantee protection against filter timing attacks, in practice combined with the network latency it makes them even less feasible.
		// A properly configured rate limiter or individual fields Hidden checks are better suited if you are really concerned about eventual information disclosure by side-channel attacks.
		//
		// In all cases it doesn't really matter that much because it doesn't affect the builtin PocketBase security sensitive fields (e.g. password and tokenKey) since they
		// are not client-side filterable and in the few places where they need to be compared against an external value, a constant time check is used.
		if !e.HasSuperuserAuth() &&
			(collection.ListRule != nil && *collection.ListRule != "") &&
			(requestInfo.Query["filter"] != "") &&
			len(e.Records) == 0 &&
			checkRateLimit(e.RequestEvent, "@pb_list_timing_check_"+collection.Id, listTimingRateLimitRule) != nil {
			e.App.Logger().Debug("Randomized throttle because of too many failed searches", "collectionId", collection.Id)
			randomizedThrottle(500)
		}

		return execAfterSuccessTx(true, e.App, func() error {
			return e.JSON(http.StatusOK, e.Result)
		})
	})
}

var listTimingRateLimitRule = core.RateLimitRule{MaxRequests: 3, Duration: 3}

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

func recordView(e *core.RequestEvent) error {
	collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
	if err != nil || collection == nil {
		return e.NotFoundError("Missing collection context.", err)
	}

	err = checkCollectionRateLimit(e, collection, "view")
	if err != nil {
		return err
	}

	recordId := e.Request.PathValue("id")
	if recordId == "" {
		return e.NotFoundError("", nil)
	}

	requestInfo, err := e.RequestInfo()
	if err != nil {
		return firstApiError(err, e.BadRequestError("", err))
	}

	if collection.ViewRule == nil && !requestInfo.HasSuperuserAuth() {
		return e.ForbiddenError("Only superusers can perform this action.", nil)
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		if !requestInfo.HasSuperuserAuth() && collection.ViewRule != nil && *collection.ViewRule != "" {
			resolver := core.NewRecordFieldResolver(e.App, collection, requestInfo, true)

			expr, err := search.FilterData(*collection.ViewRule).BuildExpr(resolver)
			if err != nil {
				return err
			}

			q.AndWhere(expr)

			err = resolver.UpdateQuery(q)
			if err != nil {
				return err
			}
		}
		return nil
	}

	record, fetchErr := e.App.FindRecordById(collection, recordId, ruleFunc)
	if fetchErr != nil || record == nil {
		return firstApiError(err, e.NotFoundError("", fetchErr))
	}

	event := new(core.RecordRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Record = record

	return e.App.OnRecordViewRequest().Trigger(event, func(e *core.RecordRequestEvent) error {
		if err := EnrichRecord(e.RequestEvent, e.Record); err != nil {
			return firstApiError(err, e.InternalServerError("Failed to enrich record", err))
		}

		return execAfterSuccessTx(true, e.App, func() error {
			return e.JSON(http.StatusOK, e.Record)
		})
	})
}

func recordCreate(responseWriteAfterTx bool, optFinalizer func(data any) error) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
		if err != nil || collection == nil {
			return e.NotFoundError("Missing collection context.", err)
		}

		if collection.IsView() {
			return e.BadRequestError("Unsupported collection type.", nil)
		}

		err = checkCollectionRateLimit(e, collection, "create")
		if err != nil {
			return err
		}

		requestInfo, err := e.RequestInfo()
		if err != nil {
			return firstApiError(err, e.BadRequestError("", err))
		}

		hasSuperuserAuth := requestInfo.HasSuperuserAuth()
		if !hasSuperuserAuth && collection.CreateRule == nil {
			return e.ForbiddenError("Only superusers can perform this action.", nil)
		}

		record := core.NewRecord(collection)

		data, err := recordDataFromRequest(e, record)
		if err != nil {
			return firstApiError(err, e.BadRequestError("Failed to read the submitted data.", err))
		}

		// set a random password for the OAuth2 ignoring its plain password validators
		var skipPlainPasswordRecordValidators bool
		if requestInfo.Context == core.RequestInfoContextOAuth2 {
			if _, ok := data[core.FieldNamePassword]; !ok {
				data[core.FieldNamePassword] = security.RandomString(30)
				data[core.FieldNamePassword+"Confirm"] = data[core.FieldNamePassword]
				skipPlainPasswordRecordValidators = true
			}
		}

		// replace modifiers fields so that the resolved value is always
		// available when accessing requestInfo.Body
		requestInfo.Body = data

		form := forms.NewRecordUpsert(e.App, record)
		if hasSuperuserAuth {
			form.GrantSuperuserAccess()
		}
		form.Load(data)

		if skipPlainPasswordRecordValidators {
			// unset the plain value to skip the plain password field validators
			if raw, ok := record.GetRaw(core.FieldNamePassword).(*core.PasswordFieldValue); ok {
				raw.Plain = ""
			}
		}

		// check the request and record data against the create and manage rules
		if !hasSuperuserAuth && collection.CreateRule != nil {
			dummyRecord := record.Clone()

			dummyRandomPart := "__pb_create__" + security.PseudorandomString(6)

			// set an id if it doesn't have already
			// (the value doesn't matter; it is used only to minimize the breaking changes with earlier versions)
			if dummyRecord.Id == "" {
				dummyRecord.Id = "__temp_id__" + dummyRandomPart
			}

			// unset the verified field to prevent manage API rule misuse in case the rule relies on it
			dummyRecord.SetVerified(false)

			// export the dummy record data into db params
			dummyExport, err := dummyRecord.DBExport(e.App)
			if err != nil {
				return e.BadRequestError("Failed to create record", fmt.Errorf("dummy DBExport error: %w", err))
			}

			dummyParams := make(dbx.Params, len(dummyExport))
			selects := make([]string, 0, len(dummyExport))
			var param string
			for k, v := range dummyExport {
				k = inflector.Columnify(k) // columnify is just as extra measure in case of custom fields
				param = "__pb_create__" + k
				dummyParams[param] = v
				selects = append(selects, "{:"+param+"} AS [["+k+"]]")
			}

			// shallow clone the current collection
			dummyCollection := *collection
			dummyCollection.Id += dummyRandomPart
			dummyCollection.Name += inflector.Columnify(dummyRandomPart)

			withFrom := fmt.Sprintf("WITH {{%s}} as (SELECT %s)", dummyCollection.Name, strings.Join(selects, ","))

			// check non-empty create rule
			if *dummyCollection.CreateRule != "" {
				ruleQuery := e.App.ConcurrentDB().Select("(1)").PreFragment(withFrom).From(dummyCollection.Name).AndBind(dummyParams)

				resolver := core.NewRecordFieldResolver(e.App, &dummyCollection, requestInfo, true)

				expr, err := search.FilterData(*dummyCollection.CreateRule).BuildExpr(resolver)
				if err != nil {
					return e.BadRequestError("Failed to create record", fmt.Errorf("create rule build expression failure: %w", err))
				}
				ruleQuery.AndWhere(expr)

				err = resolver.UpdateQuery(ruleQuery)
				if err != nil {
					return e.BadRequestError("Failed to create record", fmt.Errorf("create rule update query failure: %w", err))
				}

				var exists int
				err = ruleQuery.Limit(1).Row(&exists)
				if err != nil || exists == 0 {
					return e.BadRequestError("Failed to create record", fmt.Errorf("create rule failure: %w", err))
				}
			}

			// check for manage rule access
			manageRuleQuery := e.App.ConcurrentDB().Select("(1)").PreFragment(withFrom).From(dummyCollection.Name).AndBind(dummyParams)
			if !form.HasManageAccess() &&
				hasAuthManageAccess(e.App, requestInfo, &dummyCollection, manageRuleQuery) {
				form.GrantManagerAccess()
			}
		}

		var isOptFinalizerCalled bool

		event := new(core.RecordRequestEvent)
		event.RequestEvent = e
		event.Collection = collection
		event.Record = record

		hookErr := e.App.OnRecordCreateRequest().Trigger(event, func(e *core.RecordRequestEvent) error {
			form.SetApp(e.App)
			form.SetRecord(e.Record)

			err := form.Submit()
			if err != nil {
				return firstApiError(err, e.BadRequestError("Failed to create record", err))
			}

			err = EnrichRecord(e.RequestEvent, e.Record)
			if err != nil {
				return firstApiError(err, e.InternalServerError("Failed to enrich record", err))
			}

			err = execAfterSuccessTx(responseWriteAfterTx, e.App, func() error {
				return e.JSON(http.StatusOK, e.Record)
			})
			if err != nil {
				return err
			}

			if optFinalizer != nil {
				isOptFinalizerCalled = true
				err = optFinalizer(e.Record)
				if err != nil {
					return firstApiError(err, e.InternalServerError("", err))
				}
			}

			return nil
		})
		if hookErr != nil {
			return hookErr
		}

		// e.g. in case the regular hook chain was stopped and the finalizer cannot be executed as part of the last e.Next() task
		if !isOptFinalizerCalled && optFinalizer != nil {
			if err := optFinalizer(event.Record); err != nil {
				return firstApiError(err, e.InternalServerError("", err))
			}
		}

		return nil
	}
}

func recordUpdate(responseWriteAfterTx bool, optFinalizer func(data any) error) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
		if err != nil || collection == nil {
			return e.NotFoundError("Missing collection context.", err)
		}

		if collection.IsView() {
			return e.BadRequestError("Unsupported collection type.", nil)
		}

		err = checkCollectionRateLimit(e, collection, "update")
		if err != nil {
			return err
		}

		recordId := e.Request.PathValue("id")
		if recordId == "" {
			return e.NotFoundError("", nil)
		}

		requestInfo, err := e.RequestInfo()
		if err != nil {
			return firstApiError(err, e.BadRequestError("", err))
		}

		hasSuperuserAuth := requestInfo.HasSuperuserAuth()

		if !hasSuperuserAuth && collection.UpdateRule == nil {
			return firstApiError(err, e.ForbiddenError("Only superusers can perform this action.", nil))
		}

		// eager fetch the record so that the modifiers field values can be resolved
		record, err := e.App.FindRecordById(collection, recordId)
		if err != nil {
			return firstApiError(err, e.NotFoundError("", err))
		}

		data, err := recordDataFromRequest(e, record)
		if err != nil {
			return firstApiError(err, e.BadRequestError("Failed to read the submitted data.", err))
		}

		// replace modifiers fields so that the resolved value is always
		// available when accessing requestInfo.Body
		requestInfo.Body = data

		ruleFunc := func(q *dbx.SelectQuery) error {
			if !hasSuperuserAuth && collection.UpdateRule != nil && *collection.UpdateRule != "" {
				resolver := core.NewRecordFieldResolver(e.App, collection, requestInfo, true)

				expr, err := search.FilterData(*collection.UpdateRule).BuildExpr(resolver)
				if err != nil {
					return err
				}

				q.AndWhere(expr)

				err = resolver.UpdateQuery(q)
				if err != nil {
					return err
				}
			}
			return nil
		}

		// refetch with access checks
		record, err = e.App.FindRecordById(collection, recordId, ruleFunc)
		if err != nil {
			return firstApiError(err, e.NotFoundError("", err))
		}

		form := forms.NewRecordUpsert(e.App, record)
		if hasSuperuserAuth {
			form.GrantSuperuserAccess()
		}
		form.Load(data)

		manageRuleQuery := e.App.ConcurrentDB().Select("(1)").From(collection.Name).AndWhere(dbx.HashExp{
			collection.Name + ".id": record.Id,
		})
		if !form.HasManageAccess() &&
			hasAuthManageAccess(e.App, requestInfo, collection, manageRuleQuery) {
			form.GrantManagerAccess()
		}

		var isOptFinalizerCalled bool

		event := new(core.RecordRequestEvent)
		event.RequestEvent = e
		event.Collection = collection
		event.Record = record

		hookErr := e.App.OnRecordUpdateRequest().Trigger(event, func(e *core.RecordRequestEvent) error {
			form.SetApp(e.App)
			form.SetRecord(e.Record)

			err := form.Submit()
			if err != nil {
				return firstApiError(err, e.BadRequestError("Failed to update record.", err))
			}

			err = EnrichRecord(e.RequestEvent, e.Record)
			if err != nil {
				return firstApiError(err, e.InternalServerError("Failed to enrich record", err))
			}

			err = execAfterSuccessTx(responseWriteAfterTx, e.App, func() error {
				return e.JSON(http.StatusOK, e.Record)
			})
			if err != nil {
				return err
			}

			if optFinalizer != nil {
				isOptFinalizerCalled = true
				err = optFinalizer(e.Record)
				if err != nil {
					return firstApiError(err, e.InternalServerError("", fmt.Errorf("update optFinalizer error: %w", err)))
				}
			}

			return nil
		})
		if hookErr != nil {
			return hookErr
		}

		// e.g. in case the regular hook chain was stopped and the finalizer cannot be executed as part of the last e.Next() task
		if !isOptFinalizerCalled && optFinalizer != nil {
			if err := optFinalizer(event.Record); err != nil {
				return firstApiError(err, e.InternalServerError("", fmt.Errorf("update optFinalizer error: %w", err)))
			}
		}

		return nil
	}
}

func recordDelete(responseWriteAfterTx bool, optFinalizer func(data any) error) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
		if err != nil || collection == nil {
			return e.NotFoundError("Missing collection context.", err)
		}

		if collection.IsView() {
			return e.BadRequestError("Unsupported collection type.", nil)
		}

		err = checkCollectionRateLimit(e, collection, "delete")
		if err != nil {
			return err
		}

		recordId := e.Request.PathValue("id")
		if recordId == "" {
			return e.NotFoundError("", nil)
		}

		requestInfo, err := e.RequestInfo()
		if err != nil {
			return firstApiError(err, e.BadRequestError("", err))
		}

		if !requestInfo.HasSuperuserAuth() && collection.DeleteRule == nil {
			return e.ForbiddenError("Only superusers can perform this action.", nil)
		}

		ruleFunc := func(q *dbx.SelectQuery) error {
			if !requestInfo.HasSuperuserAuth() && collection.DeleteRule != nil && *collection.DeleteRule != "" {
				resolver := core.NewRecordFieldResolver(e.App, collection, requestInfo, true)

				expr, err := search.FilterData(*collection.DeleteRule).BuildExpr(resolver)
				if err != nil {
					return err
				}

				q.AndWhere(expr)

				err = resolver.UpdateQuery(q)
				if err != nil {
					return err
				}
			}
			return nil
		}

		record, err := e.App.FindRecordById(collection, recordId, ruleFunc)
		if err != nil || record == nil {
			return e.NotFoundError("", err)
		}

		var isOptFinalizerCalled bool

		event := new(core.RecordRequestEvent)
		event.RequestEvent = e
		event.Collection = collection
		event.Record = record

		hookErr := e.App.OnRecordDeleteRequest().Trigger(event, func(e *core.RecordRequestEvent) error {
			if err := e.App.Delete(e.Record); err != nil {
				return firstApiError(err, e.BadRequestError("Failed to delete record. Make sure that the record is not part of a required relation reference.", err))
			}

			err = execAfterSuccessTx(responseWriteAfterTx, e.App, func() error {
				return e.NoContent(http.StatusNoContent)
			})
			if err != nil {
				return err
			}

			if optFinalizer != nil {
				isOptFinalizerCalled = true
				err = optFinalizer(e.Record)
				if err != nil {
					return firstApiError(err, e.InternalServerError("", fmt.Errorf("delete optFinalizer error: %w", err)))
				}
			}

			return nil
		})
		if hookErr != nil {
			return hookErr
		}

		// e.g. in case the regular hook chain was stopped and the finalizer cannot be executed as part of the last e.Next() task
		if !isOptFinalizerCalled && optFinalizer != nil {
			if err := optFinalizer(event.Record); err != nil {
				return firstApiError(err, e.InternalServerError("", fmt.Errorf("delete optFinalizer error: %w", err)))
			}
		}

		return nil
	}
}

// -------------------------------------------------------------------

func recordDataFromRequest(e *core.RequestEvent, record *core.Record) (map[string]any, error) {
	info, err := e.RequestInfo()
	if err != nil {
		return nil, err
	}

	// resolve regular fields
	result := record.ReplaceModifiers(info.Body)

	// resolve uploaded files
	uploadedFiles, err := extractUploadedFiles(e, record.Collection(), "")
	if err != nil {
		return nil, err
	}
	if len(uploadedFiles) > 0 {
		for k, files := range uploadedFiles {
			uploaded := make([]any, 0, len(files))

			// if not remove/prepend/append -> merge with the submitted
			// info.Body values to prevent accidental old files deletion
			if info.Body[k] != nil &&
				!strings.HasPrefix(k, "+") &&
				!strings.HasSuffix(k, "+") &&
				!strings.HasSuffix(k, "-") {
				existing := list.ToUniqueStringSlice(info.Body[k])
				for _, name := range existing {
					uploaded = append(uploaded, name)
				}
			}

			for _, file := range files {
				uploaded = append(uploaded, file)
			}

			result[k] = uploaded
		}

		result = record.ReplaceModifiers(result)
	}

	isAuth := record.Collection().IsAuth()

	// unset hidden fields for non-superusers
	if !info.HasSuperuserAuth() {
		for _, f := range record.Collection().Fields {
			if f.GetHidden() {
				// exception for the auth collection "password" field
				if isAuth && f.GetName() == core.FieldNamePassword {
					continue
				}

				delete(result, f.GetName())
			}
		}
	}

	return result, nil
}

func extractUploadedFiles(re *core.RequestEvent, collection *core.Collection, prefix string) (map[string][]*filesystem.File, error) {
	contentType := re.Request.Header.Get("content-type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return nil, nil // not multipart/form-data request
	}

	result := map[string][]*filesystem.File{}

	for _, field := range collection.Fields {
		if field.Type() != core.FieldTypeFile {
			continue
		}

		baseKey := field.GetName()

		keys := []string{
			baseKey,
			// prepend and append modifiers
			"+" + baseKey,
			baseKey + "+",
		}

		for _, k := range keys {
			if prefix != "" {
				k = prefix + "." + k
			}
			files, err := re.FindUploadedFiles(k)
			if err != nil && !errors.Is(err, http.ErrMissingFile) {
				return nil, err
			}
			if len(files) > 0 {
				result[k] = files
			}
		}
	}

	return result, nil
}

// hasAuthManageAccess checks whether the client is allowed to have
// [forms.RecordUpsert] auth management permissions
// (e.g. allowing to change system auth fields without oldPassword).
func hasAuthManageAccess(app core.App, requestInfo *core.RequestInfo, collection *core.Collection, query *dbx.SelectQuery) bool {
	if !collection.IsAuth() {
		return false
	}

	manageRule := collection.ManageRule

	if manageRule == nil || *manageRule == "" {
		return false // only for superusers (manageRule can't be empty)
	}

	if requestInfo == nil || requestInfo.Auth == nil {
		return false // no auth record
	}

	resolver := core.NewRecordFieldResolver(app, collection, requestInfo, true)

	expr, err := search.FilterData(*manageRule).BuildExpr(resolver)
	if err != nil {
		app.Logger().Error("Manage rule build expression error", "error", err, "collectionId", collection.Id)
		return false
	}
	query.AndWhere(expr)

	err = resolver.UpdateQuery(query)
	if err != nil {
		return false
	}

	var exists int

	err = query.Limit(1).Row(&exists)

	return err == nil && exists > 0
}

// recordsExport 导出记录为CSV文件
func recordsExport(e *core.RequestEvent) error {
	collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
	if err != nil || collection == nil {
		return e.NotFoundError("Missing collection context.", err)
	}

	// 允许所有用户导出，不进行权限检查
	var requestInfo *core.RequestInfo
	requestInfo, err = e.RequestInfo()
	if err != nil {
		return firstApiError(err, e.BadRequestError("", err))
	}

	// 使用与recordsList方法相同的条件查询逻辑
	query := e.App.RecordQuery(collection)
	fieldsResolver := core.NewRecordFieldResolver(e.App, collection, requestInfo, true)

	// 应用权限规则过滤
	if !requestInfo.HasSuperuserAuth() && collection.ListRule != nil && *collection.ListRule != "" {
		expr, err := search.FilterData(*collection.ListRule).BuildExpr(fieldsResolver)
		if err != nil {
			return err
		}
		query.AndWhere(expr)
	}

	// 隐藏字段仅对超级用户可见
	fieldsResolver.SetAllowHiddenFields(requestInfo.HasSuperuserAuth())

	// 创建搜索提供者
	searchProvider := search.NewProvider(fieldsResolver).Query(query)

	// 对非视图集合使用rowid进行计数优化
	if !collection.IsView() {
		searchProvider.CountCol("_rowid_")
	}

	// 执行查询，确保导出所有记录（不限制30条）
	records := []*core.Record{}

	// 创建一个新的查询参数副本，移除分页参数
	queryParams := e.Request.URL.Query()
	queryParams.Del("page")
	queryParams.Del("perPage")

	// 先解析查询参数，但不包含分页
	if err := searchProvider.Parse(queryParams.Encode()); err != nil {
		return e.BadRequestError("过滤表达式错误", err)
	}

	// 直接设置perPage为一个足够大的值，确保导出所有记录
	// 注意：不能设置为0，因为在Exec方法中会被重置为DefaultPerPage
	searchProvider.PerPage(10000) // 设置一个足够大的值以导出所有记录

	result, err := searchProvider.Exec(&records)

	if err != nil {
		// 详细的调试日志，记录原始错误和查询参数
		e.App.Logger().Error("Records export query error", "error", err, "queryParams", queryParams.Encode())
		return e.BadRequestError("过滤表达式错误", err)
	}

	event := new(core.RecordsListRequestEvent)
	event.RequestEvent = e
	event.Collection = collection
	event.Records = records
	event.Result = result

	// 设置响应头
	e.Response.Header().Set("Content-Type", "text/csv")
	e.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-records-%s.csv", collection.Name, time.Now().Format("2006-01-02")))

	// 创建CSV写入器，直接使用ResponseWriter作为io.Writer
	writer := csv.NewWriter(e.Response)
	defer writer.Flush()

	// 获取所有字段名作为CSV头部
	var headers []string
	// 动态获取collection定义的所有字段名称
	headers = collection.Fields.FieldNames()

	// 确保系统字段存在
	systemFields := []string{"id", "created", "updated"}
	for _, sysField := range systemFields {
		found := false
		for _, header := range headers {
			if header == sysField {
				found = true
				break
			}
		}
		if !found {
			headers = append(headers, sysField)
		}
	}

	// 写入头部
	if err := writer.Write(headers); err != nil {
		return e.InternalServerError("Failed to write CSV headers", err)
	}

	// 写入数据行
	for _, record := range records {
		var row []string
		for _, fieldName := range headers {
			value := record.Get(fieldName)
			// 处理JSON字段
			if strings.Contains(fmt.Sprintf("%T", value), "JSON") {
				if jsonStr, ok := value.(string); ok {
					row = append(row, jsonStr)
				} else {
					row = append(row, fmt.Sprintf("%v", value))
				}
			} else {
				row = append(row, fmt.Sprintf("%v", value))
			}
		}
		if err := writer.Write(row); err != nil {
			return e.InternalServerError("Failed to write CSV row", err)
		}
	}

	return nil
}

// recordsImport 从CSV文件导入记录
func recordsImport(e *core.RequestEvent) error {
	collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
	if err != nil || collection == nil {
		return e.NotFoundError("Missing collection context.", err)
	}

	// 允许所有用户导入，不进行权限检查
	_, err = e.RequestInfo()
	if err != nil {
		return firstApiError(err, e.BadRequestError("", err))
	}

	// 解析表单文件
	file, _, err := e.Request.FormFile("file")
	if err != nil {
		return e.BadRequestError("Missing file upload", err)
	}
	defer file.Close()

	// 创建CSV读取器
	reader := csv.NewReader(file)

	// 读取头部
	headers, err := reader.Read()
	if err != nil {
		return e.BadRequestError("Invalid CSV format", err)
	}

	// 读取并处理每一行数据
	var importedCount int
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return e.BadRequestError("Invalid CSV format", err)
		}

		// 创建记录数据
		data := map[string]any{}
		for i, header := range headers {
			if i < len(row) && header != "id" && header != "created" && header != "updated" { // 跳过系统字段
				data[header] = row[i]
			}
		}

		// 使用正确的API创建和保存记录
		record := core.NewRecord(collection)

		// 设置记录数据
		for key, value := range data {
			record.Set(key, value)
		}

		// 保存记录到数据库
		if err := e.App.Save(record); err != nil {
			return e.BadRequestError("Failed to save record", err)
		}

		importedCount++
	}

	return e.JSON(http.StatusOK, map[string]any{
		"success":       true,
		"importedCount": importedCount,
		"message":       "CSV导入API已准备就绪，但实际记录保存功能需要根据系统API实现",
	})
}
