package daos_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
)

func ensureNoTempViews(app core.App, t *testing.T) {
	var total int

	err := app.Dao().DB().Select("count(*)").
		From("sqlite_schema").
		AndWhere(dbx.HashExp{"type": "view"}).
		AndWhere(dbx.NewExp(`[[name]] LIKE '%\_temp\_%' ESCAPE '\'`)).
		Limit(1).
		Row(&total)
	if err != nil {
		t.Fatalf("Failed to check for temp views: %v", err)
	}

	if total > 0 {
		t.Fatalf("Expected all temp views to be deleted, got %d", total)
	}
}

func TestDeleteView(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		viewName    string
		expectError bool
	}{
		{"", true},
		{"demo1", true},    // not a view table
		{"missing", false}, // missing or already deleted
		{"view1", false},   // existing
		{"VieW1", false},   // view names are case insensitives
	}

	for i, s := range scenarios {
		err := app.Dao().DeleteView(s.viewName)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%d - %q] Expected hasErr %v, got %v (%v)", i, s.viewName, s.expectError, hasErr, err)
		}
	}

	ensureNoTempViews(app, t)
}

func TestSaveView(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		scenarioName  string
		viewName      string
		query         string
		expectError   bool
		expectColumns []string
	}{
		{
			"empty name and query",
			"",
			"",
			true,
			nil,
		},
		{
			"empty name",
			"",
			"select * from _admins",
			true,
			nil,
		},
		{
			"empty query",
			"123Test",
			"",
			true,
			nil,
		},
		{
			"invalid query",
			"123Test",
			"123 456",
			true,
			nil,
		},
		{
			"missing table",
			"123Test",
			"select id from missing",
			true,
			nil,
		},
		{
			"non select query",
			"123Test",
			"drop table _admins",
			true,
			nil,
		},
		{
			"multiple select queries",
			"123Test",
			"select *, count(id) as c  from _admins; select * from demo1;",
			true,
			nil,
		},
		{
			"try to break the parent parenthesis",
			"123Test",
			"select *, count(id) as c  from `_admins`)",
			true,
			nil,
		},
		{
			"simple select query (+ trimmed semicolon)",
			"123Test",
			";select *, count(id) as c  from _admins;",
			false,
			[]string{
				"id", "created", "updated",
				"passwordHash", "tokenKey", "email",
				"lastResetSentAt", "avatar", "c",
			},
		},
		{
			"update old view with new query",
			"123Test",
			"select 1 as test from _admins",
			false,
			[]string{"test"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.scenarioName, func(t *testing.T) {
			err := app.Dao().SaveView(s.viewName, s.query)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			infoRows, err := app.Dao().TableInfo(s.viewName)
			if err != nil {
				t.Fatalf("Failed to fetch table info for %s: %v", s.viewName, err)
			}

			if len(s.expectColumns) != len(infoRows) {
				t.Fatalf("Expected %d columns, got %d", len(s.expectColumns), len(infoRows))
			}

			for _, row := range infoRows {
				if !list.ExistInSlice(row.Name, s.expectColumns) {
					t.Fatalf("Missing %q column in %v", row.Name, s.expectColumns)
				}
			}
		})
	}

	ensureNoTempViews(app, t)
}

func TestCreateViewSchemaWithDiscardedNestedTransaction(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		_, err := txDao.CreateViewSchema("select id from missing")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		return nil
	})

	ensureNoTempViews(app, t)
}

func TestCreateViewSchema(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	scenarios := []struct {
		name         string
		query        string
		expectError  bool
		expectFields map[string]string // name-type pairs
	}{
		{
			"empty query",
			"",
			true,
			nil,
		},
		{
			"invalid query",
			"test 123456",
			true,
			nil,
		},
		{
			"missing table",
			"select id from missing",
			true,
			nil,
		},
		{
			"query with wildcard column",
			"select a.id, a.* from demo1 a",
			true,
			nil,
		},
		{
			"query without id",
			"select text, url, created, updated from demo1",
			true,
			nil,
		},
		{
			"query with comments",
			`
				select
				-- test single line
				id,
				text,
				/* multi
					line comment */
				url, created, updated from demo1
			`,
			false,
			map[string]string{
				"text": schema.FieldTypeText,
				"url":  schema.FieldTypeUrl,
			},
		},
		{
			"query with all fields and quoted identifiers",
			`
				select
					"id",
					"created",
					"updated",
					[text],
					` + "`bool`" + `,
					"url",
					"select_one",
					"select_many",
					"file_one",
					"demo1"."file_many",
					` + "`demo1`." + "`number`" + ` number_alias,
					"email",
					"datetime",
					"json",
					"rel_one",
					"rel_many",
					'single_quoted_custom_literal' as 'single_quoted_column'
				from demo1
			`,
			false,
			map[string]string{
				"text":                 schema.FieldTypeText,
				"bool":                 schema.FieldTypeBool,
				"url":                  schema.FieldTypeUrl,
				"select_one":           schema.FieldTypeSelect,
				"select_many":          schema.FieldTypeSelect,
				"file_one":             schema.FieldTypeFile,
				"file_many":            schema.FieldTypeFile,
				"number_alias":         schema.FieldTypeNumber,
				"email":                schema.FieldTypeEmail,
				"datetime":             schema.FieldTypeDate,
				"json":                 schema.FieldTypeJson,
				"rel_one":              schema.FieldTypeRelation,
				"rel_many":             schema.FieldTypeRelation,
				"single_quoted_column": schema.FieldTypeJson,
			},
		},
		{
			"query with indirect relations fields",
			"select a.id, b.id as bid, b.created from demo1 as a left join demo2 b",
			false,
			map[string]string{
				"bid": schema.FieldTypeRelation,
			},
		},
		{
			"query with multiple froms, joins and style of aliasses",
			`
				select
					a.id as id,
					b.id as bid,
					lj.id cid,
					ij.id as did,
					a.bool,
					_admins.id as eid,
					_admins.email
				from demo1 a, demo2 as b
				left join demo3 lj on lj.id = 123
				inner join demo4 as ij on ij.id = 123
				join _admins
				where 1=1
				group by a.id
				limit 10
			`,
			false,
			map[string]string{
				"bid":   schema.FieldTypeRelation,
				"cid":   schema.FieldTypeRelation,
				"did":   schema.FieldTypeRelation,
				"bool":  schema.FieldTypeBool,
				"eid":   schema.FieldTypeJson, // not from collection
				"email": schema.FieldTypeJson, // not from collection
			},
		},
		{
			"query with casts",
			`select
				a.id,
				count(a.id) count,
				cast(a.id as int) cast_int,
				cast(a.id as integer) cast_integer,
				cast(a.id as real) cast_real,
				cast(a.id as decimal) cast_decimal,
				cast(a.id as numeric) cast_numeric,
				cast(a.id as text) cast_text,
				cast(a.id as bool) cast_bool,
				cast(a.id as boolean) cast_boolean,
				avg(a.id) avg,
				sum(a.id) sum,
				total(a.id) total,
				min(a.id) min,
				max(a.id) max
			from demo1 a`,
			false,
			map[string]string{
				"count":        schema.FieldTypeNumber,
				"total":        schema.FieldTypeNumber,
				"cast_int":     schema.FieldTypeNumber,
				"cast_integer": schema.FieldTypeNumber,
				"cast_real":    schema.FieldTypeNumber,
				"cast_decimal": schema.FieldTypeNumber,
				"cast_numeric": schema.FieldTypeNumber,
				"cast_text":    schema.FieldTypeText,
				"cast_bool":    schema.FieldTypeBool,
				"cast_boolean": schema.FieldTypeBool,
				// json because they are nullable
				"sum": schema.FieldTypeJson,
				"avg": schema.FieldTypeJson,
				"min": schema.FieldTypeJson,
				"max": schema.FieldTypeJson,
			},
		},
		{
			"query with reserved auth collection fields",
			`
				select
					a.id,
					a.username,
					a.email,
					a.emailVisibility,
					a.verified,
					demo1.id relid
				from users a
				left join demo1
			`,
			false,
			map[string]string{
				"username":        schema.FieldTypeText,
				"email":           schema.FieldTypeEmail,
				"emailVisibility": schema.FieldTypeBool,
				"verified":        schema.FieldTypeBool,
				"relid":           schema.FieldTypeRelation,
			},
		},
		{
			"query with unknown fields and aliases",
			`select
				id,
				id as id2,
				text as text_alias,
				url as url_alias,
				"demo1"."bool" as bool_alias,
				number as number_alias,
				created created_alias,
				updated updated_alias,
				123 as custom
			from demo1`,
			false,
			map[string]string{
				"id2":           schema.FieldTypeRelation,
				"text_alias":    schema.FieldTypeText,
				"url_alias":     schema.FieldTypeUrl,
				"bool_alias":    schema.FieldTypeBool,
				"number_alias":  schema.FieldTypeNumber,
				"created_alias": schema.FieldTypeDate,
				"updated_alias": schema.FieldTypeDate,
				"custom":        schema.FieldTypeJson,
			},
		},
		{
			"query with distinct and reordered id column",
			`select distinct
				id as id2,
				id,
				123 as custom
			from demo1`,
			false,
			map[string]string{
				"id2":    schema.FieldTypeRelation,
				"custom": schema.FieldTypeJson,
			},
		},
		{
			"query with aliasing the same field multiple times",
			`select
				a.id as id,
				a.text as alias1,
				a.text as alias2,
				b.text as alias3,
				b.text as alias4
			from demo1 a
			left join demo1 as b`,
			false,
			map[string]string{
				"alias1": schema.FieldTypeText,
				"alias2": schema.FieldTypeText,
				"alias3": schema.FieldTypeText,
				"alias4": schema.FieldTypeText,
			},
		},
	}

	for _, s := range scenarios {
		result, err := app.Dao().CreateViewSchema(s.query)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr %v, got %v (%v)", s.name, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		if len(s.expectFields) != len(result.Fields()) {
			serialized, _ := json.Marshal(result)
			t.Errorf("[%s] Expected %d fields, got %d: \n%s", s.name, len(s.expectFields), len(result.Fields()), serialized)
			continue
		}

		for name, typ := range s.expectFields {
			field := result.GetFieldByName(name)

			if field == nil {
				t.Errorf("[%s] Expected to find field %s, got nil", s.name, name)
				continue
			}

			if field.Type != typ {
				t.Errorf("[%s] Expected field %s to be %q, got %s", s.name, name, typ, field.Type)
				continue
			}
		}
	}

	ensureNoTempViews(app, t)
}

func TestFindRecordByViewFile(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	prevCollection, err := app.Dao().FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	totalLevels := 6

	// create collection view mocks
	fileOneAlias := "file_one one0"
	fileManyAlias := "file_many many0"
	mockCollections := make([]*models.Collection, 0, totalLevels)
	for i := 0; i <= totalLevels; i++ {
		view := new(models.Collection)
		view.Type = models.CollectionTypeView
		view.Name = fmt.Sprintf("_test_view%d", i)
		view.SetOptions(&models.CollectionViewOptions{
			Query: fmt.Sprintf(
				"select id, %s, %s from %s",
				fileOneAlias,
				fileManyAlias,
				prevCollection.Name,
			),
		})

		// save view
		if err := app.Dao().SaveCollection(view); err != nil {
			t.Fatalf("Failed to save view%d: %v", i, err)
		}

		mockCollections = append(mockCollections, view)
		prevCollection = view
		fileOneAlias = fmt.Sprintf("one%d one%d", i, i+1)
		fileManyAlias = fmt.Sprintf("many%d many%d", i, i+1)
	}

	fileOneName := "test_d61b33QdDU.txt"
	fileManyName := "test_QZFjKjXchk.txt"
	expectedRecordId := "84nmscqy84lsi1t"

	scenarios := []struct {
		name               string
		collectionNameOrId string
		fileFieldName      string
		filename           string
		expectError        bool
		expectRecordId     string
	}{
		{
			"missing collection",
			"missing",
			"a",
			fileOneName,
			true,
			"",
		},
		{
			"non-view collection",
			"demo1",
			"file_one",
			fileOneName,
			true,
			"",
		},
		{
			"view collection after the max recursion limit",
			mockCollections[totalLevels-1].Name,
			fmt.Sprintf("one%d", totalLevels-1),
			fileOneName,
			true,
			"",
		},
		{
			"first view collection (single file)",
			mockCollections[0].Name,
			"one0",
			fileOneName,
			false,
			expectedRecordId,
		},
		{
			"first view collection (many files)",
			mockCollections[0].Name,
			"many0",
			fileManyName,
			false,
			expectedRecordId,
		},

		{
			"last view collection before the recursion limit (single file)",
			mockCollections[totalLevels-2].Name,
			fmt.Sprintf("one%d", totalLevels-2),
			fileOneName,
			false,
			expectedRecordId,
		},
		{
			"last view collection before the recursion limit (many files)",
			mockCollections[totalLevels-2].Name,
			fmt.Sprintf("many%d", totalLevels-2),
			fileManyName,
			false,
			expectedRecordId,
		},
	}

	for _, s := range scenarios {
		record, err := app.Dao().FindRecordByViewFile(
			s.collectionNameOrId,
			s.fileFieldName,
			s.filename,
		)

		hasErr := err != nil
		if hasErr != s.expectError {
			t.Errorf("[%s] Expected hasErr %v, got %v (%v)", s.name, s.expectError, hasErr, err)
			continue
		}

		if hasErr {
			continue
		}

		if record.Id != s.expectRecordId {
			t.Errorf("[%s] Expected recordId %q, got %q", s.name, s.expectRecordId, record.Id)
		}
	}
}
