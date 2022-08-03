package daos_test

import (
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestNew(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	dao := daos.New(testApp.DB())

	if dao.DB() != testApp.DB() {
		t.Fatal("The 2 db instances are different")
	}
}

func TestDaoModelQuery(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	dao := daos.New(testApp.DB())

	scenarios := []struct {
		model    models.Model
		expected string
	}{
		{
			&models.Collection{},
			"SELECT {{_collections}}.* FROM `_collections`",
		},
		{
			&models.User{},
			"SELECT {{_users}}.* FROM `_users`",
		},
		{
			&models.Request{},
			"SELECT {{_requests}}.* FROM `_requests`",
		},
	}

	for i, scenario := range scenarios {
		sql := dao.ModelQuery(scenario.model).Build().SQL()
		if sql != scenario.expected {
			t.Errorf("(%d) Expected select %s, got %s", i, scenario.expected, sql)
		}
	}
}

func TestDaoFindById(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	scenarios := []struct {
		model       models.Model
		id          string
		expectError bool
	}{
		// missing id
		{
			&models.Collection{},
			"00000000-075d-49fe-9d09-ea7e951000dc",
			true,
		},
		// existing collection id
		{
			&models.Collection{},
			"3f2888f8-075d-49fe-9d09-ea7e951000dc",
			false,
		},
		// existing user id
		{
			&models.User{},
			"97cc3d3d-6ba2-383f-b42a-7bc84d27410c",
			false,
		},
	}

	for i, scenario := range scenarios {
		err := testApp.Dao().FindById(scenario.model, scenario.id)
		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expectError, err)
		}

		if !scenario.expectError && scenario.id != scenario.model.GetId() {
			t.Errorf("(%d) Expected model with id %v, got %v", i, scenario.id, scenario.model.GetId())
		}
	}
}

func TestDaoRunInTransaction(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	// failed nested transaction
	testApp.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		admin, _ := txDao.FindAdminByEmail("test@example.com")

		return txDao.RunInTransaction(func(tx2Dao *daos.Dao) error {
			if err := tx2Dao.DeleteAdmin(admin); err != nil {
				t.Fatal(err)
			}
			return errors.New("test error")
		})
	})

	// admin should still exist
	admin1, _ := testApp.Dao().FindAdminByEmail("test@example.com")
	if admin1 == nil {
		t.Fatal("Expected admin test@example.com to not be deleted")
	}

	// successful nested transaction
	testApp.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		admin, _ := txDao.FindAdminByEmail("test@example.com")

		return txDao.RunInTransaction(func(tx2Dao *daos.Dao) error {
			return tx2Dao.DeleteAdmin(admin)
		})
	})

	// admin should have been deleted
	admin2, _ := testApp.Dao().FindAdminByEmail("test@example.com")
	if admin2 != nil {
		t.Fatalf("Expected admin test@example.com to be deleted, found %v", admin2)
	}
}

func TestDaoSaveCreate(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	model := &models.Admin{}
	model.Email = "test_new@example.com"
	model.Avatar = 8
	if err := testApp.Dao().Save(model); err != nil {
		t.Fatal(err)
	}

	// refresh
	model, _ = testApp.Dao().FindAdminByEmail("test_new@example.com")

	if model.Avatar != 8 {
		t.Fatalf("Expected model avatar field to be 8, got %v", model.Avatar)
	}

	expectedHooks := []string{"OnModelBeforeCreate", "OnModelAfterCreate"}
	for _, h := range expectedHooks {
		if v, ok := testApp.EventCalls[h]; !ok || v != 1 {
			t.Fatalf("Expected event %s to be called exactly one time, got %d", h, v)
		}
	}
}

type TestAdminMapper struct {
	models.Admin
}

func (m *TestAdminMapper) ColumnValueMap() map[string]any {
	return map[string]any{
		"email":    m.Email,
		"avatar":   m.Avatar,
		"tokenKey": m.TokenKey,
		"passwordHash":m.PasswordHash,
	}
}

func TestDaoSaveCreateWithMapper(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	model := &models.Admin{}
	model.Email = "test_new@example.com"
	model.Avatar = 8
	mapper := &TestAdminMapper{
		*model,
	}
	if err := testApp.Dao().Save(mapper); err != nil {
		t.Fatal(err)
	}

	// refresh
	_, err := testApp.Dao().FindAdminByEmail("test_new@example.com")

	if err != nil {
		t.Fatalf("Expected to get saved admin, failed with error %v", err)
	}
}

func TestDaoSaveUpdate(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	model, _ := testApp.Dao().FindAdminByEmail("test@example.com")

	model.Avatar = 8
	if err := testApp.Dao().Save(model); err != nil {
		t.Fatal(err)
	}

	// refresh
	model, _ = testApp.Dao().FindAdminByEmail("test@example.com")

	if model.Avatar != 8 {
		t.Fatalf("Expected model avatar field to be updated to 8, got %v", model.Avatar)
	}

	expectedHooks := []string{"OnModelBeforeUpdate", "OnModelAfterUpdate"}
	for _, h := range expectedHooks {
		if v, ok := testApp.EventCalls[h]; !ok || v != 1 {
			t.Fatalf("Expected event %s to be called exactly one time, got %d", h, v)
		}
	}
}

func TestDaoDelete(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	model, _ := testApp.Dao().FindAdminByEmail("test@example.com")

	if err := testApp.Dao().Delete(model); err != nil {
		t.Fatal(err)
	}

	model, _ = testApp.Dao().FindAdminByEmail("test@example.com")
	if model != nil {
		t.Fatalf("Expected model to be deleted, found %v", model)
	}

	expectedHooks := []string{"OnModelBeforeDelete", "OnModelAfterDelete"}
	for _, h := range expectedHooks {
		if v, ok := testApp.EventCalls[h]; !ok || v != 1 {
			t.Fatalf("Expected event %s to be called exactly one time, got %d", h, v)
		}
	}
}

func TestDaoBeforeHooksError(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	testApp.Dao().BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		return errors.New("before_create")
	}
	testApp.Dao().BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		return errors.New("before_update")
	}
	testApp.Dao().BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		return errors.New("before_delete")
	}

	existingModel, _ := testApp.Dao().FindAdminByEmail("test@example.com")

	// try to create
	// ---
	newModel := &models.Admin{}
	newModel.Email = "test_new@example.com"
	if err := testApp.Dao().Save(newModel); err.Error() != "before_create" {
		t.Fatalf("Expected before_create error, got %v", err)
	}

	// try to update
	// ---
	if err := testApp.Dao().Save(existingModel); err.Error() != "before_update" {
		t.Fatalf("Expected before_update error, got %v", err)
	}

	// try to delete
	// ---
	if err := testApp.Dao().Delete(existingModel); err.Error() != "before_delete" {
		t.Fatalf("Expected before_delete error, got %v", err)
	}
}
