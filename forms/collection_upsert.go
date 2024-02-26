package forms

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/types"
)

var collectionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]*$`)

// CollectionUpsert is a [models.Collection] upsert (create/update) form.
type CollectionUpsert struct {
	app        core.App
	dao        *daos.Dao
	collection *models.Collection

	Id         string                  `form:"id" json:"id"`
	Type       string                  `form:"type" json:"type"`
	Name       string                  `form:"name" json:"name"`
	System     bool                    `form:"system" json:"system"`
	Schema     schema.Schema           `form:"schema" json:"schema"`
	Indexes    types.JsonArray[string] `form:"indexes" json:"indexes"`
	ListRule   *string                 `form:"listRule" json:"listRule"`
	ViewRule   *string                 `form:"viewRule" json:"viewRule"`
	CreateRule *string                 `form:"createRule" json:"createRule"`
	UpdateRule *string                 `form:"updateRule" json:"updateRule"`
	DeleteRule *string                 `form:"deleteRule" json:"deleteRule"`
	Options    types.JsonMap           `form:"options" json:"options"`
}

// NewCollectionUpsert creates a new [CollectionUpsert] form with initializer
// config created from the provided [core.App] and [models.Collection] instances
// (for create you could pass a pointer to an empty Collection - `&models.Collection{}`).
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewCollectionUpsert(app core.App, collection *models.Collection) *CollectionUpsert {
	form := &CollectionUpsert{
		app:        app,
		dao:        app.Dao(),
		collection: collection,
	}

	// load defaults
	form.Id = form.collection.Id
	form.Type = form.collection.Type
	form.Name = form.collection.Name
	form.System = form.collection.System
	form.Indexes = form.collection.Indexes
	form.ListRule = form.collection.ListRule
	form.ViewRule = form.collection.ViewRule
	form.CreateRule = form.collection.CreateRule
	form.UpdateRule = form.collection.UpdateRule
	form.DeleteRule = form.collection.DeleteRule
	form.Options = form.collection.Options

	if form.Type == "" {
		form.Type = models.CollectionTypeBase
	}

	clone, _ := form.collection.Schema.Clone()
	if clone != nil && form.Type != models.CollectionTypeView {
		form.Schema = *clone
	} else {
		form.Schema = schema.Schema{}
	}

	return form
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *CollectionUpsert) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *CollectionUpsert) Validate() error {
	isAuth := form.Type == models.CollectionTypeAuth
	isView := form.Type == models.CollectionTypeView

	// generate schema from the query (overwriting any explicit user defined schema)
	if isView {
		options := models.CollectionViewOptions{}
		if err := decodeOptions(form.Options, &options); err != nil {
			return err
		}
		form.Schema, _ = form.dao.CreateViewSchema(options.Query)
	}

	return validation.ValidateStruct(form,
		validation.Field(
			&form.Id,
			validation.When(
				form.collection.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
				validation.By(validators.UniqueId(form.dao, form.collection.TableName())),
			).Else(validation.In(form.collection.Id)),
		),
		validation.Field(
			&form.System,
			validation.By(form.ensureNoSystemFlagChange),
		),
		validation.Field(
			&form.Type,
			validation.Required,
			validation.In(
				models.CollectionTypeBase,
				models.CollectionTypeAuth,
				models.CollectionTypeView,
			),
			validation.By(form.ensureNoTypeChange),
		),
		validation.Field(
			&form.Name,
			validation.Required,
			validation.Length(1, 255),
			validation.Match(collectionNameRegex),
			validation.By(form.ensureNoSystemNameChange),
			validation.By(form.checkUniqueName),
			validation.By(form.checkForVia),
		),
		// validates using the type's own validation rules + some collection's specifics
		validation.Field(
			&form.Schema,
			validation.By(form.checkMinSchemaFields),
			validation.By(form.ensureNoSystemFieldsChange),
			validation.By(form.ensureNoFieldsTypeChange),
			validation.By(form.checkRelationFields),
			validation.When(isAuth, validation.By(form.ensureNoAuthFieldName)),
		),
		validation.Field(&form.ListRule, validation.By(form.checkRule)),
		validation.Field(&form.ViewRule, validation.By(form.checkRule)),
		validation.Field(
			&form.CreateRule,
			validation.When(isView, validation.Nil),
			validation.By(form.checkRule),
		),
		validation.Field(
			&form.UpdateRule,
			validation.When(isView, validation.Nil),
			validation.By(form.checkRule),
		),
		validation.Field(
			&form.DeleteRule,
			validation.When(isView, validation.Nil),
			validation.By(form.checkRule),
		),
		validation.Field(&form.Indexes, validation.By(form.checkIndexes)),
		validation.Field(&form.Options, validation.By(form.checkOptions)),
	)
}

func (form *CollectionUpsert) checkForVia(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil
	}

	if strings.Contains(strings.ToLower(v), "_via_") {
		return validation.NewError("validation_invalid_name", "The name of the collection cannot contain '_via_'.")
	}

	return nil
}

func (form *CollectionUpsert) checkUniqueName(value any) error {
	v, _ := value.(string)

	// ensure unique collection name
	if !form.dao.IsCollectionNameUnique(v, form.collection.Id) {
		return validation.NewError("validation_collection_name_exists", "Collection name must be unique (case insensitive).")
	}

	// ensure that the collection name doesn't collide with the id of any collection
	if form.dao.FindById(&models.Collection{}, v) == nil {
		return validation.NewError("validation_collection_name_id_duplicate", "The name must not match an existing collection id.")
	}

	return nil
}

func (form *CollectionUpsert) ensureNoSystemNameChange(value any) error {
	v, _ := value.(string)

	if !form.collection.IsNew() && form.collection.System && v != form.collection.Name {
		return validation.NewError("validation_collection_system_name_change", "System collections cannot be renamed.")
	}

	return nil
}

func (form *CollectionUpsert) ensureNoSystemFlagChange(value any) error {
	v, _ := value.(bool)

	if !form.collection.IsNew() && v != form.collection.System {
		return validation.NewError("validation_collection_system_flag_change", "System collection state cannot be changed.")
	}

	return nil
}

func (form *CollectionUpsert) ensureNoTypeChange(value any) error {
	v, _ := value.(string)

	if !form.collection.IsNew() && v != form.collection.Type {
		return validation.NewError("validation_collection_type_change", "Collection type cannot be changed.")
	}

	return nil
}

func (form *CollectionUpsert) ensureNoFieldsTypeChange(value any) error {
	v, _ := value.(schema.Schema)

	for i, field := range v.Fields() {
		oldField := form.collection.Schema.GetFieldById(field.Id)

		if oldField != nil && oldField.Type != field.Type {
			return validation.Errors{fmt.Sprint(i): validation.NewError(
				"validation_field_type_change",
				"Field type cannot be changed.",
			)}
		}
	}

	return nil
}

func (form *CollectionUpsert) checkRelationFields(value any) error {
	v, _ := value.(schema.Schema)

	for i, field := range v.Fields() {
		if field.Type != schema.FieldTypeRelation {
			continue
		}

		options, _ := field.Options.(*schema.RelationOptions)
		if options == nil {
			return validation.Errors{fmt.Sprint(i): validation.Errors{
				"options": validation.NewError(
					"validation_schema_invalid_relation_field_options",
					"The relation field has invalid field options.",
				)},
			}
		}

		// prevent collectionId change
		oldField := form.collection.Schema.GetFieldById(field.Id)
		if oldField != nil {
			oldOptions, _ := oldField.Options.(*schema.RelationOptions)
			if oldOptions != nil && oldOptions.CollectionId != options.CollectionId {
				return validation.Errors{fmt.Sprint(i): validation.Errors{
					"options": validation.Errors{
						"collectionId": validation.NewError(
							"validation_field_relation_change",
							"The relation collection cannot be changed.",
						),
					}},
				}
			}
		}

		relCollection, _ := form.dao.FindCollectionByNameOrId(options.CollectionId)

		// validate collectionId
		if relCollection == nil || relCollection.Id != options.CollectionId {
			return validation.Errors{fmt.Sprint(i): validation.Errors{
				"options": validation.Errors{
					"collectionId": validation.NewError(
						"validation_field_invalid_relation",
						"The relation collection doesn't exist.",
					),
				}},
			}
		}

		// allow only views to have relations to other views
		// (see https://github.com/pocketbase/pocketbase/issues/3000)
		if form.Type != models.CollectionTypeView && relCollection.IsView() {
			return validation.Errors{fmt.Sprint(i): validation.Errors{
				"options": validation.Errors{
					"collectionId": validation.NewError(
						"validation_field_non_view_base_relation_collection",
						"Non view collections are not allowed to have a view relation.",
					),
				}},
			}
		}
	}

	return nil
}

func (form *CollectionUpsert) ensureNoAuthFieldName(value any) error {
	v, _ := value.(schema.Schema)

	if form.Type != models.CollectionTypeAuth {
		return nil // not an auth collection
	}

	authFieldNames := schema.AuthFieldNames()
	// exclude the meta RecordUpsert form fields
	authFieldNames = append(authFieldNames, "password", "passwordConfirm", "oldPassword")

	errs := validation.Errors{}
	for i, field := range v.Fields() {
		if list.ExistInSlice(field.Name, authFieldNames) {
			errs[fmt.Sprint(i)] = validation.Errors{
				"name": validation.NewError(
					"validation_reserved_auth_field_name",
					"The field name is reserved and cannot be used.",
				),
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (form *CollectionUpsert) checkMinSchemaFields(value any) error {
	v, _ := value.(schema.Schema)

	switch form.Type {
	case models.CollectionTypeAuth, models.CollectionTypeView:
		return nil // no schema fields constraint
	default:
		if len(v.Fields()) == 0 {
			return validation.ErrRequired
		}
	}

	return nil
}

func (form *CollectionUpsert) ensureNoSystemFieldsChange(value any) error {
	v, _ := value.(schema.Schema)

	for _, oldField := range form.collection.Schema.Fields() {
		if !oldField.System {
			continue
		}

		newField := v.GetFieldById(oldField.Id)

		if newField == nil || oldField.String() != newField.String() {
			return validation.NewError("validation_system_field_change", "System fields cannot be deleted or changed.")
		}
	}

	return nil
}

func (form *CollectionUpsert) checkRule(value any) error {
	v, _ := value.(*string)
	if v == nil || *v == "" {
		return nil // nothing to check
	}

	dummy := *form.collection
	dummy.Type = form.Type
	dummy.Schema = form.Schema
	dummy.System = form.System
	dummy.Options = form.Options

	r := resolvers.NewRecordFieldResolver(form.dao, &dummy, nil, true)

	_, err := search.FilterData(*v).BuildExpr(r)
	if err != nil {
		return validation.NewError("validation_invalid_rule", "Invalid filter rule. Raw error: "+err.Error())
	}

	return nil
}

func (form *CollectionUpsert) checkIndexes(value any) error {
	v, _ := value.(types.JsonArray[string])

	if form.Type == models.CollectionTypeView && len(v) > 0 {
		return validation.NewError(
			"validation_indexes_not_supported",
			"The collection doesn't support indexes.",
		)
	}

	for i, rawIndex := range v {
		parsed := dbutils.ParseIndex(rawIndex)

		if !parsed.IsValid() {
			return validation.Errors{
				strconv.Itoa(i): validation.NewError(
					"validation_invalid_index_expression",
					"Invalid CREATE INDEX expression.",
				),
			}
		}

		// note: we don't check the index table because it is always
		// overwritten by the daos.SyncRecordTableSchema to allow
		// easier partial modifications (eg. changing only the collection name).
		// if !strings.EqualFold(parsed.TableName, form.Name) {
		// 	return validation.Errors{
		// 		strconv.Itoa(i): validation.NewError(
		// 			"validation_invalid_index_table",
		// 			fmt.Sprintf("The index table must be the same as the collection name."),
		// 		),
		// 	}
		// }
	}

	return nil
}

func (form *CollectionUpsert) checkOptions(value any) error {
	v, _ := value.(types.JsonMap)

	switch form.Type {
	case models.CollectionTypeAuth:
		options := models.CollectionAuthOptions{}
		if err := decodeOptions(v, &options); err != nil {
			return err
		}

		// check the generic validations
		if err := options.Validate(); err != nil {
			return err
		}

		// additional form specific validations
		if err := form.checkRule(options.ManageRule); err != nil {
			return validation.Errors{"manageRule": err}
		}
	case models.CollectionTypeView:
		options := models.CollectionViewOptions{}
		if err := decodeOptions(v, &options); err != nil {
			return err
		}

		// check the generic validations
		if err := options.Validate(); err != nil {
			return err
		}

		// check the query option
		if _, err := form.dao.CreateViewSchema(options.Query); err != nil {
			return validation.Errors{
				"query": validation.NewError(
					"validation_invalid_view_query",
					fmt.Sprintf("Invalid query - %s", err.Error()),
				),
			}
		}
	}

	return nil
}

func decodeOptions(options types.JsonMap, result any) error {
	raw, err := options.MarshalJSON()
	if err != nil {
		return validation.NewError("validation_invalid_options", "Invalid options.")
	}

	if err := json.Unmarshal(raw, result); err != nil {
		return validation.NewError("validation_invalid_options", "Invalid options.")
	}

	return nil
}

// Submit validates the form and upserts the form's Collection model.
//
// On success the related record table schema will be auto updated.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *CollectionUpsert) Submit(interceptors ...InterceptorFunc[*models.Collection]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	if form.collection.IsNew() {
		// type can be set only on create
		form.collection.Type = form.Type

		// system flag can be set only on create
		form.collection.System = form.System

		// custom insertion id can be set only on create
		if form.Id != "" {
			form.collection.MarkAsNew()
			form.collection.SetId(form.Id)
		}
	}

	// system collections cannot be renamed
	if form.collection.IsNew() || !form.collection.System {
		form.collection.Name = form.Name
	}

	// view schema is autogenerated on save and cannot have indexes
	if !form.collection.IsView() {
		form.collection.Schema = form.Schema

		// normalize indexes format
		form.collection.Indexes = make(types.JsonArray[string], len(form.Indexes))
		for i, rawIdx := range form.Indexes {
			form.collection.Indexes[i] = dbutils.ParseIndex(rawIdx).Build()
		}
	}

	form.collection.ListRule = form.ListRule
	form.collection.ViewRule = form.ViewRule
	form.collection.CreateRule = form.CreateRule
	form.collection.UpdateRule = form.UpdateRule
	form.collection.DeleteRule = form.DeleteRule
	form.collection.SetOptions(form.Options)

	return runInterceptors(form.collection, func(collection *models.Collection) error {
		return form.dao.SaveCollection(collection)
	}, interceptors...)
}
