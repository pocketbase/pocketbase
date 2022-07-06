package search

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
)

// FieldResolver defines an interface for managing search fields.
type FieldResolver interface {
	// UpdateQuery allows to updated the provided db query based on the
	// resolved search fields (eg. adding joins aliases, etc.).
	//
	// Called internally by `search.Provider` before executing the search request.
	UpdateQuery(query *dbx.SelectQuery) error

	// Resolve parses the provided field and returns a properly
	// formatted db identifier (eg. NULL, quoted column, placeholder parameter, etc.).
	Resolve(field string) (name string, placeholderParams dbx.Params, err error)
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
func (r *SimpleFieldResolver) Resolve(field string) (resultName string, placeholderParams dbx.Params, err error) {
	if !list.ExistInSliceWithRegex(field, r.allowedFields) {
		return "", nil, fmt.Errorf("Failed to resolve field %q.", field)
	}

	return fmt.Sprintf("[[%s]]", inflector.Columnify(field)), nil, nil
}
