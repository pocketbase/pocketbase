package search

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/ganigeorgiev/fexpr"
	"github.com/pocketbase/dbx"
)

var TokenFunctions = map[string]func(
	argTokenResolverFunc func(fexpr.Token) (*ResolverResult, error),
	args ...fexpr.Token,
) (*ResolverResult, error){
	// geoDistance(lonA, latA, lonB, latB) calculates the Haversine
	// distance between 2 points in kilometres (https://www.movable-type.co.uk/scripts/latlong.html).
	//
	// The accepted arguments at the moment could be either a plain number or a column identifier (including NULL).
	// If the column identifier cannot be resolved and converted to a numeric value, it resolves to NULL.
	//
	// Similar to the built-in SQLite functions, geoDistance doesn't apply
	// a "match-all" constraints in case there are multiple relation fields arguments.
	// Or in other words, if a collection has "orgs" multiple relation field pointing to "orgs" collection that has "office" as "geoPoint" field,
	// then the filter: `geoDistance(orgs.office.lon, orgs.office.lat, 1, 2) < 200`
	// will evaluate to true if for at-least-one of the "orgs.office" records the function result in a value satisfying the condition (aka. "result < 200").
	"geoDistance": func(argTokenResolverFunc func(fexpr.Token) (*ResolverResult, error), args ...fexpr.Token) (*ResolverResult, error) {
		if len(args) != 4 {
			return nil, fmt.Errorf("[geoDistance] expected 4 arguments, got %d", len(args))
		}

		resolvedArgs := make([]*ResolverResult, 4)
		for i, arg := range args {
			if arg.Type != fexpr.TokenIdentifier && arg.Type != fexpr.TokenNumber {
				return nil, fmt.Errorf("[geoDistance] argument %d must be an identifier or number", i)
			}
			resolved, err := argTokenResolverFunc(arg)
			if err != nil {
				return nil, fmt.Errorf("[geoDistance] failed to resolve argument %d: %w", i, err)
			}
			resolvedArgs[i] = resolved
		}

		lonA := resolvedArgs[0].Identifier
		latA := resolvedArgs[1].Identifier
		lonB := resolvedArgs[2].Identifier
		latB := resolvedArgs[3].Identifier

		return &ResolverResult{
			NullFallback: NullFallbackDisabled,
			Identifier: `(6371 * acos(` +
				`cos(radians(` + latA + `)) * cos(radians(` + latB + `)) * ` +
				`cos(radians(` + lonB + `) - radians(` + lonA + `)) + ` +
				`sin(radians(` + latA + `)) * sin(radians(` + latB + `))` +
				`))`,
			Params: mergeParams(resolvedArgs[0].Params, resolvedArgs[1].Params, resolvedArgs[2].Params, resolvedArgs[3].Params),
		}, nil
	},

	// strftime(format, [timeValue, modifier1, modifier2, ...]) returns
	// a date string formatted according to the specified format argument.
	//
	// It is similar to the builtin SQLite strftime function (https://sqlite.org/lang_datefunc.html)
	// with the main difference that NULL results will be normalized for
	// consistency with the non-nullable PocketBase "text" and "date" fields.
	//
	// The function accepts 1, 2 or 3+ arguments.
	//
	// (1) The first (format) argument must be always a formatting string
	// with valid substitutions as listed in https://sqlite.org/lang_datefunc.html.
	//
	// (2) The second (time-value) argument is optional and must be either a date string, number or collection field identifier
	// that matches one of the formats listed in https://sqlite.org/lang_datefunc.html#time_values.
	//
	// (3+) The remaining (modifiers) optional arguments are expected to be
	// string literals matching the listed modifiers in https://sqlite.org/lang_datefunc.html#modifiers.
	//
	// A multi-match constraint will be also applied in case the time-value
	// is an identifier as a result of a multi-value relation field.
	"strftime": func(argTokenResolverFunc func(fexpr.Token) (*ResolverResult, error), args ...fexpr.Token) (*ResolverResult, error) {
		totalArgs := len(args)

		if totalArgs < 1 {
			return nil, fmt.Errorf("[strftime] expected at least 1 arguments, got %d", len(args))
		}

		// limit the number of arguments to prevent abuse
		if totalArgs > 10 {
			return nil, fmt.Errorf("[strftime] too many arguments (max allowed 10, got %d)", totalArgs)
		}

		// format arg
		// -----------------------------------------------------------
		if args[0].Type != fexpr.TokenText {
			return nil, errors.New("[strftime] expects the first argument to be a format string")
		}

		formatArgResult, err := argTokenResolverFunc(args[0])
		if err != nil {
			return nil, fmt.Errorf("[strftime] failed to resolve format argument: %w", err)
		}

		// no further arguments
		if totalArgs == 1 {
			formatArgResult.NullFallback = NullFallbackEnforced
			formatArgResult.Identifier = "strftime(" + formatArgResult.Identifier + ")"
			return formatArgResult, nil
		}

		// time-value arg
		// -----------------------------------------------------------
		allowedTimeValueTokens := []fexpr.TokenType{fexpr.TokenText, fexpr.TokenIdentifier, fexpr.TokenNumber}
		if !slices.Contains(allowedTimeValueTokens, args[1].Type) {
			return nil, errors.New("[strftime] expects the second argument to be of a valid time-value type")
		}

		timeValueArgResult, err := argTokenResolverFunc(args[1])
		if err != nil {
			return nil, fmt.Errorf("[strftime] failed to resolve time-value argument: %w", err)
		}

		// modifiers args
		// -----------------------------------------------------------
		resolvedModifierArgs := make([]*ResolverResult, totalArgs-2)
		for i, arg := range args[2:] {
			if arg.Type != fexpr.TokenText {
				return nil, fmt.Errorf("[strftime] invalid modifier argument %d - can be only string", i)
			}

			resolved, err := argTokenResolverFunc(arg)
			if err != nil {
				return nil, fmt.Errorf("[strftime] failed to resolve modifier argument %d: %w", i, err)
			}

			resolvedModifierArgs[i] = resolved
		}

		// generating new ResolverResult
		// -----------------------------------------------------------
		result := &ResolverResult{
			NullFallback: NullFallbackEnforced,
			Params:       dbx.Params{},
		}

		identifiers := make([]string, 0, totalArgs)

		identifiers = append(identifiers, formatArgResult.Identifier)
		if err = concatUniqueParams(result.Params, formatArgResult.Params); err != nil {
			return nil, err
		}

		identifiers = append(identifiers, timeValueArgResult.Identifier)
		if err = concatUniqueParams(result.Params, timeValueArgResult.Params); err != nil {
			return nil, err
		}

		for _, m := range resolvedModifierArgs {
			identifiers = append(identifiers, m.Identifier)
			err = concatUniqueParams(result.Params, m.Params)
			if err != nil {
				return nil, err
			}
		}

		result.Identifier = "strftime(" + strings.Join(identifiers, ",") + ")"

		if timeValueArgResult.MultiMatchSubQuery != nil {
			// replace the regular time-value identifier with the multi-match one
			identifiers[1] = timeValueArgResult.MultiMatchSubQuery.ValueIdentifier
			result.MultiMatchSubQuery = timeValueArgResult.MultiMatchSubQuery
			result.MultiMatchSubQuery.ValueIdentifier = "strftime(" + strings.Join(identifiers, ",") + ")"

			err = concatUniqueParams(result.MultiMatchSubQuery.Params, result.Params)
			if err != nil {
				return nil, err
			}
		}

		return result, nil
	},
}

func concatUniqueParams(destParams, newParams dbx.Params) error {
	for k, v := range newParams {
		found, ok := destParams[k]
		if ok && v != found {
			return fmt.Errorf("conflicting param key %s", k)
		}

		destParams[k] = v
	}

	return nil
}
