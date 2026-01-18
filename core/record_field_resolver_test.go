package core_test

import (
	"encoding/json"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestRecordFieldResolverAllowedFields(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	r := core.NewRecordFieldResolver(app, collection, nil, false)

	fields := r.AllowedFields()
	if len(fields) != 8 {
		t.Fatalf("Expected %d original allowed fields, got %d", 8, len(fields))
	}

	// change the allowed fields
	newFields := []string{"a", "b", "c"}
	expected := slices.Clone(newFields)
	r.SetAllowedFields(newFields)

	// change the new fields to ensure that the slice was cloned
	newFields[2] = "d"

	fields = r.AllowedFields()
	if len(fields) != len(expected) {
		t.Fatalf("Expected %d changed allowed fields, got %d", len(expected), len(fields))
	}

	for i, v := range expected {
		if fields[i] != v {
			t.Errorf("[%d] Expected field %q", i, v)
		}
	}
}

func TestRecordFieldResolverAllowHiddenFields(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	r := core.NewRecordFieldResolver(app, collection, nil, false)

	allowHiddenFields := r.AllowHiddenFields()
	if allowHiddenFields {
		t.Fatalf("Expected original allowHiddenFields %v, got %v", allowHiddenFields, !allowHiddenFields)
	}

	// change the flag
	expected := !allowHiddenFields
	r.SetAllowHiddenFields(expected)

	allowHiddenFields = r.AllowHiddenFields()
	if allowHiddenFields != expected {
		t.Fatalf("Expected changed allowHiddenFields %v, got %v", expected, allowHiddenFields)
	}
}

func TestRecordFieldResolverUpdateQuery(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	authRecord, err := app.FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestInfo := &core.RequestInfo{
		Context: "ctx",
		Headers: map[string]string{
			"a": "123",
			"b": "456",
		},
		Query: map[string]string{
			"a": "", // to ensure that :isset returns true because the key exists
			"b": "123",
		},
		Body: map[string]any{
			"a":                  nil, // to ensure that :isset returns true because the key exists
			"b":                  123,
			"number":             10,
			"select_many":        []string{"optionA", "optionC"},
			"rel_one":            "test",
			"rel_many":           []string{"test1", "test2"},
			"file_one":           "test",
			"file_many":          []string{"test1", "test2", "test3"},
			"self_rel_one":       "test",
			"self_rel_many":      []string{"test1"},
			"rel_many_cascade":   []string{"test1", "test2"},
			"rel_one_cascade":    "test1",
			"rel_one_no_cascade": "test1",
		},
		Auth: authRecord,
	}

	scenarios := []struct {
		name               string
		collectionIdOrName string
		rule               string
		allowHiddenFields  bool
		expectQuery        string
	}{
		{
			"none relation field (with all default operators)",
			"demo4",
			"title = true || title != 'test' || title ~ 'test1' || title !~ '%test2' || title > 1 || title >= 2 || title < 3 || title <= 4",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE ([[demo4.title]] = 1 OR [[demo4.title]] IS NOT {:TEST} OR [[demo4.title]] LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] NOT LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] > {:TEST} OR [[demo4.title]] >= {:TEST} OR [[demo4.title]] < {:TEST} OR [[demo4.title]] <= {:TEST})",
		},
		{
			"none relation field (with all opt/any operators)",
			"demo4",
			"title ?= true || title ?!= 'test' || title ?~ 'test1' || title ?!~ '%test2' || title ?> 1 || title ?>= 2 || title ?< 3 || title ?<= 4",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE ([[demo4.title]] = 1 OR [[demo4.title]] IS NOT {:TEST} OR [[demo4.title]] LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] NOT LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] > {:TEST} OR [[demo4.title]] >= {:TEST} OR [[demo4.title]] < {:TEST} OR [[demo4.title]] <= {:TEST})",
		},
		{
			"single direct rel",
			"demo4",
			"self_rel_one > true",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE [[demo4.self_rel_one]] > 1",
		},
		{
			"single direct rel (with id)",
			"demo4",
			"self_rel_one.id > true", // should NOT have join
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE [[demo4.self_rel_one]] > 1",
		},
		{
			"multiple direct rel (with id)",
			"demo4",
			"self_rel_many.id ?> true", // should have join
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] WHERE [[demo4_self_rel_many.id]] > 1",
		},
		{
			"rel to collection with empty list rule",
			"demo4",
			"self_rel_one.created > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] WHERE [[demo4_self_rel_one.created]] > 1",
		},
		{
			"rel to collection with non-empty list rule",
			"demo4",
			"rel_one_cascade.created > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo3` `demo4_rel_one_cascade` ON [[demo4_rel_one_cascade.id]] = [[demo4.rel_one_cascade]] WHERE (({:TEST} IS NOT '' AND {:TEST} IS NOT {:TEST})) AND ([[demo4_rel_one_cascade.created]] > 1)",
		},
		{
			"rel to collection with non-empty list rule (with allowHiddenFields)",
			"demo4",
			"rel_one_cascade.created > true",
			true,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo3` `demo4_rel_one_cascade` ON [[demo4_rel_one_cascade.id]] = [[demo4.rel_one_cascade]] WHERE [[demo4_rel_one_cascade.created]] > 1",
		},
		{
			"rel to collection with superusers only list rule",
			"demo1",
			"rel_many.created ?> true",
			false,
			"",
		},
		{
			"rel to collection with superusers only list rule (with allowHiddenFields)",
			"demo1",
			"rel_many.created ?> true",
			true,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo1.rel_many]]), json_type([[demo1.rel_many]])='array', FALSE) THEN [[demo1.rel_many]] ELSE json_array([[demo1.rel_many]]) END) `__je_demo1_rel_many` LEFT JOIN `users` `demo1_rel_many` ON [[demo1_rel_many.id]] = [[__je_demo1_rel_many.value]] WHERE [[demo1_rel_many.created]] > 1",
		},
		{
			"nested rels with all empty list rules",
			"demo4",
			"self_rel_one.self_rel_one.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] LEFT JOIN `demo4` `demo4_self_rel_one_self_rel_one` ON [[demo4_self_rel_one_self_rel_one.id]] = [[demo4_self_rel_one.self_rel_one]] WHERE [[demo4_self_rel_one_self_rel_one.title]] > 1",
		},
		{
			"nested rels with non-empty list rule",
			"demo4",
			"self_rel_one.rel_one_cascade.created > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] LEFT JOIN `demo3` `demo4_self_rel_one_rel_one_cascade` ON [[demo4_self_rel_one_rel_one_cascade.id]] = [[demo4_self_rel_one.rel_one_cascade]] WHERE (({:TEST} IS NOT '' AND {:TEST} IS NOT {:TEST})) AND ([[demo4_self_rel_one_rel_one_cascade.created]] > 1)",
		},
		{
			"nested rels with non-empty list rule (joins reuse test)",
			"demo4",
			"self_rel_one.rel_one_cascade.created > true && self_rel_one.rel_one_cascade.updated > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] LEFT JOIN `demo3` `demo4_self_rel_one_rel_one_cascade` ON [[demo4_self_rel_one_rel_one_cascade.id]] = [[demo4_self_rel_one.rel_one_cascade]] WHERE (({:TEST} IS NOT '' AND {:TEST} IS NOT {:TEST})) AND (([[demo4_self_rel_one_rel_one_cascade.created]] > 1 AND [[demo4_self_rel_one_rel_one_cascade.updated]] > 1))",
		},
		{
			"nested rels with non-empty list rule (with allowHiddenFields)",
			"demo4",
			"self_rel_one.rel_one_cascade.created > true",
			true,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] LEFT JOIN `demo3` `demo4_self_rel_one_rel_one_cascade` ON [[demo4_self_rel_one_rel_one_cascade.id]] = [[demo4_self_rel_one.rel_one_cascade]] WHERE [[demo4_self_rel_one_rel_one_cascade.created]] > 1",
		},
		{
			"non-relation field + single rel",
			"demo4",
			"title > true || self_rel_one.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] WHERE ([[demo4.title]] > 1 OR [[demo4_self_rel_one.title]] > 1)",
		},
		{
			"nested incomplete relations (opt/any operator)",
			"demo4",
			"self_rel_many.self_rel_one ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] WHERE [[demo4_self_rel_many.self_rel_one]] > 1",
		},
		{
			"nested incomplete relations (multi-match operator)",
			"demo4",
			"self_rel_many.self_rel_one > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] WHERE ((([[demo4_self_rel_many.self_rel_one]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many.self_rel_one]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4.self_rel_many]]), json_type([[__mm_demo4.self_rel_many]])='array', FALSE) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))))",
		},
		{
			"nested complete relations (opt/any operator)",
			"demo4",
			"self_rel_many.self_rel_one.title ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] WHERE [[demo4_self_rel_many_self_rel_one.title]] > 1",
		},
		{
			"nested complete relations (multi-match operator)",
			"demo4",
			"self_rel_many.self_rel_one.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] WHERE ((([[demo4_self_rel_many_self_rel_one.title]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many_self_rel_one.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4.self_rel_many]]), json_type([[__mm_demo4.self_rel_many]])='array', FALSE) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one` ON [[__mm_demo4_self_rel_many_self_rel_one.id]] = [[__mm_demo4_self_rel_many.self_rel_one]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))))",
		},
		{
			"repeated nested relations (opt/any operator)",
			"demo4",
			"self_rel_many.self_rel_one.self_rel_many.self_rel_one.title ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4_self_rel_many_self_rel_one.self_rel_many]]), json_type([[demo4_self_rel_many_self_rel_one.self_rel_many]])='array', FALSE) THEN [[demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `__je_demo4_self_rel_many_self_rel_one_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many` ON [[demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[__je_demo4_self_rel_many_self_rel_one_self_rel_many.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] WHERE [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]] > 1",
		},
		{
			"repeated nested relations (multi-match operator)",
			"demo4",
			"self_rel_many.self_rel_one.self_rel_many.self_rel_one.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4_self_rel_many_self_rel_one.self_rel_many]]), json_type([[demo4_self_rel_many_self_rel_one.self_rel_many]])='array', FALSE) THEN [[demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `__je_demo4_self_rel_many_self_rel_one_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many` ON [[demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[__je_demo4_self_rel_many_self_rel_one_self_rel_many.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] WHERE ((([[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4.self_rel_many]]), json_type([[__mm_demo4.self_rel_many]])='array', FALSE) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one` ON [[__mm_demo4_self_rel_many_self_rel_one.id]] = [[__mm_demo4_self_rel_many.self_rel_one]] LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]]), json_type([[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]])='array', FALSE) THEN [[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `__mm_demo4_self_rel_many_self_rel_one_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one_self_rel_many` ON [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many_je.value]] LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))))",
		},
		{
			"multiple relations (opt/any operators)",
			"demo4",
			"self_rel_many.title ?= 'test' || self_rel_one.json_object.a ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] WHERE ([[demo4_self_rel_many.title]] = {:TEST} OR (CASE WHEN json_valid([[demo4_self_rel_one.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_one.json_object]], '$.a') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_one.json_object]]), '$.pb.a') END) > 1)",
		},
		{
			"multiple relations (multi-match operators)",
			"demo4",
			"self_rel_many.title = 'test' || self_rel_one.json_object.a > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] WHERE ((([[demo4_self_rel_many.title]] = {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4.self_rel_many]]), json_type([[__mm_demo4.self_rel_many]])='array', FALSE) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = {:TEST})))) OR (CASE WHEN json_valid([[demo4_self_rel_one.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_one.json_object]], '$.a') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_one.json_object]]), '$.pb.a') END) > 1)",
		},
		{
			"back relations via single relation field (without unique index)",
			"demo3",
			"demo4_via_rel_one_cascade.id = true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_one_cascade` ON [[demo3_demo4_via_rel_one_cascade.rel_one_cascade]] = [[demo3.id]] WHERE ((([[demo3_demo4_via_rel_one_cascade.id]] = 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo3_demo4_via_rel_one_cascade.id]] as [[multiMatchValue]] FROM `demo3` `__mm_demo3` LEFT JOIN `demo4` `__mm_demo3_demo4_via_rel_one_cascade` ON [[__mm_demo3_demo4_via_rel_one_cascade.rel_one_cascade]] = [[__mm_demo3.id]] WHERE `__mm_demo3`.`id` = `demo3`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = 1)))))",
		},
		{
			"back relations via single relation field (with unique index)",
			"demo3",
			"demo4_via_rel_one_unique.id = true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_one_unique` ON [[demo3_demo4_via_rel_one_unique.rel_one_unique]] = [[demo3.id]] WHERE [[demo3_demo4_via_rel_one_unique.id]] = 1",
		},
		{
			"back relations via multiple relation field (opt/any operators)",
			"demo3",
			"demo4_via_rel_many_cascade.id ?= true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade` ON [[demo3.id]] IN (SELECT [[__je_demo3_demo4_via_rel_many_cascade.value]] FROM json_each(CASE WHEN iif(json_valid([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]), json_type([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]])='array', FALSE) THEN [[demo3_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{__je_demo3_demo4_via_rel_many_cascade}}) WHERE [[demo3_demo4_via_rel_many_cascade.id]] = 1",
		},
		{
			"back relations via multiple relation field (multi-match operators)",
			"demo3",
			"demo4_via_rel_many_cascade.id = true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade` ON [[demo3.id]] IN (SELECT [[__je_demo3_demo4_via_rel_many_cascade.value]] FROM json_each(CASE WHEN iif(json_valid([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]), json_type([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]])='array', FALSE) THEN [[demo3_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{__je_demo3_demo4_via_rel_many_cascade}}) WHERE ((([[demo3_demo4_via_rel_many_cascade.id]] = 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo3_demo4_via_rel_many_cascade.id]] as [[multiMatchValue]] FROM `demo3` `__mm_demo3` LEFT JOIN `demo4` `__mm_demo3_demo4_via_rel_many_cascade` ON [[__mm_demo3.id]] IN (SELECT [[__je___mm_demo3_demo4_via_rel_many_cascade.value]] FROM json_each(CASE WHEN iif(json_valid([[__mm_demo3_demo4_via_rel_many_cascade.rel_many_cascade]]), json_type([[__mm_demo3_demo4_via_rel_many_cascade.rel_many_cascade]])='array', FALSE) THEN [[__mm_demo3_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[__mm_demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{__je___mm_demo3_demo4_via_rel_many_cascade}}) WHERE `__mm_demo3`.`id` = `demo3`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = 1)))))",
		},
		{
			"back relations via unique multiple relation field (should be the same as multi-match)",
			"demo3",
			"demo4_via_rel_many_unique.id = true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_many_unique` ON [[demo3.id]] IN (SELECT [[__je_demo3_demo4_via_rel_many_unique.value]] FROM json_each(CASE WHEN iif(json_valid([[demo3_demo4_via_rel_many_unique.rel_many_unique]]), json_type([[demo3_demo4_via_rel_many_unique.rel_many_unique]])='array', FALSE) THEN [[demo3_demo4_via_rel_many_unique.rel_many_unique]] ELSE json_array([[demo3_demo4_via_rel_many_unique.rel_many_unique]]) END) {{__je_demo3_demo4_via_rel_many_unique}}) WHERE ((([[demo3_demo4_via_rel_many_unique.id]] = 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo3_demo4_via_rel_many_unique.id]] as [[multiMatchValue]] FROM `demo3` `__mm_demo3` LEFT JOIN `demo4` `__mm_demo3_demo4_via_rel_many_unique` ON [[__mm_demo3.id]] IN (SELECT [[__je___mm_demo3_demo4_via_rel_many_unique.value]] FROM json_each(CASE WHEN iif(json_valid([[__mm_demo3_demo4_via_rel_many_unique.rel_many_unique]]), json_type([[__mm_demo3_demo4_via_rel_many_unique.rel_many_unique]])='array', FALSE) THEN [[__mm_demo3_demo4_via_rel_many_unique.rel_many_unique]] ELSE json_array([[__mm_demo3_demo4_via_rel_many_unique.rel_many_unique]]) END) {{__je___mm_demo3_demo4_via_rel_many_unique}}) WHERE `__mm_demo3`.`id` = `demo3`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = 1)))))",
		},
		{
			"view back relation with non-empty and superusers list rules",
			"demo1",
			"view1_via_rel_one.rel_many.created ?> true",
			false,
			"",
		},
		{
			"view back relation with non-empty and superusers list rules (with allowHiddenFields)",
			"demo1",
			"view1_via_rel_one.rel_many.created ?> true",
			true,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN `view1` `demo1_view1_via_rel_one` ON [[demo1_view1_via_rel_one.rel_one]] = [[demo1.id]] LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo1_view1_via_rel_one.rel_many]]), json_type([[demo1_view1_via_rel_one.rel_many]])='array', FALSE) THEN [[demo1_view1_via_rel_one.rel_many]] ELSE json_array([[demo1_view1_via_rel_one.rel_many]]) END) `__je_demo1_view1_via_rel_one_rel_many` LEFT JOIN `users` `demo1_view1_via_rel_one_rel_many` ON [[demo1_view1_via_rel_one_rel_many.id]] = [[__je_demo1_view1_via_rel_one_rel_many.value]] WHERE [[demo1_view1_via_rel_one_rel_many.created]] > 1",
		},
		{
			"recursive back relations with non-empty list rule",
			"demo3",
			"demo4_via_rel_many_cascade.rel_one_cascade.demo4_via_rel_many_cascade.id ?= true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade` ON [[demo3.id]] IN (SELECT [[__je_demo3_demo4_via_rel_many_cascade.value]] FROM json_each(CASE WHEN iif(json_valid([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]), json_type([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]])='array', FALSE) THEN [[demo3_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{__je_demo3_demo4_via_rel_many_cascade}}) LEFT JOIN `demo3` `demo3_demo4_via_rel_many_cascade_rel_one_cascade` ON [[demo3_demo4_via_rel_many_cascade_rel_one_cascade.id]] = [[demo3_demo4_via_rel_many_cascade.rel_one_cascade]] LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade` ON [[demo3_demo4_via_rel_many_cascade_rel_one_cascade.id]] IN (SELECT [[__je_demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.value]] FROM json_each(CASE WHEN iif(json_valid([[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]]), json_type([[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]])='array', FALSE) THEN [[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{__je_demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade}}) WHERE (({:TEST} IS NOT '' AND {:TEST} IS NOT {:TEST})) AND ([[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.id]] = 1)",
		},
		{
			"recursive back relations with non-empty list rule (with allowHiddenFields)",
			"demo3",
			"demo4_via_rel_many_cascade.rel_one_cascade.demo4_via_rel_many_cascade.id ?= true",
			true,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade` ON [[demo3.id]] IN (SELECT [[__je_demo3_demo4_via_rel_many_cascade.value]] FROM json_each(CASE WHEN iif(json_valid([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]), json_type([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]])='array', FALSE) THEN [[demo3_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{__je_demo3_demo4_via_rel_many_cascade}}) LEFT JOIN `demo3` `demo3_demo4_via_rel_many_cascade_rel_one_cascade` ON [[demo3_demo4_via_rel_many_cascade_rel_one_cascade.id]] = [[demo3_demo4_via_rel_many_cascade.rel_one_cascade]] LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade` ON [[demo3_demo4_via_rel_many_cascade_rel_one_cascade.id]] IN (SELECT [[__je_demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.value]] FROM json_each(CASE WHEN iif(json_valid([[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]]), json_type([[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]])='array', FALSE) THEN [[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{__je_demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade}}) WHERE [[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.id]] = 1",
		},
		{
			"@collection join (opt/any operators)",
			"demo4",
			"@collection.demo1.text ?> true || @collection.demo2.active ?> true || @collection.demo1:demo1_alias.file_one ?> true",
			true,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo1` `__collection_demo1` LEFT JOIN `demo2` `__collection_demo2` LEFT JOIN `demo1` `__collection_alias_demo1_alias` WHERE ([[__collection_demo1.text]] > 1 OR [[__collection_demo2.active]] > 1 OR [[__collection_alias_demo1_alias.file_one]] > 1)",
		},
		{
			"@collection join (multi-match operators)",
			"demo4",
			"@collection.demo1.text > true || @collection.demo2.active > true || @collection.demo1.file_one > true",
			true,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo1` `__collection_demo1` LEFT JOIN `demo2` `__collection_demo2` WHERE ((([[__collection_demo1.text]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___collection_demo1.text]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo1` `__mm___collection_demo1` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) OR (([[__collection_demo2.active]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___collection_demo2.active]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo2` `__mm___collection_demo2` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) OR (([[__collection_demo1.file_one]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___collection_demo1.file_one]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo1` `__mm___collection_demo1` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))))",
		},
		{
			"@request.auth fields",
			"demo4",
			"@request.auth.id > true || @request.auth.username > true || @request.auth.rel.title > true || @request.body.demo < true || @request.auth.missingA.missingB > false",
			true,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `users` `__auth_users` ON `__auth_users`.`id`={:p0} LEFT JOIN `demo2` `__auth_users_rel` ON [[__auth_users_rel.id]] = [[__auth_users.rel]] WHERE ({:TEST} > 1 OR [[__auth_users.username]] > 1 OR [[__auth_users_rel.title]] > 1 OR NULL < 1 OR NULL > 0)",
		},
		{
			"@request.* static fields",
			"demo4",
			"@request.context = true || @request.query.a = true || @request.query.b = true || @request.query.missing = true || @request.headers.a = true || @request.headers.missing = true",
			true,
			"SELECT `demo4`.* FROM `demo4` WHERE ({:TEST} = 1 OR '' = 1 OR {:TEST} = 1 OR '' = 1 OR {:TEST} = 1 OR '' = 1)",
		},
		{
			"direct hidden field (add emailVisibility)",
			"users",
			"email > true",
			false,
			"SELECT `users`.* FROM `users` WHERE ((([[users.email]] > 1) AND ([[users.emailVisibility]] = TRUE)))",
		},
		{
			"direct hidden field (force ignore emailVisibility)",
			"users",
			"email > true",
			true,
			"SELECT `users`.* FROM `users` WHERE [[users.email]] > 1",
		},
		{
			"mixed regular with hidden field and modifier (add emailVisibility)",
			"nologin",
			"id > true || email > true || email:lower > false",
			false,
			"SELECT `nologin`.* FROM `nologin` WHERE ([[nologin.id]] > 1 OR (([[nologin.email]] > 1) AND ([[nologin.emailVisibility]] = TRUE)) OR ((LOWER([[nologin.email]]) > 0) AND ([[nologin.emailVisibility]] = TRUE)))",
		},
		{
			"system filters in a public auth collection with hidden field and no allowHiddenFields (multi-match and add emailVisibility)",
			"demo4",
			"@collection.nologin.email > true || @request.auth.email > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `nologin` `__collection_nologin` WHERE ((((([[__collection_nologin.email]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___collection_nologin.email]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `nologin` `__mm___collection_nologin` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1))))) AND ([[__collection_nologin.emailVisibility]] = TRUE)) OR {:TEST} > 1)",
		},
		{
			"system filters in a superuser auth collection with hidden field and NO allowHiddenFields (multi-match and add emailVisibility)",
			"demo4",
			"@collection.users.email > true || @request.auth.email > true",
			false,
			"",
		},
		{
			"system filters in a superuser auth collection with hidden field and allowHiddenFields (multi-match and add emailVisibility)",
			"demo4",
			"@collection.users.email > true || @request.auth.email > true",
			true,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `users` `__collection_users` WHERE ((([[__collection_users.email]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___collection_users.email]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `users` `__mm___collection_users` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) OR {:TEST} > 1)",
		},
		{
			"collection filter in a non-empty list rule collection",
			"demo4",
			"@collection.demo3.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo3` `__collection_demo3` WHERE (({:TEST} IS NOT '' AND {:TEST} IS NOT {:TEST})) AND (((([[__collection_demo3.title]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___collection_demo3.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo3` `__mm___collection_demo3` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1))))))",
		},
		{
			"collection filter in a non-empty list rule collection (with allowHiddenFields)",
			"demo4",
			"@collection.demo3.title > true",
			true,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo3` `__collection_demo3` WHERE ((([[__collection_demo3.title]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___collection_demo3.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo3` `__mm___collection_demo3` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))))",
		},
		{
			"collection fields with :lower modifier",
			"demo1",
			"@request.body.rel_one:lower > true ||" +
				"@request.body.rel_many:lower > true ||" +
				"@request.body.rel_many.email:lower > true ||" +
				"text:lower > true ||" +
				"bool:lower > true ||" +
				"url:lower > true ||" +
				"select_one:lower > true ||" +
				"select_many:lower > true ||" +
				"file_one:lower > true ||" +
				"file_many:lower > true ||" +
				"number:lower > true ||" +
				"email:lower > true ||" +
				"datetime:lower > true ||" +
				"json:lower > true ||" +
				"rel_one:lower > true ||" +
				"rel_many:lower > true ||" +
				"rel_many.name:lower > true ||" +
				"created:lower > true",
			true,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN `users` `__data_users_rel_many` ON [[__data_users_rel_many.id]] IN ({:p0}, {:p1}) LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo1.rel_many]]), json_type([[demo1.rel_many]])='array', FALSE) THEN [[demo1.rel_many]] ELSE json_array([[demo1.rel_many]]) END) `__je_demo1_rel_many` LEFT JOIN `users` `demo1_rel_many` ON [[demo1_rel_many.id]] = [[__je_demo1_rel_many.value]] WHERE (LOWER({:infoLowerrel_oneTEST}) > 1 OR LOWER({:infoLowerrel_manyTEST}) > 1 OR ((LOWER([[__data_users_rel_many.email]]) > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT LOWER([[__mm___data_users_rel_many.email]]) as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN `users` `__mm___data_users_rel_many` ON [[__mm___data_users_rel_many.id]] IN ({:p4}, {:p5}) WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) OR LOWER([[demo1.text]]) > 1 OR LOWER([[demo1.bool]]) > 1 OR LOWER([[demo1.url]]) > 1 OR LOWER([[demo1.select_one]]) > 1 OR LOWER([[demo1.select_many]]) > 1 OR LOWER([[demo1.file_one]]) > 1 OR LOWER([[demo1.file_many]]) > 1 OR LOWER([[demo1.number]]) > 1 OR LOWER([[demo1.email]]) > 1 OR LOWER([[demo1.datetime]]) > 1 OR LOWER((CASE WHEN json_valid([[demo1.json]]) THEN JSON_EXTRACT([[demo1.json]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo1.json]]), '$.pb') END)) > 1 OR LOWER([[demo1.rel_one]]) > 1 OR LOWER([[demo1.rel_many]]) > 1 OR ((LOWER([[demo1_rel_many.name]]) > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT LOWER([[__mm_demo1_rel_many.name]]) as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo1.rel_many]]), json_type([[__mm_demo1.rel_many]])='array', FALSE) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) OR LOWER([[demo1.created]]) > 1)",
		},
		{
			"static @request fields with :lower modifier",
			"demo1",
			"@request.body.a:lower > true ||" +
				"@request.body.b:lower > true ||" +
				"@request.body.c:lower > true ||" +
				"@request.query.a:lower > true ||" +
				"@request.query.b:lower > true ||" +
				"@request.query.c:lower > true ||" +
				"@request.headers.a:lower > true ||" +
				"@request.headers.c:lower > true",
			false,
			"SELECT `demo1`.* FROM `demo1` WHERE (NULL > 1 OR LOWER({:TEST}) > 1 OR NULL > 1 OR LOWER({:TEST}) > 1 OR LOWER({:TEST}) > 1 OR NULL > 1 OR LOWER({:TEST}) > 1 OR NULL > 1)",
		},
		{
			"isset modifier",
			"demo1",
			"@request.body.a:isset > true ||" +
				"@request.body.b:isset > true ||" +
				"@request.body.c:isset > true ||" +
				"@request.query.a:isset > true ||" +
				"@request.query.b:isset > true ||" +
				"@request.query.c:isset > true ||" +
				"@request.headers.a:isset > true ||" +
				"@request.headers.c:isset > true",
			false,
			"SELECT `demo1`.* FROM `demo1` WHERE (TRUE > 1 OR TRUE > 1 OR FALSE > 1 OR TRUE > 1 OR TRUE > 1 OR FALSE > 1 OR TRUE > 1 OR FALSE > 1)",
		},
		{
			"@request.body.rel.* fields",
			"demo4",
			"@request.body.rel_one_cascade.title > true &&" +
				// reference the same as rel_one_cascade collection but should use a different join alias
				"@request.body.rel_one_no_cascade.title < true &&" +
				// different collection
				"@request.body.self_rel_many.title = true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo3` `__data_demo3_rel_one_cascade` ON [[__data_demo3_rel_one_cascade.id]]={:p0} LEFT JOIN `demo3` `__data_demo3_rel_one_no_cascade` ON [[__data_demo3_rel_one_no_cascade.id]]={:p1} LEFT JOIN `demo4` `__data_demo4_self_rel_many` ON [[__data_demo4_self_rel_many.id]]={:p2} WHERE ((({:TEST} IS NOT '' AND {:TEST} IS NOT {:TEST})) AND (({:TEST} IS NOT '' AND {:TEST} IS NOT {:TEST}))) AND (([[__data_demo3_rel_one_cascade.title]] > 1 AND [[__data_demo3_rel_one_no_cascade.title]] < 1 AND (([[__data_demo4_self_rel_many.title]] = 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___data_demo4_self_rel_many.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo4` `__mm___data_demo4_self_rel_many` ON [[__mm___data_demo4_self_rel_many.id]]={:p11} WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = 1))))))",
		},
		{
			"@request.body.arrayble:each fields",
			"demo1",
			"@request.body.select_one:each > true &&" +
				"@request.body.select_one:each ?< true &&" +
				"@request.body.select_many:each > true &&" +
				"@request.body.select_many:each ?< true &&" +
				"@request.body.file_one:each > true &&" +
				"@request.body.file_one:each ?< true &&" +
				"@request.body.file_many:each > true &&" +
				"@request.body.file_many:each ?< true &&" +
				"@request.body.rel_one:each > true &&" +
				"@request.body.rel_one:each ?< true &&" +
				"@request.body.rel_many:each > true &&" +
				"@request.body.rel_many:each ?< true",
			false,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_je_select_one` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_je_select_many` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_je_file_one` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_je_file_many` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_je_rel_one` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_je_rel_many` WHERE ([[__dataEach_je_select_one.value]] > 1 AND [[__dataEach_je_select_one.value]] < 1 AND (([[__dataEach_je_select_many.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___dataEach_je_select_many.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:mmdataEachTEST}) `__mm___dataEach_je_select_many` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) AND [[__dataEach_je_select_many.value]] < 1 AND [[__dataEach_je_file_one.value]] > 1 AND [[__dataEach_je_file_one.value]] < 1 AND (([[__dataEach_je_file_many.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___dataEach_je_file_many.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:mmdataEachTEST}) `__mm___dataEach_je_file_many` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) AND [[__dataEach_je_file_many.value]] < 1 AND [[__dataEach_je_rel_one.value]] > 1 AND [[__dataEach_je_rel_one.value]] < 1 AND (([[__dataEach_je_rel_many.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___dataEach_je_rel_many.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:mmdataEachTEST}) `__mm___dataEach_je_rel_many` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) AND [[__dataEach_je_rel_many.value]] < 1)",
		},
		{
			"regular arrayble:each fields",
			"view1",
			"select_one:each > true &&" +
				"select_one:each ?< true &&" +
				"select_many:each > true &&" +
				"select_many:each ?< true &&" +
				"file_one:each > true &&" +
				"file_one:each ?< true &&" +
				"file_many:each > true &&" +
				"file_many:each ?< true &&" +
				"rel_one:each > true &&" +
				"rel_one:each ?< true &&" +
				"rel_many:each > true &&" +
				"rel_many:each ?< true",
			false,
			"SELECT DISTINCT `view1`.* FROM `view1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[view1.select_one]]), json_type([[view1.select_one]])='array', FALSE) THEN [[view1.select_one]] ELSE json_array([[view1.select_one]]) END) `__je_view1_select_one` LEFT JOIN json_each(CASE WHEN iif(json_valid([[view1.select_many]]), json_type([[view1.select_many]])='array', FALSE) THEN [[view1.select_many]] ELSE json_array([[view1.select_many]]) END) `__je_view1_select_many` LEFT JOIN json_each(CASE WHEN iif(json_valid([[view1.file_one]]), json_type([[view1.file_one]])='array', FALSE) THEN [[view1.file_one]] ELSE json_array([[view1.file_one]]) END) `__je_view1_file_one` LEFT JOIN json_each(CASE WHEN iif(json_valid([[view1.file_many]]), json_type([[view1.file_many]])='array', FALSE) THEN [[view1.file_many]] ELSE json_array([[view1.file_many]]) END) `__je_view1_file_many` LEFT JOIN json_each(CASE WHEN iif(json_valid([[view1.rel_one]]), json_type([[view1.rel_one]])='array', FALSE) THEN [[view1.rel_one]] ELSE json_array([[view1.rel_one]]) END) `__je_view1_rel_one` LEFT JOIN json_each(CASE WHEN iif(json_valid([[view1.rel_many]]), json_type([[view1.rel_many]])='array', FALSE) THEN [[view1.rel_many]] ELSE json_array([[view1.rel_many]]) END) `__je_view1_rel_many` WHERE ([[__je_view1_select_one.value]] > 1 AND [[__je_view1_select_one.value]] < 1 AND (([[__je_view1_select_many.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__je___mm_view1_select_many.value]] as [[multiMatchValue]] FROM `view1` `__mm_view1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_view1.select_many]]), json_type([[__mm_view1.select_many]])='array', FALSE) THEN [[__mm_view1.select_many]] ELSE json_array([[__mm_view1.select_many]]) END) `__je___mm_view1_select_many` WHERE `__mm_view1`.`id` = `view1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) AND [[__je_view1_select_many.value]] < 1 AND [[__je_view1_file_one.value]] > 1 AND [[__je_view1_file_one.value]] < 1 AND (([[__je_view1_file_many.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__je___mm_view1_file_many.value]] as [[multiMatchValue]] FROM `view1` `__mm_view1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_view1.file_many]]), json_type([[__mm_view1.file_many]])='array', FALSE) THEN [[__mm_view1.file_many]] ELSE json_array([[__mm_view1.file_many]]) END) `__je___mm_view1_file_many` WHERE `__mm_view1`.`id` = `view1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) AND [[__je_view1_file_many.value]] < 1 AND [[__je_view1_rel_one.value]] > 1 AND [[__je_view1_rel_one.value]] < 1 AND (([[__je_view1_rel_many.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__je___mm_view1_rel_many.value]] as [[multiMatchValue]] FROM `view1` `__mm_view1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_view1.rel_many]]), json_type([[__mm_view1.rel_many]])='array', FALSE) THEN [[__mm_view1.rel_many]] ELSE json_array([[__mm_view1.rel_many]]) END) `__je___mm_view1_rel_many` WHERE `__mm_view1`.`id` = `view1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > 1)))) AND [[__je_view1_rel_many.value]] < 1)",
		},
		{
			"arrayble:each vs arrayble:each",
			"demo1",
			"select_one:each != select_many:each &&" +
				"select_many:each > select_one:each &&" +
				"select_many:each ?< select_one:each &&" +
				"select_many:each = @request.body.select_many:each",
			false,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo1.select_one]]), json_type([[demo1.select_one]])='array', FALSE) THEN [[demo1.select_one]] ELSE json_array([[demo1.select_one]]) END) `__je_demo1_select_one` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo1.select_many]]), json_type([[demo1.select_many]])='array', FALSE) THEN [[demo1.select_many]] ELSE json_array([[demo1.select_many]]) END) `__je_demo1_select_many` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_je_select_many` WHERE (((COALESCE([[__je_demo1_select_one.value]], '') IS NOT COALESCE([[__je_demo1_select_many.value]], '')) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__je___mm_demo1_select_many.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo1.select_many]]), json_type([[__mm_demo1.select_many]])='array', FALSE) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__je___mm_demo1_select_many` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT (COALESCE([[__je_demo1_select_one.value]], '') IS NOT COALESCE([[__smTEST.multiMatchValue]], ''))))) AND (([[__je_demo1_select_many.value]] > [[__je_demo1_select_one.value]]) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__je___mm_demo1_select_many.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo1.select_many]]), json_type([[__mm_demo1.select_many]])='array', FALSE) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__je___mm_demo1_select_many` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] > [[__je_demo1_select_one.value]])))) AND [[__je_demo1_select_many.value]] < [[__je_demo1_select_one.value]] AND (([[__je_demo1_select_many.value]] = [[__dataEach_je_select_many.value]]) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__je___mm_demo1_select_many.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo1.select_many]]), json_type([[__mm_demo1.select_many]])='array', FALSE) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__je___mm_demo1_select_many` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm___dataEach_je_select_many.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:mmdataEachTEST}) `__mm___dataEach_je_select_many` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE NOT (COALESCE([[__mlTEST.multiMatchValue]], '') = COALESCE([[__mrTEST.multiMatchValue]], ''))))))",
		},
		{
			"mixed multi-match vs multi-match in superuser only collections",
			"demo1",
			"rel_many.rel.active != rel_many.name &&" +
				"rel_many.rel.active ?= rel_many.name &&" +
				"rel_many.rel.title ~ rel_one.email &&" +
				"@collection.demo2.active = rel_many.rel.active &&" +
				"@collection.demo2.active ?= rel_many.rel.active &&" +
				"rel_many.verified > @request.body.rel_many.verified",
			false,
			"",
		},
		{
			"mixed multi-match vs multi-match in superuser only collections (with allowHiddenFields)",
			"demo1",
			"rel_many.rel.active != rel_many.name &&" +
				"rel_many.rel.active ?= rel_many.name &&" +
				"rel_many.rel.title ~ rel_one.email &&" +
				"@collection.demo2.active = rel_many.rel.active &&" +
				"@collection.demo2.active ?= rel_many.rel.active &&" +
				"rel_many.verified > @request.body.rel_many.verified",
			true,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo1.rel_many]]), json_type([[demo1.rel_many]])='array', FALSE) THEN [[demo1.rel_many]] ELSE json_array([[demo1.rel_many]]) END) `__je_demo1_rel_many` LEFT JOIN `users` `demo1_rel_many` ON [[demo1_rel_many.id]] = [[__je_demo1_rel_many.value]] LEFT JOIN `demo2` `demo1_rel_many_rel` ON [[demo1_rel_many_rel.id]] = [[demo1_rel_many.rel]] LEFT JOIN `demo1` `demo1_rel_one` ON [[demo1_rel_one.id]] = [[demo1.rel_one]] LEFT JOIN `demo2` `__collection_demo2` LEFT JOIN `users` `__data_users_rel_many` ON [[__data_users_rel_many.id]] IN ({:p0}, {:p1}) WHERE (((COALESCE([[demo1_rel_many_rel.active]], '') IS NOT COALESCE([[demo1_rel_many.name]], '')) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many_rel.active]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo1.rel_many]]), json_type([[__mm_demo1.rel_many]])='array', FALSE) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] LEFT JOIN `demo2` `__mm_demo1_rel_many_rel` ON [[__mm_demo1_rel_many_rel.id]] = [[__mm_demo1_rel_many.rel]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm_demo1_rel_many.name]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo1.rel_many]]), json_type([[__mm_demo1.rel_many]])='array', FALSE) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE NOT (COALESCE([[__mlTEST.multiMatchValue]], '') IS NOT COALESCE([[__mrTEST.multiMatchValue]], ''))))) AND COALESCE([[demo1_rel_many_rel.active]], '') = COALESCE([[demo1_rel_many.name]], '') AND (([[demo1_rel_many_rel.title]] LIKE ('%' || [[demo1_rel_one.email]] || '%') ESCAPE '\\') AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many_rel.title]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo1.rel_many]]), json_type([[__mm_demo1.rel_many]])='array', FALSE) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] LEFT JOIN `demo2` `__mm_demo1_rel_many_rel` ON [[__mm_demo1_rel_many_rel.id]] = [[__mm_demo1_rel_many.rel]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] LIKE ('%' || [[demo1_rel_one.email]] || '%') ESCAPE '\\')))) AND ((COALESCE([[__collection_demo2.active]], '') = COALESCE([[demo1_rel_many_rel.active]], '')) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm___collection_demo2.active]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN `demo2` `__mm___collection_demo2` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm_demo1_rel_many_rel.active]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo1.rel_many]]), json_type([[__mm_demo1.rel_many]])='array', FALSE) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] LEFT JOIN `demo2` `__mm_demo1_rel_many_rel` ON [[__mm_demo1_rel_many_rel.id]] = [[__mm_demo1_rel_many.rel]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE NOT (COALESCE([[__mlTEST.multiMatchValue]], '') = COALESCE([[__mrTEST.multiMatchValue]], ''))))) AND COALESCE([[__collection_demo2.active]], '') = COALESCE([[demo1_rel_many_rel.active]], '') AND (([[demo1_rel_many.verified]] > [[__data_users_rel_many.verified]]) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many.verified]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo1.rel_many]]), json_type([[__mm_demo1.rel_many]])='array', FALSE) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm___data_users_rel_many.verified]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN `users` `__mm___data_users_rel_many` ON [[__mm___data_users_rel_many.id]] IN ({:p2}, {:p3}) WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE NOT ([[__mlTEST.multiMatchValue]] > [[__mrTEST.multiMatchValue]])))))",
		},
		{
			"@request.body.arrayable:length fields",
			"demo1",
			"@request.body.select_one:length > 1 &&" +
				"@request.body.select_one:length ?> 2 &&" +
				"@request.body.select_many:length < 3 &&" +
				"@request.body.select_many:length ?> 4 &&" +
				"@request.body.rel_one:length = 5 &&" +
				"@request.body.rel_one:length ?= 6 &&" +
				"@request.body.rel_many:length != 7 &&" +
				"@request.body.rel_many:length ?!= 8 &&" +
				"@request.body.file_one:length = 9 &&" +
				"@request.body.file_one:length ?= 0 &&" +
				"@request.body.file_many:length != 1 &&" +
				"@request.body.file_many:length ?!= 2",
			false,
			"SELECT `demo1`.* FROM `demo1` WHERE (0 > {:TEST} AND 0 > {:TEST} AND 2 < {:TEST} AND 2 > {:TEST} AND 1 = {:TEST} AND 1 = {:TEST} AND 2 IS NOT {:TEST} AND 2 IS NOT {:TEST} AND 1 = {:TEST} AND 1 = {:TEST} AND 3 IS NOT {:TEST} AND 3 IS NOT {:TEST})",
		},
		{
			"regular arrayable:length fields",
			"demo4",
			"@request.body.self_rel_one.self_rel_many:length > 1 &&" +
				"@request.body.self_rel_one.self_rel_many:length ?> 2 &&" +
				"@request.body.rel_many_cascade.files:length ?< 3 &&" +
				"@request.body.rel_many_cascade.files:length < 4 &&" +
				"@request.body.rel_one_cascade.files:length < 4.1 &&" + // to ensure that the join to the same as above table will be aliased
				"self_rel_one.self_rel_many:length = 5 &&" +
				"self_rel_one.self_rel_many:length ?= 6 &&" +
				"self_rel_one.rel_many_cascade.files:length != 7 &&" +
				"self_rel_one.rel_many_cascade.files:length ?!= 8",
			true,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `__data_demo4_self_rel_one` ON [[__data_demo4_self_rel_one.id]]={:p0} LEFT JOIN `demo3` `__data_demo3_rel_many_cascade` ON [[__data_demo3_rel_many_cascade.id]] IN ({:p1}, {:p2}) LEFT JOIN `demo3` `__data_demo3_rel_one_cascade` ON [[__data_demo3_rel_one_cascade.id]]={:p3} LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4_self_rel_one.rel_many_cascade]]), json_type([[demo4_self_rel_one.rel_many_cascade]])='array', FALSE) THEN [[demo4_self_rel_one.rel_many_cascade]] ELSE json_array([[demo4_self_rel_one.rel_many_cascade]]) END) `__je_demo4_self_rel_one_rel_many_cascade` LEFT JOIN `demo3` `demo4_self_rel_one_rel_many_cascade` ON [[demo4_self_rel_one_rel_many_cascade.id]] = [[__je_demo4_self_rel_one_rel_many_cascade.value]] WHERE (json_array_length(CASE WHEN iif(json_valid([[__data_demo4_self_rel_one.self_rel_many]]), json_type([[__data_demo4_self_rel_one.self_rel_many]])='array', FALSE) THEN [[__data_demo4_self_rel_one.self_rel_many]] ELSE (CASE WHEN [[__data_demo4_self_rel_one.self_rel_many]] = '' OR [[__data_demo4_self_rel_one.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[__data_demo4_self_rel_one.self_rel_many]]) END) END) > {:TEST} AND json_array_length(CASE WHEN iif(json_valid([[__data_demo4_self_rel_one.self_rel_many]]), json_type([[__data_demo4_self_rel_one.self_rel_many]])='array', FALSE) THEN [[__data_demo4_self_rel_one.self_rel_many]] ELSE (CASE WHEN [[__data_demo4_self_rel_one.self_rel_many]] = '' OR [[__data_demo4_self_rel_one.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[__data_demo4_self_rel_one.self_rel_many]]) END) END) > {:TEST} AND json_array_length(CASE WHEN iif(json_valid([[__data_demo3_rel_many_cascade.files]]), json_type([[__data_demo3_rel_many_cascade.files]])='array', FALSE) THEN [[__data_demo3_rel_many_cascade.files]] ELSE (CASE WHEN [[__data_demo3_rel_many_cascade.files]] = '' OR [[__data_demo3_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__data_demo3_rel_many_cascade.files]]) END) END) < {:TEST} AND ((json_array_length(CASE WHEN iif(json_valid([[__data_demo3_rel_many_cascade.files]]), json_type([[__data_demo3_rel_many_cascade.files]])='array', FALSE) THEN [[__data_demo3_rel_many_cascade.files]] ELSE (CASE WHEN [[__data_demo3_rel_many_cascade.files]] = '' OR [[__data_demo3_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__data_demo3_rel_many_cascade.files]]) END) END) < {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT json_array_length(CASE WHEN iif(json_valid([[__mm___data_demo3_rel_many_cascade.files]]), json_type([[__mm___data_demo3_rel_many_cascade.files]])='array', FALSE) THEN [[__mm___data_demo3_rel_many_cascade.files]] ELSE (CASE WHEN [[__mm___data_demo3_rel_many_cascade.files]] = '' OR [[__mm___data_demo3_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__mm___data_demo3_rel_many_cascade.files]]) END) END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo3` `__mm___data_demo3_rel_many_cascade` ON [[__mm___data_demo3_rel_many_cascade.id]] IN ({:p8}, {:p9}) WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] < {:TEST})))) AND json_array_length(CASE WHEN iif(json_valid([[__data_demo3_rel_one_cascade.files]]), json_type([[__data_demo3_rel_one_cascade.files]])='array', FALSE) THEN [[__data_demo3_rel_one_cascade.files]] ELSE (CASE WHEN [[__data_demo3_rel_one_cascade.files]] = '' OR [[__data_demo3_rel_one_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__data_demo3_rel_one_cascade.files]]) END) END) < {:TEST} AND json_array_length(CASE WHEN iif(json_valid([[demo4_self_rel_one.self_rel_many]]), json_type([[demo4_self_rel_one.self_rel_many]])='array', FALSE) THEN [[demo4_self_rel_one.self_rel_many]] ELSE (CASE WHEN [[demo4_self_rel_one.self_rel_many]] = '' OR [[demo4_self_rel_one.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[demo4_self_rel_one.self_rel_many]]) END) END) = {:TEST} AND json_array_length(CASE WHEN iif(json_valid([[demo4_self_rel_one.self_rel_many]]), json_type([[demo4_self_rel_one.self_rel_many]])='array', FALSE) THEN [[demo4_self_rel_one.self_rel_many]] ELSE (CASE WHEN [[demo4_self_rel_one.self_rel_many]] = '' OR [[demo4_self_rel_one.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[demo4_self_rel_one.self_rel_many]]) END) END) = {:TEST} AND ((json_array_length(CASE WHEN iif(json_valid([[demo4_self_rel_one_rel_many_cascade.files]]), json_type([[demo4_self_rel_one_rel_many_cascade.files]])='array', FALSE) THEN [[demo4_self_rel_one_rel_many_cascade.files]] ELSE (CASE WHEN [[demo4_self_rel_one_rel_many_cascade.files]] = '' OR [[demo4_self_rel_one_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[demo4_self_rel_one_rel_many_cascade.files]]) END) END) IS NOT {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT json_array_length(CASE WHEN iif(json_valid([[__mm_demo4_self_rel_one_rel_many_cascade.files]]), json_type([[__mm_demo4_self_rel_one_rel_many_cascade.files]])='array', FALSE) THEN [[__mm_demo4_self_rel_one_rel_many_cascade.files]] ELSE (CASE WHEN [[__mm_demo4_self_rel_one_rel_many_cascade.files]] = '' OR [[__mm_demo4_self_rel_one_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__mm_demo4_self_rel_one_rel_many_cascade.files]]) END) END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo4` `__mm_demo4_self_rel_one` ON [[__mm_demo4_self_rel_one.id]] = [[__mm_demo4.self_rel_one]] LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4_self_rel_one.rel_many_cascade]]), json_type([[__mm_demo4_self_rel_one.rel_many_cascade]])='array', FALSE) THEN [[__mm_demo4_self_rel_one.rel_many_cascade]] ELSE json_array([[__mm_demo4_self_rel_one.rel_many_cascade]]) END) `__mm_demo4_self_rel_one_rel_many_cascade_je` LEFT JOIN `demo3` `__mm_demo4_self_rel_one_rel_many_cascade` ON [[__mm_demo4_self_rel_one_rel_many_cascade.id]] = [[__mm_demo4_self_rel_one_rel_many_cascade_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] IS NOT {:TEST})))) AND json_array_length(CASE WHEN iif(json_valid([[demo4_self_rel_one_rel_many_cascade.files]]), json_type([[demo4_self_rel_one_rel_many_cascade.files]])='array', FALSE) THEN [[demo4_self_rel_one_rel_many_cascade.files]] ELSE (CASE WHEN [[demo4_self_rel_one_rel_many_cascade.files]] = '' OR [[demo4_self_rel_one_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[demo4_self_rel_one_rel_many_cascade.files]]) END) END) IS NOT {:TEST})",
		},
		{
			"request body :changed modifier with non-existing collection field",
			"demo1",
			"@request.body.a:changed > 1",
			true,
			"",
		},
		{
			"regular body :changed modifier",
			"demo1",
			"@request.body.number:changed = false &&" +
				"@request.body.email:changed = true &&" +
				"@request.body.number:changed = @request.body.select_many:changed",
			true,
			"SELECT `demo1`.* FROM `demo1` WHERE ((TRUE = 1 AND {:TEST} IS NOT [[demo1.number]]) IS 0 AND (FALSE = 1 AND ('' IS NOT [[demo1.email]] AND [[demo1.email]] IS NOT NULL)) IS 1 AND (TRUE = 1 AND {:TEST} IS NOT [[demo1.number]]) IS (TRUE = 1 AND {:TEST} IS NOT [[demo1.select_many]]))",
		},
		{
			"json_extract and json_array_length COALESCE equal normalizations",
			"demo4",
			"json_object.a.b = '' && self_rel_many:length != 2 && json_object.a.b > 3 && self_rel_many:length <= 4",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE ((CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$.a.b') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb.a.b') END) IS {:TEST} AND json_array_length(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE (CASE WHEN [[demo4.self_rel_many]] = '' OR [[demo4.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[demo4.self_rel_many]]) END) END) IS NOT {:TEST} AND (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$.a.b') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb.a.b') END) > {:TEST} AND json_array_length(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE (CASE WHEN [[demo4.self_rel_many]] = '' OR [[demo4.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[demo4.self_rel_many]]) END) END) <= {:TEST})",
		},
		{
			"json field equal normalization checks",
			"demo4",
			"json_object = '' || json_object != '' || '' = json_object || '' != json_object ||" +
				"json_object = null || json_object != null || null = json_object || null != json_object ||" +
				"json_object = true || json_object != true || true = json_object || true != json_object ||" +
				"json_object = json_object || json_object != json_object ||" +
				"json_object = title || title != json_object ||" +
				// multimatch expressions
				"self_rel_many.json_object = '' || null = self_rel_many.json_object ||" +
				"self_rel_many.json_object = self_rel_many.json_object",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo4.self_rel_many]]), json_type([[demo4.self_rel_many]])='array', FALSE) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `__je_demo4_self_rel_many` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[__je_demo4_self_rel_many.value]] WHERE ((CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS {:TEST} OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NOT {:TEST} OR {:TEST} IS (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR {:TEST} IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NULL OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NOT NULL OR NULL IS (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR NULL IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS 1 OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NOT 1 OR 1 IS (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR 1 IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS [[demo4.title]] OR [[demo4.title]] IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (((CASE WHEN json_valid([[demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_many.json_object]]), '$.pb') END) IS {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT (CASE WHEN json_valid([[__mm_demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[__mm_demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[__mm_demo4_self_rel_many.json_object]]), '$.pb') END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4.self_rel_many]]), json_type([[__mm_demo4.self_rel_many]])='array', FALSE) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] IS {:TEST})))) OR ((NULL IS (CASE WHEN json_valid([[demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_many.json_object]]), '$.pb') END)) AND (NOT EXISTS (SELECT 1 FROM (SELECT (CASE WHEN json_valid([[__mm_demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[__mm_demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[__mm_demo4_self_rel_many.json_object]]), '$.pb') END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4.self_rel_many]]), json_type([[__mm_demo4.self_rel_many]])='array', FALSE) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT (NULL IS [[__smTEST.multiMatchValue]])))) OR (((CASE WHEN json_valid([[demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_many.json_object]]), '$.pb') END) IS (CASE WHEN json_valid([[demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_many.json_object]]), '$.pb') END)) AND (NOT EXISTS (SELECT 1 FROM (SELECT (CASE WHEN json_valid([[__mm_demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[__mm_demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[__mm_demo4_self_rel_many.json_object]]), '$.pb') END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4.self_rel_many]]), json_type([[__mm_demo4.self_rel_many]])='array', FALSE) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__mlTEST}} LEFT JOIN (SELECT (CASE WHEN json_valid([[__mm_demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[__mm_demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[__mm_demo4_self_rel_many.json_object]]), '$.pb') END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo4.self_rel_many]]), json_type([[__mm_demo4.self_rel_many]])='array', FALSE) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__mrTEST}} WHERE NOT ([[__mlTEST.multiMatchValue]] IS [[__mrTEST.multiMatchValue]])))))",
		},
		{
			"geoPoint props access",
			"view1",
			"point = '' || point.lat > 1 || point.lon < 2 || point.something > 3",
			false,
			"SELECT `view1`.* FROM `view1` WHERE (([[view1.point]] = '' OR [[view1.point]] IS NULL) OR (CASE WHEN json_valid([[view1.point]]) THEN JSON_EXTRACT([[view1.point]], '$.lat') ELSE JSON_EXTRACT(json_object('pb', [[view1.point]]), '$.pb.lat') END) > {:TEST} OR (CASE WHEN json_valid([[view1.point]]) THEN JSON_EXTRACT([[view1.point]], '$.lon') ELSE JSON_EXTRACT(json_object('pb', [[view1.point]]), '$.pb.lon') END) < {:TEST} OR (CASE WHEN json_valid([[view1.point]]) THEN JSON_EXTRACT([[view1.point]], '$.something') ELSE JSON_EXTRACT(json_object('pb', [[view1.point]]), '$.pb.something') END) > {:TEST})",
		},
		{
			"strftime with fixed string as time-value against known empty value (null normalizations)",
			"demo5",
			"strftime('%Y-%m', '2026-01-01') = ''",
			false,
			"SELECT `demo5`.* FROM `demo5` WHERE ((strftime({:TEST},{:TEST}) = '' OR strftime({:TEST},{:TEST}) IS NULL))",
		},
		{
			"strftime without multi-match",
			"demo5",
			"strftime('%Y-%m', rel_one.created) = true",
			false,
			"SELECT DISTINCT `demo5`.* FROM `demo5` LEFT JOIN `demo4` `demo5_rel_one` ON [[demo5_rel_one.id]] = [[demo5.rel_one]] WHERE strftime({:TEST},[[demo5_rel_one.created]]) = 1",
		},
		{
			"strftime with multi-match",
			"demo5",
			"strftime('%Y-%m', rel_many.created) = true",
			false,
			"SELECT DISTINCT `demo5`.* FROM `demo5` LEFT JOIN json_each(CASE WHEN iif(json_valid([[demo5.rel_many]]), json_type([[demo5.rel_many]])='array', FALSE) THEN [[demo5.rel_many]] ELSE json_array([[demo5.rel_many]]) END) `__je_demo5_rel_many` LEFT JOIN `demo4` `demo5_rel_many` ON [[demo5_rel_many.id]] = [[__je_demo5_rel_many.value]] WHERE (((strftime({:TEST},[[demo5_rel_many.created]]) = 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT strftime({:TEST},[[__mm_demo5_rel_many.created]]) as [[multiMatchValue]] FROM `demo5` `__mm_demo5` LEFT JOIN json_each(CASE WHEN iif(json_valid([[__mm_demo5.rel_many]]), json_type([[__mm_demo5.rel_many]])='array', FALSE) THEN [[__mm_demo5.rel_many]] ELSE json_array([[__mm_demo5.rel_many]]) END) `__mm_demo5_rel_many_je` LEFT JOIN `demo4` `__mm_demo5_rel_many` ON [[__mm_demo5_rel_many.id]] = [[__mm_demo5_rel_many_je.value]] WHERE `__mm_demo5`.`id` = `demo5`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = 1)))))",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			collection, err := app.FindCollectionByNameOrId(s.collectionIdOrName)
			if err != nil {
				t.Fatalf("Failed to load collection %s: %v", s.collectionIdOrName, err)
			}

			expectError := s.expectQuery == ""

			query := app.RecordQuery(collection)

			r := core.NewRecordFieldResolver(app, collection, requestInfo, s.allowHiddenFields)

			expr, err := search.FilterData(s.rule).BuildExpr(r)
			hasErr := err != nil
			if hasErr != expectError {
				t.Fatalf("BuildExpr failed: expected hasErr %v, got %v (%v)", expectError, hasErr, err)
			}

			err = r.UpdateQuery(query)
			hasErr = err != nil
			if hasErr && expectError {
				t.Fatalf("UpdateQuery failed: expected hasErr %v, got %v (%v)", expectError, hasErr, err)
			}

			if expectError {
				return
			}

			rawQuery := query.AndWhere(expr).Build().SQL()

			// replace TEST placeholder with .+ regex pattern
			expectQuery := strings.ReplaceAll(
				"^"+regexp.QuoteMeta(s.expectQuery)+"$",
				"TEST",
				`\w+`,
			)

			if !list.ExistInSliceWithRegex(rawQuery, []string{expectQuery}) {
				t.Fatalf("Expected query\n %v \ngot:\n %v", expectQuery, rawQuery)
			}
		})
	}
}

func TestRecordFieldResolverResolveCollectionFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.FindCollectionByNameOrId("demo4")
	if err != nil {
		t.Fatal(err)
	}

	authRecord, err := app.FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestInfo := &core.RequestInfo{
		Auth: authRecord,
	}

	r := core.NewRecordFieldResolver(app, collection, requestInfo, true)

	scenarios := []struct {
		fieldName   string
		expectError bool
		expectName  string
	}{
		{"", true, ""},
		{" ", true, ""},
		{"unknown", true, ""},
		{"invalid format", true, ""},
		{"id", false, "[[demo4.id]]"},
		{"created", false, "[[demo4.created]]"},
		{"updated", false, "[[demo4.updated]]"},
		{"title", false, "[[demo4.title]]"},
		{"title.test", true, ""},
		{"self_rel_many", false, "[[demo4.self_rel_many]]"},
		{"self_rel_many.", true, ""},
		{"self_rel_many.unknown", true, ""},
		{"self_rel_many.title", false, "[[demo4_self_rel_many.title]]"},
		{"self_rel_many.self_rel_one.self_rel_many.title", false, "[[demo4_self_rel_many_self_rel_one_self_rel_many.title]]"},

		// max relations limit
		{"self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.id", false, "[[demo4_self_rel_many_self_rel_many_self_rel_many_self_rel_many_self_rel_many_self_rel_many.id]]"},
		{"self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.self_rel_many.id", true, ""},

		// back relations
		{"rel_one_cascade.demo4_via_title.id", true, ""},        // not a relation field
		{"rel_one_cascade.demo4_via_self_rel_one.id", true, ""}, // relation field but to a different collection
		{"rel_one_cascade.demo4_via_rel_one_cascade.id", false, "[[demo4_rel_one_cascade_demo4_via_rel_one_cascade.id]]"},
		{"rel_one_cascade.demo4_via_rel_one_cascade.rel_one_cascade.demo4_via_rel_one_cascade.id", false, "[[demo4_rel_one_cascade_demo4_via_rel_one_cascade_rel_one_cascade_demo4_via_rel_one_cascade.id]]"},

		// json_extract
		{"json_array.0", false, "(CASE WHEN json_valid([[demo4.json_array]]) THEN JSON_EXTRACT([[demo4.json_array]], '$[0]') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_array]]), '$.pb[0]') END)"},
		{"json_object.a.b.c", false, "(CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$.a.b.c') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb.a.b.c') END)"},

		// max relations limit shouldn't apply for json paths
		{"json_object.a.b.c.e.f.g.h.i.j.k.l.m.n.o.p", false, "(CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$.a.b.c.e.f.g.h.i.j.k.l.m.n.o.p') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb.a.b.c.e.f.g.h.i.j.k.l.m.n.o.p') END)"},

		// @request.auth relation join
		{"@request.auth.rel", false, "[[__auth_users.rel]]"},
		{"@request.auth.rel.title", false, "[[__auth_users_rel.title]]"},
		{"@request.auth.demo1_via_rel_many.id", false, "[[__auth_users_demo1_via_rel_many.id]]"},
		{"@request.auth.rel.missing", false, "NULL"},
		{"@request.auth.missing_via_rel", false, "NULL"},
		{"@request.auth.demo1_via_file_one.id", false, "NULL"}, // not a relation field
		{"@request.auth.demo1_via_rel_one.id", false, "NULL"},  // relation field but to a different collection

		// @collection fieds
		{"@collect", true, ""},
		{"collection.demo4.title", true, ""},
		{"@collection", true, ""},
		{"@collection.unknown", true, ""},
		{"@collection.demo2", true, ""},
		{"@collection.demo2.", true, ""},
		{"@collection.demo2:someAlias", true, ""},
		{"@collection.demo2:someAlias.", true, ""},
		{"@collection.demo2.title", false, "[[__collection_demo2.title]]"},
		{"@collection.demo2:someAlias.title", false, "[[__collection_alias_someAlias.title]]"},
		{"@collection.demo4.id", false, "[[__collection_demo4.id]]"},
		{"@collection.demo4.created", false, "[[__collection_demo4.created]]"},
		{"@collection.demo4.updated", false, "[[__collection_demo4.updated]]"},
		{"@collection.demo4.self_rel_many.missing", true, ""},
		{"@collection.demo4.self_rel_many.self_rel_one.self_rel_many.self_rel_one.title", false, "[[__collection_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]]"},
	}

	for _, s := range scenarios {
		t.Run(s.fieldName, func(t *testing.T) {
			r, err := r.Resolve(s.fieldName)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if r.Identifier != s.expectName {
				t.Fatalf("Expected r.Identifier\n%q\ngot\n%q", s.expectName, r.Identifier)
			}

			// params should be empty for non @request fields
			if len(r.Params) != 0 {
				t.Fatalf("Expected 0 r.Params, got\n%v", r.Params)
			}
		})
	}
}

func TestRecordFieldResolverResolveStaticRequestInfoFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	authRecord, err := app.FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestInfo := &core.RequestInfo{
		Context: "ctx",
		Method:  "get",
		Query: map[string]string{
			"a": "123",
		},
		Body: map[string]any{
			"number":          "10",
			"number_unknown":  "20",
			"raw_json_obj":    types.JSONRaw(`{"a":123}`),
			"raw_json_arr1":   types.JSONRaw(`[123, 456]`),
			"raw_json_arr2":   types.JSONRaw(`[{"a":123},{"b":456}]`),
			"raw_json_simple": types.JSONRaw(`123`),
			"b":               456,
			"c":               map[string]int{"sub": 1},
		},
		Headers: map[string]string{
			"d": "789",
		},
		Auth: authRecord,
	}

	r := core.NewRecordFieldResolver(app, collection, requestInfo, true)

	scenarios := []struct {
		fieldName        string
		expectError      bool
		expectParamValue string // encoded json
	}{
		{"@request", true, ""},
		{"@request.invalid format", true, ""},
		{"@request.invalid_format2!", true, ""},
		{"@request.missing", true, ""},
		{"@request.context", false, `"ctx"`},
		{"@request.method", false, `"get"`},
		{"@request.query", true, ``},
		{"@request.query.a", false, `"123"`},
		{"@request.query.a.missing", false, ``},
		{"@request.headers", true, ``},
		{"@request.headers.missing", false, ``},
		{"@request.headers.d", false, `"789"`},
		{"@request.headers.d.sub", false, ``},
		{"@request.body", true, ``},
		{"@request.body.b", false, `456`},
		{"@request.body.number", false, `10`},           // number field normalization
		{"@request.body.number_unknown", false, `"20"`}, // no numeric normalizations for unknown fields
		{"@request.body.b.missing", false, ``},
		{"@request.body.c", false, `"{\"sub\":1}"`},
		{"@request.auth", true, ""},
		{"@request.auth.id", false, `"4q1xlclmfloku33"`},
		{"@request.auth.collectionId", false, `"` + authRecord.Collection().Id + `"`},
		{"@request.auth.collectionName", false, `"` + authRecord.Collection().Name + `"`},
		{"@request.auth.verified", false, `false`},
		{"@request.auth.emailVisibility", false, `false`},
		{"@request.auth.email", false, `"test@example.com"`}, // should always be returned no matter of the emailVisibility state
		{"@request.auth.missing", false, `NULL`},
		{"@request.body.raw_json_simple", false, `"123"`},
		{"@request.body.raw_json_simple.a", false, `NULL`},
		{"@request.body.raw_json_obj.a", false, `123`},
		{"@request.body.raw_json_obj.b", false, `NULL`},
		{"@request.body.raw_json_arr1.1", false, `456`},
		{"@request.body.raw_json_arr1.3", false, `NULL`},
		{"@request.body.raw_json_arr2.0.a", false, `123`},
		{"@request.body.raw_json_arr2.0.b", false, `NULL`},
	}

	for _, s := range scenarios {
		t.Run(s.fieldName, func(t *testing.T) {
			r, err := r.Resolve(s.fieldName)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			// missing key
			// ---
			if len(r.Params) == 0 {
				if r.Identifier != "NULL" {
					t.Fatalf("Expected 0 placeholder parameters for %v, got %v", r.Identifier, r.Params)
				}
				return
			}

			// existing key
			// ---
			if len(r.Params) != 1 {
				t.Fatalf("Expected 1 placeholder parameter for %v, got %v", r.Identifier, r.Params)
			}

			var paramName string
			var paramValue any
			for k, v := range r.Params {
				paramName = k
				paramValue = v
			}

			if r.Identifier != ("{:" + paramName + "}") {
				t.Fatalf("Expected parameter r.Identifier %q, got %q", paramName, r.Identifier)
			}

			encodedParamValue, _ := json.Marshal(paramValue)
			if string(encodedParamValue) != s.expectParamValue {
				t.Fatalf("Expected r.Params %#v for %s, got %#v", s.expectParamValue, r.Identifier, string(encodedParamValue))
			}
		})
	}

	// ensure that the original email visibility was restored
	if authRecord.EmailVisibility() {
		t.Fatal("Expected the original authRecord emailVisibility to remain unchanged")
	}
	if v, ok := authRecord.PublicExport()[core.FieldNameEmail]; ok {
		t.Fatalf("Expected the original authRecord email to not be exported, got %q", v)
	}
}
