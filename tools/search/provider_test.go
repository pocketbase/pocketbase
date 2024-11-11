package search

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/list"
	_ "modernc.org/sqlite"
)

func TestNewProvider(t *testing.T) {
	r := &testFieldResolver{}
	p := NewProvider(r)

	if p.page != 1 {
		t.Fatalf("Expected page %d, got %d", 1, p.page)
	}

	if p.perPage != DefaultPerPage {
		t.Fatalf("Expected perPage %d, got %d", DefaultPerPage, p.perPage)
	}

	if p.maxFilterExprLimit != DefaultFilterExprLimit {
		t.Fatalf("Expected maxFilterExprLimit %d, got %d", DefaultFilterExprLimit, p.maxFilterExprLimit)
	}

	if p.maxSortExprLimit != DefaultSortExprLimit {
		t.Fatalf("Expected maxSortExprLimit %d, got %d", DefaultSortExprLimit, p.maxSortExprLimit)
	}
}

func TestMaxFilterExprLimit(t *testing.T) {
	p := NewProvider(&testFieldResolver{})

	testVals := []int{0, -10, 10}

	for _, val := range testVals {
		t.Run("max_"+strconv.Itoa(val), func(t *testing.T) {
			p.MaxFilterExprLimit(val)

			if p.maxFilterExprLimit != val {
				t.Fatalf("Expected maxFilterExprLimit to change to %d, got %d", val, p.maxFilterExprLimit)
			}
		})
	}
}

func TestMaxSortExprLimit(t *testing.T) {
	p := NewProvider(&testFieldResolver{})

	testVals := []int{0, -10, 10}

	for _, val := range testVals {
		t.Run("max_"+strconv.Itoa(val), func(t *testing.T) {
			p.MaxSortExprLimit(val)

			if p.maxSortExprLimit != val {
				t.Fatalf("Expected maxSortExprLimit to change to %d, got %d", val, p.maxSortExprLimit)
			}
		})
	}
}

func TestProviderQuery(t *testing.T) {
	db := dbx.NewFromDB(nil, "")
	query := db.Select("id").From("test")
	querySql := query.Build().SQL()

	r := &testFieldResolver{}
	p := NewProvider(r).Query(query)

	expected := p.query.Build().SQL()

	if querySql != expected {
		t.Fatalf("Expected %v, got %v", expected, querySql)
	}
}

func TestProviderSkipTotal(t *testing.T) {
	p := NewProvider(&testFieldResolver{})

	if p.skipTotal {
		t.Fatalf("Expected the default skipTotal to be %v, got %v", false, p.skipTotal)
	}

	p.SkipTotal(true)

	if !p.skipTotal {
		t.Fatalf("Expected skipTotal to change to %v, got %v", true, p.skipTotal)
	}
}

func TestProviderCountCol(t *testing.T) {
	p := NewProvider(&testFieldResolver{})

	if p.countCol != "id" {
		t.Fatalf("Expected the default countCol to be %s, got %s", "id", p.countCol)
	}

	p.CountCol("test")

	if p.countCol != "test" {
		t.Fatalf("Expected colCount to change to %s, got %s", "test", p.countCol)
	}
}

func TestProviderPage(t *testing.T) {
	r := &testFieldResolver{}
	p := NewProvider(r).Page(10)

	if p.page != 10 {
		t.Fatalf("Expected page %v, got %v", 10, p.page)
	}
}

func TestProviderPerPage(t *testing.T) {
	r := &testFieldResolver{}
	p := NewProvider(r).PerPage(456)

	if p.perPage != 456 {
		t.Fatalf("Expected perPage %v, got %v", 456, p.perPage)
	}
}

func TestProviderSort(t *testing.T) {
	initialSort := []SortField{{"test1", SortAsc}, {"test2", SortAsc}}
	r := &testFieldResolver{}
	p := NewProvider(r).
		Sort(initialSort).
		AddSort(SortField{"test3", SortDesc})

	encoded, _ := json.Marshal(p.sort)
	expected := `[{"name":"test1","direction":"ASC"},{"name":"test2","direction":"ASC"},{"name":"test3","direction":"DESC"}]`

	if string(encoded) != expected {
		t.Fatalf("Expected sort %v, got \n%v", expected, string(encoded))
	}
}

func TestProviderFilter(t *testing.T) {
	initialFilter := []FilterData{"test1", "test2"}
	r := &testFieldResolver{}
	p := NewProvider(r).
		Filter(initialFilter).
		AddFilter("test3")

	encoded, _ := json.Marshal(p.filter)
	expected := `["test1","test2","test3"]`

	if string(encoded) != expected {
		t.Fatalf("Expected filter %v, got \n%v", expected, string(encoded))
	}
}

func TestProviderParse(t *testing.T) {
	initialPage := 2
	initialPerPage := 123
	initialSort := []SortField{{"test1", SortAsc}, {"test2", SortAsc}}
	initialFilter := []FilterData{"test1", "test2"}

	scenarios := []struct {
		query         string
		expectError   bool
		expectPage    int
		expectPerPage int
		expectSort    string
		expectFilter  string
	}{
		// empty
		{
			"",
			false,
			initialPage,
			initialPerPage,
			`[{"name":"test1","direction":"ASC"},{"name":"test2","direction":"ASC"}]`,
			`["test1","test2"]`,
		},
		// invalid query
		{
			"invalid;",
			true,
			initialPage,
			initialPerPage,
			`[{"name":"test1","direction":"ASC"},{"name":"test2","direction":"ASC"}]`,
			`["test1","test2"]`,
		},
		// invalid page
		{
			"page=a",
			true,
			initialPage,
			initialPerPage,
			`[{"name":"test1","direction":"ASC"},{"name":"test2","direction":"ASC"}]`,
			`["test1","test2"]`,
		},
		// invalid perPage
		{
			"perPage=a",
			true,
			initialPage,
			initialPerPage,
			`[{"name":"test1","direction":"ASC"},{"name":"test2","direction":"ASC"}]`,
			`["test1","test2"]`,
		},
		// valid query parameters
		{
			"page=3&perPage=456&filter=test3&sort=-a,b,+c&other=123",
			false,
			3,
			456,
			`[{"name":"test1","direction":"ASC"},{"name":"test2","direction":"ASC"},{"name":"a","direction":"DESC"},{"name":"b","direction":"ASC"},{"name":"c","direction":"ASC"}]`,
			`["test1","test2","test3"]`,
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.query), func(t *testing.T) {
			r := &testFieldResolver{}
			p := NewProvider(r).
				Page(initialPage).
				PerPage(initialPerPage).
				Sort(initialSort).
				Filter(initialFilter)

			err := p.Parse(s.query)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if p.page != s.expectPage {
				t.Fatalf("Expected page %v, got %v", s.expectPage, p.page)
			}

			if p.perPage != s.expectPerPage {
				t.Fatalf("Expected perPage %v, got %v", s.expectPerPage, p.perPage)
			}

			encodedSort, _ := json.Marshal(p.sort)
			if string(encodedSort) != s.expectSort {
				t.Fatalf("Expected sort %v, got \n%v", s.expectSort, string(encodedSort))
			}

			encodedFilter, _ := json.Marshal(p.filter)
			if string(encodedFilter) != s.expectFilter {
				t.Fatalf("Expected filter %v, got \n%v", s.expectFilter, string(encodedFilter))
			}
		})
	}
}

func TestProviderExecEmptyQuery(t *testing.T) {
	p := NewProvider(&testFieldResolver{}).
		Query(nil)

	_, err := p.Exec(&[]testTableStruct{})
	if err == nil {
		t.Fatalf("Expected error with empty query, got nil")
	}
}

func TestProviderExecNonEmptyQuery(t *testing.T) {
	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	query := testDB.Select("*").
		From("test").
		Where(dbx.Not(dbx.HashExp{"test1": nil})).
		OrderBy("test1 ASC")

	scenarios := []struct {
		name          string
		page          int
		perPage       int
		sort          []SortField
		filter        []FilterData
		skipTotal     bool
		expectError   bool
		expectResult  string
		expectQueries []string
	}{
		{
			"page normalization",
			-1,
			10,
			[]SortField{},
			[]FilterData{},
			false,
			false,
			`{"items":[{"test1":1,"test2":"test2.1","test3":""},{"test1":2,"test2":"test2.2","test3":""}],"page":1,"perPage":10,"totalItems":2,"totalPages":1}`,
			[]string{
				"SELECT COUNT(DISTINCT [[test.id]]) FROM `test` WHERE NOT (`test1` IS NULL)",
				"SELECT * FROM `test` WHERE NOT (`test1` IS NULL) ORDER BY `test1` ASC LIMIT 10",
			},
		},
		{
			"perPage normalization",
			10,
			0, // fallback to default
			[]SortField{},
			[]FilterData{},
			false,
			false,
			`{"items":[],"page":10,"perPage":30,"totalItems":2,"totalPages":1}`,
			[]string{
				"SELECT COUNT(DISTINCT [[test.id]]) FROM `test` WHERE NOT (`test1` IS NULL)",
				"SELECT * FROM `test` WHERE NOT (`test1` IS NULL) ORDER BY `test1` ASC LIMIT 30 OFFSET 270",
			},
		},
		{
			"invalid sort field",
			1,
			10,
			[]SortField{{"unknown", SortAsc}},
			[]FilterData{},
			false,
			true,
			"",
			nil,
		},
		{
			"invalid filter",
			1,
			10,
			[]SortField{},
			[]FilterData{"test2 = 'test2.1'", "invalid"},
			false,
			true,
			"",
			nil,
		},
		{
			"valid sort and filter fields",
			1,
			5555, // will be limited by MaxPerPage
			[]SortField{{"test2", SortDesc}},
			[]FilterData{"test2 != null", "test1 >= 2"},
			false,
			false,
			`{"items":[{"test1":2,"test2":"test2.2","test3":""}],"page":1,"perPage":` + fmt.Sprint(MaxPerPage) + `,"totalItems":1,"totalPages":1}`,
			[]string{
				"SELECT COUNT(DISTINCT [[test.id]]) FROM `test` WHERE ((NOT (`test1` IS NULL)) AND (((test2 IS NOT '' AND test2 IS NOT NULL)))) AND (test1 >= 2)",
				"SELECT * FROM `test` WHERE ((NOT (`test1` IS NULL)) AND (((test2 IS NOT '' AND test2 IS NOT NULL)))) AND (test1 >= 2) ORDER BY `test1` ASC, `test2` DESC LIMIT " + fmt.Sprint(MaxPerPage),
			},
		},
		{
			"valid sort and filter fields (skipTotal=1)",
			1,
			5555, // will be limited by MaxPerPage
			[]SortField{{"test2", SortDesc}},
			[]FilterData{"test2 != null", "test1 >= 2"},
			true,
			false,
			`{"items":[{"test1":2,"test2":"test2.2","test3":""}],"page":1,"perPage":` + fmt.Sprint(MaxPerPage) + `,"totalItems":-1,"totalPages":-1}`,
			[]string{
				"SELECT * FROM `test` WHERE ((NOT (`test1` IS NULL)) AND (((test2 IS NOT '' AND test2 IS NOT NULL)))) AND (test1 >= 2) ORDER BY `test1` ASC, `test2` DESC LIMIT " + fmt.Sprint(MaxPerPage),
			},
		},
		{
			"valid sort and filter fields (zero results)",
			1,
			10,
			[]SortField{{"test3", SortAsc}},
			[]FilterData{"test3 != ''"},
			false,
			false,
			`{"items":[],"page":1,"perPage":10,"totalItems":0,"totalPages":0}`,
			[]string{
				"SELECT COUNT(DISTINCT [[test.id]]) FROM `test` WHERE (NOT (`test1` IS NULL)) AND (((test3 IS NOT '' AND test3 IS NOT NULL)))",
				"SELECT * FROM `test` WHERE (NOT (`test1` IS NULL)) AND (((test3 IS NOT '' AND test3 IS NOT NULL))) ORDER BY `test1` ASC, `test3` ASC LIMIT 10",
			},
		},
		{
			"valid sort and filter fields (zero results; skipTotal=1)",
			1,
			10,
			[]SortField{{"test3", SortAsc}},
			[]FilterData{"test3 != ''"},
			true,
			false,
			`{"items":[],"page":1,"perPage":10,"totalItems":-1,"totalPages":-1}`,
			[]string{
				"SELECT * FROM `test` WHERE (NOT (`test1` IS NULL)) AND (((test3 IS NOT '' AND test3 IS NOT NULL))) ORDER BY `test1` ASC, `test3` ASC LIMIT 10",
			},
		},
		{
			"pagination test",
			2,
			1,
			[]SortField{},
			[]FilterData{},
			false,
			false,
			`{"items":[{"test1":2,"test2":"test2.2","test3":""}],"page":2,"perPage":1,"totalItems":2,"totalPages":2}`,
			[]string{
				"SELECT COUNT(DISTINCT [[test.id]]) FROM `test` WHERE NOT (`test1` IS NULL)",
				"SELECT * FROM `test` WHERE NOT (`test1` IS NULL) ORDER BY `test1` ASC LIMIT 1 OFFSET 1",
			},
		},
		{
			"pagination test (skipTotal=1)",
			2,
			1,
			[]SortField{},
			[]FilterData{},
			true,
			false,
			`{"items":[{"test1":2,"test2":"test2.2","test3":""}],"page":2,"perPage":1,"totalItems":-1,"totalPages":-1}`,
			[]string{
				"SELECT * FROM `test` WHERE NOT (`test1` IS NULL) ORDER BY `test1` ASC LIMIT 1 OFFSET 1",
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testDB.CalledQueries = []string{} // reset

			testResolver := &testFieldResolver{}
			p := NewProvider(testResolver).
				Query(query).
				Page(s.page).
				PerPage(s.perPage).
				Sort(s.sort).
				SkipTotal(s.skipTotal).
				Filter(s.filter)

			result, err := p.Exec(&[]testTableStruct{})

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if testResolver.UpdateQueryCalls != 1 {
				t.Fatalf("Expected resolver.Update to be called %d, got %d", 1, testResolver.UpdateQueryCalls)
			}

			encoded, _ := json.Marshal(result)
			if string(encoded) != s.expectResult {
				t.Fatalf("Expected result %v, got \n%v", s.expectResult, string(encoded))
			}

			if len(s.expectQueries) != len(testDB.CalledQueries) {
				t.Fatalf("Expected %d queries, got %d: \n%v", len(s.expectQueries), len(testDB.CalledQueries), testDB.CalledQueries)
			}

			for _, q := range testDB.CalledQueries {
				if !list.ExistInSliceWithRegex(q, s.expectQueries) {
					t.Fatalf("Didn't expect query \n%v \nin \n%v", q, s.expectQueries)
				}
			}
		})
	}
}

func TestProviderFilterAndSortLimits(t *testing.T) {
	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	query := testDB.Select("*").
		From("test").
		Where(dbx.Not(dbx.HashExp{"test1": nil})).
		OrderBy("test1 ASC")

	scenarios := []struct {
		name               string
		filter             []FilterData
		sort               []SortField
		maxFilterExprLimit int
		maxSortExprLimit   int
		expectError        bool
	}{
		// filter
		{
			"<= max filter length",
			[]FilterData{
				"1=2",
				FilterData("1='" + strings.Repeat("a", MaxFilterLength-4) + "'"),
			},
			[]SortField{},
			1,
			0,
			false,
		},
		{
			"> max filter length",
			[]FilterData{
				"1=2",
				FilterData("1='" + strings.Repeat("a", MaxFilterLength-3) + "'"),
			},
			[]SortField{},
			1,
			0,
			true,
		},
		{
			"<= max filter exprs",
			[]FilterData{
				"1=2",
				"(1=1 || 1=1) && (1=1 || (1=1 || 1=1)) && (1=1)",
			},
			[]SortField{},
			6,
			0,
			false,
		},
		{
			"> max filter exprs",
			[]FilterData{
				"1=2",
				"(1=1 || 1=1) && (1=1 || (1=1 || 1=1)) && (1=1)",
			},
			[]SortField{},
			5,
			0,
			true,
		},

		// sort
		{
			"<= max sort field length",
			[]FilterData{},
			[]SortField{
				{"id", SortAsc},
				{"test1", SortDesc},
				{strings.Repeat("a", MaxSortFieldLength), SortDesc},
			},
			0,
			10,
			false,
		},
		{
			"> max sort field length",
			[]FilterData{},
			[]SortField{
				{"id", SortAsc},
				{"test1", SortDesc},
				{strings.Repeat("b", MaxSortFieldLength+1), SortDesc},
			},
			0,
			10,
			true,
		},
		{
			"<= max sort exprs",
			[]FilterData{},
			[]SortField{
				{"id", SortAsc},
				{"test1", SortDesc},
			},
			0,
			2,
			false,
		},
		{
			"> max sort exprs",
			[]FilterData{},
			[]SortField{
				{"id", SortAsc},
				{"test1", SortDesc},
			},
			0,
			1,
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testResolver := &testFieldResolver{}
			p := NewProvider(testResolver).
				Query(query).
				Sort(s.sort).
				Filter(s.filter).
				MaxFilterExprLimit(s.maxFilterExprLimit).
				MaxSortExprLimit(s.maxSortExprLimit)

			_, err := p.Exec(&[]testTableStruct{})

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectError, hasErr)
			}
		})
	}
}

func TestProviderParseAndExec(t *testing.T) {
	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	query := testDB.Select("*").
		From("test").
		Where(dbx.Not(dbx.HashExp{"test1": nil})).
		OrderBy("test1 ASC")

	scenarios := []struct {
		name         string
		queryString  string
		expectError  bool
		expectResult string
	}{
		{
			"no extra query params (aka. use the provider presets)",
			"",
			false,
			`{"items":[],"page":2,"perPage":123,"totalItems":2,"totalPages":1}`,
		},
		{
			"invalid query",
			"invalid;",
			true,
			"",
		},
		{
			"invalid page",
			"page=a",
			true,
			"",
		},
		{
			"invalid perPage",
			"perPage=a",
			true,
			"",
		},
		{
			"invalid skipTotal",
			"skipTotal=a",
			true,
			"",
		},
		{
			"invalid sorting field",
			"sort=-unknown",
			true,
			"",
		},
		{
			"invalid filter field",
			"filter=unknown>1",
			true,
			"",
		},
		{
			"page > existing",
			"page=3&perPage=9999",
			false,
			`{"items":[],"page":3,"perPage":1000,"totalItems":2,"totalPages":1}`,
		},
		{
			"valid query params",
			"page=1&perPage=9999&filter=test1>1&sort=-test2,test3",
			false,
			`{"items":[{"test1":2,"test2":"test2.2","test3":""}],"page":1,"perPage":1000,"totalItems":1,"totalPages":1}`,
		},
		{
			"valid query params with skipTotal=1",
			"page=1&perPage=9999&filter=test1>1&sort=-test2,test3&skipTotal=1",
			false,
			`{"items":[{"test1":2,"test2":"test2.2","test3":""}],"page":1,"perPage":1000,"totalItems":-1,"totalPages":-1}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			testDB.CalledQueries = []string{} // reset

			testResolver := &testFieldResolver{}
			provider := NewProvider(testResolver).
				Query(query).
				Page(2).
				PerPage(123).
				Sort([]SortField{{"test2", SortAsc}}).
				Filter([]FilterData{"test1 > 0"})

			result, err := provider.ParseAndExec(s.queryString, &[]testTableStruct{})

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			if hasErr {
				return
			}

			if testResolver.UpdateQueryCalls != 1 {
				t.Fatalf("Expected resolver.Update to be called %d, got %d", 1, testResolver.UpdateQueryCalls)
			}

			expectedQueries := 2
			if provider.skipTotal {
				expectedQueries = 1
			}

			if len(testDB.CalledQueries) != expectedQueries {
				t.Fatalf("Expected %d db queries, got %d: \n%v", expectedQueries, len(testDB.CalledQueries), testDB.CalledQueries)
			}

			encoded, _ := json.Marshal(result)
			if string(encoded) != s.expectResult {
				t.Fatalf("Expected result \n%v\ngot\n%v", s.expectResult, string(encoded))
			}
		})
	}
}

// -------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------

type testTableStruct struct {
	Test1 int    `db:"test1" json:"test1"`
	Test2 string `db:"test2" json:"test2"`
	Test3 string `db:"test3" json:"test3"`
}

type testDB struct {
	*dbx.DB
	CalledQueries []string
}

// NB! Don't forget to call `db.Close()` at the end of the test.
func createTestDB() (*testDB, error) {
	// using a shared cache to allow multiple connections access to
	// the same in memory database https://www.sqlite.org/inmemorydb.html
	sqlDB, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}

	db := testDB{DB: dbx.NewFromDB(sqlDB, "sqlite")}
	db.CreateTable("test", map[string]string{
		"id":                                    "int default 0",
		"test1":                                 "int default 0",
		"test2":                                 "text default ''",
		"test3":                                 "text default ''",
		strings.Repeat("a", MaxSortFieldLength): "text default ''",
		strings.Repeat("b", MaxSortFieldLength+1): "text default ''",
	}).Execute()
	db.Insert("test", dbx.Params{"id": 1, "test1": 1, "test2": "test2.1"}).Execute()
	db.Insert("test", dbx.Params{"id": 2, "test1": 2, "test2": "test2.2"}).Execute()
	db.QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		db.CalledQueries = append(db.CalledQueries, sql)
	}

	return &db, nil
}

// ---

type testFieldResolver struct {
	UpdateQueryCalls int
	ResolveCalls     int
}

func (t *testFieldResolver) UpdateQuery(query *dbx.SelectQuery) error {
	t.UpdateQueryCalls++
	return nil
}

func (t *testFieldResolver) Resolve(field string) (*ResolverResult, error) {
	t.ResolveCalls++

	if field == "unknown" {
		return nil, errors.New("test error")
	}

	return &ResolverResult{Identifier: field}, nil
}
