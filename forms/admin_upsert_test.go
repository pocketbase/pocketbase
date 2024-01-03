package forms_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestNewAdminUpsert(t *testing.T) {
	t.Parallel()

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

func TestAdminUpsertValidateAndSubmit(t *testing.T) {
	t.Parallel()

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
			"sywbhecnh46rhm0",
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
			"sywbhecnh46rhm0",
			`{
				"email": "test2@example.com"
			}`,
			true,
		},
		{
			// update failure - mismatching passwords
			"sywbhecnh46rhm0",
			`{
				"password":        "1234567890",
				"passwordConfirm": "1234567891"
			}`,
			true,
		},
		{
			// update success - new email
			"sywbhecnh46rhm0",
			`{
				"email": "test_update@example.com"
			}`,
			false,
		},
		{
			// update success - new password
			"sywbhecnh46rhm0",
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

		interceptorCalls := 0

		err := form.Submit(func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
			return func(m *models.Admin) error {
				interceptorCalls++
				return next(m)
			}
		})

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v (%v)", i, s.expectError, hasErr, err)
		}

		foundAdmin, _ := app.Dao().FindAdminByEmail(form.Email)

		if !s.expectError && isCreate && foundAdmin == nil {
			t.Errorf("(%d) Expected admin to be created, got nil", i)
			continue
		}

		expectInterceptorCall := 1
		if s.expectError {
			expectInterceptorCall = 0
		}
		if interceptorCalls != expectInterceptorCall {
			t.Errorf("(%d) Expected interceptor to be called %d, got %d", i, expectInterceptorCall, interceptorCalls)
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

func TestAdminUpsertSubmitInterceptors(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	admin := &models.Admin{}
	form := forms.NewAdminUpsert(app, admin)
	form.Email = "test_new@example.com"
	form.Password = "1234567890"
	form.PasswordConfirm = form.Password

	testErr := errors.New("test_error")
	interceptorAdminEmail := ""

	interceptor1Called := false
	interceptor1 := func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
		return func(m *models.Admin) error {
			interceptor1Called = true
			return next(m)
		}
	}

	interceptor2Called := false
	interceptor2 := func(next forms.InterceptorNextFunc[*models.Admin]) forms.InterceptorNextFunc[*models.Admin] {
		return func(m *models.Admin) error {
			interceptorAdminEmail = admin.Email // to check if the record was filled
			interceptor2Called = true
			return testErr
		}
	}

	err := form.Submit(interceptor1, interceptor2)
	if err != testErr {
		t.Fatalf("Expected error %v, got %v", testErr, err)
	}

	if !interceptor1Called {
		t.Fatalf("Expected interceptor1 to be called")
	}

	if !interceptor2Called {
		t.Fatalf("Expected interceptor2 to be called")
	}

	if interceptorAdminEmail != form.Email {
		t.Fatalf("Expected the form model to be filled before calling the interceptors")
	}
}

func TestAdminUpsertWithCustomId(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	existingAdmin, err := app.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name        string
		jsonData    string
		collection  *models.Admin
		expectError bool
	}{
		{
			"empty data",
			"{}",
			&models.Admin{},
			false,
		},
		{
			"empty id",
			`{"id":""}`,
			&models.Admin{},
			false,
		},
		{
			"id < 15 chars",
			`{"id":"a23"}`,
			&models.Admin{},
			true,
		},
		{
			"id > 15 chars",
			`{"id":"a234567890123456"}`,
			&models.Admin{},
			true,
		},
		{
			"id = 15 chars (invalid chars)",
			`{"id":"a@3456789012345"}`,
			&models.Admin{},
			true,
		},
		{
			"id = 15 chars (valid chars)",
			`{"id":"a23456789012345"}`,
			&models.Admin{},
			false,
		},
		{
			"changing the id of an existing item",
			`{"id":"b23456789012345"}`,
			existingAdmin,
			true,
		},
		{
			"using the same existing item id",
			`{"id":"` + existingAdmin.Id + `"}`,
			existingAdmin,
			false,
		},
		{
			"skipping the id for existing item",
			`{}`,
			existingAdmin,
			false,
		},
	}

	for i, scenario := range scenarios {
		form := forms.NewAdminUpsert(app, scenario.collection)
		if form.Email == "" {
			form.Email = fmt.Sprintf("test_id_%d@example.com", i)
		}
		form.Password = "1234567890"
		form.PasswordConfirm = form.Password

		// load data
		loadErr := json.Unmarshal([]byte(scenario.jsonData), form)
		if loadErr != nil {
			t.Errorf("[%s] Failed to load form data: %v", scenario.name, loadErr)
			continue
		}

		submitErr := form.Submit()
		hasErr := submitErr != nil

		if hasErr != scenario.expectError {
			t.Errorf("[%s] Expected hasErr to be %v, got %v (%v)", scenario.name, scenario.expectError, hasErr, submitErr)
		}

		if !hasErr && form.Id != "" {
			_, err := app.Dao().FindAdminById(form.Id)
			if err != nil {
				t.Errorf("[%s] Expected to find record with id %s, got %v", scenario.name, form.Id, err)
			}
		}
	}
}
