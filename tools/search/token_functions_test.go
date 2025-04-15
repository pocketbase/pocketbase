package search

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ganigeorgiev/fexpr"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/security"
)

func TestTokenFunctionsGeoDistance(t *testing.T) {
	t.Parallel()

	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	fn, ok := TokenFunctions["geoDistance"]
	if !ok {
		t.Error("Expected geoDistance token function to be registered.")
	}

	baseTokenResolver := func(t fexpr.Token) (*ResolverResult, error) {
		placeholder := "t" + security.PseudorandomString(5)
		return &ResolverResult{Identifier: "{:" + placeholder + "}", Params: map[string]any{placeholder: t.Literal}}, nil
	}

	scenarios := []struct {
		name      string
		args      []fexpr.Token
		resolver  func(t fexpr.Token) (*ResolverResult, error)
		result    *ResolverResult
		expectErr bool
	}{
		{
			"no args",
			nil,
			baseTokenResolver,
			nil,
			true,
		},
		{
			"< 4 args",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenNumber},
				{Literal: "2", Type: fexpr.TokenNumber},
				{Literal: "3", Type: fexpr.TokenNumber},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"> 4 args",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenNumber},
				{Literal: "2", Type: fexpr.TokenNumber},
				{Literal: "3", Type: fexpr.TokenNumber},
				{Literal: "4", Type: fexpr.TokenNumber},
				{Literal: "5", Type: fexpr.TokenNumber},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"unsupported function argument",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenFunction},
				{Literal: "2", Type: fexpr.TokenNumber},
				{Literal: "3", Type: fexpr.TokenNumber},
				{Literal: "4", Type: fexpr.TokenNumber},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"unsupported text argument",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText},
				{Literal: "2", Type: fexpr.TokenNumber},
				{Literal: "3", Type: fexpr.TokenNumber},
				{Literal: "4", Type: fexpr.TokenNumber},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"4 valid arguments but with resolver error",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenNumber},
				{Literal: "2", Type: fexpr.TokenNumber},
				{Literal: "3", Type: fexpr.TokenNumber},
				{Literal: "4", Type: fexpr.TokenNumber},
			},
			func(t fexpr.Token) (*ResolverResult, error) {
				return nil, errors.New("test")
			},
			nil,
			true,
		},
		{
			"4 valid arguments",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenNumber},
				{Literal: "2", Type: fexpr.TokenNumber},
				{Literal: "3", Type: fexpr.TokenNumber},
				{Literal: "4", Type: fexpr.TokenNumber},
			},
			baseTokenResolver,
			&ResolverResult{
				NoCoalesce: true,
				Identifier: `(6371 * acos(cos(radians({:latA})) * cos(radians({:latB})) * cos(radians({:lonB}) - radians({:lonA})) + sin(radians({:latA})) * sin(radians({:latB}))))`,
				Params: map[string]any{
					"lonA": 1,
					"latA": 2,
					"lonB": 3,
					"latB": 4,
				},
			},
			false,
		},
		{
			"mixed arguments",
			[]fexpr.Token{
				{Literal: "null", Type: fexpr.TokenIdentifier},
				{Literal: "2", Type: fexpr.TokenNumber},
				{Literal: "false", Type: fexpr.TokenIdentifier},
				{Literal: "4", Type: fexpr.TokenNumber},
			},
			baseTokenResolver,
			&ResolverResult{
				NoCoalesce: true,
				Identifier: `(6371 * acos(cos(radians({:latA})) * cos(radians({:latB})) * cos(radians({:lonB}) - radians({:lonA})) + sin(radians({:latA})) * sin(radians({:latB}))))`,
				Params: map[string]any{
					"lonA": "null",
					"latA": 2,
					"lonB": false,
					"latB": 4,
				},
			},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result, err := fn(s.resolver, s.args...)

			hasErr := err != nil
			if hasErr != s.expectErr {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectErr, hasErr, err)
			}

			testCompareResults(t, s.result, result)
		})
	}
}

func TestTokenFunctionsGeoDistanceExec(t *testing.T) {
	t.Parallel()

	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	fn, ok := TokenFunctions["geoDistance"]
	if !ok {
		t.Error("Expected geoDistance token function to be registered.")
	}

	result, err := fn(
		func(t fexpr.Token) (*ResolverResult, error) {
			placeholder := "t" + security.PseudorandomString(5)
			return &ResolverResult{Identifier: "{:" + placeholder + "}", Params: map[string]any{placeholder: t.Literal}}, nil
		},
		fexpr.Token{Literal: "23.23033854945808", Type: fexpr.TokenNumber},
		fexpr.Token{Literal: "42.713146090563384", Type: fexpr.TokenNumber},
		fexpr.Token{Literal: "23.44920680886216", Type: fexpr.TokenNumber},
		fexpr.Token{Literal: "42.7078484153991", Type: fexpr.TokenNumber},
	)
	if err != nil {
		t.Fatal(err)
	}

	column := []float64{}
	err = testDB.NewQuery("select " + result.Identifier).Bind(result.Params).Column(&column)
	if err != nil {
		t.Fatal(err)
	}

	if len(column) != 1 {
		t.Fatalf("Expected exactly 1 column value as result, got %v", column)
	}

	expected := "17.89"
	distance := fmt.Sprintf("%.2f", column[0])
	if distance != expected {
		t.Fatalf("Expected distance value %s, got %s", expected, distance)
	}
}

// -------------------------------------------------------------------

func testCompareResults(t *testing.T, a, b *ResolverResult) {
	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	aIsNil := a == nil
	bIsNil := b == nil
	if aIsNil != bIsNil {
		t.Fatalf("Expected aIsNil and bIsNil to be the same, got %v vs %v", aIsNil, bIsNil)
	}

	if aIsNil && bIsNil {
		return
	}

	aHasAfterBuild := a.AfterBuild == nil
	bHasAfterBuild := b.AfterBuild == nil
	if aHasAfterBuild != bHasAfterBuild {
		t.Fatalf("Expected aHasAfterBuild and bHasAfterBuild to be the same, got %v vs %v", aHasAfterBuild, bHasAfterBuild)
	}

	var aAfterBuild string
	if a.AfterBuild != nil {
		aAfterBuild = a.AfterBuild(dbx.NewExp("test")).Build(testDB.DB, a.Params)
	}
	var bAfterBuild string
	if b.AfterBuild != nil {
		bAfterBuild = b.AfterBuild(dbx.NewExp("test")).Build(testDB.DB, a.Params)
	}
	if aAfterBuild != bAfterBuild {
		t.Fatalf("Expected bAfterBuild and bAfterBuild to be the same, got\n%s\nvs\n%s", aAfterBuild, bAfterBuild)
	}

	var aMultiMatchSubQuery string
	if a.MultiMatchSubQuery != nil {
		aMultiMatchSubQuery = a.MultiMatchSubQuery.Build(testDB.DB, a.Params)
	}
	var bMultiMatchSubQuery string
	if b.MultiMatchSubQuery != nil {
		bMultiMatchSubQuery = b.MultiMatchSubQuery.Build(testDB.DB, b.Params)
	}
	if aMultiMatchSubQuery != bMultiMatchSubQuery {
		t.Fatalf("Expected bMultiMatchSubQuery and bMultiMatchSubQuery to be the same, got\n%s\nvs\n%s", aMultiMatchSubQuery, bMultiMatchSubQuery)
	}

	if a.NoCoalesce != b.NoCoalesce {
		t.Fatalf("Expected NoCoalesce to match, got %v vs %v", a.NoCoalesce, b.NoCoalesce)
	}

	// loose placeholders replacement
	var aResolved = a.Identifier
	for k, v := range a.Params {
		aResolved = strings.ReplaceAll(aResolved, "{:"+k+"}", fmt.Sprintf("%v", v))
	}
	var bResolved = b.Identifier
	for k, v := range b.Params {
		bResolved = strings.ReplaceAll(bResolved, "{:"+k+"}", fmt.Sprintf("%v", v))
	}
	if aResolved != bResolved {
		t.Fatalf("Expected resolved identifiers to match, got\n%s\nvs\n%s", aResolved, bResolved)
	}
}
