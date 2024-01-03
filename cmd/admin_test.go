package cmd_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/tests"
)

func TestAdminCreateCommand(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name        string
		email       string
		password    string
		expectError bool
	}{
		{
			"empty email and password",
			"",
			"",
			true,
		},
		{
			"empty email",
			"",
			"1234567890",
			true,
		},
		{
			"invalid email",
			"invalid",
			"1234567890",
			true,
		},
		{
			"duplicated email",
			"test@example.com",
			"1234567890",
			true,
		},
		{
			"empty password",
			"test@example.com",
			"",
			true,
		},
		{
			"short password",
			"test_new@example.com",
			"1234567",
			true,
		},
		{
			"valid email and password",
			"test_new@example.com",
			"12345678",
			false,
		},
	}

	for _, s := range scenarios {
		command := cmd.NewAdminCommand(app)
		command.SetArgs([]string{"create", s.email, s.password})

		err := command.Execute()

		hasErr := err != nil
		if s.expectError != hasErr {
			t.Errorf("[%s] Expected hasErr %v, got %v (%v)", s.name, s.expectError, hasErr, err)
		}

		if hasErr {
			continue
		}

		// check whether the admin account was actually created
		admin, err := app.Dao().FindAdminByEmail(s.email)
		if err != nil {
			t.Errorf("[%s] Failed to fetch created admin %s: %v", s.name, s.email, err)
		} else if !admin.ValidatePassword(s.password) {
			t.Errorf("[%s] Expected the admin password to match", s.name)
		}
	}
}

func TestAdminUpdateCommand(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name        string
		email       string
		password    string
		expectError bool
	}{
		{
			"empty email and password",
			"",
			"",
			true,
		},
		{
			"empty email",
			"",
			"1234567890",
			true,
		},
		{
			"invalid email",
			"invalid",
			"1234567890",
			true,
		},
		{
			"nonexisting admin",
			"test_missing@example.com",
			"1234567890",
			true,
		},
		{
			"empty password",
			"test@example.com",
			"",
			true,
		},
		{
			"short password",
			"test_new@example.com",
			"1234567",
			true,
		},
		{
			"valid email and password",
			"test@example.com",
			"12345678",
			false,
		},
	}

	for _, s := range scenarios {
		command := cmd.NewAdminCommand(app)
		command.SetArgs([]string{"update", s.email, s.password})

		err := command.Execute()

		hasErr := err != nil
		if s.expectError != hasErr {
			t.Errorf("[%s] Expected hasErr %v, got %v (%v)", s.name, s.expectError, hasErr, err)
		}

		if hasErr {
			continue
		}

		// check whether the admin password was actually changed
		admin, err := app.Dao().FindAdminByEmail(s.email)
		if err != nil {
			t.Errorf("[%s] Failed to fetch admin %s: %v", s.name, s.email, err)
		} else if !admin.ValidatePassword(s.password) {
			t.Errorf("[%s] Expected the admin password to match", s.name)
		}
	}
}

func TestAdminDeleteCommand(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			"empty email",
			"",
			true,
		},
		{
			"invalid email",
			"invalid",
			true,
		},
		{
			"nonexisting admin",
			"test_missing@example.com",
			false,
		},
		{
			"existing admin",
			"test@example.com",
			false,
		},
	}

	for _, s := range scenarios {
		command := cmd.NewAdminCommand(app)
		command.SetArgs([]string{"delete", s.email})

		err := command.Execute()

		hasErr := err != nil
		if s.expectError != hasErr {
			t.Errorf("[%s] Expected hasErr %v, got %v (%v)", s.name, s.expectError, hasErr, err)
		}

		if hasErr {
			continue
		}

		// check whether the admin account was actually deleted
		if _, err := app.Dao().FindAdminByEmail(s.email); err == nil {
			t.Errorf("[%s] Expected the admin account to be deleted", s.name)
		}
	}
}
