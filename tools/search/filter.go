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

	var result dbx.Expression

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

		if group.Join == fexpr.JoinAnd {
			result = dbx.And(result, expr)
		} else {
			result = dbx.Or(result, expr)
		}
	}

	return result, nil
}

func (f FilterData) resolveTokenizedExpr(expr fexpr.Expr, fieldResolver FieldResolver) (dbx.Expression, error) {
	lName, lParams, lErr := f.resolveToken(expr.Left, fieldResolver)
	if lName == "" || lErr != nil {
		return nil, fmt.Errorf("Invalid left operand %q - %v.", expr.Left.Literal, lErr)
	}

	rName, rParams, rErr := f.resolveToken(expr.Right, fieldResolver)
	if rName == "" || rErr != nil {
		return nil, fmt.Errorf("Invalid right operand %q - %v.", expr.Right.Literal, rErr)
	}

	switch expr.Op {
	case fexpr.SignEq:
		return dbx.NewExp(fmt.Sprintf("COALESCE(%s, '') = COALESCE(%s, '')", lName, rName), mergeParams(lParams, rParams)), nil
	case fexpr.SignNeq:
		return dbx.NewExp(fmt.Sprintf("COALESCE(%s, '') != COALESCE(%s, '')", lName, rName), mergeParams(lParams, rParams)), nil
	case fexpr.SignLike:
		// the right side is a column and therefor wrap it with "%" for contains like behavior
		if len(rParams) == 0 {
			return dbx.NewExp(fmt.Sprintf("%s LIKE ('%%' || %s || '%%')", lName, rName), lParams), nil
		}

		return dbx.NewExp(fmt.Sprintf("%s LIKE %s", lName, rName), mergeParams(lParams, wrapLikeParams(rParams))), nil
	case fexpr.SignNlike:
		// the right side is a column and therefor wrap it with "%" for not-contains like behavior
		if len(rParams) == 0 {
			return dbx.NewExp(fmt.Sprintf("%s NOT LIKE ('%%' || %s || '%%')", lName, rName), lParams), nil
		}

		// normalize operands and switch sides if the left operand is a number/text, but the right one is a column
		// (usually this shouldn't be needed, but it's kept for backward compatibility)
		if len(lParams) > 0 && len(rParams) == 0 {
			return dbx.NewExp(fmt.Sprintf("%s NOT LIKE %s", rName, lName), wrapLikeParams(lParams)), nil
		}

		return dbx.NewExp(fmt.Sprintf("%s NOT LIKE %s", lName, rName), mergeParams(lParams, wrapLikeParams(rParams))), nil
	case fexpr.SignLt:
		return dbx.NewExp(fmt.Sprintf("%s < %s", lName, rName), mergeParams(lParams, rParams)), nil
	case fexpr.SignLte:
		return dbx.NewExp(fmt.Sprintf("%s <= %s", lName, rName), mergeParams(lParams, rParams)), nil
	case fexpr.SignGt:
		return dbx.NewExp(fmt.Sprintf("%s > %s", lName, rName), mergeParams(lParams, rParams)), nil
	case fexpr.SignGte:
		return dbx.NewExp(fmt.Sprintf("%s >= %s", lName, rName), mergeParams(lParams, rParams)), nil
	}

	return nil, fmt.Errorf("Unknown expression operator %q", expr.Op)
}

func (f FilterData) resolveToken(token fexpr.Token, fieldResolver FieldResolver) (name string, params dbx.Params, err error) {
	switch token.Type {
	case fexpr.TokenIdentifier:
		// current datetime constant
		// ---
		if token.Literal == "@now" {
			placeholder := "t" + security.PseudorandomString(8)
			name := fmt.Sprintf("{:%s}", placeholder)
			params := dbx.Params{placeholder: types.NowDateTime().String()}

			return name, params, nil
		}

		// custom resolver
		// ---
		name, params, err := fieldResolver.Resolve(token.Literal)

		if name == "" || err != nil {
			m := map[string]string{
				// if `null` field is missing, treat `null` identifier as NULL token
				"null": "NULL",
				// if `true` field is missing, treat `true` identifier as TRUE token
				"true": "1",
				// if `false` field is missing, treat `false` identifier as FALSE token
				"false": "0",
			}
			if v, ok := m[strings.ToLower(token.Literal)]; ok {
				return v, nil, nil
			}
			return "", nil, err
		}

		return name, params, err
	case fexpr.TokenText:
		placeholder := "t" + security.PseudorandomString(8)
		name := fmt.Sprintf("{:%s}", placeholder)
		params := dbx.Params{placeholder: token.Literal}

		return name, params, nil
	case fexpr.TokenNumber:
		placeholder := "t" + security.PseudorandomString(8)
		name := fmt.Sprintf("{:%s}", placeholder)
		params := dbx.Params{placeholder: cast.ToFloat64(token.Literal)}

		return name, params, nil
	}

	return "", nil, errors.New("Unresolvable token type.")
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
