package core_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

func TestNewRecord(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")
	collection.Fields.Add(&core.BoolField{Name: "status"})

	m := core.NewRecord(collection)

	rawData, err := json.Marshal(m.FieldsData()) // should be initialized with the defaults
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"id":"","status":false}`

	if str := string(rawData); str != expected {
		t.Fatalf("Expected schema data\n%v\ngot\n%v", expected, str)
	}
}

func TestRecordCollection(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")

	m := core.NewRecord(collection)

	if m.Collection().Name != collection.Name {
		t.Fatalf("Expected collection with name %q, got %q", collection.Name, m.Collection().Name)
	}
}

func TestRecordTableName(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")

	m := core.NewRecord(collection)

	if m.TableName() != collection.Name {
		t.Fatalf("Expected table %q, got %q", collection.Name, m.TableName())
	}
}

func TestRecordPostScan(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test_collection")
	collection.Fields.Add(&core.TextField{Name: "test"})

	m := core.NewRecord(collection)

	// calling PostScan without id
	err := m.PostScan()
	if err == nil {
		t.Fatal("Expected PostScan id error, got nil")
	}

	m.Id = "test_id"
	m.Set("test", "abc")

	if v := m.IsNew(); v != true {
		t.Fatalf("[before PostScan] Expected IsNew %v, got %v", true, v)
	}
	if v := m.Original().PK(); v != "" {
		t.Fatalf("[before PostScan] Expected the original PK to be empty string, got %v", v)
	}
	if v := m.Original().Get("test"); v != "" {
		t.Fatalf("[before PostScan] Expected the original 'test' field to be empty string, got %v", v)
	}

	err = m.PostScan()
	if err != nil {
		t.Fatalf("Expected PostScan nil error, got %v", err)
	}

	if v := m.IsNew(); v != false {
		t.Fatalf("[after PostScan] Expected IsNew %v, got %v", false, v)
	}
	if v := m.Original().PK(); v != "test_id" {
		t.Fatalf("[after PostScan] Expected the original PK to be %q, got %v", "test_id", v)
	}
	if v := m.Original().Get("test"); v != "abc" {
		t.Fatalf("[after PostScan] Expected the original 'test' field to be %q, got %v", "abc", v)
	}
}

func TestRecordHookTags(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")

	m := core.NewRecord(collection)

	tags := m.HookTags()

	expectedTags := []string{collection.Id, collection.Name}

	if len(tags) != len(expectedTags) {
		t.Fatalf("Expected tags\n%v\ngot\n%v", expectedTags, tags)
	}

	for _, tag := range tags {
		if !slices.Contains(expectedTags, tag) {
			t.Errorf("Missing expected tag %q", tag)
		}
	}
}

func TestRecordBaseFilesPath(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")

	m := core.NewRecord(collection)
	m.Id = "abc"

	result := m.BaseFilesPath()
	expected := collection.BaseFilesPath() + "/" + m.Id
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestRecordOriginal(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	originalId := record.Id
	originalName := record.GetString("name")

	extraFieldsCheck := []string{`"email":`, `"custom":`}

	// change the fields
	record.Id = "changed"
	record.Set("name", "name_new")
	record.Set("custom", "test_custom")
	record.SetExpand(map[string]any{"test": 123})
	record.IgnoreEmailVisibility(true)
	record.IgnoreUnchangedFields(true)
	record.WithCustomData(true)
	record.Unhide(record.Collection().Fields.FieldNames()...)

	// ensure that the email visibility and the custom data toggles are active
	raw, err := record.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)
	for _, f := range extraFieldsCheck {
		if !strings.Contains(rawStr, f) {
			t.Fatalf("Expected %s in\n%s", f, rawStr)
		}
	}

	// check changes
	if v := record.GetString("name"); v != "name_new" {
		t.Fatalf("Expected name to be %q, got %q", "name_new", v)
	}
	if v := record.GetString("custom"); v != "test_custom" {
		t.Fatalf("Expected custom to be %q, got %q", "test_custom", v)
	}

	// check original
	if v := record.Original().PK(); v != originalId {
		t.Fatalf("Expected the original PK to be %q, got %q", originalId, v)
	}
	if v := record.Original().Id; v != originalId {
		t.Fatalf("Expected the original id to be %q, got %q", originalId, v)
	}
	if v := record.Original().GetString("name"); v != originalName {
		t.Fatalf("Expected the original name to be %q, got %q", originalName, v)
	}
	if v := record.Original().GetString("custom"); v != "" {
		t.Fatalf("Expected the original custom to be %q, got %q", "", v)
	}
	if v := record.Original().Expand(); len(v) != 0 {
		t.Fatalf("Expected empty original expand, got\n%v", v)
	}

	// ensure that the email visibility and the custom flag toggles weren't copied
	originalRaw, err := record.Original().MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	originalRawStr := string(originalRaw)
	for _, f := range extraFieldsCheck {
		if strings.Contains(originalRawStr, f) {
			t.Fatalf("Didn't expected %s in original\n%s", f, originalRawStr)
		}
	}

	// loading new data shouldn't affect the original state
	record.Load(map[string]any{"name": "name_new2"})

	if v := record.GetString("name"); v != "name_new2" {
		t.Fatalf("Expected name to be %q, got %q", "name_new2", v)
	}

	if v := record.Original().GetString("name"); v != originalName {
		t.Fatalf("Expected the original name still to be %q, got %q", originalName, v)
	}
}

func TestRecordFresh(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	originalId := record.Id

	extraFieldsCheck := []string{`"email":`, `"custom":`}

	autodateTest := types.NowDateTime()

	// change the fields
	record.Id = "changed"
	record.Set("name", "name_new")
	record.Set("custom", "test_custom")
	record.SetRaw("created", autodateTest)
	record.SetExpand(map[string]any{"test": 123})
	record.IgnoreEmailVisibility(true)
	record.IgnoreUnchangedFields(true)
	record.WithCustomData(true)
	record.Unhide(record.Collection().Fields.FieldNames()...)

	// ensure that the email visibility and the custom data toggles are active
	raw, err := record.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)
	for _, f := range extraFieldsCheck {
		if !strings.Contains(rawStr, f) {
			t.Fatalf("Expected %s in\n%s", f, rawStr)
		}
	}

	// check changes
	if v := record.GetString("name"); v != "name_new" {
		t.Fatalf("Expected name to be %q, got %q", "name_new", v)
	}
	if v := record.GetDateTime("created").String(); v != autodateTest.String() {
		t.Fatalf("Expected created to be %q, got %q", autodateTest.String(), v)
	}
	if v := record.GetString("custom"); v != "test_custom" {
		t.Fatalf("Expected custom to be %q, got %q", "test_custom", v)
	}

	// check fresh
	if v := record.Fresh().LastSavedPK(); v != originalId {
		t.Fatalf("Expected the fresh LastSavedPK to be %q, got %q", originalId, v)
	}
	if v := record.Fresh().PK(); v != record.Id {
		t.Fatalf("Expected the fresh PK to be %q, got %q", record.Id, v)
	}
	if v := record.Fresh().Id; v != record.Id {
		t.Fatalf("Expected the fresh id to be %q, got %q", record.Id, v)
	}
	if v := record.Fresh().GetString("name"); v != record.GetString("name") {
		t.Fatalf("Expected the fresh name to be %q, got %q", record.GetString("name"), v)
	}
	if v := record.Fresh().GetDateTime("created").String(); v != autodateTest.String() {
		t.Fatalf("Expected the fresh created to be %q, got %q", autodateTest.String(), v)
	}
	if v := record.Fresh().GetDateTime("updated").String(); v != record.GetDateTime("updated").String() {
		t.Fatalf("Expected the fresh updated to be %q, got %q", record.GetDateTime("updated").String(), v)
	}
	if v := record.Fresh().GetString("custom"); v != "" {
		t.Fatalf("Expected the fresh custom to be %q, got %q", "", v)
	}
	if v := record.Fresh().Expand(); len(v) != 0 {
		t.Fatalf("Expected empty fresh expand, got\n%v", v)
	}

	// ensure that the email visibility and the custom flag toggles weren't copied
	freshRaw, err := record.Fresh().MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	freshRawStr := string(freshRaw)
	for _, f := range extraFieldsCheck {
		if strings.Contains(freshRawStr, f) {
			t.Fatalf("Didn't expected %s in fresh\n%s", f, freshRawStr)
		}
	}
}

func TestRecordClone(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	originalId := record.Id

	extraFieldsCheck := []string{`"email":`, `"custom":`}

	autodateTest := types.NowDateTime()

	// change the fields
	record.Id = "changed"
	record.Set("name", "name_new")
	record.Set("custom", "test_custom")
	record.SetRaw("created", autodateTest)
	record.SetExpand(map[string]any{"test": 123})
	record.IgnoreEmailVisibility(true)
	record.WithCustomData(true)
	record.Unhide(record.Collection().Fields.FieldNames()...)

	// ensure that the email visibility and the custom data toggles are active
	raw, err := record.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)
	for _, f := range extraFieldsCheck {
		if !strings.Contains(rawStr, f) {
			t.Fatalf("Expected %s in\n%s", f, rawStr)
		}
	}

	// check changes
	if v := record.GetString("name"); v != "name_new" {
		t.Fatalf("Expected name to be %q, got %q", "name_new", v)
	}
	if v := record.GetDateTime("created").String(); v != autodateTest.String() {
		t.Fatalf("Expected created to be %q, got %q", autodateTest.String(), v)
	}
	if v := record.GetString("custom"); v != "test_custom" {
		t.Fatalf("Expected custom to be %q, got %q", "test_custom", v)
	}

	// check clone
	if v := record.Clone().LastSavedPK(); v != originalId {
		t.Fatalf("Expected the clone LastSavedPK to be %q, got %q", originalId, v)
	}
	if v := record.Clone().PK(); v != record.Id {
		t.Fatalf("Expected the clone PK to be %q, got %q", record.Id, v)
	}
	if v := record.Clone().Id; v != record.Id {
		t.Fatalf("Expected the clone id to be %q, got %q", record.Id, v)
	}
	if v := record.Clone().GetString("name"); v != record.GetString("name") {
		t.Fatalf("Expected the clone name to be %q, got %q", record.GetString("name"), v)
	}
	if v := record.Clone().GetDateTime("created").String(); v != autodateTest.String() {
		t.Fatalf("Expected the clone created to be %q, got %q", autodateTest.String(), v)
	}
	if v := record.Clone().GetDateTime("updated").String(); v != record.GetDateTime("updated").String() {
		t.Fatalf("Expected the clone updated to be %q, got %q", record.GetDateTime("updated").String(), v)
	}
	if v := record.Clone().GetString("custom"); v != "test_custom" {
		t.Fatalf("Expected the clone custom to be %q, got %q", "test_custom", v)
	}
	if _, ok := record.Clone().Expand()["test"]; !ok {
		t.Fatalf("Expected non-empty clone expand")
	}

	// ensure that the email visibility and the custom data toggles state were copied
	cloneRaw, err := record.Clone().MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	cloneRawStr := string(cloneRaw)
	for _, f := range extraFieldsCheck {
		if !strings.Contains(cloneRawStr, f) {
			t.Fatalf("Expected %s in clone\n%s", f, cloneRawStr)
		}
	}
}

func TestRecordExpand(t *testing.T) {
	t.Parallel()

	record := core.NewRecord(core.NewBaseCollection("test"))

	expand := record.Expand()
	if expand == nil || len(expand) != 0 {
		t.Fatalf("Expected empty map expand, got %v", expand)
	}

	data1 := map[string]any{"a": 123, "b": 456}
	data2 := map[string]any{"c": 123}
	record.SetExpand(data1)
	record.SetExpand(data2) // should overwrite the previous call

	// modify the expand map to check for shallow copy
	data2["d"] = 456

	expand = record.Expand()
	if len(expand) != 1 {
		t.Fatalf("Expected empty map expand, got %v", expand)
	}
	if v := expand["c"]; v != 123 {
		t.Fatalf("Expected to find expand.c %v, got %v", 123, v)
	}
}

func TestRecordMergeExpand(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")
	collection.Id = "_pbc_123"

	m := core.NewRecord(collection)
	m.Id = "m"

	// a
	a := core.NewRecord(collection)
	a.Id = "a"
	a1 := core.NewRecord(collection)
	a1.Id = "a1"
	a2 := core.NewRecord(collection)
	a2.Id = "a2"
	a3 := core.NewRecord(collection)
	a3.Id = "a3"
	a31 := core.NewRecord(collection)
	a31.Id = "a31"
	a32 := core.NewRecord(collection)
	a32.Id = "a32"
	a.SetExpand(map[string]any{
		"a1":  a1,
		"a23": []*core.Record{a2, a3},
	})
	a3.SetExpand(map[string]any{
		"a31": a31,
		"a32": []*core.Record{a32},
	})

	// b
	b := core.NewRecord(collection)
	b.Id = "b"
	b1 := core.NewRecord(collection)
	b1.Id = "b1"
	b.SetExpand(map[string]any{
		"b1": b1,
	})

	// c
	c := core.NewRecord(collection)
	c.Id = "c"

	// load initial expand
	m.SetExpand(map[string]any{
		"a": a,
		"b": b,
		"c": []*core.Record{c},
	})

	// a (new)
	aNew := core.NewRecord(collection)
	aNew.Id = a.Id
	a3New := core.NewRecord(collection)
	a3New.Id = a3.Id
	a32New := core.NewRecord(collection)
	a32New.Id = "a32New"
	a33New := core.NewRecord(collection)
	a33New.Id = "a33New"
	a3New.SetExpand(map[string]any{
		"a32":    []*core.Record{a32New},
		"a33New": a33New,
	})
	aNew.SetExpand(map[string]any{
		"a23": []*core.Record{a2, a3New},
	})

	// b (new)
	bNew := core.NewRecord(collection)
	bNew.Id = "bNew"
	dNew := core.NewRecord(collection)
	dNew.Id = "dNew"

	// merge expands
	m.MergeExpand(map[string]any{
		"a":    aNew,
		"b":    []*core.Record{bNew},
		"dNew": dNew,
	})

	result := m.Expand()

	raw, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)

	expected := `{"a":{"collectionId":"_pbc_123","collectionName":"test","expand":{"a1":{"collectionId":"_pbc_123","collectionName":"test","id":"a1"},"a23":[{"collectionId":"_pbc_123","collectionName":"test","id":"a2"},{"collectionId":"_pbc_123","collectionName":"test","expand":{"a31":{"collectionId":"_pbc_123","collectionName":"test","id":"a31"},"a32":[{"collectionId":"_pbc_123","collectionName":"test","id":"a32"},{"collectionId":"_pbc_123","collectionName":"test","id":"a32New"}],"a33New":{"collectionId":"_pbc_123","collectionName":"test","id":"a33New"}},"id":"a3"}]},"id":"a"},"b":[{"collectionId":"_pbc_123","collectionName":"test","expand":{"b1":{"collectionId":"_pbc_123","collectionName":"test","id":"b1"}},"id":"b"},{"collectionId":"_pbc_123","collectionName":"test","id":"bNew"}],"c":[{"collectionId":"_pbc_123","collectionName":"test","id":"c"}],"dNew":{"collectionId":"_pbc_123","collectionName":"test","id":"dNew"}}`

	if expected != rawStr {
		t.Fatalf("Expected \n%v, \ngot \n%v", expected, rawStr)
	}
}

func TestRecordMergeExpandNilCheck(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")
	collection.Id = "_pbc_123"

	scenarios := []struct {
		name     string
		expand   map[string]any
		expected string
	}{
		{
			"nil expand",
			nil,
			`{"collectionId":"_pbc_123","collectionName":"test","id":""}`,
		},
		{
			"empty expand",
			map[string]any{},
			`{"collectionId":"_pbc_123","collectionName":"test","id":""}`,
		},
		{
			"non-empty expand",
			map[string]any{"test": core.NewRecord(collection)},
			`{"collectionId":"_pbc_123","collectionName":"test","expand":{"test":{"collectionId":"_pbc_123","collectionName":"test","id":""}},"id":""}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			m := core.NewRecord(collection)
			m.MergeExpand(s.expand)

			raw, err := json.Marshal(m)
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			if rawStr != s.expected {
				t.Fatalf("Expected \n%v, \ngot \n%v", s.expected, rawStr)
			}
		})
	}
}

func TestRecordExpandedOne(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")

	main := core.NewRecord(collection)

	single := core.NewRecord(collection)
	single.Id = "single"

	multiple1 := core.NewRecord(collection)
	multiple1.Id = "multiple1"

	multiple2 := core.NewRecord(collection)
	multiple2.Id = "multiple2"

	main.SetExpand(map[string]any{
		"single":   single,
		"multiple": []*core.Record{multiple1, multiple2},
	})

	if v := main.ExpandedOne("missing"); v != nil {
		t.Fatalf("Expected nil, got %v", v)
	}

	if v := main.ExpandedOne("single"); v == nil || v.Id != "single" {
		t.Fatalf("Expected record with id %q, got %v", "single", v)
	}

	if v := main.ExpandedOne("multiple"); v == nil || v.Id != "multiple1" {
		t.Fatalf("Expected record with id %q, got %v", "multiple1", v)
	}
}

func TestRecordExpandedAll(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")

	main := core.NewRecord(collection)

	single := core.NewRecord(collection)
	single.Id = "single"

	multiple1 := core.NewRecord(collection)
	multiple1.Id = "multiple1"

	multiple2 := core.NewRecord(collection)
	multiple2.Id = "multiple2"

	main.SetExpand(map[string]any{
		"single":   single,
		"multiple": []*core.Record{multiple1, multiple2},
	})

	if v := main.ExpandedAll("missing"); v != nil {
		t.Fatalf("Expected nil, got %v", v)
	}

	if v := main.ExpandedAll("single"); len(v) != 1 || v[0].Id != "single" {
		t.Fatalf("Expected [single] slice, got %v", v)
	}

	if v := main.ExpandedAll("multiple"); len(v) != 2 || v[0].Id != "multiple1" || v[1].Id != "multiple2" {
		t.Fatalf("Expected [multiple1, multiple2] slice, got %v", v)
	}
}

func TestRecordFieldsData(t *testing.T) {
	t.Parallel()

	collection := core.NewAuthCollection("test")
	collection.Fields.Add(&core.TextField{Name: "field1"})
	collection.Fields.Add(&core.TextField{Name: "field2"})

	m := core.NewRecord(collection)
	m.Id = "test_id" // direct id assignment
	m.Set("email", "test@example.com")
	m.Set("password", "123") // hidden fields should be also returned
	m.Set("tokenKey", "789")
	m.Set("field1", 123)
	m.Set("field2", 456)
	m.Set("unknown", 789)

	raw, err := json.Marshal(m.FieldsData())
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"email":"test@example.com","emailVisibility":false,"field1":"123","field2":"456","id":"test_id","password":"123","tokenKey":"789","verified":false}`

	if v := string(raw); v != expected {
		t.Fatalf("Expected\n%v\ngot\n%v", expected, v)
	}
}

func TestRecordCustomData(t *testing.T) {
	t.Parallel()

	collection := core.NewAuthCollection("test")
	collection.Fields.Add(&core.TextField{Name: "field1"})
	collection.Fields.Add(&core.TextField{Name: "field2"})

	m := core.NewRecord(collection)
	m.Id = "test_id" // direct id assignment
	m.Set("email", "test@example.com")
	m.Set("password", "123") // hidden fields should be also returned
	m.Set("tokenKey", "789")
	m.Set("field1", 123)
	m.Set("field2", 456)
	m.Set("unknown", 789)

	raw, err := json.Marshal(m.CustomData())
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"unknown":789}`

	if v := string(raw); v != expected {
		t.Fatalf("Expected\n%v\ngot\n%v", expected, v)
	}
}

func TestRecordSetGet(t *testing.T) {
	t.Parallel()

	f1 := &mockField{}
	f1.Name = "mock1"

	f2 := &mockField{}
	f2.Name = "mock2"

	f3 := &mockField{}
	f3.Name = "mock3"

	collection := core.NewBaseCollection("test")
	collection.Fields.Add(&core.TextField{Name: "text1"})
	collection.Fields.Add(&core.TextField{Name: "text2"})
	collection.Fields.Add(f1)
	collection.Fields.Add(f2)
	collection.Fields.Add(f3)

	record := core.NewRecord(collection)
	record.Set("text1", 123) // should be converted to string using the ScanValue fallback
	record.SetRaw("text2", 456)
	record.Set("mock1", 1) // should be converted to string using the setter
	record.SetRaw("mock2", 1)
	record.Set("mock3:test", "abc")
	record.Set("unknown", 789)

	t.Run("GetRaw", func(t *testing.T) {
		expected := map[string]any{
			"text1":      "123",
			"text2":      456,
			"mock1":      "1",
			"mock2":      1,
			"mock3":      "modifier_set",
			"mock3:test": nil,
			"unknown":    789,
		}

		for k, v := range expected {
			raw := record.GetRaw(k)
			if raw != v {
				t.Errorf("Expected %q to be %v, got %v", k, v, raw)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		expected := map[string]any{
			"text1":      "123",
			"text2":      456,
			"mock1":      "1",
			"mock2":      1,
			"mock3":      "modifier_set",
			"mock3:test": "modifier_get",
			"unknown":    789,
		}

		for k, v := range expected {
			get := record.Get(k)
			if get != v {
				t.Errorf("Expected %q to be %v, got %v", k, v, get)
			}
		}
	})
}

func TestRecordLoad(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")
	collection.Fields.Add(&core.TextField{Name: "text"})

	record := core.NewRecord(collection)
	record.Load(map[string]any{
		"text":   123,
		"custom": 456,
	})

	expected := map[string]any{
		"text":   "123",
		"custom": 456,
	}

	for k, v := range expected {
		get := record.Get(k)
		if get != v {
			t.Errorf("Expected %q to be %#v, got %#v", k, v, get)
		}
	}
}

func TestRecordGetBool(t *testing.T) {
	t.Parallel()

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

	collection := core.NewBaseCollection("test")
	record := core.NewRecord(collection)

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			record.Set("test", s.value)

			result := record.GetBool("test")
			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestRecordGetString(t *testing.T) {
	t.Parallel()

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

	collection := core.NewBaseCollection("test")
	record := core.NewRecord(collection)

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			record.Set("test", s.value)

			result := record.GetString("test")
			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}

func TestRecordGetInt(t *testing.T) {
	t.Parallel()

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

	collection := core.NewBaseCollection("test")
	record := core.NewRecord(collection)

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			record.Set("test", s.value)

			result := record.GetInt("test")
			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestRecordGetFloat(t *testing.T) {
	t.Parallel()

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

	collection := core.NewBaseCollection("test")
	record := core.NewRecord(collection)

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			record.Set("test", s.value)

			result := record.GetFloat("test")
			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestRecordGetDateTime(t *testing.T) {
	t.Parallel()

	nowTime := time.Now()
	testTime, _ := time.Parse(types.DefaultDateLayout, "2022-01-01 08:00:40.000Z")

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

	collection := core.NewBaseCollection("test")
	record := core.NewRecord(collection)

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			record.Set("test", s.value)

			result := record.GetDateTime("test")
			if !result.Time().Equal(s.expected) {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestRecordGetStringSlice(t *testing.T) {
	t.Parallel()

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

	collection := core.NewBaseCollection("test")
	record := core.NewRecord(collection)

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			record.Set("test", s.value)

			result := record.GetStringSlice("test")

			if len(result) != len(s.expected) {
				t.Fatalf("Expected %d elements, got %d: %v", len(s.expected), len(result), result)
			}

			for _, v := range result {
				if !slices.Contains(s.expected, v) {
					t.Fatalf("Cannot find %v in %v", v, s.expected)
				}
			}
		})
	}
}

func TestRecordGetGeoPoint(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		value    any
		expected string
	}{
		{nil, `{"lon":0,"lat":0}`},
		{"", `{"lon":0,"lat":0}`},
		{0, `{"lon":0,"lat":0}`},
		{false, `{"lon":0,"lat":0}`},
		{"{}", `{"lon":0,"lat":0}`},
		{"[]", `{"lon":0,"lat":0}`},
		{[]int{1, 2}, `{"lon":0,"lat":0}`},
		{map[string]any{"lon": 1, "lat": 2}, `{"lon":1,"lat":2}`},
		{[]byte(`{"lon":1,"lat":2}`), `{"lon":1,"lat":2}`},
		{`{"lon":1,"lat":2}`, `{"lon":1,"lat":2}`},
		{types.GeoPoint{Lon: 1, Lat: 2}, `{"lon":1,"lat":2}`},
		{&types.GeoPoint{Lon: 1, Lat: 2}, `{"lon":1,"lat":2}`},
	}

	collection := core.NewBaseCollection("test")
	record := core.NewRecord(collection)

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			record.Set("test", s.value)

			pointStr := record.GetGeoPoint("test").String()

			if pointStr != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, pointStr)
			}
		})
	}
}

func TestRecordGetUnsavedFiles(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f1, err := filesystem.NewFileFromBytes([]byte("test"), "f1")
	if err != nil {
		t.Fatal(err)
	}
	f1.Name = "f1"

	f2, err := filesystem.NewFileFromBytes([]byte("test"), "f2")
	if err != nil {
		t.Fatal(err)
	}
	f2.Name = "f2"

	record, err := app.FindRecordById("demo3", "lcl9d87w22ml6jy")
	if err != nil {
		t.Fatal(err)
	}
	record.Set("files+", []any{f1, f2})

	scenarios := []struct {
		key      string
		expected string
	}{
		{
			"",
			"null",
		},
		{
			"title",
			"null",
		},
		{
			"files",
			`[{"name":"f1","originalName":"f1","size":4},{"name":"f2","originalName":"f2","size":4}]`,
		},
		{
			"files:unsaved",
			`[{"name":"f1","originalName":"f1","size":4},{"name":"f2","originalName":"f2","size":4}]`,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.key), func(t *testing.T) {
			v := record.GetUnsavedFiles(s.key)

			raw, err := json.Marshal(v)
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			if rawStr != s.expected {
				t.Fatalf("Expected\n%s\ngot\n%s", s.expected, rawStr)
			}
		})
	}
}

func TestRecordUnmarshalJSONField(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")
	collection.Fields.Add(&core.JSONField{Name: "field"})

	record := core.NewRecord(collection)

	var testPointer *string
	var testStr string
	var testInt int
	var testBool bool
	var testSlice []int
	var testMap map[string]any

	scenarios := []struct {
		value        any
		destination  any
		expectError  bool
		expectedJSON string
	}{
		{nil, testPointer, false, `null`},
		{nil, testStr, false, `""`},
		{"", testStr, false, `""`},
		{1, testInt, false, `1`},
		{true, testBool, false, `true`},
		{[]int{1, 2, 3}, testSlice, false, `[1,2,3]`},
		{map[string]any{"test": 123}, testMap, false, `{"test":123}`},
		// json encoded values
		{`null`, testPointer, false, `null`},
		{`true`, testBool, false, `true`},
		{`456`, testInt, false, `456`},
		{`"test"`, testStr, false, `"test"`},
		{`[4,5,6]`, testSlice, false, `[4,5,6]`},
		{`{"test":456}`, testMap, false, `{"test":456}`},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%#v", i, s.value), func(t *testing.T) {
			record.Set("field", s.value)

			err := record.UnmarshalJSONField("field", &s.destination)
			hasErr := err != nil

			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectError, hasErr)
			}

			raw, _ := json.Marshal(s.destination)
			if v := string(raw); v != s.expectedJSON {
				t.Fatalf("Expected %q, got %q", s.expectedJSON, v)
			}
		})
	}
}

func TestRecordFindFileFieldByFile(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")
	collection.Fields.Add(
		&core.TextField{Name: "field1"},
		&core.FileField{Name: "field2", MaxSelect: 1, MaxSize: 1},
		&core.FileField{Name: "field3", MaxSelect: 2, MaxSize: 1},
	)

	m := core.NewRecord(collection)
	m.Set("field1", "test")
	m.Set("field2", "test.png")
	m.Set("field3", []string{"test1.png", "test2.png"})

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
		t.Run(fmt.Sprintf("%d_%#v", i, s.filename), func(t *testing.T) {
			result := m.FindFileFieldByFile(s.filename)

			var fieldName string
			if result != nil {
				fieldName = result.Name
			}

			if s.expectField != fieldName {
				t.Fatalf("Expected field %v, got %v", s.expectField, result)
			}
		})
	}
}

func TestRecordDBExport(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	f1 := &core.TextField{Name: "field1"}
	f2 := &core.FileField{Name: "field2", MaxSelect: 1, MaxSize: 1}
	f3 := &core.SelectField{Name: "field3", MaxSelect: 2, Values: []string{"test1", "test2", "test3"}}
	f4 := &core.RelationField{Name: "field4", MaxSelect: 2}

	colBase := core.NewBaseCollection("test_base")
	colBase.Fields.Add(f1, f2, f3, f4)

	colAuth := core.NewAuthCollection("test_auth")
	colAuth.Fields.Add(f1, f2, f3, f4)

	scenarios := []struct {
		collection *core.Collection
		expected   string
	}{
		{
			colBase,
			`{"field1":"test","field2":"test.png","field3":["test1","test2"],"field4":["test11","test12"],"id":"test_id"}`,
		},
		{
			colAuth,
			`{"email":"test_email","emailVisibility":true,"field1":"test","field2":"test.png","field3":["test1","test2"],"field4":["test11","test12"],"id":"test_id","password":"_TEST_","tokenKey":"test_tokenKey","verified":false}`,
		},
	}

	data := map[string]any{
		"id":              "test_id",
		"field1":          "test",
		"field2":          "test.png",
		"field3":          []string{"test1", "test2"},
		"field4":          []string{"test11", "test12", "test11"}, // strip duplicate,
		"unknown":         "test_unknown",
		"password":        "test_passwordHash",
		"username":        "test_username",
		"emailVisibility": true,
		"email":           "test_email",
		"verified":        "invalid", // should be casted
		"tokenKey":        "test_tokenKey",
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s_%s", i, s.collection.Type, s.collection.Name), func(t *testing.T) {
			record := core.NewRecord(s.collection)

			record.Load(data)

			result, err := record.DBExport(app)
			if err != nil {
				t.Fatal(err)
			}

			raw, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			// replace _TEST_ placeholder with .+ regex pattern
			pattern := regexp.MustCompile(strings.ReplaceAll(
				"^"+regexp.QuoteMeta(s.expected)+"$",
				"_TEST_",
				`.+`,
			))

			if !pattern.MatchString(rawStr) {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, rawStr)
			}
		})
	}
}

func TestRecordIgnoreUnchangedFields(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	col, err := app.FindCollectionByNameOrId("demo3")
	if err != nil {
		t.Fatal(err)
	}

	new := core.NewRecord(col)

	existing, err := app.FindRecordById(col, "mk5fmymtx4wsprk")
	if err != nil {
		t.Fatal(err)
	}
	existing.Set("title", "test_new")
	existing.Set("files", existing.Get("files")) // no change

	scenarios := []struct {
		ignoreUnchangedFields bool
		record                *core.Record
		expected              []string
	}{
		{
			false,
			new,
			[]string{"id", "created", "updated", "title", "files"},
		},
		{
			true,
			new,
			[]string{"id", "created", "updated", "title", "files"},
		},
		{
			false,
			existing,
			[]string{"id", "created", "updated", "title", "files"},
		},
		{
			true,
			existing,
			[]string{"id", "title"},
		},
	}

	for i, s := range scenarios {
		action := "create"
		if !s.record.IsNew() {
			action = "update"
		}

		t.Run(fmt.Sprintf("%d_%s_%v", i, action, s.ignoreUnchangedFields), func(t *testing.T) {
			s.record.IgnoreUnchangedFields(s.ignoreUnchangedFields)

			result, err := s.record.DBExport(app)
			if err != nil {
				t.Fatal(err)
			}

			if len(result) != len(s.expected) {
				t.Fatalf("Expected %d keys, got %d:\n%v", len(s.expected), len(result), result)
			}

			for _, key := range s.expected {
				if _, ok := result[key]; !ok {
					t.Fatalf("Missing expected key %q in\n%v", key, result)
				}
			}
		})
	}
}

func TestRecordPublicExportAndMarshalJSON(t *testing.T) {
	t.Parallel()

	f1 := &core.TextField{Name: "field1"}
	f2 := &core.FileField{Name: "field2", MaxSelect: 1, MaxSize: 1}
	f3 := &core.SelectField{Name: "field3", MaxSelect: 2, Values: []string{"test1", "test2", "test3"}}
	f4 := &core.TextField{Name: "field4", Hidden: true}
	f5 := &core.TextField{Name: "field5", Hidden: true}

	colBase := core.NewBaseCollection("test_base")
	colBase.Id = "_pbc_base_123"
	colBase.Fields.Add(f1, f2, f3, f4, f5)

	colAuth := core.NewAuthCollection("test_auth")
	colAuth.Id = "_pbc_auth_123"
	colAuth.Fields.Add(f1, f2, f3, f4, f5)

	scenarios := []struct {
		name                  string
		collection            *core.Collection
		ignoreEmailVisibility bool
		withCustomData        bool
		hideFields            []string
		unhideFields          []string
		expectedJSON          string
	}{
		// base
		{
			"[base] no extra flags",
			colBase,
			false,
			false,
			nil,
			nil,
			`{"collectionId":"_pbc_base_123","collectionName":"test_base","expand":{"test":123},"field1":"field_1","field2":"field_2.png","field3":["test1","test2"],"id":"test_id"}`,
		},
		{
			"[base] with email visibility",
			colBase,
			true, // should have no effect
			false,
			nil,
			nil,
			`{"collectionId":"_pbc_base_123","collectionName":"test_base","expand":{"test":123},"field1":"field_1","field2":"field_2.png","field3":["test1","test2"],"id":"test_id"}`,
		},
		{
			"[base] with custom data",
			colBase,
			true, // should have no effect
			true,
			nil,
			nil,
			`{"collectionId":"_pbc_base_123","collectionName":"test_base","email":"test_email","emailVisibility":"test_invalid","expand":{"test":123},"field1":"field_1","field2":"field_2.png","field3":["test1","test2"],"id":"test_id","password":"test_passwordHash","tokenKey":"test_tokenKey","unknown":"test_unknown","verified":true}`,
		},
		{
			"[base] with explicit hide and unhide fields",
			colBase,
			false,
			true,
			[]string{"field3", "field1", "expand", "collectionId", "collectionName", "email", "tokenKey", "unknown"},
			[]string{"field4", "@pbInternalAbc"},
			`{"emailVisibility":"test_invalid","field2":"field_2.png","field4":"field_4","id":"test_id","password":"test_passwordHash","verified":true}`,
		},
		{
			"[base] trying to unhide custom fields without explicit WithCustomData",
			colBase,
			false,
			true,
			nil,
			[]string{"field5", "@pbInternalAbc", "email", "tokenKey", "unknown"},
			`{"collectionId":"_pbc_base_123","collectionName":"test_base","email":"test_email","emailVisibility":"test_invalid","expand":{"test":123},"field1":"field_1","field2":"field_2.png","field3":["test1","test2"],"field5":"field_5","id":"test_id","password":"test_passwordHash","tokenKey":"test_tokenKey","unknown":"test_unknown","verified":true}`,
		},

		// auth
		{
			"[auth] no extra flags",
			colAuth,
			false,
			false,
			nil,
			nil,
			`{"collectionId":"_pbc_auth_123","collectionName":"test_auth","emailVisibility":false,"expand":{"test":123},"field1":"field_1","field2":"field_2.png","field3":["test1","test2"],"id":"test_id","verified":true}`,
		},
		{
			"[auth] with email visibility",
			colAuth,
			true,
			false,
			nil,
			nil,
			`{"collectionId":"_pbc_auth_123","collectionName":"test_auth","email":"test_email","emailVisibility":false,"expand":{"test":123},"field1":"field_1","field2":"field_2.png","field3":["test1","test2"],"id":"test_id","verified":true}`,
		},
		{
			"[auth] with custom data",
			colAuth,
			false,
			true,
			nil,
			nil,
			`{"collectionId":"_pbc_auth_123","collectionName":"test_auth","emailVisibility":false,"expand":{"test":123},"field1":"field_1","field2":"field_2.png","field3":["test1","test2"],"id":"test_id","unknown":"test_unknown","verified":true}`,
		},
		{
			"[auth] with explicit hide and unhide fields",
			colAuth,
			true,
			true,
			[]string{"field3", "field1", "expand", "collectionId", "collectionName", "email", "unknown"},
			[]string{"field4", "@pbInternalAbc"},
			`{"emailVisibility":false,"field2":"field_2.png","field4":"field_4","id":"test_id","verified":true}`,
		},
		{
			"[auth] trying to unhide custom fields without explicit WithCustomData",
			colAuth,
			false,
			true,
			nil,
			[]string{"field5", "@pbInternalAbc", "tokenKey", "unknown", "email"}, // emailVisibility:false has higher priority
			`{"collectionId":"_pbc_auth_123","collectionName":"test_auth","emailVisibility":false,"expand":{"test":123},"field1":"field_1","field2":"field_2.png","field3":["test1","test2"],"field5":"field_5","id":"test_id","unknown":"test_unknown","verified":true}`,
		},
	}

	data := map[string]any{
		"id":              "test_id",
		"field1":          "field_1",
		"field2":          "field_2.png",
		"field3":          []string{"test1", "test2"},
		"field4":          "field_4",
		"field5":          "field_5",
		"expand":          map[string]any{"test": 123},
		"collectionId":    "m_id",   // should be always ignored
		"collectionName":  "m_name", // should be always ignored
		"unknown":         "test_unknown",
		"password":        "test_passwordHash",
		"emailVisibility": "test_invalid", // for auth collections should be casted to bool
		"email":           "test_email",
		"verified":        true,
		"tokenKey":        "test_tokenKey",
		"@pbInternalAbc":  "test_custom_inter", // always hidden
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			m := core.NewRecord(s.collection)

			m.Load(data)
			m.IgnoreEmailVisibility(s.ignoreEmailVisibility)
			m.WithCustomData(s.withCustomData)
			m.Unhide(s.unhideFields...)
			m.Hide(s.hideFields...)

			exportResult, err := json.Marshal(m.PublicExport())
			if err != nil {
				t.Fatal(err)
			}
			exportResultStr := string(exportResult)

			// MarshalJSON and PublicExport should return the same
			marshalResult, err := m.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			marshalResultStr := string(marshalResult)

			if exportResultStr != marshalResultStr {
				t.Fatalf("Expected the PublicExport to be the same as MarshalJSON, but got \n%v \nvs \n%v", exportResultStr, marshalResultStr)
			}

			if exportResultStr != s.expectedJSON {
				t.Fatalf("Expected json \n%v \ngot \n%v", s.expectedJSON, exportResultStr)
			}
		})
	}
}

func TestRecordUnmarshalJSON(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")
	collection.Fields.Add(&core.TextField{Name: "text"})

	record := core.NewRecord(collection)

	data := map[string]any{
		"text":   123,
		"custom": 456.789,
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	err = record.UnmarshalJSON(rawData)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	expected := map[string]any{
		"text":   "123",
		"custom": 456.789,
	}

	for k, v := range expected {
		get := record.Get(k)
		if get != v {
			t.Errorf("Expected %q to be %#v, got %#v", k, v, get)
		}
	}
}

func TestRecordReplaceModifiers(t *testing.T) {
	t.Parallel()

	collection := core.NewBaseCollection("test")
	collection.Fields.Add(
		&mockField{core.TextField{Name: "mock"}},
		&core.NumberField{Name: "number"},
	)

	originalData := map[string]any{
		"mock":   "a",
		"number": 2.1,
	}

	record := core.NewRecord(collection)
	for k, v := range originalData {
		record.Set(k, v)
	}

	result := record.ReplaceModifiers(map[string]any{
		"mock:test": "b",
		"number+":   3,
	})

	expected := map[string]any{
		"mock":   "modifier_set",
		"number": 5.1,
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected\n%v\ngot\n%v", expected, result)
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Expected %q %#v, got %#v", k, v, result[k])
		}
	}

	// ensure that the original data hasn't changed
	for k, v := range originalData {
		rv := record.Get(k)
		if rv != v {
			t.Errorf("Expected original %q %#v, got %#v", k, v, rv)
		}
	}
}

func TestRecordValidate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// dummy collection to ensure that the specified field validators are triggered
	collection := core.NewBaseCollection("validate_test")
	collection.Fields.Add(
		&core.TextField{Name: "f1", Min: 3},
		&core.NumberField{Name: "f2", Required: true},
	)
	if err := app.Save(collection); err != nil {
		t.Fatal(err)
	}

	record := core.NewRecord(collection)
	record.Id = "!invalid"

	t.Run("no data set", func(t *testing.T) {
		tests.TestValidationErrors(t, app.Validate(record), []string{"id", "f2"})
	})

	t.Run("failing the text field min requirement", func(t *testing.T) {
		record.Set("f1", "a")
		tests.TestValidationErrors(t, app.Validate(record), []string{"id", "f1", "f2"})
	})

	t.Run("satisfying the fields validations", func(t *testing.T) {
		record.Id = strings.Repeat("b", 15)
		record.Set("f1", "abc")
		record.Set("f2", 1)
		tests.TestValidationErrors(t, app.Validate(record), nil)
	})
}

func TestRecordModelEventSync(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	col, err := app.FindCollectionByNameOrId("demo3")
	if err != nil {
		t.Fatal(err)
	}

	testRecords := make([]*core.Record, 4)
	for i := 0; i < 4; i++ {
		testRecords[i] = core.NewRecord(col)
		testRecords[i].Set("title", "sync_test_"+strconv.Itoa(i))
		if err := app.Save(testRecords[i]); err != nil {
			t.Fatal(err)
		}
	}

	createModelEvent := func() *core.ModelEvent {
		event := new(core.ModelEvent)
		event.App = app
		event.Context = context.Background()
		event.Type = "test_a"
		event.Model = testRecords[0]
		return event
	}

	createModelErrorEvent := func() *core.ModelErrorEvent {
		event := new(core.ModelErrorEvent)
		event.ModelEvent = *createModelEvent()
		event.Error = errors.New("error_a")
		return event
	}

	changeRecordEventBefore := func(e *core.RecordEvent) {
		e.Type = "test_b"
		//nolint:staticcheck
		e.Context = context.WithValue(context.Background(), "test", 123)
		e.Record = testRecords[1]
	}

	modelEventFinalizerChange := func(e *core.ModelEvent) {
		e.Type = "test_c"
		//nolint:staticcheck
		e.Context = context.WithValue(context.Background(), "test", 456)
		e.Model = testRecords[2]
	}

	changeRecordEventAfter := func(e *core.RecordEvent) {
		e.Type = "test_d"
		//nolint:staticcheck
		e.Context = context.WithValue(context.Background(), "test", 789)
		e.Record = testRecords[3]
	}

	expectedBeforeModelEventHandlerChecks := func(t *testing.T, e *core.ModelEvent) {
		if e.Type != "test_a" {
			t.Fatalf("Expected type %q, got %q", "test_a", e.Type)
		}

		if v := e.Context.Value("test"); v != nil {
			t.Fatalf("Expected context value %v, got %v", nil, v)
		}

		if e.Model.PK() != testRecords[0].Id {
			t.Fatalf("Expected record with id %q, got %q (%d)", testRecords[0].Id, e.Model.PK(), 0)
		}
	}

	expectedAfterModelEventHandlerChecks := func(t *testing.T, e *core.ModelEvent) {
		if e.Type != "test_d" {
			t.Fatalf("Expected type %q, got %q", "test_d", e.Type)
		}

		if v := e.Context.Value("test"); v != 789 {
			t.Fatalf("Expected context value %v, got %v", 789, v)
		}

		// note: currently the Model and Record values are not synced due to performance consideration
		if e.Model.PK() != testRecords[2].Id {
			t.Fatalf("Expected record with id %q, got %q (%d)", testRecords[2].Id, e.Model.PK(), 2)
		}
	}

	expectedBeforeRecordEventHandlerChecks := func(t *testing.T, e *core.RecordEvent) {
		if e.Type != "test_a" {
			t.Fatalf("Expected type %q, got %q", "test_a", e.Type)
		}

		if v := e.Context.Value("test"); v != nil {
			t.Fatalf("Expected context value %v, got %v", nil, v)
		}

		if e.Record.Id != testRecords[0].Id {
			t.Fatalf("Expected record with id %q, got %q (%d)", testRecords[0].Id, e.Record.Id, 2)
		}
	}

	expectedAfterRecordEventHandlerChecks := func(t *testing.T, e *core.RecordEvent) {
		if e.Type != "test_c" {
			t.Fatalf("Expected type %q, got %q", "test_c", e.Type)
		}

		if v := e.Context.Value("test"); v != 456 {
			t.Fatalf("Expected context value %v, got %v", 456, v)
		}

		// note: currently the Model and Record values are not synced due to performance consideration
		if e.Record.Id != testRecords[1].Id {
			t.Fatalf("Expected record with id %q, got %q (%d)", testRecords[1].Id, e.Record.Id, 1)
		}
	}

	modelEventFinalizer := func(e *core.ModelEvent) error {
		modelEventFinalizerChange(e)
		return nil
	}

	modelErrorEventFinalizer := func(e *core.ModelErrorEvent) error {
		modelEventFinalizerChange(&e.ModelEvent)
		e.Error = errors.New("error_c")
		return nil
	}

	modelEventHandler := &hook.Handler[*core.ModelEvent]{
		Priority: -999,
		Func: func(e *core.ModelEvent) error {
			t.Run("before model", func(t *testing.T) {
				expectedBeforeModelEventHandlerChecks(t, e)
			})

			_ = e.Next()

			t.Run("after model", func(t *testing.T) {
				expectedAfterModelEventHandlerChecks(t, e)
			})

			return nil
		},
	}

	modelErrorEventHandler := &hook.Handler[*core.ModelErrorEvent]{
		Priority: -999,
		Func: func(e *core.ModelErrorEvent) error {
			t.Run("before model error", func(t *testing.T) {
				expectedBeforeModelEventHandlerChecks(t, &e.ModelEvent)
				if v := e.Error.Error(); v != "error_a" {
					t.Fatalf("Expected error %q, got %q", "error_a", v)
				}
			})

			_ = e.Next()

			t.Run("after model error", func(t *testing.T) {
				expectedAfterModelEventHandlerChecks(t, &e.ModelEvent)
				if v := e.Error.Error(); v != "error_d" {
					t.Fatalf("Expected error %q, got %q", "error_d", v)
				}
			})

			return nil
		},
	}

	recordEventHandler := &hook.Handler[*core.RecordEvent]{
		Priority: -999,
		Func: func(e *core.RecordEvent) error {
			t.Run("before record", func(t *testing.T) {
				expectedBeforeRecordEventHandlerChecks(t, e)
			})

			changeRecordEventBefore(e)

			_ = e.Next()

			t.Run("after record", func(t *testing.T) {
				expectedAfterRecordEventHandlerChecks(t, e)
			})

			changeRecordEventAfter(e)

			return nil
		},
	}

	recordErrorEventHandler := &hook.Handler[*core.RecordErrorEvent]{
		Priority: -999,
		Func: func(e *core.RecordErrorEvent) error {
			t.Run("before record error", func(t *testing.T) {
				expectedBeforeRecordEventHandlerChecks(t, &e.RecordEvent)
				if v := e.Error.Error(); v != "error_a" {
					t.Fatalf("Expected error %q, got %q", "error_c", v)
				}
			})

			changeRecordEventBefore(&e.RecordEvent)
			e.Error = errors.New("error_b")

			_ = e.Next()

			t.Run("after record error", func(t *testing.T) {
				expectedAfterRecordEventHandlerChecks(t, &e.RecordEvent)
				if v := e.Error.Error(); v != "error_c" {
					t.Fatalf("Expected error %q, got %q", "error_c", v)
				}
			})

			changeRecordEventAfter(&e.RecordEvent)
			e.Error = errors.New("error_d")

			return nil
		},
	}

	// OnModelValidate
	app.OnRecordValidate().Bind(recordEventHandler)
	app.OnModelValidate().Bind(modelEventHandler)
	app.OnModelValidate().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelCreate
	app.OnRecordCreate().Bind(recordEventHandler)
	app.OnModelCreate().Bind(modelEventHandler)
	app.OnModelCreate().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelCreateExecute
	app.OnRecordCreateExecute().Bind(recordEventHandler)
	app.OnModelCreateExecute().Bind(modelEventHandler)
	app.OnModelCreateExecute().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterCreateSuccess
	app.OnRecordAfterCreateSuccess().Bind(recordEventHandler)
	app.OnModelAfterCreateSuccess().Bind(modelEventHandler)
	app.OnModelAfterCreateSuccess().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterCreateError
	app.OnRecordAfterCreateError().Bind(recordErrorEventHandler)
	app.OnModelAfterCreateError().Bind(modelErrorEventHandler)
	app.OnModelAfterCreateError().Trigger(createModelErrorEvent(), modelErrorEventFinalizer)

	// OnModelUpdate
	app.OnRecordUpdate().Bind(recordEventHandler)
	app.OnModelUpdate().Bind(modelEventHandler)
	app.OnModelUpdate().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelUpdateExecute
	app.OnRecordUpdateExecute().Bind(recordEventHandler)
	app.OnModelUpdateExecute().Bind(modelEventHandler)
	app.OnModelUpdateExecute().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterUpdateSuccess
	app.OnRecordAfterUpdateSuccess().Bind(recordEventHandler)
	app.OnModelAfterUpdateSuccess().Bind(modelEventHandler)
	app.OnModelAfterUpdateSuccess().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterUpdateError
	app.OnRecordAfterUpdateError().Bind(recordErrorEventHandler)
	app.OnModelAfterUpdateError().Bind(modelErrorEventHandler)
	app.OnModelAfterUpdateError().Trigger(createModelErrorEvent(), modelErrorEventFinalizer)

	// OnModelDelete
	app.OnRecordDelete().Bind(recordEventHandler)
	app.OnModelDelete().Bind(modelEventHandler)
	app.OnModelDelete().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelDeleteExecute
	app.OnRecordDeleteExecute().Bind(recordEventHandler)
	app.OnModelDeleteExecute().Bind(modelEventHandler)
	app.OnModelDeleteExecute().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterDeleteSuccess
	app.OnRecordAfterDeleteSuccess().Bind(recordEventHandler)
	app.OnModelAfterDeleteSuccess().Bind(modelEventHandler)
	app.OnModelAfterDeleteSuccess().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterDeleteError
	app.OnRecordAfterDeleteError().Bind(recordErrorEventHandler)
	app.OnModelAfterDeleteError().Bind(modelErrorEventHandler)
	app.OnModelAfterDeleteError().Trigger(createModelErrorEvent(), modelErrorEventFinalizer)
}

func TestRecordSave(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name        string
		record      func(app core.App) (*core.Record, error)
		expectError bool
	}{
		// trigger validators
		{
			name: "create - trigger validators",
			record: func(app core.App) (*core.Record, error) {
				c, _ := app.FindCollectionByNameOrId("demo2")
				record := core.NewRecord(c)
				return record, nil
			},
			expectError: true,
		},
		{
			name: "update - trigger validators",
			record: func(app core.App) (*core.Record, error) {
				record, _ := app.FindFirstRecordByData("demo2", "title", "test1")
				record.Set("title", "")
				return record, nil
			},
			expectError: true,
		},

		// create
		{
			name: "create base record",
			record: func(app core.App) (*core.Record, error) {
				c, _ := app.FindCollectionByNameOrId("demo2")
				record := core.NewRecord(c)
				record.Set("title", "new_test")
				return record, nil
			},
			expectError: false,
		},
		{
			name: "create auth record",
			record: func(app core.App) (*core.Record, error) {
				c, _ := app.FindCollectionByNameOrId("nologin")
				record := core.NewRecord(c)
				record.Set("email", "test_new@example.com")
				record.Set("password", "1234567890")
				return record, nil
			},
			expectError: false,
		},
		{
			name: "create view record",
			record: func(app core.App) (*core.Record, error) {
				c, _ := app.FindCollectionByNameOrId("view2")
				record := core.NewRecord(c)
				record.Set("state", true)
				return record, nil
			},
			expectError: true, // view records are read-only
		},

		// update
		{
			name: "update base record",
			record: func(app core.App) (*core.Record, error) {
				record, _ := app.FindFirstRecordByData("demo2", "title", "test1")
				record.Set("title", "test_new")
				return record, nil
			},
			expectError: false,
		},
		{
			name: "update auth record",
			record: func(app core.App) (*core.Record, error) {
				record, _ := app.FindAuthRecordByEmail("nologin", "test@example.com")
				record.Set("name", "test_new")
				record.Set("email", "test_new@example.com")
				return record, nil
			},
			expectError: false,
		},
		{
			name: "update view record",
			record: func(app core.App) (*core.Record, error) {
				record, _ := app.FindFirstRecordByData("view2", "state", true)
				record.Set("state", false)
				return record, nil
			},
			expectError: true, // view records are read-only
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			record, err := s.record(app)
			if err != nil {
				t.Fatalf("Failed to retrieve test record: %v", err)
			}

			saveErr := app.Save(record)

			hasErr := saveErr != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", hasErr, s.expectError, saveErr)
			}

			if hasErr {
				return
			}

			// the record should always have an id after successful Save
			if record.Id == "" {
				t.Fatal("Expected record id to be set")
			}

			if record.IsNew() {
				t.Fatal("Expected the record to be marked as not new")
			}

			// refetch and compare the serialization
			refreshed, err := app.FindRecordById(record.Collection(), record.Id)
			if err != nil {
				t.Fatal(err)
			}

			rawRefreshed, err := refreshed.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			raw, err := record.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(raw, rawRefreshed) {
				t.Fatalf("Expected the refreshed record to be the same as the saved one, got\n%s\nVS\n%s", raw, rawRefreshed)
			}
		})
	}
}

func TestRecordSaveIdFromOtherCollection(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	baseCollection, _ := app.FindCollectionByNameOrId("demo2")
	authCollection, _ := app.FindCollectionByNameOrId("nologin")

	// base collection test
	r1 := core.NewRecord(baseCollection)
	r1.Set("title", "test_new")
	r1.Set("id", "mk5fmymtx4wsprk") // existing id of demo3 record
	if err := app.Save(r1); err != nil {
		t.Fatalf("Expected nil, got error %v", err)
	}

	// auth collection test
	r2 := core.NewRecord(authCollection)
	r2.SetEmail("test_new@example.com")
	r2.SetPassword("1234567890")
	r2.Set("id", "gk390qegs4y47wn") // existing id of "clients" record
	if err := app.Save(r2); err == nil {
		t.Fatal("Expected error, got nil")
	}

	// try again with unique id
	r2.Set("id", strings.Repeat("a", 15))
	if err := app.Save(r2); err != nil {
		t.Fatalf("Expected nil, got error %v", err)
	}
}

func TestRecordSaveIdUpdateNoValidation(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	rec, err := app.FindRecordById("demo3", "7nwo8tuiatetxdm")
	if err != nil {
		t.Fatal(err)
	}

	rec.Id = strings.Repeat("a", 15)

	err = app.SaveNoValidate(rec)
	if err == nil {
		t.Fatal("Expected save to fail, got nil")
	}

	// no changes
	rec.Load(rec.Original().FieldsData())
	err = app.SaveNoValidate(rec)
	if err != nil {
		t.Fatalf("Expected save to succeed, got error %v", err)
	}
}

func TestRecordSaveWithAutoTokenKeyRefresh(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name           string
		payload        map[string]any
		expectedChange bool
	}{
		{
			"no email or password change",
			map[string]any{"name": "example"},
			false,
		},
		{
			"password change",
			map[string]any{"password": "1234567890"},
			true,
		},
		{
			"email change",
			map[string]any{"email": "test_update@example.com"},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			record, err := app.FindFirstRecordByFilter("nologin", "1=1")
			if err != nil {
				t.Fatal(err)
			}

			originalTokenKey := record.TokenKey()

			record.Load(s.payload)

			err = app.Save(record)
			if err != nil {
				t.Fatal(err)
			}

			newTokenKey := record.TokenKey()

			hasChange := originalTokenKey != newTokenKey

			if hasChange != s.expectedChange {
				t.Fatalf("Expected hasChange %v, got %v", s.expectedChange, hasChange)
			}
		})
	}
}

func TestRecordDelete(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demoCollection, _ := app.FindCollectionByNameOrId("demo2")

	// delete unsaved record
	// ---
	newRec := core.NewRecord(demoCollection)
	if err := app.Delete(newRec); err == nil {
		t.Fatal("(newRec) Didn't expect to succeed deleting unsaved record")
	}

	// delete view record
	// ---
	viewRec, _ := app.FindRecordById("view2", "84nmscqy84lsi1t")
	if err := app.Delete(viewRec); err == nil {
		t.Fatal("(viewRec) Didn't expect to succeed deleting view record")
	}
	// check if it still exists
	viewRec, _ = app.FindRecordById(viewRec.Collection().Id, viewRec.Id)
	if viewRec == nil {
		t.Fatal("(viewRec) Expected view record to still exists")
	}

	// delete existing record + external auths
	// ---
	rec1, _ := app.FindRecordById("users", "4q1xlclmfloku33")
	if err := app.Delete(rec1); err != nil {
		t.Fatalf("(rec1) Expected nil, got error %v", err)
	}
	// check if it was really deleted
	if refreshed, _ := app.FindRecordById(rec1.Collection().Id, rec1.Id); refreshed != nil {
		t.Fatalf("(rec1) Expected record to be deleted, got %v", refreshed)
	}
	// check if the external auths were deleted
	if auths, _ := app.FindAllExternalAuthsByRecord(rec1); len(auths) > 0 {
		t.Fatalf("(rec1) Expected external auths to be deleted, got %v", auths)
	}

	// delete existing record while being part of a non-cascade required relation
	// ---
	rec2, _ := app.FindRecordById("demo3", "7nwo8tuiatetxdm")
	if err := app.Delete(rec2); err == nil {
		t.Fatalf("(rec2) Expected error, got nil")
	}

	// delete existing record + cascade
	// ---
	calledQueries := []string{}
	app.NonconcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		calledQueries = append(calledQueries, sql)
	}
	app.DB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		calledQueries = append(calledQueries, sql)
	}
	app.NonconcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		calledQueries = append(calledQueries, sql)
	}
	app.DB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		calledQueries = append(calledQueries, sql)
	}
	rec3, _ := app.FindRecordById("users", "oap640cot4yru2s")
	// delete
	if err := app.Delete(rec3); err != nil {
		t.Fatalf("(rec3) Expected nil, got error %v", err)
	}
	// check if it was really deleted
	rec3, _ = app.FindRecordById(rec3.Collection().Id, rec3.Id)
	if rec3 != nil {
		t.Fatalf("(rec3) Expected record to be deleted, got %v", rec3)
	}
	// check if the operation cascaded
	rel, _ := app.FindRecordById("demo1", "84nmscqy84lsi1t")
	if rel != nil {
		t.Fatalf("(rec3) Expected the delete to cascade, found relation %v", rel)
	}
	// ensure that the json rel fields were prefixed
	joinedQueries := strings.Join(calledQueries, " ")
	expectedRelManyPart := "SELECT `demo1`.* FROM `demo1` WHERE EXISTS (SELECT 1 FROM json_each(CASE WHEN json_valid([[demo1.rel_many]]) THEN [[demo1.rel_many]] ELSE json_array([[demo1.rel_many]]) END) {{__je__}} WHERE [[__je__.value]]='"
	if !strings.Contains(joinedQueries, expectedRelManyPart) {
		t.Fatalf("(rec3) Expected the cascade delete to call the query \n%v, got \n%v", expectedRelManyPart, calledQueries)
	}
	expectedRelOnePart := "SELECT `demo1`.* FROM `demo1` WHERE (`demo1`.`rel_one`='"
	if !strings.Contains(joinedQueries, expectedRelOnePart) {
		t.Fatalf("(rec3) Expected the cascade delete to call the query \n%v, got \n%v", expectedRelOnePart, calledQueries)
	}
}

func TestRecordDeleteBatchProcessing(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	if err := createMockBatchProcessingData(app); err != nil {
		t.Fatal(err)
	}

	// find and delete the first c1 record to trigger cascade
	mainRecord, _ := app.FindRecordById("c1", "a")
	if err := app.Delete(mainRecord); err != nil {
		t.Fatal(err)
	}

	// check if the main record was deleted
	_, err := app.FindRecordById(mainRecord.Collection().Id, mainRecord.Id)
	if err == nil {
		t.Fatal("The main record wasn't deleted")
	}

	// check if the c1 b rel field were updated
	c1RecordB, err := app.FindRecordById("c1", "b")
	if err != nil || c1RecordB.GetString("rel") != "" {
		t.Fatalf("Expected c1RecordB.rel to be nil, got %v", c1RecordB.GetString("rel"))
	}

	// check if the c2 rel fields were updated
	c2Records, err := app.FindAllRecords("c2", nil)
	if err != nil || len(c2Records) == 0 {
		t.Fatalf("Failed to fetch c2 records: %v", err)
	}
	for _, r := range c2Records {
		ids := r.GetStringSlice("rel")
		if len(ids) != 1 || ids[0] != "b" {
			t.Fatalf("Expected only 'b' rel id, got %v", ids)
		}
	}

	// check if all c3 relations were deleted
	c3Records, err := app.FindAllRecords("c3", nil)
	if err != nil {
		t.Fatalf("Failed to fetch c3 records: %v", err)
	}
	if total := len(c3Records); total != 0 {
		t.Fatalf("Expected c3 records to be deleted, found %d", total)
	}
}

func createMockBatchProcessingData(app core.App) error {
	// create mock collection without relation
	c1 := core.NewBaseCollection("c1")
	c1.Id = "c1"
	c1.Fields.Add(
		&core.TextField{Name: "text"},
		&core.RelationField{
			Name:          "rel",
			MaxSelect:     1,
			CollectionId:  "c1",
			CascadeDelete: false, // should unset all rel fields
		},
	)
	if err := app.SaveNoValidate(c1); err != nil {
		return err
	}

	// create mock collection with a multi-rel field
	c2 := core.NewBaseCollection("c2")
	c2.Id = "c2"
	c2.Fields.Add(
		&core.TextField{Name: "text"},
		&core.RelationField{
			Name:          "rel",
			MaxSelect:     10,
			CollectionId:  "c1",
			CascadeDelete: false, // should unset all rel fields
		},
	)
	if err := app.SaveNoValidate(c2); err != nil {
		return err
	}

	// create mock collection with a single-rel field
	c3 := core.NewBaseCollection("c3")
	c3.Id = "c3"
	c3.Fields.Add(
		&core.RelationField{
			Name:          "rel",
			MaxSelect:     1,
			CollectionId:  "c1",
			CascadeDelete: true, // should delete all c3 records
		},
	)
	if err := app.SaveNoValidate(c3); err != nil {
		return err
	}

	// insert mock records
	c1RecordA := core.NewRecord(c1)
	c1RecordA.Id = "a"
	c1RecordA.Set("rel", c1RecordA.Id) // self reference
	if err := app.SaveNoValidate(c1RecordA); err != nil {
		return err
	}
	c1RecordB := core.NewRecord(c1)
	c1RecordB.Id = "b"
	c1RecordB.Set("rel", c1RecordA.Id) // rel to another record from the same collection
	if err := app.SaveNoValidate(c1RecordB); err != nil {
		return err
	}
	for i := 0; i < 4500; i++ {
		c2Record := core.NewRecord(c2)
		c2Record.Set("rel", []string{c1RecordA.Id, c1RecordB.Id})
		if err := app.SaveNoValidate(c2Record); err != nil {
			return err
		}

		c3Record := core.NewRecord(c3)
		c3Record.Set("rel", c1RecordA.Id)
		if err := app.SaveNoValidate(c3Record); err != nil {
			return err
		}
	}

	// set the same id as the relation for at least 1 record
	// to check whether the correct condition will be added
	c3Record := core.NewRecord(c3)
	c3Record.Set("rel", c1RecordA.Id)
	c3Record.Id = c1RecordA.Id
	if err := app.SaveNoValidate(c3Record); err != nil {
		return err
	}

	return nil
}

// -------------------------------------------------------------------

type mockField struct {
	core.TextField
}

func (f *mockField) FindGetter(key string) core.GetterFunc {
	switch key {
	case f.Name + ":test":
		return func(record *core.Record) any {
			return "modifier_get"
		}
	default:
		return nil
	}
}

func (f *mockField) FindSetter(key string) core.SetterFunc {
	switch key {
	case f.Name:
		return func(record *core.Record, raw any) {
			record.SetRaw(f.Name, cast.ToString(raw))
		}
	case f.Name + ":test":
		return func(record *core.Record, raw any) {
			record.SetRaw(f.Name, "modifier_set")
		}
	default:
		return nil
	}
}
