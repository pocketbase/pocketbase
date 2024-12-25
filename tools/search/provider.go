package search

import (
	"errors"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"golang.org/x/sync/errgroup"
)

const (
	// DefaultPerPage specifies the default number of returned search result items.
	DefaultPerPage int = 30

	// DefaultFilterExprLimit specifies the default filter expressions limit.
	DefaultFilterExprLimit int = 200

	// DefaultSortExprLimit specifies the default sort expressions limit.
	DefaultSortExprLimit int = 8

	// MaxPerPage specifies the max allowed search result items returned in a single page.
	MaxPerPage int = 1000

	// MaxFilterLength specifies the max allowed individual search filter parsable length.
	MaxFilterLength int = 3500

	// MaxSortFieldLength specifies the max allowed individual sort field parsable length.
	MaxSortFieldLength int = 255
)

// Common search errors.
var (
	ErrEmptyQuery           = errors.New("search query is not set")
	ErrSortExprLimit        = errors.New("max sort expressions limit reached")
	ErrFilterExprLimit      = errors.New("max filter expressions limit reached")
	ErrFilterLengthLimit    = errors.New("max filter length limit reached")
	ErrSortFieldLengthLimit = errors.New("max sort field length limit reached")
)

// URL search query params
const (
	PageQueryParam      string = "page"
	PerPageQueryParam   string = "perPage"
	SortQueryParam      string = "sort"
	FilterQueryParam    string = "filter"
	SkipTotalQueryParam string = "skipTotal"
)

// Result defines the returned search result structure.
type Result struct {
	Items      any `json:"items"`
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// Provider represents a single configured search provider instance.
type Provider struct {
	fieldResolver      FieldResolver
	query              *dbx.SelectQuery
	countCol           string
	sort               []SortField
	filter             []FilterData
	page               int
	perPage            int
	skipTotal          bool
	maxFilterExprLimit int
	maxSortExprLimit   int
}

// NewProvider initializes and returns a new search provider.
//
// Example:
//
//	baseQuery := db.Select("*").From("user")
//	fieldResolver := search.NewSimpleFieldResolver("id", "name")
//	models := []*YourDataStruct{}
//
//	result, err := search.NewProvider(fieldResolver).
//		Query(baseQuery).
//		ParseAndExec("page=2&filter=id>0&sort=-email", &models)
func NewProvider(fieldResolver FieldResolver) *Provider {
	return &Provider{
		fieldResolver:      fieldResolver,
		countCol:           "id",
		page:               1,
		perPage:            DefaultPerPage,
		sort:               []SortField{},
		filter:             []FilterData{},
		maxFilterExprLimit: DefaultFilterExprLimit,
		maxSortExprLimit:   DefaultSortExprLimit,
	}
}

// MaxFilterExprLimit changes the default max allowed filter expressions.
//
// Note that currently the limit is applied individually for each separate filter.
func (s *Provider) MaxFilterExprLimit(max int) *Provider {
	s.maxFilterExprLimit = max
	return s
}

// MaxSortExprLimit changes the default max allowed sort expressions.
func (s *Provider) MaxSortExprLimit(max int) *Provider {
	s.maxSortExprLimit = max
	return s
}

// Query sets the base query that will be used to fetch the search items.
func (s *Provider) Query(query *dbx.SelectQuery) *Provider {
	s.query = query
	return s
}

// SkipTotal changes the `skipTotal` field of the current search provider.
func (s *Provider) SkipTotal(skipTotal bool) *Provider {
	s.skipTotal = skipTotal
	return s
}

// CountCol allows changing the default column (id) that is used
// to generate the COUNT SQL query statement.
//
// This field is ignored if skipTotal is true.
func (s *Provider) CountCol(name string) *Provider {
	s.countCol = name
	return s
}

// Page sets the `page` field of the current search provider.
//
// Normalization on the `page` value is done during `Exec()`.
func (s *Provider) Page(page int) *Provider {
	s.page = page
	return s
}

// PerPage sets the `perPage` field of the current search provider.
//
// Normalization on the `perPage` value is done during `Exec()`.
func (s *Provider) PerPage(perPage int) *Provider {
	s.perPage = perPage
	return s
}

// Sort sets the `sort` field of the current search provider.
func (s *Provider) Sort(sort []SortField) *Provider {
	s.sort = sort
	return s
}

// AddSort appends the provided SortField to the existing provider's sort field.
func (s *Provider) AddSort(field SortField) *Provider {
	s.sort = append(s.sort, field)
	return s
}

// Filter sets the `filter` field of the current search provider.
func (s *Provider) Filter(filter []FilterData) *Provider {
	s.filter = filter
	return s
}

// AddFilter appends the provided FilterData to the existing provider's filter field.
func (s *Provider) AddFilter(filter FilterData) *Provider {
	if filter != "" {
		s.filter = append(s.filter, filter)
	}
	return s
}

// Parse parses the search query parameter from the provided query string
// and assigns the found fields to the current search provider.
//
// The data from the "sort" and "filter" query parameters are appended
// to the existing provider's `sort` and `filter` fields
// (aka. using `AddSort` and `AddFilter`).
func (s *Provider) Parse(urlQuery string) error {
	params, err := url.ParseQuery(urlQuery)
	if err != nil {
		return err
	}

	if raw := params.Get(SkipTotalQueryParam); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		s.SkipTotal(v)
	}

	if raw := params.Get(PageQueryParam); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return err
		}
		s.Page(v)
	}

	if raw := params.Get(PerPageQueryParam); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return err
		}
		s.PerPage(v)
	}

	if raw := params.Get(SortQueryParam); raw != "" {
		for _, sortField := range ParseSortFromString(raw) {
			s.AddSort(sortField)
		}
	}

	if raw := params.Get(FilterQueryParam); raw != "" {
		s.AddFilter(FilterData(raw))
	}

	return nil
}

// Exec executes the search provider and fills/scans
// the provided `items` slice with the found models.
func (s *Provider) Exec(items any) (*Result, error) {
	if s.query == nil {
		return nil, ErrEmptyQuery
	}

	// shallow clone the provider's query
	modelsQuery := *s.query

	// build filters
	for _, f := range s.filter {
		if len(f) > MaxFilterLength {
			return nil, ErrFilterLengthLimit
		}
		expr, err := f.BuildExprWithLimit(s.fieldResolver, s.maxFilterExprLimit)
		if err != nil {
			return nil, err
		}
		if expr != nil {
			modelsQuery.AndWhere(expr)
		}
	}

	// apply sorting
	if len(s.sort) > s.maxSortExprLimit {
		return nil, ErrSortExprLimit
	}
	for _, sortField := range s.sort {
		if len(sortField.Name) > MaxSortFieldLength {
			return nil, ErrSortFieldLengthLimit
		}
		expr, err := sortField.BuildExpr(s.fieldResolver)
		if err != nil {
			return nil, err
		}
		if expr != "" {
			// ensure that _rowid_ expressions are always prefixed with the first FROM table
			if sortField.Name == rowidSortKey && !strings.Contains(expr, ".") {
				queryInfo := modelsQuery.Info()
				if len(queryInfo.From) > 0 {
					expr = "[[" + inflector.Columnify(queryInfo.From[0]) + "]]." + expr
				}
			}

			modelsQuery.AndOrderBy(expr)
		}
	}

	// apply field resolver query modifications (if any)
	if err := s.fieldResolver.UpdateQuery(&modelsQuery); err != nil {
		return nil, err
	}

	// normalize page
	if s.page <= 0 {
		s.page = 1
	}

	// normalize perPage
	if s.perPage <= 0 {
		s.perPage = DefaultPerPage
	} else if s.perPage > MaxPerPage {
		s.perPage = MaxPerPage
	}

	// negative value to differentiate from the zero default
	totalCount := -1
	totalPages := -1

	// prepare a count query from the base one
	countQuery := modelsQuery // shallow clone
	countExec := func() error {
		queryInfo := countQuery.Info()
		countCol := s.countCol
		if len(queryInfo.From) > 0 {
			countCol = queryInfo.From[0] + "." + countCol
		}

		// note: countQuery is shallow cloned and slice/map in-place modifications should be avoided
		err := countQuery.Distinct(false).
			Select("COUNT(DISTINCT [[" + countCol + "]])").
			OrderBy( /* reset */ ).
			Row(&totalCount)
		if err != nil {
			return err
		}

		totalPages = int(math.Ceil(float64(totalCount) / float64(s.perPage)))

		return nil
	}

	// apply pagination to the original query and fetch the models
	modelsExec := func() error {
		modelsQuery.Limit(int64(s.perPage))
		modelsQuery.Offset(int64(s.perPage * (s.page - 1)))

		return modelsQuery.All(items)
	}

	if !s.skipTotal {
		// execute the 2 queries concurrently
		errg := new(errgroup.Group)
		errg.SetLimit(2)
		errg.Go(countExec)
		errg.Go(modelsExec)
		if err := errg.Wait(); err != nil {
			return nil, err
		}
	} else {
		if err := modelsExec(); err != nil {
			return nil, err
		}
	}

	result := &Result{
		Page:       s.page,
		PerPage:    s.perPage,
		TotalItems: totalCount,
		TotalPages: totalPages,
		Items:      items,
	}

	return result, nil
}

// ParseAndExec is a short convenient method to trigger both
// `Parse()` and `Exec()` in a single call.
func (s *Provider) ParseAndExec(urlQuery string, modelsSlice any) (*Result, error) {
	if err := s.Parse(urlQuery); err != nil {
		return nil, err
	}

	return s.Exec(modelsSlice)
}
