package core_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

func TestNewStaticAuthToken(t *testing.T) {
	t.Parallel()

	testRecordToken(t, core.TokenTypeAuth, func(record *core.Record) (string, error) {
		return record.NewStaticAuthToken(0)
	}, map[string]any{
		core.TokenClaimRefreshable: false,
	})
}

func TestNewStaticAuthTokenWithCustomDuration(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	var tolerance int64 = 1 // in sec

	durations := []int64{-100, 0, 100}

	for i, d := range durations {
		t.Run(fmt.Sprintf("%d_%d", i, d), func(t *testing.T) {
			now := time.Now()

			duration := time.Duration(d) * time.Second

			token, err := user.NewStaticAuthToken(duration)
			if err != nil {
				t.Fatal(err)
			}

			claims, err := security.ParseUnverifiedJWT(token)
			if err != nil {
				t.Fatal(err)
			}

			exp := cast.ToInt64(claims["exp"])

			expectedDuration := duration
			// should fallback to the collection setting
			if expectedDuration <= 0 {
				expectedDuration = user.Collection().AuthToken.DurationTime()
			}
			expectedMinExp := now.Add(expectedDuration).Unix() - tolerance
			expectedMaxExp := now.Add(expectedDuration).Unix() + tolerance

			if exp < expectedMinExp {
				t.Fatalf("Expected token exp to be greater than %d, got %d", expectedMinExp, exp)
			}

			if exp > expectedMaxExp {
				t.Fatalf("Expected token exp to be less than %d, got %d", expectedMaxExp, exp)
			}
		})
	}
}

func TestNewAuthToken(t *testing.T) {
	t.Parallel()

	testRecordToken(t, core.TokenTypeAuth, func(record *core.Record) (string, error) {
		return record.NewAuthToken()
	}, map[string]any{
		core.TokenClaimRefreshable: true,
	})
}

func TestNewVerificationToken(t *testing.T) {
	t.Parallel()

	testRecordToken(t, core.TokenTypeVerification, func(record *core.Record) (string, error) {
		return record.NewVerificationToken()
	}, nil)
}

func TestNewPasswordResetToken(t *testing.T) {
	t.Parallel()

	testRecordToken(t, core.TokenTypePasswordReset, func(record *core.Record) (string, error) {
		return record.NewPasswordResetToken()
	}, nil)
}

func TestNewEmailChangeToken(t *testing.T) {
	t.Parallel()

	testRecordToken(t, core.TokenTypeEmailChange, func(record *core.Record) (string, error) {
		return record.NewEmailChangeToken("new@example.com")
	}, nil)
}

func TestNewFileToken(t *testing.T) {
	t.Parallel()

	testRecordToken(t, core.TokenTypeFile, func(record *core.Record) (string, error) {
		return record.NewFileToken()
	}, nil)
}

func testRecordToken(
	t *testing.T,
	tokenType string,
	tokenFunc func(record *core.Record) (string, error),
	expectedClaims map[string]any,
) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	demo1, err := app.FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("non-auth record", func(t *testing.T) {
		_, err = tokenFunc(demo1)
		if err == nil {
			t.Fatal("Expected error for non-auth records")
		}
	})

	t.Run("auth record", func(t *testing.T) {
		token, err := tokenFunc(user)
		if err != nil {
			t.Fatal(err)
		}

		tokenRecord, _ := app.FindAuthRecordByToken(token, tokenType)
		if tokenRecord == nil || tokenRecord.Id != user.Id {
			t.Fatalf("Expected auth record\n%v\ngot\n%v", user, tokenRecord)
		}

		if len(expectedClaims) > 0 {
			claims, _ := security.ParseUnverifiedJWT(token)
			for k, v := range expectedClaims {
				if claims[k] != v {
					t.Errorf("Expected claim %q with value %#v, got %#v", k, v, claims[k])
				}
			}
		}
	})

	t.Run("empty signing key", func(t *testing.T) {
		user.SetTokenKey("")
		collection := user.Collection()
		*collection = core.Collection{}
		collection.Type = core.CollectionTypeAuth

		_, err := tokenFunc(user)
		if err == nil {
			t.Fatal("Expected empty signing key error")
		}
	})
}
