package core_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestFindAllMFAsByRecord(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	if err := tests.StubMFARecords(app); err != nil {
		t.Fatal(err)
	}

	demo1, err := app.FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	superuser2, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test2@example.com")
	if err != nil {
		t.Fatal(err)
	}

	superuser4, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test4@example.com")
	if err != nil {
		t.Fatal(err)
	}

	user1, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		record   *core.Record
		expected []string
	}{
		{demo1, nil},
		{superuser2, []string{"superuser2_0", "superuser2_3", "superuser2_2", "superuser2_1", "superuser2_4"}},
		{superuser4, nil},
		{user1, []string{"user1_0"}},
	}

	for _, s := range scenarios {
		t.Run(s.record.Collection().Name+"_"+s.record.Id, func(t *testing.T) {
			result, err := app.FindAllMFAsByRecord(s.record)
			if err != nil {
				t.Fatal(err)
			}

			if len(result) != len(s.expected) {
				t.Fatalf("Expected total mfas %d, got %d", len(s.expected), len(result))
			}

			for i, id := range s.expected {
				if result[i].Id != id {
					t.Errorf("[%d] Expected id %q, got %q", i, id, result[i].Id)
				}
			}
		})
	}
}

func TestFindAllMFAsByCollection(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	if err := tests.StubMFARecords(app); err != nil {
		t.Fatal(err)
	}

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
		{superusers, []string{
			"superuser2_0",
			"superuser2_3",
			"superuser3_0",
			"superuser2_2",
			"superuser3_1",
			"superuser2_1",
			"superuser2_4",
		}},
		{clients, nil},
		{users, []string{"user1_0"}},
	}

	for _, s := range scenarios {
		t.Run(s.collection.Name, func(t *testing.T) {
			result, err := app.FindAllMFAsByCollection(s.collection)
			if err != nil {
				t.Fatal(err)
			}

			if len(result) != len(s.expected) {
				t.Fatalf("Expected total mfas %d, got %d", len(s.expected), len(result))
			}

			for i, id := range s.expected {
				if result[i].Id != id {
					t.Errorf("[%d] Expected id %q, got %q", i, id, result[i].Id)
				}
			}
		})
	}
}

func TestFindMFAById(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	if err := tests.StubMFARecords(app); err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		id          string
		expectError bool
	}{
		{"", true},
		{"84nmscqy84lsi1t", true}, // non-mfa id
		{"superuser2_0", false},
		{"superuser2_4", false}, // expired
		{"user1_0", false},
	}

	for _, s := range scenarios {
		t.Run(s.id, func(t *testing.T) {
			result, err := app.FindMFAById(s.id)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if result.Id != s.id {
				t.Fatalf("Expected record with id %q, got %q", s.id, result.Id)
			}
		})
	}
}

func TestDeleteAllMFAsByRecord(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	demo1, err := testApp.FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	superuser2, err := testApp.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test2@example.com")
	if err != nil {
		t.Fatal(err)
	}

	superuser4, err := testApp.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test4@example.com")
	if err != nil {
		t.Fatal(err)
	}

	user1, err := testApp.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		record     *core.Record
		deletedIds []string
	}{
		{demo1, nil}, // non-auth record
		{superuser2, []string{"superuser2_0", "superuser2_1", "superuser2_3", "superuser2_2", "superuser2_4"}},
		{superuser4, nil},
		{user1, []string{"user1_0"}},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s_%s", i, s.record.Collection().Name, s.record.Id), func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			if err := tests.StubMFARecords(app); err != nil {
				t.Fatal(err)
			}

			deletedIds := []string{}
			app.OnRecordAfterDeleteSuccess().BindFunc(func(e *core.RecordEvent) error {
				deletedIds = append(deletedIds, e.Record.Id)
				return e.Next()
			})

			err := app.DeleteAllMFAsByRecord(s.record)
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

func TestDeleteExpiredMFAs(t *testing.T) {
	t.Parallel()

	checkDeletedIds := func(app core.App, t *testing.T, expectedDeletedIds []string) {
		if err := tests.StubMFARecords(app); err != nil {
			t.Fatal(err)
		}

		deletedIds := []string{}
		app.OnRecordDelete().BindFunc(func(e *core.RecordEvent) error {
			deletedIds = append(deletedIds, e.Record.Id)
			return e.Next()
		})

		if err := app.DeleteExpiredMFAs(); err != nil {
			t.Fatal(err)
		}

		if len(deletedIds) != len(expectedDeletedIds) {
			t.Fatalf("Expected deleted ids\n%v\ngot\n%v", expectedDeletedIds, deletedIds)
		}

		for _, id := range expectedDeletedIds {
			if !slices.Contains(deletedIds, id) {
				t.Errorf("Expected to find deleted id %q in %v", id, deletedIds)
			}
		}
	}

	t.Run("default test collections", func(t *testing.T) {
		app, _ := tests.NewTestApp()
		defer app.Cleanup()

		checkDeletedIds(app, t, []string{
			"user1_0",
			"superuser2_1",
			"superuser2_4",
		})
	})

	t.Run("mfa collection duration mock", func(t *testing.T) {
		app, _ := tests.NewTestApp()
		defer app.Cleanup()

		superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		if err != nil {
			t.Fatal(err)
		}
		superusers.MFA.Duration = 60
		if err := app.Save(superusers); err != nil {
			t.Fatalf("Failed to mock superusers mfa duration: %v", err)
		}

		checkDeletedIds(app, t, []string{
			"user1_0",
			"superuser2_1",
			"superuser2_2",
			"superuser2_4",
			"superuser3_1",
		})
	})
}
