package daos_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestRecordQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.Dao().FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("SELECT `%s`.* FROM `%s`", collection.Name, collection.Name)

	sql := app.Dao().RecordQuery(collection).Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindRecordById(t *testing.T) {
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

func TestIsRecordValueUnique(t *testing.T) {
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
	rec3, _ := app.Dao().FindRecordById("users", "oap640cot4yru2s")
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
}

func TestSyncRecordTableSchema(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	oldCollection, err := app.Dao().FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}
	updatedCollection, err := app.Dao().FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}
	updatedCollection.Name = "demo_renamed"
	updatedCollection.Schema.RemoveField(updatedCollection.Schema.GetFieldByName("active").Id)
	updatedCollection.Schema.AddField(
		&schema.SchemaField{
			Name: "new_field",
			Type: schema.FieldTypeEmail,
		},
	)
	updatedCollection.Schema.AddField(
		&schema.SchemaField{
			Id:   updatedCollection.Schema.GetFieldByName("title").Id,
			Name: "title_renamed",
			Type: schema.FieldTypeEmail,
		},
	)

	scenarios := []struct {
		newCollection     *models.Collection
		oldCollection     *models.Collection
		expectedTableName string
		expectedColumns   []string
	}{
		// new base collection
		{
			&models.Collection{
				Name: "new_table",
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name: "test",
						Type: schema.FieldTypeText,
					},
				),
			},
			nil,
			"new_table",
			[]string{"id", "created", "updated", "test"},
		},
		// new auth collection
		{
			&models.Collection{
				Name: "new_table_auth",
				Type: models.CollectionTypeAuth,
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name: "test",
						Type: schema.FieldTypeText,
					},
				),
			},
			nil,
			"new_table_auth",
			[]string{
				"id", "created", "updated", "test",
				"username", "email", "verified", "emailVisibility",
				"tokenKey", "passwordHash", "lastResetSentAt", "lastVerificationSentAt",
			},
		},
		// no changes
		{
			oldCollection,
			oldCollection,
			"demo3",
			[]string{"id", "created", "updated", "title", "active"},
		},
		// renamed table, deleted column, renamed columnd and new column
		{
			updatedCollection,
			oldCollection,
			"demo_renamed",
			[]string{"id", "created", "updated", "title_renamed", "new_field"},
		},
	}

	for i, scenario := range scenarios {
		err := app.Dao().SyncRecordTableSchema(scenario.newCollection, scenario.oldCollection)
		if err != nil {
			t.Errorf("(%d) %v", i, err)
			continue
		}

		if !app.Dao().HasTable(scenario.newCollection.Name) {
			t.Errorf("(%d) Expected table %s to exist", i, scenario.newCollection.Name)
		}

		cols, _ := app.Dao().GetTableColumns(scenario.newCollection.Name)
		if len(cols) != len(scenario.expectedColumns) {
			t.Errorf("(%d) Expected columns %v, got %v", i, scenario.expectedColumns, cols)
		}

		for _, c := range cols {
			if !list.ExistInSlice(c, scenario.expectedColumns) {
				t.Errorf("(%d) Couldn't find column %s in %v", i, c, scenario.expectedColumns)
			}
		}
	}
}
