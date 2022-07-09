package forms

import (
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/search"
)

var collectionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_]*$`)

// CollectionUpsert defines a collection upsert (create/update) form.
type CollectionUpsert struct {
	app        core.App
	collection *models.Collection
	isCreate   bool

	Name       string        `form:"name" json:"name"`
	System     bool          `form:"system" json:"system"`
	Schema     schema.Schema `form:"schema" json:"schema"`
	ListRule   *string       `form:"listRule" json:"listRule"`
	ViewRule   *string       `form:"viewRule" json:"viewRule"`
	CreateRule *string       `form:"createRule" json:"createRule"`
	UpdateRule *string       `form:"updateRule" json:"updateRule"`
	DeleteRule *string       `form:"deleteRule" json:"deleteRule"`
}

// NewCollectionUpsert creates new collection upsert form for the provided Collection model
// (pass an empty Collection model instance (`&models.Collection{}`) for create).
func NewCollectionUpsert(app core.App, collection *models.Collection) *CollectionUpsert {
	form := &CollectionUpsert{
		app:        app,
		collection: collection,
		isCreate:   !collection.HasId(),
	}

	// load defaults
	form.Name = collection.Name
	form.System = collection.System
	form.ListRule = collection.ListRule
	form.ViewRule = collection.ViewRule
	form.CreateRule = collection.CreateRule
	form.UpdateRule = collection.UpdateRule
	form.DeleteRule = collection.DeleteRule

	clone, _ := collection.Schema.Clone()
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
			validation.By(form.ensureNoFieldsNameReuse),
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

	if !form.app.Dao().IsCollectionNameUnique(v, form.collection.Id) {
		return validation.NewError("validation_collection_name_exists", "Collection name must be unique (case insensitive).")
	}

	if (form.isCreate || !strings.EqualFold(v, form.collection.Name)) && form.app.Dao().HasTable(v) {
		return validation.NewError("validation_collection_name_table_exists", "The collection name must be also unique table name.")
	}

	return nil
}

func (form *CollectionUpsert) ensureNoSystemNameChange(value any) error {
	v, _ := value.(string)

	if form.isCreate || !form.collection.System || v == form.collection.Name {
		return nil
	}

	return validation.NewError("validation_system_collection_name_change", "System collections cannot be renamed.")
}

func (form *CollectionUpsert) ensureNoSystemFlagChange(value any) error {
	v, _ := value.(bool)

	if form.isCreate || v == form.collection.System {
		return nil
	}

	return validation.NewError("validation_system_collection_flag_change", "System collection state cannot be changed.")
}

func (form *CollectionUpsert) ensureNoFieldsTypeChange(value any) error {
	v, _ := value.(schema.Schema)

	for _, field := range v.Fields() {
		oldField := form.collection.Schema.GetFieldById(field.Id)

		if oldField != nil && oldField.Type != field.Type {
			return validation.NewError("validation_field_type_change", "Field type cannot be changed.")
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

func (form *CollectionUpsert) ensureNoFieldsNameReuse(value any) error {
	v, _ := value.(schema.Schema)

	for _, field := range v.Fields() {
		oldField := form.collection.Schema.GetFieldByName(field.Name)

		if oldField != nil && oldField.Id != field.Id {
			return validation.NewError("validation_field_old_field_exist", "Cannot use existing schema field names when renaming fields.")
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
	r := resolvers.NewRecordFieldResolver(form.app.Dao(), dummy, nil)

	_, err := search.FilterData(*v).BuildExpr(r)
	if err != nil {
		return validation.NewError("validation_collection_rule", "Invalid filter rule.")
	}

	return nil
}

// Submit validates the form and upserts the form's Collection model.
//
// On success the related record table schema will be auto updated.
func (form *CollectionUpsert) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	// system flag can be set only for create
	if form.isCreate {
		form.collection.System = form.System
	}

	// system collections cannot be renamed
	if form.isCreate || !form.collection.System {
		form.collection.Name = form.Name
	}

	form.collection.Schema = form.Schema
	form.collection.ListRule = form.ListRule
	form.collection.ViewRule = form.ViewRule
	form.collection.CreateRule = form.CreateRule
	form.collection.UpdateRule = form.UpdateRule
	form.collection.DeleteRule = form.DeleteRule

	return form.app.Dao().SaveCollection(form.collection)
}
