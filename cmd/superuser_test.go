package cmd_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestSuperuserUpsertCommand(t *testing.T) {
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
			"existing user",
			"test@example.com",
			"1234567890!",
			false,
		},
		{
			"new user",
			"test_new@example.com",
			"1234567890!",
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			command := cmd.NewSuperuserCommand(app)
			command.SetArgs([]string{"upsert", s.email, s.password})

			err := command.Execute()

			hasErr := err != nil
			if s.expectError != hasErr {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			// check whether the superuser account was actually upserted
			superuser, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, s.email)
			if err != nil {
				t.Fatalf("Failed to fetch superuser %s: %v", s.email, err)
			} else if !superuser.ValidatePassword(s.password) {
				t.Fatal("Expected the superuser password to match")
			}
		})
	}
}

func TestSuperuserCreateCommand(t *testing.T) {
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
		t.Run(s.name, func(t *testing.T) {
			command := cmd.NewSuperuserCommand(app)
			command.SetArgs([]string{"create", s.email, s.password})

			err := command.Execute()

			hasErr := err != nil
			if s.expectError != hasErr {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			// check whether the superuser account was actually created
			superuser, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, s.email)
			if err != nil {
				t.Fatalf("Failed to fetch created superuser %s: %v", s.email, err)
			} else if !superuser.ValidatePassword(s.password) {
				t.Fatal("Expected the superuser password to match")
			}
		})
	}
}

func TestSuperuserUpdateCommand(t *testing.T) {
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
			"nonexisting superuser",
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
		t.Run(s.name, func(t *testing.T) {
			command := cmd.NewSuperuserCommand(app)
			command.SetArgs([]string{"update", s.email, s.password})

			err := command.Execute()

			hasErr := err != nil
			if s.expectError != hasErr {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			// check whether the superuser password was actually changed
			superuser, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, s.email)
			if err != nil {
				t.Fatalf("Failed to fetch superuser %s: %v", s.email, err)
			} else if !superuser.ValidatePassword(s.password) {
				t.Fatal("Expected the superuser password to match")
			}
		})
	}
}

func TestSuperuserDeleteCommand(t *testing.T) {
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
			"nonexisting superuser",
			"test_missing@example.com",
			false,
		},
		{
			"existing superuser",
			"test@example.com",
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			command := cmd.NewSuperuserCommand(app)
			command.SetArgs([]string{"delete", s.email})

			err := command.Execute()

			hasErr := err != nil
			if s.expectError != hasErr {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if _, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, s.email); err == nil {
				t.Fatal("Expected the superuser account to be deleted")
			}
		})
	}
}
