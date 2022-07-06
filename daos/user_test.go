package daos_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestUserQuery(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	expected := "SELECT {{_users}}.* FROM `_users`"

	sql := app.Dao().UserQuery().Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestLoadProfile(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// try to load missing profile (shouldn't return an error)
	// ---
	newUser := &models.User{}
	err1 := app.Dao().LoadProfile(newUser)
	if err1 != nil {
		t.Fatalf("Expected nil, got error %v", err1)
	}

	// try to load existing profile
	// ---
	existingUser, _ := app.Dao().FindUserByEmail("test@example.com")
	existingUser.Profile = nil // reset

	err2 := app.Dao().LoadProfile(existingUser)
	if err2 != nil {
		t.Fatal(err2)
	}

	if existingUser.Profile == nil {
		t.Fatal("Expected user profile to be loaded, got nil")
	}

	if existingUser.Profile.GetStringDataValue("name") != "test" {
		t.Fatalf("Expected profile.name to be 'test', got %s", existingUser.Profile.GetStringDataValue("name"))
	}
}

func TestLoadProfiles(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	u0 := &models.User{}
	u1, _ := app.Dao().FindUserByEmail("test@example.com")
	u2, _ := app.Dao().FindUserByEmail("test2@example.com")

	users := []*models.User{u0, u1, u2}

	err := app.Dao().LoadProfiles(users)
	if err != nil {
		t.Fatal(err)
	}

	if u0.Profile != nil {
		t.Errorf("Expected profile to be nil for u0, got %v", u0.Profile)
	}
	if u1.Profile == nil {
		t.Errorf("Expected profile to be set for u1, got nil")
	}
	if u2.Profile == nil {
		t.Errorf("Expected profile to be set for u2, got nil")
	}
}

func TestFindUserById(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		id          string
		expectError bool
	}{
		{"00000000-2b4a-a26b-4d01-42d3c3d77bc8", true},
		{"97cc3d3d-6ba2-383f-b42a-7bc84d27410c", false},
	}

	for i, scenario := range scenarios {
		user, err := app.Dao().FindUserById(scenario.id)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if user != nil && user.Id != scenario.id {
			t.Errorf("(%d) Expected user with id %s, got %s", i, scenario.id, user.Id)
		}
	}
}

func TestFindUserByEmail(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		email       string
		expectError bool
	}{
		{"invalid", true},
		{"missing@example.com", true},
		{"test@example.com", false},
	}

	for i, scenario := range scenarios {
		user, err := app.Dao().FindUserByEmail(scenario.email)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if !scenario.expectError && user.Email != scenario.email {
			t.Errorf("(%d) Expected user with email %s, got %s", i, scenario.email, user.Email)
		}
	}
}

func TestFindUserByToken(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		token         string
		baseKey       string
		expectedEmail string
		expectError   bool
	}{
		// invalid base key (password reset key for auth token)
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			app.Settings().UserPasswordResetToken.Secret,
			"",
			true,
		},
		// expired token
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxNjQwOTkxNjYxfQ.RrSG5NwysI38DEZrIQiz3lUgI6sEuYGTll_jLRbBSiw",
			app.Settings().UserAuthToken.Secret,
			"",
			true,
		},
		// valid token
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjRkMDE5N2NjLTJiNGEtM2Y4My1hMjZiLWQ3N2JjODQyM2QzYyIsInR5cGUiOiJ1c2VyIiwiZXhwIjoxODkzNDc0MDAwfQ.Wq5ac1q1f5WntIzEngXk22ydMj-eFgvfSRg7dhmPKic",
			app.Settings().UserAuthToken.Secret,
			"test@example.com",
			false,
		},
	}

	for i, scenario := range scenarios {
		user, err := app.Dao().FindUserByToken(scenario.token, scenario.baseKey)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if !scenario.expectError && user.Email != scenario.expectedEmail {
			t.Errorf("(%d) Expected user model %s, got %s", i, scenario.expectedEmail, user.Email)
		}
	}
}

func TestIsUserEmailUnique(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		email     string
		excludeId string
		expected  bool
	}{
		{"", "", false},
		{"test@example.com", "", false},
		{"new@example.com", "", true},
		{"test@example.com", "4d0197cc-2b4a-3f83-a26b-d77bc8423d3c", true},
	}

	for i, scenario := range scenarios {
		result := app.Dao().IsUserEmailUnique(scenario.email, scenario.excludeId)
		if result != scenario.expected {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expected, result)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// try to delete unsaved user
	// ---
	err1 := app.Dao().DeleteUser(&models.User{})
	if err1 == nil {
		t.Fatal("Expected error, got nil")
	}

	// try to delete existing user
	// ---
	user, _ := app.Dao().FindUserByEmail("test3@example.com")
	err2 := app.Dao().DeleteUser(user)
	if err2 != nil {
		t.Fatalf("Expected nil, got error %v", err2)
	}

	// check if the delete operation was cascaded to the profiles collection (record delete)
	profilesCol, _ := app.Dao().FindCollectionByNameOrId(models.ProfileCollectionName)
	profile, _ := app.Dao().FindRecordById(profilesCol, user.Profile.Id, nil)
	if profile != nil {
		t.Fatalf("Expected user profile to be deleted, got %v", profile)
	}

	// check if delete operation was cascaded to the related demo2 collection (null set)
	demo2Col, _ := app.Dao().FindCollectionByNameOrId("demo2")
	record, _ := app.Dao().FindRecordById(demo2Col, "94568ca2-0bee-49d7-b749-06cb97956fd9", nil)
	if record == nil {
		t.Fatal("Expected to found related record, got nil")
	}
	if record.GetStringDataValue("user") != "" {
		t.Fatalf("Expected user field to be set to empty string, got %v", record.GetStringDataValue("user"))
	}
}

func TestSaveUser(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create
	// ---
	u1 := &models.User{}
	u1.Email = "new@example.com"
	u1.SetPassword("123456")
	err1 := app.Dao().SaveUser(u1)
	if err1 != nil {
		t.Fatal(err1)
	}
	u1, refreshErr1 := app.Dao().FindUserByEmail("new@example.com")
	if refreshErr1 != nil {
		t.Fatalf("Expected user with email new@example.com to have been created, got error %v", refreshErr1)
	}
	if u1.Profile == nil {
		t.Fatalf("Expected creating a user to create also an empty profile record")
	}

	// update
	// ---
	u2, _ := app.Dao().FindUserByEmail("test@example.com")
	u2.Email = "test_update@example.com"
	err2 := app.Dao().SaveUser(u2)
	if err2 != nil {
		t.Fatal(err2)
	}
	u2, refreshErr2 := app.Dao().FindUserByEmail("test_update@example.com")
	if u2 == nil {
		t.Fatalf("Couldn't find user with email test_update@example.com (%v)", refreshErr2)
	}
}
