package apis

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/search"
)

// bindCollectionApi registers the collection api endpoints and the corresponding handlers.
func bindCollectionApi(app core.App, rg *echo.Group) {
	api := collectionApi{app: app}

	subGroup := rg.Group("/collections", ActivityLogger(app), RequireAdminAuth())
	subGroup.GET("", api.list)
	subGroup.POST("", api.create)
	subGroup.GET("/:collection", api.view)
	subGroup.PATCH("/:collection", api.update)
	subGroup.DELETE("/:collection", api.delete)
	subGroup.PUT("/import", api.bulkImport)
}

type collectionApi struct {
	app core.App
}

func (api *collectionApi) list(c echo.Context) error {
	fieldResolver := search.NewSimpleFieldResolver(
		"id", "created", "updated", "name", "system", "type",
	)

	collections := []*models.Collection{}

	result, err := search.NewProvider(fieldResolver).
		Query(api.app.Dao().CollectionQuery()).
		ParseAndExec(c.QueryParams().Encode(), &collections)

	if err != nil {
		return NewBadRequestError("", err)
	}

	event := &core.CollectionsListEvent{
		HttpContext: c,
		Collections: collections,
		Result:      result,
	}

	return api.app.OnCollectionsListRequest().Trigger(event, func(e *core.CollectionsListEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.Result)
	})
}

func (api *collectionApi) view(c echo.Context) error {
	collection, err := api.app.Dao().FindCollectionByNameOrId(c.PathParam("collection"))
	if err != nil || collection == nil {
		return NewNotFoundError("", err)
	}

	event := &core.CollectionViewEvent{
		HttpContext: c,
		Collection:  collection,
	}

	return api.app.OnCollectionViewRequest().Trigger(event, func(e *core.CollectionViewEvent) error {
		return e.HttpContext.JSON(http.StatusOK, e.Collection)
	})
}

func (api *collectionApi) create(c echo.Context) error {
	collection := &models.Collection{}

	form := forms.NewCollectionUpsert(api.app, collection)

	// load request
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.CollectionCreateEvent{
		HttpContext: c,
		Collection:  collection,
	}

	// create the collection
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnCollectionBeforeCreateRequest().Trigger(event, func(e *core.CollectionCreateEvent) error {
				if err := next(); err != nil {
					return NewBadRequestError("Failed to create the collection.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.Collection)
			})
		}
	})

	if submitErr == nil {
		api.app.OnCollectionAfterCreateRequest().Trigger(event)
	}

	return submitErr
}

func (api *collectionApi) update(c echo.Context) error {
	collection, err := api.app.Dao().FindCollectionByNameOrId(c.PathParam("collection"))
	if err != nil || collection == nil {
		return NewNotFoundError("", err)
	}

	form := forms.NewCollectionUpsert(api.app, collection)

	// load request
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.CollectionUpdateEvent{
		HttpContext: c,
		Collection:  collection,
	}

	// update the collection
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnCollectionBeforeUpdateRequest().Trigger(event, func(e *core.CollectionUpdateEvent) error {
				if err := next(); err != nil {
					return NewBadRequestError("Failed to update the collection.", err)
				}

				return e.HttpContext.JSON(http.StatusOK, e.Collection)
			})
		}
	})

	if submitErr == nil {
		api.app.OnCollectionAfterUpdateRequest().Trigger(event)
	}

	return submitErr
}

func (api *collectionApi) delete(c echo.Context) error {
	collection, err := api.app.Dao().FindCollectionByNameOrId(c.PathParam("collection"))
	if err != nil || collection == nil {
		return NewNotFoundError("", err)
	}

	event := &core.CollectionDeleteEvent{
		HttpContext: c,
		Collection:  collection,
	}

	handlerErr := api.app.OnCollectionBeforeDeleteRequest().Trigger(event, func(e *core.CollectionDeleteEvent) error {
		if err := api.app.Dao().DeleteCollection(e.Collection); err != nil {
			return NewBadRequestError("Failed to delete collection. Make sure that the collection is not referenced by other collections.", err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnCollectionAfterDeleteRequest().Trigger(event)
	}

	return handlerErr
}

func (api *collectionApi) bulkImport(c echo.Context) error {
	form := forms.NewCollectionsImport(api.app)

	// load request data
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := &core.CollectionsImportEvent{
		HttpContext: c,
		Collections: form.Collections,
	}

	// import collections
	submitErr := form.Submit(func(next forms.InterceptorNextFunc) forms.InterceptorNextFunc {
		return func() error {
			return api.app.OnCollectionsBeforeImportRequest().Trigger(event, func(e *core.CollectionsImportEvent) error {
				form.Collections = e.Collections // ensures that the form always has the latest changes

				if err := next(); err != nil {
					return NewBadRequestError("Failed to import the submitted collections.", err)
				}

				return e.HttpContext.NoContent(http.StatusNoContent)
			})
		}
	})

	if submitErr == nil {
		api.app.OnCollectionsAfterImportRequest().Trigger(event)
	}

	return submitErr
}
