package forms

import (
	"encoding/json"
	"fmt"
	"log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// CollectionsImport specifies a form model to bulk import
// (create, replace and delete) collections from a user provided list.
type CollectionsImport struct {
	config CollectionsImportConfig

	Collections  []*models.Collection `form:"collections" json:"collections"`
	DeleteOthers bool                 `form:"deleteOthers" json:"deleteOthers"`
}

// CollectionsImportConfig is the [CollectionsImport] factory initializer config.
//
// NB! App is a required struct member.
type CollectionsImportConfig struct {
	App   core.App
	TxDao *daos.Dao
}

// NewCollectionsImport creates a new [CollectionsImport] form with
// initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewCollectionsImportWithConfig] with explicitly set TxDao.
func NewCollectionsImport(app core.App) *CollectionsImport {
	return NewCollectionsImportWithConfig(CollectionsImportConfig{
		App: app,
	})
}

// NewCollectionsImportWithConfig creates a new [CollectionsImport]
// form with the provided config or panics on invalid configuration.
func NewCollectionsImportWithConfig(config CollectionsImportConfig) *CollectionsImport {
	form := &CollectionsImport{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.TxDao == nil {
		form.config.TxDao = form.config.App.Dao()
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *CollectionsImport) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Collections, validation.Required),
	)
}

// Submit applies the import, aka.:
// - imports the form collections (create or replace)
// - sync the collection changes with their related records table
// - ensures the integrity of the imported structure (aka. run validations for each collection)
// - if [form.DeleteOthers] is set, deletes all local collections that are not found in the imports list
//
// All operations are wrapped in a single transaction that are
// rollbacked on the first encountered error.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *CollectionsImport) Submit(interceptors ...InterceptorFunc) error {
	if err := form.Validate(); err != nil {
		return err
	}

	return runInterceptors(form.submit, interceptors...)
}

func (form *CollectionsImport) submit() error {
	return form.config.TxDao.RunInTransaction(func(txDao *daos.Dao) error {
		oldCollections := []*models.Collection{}
		if err := txDao.CollectionQuery().All(&oldCollections); err != nil {
			return err
		}

		mappedOldCollections := make(map[string]*models.Collection, len(oldCollections))
		for _, old := range oldCollections {
			mappedOldCollections[old.GetId()] = old
		}

		mappedFormCollections := make(map[string]*models.Collection, len(form.Collections))
		for _, collection := range form.Collections {
			mappedFormCollections[collection.GetId()] = collection
		}

		// delete all other collections not sent with the import
		if form.DeleteOthers {
			for _, old := range oldCollections {
				if mappedFormCollections[old.GetId()] == nil {
					// delete the collection
					if err := txDao.DeleteCollection(old); err != nil {
						if form.config.App.IsDebug() {
							log.Println("[CollectionsImport] DeleteOthers failure", old.Name, err)
						}
						return validation.Errors{"collections": validation.NewError(
							"collections_import_collection_delete_failure",
							fmt.Sprintf("Failed to delete collection %q (%s). Make sure that the collection is not system or referenced by other collections.", old.Name, old.Id),
						)}
					}

					delete(mappedOldCollections, old.GetId())
				}
			}
		}

		// raw insert/replace (aka. without any validations)
		// (required to make sure that all linked collections exists before running the validations)
		for _, collection := range form.Collections {
			if mappedOldCollections[collection.GetId()] == nil {
				collection.MarkAsNew()
			}

			if err := txDao.Save(collection); err != nil {
				if form.config.App.IsDebug() {
					log.Println("[CollectionsImport] Save failure", collection.Name, err)
				}
				return validation.Errors{"collections": validation.NewError(
					"collections_import_save_failure",
					fmt.Sprintf("Integrity constraints failed - the collection %q (%s) cannot be imported.", collection.Name, collection.Id),
				)}
			}
		}

		// refresh the actual persisted collections list
		refreshedCollections := []*models.Collection{}
		if err := txDao.CollectionQuery().All(&refreshedCollections); err != nil {
			return err
		}

		// trigger the validator for each existing collection to
		// ensure that the app is not left in a broken state
		for _, collection := range refreshedCollections {
			upsertModel := mappedOldCollections[collection.GetId()]
			if upsertModel == nil {
				upsertModel = collection
			}
			upsertForm := NewCollectionUpsertWithConfig(CollectionUpsertConfig{
				App:   form.config.App,
				TxDao: txDao,
			}, upsertModel)
			// load form fields with the refreshed collection state
			upsertForm.Id = collection.Id
			upsertForm.Name = collection.Name
			upsertForm.System = collection.System
			upsertForm.ListRule = collection.ListRule
			upsertForm.ViewRule = collection.ViewRule
			upsertForm.CreateRule = collection.CreateRule
			upsertForm.UpdateRule = collection.UpdateRule
			upsertForm.DeleteRule = collection.DeleteRule
			upsertForm.Schema = collection.Schema
			if err := upsertForm.Validate(); err != nil {
				// serialize the validation error(s)
				serializedErr, _ := json.Marshal(err)

				return validation.Errors{"collections": validation.NewError(
					"collections_import_validate_failure",
					fmt.Sprintf("Data validations failed for collection %q (%s): %s", collection.Name, collection.Id, serializedErr),
				)}
			}
		}

		// sync the records table for each updated collection
		for _, collection := range form.Collections {
			oldCollection := mappedOldCollections[collection.GetId()]
			if err := txDao.SyncRecordTableSchema(collection, oldCollection); err != nil {
				if form.config.App.IsDebug() {
					log.Println("[CollectionsImport] Records table sync failure", collection.Name, err)
				}
				return validation.Errors{"collections": validation.NewError(
					"collections_import_records_table_sync_failure",
					fmt.Sprintf("Failed to sync the records table changes for collection %q.", collection.Name),
				)}
			}
		}

		return nil
	})
}
