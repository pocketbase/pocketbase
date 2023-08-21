package schema_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models/schema"
)

func TestNewSchemaAndFields(t *testing.T) {
	testSchema := schema.NewSchema(
		&schema.SchemaField{Id: "id1", Name: "test1"},
		&schema.SchemaField{Name: "test2"},
		&schema.SchemaField{Id: "id1", Name: "test1_new"}, // should replace the original id1 field
	)

	fields := testSchema.Fields()

	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d (%v)", len(fields), fields)
	}

	for _, f := range fields {
		if f.Id == "" {
			t.Fatalf("Expected field id to be set, found empty id for field %v", f)
		}
	}

	if fields[0].Name != "test1_new" {
		t.Fatalf("Expected field with name test1_new, got %s", fields[0].Name)
	}

	if fields[1].Name != "test2" {
		t.Fatalf("Expected field with name test2, got %s", fields[1].Name)
	}
}

func TestSchemaInitFieldsOptions(t *testing.T) {
	f0 := &schema.SchemaField{Name: "test1", Type: "unknown"}
	schema0 := schema.NewSchema(f0)

	err0 := schema0.InitFieldsOptions()
	if err0 == nil {
		t.Fatalf("Expected unknown field schema to fail, got nil")
	}

	// ---

	f1 := &schema.SchemaField{Name: "test1", Type: schema.FieldTypeText}
	f2 := &schema.SchemaField{Name: "test2", Type: schema.FieldTypeEmail}
	schema1 := schema.NewSchema(f1, f2)

	err1 := schema1.InitFieldsOptions()
	if err1 != nil {
		t.Fatal(err1)
	}

	if _, ok := f1.Options.(*schema.TextOptions); !ok {
		t.Fatalf("Failed to init f1 options")
	}

	if _, ok := f2.Options.(*schema.EmailOptions); !ok {
		t.Fatalf("Failed to init f2 options")
	}
}

func TestSchemaClone(t *testing.T) {
	f1 := &schema.SchemaField{Name: "test1", Type: schema.FieldTypeText}
	f2 := &schema.SchemaField{Name: "test2", Type: schema.FieldTypeEmail}
	s1 := schema.NewSchema(f1, f2)

	s2, err := s1.Clone()
	if err != nil {
		t.Fatal(err)
	}

	s1Encoded, _ := s1.MarshalJSON()
	s2Encoded, _ := s2.MarshalJSON()

	if string(s1Encoded) != string(s2Encoded) {
		t.Fatalf("Expected the cloned schema to be equal, got %v VS\n %v", s1, s2)
	}

	// change in one schema shouldn't result to change in the other
	// (aka. check if it is a deep clone)
	s1.Fields()[0].Name = "test1_update"
	if s2.Fields()[0].Name != "test1" {
		t.Fatalf("Expected s2 field name to not change, got %q", s2.Fields()[0].Name)
	}
}

func TestSchemaAsMap(t *testing.T) {
	f1 := &schema.SchemaField{Name: "test1", Type: schema.FieldTypeText}
	f2 := &schema.SchemaField{Name: "test2", Type: schema.FieldTypeEmail}
	testSchema := schema.NewSchema(f1, f2)

	result := testSchema.AsMap()

	if len(result) != 2 {
		t.Fatalf("Expected 2 map elements, got %d (%v)", len(result), result)
	}

	expectedIndexes := []string{f1.Name, f2.Name}

	for _, index := range expectedIndexes {
		if _, ok := result[index]; !ok {
			t.Fatalf("Missing index %q", index)
		}
	}
}

func TestSchemaGetFieldByName(t *testing.T) {
	f1 := &schema.SchemaField{Name: "test1", Type: schema.FieldTypeText}
	f2 := &schema.SchemaField{Name: "test2", Type: schema.FieldTypeText}
	testSchema := schema.NewSchema(f1, f2)

	// missing field
	result1 := testSchema.GetFieldByName("missing")
	if result1 != nil {
		t.Fatalf("Found unexpected field %v", result1)
	}

	// existing field
	result2 := testSchema.GetFieldByName("test1")
	if result2 == nil || result2.Name != "test1" {
		t.Fatalf("Cannot find field with Name 'test1', got %v ", result2)
	}
}

func TestSchemaGetFieldById(t *testing.T) {
	f1 := &schema.SchemaField{Id: "id1", Name: "test1", Type: schema.FieldTypeText}
	f2 := &schema.SchemaField{Id: "id2", Name: "test2", Type: schema.FieldTypeText}
	testSchema := schema.NewSchema(f1, f2)

	// missing field id
	result1 := testSchema.GetFieldById("test1")
	if result1 != nil {
		t.Fatalf("Found unexpected field %v", result1)
	}

	// existing field id
	result2 := testSchema.GetFieldById("id2")
	if result2 == nil || result2.Id != "id2" {
		t.Fatalf("Cannot find field with id 'id2', got %v ", result2)
	}
}

func TestSchemaRemoveField(t *testing.T) {
	f1 := &schema.SchemaField{Id: "id1", Name: "test1", Type: schema.FieldTypeText}
	f2 := &schema.SchemaField{Id: "id2", Name: "test2", Type: schema.FieldTypeText}
	f3 := &schema.SchemaField{Id: "id3", Name: "test3", Type: schema.FieldTypeText}
	testSchema := schema.NewSchema(f1, f2, f3)

	testSchema.RemoveField("id2")
	testSchema.RemoveField("test3") // should do nothing

	expected := []string{"test1", "test3"}

	if len(testSchema.Fields()) != len(expected) {
		t.Fatalf("Expected %d, got %d (%v)", len(expected), len(testSchema.Fields()), testSchema)
	}

	for _, name := range expected {
		if f := testSchema.GetFieldByName(name); f == nil {
			t.Fatalf("Missing field %q", name)
		}
	}
}

func TestSchemaAddField(t *testing.T) {
	f1 := &schema.SchemaField{Name: "test1", Type: schema.FieldTypeText}
	f2 := &schema.SchemaField{Id: "f2Id", Name: "test2", Type: schema.FieldTypeText}
	f3 := &schema.SchemaField{Id: "f3Id", Name: "test3", Type: schema.FieldTypeText}
	testSchema := schema.NewSchema(f1, f2, f3)

	f2New := &schema.SchemaField{Id: "f2Id", Name: "test2_new", Type: schema.FieldTypeEmail}
	f4 := &schema.SchemaField{Name: "test4", Type: schema.FieldTypeUrl}

	testSchema.AddField(f2New)
	testSchema.AddField(f4)

	if len(testSchema.Fields()) != 4 {
		t.Fatalf("Expected %d, got %d (%v)", 4, len(testSchema.Fields()), testSchema)
	}

	// check if each field has id
	for _, f := range testSchema.Fields() {
		if f.Id == "" {
			t.Fatalf("Expected field id to be set, found empty id for field %v", f)
		}
	}

	// check if f2 field was replaced
	if f := testSchema.GetFieldById("f2Id"); f == nil || f.Type != schema.FieldTypeEmail {
		t.Fatalf("Expected f2 field to be replaced, found %v", f)
	}

	// check if f4 was added
	if f := testSchema.GetFieldByName("test4"); f == nil || f.Name != "test4" {
		t.Fatalf("Expected f4 field to be added, found %v", f)
	}
}

func TestSchemaValidate(t *testing.T) {
	// emulate duplicated field ids
	duplicatedIdsSchema := schema.NewSchema(
		&schema.SchemaField{Id: "id1", Name: "test1", Type: schema.FieldTypeText},
		&schema.SchemaField{Id: "id2", Name: "test2", Type: schema.FieldTypeText},
	)
	duplicatedIdsSchema.Fields()[1].Id = "id1" // manually set existing id

	scenarios := []struct {
		schema      schema.Schema
		expectError bool
	}{
		// no fields
		{
			schema.NewSchema(),
			false,
		},
		// duplicated field ids
		{
			duplicatedIdsSchema,
			true,
		},
		// duplicated field names (case insensitive)
		{
			schema.NewSchema(
				&schema.SchemaField{Name: "test", Type: schema.FieldTypeText},
				&schema.SchemaField{Name: "TeSt", Type: schema.FieldTypeText},
			),
			true,
		},
		// failure - base individual fields validation
		{
			schema.NewSchema(
				&schema.SchemaField{Name: "", Type: schema.FieldTypeText},
			),
			true,
		},
		// success - base individual fields validation
		{
			schema.NewSchema(
				&schema.SchemaField{Name: "test", Type: schema.FieldTypeText},
			),
			false,
		},
		// failure - individual field options validation
		{
			schema.NewSchema(
				&schema.SchemaField{Name: "test", Type: schema.FieldTypeFile},
			),
			true,
		},
		// success - individual field options validation
		{
			schema.NewSchema(
				&schema.SchemaField{Name: "test", Type: schema.FieldTypeFile, Options: &schema.FileOptions{MaxSelect: 1, MaxSize: 1}},
			),
			false,
		},
	}

	for i, s := range scenarios {
		err := s.schema.Validate()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}
	}
}

func TestSchemaMarshalJSON(t *testing.T) {
	f1 := &schema.SchemaField{Id: "f1id", Name: "test1", Type: schema.FieldTypeText}
	f2 := &schema.SchemaField{
		Id:      "f2id",
		Name:    "test2",
		Type:    schema.FieldTypeText,
		Options: &schema.TextOptions{Pattern: "test"},
	}
	testSchema := schema.NewSchema(f1, f2)

	result, err := testSchema.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	expected := `[{"system":false,"id":"f1id","name":"test1","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}},{"system":false,"id":"f2id","name":"test2","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}]`

	if string(result) != expected {
		t.Fatalf("Expected %s, got %s", expected, string(result))
	}
}

func TestSchemaUnmarshalJSON(t *testing.T) {
	encoded := `[{"system":false,"id":"fid1", "name":"test1","type":"text","required":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}},{"system":false,"name":"test2","type":"text","required":false,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}]`
	testSchema := schema.Schema{}
	testSchema.AddField(&schema.SchemaField{Name: "tempField", Type: schema.FieldTypeUrl})
	err := testSchema.UnmarshalJSON([]byte(encoded))
	if err != nil {
		t.Fatal(err)
	}

	fields := testSchema.Fields()
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, found %v", fields)
	}

	f1 := testSchema.GetFieldByName("test1")
	if f1 == nil {
		t.Fatal("Expected to find field 'test1', got nil")
	}
	if f1.Id != "fid1" {
		t.Fatalf("Expected fid1 id, got %s", f1.Id)
	}
	_, ok := f1.Options.(*schema.TextOptions)
	if !ok {
		t.Fatal("'test1' field options are not inited.")
	}

	f2 := testSchema.GetFieldByName("test2")
	if f2 == nil {
		t.Fatal("Expected to find field 'test2', got nil")
	}
	if f2.Id == "" {
		t.Fatal("Expected f2 id to be set, got empty string")
	}
	o2, ok := f2.Options.(*schema.TextOptions)
	if !ok {
		t.Fatal("'test2' field options are not inited.")
	}
	if o2.Pattern != "test" {
		t.Fatalf("Expected pattern to be %q, got %q", "test", o2.Pattern)
	}
}

func TestSchemaValue(t *testing.T) {
	// empty schema
	s1 := schema.Schema{}
	v1, err := s1.Value()
	if err != nil {
		t.Fatal(err)
	}
	if v1 != "[]" {
		t.Fatalf("Expected nil, got %v", v1)
	}

	// schema with fields
	f1 := &schema.SchemaField{Id: "f1id", Name: "test1", Type: schema.FieldTypeText}
	s2 := schema.NewSchema(f1)

	v2, err := s2.Value()
	if err != nil {
		t.Fatal(err)
	}
	expected := `[{"system":false,"id":"f1id","name":"test1","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}}]`

	if v2 != expected {
		t.Fatalf("Expected %v, got %v", expected, v2)
	}
}

func TestSchemaScan(t *testing.T) {
	scenarios := []struct {
		data        any
		expectError bool
		expectJson  string
	}{
		{nil, false, "[]"},
		{"", false, "[]"},
		{[]byte{}, false, "[]"},
		{"[]", false, "[]"},
		{"invalid", true, "[]"},
		{123, true, "[]"},
		// no field type
		{`[{}]`, true, `[]`},
		// unknown field type
		{
			`[{"system":false,"id":"123","name":"test1","type":"unknown","required":false,"presentable":false,"unique":false}]`,
			true,
			`[]`,
		},
		// without options
		{
			`[{"system":false,"id":"123","name":"test1","type":"text","required":false,"presentable":false,"unique":false}]`,
			false,
			`[{"system":false,"id":"123","name":"test1","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":""}}]`,
		},
		// with options
		{
			`[{"system":false,"id":"123","name":"test1","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}]`,
			false,
			`[{"system":false,"id":"123","name":"test1","type":"text","required":false,"presentable":false,"unique":false,"options":{"min":null,"max":null,"pattern":"test"}}]`,
		},
	}

	for i, s := range scenarios {
		testSchema := schema.Schema{}

		err := testSchema.Scan(s.data)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}

		json, _ := testSchema.MarshalJSON()
		if string(json) != s.expectJson {
			t.Errorf("(%d) Expected json %v, got %v", i, s.expectJson, string(json))
		}
	}
}
