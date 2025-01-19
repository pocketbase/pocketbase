package core_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func ensureNoTempViews(app core.App, t *testing.T) {
	var total int

	err := app.DB().Select("count(*)").
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
		err := app.DeleteView(s.viewName)

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
			"select * from " + core.CollectionNameSuperusers,
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
			"drop table " + core.CollectionNameSuperusers,
			true,
			nil,
		},
		{
			"multiple select queries",
			"123Test",
			"select *, count(id) as c  from " + core.CollectionNameSuperusers + "; select * from demo1;",
			true,
			nil,
		},
		{
			"try to break the parent parenthesis",
			"123Test",
			"select *, count(id) as c  from `" + core.CollectionNameSuperusers + "`)",
			true,
			nil,
		},
		{
			"simple select query (+ trimmed semicolon)",
			"123Test",
			";select *, count(id) as c  from " + core.CollectionNameSuperusers + ";",
			false,
			[]string{
				"id", "created", "updated",
				"password", "tokenKey", "email",
				"emailVisibility", "verified",
				"c",
			},
		},
		{
			"update old view with new query",
			"123Test",
			"select 1 as test from " + core.CollectionNameSuperusers,
			false,
			[]string{"test"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.scenarioName, func(t *testing.T) {
			err := app.SaveView(s.viewName, s.query)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			infoRows, err := app.TableInfo(s.viewName)
			if err != nil {
				t.Fatalf("Failed to fetch table info for %s: %v", s.viewName, err)
			}

			if len(s.expectColumns) != len(infoRows) {
				t.Fatalf("Expected %d columns, got %d", len(s.expectColumns), len(infoRows))
			}

			for _, row := range infoRows {
				if !slices.Contains(s.expectColumns, row.Name) {
					t.Fatalf("Missing %q column in %v", row.Name, s.expectColumns)
				}
			}
		})
	}

	ensureNoTempViews(app, t)
}

func TestCreateViewFieldsWithDiscardedNestedTransaction(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	app.RunInTransaction(func(txApp core.App) error {
		_, err := txApp.CreateViewFields("select id from missing")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		return nil
	})

	ensureNoTempViews(app, t)
}

func TestCreateViewFields(t *testing.T) {
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
				demo1.id,
				demo1.text,
				/* multi
				 * line comment block */
				demo1.url, demo1.created, demo1.updated from/* inline comment block with no spaces between the identifiers */demo1
				-- comment before join
				join demo2 ON (
					-- comment inside join
					demo2.id = demo1.id
				)
				-- comment before where
				where (
					-- comment inside where
					demo2.id = demo1.id
				)
			`,
			false,
			map[string]string{
				"id":      core.FieldTypeText,
				"text":    core.FieldTypeText,
				"url":     core.FieldTypeURL,
				"created": core.FieldTypeAutodate,
				"updated": core.FieldTypeAutodate,
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
				"id":                   core.FieldTypeText,
				"created":              core.FieldTypeAutodate,
				"updated":              core.FieldTypeAutodate,
				"text":                 core.FieldTypeText,
				"bool":                 core.FieldTypeBool,
				"url":                  core.FieldTypeURL,
				"select_one":           core.FieldTypeSelect,
				"select_many":          core.FieldTypeSelect,
				"file_one":             core.FieldTypeFile,
				"file_many":            core.FieldTypeFile,
				"number_alias":         core.FieldTypeNumber,
				"email":                core.FieldTypeEmail,
				"datetime":             core.FieldTypeDate,
				"json":                 core.FieldTypeJSON,
				"rel_one":              core.FieldTypeRelation,
				"rel_many":             core.FieldTypeRelation,
				"single_quoted_column": core.FieldTypeJSON,
			},
		},
		{
			"query with indirect relations fields",
			"select a.id, b.id as bid, b.created from demo1 as a left join demo2 b",
			false,
			map[string]string{
				"id":      core.FieldTypeText,
				"bid":     core.FieldTypeRelation,
				"created": core.FieldTypeAutodate,
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
					` + core.CollectionNameSuperusers + `.id as eid,
					` + core.CollectionNameSuperusers + `.email
				from demo1 a, demo2 as b
				left join demo3 lj on lj.id = 123
				inner join demo4 as ij on ij.id = 123
				join ` + core.CollectionNameSuperusers + `
				where 1=1
				group by a.id
				limit 10
			`,
			false,
			map[string]string{
				"id":    core.FieldTypeText,
				"bid":   core.FieldTypeRelation,
				"cid":   core.FieldTypeRelation,
				"did":   core.FieldTypeRelation,
				"bool":  core.FieldTypeBool,
				"eid":   core.FieldTypeRelation,
				"email": core.FieldTypeEmail,
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
				"id":           core.FieldTypeText,
				"count":        core.FieldTypeNumber,
				"total":        core.FieldTypeNumber,
				"cast_int":     core.FieldTypeNumber,
				"cast_integer": core.FieldTypeNumber,
				"cast_real":    core.FieldTypeNumber,
				"cast_decimal": core.FieldTypeNumber,
				"cast_numeric": core.FieldTypeNumber,
				"cast_text":    core.FieldTypeText,
				"cast_bool":    core.FieldTypeBool,
				"cast_boolean": core.FieldTypeBool,
				// json because they are nullable
				"sum": core.FieldTypeJSON,
				"avg": core.FieldTypeJSON,
				"min": core.FieldTypeJSON,
				"max": core.FieldTypeJSON,
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
				"id":              core.FieldTypeText,
				"username":        core.FieldTypeText,
				"email":           core.FieldTypeEmail,
				"emailVisibility": core.FieldTypeBool,
				"verified":        core.FieldTypeBool,
				"relid":           core.FieldTypeRelation,
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
				"id":            core.FieldTypeText,
				"id2":           core.FieldTypeRelation,
				"text_alias":    core.FieldTypeText,
				"url_alias":     core.FieldTypeURL,
				"bool_alias":    core.FieldTypeBool,
				"number_alias":  core.FieldTypeNumber,
				"created_alias": core.FieldTypeAutodate,
				"updated_alias": core.FieldTypeAutodate,
				"custom":        core.FieldTypeJSON,
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
				"id2":    core.FieldTypeRelation,
				"id":     core.FieldTypeText,
				"custom": core.FieldTypeJSON,
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
				"id":     core.FieldTypeText,
				"alias1": core.FieldTypeText,
				"alias2": core.FieldTypeText,
				"alias3": core.FieldTypeText,
				"alias4": core.FieldTypeText,
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result, err := app.CreateViewFields(s.query)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if len(s.expectFields) != len(result) {
				serialized, _ := json.Marshal(result)
				t.Fatalf("Expected %d fields, got %d: \n%s", len(s.expectFields), len(result), serialized)
			}

			for name, typ := range s.expectFields {
				field := result.GetByName(name)

				if field == nil {
					t.Fatalf("Expected to find field %s, got nil", name)
				}

				if field.Type() != typ {
					t.Fatalf("Expected field %s to be %q, got %q", name, typ, field.Type())
				}
			}
		})
	}

	ensureNoTempViews(app, t)
}

func TestFindRecordByViewFile(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	prevCollection, err := app.FindCollectionByNameOrId("demo1")
	if err != nil {
		t.Fatal(err)
	}

	totalLevels := 6

	// create collection view mocks
	fileOneAlias := "file_one one0"
	fileManyAlias := "file_many many0"
	mockCollections := make([]*core.Collection, 0, totalLevels)
	for i := 0; i <= totalLevels; i++ {
		view := new(core.Collection)
		view.Type = core.CollectionTypeView
		view.Name = fmt.Sprintf("_test_view%d", i)
		view.ViewQuery = fmt.Sprintf(
			"select id, %s, %s from %s",
			fileOneAlias,
			fileManyAlias,
			prevCollection.Name,
		)

		// save view
		if err := app.Save(view); err != nil {
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
		t.Run(s.name, func(t *testing.T) {
			record, err := app.FindRecordByViewFile(
				s.collectionNameOrId,
				s.fileFieldName,
				s.filename,
			)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if record.Id != s.expectRecordId {
				t.Fatalf("Expected recordId %q, got %q", s.expectRecordId, record.Id)
			}
		})
	}
}
