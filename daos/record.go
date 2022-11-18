package daos

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

// RecordQuery returns a new Record select query.
func (dao *Dao) RecordQuery(collection *models.Collection) *dbx.SelectQuery {
	tableName := collection.Name
	selectCols := fmt.Sprintf("%s.*", dao.DB().QuoteSimpleColumnName(tableName))

	return dao.DB().Select(selectCols).From(tableName)
}

// FindRecordById finds the Record model by its id.
func (dao *Dao) FindRecordById(
	collectionNameOrId string,
	recordId string,
	optFilters ...func(q *dbx.SelectQuery) error,
) (*models.Record, error) {
	collection, err := dao.FindCollectionByNameOrId(collectionNameOrId)
	if err != nil {
		return nil, err
	}

	tableName := collection.Name

	query := dao.RecordQuery(collection).
		AndWhere(dbx.HashExp{tableName + ".id": recordId})

	for _, filter := range optFilters {
		if filter == nil {
			continue
		}
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
	collectionNameOrId string,
	recordIds []string,
	optFilters ...func(q *dbx.SelectQuery) error,
) ([]*models.Record, error) {
	collection, err := dao.FindCollectionByNameOrId(collectionNameOrId)
	if err != nil {
		return nil, err
	}

	query := dao.RecordQuery(collection).
		AndWhere(dbx.In(
			collection.Name+".id",
			list.ToInterfaceSlice(recordIds)...,
		))

	for _, filter := range optFilters {
		if filter == nil {
			continue
		}
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

// FindRecordsByExpr finds all records by the specified db expression.
//
// Returns all collection records if no expressions are provided.
//
// Returns an empty slice if no records are found.
//
// Example:
//	expr1 := dbx.HashExp{"email": "test@example.com"}
//	expr2 := dbx.NewExp("LOWER(username) = {:username}", dbx.Params{"username": "test"})
//	dao.FindRecordsByExpr("example", expr1, expr2)
func (dao *Dao) FindRecordsByExpr(collectionNameOrId string, exprs ...dbx.Expression) ([]*models.Record, error) {
	collection, err := dao.FindCollectionByNameOrId(collectionNameOrId)
	if err != nil {
		return nil, err
	}

	query := dao.RecordQuery(collection)

	// add only the non-nil expressions
	for _, expr := range exprs {
		if expr != nil {
			query.AndWhere(expr)
		}
	}

	rows := []dbx.NullStringMap{}

	if err := query.All(&rows); err != nil {
		return nil, err
	}

	return models.NewRecordsFromNullStringMaps(collection, rows), nil
}

// FindFirstRecordByData returns the first found record matching
// the provided key-value pair.
func (dao *Dao) FindFirstRecordByData(collectionNameOrId string, key string, value any) (*models.Record, error) {
	collection, err := dao.FindCollectionByNameOrId(collectionNameOrId)
	if err != nil {
		return nil, err
	}

	row := dbx.NullStringMap{}

	err = dao.RecordQuery(collection).
		AndWhere(dbx.HashExp{inflector.Columnify(key): value}).
		Limit(1).
		One(row)

	if err != nil {
		return nil, err
	}

	return models.NewRecordFromNullStringMap(collection, row), nil
}

// IsRecordValueUnique checks if the provided key-value pair is a unique Record value.
//
// For correctness, if the collection is "auth" and the key is "username",
// the unique check will be case insensitive.
//
// NB! Array values (eg. from multiple select fields) are matched
// as a serialized json strings (eg. `["a","b"]`), so the value uniqueness
// depends on the elements order. Or in other words the following values
// are considered different: `[]string{"a","b"}` and `[]string{"b","a"}`
func (dao *Dao) IsRecordValueUnique(
	collectionNameOrId string,
	key string,
	value any,
	excludeIds ...string,
) bool {
	collection, err := dao.FindCollectionByNameOrId(collectionNameOrId)
	if err != nil {
		return false
	}

	var expr dbx.Expression
	if collection.IsAuth() && key == schema.FieldNameUsername {
		expr = dbx.NewExp("LOWER([["+schema.FieldNameUsername+"]])={:username}", dbx.Params{
			"username": strings.ToLower(cast.ToString(value)),
		})
	} else {
		var normalizedVal any
		switch val := value.(type) {
		case []string:
			normalizedVal = append(types.JsonArray{}, list.ToInterfaceSlice(val)...)
		case []any:
			normalizedVal = append(types.JsonArray{}, val...)
		default:
			normalizedVal = val
		}

		expr = dbx.HashExp{inflector.Columnify(key): normalizedVal}
	}

	query := dao.RecordQuery(collection).
		Select("count(*)").
		AndWhere(expr).
		Limit(1)

	if len(excludeIds) > 0 {
		uniqueExcludeIds := list.NonzeroUniques(excludeIds)
		query.AndWhere(dbx.NotIn(collection.Name+".id", list.ToInterfaceSlice(uniqueExcludeIds)...))
	}

	var exists bool

	return query.Row(&exists) == nil && !exists
}

// FindAuthRecordByToken finds the auth record associated with the provided JWT token.
//
// Returns an error if the JWT token is invalid, expired or not associated to an auth collection record.
func (dao *Dao) FindAuthRecordByToken(token string, baseTokenKey string) (*models.Record, error) {
	unverifiedClaims, err := security.ParseUnverifiedJWT(token)
	if err != nil {
		return nil, err
	}

	// check required claims
	id, _ := unverifiedClaims["id"].(string)
	collectionId, _ := unverifiedClaims["collectionId"].(string)
	if id == "" || collectionId == "" {
		return nil, errors.New("Missing or invalid token claims.")
	}

	record, err := dao.FindRecordById(collectionId, id)
	if err != nil {
		return nil, err
	}

	if !record.Collection().IsAuth() {
		return nil, errors.New("The token is not associated to an auth collection record.")
	}

	verificationKey := record.TokenKey() + baseTokenKey

	// verify token signature
	if _, err := security.ParseJWT(token, verificationKey); err != nil {
		return nil, err
	}

	return record, nil
}

// FindAuthRecordByEmail finds the auth record associated with the provided email.
//
// Returns an error if it is not an auth collection or the record is not found.
func (dao *Dao) FindAuthRecordByEmail(collectionNameOrId string, email string) (*models.Record, error) {
	collection, err := dao.FindCollectionByNameOrId(collectionNameOrId)
	if err != nil || !collection.IsAuth() {
		return nil, errors.New("Missing or not an auth collection.")
	}

	row := dbx.NullStringMap{}

	err = dao.RecordQuery(collection).
		AndWhere(dbx.HashExp{schema.FieldNameEmail: email}).
		Limit(1).
		One(row)

	if err != nil {
		return nil, err
	}

	return models.NewRecordFromNullStringMap(collection, row), nil
}

// FindAuthRecordByUsername finds the auth record associated with the provided username (case insensitive).
//
// Returns an error if it is not an auth collection or the record is not found.
func (dao *Dao) FindAuthRecordByUsername(collectionNameOrId string, username string) (*models.Record, error) {
	collection, err := dao.FindCollectionByNameOrId(collectionNameOrId)
	if err != nil || !collection.IsAuth() {
		return nil, errors.New("Missing or not an auth collection.")
	}

	row := dbx.NullStringMap{}

	err = dao.RecordQuery(collection).
		AndWhere(dbx.NewExp("LOWER([["+schema.FieldNameUsername+"]])={:username}", dbx.Params{
			"username": strings.ToLower(username),
		})).
		Limit(1).
		One(row)

	if err != nil {
		return nil, err
	}

	return models.NewRecordFromNullStringMap(collection, row), nil
}

// SuggestUniqueAuthRecordUsername checks if the provided username is unique
// and return a new "unique" username with appended random numeric part
// (eg. "existingName" -> "existingName583").
//
// The same username will be returned if the provided string is already unique.
func (dao *Dao) SuggestUniqueAuthRecordUsername(
	collectionNameOrId string,
	baseUsername string,
	excludeIds ...string,
) string {
	username := baseUsername

	for i := 0; i < 10; i++ { // max 10 attempts
		isUnique := dao.IsRecordValueUnique(
			collectionNameOrId,
			schema.FieldNameUsername,
			username,
			excludeIds...,
		)
		if isUnique {
			break // already unique
		}
		username = baseUsername + security.RandomStringWithAlphabet(3+i, "123456789")
	}

	return username
}

// SaveRecord upserts the provided Record model.
func (dao *Dao) SaveRecord(record *models.Record) error {
	if record.Collection().IsAuth() {
		if record.Username() == "" {
			return errors.New("Unable to save auth record without username.")
		}

		// Cross-check that the auth record id is unique for all auth collections.
		// This is to make sure that the filter `@request.auth.id` always returns a unique id.
		authCollections, err := dao.FindCollectionsByType(models.CollectionTypeAuth)
		if err != nil {
			return fmt.Errorf("Unable to fetch the auth collections for cross-id unique check: %v", err)
		}
		for _, collection := range authCollections {
			if record.Collection().Id == collection.Id {
				continue // skip current collection (sqlite will do the check for us)
			}
			isUnique := dao.IsRecordValueUnique(collection.Id, schema.FieldNameId, record.Id)
			if !isUnique {
				return errors.New("The auth record ID must be unique across all auth collections.")
			}
		}
	}

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
	// note: the select is outside of the transaction to prevent SQLITE_LOCKED error when mixing read&write in a single transaction.
	refs, err := dao.FindCollectionReferences(record.Collection())
	if err != nil {
		return err
	}

	// check if related records has to be deleted (if `CascadeDelete` is set)
	// OR
	// just unset the record id from any relation field values (if they are not required)
	// -----------------------------------------------------------
	return dao.RunInTransaction(func(txDao *Dao) error {
		// delete/update references
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
					ids := refRecord.GetStringSlice(field.Name)

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
					refRecord.Set(field.Name, field.PrepareValue(ids))
					if err := txDao.SaveRecord(refRecord); err != nil {
						return err
					}
				}
			}
		}

		// delete linked external auths
		if record.Collection().IsAuth() {
			_, err = txDao.DB().Delete((&models.ExternalAuth{}).TableName(), dbx.HashExp{
				"collectionId": record.Collection().Id,
				"recordId":     record.Id,
			}).Execute()
			if err != nil {
				return err
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
			schema.FieldNameId:      "TEXT PRIMARY KEY",
			schema.FieldNameCreated: "TEXT DEFAULT '' NOT NULL",
			schema.FieldNameUpdated: "TEXT DEFAULT '' NOT NULL",
		}

		if newCollection.IsAuth() {
			cols[schema.FieldNameUsername] = "TEXT NOT NULL"
			cols[schema.FieldNameEmail] = "TEXT DEFAULT '' NOT NULL"
			cols[schema.FieldNameEmailVisibility] = "BOOLEAN DEFAULT FALSE NOT NULL"
			cols[schema.FieldNameVerified] = "BOOLEAN DEFAULT FALSE NOT NULL"
			cols[schema.FieldNameTokenKey] = "TEXT NOT NULL"
			cols[schema.FieldNamePasswordHash] = "TEXT NOT NULL"
			cols[schema.FieldNameLastResetSentAt] = "TEXT DEFAULT '' NOT NULL"
			cols[schema.FieldNameLastVerificationSentAt] = "TEXT DEFAULT '' NOT NULL"
		}

		// ensure that the new collection has an id
		if !newCollection.HasId() {
			newCollection.RefreshId()
			newCollection.MarkAsNew()
		}

		tableName := newCollection.Name

		// add schema field definitions
		for _, field := range newCollection.Schema.Fields() {
			cols[field.Name] = field.ColDefinition()
		}

		// create table
		if _, err := dao.DB().CreateTable(tableName, cols).Execute(); err != nil {
			return err
		}

		// add named index on the base `created` column
		if _, err := dao.DB().CreateIndex(tableName, "_"+newCollection.Id+"_created_idx", "created").Execute(); err != nil {
			return err
		}

		// add named unique index on the email and tokenKey columns
		if newCollection.IsAuth() {
			_, err := dao.DB().NewQuery(fmt.Sprintf(
				`
				CREATE UNIQUE INDEX _%s_username_idx ON {{%s}} ([[username]]);
				CREATE UNIQUE INDEX _%s_email_idx ON {{%s}} ([[email]]) WHERE [[email]] != '';
				CREATE UNIQUE INDEX _%s_tokenKey_idx ON {{%s}} ([[tokenKey]]);
				`,
				newCollection.Id, tableName,
				newCollection.Id, tableName,
				newCollection.Id, tableName,
			)).Execute()
			if err != nil {
				return err
			}
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
			_, err := txDao.DB().RenameTable(oldTableName, newTableName).Execute()
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
				tempName := field.Name + security.PseudorandomString(5)
				toRename[tempName] = field.Name

				// add
				_, err := txDao.DB().AddColumn(newTableName, tempName, field.ColDefinition()).Execute()
				if err != nil {
					return err
				}
			} else if oldField.Name != field.Name {
				tempName := field.Name + security.PseudorandomString(5)
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
