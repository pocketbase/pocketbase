package daos

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
)

// SyncRecordTableSchema compares the two provided collections
// and applies the necessary related record table changes.
//
// If `oldCollection` is null, then only `newCollection` is used to create the record table.
func (dao *Dao) SyncRecordTableSchema(newCollection *models.Collection, oldCollection *models.Collection) error {
	// create
	if oldCollection == nil {
		cols := map[string]string{
			schema.FieldNameId:      "TEXT PRIMARY KEY NOT NULL",
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
		deletedFieldNames := []string{}
		renamedFieldNames := map[string]string{}

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

		return txDao.syncCollectionReferences(newCollection, renamedFieldNames, deletedFieldNames)
	})
}

func (dao *Dao) syncCollectionReferences(collection *models.Collection, renamedFieldNames map[string]string, deletedFieldNames []string) error {
	if len(renamedFieldNames) == 0 && len(deletedFieldNames) == 0 {
		return nil // nothing to sync
	}

	refs, err := dao.FindCollectionReferences(collection)
	if err != nil {
		return err
	}

	for refCollection, refFields := range refs {
		for _, refField := range refFields {
			options, _ := refField.Options.(*schema.RelationOptions)
			if options == nil {
				continue
			}

			// remove deleted (if any)
			newDisplayFields := list.SubtractSlice(options.DisplayFields, deletedFieldNames)

			for old, new := range renamedFieldNames {
				for i, name := range newDisplayFields {
					if name == old {
						newDisplayFields[i] = new
					}
				}
			}

			// has changes
			if len(list.SubtractSlice(options.DisplayFields, newDisplayFields)) > 0 {
				options.DisplayFields = newDisplayFields

				// direct collection save to prevent self-referencing
				// recursion and unnecessary records table sync checks
				if err := dao.Save(refCollection); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
