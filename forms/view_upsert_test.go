package forms_test

import (
	"encoding/json"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tests"
)

func TestViewUpsertPanic1(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewViewUpsert(nil, nil)
}

func TestViewUpsertPanic2(t *testing.T) {
	app, _ := tests.NewTestApp(getViewTestData())
	defer app.Cleanup()

	defer func() {
		if recover() == nil {
			t.Fatal("The form did not panic")
		}
	}()

	forms.NewCollectionUpsert(app, nil)
}

func TestNewViewUpsert(t *testing.T) {
	app, _ := tests.NewTestApp(getViewTestData())
	defer app.Cleanup()
	view := &models.View{}
	view.Name = "test"
	listRule := "testview"
	view.ListRule = &listRule
	view.Sql = "SELECT 1"
	form := forms.NewViewUpsert(app, view)

	if form.Name != view.Name {
		t.Errorf("Expected Name %q, got %q", view.Name, form.Name)
	}

	if form.ListRule != view.ListRule {
		t.Errorf("Expected ListRule %v, got %v", view.ListRule, form.ListRule)
	}

	if form.Sql != view.Sql {
		t.Errorf("Expected Sql %v, got %v", view.ListRule, form.ListRule)
	}
	// store previous state and modify the collection schema to verify
	// that the form.Schema is a deep clone
}

func getViewTestData() string {
	_, currentFile, _, _ := runtime.Caller(0)
	testDataDir := filepath.Join(path.Dir(currentFile), "../tests/view_data")
	return testDataDir
}

func TestViewUpsertValidate(t *testing.T) {
	app, _ := tests.NewTestApp(getViewTestData())
	defer app.Cleanup()

	scenarios := []struct {
		jsonData       string
		expectedErrors []string
	}{
		{"{}", []string{"name", "sql"}},
		{
			`{
				"name": "test ?!@#$",
				"sql": "select 1;delete from _admins",
				"listRule": "missing = '123'"
			}`,
			[]string{"name", "sql"},
		},
		{
			`{
				"name": "test",
				"listRule": "test='123'",
                "sql":"select 1"
			}`,
			[]string{},
		},
	}

	for i, s := range scenarios {
		form := forms.NewViewUpsert(app, &models.View{}) // load data
		loadErr := json.Unmarshal([]byte(s.jsonData), form)
		if loadErr != nil {
			t.Errorf("(%d) Failed to load form data: %v", i, loadErr)
			continue
		}

		// parse errors
		result := form.Validate()
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("(%d) Failed to parse errors %v", i, result)
			continue
		}

		// check errors
		if len(errs) != len(s.expectedErrors) {
			t.Errorf("(%d) Expected error keys %v, got %v", i, s.expectedErrors, errs)
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("(%d) Missing expected error key %q in %v", i, k, errs)
			}
		}
	}
}
