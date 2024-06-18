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
	t.Parallel()

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
				"tokenKey", "passwordHash", "lastResetSentAt", "lastVerificationSentAt", "lastLoginAlertSentAt",
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
	t.Parallel()

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
	{
		// new multivaluer field to check whether the array normalization
		// will be applied for already inserted data
		collection.Schema.AddField(&schema.SchemaField{
			Name: "new_multiple",
			Type: schema.FieldTypeSelect,
			Options: &schema.SelectOptions{
				Values:    []string{"a", "b", "c"},
				MaxSelect: 3,
			},
		})
	}

	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// ensures that the writable schema was reverted to its expected default
	var writableSchema bool
	app.Dao().DB().NewQuery("PRAGMA writable_schema").Row(&writableSchema)
	if writableSchema == true {
		t.Fatalf("Expected writable_schema to be OFF, got %v", writableSchema)
	}

	// check whether the columns DEFAULT definition was updated
	// ---------------------------------------------------------------
	tableInfo, err := app.Dao().TableInfo(collection.Name)
	if err != nil {
		t.Fatal(err)
	}

	tableInfoExpectations := map[string]string{
		"select_one":   `'[]'`,
		"select_many":  `''`,
		"file_one":     `'[]'`,
		"file_many":    `''`,
		"rel_one":      `'[]'`,
		"rel_many":     `''`,
		"new_multiple": `'[]'`,
	}
	for col, dflt := range tableInfoExpectations {
		t.Run("check default for "+col, func(t *testing.T) {
			var row *models.TableInfoRow
			for _, r := range tableInfo {
				if r.Name == col {
					row = r
					break
				}
			}
			if row == nil {
				t.Fatalf("Missing info for column %q", col)
			}

			if v := row.DefaultValue.String(); v != dflt {
				t.Fatalf("Expected default value %q, got %q", dflt, v)
			}
		})
	}

	// check whether the values were normalized
	// ---------------------------------------------------------------
	type fieldsExpectation struct {
		SelectOne   string `db:"select_one"`
		SelectMany  string `db:"select_many"`
		FileOne     string `db:"file_one"`
		FileMany    string `db:"file_many"`
		RelOne      string `db:"rel_one"`
		RelMany     string `db:"rel_many"`
		NewMultiple string `db:"new_multiple"`
	}

	fieldsScenarios := []struct {
		recordId string
		expected fieldsExpectation
	}{
		{
			"imy661ixudk5izi",
			fieldsExpectation{
				SelectOne:   `[]`,
				SelectMany:  ``,
				FileOne:     `[]`,
				FileMany:    ``,
				RelOne:      `[]`,
				RelMany:     ``,
				NewMultiple: `[]`,
			},
		},
		{
			"al1h9ijdeojtsjy",
			fieldsExpectation{
				SelectOne:   `["optionB"]`,
				SelectMany:  `optionB`,
				FileOne:     `["300_Jsjq7RdBgA.png"]`,
				FileMany:    ``,
				RelOne:      `["84nmscqy84lsi1t"]`,
				RelMany:     `oap640cot4yru2s`,
				NewMultiple: `[]`,
			},
		},
		{
			"84nmscqy84lsi1t",
			fieldsExpectation{
				SelectOne:   `["optionB"]`,
				SelectMany:  `optionC`,
				FileOne:     `["test_d61b33QdDU.txt"]`,
				FileMany:    `test_tC1Yc87DfC.txt`,
				RelOne:      `[]`,
				RelMany:     `oap640cot4yru2s`,
				NewMultiple: `[]`,
			},
		},
	}

	for _, s := range fieldsScenarios {
		t.Run("check fields for record "+s.recordId, func(t *testing.T) {
			result := new(fieldsExpectation)

			err := app.Dao().DB().Select(
				"select_one",
				"select_many",
				"file_one",
				"file_many",
				"rel_one",
				"rel_many",
				"new_multiple",
			).From(collection.Name).Where(dbx.HashExp{"id": s.recordId}).One(result)
			if err != nil {
				t.Fatalf("Failed to load record: %v", err)
			}

			encodedResult, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Failed to encode result: %v", err)
			}

			encodedExpectation, err := json.Marshal(s.expected)
			if err != nil {
				t.Fatalf("Failed to encode expectation: %v", err)
			}

			if !bytes.EqualFold(encodedExpectation, encodedResult) {
				t.Fatalf("Expected \n%s, \ngot \n%s", encodedExpectation, encodedResult)
			}
		})
	}
}
