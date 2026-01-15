package search

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
)

type NullFallbackPreference int

const (
	NullFallbackAuto     NullFallbackPreference = 0
	NullFallbackDisabled NullFallbackPreference = 1
	NullFallbackEnforced NullFallbackPreference = 2
)

// ResolverResult defines a single FieldResolver.Resolve() successfully parsed result.
type ResolverResult struct {
	// Identifier is the plain SQL identifier/column that will be used
	// in the final db expression as left or right operand.
	Identifier string

	// NullFallback specify the preference for how NULL or empty values
	// should be resolved (default to "auto").
	//
	// Set to NullFallbackDisabled to prevent any COALESCE or NULL fallbacks.
	// Set to NullFallbackEnforced to prefer COALESCE or NULL fallbacks when needed.
	NullFallback NullFallbackPreference

	// Params is a map with db placeholder->value pairs that will be added
	// to the query when building both resolved operands/sides in a single expression.
	Params dbx.Params

	// MultiMatchSubQuery is an optional sub query expression that will be added
	// in addition to the combined ResolverResult expression during build.
	MultiMatchSubQuery *MultiMatchSubquery

	// AfterBuild is an optional function that will be called after building
	// and combining the result of both resolved operands/sides in a single expression.
	AfterBuild func(expr dbx.Expression) dbx.Expression
}

// FieldResolver defines an interface for managing search fields.
type FieldResolver interface {
	// UpdateQuery allows to updated the provided db query based on the
	// resolved search fields (eg. adding joins aliases, etc.).
	//
	// Called internally by `search.Provider` before executing the search request.
	UpdateQuery(query *dbx.SelectQuery) error

	// Resolve parses the provided field and returns a properly
	// formatted db identifier (eg. NULL, quoted column, placeholder parameter, etc.).
	Resolve(field string) (*ResolverResult, error)
}

// NewSimpleFieldResolver creates a new `SimpleFieldResolver` with the
// provided `allowedFields`.
//
// Each `allowedFields` could be a plain string (eg. "name")
// or a regexp pattern (eg. `^\w+[\w\.]*$`).
func NewSimpleFieldResolver(allowedFields ...string) *SimpleFieldResolver {
	return &SimpleFieldResolver{
		allowedFields: allowedFields,
	}
}

// SimpleFieldResolver defines a generic search resolver that allows
// only its listed fields to be resolved and take part in a search query.
//
// If `allowedFields` are empty no fields filtering is applied.
type SimpleFieldResolver struct {
	allowedFields []string
}

// UpdateQuery implements `search.UpdateQuery` interface.
func (r *SimpleFieldResolver) UpdateQuery(query *dbx.SelectQuery) error {
	// nothing to update...
	return nil
}

// Resolve implements `search.Resolve` interface.
//
// Returns error if `field` is not in `r.allowedFields`.
func (r *SimpleFieldResolver) Resolve(field string) (*ResolverResult, error) {
	if !list.ExistInSliceWithRegex(field, r.allowedFields) {
		return nil, fmt.Errorf("failed to resolve field %q", field)
	}

	parts := strings.Split(field, ".")

	// single regular field
	if len(parts) == 1 {
		return &ResolverResult{
			Identifier: "[[" + inflector.Columnify(parts[0]) + "]]",
		}, nil
	}

	// treat as json path
	var jsonPath strings.Builder
	jsonPath.WriteString("$")
	for _, part := range parts[1:] {
		if _, err := strconv.Atoi(part); err == nil {
			jsonPath.WriteString("[")
			jsonPath.WriteString(inflector.Columnify(part))
			jsonPath.WriteString("]")
		} else {
			jsonPath.WriteString(".")
			jsonPath.WriteString(inflector.Columnify(part))
		}
	}

	return &ResolverResult{
		NullFallback: NullFallbackDisabled,
		Identifier: fmt.Sprintf(
			"JSON_EXTRACT([[%s]], '%s')",
			inflector.Columnify(parts[0]),
			jsonPath.String(),
		),
	}, nil
}
