package issues

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// https://github.com/fondoger/pocketbase/issues/35
func TestIssue35_DateTimeFieldWithEmptyValue(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewBaseCollection("test_issue_35")
	collection.Fields.Add(&core.DateField{
		Name: "updated",
	})

	err := app.Save(collection)
	if err != nil {
		t.Fatal(err)
	}

	record := core.NewRecord(collection)
	record.Set("updated", "2025-09-01")
	err = app.Save(record)
	if err != nil {
		t.Fatal(err)
	}

	record = core.NewRecord(collection)
	err = app.Save(record)
	if err != nil {
		t.Fatal(err)
	}
}
