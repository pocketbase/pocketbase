package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

// filter modifiers
const (
	eachModifier    string = "each"
	issetModifier   string = "isset"
	lengthModifier  string = "length"
	lowerModifier   string = "lower"
	changedModifier string = "changed"
)

// ensure that `search.FieldResolver` interface is implemented
var _ search.FieldResolver = (*RecordFieldResolver)(nil)

// RecordFieldResolver defines a custom search resolver struct for
// managing Record model search fields.
//
// Usually used together with `search.Provider`.
// Example:
//
//	resolver := resolvers.NewRecordFieldResolver(
//	    app,
//	    myCollection,
//	    &models.RequestInfo{...},
//	    true,
//	)
//	provider := search.NewProvider(resolver)
//	...
type RecordFieldResolver struct {
	app               App
	baseCollection    *Collection
	requestInfo       *RequestInfo
	staticRequestInfo map[string]any
	allowedFields     []string
	joins             []*search.Join
	allowHiddenFields bool
	// ---
	listRuleJoins       map[string]*Collection // tableAlias->collection
	joinAliasSuffix     string                 // used for uniqueness in the flatten collection list rule join
	baseCollectionAlias string
}

// AllowedFields returns a copy of the resolver's allowed fields.
func (r *RecordFieldResolver) AllowedFields() []string {
	return slices.Clone(r.allowedFields)
}

// SetAllowedFields replaces the resolver's allowed fields with the new ones.
func (r *RecordFieldResolver) SetAllowedFields(newAllowedFields []string) {
	r.allowedFields = slices.Clone(newAllowedFields)
}

// AllowHiddenFields returns whether the current resolver allows filtering hidden fields.
func (r *RecordFieldResolver) AllowHiddenFields() bool {
	return r.allowHiddenFields
}

// SetAllowHiddenFields enables or disables hidden fields filtering.
func (r *RecordFieldResolver) SetAllowHiddenFields(allowHiddenFields bool) {
	r.allowHiddenFields = allowHiddenFields
}

// NewRecordFieldResolver creates and initializes a new `RecordFieldResolver`.
func NewRecordFieldResolver(
	app App,
	baseCollection *Collection,
	requestInfo *RequestInfo,
	allowHiddenFields bool,
) *RecordFieldResolver {
	r := &RecordFieldResolver{
		app:               app,
		baseCollection:    baseCollection,
		requestInfo:       requestInfo,
		allowHiddenFields: allowHiddenFields, // note: it is not based only on the requestInfo.auth since it could be used by a non-request internal method
		joins:             []*search.Join{},
		allowedFields: []string{
			`^\w+[\w\.\:]*$`,
			`^\@request\.context$`,
			`^\@request\.method$`,
			`^\@request\.auth\.[\w\.\:]*\w+$`,
			`^\@request\.body\.[\w\.\:]*\w+$`,
			`^\@request\.query\.[\w\.\:]*\w+$`,
			`^\@request\.headers\.[\w\.\:]*\w+$`,
			`^\@collection\.\w+(\:\w+)?\.[\w\.\:]*\w+$`,
		},
	}

	r.staticRequestInfo = map[string]any{}
	if r.requestInfo != nil {
		r.staticRequestInfo["context"] = r.requestInfo.Context
		r.staticRequestInfo["method"] = r.requestInfo.Method
		r.staticRequestInfo["query"] = r.requestInfo.Query
		r.staticRequestInfo["headers"] = r.requestInfo.Headers
		r.staticRequestInfo["body"] = r.requestInfo.Body
		r.staticRequestInfo["auth"] = nil
		if r.requestInfo.Auth != nil {
			authClone := r.requestInfo.Auth.Clone()
			r.staticRequestInfo["auth"] = authClone.
				Unhide(authClone.Collection().Fields.FieldNames()...).
				IgnoreEmailVisibility(true).
				PublicExport()
		}
	}

	return r
}

// @todo think of a better a way how to call it automatically after BuildExpr
//
// UpdateQuery implements `search.FieldResolver` interface.
//
// Conditionally updates the provided search query based on the
// resolved fields (eg. dynamically joining relations).
func (r *RecordFieldResolver) UpdateQuery(query *dbx.SelectQuery) error {
	if len(r.joins) > 0 {
		r.updateQueryWithDeduplicateConstraint(query)

		for _, join := range r.joins {
			query.LeftJoin(
				(join.TableName + " " + join.TableAlias),
				join.On,
			)
		}
	}

	// note: for now the joins are not applied for multi-match conditions to avoid excessive checks
	if len(r.listRuleJoins) > 0 {
		for alias, c := range r.listRuleJoins {
			err := r.updateQueryWithCollectionListRule(c, alias, query)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *RecordFieldResolver) updateQueryWithCollectionListRule(c *Collection, tableAlias string, query *dbx.SelectQuery) error {
	if r.allowHiddenFields || c == nil || c.ListRule == nil || *c.ListRule == "" {
		return nil
	}

	cloneR := *r
	cloneR.joins = []*search.Join{}
	cloneR.baseCollection = c
	cloneR.baseCollectionAlias = tableAlias
	cloneR.allowHiddenFields = true
	cloneR.joinAliasSuffix = security.PseudorandomString(8)

	expr, err := search.FilterData(*c.ListRule).BuildExpr(&cloneR)
	if err != nil {
		return fmt.Errorf("to buld %q list rule join subquery filter expression: %w", c.Name, err)
	}

	query.AndWhere(expr)

	if len(cloneR.joins) > 0 {
		r.updateQueryWithDeduplicateConstraint(query)

		for _, j := range cloneR.joins {
			query.LeftJoin(
				(j.TableName + " " + j.TableAlias),
				j.On,
			)
		}
	}

	return nil
}

func (r *RecordFieldResolver) updateQueryWithDeduplicateConstraint(query *dbx.SelectQuery) {
	query.Distinct(true)

	// @todo Research better options for generic rows deduplication.
	//
	// Disable the GROUP BY conditional checks for now since it prevents
	// proper utilization of ORDER BY indexes (and maybe others)
	// (https://github.com/pocketbase/pocketbase/discussions/7461)

	// info := query.Info()
	// if info.Distinct {
	// 	return
	// }

	// // already has the group by registered
	// var groupByCol = r.baseCollection.Name
	// if r.baseCollectionAlias != "" {
	// 	groupByCol = r.baseCollectionAlias
	// }
	// groupByCol += ".id"
	// if len(info.GroupBy) > 0 && info.GroupBy[0] == groupByCol {
	// 	return
	// }

	// // when deemed safe (GROUP BY could have different execution order compared to DISTINCT),
	// // prefer GROUP BY to deduplicate only on the id field instead of all columns
	// // so that the size of a single row wouldn't matter that much
	// if preferGroupBy(info, groupByCol) {
	// 	query.GroupBy(groupByCol)
	// } else {
	// 	query.Distinct(true)
	// }
}

func preferGroupBy(info *dbx.QueryInfo, fullUnquotedGroupByCol string) bool {
	if len(info.GroupBy) != 0 {
		return false
	}

	if info.Having != nil {
		return false
	}

	// dbx fallbacks to * if not set
	if len(info.Selects) == 0 {
		return true
	}

	if len(info.Selects) != 1 {
		return false
	}

	identifier := info.Selects[0]

	if identifier == "*" || identifier == fullUnquotedGroupByCol {
		return true
	}

	// try again as direct col match in an unquoted column format
	identifier = inflector.Columnify(identifier)
	if identifier == fullUnquotedGroupByCol {
		return true
	}

	// remains table.* to check
	// (aliased columns for now are ignored as they could be represented by expressions)
	if !strings.HasSuffix(identifier, ".*") {
		return false
	}

	return strings.HasPrefix(fullUnquotedGroupByCol, strings.TrimSuffix(identifier, "*"))
}

// Resolve implements `search.FieldResolver` interface.
//
// Example of some resolvable fieldName formats:
//
//	id
//	someSelect.each
//	project.screen.status
//	screen.project_via_prototype.name
//	@request.context
//	@request.method
//	@request.query.filter
//	@request.headers.x_token
//	@request.auth.someRelation.name
//	@request.body.someRelation.name
//	@request.body.someField
//	@request.body.someSelect:each
//	@request.body.someField:isset
//	@collection.product.name
func (r *RecordFieldResolver) Resolve(fieldName string) (*search.ResolverResult, error) {
	return parseAndRun(fieldName, r)
}

func (r *RecordFieldResolver) resolveStaticRequestField(path ...string) (*search.ResolverResult, error) {
	if len(path) == 0 {
		return nil, errors.New("at least one path key should be provided")
	}

	lastProp, modifier, err := splitModifier(path[len(path)-1])
	if err != nil {
		return nil, err
	}

	path[len(path)-1] = lastProp

	// extract value
	resultVal, err := extractNestedVal(r.staticRequestInfo, path...)
	if err != nil {
		r.app.Logger().Debug("resolveStaticRequestField graceful fallback", "error", err.Error())
	}

	if modifier == issetModifier {
		if err != nil {
			return &search.ResolverResult{Identifier: "FALSE"}, nil
		}
		return &search.ResolverResult{Identifier: "TRUE"}, nil
	}

	// note: we are ignoring the error because requestInfo is dynamic
	// and some of the lookup keys may not be defined for the request

	switch v := resultVal.(type) {
	case nil:
		// no further processing is needed...
	case string:
		// check if it is a number field and explicitly try to cast to
		// float in case of a numeric string value was used
		// (this usually the case when the data is from a multipart/form-data request)
		field := r.baseCollection.Fields.GetByName(path[len(path)-1])
		if field != nil && field.Type() == FieldTypeNumber {
			if nv, err := strconv.ParseFloat(v, 64); err == nil {
				resultVal = nv
			}
		}
		// otherwise - no further processing is needed...
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		// no further processing is needed...
	default:
		// non-plain value
		// try casting to string (in case for exampe fmt.Stringer is implemented)
		val, castErr := cast.ToStringE(v)

		// if that doesn't work, try encoding it
		if castErr != nil {
			encoded, jsonErr := json.Marshal(v)
			if jsonErr == nil {
				val = string(encoded)
			}
		}

		resultVal = val
	}

	// unsupported modifier
	// @todo consider deprecating with the introduction of filter functions
	if modifier != "" && modifier != lowerModifier {
		return nil, fmt.Errorf("invalid modifier sequence %s:%s", lastProp, modifier)
	}

	// no need to wrap as placeholder if we already know that it is null
	if resultVal == nil {
		return &search.ResolverResult{Identifier: "NULL"}, nil
	}

	placeholder := "f" + security.PseudorandomString(10)

	// @todo consider deprecating with the introduction of filter functions
	if modifier == lowerModifier {
		return &search.ResolverResult{
			Identifier: "LOWER({:" + placeholder + "})",
			Params:     dbx.Params{placeholder: resultVal},
		}, nil
	}

	return &search.ResolverResult{
		Identifier: "{:" + placeholder + "}",
		Params:     dbx.Params{placeholder: resultVal},
	}, nil
}

func (r *RecordFieldResolver) loadCollection(collectionNameOrId string) (*Collection, error) {
	if collectionNameOrId == r.baseCollection.Name || collectionNameOrId == r.baseCollection.Id {
		return r.baseCollection, nil
	}

	return getCollectionByModelOrIdentifier(r.app, collectionNameOrId)
}

func (r *RecordFieldResolver) registerJoin(tableName string, tableAlias string, on dbx.Expression) error {
	newJoin := &search.Join{
		TableName:  tableName,
		TableAlias: tableAlias,
		On:         on,
	}

	// (see updateQueryWithCollectionListRule)
	if !r.allowHiddenFields {
		c, _ := r.loadCollection(tableName)

		// ignore non-collections since the table name could be an expression (e.g. json) or some other subquery
		if c != nil {
			// treat all fields as if they are hidden
			if c.ListRule == nil {
				return fmt.Errorf("%q fields can be accessed only when allowHiddenFields is enabled or by superusers", c.Name)
			}

			if r.listRuleJoins == nil {
				r.listRuleJoins = map[string]*Collection{}
			}
			r.listRuleJoins[newJoin.TableAlias] = c
		}
	}

	// replace existing join
	for i, j := range r.joins {
		if j.TableAlias == newJoin.TableAlias {
			r.joins[i] = newJoin
			return nil
		}
	}

	// register new join
	r.joins = append(r.joins, newJoin)
	return nil
}

type mapExtractor interface {
	AsMap() map[string]any
}

func extractNestedVal(rawData any, keys ...string) (any, error) {
	if len(keys) == 0 {
		return nil, errors.New("at least one key should be provided")
	}

	switch m := rawData.(type) {
	// maps
	case map[string]any:
		return mapVal(m, keys...)
	case map[string]string:
		return mapVal(m, keys...)
	case map[string]bool:
		return mapVal(m, keys...)
	case map[string]float32:
		return mapVal(m, keys...)
	case map[string]float64:
		return mapVal(m, keys...)
	case map[string]int:
		return mapVal(m, keys...)
	case map[string]int8:
		return mapVal(m, keys...)
	case map[string]int16:
		return mapVal(m, keys...)
	case map[string]int32:
		return mapVal(m, keys...)
	case map[string]int64:
		return mapVal(m, keys...)
	case map[string]uint:
		return mapVal(m, keys...)
	case map[string]uint8:
		return mapVal(m, keys...)
	case map[string]uint16:
		return mapVal(m, keys...)
	case map[string]uint32:
		return mapVal(m, keys...)
	case map[string]uint64:
		return mapVal(m, keys...)
	case mapExtractor:
		return mapVal(m.AsMap(), keys...)
	case types.JSONRaw:
		var raw any
		err := json.Unmarshal(m, &raw)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal raw JSON in order extract nested value from: %w", err)
		}
		return extractNestedVal(raw, keys...)

	// slices
	case []string:
		return arrVal(m, keys...)
	case []bool:
		return arrVal(m, keys...)
	case []float32:
		return arrVal(m, keys...)
	case []float64:
		return arrVal(m, keys...)
	case []int:
		return arrVal(m, keys...)
	case []int8:
		return arrVal(m, keys...)
	case []int16:
		return arrVal(m, keys...)
	case []int32:
		return arrVal(m, keys...)
	case []int64:
		return arrVal(m, keys...)
	case []uint:
		return arrVal(m, keys...)
	case []uint8:
		return arrVal(m, keys...)
	case []uint16:
		return arrVal(m, keys...)
	case []uint32:
		return arrVal(m, keys...)
	case []uint64:
		return arrVal(m, keys...)
	case []mapExtractor:
		extracted := make([]any, len(m))
		for i, v := range m {
			extracted[i] = v.AsMap()
		}
		return arrVal(extracted, keys...)
	case []any:
		return arrVal(m, keys...)
	case []types.JSONRaw:
		return arrVal(m, keys...)
	default:
		return nil, fmt.Errorf("expected map or array, got %#v", rawData)
	}
}

func mapVal[T any](m map[string]T, keys ...string) (any, error) {
	result, ok := m[keys[0]]
	if !ok {
		return nil, fmt.Errorf("invalid key path - missing key %q", keys[0])
	}

	// end key reached
	if len(keys) == 1 {
		return result, nil
	}

	return extractNestedVal(result, keys[1:]...)
}

func arrVal[T any](m []T, keys ...string) (any, error) {
	idx, err := strconv.Atoi(keys[0])
	if err != nil || idx < 0 || idx >= len(m) {
		return nil, fmt.Errorf("invalid key path - invalid or missing array index %q", keys[0])
	}

	result := m[idx]

	// end key reached
	if len(keys) == 1 {
		return result, nil
	}

	return extractNestedVal(result, keys[1:]...)
}

func splitModifier(combined string) (string, string, error) {
	parts := strings.Split(combined, ":")

	if len(parts) != 2 {
		return combined, "", nil
	}

	// validate modifier
	switch parts[1] {
	case issetModifier,
		eachModifier,
		lengthModifier,
		lowerModifier,
		changedModifier:
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unknown modifier in %q", combined)
}
