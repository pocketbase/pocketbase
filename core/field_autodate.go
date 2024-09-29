package core

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	Fields[FieldTypeAutodate] = func() Field {
		return &AutodateField{}
	}
}

const FieldTypeAutodate = "autodate"

var (
	_ Field             = (*AutodateField)(nil)
	_ SetterFinder      = (*AutodateField)(nil)
	_ RecordInterceptor = (*AutodateField)(nil)
)

// AutodateField defines an "autodate" type field, aka.
// field which datetime value could be auto set on record create/update.
//
// Requires either both or at least one of the OnCreate or OnUpdate options to be set.
type AutodateField struct {
	Id          string `form:"id" json:"id"`
	Name        string `form:"name" json:"name"`
	System      bool   `form:"system" json:"system"`
	Hidden      bool   `form:"hidden" json:"hidden"`
	Presentable bool   `form:"presentable" json:"presentable"`

	// ---

	// OnCreate auto sets the current datetime as field value on record create.
	OnCreate bool `form:"onCreate" json:"onCreate"`

	// OnUpdate auto sets the current datetime as field value on record update.
	OnUpdate bool `form:"onUpdate" json:"onUpdate"`
}

// Type implements [Field.Type] interface method.
func (f *AutodateField) Type() string {
	return FieldTypeAutodate
}

// GetId implements [Field.GetId] interface method.
func (f *AutodateField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *AutodateField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *AutodateField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *AutodateField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *AutodateField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *AutodateField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *AutodateField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *AutodateField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *AutodateField) ColumnType(app App) string {
	return "TEXT DEFAULT '' NOT NULL" // note: sqlite doesn't allow adding new columns with non-constant defaults
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *AutodateField) PrepareValue(record *Record, raw any) (any, error) {
	val, _ := types.ParseDateTime(raw)
	return val, nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *AutodateField) ValidateValue(ctx context.Context, app App, record *Record) error {
	return nil
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *AutodateField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	oldOnCreate := f.OnCreate
	oldOnUpdate := f.OnUpdate

	oldCollection, _ := app.FindCollectionByNameOrId(collection.Id)
	if oldCollection != nil {
		oldField, ok := oldCollection.Fields.GetById(f.Id).(*AutodateField)
		if ok && oldField != nil {
			oldOnCreate = oldField.OnCreate
			oldOnUpdate = oldField.OnUpdate
		}
	}

	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
		validation.Field(
			&f.OnCreate,
			validation.When(f.System, validation.By(validators.Equal(oldOnCreate))),
			validation.Required.Error("either onCreate or onUpdate must be enabled").When(!f.OnUpdate),
		),
		validation.Field(
			&f.OnUpdate,
			validation.When(f.System, validation.By(validators.Equal(oldOnUpdate))),
			validation.Required.Error("either onCreate or onUpdate must be enabled").When(!f.OnCreate),
		),
	)
}

// FindSetter implements the [SetterFinder] interface.
func (f *AutodateField) FindSetter(key string) SetterFunc {
	switch key {
	case f.Name:
		// return noopSetter to disallow updating the value with record.Set()
		return noopSetter
	default:
		return nil
	}
}

// Intercept implements the [RecordInterceptor] interface.
func (f *AutodateField) Intercept(
	ctx context.Context,
	app App,
	record *Record,
	actionName string,
	actionFunc func() error,
) error {
	switch actionName {
	case InterceptorActionCreate:
		// ignore for custom date manually set with record.SetRaw()
		if f.OnCreate && !f.hasBeenManuallyChanged(record) {
			record.SetRaw(f.Name, types.NowDateTime())
		}
	case InterceptorActionUpdate:
		// ignore for custom date manually set with record.SetRaw()
		if f.OnUpdate && !f.hasBeenManuallyChanged(record) {
			record.SetRaw(f.Name, types.NowDateTime())
		}
	}

	return actionFunc()
}

func (f *AutodateField) hasBeenManuallyChanged(record *Record) bool {
	vNew, _ := record.GetRaw(f.Name).(types.DateTime)
	vOld, _ := record.Original().GetRaw(f.Name).(types.DateTime)

	return vNew.String() != vOld.String()
}
