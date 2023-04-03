package resolvers_test

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/search"
)

func TestRecordFieldResolverUpdateQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	authRecord, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	requestData := &models.RequestData{
		Headers: map[string]any{
			"a": "123",
			"b": "456",
		},
		Query: map[string]any{
			"a": nil,
			"b": 123,
		},
		Data: map[string]any{
			"a":                nil,
			"b":                123,
			"number":           10,
			"select_many":      []string{"optionA", "optionC"},
			"rel_one":          "test",
			"rel_many":         []string{"test1", "test2"},
			"file_one":         "test",
			"file_many":        []string{"test1", "test2", "test3"},
			"self_rel_one":     "test",
			"self_rel_many":    []string{"test1"},
			"rel_many_cascade": []string{"test1", "test2"},
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
			"SELECT `demo4`.* FROM `demo4` WHERE ([[demo4.title]] = 1 OR [[demo4.title]] != {:TEST} OR [[demo4.title]] LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] NOT LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] > {:TEST} OR [[demo4.title]] >= {:TEST} OR [[demo4.title]] < {:TEST} OR [[demo4.title]] <= {:TEST})",
		},
		{
			"non relation field (with all opt/any operators)",
			"demo4",
			"title ?= true || title ?!= 'test' || title ?~ 'test1' || title ?!~ '%test2' || title ?> 1 || title ?>= 2 || title ?< 3 || title ?<= 4",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE ([[demo4.title]] = 1 OR [[demo4.title]] != {:TEST} OR [[demo4.title]] LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] NOT LIKE {:TEST} ESCAPE '\\' OR [[demo4.title]] > {:TEST} OR [[demo4.title]] >= {:TEST} OR [[demo4.title]] < {:TEST} OR [[demo4.title]] <= {:TEST})",
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
			"nested incomplete rels (opt/any operator)",
			"demo4",
			"self_rel_many.self_rel_one ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] WHERE [[demo4_self_rel_many.self_rel_one]] > 1",
		},
		{
			"nested incomplete rels (multi-match operator)",
			"demo4",
			"self_rel_many.self_rel_one > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] WHERE ((([[demo4_self_rel_many.self_rel_one]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many.self_rel_one]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{TEST}} WHERE ((NOT ([[TEST.multiMatchValue]] > 1)) OR ([[TEST.multiMatchValue]] IS NULL))))))",
		},
		{
			"nested complete rels (opt/any operator)",
			"demo4",
			"self_rel_many.self_rel_one.title ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] WHERE [[demo4_self_rel_many_self_rel_one.title]] > 1",
		},
		{
			"nested complete rels (multi-match operator)",
			"demo4",
			"self_rel_many.self_rel_one.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] WHERE ((([[demo4_self_rel_many_self_rel_one.title]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many_self_rel_one.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one` ON [[__mm_demo4_self_rel_many_self_rel_one.id]] = [[__mm_demo4_self_rel_many.self_rel_one]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))))",
		},
		{
			"repeated nested rels (opt/any operator)",
			"demo4",
			"self_rel_many.self_rel_one.self_rel_many.self_rel_one.title ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_many_self_rel_one.self_rel_many]]) THEN [[demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `demo4_self_rel_many_self_rel_one_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many` ON [[demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] WHERE [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]] > 1",
		},
		{
			"repeated nested rels (multi-match operator)",
			"demo4",
			"self_rel_many.self_rel_one.self_rel_many.self_rel_one.title > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_many_self_rel_one.self_rel_many]]) THEN [[demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `demo4_self_rel_many_self_rel_one_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many` ON [[demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] WHERE ((([[demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]] > 1) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one` ON [[__mm_demo4_self_rel_many_self_rel_one.id]] = [[__mm_demo4_self_rel_many.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]]) THEN [[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]] ELSE json_array([[__mm_demo4_self_rel_many_self_rel_one.self_rel_many]]) END) `__mm_demo4_self_rel_many_self_rel_one_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one_self_rel_many` ON [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many.id]] = [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many_je.value]] LEFT JOIN `demo4` `__mm_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one` ON [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.id]] = [[__mm_demo4_self_rel_many_self_rel_one_self_rel_many.self_rel_one]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > 1)) OR ([[__smTEST.multiMatchValue]] IS NULL))))))",
		},
		{
			"multiple rels (opt/any operators)",
			"demo4",
			"self_rel_many.title ?= 'test' || self_rel_one.json_object.a ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] WHERE ([[demo4_self_rel_many.title]] = {:TEST} OR JSON_EXTRACT([[demo4_self_rel_one.json_object]], '$.a') > 1)",
		},
		{
			"multiple rels (multi-match operators)",
			"demo4",
			"self_rel_many.title = 'test' || self_rel_one.json_object.a > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN json_each(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) `demo4_self_rel_many_je` LEFT JOIN `demo4` `demo4_self_rel_many` ON [[demo4_self_rel_many.id]] = [[demo4_self_rel_many_je.value]] LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] WHERE ((([[demo4_self_rel_many.title]] = {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo4_self_rel_many.title]] as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4.self_rel_many]]) THEN [[__mm_demo4.self_rel_many]] ELSE json_array([[__mm_demo4.self_rel_many]]) END) `__mm_demo4_self_rel_many_je` LEFT JOIN `demo4` `__mm_demo4_self_rel_many` ON [[__mm_demo4_self_rel_many.id]] = [[__mm_demo4_self_rel_many_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = {:TEST})))) OR JSON_EXTRACT([[demo4_self_rel_one.json_object]], '$.a') > 1)",
		},
		{
			"@collection join (opt/any operators)",
			"demo4",
			"@collection.demo1.text ?> true || @collection.demo2.active ?> true || @collection.demo1.file_one ?> true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo1` `__collection_demo1` LEFT JOIN `demo2` `__collection_demo2` WHERE ([[__collection_demo1.text]] > 1 OR [[__collection_demo2.active]] > 1 OR [[__collection_demo1.file_one]] > 1)",
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
			"@request.auth.id > true || @request.auth.username > true || @request.auth.rel.title > true || @request.data.demo > true",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `users` `__auth_users` ON `__auth_users`.`id`={:TEST} LEFT JOIN `demo2` `__auth_users_rel` ON [[__auth_users_rel.id]] = [[__auth_users.rel]] WHERE ({:TEST} > 1 OR {:TEST} > 1 OR [[__auth_users_rel.title]] > 1 OR NULL > 1)",
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
			"demo1",
			"@request.data.rel_one > true &&" +
				"@request.data.rel_one.text > true &&" +
				"@request.data.rel_many > true &&" +
				"@request.data.rel_many.email != 'test' &&" +
				"@request.data.rel_many.url ?= 'test' &&" +
				"@request.data.rel_many.avatar ~ 'test'",
			false,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN `demo1` `__data_demo1` ON [[__data_demo1.id]]={:TEST} LEFT JOIN `users` `__data_users` ON [[__data_users.id]] IN ({:TEST}, {:TEST}) WHERE ({:TEST} > 1 AND [[__data_demo1.text]] > 1 AND {:TEST} > 1 AND (([[__data_users.email]] != {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__data_mm_users.email]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN `users` `__data_mm_users` ON `__data_mm_users`.`id` IN ({:TEST}, {:TEST}) WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] != {:TEST})) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND '' = {:TEST} AND (([[__data_users.avatar]] LIKE {:TEST} ESCAPE '\\') AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__data_mm_users.avatar]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN `users` `__data_mm_users` ON `__data_mm_users`.`id` IN ({:TEST}, {:TEST}) WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] LIKE {:TEST} ESCAPE '\\')) OR ([[__smTEST.multiMatchValue]] IS NULL))))))",
		},
		{
			"@request.data.select:each fields",
			"demo1",
			"@request.data.select_one = 'test' &&" +
				"@request.data.select_one:each != 'test' &&" +
				"@request.data.select_one:each ?= 'test' &&" +
				"@request.data.select_many ~ 'test' &&" +
				"@request.data.select_many:each = 'test' &&" +
				"@request.data.select_many:each ?< true",
			false,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each({:TEST}) `__dataSelect_select_one_je` LEFT JOIN json_each({:TEST}) `__dataSelect_select_many_je` WHERE ('' = {:TEST} AND [[__dataSelect_select_one_je.value]] != {:TEST} AND [[__dataSelect_select_one_je.value]] = {:TEST} AND {:TEST} LIKE {:TEST} ESCAPE '\\' AND (([[__dataSelect_select_many_je.value]] = {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__dataSelect_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:TEST}) `__mm__dataSelect_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = {:TEST})))) AND [[__dataSelect_select_many_je.value]] < 1)",
		},
		{
			"regular select:each fields",
			"demo1",
			"select_one = 'test' &&" +
				"select_one:each != 'test' &&" +
				"select_one:each ?> true &&" +
				"select_many ~ 'test' &&" +
				"select_many:each = 'test' &&" +
				"select_many:each ?> true",
			false,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.select_one]]) THEN [[demo1.select_one]] ELSE json_array([[demo1.select_one]]) END) `demo1_select_one_je` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.select_many]]) THEN [[demo1.select_many]] ELSE json_array([[demo1.select_many]]) END) `demo1_select_many_je` WHERE ([[demo1.select_one]] = {:TEST} AND [[demo1_select_one_je.value]] != {:TEST} AND [[demo1_select_one_je.value]] > 1 AND [[demo1.select_many]] LIKE {:TEST} ESCAPE '\\' AND (([[demo1_select_many_je.value]] = {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.select_many]]) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__mm_demo1_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE NOT ([[__smTEST.multiMatchValue]] = {:TEST})))) AND [[demo1_select_many_je.value]] > 1)",
		},
		{
			"select:each vs select:each",
			"demo1",
			"select_one:each != select_many:each &&" +
				"select_many:each > select_one:each &&" +
				"select_many:each ?< select_one:each &&" +
				"select_many:each = @request.data.select_many:each",
			false,
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.select_one]]) THEN [[demo1.select_one]] ELSE json_array([[demo1.select_one]]) END) `demo1_select_one_je` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.select_many]]) THEN [[demo1.select_many]] ELSE json_array([[demo1.select_many]]) END) `demo1_select_many_je` LEFT JOIN json_each({:dataSelectTEST}) `__dataSelect_select_many_je` WHERE (((COALESCE([[demo1_select_one_je.value]], '') != COALESCE([[demo1_select_many_je.value]], '')) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.select_many]]) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__mm_demo1_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT (COALESCE([[demo1_select_one_je.value]], '') != COALESCE([[__smTEST.multiMatchValue]], ''))) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND (([[demo1_select_many_je.value]] > [[demo1_select_one_je.value]]) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.select_many]]) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__mm_demo1_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] > [[demo1_select_one_je.value]])) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND [[demo1_select_many_je.value]] < [[demo1_select_one_je.value]] AND (([[demo1_select_many_je.value]] = [[__dataSelect_select_many_je.value]]) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.select_many]]) THEN [[__mm_demo1.select_many]] ELSE json_array([[__mm_demo1.select_many]]) END) `__mm_demo1_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm__dataSelect_select_many_je.value]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each({:mmdataSelectTEST}) `__mm__dataSelect_select_many_je` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE NOT (COALESCE([[__mlTEST.multiMatchValue]], '') = COALESCE([[__mrTEST.multiMatchValue]], ''))))))",
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
			"SELECT DISTINCT `demo1`.* FROM `demo1` LEFT JOIN json_each(CASE WHEN json_valid([[demo1.rel_many]]) THEN [[demo1.rel_many]] ELSE json_array([[demo1.rel_many]]) END) `demo1_rel_many_je` LEFT JOIN `users` `demo1_rel_many` ON [[demo1_rel_many.id]] = [[demo1_rel_many_je.value]] LEFT JOIN `demo2` `demo1_rel_many_rel` ON [[demo1_rel_many_rel.id]] = [[demo1_rel_many.rel]] LEFT JOIN `demo1` `demo1_rel_one` ON [[demo1_rel_one.id]] = [[demo1.rel_one]] LEFT JOIN `demo2` `__collection_demo2` LEFT JOIN `users` `__data_users` ON [[__data_users.id]] IN ({:TEST}, {:TEST}) WHERE (((COALESCE([[demo1_rel_many_rel.active]], '') != COALESCE([[demo1_rel_many.name]], '')) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many_rel.active]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] LEFT JOIN `demo2` `__mm_demo1_rel_many_rel` ON [[__mm_demo1_rel_many_rel.id]] = [[__mm_demo1_rel_many.rel]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm_demo1_rel_many.name]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE ((NOT (COALESCE([[__mlTEST.multiMatchValue]], '') != COALESCE([[__mrTEST.multiMatchValue]], ''))) OR ([[__mlTEST.multiMatchValue]] IS NULL) OR ([[__mrTEST.multiMatchValue]] IS NULL))))) AND COALESCE([[demo1_rel_many_rel.active]], '') = COALESCE([[demo1_rel_many.name]], '') AND (([[demo1_rel_many_rel.title]] LIKE ('%' || [[demo1_rel_one.email]] || '%') ESCAPE '\\') AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many_rel.title]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] LEFT JOIN `demo2` `__mm_demo1_rel_many_rel` ON [[__mm_demo1_rel_many_rel.id]] = [[__mm_demo1_rel_many.rel]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] LIKE ('%' || [[demo1_rel_one.email]] || '%') ESCAPE '\\')) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND ((COALESCE([[__collection_demo2.active]], '') = COALESCE([[demo1_rel_many_rel.active]], '')) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm__collection_demo2.active]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN `demo2` `__mm__collection_demo2` WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__mm_demo1_rel_many_rel.active]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] LEFT JOIN `demo2` `__mm_demo1_rel_many_rel` ON [[__mm_demo1_rel_many_rel.id]] = [[__mm_demo1_rel_many.rel]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE NOT (COALESCE([[__mlTEST.multiMatchValue]], '') = COALESCE([[__mrTEST.multiMatchValue]], ''))))) AND COALESCE([[__collection_demo2.active]], '') = COALESCE([[demo1_rel_many_rel.active]], '') AND (((([[demo1_rel_many.email]] > [[__data_users.email]]) AND (NOT EXISTS (SELECT 1 FROM (SELECT [[__mm_demo1_rel_many.email]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo1.rel_many]]) THEN [[__mm_demo1.rel_many]] ELSE json_array([[__mm_demo1.rel_many]]) END) `__mm_demo1_rel_many_je` LEFT JOIN `users` `__mm_demo1_rel_many` ON [[__mm_demo1_rel_many.id]] = [[__mm_demo1_rel_many_je.value]] WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mlTEST}} LEFT JOIN (SELECT [[__data_mm_users.email]] as [[multiMatchValue]] FROM `demo1` `__mm_demo1` LEFT JOIN `users` `__data_mm_users` ON `__data_mm_users`.`id` IN ({:TEST}, {:TEST}) WHERE `__mm_demo1`.`id` = `demo1`.`id`) {{__mrTEST}} WHERE ((NOT ([[__mlTEST.multiMatchValue]] > [[__mrTEST.multiMatchValue]])) OR ([[__mlTEST.multiMatchValue]] IS NULL) OR ([[__mrTEST.multiMatchValue]] IS NULL)))))) AND ([[demo1_rel_many.emailVisibility]] = TRUE)))",
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
			"SELECT `demo1`.* FROM `demo1` WHERE (0 > {:TEST} AND 0 > {:TEST} AND 2 < {:TEST} AND 2 > {:TEST} AND 1 = {:TEST} AND 1 = {:TEST} AND 2 != {:TEST} AND 2 != {:TEST} AND 1 = {:TEST} AND 1 = {:TEST} AND 3 != {:TEST} AND 3 != {:TEST})",
		},
		{
			"regular arrayable:length fields",
			"demo4",
			"@request.data.self_rel_one.self_rel_many:length > 1 &&" +
				"@request.data.self_rel_one.self_rel_many:length ?> 2 &&" +
				"@request.data.rel_many_cascade.files:length ?< 3 &&" +
				"@request.data.rel_many_cascade.files:length < 4 &&" +
				"self_rel_one.self_rel_many:length = 5 &&" +
				"self_rel_one.self_rel_many:length ?= 6 &&" +
				"self_rel_one.rel_many_cascade.files:length != 7 &&" +
				"self_rel_one.rel_many_cascade.files:length ?!= 8",
			false,
			"SELECT DISTINCT `demo4`.* FROM `demo4` LEFT JOIN `demo4` `__data_demo4` ON [[__data_demo4.id]]={:TEST} LEFT JOIN `demo3` `__data_demo3` ON [[__data_demo3.id]] IN ({:TEST}, {:TEST}) LEFT JOIN `demo4` `demo4_self_rel_one` ON [[demo4_self_rel_one.id]] = [[demo4.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[demo4_self_rel_one.rel_many_cascade]]) THEN [[demo4_self_rel_one.rel_many_cascade]] ELSE json_array([[demo4_self_rel_one.rel_many_cascade]]) END) `demo4_self_rel_one_rel_many_cascade_je` LEFT JOIN `demo3` `demo4_self_rel_one_rel_many_cascade` ON [[demo4_self_rel_one_rel_many_cascade.id]] = [[demo4_self_rel_one_rel_many_cascade_je.value]] WHERE (json_array_length(CASE WHEN json_valid([[__data_demo4.self_rel_many]]) THEN [[__data_demo4.self_rel_many]] ELSE json_array([[__data_demo4.self_rel_many]]) END) > {:TEST} AND json_array_length(CASE WHEN json_valid([[__data_demo4.self_rel_many]]) THEN [[__data_demo4.self_rel_many]] ELSE json_array([[__data_demo4.self_rel_many]]) END) > {:TEST} AND json_array_length(CASE WHEN json_valid([[__data_demo3.files]]) THEN [[__data_demo3.files]] ELSE json_array([[__data_demo3.files]]) END) < {:TEST} AND ((json_array_length(CASE WHEN json_valid([[__data_demo3.files]]) THEN [[__data_demo3.files]] ELSE json_array([[__data_demo3.files]]) END) < {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT json_array_length(CASE WHEN json_valid([[__data_mm_demo3.files]]) THEN [[__data_mm_demo3.files]] ELSE json_array([[__data_mm_demo3.files]]) END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo3` `__data_mm_demo3` ON `__data_mm_demo3`.`id` IN ({:TEST}, {:TEST}) WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] < {:TEST})) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND json_array_length(CASE WHEN json_valid([[demo4_self_rel_one.self_rel_many]]) THEN [[demo4_self_rel_one.self_rel_many]] ELSE json_array([[demo4_self_rel_one.self_rel_many]]) END) = {:TEST} AND json_array_length(CASE WHEN json_valid([[demo4_self_rel_one.self_rel_many]]) THEN [[demo4_self_rel_one.self_rel_many]] ELSE json_array([[demo4_self_rel_one.self_rel_many]]) END) = {:TEST} AND ((json_array_length(CASE WHEN json_valid([[demo4_self_rel_one_rel_many_cascade.files]]) THEN [[demo4_self_rel_one_rel_many_cascade.files]] ELSE json_array([[demo4_self_rel_one_rel_many_cascade.files]]) END) != {:TEST}) AND (NOT EXISTS (SELECT 1 FROM (SELECT json_array_length(CASE WHEN json_valid([[__mm_demo4_self_rel_one_rel_many_cascade.files]]) THEN [[__mm_demo4_self_rel_one_rel_many_cascade.files]] ELSE json_array([[__mm_demo4_self_rel_one_rel_many_cascade.files]]) END) as [[multiMatchValue]] FROM `demo4` `__mm_demo4` LEFT JOIN `demo4` `__mm_demo4_self_rel_one` ON [[__mm_demo4_self_rel_one.id]] = [[__mm_demo4.self_rel_one]] LEFT JOIN json_each(CASE WHEN json_valid([[__mm_demo4_self_rel_one.rel_many_cascade]]) THEN [[__mm_demo4_self_rel_one.rel_many_cascade]] ELSE json_array([[__mm_demo4_self_rel_one.rel_many_cascade]]) END) `__mm_demo4_self_rel_one_rel_many_cascade_je` LEFT JOIN `demo3` `__mm_demo4_self_rel_one_rel_many_cascade` ON [[__mm_demo4_self_rel_one_rel_many_cascade.id]] = [[__mm_demo4_self_rel_one_rel_many_cascade_je.value]] WHERE `__mm_demo4`.`id` = `demo4`.`id`) {{__smTEST}} WHERE ((NOT ([[__smTEST.multiMatchValue]] != {:TEST})) OR ([[__smTEST.multiMatchValue]] IS NULL))))) AND json_array_length(CASE WHEN json_valid([[demo4_self_rel_one_rel_many_cascade.files]]) THEN [[demo4_self_rel_one_rel_many_cascade.files]] ELSE json_array([[demo4_self_rel_one_rel_many_cascade.files]]) END) != {:TEST})",
		},
		{
			"json_extract and json_array_length COALESCE equal normalizations",
			"demo4",
			"json_object.a.b = '' && self_rel_many:length != 2 && json_object.a.b > 3 && self_rel_many:length <= 4",
			false,
			"SELECT `demo4`.* FROM `demo4` WHERE ((JSON_EXTRACT([[demo4.json_object]], '$.a.b') = '' OR JSON_EXTRACT([[demo4.json_object]], '$.a.b') IS NULL) AND json_array_length(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) != {:TEST} AND JSON_EXTRACT([[demo4.json_object]], '$.a.b') > {:TEST} AND json_array_length(CASE WHEN json_valid([[demo4.self_rel_many]]) THEN [[demo4.self_rel_many]] ELSE json_array([[demo4.self_rel_many]]) END) <= {:TEST})",
		},
	}

	for _, s := range scenarios {
		collection, err := app.Dao().FindCollectionByNameOrId(s.collectionIdOrName)
		if err != nil {
			t.Fatalf("[%s] Failed to load collection %s: %v", s.name, s.collectionIdOrName, err)
		}

		query := app.Dao().RecordQuery(collection)

		r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestData, s.allowHiddenFields)

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

	requestData := &models.RequestData{
		AuthRecord: authRecord,
	}

	r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestData, true)

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
		// json_extract
		{"json_array.0", false, "JSON_EXTRACT([[demo4.json_array]], '$[0]')"},
		{"json_object.a.b.c", false, "JSON_EXTRACT([[demo4.json_object]], '$.a.b.c')"},
		// @request.auth relation join:
		{"@request.auth.rel", false, "[[__auth_users.rel]]"},
		{"@request.auth.rel.title", false, "[[__auth_users_rel.title]]"},
		// @collection fieds:
		{"@collect", true, ""},
		{"collection.demo4.title", true, ""},
		{"@collection", true, ""},
		{"@collection.unknown", true, ""},
		{"@collection.demo2", true, ""},
		{"@collection.demo2.", true, ""},
		{"@collection.demo2.title", false, "[[__collection_demo2.title]]"},
		{"@collection.demo4.title", false, "[[__collection_demo4.title]]"},
		{"@collection.demo4.id", false, "[[__collection_demo4.id]]"},
		{"@collection.demo4.created", false, "[[__collection_demo4.created]]"},
		{"@collection.demo4.updated", false, "[[__collection_demo4.updated]]"},
		{"@collection.demo4.self_rel_many.missing", true, ""},
		{"@collection.demo4.self_rel_many.self_rel_one.self_rel_many.self_rel_one.title", false, "[[__collection_demo4_self_rel_many_self_rel_one_self_rel_many_self_rel_one.title]]"},
	}

	for _, s := range scenarios {
		r, err := r.Resolve(s.fieldName)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%q) Expected hasErr %v, got %v (%v)", s.fieldName, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		if r.Identifier != s.expectName {
			t.Errorf("(%q) Expected r.Identifier %q, got %q", s.fieldName, s.expectName, r.Identifier)
		}

		// params should be empty for non @request fields
		if len(r.Params) != 0 {
			t.Errorf("(%q) Expected 0 r.Params, got %v", s.fieldName, r.Params)
		}
	}
}

func TestRecordFieldResolverResolveStaticRequestDataFields(t *testing.T) {
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

	requestData := &models.RequestData{
		Method: "get",
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

	r := resolvers.NewRecordFieldResolver(app.Dao(), collection, requestData, true)

	scenarios := []struct {
		fieldName        string
		expectError      bool
		expectParamValue string // encoded json
	}{
		{"@request", true, ""},
		{"@request.invalid format", true, ""},
		{"@request.invalid_format2!", true, ""},
		{"@request.missing", true, ""},
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
		{"@request.auth.email", false, `"test@example.com"`},
		{"@request.auth.username", false, `"users75657"`},
		{"@request.auth.verified", false, `false`},
		{"@request.auth.emailVisibility", false, `false`},
		{"@request.auth.missing", false, `NULL`},
	}

	for i, s := range scenarios {
		r, err := r.Resolve(s.fieldName)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		// missing key
		// ---
		if len(r.Params) == 0 {
			if r.Identifier != "NULL" {
				t.Errorf("(%d) Expected 0 placeholder parameters for %v, got %v", i, r.Identifier, r.Params)
			}
			continue
		}

		// existing key
		// ---
		if len(r.Params) != 1 {
			t.Errorf("(%d) Expected 1 placeholder parameter for %v, got %v", i, r.Identifier, r.Params)
			continue
		}

		var paramName string
		var paramValue any
		for k, v := range r.Params {
			paramName = k
			paramValue = v
		}

		if r.Identifier != ("{:" + paramName + "}") {
			t.Errorf("(%d) Expected parameter r.Identifier %q, got %q", i, paramName, r.Identifier)
		}

		encodedParamValue, _ := json.Marshal(paramValue)
		if string(encodedParamValue) != s.expectParamValue {
			t.Errorf("(%d) Expected r.Params %v for %v, got %v", i, s.expectParamValue, r.Identifier, string(encodedParamValue))
		}
	}
}
