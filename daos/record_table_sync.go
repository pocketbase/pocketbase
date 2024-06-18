package daos

import (
	"fmt"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/security"
)

// SyncRecordTableSchema compares the two provided collections
// and applies the necessary related record table changes.
//
// If `oldCollection` is null, then only `newCollection` is used to create the record table.
func (dao *Dao) SyncRecordTableSchema(newCollection *models.Collection, oldCollection *models.Collection) error {
	return dao.RunInTransaction(func(txDao *Dao) error {
		// create
		// -----------------------------------------------------------
		if oldCollection == nil {
			cols := map[string]string{
				schema.FieldNameId:      "TEXT PRIMARY KEY DEFAULT ('r'||lower(hex(randomblob(7)))) NOT NULL",
				schema.FieldNameCreated: "TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL",
				schema.FieldNameUpdated: "TEXT DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ')) NOT NULL",
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
				cols[schema.FieldNameLastLoginAlertSentAt] = "TEXT DEFAULT '' NOT NULL"
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
			if _, err := txDao.DB().CreateTable(tableName, cols).Execute(); err != nil {
				return err
			}

			// add named unique index on the email and tokenKey columns
			if newCollection.IsAuth() {
				_, err := txDao.DB().NewQuery(fmt.Sprintf(
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

			return txDao.createCollectionIndexes(newCollection)
		}

		// update
		// -----------------------------------------------------------
		oldTableName := oldCollection.Name
		newTableName := newCollection.Name
		oldSchema := oldCollection.Schema
		newSchema := newCollection.Schema
		deletedFieldNames := []string{}
		renamedFieldNames := map[string]string{}

		// drop old indexes (if any)
		if err := txDao.dropCollectionIndex(oldCollection); err != nil {
			return err
		}

		// check for renamed table
		if !strings.EqualFold(oldTableName, newTableName) {
			_, err := txDao.DB().RenameTable("{{"+oldTableName+"}}", "{{"+newTableName+"}}").Execute()
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
				return fmt.Errorf("failed to drop column %s - %w", oldField.Name, err)
			}

			deletedFieldNames = append(deletedFieldNames, oldField.Name)
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
					return fmt.Errorf("failed to add column %s - %w", field.Name, err)
				}
			} else if oldField.Name != field.Name {
				tempName := field.Name + security.PseudorandomString(5)
				toRename[tempName] = field.Name

				// rename
				_, err := txDao.DB().RenameColumn(newTableName, oldField.Name, tempName).Execute()
				if err != nil {
					return fmt.Errorf("failed to rename column %s - %w", oldField.Name, err)
				}

				renamedFieldNames[oldField.Name] = field.Name
			}
		}

		// set the actual columns name
		for tempName, actualName := range toRename {
			_, err := txDao.DB().RenameColumn(newTableName, tempName, actualName).Execute()
			if err != nil {
				return err
			}
		}

		if err := txDao.normalizeSingleVsMultipleFieldChanges(newCollection, oldCollection); err != nil {
			return err
		}

		return txDao.createCollectionIndexes(newCollection)
	})
}

func (dao *Dao) normalizeSingleVsMultipleFieldChanges(newCollection, oldCollection *models.Collection) error {
	if newCollection.IsView() || oldCollection == nil {
		return nil // view or not an update
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		// temporary disable the schema error checks to prevent view and trigger errors
		// when "altering" (aka. deleting and recreating) the non-normalized columns
		if _, err := txDao.DB().NewQuery("PRAGMA writable_schema = ON").Execute(); err != nil {
			return err
		}
		// executed with defer to make sure that the pragma is always reverted
		// in case of an error and when nested transactions are used
		defer txDao.DB().NewQuery("PRAGMA writable_schema = RESET").Execute()

		for _, newField := range newCollection.Schema.Fields() {
			// allow to continue even if there is no old field for the cases
			// when a new field is added and there are already inserted data
			var isOldMultiple bool
			if oldField := oldCollection.Schema.GetFieldById(newField.Id); oldField != nil {
				if opt, ok := oldField.Options.(schema.MultiValuer); ok {
					isOldMultiple = opt.IsMultiple()
				}
			}

			var isNewMultiple bool
			if opt, ok := newField.Options.(schema.MultiValuer); ok {
				isNewMultiple = opt.IsMultiple()
			}

			if isOldMultiple == isNewMultiple {
				continue // no change
			}

			// update the column definition by:
			// 1. inserting a new column with the new definition
			// 2. copy normalized values from the original column to the new one
			// 3. drop the original column
			// 4. rename the new column to the original column
			// -------------------------------------------------------

			originalName := newField.Name
			tempName := "_" + newField.Name + security.PseudorandomString(5)

			_, err := txDao.DB().AddColumn(newCollection.Name, tempName, newField.ColDefinition()).Execute()
			if err != nil {
				return err
			}

			var copyQuery *dbx.Query

			if !isOldMultiple && isNewMultiple {
				// single -> multiple (convert to array)
				copyQuery = txDao.DB().NewQuery(fmt.Sprintf(
					`UPDATE {{%s}} set [[%s]] = (
							CASE
								WHEN COALESCE([[%s]], '') = ''
								THEN '[]'
								ELSE (
									CASE
										WHEN json_valid([[%s]]) AND json_type([[%s]]) == 'array'
										THEN [[%s]]
										ELSE json_array([[%s]])
									END
								)
							END
						)`,
					newCollection.Name,
					tempName,
					originalName,
					originalName,
					originalName,
					originalName,
					originalName,
				))
			} else {
				// multiple -> single (keep only the last element)
				//
				// note: for file fields the actual file objects are not
				// deleted allowing additional custom handling via migration
				copyQuery = txDao.DB().NewQuery(fmt.Sprintf(
					`UPDATE {{%s}} set [[%s]] = (
						CASE
							WHEN COALESCE([[%s]], '[]') = '[]'
							THEN ''
							ELSE (
								CASE
									WHEN json_valid([[%s]]) AND json_type([[%s]]) == 'array'
									THEN COALESCE(json_extract([[%s]], '$[#-1]'), '')
									ELSE [[%s]]
								END
							)
						END
					)`,
					newCollection.Name,
					tempName,
					originalName,
					originalName,
					originalName,
					originalName,
					originalName,
				))
			}

			// copy the normalized values
			if _, err := copyQuery.Execute(); err != nil {
				return err
			}

			// drop the original column
			if _, err := txDao.DB().DropColumn(newCollection.Name, originalName).Execute(); err != nil {
				return err
			}

			// rename the new column back to the original
			if _, err := txDao.DB().RenameColumn(newCollection.Name, tempName, originalName).Execute(); err != nil {
				return err
			}
		}

		// revert the pragma and reload the schema
		_, revertErr := txDao.DB().NewQuery("PRAGMA writable_schema = RESET").Execute()

		return revertErr
	})
}

func (dao *Dao) dropCollectionIndex(collection *models.Collection) error {
	if collection.IsView() {
		return nil // views don't have indexes
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		for _, raw := range collection.Indexes {
			parsed := dbutils.ParseIndex(raw)

			if !parsed.IsValid() {
				continue
			}

			if _, err := txDao.DB().NewQuery(fmt.Sprintf("DROP INDEX IF EXISTS [[%s]]", parsed.IndexName)).Execute(); err != nil {
				return err
			}
		}

		return nil
	})
}

func (dao *Dao) createCollectionIndexes(collection *models.Collection) error {
	if collection.IsView() {
		return nil // views don't have indexes
	}

	return dao.RunInTransaction(func(txDao *Dao) error {
		// drop new indexes in case a duplicated index name is used
		if err := txDao.dropCollectionIndex(collection); err != nil {
			return err
		}

		// upsert new indexes
		//
		// note: we are returning validation errors because the indexes cannot be
		//       validated in a form, aka. before persisting the related collection
		//       record table changes
		errs := validation.Errors{}
		for i, idx := range collection.Indexes {
			parsed := dbutils.ParseIndex(idx)

			// ensure that the index is always for the current collection
			parsed.TableName = collection.Name

			if !parsed.IsValid() {
				errs[strconv.Itoa(i)] = validation.NewError(
					"validation_invalid_index_expression",
					"Invalid CREATE INDEX expression.",
				)
				continue
			}

			if _, err := txDao.DB().NewQuery(parsed.Build()).Execute(); err != nil {
				errs[strconv.Itoa(i)] = validation.NewError(
					"validation_invalid_index_expression",
					fmt.Sprintf("Failed to create index %s - %v.", parsed.IndexName, err.Error()),
				)
				continue
			}
		}

		if len(errs) > 0 {
			return validation.Errors{"indexes": errs}
		}

		return nil
	})
}
