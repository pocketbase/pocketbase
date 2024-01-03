package tokens_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tokens"
)

func TestNewRecordAuthToken(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewRecordAuthToken(app, user)
	if err != nil {
		t.Fatal(err)
	}

	tokenRecord, _ := app.Dao().FindAuthRecordByToken(
		token,
		app.Settings().RecordAuthToken.Secret,
	)
	if tokenRecord == nil || tokenRecord.Id != user.Id {
		t.Fatalf("Expected auth record %v, got %v", user, tokenRecord)
	}
}

func TestNewRecordVerifyToken(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewRecordVerifyToken(app, user)
	if err != nil {
		t.Fatal(err)
	}

	tokenRecord, _ := app.Dao().FindAuthRecordByToken(
		token,
		app.Settings().RecordVerificationToken.Secret,
	)
	if tokenRecord == nil || tokenRecord.Id != user.Id {
		t.Fatalf("Expected auth record %v, got %v", user, tokenRecord)
	}
}

func TestNewRecordResetPasswordToken(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewRecordResetPasswordToken(app, user)
	if err != nil {
		t.Fatal(err)
	}

	tokenRecord, _ := app.Dao().FindAuthRecordByToken(
		token,
		app.Settings().RecordPasswordResetToken.Secret,
	)
	if tokenRecord == nil || tokenRecord.Id != user.Id {
		t.Fatalf("Expected auth record %v, got %v", user, tokenRecord)
	}
}

func TestNewRecordChangeEmailToken(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewRecordChangeEmailToken(app, user, "test_new@example.com")
	if err != nil {
		t.Fatal(err)
	}

	tokenRecord, _ := app.Dao().FindAuthRecordByToken(
		token,
		app.Settings().RecordEmailChangeToken.Secret,
	)
	if tokenRecord == nil || tokenRecord.Id != user.Id {
		t.Fatalf("Expected auth record %v, got %v", user, tokenRecord)
	}
}

func TestNewRecordFileToken(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.Dao().FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := tokens.NewRecordFileToken(app, user)
	if err != nil {
		t.Fatal(err)
	}

	tokenRecord, _ := app.Dao().FindAuthRecordByToken(
		token,
		app.Settings().RecordFileToken.Secret,
	)
	if tokenRecord == nil || tokenRecord.Id != user.Id {
		t.Fatalf("Expected auth record %v, got %v", user, tokenRecord)
	}
}
