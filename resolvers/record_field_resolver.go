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
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

// ensure that `search.FieldResolver` interface is implemented
var _ search.FieldResolver = (*RecordFieldResolver)(nil)

// list of auth filter fields that don't require join with the auth
// collection or any other extra checks to be resolved
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

type join struct {
	id    string
	table string
	on    dbx.Expression
}

// RecordFieldResolver defines a custom search resolver struct for
// managing Record model search fields.
//
// Usually used together with `search.Provider`. Example:
//	resolver := resolvers.NewRecordFieldResolver(
//		app.Dao(),
//		myCollection,
//		&models.RequestData{...},
//		true,
//	)
//	provider := search.NewProvider(resolver)
//	...
type RecordFieldResolver struct {
	dao               *daos.Dao
	baseCollection    *models.Collection
	allowHiddenFields bool
	allowedFields     []string
	loadedCollections []*models.Collection
	joins             []join // we cannot use a map because the insertion order is not preserved
	exprs             []dbx.Expression
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
		joins:             []join{},
		exprs:             []dbx.Expression{},
		loadedCollections: []*models.Collection{baseCollection},
		allowedFields: []string{
			`^\w+[\w\.]*$`,
			`^\@request\.method$`,
			`^\@request\.auth\.\w+[\w\.]*$`,
			`^\@request\.data\.\w+[\w\.]*$`,
			`^\@request\.query\.\w+[\w\.]*$`,
			`^\@collection\.\w+\.\w+[\w\.]*$`,
		},
	}

	// @todo remove after IN operator and multi-match filter enhancements
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
			query.LeftJoin(join.table, join.on)
		}
	}

	for _, expr := range r.exprs {
		if expr != nil {
			query.AndWhere(expr)
		}
	}

	return nil
}

// Resolve implements `search.FieldResolver` interface.
//
// Example of resolvable field formats:
//	id
//	project.screen.status
//	@request.status
//	@request.auth.someRelation.name
//	@collection.product.name
func (r *RecordFieldResolver) Resolve(fieldName string) (resultName string, placeholderParams dbx.Params, err error) {
	if len(r.allowedFields) > 0 && !list.ExistInSliceWithRegex(fieldName, r.allowedFields) {
		return "", nil, fmt.Errorf("Failed to resolve field %q", fieldName)
	}

	props := strings.Split(fieldName, ".")

	currentCollectionName := r.baseCollection.Name
	currentTableAlias := inflector.Columnify(currentCollectionName)

	// flag indicating whether to return null on missing field or return on an error
	nullifyMisingField := false

	allowHiddenFields := r.allowHiddenFields

	// check for @collection field (aka. non-relational join)
	// must be in the format "@collection.COLLECTION_NAME.FIELD[.FIELD2....]"
	if props[0] == "@collection" {
		if len(props) < 3 {
			return "", nil, fmt.Errorf("Invalid @collection field path in %q.", fieldName)
		}

		currentCollectionName = props[1]
		currentTableAlias = inflector.Columnify("__collection_" + currentCollectionName)

		collection, err := r.loadCollection(currentCollectionName)
		if err != nil {
			return "", nil, fmt.Errorf("Failed to load collection %q from field path %q.", currentCollectionName, fieldName)
		}

		// always allow hidden fields since the @collection.* filter is a system one
		allowHiddenFields = true

		r.registerJoin(inflector.Columnify(collection.Name), currentTableAlias, nil)

		props = props[2:] // leave only the collection fields
	} else if props[0] == "@request" {
		if len(props) == 1 {
			return "", nil, fmt.Errorf("Invalid @request data field path in %q.", fieldName)
		}

		if r.requestData == nil {
			return "NULL", nil, nil
		}

		// plain @request.* field
		if !strings.HasPrefix(fieldName, "@request.auth.") || list.ExistInSlice(fieldName, plainRequestAuthFields) {
			return r.resolveStaticRequestField(props[1:]...)
		}

		// always allow hidden fields since the @request.* filter is a system one
		allowHiddenFields = true

		// enable the ignore flag for missing @request.auth.* fields
		// for consistency with @request.data.* and @request.query.*
		nullifyMisingField = true

		// resolve the auth collection fields
		// ---
		if r.requestData == nil || r.requestData.AuthRecord == nil || r.requestData.AuthRecord.Collection() == nil {
			return "NULL", nil, nil
		}

		collection := r.requestData.AuthRecord.Collection()
		r.loadedCollections = append(r.loadedCollections, collection)

		currentCollectionName = collection.Name
		currentTableAlias = "__auth_" + inflector.Columnify(currentCollectionName)

		authIdParamKey := "auth" + security.PseudorandomString(5)
		authIdParams := dbx.Params{authIdParamKey: r.requestData.AuthRecord.Id}
		// ---

		// join the auth collection
		r.registerJoin(
			inflector.Columnify(collection.Name),
			currentTableAlias,
			dbx.NewExp(fmt.Sprintf(
				// aka. __auth_users.id = :userId
				"[[%s.id]] = {:%s}",
				inflector.Columnify(currentTableAlias),
				authIdParamKey,
			), authIdParams),
		)

		props = props[2:] // leave only the auth relation fields
	}

	totalProps := len(props)

	for i, prop := range props {
		collection, err := r.loadCollection(currentCollectionName)
		if err != nil {
			return "", nil, fmt.Errorf("Failed to resolve field %q.", prop)
		}

		systemFieldNames := schema.BaseModelFieldNames()
		if collection.IsAuth() {
			systemFieldNames = append(
				systemFieldNames,
				schema.FieldNameUsername,
				schema.FieldNameVerified,
				schema.FieldNameEmailVisibility,
				schema.FieldNameEmail,
			)
		}

		// internal model prop (always available but not part of the collection schema)
		if list.ExistInSlice(prop, systemFieldNames) {
			// allow querying only auth records with emails marked as public
			if prop == schema.FieldNameEmail && !allowHiddenFields {
				r.registerExpr(dbx.NewExp(fmt.Sprintf(
					"[[%s.%s]] = TRUE",
					currentTableAlias,
					inflector.Columnify(schema.FieldNameEmailVisibility),
				)))
			}

			return fmt.Sprintf("[[%s.%s]]", currentTableAlias, inflector.Columnify(prop)), nil, nil
		}

		field := collection.Schema.GetFieldByName(prop)
		if field == nil {
			if nullifyMisingField {
				return "NULL", nil, nil
			}

			return "", nil, fmt.Errorf("Unrecognized field %q.", prop)
		}

		// last prop
		if i == totalProps-1 {
			return fmt.Sprintf("[[%s.%s]]", currentTableAlias, inflector.Columnify(prop)), nil, nil
		}

		// check if it is a json field
		if field.Type == schema.FieldTypeJson {
			var jsonPath strings.Builder
			jsonPath.WriteString("$")
			for _, p := range props[i+1:] {
				if _, err := strconv.Atoi(p); err == nil {
					jsonPath.WriteString("[")
					jsonPath.WriteString(inflector.Columnify(p))
					jsonPath.WriteString("]")
				} else {
					jsonPath.WriteString(".")
					jsonPath.WriteString(inflector.Columnify(p))
				}
			}
			return fmt.Sprintf(
				"JSON_EXTRACT([[%s.%s]], '%s')",
				currentTableAlias,
				inflector.Columnify(prop),
				jsonPath.String(),
			), nil, nil
		}

		// check if it is a relation field
		if field.Type != schema.FieldTypeRelation {
			return "", nil, fmt.Errorf("Field %q is not a valid relation.", prop)
		}

		// auto join the relation
		// ---
		field.InitOptions()
		options, ok := field.Options.(*schema.RelationOptions)
		if !ok {
			return "", nil, fmt.Errorf("Failed to initialize field %q options.", prop)
		}

		relCollection, relErr := r.loadCollection(options.CollectionId)
		if relErr != nil {
			return "", nil, fmt.Errorf("Failed to find field %q collection.", prop)
		}

		cleanFieldName := inflector.Columnify(field.Name)
		newCollectionName := relCollection.Name
		newTableAlias := currentTableAlias + "_" + cleanFieldName

		jeTable := currentTableAlias + "_" + cleanFieldName + "_je"
		jePair := currentTableAlias + "." + cleanFieldName

		r.registerJoin(
			fmt.Sprintf(
				// note: the case is used to normalize value access for single and multiple relations.
				`json_each(CASE WHEN json_valid([[%s]]) THEN [[%s]] ELSE json_array([[%s]]) END)`,
				jePair, jePair, jePair,
			),
			jeTable,
			nil,
		)
		r.registerJoin(
			inflector.Columnify(newCollectionName),
			newTableAlias,
			dbx.NewExp(fmt.Sprintf("[[%s.id]] = [[%s.value]]", newTableAlias, jeTable)),
		)

		currentCollectionName = newCollectionName
		currentTableAlias = newTableAlias
	}

	return "", nil, fmt.Errorf("Failed to resolve field %q.", fieldName)
}

func (r *RecordFieldResolver) resolveStaticRequestField(path ...string) (resultName string, placeholderParams dbx.Params, err error) {
	// ignore error because requestData is dynamic and some of the
	// lookup keys may not be defined for the request
	resultVal, _ := extractNestedMapVal(r.staticRequestData, path...)

	switch v := resultVal.(type) {
	case nil:
		return "NULL", nil, nil
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
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
	name := fmt.Sprintf("{:%s}", placeholder)
	params := dbx.Params{placeholder: resultVal}

	return name, params, nil
}

func extractNestedMapVal(m map[string]any, keys ...string) (result any, err error) {
	var ok bool

	if len(keys) == 0 {
		return nil, fmt.Errorf("At least one key should be provided.")
	}

	if result, ok = m[keys[0]]; !ok {
		return nil, fmt.Errorf("Invalid key path - missing key %q.", keys[0])
	}

	// end key reached
	if len(keys) == 1 {
		return result, nil
	}

	if m, ok = result.(map[string]any); !ok {
		return nil, fmt.Errorf("Expected map structure, got %#v.", result)
	}

	return extractNestedMapVal(m, keys[1:]...)
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
	tableExpr := (tableName + " " + tableAlias)

	join := join{
		id:    tableAlias,
		table: tableExpr,
		on:    on,
	}

	// replace existing join
	for i, j := range r.joins {
		if j.id == join.id {
			r.joins[i] = join
			return
		}
	}

	// register new join
	r.joins = append(r.joins, join)
}

func (r *RecordFieldResolver) registerExpr(expr dbx.Expression) {
	r.exprs = append(r.exprs, expr)
}
