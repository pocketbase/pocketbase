package resolvers

import (
	"encoding/json"
	"fmt"
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

type join struct {
	id    string
	table string
	on    dbx.Expression
}

// RecordFieldResolver defines a custom search resolver struct for
// managing Record model search fields.
//
// Usually used together with `search.Provider`. Example:
//	resolver := resolvers.NewRecordFieldResolver(app.Dao(), myCollection, map[string]any{"test": 123})
//	provider := search.NewProvider(resolver)
//	...
type RecordFieldResolver struct {
	dao               *daos.Dao
	baseCollection    *models.Collection
	allowedFields     []string
	requestData       map[string]any
	joins             []join // we cannot use a map because the insertion order is not preserved
	loadedCollections []*models.Collection
}

// NewRecordFieldResolver creates and initializes a new `RecordFieldResolver`.
func NewRecordFieldResolver(
	dao *daos.Dao,
	baseCollection *models.Collection,
	requestData map[string]any,
) *RecordFieldResolver {
	return &RecordFieldResolver{
		dao:               dao,
		baseCollection:    baseCollection,
		requestData:       requestData,
		joins:             []join{},
		loadedCollections: []*models.Collection{baseCollection},
		allowedFields: []string{
			`^\w+[\w\.]*$`,
			`^\@request\.method$`,
			`^\@request\.user\.\w+[\w\.]*$`,
			`^\@request\.data\.\w+[\w\.]*$`,
			`^\@request\.query\.\w+[\w\.]*$`,
			`^\@collection\.\w+\.\w+[\w\.]*$`,
		},
	}
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

	return nil
}

// Resolve implements `search.FieldResolver` interface.
//
// Example of resolvable field formats:
//	id
//	project.screen.status
//	@request.status
//	@request.user.profile.someRelation.name
//	@collection.product.name
func (r *RecordFieldResolver) Resolve(fieldName string) (resultName string, placeholderParams dbx.Params, err error) {
	if len(r.allowedFields) > 0 && !list.ExistInSliceWithRegex(fieldName, r.allowedFields) {
		return "", nil, fmt.Errorf("Failed to resolve field %q", fieldName)
	}

	props := strings.Split(fieldName, ".")

	currentCollectionName := r.baseCollection.Name
	currentTableAlias := inflector.Columnify(currentCollectionName)

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

		r.addJoin(inflector.Columnify(collection.Name), currentTableAlias, nil)

		props = props[2:] // leave only the collection fields
	} else if props[0] == "@request" {
		// check for @request field
		if len(props) == 1 {
			return "", nil, fmt.Errorf("Invalid @request data field path in %q.", fieldName)
		}

		// not a profile relational field
		if !strings.HasPrefix(fieldName, "@request.user.profile.") {
			return r.resolveStaticRequestField(props[1:]...)
		}

		// resolve the profile collection fields
		currentCollectionName = models.ProfileCollectionName
		currentTableAlias = inflector.Columnify("__user_" + currentCollectionName)

		collection, err := r.loadCollection(currentCollectionName)
		if err != nil {
			return "", nil, fmt.Errorf("Failed to load collection %q from field path %q.", currentCollectionName, fieldName)
		}

		profileIdPlaceholder, profileIdPlaceholderParam, err := r.resolveStaticRequestField("user", "profile", "id")
		if err != nil {
			return "", nil, fmt.Errorf("Failed to resolve @request.user.profile.id path in %q.", fieldName)
		}
		if strings.ToLower(profileIdPlaceholder) == "null" {
			// the user doesn't have an associated profile
			return "NULL", nil, nil
		}

		// join the profile collection
		r.addJoin(
			inflector.Columnify(collection.Name),
			currentTableAlias,
			dbx.NewExp(fmt.Sprintf(
				// aka. profiles.id = profileId
				"[[%s.id]] = %s",
				currentTableAlias,
				profileIdPlaceholder,
			), profileIdPlaceholderParam),
		)

		props = props[3:] // leave only the profile fields
	}

	baseModelFields := schema.ReservedFieldNames()

	totalProps := len(props)

	for i, prop := range props {
		collection, err := r.loadCollection(currentCollectionName)
		if err != nil {
			return "", nil, fmt.Errorf("Failed to resolve field %q.", prop)
		}

		// base model prop (always available but not part of the collection schema)
		if list.ExistInSlice(prop, baseModelFields) {
			return fmt.Sprintf("[[%s.%s]]", currentTableAlias, inflector.Columnify(prop)), nil, nil
		}

		field := collection.Schema.GetFieldByName(prop)
		if field == nil {
			return "", nil, fmt.Errorf("Unrecognized field %q.", prop)
		}

		// last prop
		if i == totalProps-1 {
			return fmt.Sprintf("[[%s.%s]]", currentTableAlias, inflector.Columnify(prop)), nil, nil
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

		r.addJoin(
			fmt.Sprintf(
				// note: the case is used to normalize value access for single and multiple relations.
				`json_each(CASE WHEN json_valid([[%s]]) THEN [[%s]] ELSE json_array([[%s]]) END)`,
				jePair, jePair, jePair,
			),
			jeTable,
			nil,
		)
		r.addJoin(
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
	resultVal, _ := extractNestedMapVal(r.requestData, path...)

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

	placeholder := "f" + security.RandomString(7)
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
		if collection.Name == collectionNameOrId || collection.Id == collectionNameOrId {
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

func (r *RecordFieldResolver) addJoin(tableName string, tableAlias string, on dbx.Expression) {
	tableExpr := fmt.Sprintf("%s %s", tableName, tableAlias)

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
