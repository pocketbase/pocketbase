package daos

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
)

// RecordQuery returns a new Record select query.
func (dao *Dao) RecordQuery(collection *models.Collection) *dbx.SelectQuery {
	tableName := collection.Name
	selectCols := fmt.Sprintf("%s.*", dao.DB().QuoteSimpleColumnName(tableName))

	return dao.DB().Select(selectCols).From(tableName)
}

// FindRecordById finds the Record model by its id.
func (dao *Dao) FindRecordById(
	collection *models.Collection,
	recordId string,
	filter func(q *dbx.SelectQuery) error,
) (*models.Record, error) {
	tableName := collection.Name

	query := dao.RecordQuery(collection).
		AndWhere(dbx.HashExp{tableName + ".id": recordId})

	if filter != nil {
		if err := filter(query); err != nil {
			return nil, err
		}
	}

	row := dbx.NullStringMap{}
	if err := query.Limit(1).One(row); err != nil {
		return nil, err
	}

	return models.NewRecordFromNullStringMap(collection, row), nil
}

// FindRecordsByIds finds all Record models by the provided ids.
// If no records are found, returns an empty slice.
func (dao *Dao) FindRecordsByIds(
	collection *models.Collection,
	recordIds []string,
	filter func(q *dbx.SelectQuery) error,
) ([]*models.Record, error) {
	tableName := collection.Name

	query := dao.RecordQuery(collection).
		AndWhere(dbx.In(tableName+".id", list.ToInterfaceSlice(recordIds)...))

	if filter != nil {
		if err := filter(query); err != nil {
			return nil, err
		}
	}

	rows := []dbx.NullStringMap{}
	if err := query.All(&rows); err != nil {
		return nil, err
	}

	return models.NewRecordsFromNullStringMaps(collection, rows), nil
}

// FindRecordsByExpr finds all records by the provided db expression.
// If no records are found, returns an empty slice.
//
// Example:
//	expr := dbx.HashExp{"email": "test@example.com"}
//	dao.FindRecordsByExpr(collection, expr)
func (dao *Dao) FindRecordsByExpr(collection *models.Collection, expr dbx.Expression) ([]*models.Record, error) {
	if expr == nil {
		return nil, errors.New("Missing filter expression")
	}

	rows := []dbx.NullStringMap{}

	err := dao.RecordQuery(collection).
		AndWhere(expr).
		All(&rows)

	if err != nil {
		return nil, err
	}

	return models.NewRecordsFromNullStringMaps(collection, rows), nil
}

// FindFirstRecordByData returns the first found record matching
// the provided key-value pair.
func (dao *Dao) FindFirstRecordByData(collection *models.Collection, key string, value any) (*models.Record, error) {
	row := dbx.NullStringMap{}

	err := dao.RecordQuery(collection).
		AndWhere(dbx.HashExp{key: value}).
		Limit(1).
		One(row)

	if err != nil {
		return nil, err
	}

	return models.NewRecordFromNullStringMap(collection, row), nil
}

// IsRecordValueUnique checks if the provided key-value pair is a unique Record value.
//
// NB! Array values (eg. from multiple select fields) are matched
// as a serialized json strings (eg. `["a","b"]`), so the value uniqueness
// depends on the elements order. Or in other words the following values
// are considered different: `[]string{"a","b"}` and `[]string{"b","a"}`
func (dao *Dao) IsRecordValueUnique(
	collection *models.Collection,
	key string,
	value any,
	excludeId string,
) bool {
	var exists bool

	var normalizedVal any
	switch val := value.(type) {
	case []string:
		normalizedVal = append(types.JsonArray{}, list.ToInterfaceSlice(val)...)
	case []any:
		normalizedVal = append(types.JsonArray{}, val...)
	default:
		normalizedVal = val
	}

	err := dao.RecordQuery(collection).
		Select("count(*)").
		AndWhere(dbx.Not(dbx.HashExp{"id": excludeId})).
		AndWhere(dbx.HashExp{key: normalizedVal}).
		Limit(1).
		Row(&exists)

	return err == nil && !exists
}

// FindUserRelatedRecords returns all records that has a reference
// to the provided User model (via the user shema field).
func (dao *Dao) FindUserRelatedRecords(user *models.User) ([]*models.Record, error) {
	if user.Id == "" {
		return []*models.Record{}, nil
	}

	collections, err := dao.FindCollectionsWithUserFields()
	if err != nil {
		return nil, err
	}

	result := []*models.Record{}
	for _, collection := range collections {
		userFields := []*schema.SchemaField{}

		// prepare fields options
		if err := collection.Schema.InitFieldsOptions(); err != nil {
			return nil, err
		}

		// extract user fields
		for _, field := range collection.Schema.Fields() {
			if field.Type == schema.FieldTypeUser {
				userFields = append(userFields, field)
			}
		}

		// fetch records associated to the user
		exprs := []dbx.Expression{}
		for _, field := range userFields {
			exprs = append(exprs, dbx.HashExp{field.Name: user.Id})
		}
		rows := []dbx.NullStringMap{}
		if err := dao.RecordQuery(collection).AndWhere(dbx.Or(exprs...)).All(&rows); err != nil {
			return nil, err
		}
		records := models.NewRecordsFromNullStringMaps(collection, rows)

		result = append(result, records...)
	}

	return result, nil
}

// SaveRecord upserts the provided Record model.
func (dao *Dao) SaveRecord(record *models.Record) error {
	return dao.Save(record)
}

// DeleteRecord deletes the provided Record model.
//
// This method will also cascade the delete operation to all linked
// relational records (delete or set to NULL, depending on the rel settings).
//
// The delete operation may fail if the record is part of a required
// reference in another record (aka. cannot be deleted or set to NULL).
func (dao *Dao) DeleteRecord(record *models.Record) error {
	// check for references
	// note: the select is outside of the transaction to prevent SQLITE_LOCKED error when mixing read&write in a single transaction
	refs, err := dao.FindCollectionReferences(record.Collection(), "")
	if err != nil {
		return err
	}

	// check if related records has to be deleted (if `CascadeDelete` is set)
	// OR
	// just unset the record id from any relation field values (if they are not required)
	// -----------------------------------------------------------
	return dao.RunInTransaction(func(txDao *Dao) error {
		for refCollection, fields := range refs {
			for _, field := range fields {
				options, _ := field.Options.(*schema.RelationOptions)

				rows := []dbx.NullStringMap{}

				// note: the select is not using the transaction dao to prevent SQLITE_LOCKED error when mixing read&write in a single transaction
				err := dao.RecordQuery(refCollection).
					AndWhere(dbx.Not(dbx.HashExp{"id": record.Id})).
					AndWhere(dbx.Like(field.Name, record.Id).Match(true, true)).
					All(&rows)
				if err != nil {
					return err
				}

				refRecords := models.NewRecordsFromNullStringMaps(refCollection, rows)
				for _, refRecord := range refRecords {
					ids := refRecord.GetStringSliceDataValue(field.Name)

					// unset the record id
					for i := len(ids) - 1; i >= 0; i-- {
						if ids[i] == record.Id {
							ids = append(ids[:i], ids[i+1:]...)
							break
						}
					}

					// cascade delete the reference
					// (only if there are no other active references in case of multiple select)
					if options.CascadeDelete && len(ids) == 0 {
						if err := txDao.DeleteRecord(refRecord); err != nil {
							return err
						}
						// no further action are needed (the reference is deleted)
						continue
					}

					if field.Required && len(ids) == 0 {
						return fmt.Errorf("The record cannot be deleted because it is part of a required reference in record %s (%s collection).", refRecord.Id, refCollection.Name)
					}

					// save the reference changes
					refRecord.SetDataValue(field.Name, field.PrepareValue(ids))
					if err := txDao.SaveRecord(refRecord); err != nil {
						return err
					}
				}
			}
		}

		return txDao.Delete(record)
	})
}

// SyncRecordTableSchema compares the two provided collections
// and applies the necessary related record table changes.
//
// If `oldCollection` is null, then only `newCollection` is used to create the record table.
func (dao *Dao) SyncRecordTableSchema(newCollection *models.Collection, oldCollection *models.Collection) error {
	// create
	if oldCollection == nil {
		cols := map[string]string{
			schema.ReservedFieldNameId:      "TEXT PRIMARY KEY",
			schema.ReservedFieldNameCreated: `TEXT DEFAULT "" NOT NULL`,
			schema.ReservedFieldNameUpdated: `TEXT DEFAULT "" NOT NULL`,
		}

		tableName := newCollection.Name

		// add schema field definitions
		for _, field := range newCollection.Schema.Fields() {
			cols[field.Name] = field.ColDefinition()
		}

		// create table
		_, tableErr := dao.DB().CreateTable(tableName, cols).Execute()
		if tableErr != nil {
			return tableErr
		}

		// add index on the base `created` column
		_, indexErr := dao.DB().CreateIndex(tableName, tableName+"_created_idx", "created").Execute()
		if indexErr != nil {
			return indexErr
		}

		return nil
	}

	// update
	return dao.RunInTransaction(func(txDao *Dao) error {
		oldTableName := oldCollection.Name
		newTableName := newCollection.Name
		oldSchema := oldCollection.Schema
		newSchema := newCollection.Schema

		// check for renamed table
		if !strings.EqualFold(oldTableName, newTableName) {
			_, err := dao.DB().RenameTable(oldTableName, newTableName).Execute()
			if err != nil {
				return err
			}
		}

		// check for deleted columns
		for _, oldField := range oldSchema.Fields() {
			if f := newSchema.GetFieldById(oldField.Id); f != nil {
				continue // exist
			}

			_, err := txDao.DB().DropColumn(newTableName, oldField.Name).Execute()
			if err != nil {
				return err
			}
		}

		// check for new or renamed columns
		toRename := map[string]string{}
		for _, field := range newSchema.Fields() {
			oldField := oldSchema.GetFieldById(field.Id)
			// Note:
			// We are using a temporary column name when adding or renaming columns
			// to ensure that there are no name collisions in case there is
			// names switch/reuse of existing columns (eg. name, title -> title, name).
			// This way we are always doing 1 more rename operation but it provides better dev experience.

			if oldField == nil {
				tempName := field.Name + security.RandomString(5)
				toRename[tempName] = field.Name

				// add
				_, err := txDao.DB().AddColumn(newTableName, tempName, field.ColDefinition()).Execute()
				if err != nil {
					return err
				}
			} else if oldField.Name != field.Name {
				tempName := field.Name + security.RandomString(5)
				toRename[tempName] = field.Name

				// rename
				_, err := txDao.DB().RenameColumn(newTableName, oldField.Name, tempName).Execute()
				if err != nil {
					return err
				}
			}
		}

		// set the actual columns name
		for tempName, actualName := range toRename {
			_, err := txDao.DB().RenameColumn(newTableName, tempName, actualName).Execute()
			if err != nil {
				return err
			}
		}

		return nil
	})
}
