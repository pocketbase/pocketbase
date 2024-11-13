package core_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestNewOTP(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	otp := core.NewOTP(app)

	if otp.Collection().Name != core.CollectionNameOTPs {
		t.Fatalf("Expected record with %q collection, got %q", core.CollectionNameOTPs, otp.Collection().Name)
	}
}

func TestOTPProxyRecord(t *testing.T) {
	t.Parallel()

	record := core.NewRecord(core.NewBaseCollection("test"))
	record.Id = "test_id"

	otp := core.OTP{}
	otp.SetProxyRecord(record)

	if otp.ProxyRecord() == nil || otp.ProxyRecord().Id != record.Id {
		t.Fatalf("Expected proxy record with id %q, got %v", record.Id, otp.ProxyRecord())
	}
}

func TestOTPRecordRef(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	otp := core.NewOTP(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			otp.SetRecordRef(testValue)

			if v := otp.RecordRef(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := otp.GetString("recordRef"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestOTPCollectionRef(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	otp := core.NewOTP(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			otp.SetCollectionRef(testValue)

			if v := otp.CollectionRef(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := otp.GetString("collectionRef"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestOTPSentTo(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	otp := core.NewOTP(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			otp.SetSentTo(testValue)

			if v := otp.SentTo(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := otp.GetString("sentTo"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestOTPCreated(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	otp := core.NewOTP(app)

	if v := otp.Created().String(); v != "" {
		t.Fatalf("Expected empty created, got %q", v)
	}

	now := types.NowDateTime()
	otp.SetRaw("created", now)

	if v := otp.Created().String(); v != now.String() {
		t.Fatalf("Expected %q created, got %q", now.String(), v)
	}
}

func TestOTPUpdated(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	otp := core.NewOTP(app)

	if v := otp.Updated().String(); v != "" {
		t.Fatalf("Expected empty updated, got %q", v)
	}

	now := types.NowDateTime()
	otp.SetRaw("updated", now)

	if v := otp.Updated().String(); v != now.String() {
		t.Fatalf("Expected %q updated, got %q", now.String(), v)
	}
}

func TestOTPHasExpired(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	now := types.NowDateTime()

	otp := core.NewOTP(app)
	otp.SetRaw("created", now.Add(-5*time.Minute))

	scenarios := []struct {
		maxElapsed time.Duration
		expected   bool
	}{
		{0 * time.Minute, true},
		{3 * time.Minute, true},
		{5 * time.Minute, true},
		{6 * time.Minute, false},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.maxElapsed.String()), func(t *testing.T) {
			result := otp.HasExpired(s.maxElapsed)

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestOTPPreValidate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	otpsCol, err := app.FindCollectionByNameOrId(core.CollectionNameOTPs)
	if err != nil {
		t.Fatal(err)
	}

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("no proxy record", func(t *testing.T) {
		otp := &core.OTP{}

		if err := app.Validate(otp); err == nil {
			t.Fatal("Expected collection validation error")
		}
	})

	t.Run("non-OTP collection", func(t *testing.T) {
		otp := &core.OTP{}
		otp.SetProxyRecord(core.NewRecord(core.NewBaseCollection("invalid")))
		otp.SetRecordRef(user.Id)
		otp.SetCollectionRef(user.Collection().Id)
		otp.SetPassword("test123")

		if err := app.Validate(otp); err == nil {
			t.Fatal("Expected collection validation error")
		}
	})

	t.Run("OTP collection", func(t *testing.T) {
		otp := &core.OTP{}
		otp.SetProxyRecord(core.NewRecord(otpsCol))
		otp.SetRecordRef(user.Id)
		otp.SetCollectionRef(user.Collection().Id)
		otp.SetPassword("test123")

		if err := app.Validate(otp); err != nil {
			t.Fatalf("Expected nil validation error, got %v", err)
		}
	})
}

func TestOTPValidateHook(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	demo1, err := app.FindRecordById("demo1", "84nmscqy84lsi1t")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		name         string
		otp          func() *core.OTP
		expectErrors []string
	}{
		{
			"empty",
			func() *core.OTP {
				return core.NewOTP(app)
			},
			[]string{"collectionRef", "recordRef", "password"},
		},
		{
			"non-auth collection",
			func() *core.OTP {
				otp := core.NewOTP(app)
				otp.SetCollectionRef(demo1.Collection().Id)
				otp.SetRecordRef(demo1.Id)
				otp.SetPassword("test123")
				return otp
			},
			[]string{"collectionRef"},
		},
		{
			"missing record id",
			func() *core.OTP {
				otp := core.NewOTP(app)
				otp.SetCollectionRef(user.Collection().Id)
				otp.SetRecordRef("missing")
				otp.SetPassword("test123")
				return otp
			},
			[]string{"recordRef"},
		},
		{
			"valid ref",
			func() *core.OTP {
				otp := core.NewOTP(app)
				otp.SetCollectionRef(user.Collection().Id)
				otp.SetRecordRef(user.Id)
				otp.SetPassword("test123")
				return otp
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			errs := app.Validate(s.otp())
			tests.TestValidationErrors(t, errs, s.expectErrors)
		})
	}
}
