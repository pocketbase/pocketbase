package daos

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

// CollectionQuery returns a new Collection select query.
func (dao *Dao) CollectionQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&models.Collection{})
}

// FindCollectionByNameOrId finds the first collection by its name or id.
func (dao *Dao) FindCollectionByNameOrId(nameOrId string) (*models.Collection, error) {
	model := &models.Collection{}

	err := dao.CollectionQuery().
		AndWhere(dbx.Or(
			dbx.HashExp{"id": nameOrId},
			dbx.HashExp{"name": nameOrId},
		)).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

// IsCollectionNameUnique checks that there is no existing collection
// with the provided name (case insensitive!).
//
// Note: case sensitive check because the name is used also as a table name for the records.
func (dao *Dao) IsCollectionNameUnique(name string, excludeId string) bool {
	if name == "" {
		return false
	}

	var exists bool
	err := dao.CollectionQuery().
		Select("count(*)").
		AndWhere(dbx.Not(dbx.HashExp{"id": excludeId})).
		AndWhere(dbx.NewExp("LOWER([[name]])={:name}", dbx.Params{"name": strings.ToLower(name)})).
		Limit(1).
		Row(&exists)

	return err == nil && !exists
}

// FindCollectionsWithUserFields finds all collections that has
// at least one user schema field.
func (dao *Dao) FindCollectionsWithUserFields() ([]*models.Collection, error) {
	result := []*models.Collection{}

	err := dao.CollectionQuery().
		InnerJoin(
			"json_each(schema) as jsonField",
			dbx.NewExp(
				"json_extract(jsonField.value, '$.type') = {:type}",
				dbx.Params{"type": schema.FieldTypeUser},
			),
		).
		All(&result)

	return result, err
}

// FindCollectionReferences returns information for all
// relation schema fields referencing the provided collection.
//
// If the provided collection has reference to itself then it will be
// also included in the result. To exlude it, pass the collection id
// as the excludeId argument.
func (dao *Dao) FindCollectionReferences(collection *models.Collection, excludeId string) (map[*models.Collection][]*schema.SchemaField, error) {
	collections := []*models.Collection{}

	err := dao.CollectionQuery().
		AndWhere(dbx.Not(dbx.HashExp{"id": excludeId})).
		All(&collections)
	if err != nil {
		return nil, err
	}

	result := map[*models.Collection][]*schema.SchemaField{}
	for _, c := range collections {
		for _, f := range c.Schema.Fields() {
			if f.Type != schema.FieldTypeRelation {
				continue
			}
			f.InitOptions()
			options, _ := f.Options.(*schema.RelationOptions)
			if options != nil && options.CollectionId == collection.Id {
				result[c] = append(result[c], f)
			}
		}
	}

	return result, nil
}

// DeleteCollection deletes the provided Collection model.
// This method automatically deletes the related collection records table.
//
// NB! The collection cannot be deleted, if:
// - is system collection (aka. collection.System is true)
// - is referenced as part of a relation field in another collection
func (dao *Dao) DeleteCollection(collection *models.Collection) error {
	if collection.System {
		return errors.New("System collections cannot be deleted.")
	}

	// ensure that there aren't any existing references.
	// note: the select is outside of the transaction to prevent SQLITE_LOCKED error when mixing read&write in a single transaction
	result, err := dao.FindCollectionReferences(collection, collection.Id)
	if err != nil {
		return err
	}
	if total := len(result); total > 0 {
		return fmt.Errorf("The collection has external relation field references (%d).", total)
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		// delete the related records table
		if err := txDao.DeleteTable(collection.Name); err != nil {
			return err
		}

		return txDao.Delete(collection)
	})
}

// SaveCollection upserts the provided Collection model and updates
// its related records table schema.
func (dao *Dao) SaveCollection(collection *models.Collection) error {
	var oldCollection *models.Collection

	if !collection.IsNew() {
		// get the existing collection state to compare with the new one
		// note: the select is outside of the transaction to prevent SQLITE_LOCKED error when mixing read&write in a single transaction
		var findErr error
		oldCollection, findErr = dao.FindCollectionByNameOrId(collection.Id)
		if findErr != nil {
			return findErr
		}
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		// persist the collection model
		if err := txDao.Save(collection); err != nil {
			return err
		}

		// sync the changes with the related records table
		return txDao.SyncRecordTableSchema(collection, oldCollection)
	})
}
