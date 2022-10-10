package daos_test

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
)

func getViewTestData() string {
	_, currentFile, _, _ := runtime.Caller(0)
	testDataDir := filepath.Join(path.Dir(currentFile), "../tests/view_data")
	return testDataDir
}

func TestViewQuery(t *testing.T) {
	app, _ := tests.NewTestApp(getViewTestData())
	defer app.Cleanup()

	expected := "SELECT {{_views}}.* FROM `_views`"

	sql := app.Dao().ViewQuery().Build().SQL()
	if sql != expected {
		t.Errorf("Expected sql %s, got %s", expected, sql)
	}
}

func TestFindViewByName(t *testing.T) {
	app, _ := tests.NewTestApp(getViewTestData())
	defer app.Cleanup()

	scenarios := []struct {
		name        string
		expectError bool
	}{
		{"", true},
		{"missing", true},
		{"view_demo", false},
	}

	for i, scenario := range scenarios {
		model, err := app.Dao().FindViewByName(scenario.name)
		hasErr := err != nil
		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr to be %v, got %v (%v)", i, scenario.expectError, hasErr, err)
		}

		if model != nil && model.Id != scenario.name && model.Name != scenario.name {
			t.Errorf("(%d) Expected model with identifier %s, got %v", i, scenario.name, model)
		}
	}
}

func TestIsViewNameUnique(t *testing.T) {
	app, _ := tests.NewTestApp(getViewTestData())
	defer app.Cleanup()

	scenarios := []struct {
		name      string
		excludeId string
		expected  bool
	}{
		{"", "", false},
		{"view_demo", "", false},
		{"new", "", true},
		{"view_demo", "da06mbbnabtu22j", true},
	}

	for i, scenario := range scenarios {
		result := app.Dao().IsViewNameUnique(scenario.name, scenario.excludeId)
		if result != scenario.expected {
			t.Errorf("(%d) Expected %v, got %v", i, scenario.expected, result)
		}
	}
}

func TestDeleteViews(t *testing.T) {
	app, _ := tests.NewTestApp(getViewTestData())
	defer app.Cleanup()

	c0 := &models.View{}
	c1, err := app.Dao().FindViewByName("view_demo")
	if err != nil {
		t.Fatal(err)
	}
	c2, err := app.Dao().FindViewByName("view_demo2")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		model       *models.View
		expectError bool
	}{
		{c0, true},
		{c1, false},
		{c2, false},
	}

	for i, scenario := range scenarios {
		err := app.Dao().DeleteView(scenario.model)
		hasErr := err != nil

		if hasErr != scenario.expectError {
			t.Errorf("(%d) Expected hasErr %v, got %v", i, scenario.expectError, hasErr)
		}
	}
}

func TestGetSchema(t *testing.T) {
	app, _ := tests.NewTestApp(getViewTestData())
	defer app.Cleanup()
	s, err := app.Dao().GetViewSchema("view_demo2")
	if err != nil {
		t.Fatal(err)
	}
	sMap := s.AsMap()
	if sMap["email"].Type != schema.FieldTypeText {
		t.Errorf("Expected 'email' type %v, got %v", schema.FieldTypeText, sMap["email"].Type)
	}
	if sMap["url"].Type != schema.FieldTypeText {
		t.Errorf("Expected 'url' type %v, got %v", schema.FieldTypeText, sMap["url"].Type)
	}
	if sMap["text"].Type != schema.FieldTypeText {
		t.Errorf("Expected 'text' type %v, got %v", schema.FieldTypeText, sMap["text"].Type)
	}
	if sMap["number"].Type != schema.FieldTypeNumber {
		t.Errorf("Expected 'number' type %v, got %v", schema.FieldTypeNumber, sMap["number"].Type)
	}
}

func TestSaveView(t *testing.T) {
	app, _ := tests.NewTestApp(getViewTestData())
	defer app.Cleanup()

	view := &models.View{
		Name: "new_test",
		Sql:  `SELECT * from profiles`,
	}

	err := app.Dao().SaveView(view)
	if err != nil {
		t.Fatal(err)
	}

	if view.Id == "" {
		t.Fatal("Expected view id to be set")
	}

	if !app.Dao().HasView(view.Name) {
		t.Fatalf("Expected View %s to be created", view.Name)
	}
	found := &models.View{}

	err = app.DB().Select().
		From("sqlite_master").
		AndWhere(dbx.HashExp{"type": "view"}).
		AndWhere(dbx.NewExp("LOWER([[name]])=LOWER({:viewName})", dbx.Params{"viewName": view.Name})).
		One(&found)

	if err != nil {
		t.Fatal(err)
	}

	if view.Name != found.Name {
		t.Fatalf("Expected View name to be %s instead got %s", view.Name, found.Name)
	}

	if stripViewSql(view.Sql) != stripViewSql(found.Sql) {
		t.Fatalf("Expected View sql to be %s instead got %s", stripViewSql(view.Sql), stripViewSql(found.Sql))
	}
}

// the sql statement returned from sqlite_master is the full view definition
// stripViewSql remove the `CREATE VIEW {:view_name} AS` And spaces Prefix from the provided SQL query
func stripViewSql(s string) string {
	splits := strings.SplitAfterN(s, "AS", 2)
	if len(splits) == 1 {
		return strings.TrimSpace(splits[0])
	}
	return strings.TrimSpace(splits[1])
}
