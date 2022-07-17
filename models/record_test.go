package models_test

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestNewRecord(t *testing.T) {
	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "test",
				Type: schema.FieldTypeText,
			},
		),
	}

	m := models.NewRecord(collection)

	if m.Collection().Id != collection.Id {
		t.Fatalf("Expected collection with id %v, got %v", collection.Id, m.Collection().Id)
	}

	if len(m.Data()) != 0 {
		t.Fatalf("Expected empty data, got %v", m.Data())
	}
}

func TestNewRecordFromNullStringMap(t *testing.T) {
	collection := &models.Collection{
		Name: "test",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field1",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "field2",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "field3",
				Type: schema.FieldTypeBool,
			},
			&schema.SchemaField{
				Name: "field4",
				Type: schema.FieldTypeNumber,
			},
			&schema.SchemaField{
				Name: "field5",
				Type: schema.FieldTypeSelect,
				Options: &schema.SelectOptions{
					Values:    []string{"test1", "test2"},
					MaxSelect: 1,
				},
			},
			&schema.SchemaField{
				Name: "field6",
				Type: schema.FieldTypeFile,
				Options: &schema.FileOptions{
					MaxSelect: 2,
					MaxSize:   1,
				},
			},
		),
	}

	data := dbx.NullStringMap{
		"id": sql.NullString{
			String: "c23eb053-d07e-4fbe-86b3-b8ac31982e9a",
			Valid:  true,
		},
		"created": sql.NullString{
			String: "2022-01-01 10:00:00.123",
			Valid:  true,
		},
		"updated": sql.NullString{
			String: "2022-01-01 10:00:00.456",
			Valid:  true,
		},
		"field1": sql.NullString{
			String: "test",
			Valid:  true,
		},
		"field2": sql.NullString{
			String: "test",
			Valid:  false, // test invalid db serialization
		},
		"field3": sql.NullString{
			String: "true",
			Valid:  true,
		},
		"field4": sql.NullString{
			String: "123.123",
			Valid:  true,
		},
		"field5": sql.NullString{
			String: `["test1","test2"]`, // will select only the first elem
			Valid:  true,
		},
		"field6": sql.NullString{
			String: "test", // will be converted to slice
			Valid:  true,
		},
	}

	m := models.NewRecordFromNullStringMap(collection, data)
	encoded, err := m.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"@collectionId":"","@collectionName":"test","created":"2022-01-01 10:00:00.123","field1":"test","field2":"","field3":true,"field4":123.123,"field5":"test1","field6":["test"],"id":"c23eb053-d07e-4fbe-86b3-b8ac31982e9a","updated":"2022-01-01 10:00:00.456"}`

	if string(encoded) != expected {
		t.Fatalf("Expected %v, got \n%v", expected, string(encoded))
	}
}

func TestNewRecordsFromNullStringMaps(t *testing.T) {
	collection := &models.Collection{
		Name: "test",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field1",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "field2",
				Type: schema.FieldTypeNumber,
			},
		),
	}

	data := []dbx.NullStringMap{
		{
			"id": sql.NullString{
				String: "11111111-d07e-4fbe-86b3-b8ac31982e9a",
				Valid:  true,
			},
			"created": sql.NullString{
				String: "2022-01-01 10:00:00.123",
				Valid:  true,
			},
			"updated": sql.NullString{
				String: "2022-01-01 10:00:00.456",
				Valid:  true,
			},
			"field1": sql.NullString{
				String: "test1",
				Valid:  true,
			},
			"field2": sql.NullString{
				String: "123",
				Valid:  false, // test invalid db serialization
			},
		},
		{
			"id": sql.NullString{
				String: "22222222-d07e-4fbe-86b3-b8ac31982e9a",
				Valid:  true,
			},
			"field1": sql.NullString{
				String: "test2",
				Valid:  true,
			},
			"field2": sql.NullString{
				String: "123",
				Valid:  true,
			},
		},
	}

	result := models.NewRecordsFromNullStringMaps(collection, data)
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[{"@collectionId":"","@collectionName":"test","created":"2022-01-01 10:00:00.123","field1":"test1","field2":0,"id":"11111111-d07e-4fbe-86b3-b8ac31982e9a","updated":"2022-01-01 10:00:00.456"},{"@collectionId":"","@collectionName":"test","created":"","field1":"test2","field2":123,"id":"22222222-d07e-4fbe-86b3-b8ac31982e9a","updated":""}]`

	if string(encoded) != expected {
		t.Fatalf("Expected \n%v, got \n%v", expected, string(encoded))
	}
}

func TestRecordCollection(t *testing.T) {
	collection := &models.Collection{}
	collection.RefreshId()

	m := models.NewRecord(collection)

	if m.Collection().Id != collection.Id {
		t.Fatalf("Expected collection with id %v, got %v", collection.Id, m.Collection().Id)
	}
}

func TestRecordTableName(t *testing.T) {
	collection := &models.Collection{}
	collection.Name = "test"
	collection.RefreshId()

	m := models.NewRecord(collection)

	if m.TableName() != collection.Name {
		t.Fatalf("Expected table %q, got %q", collection.Name, m.TableName())
	}
}

func TestRecordExpand(t *testing.T) {
	collection := &models.Collection{}
	m := models.NewRecord(collection)

	data := map[string]any{"test": 123}

	m.SetExpand(data)

	// change the original data to check if it was shallow copied
	data["test"] = 456

	expand := m.GetExpand()
	if v, ok := expand["test"]; !ok || v != 123 {
		t.Fatalf("Expected expand.test to be %v, got %v", 123, v)
	}
}

func TestRecordLoadAndData(t *testing.T) {
	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field",
				Type: schema.FieldTypeText,
			},
		),
	}
	m := models.NewRecord(collection)

	data := map[string]any{
		"id":      "11111111-d07e-4fbe-86b3-b8ac31982e9a",
		"created": "2022-01-01 10:00:00.123",
		"updated": "2022-01-01 10:00:00.456",
		"field":   "test",
		"unknown": "test",
	}

	m.Load(data)

	// change some of original data fields to check if they were shallow copied
	data["id"] = "22222222-d07e-4fbe-86b3-b8ac31982e9a"
	data["field"] = "new_test"

	expectedData := `{"field":"test"}`
	encodedData, _ := json.Marshal(m.Data())
	if string(encodedData) != expectedData {
		t.Fatalf("Expected data %v, got \n%v", expectedData, string(encodedData))
	}

	expectedModel := `{"@collectionId":"","@collectionName":"","created":"2022-01-01 10:00:00.123","field":"test","id":"11111111-d07e-4fbe-86b3-b8ac31982e9a","updated":"2022-01-01 10:00:00.456"}`
	encodedModel, _ := json.Marshal(m)
	if string(encodedModel) != expectedModel {
		t.Fatalf("Expected model %v, got \n%v", expectedModel, string(encodedModel))
	}
}

func TestRecordSetDataValue(t *testing.T) {
	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field",
				Type: schema.FieldTypeText,
			},
		),
	}
	m := models.NewRecord(collection)

	m.SetDataValue("unknown", 123)
	m.SetDataValue("field", 123) // test whether PrepareValue will be called and casted to string

	data := m.Data()
	if len(data) != 1 {
		t.Fatalf("Expected only 1 data field to be set, got %v", data)
	}

	if v, ok := data["field"]; !ok || v != "123" {
		t.Fatalf("Expected field to be %v, got %v", "123", v)
	}
}

func TestRecordGetDataValue(t *testing.T) {
	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field1",
				Type: schema.FieldTypeNumber,
			},
			&schema.SchemaField{
				Name: "field2",
				Type: schema.FieldTypeNumber,
			},
		),
	}
	m := models.NewRecord(collection)

	m.SetDataValue("field2", 123)

	// missing
	v0 := m.GetDataValue("missing")
	if v0 != nil {
		t.Fatalf("Unexpected value for key 'missing'")
	}

	// existing - not set
	v1 := m.GetDataValue("field1")
	if v1 != nil {
		t.Fatalf("Unexpected value for key 'field1'")
	}

	// existing - set
	v2 := m.GetDataValue("field2")
	if v2 != 123.0 {
		t.Fatalf("Expected 123.0, got %v", v2)
	}
}

func TestRecordGetBoolDataValue(t *testing.T) {
	scenarios := []struct {
		value    any
		expected bool
	}{
		{nil, false},
		{"", false},
		{0, false},
		{1, true},
		{[]string{"true"}, false},
		{time.Now(), false},
		{"test", false},
		{"false", false},
		{"true", true},
		{false, false},
		{true, true},
	}

	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "test"},
		),
	}

	for i, s := range scenarios {
		m := models.NewRecord(collection)
		m.SetDataValue("test", s.value)

		result := m.GetBoolDataValue("test")
		if result != s.expected {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestRecordGetStringDataValue(t *testing.T) {
	scenarios := []struct {
		value    any
		expected string
	}{
		{nil, ""},
		{"", ""},
		{0, "0"},
		{1.4, "1.4"},
		{[]string{"true"}, ""},
		{map[string]int{"test": 1}, ""},
		{[]byte("abc"), "abc"},
		{"test", "test"},
		{false, "false"},
		{true, "true"},
	}

	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "test"},
		),
	}

	for i, s := range scenarios {
		m := models.NewRecord(collection)
		m.SetDataValue("test", s.value)

		result := m.GetStringDataValue("test")
		if result != s.expected {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestRecordGetIntDataValue(t *testing.T) {
	scenarios := []struct {
		value    any
		expected int
	}{
		{nil, 0},
		{"", 0},
		{[]string{"true"}, 0},
		{map[string]int{"test": 1}, 0},
		{time.Now(), 0},
		{"test", 0},
		{123, 123},
		{2.4, 2},
		{"123", 123},
		{"123.5", 0},
		{false, 0},
		{true, 1},
	}

	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "test"},
		),
	}

	for i, s := range scenarios {
		m := models.NewRecord(collection)
		m.SetDataValue("test", s.value)

		result := m.GetIntDataValue("test")
		if result != s.expected {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestRecordGetFloatDataValue(t *testing.T) {
	scenarios := []struct {
		value    any
		expected float64
	}{
		{nil, 0},
		{"", 0},
		{[]string{"true"}, 0},
		{map[string]int{"test": 1}, 0},
		{time.Now(), 0},
		{"test", 0},
		{123, 123},
		{2.4, 2.4},
		{"123", 123},
		{"123.5", 123.5},
		{false, 0},
		{true, 1},
	}

	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "test"},
		),
	}

	for i, s := range scenarios {
		m := models.NewRecord(collection)
		m.SetDataValue("test", s.value)

		result := m.GetFloatDataValue("test")
		if result != s.expected {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestRecordGetTimeDataValue(t *testing.T) {
	nowTime := time.Now()
	testTime, _ := time.Parse(types.DefaultDateLayout, "2022-01-01 08:00:40.000")

	scenarios := []struct {
		value    any
		expected time.Time
	}{
		{nil, time.Time{}},
		{"", time.Time{}},
		{false, time.Time{}},
		{true, time.Time{}},
		{"test", time.Time{}},
		{[]string{"true"}, time.Time{}},
		{map[string]int{"test": 1}, time.Time{}},
		{1641024040, testTime},
		{"2022-01-01 08:00:40.000", testTime},
		{nowTime, nowTime},
	}

	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "test"},
		),
	}

	for i, s := range scenarios {
		m := models.NewRecord(collection)
		m.SetDataValue("test", s.value)

		result := m.GetTimeDataValue("test")
		if !result.Equal(s.expected) {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestRecordGetDateTimeDataValue(t *testing.T) {
	nowTime := time.Now()
	testTime, _ := time.Parse(types.DefaultDateLayout, "2022-01-01 08:00:40.000")

	scenarios := []struct {
		value    any
		expected time.Time
	}{
		{nil, time.Time{}},
		{"", time.Time{}},
		{false, time.Time{}},
		{true, time.Time{}},
		{"test", time.Time{}},
		{[]string{"true"}, time.Time{}},
		{map[string]int{"test": 1}, time.Time{}},
		{1641024040, testTime},
		{"2022-01-01 08:00:40.000", testTime},
		{nowTime, nowTime},
	}

	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "test"},
		),
	}

	for i, s := range scenarios {
		m := models.NewRecord(collection)
		m.SetDataValue("test", s.value)

		result := m.GetDateTimeDataValue("test")
		if !result.Time().Equal(s.expected) {
			t.Errorf("(%d) Expected %v, got %v", i, s.expected, result)
		}
	}
}

func TestRecordGetStringSliceDataValue(t *testing.T) {
	nowTime := time.Now()

	scenarios := []struct {
		value    any
		expected []string
	}{
		{nil, []string{}},
		{"", []string{}},
		{false, []string{"false"}},
		{true, []string{"true"}},
		{nowTime, []string{}},
		{123, []string{"123"}},
		{"test", []string{"test"}},
		{map[string]int{"test": 1}, []string{}},
		{`["test1", "test2"]`, []string{"test1", "test2"}},
		{[]int{123, 123, 456}, []string{"123", "456"}},
		{[]string{"test", "test", "123"}, []string{"test", "123"}},
	}

	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{Name: "test"},
		),
	}

	for i, s := range scenarios {
		m := models.NewRecord(collection)
		m.SetDataValue("test", s.value)

		result := m.GetStringSliceDataValue("test")

		if len(result) != len(s.expected) {
			t.Errorf("(%d) Expected %d elements, got %d: %v", i, len(s.expected), len(result), result)
			continue
		}

		for _, v := range result {
			if !list.ExistInSlice(v, s.expected) {
				t.Errorf("(%d) Cannot find %v in %v", i, v, s.expected)
			}
		}
	}
}

func TestRecordBaseFilesPath(t *testing.T) {
	collection := &models.Collection{}
	collection.RefreshId()
	collection.Name = "test"

	m := models.NewRecord(collection)
	m.RefreshId()

	expected := collection.BaseFilesPath() + "/" + m.Id
	result := m.BaseFilesPath()

	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestRecordFindFileFieldByFile(t *testing.T) {
	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field1",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "field2",
				Type: schema.FieldTypeFile,
				Options: &schema.FileOptions{
					MaxSelect: 1,
					MaxSize:   1,
				},
			},
			&schema.SchemaField{
				Name: "field3",
				Type: schema.FieldTypeFile,
				Options: &schema.FileOptions{
					MaxSelect: 2,
					MaxSize:   1,
				},
			},
		),
	}

	m := models.NewRecord(collection)
	m.SetDataValue("field1", "test")
	m.SetDataValue("field2", "test.png")
	m.SetDataValue("field3", []string{"test1.png", "test2.png"})

	scenarios := []struct {
		filename    string
		expectField string
	}{
		{"", ""},
		{"test", ""},
		{"test2", ""},
		{"test.png", "field2"},
		{"test2.png", "field3"},
	}

	for i, s := range scenarios {
		result := m.FindFileFieldByFile(s.filename)

		var fieldName string
		if result != nil {
			fieldName = result.Name
		}

		if s.expectField != fieldName {
			t.Errorf("(%d) Expected field %v, got %v", i, s.expectField, result)
			continue
		}
	}
}

func TestRecordColumnValueMap(t *testing.T) {
	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field1",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "field2",
				Type: schema.FieldTypeFile,
				Options: &schema.FileOptions{
					MaxSelect: 1,
					MaxSize:   1,
				},
			},
			&schema.SchemaField{
				Name: "#field3",
				Type: schema.FieldTypeSelect,
				Options: &schema.SelectOptions{
					MaxSelect: 2,
					Values:    []string{"test1", "test2", "test3"},
				},
			},
			&schema.SchemaField{
				Name: "field4",
				Type: schema.FieldTypeRelation,
				Options: &schema.RelationOptions{
					MaxSelect: 2,
				},
			},
		),
	}

	id1 := "11111111-1e32-4c94-ae06-90c25fcf6791"
	id2 := "22222222-1e32-4c94-ae06-90c25fcf6791"
	created, _ := types.ParseDateTime("2022-01-01 10:00:30.123")

	m := models.NewRecord(collection)
	m.Id = id1
	m.Created = created
	m.SetDataValue("field1", "test")
	m.SetDataValue("field2", "test.png")
	m.SetDataValue("#field3", []string{"test1", "test2"})
	m.SetDataValue("field4", []string{id1, id2, id1})

	result := m.ColumnValueMap()

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"#field3":["test1","test2"],"created":"2022-01-01 10:00:30.123","field1":"test","field2":"test.png","field4":["11111111-1e32-4c94-ae06-90c25fcf6791","22222222-1e32-4c94-ae06-90c25fcf6791"],"id":"11111111-1e32-4c94-ae06-90c25fcf6791","updated":""}`

	if string(encoded) != expected {
		t.Fatalf("Expected %v, got \n%v", expected, string(encoded))
	}
}

func TestRecordPublicExport(t *testing.T) {
	collection := &models.Collection{
		Name: "test",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field1",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "field2",
				Type: schema.FieldTypeFile,
				Options: &schema.FileOptions{
					MaxSelect: 1,
					MaxSize:   1,
				},
			},
			&schema.SchemaField{
				Name: "#field3",
				Type: schema.FieldTypeSelect,
				Options: &schema.SelectOptions{
					MaxSelect: 2,
					Values:    []string{"test1", "test2", "test3"},
				},
			},
		),
	}

	created, _ := types.ParseDateTime("2022-01-01 10:00:30.123")

	m := models.NewRecord(collection)
	m.Id = "210a896c-1e32-4c94-ae06-90c25fcf6791"
	m.Created = created
	m.SetDataValue("field1", "test")
	m.SetDataValue("field2", "test.png")
	m.SetDataValue("#field3", []string{"test1", "test2"})
	m.SetExpand(map[string]any{"test": 123})

	result := m.PublicExport()

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"@collectionId":"","@collectionName":"test","@expand":{"test":123},"created":"2022-01-01 10:00:30.123","field1":"test","field2":"test.png","id":"210a896c-1e32-4c94-ae06-90c25fcf6791","updated":""}`

	if string(encoded) != expected {
		t.Fatalf("Expected %v, got \n%v", expected, string(encoded))
	}
}

func TestRecordMarshalJSON(t *testing.T) {
	collection := &models.Collection{
		Name: "test",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field1",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "field2",
				Type: schema.FieldTypeFile,
				Options: &schema.FileOptions{
					MaxSelect: 1,
					MaxSize:   1,
				},
			},
			&schema.SchemaField{
				Name: "#field3",
				Type: schema.FieldTypeSelect,
				Options: &schema.SelectOptions{
					MaxSelect: 2,
					Values:    []string{"test1", "test2", "test3"},
				},
			},
		),
	}

	created, _ := types.ParseDateTime("2022-01-01 10:00:30.123")

	m := models.NewRecord(collection)
	m.Id = "210a896c-1e32-4c94-ae06-90c25fcf6791"
	m.Created = created
	m.SetDataValue("field1", "test")
	m.SetDataValue("field2", "test.png")
	m.SetDataValue("#field3", []string{"test1", "test2"})
	m.SetExpand(map[string]any{"test": 123})

	encoded, err := m.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"@collectionId":"","@collectionName":"test","@expand":{"test":123},"created":"2022-01-01 10:00:30.123","field1":"test","field2":"test.png","id":"210a896c-1e32-4c94-ae06-90c25fcf6791","updated":""}`

	if string(encoded) != expected {
		t.Fatalf("Expected %v, got \n%v", expected, string(encoded))
	}
}

func TestRecordUnmarshalJSON(t *testing.T) {
	collection := &models.Collection{
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "field",
				Type: schema.FieldTypeText,
			},
		),
	}
	m := models.NewRecord(collection)

	m.UnmarshalJSON([]byte(`{
		"id":      "11111111-d07e-4fbe-86b3-b8ac31982e9a",
		"created": "2022-01-01 10:00:00.123",
		"updated": "2022-01-01 10:00:00.456",
		"field":   "test",
		"unknown": "test"
	}`))

	expected := `{"@collectionId":"","@collectionName":"","created":"2022-01-01 10:00:00.123","field":"test","id":"11111111-d07e-4fbe-86b3-b8ac31982e9a","updated":"2022-01-01 10:00:00.456"}`
	encoded, _ := json.Marshal(m)
	if string(encoded) != expected {
		t.Fatalf("Expected model %v, got \n%v", expected, string(encoded))
	}
}
