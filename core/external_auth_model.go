package core

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/types"
)

var (
	_ Model        = (*ExternalAuth)(nil)
	_ PreValidator = (*ExternalAuth)(nil)
	_ RecordProxy  = (*ExternalAuth)(nil)
)

const CollectionNameExternalAuths = "_externalAuths"

// ExternalAuth defines a Record proxy for working with the externalAuths collection.
type ExternalAuth struct {
	*Record
}

// NewExternalAuth instantiates and returns a new blank *ExternalAuth model.
//
// Example usage:
//
//	ea := core.NewExternalAuth(app)
//	ea.SetRecordRef(user.Id)
//	ea.SetCollectionRef(user.Collection().Id)
//	ea.SetProvider("google")
//	ea.SetProviderId("...")
//	app.Save(ea)
func NewExternalAuth(app App) *ExternalAuth {
	m := &ExternalAuth{}

	c, err := app.FindCachedCollectionByNameOrId(CollectionNameExternalAuths)
	if err != nil {
		// this is just to make tests easier since it is a system collection and it is expected to be always accessible
		// (note: the loaded record is further checked on ExternalAuth.PreValidate())
		c = NewBaseCollection("@__invalid__")
	}

	m.Record = NewRecord(c)

	return m
}

// PreValidate implements the [PreValidator] interface and checks
// whether the proxy is properly loaded.
func (m *ExternalAuth) PreValidate(ctx context.Context, app App) error {
	if m.Record == nil || m.Record.Collection().Name != CollectionNameExternalAuths {
		return errors.New("missing or invalid ExternalAuth ProxyRecord")
	}

	return nil
}

// ProxyRecord returns the proxied Record model.
func (m *ExternalAuth) ProxyRecord() *Record {
	return m.Record
}

// SetProxyRecord loads the specified record model into the current proxy.
func (m *ExternalAuth) SetProxyRecord(record *Record) {
	m.Record = record
}

// CollectionRef returns the "collectionRef" field value.
func (m *ExternalAuth) CollectionRef() string {
	return m.GetString("collectionRef")
}

// SetCollectionRef updates the "collectionRef" record field value.
func (m *ExternalAuth) SetCollectionRef(collectionId string) {
	m.Set("collectionRef", collectionId)
}

// RecordRef returns the "recordRef" record field value.
func (m *ExternalAuth) RecordRef() string {
	return m.GetString("recordRef")
}

// SetRecordRef updates the "recordRef" record field value.
func (m *ExternalAuth) SetRecordRef(recordId string) {
	m.Set("recordRef", recordId)
}

// Provider returns the "provider" record field value.
func (m *ExternalAuth) Provider() string {
	return m.GetString("provider")
}

// SetProvider updates the "provider" record field value.
func (m *ExternalAuth) SetProvider(provider string) {
	m.Set("provider", provider)
}

// Provider returns the "providerId" record field value.
func (m *ExternalAuth) ProviderId() string {
	return m.GetString("providerId")
}

// SetProvider updates the "providerId" record field value.
func (m *ExternalAuth) SetProviderId(providerId string) {
	m.Set("providerId", providerId)
}

// Created returns the "created" record field value.
func (m *ExternalAuth) Created() types.DateTime {
	return m.GetDateTime("created")
}

// Updated returns the "updated" record field value.
func (m *ExternalAuth) Updated() types.DateTime {
	return m.GetDateTime("updated")
}

func (app *BaseApp) registerExternalAuthHooks() {
	recordRefHooks[*ExternalAuth](app, CollectionNameExternalAuths, CollectionTypeAuth)

	app.OnRecordValidate(CollectionNameExternalAuths).Bind(&hook.Handler[*RecordEvent]{
		Func: func(e *RecordEvent) error {
			providerNames := make([]any, 0, len(auth.Providers))
			for name := range auth.Providers {
				providerNames = append(providerNames, name)
			}

			provider := e.Record.GetString("provider")
			if err := validation.Validate(provider, validation.Required, validation.In(providerNames...)); err != nil {
				return validation.Errors{"provider": err}
			}

			return e.Next()
		},
		Priority: 99,
	})

	// delete all pre-existing external auths on verified upgrade
	app.OnRecordUpdateExecute().Bind(&hook.Handler[*RecordEvent]{
		Func: func(e *RecordEvent) error {
			if !e.Record.Collection().IsAuth() {
				return e.Next()
			}

			hasUpgradedVerified := !e.Record.Original().IsNew() && !e.Record.Original().Verified() && e.Record.Verified()

			if !hasUpgradedVerified {
				return e.Next()
			}

			originalApp := e.App
			return e.App.RunInTransaction(func(txApp App) error {
				e.App = txApp
				defer func() { e.App = originalApp }()

				externalAuths, err := txApp.FindAllExternalAuthsByRecord(e.Record)
				if err != nil {
					return err
				}
				if len(externalAuths) > 0 {
					// delete all pre-existing external auths
					if err := txApp.DeleteAllExternalAuthsByRecord(e.Record); err != nil {
						return err
					}

					// force refresh tokens reset (if not already)
					e.Record.RefreshTokenKey()
				}

				return e.Next()
			})
		},
		Priority: 99,
	})
}
