package daos_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
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
	updatedCollection.Indexes = types.JsonArray[string]{"create index idx_title_renamed on anything (title_renamed)"}

	scenarios := []struct {
		name                 string
		newCollection        *models.Collection
		oldCollection        *models.Collection
		expectedColumns      []string
		expectedIndexesCount int
	}{
		{
			"new base collection",
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
			[]string{"id", "created", "updated", "test"},
			0,
		},
		{
			"new auth collection",
			&models.Collection{
				Name: "new_table_auth",
				Type: models.CollectionTypeAuth,
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name: "test",
						Type: schema.FieldTypeText,
					},
				),
				Indexes: types.JsonArray[string]{"create index idx_auth_test on anything (email, username)"},
			},
			nil,
			[]string{
				"id", "created", "updated", "test",
				"username", "email", "verified", "emailVisibility",
				"tokenKey", "passwordHash", "lastResetSentAt", "lastVerificationSentAt",
			},
			4,
		},
		{
			"no changes",
			oldCollection,
			oldCollection,
			[]string{"id", "created", "updated", "title", "active"},
			3,
		},
		{
			"renamed table, deleted column, renamed columnd and new column",
			updatedCollection,
			oldCollection,
			[]string{"id", "created", "updated", "title_renamed", "new_field"},
			1,
		},
	}

	for _, s := range scenarios {
		err := app.Dao().SyncRecordTableSchema(s.newCollection, s.oldCollection)
		if err != nil {
			t.Errorf("[%s] %v", s.name, err)
			continue
		}

		if !app.Dao().HasTable(s.newCollection.Name) {
			t.Errorf("[%s] Expected table %s to exist", s.name, s.newCollection.Name)
		}

		cols, _ := app.Dao().TableColumns(s.newCollection.Name)
		if len(cols) != len(s.expectedColumns) {
			t.Errorf("[%s] Expected columns %v, got %v", s.name, s.expectedColumns, cols)
		}

		for _, c := range cols {
			if !list.ExistInSlice(c, s.expectedColumns) {
				t.Errorf("[%s] Couldn't find column %s in %v", s.name, c, s.expectedColumns)
			}
		}

		indexes, _ := app.Dao().TableIndexes(s.newCollection.Name)

		if totalIndexes := len(indexes); totalIndexes != s.expectedIndexesCount {
			t.Errorf("[%s] Expected %d indexes, got %d:\n%v", s.name, s.expectedIndexesCount, totalIndexes, indexes)
		}
	}
}

func TestSingleVsMultipleValuesNormalization(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	// mock field changes
	{
		selectOneField := collection.Schema.GetFieldByName("select_one")
		opt := selectOneField.Options.(*schema.SelectOptions)
		opt.MaxSelect = 2
	}
	{
		selectManyField := collection.Schema.GetFieldByName("select_many")
		opt := selectManyField.Options.(*schema.SelectOptions)
		opt.MaxSelect = 1
	}
	{
		fileOneField := collection.Schema.GetFieldByName("file_one")
		opt := fileOneField.Options.(*schema.FileOptions)
		opt.MaxSelect = 2
	}
	{
		fileManyField := collection.Schema.GetFieldByName("file_many")
		opt := fileManyField.Options.(*schema.FileOptions)
		opt.MaxSelect = 1
	}
	{
		relOneField := collection.Schema.GetFieldByName("rel_one")
		opt := relOneField.Options.(*schema.RelationOptions)
		opt.MaxSelect = types.Pointer(2)
	}
	{
		relManyField := collection.Schema.GetFieldByName("rel_many")
		opt := relManyField.Options.(*schema.RelationOptions)
		opt.MaxSelect = types.Pointer(1)
	}

	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	type expectation struct {
		SelectOne  string `db:"select_one"`
		SelectMany string `db:"select_many"`
		FileOne    string `db:"file_one"`
		FileMany   string `db:"file_many"`
		RelOne     string `db:"rel_one"`
		RelMany    string `db:"rel_many"`
	}

	scenarios := []struct {
		recordId string
		expected expectation
	}{
		{
			"imy661ixudk5izi",
			expectation{
				SelectOne:  `[]`,
				SelectMany: ``,
				FileOne:    `[]`,
				FileMany:   ``,
				RelOne:     `[]`,
				RelMany:    ``,
			},
		},
		{
			"al1h9ijdeojtsjy",
			expectation{
				SelectOne:  `["optionB"]`,
				SelectMany: `optionB`,
				FileOne:    `["300_Jsjq7RdBgA.png"]`,
				FileMany:   ``,
				RelOne:     `["84nmscqy84lsi1t"]`,
				RelMany:    `oap640cot4yru2s`,
			},
		},
		{
			"84nmscqy84lsi1t",
			expectation{
				SelectOne:  `["optionB"]`,
				SelectMany: `optionC`,
				FileOne:    `["test_d61b33QdDU.txt"]`,
				FileMany:   `test_tC1Yc87DfC.txt`,
				RelOne:     `[]`,
				RelMany:    `oap640cot4yru2s`,
			},
		},
	}

	for _, s := range scenarios {
		result := new(expectation)

		err := app.Dao().DB().Select(
			"select_one",
			"select_many",
			"file_one",
			"file_many",
			"rel_one",
			"rel_many",
		).From(collection.Name).Where(dbx.HashExp{"id": s.recordId}).One(result)
		if err != nil {
			t.Errorf("[%s] Failed to load record: %v", s.recordId, err)
			continue
		}

		encodedResult, err := json.Marshal(result)
		if err != nil {
			t.Errorf("[%s] Failed to encode result: %v", s.recordId, err)
			continue
		}

		encodedExpectation, err := json.Marshal(s.expected)
		if err != nil {
			t.Errorf("[%s] Failed to encode expectation: %v", s.recordId, err)
			continue
		}

		if !bytes.EqualFold(encodedExpectation, encodedResult) {
			t.Errorf("[%s] Expected \n%s, \ngot \n%s", s.recordId, encodedExpectation, encodedResult)
		}
	}
}
