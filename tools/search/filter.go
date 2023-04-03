package search

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ganigeorgiev/fexpr"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

// FilterData is a filter expression string following the `fexpr` package grammar.
//
// Example:
//
//	var filter FilterData = "id = null || (name = 'test' && status = true)"
//	resolver := search.NewSimpleFieldResolver("id", "name", "status")
//	expr, err := filter.BuildExpr(resolver)
type FilterData string

// parsedFilterData holds a cache with previously parsed filter data expressions
// (initialized with some preallocated empty data map)
var parsedFilterData = store.New(make(map[string][]fexpr.ExprGroup, 50))

// BuildExpr parses the current filter data and returns a new db WHERE expression.
func (f FilterData) BuildExpr(fieldResolver FieldResolver) (dbx.Expression, error) {
	raw := string(f)
	if parsedFilterData.Has(raw) {
		return f.build(parsedFilterData.Get(raw), fieldResolver)
	}
	data, err := fexpr.Parse(raw)
	if err != nil {
		return nil, err
	}
	// store in cache
	// (the limit size is arbitrary and it is there to prevent the cache growing too big)
	parsedFilterData.SetIfLessThanLimit(raw, data, 500)
	return f.build(data, fieldResolver)
}

func (f FilterData) build(data []fexpr.ExprGroup, fieldResolver FieldResolver) (dbx.Expression, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty filter expression.")
	}

	result := &concatExpr{separator: " "}

	for _, group := range data {
		var expr dbx.Expression
		var exprErr error

		switch item := group.Item.(type) {
		case fexpr.Expr:
			expr, exprErr = f.resolveTokenizedExpr(item, fieldResolver)
		case fexpr.ExprGroup:
			expr, exprErr = f.build([]fexpr.ExprGroup{item}, fieldResolver)
		case []fexpr.ExprGroup:
			expr, exprErr = f.build(item, fieldResolver)
		default:
			exprErr = errors.New("Unsupported expression item.")
		}

		if exprErr != nil {
			return nil, exprErr
		}

		if len(result.parts) > 0 {
			var op string
			if group.Join == fexpr.JoinOr {
				op = "OR"
			} else {
				op = "AND"
			}
			result.parts = append(result.parts, &opExpr{op})
		}

		result.parts = append(result.parts, expr)
	}

	return result, nil
}

func (f FilterData) resolveTokenizedExpr(expr fexpr.Expr, fieldResolver FieldResolver) (dbx.Expression, error) {
	lResult, lErr := resolveToken(expr.Left, fieldResolver)
	if lErr != nil || lResult.Identifier == "" {
		return nil, fmt.Errorf("invalid left operand %q - %v", expr.Left.Literal, lErr)
	}

	rResult, rErr := resolveToken(expr.Right, fieldResolver)
	if rErr != nil || rResult.Identifier == "" {
		return nil, fmt.Errorf("invalid right operand %q - %v", expr.Right.Literal, rErr)
	}

	return buildExpr(lResult, expr.Op, rResult)
}

func buildExpr(
	left *ResolverResult,
	op fexpr.SignOp,
	right *ResolverResult,
) (dbx.Expression, error) {
	var expr dbx.Expression

	switch op {
	case fexpr.SignEq, fexpr.SignAnyEq:
		expr = resolveEqualExpr(true, left, right)
	case fexpr.SignNeq, fexpr.SignAnyNeq:
		expr = resolveEqualExpr(false, left, right)
	case fexpr.SignLike, fexpr.SignAnyLike:
		// the right side is a column and therefor wrap it with "%" for contains like behavior
		if len(right.Params) == 0 {
			expr = dbx.NewExp(fmt.Sprintf("%s LIKE ('%%' || %s || '%%') ESCAPE '\\'", left.Identifier, right.Identifier), left.Params)
		} else {
			expr = dbx.NewExp(fmt.Sprintf("%s LIKE %s ESCAPE '\\'", left.Identifier, right.Identifier), mergeParams(left.Params, wrapLikeParams(right.Params)))
		}
	case fexpr.SignNlike, fexpr.SignAnyNlike:
		// the right side is a column and therefor wrap it with "%" for not-contains like behavior
		if len(right.Params) == 0 {
			expr = dbx.NewExp(fmt.Sprintf("%s NOT LIKE ('%%' || %s || '%%') ESCAPE '\\'", left.Identifier, right.Identifier), left.Params)
		} else {
			expr = dbx.NewExp(fmt.Sprintf("%s NOT LIKE %s ESCAPE '\\'", left.Identifier, right.Identifier), mergeParams(left.Params, wrapLikeParams(right.Params)))
		}
	case fexpr.SignLt, fexpr.SignAnyLt:
		expr = dbx.NewExp(fmt.Sprintf("%s < %s", left.Identifier, right.Identifier), mergeParams(left.Params, right.Params))
	case fexpr.SignLte, fexpr.SignAnyLte:
		expr = dbx.NewExp(fmt.Sprintf("%s <= %s", left.Identifier, right.Identifier), mergeParams(left.Params, right.Params))
	case fexpr.SignGt, fexpr.SignAnyGt:
		expr = dbx.NewExp(fmt.Sprintf("%s > %s", left.Identifier, right.Identifier), mergeParams(left.Params, right.Params))
	case fexpr.SignGte, fexpr.SignAnyGte:
		expr = dbx.NewExp(fmt.Sprintf("%s >= %s", left.Identifier, right.Identifier), mergeParams(left.Params, right.Params))
	}

	if expr == nil {
		return nil, fmt.Errorf("unknown expression operator %q", op)
	}

	// multi-match expressions
	if !isAnyMatchOp(op) {
		if left.MultiMatchSubQuery != nil && right.MultiMatchSubQuery != nil {
			mm := &manyVsManyExpr{
				leftSubQuery:  left.MultiMatchSubQuery,
				rightSubQuery: right.MultiMatchSubQuery,
				op:            op,
			}

			expr = dbx.Enclose(dbx.And(expr, mm))
		} else if left.MultiMatchSubQuery != nil {
			mm := &manyVsOneExpr{
				subQuery:     left.MultiMatchSubQuery,
				op:           op,
				otherOperand: right,
			}

			expr = dbx.Enclose(dbx.And(expr, mm))
		} else if right.MultiMatchSubQuery != nil {
			mm := &manyVsOneExpr{
				subQuery:     right.MultiMatchSubQuery,
				op:           op,
				otherOperand: left,
				inverse:      true,
			}

			expr = dbx.Enclose(dbx.And(expr, mm))
		}
	}

	if left.AfterBuild != nil {
		expr = left.AfterBuild(expr)
	}

	if right.AfterBuild != nil {
		expr = right.AfterBuild(expr)
	}

	return expr, nil
}

func resolveToken(token fexpr.Token, fieldResolver FieldResolver) (*ResolverResult, error) {
	switch token.Type {
	case fexpr.TokenIdentifier:
		// current datetime constant
		// ---
		if token.Literal == "@now" {
			placeholder := "t" + security.PseudorandomString(5)

			return &ResolverResult{
				Identifier: "{:" + placeholder + "}",
				Params:     dbx.Params{placeholder: types.NowDateTime().String()},
			}, nil
		}

		// custom resolver
		// ---
		result, err := fieldResolver.Resolve(token.Literal)

		if err != nil || result.Identifier == "" {
			m := map[string]string{
				// if `null` field is missing, treat `null` identifier as NULL token
				"null": "NULL",
				// if `true` field is missing, treat `true` identifier as TRUE token
				"true": "1",
				// if `false` field is missing, treat `false` identifier as FALSE token
				"false": "0",
			}
			if v, ok := m[strings.ToLower(token.Literal)]; ok {
				return &ResolverResult{Identifier: v}, nil
			}
			return nil, err
		}

		return result, err
	case fexpr.TokenText:
		placeholder := "t" + security.PseudorandomString(5)

		return &ResolverResult{
			Identifier: "{:" + placeholder + "}",
			Params:     dbx.Params{placeholder: token.Literal},
		}, nil
	case fexpr.TokenNumber:
		placeholder := "t" + security.PseudorandomString(5)

		return &ResolverResult{
			Identifier: "{:" + placeholder + "}",
			Params:     dbx.Params{placeholder: cast.ToFloat64(token.Literal)},
		}, nil
	}

	return nil, errors.New("unresolvable token type")
}

// Resolves = and != expressions in an attempt to minimize the COALESCE
// usage and to gracefully handle null vs empty string normalizations.
//
// The expression `a = "" OR a is null` tends to perform better than
// `COALESCE(a, "") = ""` since the direct match can be accomplished
// with a seek while the COALESCE will induce a table scan.
func resolveEqualExpr(equal bool, left, right *ResolverResult) dbx.Expression {
	isLeftEmpty := isEmptyIdentifier(left) || (len(left.Params) == 1 && hasEmptyParamValue(left))
	isRightEmpty := isEmptyIdentifier(right) || (len(right.Params) == 1 && hasEmptyParamValue(right))

	equalOp := "="
	concatOp := "OR"
	nullExpr := "IS NULL"
	if !equal {
		equalOp = "!="
		concatOp = "AND"
		nullExpr = "IS NOT NULL"
	}

	// both operands are empty
	if isLeftEmpty && isRightEmpty {
		return dbx.NewExp(fmt.Sprintf("'' %s ''", equalOp), mergeParams(left.Params, right.Params))
	}

	// direct compare since at least one of the operands is known to be non-empty
	// eg. a = 'example'
	if isKnownNonEmptyIdentifier(left) || isKnownNonEmptyIdentifier(right) {
		leftIdentifier := left.Identifier
		if isLeftEmpty {
			leftIdentifier = "''"
		}
		rightIdentifier := right.Identifier
		if isRightEmpty {
			rightIdentifier = "''"
		}
		return dbx.NewExp(
			fmt.Sprintf("%s %s %s", leftIdentifier, equalOp, rightIdentifier),
			mergeParams(left.Params, right.Params),
		)
	}

	// "" = b OR b IS NULL
	// "" != b AND b IS NOT NULL
	if isLeftEmpty {
		return dbx.NewExp(
			fmt.Sprintf("('' %s %s %s %s %s)", equalOp, right.Identifier, concatOp, right.Identifier, nullExpr),
			mergeParams(left.Params, right.Params),
		)
	}

	// a = "" OR a IS NULL
	// a != "" AND a IS NOT NULL
	if isRightEmpty {
		return dbx.NewExp(
			fmt.Sprintf("(%s %s '' %s %s %s)", left.Identifier, equalOp, concatOp, left.Identifier, nullExpr),
			mergeParams(left.Params, right.Params),
		)
	}

	// fallback to a COALESCE comparison
	return dbx.NewExp(
		fmt.Sprintf(
			"COALESCE(%s, '') %s COALESCE(%s, '')",
			left.Identifier,
			equalOp,
			right.Identifier,
		),
		mergeParams(left.Params, right.Params),
	)
}

func hasEmptyParamValue(result *ResolverResult) bool {
	for _, p := range result.Params {
		switch v := p.(type) {
		case nil:
			return true
		case string:
			if v == "" {
				return true
			}
		}
	}

	return false
}

func isKnownNonEmptyIdentifier(result *ResolverResult) bool {
	switch strings.ToLower(result.Identifier) {
	case "1", "0", "false", `true`:
		return true
	}

	return len(result.Params) > 0 && !hasEmptyParamValue(result) && !isEmptyIdentifier(result)
}

func isEmptyIdentifier(result *ResolverResult) bool {
	switch strings.ToLower(result.Identifier) {
	case "", "null", "''", `""`, "``":
		return true
	default:
		return false
	}
}

func isAnyMatchOp(op fexpr.SignOp) bool {
	switch op {
	case
		fexpr.SignAnyEq,
		fexpr.SignAnyNeq,
		fexpr.SignAnyLike,
		fexpr.SignAnyNlike,
		fexpr.SignAnyLt,
		fexpr.SignAnyLte,
		fexpr.SignAnyGt,
		fexpr.SignAnyGte:
		return true
	}

	return false
}

// mergeParams returns new dbx.Params where each provided params item
// is merged in the order they are specified.
func mergeParams(params ...dbx.Params) dbx.Params {
	result := dbx.Params{}

	for _, p := range params {
		for k, v := range p {
			result[k] = v
		}
	}

	return result
}

// wrapLikeParams wraps each provided param value string with `%`
// if the string doesn't contains the `%` char (including its escape sequence).
func wrapLikeParams(params dbx.Params) dbx.Params {
	result := dbx.Params{}

	for k, v := range params {
		vStr := cast.ToString(v)
		if !strings.Contains(vStr, "%") {
			for i := 0; i < len(dbx.DefaultLikeEscape); i += 2 {
				vStr = strings.ReplaceAll(vStr, dbx.DefaultLikeEscape[i], dbx.DefaultLikeEscape[i+1])
			}
			vStr = "%" + vStr + "%"
		}
		result[k] = vStr
	}

	return result
}

// -------------------------------------------------------------------

var _ dbx.Expression = (*opExpr)(nil)

// opExpr defines an expression that contains a raw sql operator string.
type opExpr struct {
	op string
}

// Build converts the expression into a SQL fragment.
//
// Implements [dbx.Expression] interface.
func (e *opExpr) Build(db *dbx.DB, params dbx.Params) string {
	return e.op
}

// -------------------------------------------------------------------

var _ dbx.Expression = (*concatExpr)(nil)

// concatExpr defines an expression that concatenates multiple
// other expressions with a specified separator.
type concatExpr struct {
	parts     []dbx.Expression
	separator string
}

// Build converts the expression into a SQL fragment.
//
// Implements [dbx.Expression] interface.
func (e *concatExpr) Build(db *dbx.DB, params dbx.Params) string {
	if len(e.parts) == 0 {
		return ""
	}

	stringParts := make([]string, 0, len(e.parts))

	for _, p := range e.parts {
		if p == nil {
			continue
		}

		if sql := p.Build(db, params); sql != "" {
			stringParts = append(stringParts, sql)
		}
	}

	// skip extra parenthesis for single concat expression
	if len(stringParts) == 1 &&
		// check for already concatenated raw/plain expressions
		!strings.Contains(strings.ToUpper(stringParts[0]), " AND ") &&
		!strings.Contains(strings.ToUpper(stringParts[0]), " OR ") {
		return stringParts[0]
	}

	return "(" + strings.Join(stringParts, e.separator) + ")"
}

// -------------------------------------------------------------------

var _ dbx.Expression = (*manyVsManyExpr)(nil)

// manyVsManyExpr constructs a multi-match many<->many db where expression.
//
// Expects leftSubQuery and rightSubQuery to return a subquery with a
// single "multiMatchValue" column.
type manyVsManyExpr struct {
	leftSubQuery  dbx.Expression
	rightSubQuery dbx.Expression
	op            fexpr.SignOp
}

// Build converts the expression into a SQL fragment.
//
// Implements [dbx.Expression] interface.
func (e *manyVsManyExpr) Build(db *dbx.DB, params dbx.Params) string {
	if e.leftSubQuery == nil || e.rightSubQuery == nil {
		return "0=1"
	}

	lAlias := "__ml" + security.PseudorandomString(5)
	rAlias := "__mr" + security.PseudorandomString(5)

	whereExpr, buildErr := buildExpr(
		&ResolverResult{
			Identifier: "[[" + lAlias + ".multiMatchValue]]",
		},
		e.op,
		&ResolverResult{
			Identifier: "[[" + rAlias + ".multiMatchValue]]",
			// note: the AfterBuild needs to be handled only once and it
			// doesn't matter whether it is applied on the left or right subquery operand
			AfterBuild: multiMatchAfterBuildFunc(e.op, lAlias, rAlias),
		},
	)

	if buildErr != nil {
		return "0=1"
	}

	return fmt.Sprintf(
		"NOT EXISTS (SELECT 1 FROM (%s) {{%s}} LEFT JOIN (%s) {{%s}} WHERE %s)",
		e.leftSubQuery.Build(db, params),
		lAlias,
		e.rightSubQuery.Build(db, params),
		rAlias,
		whereExpr.Build(db, params),
	)
}

// -------------------------------------------------------------------

var _ dbx.Expression = (*manyVsOneExpr)(nil)

// manyVsManyExpr constructs a multi-match many<->one db where expression.
//
// Expects subQuery to return a subquery with a single "multiMatchValue" column.
//
// You can set inverse=false to reverse the condition sides (aka. one<->many).
type manyVsOneExpr struct {
	subQuery     dbx.Expression
	op           fexpr.SignOp
	otherOperand *ResolverResult
	inverse      bool
}

// Build converts the expression into a SQL fragment.
//
// Implements [dbx.Expression] interface.
func (e *manyVsOneExpr) Build(db *dbx.DB, params dbx.Params) string {
	if e.subQuery == nil {
		return "0=1"
	}

	alias := "__sm" + security.PseudorandomString(5)

	r1 := &ResolverResult{
		Identifier: "[[" + alias + ".multiMatchValue]]",
		AfterBuild: multiMatchAfterBuildFunc(e.op, alias),
	}

	r2 := &ResolverResult{
		Identifier: e.otherOperand.Identifier,
		Params:     e.otherOperand.Params,
	}

	var whereExpr dbx.Expression
	var buildErr error

	if e.inverse {
		whereExpr, buildErr = buildExpr(r2, e.op, r1)
	} else {
		whereExpr, buildErr = buildExpr(r1, e.op, r2)
	}

	if buildErr != nil {
		return "0=1"
	}

	return fmt.Sprintf(
		"NOT EXISTS (SELECT 1 FROM (%s) {{%s}} WHERE %s)",
		e.subQuery.Build(db, params),
		alias,
		whereExpr.Build(db, params),
	)
}

func multiMatchAfterBuildFunc(op fexpr.SignOp, multiMatchAliases ...string) func(dbx.Expression) dbx.Expression {
	return func(expr dbx.Expression) dbx.Expression {
		expr = dbx.Not(expr) // inverse for the not-exist expression

		if op == fexpr.SignEq {
			return expr
		}

		orExprs := make([]dbx.Expression, len(multiMatchAliases)+1)
		orExprs[0] = expr

		// Add an optional "IS NULL" condition(s) to handle the empty rows result.
		//
		// For example, let's assume that some "rel" field is [nonemptyRel1, nonemptyRel2, emptyRel3],
		// The filter "rel.total > 0" will ensures that the above will return true only if all relations
		// are existing and match the condition.
		//
		// The "=" operator is excluded because it will never equal directly with NULL anyway
		// and also because we want in case "rel.id = ''" is specified to allow
		// matching the empty relations (they will match due to the applied COALESCE).
		for i, mAlias := range multiMatchAliases {
			orExprs[i+1] = dbx.NewExp("[[" + mAlias + ".multiMatchValue]] IS NULL")
		}

		return dbx.Enclose(dbx.Or(orExprs...))
	}
}
