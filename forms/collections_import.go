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

	Collections   []*models.Collection `form:"collections" json:"collections"`
	DeleteMissing bool                 `form:"deleteMissing" json:"deleteMissing"`
}

// CollectionsImportConfig is the [CollectionsImport] factory initializer config.
//
// NB! App is a required struct member.
type CollectionsImportConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewCollectionsImport creates a new [CollectionsImport] form with
// initializer config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewCollectionsImportWithConfig] with explicitly set Dao.
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

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
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
// - if [form.DeleteMissing] is set, deletes all local collections that are not found in the imports list
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

	return runInterceptors(func() error {
		return form.config.Dao.RunInTransaction(func(txDao *daos.Dao) error {
			importErr := txDao.ImportCollections(
				form.Collections,
				form.DeleteMissing,
				form.beforeRecordsSync,
			)
			if importErr == nil {
				return nil
			}

			// validation failure
			if err, ok := importErr.(validation.Errors); ok {
				return err
			}

			// generic/db failure
			if form.config.App.IsDebug() {
				log.Println("Internal import failure:", importErr)
			}
			return validation.Errors{"collections": validation.NewError(
				"collections_import_failure",
				"Failed to import the collections configuration.",
			)}
		})
	}, interceptors...)
}

func (form *CollectionsImport) beforeRecordsSync(txDao *daos.Dao, mappedNew, mappedOld map[string]*models.Collection) error {
	// refresh the actual persisted collections list
	refreshedCollections := []*models.Collection{}
	if err := txDao.CollectionQuery().OrderBy("created ASC").All(&refreshedCollections); err != nil {
		return err
	}

	// trigger the validator for each existing collection to
	// ensure that the app is not left in a broken state
	for _, collection := range refreshedCollections {
		upsertModel := mappedOld[collection.GetId()]
		if upsertModel == nil {
			upsertModel = collection
		}

		upsertForm := NewCollectionUpsertWithConfig(CollectionUpsertConfig{
			App: form.config.App,
			Dao: txDao,
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
			serializedErr, _ := json.MarshalIndent(err, "", "  ")

			return validation.Errors{"collections": validation.NewError(
				"collections_import_validate_failure",
				fmt.Sprintf("Data validations failed for collection %q (%s):\n%s", collection.Name, collection.Id, serializedErr),
			)}
		}
	}

	return nil
}
