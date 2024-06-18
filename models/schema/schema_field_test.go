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
	expected := 9

	if len(result) != expected {
		t.Fatalf("Expected %d auth field names, got %d (%v)", expected, len(result), result)
	}
}

func TestFieldTypes(t *testing.T) {
	result := schema.FieldTypes()
	expected := 11

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
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeNumber, Name: "test"},
			"NUMERIC DEFAULT 0 NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeBool, Name: "test"},
			"BOOLEAN DEFAULT FALSE NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeEmail, Name: "test"},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeUrl, Name: "test"},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeEditor, Name: "test"},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeDate, Name: "test"},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeJson, Name: "test"},
			"JSON DEFAULT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeSelect, Name: "test"},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeSelect, Name: "test_multiple", Options: &schema.SelectOptions{MaxSelect: 2}},
			"JSON DEFAULT '[]' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeFile, Name: "test"},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeFile, Name: "test_multiple", Options: &schema.FileOptions{MaxSelect: 2}},
			"JSON DEFAULT '[]' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeRelation, Name: "test", Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			"TEXT DEFAULT '' NOT NULL",
		},
		{
			schema.SchemaField{Type: schema.FieldTypeRelation, Name: "test_multiple", Options: &schema.RelationOptions{MaxSelect: nil}},
			"JSON DEFAULT '[]' NOT NULL",
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
		Id:          "abc",
		Name:        "test",
		Type:        schema.FieldTypeText,
		Required:    true,
		Presentable: true,
		System:      true,
		Options: &schema.TextOptions{
			Pattern: "test",
		},
	}

	result := f.String()
	expected := `{"system":true,"id":"abc","name":"test","type":"text","required":true,"presentable":true,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}`

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
			`{"system":false,"id":"","name":"","type":"","required":false,"presentable":false,"unique":false,"options":null}`,
		},
		// without defined options
		{
			schema.SchemaField{
				Id:          "abc",
				Name:        "test",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: true,
				System:      true,
			},
			`{"system":true,"id":"abc","name":"test","type":"text","required":true,"presentable":true,"unique":false,"options":{"min":null,"max":null,"pattern":""}}`,
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
			`{"system":true,"id":"","name":"test","type":"text","required":true,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}`,
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
			`{"system":false,"id":"","name":"","type":"","required":false,"presentable":false,"unique":false,"options":null}`,
		},
		{
			[]byte{},
			true,
			`{"system":false,"id":"","name":"","type":"","required":false,"presentable":false,"unique":false,"options":null}`,
		},
		{
			[]byte(`{"system": true}`),
			true,
			`{"system":true,"id":"","name":"","type":"","required":false,"presentable":false,"unique":false,"options":null}`,
		},
		{
			[]byte(`{"invalid"`),
			true,
			`{"system":false,"id":"","name":"","type":"","required":false,"presentable":false,"unique":false,"options":null}`,
		},
		{
			[]byte(`{"type":"text","system":true}`),
			false,
			`{"system":true,"id":"","name":"","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}}`,
		},
		{
			[]byte(`{"type":"text","options":{"pattern":"test"}}`),
			false,
			`{"system":false,"id":"","name":"","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}`,
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
			"name with _via_",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: "a_via_b",
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
			"reserved name (_rowid_)",
			schema.SchemaField{
				Type: schema.FieldTypeText,
				Id:   "1234567890",
				Name: "_rowid_",
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
			`{"system":false,"id":"","name":"","type":"","required":false,"presentable":false,"unique":false,"options":null}`,
		},
		{
			schema.SchemaField{Type: "unknown"},
			true,
			`{"system":false,"id":"","name":"","type":"unknown","required":false,"presentable":false,"unique":false,"options":null}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeText},
			false,
			`{"system":false,"id":"","name":"","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeNumber},
			false,
			`{"system":false,"id":"","name":"","type":"number","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"noDecimal":false}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeBool},
			false,
			`{"system":false,"id":"","name":"","type":"bool","required":false,"presentable":false,"unique":false,"options":{}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeEmail},
			false,
			`{"system":false,"id":"","name":"","type":"email","required":false,"presentable":false,"unique":false,"options":{"exceptDomains":null,"onlyDomains":null}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeUrl},
			false,
			`{"system":false,"id":"","name":"","type":"url","required":false,"presentable":false,"unique":false,"options":{"exceptDomains":null,"onlyDomains":null}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeEditor},
			false,
			`{"system":false,"id":"","name":"","type":"editor","required":false,"presentable":false,"unique":false,"options":{"convertUrls":false}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeDate},
			false,
			`{"system":false,"id":"","name":"","type":"date","required":false,"presentable":false,"unique":false,"options":{"min":"","max":""}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeSelect},
			false,
			`{"system":false,"id":"","name":"","type":"select","required":false,"presentable":false,"unique":false,"options":{"maxSelect":0,"values":null}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeJson},
			false,
			`{"system":false,"id":"","name":"","type":"json","required":false,"presentable":false,"unique":false,"options":{"maxSize":0}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeFile},
			false,
			`{"system":false,"id":"","name":"","type":"file","required":false,"presentable":false,"unique":false,"options":{"mimeTypes":null,"thumbs":null,"maxSelect":0,"maxSize":0,"protected":false}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeRelation},
			false,
			`{"system":false,"id":"","name":"","type":"relation","required":false,"presentable":false,"unique":false,"options":{"collectionId":"","cascadeDelete":false,"minSelect":null,"maxSelect":null,"displayFields":null}}`,
		},
		{
			schema.SchemaField{Type: schema.FieldTypeUser},
			false,
			`{"system":false,"id":"","name":"","type":"user","required":false,"presentable":false,"unique":false,"options":{"maxSelect":0,"cascadeDelete":false}}`,
		},
		{
			schema.SchemaField{
				Type:    schema.FieldTypeText,
				Options: &schema.TextOptions{Pattern: "test"},
			},
			false,
			`{"system":false,"id":"","name":"","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}`,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("s%d_%s", i, s.field.Type), func(t *testing.T) {
			err := s.field.InitOptions()

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if s.field.String() != s.expectJson {
				t.Fatalf(" Expected\n%v\ngot\n%v", s.expectJson, s.field.String())
			}
		})
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

		// editor
		{schema.SchemaField{Type: schema.FieldTypeEditor}, nil, `""`},
		{schema.SchemaField{Type: schema.FieldTypeEditor}, "", `""`},
		{schema.SchemaField{Type: schema.FieldTypeEditor}, []int{1, 2}, `""`},
		{schema.SchemaField{Type: schema.FieldTypeEditor}, "test", `"test"`},
		{schema.SchemaField{Type: schema.FieldTypeEditor}, 123, `"123"`},

		// json
		{schema.SchemaField{Type: schema.FieldTypeJson}, nil, "null"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, "null", "null"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, 123, "123"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, -123, "-123"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, "123", "123"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, "-123", "-123"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, 123.456, "123.456"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, -123.456, "-123.456"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, "123.456", "123.456"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, "-123.456", "-123.456"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, "123.456 abc", `"123.456 abc"`}, // invalid numeric string
		{schema.SchemaField{Type: schema.FieldTypeJson}, "-a123", `"-a123"`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, true, "true"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, "true", "true"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, false, "false"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, "false", "false"},
		{schema.SchemaField{Type: schema.FieldTypeJson}, "", `""`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, `test`, `"test"`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, `"test"`, `"test"`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, `{test":1}`, `"{test\":1}"`}, // invalid object string
		{schema.SchemaField{Type: schema.FieldTypeJson}, `[1 2 3]`, `"[1 2 3]"`},      // invalid array string
		{schema.SchemaField{Type: schema.FieldTypeJson}, map[string]int{}, `{}`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, `{}`, `{}`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, map[string]int{"test": 123}, `{"test":123}`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, `{"test":123}`, `{"test":123}`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, []int{}, `[]`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, `[]`, `[]`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, []int{1, 2, 1}, `[1,2,1]`},
		{schema.SchemaField{Type: schema.FieldTypeJson}, `[1,2,1]`, `[1,2,1]`},

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
		{schema.SchemaField{Type: schema.FieldTypeSelect}, []string{"test1", "test2"}, `"test2"`},
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
		{schema.SchemaField{Type: schema.FieldTypeFile}, []string{"test1", "test2"}, `"test2"`},
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
			`"2ba88b4f-e9da-42f0-9764-9a55c953e724"`,
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

func TestSchemaFieldPrepareValueWithModifier(t *testing.T) {
	scenarios := []struct {
		name          string
		field         schema.SchemaField
		baseValue     any
		modifier      string
		modifierValue any
		expectJson    string
	}{
		// text
		{
			"text with '+' modifier",
			schema.SchemaField{Type: schema.FieldTypeText},
			"base",
			"+",
			"new",
			`"base"`,
		},
		{
			"text with '-' modifier",
			schema.SchemaField{Type: schema.FieldTypeText},
			"base",
			"-",
			"new",
			`"base"`,
		},
		{
			"text with unknown modifier",
			schema.SchemaField{Type: schema.FieldTypeText},
			"base",
			"?",
			"new",
			`"base"`,
		},
		{
			"text cast check",
			schema.SchemaField{Type: schema.FieldTypeText},
			123,
			"?",
			"new",
			`"123"`,
		},

		// number
		{
			"number with '+' modifier",
			schema.SchemaField{Type: schema.FieldTypeNumber},
			1,
			"+",
			4,
			`5`,
		},
		{
			"number with '-' modifier",
			schema.SchemaField{Type: schema.FieldTypeNumber},
			1,
			"-",
			4,
			`-3`,
		},
		{
			"number with unknown modifier",
			schema.SchemaField{Type: schema.FieldTypeNumber},
			"1",
			"?",
			4,
			`1`,
		},
		{
			"number cast check",
			schema.SchemaField{Type: schema.FieldTypeNumber},
			"test",
			"+",
			"4",
			`4`,
		},

		// bool
		{
			"bool with '+' modifier",
			schema.SchemaField{Type: schema.FieldTypeBool},
			true,
			"+",
			false,
			`true`,
		},
		{
			"bool with '-' modifier",
			schema.SchemaField{Type: schema.FieldTypeBool},
			true,
			"-",
			false,
			`true`,
		},
		{
			"bool with unknown modifier",
			schema.SchemaField{Type: schema.FieldTypeBool},
			true,
			"?",
			false,
			`true`,
		},
		{
			"bool cast check",
			schema.SchemaField{Type: schema.FieldTypeBool},
			"true",
			"?",
			false,
			`true`,
		},

		// email
		{
			"email with '+' modifier",
			schema.SchemaField{Type: schema.FieldTypeEmail},
			"base",
			"+",
			"new",
			`"base"`,
		},
		{
			"email with '-' modifier",
			schema.SchemaField{Type: schema.FieldTypeEmail},
			"base",
			"-",
			"new",
			`"base"`,
		},
		{
			"email with unknown modifier",
			schema.SchemaField{Type: schema.FieldTypeEmail},
			"base",
			"?",
			"new",
			`"base"`,
		},
		{
			"email cast check",
			schema.SchemaField{Type: schema.FieldTypeEmail},
			123,
			"?",
			"new",
			`"123"`,
		},

		// url
		{
			"url with '+' modifier",
			schema.SchemaField{Type: schema.FieldTypeUrl},
			"base",
			"+",
			"new",
			`"base"`,
		},
		{
			"url with '-' modifier",
			schema.SchemaField{Type: schema.FieldTypeUrl},
			"base",
			"-",
			"new",
			`"base"`,
		},
		{
			"url with unknown modifier",
			schema.SchemaField{Type: schema.FieldTypeUrl},
			"base",
			"?",
			"new",
			`"base"`,
		},
		{
			"url cast check",
			schema.SchemaField{Type: schema.FieldTypeUrl},
			123,
			"-",
			"new",
			`"123"`,
		},

		// editor
		{
			"editor with '+' modifier",
			schema.SchemaField{Type: schema.FieldTypeEditor},
			"base",
			"+",
			"new",
			`"base"`,
		},
		{
			"editor with '-' modifier",
			schema.SchemaField{Type: schema.FieldTypeEditor},
			"base",
			"-",
			"new",
			`"base"`,
		},
		{
			"editor with unknown modifier",
			schema.SchemaField{Type: schema.FieldTypeEditor},
			"base",
			"?",
			"new",
			`"base"`,
		},
		{
			"editor cast check",
			schema.SchemaField{Type: schema.FieldTypeEditor},
			123,
			"-",
			"new",
			`"123"`,
		},

		// date
		{
			"date with '+' modifier",
			schema.SchemaField{Type: schema.FieldTypeDate},
			"2023-01-01 00:00:00.123",
			"+",
			"2023-02-01 00:00:00.456",
			`"2023-01-01 00:00:00.123Z"`,
		},
		{
			"date with '-' modifier",
			schema.SchemaField{Type: schema.FieldTypeDate},
			"2023-01-01 00:00:00.123Z",
			"-",
			"2023-02-01 00:00:00.456Z",
			`"2023-01-01 00:00:00.123Z"`,
		},
		{
			"date with unknown modifier",
			schema.SchemaField{Type: schema.FieldTypeDate},
			"2023-01-01 00:00:00.123",
			"?",
			"2023-01-01 00:00:00.456",
			`"2023-01-01 00:00:00.123Z"`,
		},
		{
			"date cast check",
			schema.SchemaField{Type: schema.FieldTypeDate},
			1672524000, // 2022-12-31 22:00:00.000Z
			"+",
			100,
			`"2022-12-31 22:00:00.000Z"`,
		},

		// json
		{
			"json with '+' modifier",
			schema.SchemaField{Type: schema.FieldTypeJson},
			10,
			"+",
			5,
			`10`,
		},
		{
			"json with '+' modifier (slice)",
			schema.SchemaField{Type: schema.FieldTypeJson},
			[]string{"a", "b"},
			"+",
			"c",
			`["a","b"]`,
		},
		{
			"json with '-' modifier",
			schema.SchemaField{Type: schema.FieldTypeJson},
			10,
			"-",
			5,
			`10`,
		},
		{
			"json with '-' modifier (slice)",
			schema.SchemaField{Type: schema.FieldTypeJson},
			`["a","b"]`,
			"-",
			"c",
			`["a","b"]`,
		},
		{
			"json with unknown modifier",
			schema.SchemaField{Type: schema.FieldTypeJson},
			`"base"`,
			"?",
			`"new"`,
			`"base"`,
		},

		// single select
		{
			"single select with '+' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 1}},
			"",
			"+",
			"b",
			`"b"`,
		},
		{
			"single select with '+' modifier (nonempty base)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 1}},
			"a",
			"+",
			"b",
			`"b"`,
		},
		{
			"single select with '-' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 1}},
			"",
			"-",
			"a",
			`""`,
		},
		{
			"single select with '-' modifier (nonempty base and empty modifier value)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 1}},
			"a",
			"-",
			"",
			`"a"`,
		},
		{
			"single select with '-' modifier (nonempty base and different value)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 1}},
			"a",
			"-",
			"b",
			`"a"`,
		},
		{
			"single select with '-' modifier (nonempty base and matching value)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 1}},
			"a",
			"-",
			"a",
			`""`,
		},
		{
			"single select with '-' modifier (nonempty base and matching value in a slice)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 1}},
			"a",
			"-",
			[]string{"b", "a", "c", "123"},
			`""`,
		},
		{
			"single select with unknown modifier (nonempty)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 1}},
			"",
			"?",
			"a",
			`""`,
		},

		// multi select
		{
			"multi select with '+' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
			nil,
			"+",
			"b",
			`["b"]`,
		},
		{
			"multi select with '+' modifier (nonempty base)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
			[]string{"a"},
			"+",
			[]string{"b", "c"},
			`["a","b","c"]`,
		},
		{
			"multi select with '+' modifier (nonempty base; already existing value)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
			[]string{"a", "b"},
			"+",
			"b",
			`["a","b"]`,
		},
		{
			"multi select with '-' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
			nil,
			"-",
			[]string{"a"},
			`[]`,
		},
		{
			"multi select with '-' modifier (nonempty base and empty modifier value)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
			"a",
			"-",
			"",
			`["a"]`,
		},
		{
			"multi select with '-' modifier (nonempty base and different value)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
			"a",
			"-",
			"b",
			`["a"]`,
		},
		{
			"multi select with '-' modifier (nonempty base and matching value)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
			[]string{"a", "b", "c", "d"},
			"-",
			"c",
			`["a","b","d"]`,
		},
		{
			"multi select with '-' modifier (nonempty base and matching value in a slice)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
			[]string{"a", "b", "c", "d"},
			"-",
			[]string{"b", "a", "123"},
			`["c","d"]`,
		},
		{
			"multi select with unknown modifier (nonempty)",
			schema.SchemaField{Type: schema.FieldTypeSelect, Options: &schema.SelectOptions{MaxSelect: 10}},
			[]string{"a", "b"},
			"?",
			"a",
			`["a","b"]`,
		},

		// single relation
		{
			"single relation with '+' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			"",
			"+",
			"b",
			`"b"`,
		},
		{
			"single relation with '+' modifier (nonempty base)",
			schema.SchemaField{Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			"a",
			"+",
			"b",
			`"b"`,
		},
		{
			"single relation with '-' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			"",
			"-",
			"a",
			`""`,
		},
		{
			"single relation with '-' modifier (nonempty base and empty modifier value)",
			schema.SchemaField{Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			"a",
			"-",
			"",
			`"a"`,
		},
		{
			"single relation with '-' modifier (nonempty base and different value)",
			schema.SchemaField{Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			"a",
			"-",
			"b",
			`"a"`,
		},
		{
			"single relation with '-' modifier (nonempty base and matching value)",
			schema.SchemaField{Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			"a",
			"-",
			"a",
			`""`,
		},
		{
			"single relation with '-' modifier (nonempty base and matching value in a slice)",
			schema.SchemaField{Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			"a",
			"-",
			[]string{"b", "a", "c", "123"},
			`""`,
		},
		{
			"single relation with unknown modifier (nonempty)",
			schema.SchemaField{Type: schema.FieldTypeRelation, Options: &schema.RelationOptions{MaxSelect: types.Pointer(1)}},
			"",
			"?",
			"a",
			`""`,
		},

		// multi relation
		{
			"multi relation with '+' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeRelation},
			nil,
			"+",
			"b",
			`["b"]`,
		},
		{
			"multi relation with '+' modifier (nonempty base)",
			schema.SchemaField{Type: schema.FieldTypeRelation},
			[]string{"a"},
			"+",
			[]string{"b", "c"},
			`["a","b","c"]`,
		},
		{
			"multi relation with '+' modifier (nonempty base; already existing value)",
			schema.SchemaField{Type: schema.FieldTypeRelation},
			[]string{"a", "b"},
			"+",
			"b",
			`["a","b"]`,
		},
		{
			"multi relation with '-' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeRelation},
			nil,
			"-",
			[]string{"a"},
			`[]`,
		},
		{
			"multi relation with '-' modifier (nonempty base and empty modifier value)",
			schema.SchemaField{Type: schema.FieldTypeRelation},
			"a",
			"-",
			"",
			`["a"]`,
		},
		{
			"multi relation with '-' modifier (nonempty base and different value)",
			schema.SchemaField{Type: schema.FieldTypeRelation},
			"a",
			"-",
			"b",
			`["a"]`,
		},
		{
			"multi relation with '-' modifier (nonempty base and matching value)",
			schema.SchemaField{Type: schema.FieldTypeRelation},
			[]string{"a", "b", "c", "d"},
			"-",
			"c",
			`["a","b","d"]`,
		},
		{
			"multi relation with '-' modifier (nonempty base and matching value in a slice)",
			schema.SchemaField{Type: schema.FieldTypeRelation},
			[]string{"a", "b", "c", "d"},
			"-",
			[]string{"b", "a", "123"},
			`["c","d"]`,
		},
		{
			"multi relation with unknown modifier (nonempty)",
			schema.SchemaField{Type: schema.FieldTypeRelation},
			[]string{"a", "b"},
			"?",
			"a",
			`["a","b"]`,
		},

		// single file
		{
			"single file with '+' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 1}},
			"",
			"+",
			"b",
			`""`,
		},
		{
			"single file with '+' modifier (nonempty base)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 1}},
			"a",
			"+",
			"b",
			`"a"`,
		},
		{
			"single file with '-' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 1}},
			"",
			"-",
			"a",
			`""`,
		},
		{
			"single file with '-' modifier (nonempty base and empty modifier value)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 1}},
			"a",
			"-",
			"",
			`"a"`,
		},
		{
			"single file with '-' modifier (nonempty base and different value)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 1}},
			"a",
			"-",
			"b",
			`"a"`,
		},
		{
			"single file with '-' modifier (nonempty base and matching value)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 1}},
			"a",
			"-",
			"a",
			`""`,
		},
		{
			"single file with '-' modifier (nonempty base and matching value in a slice)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 1}},
			"a",
			"-",
			[]string{"b", "a", "c", "123"},
			`""`,
		},
		{
			"single file with unknown modifier (nonempty)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 1}},
			"",
			"?",
			"a",
			`""`,
		},

		// multi file
		{
			"multi file with '+' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 10}},
			nil,
			"+",
			"b",
			`[]`,
		},
		{
			"multi file with '+' modifier (nonempty base)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 10}},
			[]string{"a"},
			"+",
			[]string{"b", "c"},
			`["a"]`,
		},
		{
			"multi file with '+' modifier (nonempty base; already existing value)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 10}},
			[]string{"a", "b"},
			"+",
			"b",
			`["a","b"]`,
		},
		{
			"multi file with '-' modifier (empty base)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 10}},
			nil,
			"-",
			[]string{"a"},
			`[]`,
		},
		{
			"multi file with '-' modifier (nonempty base and empty modifier value)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 10}},
			"a",
			"-",
			"",
			`["a"]`,
		},
		{
			"multi file with '-' modifier (nonempty base and different value)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 10}},
			"a",
			"-",
			"b",
			`["a"]`,
		},
		{
			"multi file with '-' modifier (nonempty base and matching value)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 10}},
			[]string{"a", "b", "c", "d"},
			"-",
			"c",
			`["a","b","d"]`,
		},
		{
			"multi file with '-' modifier (nonempty base and matching value in a slice)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 10}},
			[]string{"a", "b", "c", "d"},
			"-",
			[]string{"b", "a", "123"},
			`["c","d"]`,
		},
		{
			"multi file with unknown modifier (nonempty)",
			schema.SchemaField{Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 10}},
			[]string{"a", "b"},
			"?",
			"a",
			`["a","b"]`,
		},
	}

	for _, s := range scenarios {
		result := s.field.PrepareValueWithModifier(s.baseValue, s.modifier, s.modifierValue)

		encoded, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("[%s] %v", s.name, err)
		}

		if string(encoded) != s.expectJson {
			t.Fatalf("[%s], Expected %v, got %v", s.name, s.expectJson, string(encoded))
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
	int1 := 10.0
	int2 := 20.0

	decimal1 := 10.5
	decimal2 := 20.5

	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.NumberOptions{},
			[]string{},
		},
		{
			"max - without min",
			schema.NumberOptions{
				Max: &int1,
			},
			[]string{},
		},
		{
			"max - failure with min",
			schema.NumberOptions{
				Min: &int2,
				Max: &int1,
			},
			[]string{"max"},
		},
		{
			"max - success with min",
			schema.NumberOptions{
				Min: &int1,
				Max: &int2,
			},
			[]string{},
		},
		{
			"NoDecimal range failure",
			schema.NumberOptions{
				Min:       &decimal1,
				Max:       &decimal2,
				NoDecimal: true,
			},
			[]string{"min", "max"},
		},
		{
			"NoDecimal range success",
			schema.NumberOptions{
				Min:       &int1,
				Max:       &int2,
				NoDecimal: true,
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

func TestEditorOptionsValidate(t *testing.T) {
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.EditorOptions{},
			[]string{},
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

func TestSelectOptionsIsMultiple(t *testing.T) {
	scenarios := []struct {
		maxSelect int
		expect    bool
	}{
		{-1, false},
		{0, false},
		{1, false},
		{2, true},
	}

	for i, s := range scenarios {
		opt := schema.SelectOptions{
			MaxSelect: s.maxSelect,
		}

		if v := opt.IsMultiple(); v != s.expect {
			t.Errorf("[%d] Expected %v, got %v", i, s.expect, v)
		}
	}
}

func TestJsonOptionsValidate(t *testing.T) {
	scenarios := []fieldOptionsScenario{
		{
			"empty",
			schema.JsonOptions{},
			[]string{"maxSize"},
		},
		{
			"MaxSize < 0",
			schema.JsonOptions{MaxSize: -1},
			[]string{"maxSize"},
		},
		{
			"MaxSize > 0",
			schema.JsonOptions{MaxSize: 1},
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

func TestFileOptionsIsMultiple(t *testing.T) {
	scenarios := []struct {
		maxSelect int
		expect    bool
	}{
		{-1, false},
		{0, false},
		{1, false},
		{2, true},
	}

	for i, s := range scenarios {
		opt := schema.FileOptions{
			MaxSelect: s.maxSelect,
		}

		if v := opt.IsMultiple(); v != s.expect {
			t.Errorf("[%d] Expected %v, got %v", i, s.expect, v)
		}
	}
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
			"MinSelect < 0",
			schema.RelationOptions{
				CollectionId: "abc",
				MinSelect:    types.Pointer(-1),
			},
			[]string{"minSelect"},
		},
		{
			"MinSelect >= 0",
			schema.RelationOptions{
				CollectionId: "abc",
				MinSelect:    types.Pointer(0),
			},
			[]string{},
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
			"MaxSelect > 0 && nonempty CollectionId",
			schema.RelationOptions{
				CollectionId: "abc",
				MaxSelect:    types.Pointer(1),
			},
			[]string{},
		},
		{
			"MinSelect < MaxSelect",
			schema.RelationOptions{
				CollectionId: "abc",
				MinSelect:    nil,
				MaxSelect:    types.Pointer(1),
			},
			[]string{},
		},
		{
			"MinSelect = MaxSelect (non-zero)",
			schema.RelationOptions{
				CollectionId: "abc",
				MinSelect:    types.Pointer(1),
				MaxSelect:    types.Pointer(1),
			},
			[]string{},
		},
		{
			"MinSelect = MaxSelect (both zero)",
			schema.RelationOptions{
				CollectionId: "abc",
				MinSelect:    types.Pointer(0),
				MaxSelect:    types.Pointer(0),
			},
			[]string{"maxSelect"},
		},
		{
			"MinSelect > MaxSelect",
			schema.RelationOptions{
				CollectionId: "abc",
				MinSelect:    types.Pointer(2),
				MaxSelect:    types.Pointer(1),
			},
			[]string{"maxSelect"},
		},
	}

	checkFieldOptionsScenarios(t, scenarios)
}

func TestRelationOptionsIsMultiple(t *testing.T) {
	scenarios := []struct {
		maxSelect *int
		expect    bool
	}{
		{nil, true},
		{types.Pointer(-1), false},
		{types.Pointer(0), false},
		{types.Pointer(1), false},
		{types.Pointer(2), true},
	}

	for i, s := range scenarios {
		opt := schema.RelationOptions{
			MaxSelect: s.maxSelect,
		}

		if v := opt.IsMultiple(); v != s.expect {
			t.Errorf("[%d] Expected %v, got %v", i, s.expect, v)
		}
	}
}
