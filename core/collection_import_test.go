package core_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestImportCollections(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	var regularCollections []*core.Collection
	err := testApp.CollectionQuery().AndWhere(dbx.HashExp{"system": false}).All(&regularCollections)
	if err != nil {
		t.Fatal(err)
	}

	var systemCollections []*core.Collection
	err = testApp.CollectionQuery().AndWhere(dbx.HashExp{"system": true}).All(&systemCollections)
	if err != nil {
		t.Fatal(err)
	}

	totalRegularCollections := len(regularCollections)
	totalSystemCollections := len(systemCollections)
	totalCollections := totalRegularCollections + totalSystemCollections

	scenarios := []struct {
		name                   string
		data                   []map[string]any
		deleteMissing          bool
		expectError            bool
		expectCollectionsCount int
		afterTestFunc          func(testApp *tests.TestApp, resultCollections []*core.Collection)
	}{
		{
			name:                   "empty collections",
			data:                   []map[string]any{},
			expectError:            true,
			expectCollectionsCount: totalCollections,
		},
		{
			name: "minimal collection import (with missing system fields)",
			data: []map[string]any{
				{"name": "import_test1", "type": "auth"},
				{
					"name": "import_test2", "fields": []map[string]any{
						{"name": "test", "type": "text"},
					},
				},
			},
			deleteMissing:          false,
			expectError:            false,
			expectCollectionsCount: totalCollections + 2,
		},
		{
			name: "minimal collection import (trigger collection model validations)",
			data: []map[string]any{
				{"name": ""},
				{
					"name": "import_test2", "fields": []map[string]any{
						{"name": "test", "type": "text"},
					},
				},
			},
			deleteMissing:          false,
			expectError:            true,
			expectCollectionsCount: totalCollections,
		},
		{
			name: "minimal collection import (trigger field settings validation)",
			data: []map[string]any{
				{"name": "import_test", "fields": []map[string]any{{"name": "test", "type": "text", "min": -1}}},
			},
			deleteMissing:          false,
			expectError:            true,
			expectCollectionsCount: totalCollections,
		},
		{
			name: "new + update + delete (system collections delete should be ignored)",
			data: []map[string]any{
				{
					"id":   "wsmn24bux7wo113",
					"name": "demo",
					"fields": []map[string]any{
						{
							"id":       "_2hlxbmp",
							"name":     "title",
							"type":     "text",
							"system":   false,
							"required": true,
							"min":      3,
							"max":      nil,
							"pattern":  "",
						},
					},
					"indexes": []string{},
				},
				{
					"name": "import1",
					"fields": []map[string]any{
						{
							"name": "active",
							"type": "bool",
						},
					},
				},
			},
			deleteMissing:          true,
			expectError:            false,
			expectCollectionsCount: totalSystemCollections + 2,
		},
		{
			name: "test with deleteMissing: false",
			data: []map[string]any{
				{
					// "id":   "wsmn24bux7wo113", // test update with only name as identifier
					"name": "demo1",
					"fields": []map[string]any{
						{
							"id":       "_2hlxbmp",
							"name":     "title",
							"type":     "text",
							"system":   false,
							"required": true,
							"min":      3,
							"max":      nil,
							"pattern":  "",
						},
						{
							"id":       "_2hlxbmp",
							"name":     "field_with_duplicate_id",
							"type":     "text",
							"system":   false,
							"required": true,
							"unique":   false,
							"min":      4,
							"max":      nil,
							"pattern":  "",
						},
						{
							"id":   "abcd_import",
							"name": "new_field",
							"type": "text",
						},
					},
				},
				{
					"name": "new_import",
					"fields": []map[string]any{
						{
							"id":   "abcd_import",
							"name": "active",
							"type": "bool",
						},
					},
				},
			},
			deleteMissing:          false,
			expectError:            false,
			expectCollectionsCount: totalCollections + 1,
			afterTestFunc: func(testApp *tests.TestApp, resultCollections []*core.Collection) {
				expectedCollectionFields := map[string]int{
					core.CollectionNameAuthOrigins: 6,
					"nologin":                      10,
					"demo1":                        19,
					"demo2":                        5,
					"demo3":                        5,
					"demo4":                        16,
					"demo5":                        9,
					"new_import":                   2,
				}
				for name, expectedCount := range expectedCollectionFields {
					collection, err := testApp.FindCollectionByNameOrId(name)
					if err != nil {
						t.Fatal(err)
					}

					if totalFields := len(collection.Fields); totalFields != expectedCount {
						t.Errorf("Expected %d %q fields, got %d", expectedCount, collection.Name, totalFields)
					}
				}
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testApp, _ := tests.NewTestApp()
			defer testApp.Cleanup()

			err := testApp.ImportCollections(s.data, s.deleteMissing)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr to be %v, got %v (%v)", s.expectError, hasErr, err)
			}

			// check collections count
			collections := []*core.Collection{}
			if err := testApp.CollectionQuery().All(&collections); err != nil {
				t.Fatal(err)
			}
			if len(collections) != s.expectCollectionsCount {
				t.Fatalf("Expected %d collections, got %d", s.expectCollectionsCount, len(collections))
			}

			if s.afterTestFunc != nil {
				s.afterTestFunc(testApp, collections)
			}
		})
	}
}

func TestImportCollectionsByMarshaledJSON(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	var regularCollections []*core.Collection
	err := testApp.CollectionQuery().AndWhere(dbx.HashExp{"system": false}).All(&regularCollections)
	if err != nil {
		t.Fatal(err)
	}

	var systemCollections []*core.Collection
	err = testApp.CollectionQuery().AndWhere(dbx.HashExp{"system": true}).All(&systemCollections)
	if err != nil {
		t.Fatal(err)
	}

	totalRegularCollections := len(regularCollections)
	totalSystemCollections := len(systemCollections)
	totalCollections := totalRegularCollections + totalSystemCollections

	scenarios := []struct {
		name                   string
		data                   string
		deleteMissing          bool
		expectError            bool
		expectCollectionsCount int
		afterTestFunc          func(testApp *tests.TestApp, resultCollections []*core.Collection)
	}{
		{
			name:                   "invalid json array",
			data:                   `{"test":123}`,
			expectError:            true,
			expectCollectionsCount: totalCollections,
		},
		{
			name: "new + update + delete (system collections delete should be ignored)",
			data: `[
				{
					"id":   "wsmn24bux7wo113",
					"name": "demo",
					"fields": [
						{
							"id":       "_2hlxbmp",
							"name":     "title",
							"type":     "text",
							"system":   false,
							"required": true,
							"min":      3,
							"max":      null,
							"pattern":  ""
						}
					],
					"indexes": []
				},
				{
					"name": "import1",
					"fields": [
						{
							"name": "active",
							"type": "bool"
						}
					]
				}
			]`,
			deleteMissing:          true,
			expectError:            false,
			expectCollectionsCount: totalSystemCollections + 2,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testApp, _ := tests.NewTestApp()
			defer testApp.Cleanup()

			err := testApp.ImportCollectionsByMarshaledJSON([]byte(s.data), s.deleteMissing)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr to be %v, got %v (%v)", s.expectError, hasErr, err)
			}

			// check collections count
			collections := []*core.Collection{}
			if err := testApp.CollectionQuery().All(&collections); err != nil {
				t.Fatal(err)
			}
			if len(collections) != s.expectCollectionsCount {
				t.Fatalf("Expected %d collections, got %d", s.expectCollectionsCount, len(collections))
			}

			if s.afterTestFunc != nil {
				s.afterTestFunc(testApp, collections)
			}
		})
	}
}

func TestImportCollectionsUpdateRules(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name          string
		data          map[string]any
		deleteMissing bool
	}{
		{
			"extend existing by name (without deleteMissing)",
			map[string]any{"name": "clients", "authToken": map[string]any{"duration": 100}, "fields": []map[string]any{{"name": "test", "type": "text"}}},
			false,
		},
		{
			"extend existing by id (without deleteMissing)",
			map[string]any{"id": "v851q4r790rhknl", "authToken": map[string]any{"duration": 100}, "fields": []map[string]any{{"name": "test", "type": "text"}}},
			false,
		},
		{
			"extend with delete missing",
			map[string]any{
				"id":           "v851q4r790rhknl",
				"authToken":    map[string]any{"duration": 100},
				"fields":       []map[string]any{{"name": "test", "type": "text"}},
				"passwordAuth": map[string]any{"identityFields": []string{"email"}},
				"indexes": []string{
					// min required system fields indexes
					"CREATE UNIQUE INDEX `_v851q4r790rhknl_email_idx` ON `clients` (email) WHERE email != ''",
					"CREATE UNIQUE INDEX `_v851q4r790rhknl_tokenKey_idx` ON `clients` (tokenKey)",
				},
			},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testApp, _ := tests.NewTestApp()
			defer testApp.Cleanup()

			beforeCollection, err := testApp.FindCollectionByNameOrId("clients")
			if err != nil {
				t.Fatal(err)
			}

			err = testApp.ImportCollections([]map[string]any{s.data}, s.deleteMissing)
			if err != nil {
				t.Fatal(err)
			}

			afterCollection, err := testApp.FindCollectionByNameOrId("clients")
			if err != nil {
				t.Fatal(err)
			}

			if afterCollection.AuthToken.Duration != 100 {
				t.Fatalf("Expected AuthToken duration to be %d, got %d", 100, afterCollection.AuthToken.Duration)
			}
			if beforeCollection.AuthToken.Secret != afterCollection.AuthToken.Secret {
				t.Fatalf("Expected AuthToken secrets to remain the same, got\n%q\nVS\n%q", beforeCollection.AuthToken.Secret, afterCollection.AuthToken.Secret)
			}
			if beforeCollection.Name != afterCollection.Name {
				t.Fatalf("Expected Name to remain the same, got\n%q\nVS\n%q", beforeCollection.Name, afterCollection.Name)
			}
			if beforeCollection.Id != afterCollection.Id {
				t.Fatalf("Expected Id to remain the same, got\n%q\nVS\n%q", beforeCollection.Id, afterCollection.Id)
			}

			if !s.deleteMissing {
				totalExpectedFields := len(beforeCollection.Fields) + 1
				if v := len(afterCollection.Fields); v != totalExpectedFields {
					t.Fatalf("Expected %d total fields, got %d", totalExpectedFields, v)
				}

				if afterCollection.Fields.GetByName("test") == nil {
					t.Fatalf("Missing new field %q", "test")
				}

				// ensure that the old fields still exist
				oldFields := beforeCollection.Fields.FieldNames()
				for _, name := range oldFields {
					if afterCollection.Fields.GetByName(name) == nil {
						t.Fatalf("Missing expected old field %q", name)
					}
				}
			} else {
				totalExpectedFields := 1
				for _, f := range beforeCollection.Fields {
					if f.GetSystem() {
						totalExpectedFields++
					}
				}

				if v := len(afterCollection.Fields); v != totalExpectedFields {
					t.Fatalf("Expected %d total fields, got %d", totalExpectedFields, v)
				}

				if afterCollection.Fields.GetByName("test") == nil {
					t.Fatalf("Missing new field %q", "test")
				}

				// ensure that the old system fields still exist
				for _, f := range beforeCollection.Fields {
					if f.GetSystem() && afterCollection.Fields.GetByName(f.GetName()) == nil {
						t.Fatalf("Missing expected old field %q", f.GetName())
					}
				}
			}
		})
	}
}

func TestImportCollectionsCreateRules(t *testing.T) {
	t.Parallel()

	testApp, _ := tests.NewTestApp()
	defer testApp.Cleanup()

	err := testApp.ImportCollections([]map[string]any{
		{"name": "new_test", "type": "auth", "authToken": map[string]any{"duration": 123}, "fields": []map[string]any{{"name": "test", "type": "text"}}},
	}, false)
	if err != nil {
		t.Fatal(err)
	}

	collection, err := testApp.FindCollectionByNameOrId("new_test")
	if err != nil {
		t.Fatal(err)
	}

	raw, err := json.Marshal(collection)
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)

	expectedParts := []string{
		`"name":"new_test"`,
		`"fields":[`,
		`"name":"id"`,
		`"name":"email"`,
		`"name":"tokenKey"`,
		`"name":"password"`,
		`"name":"test"`,
		`"indexes":[`,
		`CREATE UNIQUE INDEX`,
		`"duration":123`,
	}

	for _, part := range expectedParts {
		if !strings.Contains(rawStr, part) {
			t.Errorf("Missing %q in\n%s", part, rawStr)
		}
	}
}
