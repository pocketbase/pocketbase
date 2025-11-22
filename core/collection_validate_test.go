package core_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestCollectionValidate(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name           string
		collection     func(app core.App) (*core.Collection, error)
		expectedErrors []string
	}{
		{
			name: "empty collection",
			collection: func(app core.App) (*core.Collection, error) {
				return &core.Collection{}, nil
			},
			expectedErrors: []string{
				"id", "name", "type", "fields", // no default fields because the type is unknown
			},
		},
		{
			name: "unknown type with all invalid fields",
			collection: func(app core.App) (*core.Collection, error) {
				c := &core.Collection{}
				c.Id = "invalid_id ?!@#$"
				c.Name = "invalid_name ?!@#$"
				c.Type = "invalid_type"
				c.ListRule = types.Pointer("missing = '123'")
				c.ViewRule = types.Pointer("missing = '123'")
				c.CreateRule = types.Pointer("missing = '123'")
				c.UpdateRule = types.Pointer("missing = '123'")
				c.DeleteRule = types.Pointer("missing = '123'")
				c.Indexes = []string{"create index '' on '' ()"}

				// type specific fields
				c.ViewQuery = "invalid"                       // should be ignored
				c.AuthRule = types.Pointer("missing = '123'") // should be ignored

				return c, nil
			},
			expectedErrors: []string{
				"id", "name", "type", "indexes",
				"listRule", "viewRule", "createRule", "updateRule", "deleteRule",
				"fields", // no default fields because the type is unknown
			},
		},
		{
			name: "base with invalid fields",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("invalid_name ?!@#$")
				c.Indexes = []string{"create index '' on '' ()"}

				// type specific fields
				c.ViewQuery = "invalid"                       // should be ignored
				c.AuthRule = types.Pointer("missing = '123'") // should be ignored

				return c, nil
			},
			expectedErrors: []string{"name", "indexes"},
		},
		{
			name: "view with invalid fields",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewViewCollection("invalid_name ?!@#$")
				c.Indexes = []string{"create index '' on '' ()"}

				// type specific fields
				c.ViewQuery = "invalid"
				c.AuthRule = types.Pointer("missing = '123'") // should be ignored

				return c, nil
			},
			expectedErrors: []string{"indexes", "name", "fields", "viewQuery"},
		},
		{
			name: "auth with invalid fields",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("invalid_name ?!@#$")
				c.Indexes = []string{"create index '' on '' ()"}

				// type specific fields
				c.ViewQuery = "invalid" // should be ignored
				c.AuthRule = types.Pointer("missing = '123'")

				return c, nil
			},
			expectedErrors: []string{"indexes", "name", "authRule"},
		},

		// type checks
		{
			name: "empty type",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("test")
				c.Type = ""
				return c, nil
			},
			expectedErrors: []string{"type"},
		},
		{
			name: "unknown type",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("test")
				c.Type = "unknown"
				return c, nil
			},
			expectedErrors: []string{"type"},
		},
		{
			name: "base type",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("test")
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "view type",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewViewCollection("test")
				c.ViewQuery = "select 1 as id"
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "auth type",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("test")
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "changing type",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("users")
				c.Type = core.CollectionTypeBase
				return c, nil
			},
			expectedErrors: []string{"type"},
		},

		// system checks
		{
			name: "change from system to regular",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
				c.System = false
				return c, nil
			},
			expectedErrors: []string{"system"},
		},
		{
			name: "change from regular to system",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				c.System = true
				return c, nil
			},
			expectedErrors: []string{"system"},
		},
		{
			name: "create system",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("new_system")
				c.System = true
				return c, nil
			},
			expectedErrors: []string{},
		},

		// id checks
		{
			name: "empty id",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("test")
				c.Id = ""
				return c, nil
			},
			expectedErrors: []string{"id"},
		},
		{
			name: "invalid id",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("test")
				c.Id = "!invalid"
				return c, nil
			},
			expectedErrors: []string{"id"},
		},
		{
			name: "existing id",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("test")
				c.Id = "_pb_users_auth_"
				return c, nil
			},
			expectedErrors: []string{"id"},
		},
		{
			name: "changing id",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo3")
				c.Id = "anything"
				return c, nil
			},
			expectedErrors: []string{"id"},
		},
		{
			name: "valid id",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("test")
				c.Id = "anything"
				return c, nil
			},
			expectedErrors: []string{},
		},

		// name checks
		{
			name: "empty name",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("")
				c.Id = "test"
				return c, nil
			},
			expectedErrors: []string{"name"},
		},
		{
			name: "invalid name",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("!invalid")
				return c, nil
			},
			expectedErrors: []string{"name"},
		},
		{
			name: "name with _via_",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("a_via_b")
				return c, nil
			},
			expectedErrors: []string{"name"},
		},
		{
			name: "create with existing collection name",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("demo1")
				return c, nil
			},
			expectedErrors: []string{"name"},
		},
		{
			name: "create with existing internal table name",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("_collections")
				return c, nil
			},
			expectedErrors: []string{"name"},
		},
		{
			name: "update with existing collection name",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("users")
				c.Name = "demo1"
				return c, nil
			},
			expectedErrors: []string{"name"},
		},
		{
			name: "update with existing internal table name",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("users")
				c.Name = "_collections"
				return c, nil
			},
			expectedErrors: []string{"name"},
		},
		{
			name: "system collection name change",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
				c.Name = "superusers_new"
				return c, nil
			},
			expectedErrors: []string{"name"},
		},
		{
			name: "create with valid name",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("new_col")
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "update with valid name",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				c.Name = "demo1_new"
				return c, nil
			},
			expectedErrors: []string{},
		},

		// rule checks
		{
			name: "invalid base rules",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("new")
				c.ListRule = types.Pointer("!invalid")
				c.ViewRule = types.Pointer("missing = 123")
				c.CreateRule = types.Pointer("id = 123 && missing = 456")
				c.UpdateRule = types.Pointer("@request.body.missing:changed = false")
				c.DeleteRule = types.Pointer("(id=123")
				return c, nil
			},
			expectedErrors: []string{"listRule", "viewRule", "createRule", "updateRule", "deleteRule"},
		},
		{
			name: "valid base rules",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("new")
				c.Fields.Add(&core.TextField{Name: "f1"}) // dummy field to ensure that new fields can be referenced
				c.ListRule = types.Pointer("")
				c.ViewRule = types.Pointer("f1 = 123")
				c.CreateRule = types.Pointer("id = 123 && f1 = 456")
				c.UpdateRule = types.Pointer("(id = 123 && @request.body.f1:changed = false)")
				c.DeleteRule = types.Pointer("f1 = 123")
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "view with non-nil create/update/delete rules",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewViewCollection("new")
				c.ViewQuery = "select 1 as id, 'text' as f1"
				c.ListRule = types.Pointer("id = 123")
				c.ViewRule = types.Pointer("f1 = 456")
				c.CreateRule = types.Pointer("")
				c.UpdateRule = types.Pointer("")
				c.DeleteRule = types.Pointer("")
				return c, nil
			},
			expectedErrors: []string{"createRule", "updateRule", "deleteRule"},
		},
		{
			name: "view with nil create/update/delete rules",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewViewCollection("new")
				c.ViewQuery = "select 1 as id, 'text' as f1"
				c.ListRule = types.Pointer("id = 1")
				c.ViewRule = types.Pointer("f1 = 456")
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "changing api rules",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("users")
				c.Fields.Add(&core.TextField{Name: "f1"}) // dummy field to ensure that new fields can be referenced
				c.ListRule = types.Pointer("id = 1")
				c.ViewRule = types.Pointer("f1 = 456")
				c.CreateRule = types.Pointer("id = 123 && f1 = 456")
				c.UpdateRule = types.Pointer("(id = 123)")
				c.DeleteRule = types.Pointer("f1 = 123")
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "changing system collection api rules",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
				c.ListRule = types.Pointer("1 = 1")
				c.ViewRule = types.Pointer("1 = 1")
				c.CreateRule = types.Pointer("1 = 1")
				c.UpdateRule = types.Pointer("1 = 1")
				c.DeleteRule = types.Pointer("1 = 1")
				c.ManageRule = types.Pointer("1 = 1")
				c.AuthRule = types.Pointer("1 = 1")
				return c, nil
			},
			expectedErrors: []string{
				"listRule", "viewRule", "createRule", "updateRule",
				"deleteRule", "manageRule", "authRule",
			},
		},

		// indexes checks
		{
			name: "invalid index expression",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				c.Indexes = []string{
					"create index invalid",
					"create index idx_test_demo2 on anything (text)", // the name of table shouldn't matter
				}
				return c, nil
			},
			expectedErrors: []string{"indexes"},
		},
		{
			name: "index name used in other table",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				c.Indexes = []string{
					"create index `idx_test_demo1` on demo1 (id)",
					"create index `__pb_USERS_auth__username_idx` on anything (text)", // should be case-insensitive
				}
				return c, nil
			},
			expectedErrors: []string{"indexes"},
		},
		{
			name: "duplicated index names",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				c.Indexes = []string{
					"create index idx_test_demo1 on demo1 (id)",
					"create index idx_test_demo1 on anything (text)",
				}
				return c, nil
			},
			expectedErrors: []string{"indexes"},
		},
		{
			name: "duplicated index definitions",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				c.Indexes = []string{
					"create index idx_test_demo1 on demo1 (id)",
					"create index idx_test_demo2 on demo1 (id)",
				}
				return c, nil
			},
			expectedErrors: []string{"indexes"},
		},
		{
			name: "try to add index to a view collection",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("view1")
				c.Indexes = []string{"create index idx_test_view1 on view1 (id)"}
				return c, nil
			},
			expectedErrors: []string{"indexes"},
		},
		{
			name: "replace old with new indexes",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				c.Indexes = []string{
					"create index idx_test_demo1 on demo1 (id)",
					"create index idx_test_demo2 on anything (text)", // the name of table shouldn't matter
				}
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "old + new indexes",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				c.Indexes = []string{
					"CREATE INDEX `_wsmn24bux7wo113_created_idx` ON `demo1` (`created`)",
					"create index idx_test_demo1 on anything (id)",
				}
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "index for missing field",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				c.Indexes = []string{
					"create index idx_test_demo1 on anything (missing)", // still valid because it is checked on db persist
				}
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "auth collection with missing required unique indexes",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Indexes = []string{}
				return c, nil
			},
			expectedErrors: []string{"indexes", "passwordAuth"},
		},
		{
			name: "auth collection with non-unique required indexes",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Indexes = []string{
					"create index test_idx1 on new_auth (tokenKey)",
					"create index test_idx2 on new_auth (email)",
				}
				return c, nil
			},
			expectedErrors: []string{"indexes", "passwordAuth"},
		},
		{
			name: "auth collection with unique required indexes",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Indexes = []string{
					"create unique index test_idx1 on new_auth (tokenKey)",
					"create unique index test_idx2 on new_auth (email)",
				}
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "removing index on system field",
			collection: func(app core.App) (*core.Collection, error) {
				demo2, err := app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// mark the title field as system
				demo2.Fields.GetByName("title").SetSystem(true)
				if err = app.Save(demo2); err != nil {
					return nil, err
				}

				// refresh
				demo2, err = app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				demo2.RemoveIndex("idx_unique_demo2_title")

				return demo2, nil
			},
			expectedErrors: []string{"indexes"},
		},
		{
			name: "changing partial constraint of existing index on system field",
			collection: func(app core.App) (*core.Collection, error) {
				demo2, err := app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// mark the title field as system
				demo2.Fields.GetByName("title").SetSystem(true)
				if err = app.Save(demo2); err != nil {
					return nil, err
				}

				// refresh
				demo2, err = app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// replace the index with a partial one
				demo2.RemoveIndex("idx_unique_demo2_title")
				demo2.AddIndex("idx_new_demo2_title", true, "title", "1 = 1")

				return demo2, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "changing column sort and collate of existing index on system field",
			collection: func(app core.App) (*core.Collection, error) {
				demo2, err := app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// mark the title field as system
				demo2.Fields.GetByName("title").SetSystem(true)
				if err = app.Save(demo2); err != nil {
					return nil, err
				}

				// refresh
				demo2, err = app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// replace the index with a new one for the same column but with collate and sort
				demo2.RemoveIndex("idx_unique_demo2_title")
				demo2.AddIndex("idx_new_demo2_title", true, "title COLLATE test ASC", "")

				return demo2, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "adding new column to index on system field",
			collection: func(app core.App) (*core.Collection, error) {
				demo2, err := app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// mark the title field as system
				demo2.Fields.GetByName("title").SetSystem(true)
				if err = app.Save(demo2); err != nil {
					return nil, err
				}

				// refresh
				demo2, err = app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// replace the index with a non-unique one
				demo2.RemoveIndex("idx_unique_demo2_title")
				demo2.AddIndex("idx_new_title", false, "title, id", "")

				return demo2, nil
			},
			expectedErrors: []string{"indexes"},
		},
		{
			name: "changing index type on system field",
			collection: func(app core.App) (*core.Collection, error) {
				demo2, err := app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// mark the title field as system
				demo2.Fields.GetByName("title").SetSystem(true)
				if err = app.Save(demo2); err != nil {
					return nil, err
				}

				// refresh
				demo2, err = app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// replace the index with a non-unique one (partial constraints are ignored)
				demo2.RemoveIndex("idx_unique_demo2_title")
				demo2.AddIndex("idx_new_title", false, "title", "1=1")

				return demo2, nil
			},
			expectedErrors: []string{"indexes"},
		},
		{
			name: "changing index on non-system field",
			collection: func(app core.App) (*core.Collection, error) {
				demo2, err := app.FindCollectionByNameOrId("demo2")
				if err != nil {
					return nil, err
				}

				// replace the index with a partial one
				demo2.RemoveIndex("idx_demo2_active")
				demo2.AddIndex("idx_demo2_active", true, "active", "1 = 1")

				return demo2, nil
			},
			expectedErrors: []string{},
		},

		// fields list checks
		{
			name: "empty fields",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("new_auth")
				c.Fields = nil // the minimum fields should auto added
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "no id primay key field",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("new_auth")
				c.Fields = core.NewFieldsList(
					&core.TextField{Name: "id"},
				)
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "with id primay key field",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("new_auth")
				c.Fields = core.NewFieldsList(
					&core.TextField{Name: "id", PrimaryKey: true, Required: true, Pattern: `\w+`},
				)
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "duplicated field names",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("new_auth")
				c.Fields = core.NewFieldsList(
					&core.TextField{Name: "id", PrimaryKey: true, Required: true, Pattern: `\w+`},
					&core.TextField{Id: "f1", Name: "Test"}, // case-insensitive
					&core.BoolField{Id: "f2", Name: "test"},
				)
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "changing field type",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo1")
				f := c.Fields.GetByName("text")
				c.Fields.Add(&core.BoolField{Id: f.GetId(), Name: f.GetName()})
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "renaming system field",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId(core.CollectionNameAuthOrigins)
				f := c.Fields.GetByName("fingerprint")
				f.SetName("fingerprint_new")
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "deleting system field",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId(core.CollectionNameAuthOrigins)
				c.Fields.RemoveByName("fingerprint")
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "invalid field setting",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("test_new")
				c.Fields.Add(&core.TextField{Name: "f1", Min: -10})
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "valid field setting",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("test_new")
				c.Fields.Add(&core.TextField{Name: "f1", Min: 10})
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "fields view changes should be ignored",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("view1")
				c.Fields = nil
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "with reserved auth only field name (passwordConfirm)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.TextField{Name: "passwordConfirm"},
				)
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "with reserved auth only field name (oldPassword)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.TextField{Name: "oldPassword"},
				)
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "with invalid password auth field options (1)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.TextField{Name: "password", System: true, Hidden: true}, // should be PasswordField
				)
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "with valid password auth field options (2)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.PasswordField{Name: "password", System: true, Hidden: true},
				)
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "with invalid tokenKey auth field options (1)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.TextField{Name: "tokenKey", System: true}, // should be also hidden
				)
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "with valid tokenKey auth field options (2)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.TextField{Name: "tokenKey", System: true, Hidden: true},
				)
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "with invalid email auth field options (1)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.TextField{Name: "email", System: true}, // should be EmailField
				)
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "with valid email auth field options (2)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.EmailField{Name: "email", System: true},
				)
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "with invalid verified auth field options (1)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.TextField{Name: "verified", System: true}, // should be BoolField
				)
				return c, nil
			},
			expectedErrors: []string{"fields"},
		},
		{
			name: "with valid verified auth field options (2)",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(
					&core.BoolField{Name: "verified", System: true},
				)
				return c, nil
			},
			expectedErrors: []string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			collection, err := s.collection(app)
			if err != nil {
				t.Fatalf("Failed to retrieve test collection: %v", err)
			}

			result := app.Validate(collection)

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}
