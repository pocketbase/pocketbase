package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maps"
	"slices"
	"sort"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

// used as a workaround by some fields for persisting local state between various events
// (for now is kept private and cannot be changed or cloned outside of the core package)
const internalCustomFieldKeyPrefix = "@pbInternal"

var (
	_ Model        = (*Record)(nil)
	_ HookTagger   = (*Record)(nil)
	_ DBExporter   = (*Record)(nil)
	_ FilesManager = (*Record)(nil)
)

type Record struct {
	collection       *Collection
	originalData     map[string]any
	customVisibility *store.Store[string, bool]
	data             *store.Store[string, any]
	expand           *store.Store[string, any]

	BaseModel

	exportCustomData      bool
	ignoreEmailVisibility bool
	ignoreUnchangedFields bool
}

const systemHookIdRecord = "__pbRecordSystemHook__"

func (app *BaseApp) registerRecordHooks() {
	app.OnModelValidate().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordValidate().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelCreate().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordCreate().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelCreateExecute().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordCreateExecute().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterCreateSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordAfterCreateSuccess().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterCreateError().Bind(&hook.Handler[*ModelErrorEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelErrorEvent) error {
			if re, ok := newRecordErrorEventFromModelErrorEvent(me); ok {
				err := me.App.OnRecordAfterCreateError().Trigger(re, func(re *RecordErrorEvent) error {
					syncModelErrorEventWithRecordErrorEvent(me, re)
					defer syncRecordErrorEventWithModelErrorEvent(re, me)
					return me.Next()
				})
				syncModelErrorEventWithRecordErrorEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelUpdate().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordUpdate().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelUpdateExecute().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordUpdateExecute().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterUpdateSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordAfterUpdateSuccess().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterUpdateError().Bind(&hook.Handler[*ModelErrorEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelErrorEvent) error {
			if re, ok := newRecordErrorEventFromModelErrorEvent(me); ok {
				err := me.App.OnRecordAfterUpdateError().Trigger(re, func(re *RecordErrorEvent) error {
					syncModelErrorEventWithRecordErrorEvent(me, re)
					defer syncRecordErrorEventWithModelErrorEvent(re, me)
					return me.Next()
				})
				syncModelErrorEventWithRecordErrorEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelDelete().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordDelete().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelDeleteExecute().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordDeleteExecute().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterDeleteSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelEvent) error {
			if re, ok := newRecordEventFromModelEvent(me); ok {
				err := me.App.OnRecordAfterDeleteSuccess().Trigger(re, func(re *RecordEvent) error {
					syncModelEventWithRecordEvent(me, re)
					defer syncRecordEventWithModelEvent(re, me)
					return me.Next()
				})
				syncModelEventWithRecordEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterDeleteError().Bind(&hook.Handler[*ModelErrorEvent]{
		Id: systemHookIdRecord,
		Func: func(me *ModelErrorEvent) error {
			if re, ok := newRecordErrorEventFromModelErrorEvent(me); ok {
				err := me.App.OnRecordAfterDeleteError().Trigger(re, func(re *RecordErrorEvent) error {
					syncModelErrorEventWithRecordErrorEvent(me, re)
					defer syncRecordErrorEventWithModelErrorEvent(re, me)
					return me.Next()
				})
				syncModelErrorEventWithRecordErrorEvent(me, re)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	// ---------------------------------------------------------------

	app.OnRecordValidate().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionValidate,
				func() error {
					return onRecordValidate(e)
				},
			)
		},
		Priority: 99,
	})

	app.OnRecordCreate().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionCreate,
				e.Next,
			)
		},
		Priority: -99,
	})

	app.OnRecordCreateExecute().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionCreateExecute,
				func() error {
					return onRecordSaveExecute(e)
				},
			)
		},
		Priority: 99,
	})

	app.OnRecordAfterCreateSuccess().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionAfterCreate,
				e.Next,
			)
		},
		Priority: -99,
	})

	app.OnRecordAfterCreateError().Bind(&hook.Handler[*RecordErrorEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordErrorEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionAfterCreateError,
				e.Next,
			)
		},
		Priority: -99,
	})

	app.OnRecordUpdate().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionUpdate,
				e.Next,
			)
		},
		Priority: -99,
	})

	app.OnRecordUpdateExecute().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionUpdateExecute,
				func() error {
					return onRecordSaveExecute(e)
				},
			)
		},
		Priority: 99,
	})

	app.OnRecordAfterUpdateSuccess().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionAfterUpdate,
				e.Next,
			)
		},
		Priority: -99,
	})

	app.OnRecordAfterUpdateError().Bind(&hook.Handler[*RecordErrorEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordErrorEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionAfterUpdateError,
				e.Next,
			)
		},
		Priority: -99,
	})

	app.OnRecordDelete().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionDelete,
				func() error {
					if e.Record.Collection().IsView() {
						return errors.New("view records cannot be deleted")
					}

					return e.Next()
				},
			)
		},
		Priority: -99,
	})

	app.OnRecordDeleteExecute().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionDeleteExecute,
				func() error {
					return onRecordDeleteExecute(e)
				},
			)
		},
		Priority: 99,
	})

	app.OnRecordAfterDeleteSuccess().Bind(&hook.Handler[*RecordEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionAfterDelete,
				e.Next,
			)
		},
		Priority: -99,
	})

	app.OnRecordAfterDeleteError().Bind(&hook.Handler[*RecordErrorEvent]{
		Id: systemHookIdRecord,
		Func: func(e *RecordErrorEvent) error {
			return e.Record.callFieldInterceptors(
				e.Context,
				e.App,
				InterceptorActionAfterDeleteError,
				e.Next,
			)
		},
		Priority: -99,
	})
}

// -------------------------------------------------------------------

// newRecordFromNullStringMap initializes a single new Record model
// with data loaded from the provided NullStringMap.
//
// Note that this method is intended to load and Scan data from a database row result.
func newRecordFromNullStringMap(collection *Collection, data dbx.NullStringMap) (*Record, error) {
	record := NewRecord(collection)

	var fieldName string
	for _, field := range collection.Fields {
		fieldName = field.GetName()

		nullString, ok := data[fieldName]

		var value any
		var err error

		if ok && nullString.Valid {
			value, err = field.PrepareValue(record, nullString.String)
		} else {
			value, err = field.PrepareValue(record, nil)
		}

		if err != nil {
			return nil, err
		}

		// we load only the original data to avoid unnecessary copying the same data into the record.data store
		// (it is also the reason why we don't invoke PostScan on the record itself)
		record.originalData[fieldName] = value

		if fieldName == FieldNameId {
			record.Id = cast.ToString(value)
		}
	}

	record.BaseModel.PostScan()

	return record, nil
}

// newRecordsFromNullStringMaps initializes a new Record model for
// each row in the provided NullStringMap slice.
//
// Note that this method is intended to load and Scan data from a database rows result.
func newRecordsFromNullStringMaps(collection *Collection, rows []dbx.NullStringMap) ([]*Record, error) {
	result := make([]*Record, len(rows))

	var err error
	for i, row := range rows {
		result[i], err = newRecordFromNullStringMap(collection, row)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// -------------------------------------------------------------------

// NewRecord initializes a new empty Record model.
func NewRecord(collection *Collection) *Record {
	record := &Record{
		collection:       collection,
		data:             store.New[string, any](nil),
		customVisibility: store.New[string, bool](nil),
		originalData:     make(map[string]any, len(collection.Fields)),
	}

	// initialize default field values
	var fieldName string
	for _, field := range collection.Fields {
		fieldName = field.GetName()

		if fieldName == FieldNameId {
			continue
		}

		value, _ := field.PrepareValue(record, nil)
		record.originalData[fieldName] = value
	}

	return record
}

// Collection returns the Collection model associated with the current Record model.
//
// NB! The returned collection is only for read purposes and it shouldn't be modified
// because it could have unintended side-effects on other Record models from the same collection.
func (m *Record) Collection() *Collection {
	return m.collection
}

// TableName returns the table name associated with the current Record model.
func (m *Record) TableName() string {
	return m.collection.Name
}

// PostScan implements the [dbx.PostScanner] interface.
//
// It essentially refreshes/updates the current Record original state
// as if the model was fetched from the databases for the first time.
//
// Or in other words, it means that m.Original().FieldsData() will have
// the same values as m.Record().FieldsData().
func (m *Record) PostScan() error {
	if m.Id == "" {
		return errors.New("missing record primary key")
	}

	if err := m.BaseModel.PostScan(); err != nil {
		return err
	}

	m.originalData = m.FieldsData()

	return nil
}

// HookTags returns the hook tags associated with the current record.
func (m *Record) HookTags() []string {
	return []string{m.collection.Name, m.collection.Id}
}

// BaseFilesPath returns the storage dir path used by the record.
func (m *Record) BaseFilesPath() string {
	id := cast.ToString(m.LastSavedPK())
	if id == "" {
		id = m.Id
	}

	return m.collection.BaseFilesPath() + "/" + id
}

// Original returns a shallow copy of the current record model populated
// with its ORIGINAL db data state (aka. right after PostScan())
// and everything else reset to the defaults.
//
// If record was created using NewRecord() the original will be always
// a blank record (until PostScan() is invoked).
func (m *Record) Original() *Record {
	newRecord := NewRecord(m.collection)

	newRecord.originalData = maps.Clone(m.originalData)

	if newRecord.originalData[FieldNameId] != nil {
		newRecord.lastSavedPK = cast.ToString(newRecord.originalData[FieldNameId])
		newRecord.Id = newRecord.lastSavedPK
	}

	return newRecord
}

// Fresh returns a shallow copy of the current record model populated
// with its LATEST data state and everything else reset to the defaults
// (aka. no expand, no unknown fields and with default visibility flags).
func (m *Record) Fresh() *Record {
	newRecord := m.Original()

	// note: this will also load the Id field through m.GetRaw
	var fieldName string
	for _, field := range m.collection.Fields {
		fieldName = field.GetName()
		newRecord.SetRaw(fieldName, m.GetRaw(fieldName))
	}

	return newRecord
}

// Clone returns a shallow copy of the current record model with all of
// its collection and unknown fields data, expand and flags copied.
//
// use [Record.Fresh()] instead if you want a copy with only the latest
// collection fields data and everything else reset to the defaults.
func (m *Record) Clone() *Record {
	newRecord := m.Original()

	newRecord.Id = m.Id
	newRecord.exportCustomData = m.exportCustomData
	newRecord.ignoreEmailVisibility = m.ignoreEmailVisibility
	newRecord.ignoreUnchangedFields = m.ignoreUnchangedFields
	newRecord.customVisibility.Reset(m.customVisibility.GetAll())

	data := m.data.GetAll()
	for k, v := range data {
		newRecord.SetRaw(k, v)
	}

	if m.expand != nil {
		newRecord.SetExpand(m.expand.GetAll())
	}

	return newRecord
}

// Expand returns a shallow copy of the current Record model expand data (if any).
func (m *Record) Expand() map[string]any {
	if m.expand == nil {
		// return a dummy initialized map to avoid assignment to nil map errors
		return map[string]any{}
	}

	return m.expand.GetAll()
}

// SetExpand replaces the current Record's expand with the provided expand arg data (shallow copied).
func (m *Record) SetExpand(expand map[string]any) {
	if m.expand == nil {
		m.expand = store.New[string, any](nil)
	}

	m.expand.Reset(expand)
}

// MergeExpand merges recursively the provided expand data into
// the current model's expand (if any).
//
// Note that if an expanded prop with the same key is a slice (old or new expand)
// then both old and new records will be merged into a new slice (aka. a :merge: [b,c] => [a,b,c]).
// Otherwise the "old" expanded record will be replace with the "new" one (aka. a :merge: aNew => aNew).
func (m *Record) MergeExpand(expand map[string]any) {
	// nothing to merge
	if len(expand) == 0 {
		return
	}

	// no old expand
	if m.expand == nil {
		m.expand = store.New(expand)
		return
	}

	oldExpand := m.expand.GetAll()

	for key, new := range expand {
		old, ok := oldExpand[key]
		if !ok {
			oldExpand[key] = new
			continue
		}

		var wasOldSlice bool
		var oldSlice []*Record
		switch v := old.(type) {
		case *Record:
			oldSlice = []*Record{v}
		case []*Record:
			wasOldSlice = true
			oldSlice = v
		default:
			// invalid old expand data -> assign directly the new
			// (no matter whether new is valid or not)
			oldExpand[key] = new
			continue
		}

		var wasNewSlice bool
		var newSlice []*Record
		switch v := new.(type) {
		case *Record:
			newSlice = []*Record{v}
		case []*Record:
			wasNewSlice = true
			newSlice = v
		default:
			// invalid new expand data -> skip
			continue
		}

		oldIndexed := make(map[string]*Record, len(oldSlice))
		for _, oldRecord := range oldSlice {
			oldIndexed[oldRecord.Id] = oldRecord
		}

		for _, newRecord := range newSlice {
			oldRecord := oldIndexed[newRecord.Id]
			if oldRecord != nil {
				// note: there is no need to update oldSlice since oldRecord is a reference
				oldRecord.MergeExpand(newRecord.Expand())
			} else {
				// missing new entry
				oldSlice = append(oldSlice, newRecord)
			}
		}

		if wasOldSlice || wasNewSlice || len(oldSlice) == 0 {
			oldExpand[key] = oldSlice
		} else {
			oldExpand[key] = oldSlice[0]
		}
	}

	m.expand.Reset(oldExpand)
}

// FieldsData returns a shallow copy ONLY of the collection's fields record's data.
func (m *Record) FieldsData() map[string]any {
	result := make(map[string]any, len(m.collection.Fields))

	var fieldName string
	for _, field := range m.collection.Fields {
		fieldName = field.GetName()
		result[fieldName] = m.Get(fieldName)
	}

	return result
}

// CustomData returns a shallow copy ONLY of the custom record fields data,
// aka. fields that are neither defined by the collection, nor special system ones.
//
// Note that custom fields prefixed with "@pbInternal" are always skipped.
func (m *Record) CustomData() map[string]any {
	if m.data == nil {
		return nil
	}

	fields := m.Collection().Fields

	knownFields := make(map[string]struct{}, len(fields))

	for _, f := range fields {
		knownFields[f.GetName()] = struct{}{}
	}

	result := map[string]any{}

	rawData := m.data.GetAll()
	for k, v := range rawData {
		if _, ok := knownFields[k]; !ok {
			// skip internal custom fields
			if strings.HasPrefix(k, internalCustomFieldKeyPrefix) {
				continue
			}

			result[k] = v
		}
	}

	return result
}

// WithCustomData toggles the export/serialization of custom data fields
// (false by default).
func (m *Record) WithCustomData(state bool) *Record {
	m.exportCustomData = state
	return m
}

// IgnoreEmailVisibility toggles the flag to ignore the auth record email visibility check.
func (m *Record) IgnoreEmailVisibility(state bool) *Record {
	m.ignoreEmailVisibility = state
	return m
}

// IgnoreUnchangedFields toggles the flag to ignore the unchanged fields
// from the DB export for the UPDATE SQL query.
//
// This could be used if you want to save only the record fields that you've changed
// without overwrite other untouched fields in case of concurrent update.
//
// Note that the fields change comparison is based on the current fields against m.Original()
// (aka. if you have performed save on the same Record instance multiple times you may have to refetch it,
// so that m.Original() could reflect the last saved change).
func (m *Record) IgnoreUnchangedFields(state bool) *Record {
	m.ignoreUnchangedFields = state
	return m
}

// Set sets the provided key-value data pair into the current Record
// model directly as it is WITHOUT NORMALIZATIONS.
//
// See also [Record.Set].
func (m *Record) SetRaw(key string, value any) {
	if key == FieldNameId {
		m.Id = cast.ToString(value)
	}

	m.data.Set(key, value)
}

// SetIfFieldExists sets the provided key-value data pair into the current Record model
// ONLY if key is existing Collection field name/modifier.
//
// This method does nothing if key is not a known Collection field name/modifier.
//
// On success returns the matched Field, otherwise - nil.
//
// To set any key-value, including custom/unknown fields, use the [Record.Set] method.
func (m *Record) SetIfFieldExists(key string, value any) Field {
	for _, field := range m.Collection().Fields {
		ff, ok := field.(SetterFinder)
		if ok {
			setter := ff.FindSetter(key)
			if setter != nil {
				setter(m, value)
				return field
			}
		}

		// fallback to the default field PrepareValue method for direct match
		if key == field.GetName() {
			value, _ = field.PrepareValue(m, value)
			m.SetRaw(key, value)
			return field
		}
	}

	return nil
}

// Set sets the provided key-value data pair into the current Record model.
//
// If the record collection has field with name matching the provided "key",
// the value will be further normalized according to the field setter(s).
func (m *Record) Set(key string, value any) {
	switch key {
	case FieldNameExpand: // for backward-compatibility with earlier versions
		m.SetExpand(cast.ToStringMap(value))
	default:
		field := m.SetIfFieldExists(key, value)
		if field == nil {
			// custom key - set it without any transformations
			m.SetRaw(key, value)
		}
	}
}

func (m *Record) GetRaw(key string) any {
	if key == FieldNameId {
		return m.Id
	}

	if v, ok := m.data.GetOk(key); ok {
		return v
	}

	return m.originalData[key]
}

// Get returns a normalized single record model data value for "key".
func (m *Record) Get(key string) any {
	switch key {
	case FieldNameExpand: // for backward-compatibility with earlier versions
		return m.Expand()
	default:
		for _, field := range m.Collection().Fields {
			gm, ok := field.(GetterFinder)
			if !ok {
				continue // no custom getters
			}

			getter := gm.FindGetter(key)
			if getter != nil {
				return getter(m)
			}
		}

		return m.GetRaw(key)
	}
}

// Load bulk loads the provided data into the current Record model.
func (m *Record) Load(data map[string]any) {
	for k, v := range data {
		m.Set(k, v)
	}
}

// GetBool returns the data value for "key" as a bool.
func (m *Record) GetBool(key string) bool {
	return cast.ToBool(m.Get(key))
}

// GetString returns the data value for "key" as a string.
func (m *Record) GetString(key string) string {
	return cast.ToString(m.Get(key))
}

// GetInt returns the data value for "key" as an int.
func (m *Record) GetInt(key string) int {
	return cast.ToInt(m.Get(key))
}

// GetFloat returns the data value for "key" as a float64.
func (m *Record) GetFloat(key string) float64 {
	return cast.ToFloat64(m.Get(key))
}

// GetDateTime returns the data value for "key" as a DateTime instance.
func (m *Record) GetDateTime(key string) types.DateTime {
	d, _ := types.ParseDateTime(m.Get(key))
	return d
}

// GetGeoPoint returns the data value for "key" as a GeoPoint instance.
func (m *Record) GetGeoPoint(key string) types.GeoPoint {
	point := types.GeoPoint{}
	_ = point.Scan(m.Get(key))
	return point
}

// GetStringSlice returns the data value for "key" as a slice of non-zero unique strings.
func (m *Record) GetStringSlice(key string) []string {
	return list.ToUniqueStringSlice(m.Get(key))
}

// GetUnsavedFiles returns the uploaded files for the provided "file" field key,
// (aka. the current [*filesytem.File] values) so that you can apply further
// validations or modifications (including changing the file name or content before persisting).
//
// Example:
//
//	files := record.GetUnsavedFiles("documents")
//	for _, f := range files {
//	    f.Name = "doc_" + f.Name // add a prefix to each file name
//	}
//	app.Save(record) // the files are pointers so the applied changes will transparently reflect on the record value
func (m *Record) GetUnsavedFiles(key string) []*filesystem.File {
	if !strings.HasSuffix(key, ":unsaved") {
		key += ":unsaved"
	}

	values, _ := m.Get(key).([]*filesystem.File)

	return values
}

// Deprecated: replaced with GetUnsavedFiles.
func (m *Record) GetUploadedFiles(key string) []*filesystem.File {
	log.Println("Please replace GetUploadedFiles with GetUnsavedFiles")
	return m.GetUnsavedFiles(key)
}

// Retrieves the "key" json field value and unmarshals it into "result".
//
// Example
//
//	result := struct {
//	    FirstName string `json:"first_name"`
//	}{}
//	err := m.UnmarshalJSONField("my_field_name", &result)
func (m *Record) UnmarshalJSONField(key string, result any) error {
	return json.Unmarshal([]byte(m.GetString(key)), &result)
}

// ExpandedOne retrieves a single relation Record from the already
// loaded expand data of the current model.
//
// If the requested expand relation is multiple, this method returns
// only first available Record from the expanded relation.
//
// Returns nil if there is no such expand relation loaded.
func (m *Record) ExpandedOne(relField string) *Record {
	if m.expand == nil {
		return nil
	}

	rel := m.expand.Get(relField)

	switch v := rel.(type) {
	case *Record:
		return v
	case []*Record:
		if len(v) > 0 {
			return v[0]
		}
	}

	return nil
}

// ExpandedAll retrieves a slice of relation Records from the already
// loaded expand data of the current model.
//
// If the requested expand relation is single, this method normalizes
// the return result and will wrap the single model as a slice.
//
// Returns nil slice if there is no such expand relation loaded.
func (m *Record) ExpandedAll(relField string) []*Record {
	if m.expand == nil {
		return nil
	}

	rel := m.expand.Get(relField)

	switch v := rel.(type) {
	case *Record:
		return []*Record{v}
	case []*Record:
		return v
	}

	return nil
}

// FindFileFieldByFile returns the first file type field for which
// any of the record's data contains the provided filename.
func (m *Record) FindFileFieldByFile(filename string) *FileField {
	for _, field := range m.Collection().Fields {
		if field.Type() != FieldTypeFile {
			continue
		}

		f, ok := field.(*FileField)
		if !ok {
			continue
		}

		filenames := m.GetStringSlice(f.GetName())
		if slices.Contains(filenames, filename) {
			return f
		}
	}

	return nil
}

// DBExport implements the [DBExporter] interface and returns a key-value
// map with the data to be persisted when saving the Record in the database.
func (m *Record) DBExport(app App) (map[string]any, error) {
	result, err := m.dbExport()
	if err != nil {
		return nil, err
	}

	// remove exported fields that haven't changed
	// (with exception of the id column)
	if !m.IsNew() && m.ignoreUnchangedFields {
		oldResult, err := m.Original().dbExport()
		if err != nil {
			return nil, err
		}

		for oldK, oldV := range oldResult {
			if oldK == idColumn {
				continue
			}
			newV, ok := result[oldK]
			if ok && areValuesEqual(newV, oldV) {
				delete(result, oldK)
			}
		}
	}

	return result, nil
}

func (m *Record) dbExport() (map[string]any, error) {
	fields := m.Collection().Fields

	result := make(map[string]any, len(fields))

	var fieldName string
	for _, field := range fields {
		fieldName = field.GetName()

		if f, ok := field.(DriverValuer); ok {
			v, err := f.DriverValue(m)
			if err != nil {
				return nil, err
			}
			result[fieldName] = v
		} else {
			result[fieldName] = m.GetRaw(fieldName)
		}
	}

	return result, nil
}

func areValuesEqual(a any, b any) bool {
	switch av := a.(type) {
	case string:
		bv, ok := b.(string)
		return ok && bv == av
	case bool:
		bv, ok := b.(bool)
		return ok && bv == av
	case float32:
		bv, ok := b.(float32)
		return ok && bv == av
	case float64:
		bv, ok := b.(float64)
		return ok && bv == av
	case uint:
		bv, ok := b.(uint)
		return ok && bv == av
	case uint8:
		bv, ok := b.(uint8)
		return ok && bv == av
	case uint16:
		bv, ok := b.(uint16)
		return ok && bv == av
	case uint32:
		bv, ok := b.(uint32)
		return ok && bv == av
	case uint64:
		bv, ok := b.(uint64)
		return ok && bv == av
	case int:
		bv, ok := b.(int)
		return ok && bv == av
	case int8:
		bv, ok := b.(int8)
		return ok && bv == av
	case int16:
		bv, ok := b.(int16)
		return ok && bv == av
	case int32:
		bv, ok := b.(int32)
		return ok && bv == av
	case int64:
		bv, ok := b.(int64)
		return ok && bv == av
	case []byte:
		bv, ok := b.([]byte)
		return ok && bytes.Equal(av, bv)
	case []string:
		bv, ok := b.([]string)
		return ok && slices.Equal(av, bv)
	case []int:
		bv, ok := b.([]int)
		return ok && slices.Equal(av, bv)
	case []int32:
		bv, ok := b.([]int32)
		return ok && slices.Equal(av, bv)
	case []int64:
		bv, ok := b.([]int64)
		return ok && slices.Equal(av, bv)
	case []float32:
		bv, ok := b.([]float32)
		return ok && slices.Equal(av, bv)
	case []float64:
		bv, ok := b.([]float64)
		return ok && slices.Equal(av, bv)
	case types.JSONArray[string]:
		bv, ok := b.(types.JSONArray[string])
		return ok && slices.Equal(av, bv)
	case types.JSONRaw:
		bv, ok := b.(types.JSONRaw)
		return ok && bytes.Equal(av, bv)
	default:
		aRaw, err := json.Marshal(a)
		if err != nil {
			return false
		}

		bRaw, err := json.Marshal(b)
		if err != nil {
			return false
		}

		return bytes.Equal(aRaw, bRaw)
	}
}

// Hide hides the specified fields from the public safe serialization of the record.
func (record *Record) Hide(fieldNames ...string) *Record {
	for _, name := range fieldNames {
		record.customVisibility.Set(name, false)
	}

	return record
}

// Unhide forces to unhide the specified fields from the public safe serialization
// of the record (even when the collection field itself is marked as hidden).
func (record *Record) Unhide(fieldNames ...string) *Record {
	for _, name := range fieldNames {
		record.customVisibility.Set(name, true)
	}

	return record
}

// PublicExport exports only the record fields that are safe to be public.
//
// To export unknown data fields you need to set record.WithCustomData(true).
//
// For auth records, to force the export of the email field you need to set
// record.IgnoreEmailVisibility(true).
func (record *Record) PublicExport() map[string]any {
	export := make(map[string]any, len(record.collection.Fields)+3)

	var isVisible, hasCustomVisibility bool

	customVisibility := record.customVisibility.GetAll()

	// export schema fields
	var fieldName string
	for _, f := range record.collection.Fields {
		fieldName = f.GetName()

		isVisible, hasCustomVisibility = customVisibility[fieldName]
		if !hasCustomVisibility {
			isVisible = !f.GetHidden()
		}

		if !isVisible {
			continue
		}

		export[fieldName] = record.Get(fieldName)
	}

	// export custom fields
	if record.exportCustomData {
		for k, v := range record.CustomData() {
			isVisible, hasCustomVisibility = customVisibility[k]
			if !hasCustomVisibility || isVisible {
				export[k] = v
			}
		}
	}

	if record.Collection().IsAuth() {
		// always hide the password and tokenKey fields
		delete(export, FieldNamePassword)
		delete(export, FieldNameTokenKey)

		if !record.ignoreEmailVisibility && !record.GetBool(FieldNameEmailVisibility) {
			delete(export, FieldNameEmail)
		}
	}

	// add helper collection reference fields
	isVisible, hasCustomVisibility = customVisibility[FieldNameCollectionId]
	if !hasCustomVisibility || isVisible {
		export[FieldNameCollectionId] = record.collection.Id
	}
	isVisible, hasCustomVisibility = customVisibility[FieldNameCollectionName]
	if !hasCustomVisibility || isVisible {
		export[FieldNameCollectionName] = record.collection.Name
	}

	// add expand (if non-nil)
	isVisible, hasCustomVisibility = customVisibility[FieldNameExpand]
	if (!hasCustomVisibility || isVisible) && record.expand != nil {
		export[FieldNameExpand] = record.expand.GetAll()
	}

	return export
}

// MarshalJSON implements the [json.Marshaler] interface.
//
// Only the data exported by `PublicExport()` will be serialized.
func (m Record) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.PublicExport())
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (m *Record) UnmarshalJSON(data []byte) error {
	result := map[string]any{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	m.Load(result)

	return nil
}

// ReplaceModifiers returns a new map with applied modifier
// values based on the current record and the specified data.
//
// The resolved modifier keys will be removed.
//
// Multiple modifiers will be applied one after another,
// while reusing the previous base key value result (ex. 1; -5; +2 => -2).
//
// Note that because Go doesn't guaranteed the iteration order of maps,
// we would explicitly apply shorter keys first for a more consistent and reproducible behavior.
//
// Example usage:
//
//	 newData := record.ReplaceModifiers(data)
//		// record: {"field": 10}
//		// data:   {"field+": 5}
//		// result: {"field": 15}
func (m *Record) ReplaceModifiers(data map[string]any) map[string]any {
	if len(data) == 0 {
		return data
	}

	dataCopy := maps.Clone(data)

	recordCopy := m.Fresh()

	// key orders is not guaranteed so
	sortedDataKeys := make([]string, 0, len(data))
	for k := range data {
		sortedDataKeys = append(sortedDataKeys, k)
	}
	sort.SliceStable(sortedDataKeys, func(i int, j int) bool {
		return len(sortedDataKeys[i]) < len(sortedDataKeys[j])
	})

	for _, k := range sortedDataKeys {
		field := recordCopy.SetIfFieldExists(k, data[k])
		if field != nil {
			// delete the original key in case it is with a modifer (ex. "items+")
			delete(dataCopy, k)

			// store the transformed value under the field name
			dataCopy[field.GetName()] = recordCopy.Get(field.GetName())
		}
	}

	return dataCopy
}

// -------------------------------------------------------------------

func (m *Record) callFieldInterceptors(
	ctx context.Context,
	app App,
	actionName string,
	actionFunc func() error,
) error {
	// the firing order of the fields doesn't matter
	for _, field := range m.Collection().Fields {
		if f, ok := field.(RecordInterceptor); ok {
			oldfn := actionFunc
			actionFunc = func() error {
				return f.Intercept(ctx, app, m, actionName, oldfn)
			}
		}
	}

	return actionFunc()
}

func onRecordValidate(e *RecordEvent) error {
	errs := validation.Errors{}

	for _, f := range e.Record.Collection().Fields {
		if err := f.ValidateValue(e.Context, e.App, e.Record); err != nil {
			errs[f.GetName()] = err
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return e.Next()
}

func onRecordSaveExecute(e *RecordEvent) error {
	if e.Record.Collection().IsAuth() {
		// ensure that the token key is regenerated on password change or email change
		if !e.Record.IsNew() {
			lastSavedRecord, err := e.App.FindRecordById(e.Record.Collection(), e.Record.Id)
			if err != nil {
				return err
			}

			if lastSavedRecord.TokenKey() == e.Record.TokenKey() &&
				(lastSavedRecord.Get(FieldNamePassword) != e.Record.Get(FieldNamePassword) ||
					lastSavedRecord.Email() != e.Record.Email()) {
				e.Record.RefreshTokenKey()
			}
		}

		// cross-check that the auth record id is unique across all auth collections.
		authCollections, err := e.App.FindAllCollections(CollectionTypeAuth)
		if err != nil {
			return fmt.Errorf("unable to fetch the auth collections for cross-id unique check: %w", err)
		}
		for _, collection := range authCollections {
			if e.Record.Collection().Id == collection.Id {
				continue // skip current collection (sqlite will do the check for us)
			}
			record, _ := e.App.FindRecordById(collection, e.Record.Id)
			if record != nil {
				return validation.Errors{
					FieldNameId: validation.NewError("validation_invalid_auth_id", "Invalid or duplicated auth record id."),
				}
			}
		}
	}

	err := e.Next()
	if err == nil {
		return nil
	}

	return validators.NormalizeUniqueIndexError(
		err,
		e.Record.Collection().Name,
		e.Record.Collection().Fields.FieldNames(),
	)
}

func onRecordDeleteExecute(e *RecordEvent) error {
	// fetch rel references (if any)
	//
	// note: the select is outside of the transaction to minimize
	// SQLITE_BUSY errors when mixing read&write in a single transaction
	refs, err := e.App.FindCachedCollectionReferences(e.Record.Collection())
	if err != nil {
		return err
	}

	originalApp := e.App
	txErr := e.App.RunInTransaction(func(txApp App) error {
		e.App = txApp

		// delete the record before the relation references to ensure that there
		// will be no "A<->B" relations to prevent deadlock when calling DeleteRecord recursively
		if err := e.Next(); err != nil {
			return err
		}

		return cascadeRecordDelete(txApp, e.Record, refs)
	})
	e.App = originalApp

	return txErr
}

// cascadeRecordDelete triggers cascade deletion for the provided references.
//
// NB! This method is expected to be called from inside of a transaction.
func cascadeRecordDelete(app App, mainRecord *Record, refs map[*Collection][]Field) error {
	// Sort the refs keys to ensure that the cascade events firing order is always the same.
	// This is not necessary for the operation to function correctly but it helps having deterministic output during testing.
	sortedRefKeys := make([]*Collection, 0, len(refs))
	for k := range refs {
		sortedRefKeys = append(sortedRefKeys, k)
	}
	sort.Slice(sortedRefKeys, func(i, j int) bool {
		return sortedRefKeys[i].Name < sortedRefKeys[j].Name
	})

	for _, refCollection := range sortedRefKeys {
		fields, ok := refs[refCollection]

		if !ok || refCollection.IsView() {
			continue // skip missing or view collections
		}

		recordTableName := inflector.Columnify(refCollection.Name)

		for _, field := range fields {
			prefixedFieldName := recordTableName + "." + inflector.Columnify(field.GetName())

			query := app.RecordQuery(refCollection)

			if opt, ok := field.(MultiValuer); !ok || !opt.IsMultiple() {
				query.AndWhere(dbx.HashExp{prefixedFieldName: mainRecord.Id})
			} else {
				query.AndWhere(dbx.Exists(dbx.NewExp(fmt.Sprintf(
					`SELECT 1 FROM json_each(CASE WHEN json_valid([[%s]]) THEN [[%s]] ELSE json_array([[%s]]) END) {{__je__}} WHERE [[__je__.value]]={:jevalue}`,
					prefixedFieldName, prefixedFieldName, prefixedFieldName,
				), dbx.Params{
					"jevalue": mainRecord.Id,
				})))
			}

			if refCollection.Id == mainRecord.Collection().Id {
				query.AndWhere(dbx.Not(dbx.HashExp{recordTableName + ".id": mainRecord.Id}))
			}

			// trigger cascade for each batchSize rel items until there is none
			batchSize := 4000
			rows := make([]*Record, 0, batchSize)
			for {
				if err := query.Limit(int64(batchSize)).All(&rows); err != nil {
					return err
				}

				total := len(rows)
				if total == 0 {
					break
				}

				err := deleteRefRecords(app, mainRecord, rows, field)
				if err != nil {
					return err
				}

				if total < batchSize {
					break // no more items
				}

				rows = rows[:0] // keep allocated memory
			}
		}
	}

	return nil
}

// deleteRefRecords checks if related records has to be deleted (if `CascadeDelete` is set)
// OR
// just unset the record id from any relation field values (if they are not required).
//
// NB! This method is expected to be called from inside of a transaction.
func deleteRefRecords(app App, mainRecord *Record, refRecords []*Record, field Field) error {
	relField, _ := field.(*RelationField)
	if relField == nil {
		return errors.New("only RelationField is supported at the moment, got " + field.Type())
	}

	for _, refRecord := range refRecords {
		ids := refRecord.GetStringSlice(relField.Name)

		// unset the record id
		for i := len(ids) - 1; i >= 0; i-- {
			if ids[i] == mainRecord.Id {
				ids = append(ids[:i], ids[i+1:]...)
				break
			}
		}

		// cascade delete the reference
		// (only if there are no other active references in case of multiple select)
		if relField.CascadeDelete && len(ids) == 0 {
			if err := app.Delete(refRecord); err != nil {
				return err
			}
			// no further actions are needed (the reference is deleted)
			continue
		}

		if relField.Required && len(ids) == 0 {
			return fmt.Errorf("the record cannot be deleted because it is part of a required reference in record %s (%s collection)", refRecord.Id, refRecord.Collection().Name)
		}

		// save the reference changes
		// (without validation because it is possible that another relation field to have a reference to a previous deleted record)
		refRecord.Set(relField.Name, ids)
		if err := app.SaveNoValidate(refRecord); err != nil {
			return err
		}
	}

	return nil
}
