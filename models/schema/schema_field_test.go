package schema_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestBaseModelFieldNames(t *testing.T) {
	result := schema.BaseModelFieldNames()
	expected := 3

	if len(result) != expected {
		t.Fatalf("Expected %d field names, got %d (%v)", expected, len(result), result)
	}
}

func TestSystemFieldNames(t *testing.T) {
	result := schema.SystemFieldNames()
	expected := 3

	if len(result) != expected {
		t.Fatalf("Expected %d field names, got %d (%v)", expected, len(result), result)
	}
}

func TestAuthFieldNames(t *testing.T) {
	result := schema.AuthFieldNames()
	expected := 8

	if len(result) != expected {
		t.Fatalf("Expected %d auth field names, got %d (%v)", expected, len(result), result)
	}
}

func TestFieldTypes(t *testing.T) {
	result := schema.FieldTypes()
	expected := 10

	if len(result) != expected {
		t.Fatalf("Expected %d types, got %d (%v)", expected, len(result), result)
	}
}

func TestArraybleFieldTypes(t *testing.T) {
	result := schema.ArraybleFieldTypes()
	expected := 3

	if len(result) != expected {
		t.Fatalf("Expected %d arrayble types, got %d (%v)", expected, len(result), result)
	}
}

func TestSchemaFieldColDefinition(t *testing.T) {
	scenarios := []struct {
		field    schema.SchemaField
		expected string
	}{
		{
			schema.SchemaField{Type: schema.FieldTypeText, Name: "test"},
			"TEXT DEFAULT ''",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeNumber, Name: "test"},
			"REAL DEFAULT 0",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeBool, Name: "test"},
			"BOOLEAN DEFAULT FALSE",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeEmail, Name: "test"},
			"TEXT DEFAULT ''",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeUrl, Name: "test"},
			"TEXT DEFAULT ''",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeDate, Name: "test"},
			"TEXT DEFAULT ''",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeSelect, Name: "test"},
			"TEXT DEFAULT ''",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeJson, Name: "test"},
			"JSON DEFAULT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeFile, Name: "test"},
			"TEXT DEFAULT ''",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeRelation, Name: "test"},
			"TEXT DEFAULT ''",
		},
	}

	for i, s := range scenarios {
		def := s.field.ColDefinition()
		if def != s.expected {
			t.Errorf("(%d) Expected definition %q, got %q", i, s.expected, def)
		}
	}
}

func TestSchemaFieldString(t *testing.T) {
	f := schema.SchemaField{
		Id:       "abc",
		Name:     "test",
		Type:     schema.FieldTypeText,
		Required: true,
		Unique:   false,
		System:   true,
		Options: &schema.TextOptions{
			Pattern: "test",
		},
	}

	result := f.String()
	expected := `{"system":true,"id":"abc","name":"test","type":"text","required":true,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}`

	if result != expected {
		t.Errorf("Expected \n%v, got \n%v", expected, result)
	}
}

func TestSchemaFieldMarshalJSON(t *testing.T) {
	scenarios := []struct {
		field    schema.SchemaField
		expected string
	}{
		// empty
		{
			schema.SchemaField{},
			`{"system":false,"id":"","name":"","type":"","required":false,"unique":false,"options":null}`,
		},
		// without defined options
		{
			schema.SchemaField{
				Id:       "abc",
				Name:     "test",
				Type:     schema.FieldTypeText,
				Required: true,
				Unique:   false,
				System:   true,
			},
			`{"system":true,"id":"abc","name":"test","type":"text","required":true,"unique":false,"options":{"min":null,"max":null,"pattern":""}}`,
		},
		// with defined options
		{
			schema.SchemaField{
				Name:     "test",
				Type:     schema.FieldTypeText,
				Required: true,
				Unique:   false,
				System:   true,
				Options: &schema.TextOptions{
					Pattern: "test",
				},
			},
			`{"system":true,"id":"","name":"test","type":"text","required":true,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}`,
		},
	}

	for i, s := range scenarios {
		result, err := s.field.MarshalJSON()
		if err != nil {
			t.Fatalf("(%d) %v", i, err)
		}

		if string(result) != s.expected {
			t.Errorf("(%d), Expected \n%v, got \n%v", i, s.expected, string(result))
		}
	}
}

func TestSchemaFieldUnmarshalJSON(t *testing.T) {
	scenarios := []struct {
		data        []byte
		expectError bool
		expectJson  string
	}{
		{
			nil,
			true,
			`{"system":false,"id":"","name":"","type":"","required":false,"unique":false,"options":null}`,
		},
		{
			[]byte{},
			true,
			`{"system":false,"id":"","name":"","type":"","required":false,"unique":false,"options":null}`,
		},
		{
			[]byte(`{"system": true}`),
			true,
			`{"system":true,"id":"","name":"","type":"","required":false,"unique":false,"options":null}`,
		},
		{
			[]byte(`{"invalid"`),
			true,
			`{"system":false,"id":"","name":"","type":"","required":false,"unique":false,"options":null}`,
		},
		{
			[]byte(`{"type":"text","system":true}`),
			false,
			`{"system":true,"id":"","name":"","type":"text","required":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}}`,
		},
		{
			[]byte(`{"type":"text","options":{"pattern":"test"}}`),
			false,
			`{"system":false,"id":"","name":"","type":"text","required":false,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}`,
		},
	}

	for i, s := range scenarios {
		f := schema.SchemaField{}
		err := f.UnmarshalJSON(s.data)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		if f.String() != s.expectJson {
			t.Errorf("(%d), Expected json \n%v, got \n%v", i, s.expectJson, f.String())
		}
	}
}

func TestSchemaFieldValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		field          schema.SchemaField
		expectedErrors []string
	}{
		{
			"empty field",
			schema.SchemaField{},
			[]string{"id", "options", "name", "type"},
		},
		{
			"missing id",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "",
				Name: "test",
			},
			[]string{"id"},
		},
		{
			"invalid id length check",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234",
				Name: "test",
			},
			[]string{"id"},
		},
		{
			"valid id length check",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "12345",
				Name: "test",
			},
			[]string{},
		},
		{
			"invalid name format",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: "test!@#",
			},
			[]string{"name"},
		},
		{
			"reserved name (null)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: "null",
			},
			[]string{"name"},
		},
		{
			"reserved name (true)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: "null",
			},
			[]string{"name"},
		},
		{
			"reserved name (false)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: "false",
			},
			[]string{"name"},
		},
		{
			"reserved name (id)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: schema.FieldNameId,
			},
			[]string{"name"},
		},
		{
			"reserved name (created)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: schema.FieldNameCreated,
			},
			[]string{"name"},
		},
		{
			"reserved name (updated)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: schema.FieldNameUpdated,
			},
			[]string{"name"},
		},
		{
			"reserved name (collectionId)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: schema.FieldNameCollectionId,
			},
			[]string{"name"},
		},
		{
			"reserved name (collectionName)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: schema.FieldNameCollectionName,
			},
			[]string{"name"},
		},
		{
			"reserved name (expand)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: schema.FieldNameExpand,
			},
			[]string{"name"},
		},
		{
			"valid name",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: "test",
			},
			[]string{},
		},
		{
			"unique check for type file",
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Id:      "1234567890",
				Name:    "test",
				Unique:  true,
				Options: &schema.FileOptions{MaxSelect: 1, MaxSize: 1},
			},
			[]string{"unique"},
		},
		{
			"trigger options validator (auto init)",
			schema.SchemaField{
				Type: schema.FieldTypeFile,
				Id:   "1234567890",
				Name: "test",
			},
			[]string{"options"},
		},
		{
			"trigger options validator (invalid option field value)",
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Id:      "1234567890",
				Name:    "test",
				Options: &schema.FileOptions{MaxSelect: 0, MaxSize: 0},
			},
			[]string{"options"},
		},
		{
			"trigger options validator (valid option field value)",
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Id:      "1234567890",
				Name:    "test",
				Options: &schema.FileOptions{MaxSelect: 1, MaxSize: 1},
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		result := s.field.Validate()

		// parse errors
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("[%s] Failed to parse errors %v", s.name, result)
			continue
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("[%s] Expected error keys %v, got %v", s.name, s.expectedErrors, errs)
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("[%s] Missing expected error key %q in %v", s.name, k, errs)
			}
		}
	}
}

func TestSchemaFieldInitOptions(t *testing.T) {
	scenarios := []struct {
		field       schema.SchemaField
		expectError bool
		expectJson  string
	}{
		{
			schema.SchemaField{},
			true,
			`{"system":false,"id":"","name":"","type":"","required":false,"unique":false,"options":null}`,
		},
		{
			schema.SchemaField{Type: "unknown"},
			true,
			`{"system":false,"id":"","name":"","type":"unknown","required":false,"unique":false,"options":null}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeText},
			false,
			`{"system":false,"id":"","name":"","type":"text","required":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeNumber},
			false,
			`{"system":false,"id":"","name":"","type":"number","required":false,"unique":false,"options":{"min":null,"max":null}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeBool},
			false,
			`{"system":false,"id":"","name":"","type":"bool","required":false,"unique":false,"options":{}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeEmail},
			false,
			`{"system":false,"id":"","name":"","type":"email","required":false,"unique":false,"options":{"exceptDomains":null,"onlyDomains":null}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeUrl},
			false,
			`{"system":false,"id":"","name":"","type":"url","required":false,"unique":false,"options":{"exceptDomains":null,"onlyDomains":null}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeDate},
			false,
			`{"system":false,"id":"","name":"","type":"date","required":false,"unique":false,"options":{"min":"","max":""}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeSelect},
			false,
			`{"system":false,"id":"","name":"","type":"select","required":false,"unique":false,"options":{"maxSelect":0,"values":null}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeJson},
			false,
			`{"system":false,"id":"","name":"","type":"json","required":false,"unique":false,"options":{}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeFile},
			false,
			`{"system":false,"id":"","name":"","type":"file","required":false,"unique":false,"options":{"maxSelect":0,"maxSize":0,"mimeTypes":null,"thumbs":null}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeRelation},
			false,
			`{"system":false,"id":"","name":"","type":"relation","required":false,"unique":false,"options":{"maxSelect":null,"collectionId":"","cascadeDelete":false}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeUser},
			false,
			`{"system":false,"id":"","name":"","type":"user","required":false,"unique":false,"options":{"maxSelect":0,"cascadeDelete":false}}`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeText,
				Options: &schema.TextOptions{Pattern: "test"},
			},
			false,
			`{"system":false,"id":"","name":"","type":"text","required":false,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}`,
		},
	}

	for i, s := range scenarios {
		err := s.field.InitOptions()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		if s.field.String() != s.expectJson {
			t.Errorf("(%d), Expected %v, got %v", i, s.expectJson, s.field.String())
		}
	}
}

func TestSchemaFieldPrepareValue(t *testing.T) {
	scenarios := []struct {
		field      schema.SchemaField
		value      any
		expectJson string
	}{
		{schema.SchemaField{Type: "unknown"}, "test", `"test"`},
		{schema.SchemaField{Type: "unknown"}, 123, "123"},
		{schema.SchemaField{Type: "unknown"}, []int{1, 2, 1}, "[1,2,1]"},

		// text
		{schema.SchemaField{Type: schema.FieldTypeText}, nil, `""`},
		{schema.SchemaField{Type: schema.FieldTypeText}, "", `""`},
		{schema.SchemaField{Type: schema.FieldTypeText}, []int{1, 2}, `""`},
		{schema.SchemaField{Type: schema.FieldTypeText}, "test", `"test"`},
		{schema.SchemaField{Type: schema.FieldTypeText}, 123, `"123"`},

		// email
		{schema.SchemaField{Type: schema.FieldTypeEmail}, nil, `""`},
		{schema.SchemaField{Type: schema.FieldTypeEmail}, "", `""`},
		{schema.SchemaField{Type: schema.FieldTypeEmail}, []int{1, 2}, `""`},
		{schema.SchemaField{Type: schema.FieldTypeEmail}, "test", `"test"`},
		{schema.SchemaField{Type: schema.FieldTypeEmail}, 123, `"123"`},

		// url
		{schema.SchemaField{Type: schema.FieldTypeUrl}, nil, `""`},
		{schema.SchemaField{Type: schema.FieldTypeUrl}, "", `""`},
		{schema.SchemaField{Type: schema.FieldTypeUrl}, []int{1, 2}, `""`},
		{schema.SchemaField{Type: schema.FieldTypeUrl}, "test", `"test"`},
		{schema.SchemaField{Type: schema.FieldTypeUrl}, 123, `"123"`},

		// json
		{schema.SchemaField{Type: schema.FieldTypeJson}, nil, "null"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, 123, "123"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, `"test"`, `"test"`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, map[string]int{"test": 123}, `{"test":123}`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, []int{1, 2, 1}, `[1,2,1]`},

		// number
		{schema.SchemaField{Type: schema.FieldTypeNumber}, nil, "0"},
		{schema.SchemaField{Type: schema.FieldTypeNumber}, "", "0"},
		{schema.SchemaField{Type: schema.FieldTypeNumber}, "test", "0"},
		{schema.SchemaField{Type: schema.FieldTypeNumber}, 1, "1"},
		{schema.SchemaField{Type: schema.FieldTypeNumber}, 1.5, "1.5"},
		{schema.SchemaField{Type: schema.FieldTypeNumber}, "1.5", "1.5"},

		// bool
		{schema.SchemaField{Type: schema.FieldTypeBool}, nil, "false"},
		{schema.SchemaField{Type: schema.FieldTypeBool}, 1, "true"},
		{schema.SchemaField{Type: schema.FieldTypeBool}, 0, "false"},
		{schema.SchemaField{Type: schema.FieldTypeBool}, "", "false"},
		{schema.SchemaField{Type: schema.FieldTypeBool}, "test", "false"},
		{schema.SchemaField{Type: schema.FieldTypeBool}, "false", "false"},
		{schema.SchemaField{Type: schema.FieldTypeBool}, "true", "true"},
		{schema.SchemaField{Type: schema.FieldTypeBool}, false, "false"},
		{schema.SchemaField{Type: schema.FieldTypeBool}, true, "true"},

		// date
		{schema.SchemaField{Type: schema.FieldTypeDate}, nil, `""`},
		{schema.SchemaField{Type: schema.FieldTypeDate}, "", `""`},
		{schema.SchemaField{Type: schema.FieldTypeDate}, "test", `""`},
		{schema.SchemaField{Type: schema.FieldTypeDate}, 1641024040, `"2022-01-01 08:00:40.000Z"`},
		{schema.SchemaField{Type: schema.FieldTypeDate}, "2022-01-01 11:27:10.123", `"2022-01-01 11:27:10.123Z"`},
		{schema.SchemaField{Type: schema.FieldTypeDate}, "2022-01-01 11:27:10.123Z", `"2022-01-01 11:27:10.123Z"`},
		{schema.SchemaField{Type: schema.FieldTypeDate}, types.DateTime{}, `""`},
		{schema.SchemaField{Type: schema.FieldTypeDate}, time.Time{}, `""`},

		// select (single)
		{schema.SchemaField{Type: schema.FieldTypeSelect}, nil, `""`},
		{schema.SchemaField{Type: schema.FieldTypeSelect}, "", `""`},
		{schema.SchemaField{Type: schema.FieldTypeSelect}, 123, `"123"`},
		{schema.SchemaField{Type: schema.FieldTypeSelect}, "test", `"test"`},
		{schema.SchemaField{Type: schema.FieldTypeSelect}, []string{"test1", "test2"}, `"test1"`},
		{
			// no values validation/filtering
			schema.SchemaField{
				Type: schema.FieldTypeSelect,
				Options: &schema.SelectOptions{
					Values: []string{"test1", "test2"},
				},
			},
			"test",
			`"test"`,
		},
		// select (multiple)
		{
			schema.SchemaField{
				Type:    schema.FieldTypeSelect,
				Options: &schema.SelectOptions{MaxSelect: 2},
			},
			nil,
			`[]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeSelect,
				Options: &schema.SelectOptions{MaxSelect: 2},
			},
			"",
			`[]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeSelect,
				Options: &schema.SelectOptions{MaxSelect: 2},
			},
			[]string{},
			`[]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeSelect,
				Options: &schema.SelectOptions{MaxSelect: 2},
			},
			123,
			`["123"]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeSelect,
				Options: &schema.SelectOptions{MaxSelect: 2},
			},
			"test",
			`["test"]`,
		},
		{
			// no values validation
			schema.SchemaField{
				Type:    schema.FieldTypeSelect,
				Options: &schema.SelectOptions{MaxSelect: 2},
			},
			[]string{"test1", "test2", "test3"},
			`["test1","test2","test3"]`,
		},
		{
			// duplicated values
			schema.SchemaField{
				Type:    schema.FieldTypeSelect,
				Options: &schema.SelectOptions{MaxSelect: 2},
			},
			[]string{"test1", "test2", "test1"},
			`["test1","test2"]`,
		},

		// file (single)
		{schema.SchemaField{Type: schema.FieldTypeFile}, nil, `""`},
		{schema.SchemaField{Type: schema.FieldTypeFile}, "", `""`},
		{schema.SchemaField{Type: schema.FieldTypeFile}, 123, `"123"`},
		{schema.SchemaField{Type: schema.FieldTypeFile}, "test", `"test"`},
		{schema.SchemaField{Type: schema.FieldTypeFile}, []string{"test1", "test2"}, `"test1"`},
		// file (multiple)
		{
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Options: &schema.FileOptions{MaxSelect: 2},
			},
			nil,
			`[]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Options: &schema.FileOptions{MaxSelect: 2},
			},
			"",
			`[]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Options: &schema.FileOptions{MaxSelect: 2},
			},
			[]string{},
			`[]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Options: &schema.FileOptions{MaxSelect: 2},
			},
			123,
			`["123"]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Options: &schema.FileOptions{MaxSelect: 2},
			},
			"test",
			`["test"]`,
		},
		{
			// no values validation
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Options: &schema.FileOptions{MaxSelect: 2},
			},
			[]string{"test1", "test2", "test3"},
			`["test1","test2","test3"]`,
		},
		{
			// duplicated values
			schema.SchemaField{
				Type:    schema.FieldTypeFile,
				Options: &schema.FileOptions{MaxSelect: 2},
			},
			[]string{"test1", "test2", "test1"},
			`["test1","test2"]`,
		},

		// relation (single)
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)},
			},
			nil,
			`""`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)},
			},
			"",
			`""`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)},
			},
			123,
			`"123"`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)},
			},
			"abc",
			`"abc"`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)},
			},
			"1ba88b4f-e9da-42f0-9764-9a55c953e724",
			`"1ba88b4f-e9da-42f0-9764-9a55c953e724"`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			[]string{"1ba88b4f-e9da-42f0-9764-9a55c953e724", "2ba88b4f-e9da-42f0-9764-9a55c953e724"},
			`"1ba88b4f-e9da-42f0-9764-9a55c953e724"`,
		},
		// relation (multiple)
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(2)},
			},
			nil,
			`[]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(2)},
			},
			"",
			`[]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(2)},
			},
			[]string{},
			`[]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(2)},
			},
			123,
			`["123"]`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(2)},
			},
			[]string{"", "abc"},
			`["abc"]`,
		},
		{
			// no values validation
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(2)},
			},
			[]string{"1ba88b4f-e9da-42f0-9764-9a55c953e724", "2ba88b4f-e9da-42f0-9764-9a55c953e724"},
			`["1ba88b4f-e9da-42f0-9764-9a55c953e724","2ba88b4f-e9da-42f0-9764-9a55c953e724"]`,
		},
		{
			// duplicated values
			schema.SchemaField{
				Type:    schema.FieldTypeRelation,
				Options: &schema.RelationOptions{MaxSelect: types.Pointer(2)},
			},
			[]string{"1ba88b4f-e9da-42f0-9764-9a55c953e724", "2ba88b4f-e9da-42f0-9764-9a55c953e724", "1ba88b4f-e9da-42f0-9764-9a55c953e724"},
			`["1ba88b4f-e9da-42f0-9764-9a55c953e724","2ba88b4f-e9da-42f0-9764-9a55c953e724"]`,
		},
	}

	for i, s := range scenarios {
		result := s.field.PrepareValue(s.value)

		encoded, err := json.Marshal(result)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
			continue
		}

		if string(encoded) != s.expectJson {
			t.Errorf("(%d), Expected %v, got %v", i, s.expectJson, string(encoded))
		}
	}
}

// -------------------------------------------------------------------

type fieldOptionsScenario struct {
	name           string
	options        schema.FieldOptions
	expectedErrors []string
}

func checkFieldOptionsScenarios(t *testing.T, scenarios []fieldOptionsScenario) {
	for i, s := range scenarios {
		result := s.options.Validate()

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

func TestTextOptionsValidate(t *testing.T) {
	minus := -1
	number0 := 0
	number1 := 10
	number2 := 20
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.TextOptions{},
			[]string{},
		},
		{
			"min - failure",
			schema.TextOptions{
				Min: &minus,
			},
			[]string{"min"},
		},
		{
			"min - success",
			schema.TextOptions{
				Min: &number0,
			},
			[]string{},
		},
		{
			"max - failure without min",
			schema.TextOptions{
				Max: &minus,
			},
			[]string{"max"},
		},
		{
			"max - failure with min",
			schema.TextOptions{
				Min: &number2,
				Max: &number1,
			},
			[]string{"max"},
		},
		{
			"max - success",
			schema.TextOptions{
				Min: &number1,
				Max: &number2,
			},
			[]string{},
		},
		{
			"pattern - failure",
			schema.TextOptions{Pattern: "(test"},
			[]string{"pattern"},
		},
		{
			"pattern - success",
			schema.TextOptions{Pattern: `^\#?\w+$`},
			[]string{},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestNumberOptionsValidate(t *testing.T) {
	number1 := 10.0
	number2 := 20.0
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.NumberOptions{},
			[]string{},
		},
		{
			"max - without min",
			schema.NumberOptions{
				Max: &number1,
			},
			[]string{},
		},
		{
			"max - failure with min",
			schema.NumberOptions{
				Min: &number2,
				Max: &number1,
			},
			[]string{"max"},
		},
		{
			"max - success with min",
			schema.NumberOptions{
				Min: &number1,
				Max: &number2,
			},
			[]string{},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestBoolOptionsValidate(t *testing.T) {
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.BoolOptions{},
			[]string{},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestEmailOptionsValidate(t *testing.T) {
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.EmailOptions{},
			[]string{},
		},
		{
			"ExceptDomains failure",
			schema.EmailOptions{
				ExceptDomains: []string{"invalid"},
			},
			[]string{"exceptDomains"},
		},
		{
			"ExceptDomains success",
			schema.EmailOptions{
				ExceptDomains: []string{"example.com", "sub.example.com"},
			},
			[]string{},
		},
		{
			"OnlyDomains check",
			schema.EmailOptions{
				OnlyDomains: []string{"invalid"},
			},
			[]string{"onlyDomains"},
		},
		{
			"OnlyDomains success",
			schema.EmailOptions{
				OnlyDomains: []string{"example.com", "sub.example.com"},
			},
			[]string{},
		},
		{
			"OnlyDomains + ExceptDomains at the same time",
			schema.EmailOptions{
				ExceptDomains: []string{"test1.com"},
				OnlyDomains:   []string{"test2.com"},
			},
			[]string{"exceptDomains", "onlyDomains"},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestUrlOptionsValidate(t *testing.T) {
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.UrlOptions{},
			[]string{},
		},
		{
			"ExceptDomains failure",
			schema.UrlOptions{
				ExceptDomains: []string{"invalid"},
			},
			[]string{"exceptDomains"},
		},
		{
			"ExceptDomains success",
			schema.UrlOptions{
				ExceptDomains: []string{"example.com", "sub.example.com"},
			},
			[]string{},
		},
		{
			"OnlyDomains check",
			schema.UrlOptions{
				OnlyDomains: []string{"invalid"},
			},
			[]string{"onlyDomains"},
		},
		{
			"OnlyDomains success",
			schema.UrlOptions{
				OnlyDomains: []string{"example.com", "sub.example.com"},
			},
			[]string{},
		},
		{
			"OnlyDomains + ExceptDomains at the same time",
			schema.UrlOptions{
				ExceptDomains: []string{"test1.com"},
				OnlyDomains:   []string{"test2.com"},
			},
			[]string{"exceptDomains", "onlyDomains"},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestDateOptionsValidate(t *testing.T) {
	date1 := types.NowDateTime()
	date2, _ := types.ParseDateTime(date1.Time().AddDate(1, 0, 0))

	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.DateOptions{},
			[]string{},
		},
		{
			"min only",
			schema.DateOptions{
				Min: date1,
			},
			[]string{},
		},
		{
			"max only",
			schema.DateOptions{
				Min: date1,
			},
			[]string{},
		},
		{
			"zero min + max",
			schema.DateOptions{
				Min: types.DateTime{},
				Max: date1,
			},
			[]string{},
		},
		{
			"min + zero max",
			schema.DateOptions{
				Min: date1,
				Max: types.DateTime{},
			},
			[]string{},
		},
		{
			"min > max",
			schema.DateOptions{
				Min: date2,
				Max: date1,
			},
			[]string{"max"},
		},
		{
			"min == max",
			schema.DateOptions{
				Min: date1,
				Max: date1,
			},
			[]string{"max"},
		},
		{
			"min < max",
			schema.DateOptions{
				Min: date1,
				Max: date2,
			},
			[]string{},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestSelectOptionsValidate(t *testing.T) {
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.SelectOptions{},
			[]string{"values", "maxSelect"},
		},
		{
			"MaxSelect <= 0",
			schema.SelectOptions{
				Values:    []string{"test1", "test2"},
				MaxSelect: 0,
			},
			[]string{"maxSelect"},
		},
		{
			"MaxSelect > Values",
			schema.SelectOptions{
				Values:    []string{"test1", "test2"},
				MaxSelect: 3,
			},
			[]string{"maxSelect"},
		},
		{
			"MaxSelect <= Values",
			schema.SelectOptions{
				Values:    []string{"test1", "test2"},
				MaxSelect: 2,
			},
			[]string{},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestJsonOptionsValidate(t *testing.T) {
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.JsonOptions{},
			[]string{},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestFileOptionsValidate(t *testing.T) {
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.FileOptions{},
			[]string{"maxSelect", "maxSize"},
		},
		{
			"MaxSelect <= 0 && maxSize <= 0",
			schema.FileOptions{
				MaxSize:   0,
				MaxSelect: 0,
			},
			[]string{"maxSelect", "maxSize"},
		},
		{
			"MaxSelect > 0 && maxSize > 0",
			schema.FileOptions{
				MaxSize:   2,
				MaxSelect: 1,
			},
			[]string{},
		},
		{
			"invalid thumbs format",
			schema.FileOptions{
				MaxSize:   1,
				MaxSelect: 2,
				Thumbs:    []string{"100", "200x100"},
			},
			[]string{"thumbs"},
		},
		{
			"invalid thumbs format - zero width and height",
			schema.FileOptions{
				MaxSize:   1,
				MaxSelect: 2,
				Thumbs:    []string{"0x0", "0x0t", "0x0b", "0x0f"},
			},
			[]string{"thumbs"},
		},
		{
			"valid thumbs format",
			schema.FileOptions{
				MaxSize:   1,
				MaxSelect: 2,
				Thumbs: []string{
					"100x100", "200x100", "0x100", "100x0",
					"10x10t", "10x10b", "10x10f",
				},
			},
			[]string{},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestRelationOptionsValidate(t *testing.T) {
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.RelationOptions{},
			[]string{"collectionId"},
		},
		{
			"empty CollectionId",
			schema.RelationOptions{
				CollectionId: "",
				MaxSelect:    types.Pointer(1),
			},
			[]string{"collectionId"},
		},
		{
			"MaxSelect <= 0",
			schema.RelationOptions{
				CollectionId: "abc",
				MaxSelect:    types.Pointer(0),
			},
			[]string{"maxSelect"},
		},
		{
			"MaxSelect > 0 && non-empty CollectionId",
			schema.RelationOptions{
				CollectionId: "abc",
				MaxSelect:    types.Pointer(1),
			},
			[]string{},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}
