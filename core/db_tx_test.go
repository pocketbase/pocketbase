package core_test

import (
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRunInTransaction(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	t.Run("failed nested transaction", func(t *testing.T) {
		app.RunInTransaction(func(txApp core.App) error {
			superuser, _ := txApp.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")

			return txApp.RunInTransaction(func(tx2Dao core.App) error {
				if err := tx2Dao.Delete(superuser); err != nil {
					t.Fatal(err)
				}
				return errors.New("test error")
			})
		})

		// superuser should still exist
		superuser, _ := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")
		if superuser == nil {
			t.Fatal("Expected superuser test@example.com to not be deleted")
		}
	})

	t.Run("successful nested transaction", func(t *testing.T) {
		app.RunInTransaction(func(txApp core.App) error {
			superuser, _ := txApp.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")

			return txApp.RunInTransaction(func(tx2Dao core.App) error {
				return tx2Dao.Delete(superuser)
			})
		})

		// superuser should have been deleted
		superuser, _ := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")
		if superuser != nil {
			t.Fatalf("Expected superuser test@example.com to be deleted, found %v", superuser)
		}
	})
}

func TestTransactionHooksCallsOnFailure(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	createHookCalls := 0
	updateHookCalls := 0
	deleteHookCalls := 0
	afterCreateHookCalls := 0
	afterUpdateHookCalls := 0
	afterDeleteHookCalls := 0

	app.OnModelCreate().BindFunc(func(e *core.ModelEvent) error {
		createHookCalls++
		return e.Next()
	})

	app.OnModelUpdate().BindFunc(func(e *core.ModelEvent) error {
		updateHookCalls++
		return e.Next()
	})

	app.OnModelDelete().BindFunc(func(e *core.ModelEvent) error {
		deleteHookCalls++
		return e.Next()
	})

	app.OnModelAfterCreateSuccess().BindFunc(func(e *core.ModelEvent) error {
		afterCreateHookCalls++
		return e.Next()
	})

	app.OnModelAfterUpdateSuccess().BindFunc(func(e *core.ModelEvent) error {
		afterUpdateHookCalls++
		return e.Next()
	})

	app.OnModelAfterDeleteSuccess().BindFunc(func(e *core.ModelEvent) error {
		afterDeleteHookCalls++
		return e.Next()
	})

	existingModel, _ := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")

	app.RunInTransaction(func(txApp1 core.App) error {
		return txApp1.RunInTransaction(func(txApp2 core.App) error {
			// test create
			// ---
			newModel := core.NewRecord(existingModel.Collection())
			newModel.SetEmail("test_new1@example.com")
			newModel.SetPassword("1234567890")
			if err := txApp2.Save(newModel); err != nil {
				t.Fatal(err)
			}

			// test update (twice)
			// ---
			if err := txApp2.Save(existingModel); err != nil {
				t.Fatal(err)
			}
			if err := txApp2.Save(existingModel); err != nil {
				t.Fatal(err)
			}

			// test delete
			// ---
			if err := txApp2.Delete(newModel); err != nil {
				t.Fatal(err)
			}

			return errors.New("test_tx_error")
		})
	})

	if createHookCalls != 1 {
		t.Errorf("Expected createHookCalls to be called 1 time, got %d", createHookCalls)
	}
	if updateHookCalls != 2 {
		t.Errorf("Expected updateHookCalls to be called 2 times, got %d", updateHookCalls)
	}
	if deleteHookCalls != 1 {
		t.Errorf("Expected deleteHookCalls to be called 1 time, got %d", deleteHookCalls)
	}
	if afterCreateHookCalls != 0 {
		t.Errorf("Expected afterCreateHookCalls to be called 0 times, got %d", afterCreateHookCalls)
	}
	if afterUpdateHookCalls != 0 {
		t.Errorf("Expected afterUpdateHookCalls to be called 0 times, got %d", afterUpdateHookCalls)
	}
	if afterDeleteHookCalls != 0 {
		t.Errorf("Expected afterDeleteHookCalls to be called 0 times, got %d", afterDeleteHookCalls)
	}
}

func TestTransactionHooksCallsOnSuccess(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	createHookCalls := 0
	updateHookCalls := 0
	deleteHookCalls := 0
	afterCreateHookCalls := 0
	afterUpdateHookCalls := 0
	afterDeleteHookCalls := 0

	app.OnModelCreate().BindFunc(func(e *core.ModelEvent) error {
		createHookCalls++
		return e.Next()
	})

	app.OnModelUpdate().BindFunc(func(e *core.ModelEvent) error {
		updateHookCalls++
		return e.Next()
	})

	app.OnModelDelete().BindFunc(func(e *core.ModelEvent) error {
		deleteHookCalls++
		return e.Next()
	})

	app.OnModelAfterCreateSuccess().BindFunc(func(e *core.ModelEvent) error {
		if e.App.IsTransactional() {
			t.Fatal("Expected e.App to be non-transactional")
		}

		afterCreateHookCalls++
		return e.Next()
	})

	app.OnModelAfterUpdateSuccess().BindFunc(func(e *core.ModelEvent) error {
		if e.App.IsTransactional() {
			t.Fatal("Expected e.App to be non-transactional")
		}

		afterUpdateHookCalls++
		return e.Next()
	})

	app.OnModelAfterDeleteSuccess().BindFunc(func(e *core.ModelEvent) error {
		if e.App.IsTransactional() {
			t.Fatal("Expected e.App to be non-transactional")
		}

		afterDeleteHookCalls++
		return e.Next()
	})

	existingModel, _ := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")

	app.RunInTransaction(func(txApp1 core.App) error {
		return txApp1.RunInTransaction(func(txApp2 core.App) error {
			// test create
			// ---
			newModel := core.NewRecord(existingModel.Collection())
			newModel.SetEmail("test_new1@example.com")
			newModel.SetPassword("1234567890")
			if err := txApp2.Save(newModel); err != nil {
				t.Fatal(err)
			}

			// test update (twice)
			// ---
			if err := txApp2.Save(existingModel); err != nil {
				t.Fatal(err)
			}
			if err := txApp2.Save(existingModel); err != nil {
				t.Fatal(err)
			}

			// test delete
			// ---
			if err := txApp2.Delete(newModel); err != nil {
				t.Fatal(err)
			}

			return nil
		})
	})

	if createHookCalls != 1 {
		t.Errorf("Expected createHookCalls to be called 1 time, got %d", createHookCalls)
	}
	if updateHookCalls != 2 {
		t.Errorf("Expected updateHookCalls to be called 2 times, got %d", updateHookCalls)
	}
	if deleteHookCalls != 1 {
		t.Errorf("Expected deleteHookCalls to be called 1 time, got %d", deleteHookCalls)
	}
	if afterCreateHookCalls != 1 {
		t.Errorf("Expected afterCreateHookCalls to be called 1 time, got %d", afterCreateHookCalls)
	}
	if afterUpdateHookCalls != 2 {
		t.Errorf("Expected afterUpdateHookCalls to be called 2 times, got %d", afterUpdateHookCalls)
	}
	if afterDeleteHookCalls != 1 {
		t.Errorf("Expected afterDeleteHookCalls to be called 1 time, got %d", afterDeleteHookCalls)
	}
}

func TestTransactionFromInnerCreateHook(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	app.OnRecordCreateExecute("demo2").BindFunc(func(e *core.RecordEvent) error {
		originalApp := e.App
		return e.App.RunInTransaction(func(txApp core.App) error {
			e.App = txApp
			defer func() {
				e.App = originalApp
			}()

			nextErr := e.Next()

			return nextErr
		})
	})

	app.OnRecordAfterCreateSuccess("demo2").BindFunc(func(e *core.RecordEvent) error {
		if e.App.IsTransactional() {
			t.Fatal("Expected e.App to be non-transactional")
		}

		// perform a db query with the app instance to ensure that it is still valid
		_, err := e.App.FindFirstRecordByFilter("demo2", "1=1")
		if err != nil {
			t.Fatalf("Failed to perform a db query after tx success: %v", err)
		}

		return e.Next()
	})

	collection, err := app.FindCollectionByNameOrId("demo2")
	if err != nil {
		t.Fatal(err)
	}

	record := core.NewRecord(collection)

	record.Set("title", "test_inner_tx")

	if err = app.Save(record); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	expectedHookCalls := map[string]int{
		"OnRecordCreateExecute":      1,
		"OnRecordAfterCreateSuccess": 1,
	}
	for k, total := range expectedHookCalls {
		if found, ok := app.EventCalls[k]; !ok || total != found {
			t.Fatalf("Expected %q %d calls, got %d", k, total, found)
		}
	}
}

func TestTransactionFromInnerUpdateHook(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	app.OnRecordUpdateExecute("demo2").BindFunc(func(e *core.RecordEvent) error {
		originalApp := e.App
		return e.App.RunInTransaction(func(txApp core.App) error {
			e.App = txApp
			defer func() {
				e.App = originalApp
			}()

			nextErr := e.Next()

			return nextErr
		})
	})

	app.OnRecordAfterUpdateSuccess("demo2").BindFunc(func(e *core.RecordEvent) error {
		if e.App.IsTransactional() {
			t.Fatal("Expected e.App to be non-transactional")
		}

		// perform a db query with the app instance to ensure that it is still valid
		_, err := e.App.FindFirstRecordByFilter("demo2", "1=1")
		if err != nil {
			t.Fatalf("Failed to perform a db query after tx success: %v", err)
		}

		return e.Next()
	})

	existingModel, err := app.FindFirstRecordByFilter("demo2", "1=1")
	if err != nil {
		t.Fatal(err)
	}

	if err = app.Save(existingModel); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	expectedHookCalls := map[string]int{
		"OnRecordUpdateExecute":      1,
		"OnRecordAfterUpdateSuccess": 1,
	}
	for k, total := range expectedHookCalls {
		if found, ok := app.EventCalls[k]; !ok || total != found {
			t.Fatalf("Expected %q %d calls, got %d", k, total, found)
		}
	}
}

func TestTransactionFromInnerDeleteHook(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	app.OnRecordDeleteExecute("demo2").BindFunc(func(e *core.RecordEvent) error {
		originalApp := e.App
		return e.App.RunInTransaction(func(txApp core.App) error {
			e.App = txApp
			defer func() {
				e.App = originalApp
			}()

			nextErr := e.Next()

			return nextErr
		})
	})

	app.OnRecordAfterDeleteSuccess("demo2").BindFunc(func(e *core.RecordEvent) error {
		if e.App.IsTransactional() {
			t.Fatal("Expected e.App to be non-transactional")
		}

		// perform a db query with the app instance to ensure that it is still valid
		_, err := e.App.FindFirstRecordByFilter("demo2", "1=1")
		if err != nil {
			t.Fatalf("Failed to perform a db query after tx success: %v", err)
		}

		return e.Next()
	})

	existingModel, err := app.FindFirstRecordByFilter("demo2", "1=1")
	if err != nil {
		t.Fatal(err)
	}

	if err = app.Delete(existingModel); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	expectedHookCalls := map[string]int{
		"OnRecordDeleteExecute":      1,
		"OnRecordAfterDeleteSuccess": 1,
	}
	for k, total := range expectedHookCalls {
		if found, ok := app.EventCalls[k]; !ok || total != found {
			t.Fatalf("Expected %q %d calls, got %d", k, total, found)
		}
	}
}
