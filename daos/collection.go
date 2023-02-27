package daos

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
)

// CollectionQuery returns a new Collection select query.
func (dao *Dao) CollectionQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&models.Collection{})
}

// FindCollectionsByType finds all collections by the given type.
func (dao *Dao) FindCollectionsByType(collectionType string) ([]*models.Collection, error) {
	collections := []*models.Collection{}

	err := dao.CollectionQuery().
		AndWhere(dbx.HashExp{"type": collectionType}).
		OrderBy("created ASC").
		All(&collections)

	if err != nil {
		return nil, err
	}

	return collections, nil
}

// FindCollectionByNameOrId finds a single collection by its name (case insensitive) or id.
func (dao *Dao) FindCollectionByNameOrId(nameOrId string) (*models.Collection, error) {
	model := &models.Collection{}

	err := dao.CollectionQuery().
		AndWhere(dbx.NewExp("[[id]] = {:id} OR LOWER([[name]])={:name}", dbx.Params{
			"id":   nameOrId,
			"name": strings.ToLower(nameOrId),
		})).
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
// Note: case insensitive check because the name is used also as a table name for the records.
func (dao *Dao) IsCollectionNameUnique(name string, excludeIds ...string) bool {
	if name == "" {
		return false
	}

	query := dao.CollectionQuery().
		Select("count(*)").
		AndWhere(dbx.NewExp("LOWER([[name]])={:name}", dbx.Params{"name": strings.ToLower(name)})).
		Limit(1)

	if uniqueExcludeIds := list.NonzeroUniques(excludeIds); len(uniqueExcludeIds) > 0 {
		query.AndWhere(dbx.NotIn("id", list.ToInterfaceSlice(uniqueExcludeIds)...))
	}

	var exists bool

	return query.Row(&exists) == nil && !exists
}

// FindCollectionReferences returns information for all
// relation schema fields referencing the provided collection.
//
// If the provided collection has reference to itself then it will be
// also included in the result. To exclude it, pass the collection id
// as the excludeId argument.
func (dao *Dao) FindCollectionReferences(collection *models.Collection, excludeIds ...string) (map[*models.Collection][]*schema.SchemaField, error) {
	collections := []*models.Collection{}

	query := dao.CollectionQuery()

	if uniqueExcludeIds := list.NonzeroUniques(excludeIds); len(uniqueExcludeIds) > 0 {
		query.AndWhere(dbx.NotIn("id", list.ToInterfaceSlice(uniqueExcludeIds)...))
	}

	if err := query.All(&collections); err != nil {
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
		return fmt.Errorf("System collection %q cannot be deleted.", collection.Name)
	}

	// ensure that there aren't any existing references.
	// note: the select is outside of the transaction to prevent SQLITE_LOCKED error when mixing read&write in a single transaction
	result, err := dao.FindCollectionReferences(collection, collection.Id)
	if err != nil {
		return err
	}
	if total := len(result); total > 0 {
		names := make([]string, 0, len(result))
		for ref := range result {
			names = append(names, ref.Name)
		}
		return fmt.Errorf("The collection %q has external relation field references (%s).", collection.Name, strings.Join(names, ", "))
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		// delete the related view or records table
		if collection.IsView() {
			if err := txDao.DeleteView(collection.Name); err != nil {
				return err
			}
		} else {
			if err := txDao.DeleteTable(collection.Name); err != nil {
				return err
			}
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
		// set default collection type
		if collection.Type == "" {
			collection.Type = models.CollectionTypeBase
		}

		switch collection.Type {
		case models.CollectionTypeView:
			return txDao.saveViewCollection(collection, oldCollection)
		default:
			// persist the collection model
			if err := txDao.Save(collection); err != nil {
				return err
			}

			// sync the changes with the related records table
			return txDao.SyncRecordTableSchema(collection, oldCollection)
		}
	})
}

// ImportCollections imports the provided collections list within a single transaction.
//
// NB1! If deleteMissing is set, all local collections and schema fields, that are not present
// in the imported configuration, WILL BE DELETED (including their related records data).
//
// NB2! This method doesn't perform validations on the imported collections data!
// If you need validations, use [forms.CollectionsImport].
func (dao *Dao) ImportCollections(
	importedCollections []*models.Collection,
	deleteMissing bool,
	beforeRecordsSync func(txDao *Dao, mappedImported, mappedExisting map[string]*models.Collection) error,
) error {
	if len(importedCollections) == 0 {
		return errors.New("No collections to import")
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		existingCollections := []*models.Collection{}
		if err := txDao.CollectionQuery().OrderBy("created ASC").All(&existingCollections); err != nil {
			return err
		}
		mappedExisting := make(map[string]*models.Collection, len(existingCollections))
		for _, existing := range existingCollections {
			mappedExisting[existing.GetId()] = existing
		}

		mappedImported := make(map[string]*models.Collection, len(importedCollections))
		for _, imported := range importedCollections {
			// generate id if not set
			if !imported.HasId() {
				imported.MarkAsNew()
				imported.RefreshId()
			}

			// set default type if missing
			if imported.Type == "" {
				imported.Type = models.CollectionTypeBase
			}

			if existing, ok := mappedExisting[imported.GetId()]; ok {
				imported.MarkAsNotNew()

				// preserve original created date
				if !existing.Created.IsZero() {
					imported.Created = existing.Created
				}

				// extend existing schema
				if !deleteMissing {
					schema, _ := existing.Schema.Clone()
					for _, f := range imported.Schema.Fields() {
						schema.AddField(f) // add or replace
					}
					imported.Schema = *schema
				}
			} else {
				imported.MarkAsNew()
			}

			mappedImported[imported.GetId()] = imported
		}

		// delete old collections not available in the new configuration
		// (before saving the imports in case a deleted collection name is being reused)
		if deleteMissing {
			for _, existing := range existingCollections {
				if mappedImported[existing.GetId()] != nil {
					continue // exist
				}

				if existing.System {
					return fmt.Errorf("System collection %q cannot be deleted.", existing.Name)
				}

				// delete the collection
				if err := txDao.Delete(existing); err != nil {
					return err
				}
			}
		}

		// upsert imported collections
		for _, imported := range importedCollections {
			if err := txDao.Save(imported); err != nil {
				return err
			}
		}

		if beforeRecordsSync != nil {
			if err := beforeRecordsSync(txDao, mappedImported, mappedExisting); err != nil {
				return err
			}
		}

		// delete the record tables of the deleted collections
		if deleteMissing {
			for _, existing := range existingCollections {
				if mappedImported[existing.GetId()] != nil {
					continue // exist
				}

				if existing.IsView() {
					if err := txDao.DeleteView(existing.Name); err != nil {
						return err
					}
				} else {
					if err := txDao.DeleteTable(existing.Name); err != nil {
						return err
					}
				}
			}
		}

		// sync the upserted collections with the related records table
		for _, imported := range importedCollections {
			existing := mappedExisting[imported.GetId()]

			if imported.IsView() {
				if err := txDao.saveViewCollection(imported, existing); err != nil {
					return err
				}
			} else {
				if err := txDao.SyncRecordTableSchema(imported, existing); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// saveViewCollection persists the provided View collection changes:
//   - deletes the old related SQL view (if any)
//   - creates a new SQL view with the latest newCollection.Options.Query
//   - generates a new schema based on newCollection.Options.Query
//   - updates newCollection.Schema based on the generated view table info and query
//   - saves the newCollection
//
// This method returns an error if newCollection is not a "view".
func (dao *Dao) saveViewCollection(newCollection *models.Collection, oldCollection *models.Collection) error {
	if newCollection.IsAuth() {
		return errors.New("not a view collection")
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		query := newCollection.ViewOptions().Query

		// generate collection schema from the query
		schema, err := txDao.CreateViewSchema(query)
		if err != nil {
			return err
		}

		// delete old renamed view
		if oldCollection != nil && newCollection.Name != oldCollection.Name {
			if err := txDao.DeleteView(oldCollection.Name); err != nil {
				return err
			}
		}

		// (re)create the view
		if err := txDao.SaveView(newCollection.Name, query); err != nil {
			return err
		}

		newCollection.Schema = schema

		return txDao.Save(newCollection)
	})
}
