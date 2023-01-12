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

	rows := make([]dbx.NullStringMap, 0, len(recordIds))
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

	if uniqueExcludeIds := list.NonzeroUniques(excludeIds); len(uniqueExcludeIds) > 0 {
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
		return nil, errors.New("missing or invalid token claims")
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
	if err != nil {
		return nil, fmt.Errorf("failed to fetch auth collection %q (%w)", collectionNameOrId, err)
	}
	if !collection.IsAuth() {
		return nil, fmt.Errorf("%q is not an auth collection", collectionNameOrId)
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
	if err != nil {
		return nil, fmt.Errorf("failed to fetch auth collection %q (%w)", collectionNameOrId, err)
	}
	if !collection.IsAuth() {
		return nil, fmt.Errorf("%q is not an auth collection", collectionNameOrId)
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
			return errors.New("unable to save auth record without username")
		}

		// Cross-check that the auth record id is unique for all auth collections.
		// This is to make sure that the filter `@request.auth.id` always returns a unique id.
		authCollections, err := dao.FindCollectionsByType(models.CollectionTypeAuth)
		if err != nil {
			return fmt.Errorf("unable to fetch the auth collections for cross-id unique check: %w", err)
		}
		for _, collection := range authCollections {
			if record.Collection().Id == collection.Id {
				continue // skip current collection (sqlite will do the check for us)
			}
			isUnique := dao.IsRecordValueUnique(collection.Id, schema.FieldNameId, record.Id)
			if !isUnique {
				return errors.New("the auth record ID must be unique across all auth collections")
			}
		}
	}

	return dao.Save(record)
}

// DeleteRecord deletes the provided Record model.
//
// This method will also cascade the delete operation to all linked
// relational records (delete or unset, depending on the rel settings).
//
// The delete operation may fail if the record is part of a required
// reference in another record (aka. cannot be deleted or unset).
func (dao *Dao) DeleteRecord(record *models.Record) error {
	// fetch rel references (if any)
	//
	// note: the select is outside of the transaction to minimize
	// SQLITE_BUSY errors when mixing read&write in a single transaction
	refs, err := dao.FindCollectionReferences(record.Collection())
	if err != nil {
		return err
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		// manually trigger delete on any linked external auth to ensure
		// that the `OnModel*` hooks are triggered
		if record.Collection().IsAuth() {
			// note: the select is outside of the transaction to minimize
			// SQLITE_BUSY errors when mixing read&write in a single transaction
			externalAuths, err := dao.FindAllExternalAuthsByRecord(record)
			if err != nil {
				return err
			}
			for _, auth := range externalAuths {
				if err := txDao.DeleteExternalAuth(auth); err != nil {
					return err
				}
			}
		}

		// delete the record before the relation references to ensure that there
		// will be no "A<->B" relations to prevent deadlock when calling DeleteRecord recursively
		if err := txDao.Delete(record); err != nil {
			return err
		}

		return txDao.cascadeRecordDelete(record, refs)
	})
}

// cascadeRecordDelete triggers cascade deletion for the provided references.
//
// NB! This method is expected to be called inside a transaction.
func (dao *Dao) cascadeRecordDelete(mainRecord *models.Record, refs map[*models.Collection][]*schema.SchemaField) error {
	uniqueJsonEachAlias := "__je__" + security.PseudorandomString(4)

	for refCollection, fields := range refs {
		for _, field := range fields {
			recordTableName := inflector.Columnify(refCollection.Name)
			prefixedFieldName := recordTableName + "." + inflector.Columnify(field.Name)

			// @todo optimize single relation lookup in v0.12+
			query := dao.RecordQuery(refCollection).
				Distinct(true).
				AndWhere(dbx.Not(dbx.HashExp{recordTableName + ".id": mainRecord.Id})).
				InnerJoin(fmt.Sprintf(
					// note: the case is used to normalize the value access
					`json_each(CASE WHEN json_valid([[%s]]) THEN [[%s]] ELSE json_array([[%s]]) END) as {{%s}}`,
					prefixedFieldName, prefixedFieldName, prefixedFieldName, uniqueJsonEachAlias,
				), dbx.HashExp{uniqueJsonEachAlias + ".value": mainRecord.Id})

			// trigger cascade for each batchSize rel items until there is none
			batchSize := 4000
			rows := make([]dbx.NullStringMap, 0, batchSize)
			for {
				if err := query.Limit(int64(batchSize)).All(&rows); err != nil {
					return err
				}

				total := len(rows)
				if total == 0 {
					break
				}

				refRecords := models.NewRecordsFromNullStringMaps(refCollection, rows)

				err := dao.deleteRefRecords(mainRecord, refRecords, field)
				if err != nil {
					return err
				}

				if total < batchSize {
					break // no more items
				}

				rows = rows[:0] // keep allocated memory
			}
		}
	}

	return nil
}

// deleteRefRecords checks if related records has to be deleted (if `CascadeDelete` is set)
// OR
// just unset the record id from any relation field values (if they are not required).
//
// NB! This method is expected to be called inside a transaction.
func (dao *Dao) deleteRefRecords(mainRecord *models.Record, refRecords []*models.Record, field *schema.SchemaField) error {
	options, _ := field.Options.(*schema.RelationOptions)
	if options == nil {
		return errors.New("relation field options are not initialized")
	}

	for _, refRecord := range refRecords {
		ids := refRecord.GetStringSlice(field.Name)

		// unset the record id
		for i := len(ids) - 1; i >= 0; i-- {
			if ids[i] == mainRecord.Id {
				ids = append(ids[:i], ids[i+1:]...)
				break
			}
		}

		// cascade delete the reference
		// (only if there are no other active references in case of multiple select)
		if options.CascadeDelete && len(ids) == 0 {
			if err := dao.DeleteRecord(refRecord); err != nil {
				return err
			}
			// no further actions are needed (the reference is deleted)
			continue
		}

		if field.Required && len(ids) == 0 {
			return fmt.Errorf("the record cannot be deleted because it is part of a required reference in record %s (%s collection)", refRecord.Id, refRecord.Collection().Name)
		}

		// save the reference changes
		refRecord.Set(field.Name, field.PrepareValue(ids))
		if err := dao.SaveRecord(refRecord); err != nil {
			return err
		}
	}

	return nil
}
