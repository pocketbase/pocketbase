package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestGenerateDefaultRandomId(t *testing.T) {
	t.Parallel()

	id1 := core.GenerateDefaultRandomId()
	id2 := core.GenerateDefaultRandomId()

	if id1 == id2 {
		t.Fatalf("Expected id1 and id2 to differ, got %q", id1)
	}

	if l := len(id1); l != 15 {
		t.Fatalf("Expected id1 length %d, got %d", 15, l)
	}

	if l := len(id2); l != 15 {
		t.Fatalf("Expected id2 length %d, got %d", 15, l)
	}
}

func TestModelQuery(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	modelsQuery := app.ModelQuery(&core.Collection{})
	logsModelQuery := app.AuxModelQuery(&core.Collection{})

	if app.DB() == modelsQuery.Info().Builder {
		t.Fatalf("ModelQuery() is not using app.DB()")
	}

	if app.AuxDB() == logsModelQuery.Info().Builder {
		t.Fatalf("AuxModelQuery() is not using app.AuxDB()")
	}

	expectedSQL := "SELECT {{_collections}}.* FROM `_collections`"
	for i, q := range []*dbx.SelectQuery{modelsQuery, logsModelQuery} {
		sql := q.Build().SQL()
		if sql != expectedSQL {
			t.Fatalf("[%d] Expected select\n%s\ngot\n%s", i, expectedSQL, sql)
		}
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	u := &mockSuperusers{}

	testErr := errors.New("test")

	app.OnModelValidate().BindFunc(func(e *core.ModelEvent) error {
		return testErr
	})

	err := app.Validate(u)
	if err != testErr {
		t.Fatalf("Expected error %v, got %v", testErr, err)
	}
}

func TestValidateWithContext(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	u := &mockSuperusers{}

	testErr := errors.New("test")

	app.OnModelValidate().BindFunc(func(e *core.ModelEvent) error {
		if v := e.Context.Value("test"); v != 123 {
			t.Fatalf("Expected 'test' context value %#v, got %#v", 123, v)
		}
		return testErr
	})

	//nolint:staticcheck
	ctx := context.WithValue(context.Background(), "test", 123)

	err := app.ValidateWithContext(ctx, u)
	if err != testErr {
		t.Fatalf("Expected error %v, got %v", testErr, err)
	}
}

// -------------------------------------------------------------------

type mockSuperusers struct {
	core.BaseModel
	Email string `db:"email"`
}

func (m *mockSuperusers) TableName() string {
	return core.CollectionNameSuperusers
}
