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
				NullFallback: NullFallbackDisabled,
				Identifier:   `(6371 * acos(cos(radians({:latA})) * cos(radians({:latB})) * cos(radians({:lonB}) - radians({:lonA})) + sin(radians({:latA})) * sin(radians({:latB}))))`,
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
				NullFallback: NullFallbackDisabled,
				Identifier:   `(6371 * acos(cos(radians({:latA})) * cos(radians({:latB})) * cos(radians({:lonB}) - radians({:lonA})) + sin(radians({:latA})) * sin(radians({:latB}))))`,
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

func TestTokenFunctionsStrftime(t *testing.T) {
	t.Parallel()

	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	fn, ok := TokenFunctions["strftime"]
	if !ok {
		t.Error("Expected strftime token function to be registered.")
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

		// format arg
		// -----------------------------------------------------------
		{
			"(format arg) invalid token type function",
			[]fexpr.Token{
				{Literal: "abc", Type: fexpr.TokenFunction},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(format arg) invalid token type ws",
			[]fexpr.Token{
				{Literal: "abc", Type: fexpr.TokenWS},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(format arg) invalid token type number",
			[]fexpr.Token{
				{Literal: "abc", Type: fexpr.TokenNumber},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(format arg) invalid token type identifier",
			[]fexpr.Token{
				{Literal: "abc", Type: fexpr.TokenIdentifier},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(format arg) valid token type text",
			[]fexpr.Token{
				{Literal: "abc", Type: fexpr.TokenText},
			},
			baseTokenResolver,
			&ResolverResult{
				NullFallback: NullFallbackEnforced,
				Identifier:   `strftime({:format})`,
				Params:       map[string]any{"format": "abc"},
			},
			false,
		},

		// format + time-value args
		// -----------------------------------------------------------
		{
			"(format arg) invalid token type function",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText},
				{Literal: "2", Type: fexpr.TokenFunction},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(format arg) invalid token type ws",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText},
				{Literal: "2", Type: fexpr.TokenWS},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(format arg) valid token type number",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText},
				{Literal: "2", Type: fexpr.TokenNumber},
			},
			baseTokenResolver,
			&ResolverResult{
				NullFallback: NullFallbackEnforced,
				Identifier:   `strftime({:format},{:time})`,
				Params:       map[string]any{"format": "1", "time": "2"},
			},
			false,
		},
		{
			"(format arg) valid token type identifier",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText},
				{Literal: "2", Type: fexpr.TokenIdentifier},
			},
			baseTokenResolver,
			&ResolverResult{
				NullFallback: NullFallbackEnforced,
				Identifier:   `strftime({:format},{:time})`,
				Params:       map[string]any{"format": "1", "time": "2"},
			},
			false,
		},
		{
			"(format arg) valid token type text",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText},
				{Literal: "2", Type: fexpr.TokenText},
			},
			baseTokenResolver,
			&ResolverResult{
				NullFallback: NullFallbackEnforced,
				Identifier:   `strftime({:format},{:time})`,
				Params:       map[string]any{"format": "1", "time": "2"},
			},
			false,
		},

		// format + time-value + modifier args
		// -----------------------------------------------------------
		{
			"(modifiers arg) invalid token type function",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText}, // valid format
				{Literal: "2", Type: fexpr.TokenText}, // valid time-value
				{Literal: "3", Type: fexpr.TokenText}, // valid modifier
				{Literal: "4", Type: fexpr.TokenFunction},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(modifiers arg) invalid token type ws",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText}, // valid format
				{Literal: "2", Type: fexpr.TokenText}, // valid time-value
				{Literal: "3", Type: fexpr.TokenText}, // valid modifier
				{Literal: "4", Type: fexpr.TokenWS},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(modifiers arg) valid token type number",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText}, // valid format
				{Literal: "2", Type: fexpr.TokenText}, // valid time-value
				{Literal: "3", Type: fexpr.TokenText}, // valid modifier
				{Literal: "4", Type: fexpr.TokenNumber},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(modifiers arg) valid token type identifier",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText}, // valid format
				{Literal: "2", Type: fexpr.TokenText}, // valid time-value
				{Literal: "3", Type: fexpr.TokenText}, // valid modifier
				{Literal: "4", Type: fexpr.TokenIdentifier},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"(modifiers arg) valid token type text",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText}, // valid format
				{Literal: "2", Type: fexpr.TokenText}, // valid time-value
				{Literal: "3", Type: fexpr.TokenText}, // valid modifier
				{Literal: "4", Type: fexpr.TokenText}, // valid modifier
			},
			baseTokenResolver,
			&ResolverResult{
				NullFallback: NullFallbackEnforced,
				Identifier:   `strftime({:format},{:time},{:m1},{:m2})`,
				Params:       map[string]any{"format": "1", "time": "2", "m1": "3", "m2": "4"},
			},
			false,
		},

		// -----------------------------------------------------------

		{
			"= 10 args limit",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText},
				{Literal: "2", Type: fexpr.TokenText},
				{Literal: "3", Type: fexpr.TokenText},
				{Literal: "4", Type: fexpr.TokenText},
				{Literal: "5", Type: fexpr.TokenText},
				{Literal: "6", Type: fexpr.TokenText},
				{Literal: "7", Type: fexpr.TokenText},
				{Literal: "8", Type: fexpr.TokenText},
				{Literal: "9", Type: fexpr.TokenText},
				{Literal: "10", Type: fexpr.TokenText},
			},
			baseTokenResolver,
			&ResolverResult{
				NullFallback: NullFallbackEnforced,
				Identifier:   `strftime({:format},{:time},{:m1},{:m2},{:m3},{:m4},{:m5},{:m6},{:m7},{:m8})`,
				Params: map[string]any{
					"format": "1",
					"time":   "2",
					"m1":     "3",
					"m2":     "4",
					"m3":     "5",
					"m4":     "6",
					"m5":     "7",
					"m6":     "8",
					"m7":     "9",
					"m8":     "10",
				},
			},
			false,
		},
		{
			"> 10 args limit",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText},
				{Literal: "2", Type: fexpr.TokenText},
				{Literal: "3", Type: fexpr.TokenText},
				{Literal: "4", Type: fexpr.TokenText},
				{Literal: "5", Type: fexpr.TokenText},
				{Literal: "6", Type: fexpr.TokenText},
				{Literal: "7", Type: fexpr.TokenText},
				{Literal: "8", Type: fexpr.TokenText},
				{Literal: "9", Type: fexpr.TokenText},
				{Literal: "10", Type: fexpr.TokenText},
				{Literal: "11", Type: fexpr.TokenText},
			},
			baseTokenResolver,
			nil,
			true,
		},
		{
			"valid arguments but with resolver error",
			[]fexpr.Token{
				{Literal: "1", Type: fexpr.TokenText},
				{Literal: "2", Type: fexpr.TokenText},
				{Literal: "3", Type: fexpr.TokenText},
			},
			func(t fexpr.Token) (*ResolverResult, error) {
				return nil, errors.New("test")
			},
			nil,
			true,
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

func TestTokenFunctionsStrftimeExec(t *testing.T) {
	t.Parallel()

	testDB, err := createTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	fn, ok := TokenFunctions["strftime"]
	if !ok {
		t.Error("Expected strftime token function to be registered.")
	}

	result, err := fn(
		func(t fexpr.Token) (*ResolverResult, error) {
			placeholder := "t" + security.PseudorandomString(5)
			return &ResolverResult{Identifier: "{:" + placeholder + "}", Params: map[string]any{placeholder: t.Literal}}, nil
		},
		fexpr.Token{Literal: "%Y-%m", Type: fexpr.TokenText},
		fexpr.Token{Literal: "2026-01-02 01:02:03.456Z", Type: fexpr.TokenText},
		fexpr.Token{Literal: "+1 years", Type: fexpr.TokenText},
		fexpr.Token{Literal: "+5 months", Type: fexpr.TokenText},
	)
	if err != nil {
		t.Fatal(err)
	}

	column := []string{}
	err = testDB.NewQuery("select " + result.Identifier).Bind(result.Params).Column(&column)
	if err != nil {
		t.Fatal(err)
	}

	if len(column) != 1 {
		t.Fatalf("Expected exactly 1 column value as result, got %v", column)
	}

	expected := "2027-06"
	if column[0] != expected {
		t.Fatalf("Expected date value %s, got %s", expected, column[0])
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

	if a.NullFallback != b.NullFallback {
		t.Fatalf("Expected NullFallback to match, got %v vs %v", a.NullFallback, b.NullFallback)
	}

	if len(a.Params) != len(b.Params) {
		t.Fatalf("Expected equal number of params, got %v vs %v", len(a.Params), len(b.Params))
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
