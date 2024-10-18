package core_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestSyncRecordTableSchema(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	oldCollection, err := app.FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}
	updatedCollection, err := app.FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}
	updatedCollection.Name = "demo_renamed"
	updatedCollection.Fields.RemoveByName("active")
	updatedCollection.Fields.Add(&core.EmailField{
		Name: "new_field",
	})
	updatedCollection.Fields.Add(&core.EmailField{
		Id:   updatedCollection.Fields.GetByName("title").GetId(),
		Name: "title_renamed",
	})
	updatedCollection.Indexes = types.JSONArray[string]{"create index idx_title_renamed on anything (title_renamed)"}

	baseCol := core.NewBaseCollection("new_base")
	baseCol.Fields.Add(&core.TextField{Name: "test"})

	authCol := core.NewAuthCollection("new_auth")
	authCol.Fields.Add(&core.TextField{Name: "test"})
	authCol.AddIndex("idx_auth_test", false, "email, id", "")

	scenarios := []struct {
		name                 string
		newCollection        *core.Collection
		oldCollection        *core.Collection
		expectedColumns      []string
		expectedIndexesCount int
	}{
		{
			"new base collection",
			baseCol,
			nil,
			[]string{"id", "test"},
			0,
		},
		{
			"new auth collection",
			authCol,
			nil,
			[]string{
				"id", "test", "email", "verified",
				"emailVisibility", "tokenKey", "password",
			},
			3,
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
		t.Run(s.name, func(t *testing.T) {
			err := app.SyncRecordTableSchema(s.newCollection, s.oldCollection)
			if err != nil {
				t.Fatal(err)
			}

			if !app.HasTable(s.newCollection.Name) {
				t.Fatalf("Expected table %s to exist", s.newCollection.Name)
			}

			cols, _ := app.TableColumns(s.newCollection.Name)
			if len(cols) != len(s.expectedColumns) {
				t.Fatalf("Expected columns %v, got %v", s.expectedColumns, cols)
			}

			for _, col := range cols {
				if !list.ExistInSlice(col, s.expectedColumns) {
					t.Fatalf("Couldn't find column %s in %v", col, s.expectedColumns)
				}
			}

			indexes, _ := app.TableIndexes(s.newCollection.Name)

			if totalIndexes := len(indexes); totalIndexes != s.expectedIndexesCount {
				t.Fatalf("Expected %d indexes, got %d:\n%v", s.expectedIndexesCount, totalIndexes, indexes)
			}
		})
	}
}

func getTotalViews(app core.App) (int, error) {
	var total int

	err := app.DB().Select("count(*)").
		From("sqlite_master").
		AndWhere(dbx.NewExp("sql is not null")).
		AndWhere(dbx.HashExp{"type": "view"}).
		Row(&total)

	return total, err
}

func TestSingleVsMultipleValuesNormalization(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	beforeTotalViews, err := getTotalViews(app)
	if err != nil {
		t.Fatal(err)
	}

	// mock field changes
	collection.Fields.GetByName("select_one").(*core.SelectField).MaxSelect = 2
	collection.Fields.GetByName("select_many").(*core.SelectField).MaxSelect = 1
	collection.Fields.GetByName("file_one").(*core.FileField).MaxSelect = 2
	collection.Fields.GetByName("file_many").(*core.FileField).MaxSelect = 1
	collection.Fields.GetByName("rel_one").(*core.RelationField).MaxSelect = 2
	collection.Fields.GetByName("rel_many").(*core.RelationField).MaxSelect = 1

	// new multivaluer field to check whether the array normalization
	// will be applied for already inserted data
	collection.Fields.Add(&core.SelectField{
		Name:      "new_multiple",
		Values:    []string{"a", "b", "c"},
		MaxSelect: 3,
	})

	if err := app.Save(collection); err != nil {
		t.Fatal(err)
	}

	// ensure that the views were reinserted
	afterTotalViews, err := getTotalViews(app)
	if err != nil {
		t.Fatal(err)
	}
	if afterTotalViews != beforeTotalViews {
		t.Fatalf("Expected total views %d, got %d", beforeTotalViews, afterTotalViews)
	}

	// check whether the columns DEFAULT definition was updated
	// ---------------------------------------------------------------
	tableInfo, err := app.TableInfo(collection.Name)
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
			var row *core.TableInfoRow
			for _, r := range tableInfo {
				if r.Name == col {
					row = r
					break
				}
			}
			if row == nil {
				t.Fatalf("Missing info for column %q", col)
			}

			if v := row.DefaultValue.String; v != dflt {
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

			err := app.DB().Select(
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
