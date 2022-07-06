package tokens_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tokens"
)

func TestNewUserAuthToken(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindUserByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewUserAuthToken(app, user)
	if err != nil {
		t.Fatal(err)
	}

	tokenUser, _ := app.Dao().FindUserByToken(
		token,
		app.Settings().UserAuthToken.Secret,
	)
	if tokenUser == nil || tokenUser.Id != user.Id {
		t.Fatalf("Expected user %v, got %v", user, tokenUser)
	}
}

func TestNewUserVerifyToken(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindUserByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewUserVerifyToken(app, user)
	if err != nil {
		t.Fatal(err)
	}

	tokenUser, _ := app.Dao().FindUserByToken(
		token,
		app.Settings().UserVerificationToken.Secret,
	)
	if tokenUser == nil || tokenUser.Id != user.Id {
		t.Fatalf("Expected user %v, got %v", user, tokenUser)
	}
}

func TestNewUserResetPasswordToken(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindUserByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewUserResetPasswordToken(app, user)
	if err != nil {
		t.Fatal(err)
	}

	tokenUser, _ := app.Dao().FindUserByToken(
		token,
		app.Settings().UserPasswordResetToken.Secret,
	)
	if tokenUser == nil || tokenUser.Id != user.Id {
		t.Fatalf("Expected user %v, got %v", user, tokenUser)
	}
}

func TestNewUserChangeEmailToken(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindUserByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewUserChangeEmailToken(app, user, "test_new@example.com")
	if err != nil {
		t.Fatal(err)
	}

	tokenUser, _ := app.Dao().FindUserByToken(
		token,
		app.Settings().UserEmailChangeToken.Secret,
	)
	if tokenUser == nil || tokenUser.Id != user.Id {
		t.Fatalf("Expected user %v, got %v", user, tokenUser)
	}
}
