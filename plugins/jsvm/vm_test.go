package jsvm_test

import (
	"reflect"
	"testing"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/tests"
)

func TestBaseVMUnmarshal(t *testing.T) {
	vm := jsvm.NewBaseVM()

	v, err := vm.RunString(`unmarshal({ name: "test" }, new Collection())`)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := v.Export().(*models.Collection)
	if !ok {
		t.Fatalf("Expected models.Collection, got %v", m)
	}

	if m.Name != "test" {
		t.Fatalf("Expected collection with name %q, got %q", "test", m.Name)
	}
}

func TestBaseVMRecordBind(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	vm := jsvm.NewBaseVM()
	vm.Set("collection", collection)

	// without record data
	// ---
	v1, err := vm.RunString(`new Record(collection)`)
	if err != nil {
		t.Fatal(err)
	}

	m1, ok := v1.Export().(*models.Record)
	if !ok {
		t.Fatalf("Expected m1 to be models.Record, got \n%v", m1)
	}

	// with record data
	// ---
	v2, err := vm.RunString(`new Record(collection, { email: "test@example.com" })`)
	if err != nil {
		t.Fatal(err)
	}

	m2, ok := v2.Export().(*models.Record)
	if !ok {
		t.Fatalf("Expected m2 to be models.Record, got \n%v", m2)
	}

	if m2.Collection().Name != "users" {
		t.Fatalf("Expected record with collection %q, got \n%v", "users", m2.Collection())
	}

	if m2.Email() != "test@example.com" {
		t.Fatalf("Expected record with email field set to %q, got \n%v", "test@example.com", m2)
	}
}

// @todo enable after https://github.com/dop251/goja/issues/426
// func TestBaseVMRecordGetAndSetBind(t *testing.T) {
// 	app, _ := tests.NewTestApp()
// 	defer app.Cleanup()

// 	collection, err := app.Dao().FindCollectionByNameOrId("users")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	vm := jsvm.NewBaseVM()
// 	vm.Set("collection", collection)
// 	vm.Set("getRecord", func() *models.Record {
// 		return models.NewRecord(collection)
// 	})

// 	_, runErr := vm.RunString(`
// 		const jsRecord = new Record(collection);
// 		jsRecord.email = "test@example.com"; // test js record setter
// 		const email    = jsRecord.email; // test js record getter

// 		const goRecord = getRecord()
// 		goRecord.name  = "test" // test go record setter
// 		const name     = goRecord.name; // test go record getter
// 	`)
// 	if runErr != nil {
// 		t.Fatal(runErr)
// 	}

// 	expectedEmail := "test@example.com"
// 	expectedName := "test"

// 	jsRecord, ok := vm.Get("jsRecord").Export().(*models.Record)
// 	if !ok {
// 		t.Fatalf("Failed to export jsRecord")
// 	}
// 	if v := jsRecord.Email(); v != expectedEmail {
// 		t.Fatalf("Expected the js created record to have email %q, got %q", expectedEmail, v)
// 	}

// 	email := vm.Get("email").Export().(string)
// 	if email != expectedEmail {
// 		t.Fatalf("Expected exported email %q, got %q", expectedEmail, email)
// 	}

// 	goRecord, ok := vm.Get("goRecord").Export().(*models.Record)
// 	if !ok {
// 		t.Fatalf("Failed to export goRecord")
// 	}
// 	if v := goRecord.GetString("name"); v != expectedName {
// 		t.Fatalf("Expected the go created record to have name %q, got %q", expectedName, v)
// 	}

// 	name := vm.Get("name").Export().(string)
// 	if name != expectedName {
// 		t.Fatalf("Expected exported name %q, got %q", expectedName, name)
// 	}

// 	// ensure that the two record instances are not mixed
// 	if v := goRecord.Email(); v != "" {
// 		t.Fatalf("Expected the go created record to not have an email, got %q", v)
// 	}
// 	if v := jsRecord.GetString("name"); v != "" {
// 		t.Fatalf("Expected the js created record to not have a name, got %q", v)
// 	}
// }

func TestBaseVMCollectionBind(t *testing.T) {
	vm := jsvm.NewBaseVM()

	v, err := vm.RunString(`new Collection({ name: "test", schema: [{name: "title", "type": "text"}] })`)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := v.Export().(*models.Collection)
	if !ok {
		t.Fatalf("Expected models.Collection, got %v", m)
	}

	if m.Name != "test" {
		t.Fatalf("Expected collection with name %q, got %q", "test", m.Name)
	}

	if f := m.Schema.GetFieldByName("title"); f == nil {
		t.Fatalf("Expected schema to be set, got %v", m.Schema)
	}
}

func TestBaseVMAdminBind(t *testing.T) {
	vm := jsvm.NewBaseVM()

	v, err := vm.RunString(`new Admin({ email: "test@example.com" })`)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := v.Export().(*models.Admin)
	if !ok {
		t.Fatalf("Expected models.Admin, got %v", m)
	}
}

func TestBaseVMSchemaBind(t *testing.T) {
	vm := jsvm.NewBaseVM()

	v, err := vm.RunString(`new Schema([{name: "title", "type": "text"}])`)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := v.Export().(*schema.Schema)
	if !ok {
		t.Fatalf("Expected schema.Schema, got %v", m)
	}

	if f := m.GetFieldByName("title"); f == nil {
		t.Fatalf("Expected schema fields to be loaded, got %v", m.Fields())
	}
}

func TestBaseVMSchemaFieldBind(t *testing.T) {
	vm := jsvm.NewBaseVM()

	v, err := vm.RunString(`new SchemaField({name: "title", "type": "text"})`)
	if err != nil {
		t.Fatal(err)
	}

	f, ok := v.Export().(*schema.SchemaField)
	if !ok {
		t.Fatalf("Expected schema.SchemaField, got %v", f)
	}

	if f.Name != "title" {
		t.Fatalf("Expected field %q, got %v", "title", f)
	}
}

func TestBaseVMDaoBind(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	vm := jsvm.NewBaseVM()
	vm.Set("db", app.DB())

	v, err := vm.RunString(`new Dao(db)`)
	if err != nil {
		t.Fatal(err)
	}

	d, ok := v.Export().(*daos.Dao)
	if !ok {
		t.Fatalf("Expected daos.Dao, got %v", d)
	}

	if d.DB() != app.DB() {
		t.Fatalf("The db instances doesn't match")
	}
}

func TestFieldMapper(t *testing.T) {
	mapper := jsvm.FieldMapper{}

	scenarios := []struct {
		name     string
		expected string
	}{
		{"", ""},
		{"test", "test"},
		{"Test", "test"},
		{"miXeD", "miXeD"},
		{"MiXeD", "miXeD"},
		{"ResolveRequestAsJSON", "resolveRequestAsJSON"},
		{"Variable_with_underscore", "variable_with_underscore"},
		{"ALLCAPS", "allcaps"},
		{"NOTALLCAPs", "nOTALLCAPs"},
		{"ALL_CAPS_WITH_UNDERSCORE", "all_caps_with_underscore"},
	}

	for i, s := range scenarios {
		field := reflect.StructField{Name: s.name}
		if v := mapper.FieldName(nil, field); v != s.expected {
			t.Fatalf("[%d] Expected FieldName %q, got %q", i, s.expected, v)
		}

		method := reflect.Method{Name: s.name}
		if v := mapper.MethodName(nil, method); v != s.expected {
			t.Fatalf("[%d] Expected MethodName %q, got %q", i, s.expected, v)
		}
	}
}
