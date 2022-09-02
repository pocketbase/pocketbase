package daos_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestExternalAuthQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	expected := "SELECT {{_externalAuths}}.* FROM `_externalAuths`"

	sql := app.Dao().ExternalAuthQuery().Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindAllExternalAuthsByUserId(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		userId        string
		expectedCount int
	}{
		{"", 0},
		{"missing", 0},
		{"97cc3d3d-6ba2-383f-b42a-7bc84d27410c", 0},
		{"cx9u0dh2udo8xol", 2},
	}

	for i, s := range scenarios {
		auths, err := app.Dao().FindAllExternalAuthsByUserId(s.userId)
		if err != nil {
			t.Errorf("(%d) Unexpected error %v", i, err)
			continue
		}

		if len(auths) != s.expectedCount {
			t.Errorf("(%d) Expected %d auths, got %d", i, s.expectedCount, len(auths))
		}

		for _, auth := range auths {
			if auth.UserId != s.userId {
				t.Errorf("(%d) Expected all auths to be linked to userId %s, got %v", i, s.userId, auth)
			}
		}
	}
}

func TestFindExternalAuthByProvider(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		provider   string
		providerId string
		expectedId string
	}{
		{"", "", ""},
		{"github", "", ""},
		{"github", "id1", ""},
		{"github", "id2", ""},
		{"google", "id1", "abcdefghijklmn0"},
		{"gitlab", "id2", "abcdefghijklmn1"},
	}

	for i, s := range scenarios {
		auth, err := app.Dao().FindExternalAuthByProvider(s.provider, s.providerId)

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

func TestFindExternalAuthByUserIdAndProvider(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		userId     string
		provider   string
		expectedId string
	}{
		{"", "", ""},
		{"", "github", ""},
		{"123456", "github", ""}, // missing user and provider record
		{"123456", "google", ""}, // missing user but existing provider record
		{"97cc3d3d-6ba2-383f-b42a-7bc84d27410c", "google", ""},
		{"cx9u0dh2udo8xol", "google", "abcdefghijklmn0"},
		{"cx9u0dh2udo8xol", "gitlab", "abcdefghijklmn1"},
	}

	for i, s := range scenarios {
		auth, err := app.Dao().FindExternalAuthByUserIdAndProvider(s.userId, s.provider)

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
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// save with empty provider data
	emptyAuth := &models.ExternalAuth{}
	if err := app.Dao().SaveExternalAuth(emptyAuth); err == nil {
		t.Fatal("Expected error, got nil")
	}

	auth := &models.ExternalAuth{
		UserId:     "97cc3d3d-6ba2-383f-b42a-7bc84d27410c",
		Provider:   "test",
		ProviderId: "test_id",
	}

	if err := app.Dao().SaveExternalAuth(auth); err != nil {
		t.Fatal(err)
	}

	// check if it was really saved
	foundAuth, err := app.Dao().FindExternalAuthByProvider("test", "test_id")
	if err != nil {
		t.Fatal(err)
	}

	if auth.Id != foundAuth.Id {
		t.Fatalf("Expected ExternalAuth with id %s, got \n%v", auth.Id, foundAuth)
	}
}

func TestDeleteExternalAuth(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindUserById("cx9u0dh2udo8xol")
	if err != nil {
		t.Fatal(err)
	}

	auths, err := app.Dao().FindAllExternalAuthsByUserId(user.Id)
	if err != nil {
		t.Fatal(err)
	}

	if err := app.Dao().DeleteExternalAuth(auths[0]); err != nil {
		t.Fatalf("Failed to delete the first ExternalAuth relation, got \n%v", err)
	}

	if err := app.Dao().DeleteExternalAuth(auths[1]); err == nil {
		t.Fatal("Expected delete to fail, got nil")
	}

	// update the user model and try again
	user.Email = "test_new@example.com"
	if err := app.Dao().SaveUser(user); err != nil {
		t.Fatal(err)
	}

	// try to delete auths[1] again
	if err := app.Dao().DeleteExternalAuth(auths[1]); err != nil {
		t.Fatalf("Failed to delete the last ExternalAuth relation, got \n%v", err)
	}

	// check if the relations were really deleted
	newAuths, err := app.Dao().FindAllExternalAuthsByUserId(user.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(newAuths) != 0 {
		t.Fatalf("Expected all user %s ExternalAuth relations to be deleted, got \n%v", user.Id, newAuths)
	}
}
