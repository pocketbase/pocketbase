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
			&models.Admin{},
			"SELECT {{_admins}}.* FROM `_admins`",
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
			"missing",
			true,
		},
		// existing collection id
		{
			&models.Collection{},
			"wsmn24bux7wo113",
			false,
		},
		// existing admin id
		{
			&models.Admin{},
			"sbmbsdb40jyxf7h",
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

func TestDaoSaveWithInsertId(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	model := &models.Admin{}
	model.Id = "test"
	model.Email = "test_new@example.com"
	model.MarkAsNew()
	if err := testApp.Dao().Save(model); err != nil {
		t.Fatal(err)
	}

	// refresh
	model, _ = testApp.Dao().FindAdminById("test")

	if model == nil {
		t.Fatal("Failed to find admin with id 'test'")
	}

	expectedHooks := []string{"OnModelBeforeCreate", "OnModelAfterCreate"}
	for _, h := range expectedHooks {
		if v, ok := testApp.EventCalls[h]; !ok || v != 1 {
			t.Fatalf("Expected event %s to be called exactly one time, got %d", h, v)
		}
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

type dummyColumnValueMapper struct {
	models.Admin
}

func (a *dummyColumnValueMapper) ColumnValueMap() map[string]any {
	return map[string]any{
		"email":        a.Email,
		"passwordHash": a.PasswordHash,
		"tokenKey":     "custom_token_key",
	}
}

func TestDaoSaveWithColumnValueMapper(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	model := &dummyColumnValueMapper{}
	model.Id = "test_mapped_id" // explicitly set an id
	model.Email = "test_mapped_create@example.com"
	model.TokenKey = "test_unmapped_token_key" // not used in the map
	model.SetPassword("123456")
	model.MarkAsNew()
	if err := testApp.Dao().Save(model); err != nil {
		t.Fatal(err)
	}

	createdModel, _ := testApp.Dao().FindAdminById("test_mapped_id")
	if createdModel == nil {
		t.Fatal("[create] Failed to find model with id 'test_mapped_id'")
	}
	if createdModel.Email != model.Email {
		t.Fatalf("Expected model with email %q, got %q", model.Email, createdModel.Email)
	}
	if createdModel.TokenKey != "custom_token_key" {
		t.Fatalf("Expected model with tokenKey %q, got %q", "custom_token_key", createdModel.TokenKey)
	}

	model.Email = "test_mapped_update@example.com"
	model.Avatar = 9 // not mapped and expect to be ignored
	if err := testApp.Dao().Save(model); err != nil {
		t.Fatal(err)
	}

	updatedModel, _ := testApp.Dao().FindAdminById("test_mapped_id")
	if updatedModel == nil {
		t.Fatal("[update] Failed to find model with id 'test_mapped_id'")
	}
	if updatedModel.Email != model.Email {
		t.Fatalf("Expected model with email %q, got %q", model.Email, createdModel.Email)
	}
	if updatedModel.Avatar != 0 {
		t.Fatalf("Expected model avatar 0, got %v", updatedModel.Avatar)
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

	baseDao := testApp.Dao()

	baseDao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		return errors.New("before_create")
	}
	baseDao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		return errors.New("before_update")
	}
	baseDao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		return errors.New("before_delete")
	}

	existingModel, _ := testApp.Dao().FindAdminByEmail("test@example.com")

	// test create error
	// ---
	newModel := &models.Admin{}
	if err := baseDao.Save(newModel); err.Error() != "before_create" {
		t.Fatalf("Expected before_create error, got %v", err)
	}

	// test update error
	// ---
	if err := baseDao.Save(existingModel); err.Error() != "before_update" {
		t.Fatalf("Expected before_update error, got %v", err)
	}

	// test delete error
	// ---
	if err := baseDao.Delete(existingModel); err.Error() != "before_delete" {
		t.Fatalf("Expected before_delete error, got %v", err)
	}
}

func TestDaoTransactionHooksCallsOnFailure(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	beforeCreateFuncCalls := 0
	beforeUpdateFuncCalls := 0
	beforeDeleteFuncCalls := 0
	afterCreateFuncCalls := 0
	afterUpdateFuncCalls := 0
	afterDeleteFuncCalls := 0

	baseDao := testApp.Dao()

	baseDao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		beforeCreateFuncCalls++
		return nil
	}
	baseDao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		beforeUpdateFuncCalls++
		return nil
	}
	baseDao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		beforeDeleteFuncCalls++
		return nil
	}

	baseDao.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) {
		afterCreateFuncCalls++
	}
	baseDao.AfterUpdateFunc = func(eventDao *daos.Dao, m models.Model) {
		afterUpdateFuncCalls++
	}
	baseDao.AfterDeleteFunc = func(eventDao *daos.Dao, m models.Model) {
		afterDeleteFuncCalls++
	}

	existingModel, _ := testApp.Dao().FindAdminByEmail("test@example.com")

	baseDao.RunInTransaction(func(txDao1 *daos.Dao) error {
		return txDao1.RunInTransaction(func(txDao2 *daos.Dao) error {
			// test create
			// ---
			newModel := &models.Admin{}
			newModel.Email = "test_new1@example.com"
			newModel.SetPassword("123456")
			if err := txDao2.Save(newModel); err != nil {
				t.Fatal(err)
			}

			// test update (twice)
			// ---
			if err := txDao2.Save(existingModel); err != nil {
				t.Fatal(err)
			}
			if err := txDao2.Save(existingModel); err != nil {
				t.Fatal(err)
			}

			// test delete
			// ---
			if err := txDao2.Delete(existingModel); err != nil {
				t.Fatal(err)
			}

			return errors.New("test_tx_error")
		})
	})

	if beforeCreateFuncCalls != 1 {
		t.Fatalf("Expected beforeCreateFuncCalls to be called 1 times, got %d", beforeCreateFuncCalls)
	}
	if beforeUpdateFuncCalls != 2 {
		t.Fatalf("Expected beforeUpdateFuncCalls to be called 2 times, got %d", beforeUpdateFuncCalls)
	}
	if beforeDeleteFuncCalls != 1 {
		t.Fatalf("Expected beforeDeleteFuncCalls to be called 1 times, got %d", beforeDeleteFuncCalls)
	}
	if afterCreateFuncCalls != 0 {
		t.Fatalf("Expected afterCreateFuncCalls to be called 0 times, got %d", afterCreateFuncCalls)
	}
	if afterUpdateFuncCalls != 0 {
		t.Fatalf("Expected afterUpdateFuncCalls to be called 0 times, got %d", afterUpdateFuncCalls)
	}
	if afterDeleteFuncCalls != 0 {
		t.Fatalf("Expected afterDeleteFuncCalls to be called 0 times, got %d", afterDeleteFuncCalls)
	}
}

func TestDaoTransactionHooksCallsOnSuccess(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	beforeCreateFuncCalls := 0
	beforeUpdateFuncCalls := 0
	beforeDeleteFuncCalls := 0
	afterCreateFuncCalls := 0
	afterUpdateFuncCalls := 0
	afterDeleteFuncCalls := 0

	baseDao := testApp.Dao()

	baseDao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		beforeCreateFuncCalls++
		return nil
	}
	baseDao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		beforeUpdateFuncCalls++
		return nil
	}
	baseDao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		beforeDeleteFuncCalls++
		return nil
	}

	baseDao.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) {
		afterCreateFuncCalls++
	}
	baseDao.AfterUpdateFunc = func(eventDao *daos.Dao, m models.Model) {
		afterUpdateFuncCalls++
	}
	baseDao.AfterDeleteFunc = func(eventDao *daos.Dao, m models.Model) {
		afterDeleteFuncCalls++
	}

	existingModel, _ := testApp.Dao().FindAdminByEmail("test@example.com")

	baseDao.RunInTransaction(func(txDao1 *daos.Dao) error {
		return txDao1.RunInTransaction(func(txDao2 *daos.Dao) error {
			// test create
			// ---
			newModel := &models.Admin{}
			newModel.Email = "test_new1@example.com"
			newModel.SetPassword("123456")
			if err := txDao2.Save(newModel); err != nil {
				t.Fatal(err)
			}

			// test update (twice)
			// ---
			if err := txDao2.Save(existingModel); err != nil {
				t.Fatal(err)
			}
			if err := txDao2.Save(existingModel); err != nil {
				t.Fatal(err)
			}

			// test delete
			// ---
			if err := txDao2.Delete(existingModel); err != nil {
				t.Fatal(err)
			}

			return nil
		})
	})

	if beforeCreateFuncCalls != 1 {
		t.Fatalf("Expected beforeCreateFuncCalls to be called 1 times, got %d", beforeCreateFuncCalls)
	}
	if beforeUpdateFuncCalls != 2 {
		t.Fatalf("Expected beforeUpdateFuncCalls to be called 2 times, got %d", beforeUpdateFuncCalls)
	}
	if beforeDeleteFuncCalls != 1 {
		t.Fatalf("Expected beforeDeleteFuncCalls to be called 1 times, got %d", beforeDeleteFuncCalls)
	}
	if afterCreateFuncCalls != 1 {
		t.Fatalf("Expected afterCreateFuncCalls to be called 1 times, got %d", afterCreateFuncCalls)
	}
	if afterUpdateFuncCalls != 2 {
		t.Fatalf("Expected afterUpdateFuncCalls to be called 2 times, got %d", afterUpdateFuncCalls)
	}
	if afterDeleteFuncCalls != 1 {
		t.Fatalf("Expected afterDeleteFuncCalls to be called 1 times, got %d", afterDeleteFuncCalls)
	}
}
