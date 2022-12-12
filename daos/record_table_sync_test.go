package daos_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestSyncRecordTableSchema(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	oldCollection, err := app.Dao().FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}
	updatedCollection, err := app.Dao().FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}
	updatedCollection.Name = "demo_renamed"
	updatedCollection.Schema.RemoveField(updatedCollection.Schema.GetFieldByName("active").Id)
	updatedCollection.Schema.AddField(
		&schema.SchemaField{
			Name: "new_field",
			Type: schema.FieldTypeEmail,
		},
	)
	updatedCollection.Schema.AddField(
		&schema.SchemaField{
			Id:   updatedCollection.Schema.GetFieldByName("title").Id,
			Name: "title_renamed",
			Type: schema.FieldTypeEmail,
		},
	)

	scenarios := []struct {
		newCollection     *models.Collection
		oldCollection     *models.Collection
		expectedTableName string
		expectedColumns   []string
	}{
		// new base collection
		{
			&models.Collection{
				Name: "new_table",
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name: "test",
						Type: schema.FieldTypeText,
					},
				),
			},
			nil,
			"new_table",
			[]string{"id", "created", "updated", "test"},
		},
		// new auth collection
		{
			&models.Collection{
				Name: "new_table_auth",
				Type: models.CollectionTypeAuth,
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name: "test",
						Type: schema.FieldTypeText,
					},
				),
			},
			nil,
			"new_table_auth",
			[]string{
				"id", "created", "updated", "test",
				"username", "email", "verified", "emailVisibility",
				"tokenKey", "passwordHash", "lastResetSentAt", "lastVerificationSentAt",
			},
		},
		// no changes
		{
			oldCollection,
			oldCollection,
			"demo3",
			[]string{"id", "created", "updated", "title", "active"},
		},
		// renamed table, deleted column, renamed columnd and new column
		{
			updatedCollection,
			oldCollection,
			"demo_renamed",
			[]string{"id", "created", "updated", "title_renamed", "new_field"},
		},
	}

	for i, scenario := range scenarios {
		err := app.Dao().SyncRecordTableSchema(scenario.newCollection, scenario.oldCollection)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
			continue
		}

		if !app.Dao().HasTable(scenario.newCollection.Name) {
			t.Errorf("(%d) Expected table %s to exist", i, scenario.newCollection.Name)
		}

		cols, _ := app.Dao().GetTableColumns(scenario.newCollection.Name)
		if len(cols) != len(scenario.expectedColumns) {
			t.Errorf("(%d) Expected columns %v, got %v", i, scenario.expectedColumns, cols)
		}

		for _, c := range cols {
			if !list.ExistInSlice(c, scenario.expectedColumns) {
				t.Errorf("(%d) Couldn't find column %s in %v", i, c, scenario.expectedColumns)
			}
		}
	}
}
