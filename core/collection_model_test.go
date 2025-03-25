package core_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestNewCollection(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		typ      string
		name     string
		expected []string
	}{
		{
			"",
			"",
			[]string{
				`"id":"pbc_`,
				`"name":""`,
				`"type":"base"`,
				`"system":false`,
				`"indexes":[]`,
				`"fields":[{`,
				`"name":"id"`,
				`"type":"text"`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
			},
		},
		{
			"unknown",
			"test",
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"base"`,
				`"system":false`,
				`"indexes":[]`,
				`"fields":[{`,
				`"name":"id"`,
				`"type":"text"`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
			},
		},
		{
			"base",
			"test",
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"base"`,
				`"system":false`,
				`"indexes":[]`,
				`"fields":[{`,
				`"name":"id"`,
				`"type":"text"`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
			},
		},
		{
			"view",
			"test",
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"view"`,
				`"indexes":[]`,
				`"fields":[]`,
				`"system":false`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
			},
		},
		{
			"auth",
			"test",
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"auth"`,
				`"fields":[{`,
				`"system":false`,
				`"type":"text"`,
				`"type":"email"`,
				`"name":"id"`,
				`"name":"email"`,
				`"name":"password"`,
				`"name":"tokenKey"`,
				`"name":"emailVisibility"`,
				`"name":"verified"`,
				`idx_email`,
				`idx_tokenKey`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
				`"identityFields":["email"]`,
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s_%s", i, s.typ, s.name), func(t *testing.T) {
			result := core.NewCollection(s.typ, s.name).String()

			for _, part := range s.expected {
				if !strings.Contains(result, part) {
					t.Fatalf("Missing part %q in\n%v", part, result)
				}
			}
		})
	}
}

func TestNewBaseCollection(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name     string
		expected []string
	}{
		{
			"",
			[]string{
				`"id":"pbc_`,
				`"name":""`,
				`"type":"base"`,
				`"system":false`,
				`"indexes":[]`,
				`"fields":[{`,
				`"name":"id"`,
				`"type":"text"`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
			},
		},
		{
			"test",
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"base"`,
				`"system":false`,
				`"indexes":[]`,
				`"fields":[{`,
				`"name":"id"`,
				`"type":"text"`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.name), func(t *testing.T) {
			result := core.NewBaseCollection(s.name).String()

			for _, part := range s.expected {
				if !strings.Contains(result, part) {
					t.Fatalf("Missing part %q in\n%v", part, result)
				}
			}
		})
	}
}

func TestNewViewCollection(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name     string
		expected []string
	}{
		{
			"",
			[]string{
				`"id":"pbc_`,
				`"name":""`,
				`"type":"view"`,
				`"indexes":[]`,
				`"fields":[]`,
				`"system":false`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
			},
		},
		{
			"test",
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"view"`,
				`"indexes":[]`,
				`"fields":[]`,
				`"system":false`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.name), func(t *testing.T) {
			result := core.NewViewCollection(s.name).String()

			for _, part := range s.expected {
				if !strings.Contains(result, part) {
					t.Fatalf("Missing part %q in\n%v", part, result)
				}
			}
		})
	}
}

func TestNewAuthCollection(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name     string
		expected []string
	}{
		{
			"",
			[]string{
				`"id":""`,
				`"name":""`,
				`"type":"auth"`,
				`"fields":[{`,
				`"system":false`,
				`"type":"text"`,
				`"type":"email"`,
				`"name":"id"`,
				`"name":"email"`,
				`"name":"password"`,
				`"name":"tokenKey"`,
				`"name":"emailVisibility"`,
				`"name":"verified"`,
				`idx_email`,
				`idx_tokenKey`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
				`"identityFields":["email"]`,
			},
		},
		{
			"test",
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"auth"`,
				`"fields":[{`,
				`"system":false`,
				`"type":"text"`,
				`"type":"email"`,
				`"name":"id"`,
				`"name":"email"`,
				`"name":"password"`,
				`"name":"tokenKey"`,
				`"name":"emailVisibility"`,
				`"name":"verified"`,
				`idx_email`,
				`idx_tokenKey`,
				`"listRule":null`,
				`"viewRule":null`,
				`"createRule":null`,
				`"updateRule":null`,
				`"deleteRule":null`,
				`"identityFields":["email"]`,
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.name), func(t *testing.T) {
			result := core.NewAuthCollection(s.name).String()

			for _, part := range s.expected {
				if !strings.Contains(result, part) {
					t.Fatalf("Missing part %q in\n%v", part, result)
				}
			}
		})
	}
}

func TestCollectionTableName(t *testing.T) {
	t.Parallel()

	c := core.NewBaseCollection("test")
	if c.TableName() != "_collections" {
		t.Fatalf("Expected tableName %q, got %q", "_collections", c.TableName())
	}
}

func TestCollectionBaseFilesPath(t *testing.T) {
	t.Parallel()

	c := core.Collection{}

	if c.BaseFilesPath() != "" {
		t.Fatalf("Expected empty string, got %q", c.BaseFilesPath())
	}

	c.Id = "test"

	if c.BaseFilesPath() != c.Id {
		t.Fatalf("Expected %q, got %q", c.Id, c.BaseFilesPath())
	}
}

func TestCollectionIsBase(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		typ      string
		expected bool
	}{
		{"unknown", false},
		{core.CollectionTypeBase, true},
		{core.CollectionTypeView, false},
		{core.CollectionTypeAuth, false},
	}

	for _, s := range scenarios {
		t.Run(s.typ, func(t *testing.T) {
			c := core.Collection{}
			c.Type = s.typ

			if v := c.IsBase(); v != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, v)
			}
		})
	}
}

func TestCollectionIsView(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		typ      string
		expected bool
	}{
		{"unknown", false},
		{core.CollectionTypeBase, false},
		{core.CollectionTypeView, true},
		{core.CollectionTypeAuth, false},
	}

	for _, s := range scenarios {
		t.Run(s.typ, func(t *testing.T) {
			c := core.Collection{}
			c.Type = s.typ

			if v := c.IsView(); v != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, v)
			}
		})
	}
}

func TestCollectionIsAuth(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		typ      string
		expected bool
	}{
		{"unknown", false},
		{core.CollectionTypeBase, false},
		{core.CollectionTypeView, false},
		{core.CollectionTypeAuth, true},
	}

	for _, s := range scenarios {
		t.Run(s.typ, func(t *testing.T) {
			c := core.Collection{}
			c.Type = s.typ

			if v := c.IsAuth(); v != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, v)
			}
		})
	}
}

func TestCollectionPostScan(t *testing.T) {
	t.Parallel()

	rawOptions := types.JSONRaw(`{
		"viewQuery":"select 1",
		"authRule":"1=2"
	}`)

	scenarios := []struct {
		typ        string
		rawOptions types.JSONRaw
		expected   []string
	}{
		{
			core.CollectionTypeBase,
			rawOptions,
			[]string{
				`lastSavedPK:"test"`,
				`ViewQuery:""`,
				`AuthRule:(*string)(nil)`,
			},
		},
		{
			core.CollectionTypeView,
			rawOptions,
			[]string{
				`lastSavedPK:"test"`,
				`ViewQuery:"select 1"`,
				`AuthRule:(*string)(nil)`,
			},
		},
		{
			core.CollectionTypeAuth,
			rawOptions,
			[]string{
				`lastSavedPK:"test"`,
				`ViewQuery:""`,
				`AuthRule:(*string)(0x`,
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.typ), func(t *testing.T) {
			c := core.Collection{}
			c.Id = "test"
			c.Type = s.typ
			c.RawOptions = s.rawOptions

			err := c.PostScan()
			if err != nil {
				t.Fatal(err)
			}

			if c.IsNew() {
				t.Fatal("Expected the collection to be marked as not new")
			}

			rawModel := fmt.Sprintf("%#v", c)

			for _, part := range s.expected {
				if !strings.Contains(rawModel, part) {
					t.Fatalf("Missing part %q in\n%v", part, rawModel)
				}
			}
		})
	}
}

func TestCollectionUnmarshalJSON(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name        string
		raw         string
		collection  func() *core.Collection
		expected    []string
		notExpected []string
	}{
		{
			"base new empty",
			`{"type":"base","name":"test","listRule":"1=2","authRule":"1=3","viewQuery":"abc"}`,
			func() *core.Collection {
				return &core.Collection{}
			},
			[]string{
				`"type":"base"`,
				`"id":"pbc_`,
				`"name":"test"`,
				`"listRule":"1=2"`,
				`"fields":[`,
				`"name":"id"`,
				`"indexes":[]`,
			},
			[]string{
				`"authRule":"1=3"`,
				`"viewQuery":"abc"`,
			},
		},
		{
			"view new empty",
			`{"type":"view","name":"test","listRule":"1=2","authRule":"1=3","viewQuery":"abc"}`,
			func() *core.Collection {
				return &core.Collection{}
			},
			[]string{
				`"type":"view"`,
				`"id":"pbc_`,
				`"name":"test"`,
				`"listRule":"1=2"`,
				`"fields":[]`,
				`"viewQuery":"abc"`,
				`"indexes":[]`,
			},
			[]string{
				`"authRule":"1=3"`,
			},
		},
		{
			"auth new empty",
			`{"type":"auth","name":"test","listRule":"1=2","authRule":"1=3","viewQuery":"abc"}`,
			func() *core.Collection {
				return &core.Collection{}
			},
			[]string{
				`"type":"auth"`,
				`"id":"pbc_`,
				`"name":"test"`,
				`"listRule":"1=2"`,
				`"authRule":"1=3"`,
				`"fields":[`,
				`"name":"id"`,
			},
			[]string{
				`"indexes":[]`,
				`"viewQuery":"abc"`,
			},
		},
		{
			"new but with set type (no default fields load)",
			`{"type":"base","name":"test","listRule":"1=2","authRule":"1=3","viewQuery":"abc"}`,
			func() *core.Collection {
				c := &core.Collection{}
				c.Type = core.CollectionTypeBase
				return c
			},
			[]string{
				`"type":"base"`,
				`"id":""`,
				`"name":"test"`,
				`"listRule":"1=2"`,
				`"fields":[]`,
			},
			[]string{
				`"authRule":"1=3"`,
				`"viewQuery":"abc"`,
			},
		},
		{
			"existing (no default fields load)",
			`{"type":"auth","name":"test","listRule":"1=2","authRule":"1=3","viewQuery":"abc"}`,
			func() *core.Collection {
				c, _ := app.FindCollectionByNameOrId("demo1")
				return c
			},
			[]string{
				`"type":"auth"`,
				`"name":"test"`,
				`"listRule":"1=2"`,
				`"authRule":"1=3"`,
				`"fields":[`,
				`"name":"id"`,
			},
			[]string{
				`"name":"tokenKey"`,
				`"viewQuery":"abc"`,
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			collection := s.collection()

			err := json.Unmarshal([]byte(s.raw), collection)
			if err != nil {
				t.Fatal(err)
			}

			rawResult, err := json.Marshal(collection)
			if err != nil {
				t.Fatal(err)
			}
			rawResultStr := string(rawResult)

			for _, part := range s.expected {
				if !strings.Contains(rawResultStr, part) {
					t.Fatalf("Missing expected %q in\n%v", part, rawResultStr)
				}
			}

			for _, part := range s.notExpected {
				if strings.Contains(rawResultStr, part) {
					t.Fatalf("Didn't expected %q in\n%v", part, rawResultStr)
				}
			}
		})
	}
}

func TestCollectionSerialize(t *testing.T) {
	scenarios := []struct {
		name        string
		collection  func() *core.Collection
		expected    []string
		notExpected []string
	}{
		{
			"base",
			func() *core.Collection {
				c := core.NewCollection(core.CollectionTypeBase, "test")
				c.ViewQuery = "1=1"
				c.OAuth2.Providers = []core.OAuth2ProviderConfig{
					{Name: "test1", ClientId: "test_client_id1", ClientSecret: "test_client_secret1"},
					{Name: "test2", ClientId: "test_client_id2", ClientSecret: "test_client_secret2"},
				}

				return c
			},
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"base"`,
			},
			[]string{
				"verificationTemplate",
				"manageRule",
				"authRule",
				"secret",
				"oauth2",
				"clientId",
				"clientSecret",
				"viewQuery",
			},
		},
		{
			"view",
			func() *core.Collection {
				c := core.NewCollection(core.CollectionTypeView, "test")
				c.ViewQuery = "1=1"
				c.OAuth2.Providers = []core.OAuth2ProviderConfig{
					{Name: "test1", ClientId: "test_client_id1", ClientSecret: "test_client_secret1"},
					{Name: "test2", ClientId: "test_client_id2", ClientSecret: "test_client_secret2"},
				}

				return c
			},
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"view"`,
				`"viewQuery":"1=1"`,
			},
			[]string{
				"verificationTemplate",
				"manageRule",
				"authRule",
				"secret",
				"oauth2",
				"clientId",
				"clientSecret",
			},
		},
		{
			"auth",
			func() *core.Collection {
				c := core.NewCollection(core.CollectionTypeAuth, "test")
				c.ViewQuery = "1=1"
				c.OAuth2.Providers = []core.OAuth2ProviderConfig{
					{Name: "test1", ClientId: "test_client_id1", ClientSecret: "test_client_secret1"},
					{Name: "test2", ClientId: "test_client_id2", ClientSecret: "test_client_secret2"},
				}

				return c
			},
			[]string{
				`"id":"pbc_`,
				`"name":"test"`,
				`"type":"auth"`,
				`"oauth2":{`,
				`"providers":[{`,
				`"clientId":"test_client_id1"`,
				`"clientId":"test_client_id2"`,
			},
			[]string{
				"viewQuery",
				"secret",
				"clientSecret",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			collection := s.collection()

			raw, err := collection.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			rawStr := string(raw)

			if rawStr != collection.String() {
				t.Fatalf("Expected the same serialization, got\n%v\nVS\n%v", collection.String(), rawStr)
			}

			for _, part := range s.expected {
				if !strings.Contains(rawStr, part) {
					t.Fatalf("Missing part %q in\n%v", part, rawStr)
				}
			}

			for _, part := range s.notExpected {
				if strings.Contains(rawStr, part) {
					t.Fatalf("Didn't expect part %q in\n%v", part, rawStr)
				}
			}
		})
	}
}

func TestCollectionDBExport(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	date, err := types.ParseDateTime("2024-07-01 01:02:03.456Z")
	if err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		typ      string
		expected string
	}{
		{
			"unknown",
			`{"createRule":"1=3","created":"2024-07-01 01:02:03.456Z","deleteRule":"1=5","fields":[{"hidden":false,"id":"f1_id","name":"f1","presentable":false,"required":false,"system":true,"type":"bool"},{"hidden":false,"id":"f2_id","name":"f2","presentable":false,"required":true,"system":false,"type":"bool"}],"id":"test_id","indexes":["CREATE INDEX idx1 on test_name(id)","CREATE INDEX idx2 on test_name(id)"],"listRule":"1=1","name":"test_name","options":"{}","system":true,"type":"unknown","updateRule":"1=4","updated":"2024-07-01 01:02:03.456Z","viewRule":"1=7"}`,
		},
		{
			core.CollectionTypeBase,
			`{"createRule":"1=3","created":"2024-07-01 01:02:03.456Z","deleteRule":"1=5","fields":[{"hidden":false,"id":"f1_id","name":"f1","presentable":false,"required":false,"system":true,"type":"bool"},{"hidden":false,"id":"f2_id","name":"f2","presentable":false,"required":true,"system":false,"type":"bool"}],"id":"test_id","indexes":["CREATE INDEX idx1 on test_name(id)","CREATE INDEX idx2 on test_name(id)"],"listRule":"1=1","name":"test_name","options":"{}","system":true,"type":"base","updateRule":"1=4","updated":"2024-07-01 01:02:03.456Z","viewRule":"1=7"}`,
		},
		{
			core.CollectionTypeView,
			`{"createRule":"1=3","created":"2024-07-01 01:02:03.456Z","deleteRule":"1=5","fields":[{"hidden":false,"id":"f1_id","name":"f1","presentable":false,"required":false,"system":true,"type":"bool"},{"hidden":false,"id":"f2_id","name":"f2","presentable":false,"required":true,"system":false,"type":"bool"}],"id":"test_id","indexes":["CREATE INDEX idx1 on test_name(id)","CREATE INDEX idx2 on test_name(id)"],"listRule":"1=1","name":"test_name","options":{"viewQuery":"select 1"},"system":true,"type":"view","updateRule":"1=4","updated":"2024-07-01 01:02:03.456Z","viewRule":"1=7"}`,
		},
		{
			core.CollectionTypeAuth,
			`{"createRule":"1=3","created":"2024-07-01 01:02:03.456Z","deleteRule":"1=5","fields":[{"hidden":false,"id":"f1_id","name":"f1","presentable":false,"required":false,"system":true,"type":"bool"},{"hidden":false,"id":"f2_id","name":"f2","presentable":false,"required":true,"system":false,"type":"bool"}],"id":"test_id","indexes":["CREATE INDEX idx1 on test_name(id)","CREATE INDEX idx2 on test_name(id)"],"listRule":"1=1","name":"test_name","options":{"authRule":null,"manageRule":"1=6","authAlert":{"enabled":false,"emailTemplate":{"subject":"","body":""}},"oauth2":{"providers":null,"mappedFields":{"id":"","name":"","username":"","avatarURL":""},"enabled":false},"passwordAuth":{"enabled":false,"identityFields":null},"mfa":{"enabled":false,"duration":0,"rule":""},"otp":{"enabled":false,"duration":0,"length":0,"emailTemplate":{"subject":"","body":""}},"authToken":{"duration":0},"passwordResetToken":{"duration":0},"emailChangeToken":{"duration":0},"verificationToken":{"duration":0},"fileToken":{"duration":0},"verificationTemplate":{"subject":"","body":""},"resetPasswordTemplate":{"subject":"","body":""},"confirmEmailChangeTemplate":{"subject":"","body":""}},"system":true,"type":"auth","updateRule":"1=4","updated":"2024-07-01 01:02:03.456Z","viewRule":"1=7"}`,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.typ), func(t *testing.T) {
			c := core.Collection{}
			c.Type = s.typ
			c.Id = "test_id"
			c.Name = "test_name"
			c.System = true
			c.ListRule = types.Pointer("1=1")
			c.ViewRule = types.Pointer("1=2")
			c.CreateRule = types.Pointer("1=3")
			c.UpdateRule = types.Pointer("1=4")
			c.DeleteRule = types.Pointer("1=5")
			c.ManageRule = types.Pointer("1=6")
			c.ViewRule = types.Pointer("1=7")
			c.Created = date
			c.Updated = date
			c.Indexes = types.JSONArray[string]{"CREATE INDEX idx1 on test_name(id)", "CREATE INDEX idx2 on test_name(id)"}
			c.ViewQuery = "select 1"
			c.Fields.Add(&core.BoolField{Id: "f1_id", Name: "f1", System: true})
			c.Fields.Add(&core.BoolField{Id: "f2_id", Name: "f2", Required: true})
			c.RawOptions = types.JSONRaw(`{"viewQuery": "select 2"}`) // should be ignored

			result, err := c.DBExport(app)
			if err != nil {
				t.Fatal(err)
			}

			raw, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}

			if str := string(raw); str != s.expected {
				t.Fatalf("Expected\n%v\ngot\n%v", s.expected, str)
			}
		})
	}
}

func TestCollectionIndexHelpers(t *testing.T) {
	t.Parallel()

	checkIndexes := func(t *testing.T, indexes, expectedIndexes []string) {
		if len(indexes) != len(expectedIndexes) {
			t.Fatalf("Expected %d indexes, got %d\n%v", len(expectedIndexes), len(indexes), indexes)
		}

		for _, idx := range expectedIndexes {
			if !slices.Contains(indexes, idx) {
				t.Fatalf("Missing index\n%v\nin\n%v", idx, indexes)
			}
		}
	}

	c := core.NewBaseCollection("test")
	checkIndexes(t, c.Indexes, nil)

	c.AddIndex("idx1", false, "colA,colB", "colA != 1")
	c.AddIndex("idx2", true, "colA", "")
	c.AddIndex("idx3", false, "colA", "")
	c.AddIndex("idx3", false, "colB", "") // should overwrite the previous one

	idx1 := "CREATE INDEX `idx1` ON `test` (colA,colB) WHERE colA != 1"
	idx2 := "CREATE UNIQUE INDEX `idx2` ON `test` (colA)"
	idx3 := "CREATE INDEX `idx3` ON `test` (colB)"

	checkIndexes(t, c.Indexes, []string{idx1, idx2, idx3})

	c.RemoveIndex("iDx2")    // case-insensitive
	c.RemoveIndex("missing") // noop

	checkIndexes(t, c.Indexes, []string{idx1, idx3})

	expectedIndexes := map[string]string{
		"missing": "",
		"idx1":    idx1,
		// the name is case insensitive
		"iDX3": idx3,
	}
	for key, expectedIdx := range expectedIndexes {
		idx := c.GetIndex(key)
		if idx != expectedIdx {
			t.Errorf("Expected index %q to be\n%v\ngot\n%v", key, expectedIdx, idx)
		}
	}
}

// -------------------------------------------------------------------

func TestCollectionDelete(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name                   string
		collection             string
		disableIntegrityChecks bool
		expectError            bool
	}{
		{
			name:        "unsaved",
			collection:  "",
			expectError: true,
		},
		{
			name:        "system",
			collection:  core.CollectionNameSuperusers,
			expectError: true,
		},
		{
			name:        "base with references",
			collection:  "demo1",
			expectError: true,
		},
		{
			name:                   "base with references with disabled integrity checks",
			collection:             "demo1",
			disableIntegrityChecks: true,
			expectError:            false,
		},
		{
			name:        "base without references",
			collection:  "demo1",
			expectError: true,
		},
		{
			name:        "view with reference",
			collection:  "view1",
			expectError: true,
		},
		{
			name:                   "view with references with disabled integrity checks",
			collection:             "view1",
			disableIntegrityChecks: true,
			expectError:            false,
		},
		{
			name:                   "view without references",
			collection:             "view2",
			disableIntegrityChecks: true,
			expectError:            false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			var col *core.Collection

			if s.collection == "" {
				col = core.NewBaseCollection("test")
			} else {
				var err error
				col, err = app.FindCollectionByNameOrId(s.collection)
				if err != nil {
					t.Fatal(err)
				}
			}

			if s.disableIntegrityChecks {
				col.IntegrityChecks(!s.disableIntegrityChecks)
			}

			err := app.Delete(col)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			exists := app.HasTable(col.Name)

			if !col.IsNew() && exists != hasErr {
				t.Fatalf("Expected HasTable %v, got %v", hasErr, exists)
			}

			if !hasErr {
				cache, _ := app.FindCachedCollectionByNameOrId(col.Id)
				if cache != nil {
					t.Fatal("Expected the collection to be removed from the cache.")
				}
			}
		})
	}
}

func TestCollectionModelEventSync(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	testCollections := make([]*core.Collection, 4)
	for i := 0; i < 4; i++ {
		testCollections[i] = core.NewBaseCollection("sync_test_" + strconv.Itoa(i))
		if err := app.Save(testCollections[i]); err != nil {
			t.Fatal(err)
		}
	}

	createModelEvent := func() *core.ModelEvent {
		event := new(core.ModelEvent)
		event.App = app
		event.Context = context.Background()
		event.Type = "test_a"
		event.Model = testCollections[0]
		return event
	}

	createModelErrorEvent := func() *core.ModelErrorEvent {
		event := new(core.ModelErrorEvent)
		event.ModelEvent = *createModelEvent()
		event.Error = errors.New("error_a")
		return event
	}

	changeCollectionEventBefore := func(e *core.CollectionEvent) {
		e.Type = "test_b"
		//nolint:staticcheck
		e.Context = context.WithValue(context.Background(), "test", 123)
		e.Collection = testCollections[1]
	}

	modelEventFinalizerChange := func(e *core.ModelEvent) {
		e.Type = "test_c"
		//nolint:staticcheck
		e.Context = context.WithValue(context.Background(), "test", 456)
		e.Model = testCollections[2]
	}

	changeCollectionEventAfter := func(e *core.CollectionEvent) {
		e.Type = "test_d"
		//nolint:staticcheck
		e.Context = context.WithValue(context.Background(), "test", 789)
		e.Collection = testCollections[3]
	}

	expectedBeforeModelEventHandlerChecks := func(t *testing.T, e *core.ModelEvent) {
		if e.Type != "test_a" {
			t.Fatalf("Expected type %q, got %q", "test_a", e.Type)
		}

		if v := e.Context.Value("test"); v != nil {
			t.Fatalf("Expected context value %v, got %v", nil, v)
		}

		if e.Model.PK() != testCollections[0].Id {
			t.Fatalf("Expected collection with id %q, got %q (%d)", testCollections[0].Id, e.Model.PK(), 0)
		}
	}

	expectedAfterModelEventHandlerChecks := func(t *testing.T, e *core.ModelEvent) {
		if e.Type != "test_d" {
			t.Fatalf("Expected type %q, got %q", "test_d", e.Type)
		}

		if v := e.Context.Value("test"); v != 789 {
			t.Fatalf("Expected context value %v, got %v", 789, v)
		}

		if e.Model.PK() != testCollections[3].Id {
			t.Fatalf("Expected collection with id %q, got %q (%d)", testCollections[3].Id, e.Model.PK(), 3)
		}
	}

	expectedBeforeCollectionEventHandlerChecks := func(t *testing.T, e *core.CollectionEvent) {
		if e.Type != "test_a" {
			t.Fatalf("Expected type %q, got %q", "test_a", e.Type)
		}

		if v := e.Context.Value("test"); v != nil {
			t.Fatalf("Expected context value %v, got %v", nil, v)
		}

		if e.Collection.Id != testCollections[0].Id {
			t.Fatalf("Expected collection with id %q, got %q (%d)", testCollections[0].Id, e.Collection.Id, 0)
		}
	}

	expectedAfterCollectionEventHandlerChecks := func(t *testing.T, e *core.CollectionEvent) {
		if e.Type != "test_c" {
			t.Fatalf("Expected type %q, got %q", "test_c", e.Type)
		}

		if v := e.Context.Value("test"); v != 456 {
			t.Fatalf("Expected context value %v, got %v", 456, v)
		}

		if e.Collection.Id != testCollections[2].Id {
			t.Fatalf("Expected collection with id %q, got %q (%d)", testCollections[2].Id, e.Collection.Id, 2)
		}
	}

	modelEventFinalizer := func(e *core.ModelEvent) error {
		modelEventFinalizerChange(e)
		return nil
	}

	modelErrorEventFinalizer := func(e *core.ModelErrorEvent) error {
		modelEventFinalizerChange(&e.ModelEvent)
		e.Error = errors.New("error_c")
		return nil
	}

	modelEventHandler := &hook.Handler[*core.ModelEvent]{
		Priority: -999,
		Func: func(e *core.ModelEvent) error {
			t.Run("before model", func(t *testing.T) {
				expectedBeforeModelEventHandlerChecks(t, e)
			})

			_ = e.Next()

			t.Run("after model", func(t *testing.T) {
				expectedAfterModelEventHandlerChecks(t, e)
			})

			return nil
		},
	}

	modelErrorEventHandler := &hook.Handler[*core.ModelErrorEvent]{
		Priority: -999,
		Func: func(e *core.ModelErrorEvent) error {
			t.Run("before model error", func(t *testing.T) {
				expectedBeforeModelEventHandlerChecks(t, &e.ModelEvent)
				if v := e.Error.Error(); v != "error_a" {
					t.Fatalf("Expected error %q, got %q", "error_a", v)
				}
			})

			_ = e.Next()

			t.Run("after model error", func(t *testing.T) {
				expectedAfterModelEventHandlerChecks(t, &e.ModelEvent)
				if v := e.Error.Error(); v != "error_d" {
					t.Fatalf("Expected error %q, got %q", "error_d", v)
				}
			})

			return nil
		},
	}

	recordEventHandler := &hook.Handler[*core.CollectionEvent]{
		Priority: -999,
		Func: func(e *core.CollectionEvent) error {
			t.Run("before collection", func(t *testing.T) {
				expectedBeforeCollectionEventHandlerChecks(t, e)
			})

			changeCollectionEventBefore(e)

			_ = e.Next()

			t.Run("after collection", func(t *testing.T) {
				expectedAfterCollectionEventHandlerChecks(t, e)
			})

			changeCollectionEventAfter(e)

			return nil
		},
	}

	collectionErrorEventHandler := &hook.Handler[*core.CollectionErrorEvent]{
		Priority: -999,
		Func: func(e *core.CollectionErrorEvent) error {
			t.Run("before collection error", func(t *testing.T) {
				expectedBeforeCollectionEventHandlerChecks(t, &e.CollectionEvent)
				if v := e.Error.Error(); v != "error_a" {
					t.Fatalf("Expected error %q, got %q", "error_c", v)
				}
			})

			changeCollectionEventBefore(&e.CollectionEvent)
			e.Error = errors.New("error_b")

			_ = e.Next()

			t.Run("after collection error", func(t *testing.T) {
				expectedAfterCollectionEventHandlerChecks(t, &e.CollectionEvent)
				if v := e.Error.Error(); v != "error_c" {
					t.Fatalf("Expected error %q, got %q", "error_c", v)
				}
			})

			changeCollectionEventAfter(&e.CollectionEvent)
			e.Error = errors.New("error_d")

			return nil
		},
	}

	// OnModelValidate
	app.OnCollectionValidate().Bind(recordEventHandler)
	app.OnModelValidate().Bind(modelEventHandler)
	app.OnModelValidate().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelCreate
	app.OnCollectionCreate().Bind(recordEventHandler)
	app.OnModelCreate().Bind(modelEventHandler)
	app.OnModelCreate().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelCreateExecute
	app.OnCollectionCreateExecute().Bind(recordEventHandler)
	app.OnModelCreateExecute().Bind(modelEventHandler)
	app.OnModelCreateExecute().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterCreateSuccess
	app.OnCollectionAfterCreateSuccess().Bind(recordEventHandler)
	app.OnModelAfterCreateSuccess().Bind(modelEventHandler)
	app.OnModelAfterCreateSuccess().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterCreateError
	app.OnCollectionAfterCreateError().Bind(collectionErrorEventHandler)
	app.OnModelAfterCreateError().Bind(modelErrorEventHandler)
	app.OnModelAfterCreateError().Trigger(createModelErrorEvent(), modelErrorEventFinalizer)

	// OnModelUpdate
	app.OnCollectionUpdate().Bind(recordEventHandler)
	app.OnModelUpdate().Bind(modelEventHandler)
	app.OnModelUpdate().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelUpdateExecute
	app.OnCollectionUpdateExecute().Bind(recordEventHandler)
	app.OnModelUpdateExecute().Bind(modelEventHandler)
	app.OnModelUpdateExecute().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterUpdateSuccess
	app.OnCollectionAfterUpdateSuccess().Bind(recordEventHandler)
	app.OnModelAfterUpdateSuccess().Bind(modelEventHandler)
	app.OnModelAfterUpdateSuccess().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterUpdateError
	app.OnCollectionAfterUpdateError().Bind(collectionErrorEventHandler)
	app.OnModelAfterUpdateError().Bind(modelErrorEventHandler)
	app.OnModelAfterUpdateError().Trigger(createModelErrorEvent(), modelErrorEventFinalizer)

	// OnModelDelete
	app.OnCollectionDelete().Bind(recordEventHandler)
	app.OnModelDelete().Bind(modelEventHandler)
	app.OnModelDelete().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelDeleteExecute
	app.OnCollectionDeleteExecute().Bind(recordEventHandler)
	app.OnModelDeleteExecute().Bind(modelEventHandler)
	app.OnModelDeleteExecute().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterDeleteSuccess
	app.OnCollectionAfterDeleteSuccess().Bind(recordEventHandler)
	app.OnModelAfterDeleteSuccess().Bind(modelEventHandler)
	app.OnModelAfterDeleteSuccess().Trigger(createModelEvent(), modelEventFinalizer)

	// OnModelAfterDeleteError
	app.OnCollectionAfterDeleteError().Bind(collectionErrorEventHandler)
	app.OnModelAfterDeleteError().Bind(modelErrorEventHandler)
	app.OnModelAfterDeleteError().Trigger(createModelErrorEvent(), modelErrorEventFinalizer)
}

func TestCollectionSaveModel(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name          string
		collection    func(app core.App) (*core.Collection, error)
		expectError   bool
		expectColumns []string
	}{
		// trigger validators
		{
			name: "create - trigger validators",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("!invalid")
				c.Fields.Add(&core.TextField{Name: "example"})
				c.AddIndex("test_save_idx", false, "example", "")
				return c, nil
			},
			expectError: true,
		},
		{
			name: "update - trigger validators",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo5")
				c.Name = "demo1"
				c.Fields.Add(&core.TextField{Name: "example"})
				c.Fields.RemoveByName("file")
				c.AddIndex("test_save_idx", false, "example", "")
				return c, nil
			},
			expectError: true,
		},

		// create
		{
			name: "create base collection",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewBaseCollection("new")
				c.Type = ""                 // should be auto set to "base"
				c.Fields.RemoveByName("id") // ensure that the default fields will be loaded
				c.Fields.Add(&core.TextField{Name: "example"})
				c.AddIndex("test_save_idx", false, "example", "")
				return c, nil
			},
			expectError: false,
			expectColumns: []string{
				"id", "example",
			},
		},
		{
			name: "create auth collection",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new")
				c.Fields.RemoveByName("id")    // ensure that the default fields will be loaded
				c.Fields.RemoveByName("email") // ensure that the default fields will be loaded
				c.Fields.Add(&core.TextField{Name: "example"})
				c.AddIndex("test_save_idx", false, "example", "")
				return c, nil
			},
			expectError: false,
			expectColumns: []string{
				"id", "email", "tokenKey", "password",
				"verified", "emailVisibility", "example",
			},
		},
		{
			name: "create view collection",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewViewCollection("new")
				c.Fields.Add(&core.TextField{Name: "ignored"}) // should be ignored
				c.ViewQuery = "select 1 as id, 2 as example"
				return c, nil
			},
			expectError: false,
			expectColumns: []string{
				"id", "example",
			},
		},

		// update
		{
			name: "update base collection",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("demo5")
				c.Fields.Add(&core.TextField{Name: "example"})
				c.Fields.RemoveByName("file")
				c.Fields.GetByName("total").SetName("total_updated")
				c.AddIndex("test_save_idx", false, "example", "")
				return c, nil
			},
			expectError: false,
			expectColumns: []string{
				"id", "select_one", "select_many", "rel_one", "rel_many",
				"total_updated", "created", "updated", "example",
			},
		},
		{
			name: "update auth collection",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("clients")
				c.Fields.Add(&core.TextField{Name: "example"})
				c.Fields.RemoveByName("file")
				c.Fields.GetByName("name").SetName("name_updated")
				c.AddIndex("test_save_idx", false, "example", "")
				return c, nil
			},
			expectError: false,
			expectColumns: []string{
				"id", "email", "emailVisibility", "password", "tokenKey",
				"verified", "username", "name_updated", "created", "updated", "example",
			},
		},
		{
			name: "update view collection",
			collection: func(app core.App) (*core.Collection, error) {
				c, _ := app.FindCollectionByNameOrId("view2")
				c.Fields.Add(&core.TextField{Name: "example"}) // should be ignored
				c.ViewQuery = "select 1 as id, 2 as example"
				return c, nil
			},
			expectError: false,
			expectColumns: []string{
				"id", "example",
			},
		},

		// auth normalization
		{
			name: "unset missing oauth2 mapped fields",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new")
				c.OAuth2.Enabled = true
				// shouldn't fail
				c.OAuth2.MappedFields = core.OAuth2KnownFields{
					Id:        "missing",
					Name:      "missing",
					Username:  "missing",
					AvatarURL: "missing",
				}
				return c, nil
			},
			expectError: false,
			expectColumns: []string{
				"id", "email", "emailVisibility", "password", "tokenKey", "verified",
			},
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

			saveErr := app.Save(collection)

			hasErr := saveErr != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", hasErr, s.expectError, saveErr)
			}

			if hasErr {
				return
			}

			// the collection should always have an id after successful Save
			if collection.Id == "" {
				t.Fatal("Expected collection id to be set")
			}

			// the timestamp fields should be non-empty after successful Save
			if collection.Created.String() == "" {
				t.Fatal("Expected collection created to be set")
			}
			if collection.Updated.String() == "" {
				t.Fatal("Expected collection updated to be set")
			}

			// check if the records table was synced
			hasTable := app.HasTable(collection.Name)
			if !hasTable {
				t.Fatalf("Expected records table %s to be created", collection.Name)
			}

			// check if the records table has the fields fields
			columns, err := app.TableColumns(collection.Name)
			if err != nil {
				t.Fatal(err)
			}
			if len(columns) != len(s.expectColumns) {
				t.Fatalf("Expected columns\n%v\ngot\n%v", s.expectColumns, columns)
			}
			for i, c := range columns {
				if !slices.Contains(s.expectColumns, c) {
					t.Fatalf("[%d] Didn't expect record column %q", i, c)
				}
			}

			// make sure that all collection indexes exists
			indexes, err := app.TableIndexes(collection.Name)
			if err != nil {
				t.Fatal(err)
			}
			if len(indexes) != len(collection.Indexes) {
				t.Fatalf("Expected %d indexes, got %d", len(collection.Indexes), len(indexes))
			}
			for _, idx := range collection.Indexes {
				parsed := dbutils.ParseIndex(idx)
				if _, ok := indexes[parsed.IndexName]; !ok {
					t.Fatalf("Missing index %q in\n%v", idx, indexes)
				}
			}
		})
	}
}

// indirect update of a field used in view should cause view(s) update
func TestCollectionSaveIndirectViewsUpdate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	// update MaxSelect fields
	{
		relMany := collection.Fields.GetByName("rel_many").(*core.RelationField)
		relMany.MaxSelect = 1

		fileOne := collection.Fields.GetByName("file_one").(*core.FileField)
		fileOne.MaxSelect = 10

		if err := app.Save(collection); err != nil {
			t.Fatal(err)
		}
	}

	// check view1 fields
	{
		view1, err := app.FindCollectionByNameOrId("view1")
		if err != nil {
			t.Fatal(err)
		}

		relMany := view1.Fields.GetByName("rel_many").(*core.RelationField)
		if relMany.MaxSelect != 1 {
			t.Fatalf("Expected view1.rel_many MaxSelect to be %d, got %v", 1, relMany.MaxSelect)
		}

		fileOne := view1.Fields.GetByName("file_one").(*core.FileField)
		if fileOne.MaxSelect != 10 {
			t.Fatalf("Expected view1.file_one MaxSelect to be %d, got %v", 10, fileOne.MaxSelect)
		}
	}

	// check view2 fields
	{
		view2, err := app.FindCollectionByNameOrId("view2")
		if err != nil {
			t.Fatal(err)
		}

		relMany := view2.Fields.GetByName("rel_many").(*core.RelationField)
		if relMany.MaxSelect != 1 {
			t.Fatalf("Expected view2.rel_many MaxSelect to be %d, got %v", 1, relMany.MaxSelect)
		}
	}
}

func TestCollectionSaveViewWrapping(t *testing.T) {
	t.Parallel()

	viewName := "test_wrapping"

	scenarios := []struct {
		name     string
		query    string
		expected string
	}{
		{
			"no wrapping - text field",
			"select text as id, bool from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (select text as id, bool from demo1)",
		},
		{
			"no wrapping - id field",
			"select text as id, bool from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (select text as id, bool from demo1)",
		},
		{
			"no wrapping - relation field",
			"select rel_one as id, bool from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (select rel_one as id, bool from demo1)",
		},
		{
			"no wrapping - select field",
			"select select_many as id, bool from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (select select_many as id, bool from demo1)",
		},
		{
			"no wrapping - email field",
			"select email as id, bool from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (select email as id, bool from demo1)",
		},
		{
			"no wrapping - datetime field",
			"select datetime as id, bool from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (select datetime as id, bool from demo1)",
		},
		{
			"no wrapping - url field",
			"select url as id, bool from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (select url as id, bool from demo1)",
		},
		{
			"wrapping - bool field",
			"select bool as id, text as txt, url from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (SELECT CAST(`id` as TEXT) `id`,`txt`,`url` FROM (select bool as id, text as txt, url from demo1))",
		},
		{
			"wrapping - bool field (different order)",
			"select text as txt, url, bool as id from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (SELECT `txt`,`url`,CAST(`id` as TEXT) `id` FROM (select text as txt, url, bool as id from demo1))",
		},
		{
			"wrapping - json field",
			"select json as id, text, url from demo1",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (SELECT CAST(`id` as TEXT) `id`,`text`,`url` FROM (select json as id, text, url from demo1))",
		},
		{
			"wrapping - numeric id",
			"select 1 as id",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (SELECT CAST(`id` as TEXT) `id` FROM (select 1 as id))",
		},
		{
			"wrapping - expresion",
			"select ('test') as id",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (SELECT CAST(`id` as TEXT) `id` FROM (select ('test') as id))",
		},
		{
			"no wrapping - cast as text",
			"select cast('test' as text) as id",
			"CREATE VIEW `test_wrapping` AS SELECT * FROM (select cast('test' as text) as id)",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			collection := core.NewViewCollection(viewName)
			collection.ViewQuery = s.query

			err := app.Save(collection)
			if err != nil {
				t.Fatal(err)
			}

			var sql string

			rowErr := app.DB().NewQuery("SELECT sql FROM sqlite_master WHERE type='view' AND name={:name}").
				Bind(dbx.Params{"name": viewName}).
				Row(&sql)
			if rowErr != nil {
				t.Fatalf("Failed to retrieve view sql: %v", rowErr)
			}

			if sql != s.expected {
				t.Fatalf("Expected query \n%v, \ngot \n%v", s.expected, sql)
			}
		})
	}
}
