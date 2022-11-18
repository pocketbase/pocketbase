package apis

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/search"
)

// BindViewApi registers the view api endpoints and the corresponding handlers.
func bindViewApi(app core.App, rg *echo.Group) {
	api := ViewApi{app: app}

	// client path
	rg.GET("/views/:view/records", api.list, ActivityLogger(app))

	// admin subgroup
	adminSG := rg.Group("/views", ActivityLogger(app), RequireAdminAuth())
	adminSG.GET("", api.adminList)
	adminSG.POST("", api.create)
	adminSG.GET("/:view", api.view)
	adminSG.PATCH("/:view", api.update)
	adminSG.DELETE("/:view", api.delete)
}

type ViewApi struct {
	app core.App
}

func (api *ViewApi) list(c echo.Context) error {
	v, err := api.app.Dao().FindViewByIdOrName(c.PathParam("view"))
	if err != nil {
		return NewNotFoundError("View was not found", err)
	}

	admin, _ := c.Get(ContextAdminKey).(*models.Admin)
	if admin == nil && v.ListRule == nil {
		// only admins can access if the rule is nil
		return NewForbiddenError("Only admins can perform this action.", nil)
	}

	// forbid users and guests to query special filter/sort fields
	RA := recordApi{}
	if err := RA.checkForForbiddenQueryFields(c); err != nil {
		return err
	}

	requestData := RequestData(c)


	fieldsResolver := resolvers.NewRecordFieldResolver(api.app.Dao(), &models.Collection{Schema: v.Schema, Name: v.Name},
		 requestData,false)

	searchProvider := search.NewProvider(fieldsResolver).
		Query(api.app.Dao().RecordFromViewQuery(v))

	if admin == nil && v.ListRule != nil {
		searchProvider.AddFilter(search.FilterData(*v.ListRule))
	}

	rawRecords := []dbx.NullStringMap{}
	result, err := searchProvider.ParseAndExec(c.QueryString(), &rawRecords)
	if err != nil {
		return NewBadRequestError("Invalid filter parameters.", err)
	}
	result.Items = models.NewRecordsFromViewSchema(v, rawRecords)

	event := &core.RecordsFromViewListEvent{
		HttpContext: c,
		Result:      result,
	}

	return api.app.OnRecordsFromViewListRequest().Trigger(event, func(e *core.RecordsFromViewListEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.Result)
	})
}

func (api *ViewApi) adminList(c echo.Context) error {
	fieldResolver := search.NewSimpleFieldResolver()

	View := []*models.View{}

	result, err := search.NewProvider(fieldResolver).Query(api.app.Dao().ViewQuery()).ParseAndExec(c.QueryString(), &View)
	if err != nil {
		return NewBadRequestError("", err)
	}

	event := &core.ViewListEvent{
		HttpContext: c,
		Result:      result,
	}

	return api.app.OnViewListRequest().Trigger(event, func(e *core.ViewListEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.Result)
	})
}

func (api *ViewApi) view(c echo.Context) error {
	v, err := api.app.Dao().FindViewByIdOrName(c.PathParam("view"))
	if err != nil || v == nil {
		return NewNotFoundError("", err)
	}

	event := &core.ViewViewEvent{
		HttpContext: c,
		View:        v,
	}

	return api.app.OnViewViewRequest().Trigger(event, func(e *core.ViewViewEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.View)
	})
}

func (api *ViewApi) create(c echo.Context) error {
	v := &models.View{}

	form := forms.NewViewUpsert(api.app, v)

	// load request
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.ViewCreateEvent{
		HttpContext: c,
		View:        v,
	}

	// create the view
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnViewBeforeCreateRequest().Trigger(event, func(e *core.ViewCreateEvent) error {
				if err := next(); err != nil {
					return NewBadRequestError("Failed to create the view.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.View)
			})
		}
	})

	if submitErr == nil {
		api.app.OnViewAfterCreateRequest().Trigger(event)
	}

	return submitErr
}

func (api *ViewApi) update(c echo.Context) error {
	v, err := api.app.Dao().FindViewByIdOrName(c.PathParam("view"))
	if err != nil || v == nil {
		return NewNotFoundError("", err)
	}

	form := forms.NewViewUpsert(api.app, v)

	// load request
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.ViewUpdateEvent{
		HttpContext: c,
		View:        v,
	}

	// update the view
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnViewBeforeUpdateRequest().Trigger(event, func(e *core.ViewUpdateEvent) error {
				if err := next(); err != nil {
					return NewBadRequestError("Failed to update the view.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.View)
			})
		}
	})

	if submitErr == nil {
		api.app.OnViewAfterUpdateRequest().Trigger(event)
	}

	return submitErr
}

func (api *ViewApi) delete(c echo.Context) error {
	v, err := api.app.Dao().FindViewByIdOrName(c.PathParam("view"))
	if err != nil || v == nil {
		api.app.Dao().DeleteView(c.PathParam("view"))
		return NewNotFoundError("", err)
	}

	event := &core.ViewDeleteEvent{
		HttpContext: c,
		View:        v,
	}

	handlerErr := api.app.OnViewBeforeDeleteRequest().Trigger(event, func(e *core.ViewDeleteEvent) error {
		if err := api.app.Dao().DeleteViewModel(e.View); err != nil {
			return NewBadRequestError("Failed to delete View. ", err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnViewAfterDeleteRequest().Trigger(event)
	}

	return handlerErr
}
