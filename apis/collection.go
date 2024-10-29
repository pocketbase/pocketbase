package apis

import (
	"errors"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/search"
)

// bindCollectionApi registers the collection api endpoints and the corresponding handlers.
func bindCollectionApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/collections").Bind(RequireSuperuserAuth())
	subGroup.GET("", collectionsList)
	subGroup.POST("", collectionCreate)
	subGroup.GET("/{collection}", collectionView)
	subGroup.PATCH("/{collection}", collectionUpdate)
	subGroup.DELETE("/{collection}", collectionDelete)
	subGroup.DELETE("/{collection}/truncate", collectionTruncate)
	subGroup.PUT("/import", collectionsImport)
	subGroup.GET("/meta/scaffolds", collectionScaffolds)
}

func collectionsList(e *core.RequestEvent) error {
	fieldResolver := search.NewSimpleFieldResolver(
		"id", "created", "updated", "name", "system", "type",
	)

	collections := []*core.Collection{}

	result, err := search.NewProvider(fieldResolver).
		Query(e.App.CollectionQuery()).
		ParseAndExec(e.Request.URL.Query().Encode(), &collections)

	if err != nil {
		return e.BadRequestError("", err)
	}

	event := new(core.CollectionsListRequestEvent)
	event.RequestEvent = e
	event.Collections = collections
	event.Result = result

	return event.App.OnCollectionsListRequest().Trigger(event, func(e *core.CollectionsListRequestEvent) error {
		return e.JSON(http.StatusOK, e.Result)
	})
}

func collectionView(e *core.RequestEvent) error {
	collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
	if err != nil || collection == nil {
		return e.NotFoundError("", err)
	}

	event := new(core.CollectionRequestEvent)
	event.RequestEvent = e
	event.Collection = collection

	return e.App.OnCollectionViewRequest().Trigger(event, func(e *core.CollectionRequestEvent) error {
		return e.JSON(http.StatusOK, e.Collection)
	})
}

func collectionCreate(e *core.RequestEvent) error {
	// populate the minimal required factory collection data (if any)
	factoryExtract := struct {
		Type string `form:"type" json:"type"`
		Name string `form:"name" json:"name"`
	}{}
	if err := e.BindBody(&factoryExtract); err != nil {
		return e.BadRequestError("Failed to load the collection type data due to invalid formatting.", err)
	}

	// create scaffold
	collection := core.NewCollection(factoryExtract.Type, factoryExtract.Name)

	// merge the scaffold with the submitted request data
	if err := e.BindBody(collection); err != nil {
		return e.BadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := new(core.CollectionRequestEvent)
	event.RequestEvent = e
	event.Collection = collection

	return e.App.OnCollectionCreateRequest().Trigger(event, func(e *core.CollectionRequestEvent) error {
		if err := e.App.Save(e.Collection); err != nil {
			// validation failure
			var validationErrors validation.Errors
			if errors.As(err, &validationErrors) {
				return e.BadRequestError("Failed to create collection.", validationErrors)
			}

			// other generic db error
			return e.BadRequestError("Failed to create collection. Raw error: \n"+err.Error(), nil)
		}

		return e.JSON(http.StatusOK, e.Collection)
	})
}

func collectionUpdate(e *core.RequestEvent) error {
	collection, err := e.App.FindCollectionByNameOrId(e.Request.PathValue("collection"))
	if err != nil || collection == nil {
		return e.NotFoundError("", err)
	}

	if err := e.BindBody(collection); err != nil {
		return e.BadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	event := new(core.CollectionRequestEvent)
	event.RequestEvent = e
	event.Collection = collection

	return event.App.OnCollectionUpdateRequest().Trigger(event, func(e *core.CollectionRequestEvent) error {
		if err := e.App.Save(e.Collection); err != nil {
			// validation failure
			var validationErrors validation.Errors
			if errors.As(err, &validationErrors) {
				return e.BadRequestError("Failed to update collection.", validationErrors)
			}

			// other generic db error
			return e.BadRequestError("Failed to update collection. Raw error: \n"+err.Error(), nil)
		}

		return e.JSON(http.StatusOK, e.Collection)
	})
}

func collectionDelete(e *core.RequestEvent) error {
	collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
	if err != nil || collection == nil {
		return e.NotFoundError("", err)
	}

	event := new(core.CollectionRequestEvent)
	event.RequestEvent = e
	event.Collection = collection

	return e.App.OnCollectionDeleteRequest().Trigger(event, func(e *core.CollectionRequestEvent) error {
		if err := e.App.Delete(e.Collection); err != nil {
			msg := "Failed to delete collection"

			// check fo references
			refs, _ := e.App.FindCollectionReferences(e.Collection, e.Collection.Id)
			if len(refs) > 0 {
				names := make([]string, 0, len(refs))
				for ref := range refs {
					names = append(names, ref.Name)
				}
				msg += " probably due to existing reference in " + strings.Join(names, ", ")
			}

			return e.BadRequestError(msg, err)
		}

		return e.NoContent(http.StatusNoContent)
	})
}

func collectionTruncate(e *core.RequestEvent) error {
	collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("collection"))
	if err != nil || collection == nil {
		return e.NotFoundError("", err)
	}

	if collection.IsView() {
		return e.BadRequestError("View collections cannot be truncated since they don't store their own records.", nil)
	}

	err = e.App.TruncateCollection(collection)
	if err != nil {
		return e.BadRequestError("Failed to truncate collection (most likely due to required cascade delete record references).", err)
	}

	return e.NoContent(http.StatusNoContent)
}

func collectionScaffolds(e *core.RequestEvent) error {
	collections := map[string]*core.Collection{
		core.CollectionTypeBase: core.NewBaseCollection(""),
		core.CollectionTypeAuth: core.NewAuthCollection(""),
		core.CollectionTypeView: core.NewViewCollection(""),
	}

	for _, c := range collections {
		c.Id = "" // clear autogenerated id
	}

	return e.JSON(http.StatusOK, collections)
}
