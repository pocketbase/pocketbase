package models_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func TestViewsTableName(t *testing.T) {
	m := models.View{}
	if m.TableName() != "_views" {
		t.Fatalf("Unexpected table name, got %q", m.TableName())
	}
}

func TestNewRecordsFromViewSchema(t *testing.T) {
	v := &models.View{
		Name: "test",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "created",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "field1",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "field2",
				Type: schema.FieldTypeNumber,
			},
			&schema.SchemaField{
				Name: "id",
				Type: schema.FieldTypeText,
			},
			&schema.SchemaField{
				Name: "updated",
				Type: schema.FieldTypeText,
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

	result := models.NewRecordsFromViewSchema(v, data)
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `[{"created":"2022-01-01 10:00:00.123","field1":"test1","field2":null,"id":"11111111-d07e-4fbe-86b3-b8ac31982e9a","updated":"2022-01-01 10:00:00.456"},{"created":null,"field1":"test2","field2":123,"id":"22222222-d07e-4fbe-86b3-b8ac31982e9a","updated":null}]`

	if string(encoded) != expected {
		t.Fatalf("Expected \n%v, got \n%v", expected, string(encoded))
	}
}
