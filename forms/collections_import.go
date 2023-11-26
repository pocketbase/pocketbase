package forms

import (
	"encoding/json"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// CollectionsImport is a form model to bulk import
// (create, replace and delete) collections from a user provided list.
type CollectionsImport struct {
	app core.App
	dao *daos.Dao

	Collections   []*models.Collection `form:"collections" json:"collections"`
	DeleteMissing bool                 `form:"deleteMissing" json:"deleteMissing"`
}

// NewCollectionsImport creates a new [CollectionsImport] form with
// initialized with from the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewCollectionsImport(app core.App) *CollectionsImport {
	return &CollectionsImport{
		app: app,
		dao: app.Dao(),
	}
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *CollectionsImport) SetDao(dao *daos.Dao) {
	form.dao = dao
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
func (form *CollectionsImport) Submit(interceptors ...InterceptorFunc[[]*models.Collection]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	return runInterceptors(form.Collections, func(collections []*models.Collection) error {
		return form.dao.RunInTransaction(func(txDao *daos.Dao) error {
			importErr := txDao.ImportCollections(
				collections,
				form.DeleteMissing,
				form.afterSync,
			)
			if importErr == nil {
				return nil
			}

			// validation failure
			if err, ok := importErr.(validation.Errors); ok {
				return err
			}

			// generic/db failure
			return validation.Errors{"collections": validation.NewError(
				"collections_import_failure",
				"Failed to import the collections configuration. Raw error:\n"+importErr.Error(),
			)}
		})
	}, interceptors...)
}

func (form *CollectionsImport) afterSync(txDao *daos.Dao, mappedNew, mappedOld map[string]*models.Collection) error {
	// refresh the actual persisted collections list
	refreshedCollections := []*models.Collection{}
	if err := txDao.CollectionQuery().OrderBy("updated ASC").All(&refreshedCollections); err != nil {
		return err
	}

	// trigger the validator for each existing collection to
	// ensure that the app is not left in a broken state
	for _, collection := range refreshedCollections {
		upsertModel := mappedOld[collection.GetId()]
		if upsertModel == nil {
			upsertModel = collection
		}
		upsertModel.MarkAsNotNew()

		upsertForm := NewCollectionUpsert(form.app, upsertModel)
		upsertForm.SetDao(txDao)

		// load form fields with the refreshed collection state
		upsertForm.Id = collection.Id
		upsertForm.Type = collection.Type
		upsertForm.Name = collection.Name
		upsertForm.System = collection.System
		upsertForm.ListRule = collection.ListRule
		upsertForm.ViewRule = collection.ViewRule
		upsertForm.CreateRule = collection.CreateRule
		upsertForm.UpdateRule = collection.UpdateRule
		upsertForm.DeleteRule = collection.DeleteRule
		upsertForm.Schema = collection.Schema
		upsertForm.Options = collection.Options

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
