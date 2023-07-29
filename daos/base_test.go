package daos_test

import (
	"errors"
	"testing"
	"time"

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

func TestNewMultiDB(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	dao := daos.NewMultiDB(testApp.Dao().ConcurrentDB(), testApp.Dao().NonconcurrentDB())

	if dao.DB() != testApp.Dao().ConcurrentDB() {
		t.Fatal("[db-concurrentDB] The 2 db instances are different")
	}

	if dao.ConcurrentDB() != testApp.Dao().ConcurrentDB() {
		t.Fatal("[concurrentDB-concurrentDB] The 2 db instances are different")
	}

	if dao.NonconcurrentDB() != testApp.Dao().NonconcurrentDB() {
		t.Fatal("[nonconcurrentDB-nonconcurrentDB] The 2 db instances are different")
	}
}

func TestDaoClone(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	hookCalls := map[string]int{}

	dao := daos.NewMultiDB(testApp.Dao().ConcurrentDB(), testApp.Dao().NonconcurrentDB())
	dao.MaxLockRetries = 1
	dao.ModelQueryTimeout = 2
	dao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		hookCalls["BeforeDeleteFunc"]++
		return action()
	}
	dao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		hookCalls["BeforeUpdateFunc"]++
		return action()
	}
	dao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		hookCalls["BeforeCreateFunc"]++
		return action()
	}
	dao.AfterDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		hookCalls["AfterDeleteFunc"]++
		return nil
	}
	dao.AfterUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		hookCalls["AfterUpdateFunc"]++
		return nil
	}
	dao.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		hookCalls["AfterCreateFunc"]++
		return nil
	}

	clone := dao.Clone()
	clone.MaxLockRetries = 3
	clone.ModelQueryTimeout = 4
	clone.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		hookCalls["NewAfterCreateFunc"]++
		return nil
	}

	if dao.MaxLockRetries == clone.MaxLockRetries {
		t.Fatal("Expected different MaxLockRetries")
	}

	if dao.ModelQueryTimeout == clone.ModelQueryTimeout {
		t.Fatal("Expected different ModelQueryTimeout")
	}

	emptyAction := func() error { return nil }

	// trigger hooks
	dao.BeforeDeleteFunc(nil, nil, emptyAction)
	dao.BeforeUpdateFunc(nil, nil, emptyAction)
	dao.BeforeCreateFunc(nil, nil, emptyAction)
	dao.AfterDeleteFunc(nil, nil)
	dao.AfterUpdateFunc(nil, nil)
	dao.AfterCreateFunc(nil, nil)
	clone.BeforeDeleteFunc(nil, nil, emptyAction)
	clone.BeforeUpdateFunc(nil, nil, emptyAction)
	clone.BeforeCreateFunc(nil, nil, emptyAction)
	clone.AfterDeleteFunc(nil, nil)
	clone.AfterUpdateFunc(nil, nil)
	clone.AfterCreateFunc(nil, nil)

	expectations := []struct {
		hook  string
		total int
	}{
		{"BeforeDeleteFunc", 2},
		{"BeforeUpdateFunc", 2},
		{"BeforeCreateFunc", 2},
		{"AfterDeleteFunc", 2},
		{"AfterUpdateFunc", 2},
		{"AfterCreateFunc", 1},
		{"NewAfterCreateFunc", 1},
	}

	for _, e := range expectations {
		if hookCalls[e.hook] != e.total {
			t.Errorf("Expected %s to be caleed %d", e.hook, e.total)
		}
	}
}

func TestDaoWithoutHooks(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	hookCalls := map[string]int{}

	dao := daos.NewMultiDB(testApp.Dao().ConcurrentDB(), testApp.Dao().NonconcurrentDB())
	dao.MaxLockRetries = 1
	dao.ModelQueryTimeout = 2
	dao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		hookCalls["BeforeDeleteFunc"]++
		return action()
	}
	dao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		hookCalls["BeforeUpdateFunc"]++
		return action()
	}
	dao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		hookCalls["BeforeCreateFunc"]++
		return action()
	}
	dao.AfterDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		hookCalls["AfterDeleteFunc"]++
		return nil
	}
	dao.AfterUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		hookCalls["AfterUpdateFunc"]++
		return nil
	}
	dao.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		hookCalls["AfterCreateFunc"]++
		return nil
	}

	new := dao.WithoutHooks()

	if new.MaxLockRetries != dao.MaxLockRetries {
		t.Fatalf("Expected MaxLockRetries %d, got %d", new.Clone().MaxLockRetries, dao.MaxLockRetries)
	}

	if new.ModelQueryTimeout != dao.ModelQueryTimeout {
		t.Fatalf("Expected ModelQueryTimeout %d, got %d", new.Clone().ModelQueryTimeout, dao.ModelQueryTimeout)
	}

	if new.BeforeDeleteFunc != nil {
		t.Fatal("Expected BeforeDeleteFunc to be nil")
	}

	if new.BeforeUpdateFunc != nil {
		t.Fatal("Expected BeforeUpdateFunc to be nil")
	}

	if new.BeforeCreateFunc != nil {
		t.Fatal("Expected BeforeCreateFunc to be nil")
	}

	if new.AfterDeleteFunc != nil {
		t.Fatal("Expected AfterDeleteFunc to be nil")
	}

	if new.AfterUpdateFunc != nil {
		t.Fatal("Expected AfterUpdateFunc to be nil")
	}

	if new.AfterCreateFunc != nil {
		t.Fatal("Expected AfterCreateFunc to be nil")
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

func TestDaoModelQueryCancellation(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	dao := daos.New(testApp.DB())

	m := &models.Admin{}

	if err := dao.ModelQuery(m).One(m); err != nil {
		t.Fatalf("Failed to execute control query: %v", err)
	}

	dao.ModelQueryTimeout = 0 * time.Millisecond
	if err := dao.ModelQuery(m).One(m); err == nil {
		t.Fatal("Expected to be cancelled, got nil")
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

func TestDaoRetryCreate(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	// init mock retry dao
	retryBeforeCreateHookCalls := 0
	retryAfterCreateHookCalls := 0
	retryDao := daos.New(testApp.DB())
	retryDao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		retryBeforeCreateHookCalls++
		return errors.New("database is locked")
	}
	retryDao.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		retryAfterCreateHookCalls++
		return nil
	}

	model := &models.Admin{Email: "new@example.com"}
	if err := retryDao.Save(model); err != nil {
		t.Fatalf("Expected nil after retry, got error: %v", err)
	}

	// the before hook is expected to be called only once because
	// it is ignored after the first "database is locked" error
	if retryBeforeCreateHookCalls != 1 {
		t.Fatalf("Expected before hook calls to be 1, got %d", retryBeforeCreateHookCalls)
	}

	if retryAfterCreateHookCalls != 1 {
		t.Fatalf("Expected after hook calls to be 1, got %d", retryAfterCreateHookCalls)
	}

	// with non-locking error
	retryBeforeCreateHookCalls = 0
	retryAfterCreateHookCalls = 0
	retryDao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		retryBeforeCreateHookCalls++
		return errors.New("non-locking error")
	}

	dummy := &models.Admin{Email: "test@example.com"}
	if err := retryDao.Save(dummy); err == nil {
		t.Fatal("Expected error, got nil")
	}

	if retryBeforeCreateHookCalls != 1 {
		t.Fatalf("Expected before hook calls to be 1, got %d", retryBeforeCreateHookCalls)
	}

	if retryAfterCreateHookCalls != 0 {
		t.Fatalf("Expected after hook calls to be 0, got %d", retryAfterCreateHookCalls)
	}
}

func TestDaoRetryUpdate(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	model, err := testApp.Dao().FindAdminByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}

	// init mock retry dao
	retryBeforeUpdateHookCalls := 0
	retryAfterUpdateHookCalls := 0
	retryDao := daos.New(testApp.DB())
	retryDao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		retryBeforeUpdateHookCalls++
		return errors.New("database is locked")
	}
	retryDao.AfterUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		retryAfterUpdateHookCalls++
		return nil
	}

	if err := retryDao.Save(model); err != nil {
		t.Fatalf("Expected nil after retry, got error: %v", err)
	}

	// the before hook is expected to be called only once because
	// it is ignored after the first "database is locked" error
	if retryBeforeUpdateHookCalls != 1 {
		t.Fatalf("Expected before hook calls to be 1, got %d", retryBeforeUpdateHookCalls)
	}

	if retryAfterUpdateHookCalls != 1 {
		t.Fatalf("Expected after hook calls to be 1, got %d", retryAfterUpdateHookCalls)
	}

	// with non-locking error
	retryBeforeUpdateHookCalls = 0
	retryAfterUpdateHookCalls = 0
	retryDao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		retryBeforeUpdateHookCalls++
		return errors.New("non-locking error")
	}

	if err := retryDao.Save(model); err == nil {
		t.Fatal("Expected error, got nil")
	}

	if retryBeforeUpdateHookCalls != 1 {
		t.Fatalf("Expected before hook calls to be 1, got %d", retryBeforeUpdateHookCalls)
	}

	if retryAfterUpdateHookCalls != 0 {
		t.Fatalf("Expected after hook calls to be 0, got %d", retryAfterUpdateHookCalls)
	}
}

func TestDaoRetryDelete(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	// init mock retry dao
	retryBeforeDeleteHookCalls := 0
	retryAfterDeleteHookCalls := 0
	retryDao := daos.New(testApp.DB())
	retryDao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		retryBeforeDeleteHookCalls++
		return errors.New("database is locked")
	}
	retryDao.AfterDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		retryAfterDeleteHookCalls++
		return nil
	}

	model, _ := retryDao.FindAdminByEmail("test@example.com")
	if err := retryDao.Delete(model); err != nil {
		t.Fatalf("Expected nil after retry, got error: %v", err)
	}

	// the before hook is expected to be called only once because
	// it is ignored after the first "database is locked" error
	if retryBeforeDeleteHookCalls != 1 {
		t.Fatalf("Expected before hook calls to be 1, got %d", retryBeforeDeleteHookCalls)
	}

	if retryAfterDeleteHookCalls != 1 {
		t.Fatalf("Expected after hook calls to be 1, got %d", retryAfterDeleteHookCalls)
	}

	// with non-locking error
	retryBeforeDeleteHookCalls = 0
	retryAfterDeleteHookCalls = 0
	retryDao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		retryBeforeDeleteHookCalls++
		return errors.New("non-locking error")
	}

	dummy := &models.Admin{}
	dummy.RefreshId()
	dummy.MarkAsNotNew()
	if err := retryDao.Delete(dummy); err == nil {
		t.Fatal("Expected error, got nil")
	}

	if retryBeforeDeleteHookCalls != 1 {
		t.Fatalf("Expected before hook calls to be 1, got %d", retryBeforeDeleteHookCalls)
	}

	if retryAfterDeleteHookCalls != 0 {
		t.Fatalf("Expected after hook calls to be 0, got %d", retryAfterDeleteHookCalls)
	}
}

func TestDaoBeforeHooksError(t *testing.T) {
	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	baseDao := testApp.Dao()

	baseDao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		return errors.New("before_create")
	}
	baseDao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		return errors.New("before_update")
	}
	baseDao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
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

	baseDao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		beforeCreateFuncCalls++
		return action()
	}
	baseDao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		beforeUpdateFuncCalls++
		return action()
	}
	baseDao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		beforeDeleteFuncCalls++
		return action()
	}

	baseDao.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		afterCreateFuncCalls++
		return nil
	}
	baseDao.AfterUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		afterUpdateFuncCalls++
		return nil
	}
	baseDao.AfterDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		afterDeleteFuncCalls++
		return nil
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

	baseDao.BeforeCreateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		beforeCreateFuncCalls++
		return action()
	}
	baseDao.BeforeUpdateFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		beforeUpdateFuncCalls++
		return action()
	}
	baseDao.BeforeDeleteFunc = func(eventDao *daos.Dao, m models.Model, action func() error) error {
		beforeDeleteFuncCalls++
		return action()
	}

	baseDao.AfterCreateFunc = func(eventDao *daos.Dao, m models.Model) error {
		afterCreateFuncCalls++
		return nil
	}
	baseDao.AfterUpdateFunc = func(eventDao *daos.Dao, m models.Model) error {
		afterUpdateFuncCalls++
		return nil
	}
	baseDao.AfterDeleteFunc = func(eventDao *daos.Dao, m models.Model) error {
		afterDeleteFuncCalls++
		return nil
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
