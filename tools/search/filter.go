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

	// merge both operands parameters (if any)
	params := dbx.Params{}
	for k, v := range lParams {
		params[k] = v
	}
	for k, v := range rParams {
		params[k] = v
	}

	switch expr.Op {
	case fexpr.SignEq:
		return dbx.NewExp(fmt.Sprintf("COALESCE(%s, '') = COALESCE(%s, '')", lName, rName), params), nil
	case fexpr.SignNeq:
		return dbx.NewExp(fmt.Sprintf("COALESCE(%s, '') != COALESCE(%s, '')", lName, rName), params), nil
	case fexpr.SignLike:
		// both sides are columns and therefore wrap the right side with "%" for contains like behavior
		if len(params) == 0 {
			return dbx.NewExp(fmt.Sprintf("%s LIKE ('%%' || %s || '%%')", lName, rName), params), nil
		}

		// normalize operands and switch sides if the left operand is a number or text
		if len(lParams) > 0 {
			return dbx.NewExp(fmt.Sprintf("%s LIKE %s", rName, lName), f.normalizeLikeParams(params)), nil
		}

		return dbx.NewExp(fmt.Sprintf("%s LIKE %s", lName, rName), f.normalizeLikeParams(params)), nil
	case fexpr.SignNlike:
		// both sides are columns and therefore wrap the right side with "%" for not-contains like behavior
		if len(params) == 0 {
			return dbx.NewExp(fmt.Sprintf("%s NOT LIKE ('%%' || %s || '%%')", lName, rName), params), nil
		}

		// normalize operands and switch sides if the left operand is a number or text
		if len(lParams) > 0 {
			return dbx.NewExp(fmt.Sprintf("%s NOT LIKE %s", rName, lName), f.normalizeLikeParams(params)), nil
		}

		return dbx.NewExp(fmt.Sprintf("%s NOT LIKE %s", lName, rName), f.normalizeLikeParams(params)), nil
	case fexpr.SignLt:
		return dbx.NewExp(fmt.Sprintf("%s < %s", lName, rName), params), nil
	case fexpr.SignLte:
		return dbx.NewExp(fmt.Sprintf("%s <= %s", lName, rName), params), nil
	case fexpr.SignGt:
		return dbx.NewExp(fmt.Sprintf("%s > %s", lName, rName), params), nil
	case fexpr.SignGte:
		return dbx.NewExp(fmt.Sprintf("%s >= %s", lName, rName), params), nil
	}

	return nil, fmt.Errorf("Unknown expression operator %q", expr.Op)
}

func (f FilterData) resolveToken(token fexpr.Token, fieldResolver FieldResolver) (name string, params dbx.Params, err error) {
	switch token.Type {
	case fexpr.TokenIdentifier:
		// current datetime constant
		// ---
		if token.Literal == "@now" {
			placeholder := "t" + security.RandomString(7)
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
		placeholder := "t" + security.RandomString(7)
		name := fmt.Sprintf("{:%s}", placeholder)
		params := dbx.Params{placeholder: token.Literal}

		return name, params, nil
	case fexpr.TokenNumber:
		placeholder := "t" + security.RandomString(7)
		name := fmt.Sprintf("{:%s}", placeholder)
		params := dbx.Params{placeholder: cast.ToFloat64(token.Literal)}

		return name, params, nil
	}

	return "", nil, errors.New("Unresolvable token type.")
}

func (f FilterData) normalizeLikeParams(params dbx.Params) dbx.Params {
	result := dbx.Params{}

	for k, v := range params {
		vStr := cast.ToString(v)
		if !strings.Contains(vStr, "%") {
			vStr = "%" + vStr + "%"
		}
		result[k] = vStr
	}

	return result
}
