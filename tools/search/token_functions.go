package search

import (
	"fmt"

	"github.com/ganigeorgiev/fexpr"
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
			NoCoalesce: true,
			Identifier: `(6371 * acos(` +
				`cos(radians(` + latA + `)) * cos(radians(` + latB + `)) * ` +
				`cos(radians(` + lonB + `) - radians(` + lonA + `)) + ` +
				`sin(radians(` + latA + `)) * sin(radians(` + latB + `))` +
				`))`,
			Params: mergeParams(resolvedArgs[0].Params, resolvedArgs[1].Params, resolvedArgs[2].Params, resolvedArgs[3].Params),
		}, nil
	},
}
