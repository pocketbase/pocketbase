package core

import (
	"context"
	"errors"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/types"
)

const CollectionNameAuthOrigins = "_authOrigins"

var (
	_ Model        = (*AuthOrigin)(nil)
	_ PreValidator = (*AuthOrigin)(nil)
	_ RecordProxy  = (*AuthOrigin)(nil)
)

// AuthOrigin defines a Record proxy for working with the authOrigins collection.
type AuthOrigin struct {
	*Record
}

// NewAuthOrigin instantiates and returns a new blank *AuthOrigin model.
//
// Example usage:
//
//	origin := core.NewOrigin(app)
//	origin.SetRecordRef(user.Id)
//	origin.SetCollectionRef(user.Collection().Id)
//	origin.SetFingerprint("...")
//	app.Save(origin)
func NewAuthOrigin(app App) *AuthOrigin {
	m := &AuthOrigin{}

	c, err := app.FindCachedCollectionByNameOrId(CollectionNameAuthOrigins)
	if err != nil {
		// this is just to make tests easier since authOrigins is a system collection and it is expected to be always accessible
		// (note: the loaded record is further checked on AuthOrigin.PreValidate())
		c = NewBaseCollection("@___invalid___")
	}

	m.Record = NewRecord(c)

	return m
}

// PreValidate implements the [PreValidator] interface and checks
// whether the proxy is properly loaded.
func (m *AuthOrigin) PreValidate(ctx context.Context, app App) error {
	if m.Record == nil || m.Record.Collection().Name != CollectionNameAuthOrigins {
		return errors.New("missing or invalid AuthOrigin ProxyRecord")
	}

	return nil
}

// ProxyRecord returns the proxied Record model.
func (m *AuthOrigin) ProxyRecord() *Record {
	return m.Record
}

// SetProxyRecord loads the specified record model into the current proxy.
func (m *AuthOrigin) SetProxyRecord(record *Record) {
	m.Record = record
}

// CollectionRef returns the "collectionRef" field value.
func (m *AuthOrigin) CollectionRef() string {
	return m.GetString("collectionRef")
}

// SetCollectionRef updates the "collectionRef" record field value.
func (m *AuthOrigin) SetCollectionRef(collectionId string) {
	m.Set("collectionRef", collectionId)
}

// RecordRef returns the "recordRef" record field value.
func (m *AuthOrigin) RecordRef() string {
	return m.GetString("recordRef")
}

// SetRecordRef updates the "recordRef" record field value.
func (m *AuthOrigin) SetRecordRef(recordId string) {
	m.Set("recordRef", recordId)
}

// Fingerprint returns the "fingerprint" record field value.
func (m *AuthOrigin) Fingerprint() string {
	return m.GetString("fingerprint")
}

// SetFingerprint updates the "fingerprint" record field value.
func (m *AuthOrigin) SetFingerprint(fingerprint string) {
	m.Set("fingerprint", fingerprint)
}

// Created returns the "created" record field value.
func (m *AuthOrigin) Created() types.DateTime {
	return m.GetDateTime("created")
}

// Updated returns the "updated" record field value.
func (m *AuthOrigin) Updated() types.DateTime {
	return m.GetDateTime("updated")
}

func (app *BaseApp) registerAuthOriginHooks() {
	recordRefHooks[*AuthOrigin](app, CollectionNameAuthOrigins, CollectionTypeAuth)

	// delete existing auth origins on password change
	app.OnRecordUpdate().Bind(&hook.Handler[*RecordEvent]{
		Func: func(e *RecordEvent) error {
			err := e.Next()
			if err != nil || !e.Record.Collection().IsAuth() {
				return err
			}

			old := e.Record.Original().GetString(FieldNamePassword + ":hash")
			new := e.Record.GetString(FieldNamePassword + ":hash")
			if old != new {
				err = e.App.DeleteAllAuthOriginsByRecord(e.Record)
				if err != nil {
					e.App.Logger().Warn(
						"Failed to delete all previous auth origin fingerprints",
						"error", err,
						"recordId", e.Record.Id,
						"collectionId", e.Record.Collection().Id,
					)
				}
			}

			return nil
		},
		Priority: 99,
	})
}

// -------------------------------------------------------------------

// recordRefHooks registers common hooks that are usually used with record proxies
// that have polymorphic record relations (aka. "collectionRef" and "recordRef" fields).
func recordRefHooks[T RecordProxy](app App, collectionName string, optCollectionTypes ...string) {
	app.OnRecordValidate(collectionName).Bind(&hook.Handler[*RecordEvent]{
		Func: func(e *RecordEvent) error {
			collectionId := e.Record.GetString("collectionRef")
			err := validation.Validate(collectionId, validation.Required, validation.By(validateCollectionId(e.App, optCollectionTypes...)))
			if err != nil {
				return validation.Errors{"collectionRef": err}
			}

			recordId := e.Record.GetString("recordRef")
			err = validation.Validate(recordId, validation.Required, validation.By(validateRecordId(e.App, collectionId)))
			if err != nil {
				return validation.Errors{"recordRef": err}
			}

			return e.Next()
		},
		Priority: 99,
	})

	// delete on collection ref delete
	app.OnCollectionDeleteExecute().Bind(&hook.Handler[*CollectionEvent]{
		Func: func(e *CollectionEvent) error {
			if e.Collection.Name == collectionName || (len(optCollectionTypes) > 0 && !slices.Contains(optCollectionTypes, e.Collection.Type)) {
				return e.Next()
			}

			originalApp := e.App
			txErr := e.App.RunInTransaction(func(txApp App) error {
				e.App = txApp

				if err := e.Next(); err != nil {
					return err
				}

				rels, err := txApp.FindAllRecords(collectionName, dbx.HashExp{"collectionRef": e.Collection.Id})
				if err != nil {
					return err
				}

				for _, mfa := range rels {
					if err := txApp.Delete(mfa); err != nil {
						return err
					}
				}

				return nil
			})
			e.App = originalApp

			return txErr
		},
		Priority: 99,
	})

	// delete on record ref delete
	app.OnRecordDeleteExecute().Bind(&hook.Handler[*RecordEvent]{
		Func: func(e *RecordEvent) error {
			if e.Record.Collection().Name == collectionName ||
				(len(optCollectionTypes) > 0 && !slices.Contains(optCollectionTypes, e.Record.Collection().Type)) {
				return e.Next()
			}

			originalApp := e.App
			txErr := e.App.RunInTransaction(func(txApp App) error {
				e.App = txApp

				if err := e.Next(); err != nil {
					return err
				}

				rels, err := txApp.FindAllRecords(collectionName, dbx.HashExp{
					"collectionRef": e.Record.Collection().Id,
					"recordRef":     e.Record.Id,
				})
				if err != nil {
					return err
				}

				for _, rel := range rels {
					if err := txApp.Delete(rel); err != nil {
						return err
					}
				}

				return nil
			})
			e.App = originalApp

			return txErr
		},
		Priority: 99,
	})
}
