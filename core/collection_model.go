package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

var (
	_ Model        = (*Collection)(nil)
	_ DBExporter   = (*Collection)(nil)
	_ FilesManager = (*Collection)(nil)
)

const (
	CollectionTypeBase = "base"
	CollectionTypeAuth = "auth"
	CollectionTypeView = "view"
)

const systemHookIdCollection = "__pbCollectionSystemHook__"

const defaultLowercaseRecordIdPattern = "^[a-z0-9]+$"

func (app *BaseApp) registerCollectionHooks() {
	app.OnModelValidate().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionValidate().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelCreate().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionCreate().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelCreateExecute().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionCreateExecute().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterCreateSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionAfterCreateSuccess().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterCreateError().Bind(&hook.Handler[*ModelErrorEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelErrorEvent) error {
			if ce, ok := newCollectionErrorEventFromModelErrorEvent(me); ok {
				err := me.App.OnCollectionAfterCreateError().Trigger(ce, func(ce *CollectionErrorEvent) error {
					syncModelErrorEventWithCollectionErrorEvent(me, ce)
					defer syncCollectionErrorEventWithModelErrorEvent(ce, me)
					return me.Next()
				})
				syncModelErrorEventWithCollectionErrorEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelUpdate().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionUpdate().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelUpdateExecute().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionUpdateExecute().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterUpdateSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionAfterUpdateSuccess().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterUpdateError().Bind(&hook.Handler[*ModelErrorEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelErrorEvent) error {
			if ce, ok := newCollectionErrorEventFromModelErrorEvent(me); ok {
				err := me.App.OnCollectionAfterUpdateError().Trigger(ce, func(ce *CollectionErrorEvent) error {
					syncModelErrorEventWithCollectionErrorEvent(me, ce)
					defer syncCollectionErrorEventWithModelErrorEvent(ce, me)
					return me.Next()
				})
				syncModelErrorEventWithCollectionErrorEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelDelete().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionDelete().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelDeleteExecute().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionDeleteExecute().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterDeleteSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelEvent) error {
			if ce, ok := newCollectionEventFromModelEvent(me); ok {
				err := me.App.OnCollectionAfterDeleteSuccess().Trigger(ce, func(ce *CollectionEvent) error {
					syncModelEventWithCollectionEvent(me, ce)
					defer syncCollectionEventWithModelEvent(ce, me)
					return me.Next()
				})
				syncModelEventWithCollectionEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterDeleteError().Bind(&hook.Handler[*ModelErrorEvent]{
		Id: systemHookIdCollection,
		Func: func(me *ModelErrorEvent) error {
			if ce, ok := newCollectionErrorEventFromModelErrorEvent(me); ok {
				err := me.App.OnCollectionAfterDeleteError().Trigger(ce, func(ce *CollectionErrorEvent) error {
					syncModelErrorEventWithCollectionErrorEvent(me, ce)
					defer syncCollectionErrorEventWithModelErrorEvent(ce, me)
					return me.Next()
				})
				syncModelErrorEventWithCollectionErrorEvent(me, ce)
				return err
			}

			return me.Next()
		},
		Priority: -99,
	})

	//  --------------------------------------------------------------

	app.OnCollectionValidate().Bind(&hook.Handler[*CollectionEvent]{
		Id:       systemHookIdCollection,
		Func:     onCollectionValidate,
		Priority: 99,
	})

	app.OnCollectionCreate().Bind(&hook.Handler[*CollectionEvent]{
		Id:       systemHookIdCollection,
		Func:     onCollectionSave,
		Priority: -99,
	})

	app.OnCollectionUpdate().Bind(&hook.Handler[*CollectionEvent]{
		Id:       systemHookIdCollection,
		Func:     onCollectionSave,
		Priority: -99,
	})

	app.OnCollectionCreateExecute().Bind(&hook.Handler[*CollectionEvent]{
		Id:   systemHookIdCollection,
		Func: onCollectionSaveExecute,
		// execute as latest as possible, aka. closer to the db action to minimize the transactions lock time
		Priority: 99,
	})

	app.OnCollectionUpdateExecute().Bind(&hook.Handler[*CollectionEvent]{
		Id:       systemHookIdCollection,
		Func:     onCollectionSaveExecute,
		Priority: 99, // execute as latest as possible, aka. closer to the db action to minimize the transactions lock time
	})

	app.OnCollectionDeleteExecute().Bind(&hook.Handler[*CollectionEvent]{
		Id:       systemHookIdCollection,
		Func:     onCollectionDeleteExecute,
		Priority: 99, // execute as latest as possible, aka. closer to the db action to minimize the transactions lock time
	})

	// reload cache on failure
	// ---
	onErrorReloadCachedCollections := func(ce *CollectionErrorEvent) error {
		if err := ce.App.ReloadCachedCollections(); err != nil {
			ce.App.Logger().Warn("Failed to reload collections cache after collection change error", "error", err)
		}

		return ce.Next()
	}
	app.OnCollectionAfterCreateError().Bind(&hook.Handler[*CollectionErrorEvent]{
		Id:       systemHookIdCollection,
		Func:     onErrorReloadCachedCollections,
		Priority: -99,
	})
	app.OnCollectionAfterUpdateError().Bind(&hook.Handler[*CollectionErrorEvent]{
		Id:       systemHookIdCollection,
		Func:     onErrorReloadCachedCollections,
		Priority: -99,
	})
	app.OnCollectionAfterDeleteError().Bind(&hook.Handler[*CollectionErrorEvent]{
		Id:       systemHookIdCollection,
		Func:     onErrorReloadCachedCollections,
		Priority: -99,
	})
	// ---

	app.OnBootstrap().Bind(&hook.Handler[*BootstrapEvent]{
		Id: systemHookIdCollection,
		Func: func(e *BootstrapEvent) error {
			if err := e.Next(); err != nil {
				return err
			}

			if err := e.App.ReloadCachedCollections(); err != nil {
				return fmt.Errorf("failed to load collections cache: %w", err)
			}

			return nil
		},
		Priority: 99, // execute as latest as possible
	})
}

// @todo experiment eventually replacing the rules *string with a struct?
type baseCollection struct {
	BaseModel

	disableIntegrityChecks bool
	autogeneratedId        string

	ListRule   *string `db:"listRule" json:"listRule" form:"listRule"`
	ViewRule   *string `db:"viewRule" json:"viewRule" form:"viewRule"`
	CreateRule *string `db:"createRule" json:"createRule" form:"createRule"`
	UpdateRule *string `db:"updateRule" json:"updateRule" form:"updateRule"`
	DeleteRule *string `db:"deleteRule" json:"deleteRule" form:"deleteRule"`

	// RawOptions represents the raw serialized collection option loaded from the DB.
	// NB! This field shouldn't be modified manually. It is automatically updated
	// with the collection type specific option before save.
	RawOptions types.JSONRaw `db:"options" json:"-" xml:"-" form:"-"`

	Name    string                  `db:"name" json:"name" form:"name"`
	Type    string                  `db:"type" json:"type" form:"type"`
	Fields  FieldsList              `db:"fields" json:"fields" form:"fields"`
	Indexes types.JSONArray[string] `db:"indexes" json:"indexes" form:"indexes"`
	Created types.DateTime          `db:"created" json:"created"`
	Updated types.DateTime          `db:"updated" json:"updated"`

	// System prevents the collection rename, deletion and rules change.
	// It is used primarily for internal purposes for collections like "_superusers", "_externalAuths", etc.
	System bool `db:"system" json:"system" form:"system"`
}

// Collection defines the table, fields and various options related to a set of records.
type Collection struct {
	baseCollection
	collectionAuthOptions
	collectionViewOptions
}

// NewCollection initializes and returns a new Collection model with the specified type and name.
//
// It also loads the minimal default configuration for the collection
// (eg. system fields, indexes, type specific options, etc.).
func NewCollection(typ, name string, optId ...string) *Collection {
	switch typ {
	case CollectionTypeAuth:
		return NewAuthCollection(name, optId...)
	case CollectionTypeView:
		return NewViewCollection(name, optId...)
	default:
		return NewBaseCollection(name, optId...)
	}
}

// NewBaseCollection initializes and returns a new "base" Collection model.
//
// It also loads the minimal default configuration for the collection
// (eg. system fields, indexes, type specific options, etc.).
func NewBaseCollection(name string, optId ...string) *Collection {
	m := &Collection{}
	m.Name = name
	m.Type = CollectionTypeBase

	// @todo consider removing once inferred composite literals are supported
	if len(optId) > 0 {
		m.Id = optId[0]
	}

	m.initDefaultId()
	m.initDefaultFields()

	return m
}

// NewViewCollection initializes and returns a new "view" Collection model.
//
// It also loads the minimal default configuration for the collection
// (eg. system fields, indexes, type specific options, etc.).
func NewViewCollection(name string, optId ...string) *Collection {
	m := &Collection{}
	m.Name = name
	m.Type = CollectionTypeView

	// @todo consider removing once inferred composite literals are supported
	if len(optId) > 0 {
		m.Id = optId[0]
	}

	m.initDefaultId()
	m.initDefaultFields()

	return m
}

// NewAuthCollection initializes and returns a new "auth" Collection model.
//
// It also loads the minimal default configuration for the collection
// (eg. system fields, indexes, type specific options, etc.).
func NewAuthCollection(name string, optId ...string) *Collection {
	m := &Collection{}
	m.Name = name
	m.Type = CollectionTypeAuth

	// @todo consider removing once inferred composite literals are supported
	if len(optId) > 0 {
		m.Id = optId[0]
	}

	m.initDefaultId()
	m.initDefaultFields()
	m.setDefaultAuthOptions()

	return m
}

// TableName returns the Collection model SQL table name.
func (m *Collection) TableName() string {
	return "_collections"
}

// BaseFilesPath returns the storage dir path used by the collection.
func (m *Collection) BaseFilesPath() string {
	return m.Id
}

// IsBase checks if the current collection has "base" type.
func (m *Collection) IsBase() bool {
	return m.Type == CollectionTypeBase
}

// IsAuth checks if the current collection has "auth" type.
func (m *Collection) IsAuth() bool {
	return m.Type == CollectionTypeAuth
}

// IsView checks if the current collection has "view" type.
func (m *Collection) IsView() bool {
	return m.Type == CollectionTypeView
}

// IntegrityChecks toggles the current collection integrity checks (ex. checking references on delete).
func (m *Collection) IntegrityChecks(enable bool) {
	m.disableIntegrityChecks = !enable
}

// PostScan implements the [dbx.PostScanner] interface to auto unmarshal
// the raw serialized options into the concrete type specific fields.
func (m *Collection) PostScan() error {
	if err := m.BaseModel.PostScan(); err != nil {
		return err
	}

	return m.unmarshalRawOptions()
}

func (m *Collection) unmarshalRawOptions() error {
	raw, err := m.RawOptions.MarshalJSON()
	if err != nil {
		return nil
	}

	switch m.Type {
	case CollectionTypeView:
		return json.Unmarshal(raw, &m.collectionViewOptions)
	case CollectionTypeAuth:
		return json.Unmarshal(raw, &m.collectionAuthOptions)
	}

	return nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
//
// For new/"blank" Collection models it replaces the model with a factory
// instance and then unmarshal the provided data one on top of it.
func (m *Collection) UnmarshalJSON(b []byte) error {
	type alias *Collection

	// initialize the default fields
	// (e.g. in case the collection was NOT created using the designated factories)
	if m.IsNew() && m.Type == "" {
		minimal := &struct {
			Type string `json:"type"`
			Name string `json:"name"`
			Id   string `json:"id"`
		}{}
		if err := json.Unmarshal(b, minimal); err != nil {
			return err
		}

		blank := NewCollection(minimal.Type, minimal.Name, minimal.Id)
		*m = *blank
	}

	return json.Unmarshal(b, alias(m))
}

// MarshalJSON implements the [json.Marshaler] interface.
//
// Note that non-type related fields are ignored from the serialization
// (ex. for "view" colections the "auth" fields are skipped).
func (m Collection) MarshalJSON() ([]byte, error) {
	switch m.Type {
	case CollectionTypeView:
		return json.Marshal(struct {
			baseCollection
			collectionViewOptions
		}{m.baseCollection, m.collectionViewOptions})
	case CollectionTypeAuth:
		alias := struct {
			baseCollection
			collectionAuthOptions
		}{m.baseCollection, m.collectionAuthOptions}

		// ensure that it is always returned as array
		if alias.OAuth2.Providers == nil {
			alias.OAuth2.Providers = []OAuth2ProviderConfig{}
		}

		// hide secret keys from the serialization
		alias.AuthToken.Secret = ""
		alias.FileToken.Secret = ""
		alias.PasswordResetToken.Secret = ""
		alias.EmailChangeToken.Secret = ""
		alias.VerificationToken.Secret = ""
		for i := range alias.OAuth2.Providers {
			alias.OAuth2.Providers[i].ClientSecret = ""
		}

		return json.Marshal(alias)
	default:
		return json.Marshal(m.baseCollection)
	}
}

// String returns a string representation of the current collection.
func (m Collection) String() string {
	raw, _ := json.Marshal(m)
	return string(raw)
}

// DBExport prepares and exports the current collection data for db persistence.
func (m *Collection) DBExport(app App) (map[string]any, error) {
	result := map[string]any{
		"id":         m.Id,
		"type":       m.Type,
		"listRule":   m.ListRule,
		"viewRule":   m.ViewRule,
		"createRule": m.CreateRule,
		"updateRule": m.UpdateRule,
		"deleteRule": m.DeleteRule,
		"name":       m.Name,
		"fields":     m.Fields,
		"indexes":    m.Indexes,
		"system":     m.System,
		"created":    m.Created,
		"updated":    m.Updated,
		"options":    `{}`,
	}

	switch m.Type {
	case CollectionTypeView:
		if raw, err := types.ParseJSONRaw(m.collectionViewOptions); err == nil {
			result["options"] = raw
		} else {
			return nil, err
		}
	case CollectionTypeAuth:
		if raw, err := types.ParseJSONRaw(m.collectionAuthOptions); err == nil {
			result["options"] = raw
		} else {
			return nil, err
		}
	}

	return result, nil
}

// GetIndex returns s single Collection index expression by its name.
func (m *Collection) GetIndex(name string) string {
	for _, idx := range m.Indexes {
		if strings.EqualFold(dbutils.ParseIndex(idx).IndexName, name) {
			return idx
		}
	}

	return ""
}

// AddIndex adds a new index into the current collection.
//
// If the collection has an existing index matching the new name it will be replaced with the new one.
func (m *Collection) AddIndex(name string, unique bool, columnsExpr string, optWhereExpr string) {
	m.RemoveIndex(name)

	var idx strings.Builder

	idx.WriteString("CREATE ")
	if unique {
		idx.WriteString("UNIQUE ")
	}
	idx.WriteString("INDEX `")
	idx.WriteString(name)
	idx.WriteString("` ")
	idx.WriteString("ON `")
	idx.WriteString(m.Name)
	idx.WriteString("` (")
	idx.WriteString(columnsExpr)
	idx.WriteString(")")
	if optWhereExpr != "" {
		idx.WriteString(" WHERE ")
		idx.WriteString(optWhereExpr)
	}

	m.Indexes = append(m.Indexes, idx.String())
}

// RemoveIndex removes a single index with the specified name from the current collection.
func (m *Collection) RemoveIndex(name string) {
	for i, idx := range m.Indexes {
		if strings.EqualFold(dbutils.ParseIndex(idx).IndexName, name) {
			m.Indexes = append(m.Indexes[:i], m.Indexes[i+1:]...)
			return
		}
	}
}

// delete hook
// -------------------------------------------------------------------

func onCollectionDeleteExecute(e *CollectionEvent) error {
	if e.Collection.System {
		return fmt.Errorf("[%s] system collections cannot be deleted", e.Collection.Name)
	}

	defer func() {
		if err := e.App.ReloadCachedCollections(); err != nil {
			e.App.Logger().Warn("Failed to reload collections cache", "error", err)
		}
	}()

	if !e.Collection.disableIntegrityChecks {
		// ensure that there aren't any existing references.
		// note: the select is outside of the transaction to prevent SQLITE_LOCKED error when mixing read&write in a single transaction
		references, err := e.App.FindCollectionReferences(e.Collection, e.Collection.Id)
		if err != nil {
			return fmt.Errorf("[%s] failed to check collection references: %w", e.Collection.Name, err)
		}
		if total := len(references); total > 0 {
			names := make([]string, 0, len(references))
			for ref := range references {
				names = append(names, ref.Name)
			}
			return fmt.Errorf("[%s] failed to delete due to existing relation references: %s", e.Collection.Name, strings.Join(names, ", "))
		}
	}

	originalApp := e.App

	txErr := e.App.RunInTransaction(func(txApp App) error {
		e.App = txApp

		// delete the related view or records table
		if e.Collection.IsView() {
			if err := txApp.DeleteView(e.Collection.Name); err != nil {
				return err
			}
		} else {
			if err := txApp.DeleteTable(e.Collection.Name); err != nil {
				return err
			}
		}

		if !e.Collection.disableIntegrityChecks {
			// trigger views resave to check for dependencies
			if err := resaveViewsWithChangedFields(txApp, e.Collection.Id); err != nil {
				return fmt.Errorf("[%s] failed to delete due to existing view dependency: %w", e.Collection.Name, err)
			}
		}

		// delete
		return e.Next()
	})

	e.App = originalApp

	return txErr
}

// save hook
// -------------------------------------------------------------------

func (c *Collection) idChecksum() string {
	return "pbc_" + crc32Checksum(c.Type+c.Name)
}

func (c *Collection) initDefaultId() {
	if c.Id == "" {
		c.Id = c.idChecksum()
		c.autogeneratedId = c.Id
	}
}

func (c *Collection) updateGeneratedIdIfExists(app App) {
	if !c.IsNew() ||
		// the id was explicitly cleared
		c.Id == "" ||
		// the id was manually set
		c.Id != c.autogeneratedId {
		return
	}

	// generate an up-to-date checksum
	newId := c.idChecksum()

	// add a number to the current id (if already exists)
	for i := 2; i < 1000; i++ {
		var exists int
		_ = app.CollectionQuery().Select("(1)").AndWhere(dbx.HashExp{"id": newId}).Limit(1).Row(&exists)
		if exists == 0 {
			break
		}
		newId = c.idChecksum() + strconv.Itoa(i)
	}

	// no change
	if c.Id == newId {
		return
	}

	// replace the old id in the index names (if any)
	for i, idx := range c.Indexes {
		parsed := dbutils.ParseIndex(idx)
		original := parsed.IndexName
		parsed.IndexName = strings.ReplaceAll(parsed.IndexName, c.Id, newId)
		if parsed.IndexName != original {
			c.Indexes[i] = parsed.Build()
		}
	}

	// update model id
	c.Id = newId
}

func onCollectionSave(e *CollectionEvent) error {
	if e.Collection.Type == "" {
		e.Collection.Type = CollectionTypeBase
	}

	if e.Collection.IsNew() {
		e.Collection.initDefaultId()
		e.Collection.Created = types.NowDateTime()
	}

	e.Collection.Updated = types.NowDateTime()

	// recreate the fields list to ensure that all normalizations
	// like default field id are applied
	e.Collection.Fields = NewFieldsList(e.Collection.Fields...)

	e.Collection.initDefaultFields()

	if e.Collection.IsAuth() {
		e.Collection.unsetMissingOAuth2MappedFields()
	}

	e.Collection.updateGeneratedIdIfExists(e.App)

	return e.Next()
}

func onCollectionSaveExecute(e *CollectionEvent) error {
	defer func() {
		if err := e.App.ReloadCachedCollections(); err != nil {
			e.App.Logger().Warn("Failed to reload collections cache", "error", err)
		}
	}()

	var oldCollection *Collection
	if !e.Collection.IsNew() {
		var err error
		oldCollection, err = e.App.FindCachedCollectionByNameOrId(e.Collection.Id)
		if err != nil {
			return err
		}

		// invalidate previously issued auth tokens on auth rule change
		if oldCollection.AuthRule != e.Collection.AuthRule &&
			cast.ToString(oldCollection.AuthRule) != cast.ToString(e.Collection.AuthRule) {
			e.Collection.AuthToken.Secret = security.RandomString(50)
		}
	}

	originalApp := e.App
	txErr := e.App.RunInTransaction(func(txApp App) error {
		e.App = txApp

		isView := e.Collection.IsView()

		// ensures that the view collection shema is properly loaded
		if isView {
			query := e.Collection.ViewQuery

			// generate collection fields list from the query
			viewFields, err := e.App.CreateViewFields(query)
			if err != nil {
				return err
			}

			// delete old renamed view
			if oldCollection != nil {
				if err := e.App.DeleteView(oldCollection.Name); err != nil {
					return err
				}
			}

			// wrap view query if necessary
			query, err = normalizeViewQueryId(e.App, query)
			if err != nil {
				return fmt.Errorf("failed to normalize view query id: %w", err)
			}

			// (re)create the view
			if err := e.App.SaveView(e.Collection.Name, query); err != nil {
				return err
			}

			// updates newCollection.Fields based on the generated view table info and query
			e.Collection.Fields = viewFields
		}

		// save the Collection model
		if err := e.Next(); err != nil {
			return err
		}

		// sync the changes with the related records table
		if !isView {
			if err := e.App.SyncRecordTableSchema(e.Collection, oldCollection); err != nil {
				// note: don't wrap to allow propagating indexes validation.Errors
				return err
			}
		}

		return nil
	})
	e.App = originalApp

	if txErr != nil {
		return txErr
	}

	// trigger an update for all views with changed fields as a result of the current collection save
	// (ignoring view errors to allow users to update the query from the UI)
	resaveViewsWithChangedFields(e.App, e.Collection.Id)

	return nil
}

func (c *Collection) initDefaultFields() {
	switch c.Type {
	case CollectionTypeBase:
		c.initIdField()
	case CollectionTypeAuth:
		c.initIdField()
		c.initPasswordField()
		c.initTokenKeyField()
		c.initEmailField()
		c.initEmailVisibilityField()
		c.initVerifiedField()
	case CollectionTypeView:
		// view fields are autogenerated
	}
}

func (c *Collection) initIdField() {
	field, _ := c.Fields.GetByName(FieldNameId).(*TextField)
	if field == nil {
		// create default field
		field = &TextField{
			Name:                FieldNameId,
			System:              true,
			PrimaryKey:          true,
			Required:            true,
			Min:                 15,
			Max:                 15,
			Pattern:             defaultLowercaseRecordIdPattern,
			AutogeneratePattern: `[a-z0-9]{15}`,
		}

		// prepend it
		c.Fields = NewFieldsList(append([]Field{field}, c.Fields...)...)
	} else {
		// enforce system defaults
		field.System = true
		field.Required = true
		field.PrimaryKey = true
		field.Hidden = false
		if field.Pattern == "" {
			field.Pattern = defaultLowercaseRecordIdPattern
		}
	}
}

func (c *Collection) initPasswordField() {
	field, _ := c.Fields.GetByName(FieldNamePassword).(*PasswordField)
	if field == nil {
		// load default field
		c.Fields.Add(&PasswordField{
			Name:     FieldNamePassword,
			System:   true,
			Hidden:   true,
			Required: true,
			Min:      8,
		})
	} else {
		// enforce system defaults
		field.System = true
		field.Hidden = true
		field.Required = true
	}
}

func (c *Collection) initTokenKeyField() {
	field, _ := c.Fields.GetByName(FieldNameTokenKey).(*TextField)
	if field == nil {
		// load default field
		c.Fields.Add(&TextField{
			Name:                FieldNameTokenKey,
			System:              true,
			Hidden:              true,
			Min:                 30,
			Max:                 60,
			Required:            true,
			AutogeneratePattern: `[a-zA-Z0-9]{50}`,
		})
	} else {
		// enforce system defaults
		field.System = true
		field.Hidden = true
		field.Required = true
	}

	// ensure that there is a unique index for the field
	if _, ok := dbutils.FindSingleColumnUniqueIndex(c.Indexes, FieldNameTokenKey); !ok {
		c.Indexes = append(c.Indexes, fmt.Sprintf(
			"CREATE UNIQUE INDEX `%s` ON `%s` (`%s`)",
			c.fieldIndexName(FieldNameTokenKey),
			c.Name,
			FieldNameTokenKey,
		))
	}
}

func (c *Collection) initEmailField() {
	field, _ := c.Fields.GetByName(FieldNameEmail).(*EmailField)
	if field == nil {
		// load default field
		c.Fields.Add(&EmailField{
			Name:     FieldNameEmail,
			System:   true,
			Required: true,
		})
	} else {
		// enforce system defaults
		field.System = true
		field.Hidden = false // managed by the emailVisibility flag
	}

	// ensure that there is a unique index for the email field
	if _, ok := dbutils.FindSingleColumnUniqueIndex(c.Indexes, FieldNameEmail); !ok {
		c.Indexes = append(c.Indexes, fmt.Sprintf(
			"CREATE UNIQUE INDEX `%s` ON `%s` (`%s`) WHERE `%s` != ''",
			c.fieldIndexName(FieldNameEmail),
			c.Name,
			FieldNameEmail,
			FieldNameEmail,
		))
	}
}

func (c *Collection) initEmailVisibilityField() {
	field, _ := c.Fields.GetByName(FieldNameEmailVisibility).(*BoolField)
	if field == nil {
		// load default field
		c.Fields.Add(&BoolField{
			Name:   FieldNameEmailVisibility,
			System: true,
		})
	} else {
		// enforce system defaults
		field.System = true
	}
}

func (c *Collection) initVerifiedField() {
	field, _ := c.Fields.GetByName(FieldNameVerified).(*BoolField)
	if field == nil {
		// load default field
		c.Fields.Add(&BoolField{
			Name:   FieldNameVerified,
			System: true,
		})
	} else {
		// enforce system defaults
		field.System = true
	}
}

func (c *Collection) fieldIndexName(field string) string {
	name := "idx_" + field + "_"

	if c.Id != "" {
		name += c.Id
	} else if c.Name != "" {
		name += c.Name
	} else {
		name += security.PseudorandomString(10)
	}

	if len(name) > 64 {
		return name[:64]
	}

	return name
}
