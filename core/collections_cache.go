package core

// -------------------------------------------------------------------
// This is a small optimization ported from the [ongoing refactoring branch](https://github.com/pocketbase/pocketbase/discussions/4355).
//
// @todo remove after the refactoring is finalized.
// -------------------------------------------------------------------

import (
	"strings"

	"github.com/pocketbase/pocketbase/models"
)

const storeCachedCollectionsKey = "@cachedCollectionsContext"

func registerCachedCollectionsAppHooks(app App) {
	collectionsChangeFunc := func(e *ModelEvent) error {
		if _, ok := e.Model.(*models.Collection); !ok {
			return nil
		}

		_ = ReloadCachedCollections(app)

		return nil
	}
	app.OnModelAfterCreate().Add(collectionsChangeFunc)
	app.OnModelAfterUpdate().Add(collectionsChangeFunc)
	app.OnModelAfterDelete().Add(collectionsChangeFunc)
	app.OnBeforeServe().Add(func(e *ServeEvent) error {
		_ = ReloadCachedCollections(e.App)
		return nil
	})
}

func ReloadCachedCollections(app App) error {
	collections := []*models.Collection{}

	err := app.Dao().CollectionQuery().All(&collections)
	if err != nil {
		return err
	}

	app.Store().Set(storeCachedCollectionsKey, collections)

	return nil
}

func FindCachedCollectionByNameOrId(app App, nameOrId string) (*models.Collection, error) {
	// retrieve from the app cache
	// ---
	collections, _ := app.Store().Get(storeCachedCollectionsKey).([]*models.Collection)
	for _, c := range collections {
		if strings.EqualFold(c.Name, nameOrId) || c.Id == nameOrId {
			return c, nil
		}
	}

	// retrieve from the database
	// ---
	found, err := app.Dao().FindCollectionByNameOrId(nameOrId)
	if err != nil {
		return nil, err
	}

	err = ReloadCachedCollections(app)
	if err != nil {
		app.Logger().Warn("Failed to reload collections cache", "error", err)
	}

	return found, nil
}
