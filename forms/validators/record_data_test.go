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
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/types"
)

type testDataFieldScenario struct {
	name           string
	data           map[string]any
	files          map[string][]*filesystem.File
	expectedErrors []string
}

func TestRecordDataValidatorEmptyAndUnknown(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo2")
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
	t.Parallel()

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
			Options: &schema.TextOptions{
				Pattern: pattern,
			},
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeText,
			Options: &schema.TextOptions{
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
	dummy.Set("field1", "test")
	dummy.Set("field2", "test")
	dummy.Set("field3", "test")
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"(text) check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"(text) check min constraint",
			map[string]any{
				"field1": "test",
				"field2": "test",
				"field3": strings.Repeat("a", min-1),
			},
			nil,
			[]string{"field3"},
		},
		{
			"(text) check min constraint with multi-bytes char",
			map[string]any{
				"field1": "test",
				"field2": "test",
				"field3": "𝌆", // 4 bytes should be counted as 1 char
			},
			nil,
			[]string{"field3"},
		},
		{
			"(text) check max constraint",
			map[string]any{
				"field1": "test",
				"field2": "test",
				"field3": strings.Repeat("a", max+1),
			},
			nil,
			[]string{"field3"},
		},
		{
			"(text) check max constraint with multi-bytes chars",
			map[string]any{
				"field1": "test",
				"field2": "test",
				"field3": strings.Repeat("𝌆", max), // shouldn't exceed the max limit even though max*4bytes chars are used
			},
			nil,
			[]string{},
		},
		{
			"(text) check pattern constraint",
			map[string]any{
				"field1": nil,
				"field2": "test!",
				"field3": "test",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(text) valid data (only required)",
			map[string]any{
				"field2": "test",
			},
			nil,
			[]string{},
		},
		{
			"(text) valid data (all)",
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
	min := 2.0
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
		&schema.SchemaField{
			Name: "field4",
			Type: schema.FieldTypeNumber,
			Options: &schema.NumberOptions{
				NoDecimal: true,
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.Set("field1", 123)
	dummy.Set("field2", 123)
	dummy.Set("field3", 123)
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"(number) check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
				"field4": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"(number) check required constraint + casting",
			map[string]any{
				"field1": "invalid",
				"field2": "invalid",
				"field3": "invalid",
				"field4": "invalid",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(number) check min constraint",
			map[string]any{
				"field1": 0.5,
				"field2": 1,
				"field3": min - 0.5,
			},
			nil,
			[]string{"field3"},
		},
		{
			"(number) check min with zero-default",
			map[string]any{
				"field2": 1,
				"field3": 0,
			},
			nil,
			[]string{},
		},
		{
			"(number) check max constraint",
			map[string]any{
				"field1": nil,
				"field2": max,
				"field3": max + 0.5,
			},
			nil,
			[]string{"field3"},
		},
		{
			"(number) check NoDecimal",
			map[string]any{
				"field2": 1,
				"field4": 456.789,
			},
			nil,
			[]string{"field4"},
		},
		{
			"(number) valid data (only required)",
			map[string]any{
				"field2": 1,
			},
			nil,
			[]string{},
		},
		{
			"(number) valid data (all)",
			map[string]any{
				"field1": nil,
				"field2": 123, // test value cast
				"field3": max,
				"field4": 456,
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
	dummy.Set("field1", false)
	dummy.Set("field2", true)
	dummy.Set("field3", true)
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"(bool) check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"(bool) check required constraint + casting",
			map[string]any{
				"field1": "invalid",
				"field2": "invalid",
				"field3": "invalid",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(bool) valid data (only required)",
			map[string]any{
				"field2": 1,
			},
			nil,
			[]string{},
		},
		{
			"(bool) valid data (all)",
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
	dummy.Set("field1", "test@demo.com")
	dummy.Set("field2", "test@test.com")
	dummy.Set("field3", "test@example.com")
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"(email) check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"(email) check email format validator",
			map[string]any{
				"field1": "test",
				"field2": "test.com",
				"field3": 123,
			},
			nil,
			[]string{"field1", "field2", "field3"},
		},
		{
			"(email) check ExceptDomains constraint",
			map[string]any{
				"field1": "test@example.com",
				"field2": "test@example.com",
				"field3": "test2@example.com",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(email) check OnlyDomains constraint",
			map[string]any{
				"field1": "test@test.com",
				"field2": "test@test.com",
				"field3": "test@test.com",
			},
			nil,
			[]string{"field3"},
		},
		{
			"(email) valid data (only required)",
			map[string]any{
				"field2": "test@test.com",
			},
			nil,
			[]string{},
		},
		{
			"(email) valid data (all)",
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
	dummy.Set("field1", "http://demo.com")
	dummy.Set("field2", "http://test.com")
	dummy.Set("field3", "http://example.com")
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"(url) check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"(url) check url format validator",
			map[string]any{
				"field1": "/abc",
				"field2": "test.com", // valid
				"field3": "test@example.com",
			},
			nil,
			[]string{"field1", "field3"},
		},
		{
			"(url) check ExceptDomains constraint",
			map[string]any{
				"field1": "http://example.com",
				"field2": "http://example.com",
				"field3": "https://example.com",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(url) check OnlyDomains constraint",
			map[string]any{
				"field1": "http://test.com/abc",
				"field2": "http://test.com/abc",
				"field3": "http://test.com/abc",
			},
			nil,
			[]string{"field3"},
		},
		{
			"(url) check subdomains constraint",
			map[string]any{
				"field1": "http://test.test.com",
				"field2": "http://test.example.com",
				"field3": "http://test.example.com",
			},
			nil,
			[]string{"field3"},
		},
		{
			"(url) valid data (only required)",
			map[string]any{
				"field2": "http://sub.test.com/abc",
			},
			nil,
			[]string{},
		},
		{
			"(url) valid data (all)",
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
	dummy.Set("field1", "2022-01-01 01:01:01")
	dummy.Set("field2", "2029-01-01 01:01:01.123")
	dummy.Set("field3", "2029-01-01 01:01:01.123")
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"(date) check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"(date) check required constraint + cast",
			map[string]any{
				"field1": "invalid",
				"field2": "invalid",
				"field3": "invalid",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(date) check required constraint + zero datetime",
			map[string]any{
				"field1": "January 1, year 1, 00:00:00 UTC",
				"field2": "0001-01-01 00:00:00",
				"field3": "0001-01-01 00:00:00 +0000 UTC",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(date) check min date constraint",
			map[string]any{
				"field1": "2021-01-01 01:01:01",
				"field2": "2021-01-01 01:01:01",
				"field3": "2021-01-01 01:01:01",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(date) check max date constraint",
			map[string]any{
				"field1": "2030-02-01 01:01:01",
				"field2": "2030-02-01 01:01:01",
				"field3": "2030-02-01 01:01:01",
			},
			nil,
			[]string{"field3"},
		},
		{
			"(date) valid data (only required)",
			map[string]any{
				"field2": "2029-01-01 01:01:01",
			},
			nil,
			[]string{},
		},
		{
			"(date) valid data (all)",
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
	dummy.Set("field1", "a")
	dummy.Set("field2", []string{"a", "b"})
	dummy.Set("field3", []string{"a", "b", "c"})
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"(select) check required constraint",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"(select) check required constraint - empty values",
			map[string]any{
				"field1": "",
				"field2": "",
				"field3": "",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(select) check required constraint - multiple select cast",
			map[string]any{
				"field1": "a",
				"field2": "a",
				"field3": "a",
			},
			nil,
			[]string{},
		},
		{
			"(select) check Values constraint",
			map[string]any{
				"field1": 1,
				"field2": "d",
				"field3": 123,
			},
			nil,
			[]string{"field2", "field3"},
		},
		{
			"(select) check MaxSelect constraint",
			map[string]any{
				"field1": []string{"a", "b"}, // this will be normalized to a single string value
				"field2": []string{"a", "b", "c"},
				"field3": []string{"a", "b", "b", "b"}, // repeating values will be merged
			},
			nil,
			[]string{"field2"},
		},
		{
			"(select) valid data - only required fields",
			map[string]any{
				"field2": []string{"a", "b"},
			},
			nil,
			[]string{},
		},
		{
			"(select) valid data - all fields with normalizations",
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
			Options: &schema.JsonOptions{
				MaxSize: 10,
			},
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeJson,
			Options: &schema.JsonOptions{
				MaxSize: 9999,
			},
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeJson,
			Options: &schema.JsonOptions{
				MaxSize: 9999,
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.Set("field1", `{"test":123}`)
	dummy.Set("field2", `{"test":123}`)
	dummy.Set("field3", `{"test":123}`)
	if err := app.Dao().SaveRecord(dummy); err != nil {
		t.Fatal(err)
	}

	scenarios := []testDataFieldScenario{
		{
			"(json) check required constraint - nil",
			map[string]any{
				"field1": nil,
				"field2": nil,
				"field3": nil,
			},
			nil,
			[]string{"field2"},
		},
		{
			"(json) check required constraint - zero string",
			map[string]any{
				"field1": "",
				"field2": "",
				"field3": "",
			},
			nil,
			[]string{"field2"},
		},
		{
			"(json) check required constraint - zero number",
			map[string]any{
				"field1": 0,
				"field2": 0,
				"field3": 0,
			},
			nil,
			[]string{},
		},
		{
			"(json) check required constraint - zero slice",
			map[string]any{
				"field1": []string{},
				"field2": []string{},
				"field3": []string{},
			},
			nil,
			[]string{"field2"},
		},
		{
			"(json) check required constraint - zero map",
			map[string]any{
				"field1": map[string]string{},
				"field2": map[string]string{},
				"field3": map[string]string{},
			},
			nil,
			[]string{"field2"},
		},
		{
			"(json) check MaxSize constraint",
			map[string]any{
				"field1": `"123456789"`, // max 10bytes
				"field2": 123,
			},
			nil,
			[]string{"field1"},
		},
		{
			"(json) check json text invalid obj, array and number normalizations",
			map[string]any{
				"field1": `[1 2 3]`,
				"field2": `{a: 123}`,
				"field3": `123.456 abc`,
			},
			nil,
			[]string{},
		},
		{
			"(json) check json text reserved literals normalizations",
			map[string]any{
				"field1": `true`,
				"field2": `false`,
				"field3": `null`,
			},
			nil,
			[]string{},
		},
		{
			"(json) valid data - only required fields",
			map[string]any{
				"field2": `{"test":123}`,
			},
			nil,
			[]string{},
		},
		{
			"(json) valid data - all fields with normalizations",
			map[string]any{
				"field1": `"12345678"`,
				"field2": 123,
				"field3": []string{"a", "b", "c"},
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
				"field2": []string{"test1", testFiles[0].Name, testFiles[3].Name},
				"field3": []string{"test1", "test2", "test3", "test4"},
			},
			map[string][]*filesystem.File{
				"field2": {testFiles[0], testFiles[3]},
			},
			[]string{"field2", "field3"},
		},
		{
			"check MaxSize constraint",
			map[string]any{
				"field1": testFiles[0].Name,
				"field2": []string{"test1", testFiles[0].Name},
				"field3": []string{"test1", "test2", "test3"},
			},
			map[string][]*filesystem.File{
				"field1": {testFiles[0]},
				"field2": {testFiles[0]},
			},
			[]string{"field1"},
		},
		{
			"check MimeTypes constraint",
			map[string]any{
				"field1": "test1",
				"field2": []string{"test1", testFiles[0].Name},
				"field3": []string{testFiles[1].Name, testFiles[2].Name},
			},
			map[string][]*filesystem.File{
				"field2": {testFiles[0], testFiles[1], testFiles[2]},
				"field3": {testFiles[1], testFiles[2]},
			},
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
				"field2": []string{testFiles[0].Name, testFiles[1].Name},
				"field3": nil,
			},
			map[string][]*filesystem.File{
				"field2": {testFiles[0], testFiles[1]},
			},
			[]string{},
		},
		{
			"valid data - mixed existing and new files",
			map[string]any{
				"field1": "test1",
				"field2": []string{"test1", testFiles[0].Name},
				"field3": "test1", // will be casted
			},
			map[string][]*filesystem.File{
				"field2": {testFiles[0], testFiles[1], testFiles[2]},
			},
			[]string{},
		},
	}

	checkValidatorErrors(t, app.Dao(), models.NewRecord(collection), scenarios)
}

func TestRecordDataValidatorValidateRelation(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demo, _ := app.Dao().FindCollectionByNameOrId("demo3")

	// demo3 rel ids
	relId1 := "mk5fmymtx4wsprk"
	relId2 := "7nwo8tuiatetxdm"
	relId3 := "lcl9d87w22ml6jy"
	relId4 := "1tmknxy2868d869"

	// record rel ids from different collections
	diffRelId1 := "0yxhwia2amd8gec"
	diffRelId2 := "llvuca81nly1qls"

	// create new test collection
	collection := &models.Collection{}
	collection.Name = "validate_test"
	collection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "field1",
			Type: schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:    types.Pointer(1),
				CollectionId: demo.Id,
			},
		},
		&schema.SchemaField{
			Name:     "field2",
			Required: true,
			Type:     schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:    types.Pointer(2),
				CollectionId: demo.Id,
			},
		},
		&schema.SchemaField{
			Name:   "field3",
			Unique: true,
			Type:   schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MinSelect:    types.Pointer(2),
				CollectionId: demo.Id,
			},
		},
		&schema.SchemaField{
			Name: "field4",
			Type: schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:    types.Pointer(3),
				CollectionId: "", // missing or non-existing collection id
			},
		},
	)
	if err := app.Dao().SaveCollection(collection); err != nil {
		t.Fatal(err)
	}

	// create dummy record (used for the unique check)
	dummy := models.NewRecord(collection)
	dummy.Set("field1", relId1)
	dummy.Set("field2", []string{relId1, relId2})
	dummy.Set("field3", []string{relId1, relId2, relId3})
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
			"check min constraint",
			map[string]any{
				"field2": relId2,
				"field3": []string{relId1},
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
			[]string{"field2"},
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

func checkValidatorErrors(t *testing.T, dao *daos.Dao, record *models.Record, scenarios []testDataFieldScenario) {
	for i, s := range scenarios {
		prefix := s.name
		if prefix == "" {
			prefix = fmt.Sprintf("%d", i)
		}

		t.Run(prefix, func(t *testing.T) {
			validator := validators.NewRecordDataValidator(dao, record, s.files)
			result := validator.Validate(s.data)

			// parse errors
			errs, ok := result.(validation.Errors)
			if !ok && result != nil {
				t.Fatalf("Failed to parse errors %v", result)
			}

			// check errors
			if len(errs) > len(s.expectedErrors) {
				t.Fatalf("Expected error keys %v, got %v", s.expectedErrors, errs)
			}
			for _, k := range s.expectedErrors {
				if _, ok := errs[k]; !ok {
					t.Fatalf("Missing expected error key %q in %v", k, errs)
				}
			}
		})
	}
}
