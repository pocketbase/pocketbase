package daos_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestRecordQueryWithDifferentCollectionValues(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name          any
		collection    any
		expectedTotal int
		expectError   bool
	}{
		{"with nil value", nil, 0, true},
		{"with invalid or missing collection id/name", "missing", 0, true},
		{"with pointer model", collection, 3, false},
		{"with value model", *collection, 3, false},
		{"with name", "demo1", 3, false},
		{"with id", "wsmn24bux7wo113", 3, false},
	}

	for _, s := range scenarios {
		var records []*models.Record
		err := app.Dao().RecordQuery(s.collection).All(&records)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasError %v, got %v", s.name, s.expectError, hasErr)
			continue
		}

		if total := len(records); total != s.expectedTotal {
			t.Errorf("[%s] Expected %d records, got %d", s.name, s.expectedTotal, total)
		}
	}
}

func TestRecordQueryOneWithRecord(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	id := "84nmscqy84lsi1t"

	q := app.Dao().RecordQuery(collection).
		Where(dbx.HashExp{"id": id})

	record := &models.Record{}
	if err := q.One(record); err != nil {
		t.Fatal(err)
	}

	if record.GetString("id") != id {
		t.Fatalf("Expected record with id %q, got %q", id, record.GetString("id"))
	}
}

func TestRecordQueryAllWithRecordsSlices(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	id1 := "84nmscqy84lsi1t"
	id2 := "al1h9ijdeojtsjy"

	{
		records := []models.Record{}

		q := app.Dao().RecordQuery(collection).
			Where(dbx.HashExp{"id": []any{id1, id2}}).
			OrderBy("created asc")

		if err := q.All(&records); err != nil {
			t.Fatal(err)
		}

		if len(records) != 2 {
			t.Fatalf("Expected %d records, got %d", 2, len(records))
		}

		if records[0].Id != id1 {
			t.Fatalf("Expected record with id %q, got %q", id1, records[0].Id)
		}

		if records[1].Id != id2 {
			t.Fatalf("Expected record with id %q, got %q", id2, records[1].Id)
		}
	}

	{
		records := []*models.Record{}

		q := app.Dao().RecordQuery(collection).
			Where(dbx.HashExp{"id": []any{id1, id2}}).
			OrderBy("created asc")

		if err := q.All(&records); err != nil {
			t.Fatal(err)
		}

		if len(records) != 2 {
			t.Fatalf("Expected %d records, got %d", 2, len(records))
		}

		if records[0].Id != id1 {
			t.Fatalf("Expected record with id %q, got %q", id1, records[0].Id)
		}

		if records[1].Id != id2 {
			t.Fatalf("Expected record with id %q, got %q", id2, records[1].Id)
		}
	}
}

func TestFindRecordById(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		collectionIdOrName string
		id                 string
		filter1            func(q *dbx.SelectQuery) error
		filter2            func(q *dbx.SelectQuery) error
		expectError        bool
	}{
		{"demo2", "missing", nil, nil, true},
		{"missing", "0yxhwia2amd8gec", nil, nil, true},
		{"demo2", "0yxhwia2amd8gec", nil, nil, false},
		{"demo2", "0yxhwia2amd8gec", func(q *dbx.SelectQuery) error {
			q.AndWhere(dbx.HashExp{"title": "missing"})
			return nil
		}, nil, true},
		{"demo2", "0yxhwia2amd8gec", func(q *dbx.SelectQuery) error {
			return errors.New("test error")
		}, nil, true},
		{"demo2", "0yxhwia2amd8gec", func(q *dbx.SelectQuery) error {
			q.AndWhere(dbx.HashExp{"title": "test3"})
			return nil
		}, nil, false},
		{"demo2", "0yxhwia2amd8gec", func(q *dbx.SelectQuery) error {
			q.AndWhere(dbx.HashExp{"title": "test3"})
			return nil
		}, func(q *dbx.SelectQuery) error {
			q.AndWhere(dbx.HashExp{"active": false})
			return nil
		}, true},
		{"sz5l5z67tg7gku0", "0yxhwia2amd8gec", func(q *dbx.SelectQuery) error {
			q.AndWhere(dbx.HashExp{"title": "test3"})
			return nil
		}, func(q *dbx.SelectQuery) error {
			q.AndWhere(dbx.HashExp{"active": true})
			return nil
		}, false},
	}

	for i, scenario := range scenarios {
		record, err := app.Dao().FindRecordById(
			scenario.collectionIdOrName,
			scenario.id,
			scenario.filter1,
			scenario.filter2,
		)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if record != nil && record.Id != scenario.id {
			t.Errorf("(%d) Expected record with id %s, got %s", i, scenario.id, record.Id)
		}
	}
}

func TestFindRecordsByIds(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		collectionIdOrName string
		ids                []string
		filter1            func(q *dbx.SelectQuery) error
		filter2            func(q *dbx.SelectQuery) error
		expectTotal        int
		expectError        bool
	}{
		{"demo2", []string{}, nil, nil, 0, false},
		{"demo2", []string{""}, nil, nil, 0, false},
		{"demo2", []string{"missing"}, nil, nil, 0, false},
		{"missing", []string{"0yxhwia2amd8gec"}, nil, nil, 0, true},
		{"demo2", []string{"0yxhwia2amd8gec"}, nil, nil, 1, false},
		{"sz5l5z67tg7gku0", []string{"0yxhwia2amd8gec"}, nil, nil, 1, false},
		{
			"demo2",
			[]string{"0yxhwia2amd8gec", "llvuca81nly1qls"},
			nil,
			nil,
			2,
			false,
		},
		{
			"demo2",
			[]string{"0yxhwia2amd8gec", "llvuca81nly1qls"},
			func(q *dbx.SelectQuery) error {
				return nil // empty filter
			},
			func(q *dbx.SelectQuery) error {
				return errors.New("test error")
			},
			0,
			true,
		},
		{
			"demo2",
			[]string{"0yxhwia2amd8gec", "llvuca81nly1qls"},
			func(q *dbx.SelectQuery) error {
				q.AndWhere(dbx.HashExp{"active": true})
				return nil
			},
			nil,
			1,
			false,
		},
		{
			"sz5l5z67tg7gku0",
			[]string{"0yxhwia2amd8gec", "llvuca81nly1qls"},
			func(q *dbx.SelectQuery) error {
				q.AndWhere(dbx.HashExp{"active": true})
				return nil
			},
			func(q *dbx.SelectQuery) error {
				q.AndWhere(dbx.Not(dbx.HashExp{"title": ""}))
				return nil
			},
			1,
			false,
		},
	}

	for i, scenario := range scenarios {
		records, err := app.Dao().FindRecordsByIds(
			scenario.collectionIdOrName,
			scenario.ids,
			scenario.filter1,
			scenario.filter2,
		)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if len(records) != scenario.expectTotal {
			t.Errorf("(%d) Expected %d records, got %d", i, scenario.expectTotal, len(records))
			continue
		}

		for _, r := range records {
			if !list.ExistInSlice(r.Id, scenario.ids) {
				t.Errorf("(%d) Couldn't find id %s in %v", i, r.Id, scenario.ids)
			}
		}
	}
}

func TestFindRecordsByExpr(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		collectionIdOrName string
		expressions        []dbx.Expression
		expectIds          []string
		expectError        bool
	}{
		{
			"missing",
			nil,
			[]string{},
			true,
		},
		{
			"demo2",
			nil,
			[]string{
				"achvryl401bhse3",
				"llvuca81nly1qls",
				"0yxhwia2amd8gec",
			},
			false,
		},
		{
			"demo2",
			[]dbx.Expression{
				nil,
				dbx.HashExp{"id": "123"},
			},
			[]string{},
			false,
		},
		{
			"sz5l5z67tg7gku0",
			[]dbx.Expression{
				dbx.Like("title", "test").Match(true, true),
				dbx.HashExp{"active": true},
			},
			[]string{
				"achvryl401bhse3",
				"0yxhwia2amd8gec",
			},
			false,
		},
	}

	for i, scenario := range scenarios {
		records, err := app.Dao().FindRecordsByExpr(scenario.collectionIdOrName, scenario.expressions...)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if len(records) != len(scenario.expectIds) {
			t.Errorf("(%d) Expected %d records, got %d", i, len(scenario.expectIds), len(records))
			continue
		}

		for _, r := range records {
			if !list.ExistInSlice(r.Id, scenario.expectIds) {
				t.Errorf("(%d) Couldn't find id %s in %v", i, r.Id, scenario.expectIds)
			}
		}
	}
}

func TestFindFirstRecordByData(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		collectionIdOrName string
		key                string
		value              any
		expectId           string
		expectError        bool
	}{
		{
			"missing",
			"id",
			"llvuca81nly1qls",
			"llvuca81nly1qls",
			true,
		},
		{
			"demo2",
			"",
			"llvuca81nly1qls",
			"",
			true,
		},
		{
			"demo2",
			"id",
			"invalid",
			"",
			true,
		},
		{
			"demo2",
			"id",
			"llvuca81nly1qls",
			"llvuca81nly1qls",
			false,
		},
		{
			"sz5l5z67tg7gku0",
			"title",
			"test3",
			"0yxhwia2amd8gec",
			false,
		},
	}

	for i, scenario := range scenarios {
		record, err := app.Dao().FindFirstRecordByData(scenario.collectionIdOrName, scenario.key, scenario.value)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if !scenario.expectError && record.Id != scenario.expectId {
			t.Errorf("(%d) Expected record with id %s, got %v", i, scenario.expectId, record.Id)
		}
	}
}

func TestFindRecordsByFilter(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name               string
		collectionIdOrName string
		filter             string
		sort               string
		limit              int
		offset             int
		params             []dbx.Params
		expectError        bool
		expectRecordIds    []string
	}{
		{
			"missing collection",
			"missing",
			"id != ''",
			"",
			0,
			0,
			nil,
			true,
			nil,
		},
		{
			"missing filter",
			"demo2",
			"",
			"",
			0,
			0,
			nil,
			true,
			nil,
		},
		{
			"invalid filter",
			"demo2",
			"someMissingField > 1",
			"",
			0,
			0,
			nil,
			true,
			nil,
		},
		{
			"simple filter",
			"demo2",
			"id != ''",
			"",
			0,
			0,
			nil,
			false,
			[]string{
				"llvuca81nly1qls",
				"achvryl401bhse3",
				"0yxhwia2amd8gec",
			},
		},
		{
			"multi-condition filter with sort",
			"demo2",
			"id != '' && active=true",
			"-created,title",
			-1, // should behave the same as 0
			0,
			nil,
			false,
			[]string{
				"0yxhwia2amd8gec",
				"achvryl401bhse3",
			},
		},
		{
			"with limit and offset",
			"demo2",
			"id != ''",
			"title",
			2,
			1,
			nil,
			false,
			[]string{
				"achvryl401bhse3",
				"0yxhwia2amd8gec",
			},
		},
		{
			"with placeholder params",
			"demo2",
			"active = {:active}",
			"",
			10,
			0,
			[]dbx.Params{{"active": false}},
			false,
			[]string{
				"llvuca81nly1qls",
			},
		},
		{
			"with json filter and sort",
			"demo4",
			"json_object != null && json_object.a.b = 'test'",
			"-json_object.a",
			10,
			0,
			[]dbx.Params{{"active": false}},
			false,
			[]string{
				"i9naidtvr6qsgb4",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			records, err := app.Dao().FindRecordsByFilter(
				s.collectionIdOrName,
				s.filter,
				s.sort,
				s.limit,
				s.offset,
				s.params...,
			)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("[%s] Expected hasErr to be %v, got %v (%v)", s.name, s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if len(records) != len(s.expectRecordIds) {
				t.Fatalf("[%s] Expected %d records, got %d", s.name, len(s.expectRecordIds), len(records))
			}

			for i, id := range s.expectRecordIds {
				if id != records[i].Id {
					t.Fatalf("[%s] Expected record with id %q, got %q at index %d", s.name, id, records[i].Id, i)
				}
			}
		})
	}
}

func TestFindFirstRecordByFilter(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name               string
		collectionIdOrName string
		filter             string
		params             []dbx.Params
		expectError        bool
		expectRecordId     string
	}{
		{
			"missing collection",
			"missing",
			"id != ''",
			nil,
			true,
			"",
		},
		{
			"missing filter",
			"demo2",
			"",
			nil,
			true,
			"",
		},
		{
			"invalid filter",
			"demo2",
			"someMissingField > 1",
			nil,
			true,
			"",
		},
		{
			"valid filter but no matches",
			"demo2",
			"id = 'test'",
			nil,
			true,
			"",
		},
		{
			"valid filter and multiple matches",
			"demo2",
			"id != ''",
			nil,
			false,
			"llvuca81nly1qls",
		},
		{
			"with placeholder params",
			"demo2",
			"active = {:active}",
			[]dbx.Params{{"active": false}},
			false,
			"llvuca81nly1qls",
		},
	}

	for _, s := range scenarios {
		record, err := app.Dao().FindFirstRecordByFilter(s.collectionIdOrName, s.filter, s.params...)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr to be %v, got %v (%v)", s.name, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		if record.Id != s.expectRecordId {
			t.Errorf("[%s] Expected record with id %q, got %q", s.name, s.expectRecordId, record.Id)
		}
	}
}

func TestCanAccessRecord(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	admin, err := app.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	authRecord, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	record, err := app.Dao().FindRecordById("demo1", "imy661ixudk5izi")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name        string
		record      *models.Record
		requestInfo *models.RequestInfo
		rule        *string
		expected    bool
		expectError bool
	}{
		{
			"as admin with nil rule",
			record,
			&models.RequestInfo{
				Admin: admin,
			},
			nil,
			true,
			false,
		},
		{
			"as admin with non-empty rule",
			record,
			&models.RequestInfo{
				Admin: admin,
			},
			types.Pointer("id = ''"), // the filter rule should be ignored
			true,
			false,
		},
		{
			"as admin with invalid rule",
			record,
			&models.RequestInfo{
				Admin: admin,
			},
			types.Pointer("id ?!@ 1"), // the filter rule should be ignored
			true,
			false,
		},
		{
			"as guest with nil rule",
			record,
			&models.RequestInfo{},
			nil,
			false,
			false,
		},
		{
			"as guest with empty rule",
			record,
			&models.RequestInfo{},
			types.Pointer(""),
			true,
			false,
		},
		{
			"as guest with invalid rule",
			record,
			&models.RequestInfo{},
			types.Pointer("id ?!@ 1"),
			false,
			true,
		},
		{
			"as guest with mismatched rule",
			record,
			&models.RequestInfo{},
			types.Pointer("@request.auth.id != ''"),
			false,
			false,
		},
		{
			"as guest with matched rule",
			record,
			&models.RequestInfo{
				Data: map[string]any{"test": 1},
			},
			types.Pointer("@request.auth.id != '' || @request.data.test = 1"),
			true,
			false,
		},
		{
			"as auth record with nil rule",
			record,
			&models.RequestInfo{
				AuthRecord: authRecord,
			},
			nil,
			false,
			false,
		},
		{
			"as auth record with empty rule",
			record,
			&models.RequestInfo{
				AuthRecord: authRecord,
			},
			types.Pointer(""),
			true,
			false,
		},
		{
			"as auth record with invalid rule",
			record,
			&models.RequestInfo{
				AuthRecord: authRecord,
			},
			types.Pointer("id ?!@ 1"),
			false,
			true,
		},
		{
			"as auth record with mismatched rule",
			record,
			&models.RequestInfo{
				AuthRecord: authRecord,
				Data:       map[string]any{"test": 1},
			},
			types.Pointer("@request.auth.id != '' && @request.data.test > 1"),
			false,
			false,
		},
		{
			"as auth record with matched rule",
			record,
			&models.RequestInfo{
				AuthRecord: authRecord,
				Data:       map[string]any{"test": 2},
			},
			types.Pointer("@request.auth.id != '' && @request.data.test > 1"),
			true,
			false,
		},
	}

	for _, s := range scenarios {
		result, err := app.Dao().CanAccessRecord(s.record, s.requestInfo, s.rule)

		if result != s.expected {
			t.Errorf("[%s] Expected %v, got %v", s.name, s.expected, result)
		}

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr %v, got %v (%v)", s.name, s.expectError, hasErr, err)
		}
	}
}

func TestIsRecordValueUnique(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	testManyRelsId1 := "bgs820n361vj1qd"
	testManyRelsId2 := "4q1xlclmfloku33"
	testManyRelsId3 := "oap640cot4yru2s"

	scenarios := []struct {
		collectionIdOrName string
		key                string
		value              any
		excludeIds         []string
		expected           bool
	}{
		{"demo2", "", "", nil, false},
		{"demo2", "", "", []string{""}, false},
		{"demo2", "missing", "unique", nil, false},
		{"demo2", "title", "unique", nil, true},
		{"demo2", "title", "unique", []string{}, true},
		{"demo2", "title", "unique", []string{""}, true},
		{"demo2", "title", "test1", []string{""}, false},
		{"demo2", "title", "test1", []string{"llvuca81nly1qls"}, true},
		{"demo1", "rel_many", []string{testManyRelsId3}, nil, false},
		{"wsmn24bux7wo113", "rel_many", []any{testManyRelsId3}, []string{""}, false},
		{"wsmn24bux7wo113", "rel_many", []any{testManyRelsId3}, []string{"84nmscqy84lsi1t"}, true},
		// mixed json array order
		{"demo1", "rel_many", []string{testManyRelsId1, testManyRelsId3, testManyRelsId2}, nil, true},
		// username special case-insensitive match
		{"users", "username", "test2_username", nil, false},
		{"users", "username", "TEST2_USERNAME", nil, false},
		{"users", "username", "new_username", nil, true},
		{"users", "username", "TEST2_USERNAME", []string{"oap640cot4yru2s"}, true},
	}

	for i, scenario := range scenarios {
		result := app.Dao().IsRecordValueUnique(
			scenario.collectionIdOrName,
			scenario.key,
			scenario.value,
			scenario.excludeIds...,
		)

		if result != scenario.expected {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expected, result)
		}
	}
}

func TestFindAuthRecordByToken(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		token         string
		baseKey       string
		expectedEmail string
		expectError   bool
	}{
		// invalid auth token
		{
			"eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.H2KKcIXiAfxvuXMFzizo1SgsinDP4hcWhD3pYoP4Nqw",
			app.Settings().RecordAuthToken.Secret,
			"",
			true,
		},
		// expired token
		{
			"eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoxNjQwOTkxNjYxfQ.HqvpCpM0RAk3Qu9PfCMuZsk_DKh9UYuzFLwXBMTZd1w",
			app.Settings().RecordAuthToken.Secret,
			"",
			true,
		},
		// wrong base key (password reset token secret instead of auth secret)
		{
			"eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			app.Settings().RecordPasswordResetToken.Secret,
			"",
			true,
		},
		// valid token and base key but with deleted/missing collection
		{
			"eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoibWlzc2luZyIsImV4cCI6MjIwODk4NTI2MX0.0oEHQpdpHp0Nb3VN8La0ssg-SjwWKiRl_k1mUGxdKlU",
			app.Settings().RecordAuthToken.Secret,
			"test@example.com",
			true,
		},
		// valid token
		{
			"eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoUmVjb3JkIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyMjA4OTg1MjYxfQ.UwD8JvkbQtXpymT09d7J6fdA0aP9g4FJ1GPh_ggEkzc",
			app.Settings().RecordAuthToken.Secret,
			"test@example.com",
			false,
		},
	}

	for i, scenario := range scenarios {
		record, err := app.Dao().FindAuthRecordByToken(scenario.token, scenario.baseKey)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if !scenario.expectError && record.Email() != scenario.expectedEmail {
			t.Errorf("(%d) Expected record model %s, got %s", i, scenario.expectedEmail, record.Email())
		}
	}
}

func TestFindAuthRecordByEmail(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		collectionIdOrName string
		email              string
		expectError        bool
	}{
		{"missing", "test@example.com", true},
		{"demo2", "test@example.com", true},
		{"users", "missing@example.com", true},
		{"users", "test@example.com", false},
		{"clients", "test2@example.com", false},
	}

	for i, scenario := range scenarios {
		record, err := app.Dao().FindAuthRecordByEmail(scenario.collectionIdOrName, scenario.email)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if !scenario.expectError && record.Email() != scenario.email {
			t.Errorf("(%d) Expected record with email %s, got %s", i, scenario.email, record.Email())
		}
	}
}

func TestFindAuthRecordByUsername(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		collectionIdOrName string
		username           string
		expectError        bool
	}{
		{"missing", "test_username", true},
		{"demo2", "test_username", true},
		{"users", "missing", true},
		{"users", "test2_username", false},
		{"users", "TEST2_USERNAME", false}, // case insensitive check
		{"clients", "clients43362", false},
	}

	for i, scenario := range scenarios {
		record, err := app.Dao().FindAuthRecordByUsername(scenario.collectionIdOrName, scenario.username)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if !scenario.expectError && !strings.EqualFold(record.Username(), scenario.username) {
			t.Errorf("(%d) Expected record with username %s, got %s", i, scenario.username, record.Username())
		}
	}
}

func TestSuggestUniqueAuthRecordUsername(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		collectionIdOrName string
		baseUsername       string
		expectedPattern    string
	}{
		// missing collection
		{"missing", "test2_username", `^test2_username\d{12}$`},
		// not an auth collection
		{"demo2", "test2_username", `^test2_username\d{12}$`},
		// auth collection with unique base username
		{"users", "new_username", `^new_username$`},
		{"users", "NEW_USERNAME", `^NEW_USERNAME$`},
		// auth collection with existing username
		{"users", "test2_username", `^test2_username\d{3}$`},
		{"users", "TEST2_USERNAME", `^TEST2_USERNAME\d{3}$`},
	}

	for i, scenario := range scenarios {
		username := app.Dao().SuggestUniqueAuthRecordUsername(
			scenario.collectionIdOrName,
			scenario.baseUsername,
		)

		pattern, err := regexp.Compile(scenario.expectedPattern)
		if err != nil {
			t.Errorf("[%d] Invalid username pattern %q: %v", i, scenario.expectedPattern, err)
		}
		if !pattern.MatchString(username) {
			t.Fatalf("Expected username to match %s, got username %s", scenario.expectedPattern, username)
		}
	}
}

func TestSaveRecord(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, _ := app.Dao().FindCollectionByNameOrId("demo2")

	// create
	// ---
	r1 := models.NewRecord(collection)
	r1.Set("title", "test_new")
	err1 := app.Dao().SaveRecord(r1)
	if err1 != nil {
		t.Fatal(err1)
	}
	newR1, _ := app.Dao().FindFirstRecordByData(collection.Id, "title", "test_new")
	if newR1 == nil || newR1.Id != r1.Id || newR1.GetString("title") != r1.GetString("title") {
		t.Fatalf("Expected to find record %v, got %v", r1, newR1)
	}

	// update
	// ---
	r2, _ := app.Dao().FindFirstRecordByData(collection.Id, "id", "0yxhwia2amd8gec")
	r2.Set("title", "test_update")
	err2 := app.Dao().SaveRecord(r2)
	if err2 != nil {
		t.Fatal(err2)
	}
	newR2, _ := app.Dao().FindFirstRecordByData(collection.Id, "title", "test_update")
	if newR2 == nil || newR2.Id != r2.Id || newR2.GetString("title") != r2.GetString("title") {
		t.Fatalf("Expected to find record %v, got %v", r2, newR2)
	}
}

func TestSaveRecordWithIdFromOtherCollection(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	baseCollection, _ := app.Dao().FindCollectionByNameOrId("demo2")
	authCollection, _ := app.Dao().FindCollectionByNameOrId("nologin")

	// base collection test
	r1 := models.NewRecord(baseCollection)
	r1.Set("title", "test_new")
	r1.Set("id", "mk5fmymtx4wsprk") // existing id of demo3 record
	r1.MarkAsNew()
	if err := app.Dao().SaveRecord(r1); err != nil {
		t.Fatalf("Expected nil, got error %v", err)
	}

	// auth collection test
	r2 := models.NewRecord(authCollection)
	r2.Set("username", "test_new")
	r2.Set("id", "gk390qegs4y47wn") // existing id of "clients" record
	r2.MarkAsNew()
	if err := app.Dao().SaveRecord(r2); err == nil {
		t.Fatal("Expected error, got nil")
	}

	// try again with unique id
	r2.Set("id", "unique_id")
	if err := app.Dao().SaveRecord(r2); err != nil {
		t.Fatalf("Expected nil, got error %v", err)
	}
}

func TestDeleteRecord(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demoCollection, _ := app.Dao().FindCollectionByNameOrId("demo2")

	// delete unsaved record
	// ---
	rec0 := models.NewRecord(demoCollection)
	if err := app.Dao().DeleteRecord(rec0); err == nil {
		t.Fatal("(rec0) Didn't expect to succeed deleting unsaved record")
	}

	// delete existing record + external auths
	// ---
	rec1, _ := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err := app.Dao().DeleteRecord(rec1); err != nil {
		t.Fatalf("(rec1) Expected nil, got error %v", err)
	}
	// check if it was really deleted
	if refreshed, _ := app.Dao().FindRecordById(rec1.Collection().Id, rec1.Id); refreshed != nil {
		t.Fatalf("(rec1) Expected record to be deleted, got %v", refreshed)
	}
	// check if the external auths were deleted
	if auths, _ := app.Dao().FindAllExternalAuthsByRecord(rec1); len(auths) > 0 {
		t.Fatalf("(rec1) Expected external auths to be deleted, got %v", auths)
	}

	// delete existing record while being part of a non-cascade required relation
	// ---
	rec2, _ := app.Dao().FindRecordById("demo3", "7nwo8tuiatetxdm")
	if err := app.Dao().DeleteRecord(rec2); err == nil {
		t.Fatalf("(rec2) Expected error, got nil")
	}

	// delete existing record + cascade
	// ---
	calledQueries := []string{}
	app.Dao().NonconcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		calledQueries = append(calledQueries, sql)
	}
	app.Dao().ConcurrentDB().(*dbx.DB).QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		calledQueries = append(calledQueries, sql)
	}
	app.Dao().NonconcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		calledQueries = append(calledQueries, sql)
	}
	app.Dao().ConcurrentDB().(*dbx.DB).ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		calledQueries = append(calledQueries, sql)
	}
	rec3, _ := app.Dao().FindRecordById("users", "oap640cot4yru2s")
	// delete
	if err := app.Dao().DeleteRecord(rec3); err != nil {
		t.Fatalf("(rec3) Expected nil, got error %v", err)
	}
	// check if it was really deleted
	rec3, _ = app.Dao().FindRecordById(rec3.Collection().Id, rec3.Id)
	if rec3 != nil {
		t.Fatalf("(rec3) Expected record to be deleted, got %v", rec3)
	}
	// check if the operation cascaded
	rel, _ := app.Dao().FindRecordById("demo1", "84nmscqy84lsi1t")
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

func TestDeleteRecordBatchProcessing(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	if err := createMockBatchProcessingData(app.Dao()); err != nil {
		t.Fatal(err)
	}

	// find and delete the first c1 record to trigger cascade
	mainRecord, _ := app.Dao().FindRecordById("c1", "a")
	if err := app.Dao().DeleteRecord(mainRecord); err != nil {
		t.Fatal(err)
	}

	// check if the main record was deleted
	_, err := app.Dao().FindRecordById(mainRecord.Collection().Id, mainRecord.Id)
	if err == nil {
		t.Fatal("The main record wasn't deleted")
	}

	// check if the c1 b rel field were updated
	c1RecordB, err := app.Dao().FindRecordById("c1", "b")
	if err != nil || c1RecordB.GetString("rel") != "" {
		t.Fatalf("Expected c1RecordB.rel to be nil, got %v", c1RecordB.GetString("rel"))
	}

	// check if the c2 rel fields were updated
	c2Records, err := app.Dao().FindRecordsByExpr("c2", nil)
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
	c3Records, err := app.Dao().FindRecordsByExpr("c3", nil)
	if err != nil {
		t.Fatalf("Failed to fetch c3 records: %v", err)
	}
	if total := len(c3Records); total != 0 {
		t.Fatalf("Expected c3 records to be deleted, found %d", total)
	}
}

func createMockBatchProcessingData(dao *daos.Dao) error {
	// create mock collection without relation
	c1 := &models.Collection{}
	c1.Id = "c1"
	c1.Name = c1.Id
	c1.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "text",
			Type: schema.FieldTypeText,
		},
		// self reference
		&schema.SchemaField{
			Name: "rel",
			Type: schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:     types.Pointer(1),
				CollectionId:  "c1",
				CascadeDelete: false, // should unset all rel fields
			},
		},
	)
	if err := dao.SaveCollection(c1); err != nil {
		return err
	}

	// create mock collection with a multi-rel field
	c2 := &models.Collection{}
	c2.Id = "c2"
	c2.Name = c2.Id
	c2.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "rel",
			Type: schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:     types.Pointer(10),
				CollectionId:  "c1",
				CascadeDelete: false, // should unset all rel fields
			},
		},
	)
	if err := dao.SaveCollection(c2); err != nil {
		return err
	}

	// create mock collection with a single-rel field
	c3 := &models.Collection{}
	c3.Id = "c3"
	c3.Name = c3.Id
	c3.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name: "rel",
			Type: schema.FieldTypeRelation,
			Options: &schema.RelationOptions{
				MaxSelect:     types.Pointer(1),
				CollectionId:  "c1",
				CascadeDelete: true, // should delete all c3 records
			},
		},
	)
	if err := dao.SaveCollection(c3); err != nil {
		return err
	}

	// insert mock records
	c1RecordA := models.NewRecord(c1)
	c1RecordA.Id = "a"
	c1RecordA.Set("rel", c1RecordA.Id) // self reference
	if err := dao.Save(c1RecordA); err != nil {
		return err
	}
	c1RecordB := models.NewRecord(c1)
	c1RecordB.Id = "b"
	c1RecordB.Set("rel", c1RecordA.Id) // rel to another record from the same collection
	if err := dao.Save(c1RecordB); err != nil {
		return err
	}
	for i := 0; i < 4500; i++ {
		c2Record := models.NewRecord(c2)
		c2Record.Set("rel", []string{c1RecordA.Id, c1RecordB.Id})
		if err := dao.Save(c2Record); err != nil {
			return err
		}

		c3Record := models.NewRecord(c3)
		c3Record.Set("rel", c1RecordA.Id)
		if err := dao.Save(c3Record); err != nil {
			return err
		}
	}

	// set the same id as the relation for at least 1 record
	// to check whether the correct condition will be added
	c3Record := models.NewRecord(c3)
	c3Record.Set("rel", c1RecordA.Id)
	c3Record.Id = c1RecordA.Id
	if err := dao.Save(c3Record); err != nil {
		return err
	}

	return nil
}
