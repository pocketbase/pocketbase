package daos_test

import (
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestExternalAuthQuery(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	expected := "SELECT {{_externalAuths}}.* FROM `_externalAuths`"

	sql := app.Dao().ExternalAuthQuery().Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindAllExternalAuthsByRecord(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		userId        string
		expectedCount int
	}{
		{"oap640cot4yru2s", 0},
		{"4q1xlclmfloku33", 2},
	}

	for i, s := range scenarios {
		record, err := app.Dao().FindRecordById("users", s.userId)
		if err != nil {
			t.Errorf("(%d) Unexpected record fetch error %v", i, err)
			continue
		}

		auths, err := app.Dao().FindAllExternalAuthsByRecord(record)
		if err != nil {
			t.Errorf("(%d) Unexpected auths fetch error %v", i, err)
			continue
		}

		if len(auths) != s.expectedCount {
			t.Errorf("(%d) Expected %d auths, got %d", i, s.expectedCount, len(auths))
		}

		for _, auth := range auths {
			if auth.RecordId != record.Id {
				t.Errorf("(%d) Expected all auths to be linked to record id %s, got %v", i, record.Id, auth)
			}
		}
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
		{dbx.HashExp{"provider": "github", "providerId": ""}, ""},
		{dbx.HashExp{"provider": "github", "providerId": "id1"}, ""},
		{dbx.HashExp{"provider": "github", "providerId": "id2"}, ""},
		{dbx.HashExp{"provider": "google", "providerId": "test123"}, "clmflokuq1xl341"},
		{dbx.HashExp{"provider": "gitlab", "providerId": "test123"}, "dlmflokuq1xl342"},
	}

	for i, s := range scenarios {
		auth, err := app.Dao().FindFirstExternalAuthByExpr(s.expr)

		hasErr := err != nil
		expectErr := s.expectedId == ""
		if hasErr != expectErr {
			t.Errorf("(%d) Expected hasErr %v, got %v", i, expectErr, err)
			continue
		}

		if auth != nil && auth.Id != s.expectedId {
			t.Errorf("(%d) Expected external auth with ID %s, got \n%v", i, s.expectedId, auth)
		}
	}
}

func TestFindExternalAuthByRecordAndProvider(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		userId     string
		provider   string
		expectedId string
	}{
		{"bgs820n361vj1qd", "google", ""},
		{"4q1xlclmfloku33", "google", "clmflokuq1xl341"},
		{"4q1xlclmfloku33", "gitlab", "dlmflokuq1xl342"},
	}

	for i, s := range scenarios {
		record, err := app.Dao().FindRecordById("users", s.userId)
		if err != nil {
			t.Errorf("(%d) Unexpected record fetch error %v", i, err)
			continue
		}

		auth, err := app.Dao().FindExternalAuthByRecordAndProvider(record, s.provider)

		hasErr := err != nil
		expectErr := s.expectedId == ""
		if hasErr != expectErr {
			t.Errorf("(%d) Expected hasErr %v, got %v", i, expectErr, err)
			continue
		}

		if auth != nil && auth.Id != s.expectedId {
			t.Errorf("(%d) Expected external auth with ID %s, got \n%v", i, s.expectedId, auth)
		}
	}
}

func TestSaveExternalAuth(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// save with empty provider data
	emptyAuth := &models.ExternalAuth{}
	if err := app.Dao().SaveExternalAuth(emptyAuth); err == nil {
		t.Fatal("Expected error, got nil")
	}

	auth := &models.ExternalAuth{
		RecordId:     "o1y0dd0spd786md",
		CollectionId: "v851q4r790rhknl",
		Provider:     "test",
		ProviderId:   "test_id",
	}

	if err := app.Dao().SaveExternalAuth(auth); err != nil {
		t.Fatal(err)
	}

	// check if it was really saved
	foundAuth, err := app.Dao().FindFirstExternalAuthByExpr(dbx.HashExp{
		"collectionId": "v851q4r790rhknl",
		"provider":     "test",
		"providerId":   "test_id",
	})
	if err != nil {
		t.Fatal(err)
	}

	if auth.Id != foundAuth.Id {
		t.Fatalf("Expected ExternalAuth with id %s, got \n%v", auth.Id, foundAuth)
	}
}

func TestDeleteExternalAuth(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	record, err := app.Dao().FindRecordById("users", "4q1xlclmfloku33")
	if err != nil {
		t.Fatal(err)
	}

	auths, err := app.Dao().FindAllExternalAuthsByRecord(record)
	if err != nil {
		t.Fatal(err)
	}

	for _, auth := range auths {
		if err := app.Dao().DeleteExternalAuth(auth); err != nil {
			t.Fatalf("Failed to delete the ExternalAuth relation, got \n%v", err)
		}
	}

	// check if the relations were really deleted
	newAuths, err := app.Dao().FindAllExternalAuthsByRecord(record)
	if err != nil {
		t.Fatal(err)
	}

	if len(newAuths) != 0 {
		t.Fatalf("Expected all record %s ExternalAuth relations to be deleted, got \n%v", record.Id, newAuths)
	}
}
