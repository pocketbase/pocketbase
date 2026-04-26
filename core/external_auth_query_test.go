package core_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestFindAllExternalAuthsByRecord(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demo1, err := app.FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	superuser1, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	user1, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	user2, err := app.FindAuthRecordByEmail("users", "test2@example.com")
	if err != nil {
		t.Fatal(err)
	}

	user3, err := app.FindAuthRecordByEmail("users", "test3@example.com")
	if err != nil {
		t.Fatal(err)
	}

	client1, err := app.FindAuthRecordByEmail("clients", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		record   *core.Record
		expected []string
	}{
		{demo1, nil},
		{superuser1, nil},
		{client1, []string{"f1z5b3843pzc964"}},
		{user1, []string{"clmflokuq1xl341", "dlmflokuq1xl342"}},
		{user2, nil},
		{user3, []string{"5eto7nmys833164"}},
	}

	for _, s := range scenarios {
		t.Run(s.record.Collection().Name+"_"+s.record.Id, func(t *testing.T) {
			result, err := app.FindAllExternalAuthsByRecord(s.record)
			if err != nil {
				t.Fatal(err)
			}

			if len(result) != len(s.expected) {
				t.Fatalf("Expected total models %d, got %d", len(s.expected), len(result))
			}

			for i, id := range s.expected {
				if result[i].Id != id {
					t.Errorf("[%d] Expected id %q, got %q", i, id, result[i].Id)
				}
			}
		})
	}
}

func TestFindAllExternalAuthsByCollection(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demo1, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
	if err != nil {
		t.Fatal(err)
	}

	clients, err := app.FindCollectionByNameOrId("clients")
	if err != nil {
		t.Fatal(err)
	}

	users, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		collection *core.Collection
		expected   []string
	}{
		{demo1, nil},
		{superusers, nil},
		{clients, []string{
			"f1z5b3843pzc964",
		}},
		{users, []string{
			"5eto7nmys833164",
			"clmflokuq1xl341",
			"dlmflokuq1xl342",
		}},
	}

	for _, s := range scenarios {
		t.Run(s.collection.Name, func(t *testing.T) {
			result, err := app.FindAllExternalAuthsByCollection(s.collection)
			if err != nil {
				t.Fatal(err)
			}

			if len(result) != len(s.expected) {
				t.Fatalf("Expected total models %d, got %d", len(s.expected), len(result))
			}

			for i, id := range s.expected {
				if result[i].Id != id {
					t.Errorf("[%d] Expected id %q, got %q", i, id, result[i].Id)
				}
			}
		})
	}
}

func TestFindFirstExternalAuthByExpr(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		expr       dbx.Expression
		expectedId string
	}{
		{dbx.HashExp{"collectionRef": "invalid"}, ""},
		{dbx.HashExp{"collectionRef": "_pb_users_auth_"}, "5eto7nmys833164"},
		{dbx.HashExp{"collectionRef": "_pb_users_auth_", "provider": "gitlab"}, "dlmflokuq1xl342"},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%v", i, s.expr.Build(app.ConcurrentDB().(*dbx.DB), dbx.Params{})), func(t *testing.T) {
			result, err := app.FindFirstExternalAuthByExpr(s.expr)

			hasErr := err != nil
			expectErr := s.expectedId == ""
			if hasErr != expectErr {
				t.Fatalf("Expected hasErr %v, got %v", expectErr, hasErr)
			}

			if hasErr {
				return
			}

			if result.Id != s.expectedId {
				t.Errorf("Expected id %q, got %q", s.expectedId, result.Id)
			}
		})
	}
}

func TestDeleteAllExternalAuthsByRecord(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	demo1, err := testApp.FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	user1, err := testApp.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	client1, err := testApp.FindAuthRecordByEmail("clients", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	client2, err := testApp.FindAuthRecordByEmail("clients", "test2@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		record     *core.Record
		deletedIds []string
	}{
		{demo1, nil}, // non-auth record
		{user1, []string{"dlmflokuq1xl342", "clmflokuq1xl341"}},
		{client1, []string{"f1z5b3843pzc964"}},
		{client2, nil},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s_%s", i, s.record.Collection().Name, s.record.Id), func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			deletedIds := []string{}
			app.OnRecordDelete().BindFunc(func(e *core.RecordEvent) error {
				deletedIds = append(deletedIds, e.Record.Id)
				return e.Next()
			})

			err := app.DeleteAllExternalAuthsByRecord(s.record)
			if err != nil {
				t.Fatal(err)
			}

			if len(deletedIds) != len(s.deletedIds) {
				t.Fatalf("Expected deleted ids\n%v\ngot\n%v", s.deletedIds, deletedIds)
			}

			for _, id := range s.deletedIds {
				if !slices.Contains(deletedIds, id) {
					t.Errorf("Expected to find deleted id %q in %v", id, deletedIds)
				}
			}
		})
	}
}
