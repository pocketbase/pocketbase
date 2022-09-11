package forms

import (
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/search"
)

var collectionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]*$`)

// CollectionUpsert specifies a [models.Collection] upsert (create/update) form.
type CollectionUpsert struct {
	config     CollectionUpsertConfig
	collection *models.Collection

	Id         string        `form:"id" json:"id"`
	Name       string        `form:"name" json:"name"`
	System     bool          `form:"system" json:"system"`
	Schema     schema.Schema `form:"schema" json:"schema"`
	ListRule   *string       `form:"listRule" json:"listRule"`
	ViewRule   *string       `form:"viewRule" json:"viewRule"`
	CreateRule *string       `form:"createRule" json:"createRule"`
	UpdateRule *string       `form:"updateRule" json:"updateRule"`
	DeleteRule *string       `form:"deleteRule" json:"deleteRule"`
}

// CollectionUpsertConfig is the [CollectionUpsert] factory initializer config.
//
// NB! App is a required struct member.
type CollectionUpsertConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewCollectionUpsert creates a new [CollectionUpsert] form with initializer
// config created from the provided [core.App] and [models.Collection] instances
// (for create you could pass a pointer to an empty Collection - `&models.Collection{}`).
//
// If you want to submit the form as part of another transaction, use
// [NewCollectionUpsertWithConfig] with explicitly set Dao.
func NewCollectionUpsert(app core.App, collection *models.Collection) *CollectionUpsert {
	return NewCollectionUpsertWithConfig(CollectionUpsertConfig{
		App: app,
	}, collection)
}

// NewCollectionUpsertWithConfig creates a new [CollectionUpsert] form
// with the provided config and [models.Collection] instance or panics on invalid configuration
// (for create you could pass a pointer to an empty Collection - `&models.Collection{}`).
func NewCollectionUpsertWithConfig(config CollectionUpsertConfig, collection *models.Collection) *CollectionUpsert {
	form := &CollectionUpsert{
		config:     config,
		collection: collection,
	}

	if form.config.App == nil || form.collection == nil {
		panic("Invalid initializer config or nil upsert model.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	// load defaults
	form.Id = form.collection.Id
	form.Name = form.collection.Name
	form.System = form.collection.System
	form.ListRule = form.collection.ListRule
	form.ViewRule = form.collection.ViewRule
	form.CreateRule = form.collection.CreateRule
	form.UpdateRule = form.collection.UpdateRule
	form.DeleteRule = form.collection.DeleteRule

	clone, _ := form.collection.Schema.Clone()
	if clone != nil {
		form.Schema = *clone
	} else {
		form.Schema = schema.Schema{}
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *CollectionUpsert) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Id,
			validation.When(
				form.collection.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
			).Else(validation.In(form.collection.Id)),
		),
		validation.Field(
			&form.System,
			validation.By(form.ensureNoSystemFlagChange),
		),
		validation.Field(
			&form.Name,
			validation.Required,
			validation.Length(1, 255),
			validation.Match(collectionNameRegex),
			validation.By(form.ensureNoSystemNameChange),
			validation.By(form.checkUniqueName),
		),
		// validates using the type's own validation rules + some collection's specific
		validation.Field(
			&form.Schema,
			validation.By(form.ensureNoSystemFieldsChange),
			validation.By(form.ensureNoFieldsTypeChange),
			validation.By(form.ensureExistingRelationCollectionId),
		),
		validation.Field(&form.ListRule, validation.By(form.checkRule)),
		validation.Field(&form.ViewRule, validation.By(form.checkRule)),
		validation.Field(&form.CreateRule, validation.By(form.checkRule)),
		validation.Field(&form.UpdateRule, validation.By(form.checkRule)),
		validation.Field(&form.DeleteRule, validation.By(form.checkRule)),
	)
}

func (form *CollectionUpsert) checkUniqueName(value any) error {
	v, _ := value.(string)

	if !form.config.Dao.IsCollectionNameUnique(v, form.collection.Id) {
		return validation.NewError("validation_collection_name_exists", "Collection name must be unique (case insensitive).")
	}

	if (form.collection.IsNew() || !strings.EqualFold(v, form.collection.Name)) && form.config.Dao.HasTable(v) {
		return validation.NewError("validation_collection_name_table_exists", "The collection name must be also unique table name.")
	}

	return nil
}

func (form *CollectionUpsert) ensureNoSystemNameChange(value any) error {
	v, _ := value.(string)

	if form.collection.IsNew() || !form.collection.System || v == form.collection.Name {
		return nil
	}

	return validation.NewError("validation_system_collection_name_change", "System collections cannot be renamed.")
}

func (form *CollectionUpsert) ensureNoSystemFlagChange(value any) error {
	v, _ := value.(bool)

	if form.collection.IsNew() || v == form.collection.System {
		return nil
	}

	return validation.NewError("validation_system_collection_flag_change", "System collection state cannot be changed.")
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

func (form *CollectionUpsert) ensureExistingRelationCollectionId(value any) error {
	v, _ := value.(schema.Schema)

	for i, field := range v.Fields() {
		if field.Type != schema.FieldTypeRelation {
			continue
		}

		options, _ := field.Options.(*schema.RelationOptions)
		if options == nil {
			continue
		}

		if _, err := form.config.Dao.FindCollectionByNameOrId(options.CollectionId); err != nil {
			return validation.Errors{fmt.Sprint(i): validation.NewError(
				"validation_field_invalid_relation",
				"The relation collection doesn't exist.",
			)}
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

	dummy := &models.Collection{Schema: form.Schema}
	r := resolvers.NewRecordFieldResolver(form.config.Dao, dummy, nil)

	_, err := search.FilterData(*v).BuildExpr(r)
	if err != nil {
		return validation.NewError("validation_collection_rule", "Invalid filter rule.")
	}

	return nil
}

// Submit validates the form and upserts the form's Collection model.
//
// On success the related record table schema will be auto updated.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *CollectionUpsert) Submit(interceptors ...InterceptorFunc) error {
	if err := form.Validate(); err != nil {
		return err
	}

	if form.collection.IsNew() {
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

	form.collection.Schema = form.Schema
	form.collection.ListRule = form.ListRule
	form.collection.ViewRule = form.ViewRule
	form.collection.CreateRule = form.CreateRule
	form.collection.UpdateRule = form.UpdateRule
	form.collection.DeleteRule = form.DeleteRule

	return runInterceptors(func() error {
		return form.config.Dao.SaveCollection(form.collection)
	}, interceptors...)
}
