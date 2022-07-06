package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestNewAdminUpsert(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	admin := &models.Admin{}
	admin.Avatar = 3
	admin.Email = "new@example.com"

	form := forms.NewAdminUpsert(app, admin)

	// test defaults
	if form.Avatar != admin.Avatar {
		t.Errorf("Expected Avatar %d, got %d", admin.Avatar, form.Avatar)
	}
	if form.Email != admin.Email {
		t.Errorf("Expected Email %q, got %q", admin.Email, form.Email)
	}
}

func TestAdminUpsertValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		id              string
		avatar          int
		email           string
		password        string
		passwordConfirm string
		expectedErrors  int
	}{
		{
			"",
			-1,
			"",
			"",
			"",
			3,
		},
		{
			"",
			10,
			"invalid",
			"12345678",
			"87654321",
			4,
		},
		{
			// existing email
			"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			3,
			"test2@example.com",
			"1234567890",
			"1234567890",
			1,
		},
		{
			// mismatching passwords
			"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			3,
			"test@example.com",
			"1234567890",
			"1234567891",
			1,
		},
		{
			// create without setting password
			"",
			9,
			"test_create@example.com",
			"",
			"",
			1,
		},
		{
			// create with existing email
			"",
			9,
			"test@example.com",
			"1234567890!",
			"1234567890!",
			1,
		},
		{
			// update without setting password
			"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			3,
			"test_update@example.com",
			"",
			"",
			0,
		},
		{
			// create with password
			"",
			9,
			"test_create@example.com",
			"1234567890!",
			"1234567890!",
			0,
		},
		{
			// update with password
			"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			4,
			"test_update@example.com",
			"1234567890",
			"1234567890",
			0,
		},
	}

	for i, s := range scenarios {
		admin := &models.Admin{}
		if s.id != "" {
			admin, _ = app.Dao().FindAdminById(s.id)
		}

		form := forms.NewAdminUpsert(app, admin)
		form.Avatar = s.avatar
		form.Email = s.email
		form.Password = s.password
		form.PasswordConfirm = s.passwordConfirm

		result := form.Validate()
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("(%d) Failed to parse errors %v", i, result)
			continue
		}

		if len(errs) != s.expectedErrors {
			t.Errorf("(%d) Expected %d errors, got %d (%v)", i, s.expectedErrors, len(errs), errs)
		}
	}
}

func TestAdminUpsertSubmit(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		id          string
		jsonData    string
		expectError bool
	}{
		{
			// create empty
			"",
			`{}`,
			true,
		},
		{
			// update empty
			"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			`{}`,
			false,
		},
		{
			// create failure - existing email
			"",
			`{
				"email":           "test@example.com",
				"password":        "1234567890",
				"passwordConfirm": "1234567890"
			}`,
			true,
		},
		{
			// create failure - passwords mismatch
			"",
			`{
				"email":           "test_new@example.com",
				"password":        "1234567890",
				"passwordConfirm": "1234567891"
			}`,
			true,
		},
		{
			// create success
			"",
			`{
				"email":           "test_new@example.com",
				"password":        "1234567890",
				"passwordConfirm": "1234567890"
			}`,
			false,
		},
		{
			// update failure - existing email
			"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			`{
				"email": "test2@example.com"
			}`,
			true,
		},
		{
			// update failure - mismatching passwords
			"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			`{
				"password":        "1234567890",
				"passwordConfirm": "1234567891"
			}`,
			true,
		},
		{
			// update succcess - new email
			"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			`{
				"email": "test_update@example.com"
			}`,
			false,
		},
		{
			// update succcess - new password
			"2b4a97cc-3f83-4d01-a26b-3d77bc842d3c",
			`{
				"password":        "1234567890",
				"passwordConfirm": "1234567890"
			}`,
			false,
		},
	}

	for i, s := range scenarios {
		isCreate := true
		admin := &models.Admin{}
		if s.id != "" {
			isCreate = false
			admin, _ = app.Dao().FindAdminById(s.id)
		}
		initialTokenKey := admin.TokenKey

		form := forms.NewAdminUpsert(app, admin)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		foundAdmin, _ := app.Dao().FindAdminByEmail(form.Email)

		if !s.expectError && isCreate && foundAdmin == nil {
			t.Errorf("(%d) Expected admin to be created, got nil", i)
			continue
		}

		if s.expectError {
			continue // skip persistence check
		}

		if foundAdmin.Email != form.Email {
			t.Errorf("(%d) Expected email %s, got %s", i, form.Email, foundAdmin.Email)
		}

		if foundAdmin.Avatar != form.Avatar {
			t.Errorf("(%d) Expected avatar %d, got %d", i, form.Avatar, foundAdmin.Avatar)
		}

		if form.Password != "" && initialTokenKey == foundAdmin.TokenKey {
			t.Errorf("(%d) Expected token key to be renewed when setting a new password", i)
		}
	}
}
