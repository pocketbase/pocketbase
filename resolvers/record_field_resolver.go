package resolvers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

// filter modifiers
const (
	eachModifier   string = "each"
	issetModifier  string = "isset"
	lengthModifier string = "length"
)

// list of auth filter fields that don't require join with the auth
// collection or any other extra checks to be resolved.
var plainRequestAuthFields = []string{
	"@request.auth." + schema.FieldNameId,
	"@request.auth." + schema.FieldNameCollectionId,
	"@request.auth." + schema.FieldNameCollectionName,
	"@request.auth." + schema.FieldNameUsername,
	"@request.auth." + schema.FieldNameEmail,
	"@request.auth." + schema.FieldNameEmailVisibility,
	"@request.auth." + schema.FieldNameVerified,
	"@request.auth." + schema.FieldNameCreated,
	"@request.auth." + schema.FieldNameUpdated,
}

// ensure that `search.FieldResolver` interface is implemented
var _ search.FieldResolver = (*RecordFieldResolver)(nil)

// CollectionsFinder defines a common interface for retrieving
// collections and other related models.
//
// The interface at the moment is primarily used to avoid circular
// dependency with the daos.Dao package.
type CollectionsFinder interface {
	FindCollectionByNameOrId(collectionNameOrId string) (*models.Collection, error)
}

// RecordFieldResolver defines a custom search resolver struct for
// managing Record model search fields.
//
// Usually used together with `search.Provider`.
// Example:
//
//	resolver := resolvers.NewRecordFieldResolver(
//	    app.Dao(),
//	    myCollection,
//	    &models.RequestInfo{...},
//	    true,
//	)
//	provider := search.NewProvider(resolver)
//	...
type RecordFieldResolver struct {
	dao               CollectionsFinder
	baseCollection    *models.Collection
	requestInfo       *models.RequestInfo
	staticRequestInfo map[string]any
	allowedFields     []string
	loadedCollections []*models.Collection
	joins             []*join
	allowHiddenFields bool
}

// NewRecordFieldResolver creates and initializes a new `RecordFieldResolver`.
func NewRecordFieldResolver(
	dao CollectionsFinder,
	baseCollection *models.Collection,
	requestInfo *models.RequestInfo,
	// @todo consider moving per filter basis
	allowHiddenFields bool,
) *RecordFieldResolver {
	r := &RecordFieldResolver{
		dao:               dao,
		baseCollection:    baseCollection,
		requestInfo:       requestInfo,
		allowHiddenFields: allowHiddenFields,
		joins:             []*join{},
		loadedCollections: []*models.Collection{baseCollection},
		allowedFields: []string{
			`^\w+[\w\.\:]*$`,
			`^\@request\.context$`,
			`^\@request\.method$`,
			`^\@request\.auth\.[\w\.\:]*\w+$`,
			`^\@request\.data\.[\w\.\:]*\w+$`,
			`^\@request\.query\.[\w\.\:]*\w+$`,
			`^\@request\.headers\.\w+$`,
			`^\@collection\.\w+(\:\w+)?\.[\w\.\:]*\w+$`,
		},
	}

	r.staticRequestInfo = map[string]any{}
	if r.requestInfo != nil {
		r.staticRequestInfo["context"] = r.requestInfo.Context
		r.staticRequestInfo["method"] = r.requestInfo.Method
		r.staticRequestInfo["query"] = r.requestInfo.Query
		r.staticRequestInfo["headers"] = r.requestInfo.Headers
		r.staticRequestInfo["data"] = r.requestInfo.Data
		r.staticRequestInfo["auth"] = nil
		if r.requestInfo.AuthRecord != nil {
			authData := r.requestInfo.AuthRecord.PublicExport()
			// always add the record email no matter of the emailVisibility field
			authData[schema.FieldNameEmail] = r.requestInfo.AuthRecord.Email()
			r.staticRequestInfo["auth"] = authData
		}
	}

	return r
}

// UpdateQuery implements `search.FieldResolver` interface.
//
// Conditionally updates the provided search query based on the
// resolved fields (eg. dynamically joining relations).
func (r *RecordFieldResolver) UpdateQuery(query *dbx.SelectQuery) error {
	if len(r.joins) > 0 {
		query.Distinct(true)

		for _, join := range r.joins {
			query.LeftJoin(
				(join.tableName + " " + join.tableAlias),
				join.on,
			)
		}
	}

	return nil
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
//	@request.data.someRelation.name
//	@request.data.someField
//	@request.data.someSelect:each
//	@request.data.someField:isset
//	@collection.product.name
func (r *RecordFieldResolver) Resolve(fieldName string) (*search.ResolverResult, error) {
	return parseAndRun(fieldName, r)
}

func (r *RecordFieldResolver) resolveStaticRequestField(path ...string) (*search.ResolverResult, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("at least one path key should be provided")
	}

	lastProp, modifier, err := splitModifier(path[len(path)-1])
	if err != nil {
		return nil, err
	}

	path[len(path)-1] = lastProp

	// extract value
	resultVal, err := extractNestedMapVal(r.staticRequestInfo, path...)

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
		return &search.ResolverResult{Identifier: "NULL"}, nil
	case string:
		// check if it is a number field and explicitly try to cast to
		// float in case of a numeric string value was used
		// (this usually the case when the data is from a multipart/form-data request)
		field := r.baseCollection.Schema.GetFieldByName(path[len(path)-1])
		if field != nil && field.Type == schema.FieldTypeNumber {
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

	placeholder := "f" + security.PseudorandomString(5)

	return &search.ResolverResult{
		Identifier: "{:" + placeholder + "}",
		Params:     dbx.Params{placeholder: resultVal},
	}, nil
}

func (r *RecordFieldResolver) loadCollection(collectionNameOrId string) (*models.Collection, error) {
	// return already loaded
	for _, collection := range r.loadedCollections {
		if collection.Id == collectionNameOrId || strings.EqualFold(collection.Name, collectionNameOrId) {
			return collection, nil
		}
	}

	// load collection
	collection, err := r.dao.FindCollectionByNameOrId(collectionNameOrId)
	if err != nil {
		return nil, err
	}
	r.loadedCollections = append(r.loadedCollections, collection)

	return collection, nil
}

func (r *RecordFieldResolver) registerJoin(tableName string, tableAlias string, on dbx.Expression) {
	join := &join{
		tableName:  tableName,
		tableAlias: tableAlias,
		on:         on,
	}

	// replace existing join
	for i, j := range r.joins {
		if j.tableAlias == join.tableAlias {
			r.joins[i] = join
			return
		}
	}

	// register new join
	r.joins = append(r.joins, join)
}

func extractNestedMapVal(m map[string]any, keys ...string) (any, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("at least one key should be provided")
	}

	result, ok := m[keys[0]]
	if !ok {
		return nil, fmt.Errorf("invalid key path - missing key %q", keys[0])
	}

	// end key reached
	if len(keys) == 1 {
		return result, nil
	}

	if m, ok = result.(map[string]any); !ok {
		return nil, fmt.Errorf("expected map, got %#v", result)
	}

	return extractNestedMapVal(m, keys[1:]...)
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
		lengthModifier:
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unknown modifier in %q", combined)
}
