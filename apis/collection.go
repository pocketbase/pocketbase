package apis

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/search"
)

// BindCollectionApi registers the collection api endpoints and the corresponding handlers.
func BindCollectionApi(app core.App, rg *echo.Group) {
	api := collectionApi{app: app}

	subGroup := rg.Group("/collections", ActivityLogger(app), RequireAdminAuth())
	subGroup.GET("", api.list)
	subGroup.POST("", api.create)
	subGroup.GET("/:collection", api.view)
	subGroup.PATCH("/:collection", api.update)
	subGroup.DELETE("/:collection", api.delete)
}

type collectionApi struct {
	app core.App
}

func (api *collectionApi) list(c echo.Context) error {
	fieldResolver := search.NewSimpleFieldResolver(
		"id", "created", "updated", "name", "system",
	)

	collections := []*models.Collection{}

	result, err := search.NewProvider(fieldResolver).
		Query(api.app.Dao().CollectionQuery()).
		ParseAndExec(c.QueryString(), &collections)

	if err != nil {
		return rest.NewBadRequestError("", err)
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
		return rest.NewNotFoundError("", err)
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

	// read
	if err := c.Bind(form); err != nil {
		return rest.NewBadRequestError("Failed to read the submitted data due to invalid formatting.", err)
	}

	event := &core.CollectionCreateEvent{
		HttpContext: c,
		Collection:  collection,
	}

	handlerErr := api.app.OnCollectionBeforeCreateRequest().Trigger(event, func(e *core.CollectionCreateEvent) error {
		// submit
		if err := form.Submit(); err != nil {
			return rest.NewBadRequestError("Failed to create the collection.", err)
		}

		return e.HttpContext.JSON(http.StatusOK, e.Collection)
	})

	if handlerErr == nil {
		api.app.OnCollectionAfterCreateRequest().Trigger(event)
	}

	return handlerErr
}

func (api *collectionApi) update(c echo.Context) error {
	collection, err := api.app.Dao().FindCollectionByNameOrId(c.PathParam("collection"))
	if err != nil || collection == nil {
		return rest.NewNotFoundError("", err)
	}

	form := forms.NewCollectionUpsert(api.app, collection)

	// read
	if err := c.Bind(form); err != nil {
		return rest.NewBadRequestError("Failed to read the submitted data due to invalid formatting.", err)
	}

	event := &core.CollectionUpdateEvent{
		HttpContext: c,
		Collection:  collection,
	}

	handlerErr := api.app.OnCollectionBeforeUpdateRequest().Trigger(event, func(e *core.CollectionUpdateEvent) error {
		// submit
		if err := form.Submit(); err != nil {
			return rest.NewBadRequestError("Failed to update the collection.", err)
		}

		return e.HttpContext.JSON(http.StatusOK, e.Collection)
	})

	if handlerErr == nil {
		api.app.OnCollectionAfterUpdateRequest().Trigger(event)
	}

	return handlerErr
}

func (api *collectionApi) delete(c echo.Context) error {
	collection, err := api.app.Dao().FindCollectionByNameOrId(c.PathParam("collection"))
	if err != nil || collection == nil {
		return rest.NewNotFoundError("", err)
	}

	event := &core.CollectionDeleteEvent{
		HttpContext: c,
		Collection:  collection,
	}

	handlerErr := api.app.OnCollectionBeforeDeleteRequest().Trigger(event, func(e *core.CollectionDeleteEvent) error {
		if err := api.app.Dao().DeleteCollection(e.Collection); err != nil {
			return rest.NewBadRequestError("Failed to delete collection. Make sure that the collection is not referenced by other collections.", err)
		}

		// try to delete the collection files
		if err := api.deleteCollectionFiles(e.Collection); err != nil && api.app.IsDebug() {
			// non critical error - only log for debug
			// (usually could happen because of S3 api limits)
			log.Println(err)
		}

		return e.HttpContext.NoContent(http.StatusNoContent)
	})

	if handlerErr == nil {
		api.app.OnCollectionAfterDeleteRequest().Trigger(event)
	}

	return handlerErr
}

func (api *collectionApi) deleteCollectionFiles(collection *models.Collection) error {
	fs, err := api.app.NewFilesystem()
	if err != nil {
		return err
	}
	defer fs.Close()

	failed := fs.DeletePrefix(collection.BaseFilesPath())
	if len(failed) > 0 {
		return errors.New("Failed to delete all record files.")
	}

	return nil
}
