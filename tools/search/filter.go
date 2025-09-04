package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ganigeorgiev/fexpr"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/spf13/cast"
)

// FilterData is a filter expression string following the `fexpr` package grammar.
//
// The filter string can also contain dbx placeholder parameters (eg. "title = {:name}"),
// that will be safely replaced and properly quoted inplace with the placeholderReplacements values.
//
// Example:
//
//	var filter FilterData = "id = null || (name = 'test' && status = true) || (total >= {:min} && total <= {:max})"
//	resolver := search.NewSimpleFieldResolver("id", "name", "status")
//	expr, err := filter.BuildExpr(resolver, dbx.Params{"min": 100, "max": 200})
type FilterData string

// parsedFilterData holds a cache with previously parsed filter data expressions
// (initialized with some preallocated empty data map)
var parsedFilterData = store.New(make(map[string][]fexpr.ExprGroup, 50))

// BuildExpr parses the current filter data and returns a new db WHERE expression.
//
// The filter string can also contain dbx placeholder parameters (eg. "title = {:name}"),
// that will be safely replaced and properly quoted inplace with the placeholderReplacements values.
//
// The parsed expressions are limited up to DefaultFilterExprLimit.
// Use [FilterData.BuildExprWithLimit] if you want to set a custom limit.
func (f FilterData) BuildExpr(
	fieldResolver FieldResolver,
	placeholderReplacements ...dbx.Params,
) (dbx.Expression, error) {
	return f.BuildExprWithLimit(fieldResolver, DefaultFilterExprLimit, placeholderReplacements...)
}

// BuildExpr parses the current filter data and returns a new db WHERE expression.
//
// The filter string can also contain dbx placeholder parameters (eg. "title = {:name}"),
// that will be safely replaced and properly quoted inplace with the placeholderReplacements values.
func (f FilterData) BuildExprWithLimit(
	fieldResolver FieldResolver,
	maxExpressions int,
	placeholderReplacements ...dbx.Params,
) (dbx.Expression, error) {
	raw := string(f)

	// replace the placeholder params in the raw string filter
	for _, p := range placeholderReplacements {
		for key, value := range p {
			var replacement string
			switch v := value.(type) {
			case nil:
				replacement = "null"
			case bool, float64, float32, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
				replacement = cast.ToString(v)
			default:
				replacement = cast.ToString(v)

				// try to json serialize as fallback
				if replacement == "" {
					raw, _ := json.Marshal(v)
					replacement = string(raw)
				}

				replacement = strconv.Quote(replacement)
			}
			raw = strings.ReplaceAll(raw, "{:"+key+"}", replacement)
		}
	}

	cacheKey := raw + "/" + strconv.Itoa(maxExpressions)

	if data, ok := parsedFilterData.GetOk(cacheKey); ok {
		return buildParsedFilterExpr(data, fieldResolver, &maxExpressions)
	}

	data, err := fexpr.Parse(raw)
	if err != nil {
		// depending on the users demand we may allow empty expressions
		// (aka. expressions consisting only of whitespaces or comments)
		// but for now disallow them as it seems unnecessary
		// if errors.Is(err, fexpr.ErrEmpty) {
		// return dbx.NewExp("1=1"), nil
		// }

		return nil, err
	}

	// store in cache
	// (the limit size is arbitrary and it is there to prevent the cache growing too big)
	parsedFilterData.SetIfLessThanLimit(cacheKey, data, 500)

	return buildParsedFilterExpr(data, fieldResolver, &maxExpressions)
}

func buildParsedFilterExpr(data []fexpr.ExprGroup, fieldResolver FieldResolver, maxExpressions *int) (dbx.Expression, error) {
	if len(data) == 0 {
		return nil, fexpr.ErrEmpty
	}

	result := &concatExpr{separator: " "}

	for _, group := range data {
		var expr dbx.Expression
		var exprErr error

		switch item := group.Item.(type) {
		case fexpr.Expr:
			if *maxExpressions <= 0 {
				return nil, ErrFilterExprLimit
			}

			*maxExpressions--

			expr, exprErr = resolveTokenizedExpr(item, fieldResolver)
		case fexpr.ExprGroup:
			expr, exprErr = buildParsedFilterExpr([]fexpr.ExprGroup{item}, fieldResolver, maxExpressions)
		case []fexpr.ExprGroup:
			expr, exprErr = buildParsedFilterExpr(item, fieldResolver, maxExpressions)
		default:
			exprErr = errors.New("unsupported expression item")
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

func resolveTokenizedExpr(expr fexpr.Expr, fieldResolver FieldResolver) (dbx.Expression, error) {
	lResult, lErr := resolveToken(expr.Left, fieldResolver)
	if lErr != nil || lResult.Identifier == "" {
		return nil, fmt.Errorf("invalid left operand %q - %v", expr.Left.Literal, lErr)
	}

	rResult, rErr := resolveToken(expr.Right, fieldResolver)
	if rErr != nil || rResult.Identifier == "" {
		return nil, fmt.Errorf("invalid right operand %q - %v", expr.Right.Literal, rErr)
	}

	return buildResolversExpr(lResult, expr.Op, rResult)
}

func buildResolversExpr(
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
		/* SQLite:
		expr = dbx.NewExp(fmt.Sprintf("%s < %s", left.Identifier, right.Identifier), mergeParams(left.Params, right.Params))
		*/
		// PostgreSQL:
		expr = dbx.NewExp(numericJoin(left, "<", right), mergeParams(left.Params, right.Params))
	case fexpr.SignLte, fexpr.SignAnyLte:
		/* SQLite:
		expr = dbx.NewExp(fmt.Sprintf("%s <= %s", left.Identifier, right.Identifier), mergeParams(left.Params, right.Params))
		*/
		// PostgreSQL:
		expr = dbx.NewExp(numericJoin(left, "<=", right), mergeParams(left.Params, right.Params))
	case fexpr.SignGt, fexpr.SignAnyGt:
		/* SQLite:
		expr = dbx.NewExp(fmt.Sprintf("%s > %s", left.Identifier, right.Identifier), mergeParams(left.Params, right.Params))
		*/
		// PostgreSQL:
		expr = dbx.NewExp(numericJoin(left, ">", right), mergeParams(left.Params, right.Params))
	case fexpr.SignGte, fexpr.SignAnyGte:
		/* SQLite:
		expr = dbx.NewExp(fmt.Sprintf("%s >= %s", left.Identifier, right.Identifier), mergeParams(left.Params, right.Params))
		*/
		// PostgreSQL:
		expr = dbx.NewExp(numericJoin(left, ">=", right), mergeParams(left.Params, right.Params))
	}

	if expr == nil {
		return nil, fmt.Errorf("unknown expression operator %q", op)
	}

	// multi-match expressions
	if !isAnyMatchOp(op) {
		if left.MultiMatchSubQuery != nil && right.MultiMatchSubQuery != nil {
			mm := &manyVsManyExpr{
				left:  left,
				right: right,
				op:    op,
			}

			expr = dbx.Enclose(dbx.And(expr, mm))
		} else if left.MultiMatchSubQuery != nil {
			mm := &manyVsOneExpr{
				noCoalesce:   left.NoCoalesce,
				subQuery:     left.MultiMatchSubQuery,
				op:           op,
				otherOperand: right,
			}

			expr = dbx.Enclose(dbx.And(expr, mm))
		} else if right.MultiMatchSubQuery != nil {
			mm := &manyVsOneExpr{
				noCoalesce:   right.NoCoalesce,
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

var normalizedIdentifiers = map[string]string{
	/* SQLite:
	// if `null` field is missing, treat `null` identifier as NULL token
	"null": "NULL",
	// if `true` field is missing, treat `true` identifier as TRUE token
	"true": "1",
	// if `false` field is missing, treat `false` identifier as FALSE token
	"false": "0",
	*/
	// PostgreSQL:
	"null":  "NULL",
	"true":  "TRUE",
	"false": "FALSE",
}

func resolveToken(token fexpr.Token, fieldResolver FieldResolver) (*ResolverResult, error) {
	switch token.Type {
	case fexpr.TokenIdentifier:
		// check for macros
		// ---
		if macroFunc, ok := identifierMacros[token.Literal]; ok {
			placeholder := "t" + security.PseudorandomString(8)

			macroValue, err := macroFunc()
			if err != nil {
				return nil, err
			}

			return &ResolverResult{
				Identifier: "{:" + placeholder + "}",
				Params:     dbx.Params{placeholder: macroValue},
			}, nil
		}

		// custom resolver
		// ---
		result, err := fieldResolver.Resolve(token.Literal)
		if err != nil || result.Identifier == "" {
			for k, v := range normalizedIdentifiers {
				if strings.EqualFold(k, token.Literal) {
					return &ResolverResult{Identifier: v}, nil
				}
			}
			return nil, err
		}

		return result, err
	case fexpr.TokenText:
		// PostgreSQL only:
		// if we know it is a emty string, use the empty string directly.
		if token.Literal == "" {
			return &ResolverResult{Identifier: `''`}, nil
		}

		placeholder := "t" + security.PseudorandomString(8)

		return &ResolverResult{
			Identifier: "{:" + placeholder + "}",
			Params:     dbx.Params{placeholder: token.Literal},
		}, nil
	case fexpr.TokenNumber:
		/* SQLite:
		placeholder := "t" + security.PseudorandomString(8)

		return &ResolverResult{
			Identifier: "{:" + placeholder + "}",
			Params:     dbx.Params{placeholder: cast.ToFloat64(token.Literal)},
		}, nil
		*/
		// PostgreSQL:
		// handle a special case (where 1 = 1) where both left and right identifiers are numeric numbers.
		// Eg: To prevent SQL injection, for query "1=1", dbx will generate "select xxx where $1 = $2" (prepared statement) with params [1, 1].
		// because we didn't specify the type for both $1 and $2, so PostgreSQL will treat them as text, and expect all params to be text types.
		// And it failed to cast numeric type `1` to text `"1"` and throws an error:
		// Error: `failed to encode args[0]: unable to encode 1 into text format for text (OID 25): cannot find encode plan;`
		// Related Issue:
		// - https://github.com/jackc/pgx/issues/798,
		// - https://github.com/jackc/pgx/issues/2307
		// This is not caused by an issue of pgx, but by the strong type validation of PostgreSQL.
		//
		// To fix it, we have two options:
		// Option 1: add a explict type cast: "{:" + placeholder + "}::numeric",
		// Option 2: use the number literal directly without a param placeholder.
		// We have to convert user input to float64 to remove any harmful characters to avoid SQL injection.
		safeNumberStr := strconv.FormatFloat(cast.ToFloat64(token.Literal), 'f', -1, 64)
		return &ResolverResult{
			Identifier: safeNumberStr,
			Params:     dbx.Params{},
		}, nil
	case fexpr.TokenFunction:
		fn, ok := TokenFunctions[token.Literal]
		if !ok {
			return nil, fmt.Errorf("unknown function %q", token.Literal)
		}

		args, _ := token.Meta.([]fexpr.Token)
		return fn(func(argToken fexpr.Token) (*ResolverResult, error) {
			return resolveToken(argToken, fieldResolver)
		}, args...)
	}

	return nil, fmt.Errorf("unsupported token type %q", token.Type)
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

	/* SQLite:
	equalOp := "="
	nullEqualOp := "IS"
	*/
	// PostgreSQL:
	equalOp := "="
	nullEqualOp := "IS NOT DISTINCT FROM"
	concatOp := "OR"
	nullExpr := "IS NULL"
	if !equal {
		/* SQLite:
		// always use `IS NOT` instead of `!=` because direct non-equal comparisons
		// to nullable column values that are actually NULL yields to NULL instead of TRUE, eg.:
		// `'example' != nullableColumn` -> NULL even if nullableColumn row value is NULL
		// Note: `select 'non-null-string' != NULL` returns NULL instead of True.
		equalOp = "IS NOT"
		nullEqualOp = equalOp
		*/
		// PostgreSQL:
		// In PostgreSQL, `IS NOT` only works for NULL values, but not for empty strings.
		// `IS DISTINCT FROM` works like SQLite's `IS NOT`.
		equalOp = "IS DISTINCT FROM"
		nullEqualOp = equalOp
		concatOp = "AND"
		nullExpr = "IS NOT NULL"
	}

	// no coalesce (eg. compare to a json field)
	// a IS b
	// a IS NOT b
	if left.NoCoalesce || right.NoCoalesce {
		return dbx.NewExp(
			/* SQLite:
			fmt.Sprintf("%s %s %s", left.Identifier, nullEqualOp, right.Identifier),
			*/
			typeAwareJoinNoCoalesce(left, nullEqualOp, right),
			mergeParams(left.Params, right.Params),
		)
	}

	// both operands are empty
	if isLeftEmpty && isRightEmpty {
		return dbx.NewExp(fmt.Sprintf("'' %s ''", equalOp), mergeParams(left.Params, right.Params))
	}

	// direct compare since at least one of the operands is known to be non-empty
	// eg. a = 'example'
	if isKnownNonEmptyIdentifier(left) || isKnownNonEmptyIdentifier(right) {
		/* SQLite:

		leftIdentifier := left.Identifier
		if isLeftEmpty {
			leftIdentifier = "''"
		}
		rightIdentifier := right.Identifier
		if isRightEmpty {
			rightIdentifier = "''"
		}
		*/
		// PostgreSQL:
		// TODO：
		// create a copy of ResolvedResult.
		// If it is empty string, show a empty string.
		// Remember to remove the params from the shadow copy if it is empty or null
		// leftIdentifier := left.Identifier
		// if isLeftEmpty {
		// 	leftIdentifier = "''"
		// }
		// rightIdentifier := right.Identifier
		// if isRightEmpty {
		// 	rightIdentifier = "''"
		// }

		return dbx.NewExp(
			/* SQLite:
			fmt.Sprintf("%s %s %s", leftIdentifier, equalOp, rightIdentifier),
			*/
			// PostgreSQL:
			typeAwareJoinNoCoalesce(left, equalOp, right),
			mergeParams(left.Params, right.Params),
		)
	}

	// Hint: In PocketBase's world, NULL is treated the same as empty.
	// "" = b OR b IS NULL
	// "" IS NOT b AND b IS NOT NULL
	if isLeftEmpty {
		return dbx.NewExp(
			/* SQLite:
			fmt.Sprintf("('' %s %s %s %s %s)", equalOp, right.Identifier, concatOp, right.Identifier, nullExpr),
			*/
			// PostgreSQL:
			fmt.Sprintf("('' %s %s %s %s %s)", equalOp, withNonJsonbType(right.Identifier, "text"), concatOp, right.Identifier, nullExpr),
			mergeParams(left.Params, right.Params),
		)
	}

	// a = "" OR a IS NULL
	// a IS NOT "" AND a IS NOT NULL
	if isRightEmpty {
		return dbx.NewExp(
			/* SQLite:
			fmt.Sprintf("(%s %s '' %s %s %s)", left.Identifier, equalOp, concatOp, left.Identifier, nullExpr),
			*/
			// PostgreSQL:
			// Note: pocketbase treats empty string the same as NULL.
			// eg: WHERE col_int::text = '' OR col_int IS NULL
			fmt.Sprintf("(%s %s '' %s %s %s)", withNonJsonbType(left.Identifier, "text"), equalOp, concatOp, left.Identifier, nullExpr),
			mergeParams(left.Params, right.Params),
		)
	}

	/* SQLite:
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
	*/
	// PostgreSQL:
	// 1. We can't use COALESCE() here, because we never know the type of the column to be compared.
	//    Otherwise, PostgreSQL will throw a type mismatch error if we use default empty string.
	// 2. to_jsonb() erase the type so that different types can be compared safely.
	// 3. Use `nullEqualOp` instead of `equalOp` to safely compare null values, similar to COALESCE(),
	//    because NULL::jsonb behaves same as NULL. If either part of the equal operation is NULL,
	//    then it will produce a NULL output, and we need something like COALESCE() to avoid NULL output.
	return dbx.NewExp(
		fmt.Sprintf("%s %s %s", castToJsonb(left), nullEqualOp, castToJsonb(right)),
		mergeParams(left.Params, right.Params),
	)
}

// PostgreSQL only:
// PostgreSQL lets us write '2024-09-03' and use it as a date, timestamp, text, etc., without explicit casts every time.
// Normally, when we use `SELECT col_text = 'abc'`, the type of 'abc' can be automatically infered to `text`.
// However, when used with `to_jsonb('abc')` function, the type of 'abc' is not determistic, because to_jsonb() can
// handle many different types. So we need to add explicit type hints before using in to_jsonb().
//
// Currently, it only affects:
// 1. NULL
// 2. String Params in PreparedStatements.
// 3. Numeric Params in PreparedStatements.
//
// Only used with `to_jsonb`
func castToJsonb(identifier *ResolverResult) string {
	if strings.ToLower(identifier.Identifier) == "null" {
		return "to_jsonb(NULL::text)"
	}
	if tp := inferPolymorphicLiteral(identifier); tp != "" {
		return fmt.Sprintf("to_jsonb(%s::%s)", identifier.Identifier, tp)
	}
	return fmt.Sprintf("to_jsonb(%s)", identifier.Identifier)
}

// There are some json types:
// 1. null    -> Undetermine Polymorphic Type, can be any PostgreSQL types
// 2. text    -> Undetermine Polymorphic Type, can be Date, TimeStamp, text, etc.
// 3. numbers -> Deterministic type, always numeric, no type cast needed
// 4. bool    -> Deterministic type, always boolean, no type cast needed
func inferPolymorphicLiteral(result *ResolverResult) string {
	// Note: result cannot be "NULL" identifier when called in [inferPolymorphicLiteral],
	// because we already handled "NULL" seperately before calling this function.
	// See [resolveEqualExpr] for details.
	if strings.ToLower(result.Identifier) == "null" {
		return "null"
	}

	if result.Identifier == `''` {
		return "text"
	}

	if len(result.Params) == 1 {
		for _, p := range result.Params {
			switch p.(type) {
			case nil:
				panic("Unexpected nil type, nil is supposed to be parsed as NULL identifier")
			case string:
				return "text"
			}
		}
	}
	return ""
}

// There are some json types:
// 1. null    -> Undetermine Polymorphic Type, can be any PostgreSQL types
// 2. text    -> Undetermine Polymorphic Type, can be Date, TimeStamp, text, etc.
// 3. numbers -> Deterministic type, always numeric, no type cast needed
// 4. bool    -> Deterministic type, always boolean, no type cast needed
func inferDeterministicType(result *ResolverResult) string {
	// If there is a explict type cast suffix, then we can use it to determine the type.
	match := regexRightMostTypeCast.FindStringSubmatch(strings.TrimRight(result.Identifier, " "))
	if len(match) > 0 {
		return match[1]
	}

	// If the type is boolean, we can use it directly.
	if strings.ToLower(result.Identifier) == "true" || strings.ToLower(result.Identifier) == "false" {
		return "boolean"
	}

	// If the type is numbers, we can use it directly.
	if _, err := strconv.ParseFloat(result.Identifier, 64); err == nil {
		return "numeric"
	}
	if strings.HasPrefix(result.Identifier, "{:") && len(result.Params) == 1 {
		for _, p := range result.Params {
			switch p.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				return "numeric"
			}
		}
	}

	return ""
}

var regexRightMostTypeCast = regexp.MustCompile(`::(\w+)$`)

// PostgreSQL only:
// If either left or right identifier has a specific type cast, we need to add the same type cast to the other identifier.
func typeAwareJoinNoCoalesce(l *ResolverResult, op string, r *ResolverResult) string {
	left := strings.TrimRight(l.Identifier, " ")
	right := strings.TrimRight(r.Identifier, " ")

	leftType := inferDeterministicType(l)
	rightType := inferDeterministicType(r)
	if len(leftType) > 0 && len(rightType) > 0 {
		// If left and right identifiers have different type cast, force cast both identifiers
		// to `jsonb`` type to bypass PostgreSQL's strict type validation error.
		if leftType != rightType {
			if leftType != "jsonb" {
				left = castToJsonb(l)
			}
			if rightType != "jsonb" {
				right = castToJsonb(r)
			}
		}
		// If both identifiers have the same type cast, return it directly.
		return fmt.Sprintf("%s %s %s", left, op, right)
	}
	// If none of the identifiers have type cast
	if len(leftType) == 0 && len(rightType) == 0 {
		return fmt.Sprintf("%s %s %s", left, op, right)
	}
	if len(leftType) > 0 {
		if leftType == "jsonb" {
			// implict cast is not possible for jsonb type
			right = castToJsonb(r)
		}

		// LeftType is Deterministic, RightType is Polymorphic, allow PostgreSQL to do auto implict cast.
		return fmt.Sprintf("%s %s %s", left, op, right)
	}
	if len(rightType) > 0 {
		if rightType == "jsonb" {
			left = castToJsonb(l)
		}

		return fmt.Sprintf("%s %s %s", left, op, right)
	}
	panic("should not reach here")
}

// PostgreSQL only:
// Force cast both identifiers to numeric type.
func numericJoin(l *ResolverResult, op string, r *ResolverResult) string {
	left := l.Identifier
	right := r.Identifier

	// Note: Polyphormic literal such as "2" can be automatic casted to numeric type by PostgreSQL.
	// Eg: SELECT "2" > 1  -- works fine.
	if inferDeterministicType(l) != "numeric" && inferPolymorphicLiteral(l) == "" {
		left = withNonJsonbType(left, "numeric")
	}
	if inferDeterministicType(r) != "numeric" && inferPolymorphicLiteral(r) == "" {
		right = withNonJsonbType(right, "numeric")
	}

	return fmt.Sprintf("%s %s %s", left, op, right)
}

// PostgreSQL only:
func withNonJsonbType(identifier string, targetType string) string {
	// Note:
	// DO NOT drop existing type cast before adding a new cast.
	// Reason: `1::numeric::text` is valid but `1::text` is invalid.
	suffix := "::" + targetType
	if strings.HasSuffix(identifier, suffix) {
		return identifier
	}
	return identifier + suffix
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

	if len(result.Params) == 0 {
		if _, err := strconv.ParseFloat(result.Identifier, 64); err == nil {
			return true
		}
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

// @todo consider adding support for custom single character wildcard
//
// wrapLikeParams wraps each provided param value string with `%`
// if the param doesn't contain an explicit wildcard (`%`) character already.
func wrapLikeParams(params dbx.Params) dbx.Params {
	result := dbx.Params{}

	for k, v := range params {
		vStr := cast.ToString(v)
		if !containsUnescapedChar(vStr, '%') {
			// note: this is done to minimize the breaking changes and to preserve the original autoescape behavior
			vStr = escapeUnescapedChars(vStr, '\\', '%', '_')
			vStr = "%" + vStr + "%"
		}
		result[k] = vStr
	}

	return result
}

func escapeUnescapedChars(str string, escapeChars ...rune) string {
	rs := []rune(str)
	total := len(rs)
	result := make([]rune, 0, total)

	var match bool

	for i := total - 1; i >= 0; i-- {
		if match {
			// check if already escaped
			if rs[i] != '\\' {
				result = append(result, '\\')
			}
			match = false
		} else {
			for _, ec := range escapeChars {
				if rs[i] == ec {
					match = true
					break
				}
			}
		}

		result = append(result, rs[i])

		// in case the matching char is at the beginning
		if i == 0 && match {
			result = append(result, '\\')
		}
	}

	// reverse
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

func containsUnescapedChar(str string, ch rune) bool {
	var prev rune

	for _, c := range str {
		if c == ch && prev != '\\' {
			return true
		}

		if c == '\\' && prev == '\\' {
			prev = rune(0) // reset escape sequence
		} else {
			prev = c
		}
	}

	return false
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
	separator string
	parts     []dbx.Expression
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
	left  *ResolverResult
	right *ResolverResult
	op    fexpr.SignOp
}

// Build converts the expression into a SQL fragment.
//
// Implements [dbx.Expression] interface.
func (e *manyVsManyExpr) Build(db *dbx.DB, params dbx.Params) string {
	if e.left.MultiMatchSubQuery == nil || e.right.MultiMatchSubQuery == nil {
		return "0=1"
	}

	lAlias := "__ml" + security.PseudorandomString(8)
	rAlias := "__mr" + security.PseudorandomString(8)

	whereExpr, buildErr := buildResolversExpr(
		&ResolverResult{
			NoCoalesce: e.left.NoCoalesce,
			Identifier: "[[" + lAlias + ".multiMatchValue]]",
		},
		e.op,
		&ResolverResult{
			NoCoalesce: e.right.NoCoalesce,
			Identifier: "[[" + rAlias + ".multiMatchValue]]",
			// note: the AfterBuild needs to be handled only once and it
			// doesn't matter whether it is applied on the left or right subquery operand
			AfterBuild: dbx.Not, // inverse for the not-exist expression
		},
	)

	if buildErr != nil {
		return "0=1"
	}

	return fmt.Sprintf(
		/* SQLite:
		"NOT EXISTS (SELECT 1 FROM (%s) {{%s}} LEFT JOIN (%s) {{%s}} WHERE %s)",
		*/
		// PostgreSQL:
		"NOT EXISTS (SELECT 1 FROM (%s) {{%s}} LEFT JOIN (%s) {{%s}} ON 1 = 1 WHERE %s)",
		e.left.MultiMatchSubQuery.Build(db, params),
		lAlias,
		e.right.MultiMatchSubQuery.Build(db, params),
		rAlias,
		whereExpr.Build(db, params),
	)
}

// -------------------------------------------------------------------

var _ dbx.Expression = (*manyVsOneExpr)(nil)

// manyVsOneExpr constructs a multi-match many<->one db where expression.
//
// Expects subQuery to return a subquery with a single "multiMatchValue" column.
//
// You can set inverse=false to reverse the condition sides (aka. one<->many).
type manyVsOneExpr struct {
	otherOperand *ResolverResult
	subQuery     dbx.Expression
	op           fexpr.SignOp
	inverse      bool
	noCoalesce   bool
}

// Build converts the expression into a SQL fragment.
//
// Implements [dbx.Expression] interface.
func (e *manyVsOneExpr) Build(db *dbx.DB, params dbx.Params) string {
	if e.subQuery == nil {
		return "0=1"
	}

	alias := "__sm" + security.PseudorandomString(8)

	r1 := &ResolverResult{
		NoCoalesce: e.noCoalesce,
		Identifier: "[[" + alias + ".multiMatchValue]]",
		AfterBuild: dbx.Not, // inverse for the not-exist expression
	}

	r2 := &ResolverResult{
		Identifier: e.otherOperand.Identifier,
		Params:     e.otherOperand.Params,
	}

	var whereExpr dbx.Expression
	var buildErr error

	if e.inverse {
		whereExpr, buildErr = buildResolversExpr(r2, e.op, r1)
	} else {
		whereExpr, buildErr = buildResolversExpr(r1, e.op, r2)
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
