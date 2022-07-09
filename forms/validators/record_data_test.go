package validators_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/types"
)

type testDataFieldScenario struct {
	name           string
	data           map[string]any
	files          []*rest.UploadedFile
	expectedErrors []string
}

func TestRecordDataValidatorEmptyAndUnknown(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo")
	record := models.NewRecord(collection)
	validator := validators.NewRecordDataValidator(app.Dao(), record, nil)

	emptyErr := validator.Validate(map[string]any{})
	if emptyErr == nil {
		t.Fatal("Expected error for empty data, got nil")
	}

	unknownErr := validator.Validate(map[string]any{"unknown": 123})
	if unknownErr == nil {
		t.Fatal("Expected error for unknown data, got nil")
	}
}

func TestRecordDataValidatorValidateText(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	min := 3
	max := 10
	pattern := `^\w+$`
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeText,
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeText,
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeText,
			Options: &schema.TextOptions{
				Min:     &min,
				Max:     &max,
				Pattern: pattern,
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", "test")
	dummy.SetDataValue("field2", "test")
	dummy.SetDataValue("field3", "test")
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": "test",
				"field2": "test",
				"field3": "test",
			},
			nil,
			[]string{"field3"},
		},
		{
			"check min constraint",
			map[string]any{
				"field1": "test",
				"field2": "test",
				"field3": strings.Repeat("a", min-1),
			},
			nil,
			[]string{"field3"},
		},
		{
			"check max constraint",
			map[string]any{
				"field1": "test",
				"field2": "test",
				"field3": strings.Repeat("a", max+1),
			},
			nil,
			[]string{"field3"},
		},
		{
			"check pattern constraint",
			map[string]any{
				"field1": nil,
				"field2": "test",
				"field3": "test!",
			},
			nil,
			[]string{"field3"},
		},
		{
			"valid data (only required)",
			map[string]any{
				"field2": "test",
			},
			nil,
			[]string{},
		},
		{
			"valid data (all)",
			map[string]any{
				"field1": "test",
				"field2": 12345, // test value cast
				"field3": "test2",
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateNumber(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	min := 1.0
	max := 150.0
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeNumber,
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeNumber,
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeNumber,
			Options: &schema.NumberOptions{
				Min: &min,
				Max: &max,
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", 123)
	dummy.SetDataValue("field2", 123)
	dummy.SetDataValue("field3", 123)
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint + casting",
			map[string]any{
				"field1": "invalid",
				"field2": "invalid",
				"field3": "invalid",
			},
			nil,
			[]string{"field2", "field3"},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": 123,
				"field2": 123,
				"field3": 123,
			},
			nil,
			[]string{"field3"},
		},
		{
			"check min constraint",
			map[string]any{
				"field1": 0.5,
				"field2": 1,
				"field3": min - 0.5,
			},
			nil,
			[]string{"field3"},
		},
		{
			"check max constraint",
			map[string]any{
				"field1": nil,
				"field2": max,
				"field3": max + 0.5,
			},
			nil,
			[]string{"field3"},
		},
		{
			"valid data (only required)",
			map[string]any{
				"field2": 1,
			},
			nil,
			[]string{},
		},
		{
			"valid data (all)",
			map[string]any{
				"field1": nil,
				"field2": 123, // test value cast
				"field3": max,
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateBool(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeBool,
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeBool,
		},
		&schema.SchemaField{
			Name:    "field3",
			Unique:  true,
			Type:    schema.FieldTypeBool,
			Options: &schema.BoolOptions{},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", false)
	dummy.SetDataValue("field2", true)
	dummy.SetDataValue("field3", true)
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint + casting",
			map[string]any{
				"field1": "invalid",
				"field2": "invalid",
				"field3": "invalid",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": true,
				"field2": true,
				"field3": true,
			},
			nil,
			[]string{"field3"},
		},
		{
			"valid data (only required)",
			map[string]any{
				"field2": 1,
			},
			nil,
			[]string{},
		},
		{
			"valid data (all)",
			map[string]any{
				"field1": false,
				"field2": true,
				"field3": false,
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateEmail(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeEmail,
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeEmail,
			Options: &schema.EmailOptions{
				ExceptDomains: []string{"example.com"},
			},
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeEmail,
			Options: &schema.EmailOptions{
				OnlyDomains: []string{"example.com"},
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", "test@demo.com")
	dummy.SetDataValue("field2", "test@test.com")
	dummy.SetDataValue("field3", "test@example.com")
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check email format validator",
			map[string]any{
				"field1": "test",
				"field2": "test.com",
				"field3": 123,
			},
			nil,
			[]string{"field1", "field2", "field3"},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": "test@example.com",
				"field2": "test@test.com",
				"field3": "test@example.com",
			},
			nil,
			[]string{"field3"},
		},
		{
			"check ExceptDomains constraint",
			map[string]any{
				"field1": "test@example.com",
				"field2": "test@example.com",
				"field3": "test2@example.com",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check OnlyDomains constraint",
			map[string]any{
				"field1": "test@test.com",
				"field2": "test@test.com",
				"field3": "test@test.com",
			},
			nil,
			[]string{"field3"},
		},
		{
			"valid data (only required)",
			map[string]any{
				"field2": "test@test.com",
			},
			nil,
			[]string{},
		},
		{
			"valid data (all)",
			map[string]any{
				"field1": "123@example.com",
				"field2": "test@test.com",
				"field3": "test2@example.com",
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateUrl(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeUrl,
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeUrl,
			Options: &schema.UrlOptions{
				ExceptDomains: []string{"example.com"},
			},
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeUrl,
			Options: &schema.UrlOptions{
				OnlyDomains: []string{"example.com"},
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", "http://demo.com")
	dummy.SetDataValue("field2", "http://test.com")
	dummy.SetDataValue("field3", "http://example.com")
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check url format validator",
			map[string]any{
				"field1": "/abc",
				"field2": "test.com", // valid
				"field3": "test@example.com",
			},
			nil,
			[]string{"field1", "field3"},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": "http://example.com",
				"field2": "http://test.com",
				"field3": "http://example.com",
			},
			nil,
			[]string{"field3"},
		},
		{
			"check ExceptDomains constraint",
			map[string]any{
				"field1": "http://example.com",
				"field2": "http://example.com",
				"field3": "https://example.com",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check OnlyDomains constraint",
			map[string]any{
				"field1": "http://test.com/abc",
				"field2": "http://test.com/abc",
				"field3": "http://test.com/abc",
			},
			nil,
			[]string{"field3"},
		},
		{
			"check subdomains constraint",
			map[string]any{
				"field1": "http://test.test.com",
				"field2": "http://test.example.com",
				"field3": "http://test.example.com",
			},
			nil,
			[]string{"field3"},
		},
		{
			"valid data (only required)",
			map[string]any{
				"field2": "http://sub.test.com/abc",
			},
			nil,
			[]string{},
		},
		{
			"valid data (all)",
			map[string]any{
				"field1": "http://example.com/123",
				"field2": "http://test.com/",
				"field3": "http://example.com/test2",
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateDate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	min, _ := types.ParseDateTime("2022-01-01 01:01:01.123")
	max, _ := types.ParseDateTime("2030-01-01 01:01:01")
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeDate,
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeDate,
			Options: &schema.DateOptions{
				Min: min,
			},
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeDate,
			Options: &schema.DateOptions{
				Max: max,
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", "2022-01-01 01:01:01")
	dummy.SetDataValue("field2", "2029-01-01 01:01:01.123")
	dummy.SetDataValue("field3", "2029-01-01 01:01:01.123")
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint + cast",
			map[string]any{
				"field1": "invalid",
				"field2": "invalid",
				"field3": "invalid",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint + zero datetime",
			map[string]any{
				"field1": "January 1, year 1, 00:00:00 UTC",
				"field2": "0001-01-01 00:00:00",
				"field3": "0001-01-01 00:00:00 +0000 UTC",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": "2029-01-01 01:01:01.123",
				"field2": "2029-01-01 01:01:01.123",
				"field3": "2029-01-01 01:01:01.123",
			},
			nil,
			[]string{"field3"},
		},
		{
			"check min date constraint",
			map[string]any{
				"field1": "2021-01-01 01:01:01",
				"field2": "2021-01-01 01:01:01",
				"field3": "2021-01-01 01:01:01",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check max date constraint",
			map[string]any{
				"field1": "2030-02-01 01:01:01",
				"field2": "2030-02-01 01:01:01",
				"field3": "2030-02-01 01:01:01",
			},
			nil,
			[]string{"field3"},
		},
		{
			"valid data (only required)",
			map[string]any{
				"field2": "2029-01-01 01:01:01",
			},
			nil,
			[]string{},
		},
		{
			"valid data (all)",
			map[string]any{
				"field1": "2029-01-01 01:01:01.000",
				"field2": "2029-01-01 01:01:01",
				"field3": "2029-01-01 01:01:01.456",
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateSelect(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeSelect,
			Options: &schema.SelectOptions{
				Values:    []string{"1", "a", "b", "c"},
				MaxSelect: 1,
			},
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeSelect,
			Options: &schema.SelectOptions{
				Values:    []string{"a", "b", "c"},
				MaxSelect: 2,
			},
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeSelect,
			Options: &schema.SelectOptions{
				Values:    []string{"a", "b", "c"},
				MaxSelect: 99,
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", "a")
	dummy.SetDataValue("field2", []string{"a", "b"})
	dummy.SetDataValue("field3", []string{"a", "b", "c"})
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint - empty values",
			map[string]any{
				"field1": "",
				"field2": "",
				"field3": "",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint - multiple select cast",
			map[string]any{
				"field1": "a",
				"field2": "a",
				"field3": "a",
			},
			nil,
			[]string{},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": "a",
				"field2": "b",
				"field3": []string{"a", "b", "c"},
			},
			nil,
			[]string{"field3"},
		},
		{
			"check unique constraint - same elements but different order",
			map[string]any{
				"field1": "a",
				"field2": "b",
				"field3": []string{"a", "c", "b"},
			},
			nil,
			[]string{},
		},
		{
			"check Values constraint",
			map[string]any{
				"field1": 1,
				"field2": "d",
				"field3": 123,
			},
			nil,
			[]string{"field2", "field3"},
		},
		{
			"check MaxSelect constraint",
			map[string]any{
				"field1": []string{"a", "b"}, // this will be normalized to a single string value
				"field2": []string{"a", "b", "c"},
				"field3": []string{"a", "b", "b", "b"}, // repeating values will be merged
			},
			nil,
			[]string{"field2"},
		},
		{
			"valid data - only required fields",
			map[string]any{
				"field2": []string{"a", "b"},
			},
			nil,
			[]string{},
		},
		{
			"valid data - all fields with normalizations",
			map[string]any{
				"field1": "a",
				"field2": []string{"a", "b", "b"}, // will be collapsed
				"field3": "b",                     // will be normalzied to slice
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateJson(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeJson,
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeJson,
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeJson,
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", `{"test":123}`)
	dummy.SetDataValue("field2", `{"test":123}`)
	dummy.SetDataValue("field3", `{"test":123}`)
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint - nil",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint - zero string",
			map[string]any{
				"field1": "",
				"field2": "",
				"field3": "",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint - zero number",
			map[string]any{
				"field1": 0,
				"field2": 0,
				"field3": 0,
			},
			nil,
			[]string{},
		},
		{
			"check required constraint - zero slice",
			map[string]any{
				"field1": []string{},
				"field2": []string{},
				"field3": []string{},
			},
			nil,
			[]string{},
		},
		{
			"check required constraint - zero map",
			map[string]any{
				"field1": map[string]string{},
				"field2": map[string]string{},
				"field3": map[string]string{},
			},
			nil,
			[]string{},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": `{"test":123}`,
				"field2": `{"test":123}`,
				"field3": map[string]any{"test": 123},
			},
			nil,
			[]string{"field3"},
		},
		{
			"check json text validator",
			map[string]any{
				"field1": `[1, 2, 3`,
				"field2": `invalid`,
				"field3": `null`, // valid
			},
			nil,
			[]string{"field1", "field2"},
		},
		{
			"valid data - only required fields",
			map[string]any{
				"field2": `{"test":123}`,
			},
			nil,
			[]string{},
		},
		{
			"valid data - all fields with normalizations",
			map[string]any{
				"field1": []string{"a", "b", "c"},
				"field2": 123,
				"field3": `"test"`,
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateFile(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeFile,
			Options: &schema.FileOptions{
				MaxSelect: 1,
				MaxSize:   3,
			},
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeFile,
			Options: &schema.FileOptions{
				MaxSelect: 2,
				MaxSize:   10,
				MimeTypes: []string{"image/jpeg", "text/plain; charset=utf-8"},
			},
		},
		&schema.SchemaField{
			Name: "field3",
			Type: schema.FieldTypeFile,
			Options: &schema.FileOptions{
				MaxSelect: 3,
				MaxSize:   10,
				MimeTypes: []string{"image/jpeg"},
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// stub uploaded files
	data, mp, err := tests.MockMultipartData(nil, "test", "test", "test", "test", "test")
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/", data)
	req.Header.Add("Content-Type", mp.FormDataContentType())
	testFiles, err := rest.FindUploadedFiles(req, "test")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint - nil",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check MaxSelect constraint",
			map[string]any{
				"field1": "test1",
				"field2": []string{"test1", testFiles[0].Name(), testFiles[3].Name()},
				"field3": []string{"test1", "test2", "test3", "test4"},
			},
			[]*rest.UploadedFile{testFiles[0], testFiles[1], testFiles[2]},
			[]string{"field2", "field3"},
		},
		{
			"check MaxSize constraint",
			map[string]any{
				"field1": testFiles[0].Name(),
				"field2": []string{"test1", testFiles[0].Name()},
				"field3": []string{"test1", "test2", "test3"},
			},
			[]*rest.UploadedFile{testFiles[0], testFiles[1], testFiles[2]},
			[]string{"field1"},
		},
		{
			"check MimeTypes constraint",
			map[string]any{
				"field1": "test1",
				"field2": []string{"test1", testFiles[0].Name()},
				"field3": []string{testFiles[1].Name(), testFiles[2].Name()},
			},
			[]*rest.UploadedFile{testFiles[0], testFiles[1], testFiles[2]},
			[]string{"field3"},
		},
		{
			"valid data - no new files (just file ids)",
			map[string]any{
				"field1": "test1",
				"field2": []string{"test1", "test2"},
				"field3": []string{"test1", "test2", "test3"},
			},
			nil,
			[]string{},
		},
		{
			"valid data - just new files",
			map[string]any{
				"field1": nil,
				"field2": []string{testFiles[0].Name(), testFiles[1].Name()},
				"field3": nil,
			},
			[]*rest.UploadedFile{testFiles[0], testFiles[1], testFiles[2]},
			[]string{},
		},
		{
			"valid data - mixed existing and new files",
			map[string]any{
				"field1": "test1",
				"field2": []string{"test1", testFiles[0].Name()},
				"field3": "test1", // will be casted
			},
			[]*rest.UploadedFile{testFiles[0], testFiles[1], testFiles[2]},
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateRelation(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demo, _ := app.Dao().FindCollectionByNameOrId("demo4")

	// demo4 rel ids
	relId1 := "b8ba58f9-e2d7-42a0-b0e7-a11efd98236b"
	relId2 := "df55c8ff-45ef-4c82-8aed-6e2183fe1125"
	relId3 := "b84cd893-7119-43c9-8505-3c4e22da28a9"
	relId4 := "054f9f24-0a0a-4e09-87b1-bc7ff2b336a2"

	// record rel ids from different collections
	diffRelId1 := "63c2ab80-84ab-4057-a592-4604a731f78f"
	diffRelId2 := "2c542824-9de1-42fe-8924-e57c86267760"

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:    1,
				CollectionId: demo.Id,
			},
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:    2,
				CollectionId: demo.Id,
			},
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:    3,
				CollectionId: demo.Id,
			},
		},
		&schema.SchemaField{
			Name: "field4",
			Type: schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:    3,
				CollectionId: "", // missing or non-existing collection id
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", relId1)
	dummy.SetDataValue("field2", []string{relId1, relId2})
	dummy.SetDataValue("field3", []string{relId1, relId2, relId3})
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint - nil",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint - zero id",
			map[string]any{
				"field1": "",
				"field2": "",
				"field3": "",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": relId1,
				"field2": relId2,
				"field3": []string{relId1, relId2, relId3, relId3}, // repeating values are collapsed
			},
			nil,
			[]string{"field3"},
		},
		{
			"check nonexisting collection id",
			map[string]any{
				"field2": relId1,
				"field4": relId1,
			},
			nil,
			[]string{"field4"},
		},
		{
			"check MaxSelect constraint",
			map[string]any{
				"field1": []string{relId1, relId2}, // will be normalized to relId1 only
				"field2": []string{relId1, relId2, relId3},
				"field3": []string{relId1, relId2, relId3, relId4},
			},
			nil,
			[]string{"field2", "field3"},
		},
		{
			"check with ids from different collections",
			map[string]any{
				"field1": diffRelId1,
				"field2": []string{relId2, diffRelId1},
				"field3": []string{diffRelId1, diffRelId2},
			},
			nil,
			[]string{"field1", "field2", "field3"},
		},
		{
			"valid data - only required fields",
			map[string]any{
				"field2": []string{relId1, relId2},
			},
			nil,
			[]string{},
		},
		{
			"valid data - all fields with normalization",
			map[string]any{
				"field1": []string{relId1, relId2},
				"field2": relId2,
				"field3": []string{relId3, relId2, relId1}, // unique is not triggered because the order is different
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateUser(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	userId1 := "97cc3d3d-6ba2-383f-b42a-7bc84d27410c"
	userId2 := "7bc84d27-6ba2-b42a-383f-4197cc3d3d0c"
	userId3 := "4d0197cc-2b4a-3f83-a26b-d77bc8423d3c"
	missingUserId := "00000000-84ab-4057-a592-4604a731f78f"

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeUser,
			Options: &schema.UserOptions{
				MaxSelect: 1,
			},
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeUser,
			Options: &schema.UserOptions{
				MaxSelect: 2,
			},
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeUser,
			Options: &schema.UserOptions{
				MaxSelect: 3,
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.SetDataValue("field1", userId1)
	dummy.SetDataValue("field2", []string{userId1, userId2})
	dummy.SetDataValue("field3", []string{userId1, userId2, userId3})
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"check required constraint - nil",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"check required constraint - zero id",
			map[string]any{
				"field1": "",
				"field2": "",
				"field3": "",
			},
			nil,
			[]string{"field2"},
		},
		{
			"check unique constraint",
			map[string]any{
				"field1": nil,
				"field2": userId1,
				"field3": []string{userId1, userId2, userId3, userId3}, // repeating values are collapsed
			},
			nil,
			[]string{"field3"},
		},
		{
			"check MaxSelect constraint",
			map[string]any{
				"field1": []string{userId1, userId2}, // maxSelect is 1 and will be normalized to userId1 only
				"field2": []string{userId1, userId2, userId3},
				"field3": []string{userId1, userId3, userId2},
			},
			nil,
			[]string{"field2"},
		},
		{
			"check with mixed existing and nonexisting user ids",
			map[string]any{
				"field1": missingUserId,
				"field2": []string{missingUserId, userId1},
				"field3": []string{userId1, missingUserId},
			},
			nil,
			[]string{"field1", "field2", "field3"},
		},
		{
			"valid data - only required fields",
			map[string]any{
				"field2": []string{userId1, userId2},
			},
			nil,
			[]string{},
		},
		{
			"valid data - all fields with normalization",
			map[string]any{
				"field1": []string{userId1, userId2},
				"field2": userId2,
				"field3": []string{userId3, userId2, userId1}, // unique is not triggered because the order is different
			},
			nil,
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func checkValidatorErrors(t *testing.T, dao *daos.Dao, record *models.Record, scenarios []testDataFieldScenario) {
	for i, s := range scenarios {
		validator := validators.NewRecordDataValidator(dao, record, s.files)
		result := validator.Validate(s.data)

		prefix := fmt.Sprintf("%d", i)
		if s.name != "" {
			prefix = s.name
		}

		// parse errors
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("[%s] Failed to parse errors %v", prefix, result)
			continue
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("[%s] Expected error keys %v, got %v", prefix, s.expectedErrors, errs)
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("[%s] Missing expected error key %q in %v", prefix, k, errs)
			}
		}
	}
}
