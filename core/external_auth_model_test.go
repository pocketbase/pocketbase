package core_test

import (
	"fmt"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestNewExternalAuth(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	ea := core.NewExternalAuth(app)

	if ea.Collection().Name != core.CollectionNameExternalAuths {
		t.Fatalf("Expected record with %q collection, got %q", core.CollectionNameExternalAuths, ea.Collection().Name)
	}
}

func TestExternalAuthProxyRecord(t *testing.T) {
	t.Parallel()

	record := core.NewRecord(core.NewBaseCollection("test"))
	record.Id = "test_id"

	ea := core.ExternalAuth{}
	ea.SetProxyRecord(record)

	if ea.ProxyRecord() == nil || ea.ProxyRecord().Id != record.Id {
		t.Fatalf("Expected proxy record with id %q, got %v", record.Id, ea.ProxyRecord())
	}
}

func TestExternalAuthRecordRef(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	ea := core.NewExternalAuth(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			ea.SetRecordRef(testValue)

			if v := ea.RecordRef(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := ea.GetString("recordRef"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestExternalAuthCollectionRef(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	ea := core.NewExternalAuth(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			ea.SetCollectionRef(testValue)

			if v := ea.CollectionRef(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := ea.GetString("collectionRef"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestExternalAuthProvider(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	ea := core.NewExternalAuth(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			ea.SetProvider(testValue)

			if v := ea.Provider(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := ea.GetString("provider"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestExternalAuthProviderId(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	ea := core.NewExternalAuth(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			ea.SetProviderId(testValue)

			if v := ea.ProviderId(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := ea.GetString("providerId"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestExternalAuthCreated(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	ea := core.NewExternalAuth(app)

	if v := ea.Created().String(); v != "" {
		t.Fatalf("Expected empty created, got %q", v)
	}

	now := types.NowDateTime()
	ea.SetRaw("created", now)

	if v := ea.Created().String(); v != now.String() {
		t.Fatalf("Expected %q created, got %q", now.String(), v)
	}
}

func TestExternalAuthUpdated(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	ea := core.NewExternalAuth(app)

	if v := ea.Updated().String(); v != "" {
		t.Fatalf("Expected empty updated, got %q", v)
	}

	now := types.NowDateTime()
	ea.SetRaw("updated", now)

	if v := ea.Updated().String(); v != now.String() {
		t.Fatalf("Expected %q updated, got %q", now.String(), v)
	}
}

func TestExternalAuthPreValidate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	externalAuthsCol, err := app.FindCollectionByNameOrId(core.CollectionNameExternalAuths)
	if err != nil {
		t.Fatal(err)
	}

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("no proxy record", func(t *testing.T) {
		externalAuth := &core.ExternalAuth{}

		if err := app.Validate(externalAuth); err == nil {
			t.Fatal("Expected collection validation error")
		}
	})

	t.Run("non-ExternalAuth collection", func(t *testing.T) {
		externalAuth := &core.ExternalAuth{}
		externalAuth.SetProxyRecord(core.NewRecord(core.NewBaseCollection("invalid")))
		externalAuth.SetRecordRef(user.Id)
		externalAuth.SetCollectionRef(user.Collection().Id)
		externalAuth.SetProvider("gitlab")
		externalAuth.SetProviderId("test123")

		if err := app.Validate(externalAuth); err == nil {
			t.Fatal("Expected collection validation error")
		}
	})

	t.Run("ExternalAuth collection", func(t *testing.T) {
		externalAuth := &core.ExternalAuth{}
		externalAuth.SetProxyRecord(core.NewRecord(externalAuthsCol))
		externalAuth.SetRecordRef(user.Id)
		externalAuth.SetCollectionRef(user.Collection().Id)
		externalAuth.SetProvider("gitlab")
		externalAuth.SetProviderId("test123")

		if err := app.Validate(externalAuth); err != nil {
			t.Fatalf("Expected nil validation error, got %v", err)
		}
	})
}

func TestExternalAuthValidateHook(t *testing.T) {
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
		externalAuth func() *core.ExternalAuth
		expectErrors []string
	}{
		{
			"empty",
			func() *core.ExternalAuth {
				return core.NewExternalAuth(app)
			},
			[]string{"collectionRef", "recordRef", "provider", "providerId"},
		},
		{
			"non-auth collection",
			func() *core.ExternalAuth {
				ea := core.NewExternalAuth(app)
				ea.SetCollectionRef(demo1.Collection().Id)
				ea.SetRecordRef(demo1.Id)
				ea.SetProvider("gitlab")
				ea.SetProviderId("test123")
				return ea
			},
			[]string{"collectionRef"},
		},
		{
			"disabled provider",
			func() *core.ExternalAuth {
				ea := core.NewExternalAuth(app)
				ea.SetCollectionRef(user.Collection().Id)
				ea.SetRecordRef("missing")
				ea.SetProvider("apple")
				ea.SetProviderId("test123")
				return ea
			},
			[]string{"recordRef"},
		},
		{
			"missing record id",
			func() *core.ExternalAuth {
				ea := core.NewExternalAuth(app)
				ea.SetCollectionRef(user.Collection().Id)
				ea.SetRecordRef("missing")
				ea.SetProvider("gitlab")
				ea.SetProviderId("test123")
				return ea
			},
			[]string{"recordRef"},
		},
		{
			"valid ref",
			func() *core.ExternalAuth {
				ea := core.NewExternalAuth(app)
				ea.SetCollectionRef(user.Collection().Id)
				ea.SetRecordRef(user.Id)
				ea.SetProvider("gitlab")
				ea.SetProviderId("test123")
				return ea
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			errs := app.Validate(s.externalAuth())
			tests.TestValidationErrors(t, errs, s.expectErrors)
		})
	}
}
