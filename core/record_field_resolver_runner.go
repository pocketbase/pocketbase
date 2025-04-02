package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

// maxNestedRels defines the max allowed nested relations depth.
const maxNestedRels = 6

// list of auth filter fields that don't require join with the auth
// collection or any other extra checks to be resolved.
var plainRequestAuthFields = map[string]struct{}{
	"@request.auth." + FieldNameId:              {},
	"@request.auth." + FieldNameCollectionId:    {},
	"@request.auth." + FieldNameCollectionName:  {},
	"@request.auth." + FieldNameEmail:           {},
	"@request.auth." + FieldNameEmailVisibility: {},
	"@request.auth." + FieldNameVerified:        {},
}

// parseAndRun starts a new one-off RecordFieldResolver.Resolve execution.
func parseAndRun(fieldName string, resolver *RecordFieldResolver) (*search.ResolverResult, error) {
	r := &runner{
		fieldName: fieldName,
		resolver:  resolver,
	}

	return r.run()
}

type runner struct {
	used      bool                 // indicates whether the runner was already executed
	resolver  *RecordFieldResolver // resolver is the shared expression fields resolver
	fieldName string               // the name of the single field expression the runner is responsible for

	// shared processing state
	// ---------------------------------------------------------------
	activeProps                []string            // holds the active props that remains to be processed
	activeCollectionName       string              // the last used collection name
	activeTableAlias           string              // the last used table alias
	allowHiddenFields          bool                // indicates whether hidden fields (eg. email) should be allowed without extra conditions
	nullifyMisingField         bool                // indicating whether to return null on missing field or return an error
	withMultiMatch             bool                // indicates whether to attach a multiMatchSubquery condition to the ResolverResult
	multiMatchActiveTableAlias string              // the last used multi-match table alias
	multiMatch                 *multiMatchSubquery // the multi-match subquery expression generated from the fieldName
}

func (r *runner) run() (*search.ResolverResult, error) {
	if r.used {
		return nil, errors.New("the runner was already used")
	}

	if len(r.resolver.allowedFields) > 0 && !list.ExistInSliceWithRegex(r.fieldName, r.resolver.allowedFields) {
		return nil, fmt.Errorf("failed to resolve field %q", r.fieldName)
	}

	defer func() {
		r.used = true
	}()

	r.prepare()

	// check for @collection field (aka. non-relational join)
	// must be in the format "@collection.COLLECTION_NAME.FIELD[.FIELD2....]"
	if r.activeProps[0] == "@collection" {
		return r.processCollectionField()
	}

	if r.activeProps[0] == "@request" {
		if r.resolver.requestInfo == nil {
			return &search.ResolverResult{Identifier: "NULL"}, nil
		}

		if strings.HasPrefix(r.fieldName, "@request.auth.") {
			return r.processRequestAuthField()
		}

		if strings.HasPrefix(r.fieldName, "@request.body.") && len(r.activeProps) > 2 {
			name, modifier, err := splitModifier(r.activeProps[2])
			if err != nil {
				return nil, err
			}

			bodyField := r.resolver.baseCollection.Fields.GetByName(name)
			if bodyField == nil {
				return r.resolver.resolveStaticRequestField(r.activeProps[1:]...)
			}

			// check for body relation field
			if bodyField.Type() == FieldTypeRelation && len(r.activeProps) > 3 {
				return r.processRequestInfoRelationField(bodyField)
			}

			// check for body arrayble fields ":each" modifier
			if modifier == eachModifier && len(r.activeProps) == 3 {
				return r.processRequestInfoEachModifier(bodyField)
			}

			// check for body arrayble fields ":length" modifier
			if modifier == lengthModifier && len(r.activeProps) == 3 {
				return r.processRequestInfoLengthModifier(bodyField)
			}

			// check for body arrayble fields ":lower" modifier
			if modifier == lowerModifier && len(r.activeProps) == 3 {
				return r.processRequestInfoLowerModifier(bodyField)
			}
		}

		// some other @request.* static field
		return r.resolver.resolveStaticRequestField(r.activeProps[1:]...)
	}

	// regular field
	return r.processActiveProps()
}

func (r *runner) prepare() {
	r.activeProps = strings.Split(r.fieldName, ".")

	r.activeCollectionName = r.resolver.baseCollection.Name
	r.activeTableAlias = inflector.Columnify(r.activeCollectionName)

	r.allowHiddenFields = r.resolver.allowHiddenFields
	// always allow hidden fields since the @.* filter is a system one
	if r.activeProps[0] == "@collection" || r.activeProps[0] == "@request" {
		r.allowHiddenFields = true
	}

	// enable the ignore flag for missing @request.* fields for backward
	// compatibility and consistency with all @request.* filter fields and types
	r.nullifyMisingField = r.activeProps[0] == "@request"

	// prepare a multi-match subquery
	r.multiMatch = &multiMatchSubquery{
		baseTableAlias: r.activeTableAlias,
		params:         dbx.Params{},
	}
	r.multiMatch.fromTableName = inflector.Columnify(r.activeCollectionName)
	r.multiMatch.fromTableAlias = "__mm_" + r.activeTableAlias
	r.multiMatchActiveTableAlias = r.multiMatch.fromTableAlias
	r.withMultiMatch = false
}

func (r *runner) processCollectionField() (*search.ResolverResult, error) {
	if len(r.activeProps) < 3 {
		return nil, fmt.Errorf("invalid @collection field path in %q", r.fieldName)
	}

	// nameOrId or nameOrId:alias
	collectionParts := strings.SplitN(r.activeProps[1], ":", 2)

	collection, err := r.resolver.loadCollection(collectionParts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to load collection %q from field path %q", r.activeProps[1], r.fieldName)
	}

	r.activeCollectionName = collection.Name

	if len(collectionParts) == 2 && collectionParts[1] != "" {
		r.activeTableAlias = inflector.Columnify("__collection_alias_" + collectionParts[1])
	} else {
		r.activeTableAlias = inflector.Columnify("__collection_" + r.activeCollectionName)
	}

	r.withMultiMatch = true

	// join the collection to the main query
	r.resolver.registerJoin(inflector.Columnify(collection.Name), r.activeTableAlias, nil)

	// join the collection to the multi-match subquery
	r.multiMatchActiveTableAlias = "__mm" + r.activeTableAlias
	r.multiMatch.joins = append(r.multiMatch.joins, &join{
		tableName:  inflector.Columnify(collection.Name),
		tableAlias: r.multiMatchActiveTableAlias,
	})

	// leave only the collection fields
	// aka. @collection.someCollection.fieldA.fieldB -> fieldA.fieldB
	r.activeProps = r.activeProps[2:]

	return r.processActiveProps()
}

func (r *runner) processRequestAuthField() (*search.ResolverResult, error) {
	if r.resolver.requestInfo == nil || r.resolver.requestInfo.Auth == nil || r.resolver.requestInfo.Auth.Collection() == nil {
		return &search.ResolverResult{Identifier: "NULL"}, nil
	}

	// plain auth field
	// ---
	if _, ok := plainRequestAuthFields[r.fieldName]; ok {
		return r.resolver.resolveStaticRequestField(r.activeProps[1:]...)
	}

	// resolve the auth collection field
	// ---
	collection := r.resolver.requestInfo.Auth.Collection()

	r.activeCollectionName = collection.Name
	r.activeTableAlias = "__auth_" + inflector.Columnify(r.activeCollectionName)

	// join the auth collection to the main query
	r.resolver.registerJoin(
		inflector.Columnify(r.activeCollectionName),
		r.activeTableAlias,
		dbx.HashExp{
			// aka. __auth_users.id = :userId
			(r.activeTableAlias + ".id"): r.resolver.requestInfo.Auth.Id,
		},
	)

	// join the auth collection to the multi-match subquery
	r.multiMatchActiveTableAlias = "__mm_" + r.activeTableAlias
	r.multiMatch.joins = append(
		r.multiMatch.joins,
		&join{
			tableName:  inflector.Columnify(r.activeCollectionName),
			tableAlias: r.multiMatchActiveTableAlias,
			on: dbx.HashExp{
				(r.multiMatchActiveTableAlias + ".id"): r.resolver.requestInfo.Auth.Id,
			},
		},
	)

	// leave only the auth relation fields
	// aka. @request.auth.fieldA.fieldB -> fieldA.fieldB
	r.activeProps = r.activeProps[2:]

	return r.processActiveProps()
}

// note: nil value is returned as empty slice
func toSlice(value any) []any {
	if value == nil {
		return []any{}
	}

	rv := reflect.ValueOf(value)

	kind := rv.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return []any{value}
	}

	rvLen := rv.Len()

	result := make([]interface{}, rvLen)

	for i := 0; i < rvLen; i++ {
		result[i] = rv.Index(i).Interface()
	}

	return result
}

func (r *runner) processRequestInfoLowerModifier(bodyField Field) (*search.ResolverResult, error) {
	rawValue := cast.ToString(r.resolver.requestInfo.Body[bodyField.GetName()])

	placeholder := "infoLower" + bodyField.GetName() + security.PseudorandomString(6)

	result := &search.ResolverResult{
		Identifier: "LOWER({:" + placeholder + "})",
		Params:     dbx.Params{placeholder: rawValue},
	}

	return result, nil
}

func (r *runner) processRequestInfoLengthModifier(bodyField Field) (*search.ResolverResult, error) {
	if _, ok := bodyField.(MultiValuer); !ok {
		return nil, fmt.Errorf("field %q doesn't support multivalue operations", bodyField.GetName())
	}

	bodyItems := toSlice(r.resolver.requestInfo.Body[bodyField.GetName()])

	result := &search.ResolverResult{
		Identifier: strconv.Itoa(len(bodyItems)),
	}

	return result, nil
}

func (r *runner) processRequestInfoEachModifier(bodyField Field) (*search.ResolverResult, error) {
	multiValuer, ok := bodyField.(MultiValuer)
	if !ok {
		return nil, fmt.Errorf("field %q doesn't support multivalue operations", bodyField.GetName())
	}

	bodyItems := toSlice(r.resolver.requestInfo.Body[bodyField.GetName()])
	bodyItemsRaw, err := json.Marshal(bodyItems)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize the data for field %q", r.activeProps[2])
	}

	placeholder := "dataEach" + security.PseudorandomString(6)
	cleanFieldName := inflector.Columnify(bodyField.GetName())
	jeTable := fmt.Sprintf("json_each({:%s})", placeholder)
	jeAlias := "__dataEach_" + cleanFieldName + "_je"
	r.resolver.registerJoin(jeTable, jeAlias, nil)

	result := &search.ResolverResult{
		Identifier: fmt.Sprintf("[[%s.value]]", jeAlias),
		Params:     dbx.Params{placeholder: bodyItemsRaw},
	}

	if multiValuer.IsMultiple() {
		r.withMultiMatch = true
	}

	if r.withMultiMatch {
		placeholder2 := "mm" + placeholder
		jeTable2 := fmt.Sprintf("json_each({:%s})", placeholder2)
		jeAlias2 := "__mm" + jeAlias

		r.multiMatch.joins = append(r.multiMatch.joins, &join{
			tableName:  jeTable2,
			tableAlias: jeAlias2,
		})
		r.multiMatch.params[placeholder2] = bodyItemsRaw
		r.multiMatch.valueIdentifier = fmt.Sprintf("[[%s.value]]", jeAlias2)

		result.MultiMatchSubQuery = r.multiMatch
	}

	return result, nil
}

func (r *runner) processRequestInfoRelationField(bodyField Field) (*search.ResolverResult, error) {
	relField, ok := bodyField.(*RelationField)
	if !ok {
		return nil, fmt.Errorf("failed to initialize data relation field %q", bodyField.GetName())
	}

	dataRelCollection, err := r.resolver.loadCollection(relField.CollectionId)
	if err != nil {
		return nil, fmt.Errorf("failed to load collection %q from data field %q", relField.CollectionId, relField.Name)
	}

	var dataRelIds []string
	if r.resolver.requestInfo != nil && len(r.resolver.requestInfo.Body) != 0 {
		dataRelIds = list.ToUniqueStringSlice(r.resolver.requestInfo.Body[relField.Name])
	}
	if len(dataRelIds) == 0 {
		return &search.ResolverResult{Identifier: "NULL"}, nil
	}

	r.activeCollectionName = dataRelCollection.Name
	r.activeTableAlias = inflector.Columnify("__data_" + dataRelCollection.Name + "_" + relField.Name)

	// join the data rel collection to the main collection
	r.resolver.registerJoin(
		r.activeCollectionName,
		r.activeTableAlias,
		dbx.In(
			fmt.Sprintf("[[%s.id]]", r.activeTableAlias),
			list.ToInterfaceSlice(dataRelIds)...,
		),
	)

	if relField.IsMultiple() {
		r.withMultiMatch = true
	}

	// join the data rel collection to the multi-match subquery
	r.multiMatchActiveTableAlias = inflector.Columnify("__data_mm_" + dataRelCollection.Name + "_" + relField.Name)
	r.multiMatch.joins = append(
		r.multiMatch.joins,
		&join{
			tableName:  r.activeCollectionName,
			tableAlias: r.multiMatchActiveTableAlias,
			on: dbx.In(
				fmt.Sprintf("[[%s.id]]", r.multiMatchActiveTableAlias),
				list.ToInterfaceSlice(dataRelIds)...,
			),
		},
	)

	// leave only the data relation fields
	// aka. @request.body.someRel.fieldA.fieldB -> fieldA.fieldB
	r.activeProps = r.activeProps[3:]

	return r.processActiveProps()
}

var viaRegex = regexp.MustCompile(`^(\w+)_via_(\w+)$`)

func (r *runner) processActiveProps() (*search.ResolverResult, error) {
	totalProps := len(r.activeProps)

	for i, prop := range r.activeProps {
		collection, err := r.resolver.loadCollection(r.activeCollectionName)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve field %q", prop)
		}

		// last prop
		if i == totalProps-1 {
			return r.processLastProp(collection, prop)
		}

		field := collection.Fields.GetByName(prop)

		if field != nil && field.GetHidden() && !r.allowHiddenFields {
			return nil, fmt.Errorf("non-filterable field %q", prop)
		}

		// json or geoPoint field -> treat the rest of the props as json path
		// @todo consider converting to "JSONExtractable" interface with optional extra validation for the remaining props?
		if field != nil && (field.Type() == FieldTypeJSON || field.Type() == FieldTypeGeoPoint) {
			var jsonPath strings.Builder
			for j, p := range r.activeProps[i+1:] {
				if _, err := strconv.Atoi(p); err == nil {
					jsonPath.WriteString("[")
					jsonPath.WriteString(inflector.Columnify(p))
					jsonPath.WriteString("]")
				} else {
					if j > 0 {
						jsonPath.WriteString(".")
					}
					jsonPath.WriteString(inflector.Columnify(p))
				}
			}
			jsonPathStr := jsonPath.String()

			result := &search.ResolverResult{
				NoCoalesce: true,
				Identifier: dbutils.JSONExtract(r.activeTableAlias+"."+inflector.Columnify(prop), jsonPathStr),
			}

			if r.withMultiMatch {
				r.multiMatch.valueIdentifier = dbutils.JSONExtract(r.multiMatchActiveTableAlias+"."+inflector.Columnify(prop), jsonPathStr)
				result.MultiMatchSubQuery = r.multiMatch
			}

			return result, nil
		}

		if i >= maxNestedRels {
			return nil, fmt.Errorf("max nested relations reached for field %q", prop)
		}

		// check for back relation (eg. yourCollection_via_yourRelField)
		// -----------------------------------------------------------
		if field == nil {
			parts := viaRegex.FindStringSubmatch(prop)
			if len(parts) != 3 {
				if r.nullifyMisingField {
					return &search.ResolverResult{Identifier: "NULL"}, nil
				}
				return nil, fmt.Errorf("failed to resolve field %q", prop)
			}

			backCollection, err := r.resolver.loadCollection(parts[1])
			if err != nil {
				if r.nullifyMisingField {
					return &search.ResolverResult{Identifier: "NULL"}, nil
				}
				return nil, fmt.Errorf("failed to load back relation field %q collection", prop)
			}

			backField := backCollection.Fields.GetByName(parts[2])
			if backField == nil {
				if r.nullifyMisingField {
					return &search.ResolverResult{Identifier: "NULL"}, nil
				}
				return nil, fmt.Errorf("missing back relation field %q", parts[2])
			}

			if backField.Type() != FieldTypeRelation {
				if r.nullifyMisingField {
					return &search.ResolverResult{Identifier: "NULL"}, nil
				}
				return nil, fmt.Errorf("invalid back relation field %q", parts[2])
			}

			if backField.GetHidden() && !r.allowHiddenFields {
				return nil, fmt.Errorf("non-filterable back relation field %q", backField.GetName())
			}

			backRelField, ok := backField.(*RelationField)
			if !ok {
				return nil, fmt.Errorf("failed to initialize back relation field %q", backField.GetName())
			}
			if backRelField.CollectionId != collection.Id {
				// https://github.com/pocketbase/pocketbase/discussions/6590#discussioncomment-12496581
				if r.nullifyMisingField {
					return &search.ResolverResult{Identifier: "NULL"}, nil
				}
				return nil, fmt.Errorf("invalid collection reference of a back relation field %q", backField.GetName())
			}

			// join the back relation to the main query
			// ---
			cleanProp := inflector.Columnify(prop)
			cleanBackFieldName := inflector.Columnify(backRelField.Name)
			newTableAlias := r.activeTableAlias + "_" + cleanProp
			newCollectionName := inflector.Columnify(backCollection.Name)

			isBackRelMultiple := backRelField.IsMultiple()
			if !isBackRelMultiple {
				// additionally check if the rel field has a single column unique index
				_, hasUniqueIndex := dbutils.FindSingleColumnUniqueIndex(backCollection.Indexes, backRelField.Name)
				isBackRelMultiple = !hasUniqueIndex
			}

			if !isBackRelMultiple {
				r.resolver.registerJoin(
					newCollectionName,
					newTableAlias,
					dbx.NewExp(fmt.Sprintf("[[%s.%s]] = [[%s.id]]", newTableAlias, cleanBackFieldName, r.activeTableAlias)),
				)
			} else {
				jeAlias := r.activeTableAlias + "_" + cleanProp + "_je"
				r.resolver.registerJoin(
					newCollectionName,
					newTableAlias,
					dbx.NewExp(fmt.Sprintf(
						"[[%s.id]] IN (SELECT [[%s.value]] FROM %s {{%s}})",
						r.activeTableAlias,
						jeAlias,
						dbutils.JSONEach(newTableAlias+"."+cleanBackFieldName),
						jeAlias,
					)),
				)
			}

			r.activeCollectionName = newCollectionName
			r.activeTableAlias = newTableAlias
			// ---

			// join the back relation to the multi-match subquery
			// ---
			if isBackRelMultiple {
				r.withMultiMatch = true // enable multimatch if not already
			}

			newTableAlias2 := r.multiMatchActiveTableAlias + "_" + cleanProp

			if !isBackRelMultiple {
				r.multiMatch.joins = append(
					r.multiMatch.joins,
					&join{
						tableName:  newCollectionName,
						tableAlias: newTableAlias2,
						on:         dbx.NewExp(fmt.Sprintf("[[%s.%s]] = [[%s.id]]", newTableAlias2, cleanBackFieldName, r.multiMatchActiveTableAlias)),
					},
				)
			} else {
				jeAlias2 := r.multiMatchActiveTableAlias + "_" + cleanProp + "_je"
				r.multiMatch.joins = append(
					r.multiMatch.joins,
					&join{
						tableName:  newCollectionName,
						tableAlias: newTableAlias2,
						on: dbx.NewExp(fmt.Sprintf(
							"[[%s.id]] IN (SELECT [[%s.value]] FROM %s {{%s}})",
							r.multiMatchActiveTableAlias,
							jeAlias2,
							dbutils.JSONEach(newTableAlias2+"."+cleanBackFieldName),
							jeAlias2,
						)),
					},
				)
			}

			r.multiMatchActiveTableAlias = newTableAlias2
			// ---

			continue
		}
		// -----------------------------------------------------------

		// check for direct relation
		if field.Type() != FieldTypeRelation {
			return nil, fmt.Errorf("field %q is not a valid relation", prop)
		}

		// join the relation to the main query
		// ---
		relField, ok := field.(*RelationField)
		if !ok {
			return nil, fmt.Errorf("failed to initialize relation field %q", prop)
		}

		relCollection, relErr := r.resolver.loadCollection(relField.CollectionId)
		if relErr != nil {
			return nil, fmt.Errorf("failed to load field %q collection", prop)
		}

		// "id" lookups optimization for single relations to avoid unnecessary joins,
		// aka. "user.id" and "user" should produce the same query identifier
		if !relField.IsMultiple() &&
			// the penultimate prop is "id"
			i == totalProps-2 && r.activeProps[i+1] == FieldNameId {
			return r.processLastProp(collection, relField.Name)
		}

		cleanFieldName := inflector.Columnify(relField.Name)
		prefixedFieldName := r.activeTableAlias + "." + cleanFieldName
		newTableAlias := r.activeTableAlias + "_" + cleanFieldName
		newCollectionName := relCollection.Name

		if !relField.IsMultiple() {
			r.resolver.registerJoin(
				inflector.Columnify(newCollectionName),
				newTableAlias,
				dbx.NewExp(fmt.Sprintf("[[%s.id]] = [[%s]]", newTableAlias, prefixedFieldName)),
			)
		} else {
			jeAlias := r.activeTableAlias + "_" + cleanFieldName + "_je"
			r.resolver.registerJoin(dbutils.JSONEach(prefixedFieldName), jeAlias, nil)
			r.resolver.registerJoin(
				inflector.Columnify(newCollectionName),
				newTableAlias,
				dbx.NewExp(fmt.Sprintf("[[%s.id]] = [[%s.value]]", newTableAlias, jeAlias)),
			)
		}

		r.activeCollectionName = newCollectionName
		r.activeTableAlias = newTableAlias
		// ---

		// join the relation to the multi-match subquery
		// ---
		if relField.IsMultiple() {
			r.withMultiMatch = true // enable multimatch if not already
		}

		newTableAlias2 := r.multiMatchActiveTableAlias + "_" + cleanFieldName
		prefixedFieldName2 := r.multiMatchActiveTableAlias + "." + cleanFieldName

		if !relField.IsMultiple() {
			r.multiMatch.joins = append(
				r.multiMatch.joins,
				&join{
					tableName:  inflector.Columnify(newCollectionName),
					tableAlias: newTableAlias2,
					on:         dbx.NewExp(fmt.Sprintf("[[%s.id]] = [[%s]]", newTableAlias2, prefixedFieldName2)),
				},
			)
		} else {
			jeAlias2 := r.multiMatchActiveTableAlias + "_" + cleanFieldName + "_je"
			r.multiMatch.joins = append(
				r.multiMatch.joins,
				&join{
					tableName:  dbutils.JSONEach(prefixedFieldName2),
					tableAlias: jeAlias2,
				},
				&join{
					tableName:  inflector.Columnify(newCollectionName),
					tableAlias: newTableAlias2,
					on:         dbx.NewExp(fmt.Sprintf("[[%s.id]] = [[%s.value]]", newTableAlias2, jeAlias2)),
				},
			)
		}

		r.multiMatchActiveTableAlias = newTableAlias2
		// ---
	}

	return nil, fmt.Errorf("failed to resolve field %q", r.fieldName)
}

func (r *runner) processLastProp(collection *Collection, prop string) (*search.ResolverResult, error) {
	name, modifier, err := splitModifier(prop)
	if err != nil {
		return nil, err
	}

	field := collection.Fields.GetByName(name)
	if field == nil {
		if r.nullifyMisingField {
			return &search.ResolverResult{Identifier: "NULL"}, nil
		}
		return nil, fmt.Errorf("unknown field %q", name)
	}

	if field.GetHidden() && !r.allowHiddenFields {
		return nil, fmt.Errorf("non-filterable field %q", name)
	}

	multvaluer, isMultivaluer := field.(MultiValuer)

	cleanFieldName := inflector.Columnify(field.GetName())

	// arrayable fields with ":length" modifier
	// -------------------------------------------------------
	if modifier == lengthModifier && isMultivaluer {
		jePair := r.activeTableAlias + "." + cleanFieldName

		result := &search.ResolverResult{
			Identifier: dbutils.JSONArrayLength(jePair),
		}

		if r.withMultiMatch {
			jePair2 := r.multiMatchActiveTableAlias + "." + cleanFieldName
			r.multiMatch.valueIdentifier = dbutils.JSONArrayLength(jePair2)
			result.MultiMatchSubQuery = r.multiMatch
		}

		return result, nil
	}

	// arrayable fields with ":each" modifier
	// -------------------------------------------------------
	if modifier == eachModifier && isMultivaluer {
		jePair := r.activeTableAlias + "." + cleanFieldName
		jeAlias := r.activeTableAlias + "_" + cleanFieldName + "_je"
		r.resolver.registerJoin(dbutils.JSONEach(jePair), jeAlias, nil)

		result := &search.ResolverResult{
			Identifier: fmt.Sprintf("[[%s.value]]", jeAlias),
		}

		if multvaluer.IsMultiple() {
			r.withMultiMatch = true
		}

		if r.withMultiMatch {
			jePair2 := r.multiMatchActiveTableAlias + "." + cleanFieldName
			jeAlias2 := r.multiMatchActiveTableAlias + "_" + cleanFieldName + "_je"

			r.multiMatch.joins = append(r.multiMatch.joins, &join{
				tableName:  dbutils.JSONEach(jePair2),
				tableAlias: jeAlias2,
			})
			r.multiMatch.valueIdentifier = fmt.Sprintf("[[%s.value]]", jeAlias2)

			result.MultiMatchSubQuery = r.multiMatch
		}

		return result, nil
	}

	// default
	// -------------------------------------------------------
	result := &search.ResolverResult{
		Identifier: "[[" + r.activeTableAlias + "." + cleanFieldName + "]]",
	}

	if r.withMultiMatch {
		r.multiMatch.valueIdentifier = "[[" + r.multiMatchActiveTableAlias + "." + cleanFieldName + "]]"
		result.MultiMatchSubQuery = r.multiMatch
	}

	// allow querying only auth records with emails marked as public
	if field.GetName() == FieldNameEmail && !r.allowHiddenFields && collection.IsAuth() {
		result.AfterBuild = func(expr dbx.Expression) dbx.Expression {
			return dbx.Enclose(dbx.And(expr, dbx.NewExp(fmt.Sprintf(
				"[[%s.%s]] = TRUE",
				r.activeTableAlias,
				FieldNameEmailVisibility,
			))))
		}
	}

	// wrap in json_extract to ensure that top-level primitives
	// stored as json work correctly when compared to their SQL equivalent
	// (https://github.com/pocketbase/pocketbase/issues/4068)
	if field.Type() == FieldTypeJSON {
		result.NoCoalesce = true
		result.Identifier = dbutils.JSONExtract(r.activeTableAlias+"."+cleanFieldName, "")
		if r.withMultiMatch {
			r.multiMatch.valueIdentifier = dbutils.JSONExtract(r.multiMatchActiveTableAlias+"."+cleanFieldName, "")
		}
	}

	// account for the ":lower" modifier
	if modifier == lowerModifier {
		result.Identifier = "LOWER(" + result.Identifier + ")"
		if r.withMultiMatch {
			r.multiMatch.valueIdentifier = "LOWER(" + r.multiMatch.valueIdentifier + ")"
		}
	}

	return result, nil
}
