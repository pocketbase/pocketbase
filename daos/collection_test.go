package daos_test

import (
	"encoding/json"
	"errors"
	"strings"
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

func TestFindCollectionsByType(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		collectionType string
		expectError    bool
		expectTotal    int
	}{
		{"", false, 0},
		{"unknown", false, 0},
		{models.CollectionTypeAuth, false, 3},
		{models.CollectionTypeBase, false, 5},
	}

	for i, scenario := range scenarios {
		collections, err := app.Dao().FindCollectionsByType(scenario.collectionType)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("[%d] Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if len(collections) != scenario.expectTotal {
			t.Errorf("[%d] Expected %d collections, got %d", i, scenario.expectTotal, len(collections))
		}

		for _, c := range collections {
			if c.Type != scenario.collectionType {
				t.Errorf("[%d] Expected collection with type %s, got %s: \n%v", i, scenario.collectionType, c.Type, c)
			}
		}
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
		{"wsmn24bux7wo113", false},
		{"demo1", false},
		{"DEMO1", false}, // case insensitive check
	}

	for i, scenario := range scenarios {
		model, err := app.Dao().FindCollectionByNameOrId(scenario.nameOrId)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("[%d] Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if model != nil && model.Id != scenario.nameOrId && !strings.EqualFold(model.Name, scenario.nameOrId) {
			t.Errorf("[%d] Expected model with identifier %s, got %v", i, scenario.nameOrId, model)
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
		{"demo1", "", false},
		{"Demo1", "", false},
		{"new", "", true},
		{"demo1", "wsmn24bux7wo113", true},
	}

	for i, scenario := range scenarios {
		result := app.Dao().IsCollectionNameUnique(scenario.name, scenario.excludeId)
		if result != scenario.expected {
			t.Errorf("[%d] Expected %v, got %v", i, scenario.expected, result)
		}
	}
}

func TestFindCollectionReferences(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo3")
	if err != nil {
		t.Fatal(err)
	}

	result, err := app.Dao().FindCollectionReferences(
		collection,
		collection.Id,
		// test whether "nonempty" exclude ids condition will be skipped
		"",
		"",
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 collection, got %d: %v", len(result), result)
	}

	expectedFields := []string{
		"rel_one_no_cascade",
		"rel_one_no_cascade_required",
		"rel_one_cascade",
		"rel_many_no_cascade",
		"rel_many_no_cascade_required",
		"rel_many_cascade",
	}

	for col, fields := range result {
		if col.Name != "demo4" {
			t.Fatalf("Expected collection demo4, got %s", col.Name)
		}
		if len(fields) != len(expectedFields) {
			t.Fatalf("Expected fields %v, got %v", expectedFields, fields)
		}
		for i, f := range fields {
			if !list.ExistInSlice(f.Name, expectedFields) {
				t.Fatalf("[%d] Didn't expect field %v", i, f)
			}
		}
	}
}

func TestDeleteCollection(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	colEmpty := &models.Collection{}

	colAuth, err := app.Dao().FindCollectionByNameOrId("clients")
	if err != nil {
		t.Fatal(err)
	}

	colReferenced, err := app.Dao().FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}

	colSystem, err := app.Dao().FindCollectionByNameOrId("demo3")
	if err != nil {
		t.Fatal(err)
	}
	colSystem.System = true
	if err := app.Dao().Save(colSystem); err != nil {
		t.Fatal(err)
	}

	colView, err := app.Dao().FindCollectionByNameOrId("view1")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		model       *models.Collection
		expectError bool
	}{
		{colEmpty, true},
		{colAuth, false},
		{colReferenced, true},
		{colSystem, true},
		{colView, false},
	}

	for i, s := range scenarios {
		err := app.Dao().DeleteCollection(s.model)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%d] Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		if app.Dao().HasTable(s.model.Name) {
			t.Errorf("[%d] Expected table/view %s to be deleted", i, s.model.Name)
		}
	}
}

func TestSaveCollectionCreate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := &models.Collection{
		Name: "new_test",
		Type: models.CollectionTypeBase,
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
			t.Fatalf("[%d] Didn't expect record column %s", i, c)
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
	expectedColumns := []string{"id", "created", "updated", "title_update", "test", "files"}
	columns, err := app.Dao().GetTableColumns(collection.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(columns) != len(expectedColumns) {
		t.Fatalf("Expected columns %v, got %v", expectedColumns, columns)
	}
	for i, c := range columns {
		if !list.ExistInSlice(c, expectedColumns) {
			t.Fatalf("[%d] Didn't expect record column %s", i, c)
		}
	}
}

func TestImportCollections(t *testing.T) {
	totalCollections := 10

	scenarios := []struct {
		name                   string
		jsonData               string
		deleteMissing          bool
		beforeRecordsSync      func(txDao *daos.Dao, mappedImported, mappedExisting map[string]*models.Collection) error
		expectError            bool
		expectCollectionsCount int
		beforeTestFunc         func(testApp *tests.TestApp, resultCollections []*models.Collection)
		afterTestFunc          func(testApp *tests.TestApp, resultCollections []*models.Collection)
	}{
		{
			name:                   "empty collections",
			jsonData:               `[]`,
			expectError:            true,
			expectCollectionsCount: totalCollections,
		},
		{
			name: "minimal collection import",
			jsonData: `[
				{"name": "import_test1", "schema": [{"name":"test", "type": "text"}]},
				{"name": "import_test2", "type": "auth"}
			]`,
			deleteMissing:          false,
			expectError:            false,
			expectCollectionsCount: totalCollections + 2,
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
			expectCollectionsCount: totalCollections,
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
			expectCollectionsCount: totalCollections + 1,
		},
		{
			name: "new + update + delete system collection",
			jsonData: `[
				{
					"id":"wsmn24bux7wo113",
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
			expectCollectionsCount: totalCollections,
		},
		{
			name: "new + update + delete non-system collection",
			jsonData: `[
				{
					"id": "kpv709sk2lqbqk8",
					"system": true,
					"name": "nologin",
					"type": "auth",
					"options": {
						"allowEmailAuth": false,
						"allowOAuth2Auth": false,
						"allowUsernameAuth": false,
						"exceptEmailDomains": [],
						"manageRule": "@request.auth.collectionName = 'users'",
						"minPasswordLength": 8,
						"onlyEmailDomains": [],
						"requireEmail": true
					},
					"listRule": "",
					"viewRule": "",
					"createRule": "",
					"updateRule": "",
					"deleteRule": "",
					"schema": [
						{
							"id": "x8zzktwe",
							"name": "name",
							"type": "text",
							"system": false,
							"required": false,
							"unique": false,
							"options": {
								"min": null,
								"max": null,
								"pattern": ""
							}
						}
					]
				},
				{
					"id":"wsmn24bux7wo113",
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
				},
				{
					"id": "test_new_view",
					"name": "new_view",
					"type": "view",
					"options": {
						"query": "select id from demo2"
					}
				}
			]`,
			deleteMissing:          true,
			expectError:            false,
			expectCollectionsCount: 4,
		},
		{
			name: "test with deleteMissing: false",
			jsonData: `[
				{
					"id":"wsmn24bux7wo113",
					"name":"demo1",
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
			expectCollectionsCount: totalCollections + 1,
			afterTestFunc: func(testApp *tests.TestApp, resultCollections []*models.Collection) {
				expectedCollectionFields := map[string]int{
					"nologin":    1,
					"demo1":      15,
					"demo2":      2,
					"demo3":      2,
					"demo4":      11,
					"demo5":      6,
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
