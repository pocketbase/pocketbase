package resolvers_test

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/thinkonmay/pocketbase/models"
	"github.com/thinkonmay/pocketbase/models/schema"
	"github.com/thinkonmay/pocketbase/resolvers"
	"github.com/thinkonmay/pocketbase/tests"
	"github.com/thinkonmay/pocketbase/tools/list"
	"github.com/thinkonmay/pocketbase/tools/search"
)

func TestRecordFieldResolverUpdateQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	authRecord, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestInfo := &models.RequestInfo{
		Context: "ctx",
		Headers: map[string]any{
			"a": "123",
			"b": "456",
		},
		Query: map[string]any{
			"a": nil,
			"b": 123,
		},
		Data: map[string]any{
			"a":                  nil,
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
		AuthRecord: authRecord,
	}

	scenarios := []struct {
		name               string
		collectionIdOrName string
		rule               string
		allowHiddenFields  bool
		expectQuery        string
	}{
		{
			"non relation field (with all default operators)",
			"demo4",
			"title = true || title != 'test' || title ~ 'test1' || title !~ '%test2' || title > 1 || title >= 2 || title < 3 || title <= 4",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE ([[demo4.title]] = 1 OR [[demo4.title]] IS NOT {:TEST} OR [[demo4.title]] LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] NOT LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] > {:TEST} OR [[demo4.title]] >= {:TEST} OR [[demo4.title]] < {:TEST} OR [[demo4.title]] <= {:TEST})",
		},
		{
			"non relation field (with all opt/any operators)",
			"demo4",
			"title ?= true || title ?!= 'test' || title ?~ 'test1' || title ?!~ '%test2' || title ?> 1 || title ?>= 2 || title ?< 3 || title ?<= 4",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE ([[demo4.title]] = 1 OR [[demo4.title]] IS NOT {:TEST} OR [[demo4.title]] LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] NOT LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] > {:TEST} OR [[demo4.title]] >= {:TEST} OR [[demo4.title]] < {:TEST} OR [[demo4.title]] <= {:TEST})",
		},
		{
			"incomplete rel",
			"demo4",
			"self_rel_one > true",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE [[demo4.self_rel_one]] > 1",
		},
		{
			"single rel (self rel)",
			"demo4",
			"self_rel_one.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] WHERE [[demo4_self_rel_one.title]] > 1",
		},
		{
			"single rel (other collection)",
			"demo4",
			"rel_one_cascade.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo3` `demo4_rel_one_cascade` ON [[demo4_rel_one_cascade.id]] = [[demo4.rel_one_cascade]] WHERE [[demo4_rel_one_cascade.title]] > 1",
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
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] WHERE [[demo4_self_rel_many.self_rel_one]] > 1",
		},
		{
			"nested incomplete relations (multi-match operator)",
			"demo4",
			"self_rel_many.self_rel_one > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] WHERE ((([[demo4_self_rel_many.self_rel_one]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many.self_rel_one]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{TEST}} WHERE ((NOT ([[TEST.multiMatchValue]] > 1)) OR ([[TEST.multiMatchValue]] IS NULL))))))",
		},
		{
			"nested complete relations (opt/any operator)",
			"demo4",
			"self_rel_many.self_rel_one.title ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] WHERE [[demo4_self_rel_many_self_rel_one.title]] > 1",
		},
		{
			"nested complete relations (multi-match operator)",
			"demo4",
			"self_rel_many.self_rel_one.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] WHERE ((([[demo4_self_rel_many_self_rel_one.title]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many_self_rel_one.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one` ON [[__mm_demo4_self_rel_many_self_rel_one.id]] = [[__mm_demo4_self_rel_many.self_rel_one]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))))",
		},
		{
			"repeated nested relations (opt/any operator)",
			"demo4",
			"self_rel_many.self_rel_one.self_rel_many.self_rel_one.title ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_many_self_rel_one.self_rel_many]]) THEN [[demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `demo4_self_rel_many_self_rel_one_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many` ON [[demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] WHERE [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]] > 1",
		},
		{
			"repeated nested relations (multi-match operator)",
			"demo4",
			"self_rel_many.self_rel_one.self_rel_many.self_rel_one.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_many_self_rel_one.self_rel_many]]) THEN [[demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `demo4_self_rel_many_self_rel_one_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many` ON [[demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] WHERE ((([[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one` ON [[__mm_demo4_self_rel_many_self_rel_one.id]] = [[__mm_demo4_self_rel_many.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]]) THEN [[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `__mm_demo4_self_rel_many_self_rel_one_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one_self_rel_many` ON [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many_je.value]] LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))))",
		},
		{
			"multiple relations (opt/any operators)",
			"demo4",
			"self_rel_many.title ?= 'test' || self_rel_one.json_object.a ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] WHERE ([[demo4_self_rel_many.title]] = {:TEST} OR (CASE WHEN json_valid([[demo4_self_rel_one.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_one.json_object]], '$.a') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_one.json_object]]), '$.pb.a') END) > 1)",
		},
		{
			"multiple relations (multi-match operators)",
			"demo4",
			"self_rel_many.title = 'test' || self_rel_one.json_object.a > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] WHERE ((([[demo4_self_rel_many.title]] = {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = {:TEST})))) OR (CASE WHEN json_valid([[demo4_self_rel_one.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_one.json_object]], '$.a') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_one.json_object]]), '$.pb.a') END) > 1)",
		},
		{
			"back relations via single relation field (without unique index)",
			"demo3",
			"demo4_via_rel_one_cascade.id = true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_one_cascade` ON [[demo3.id]] IN (SELECT [[demo3_demo4_via_rel_one_cascade_je.value]] FROM json_each(CASE WHEN json_valid([[demo3_demo4_via_rel_one_cascade.rel_one_cascade]]) THEN [[demo3_demo4_via_rel_one_cascade.rel_one_cascade]] ELSE json_array([[demo3_demo4_via_rel_one_cascade.rel_one_cascade]]) END) {{demo3_demo4_via_rel_one_cascade_je}}) WHERE ((([[demo3_demo4_via_rel_one_cascade.id]] = 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo3_demo4_via_rel_one_cascade.id]] as [[multiMatchValue]] FROM `demo3` `__mm_demo3` LEFT JOIN `demo4` `__mm_demo3_demo4_via_rel_one_cascade` ON [[__mm_demo3.id]] IN (SELECT [[__mm_demo3_demo4_via_rel_one_cascade_je.value]] FROM json_each(CASE WHEN json_valid([[__mm_demo3_demo4_via_rel_one_cascade.rel_one_cascade]]) THEN [[__mm_demo3_demo4_via_rel_one_cascade.rel_one_cascade]] ELSE json_array([[__mm_demo3_demo4_via_rel_one_cascade.rel_one_cascade]]) END) {{__mm_demo3_demo4_via_rel_one_cascade_je}}) WHERE `__mm_demo3`.`id` = `demo3`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = 1)))))",
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
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade` ON [[demo3.id]] IN (SELECT [[demo3_demo4_via_rel_many_cascade_je.value]] FROM json_each(CASE WHEN json_valid([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) THEN [[demo3_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{demo3_demo4_via_rel_many_cascade_je}}) WHERE [[demo3_demo4_via_rel_many_cascade.id]] = 1",
		},
		{
			"back relations via multiple relation field (multi-match operators)",
			"demo3",
			"demo4_via_rel_many_cascade.id = true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade` ON [[demo3.id]] IN (SELECT [[demo3_demo4_via_rel_many_cascade_je.value]] FROM json_each(CASE WHEN json_valid([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) THEN [[demo3_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{demo3_demo4_via_rel_many_cascade_je}}) WHERE ((([[demo3_demo4_via_rel_many_cascade.id]] = 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo3_demo4_via_rel_many_cascade.id]] as [[multiMatchValue]] FROM `demo3` `__mm_demo3` LEFT JOIN `demo4` `__mm_demo3_demo4_via_rel_many_cascade` ON [[__mm_demo3.id]] IN (SELECT [[__mm_demo3_demo4_via_rel_many_cascade_je.value]] FROM json_each(CASE WHEN json_valid([[__mm_demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) THEN [[__mm_demo3_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[__mm_demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{__mm_demo3_demo4_via_rel_many_cascade_je}}) WHERE `__mm_demo3`.`id` = `demo3`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = 1)))))",
		},
		{
			"back relations via unique multiple relation field (should be the same as multi-match)",
			"demo3",
			"demo4_via_rel_many_unique.id = true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_many_unique` ON [[demo3.id]] IN (SELECT [[demo3_demo4_via_rel_many_unique_je.value]] FROM json_each(CASE WHEN json_valid([[demo3_demo4_via_rel_many_unique.rel_many_unique]]) THEN [[demo3_demo4_via_rel_many_unique.rel_many_unique]] ELSE json_array([[demo3_demo4_via_rel_many_unique.rel_many_unique]]) END) {{demo3_demo4_via_rel_many_unique_je}}) WHERE ((([[demo3_demo4_via_rel_many_unique.id]] = 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo3_demo4_via_rel_many_unique.id]] as [[multiMatchValue]] FROM `demo3` `__mm_demo3` LEFT JOIN `demo4` `__mm_demo3_demo4_via_rel_many_unique` ON [[__mm_demo3.id]] IN (SELECT [[__mm_demo3_demo4_via_rel_many_unique_je.value]] FROM json_each(CASE WHEN json_valid([[__mm_demo3_demo4_via_rel_many_unique.rel_many_unique]]) THEN [[__mm_demo3_demo4_via_rel_many_unique.rel_many_unique]] ELSE json_array([[__mm_demo3_demo4_via_rel_many_unique.rel_many_unique]]) END) {{__mm_demo3_demo4_via_rel_many_unique_je}}) WHERE `__mm_demo3`.`id` = `demo3`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = 1)))))",
		},
		{
			"recursive back relations",
			"demo3",
			"demo4_via_rel_many_cascade.rel_one_cascade.demo4_via_rel_many_cascade.id ?= true",
			false,
			"SELECT DISTINCT `demo3`.* FROM `demo3` LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade` ON [[demo3.id]] IN (SELECT [[demo3_demo4_via_rel_many_cascade_je.value]] FROM json_each(CASE WHEN json_valid([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) THEN [[demo3_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{demo3_demo4_via_rel_many_cascade_je}}) LEFT JOIN `demo3` `demo3_demo4_via_rel_many_cascade_rel_one_cascade` ON [[demo3_demo4_via_rel_many_cascade_rel_one_cascade.id]] = [[demo3_demo4_via_rel_many_cascade.rel_one_cascade]] LEFT JOIN `demo4` `demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade` ON [[demo3_demo4_via_rel_many_cascade_rel_one_cascade.id]] IN (SELECT [[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade_je.value]] FROM json_each(CASE WHEN json_valid([[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]]) THEN [[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]] ELSE json_array([[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.rel_many_cascade]]) END) {{demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade_je}}) WHERE [[demo3_demo4_via_rel_many_cascade_rel_one_cascade_demo4_via_rel_many_cascade.id]] = 1",
		},
		{
			"@collection join (opt/any operators)",
			"demo4",
			"@collection.demo1.text ?> true || @collection.demo2.active ?> true || @collection.demo1:demo1_alias.file_one ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo1` `__collection_demo1` LEFT JOIN `demo2` `__collection_demo2` LEFT JOIN `demo1` `__collection_alias_demo1_alias` WHERE ([[__collection_demo1.text]] > 1 OR [[__collection_demo2.active]] > 1 OR [[__collection_alias_demo1_alias.file_one]] > 1)",
		},
		{
			"@collection join (multi-match operators)",
			"demo4",
			"@collection.demo1.text > true || @collection.demo2.active > true || @collection.demo1.file_one > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo1` `__collection_demo1` LEFT JOIN `demo2` `__collection_demo2` WHERE ((([[__collection_demo1.text]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__collection_demo1.text]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo1` `__mm__collection_demo1` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{TEST}} WHERE ((NOT ([[TEST.multiMatchValue]] > 1)) OR ([[TEST.multiMatchValue]] IS NULL))))) OR (([[__collection_demo2.active]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__collection_demo2.active]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo2` `__mm__collection_demo2` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{TEST}} WHERE ((NOT ([[TEST.multiMatchValue]] > 1)) OR ([[TEST.multiMatchValue]] IS NULL))))) OR (([[__collection_demo1.file_one]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__collection_demo1.file_one]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo1` `__mm__collection_demo1` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{TEST}} WHERE ((NOT ([[TEST.multiMatchValue]] > 1)) OR ([[TEST.multiMatchValue]] IS NULL))))))",
		},
		{
			"@request.auth fields",
			"demo4",
			"@request.auth.id > true || @request.auth.username > true || @request.auth.rel.title > true || @request.data.demo < true || @request.auth.missingA.missingB > false",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `users` `__auth_users` ON `__auth_users`.`id`={:p0} LEFT JOIN `demo2` `__auth_users_rel` ON [[__auth_users_rel.id]] = [[__auth_users.rel]] WHERE ({:TEST} > 1 OR {:TEST} > 1 OR [[__auth_users_rel.title]] > 1 OR NULL < 1 OR NULL > 0)",
		},
		{
			"@request.* static fields",
			"demo4",
			"@request.context = true || @request.query.a = true || @request.query.b = true || @request.query.missing = true || @request.headers.a = true || @request.headers.missing = true",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE ({:TEST} = 1 OR '' = 1 OR {:TEST} = 1 OR '' = 1 OR {:TEST} = 1 OR '' = 1)",
		},
		{
			"hidden field with system filters (multi-match and ignore emailVisibility)",
			"demo4",
			"@collection.users.email > true || @request.auth.email > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `users` `__collection_users` WHERE ((([[__collection_users.email]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__collection_users.email]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `users` `__mm__collection_users` WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))) OR {:TEST} > 1)",
		},
		{
			"hidden field (add emailVisibility)",
			"users",
			"id > true || email > true",
			false,
			"SELECT `users`.* FROM `users` WHERE ([[users.id]] > 1 OR (([[users.email]] > 1) AND ([[users.emailVisibility]] = TRUE)))",
		},
		{
			"hidden field (force ignore emailVisibility)",
			"users",
			"email > true",
			true,
			"SELECT `users`.* FROM `users` WHERE [[users.email]] > 1",
		},
		{
			"isset key",
			"demo1",
			"@request.data.a:isset > true ||" +
				"@request.data.b:isset > true ||" +
				"@request.data.c:isset > true ||" +
				"@request.query.a:isset > true ||" +
				"@request.query.b:isset > true ||" +
				"@request.query.c:isset > true",
			false,
			"SELECT `demo1`.* FROM `demo1` WHERE (TRUE > 1 OR TRUE > 1 OR FALSE > 1 OR TRUE > 1 OR TRUE > 1 OR FALSE > 1)",
		},
		{
			"@request.data.rel.* fields",
			"demo4",
			"@request.data.rel_one_cascade.title > true &&" +
				// reference the same as rel_one_cascade collection but should use a different join alias
				"@request.data.rel_one_no_cascade.title < true &&" +
				// different collection
				"@request.data.self_rel_many.title = true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo3` `__data_demo3_rel_one_cascade` ON [[__data_demo3_rel_one_cascade.id]]={:p0} LEFT JOIN `demo3` `__data_demo3_rel_one_no_cascade` ON [[__data_demo3_rel_one_no_cascade.id]]={:p1} LEFT JOIN `demo4` `__data_demo4_self_rel_many` ON [[__data_demo4_self_rel_many.id]]={:p2} WHERE ([[__data_demo3_rel_one_cascade.title]] > 1 AND [[__data_demo3_rel_one_no_cascade.title]] < 1 AND (([[__data_demo4_self_rel_many.title]] = 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__data_mm_demo4_self_rel_many.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo4` `__data_mm_demo4_self_rel_many` ON [[__data_mm_demo4_self_rel_many.id]]={:p3} WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = 1)))))",
		},
		{
			"@request.data.arrayble:each fields",
			"demo1",
			"@request.data.select_one:each > true &&" +
				"@request.data.select_one:each ?< true &&" +
				"@request.data.select_many:each > true &&" +
				"@request.data.select_many:each ?< true &&" +
				"@request.data.file_one:each > true &&" +
				"@request.data.file_one:each ?< true &&" +
				"@request.data.file_many:each > true &&" +
				"@request.data.file_many:each ?< true &&" +
				"@request.data.rel_one:each > true &&" +
				"@request.data.rel_one:each ?< true &&" +
				"@request.data.rel_many:each > true &&" +
				"@request.data.rel_many:each ?< true",
			false,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_select_one_je` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_select_many_je` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_file_one_je` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_file_many_je` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_rel_one_je` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_rel_many_je` WHERE ([[__dataEach_select_one_je.value]] > 1 AND [[__dataEach_select_one_je.value]] < 1 AND (([[__dataEach_select_many_je.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__dataEach_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:mmdataEachTEST}) `__mm__dataEach_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND [[__dataEach_select_many_je.value]] < 1 AND [[__dataEach_file_one_je.value]] > 1 AND [[__dataEach_file_one_je.value]] < 1 AND (([[__dataEach_file_many_je.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__dataEach_file_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:mmdataEachTEST}) `__mm__dataEach_file_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND [[__dataEach_file_many_je.value]] < 1 AND [[__dataEach_rel_one_je.value]] > 1 AND [[__dataEach_rel_one_je.value]] < 1 AND (([[__dataEach_rel_many_je.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__dataEach_rel_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:mmdataEachTEST}) `__mm__dataEach_rel_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND [[__dataEach_rel_many_je.value]] < 1)",
		},
		{
			"regular arrayble:each fields",
			"demo1",
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
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.select_one]]) THEN [[demo1.select_one]] ELSE json_array([[demo1.select_one]]) END) `demo1_select_one_je` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.select_many]]) THEN [[demo1.select_many]] ELSE json_array([[demo1.select_many]]) END) `demo1_select_many_je` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.file_one]]) THEN [[demo1.file_one]] ELSE json_array([[demo1.file_one]]) END) `demo1_file_one_je` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.file_many]]) THEN [[demo1.file_many]] ELSE json_array([[demo1.file_many]]) END) `demo1_file_many_je` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.rel_one]]) THEN [[demo1.rel_one]] ELSE json_array([[demo1.rel_one]]) END) `demo1_rel_one_je` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.rel_many]]) THEN [[demo1.rel_many]] ELSE json_array([[demo1.rel_many]]) END) `demo1_rel_many_je` WHERE ([[demo1_select_one_je.value]] > 1 AND [[demo1_select_one_je.value]] < 1 AND (([[demo1_select_many_je.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.select_many]]) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__mm_demo1_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND [[demo1_select_many_je.value]] < 1 AND [[demo1_file_one_je.value]] > 1 AND [[demo1_file_one_je.value]] < 1 AND (([[demo1_file_many_je.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_file_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.file_many]]) THEN [[__mm_demo1.file_many]] ELSE json_array([[__mm_demo1.file_many]]) END) `__mm_demo1_file_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND [[demo1_file_many_je.value]] < 1 AND [[demo1_rel_one_je.value]] > 1 AND [[demo1_rel_one_je.value]] < 1 AND (([[demo1_rel_many_je.value]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND [[demo1_rel_many_je.value]] < 1)",
		},
		{
			"arrayble:each vs arrayble:each",
			"demo1",
			"select_one:each != select_many:each &&" +
				"select_many:each > select_one:each &&" +
				"select_many:each ?< select_one:each &&" +
				"select_many:each = @request.data.select_many:each",
			false,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.select_one]]) THEN [[demo1.select_one]] ELSE json_array([[demo1.select_one]]) END) `demo1_select_one_je` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.select_many]]) THEN [[demo1.select_many]] ELSE json_array([[demo1.select_many]]) END) `demo1_select_many_je` LEFT JOIN json_each({:dataEachTEST}) `__dataEach_select_many_je` WHERE (((COALESCE([[demo1_select_one_je.value]], '') IS NOT COALESCE([[demo1_select_many_je.value]], '')) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.select_many]]) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__mm_demo1_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT (COALESCE([[demo1_select_one_je.value]], '') IS NOT COALESCE([[__smTEST.multiMatchValue]], ''))) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND (([[demo1_select_many_je.value]] > [[demo1_select_one_je.value]]) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.select_many]]) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__mm_demo1_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > [[demo1_select_one_je.value]])) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND [[demo1_select_many_je.value]] < [[demo1_select_one_je.value]] AND (([[demo1_select_many_je.value]] = [[__dataEach_select_many_je.value]]) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.select_many]]) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__mm_demo1_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm__dataEach_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:mmdataEachTEST}) `__mm__dataEach_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE NOT (COALESCE([[__mlTEST.multiMatchValue]], '') = COALESCE([[__mrTEST.multiMatchValue]], ''))))))",
		},
		{
			"mixed multi-match vs multi-match",
			"demo1",
			"rel_many.rel.active != rel_many.name &&" +
				"rel_many.rel.active ?= rel_many.name &&" +
				"rel_many.rel.title ~ rel_one.email &&" +
				"@collection.demo2.active = rel_many.rel.active &&" +
				"@collection.demo2.active ?= rel_many.rel.active &&" +
				"rel_many.email > @request.data.rel_many.email",
			false,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.rel_many]]) THEN [[demo1.rel_many]] ELSE json_array([[demo1.rel_many]]) END) `demo1_rel_many_je` LEFT JOIN `users` `demo1_rel_many` ON [[demo1_rel_many.id]] = [[demo1_rel_many_je.value]] LEFT JOIN `demo2` `demo1_rel_many_rel` ON [[demo1_rel_many_rel.id]] = [[demo1_rel_many.rel]] LEFT JOIN `demo1` `demo1_rel_one` ON [[demo1_rel_one.id]] = [[demo1.rel_one]] LEFT JOIN `demo2` `__collection_demo2` LEFT JOIN `users` `__data_users_rel_many` ON [[__data_users_rel_many.id]] IN ({:p0}, {:p1}) WHERE (((COALESCE([[demo1_rel_many_rel.active]], '') IS NOT COALESCE([[demo1_rel_many.name]], '')) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many_rel.active]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] LEFT JOIN `demo2` `__mm_demo1_rel_many_rel` ON [[__mm_demo1_rel_many_rel.id]] = [[__mm_demo1_rel_many.rel]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm_demo1_rel_many.name]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE ((NOT (COALESCE([[__mlTEST.multiMatchValue]], '') IS NOT COALESCE([[__mrTEST.multiMatchValue]], ''))) OR ([[__mlTEST.multiMatchValue]] IS NULL) OR ([[__mrTEST.multiMatchValue]] IS NULL))))) AND COALESCE([[demo1_rel_many_rel.active]], '') = COALESCE([[demo1_rel_many.name]], '') AND (([[demo1_rel_many_rel.title]] LIKE ('%' || [[demo1_rel_one.email]] || '%') ESCAPE '\\') AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many_rel.title]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] LEFT JOIN `demo2` `__mm_demo1_rel_many_rel` ON [[__mm_demo1_rel_many_rel.id]] = [[__mm_demo1_rel_many.rel]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] LIKE ('%' || [[demo1_rel_one.email]] || '%') ESCAPE '\\')) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND ((COALESCE([[__collection_demo2.active]], '') = COALESCE([[demo1_rel_many_rel.active]], '')) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__collection_demo2.active]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN `demo2` `__mm__collection_demo2` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm_demo1_rel_many_rel.active]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] LEFT JOIN `demo2` `__mm_demo1_rel_many_rel` ON [[__mm_demo1_rel_many_rel.id]] = [[__mm_demo1_rel_many.rel]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE NOT (COALESCE([[__mlTEST.multiMatchValue]], '') = COALESCE([[__mrTEST.multiMatchValue]], ''))))) AND COALESCE([[__collection_demo2.active]], '') = COALESCE([[demo1_rel_many_rel.active]], '') AND (((([[demo1_rel_many.email]] > [[__data_users_rel_many.email]]) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many.email]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__data_mm_users_rel_many.email]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN `users` `__data_mm_users_rel_many` ON [[__data_mm_users_rel_many.id]] IN ({:p2}, {:p3}) WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE ((NOT ([[__mlTEST.multiMatchValue]] > [[__mrTEST.multiMatchValue]])) OR ([[__mlTEST.multiMatchValue]] IS NULL) OR ([[__mrTEST.multiMatchValue]] IS NULL)))))) AND ([[demo1_rel_many.emailVisibility]] = TRUE)))",
		},
		{
			"@request.data.arrayable:length fields",
			"demo1",
			"@request.data.select_one:length > 1 &&" +
				"@request.data.select_one:length ?> 2 &&" +
				"@request.data.select_many:length < 3 &&" +
				"@request.data.select_many:length ?> 4 &&" +
				"@request.data.rel_one:length = 5 &&" +
				"@request.data.rel_one:length ?= 6 &&" +
				"@request.data.rel_many:length != 7 &&" +
				"@request.data.rel_many:length ?!= 8 &&" +
				"@request.data.file_one:length = 9 &&" +
				"@request.data.file_one:length ?= 0 &&" +
				"@request.data.file_many:length != 1 &&" +
				"@request.data.file_many:length ?!= 2",
			false,
			"SELECT `demo1`.* FROM `demo1` WHERE (0 > {:TEST} AND 0 > {:TEST} AND 2 < {:TEST} AND 2 > {:TEST} AND 1 = {:TEST} AND 1 = {:TEST} AND 2 IS NOT {:TEST} AND 2 IS NOT {:TEST} AND 1 = {:TEST} AND 1 = {:TEST} AND 3 IS NOT {:TEST} AND 3 IS NOT {:TEST})",
		},
		{
			"regular arrayable:length fields",
			"demo4",
			"@request.data.self_rel_one.self_rel_many:length > 1 &&" +
				"@request.data.self_rel_one.self_rel_many:length ?> 2 &&" +
				"@request.data.rel_many_cascade.files:length ?< 3 &&" +
				"@request.data.rel_many_cascade.files:length < 4 &&" +
				"@request.data.rel_one_cascade.files:length < 4.1 &&" + // to ensure that the join to the same as above table will be aliased
				"self_rel_one.self_rel_many:length = 5 &&" +
				"self_rel_one.self_rel_many:length ?= 6 &&" +
				"self_rel_one.rel_many_cascade.files:length != 7 &&" +
				"self_rel_one.rel_many_cascade.files:length ?!= 8",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `__data_demo4_self_rel_one` ON [[__data_demo4_self_rel_one.id]]={:p0} LEFT JOIN `demo3` `__data_demo3_rel_many_cascade` ON [[__data_demo3_rel_many_cascade.id]] IN ({:p1}, {:p2}) LEFT JOIN `demo3` `__data_demo3_rel_one_cascade` ON [[__data_demo3_rel_one_cascade.id]]={:p3} LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_one.rel_many_cascade]]) THEN [[demo4_self_rel_one.rel_many_cascade]] ELSE json_array([[demo4_self_rel_one.rel_many_cascade]]) END) `demo4_self_rel_one_rel_many_cascade_je` LEFT JOIN `demo3` `demo4_self_rel_one_rel_many_cascade` ON [[demo4_self_rel_one_rel_many_cascade.id]] = [[demo4_self_rel_one_rel_many_cascade_je.value]] WHERE (json_array_length(CASE WHEN json_valid([[__data_demo4_self_rel_one.self_rel_many]]) THEN [[__data_demo4_self_rel_one.self_rel_many]] ELSE (CASE WHEN [[__data_demo4_self_rel_one.self_rel_many]] = '' OR [[__data_demo4_self_rel_one.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[__data_demo4_self_rel_one.self_rel_many]]) END) END) > {:TEST} AND json_array_length(CASE WHEN json_valid([[__data_demo4_self_rel_one.self_rel_many]]) THEN [[__data_demo4_self_rel_one.self_rel_many]] ELSE (CASE WHEN [[__data_demo4_self_rel_one.self_rel_many]] = '' OR [[__data_demo4_self_rel_one.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[__data_demo4_self_rel_one.self_rel_many]]) END) END) > {:TEST} AND json_array_length(CASE WHEN json_valid([[__data_demo3_rel_many_cascade.files]]) THEN [[__data_demo3_rel_many_cascade.files]] ELSE (CASE WHEN [[__data_demo3_rel_many_cascade.files]] = '' OR [[__data_demo3_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__data_demo3_rel_many_cascade.files]]) END) END) < {:TEST} AND ((json_array_length(CASE WHEN json_valid([[__data_demo3_rel_many_cascade.files]]) THEN [[__data_demo3_rel_many_cascade.files]] ELSE (CASE WHEN [[__data_demo3_rel_many_cascade.files]] = '' OR [[__data_demo3_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__data_demo3_rel_many_cascade.files]]) END) END) < {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT json_array_length(CASE WHEN json_valid([[__data_mm_demo3_rel_many_cascade.files]]) THEN [[__data_mm_demo3_rel_many_cascade.files]] ELSE (CASE WHEN [[__data_mm_demo3_rel_many_cascade.files]] = '' OR [[__data_mm_demo3_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__data_mm_demo3_rel_many_cascade.files]]) END) END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo3` `__data_mm_demo3_rel_many_cascade` ON [[__data_mm_demo3_rel_many_cascade.id]] IN ({:p8}, {:p9}) WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] < {:TEST})) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND json_array_length(CASE WHEN json_valid([[__data_demo3_rel_one_cascade.files]]) THEN [[__data_demo3_rel_one_cascade.files]] ELSE (CASE WHEN [[__data_demo3_rel_one_cascade.files]] = '' OR [[__data_demo3_rel_one_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__data_demo3_rel_one_cascade.files]]) END) END) < {:TEST} AND json_array_length(CASE WHEN json_valid([[demo4_self_rel_one.self_rel_many]]) THEN [[demo4_self_rel_one.self_rel_many]] ELSE (CASE WHEN [[demo4_self_rel_one.self_rel_many]] = '' OR [[demo4_self_rel_one.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[demo4_self_rel_one.self_rel_many]]) END) END) = {:TEST} AND json_array_length(CASE WHEN json_valid([[demo4_self_rel_one.self_rel_many]]) THEN [[demo4_self_rel_one.self_rel_many]] ELSE (CASE WHEN [[demo4_self_rel_one.self_rel_many]] = '' OR [[demo4_self_rel_one.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[demo4_self_rel_one.self_rel_many]]) END) END) = {:TEST} AND ((json_array_length(CASE WHEN json_valid([[demo4_self_rel_one_rel_many_cascade.files]]) THEN [[demo4_self_rel_one_rel_many_cascade.files]] ELSE (CASE WHEN [[demo4_self_rel_one_rel_many_cascade.files]] = '' OR [[demo4_self_rel_one_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[demo4_self_rel_one_rel_many_cascade.files]]) END) END) IS NOT {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT json_array_length(CASE WHEN json_valid([[__mm_demo4_self_rel_one_rel_many_cascade.files]]) THEN [[__mm_demo4_self_rel_one_rel_many_cascade.files]] ELSE (CASE WHEN [[__mm_demo4_self_rel_one_rel_many_cascade.files]] = '' OR [[__mm_demo4_self_rel_one_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[__mm_demo4_self_rel_one_rel_many_cascade.files]]) END) END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo4` `__mm_demo4_self_rel_one` ON [[__mm_demo4_self_rel_one.id]] = [[__mm_demo4.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4_self_rel_one.rel_many_cascade]]) THEN [[__mm_demo4_self_rel_one.rel_many_cascade]] ELSE json_array([[__mm_demo4_self_rel_one.rel_many_cascade]]) END) `__mm_demo4_self_rel_one_rel_many_cascade_je` LEFT JOIN `demo3` `__mm_demo4_self_rel_one_rel_many_cascade` ON [[__mm_demo4_self_rel_one_rel_many_cascade.id]] = [[__mm_demo4_self_rel_one_rel_many_cascade_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] IS NOT {:TEST})) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND json_array_length(CASE WHEN json_valid([[demo4_self_rel_one_rel_many_cascade.files]]) THEN [[demo4_self_rel_one_rel_many_cascade.files]] ELSE (CASE WHEN [[demo4_self_rel_one_rel_many_cascade.files]] = '' OR [[demo4_self_rel_one_rel_many_cascade.files]] IS NULL THEN json_array() ELSE json_array([[demo4_self_rel_one_rel_many_cascade.files]]) END) END) IS NOT {:TEST})",
		},
		{
			"json_extract and json_array_length COALESCE equal normalizations",
			"demo4",
			"json_object.a.b = '' && self_rel_many:length != 2 && json_object.a.b > 3 && self_rel_many:length <= 4",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE ((CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$.a.b') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb.a.b') END) IS {:TEST} AND json_array_length(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE (CASE WHEN [[demo4.self_rel_many]] = '' OR [[demo4.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[demo4.self_rel_many]]) END) END) IS NOT {:TEST} AND (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$.a.b') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb.a.b') END) > {:TEST} AND json_array_length(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE (CASE WHEN [[demo4.self_rel_many]] = '' OR [[demo4.self_rel_many]] IS NULL THEN json_array() ELSE json_array([[demo4.self_rel_many]]) END) END) <= {:TEST})",
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
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] WHERE ((CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS {:TEST} OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NOT {:TEST} OR {:TEST} IS (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR {:TEST} IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NULL OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NOT NULL OR NULL IS (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR NULL IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS 1 OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NOT 1 OR 1 IS (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR 1 IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) IS [[demo4.title]] OR [[demo4.title]] IS NOT (CASE WHEN json_valid([[demo4.json_object]]) THEN JSON_EXTRACT([[demo4.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4.json_object]]), '$.pb') END) OR (((CASE WHEN json_valid([[demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_many.json_object]]), '$.pb') END) IS {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT (CASE WHEN json_valid([[__mm_demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[__mm_demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[__mm_demo4_self_rel_many.json_object]]), '$.pb') END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] IS {:TEST})))) OR ((NULL IS (CASE WHEN json_valid([[demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_many.json_object]]), '$.pb') END)) AND (NOT EXISTS (SELECT 1 FROM (SELECT (CASE WHEN json_valid([[__mm_demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[__mm_demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[__mm_demo4_self_rel_many.json_object]]), '$.pb') END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT (NULL IS [[__smTEST.multiMatchValue]])))) OR (((CASE WHEN json_valid([[demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_many.json_object]]), '$.pb') END) IS (CASE WHEN json_valid([[demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[demo4_self_rel_many.json_object]]), '$.pb') END)) AND (NOT EXISTS (SELECT 1 FROM (SELECT (CASE WHEN json_valid([[__mm_demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[__mm_demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[__mm_demo4_self_rel_many.json_object]]), '$.pb') END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__mlTEST}} LEFT JOIN (SELECT (CASE WHEN json_valid([[__mm_demo4_self_rel_many.json_object]]) THEN JSON_EXTRACT([[__mm_demo4_self_rel_many.json_object]], '$') ELSE JSON_EXTRACT(json_object('pb', [[__mm_demo4_self_rel_many.json_object]]), '$.pb') END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__mrTEST}} WHERE NOT ([[__mlTEST.multiMatchValue]] IS [[__mrTEST.multiMatchValue]])))))",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			collection, err := app.Dao().FindCollectionByNameOrId(s.collectionIdOrName)
			if err != nil {
				t.Fatalf("[%s] Failed to load collection %s: %v", s.name, s.collectionIdOrName, err)
			}

			query := app.Dao().RecordQuery(collection)

			r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestInfo, s.allowHiddenFields)

			expr, err := search.FilterData(s.rule).BuildExpr(r)
			if err != nil {
				t.Fatalf("[%s] BuildExpr failed with error %v", s.name, err)
			}

			if err := r.UpdateQuery(query); err != nil {
				t.Fatalf("[%s] UpdateQuery failed with error %v", s.name, err)
			}

			rawQuery := query.AndWhere(expr).Build().SQL()

			// replace TEST placeholder with .+ regex pattern
			expectQuery := strings.ReplaceAll(
				"^"+regexp.QuoteMeta(s.expectQuery)+"$",
				"TEST",
				`\w+`,
			)

			if !list.ExistInSliceWithRegex(rawQuery, []string{expectQuery}) {
				t.Fatalf("[%s] Expected query\n %v \ngot:\n %v", s.name, expectQuery, rawQuery)
			}
		})
	}
}

func TestRecordFieldResolverResolveSchemaFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo4")
	if err != nil {
		t.Fatal(err)
	}

	authRecord, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestInfo := &models.RequestInfo{
		AuthRecord: authRecord,
	}

	r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestInfo, true)

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
		{"rel_one_cascade.demo4_via_title.id", true, ""}, // non-relation via field
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
		{"@request.auth.rel.missing", false, "NULL"},
		{"@request.auth.missing_via_rel", false, "NULL"},

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

	collection, err := app.Dao().FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	authRecord, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestInfo := &models.RequestInfo{
		Context: "ctx",
		Method:  "get",
		Query: map[string]any{
			"a": 123,
		},
		Data: map[string]any{
			"number":         "10",
			"number_unknown": "20",
			"b":              456,
			"c":              map[string]int{"sub": 1},
		},
		Headers: map[string]any{
			"d": "789",
		},
		AuthRecord: authRecord,
	}

	r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestInfo, true)

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
		{"@request.query.a", false, `123`},
		{"@request.query.a.missing", false, ``},
		{"@request.headers", true, ``},
		{"@request.headers.missing", false, ``},
		{"@request.headers.d", false, `"789"`},
		{"@request.headers.d.sub", true, ``},
		{"@request.data", true, ``},
		{"@request.data.b", false, `456`},
		{"@request.data.number", false, `10`},           // number field normalization
		{"@request.data.number_unknown", false, `"20"`}, // no numeric normalizations for unknown fields
		{"@request.data.b.missing", false, ``},
		{"@request.data.c", false, `"{\"sub\":1}"`},
		{"@request.auth", true, ""},
		{"@request.auth.id", false, `"4q1xlclmfloku33"`},
		{"@request.auth.username", false, `"users75657"`},
		{"@request.auth.verified", false, `false`},
		{"@request.auth.emailVisibility", false, `false`},
		{"@request.auth.email", false, `"test@example.com"`}, // should always be returned no matter of the emailVisibility state
		{"@request.auth.missing", false, `NULL`},
	}

	for i, s := range scenarios {
		t.Run(s.fieldName, func(t *testing.T) {
			r, err := r.Resolve(s.fieldName)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			// missing key
			// ---
			if len(r.Params) == 0 {
				if r.Identifier != "NULL" {
					t.Fatalf("(%d) Expected 0 placeholder parameters for %v, got %v", i, r.Identifier, r.Params)
				}
				return
			}

			// existing key
			// ---
			if len(r.Params) != 1 {
				t.Fatalf("(%d) Expected 1 placeholder parameter for %v, got %v", i, r.Identifier, r.Params)
			}

			var paramName string
			var paramValue any
			for k, v := range r.Params {
				paramName = k
				paramValue = v
			}

			if r.Identifier != ("{:" + paramName + "}") {
				t.Fatalf("(%d) Expected parameter r.Identifier %q, got %q", i, paramName, r.Identifier)
			}

			encodedParamValue, _ := json.Marshal(paramValue)
			if string(encodedParamValue) != s.expectParamValue {
				t.Fatalf("(%d) Expected r.Params %v for %v, got %v", i, s.expectParamValue, r.Identifier, string(encodedParamValue))
			}
		})
	}

	// ensure that the original email visibility was restored
	if authRecord.EmailVisibility() {
		t.Fatal("Expected the original authRecord emailVisibility to remain unchanged")
	}
	if v, ok := authRecord.PublicExport()[schema.FieldNameEmail]; ok {
		t.Fatalf("Expected the original authRecord email to not be exported, got %q", v)
	}
}
