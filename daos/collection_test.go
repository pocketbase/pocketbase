package daos_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestCollectionQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	expected := "SELECT {{_collections}}.* FROM `_collections`"

	sql := app.Dao().CollectionQuery().Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindCollectionByNameOrId(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		nameOrId    string
		expectError bool
	}{
		{"", true},
		{"missing", true},
		{"00000000-075d-49fe-9d09-ea7e951000dc", true},
		{"3f2888f8-075d-49fe-9d09-ea7e951000dc", false},
		{"demo", false},
	}

	for i, scenario := range scenarios {
		model, err := app.Dao().FindCollectionByNameOrId(scenario.nameOrId)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if model != nil && model.Id != scenario.nameOrId && model.Name != scenario.nameOrId {
			t.Errorf("(%d) Expected model with identifier %s, got %v", i, scenario.nameOrId, model)
		}
	}
}

func TestIsCollectionNameUnique(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name      string
		excludeId string
		expected  bool
	}{
		{"", "", false},
		{"demo", "", false},
		{"new", "", true},
		{"demo", "3f2888f8-075d-49fe-9d09-ea7e951000dc", true},
	}

	for i, scenario := range scenarios {
		result := app.Dao().IsCollectionNameUnique(scenario.name, scenario.excludeId)
		if result != scenario.expected {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expected, result)
		}
	}
}

func TestFindCollectionsWithUserFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	result, err := app.Dao().FindCollectionsWithUserFields()
	if err != nil {
		t.Fatal(err)
	}

	expectedNames := []string{"demo2", models.ProfileCollectionName}

	if len(result) != len(expectedNames) {
		t.Fatalf("Expected collections %v, got %v", expectedNames, result)
	}

	for i, col := range result {
		if !list.ExistInSlice(col.Name, expectedNames) {
			t.Errorf("(%d) Couldn't find %s in %v", i, col.Name, expectedNames)
		}
	}
}

func TestFindCollectionReferences(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo")
	if err != nil {
		t.Fatal(err)
	}

	result, err := app.Dao().FindCollectionReferences(collection, collection.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 collection, got %d: %v", len(result), result)
	}

	expectedFields := []string{"onerel", "manyrels", "rel_cascade"}

	for col, fields := range result {
		if col.Name != "demo2" {
			t.Fatalf("Expected collection demo2, got %s", col.Name)
		}
		if len(fields) != len(expectedFields) {
			t.Fatalf("Expected fields %v, got %v", expectedFields, fields)
		}
		for i, f := range fields {
			if !list.ExistInSlice(f.Name, expectedFields) {
				t.Fatalf("(%d) Didn't expect field %v", i, f)
			}
		}
	}
}

func TestDeleteCollection(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	c0 := &models.Collection{}
	c1, err := app.Dao().FindCollectionByNameOrId("demo")
	if err != nil {
		t.Fatal(err)
	}
	c2, err := app.Dao().FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}
	c3, err := app.Dao().FindCollectionByNameOrId(models.ProfileCollectionName)
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		model       *models.Collection
		expectError bool
	}{
		{c0, true},
		{c1, true}, // is part of a reference
		{c2, false},
		{c3, true}, // system
	}

	for i, scenario := range scenarios {
		err := app.Dao().DeleteCollection(scenario.model)
		hasErr := err != nil

		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v", i, scenario.expectError, hasErr)
		}
	}
}

func TestSaveCollectionCreate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := &models.Collection{
		Name: "new_test",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Type: schema.FieldTypeText,
				Name: "test",
			},
		),
	}

	err := app.Dao().SaveCollection(collection)
	if err != nil {
		t.Fatal(err)
	}

	if collection.Id == "" {
		t.Fatal("Expected collection id to be set")
	}

	// check if the records table was created
	hasTable := app.Dao().HasTable(collection.Name)
	if !hasTable {
		t.Fatalf("Expected records table %s to be created", collection.Name)
	}

	// check if the records table has the schema fields
	columns, err := app.Dao().GetTableColumns(collection.Name)
	if err != nil {
		t.Fatal(err)
	}
	expectedColumns := []string{"id", "created", "updated", "test"}
	if len(columns) != len(expectedColumns) {
		t.Fatalf("Expected columns %v, got %v", expectedColumns, columns)
	}
	for i, c := range columns {
		if !list.ExistInSlice(c, expectedColumns) {
			t.Fatalf("(%d) Didn't expect record column %s", i, c)
		}
	}
}

func TestSaveCollectionUpdate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo3")
	if err != nil {
		t.Fatal(err)
	}

	// rename an existing schema field and add a new one
	oldField := collection.Schema.GetFieldByName("title")
	oldField.Name = "title_update"
	collection.Schema.AddField(&schema.SchemaField{
		Type: schema.FieldTypeText,
		Name: "test",
	})

	saveErr := app.Dao().SaveCollection(collection)
	if saveErr != nil {
		t.Fatal(saveErr)
	}

	// check if the records table has the schema fields
	expectedColumns := []string{"id", "created", "updated", "title_update", "test"}
	columns, err := app.Dao().GetTableColumns(collection.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(columns) != len(expectedColumns) {
		t.Fatalf("Expected columns %v, got %v", expectedColumns, columns)
	}
	for i, c := range columns {
		if !list.ExistInSlice(c, expectedColumns) {
			t.Fatalf("(%d) Didn't expect record column %s", i, c)
		}
	}
}
