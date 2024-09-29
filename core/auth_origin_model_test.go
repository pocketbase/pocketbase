package core_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestNewAuthOrigin(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	origin := core.NewAuthOrigin(app)

	if origin.Collection().Name != core.CollectionNameAuthOrigins {
		t.Fatalf("Expected record with %q collection, got %q", core.CollectionNameAuthOrigins, origin.Collection().Name)
	}
}

func TestAuthOriginProxyRecord(t *testing.T) {
	t.Parallel()

	record := core.NewRecord(core.NewBaseCollection("test"))
	record.Id = "test_id"

	origin := core.AuthOrigin{}
	origin.SetProxyRecord(record)

	if origin.ProxyRecord() == nil || origin.ProxyRecord().Id != record.Id {
		t.Fatalf("Expected proxy record with id %q, got %v", record.Id, origin.ProxyRecord())
	}
}

func TestAuthOriginRecordRef(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	origin := core.NewAuthOrigin(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			origin.SetRecordRef(testValue)

			if v := origin.RecordRef(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := origin.GetString("recordRef"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestAuthOriginCollectionRef(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	origin := core.NewAuthOrigin(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			origin.SetCollectionRef(testValue)

			if v := origin.CollectionRef(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := origin.GetString("collectionRef"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestAuthOriginFingerprint(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	origin := core.NewAuthOrigin(app)

	testValues := []string{"test_1", "test2", ""}
	for i, testValue := range testValues {
		t.Run(fmt.Sprintf("%d_%q", i, testValue), func(t *testing.T) {
			origin.SetFingerprint(testValue)

			if v := origin.Fingerprint(); v != testValue {
				t.Fatalf("Expected getter %q, got %q", testValue, v)
			}

			if v := origin.GetString("fingerprint"); v != testValue {
				t.Fatalf("Expected field value %q, got %q", testValue, v)
			}
		})
	}
}

func TestAuthOriginCreated(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	origin := core.NewAuthOrigin(app)

	if v := origin.Created().String(); v != "" {
		t.Fatalf("Expected empty created, got %q", v)
	}

	now := types.NowDateTime()
	origin.SetRaw("created", now)

	if v := origin.Created().String(); v != now.String() {
		t.Fatalf("Expected %q created, got %q", now.String(), v)
	}
}

func TestAuthOriginUpdated(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	origin := core.NewAuthOrigin(app)

	if v := origin.Updated().String(); v != "" {
		t.Fatalf("Expected empty updated, got %q", v)
	}

	now := types.NowDateTime()
	origin.SetRaw("updated", now)

	if v := origin.Updated().String(); v != now.String() {
		t.Fatalf("Expected %q updated, got %q", now.String(), v)
	}
}

func TestAuthOriginPreValidate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	originsCol, err := app.FindCollectionByNameOrId(core.CollectionNameAuthOrigins)
	if err != nil {
		t.Fatal(err)
	}

	user, err := app.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("no proxy record", func(t *testing.T) {
		origin := &core.AuthOrigin{}

		if err := app.Validate(origin); err == nil {
			t.Fatal("Expected collection validation error")
		}
	})

	t.Run("non-AuthOrigin collection", func(t *testing.T) {
		origin := &core.AuthOrigin{}
		origin.SetProxyRecord(core.NewRecord(core.NewBaseCollection("invalid")))
		origin.SetRecordRef(user.Id)
		origin.SetCollectionRef(user.Collection().Id)
		origin.SetFingerprint("abc")

		if err := app.Validate(origin); err == nil {
			t.Fatal("Expected collection validation error")
		}
	})

	t.Run("AuthOrigin collection", func(t *testing.T) {
		origin := &core.AuthOrigin{}
		origin.SetProxyRecord(core.NewRecord(originsCol))
		origin.SetRecordRef(user.Id)
		origin.SetCollectionRef(user.Collection().Id)
		origin.SetFingerprint("abc")

		if err := app.Validate(origin); err != nil {
			t.Fatalf("Expected nil validation error, got %v", err)
		}
	})
}

func TestAuthOriginValidateHook(t *testing.T) {
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
		origin       func() *core.AuthOrigin
		expectErrors []string
	}{
		{
			"empty",
			func() *core.AuthOrigin {
				return core.NewAuthOrigin(app)
			},
			[]string{"collectionRef", "recordRef", "fingerprint"},
		},
		{
			"non-auth collection",
			func() *core.AuthOrigin {
				origin := core.NewAuthOrigin(app)
				origin.SetCollectionRef(demo1.Collection().Id)
				origin.SetRecordRef(demo1.Id)
				origin.SetFingerprint("abc")
				return origin
			},
			[]string{"collectionRef"},
		},
		{
			"missing record id",
			func() *core.AuthOrigin {
				origin := core.NewAuthOrigin(app)
				origin.SetCollectionRef(user.Collection().Id)
				origin.SetRecordRef("missing")
				origin.SetFingerprint("abc")
				return origin
			},
			[]string{"recordRef"},
		},
		{
			"valid ref",
			func() *core.AuthOrigin {
				origin := core.NewAuthOrigin(app)
				origin.SetCollectionRef(user.Collection().Id)
				origin.SetRecordRef(user.Id)
				origin.SetFingerprint("abc")
				return origin
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			errs := app.Validate(s.origin())
			tests.TestValidationErrors(t, errs, s.expectErrors)
		})
	}
}

func TestAuthOriginPasswordChangeDeletion(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	// no auth origin associated with it
	user1, err := testApp.FindAuthRecordByEmail("users", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	superuser2, err := testApp.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test2@example.com")
	if err != nil {
		t.Fatal(err)
	}

	client1, err := testApp.FindAuthRecordByEmail("clients", "test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		record     *core.Record
		deletedIds []string
	}{
		{user1, nil},
		{superuser2, []string{"5798yh833k6w6w0", "ic55o70g4f8pcl4", "dmy260k6ksjr4ib"}},
		{client1, []string{"9r2j0m74260ur8i"}},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s_%s", i, s.record.Collection().Name, s.record.Id), func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			deletedIds := []string{}
			app.OnRecordDelete().BindFunc(func(e *core.RecordEvent) error {
				deletedIds = append(deletedIds, e.Record.Id)
				return e.Next()
			})

			s.record.SetPassword("new_password")

			err := app.Save(s.record)
			if err != nil {
				t.Fatal(err)
			}

			if len(deletedIds) != len(s.deletedIds) {
				t.Fatalf("Expected deleted ids\n%v\ngot\n%v", s.deletedIds, deletedIds)
			}

			for _, id := range s.deletedIds {
				if !slices.Contains(deletedIds, id) {
					t.Errorf("Expected to find deleted id %q in %v", id, deletedIds)
				}
			}
		})
	}
}
