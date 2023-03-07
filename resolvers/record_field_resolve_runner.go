package resolvers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
)

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
		return nil, fmt.Errorf("the runner was already used")
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
		if r.resolver.requestData == nil {
			return &search.ResolverResult{Identifier: "NULL"}, nil
		}

		if strings.HasPrefix(r.fieldName, "@request.auth.") {
			return r.processRequestAuthField()
		}

		if strings.HasPrefix(r.fieldName, "@request.data.") && len(r.activeProps) > 2 {
			name, modifier, err := splitModifier(r.activeProps[2])
			if err != nil {
				return nil, err
			}

			dataField := r.resolver.baseCollection.Schema.GetFieldByName(name)
			if dataField == nil {
				return r.resolver.resolveStaticRequestField(r.activeProps[1:]...)
			}

			dataField.InitOptions()

			// check for data relation field
			if dataField.Type == schema.FieldTypeRelation && len(r.activeProps) > 3 {
				return r.processRequestDataRelationField(dataField)
			}

			// check for select:each field
			if modifier == eachModifier && dataField.Type == schema.FieldTypeSelect && len(r.activeProps) == 3 {
				return r.processRequestDataSelectEachModifier(dataField)
			}

			// check for data arrayble fields ":length" modifier
			if modifier == lengthModifier && list.ExistInSlice(dataField.Type, schema.ArraybleFieldTypes()) && len(r.activeProps) == 3 {
				return r.processRequestDataLengthModifier(dataField)
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

	collection, err := r.resolver.loadCollection(r.activeProps[1])
	if err != nil {
		return nil, fmt.Errorf("failed to load collection %q from field path %q", r.activeProps[1], r.fieldName)
	}

	r.activeCollectionName = collection.Name
	r.activeTableAlias = inflector.Columnify("__collection_" + r.activeCollectionName)

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
	// plain auth field
	// ---
	if list.ExistInSlice(r.fieldName, plainRequestAuthFields) {
		return r.resolver.resolveStaticRequestField(r.activeProps[1:]...)
	}

	// resolve the auth collection field
	// ---
	if r.resolver.requestData == nil || r.resolver.requestData.AuthRecord == nil || r.resolver.requestData.AuthRecord.Collection() == nil {
		return &search.ResolverResult{Identifier: "NULL"}, nil
	}

	collection := r.resolver.requestData.AuthRecord.Collection()
	r.resolver.loadedCollections = append(r.resolver.loadedCollections, collection)

	r.activeCollectionName = collection.Name
	r.activeTableAlias = "__auth_" + inflector.Columnify(r.activeCollectionName)

	// join the auth collection to the main query
	r.resolver.registerJoin(
		inflector.Columnify(r.activeCollectionName),
		r.activeTableAlias,
		dbx.HashExp{
			// aka. __auth_users.id = :userId
			(r.activeTableAlias + ".id"): r.resolver.requestData.AuthRecord.Id,
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
				(r.multiMatchActiveTableAlias + ".id"): r.resolver.requestData.AuthRecord.Id,
			},
		},
	)

	// leave only the auth relation fields
	// aka. @request.auth.fieldA.fieldB -> fieldA.fieldB
	r.activeProps = r.activeProps[2:]

	return r.processActiveProps()
}

func (r *runner) processRequestDataLengthModifier(dataField *schema.SchemaField) (*search.ResolverResult, error) {
	dataItems := list.ToUniqueStringSlice(r.resolver.requestData.Data[dataField.Name])

	result := &search.ResolverResult{
		Identifier: fmt.Sprintf("%d", len(dataItems)),
	}

	return result, nil
}

func (r *runner) processRequestDataSelectEachModifier(dataField *schema.SchemaField) (*search.ResolverResult, error) {
	options, ok := dataField.Options.(*schema.SelectOptions)
	if !ok {
		return nil, fmt.Errorf("failed to initialize field %q options", dataField.Name)
	}

	dataItems := list.ToUniqueStringSlice(r.resolver.requestData.Data[dataField.Name])
	rawJson, err := json.Marshal(dataItems)
	if err != nil {
		return nil, fmt.Errorf("cannot marshalize the data select item for field %q", r.activeProps[2])
	}

	placeholder := "dataSelect" + security.PseudorandomString(4)
	cleanFieldName := inflector.Columnify(dataField.Name)
	jeTable := fmt.Sprintf("json_each({:%s})", placeholder)
	jeAlias := "__dataSelect_" + cleanFieldName + "_je"
	r.resolver.registerJoin(jeTable, jeAlias, nil)

	result := &search.ResolverResult{
		Identifier: fmt.Sprintf("[[%s.value]]", jeAlias),
		Params:     dbx.Params{placeholder: rawJson},
	}

	if options.MaxSelect != 1 {
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
		r.multiMatch.params[placeholder2] = rawJson
		r.multiMatch.valueIdentifier = fmt.Sprintf("[[%s.value]]", jeAlias2)

		result.MultiMatchSubQuery = r.multiMatch
	}

	return result, nil
}

func (r *runner) processRequestDataRelationField(dataField *schema.SchemaField) (*search.ResolverResult, error) {
	options, ok := dataField.Options.(*schema.RelationOptions)
	if !ok {
		return nil, fmt.Errorf("failed to initialize data field %q options", dataField.Name)
	}

	dataRelCollection, err := r.resolver.loadCollection(options.CollectionId)
	if err != nil {
		return nil, fmt.Errorf("failed to load collection %q from data field %q", options.CollectionId, dataField.Name)
	}

	var dataRelIds []string
	if r.resolver.requestData != nil && len(r.resolver.requestData.Data) != 0 {
		dataRelIds = list.ToUniqueStringSlice(r.resolver.requestData.Data[dataField.Name])
	}
	if len(dataRelIds) == 0 {
		return &search.ResolverResult{Identifier: "NULL"}, nil
	}

	r.activeCollectionName = dataRelCollection.Name
	r.activeTableAlias = inflector.Columnify("__data_" + dataRelCollection.Name)

	// join the data rel collection to the main collection
	r.resolver.registerJoin(
		inflector.Columnify(r.activeCollectionName),
		r.activeTableAlias,
		dbx.In(
			fmt.Sprintf("[[%s.id]]", inflector.Columnify(r.activeTableAlias)),
			list.ToInterfaceSlice(dataRelIds)...,
		),
	)

	if options.MaxSelect == nil || *options.MaxSelect != 1 {
		r.withMultiMatch = true
	}

	// join the data rel collection to the multi-match subquery
	r.multiMatchActiveTableAlias = inflector.Columnify("__data_mm_" + dataRelCollection.Name)
	r.multiMatch.joins = append(
		r.multiMatch.joins,
		&join{
			tableName:  inflector.Columnify(r.activeCollectionName),
			tableAlias: r.multiMatchActiveTableAlias,
			on:         dbx.In(r.multiMatchActiveTableAlias+".id", list.ToInterfaceSlice(dataRelIds)...),
		},
	)

	// leave only the data relation fields
	// aka. @request.data.someRel.fieldA.fieldB -> fieldA.fieldB
	r.activeProps = r.activeProps[3:]

	return r.processActiveProps()
}

func (r *runner) processActiveProps() (*search.ResolverResult, error) {
	totalProps := len(r.activeProps)

	for i, prop := range r.activeProps {
		collection, err := r.resolver.loadCollection(r.activeCollectionName)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve field %q", prop)
		}

		// last prop
		if i == totalProps-1 {
			// system field, aka. internal model prop
			// (always available but not part of the collection schema)
			// -------------------------------------------------------
			if list.ExistInSlice(prop, resolvableSystemFieldNames(collection)) {
				result := &search.ResolverResult{
					Identifier: fmt.Sprintf("[[%s.%s]]", r.activeTableAlias, inflector.Columnify(prop)),
				}

				// allow querying only auth records with emails marked as public
				if prop == schema.FieldNameEmail && !r.allowHiddenFields {
					result.AfterBuild = func(expr dbx.Expression) dbx.Expression {
						return dbx.Enclose(dbx.And(expr, dbx.NewExp(fmt.Sprintf(
							"[[%s.%s]] = TRUE",
							r.activeTableAlias,
							schema.FieldNameEmailVisibility,
						))))
					}
				}

				if r.withMultiMatch {
					r.multiMatch.valueIdentifier = fmt.Sprintf("[[%s.%s]]", r.multiMatchActiveTableAlias, inflector.Columnify(prop))
					result.MultiMatchSubQuery = r.multiMatch
				}

				return result, nil
			}

			name, modifier, err := splitModifier(prop)
			if err != nil {
				return nil, err
			}

			field := collection.Schema.GetFieldByName(name)
			if field == nil {
				if r.nullifyMisingField {
					return &search.ResolverResult{Identifier: "NULL"}, nil
				}
				return nil, fmt.Errorf("unknown field %q", name)
			}

			cleanFieldName := inflector.Columnify(field.Name)

			// arrayble fields ":length" modifier
			// -------------------------------------------------------
			if modifier == lengthModifier && list.ExistInSlice(field.Type, schema.ArraybleFieldTypes()) {
				jePair := r.activeTableAlias + "." + cleanFieldName

				result := &search.ResolverResult{
					Identifier: jsonArrayLength(jePair),
				}

				if r.withMultiMatch {
					jePair2 := r.multiMatchActiveTableAlias + "." + cleanFieldName
					r.multiMatch.valueIdentifier = jsonArrayLength(jePair2)
					result.MultiMatchSubQuery = r.multiMatch
				}

				return result, nil
			}

			// select field with ":each" modifier
			// -------------------------------------------------------
			if field.Type == schema.FieldTypeSelect && modifier == eachModifier {
				jePair := r.activeTableAlias + "." + cleanFieldName
				jeAlias := r.activeTableAlias + "_" + cleanFieldName + "_je"
				r.resolver.registerJoin(jsonEach(jePair), jeAlias, nil)

				result := &search.ResolverResult{
					Identifier: fmt.Sprintf("[[%s.value]]", jeAlias),
				}

				field.InitOptions()
				options, ok := field.Options.(*schema.SelectOptions)
				if !ok {
					return nil, fmt.Errorf("failed to initialize field %q options", prop)
				}

				if options.MaxSelect != 1 {
					r.withMultiMatch = true
				}

				if r.withMultiMatch {
					jePair2 := r.multiMatchActiveTableAlias + "." + cleanFieldName
					jeAlias2 := r.multiMatchActiveTableAlias + "_" + cleanFieldName + "_je"

					r.multiMatch.joins = append(r.multiMatch.joins, &join{
						tableName:  jsonEach(jePair2),
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
				Identifier: fmt.Sprintf("[[%s.%s]]", r.activeTableAlias, cleanFieldName),
			}

			if r.withMultiMatch {
				r.multiMatch.valueIdentifier = fmt.Sprintf("[[%s.%s]]", r.multiMatchActiveTableAlias, cleanFieldName)
				result.MultiMatchSubQuery = r.multiMatch
			}

			return result, nil
		}

		field := collection.Schema.GetFieldByName(prop)
		if field == nil {
			if r.nullifyMisingField {
				return &search.ResolverResult{Identifier: "NULL"}, nil
			}
			return nil, fmt.Errorf("unknown field %q", prop)
		}

		// check if it is a json field
		if field.Type == schema.FieldTypeJson {
			var jsonPath strings.Builder
			jsonPath.WriteString("$")
			for _, p := range r.activeProps[i+1:] {
				if _, err := strconv.Atoi(p); err == nil {
					jsonPath.WriteString("[")
					jsonPath.WriteString(inflector.Columnify(p))
					jsonPath.WriteString("]")
				} else {
					jsonPath.WriteString(".")
					jsonPath.WriteString(inflector.Columnify(p))
				}
			}

			result := &search.ResolverResult{
				Identifier: fmt.Sprintf(
					"JSON_EXTRACT([[%s.%s]], '%s')",
					r.activeTableAlias,
					inflector.Columnify(prop),
					jsonPath.String(),
				),
			}

			if r.withMultiMatch {
				r.multiMatch.valueIdentifier = fmt.Sprintf(
					"JSON_EXTRACT([[%s.%s]], '%s')",
					r.multiMatchActiveTableAlias,
					inflector.Columnify(prop),
					jsonPath.String(),
				)
				result.MultiMatchSubQuery = r.multiMatch
			}

			return result, nil
		}

		// check if it is a relation field
		if field.Type != schema.FieldTypeRelation {
			return nil, fmt.Errorf("field %q is not a valid relation", prop)
		}

		// join the relation to the main query
		// ---
		field.InitOptions()
		options, ok := field.Options.(*schema.RelationOptions)
		if !ok {
			return nil, fmt.Errorf("failed to initialize field %q options", prop)
		}

		relCollection, relErr := r.resolver.loadCollection(options.CollectionId)
		if relErr != nil {
			return nil, fmt.Errorf("failed to find field %q collection", prop)
		}

		cleanFieldName := inflector.Columnify(field.Name)
		prefixedFieldName := r.activeTableAlias + "." + cleanFieldName
		newTableAlias := r.activeTableAlias + "_" + cleanFieldName
		newCollectionName := relCollection.Name

		if !options.IsMultiple() {
			r.resolver.registerJoin(
				inflector.Columnify(newCollectionName),
				newTableAlias,
				dbx.NewExp(fmt.Sprintf("[[%s.id]] = [[%s]]", newTableAlias, prefixedFieldName)),
			)
		} else {
			jeAlias := r.activeTableAlias + "_" + cleanFieldName + "_je"
			r.resolver.registerJoin(jsonEach(prefixedFieldName), jeAlias, nil)
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
		if options.IsMultiple() {
			r.withMultiMatch = true
		}

		newTableAlias2 := r.multiMatchActiveTableAlias + "_" + cleanFieldName
		prefixedFieldName2 := r.multiMatchActiveTableAlias + "." + cleanFieldName

		if !options.IsMultiple() {
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
					tableName:  jsonEach(prefixedFieldName2),
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

func jsonArrayLength(tableColumnPair string) string {
	return fmt.Sprintf(
		// note: the case is used to normalize value access for single and multiple relations.
		`json_array_length(CASE WHEN json_valid([[%s]]) THEN [[%s]] ELSE json_array([[%s]]) END)`,
		tableColumnPair, tableColumnPair, tableColumnPair,
	)
}

func jsonEach(tableColumnPair string) string {
	return fmt.Sprintf(
		// note: the case is used to normalize value access for single and multiple relations.
		`json_each(CASE WHEN json_valid([[%s]]) THEN [[%s]] ELSE json_array([[%s]]) END)`,
		tableColumnPair, tableColumnPair, tableColumnPair,
	)
}

func resolvableSystemFieldNames(collection *models.Collection) []string {
	result := schema.BaseModelFieldNames()

	if collection.IsAuth() {
		result = append(
			result,
			schema.FieldNameUsername,
			schema.FieldNameVerified,
			schema.FieldNameEmailVisibility,
			schema.FieldNameEmail,
		)
	}

	return result
}
