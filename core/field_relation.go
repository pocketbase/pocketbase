package core

import (
	"context"
	"database/sql/driver"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	Fields[FieldTypeRelation] = func() Field {
		return &RelationField{}
	}
}

const FieldTypeRelation = "relation"

var (
	_ Field        = (*RelationField)(nil)
	_ MultiValuer  = (*RelationField)(nil)
	_ DriverValuer = (*RelationField)(nil)
	_ SetterFinder = (*RelationField)(nil)
)

// RelationField defines "relation" type field for storing single or
// multiple collection record references.
//
// Requires the CollectionId option to be set.
//
// If MaxSelect is not set or <= 1, then the field value is expected to be a single record id.
//
// If MaxSelect is > 1, then the field value is expected to be a slice of record ids.
//
// The respective zero record field value is either empty string (single) or empty string slice (multiple).
//
// ---
//
// The following additional setter keys are available:
//
//   - "fieldName+" - append one or more values to the existing record one. For example:
//
//     record.Set("categories+", []string{"new1", "new2"}) // []string{"old1", "old2", "new1", "new2"}
//
//   - "+fieldName" - prepend one or more values to the existing record one. For example:
//
//     record.Set("+categories", []string{"new1", "new2"}) // []string{"new1", "new2", "old1", "old2"}
//
//   - "fieldName-" - subtract one or more values from the existing record one. For example:
//
//     record.Set("categories-", "old1") // []string{"old2"}
type RelationField struct {
	// Name (required) is the unique name of the field.
	Name string `form:"name" json:"name"`

	// Id is the unique stable field identifier.
	//
	// It is automatically generated from the name when adding to a collection FieldsList.
	Id string `form:"id" json:"id"`

	// System prevents the renaming and removal of the field.
	System bool `form:"system" json:"system"`

	// Hidden hides the field from the API response.
	Hidden bool `form:"hidden" json:"hidden"`

	// Presentable hints the Dashboard UI to use the underlying
	// field record value in the relation preview label.
	Presentable bool `form:"presentable" json:"presentable"`

	// ---

	// CollectionId is the id of the related collection.
	CollectionId string `form:"collectionId" json:"collectionId"`

	// CascadeDelete indicates whether the root model should be deleted
	// in case of delete of all linked relations.
	CascadeDelete bool `form:"cascadeDelete" json:"cascadeDelete"`

	// MinSelect indicates the min number of allowed relation records
	// that could be linked to the main model.
	//
	// No min limit is applied if it is zero or negative value.
	MinSelect int `form:"minSelect" json:"minSelect"`

	// MaxSelect indicates the max number of allowed relation records
	// that could be linked to the main model.
	//
	// For multiple select the value must be > 1, otherwise fallbacks to single (default).
	//
	// If MinSelect is set, MaxSelect must be at least >= MinSelect.
	MaxSelect int `form:"maxSelect" json:"maxSelect"`

	// Required will require the field value to be non-empty.
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *RelationField) Type() string {
	return FieldTypeRelation
}

// GetId implements [Field.GetId] interface method.
func (f *RelationField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *RelationField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *RelationField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *RelationField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *RelationField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *RelationField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *RelationField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *RelationField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// IsMultiple implements [MultiValuer] interface and checks whether the
// current field options support multiple values.
func (f *RelationField) IsMultiple() bool {
	return f.MaxSelect > 1
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *RelationField) ColumnType(app App) string {
	if f.IsMultiple() {
		return "JSON DEFAULT '[]' NOT NULL"
	}

	return "TEXT DEFAULT '' NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *RelationField) PrepareValue(record *Record, raw any) (any, error) {
	return f.normalizeValue(raw), nil
}

func (f *RelationField) normalizeValue(raw any) any {
	val := list.ToUniqueStringSlice(raw)

	if !f.IsMultiple() {
		if len(val) > 0 {
			return val[len(val)-1] // the last selected
		}
		return ""
	}

	return val
}

// DriverValue implements the [DriverValuer] interface.
func (f *RelationField) DriverValue(record *Record) (driver.Value, error) {
	val := list.ToUniqueStringSlice(record.GetRaw(f.Name))

	if !f.IsMultiple() {
		if len(val) > 0 {
			return val[len(val)-1], nil // the last selected
		}
		return "", nil
	}

	// serialize as json string array
	return append(types.JSONArray[string]{}, val...), nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *RelationField) ValidateValue(ctx context.Context, app App, record *Record) error {
	ids := list.ToUniqueStringSlice(record.GetRaw(f.Name))
	if len(ids) == 0 {
		if f.Required {
			return validation.ErrRequired
		}
		return nil // nothing to check
	}

	if f.MinSelect > 0 && len(ids) < f.MinSelect {
		return validation.NewError("validation_not_enough_values", "Select at least {{.minSelect}}").
			SetParams(map[string]any{"minSelect": f.MinSelect})
	}

	maxSelect := max(f.MaxSelect, 1)
	if len(ids) > maxSelect {
		return validation.NewError("validation_too_many_values", "Select no more than {{.maxSelect}}").
			SetParams(map[string]any{"maxSelect": maxSelect})
	}

	// check if the related records exist
	// ---
	relCollection, err := app.FindCachedCollectionByNameOrId(f.CollectionId)
	if err != nil {
		return validation.NewError("validation_missing_rel_collection", "Relation connection is missing or cannot be accessed")
	}

	var total int
	_ = app.DB().
		Select("count(*)").
		From(relCollection.Name).
		AndWhere(dbx.In("id", list.ToInterfaceSlice(ids)...)).
		Row(&total)
	if total != len(ids) {
		return validation.NewError("validation_missing_rel_records", "Failed to find all relation records with the provided ids")
	}
	// ---

	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *RelationField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
		validation.Field(&f.CollectionId, validation.Required, validation.By(f.checkCollectionId(app, collection))),
		validation.Field(&f.MinSelect, validation.Min(0)),
		validation.Field(&f.MaxSelect, validation.When(f.MinSelect > 0, validation.Required), validation.Min(f.MinSelect)),
	)
}

func (f *RelationField) checkCollectionId(app App, collection *Collection) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(string)
		if v == "" {
			return nil // nothing to check
		}

		var oldCollection *Collection

		if !collection.IsNew() {
			var err error
			oldCollection, err = app.FindCachedCollectionByNameOrId(collection.Id)
			if err != nil {
				return err
			}
		}

		// prevent collectionId change
		if oldCollection != nil {
			oldField, _ := oldCollection.Fields.GetById(f.Id).(*RelationField)
			if oldField != nil && oldField.CollectionId != v {
				return validation.NewError(
					"validation_field_relation_change",
					"The relation collection cannot be changed.",
				)
			}
		}

		relCollection, _ := app.FindCachedCollectionByNameOrId(v)

		// validate collectionId
		if relCollection == nil || relCollection.Id != v {
			return validation.NewError(
				"validation_field_relation_missing_collection",
				"The relation collection doesn't exist.",
			)
		}

		// allow only views to have relations to other views
		// (see https://github.com/pocketbase/pocketbase/issues/3000)
		if !collection.IsView() && relCollection.IsView() {
			return validation.NewError(
				"validation_relation_field_non_view_base_collection",
				"Only view collections are allowed to have relations to other views.",
			)
		}

		return nil
	}
}

// ---

// FindSetter implements [SetterFinder] interface method.
func (f *RelationField) FindSetter(key string) SetterFunc {
	switch key {
	case f.Name:
		return f.setValue
	case "+" + f.Name:
		return f.prependValue
	case f.Name + "+":
		return f.appendValue
	case f.Name + "-":
		return f.subtractValue
	default:
		return nil
	}
}

func (f *RelationField) setValue(record *Record, raw any) {
	record.SetRaw(f.Name, f.normalizeValue(raw))
}

func (f *RelationField) appendValue(record *Record, modifierValue any) {
	val := record.GetRaw(f.Name)

	val = append(
		list.ToUniqueStringSlice(val),
		list.ToUniqueStringSlice(modifierValue)...,
	)

	f.setValue(record, val)
}

func (f *RelationField) prependValue(record *Record, modifierValue any) {
	val := record.GetRaw(f.Name)

	val = append(
		list.ToUniqueStringSlice(modifierValue),
		list.ToUniqueStringSlice(val)...,
	)

	f.setValue(record, val)
}

func (f *RelationField) subtractValue(record *Record, modifierValue any) {
	val := record.GetRaw(f.Name)

	val = list.SubtractSlice(
		list.ToUniqueStringSlice(val),
		list.ToUniqueStringSlice(modifierValue),
	)

	f.setValue(record, val)
}
