package resolvers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
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

// RecordFieldResolver defines a custom search resolver struct for
// managing Record model search fields.
//
// Usually used together with `search.Provider`. Example:
//  resolver := resolvers.NewRecordFieldResolver(
//      app.Dao(),
//      myCollection,
//      &models.RequestData{...},
//      true,
//  )
//  provider := search.NewProvider(resolver)
//  ...
type RecordFieldResolver struct {
	dao               *daos.Dao
	baseCollection    *models.Collection
	allowHiddenFields bool
	allowedFields     []string
	loadedCollections []*models.Collection
	joins             []*join // we cannot use a map because the insertion order is not preserved
	requestData       *models.RequestData
	staticRequestData map[string]any
}

// NewRecordFieldResolver creates and initializes a new `RecordFieldResolver`.
func NewRecordFieldResolver(
	dao *daos.Dao,
	baseCollection *models.Collection,
	requestData *models.RequestData,
	allowHiddenFields bool,
) *RecordFieldResolver {
	r := &RecordFieldResolver{
		dao:               dao,
		baseCollection:    baseCollection,
		requestData:       requestData,
		allowHiddenFields: allowHiddenFields,
		joins:             []*join{},
		loadedCollections: []*models.Collection{baseCollection},
		allowedFields: []string{
			`^\w+[\w\.\:]*$`,
			`^\@request\.method$`,
			`^\@request\.auth\.[\w\.\:]*\w+$`,
			`^\@request\.data\.[\w\.\:]*\w+$`,
			`^\@request\.query\.[\w\.\:]*\w+$`,
			`^\@collection\.\w+\.[\w\.\:]*\w+$`,
		},
	}

	r.staticRequestData = map[string]any{}
	if r.requestData != nil {
		r.staticRequestData["method"] = r.requestData.Method
		r.staticRequestData["query"] = r.requestData.Query
		r.staticRequestData["data"] = r.requestData.Data
		r.staticRequestData["auth"] = nil
		if r.requestData.AuthRecord != nil {
			r.requestData.AuthRecord.IgnoreEmailVisibility(true)
			r.staticRequestData["auth"] = r.requestData.AuthRecord.PublicExport()
			r.requestData.AuthRecord.IgnoreEmailVisibility(false)
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
//  id
//  someSelect.each
//  project.screen.status
//  @request.status
//  @request.query.filter
//  @request.auth.someRelation.name
//  @request.data.someRelation.name
//  @request.data.someField
//  @request.data.someSelect:each
//  @request.data.someField:isset
//  @collection.product.name
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
	resultVal, err := extractNestedMapVal(r.staticRequestData, path...)

	if modifier == issetModifier {
		if err != nil {
			return &search.ResolverResult{Identifier: "FALSE"}, nil
		}
		return &search.ResolverResult{Identifier: "TRUE"}, nil
	}

	// note: we are ignoring the error because requestData is dynamic
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

	var result any
	var ok bool

	if result, ok = m[keys[0]]; !ok {
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
