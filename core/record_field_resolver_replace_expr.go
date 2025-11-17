package core

import (
	"strings"

	"github.com/pocketbase/dbx"
)

var _ dbx.Expression = (*replaceWithExpression)(nil)

// replaceWithExpression defines a custom expression that will replace
// a placeholder identifier found in "old" with the result of "new".
type replaceWithExpression struct {
	placeholder string
	old         dbx.Expression
	new         dbx.Expression
}

// Build converts the expression into a SQL fragment.
//
// Implements [dbx.Expression] interface.
func (e *replaceWithExpression) Build(db *dbx.DB, params dbx.Params) string {
	if e.placeholder == "" || e.old == nil || e.new == nil {
		return "0=1"
	}

	oldResult := e.old.Build(db, params)
	newResult := e.new.Build(db, params)

	return strings.ReplaceAll(oldResult, e.placeholder, newResult)
}
