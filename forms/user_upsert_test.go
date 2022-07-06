package forms_test

import (
	"encoding/json"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestNewUserUpsert(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user := &models.User{}
	user.Email = "new@example.com"

	form := forms.NewUserUpsert(app, user)

	// check defaults loading
	if form.Email != user.Email {
		t.Fatalf("Expected email %q, got %q", user.Email, form.Email)
	}
}

func TestUserUpsertValidate(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// mock app constraints
	app.Settings().EmailAuth.MinPasswordLength = 5
	app.Settings().EmailAuth.ExceptDomains = []string{"test.com"}
	app.Settings().EmailAuth.OnlyDomains = []string{"example.com", "test.com"}

	scenarios := []struct {
		id             string
		jsonData       string
		expectedErrors []string
	}{
		// empty data - create
		{
			"",
			`{}`,
			[]string{"email", "password", "passwordConfirm"},
		},
		// empty data - update
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{}`,
			[]string{},
		},
		// invalid email address
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"email":"invalid"}`,
			[]string{"email"},
		},
		// unique email constraint check (same email, aka. no changes)
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"email":"test@example.com"}`,
			[]string{},
		},
		// unique email constraint check (existing email)
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"email":"test2@something.com"}`,
			[]string{"email"},
		},
		// unique email constraint check (new email)
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"email":"new@example.com"}`,
			[]string{},
		},
		// EmailAuth.OnlyDomains constraints check
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"email":"test@something.com"}`,
			[]string{"email"},
		},
		// EmailAuth.ExceptDomains constraints check
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"email":"test@test.com"}`,
			[]string{"email"},
		},
		// password length constraint check
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"password":"1234", "passwordConfirm": "1234"}`,
			[]string{"password"},
		},
		// passwords mismatch
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"password":"12345", "passwordConfirm": "54321"}`,
			[]string{"passwordConfirm"},
		},
		// valid data - all fields
		{
			"",
			`{"email":"new@example.com","password":"12345","passwordConfirm":"12345"}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		user := &models.User{}
		if s.id != "" {
			user, _ = app.Dao().FindUserById(s.id)
		}

		form := forms.NewUserUpsert(app, user)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		// parse errors
		result := form.Validate()
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("(%d) Failed to parse errors %v", i, result)
			continue
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("(%d) Expected error keys %v, got %v", i, s.expectedErrors, errs)
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("(%d) Missing expected error key %q in %v", i, k, errs)
			}
		}
	}
}

func TestUserUpsertSubmit(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		id          string
		jsonData    string
		expectError bool
	}{
		// empty fields - create (Validate call check)
		{
			"",
			`{}`,
			true,
		},
		// empty fields - update (Validate call check)
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{}`,
			false,
		},
		// updating with existing user email
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"email":"test2@example.com"}`,
			true,
		},
		// updating with nonexisting user email
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"email":"update_new@example.com"}`,
			false,
		},
		// changing password
		{
			"4d0197cc-2b4a-3f83-a26b-d77bc8423d3c",
			`{"password":"123456789","passwordConfirm":"123456789"}`,
			false,
		},
		// creating user (existing email)
		{
			"",
			`{"email":"test3@example.com","password":"123456789","passwordConfirm":"123456789"}`,
			true,
		},
		// creating user (new email)
		{
			"",
			`{"email":"create_new@example.com","password":"123456789","passwordConfirm":"123456789"}`,
			false,
		},
	}

	for i, s := range scenarios {
		user := &models.User{}
		originalUser := &models.User{}
		if s.id != "" {
			user, _ = app.Dao().FindUserById(s.id)
			originalUser, _ = app.Dao().FindUserById(s.id)
		}

		form := forms.NewUserUpsert(app, user)

		// load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		err := form.Submit()

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		if s.expectError {
			continue
		}

		if user.Email != form.Email {
			t.Errorf("(%d) Expected email %q, got %q", i, form.Email, user.Email)
		}

		// on email change Verified should reset
		if user.Email != originalUser.Email && user.Verified {
			t.Errorf("(%d) Expected Verified to be false, got true", i)
		}

		if form.Password != "" && !user.ValidatePassword(form.Password) {
			t.Errorf("(%d) Expected password to be updated to %q", i, form.Password)
		}
		if form.Password != "" && originalUser.TokenKey == user.TokenKey {
			t.Errorf("(%d) Expected TokenKey to change, got %q", i, user.TokenKey)
		}
	}
}
