package daos_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestRecordQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo")

	expected := fmt.Sprintf("SELECT `%s`.* FROM `%s`", collection.Name, collection.Name)

	sql := app.Dao().RecordQuery(collection).Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindRecordById(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo")

	scenarios := []struct {
		id          string
		filter      func(q *dbx.SelectQuery) error
		expectError bool
	}{
		{"00000000-bafd-48f7-b8b7-090638afe209", nil, true},
		{"b5c2ffc2-bafd-48f7-b8b7-090638afe209", nil, false},
		{"b5c2ffc2-bafd-48f7-b8b7-090638afe209", func(q *dbx.SelectQuery) error {
			q.AndWhere(dbx.HashExp{"title": "missing"})
			return nil
		}, true},
		{"b5c2ffc2-bafd-48f7-b8b7-090638afe209", func(q *dbx.SelectQuery) error {
			return errors.New("test error")
		}, true},
		{"b5c2ffc2-bafd-48f7-b8b7-090638afe209", func(q *dbx.SelectQuery) error {
			q.AndWhere(dbx.HashExp{"title": "lorem"})
			return nil
		}, false},
	}

	for i, scenario := range scenarios {
		record, err := app.Dao().FindRecordById(collection, scenario.id, scenario.filter)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if record != nil && record.Id != scenario.id {
			t.Errorf("(%d) Expected record with id %s, got %s", i, scenario.id, record.Id)
		}
	}
}

func TestFindRecordsByIds(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo")

	scenarios := []struct {
		ids         []string
		filter      func(q *dbx.SelectQuery) error
		expectTotal int
		expectError bool
	}{
		{[]string{}, nil, 0, false},
		{[]string{"00000000-bafd-48f7-b8b7-090638afe209"}, nil, 0, false},
		{[]string{"b5c2ffc2-bafd-48f7-b8b7-090638afe209"}, nil, 1, false},
		{
			[]string{"b5c2ffc2-bafd-48f7-b8b7-090638afe209", "848a1dea-5ddd-42d6-a00d-030547bffcfe"},
			nil,
			2,
			false,
		},
		{
			[]string{"b5c2ffc2-bafd-48f7-b8b7-090638afe209", "848a1dea-5ddd-42d6-a00d-030547bffcfe"},
			func(q *dbx.SelectQuery) error {
				return errors.New("test error")
			},
			0,
			true,
		},
		{
			[]string{"b5c2ffc2-bafd-48f7-b8b7-090638afe209", "848a1dea-5ddd-42d6-a00d-030547bffcfe"},
			func(q *dbx.SelectQuery) error {
				q.AndWhere(dbx.Like("title", "test").Match(true, true))
				return nil
			},
			1,
			false,
		},
	}

	for i, scenario := range scenarios {
		records, err := app.Dao().FindRecordsByIds(collection, scenario.ids, scenario.filter)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if len(records) != scenario.expectTotal {
			t.Errorf("(%d) Expected %d records, got %d", i, scenario.expectTotal, len(records))
			continue
		}

		for _, r := range records {
			if !list.ExistInSlice(r.Id, scenario.ids) {
				t.Errorf("(%d) Couldn't find id %s in %v", i, r.Id, scenario.ids)
			}
		}
	}
}

func TestFindRecordsByExpr(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo")

	scenarios := []struct {
		expression  dbx.Expression
		expectIds   []string
		expectError bool
	}{
		{
			nil,
			[]string{},
			true,
		},
		{
			dbx.HashExp{"id": 123},
			[]string{},
			false,
		},
		{
			dbx.Like("title", "test").Match(true, true),
			[]string{
				"848a1dea-5ddd-42d6-a00d-030547bffcfe",
				"577bd676-aacb-4072-b7da-99d00ee210a4",
			},
			false,
		},
	}

	for i, scenario := range scenarios {
		records, err := app.Dao().FindRecordsByExpr(collection, scenario.expression)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if len(records) != len(scenario.expectIds) {
			t.Errorf("(%d) Expected %d records, got %d", i, len(scenario.expectIds), len(records))
			continue
		}

		for _, r := range records {
			if !list.ExistInSlice(r.Id, scenario.expectIds) {
				t.Errorf("(%d) Couldn't find id %s in %v", i, r.Id, scenario.expectIds)
			}
		}
	}
}

func TestFindFirstRecordByData(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo")

	scenarios := []struct {
		key         string
		value       any
		expectId    string
		expectError bool
	}{
		{
			"",
			"848a1dea-5ddd-42d6-a00d-030547bffcfe",
			"",
			true,
		},
		{
			"id",
			"invalid",
			"",
			true,
		},
		{
			"id",
			"848a1dea-5ddd-42d6-a00d-030547bffcfe",
			"848a1dea-5ddd-42d6-a00d-030547bffcfe",
			false,
		},
		{
			"title",
			"lorem",
			"b5c2ffc2-bafd-48f7-b8b7-090638afe209",
			false,
		},
	}

	for i, scenario := range scenarios {
		record, err := app.Dao().FindFirstRecordByData(collection, scenario.key, scenario.value)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if !scenario.expectError && record.Id != scenario.expectId {
			t.Errorf("(%d) Expected record with id %s, got %v", i, scenario.expectId, record.Id)
		}
	}
}

func TestIsRecordValueUnique(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo4")

	testManyRelsId1 := "df55c8ff-45ef-4c82-8aed-6e2183fe1125"
	testManyRelsId2 := "b84cd893-7119-43c9-8505-3c4e22da28a9"

	scenarios := []struct {
		key       string
		value     any
		excludeId string
		expected  bool
	}{
		{"", "", "", false},
		{"missing", "unique", "", false},
		{"title", "unique", "", true},
		{"title", "demo1", "", false},
		{"title", "demo1", "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2", true},
		{"manyrels", []string{testManyRelsId2}, "", false},
		{"manyrels", []any{testManyRelsId2}, "", false},
		// with exclude
		{"manyrels", []string{testManyRelsId1, testManyRelsId2}, "b8ba58f9-e2d7-42a0-b0e7-a11efd98236b", true},
		// reverse order
		{"manyrels", []string{testManyRelsId2, testManyRelsId1}, "", true},
	}

	for i, scenario := range scenarios {
		result := app.Dao().IsRecordValueUnique(collection, scenario.key, scenario.value, scenario.excludeId)

		if result != scenario.expected {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expected, result)
		}
	}
}

func TestFindUserRelatedRecords(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	u0 := &models.User{}
	u1, _ := app.Dao().FindUserByEmail("test3@example.com")
	u2, _ := app.Dao().FindUserByEmail("test2@example.com")

	scenarios := []struct {
		user        *models.User
		expectedIds []string
	}{
		{u0, []string{}},
		{u1, []string{
			"94568ca2-0bee-49d7-b749-06cb97956fd9", // demo2
			"fc69274d-ca5c-416a-b9ef-561b101cfbb1", // profile
		}},
		{u2, []string{
			"b2d5e39d-f569-4cc1-b593-3f074ad026bf", // profile
		}},
	}

	for i, scenario := range scenarios {
		records, err := app.Dao().FindUserRelatedRecords(scenario.user)
		if err != nil {
			t.Fatal(err)
		}

		if len(records) != len(scenario.expectedIds) {
			t.Errorf("(%d) Expected %d records, got %d (%v)", i, len(scenario.expectedIds), len(records), records)
			continue
		}

		for _, r := range records {
			if !list.ExistInSlice(r.Id, scenario.expectedIds) {
				t.Errorf("(%d) Couldn't find %s in %v", i, r.Id, scenario.expectedIds)
			}
		}
	}
}

func TestSaveRecord(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo")

	// create
	// ---
	r1 := models.NewRecord(collection)
	r1.SetDataValue("title", "test_new")
	err1 := app.Dao().SaveRecord(r1)
	if err1 != nil {
		t.Fatal(err1)
	}
	newR1, _ := app.Dao().FindFirstRecordByData(collection, "title", "test_new")
	if newR1 == nil || newR1.Id != r1.Id || newR1.GetStringDataValue("title") != r1.GetStringDataValue("title") {
		t.Errorf("Expected to find record %v, got %v", r1, newR1)
	}

	// update
	// ---
	r2, _ := app.Dao().FindFirstRecordByData(collection, "id", "b5c2ffc2-bafd-48f7-b8b7-090638afe209")
	r2.SetDataValue("title", "test_update")
	err2 := app.Dao().SaveRecord(r2)
	if err2 != nil {
		t.Fatal(err2)
	}
	newR2, _ := app.Dao().FindFirstRecordByData(collection, "title", "test_update")
	if newR2 == nil || newR2.Id != r2.Id || newR2.GetStringDataValue("title") != r2.GetStringDataValue("title") {
		t.Errorf("Expected to find record %v, got %v", r2, newR2)
	}
}

func TestDeleteRecord(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demo, _ := app.Dao().FindCollectionByNameOrId("demo")
	demo2, _ := app.Dao().FindCollectionByNameOrId("demo2")

	// delete unsaved record
	// ---
	rec1 := models.NewRecord(demo)
	err1 := app.Dao().DeleteRecord(rec1)
	if err1 == nil {
		t.Fatal("(rec1) Didn't expect to succeed deleting new record")
	}

	// delete existing record while being part of a non-cascade required relation
	// ---
	rec2, _ := app.Dao().FindFirstRecordByData(demo, "id", "848a1dea-5ddd-42d6-a00d-030547bffcfe")
	err2 := app.Dao().DeleteRecord(rec2)
	if err2 == nil {
		t.Fatalf("(rec2) Expected error, got nil")
	}

	// delete existing record
	// ---
	rec3, _ := app.Dao().FindFirstRecordByData(demo, "id", "577bd676-aacb-4072-b7da-99d00ee210a4")
	err3 := app.Dao().DeleteRecord(rec3)
	if err3 != nil {
		t.Fatalf("(rec3) Expected nil, got error %v", err3)
	}

	// check if it was really deleted
	rec3, _ = app.Dao().FindRecordById(demo, rec3.Id, nil)
	if rec3 != nil {
		t.Fatalf("(rec3) Expected record to be deleted, got %v", rec3)
	}

	// check if the operation cascaded
	rel, _ := app.Dao().FindFirstRecordByData(demo2, "id", "63c2ab80-84ab-4057-a592-4604a731f78f")
	if rel != nil {
		t.Fatalf("(rec3) Expected the delete to cascade, found relation %v", rel)
	}
}

func TestSyncRecordTableSchema(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	oldCollection, err := app.Dao().FindCollectionByNameOrId("demo")
	if err != nil {
		t.Fatal(err)
	}
	updatedCollection, err := app.Dao().FindCollectionByNameOrId("demo")
	if err != nil {
		t.Fatal(err)
	}
	updatedCollection.Name = "demo_renamed"
	updatedCollection.Schema.RemoveField(updatedCollection.Schema.GetFieldByName("file").Id)
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
		// no changes
		{
			oldCollection,
			oldCollection,
			"demo",
			[]string{"id", "created", "updated", "title", "file"},
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
