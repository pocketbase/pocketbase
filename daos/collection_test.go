package daos_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/daos"
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

	expectedFields := []string{"onerel", "manyrels", "cascaderel"}

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

func TestImportCollections(t *testing.T) {
	scenarios := []struct {
		name                   string
		jsonData               string
		deleteMissing          bool
		beforeRecordsSync      func(txDao *daos.Dao, mappedImported, mappedExisting map[string]*models.Collection) error
		expectError            bool
		expectCollectionsCount int
		afterTestFunc          func(testApp *tests.TestApp, resultCollections []*models.Collection)
	}{
		{
			name:                   "empty collections",
			jsonData:               `[]`,
			expectError:            true,
			expectCollectionsCount: 5,
		},
		{
			name: "check db constraints",
			jsonData: `[
				{"name": "import_test", "schema": []}
			]`,
			deleteMissing:          false,
			expectError:            true,
			expectCollectionsCount: 5,
		},
		{
			name: "minimal collection import",
			jsonData: `[
				{"name": "import_test", "schema": [{"name":"test", "type": "text"}]}
			]`,
			deleteMissing:          false,
			expectError:            false,
			expectCollectionsCount: 6,
		},
		{
			name: "minimal collection import + failed beforeRecordsSync",
			jsonData: `[
				{"name": "import_test", "schema": [{"name":"test", "type": "text"}]}
			]`,
			beforeRecordsSync: func(txDao *daos.Dao, mappedImported, mappedExisting map[string]*models.Collection) error {
				return errors.New("test_error")
			},
			deleteMissing:          false,
			expectError:            true,
			expectCollectionsCount: 5,
		},
		{
			name: "minimal collection import + successful beforeRecordsSync",
			jsonData: `[
				{"name": "import_test", "schema": [{"name":"test", "type": "text"}]}
			]`,
			beforeRecordsSync: func(txDao *daos.Dao, mappedImported, mappedExisting map[string]*models.Collection) error {
				return nil
			},
			deleteMissing:          false,
			expectError:            false,
			expectCollectionsCount: 6,
		},
		{
			name: "new + update + delete system collection",
			jsonData: `[
				{
					"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc",
					"name":"demo",
					"schema":[
						{
							"id":"_2hlxbmp",
							"name":"title",
							"type":"text",
							"system":false,
							"required":true,
							"unique":false,
							"options":{
								"min":3,
								"max":null,
								"pattern":""
							}
						}
					]
				},
				{
					"name": "import1",
					"schema": [
						{
							"name":"active",
							"type":"bool"
						}
					]
				}
			]`,
			deleteMissing:          true,
			expectError:            true,
			expectCollectionsCount: 5,
		},
		{
			name: "new + update + delete non-system collection",
			jsonData: `[
				{
					"id":"abe78266-fd4d-4aea-962d-8c0138ac522b",
					"name":"profiles",
					"system":true,
					"listRule":"userId = @request.user.id",
					"viewRule":"created > 'test_change'",
					"createRule":"userId = @request.user.id",
					"updateRule":"userId = @request.user.id",
					"deleteRule":"userId = @request.user.id",
					"schema":[
						{
							"id":"koih1lqx",
							"name":"userId",
							"type":"user",
							"system":true,
							"required":true,
							"unique":true,
							"options":{
								"maxSelect":1,
								"cascadeDelete":true
							}
						},
						{
							"id":"69ycbg3q",
							"name":"rel",
							"type":"relation",
							"system":false,
							"required":false,
							"unique":false,
							"options":{
								"maxSelect":2,
								"collectionId":"abe78266-fd4d-4aea-962d-8c0138ac522b",
								"cascadeDelete":false
							}
						}
					]
				},
				{
					"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc",
					"name":"demo",
					"schema":[
						{
							"id":"_2hlxbmp",
							"name":"title",
							"type":"text",
							"system":false,
							"required":true,
							"unique":false,
							"options":{
								"min":3,
								"max":null,
								"pattern":""
							}
						}
					]
				},
				{
					"id": "test_deleted_collection_name_reuse",
					"name": "demo2",
					"schema": [
						{
							"id":"fz6iql2m",
							"name":"active",
							"type":"bool"
						}
					]
				}
			]`,
			deleteMissing:          true,
			expectError:            false,
			expectCollectionsCount: 3,
		},
		{
			name: "test with deleteMissing: false",
			jsonData: `[
				{
					"id":"abe78266-fd4d-4aea-962d-8c0138ac522b",
					"name":"profiles",
					"system":true,
					"listRule":"userId = @request.user.id",
					"viewRule":"created > 'test_change'",
					"createRule":"userId = @request.user.id",
					"updateRule":"userId = @request.user.id",
					"deleteRule":"userId = @request.user.id",
					"schema":[
						{
							"id":"69ycbg3q",
							"name":"rel",
							"type":"relation",
							"system":false,
							"required":false,
							"unique":false,
							"options":{
								"maxSelect":2,
								"collectionId":"abe78266-fd4d-4aea-962d-8c0138ac522b",
								"cascadeDelete":true
							}
						},
						{
							"id":"abcd_import",
							"name":"new_field",
							"type":"bool"
						}
					]
				},
				{
					"id":"3f2888f8-075d-49fe-9d09-ea7e951000dc",
					"name":"demo",
					"schema":[
						{
							"id":"_2hlxbmp",
							"name":"title",
							"type":"text",
							"system":false,
							"required":true,
							"unique":false,
							"options":{
								"min":3,
								"max":null,
								"pattern":""
							}
						},
						{
							"id":"_2hlxbmp",
							"name":"field_with_duplicate_id",
							"type":"text",
							"system":false,
							"required":true,
							"unique":false,
							"options":{
								"min":3,
								"max":null,
								"pattern":""
							}
						},
						{
							"id":"abcd_import",
							"name":"new_field",
							"type":"text"
						}
					]
				},
				{
					"name": "new_import",
					"schema": [
						{
							"id":"abcd_import",
							"name":"active",
							"type":"bool"
						}
					]
				}
			]`,
			deleteMissing:          false,
			expectError:            false,
			expectCollectionsCount: 6,
			afterTestFunc: func(testApp *tests.TestApp, resultCollections []*models.Collection) {
				expectedCollectionFields := map[string]int{
					"profiles":   6,
					"demo":       3,
					"demo2":      14,
					"demo3":      1,
					"demo4":      6,
					"new_import": 1,
				}
				for name, expectedCount := range expectedCollectionFields {
					collection, err := testApp.Dao().FindCollectionByNameOrId(name)
					if err != nil {
						t.Fatal(err)
					}

					if totalFields := len(collection.Schema.Fields()); totalFields != expectedCount {
						t.Errorf("Expected %d %q fields, got %d", expectedCount, collection.Name, totalFields)
					}
				}
			},
		},
	}

	for _, scenario := range scenarios {
		testApp, _ := tests.NewTestApp()
		defer testApp.Cleanup()

		importedCollections := []*models.Collection{}

		// load data
		loadErr := json.Unmarshal([]byte(scenario.jsonData), &importedCollections)
		if loadErr != nil {
			t.Fatalf("[%s] Failed to load  data: %v", scenario.name, loadErr)
			continue
		}

		err := testApp.Dao().ImportCollections(importedCollections, scenario.deleteMissing, scenario.beforeRecordsSync)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("[%s] Expected hasErr to be %v, got %v (%v)", scenario.name, scenario.expectError, hasErr, err)
		}

		// check collections count
		collections := []*models.Collection{}
		if err := testApp.Dao().CollectionQuery().All(&collections); err != nil {
			t.Fatal(err)
		}
		if len(collections) != scenario.expectCollectionsCount {
			t.Errorf("[%s] Expected %d collections, got %d", scenario.name, scenario.expectCollectionsCount, len(collections))
		}

		if scenario.afterTestFunc != nil {
			scenario.afterTestFunc(testApp, collections)
		}
	}
}
