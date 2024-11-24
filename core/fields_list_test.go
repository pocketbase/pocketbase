package core_test

import (
	"bytes"
	"encoding/json"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

func TestNewFieldsList(t *testing.T) {
	fields := core.NewFieldsList(
		&core.TextField{Id: "id1", Name: "test1"},
		&core.TextField{Name: "test2"},
		&core.TextField{Id: "id1", Name: "test1_new"}, // should replace the original id1 field
	)

	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d (%v)", len(fields), fields)
	}

	for _, f := range fields {
		if f.GetId() == "" {
			t.Fatalf("Expected field id to be set, found empty id for field %v", f)
		}
	}

	if fields[0].GetName() != "test1_new" {
		t.Fatalf("Expected field with name test1_new, got %s", fields[0].GetName())
	}

	if fields[1].GetName() != "test2" {
		t.Fatalf("Expected field with name test2, got %s", fields[1].GetName())
	}
}

func TestFieldsListClone(t *testing.T) {
	f1 := &core.TextField{Name: "test1"}
	f2 := &core.EmailField{Name: "test2"}
	s1 := core.NewFieldsList(f1, f2)

	s2, err := s1.Clone()
	if err != nil {
		t.Fatal(err)
	}

	s1Str := s1.String()
	s2Str := s2.String()

	if s1Str != s2Str {
		t.Fatalf("Expected the cloned list to be equal, got \n%v\nVS\n%v", s1, s2)
	}

	// change in one list shouldn't result to change in the other
	// (aka. check if it is a deep clone)
	s1[0].SetName("test1_update")
	if s2[0].GetName() != "test1" {
		t.Fatalf("Expected s2 field name to not change, got %q", s2[0].GetName())
	}
}

func TestFieldsListFieldNames(t *testing.T) {
	f1 := &core.TextField{Name: "test1"}
	f2 := &core.EmailField{Name: "test2"}
	testFieldsList := core.NewFieldsList(f1, f2)

	result := testFieldsList.FieldNames()

	expected := []string{f1.Name, f2.Name}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d slice elements, got %d\n%v", len(expected), len(result), result)
	}

	for _, name := range expected {
		if !slices.Contains(result, name) {
			t.Fatalf("Missing name %q in %v", name, result)
		}
	}
}

func TestFieldsListAsMap(t *testing.T) {
	f1 := &core.TextField{Name: "test1"}
	f2 := &core.EmailField{Name: "test2"}
	testFieldsList := core.NewFieldsList(f1, f2)

	result := testFieldsList.AsMap()

	expectedIndexes := []string{f1.Name, f2.Name}

	if len(result) != len(expectedIndexes) {
		t.Fatalf("Expected %d map elements, got %d\n%v", len(expectedIndexes), len(result), result)
	}

	for _, index := range expectedIndexes {
		if _, ok := result[index]; !ok {
			t.Fatalf("Missing index %q", index)
		}
	}
}

func TestFieldsListGetById(t *testing.T) {
	f1 := &core.TextField{Id: "id1", Name: "test1"}
	f2 := &core.EmailField{Id: "id2", Name: "test2"}
	testFieldsList := core.NewFieldsList(f1, f2)

	// missing field id
	result1 := testFieldsList.GetById("test1")
	if result1 != nil {
		t.Fatalf("Found unexpected field %v", result1)
	}

	// existing field id
	result2 := testFieldsList.GetById("id2")
	if result2 == nil || result2.GetId() != "id2" {
		t.Fatalf("Cannot find field with id %q, got %v ", "id2", result2)
	}
}

func TestFieldsListGetByName(t *testing.T) {
	f1 := &core.TextField{Id: "id1", Name: "test1"}
	f2 := &core.EmailField{Id: "id2", Name: "test2"}
	testFieldsList := core.NewFieldsList(f1, f2)

	// missing field name
	result1 := testFieldsList.GetByName("id1")
	if result1 != nil {
		t.Fatalf("Found unexpected field %v", result1)
	}

	// existing field name
	result2 := testFieldsList.GetByName("test2")
	if result2 == nil || result2.GetName() != "test2" {
		t.Fatalf("Cannot find field with name %q, got %v ", "test2", result2)
	}
}

func TestFieldsListRemove(t *testing.T) {
	testFieldsList := core.NewFieldsList(
		&core.TextField{Id: "id1", Name: "test1"},
		&core.TextField{Id: "id2", Name: "test2"},
		&core.TextField{Id: "id3", Name: "test3"},
		&core.TextField{Id: "id4", Name: "test4"},
		&core.TextField{Id: "id5", Name: "test5"},
		&core.TextField{Id: "id6", Name: "test6"},
	)

	// remove by id
	testFieldsList.RemoveById("id2")
	testFieldsList.RemoveById("test3") // should do nothing

	// remove by name
	testFieldsList.RemoveByName("test5")
	testFieldsList.RemoveByName("id6") // should do nothing

	expected := []string{"test1", "test3", "test4", "test6"}

	if len(testFieldsList) != len(expected) {
		t.Fatalf("Expected %d, got %d\n%v", len(expected), len(testFieldsList), testFieldsList)
	}

	for _, name := range expected {
		if f := testFieldsList.GetByName(name); f == nil {
			t.Fatalf("Missing field %q", name)
		}
	}
}

func TestFieldsListAdd(t *testing.T) {
	f0 := &core.TextField{}
	f1 := &core.TextField{Name: "test1"}
	f2 := &core.TextField{Id: "f2Id", Name: "test2"}
	f3 := &core.TextField{Id: "f3Id", Name: "test3"}
	testFieldsList := core.NewFieldsList(f0, f1, f2, f3)

	f2New := &core.EmailField{Id: "f2Id", Name: "test2_new"}
	f4 := &core.URLField{Name: "test4"}

	testFieldsList.Add(f2New)
	testFieldsList.Add(f4)

	if len(testFieldsList) != 5 {
		t.Fatalf("Expected %d, got %d\n%v", 5, len(testFieldsList), testFieldsList)
	}

	// check if each field has id
	for _, f := range testFieldsList {
		if f.GetId() == "" {
			t.Fatalf("Expected field id to be set, found empty id for field %v", f)
		}
	}

	// check if f2 field was replaced
	if f := testFieldsList.GetById("f2Id"); f == nil || f.Type() != core.FieldTypeEmail {
		t.Fatalf("Expected f2 field to be replaced, found %v", f)
	}

	// check if f4 was added
	if f := testFieldsList.GetByName("test4"); f == nil || f.GetName() != "test4" {
		t.Fatalf("Expected f4 field to be added, found %v", f)
	}
}

func TestFieldsListAddMarshaledJSON(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name           string
		raw            []byte
		expectError    bool
		expectedFields map[string]string
	}{
		{
			"nil",
			nil,
			false,
			map[string]string{"abc": core.FieldTypeNumber},
		},
		{
			"empty array",
			[]byte(`[]`),
			false,
			map[string]string{"abc": core.FieldTypeNumber},
		},
		{
			"empty object",
			[]byte(`{}`),
			true,
			map[string]string{"abc": core.FieldTypeNumber},
		},
		{
			"array with empty object",
			[]byte(`[{}]`),
			true,
			map[string]string{"abc": core.FieldTypeNumber},
		},
		{
			"single object with invalid type",
			[]byte(`{"type":"missing","name":"test"}`),
			true,
			map[string]string{"abc": core.FieldTypeNumber},
		},
		{
			"single object with valid type",
			[]byte(`{"type":"text","name":"test"}`),
			false,
			map[string]string{
				"abc":  core.FieldTypeNumber,
				"test": core.FieldTypeText,
			},
		},
		{
			"array of object with valid types",
			[]byte(`[{"type":"text","name":"test1"},{"type":"url","name":"test2"}]`),
			false,
			map[string]string{
				"abc":   core.FieldTypeNumber,
				"test1": core.FieldTypeText,
				"test2": core.FieldTypeURL,
			},
		},
		{
			"fields with duplicated ids should replace existing fields",
			[]byte(`[{"type":"text","name":"test1"},{"type":"url","name":"test2"},{"type":"text","name":"abc2", "id":"abc_id"}]`),
			false,
			map[string]string{
				"abc2":  core.FieldTypeText,
				"test1": core.FieldTypeText,
				"test2": core.FieldTypeURL,
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testList := core.NewFieldsList(&core.NumberField{Name: "abc", Id: "abc_id"})
			err := testList.AddMarshaledJSON(s.raw)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectError, hasErr)
			}

			if len(s.expectedFields) != len(testList) {
				t.Fatalf("Expected %d fields, got %d", len(s.expectedFields), len(testList))
			}

			for fieldName, typ := range s.expectedFields {
				f := testList.GetByName(fieldName)

				if f == nil {
					t.Errorf("Missing expected field %q", fieldName)
					continue
				}

				if f.Type() != typ {
					t.Errorf("Expect field %q to has type %q, got %q", fieldName, typ, f.Type())
				}
			}
		})
	}
}

func TestFieldsListAddAt(t *testing.T) {
	scenarios := []struct {
		position int
		expected []string
	}{
		{-2, []string{"test1", "test2_new", "test3", "test4"}},
		{-1, []string{"test1", "test2_new", "test3", "test4"}},
		{0, []string{"test2_new", "test4", "test1", "test3"}},
		{1, []string{"test1", "test2_new", "test4", "test3"}},
		{2, []string{"test1", "test3", "test2_new", "test4"}},
		{3, []string{"test1", "test3", "test2_new", "test4"}},
		{4, []string{"test1", "test3", "test2_new", "test4"}},
		{5, []string{"test1", "test3", "test2_new", "test4"}},
	}

	for _, s := range scenarios {
		t.Run(strconv.Itoa(s.position), func(t *testing.T) {
			f1 := &core.TextField{Id: "f1Id", Name: "test1"}
			f2 := &core.TextField{Id: "f2Id", Name: "test2"}
			f3 := &core.TextField{Id: "f3Id", Name: "test3"}
			testFieldsList := core.NewFieldsList(f1, f2, f3)

			f2New := &core.EmailField{Id: "f2Id", Name: "test2_new"}
			f4 := &core.URLField{Name: "test4"}
			testFieldsList.AddAt(s.position, f2New, f4)

			rawNames, err := json.Marshal(testFieldsList.FieldNames())
			if err != nil {
				t.Fatal(err)
			}

			rawExpected, err := json.Marshal(s.expected)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(rawNames, rawExpected) {
				t.Fatalf("Expected fields\n%s\ngot\n%s", rawExpected, rawNames)
			}
		})
	}
}

func TestFieldsListAddMarshaledJSONAt(t *testing.T) {
	scenarios := []struct {
		position int
		expected []string
	}{
		{-2, []string{"test1", "test2_new", "test3", "test4"}},
		{-1, []string{"test1", "test2_new", "test3", "test4"}},
		{0, []string{"test2_new", "test4", "test1", "test3"}},
		{1, []string{"test1", "test2_new", "test4", "test3"}},
		{2, []string{"test1", "test3", "test2_new", "test4"}},
		{3, []string{"test1", "test3", "test2_new", "test4"}},
		{4, []string{"test1", "test3", "test2_new", "test4"}},
		{5, []string{"test1", "test3", "test2_new", "test4"}},
	}

	for _, s := range scenarios {
		t.Run(strconv.Itoa(s.position), func(t *testing.T) {
			f1 := &core.TextField{Id: "f1Id", Name: "test1"}
			f2 := &core.TextField{Id: "f2Id", Name: "test2"}
			f3 := &core.TextField{Id: "f3Id", Name: "test3"}
			testFieldsList := core.NewFieldsList(f1, f2, f3)

			err := testFieldsList.AddMarshaledJSONAt(s.position, []byte(`[
				{"id":"f2Id", "name":"test2_new", "type": "text"},
				{"name": "test4", "type": "text"}
			]`))
			if err != nil {
				t.Fatal(err)
			}

			rawNames, err := json.Marshal(testFieldsList.FieldNames())
			if err != nil {
				t.Fatal(err)
			}

			rawExpected, err := json.Marshal(s.expected)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(rawNames, rawExpected) {
				t.Fatalf("Expected fields\n%s\ngot\n%s", rawExpected, rawNames)
			}
		})
	}
}

func TestFieldsListStringAndValue(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		testFieldsList := core.NewFieldsList()

		str := testFieldsList.String()
		if str != "[]" {
			t.Fatalf("Expected empty slice, got\n%q", str)
		}

		v, err := testFieldsList.Value()
		if err != nil {
			t.Fatal(err)
		}
		if v != str {
			t.Fatalf("Expected String and Value to match")
		}
	})

	t.Run("list with fields", func(t *testing.T) {
		testFieldsList := core.NewFieldsList(
			&core.TextField{Id: "f1id", Name: "test1"},
			&core.BoolField{Id: "f2id", Name: "test2"},
			&core.URLField{Id: "f3id", Name: "test3"},
		)

		str := testFieldsList.String()

		v, err := testFieldsList.Value()
		if err != nil {
			t.Fatal(err)
		}
		if v != str {
			t.Fatalf("Expected String and Value to match")
		}

		expectedParts := []string{
			`"type":"bool"`,
			`"type":"url"`,
			`"type":"text"`,
			`"id":"f1id"`,
			`"id":"f2id"`,
			`"id":"f3id"`,
			`"name":"test1"`,
			`"name":"test2"`,
			`"name":"test3"`,
		}

		for _, part := range expectedParts {
			if !strings.Contains(str, part) {
				t.Fatalf("Missing %q in\nn%v", part, str)
			}
		}
	})
}

func TestFieldsListScan(t *testing.T) {
	scenarios := []struct {
		name        string
		data        any
		expectError bool
		expectJSON  string
	}{
		{"nil", nil, false, "[]"},
		{"empty string", "", false, "[]"},
		{"empty byte", []byte{}, false, "[]"},
		{"empty string array", "[]", false, "[]"},
		{"invalid string", "invalid", true, "[]"},
		{"non-string", 123, true, "[]"},
		{"item with no field type", `[{}]`, true, "[]"},
		{
			"unknown field type",
			`[{"id":"123","name":"test1","type":"unknown"},{"id":"456","name":"test2","type":"bool"}]`,
			true,
			`[]`,
		},
		{
			"only the minimum field options",
			`[{"id":"123","name":"test1","type":"text","required":true},{"id":"456","name":"test2","type":"bool"}]`,
			false,
			`[{"autogeneratePattern":"","hidden":false,"id":"123","max":0,"min":0,"name":"test1","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":false,"type":"text"},{"hidden":false,"id":"456","name":"test2","presentable":false,"required":false,"system":false,"type":"bool"}]`,
		},
		{
			"all field options",
			`[{"autogeneratePattern":"","hidden":true,"id":"123","max":12,"min":0,"name":"test1","pattern":"","presentable":true,"primaryKey":false,"required":true,"system":false,"type":"text"},{"hidden":false,"id":"456","name":"test2","presentable":false,"required":false,"system":true,"type":"bool"}]`,
			false,
			`[{"autogeneratePattern":"","hidden":true,"id":"123","max":12,"min":0,"name":"test1","pattern":"","presentable":true,"primaryKey":false,"required":true,"system":false,"type":"text"},{"hidden":false,"id":"456","name":"test2","presentable":false,"required":false,"system":true,"type":"bool"}]`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testFieldsList := core.FieldsList{}

			err := testFieldsList.Scan(s.data)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			str := testFieldsList.String()
			if str != s.expectJSON {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expectJSON, str)
			}
		})
	}
}

func TestFieldsListJSON(t *testing.T) {
	scenarios := []struct {
		name        string
		data        string
		expectError bool
		expectJSON  string
	}{
		{"empty string", "", true, "[]"},
		{"invalid string", "invalid", true, "[]"},
		{"empty string array", "[]", false, "[]"},
		{"item with no field type", `[{}]`, true, "[]"},
		{
			"unknown field type",
			`[{"id":"123","name":"test1","type":"unknown"},{"id":"456","name":"test2","type":"bool"}]`,
			true,
			`[]`,
		},
		{
			"only the minimum field options",
			`[{"id":"123","name":"test1","type":"text","required":true},{"id":"456","name":"test2","type":"bool"}]`,
			false,
			`[{"autogeneratePattern":"","hidden":false,"id":"123","max":0,"min":0,"name":"test1","pattern":"","presentable":false,"primaryKey":false,"required":true,"system":false,"type":"text"},{"hidden":false,"id":"456","name":"test2","presentable":false,"required":false,"system":false,"type":"bool"}]`,
		},
		{
			"all field options",
			`[{"autogeneratePattern":"","hidden":true,"id":"123","max":12,"min":0,"name":"test1","pattern":"","presentable":true,"primaryKey":false,"required":true,"system":false,"type":"text"},{"hidden":false,"id":"456","name":"test2","presentable":false,"required":false,"system":true,"type":"bool"}]`,
			false,
			`[{"autogeneratePattern":"","hidden":true,"id":"123","max":12,"min":0,"name":"test1","pattern":"","presentable":true,"primaryKey":false,"required":true,"system":false,"type":"text"},{"hidden":false,"id":"456","name":"test2","presentable":false,"required":false,"system":true,"type":"bool"}]`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testFieldsList := core.FieldsList{}

			err := testFieldsList.UnmarshalJSON([]byte(s.data))

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			raw, err := testFieldsList.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			str := string(raw)
			if str != s.expectJSON {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expectJSON, str)
			}
		})
	}
}
