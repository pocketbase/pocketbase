package daos_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestAdminQuery(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	expected := "SELECT {{_admins}}.* FROM `_admins`"

	sql := app.Dao().AdminQuery().Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindAdminById(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		id          string
		expectError bool
	}{
		{" ", true},
		{"missing", true},
		{"9q2trqumvlyr3bd", false},
	}

	for i, scenario := range scenarios {
		admin, err := app.Dao().FindAdminById(scenario.id)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if admin != nil && admin.Id != scenario.id {
			t.Errorf("(%d) Expected admin with id %s, got %s", i, scenario.id, admin.Id)
		}
	}
}

func TestFindAdminByEmail(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		email       string
		expectError bool
	}{
		{"", true},
		{"invalid", true},
		{"missing@example.com", true},
		{"test@example.com", false},
	}

	for i, scenario := range scenarios {
		admin, err := app.Dao().FindAdminByEmail(scenario.email)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if !scenario.expectError && admin.Email != scenario.email {
			t.Errorf("(%d) Expected admin with email %s, got %s", i, scenario.email, admin.Email)
		}
	}
}

func TestFindAdminByToken(t *testing.T) {
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
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MDk5MTY2MX0.qrbkI2TITtFKMP6vrATrBVKPGjEiDIBeQ0mlqPGMVeY",
			app.Settings().AdminAuthToken.Secret,
			"",
			true,
		},
		// expired token
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MTY0MDk5MTY2MX0.I7w8iktkleQvC7_UIRpD7rNzcU4OnF7i7SFIUu6lD_4",
			app.Settings().AdminAuthToken.Secret,
			"",
			true,
		},
		// wrong base token (password reset token secret instead of auth secret)
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			app.Settings().AdminPasswordResetToken.Secret,
			"",
			true,
		},
		// valid token
		{
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhZG1pbiIsImV4cCI6MjIwODk4NTI2MX0.M1m--VOqGyv0d23eeUc0r9xE8ZzHaYVmVFw1VZW6gT8",
			app.Settings().AdminAuthToken.Secret,
			"test@example.com",
			false,
		},
	}

	for i, scenario := range scenarios {
		admin, err := app.Dao().FindAdminByToken(scenario.token, scenario.baseKey)

		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
			continue
		}

		if !scenario.expectError && admin.Email != scenario.expectedEmail {
			t.Errorf("(%d) Expected admin model %s, got %s", i, scenario.expectedEmail, admin.Email)
		}
	}
}

func TestTotalAdmins(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	result1, err := app.Dao().TotalAdmins()
	if err != nil {
		t.Fatal(err)
	}
	if result1 != 3 {
		t.Fatalf("Expected 3 admins, got %d", result1)
	}

	// delete all
	app.Dao().DB().NewQuery("delete from {{_admins}}").Execute()

	result2, err := app.Dao().TotalAdmins()
	if err != nil {
		t.Fatal(err)
	}
	if result2 != 0 {
		t.Fatalf("Expected 0 admins, got %d", result2)
	}
}

func TestIsAdminEmailUnique(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		email     string
		excludeId string
		expected  bool
	}{
		{"", "", false},
		{"test@example.com", "", false},
		{"test2@example.com", "", false},
		{"test3@example.com", "", false},
		{"new@example.com", "", true},
		{"test@example.com", "sywbhecnh46rhm0", true},
	}

	for i, scenario := range scenarios {
		result := app.Dao().IsAdminEmailUnique(scenario.email, scenario.excludeId)
		if result != scenario.expected {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expected, result)
		}
	}
}

func TestDeleteAdmin(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// try to delete unsaved admin model
	deleteErr0 := app.Dao().DeleteAdmin(&models.Admin{})
	if deleteErr0 == nil {
		t.Fatal("Expected error, got nil")
	}

	admin1, err := app.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	admin2, err := app.Dao().FindAdminByEmail("test2@example.com")
	if err != nil {
		t.Fatal(err)
	}
	admin3, err := app.Dao().FindAdminByEmail("test3@example.com")
	if err != nil {
		t.Fatal(err)
	}

	deleteErr1 := app.Dao().DeleteAdmin(admin1)
	if deleteErr1 != nil {
		t.Fatal(deleteErr1)
	}

	deleteErr2 := app.Dao().DeleteAdmin(admin2)
	if deleteErr2 != nil {
		t.Fatal(deleteErr2)
	}

	// cannot delete the only remaining admin
	deleteErr3 := app.Dao().DeleteAdmin(admin3)
	if deleteErr3 == nil {
		t.Fatal("Expected delete error, got nil")
	}

	total, _ := app.Dao().TotalAdmins()
	if total != 1 {
		t.Fatalf("Expected only 1 admin, got %d", total)
	}
}

func TestSaveAdmin(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create
	newAdmin := &models.Admin{}
	newAdmin.Email = "new@example.com"
	newAdmin.SetPassword("123456")
	saveErr1 := app.Dao().SaveAdmin(newAdmin)
	if saveErr1 != nil {
		t.Fatal(saveErr1)
	}
	if newAdmin.Id == "" {
		t.Fatal("Expected admin id to be set")
	}

	// update
	existingAdmin, err := app.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	updatedEmail := "test_update@example.com"
	existingAdmin.Email = updatedEmail
	saveErr2 := app.Dao().SaveAdmin(existingAdmin)
	if saveErr2 != nil {
		t.Fatal(saveErr2)
	}
	existingAdmin, _ = app.Dao().FindAdminById(existingAdmin.Id)
	if existingAdmin.Email != updatedEmail {
		t.Fatalf("Expected admin email to be %s, got %s", updatedEmail, existingAdmin.Email)
	}
}
