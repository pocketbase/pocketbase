package core

import (
	"context"
	"errors"
	"time"

	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	MFAMethodPassword = "password"
	MFAMethodOAuth2   = "oauth2"
	MFAMethodOTP      = "otp"
)

const CollectionNameMFAs = "_mfas"

var (
	_ Model        = (*MFA)(nil)
	_ PreValidator = (*MFA)(nil)
	_ RecordProxy  = (*MFA)(nil)
)

// MFA defines a Record proxy for working with the mfas collection.
type MFA struct {
	*Record
}

// NewMFA instantiates and returns a new blank *MFA model.
//
// Example usage:
//
//	mfa := core.NewMFA(app)
//	mfa.SetRecordRef(user.Id)
//	mfa.SetCollectionRef(user.Collection().Id)
//	mfa.SetMethod(core.MFAMethodPassword)
//	app.Save(mfa)
func NewMFA(app App) *MFA {
	m := &MFA{}

	c, err := app.FindCachedCollectionByNameOrId(CollectionNameMFAs)
	if err != nil {
		// this is just to make tests easier since mfa is a system collection and it is expected to be always accessible
		// (note: the loaded record is further checked on MFA.PreValidate())
		c = NewBaseCollection("@__invalid__")
	}

	m.Record = NewRecord(c)

	return m
}

// PreValidate implements the [PreValidator] interface and checks
// whether the proxy is properly loaded.
func (m *MFA) PreValidate(ctx context.Context, app App) error {
	if m.Record == nil || m.Record.Collection().Name != CollectionNameMFAs {
		return errors.New("missing or invalid mfa ProxyRecord")
	}

	return nil
}

// ProxyRecord returns the proxied Record model.
func (m *MFA) ProxyRecord() *Record {
	return m.Record
}

// SetProxyRecord loads the specified record model into the current proxy.
func (m *MFA) SetProxyRecord(record *Record) {
	m.Record = record
}

// CollectionRef returns the "collectionRef" field value.
func (m *MFA) CollectionRef() string {
	return m.GetString("collectionRef")
}

// SetCollectionRef updates the "collectionRef" record field value.
func (m *MFA) SetCollectionRef(collectionId string) {
	m.Set("collectionRef", collectionId)
}

// RecordRef returns the "recordRef" record field value.
func (m *MFA) RecordRef() string {
	return m.GetString("recordRef")
}

// SetRecordRef updates the "recordRef" record field value.
func (m *MFA) SetRecordRef(recordId string) {
	m.Set("recordRef", recordId)
}

// Method returns the "method" record field value.
func (m *MFA) Method() string {
	return m.GetString("method")
}

// SetMethod updates the "method" record field value.
func (m *MFA) SetMethod(method string) {
	m.Set("method", method)
}

// Created returns the "created" record field value.
func (m *MFA) Created() types.DateTime {
	return m.GetDateTime("created")
}

// Updated returns the "updated" record field value.
func (m *MFA) Updated() types.DateTime {
	return m.GetDateTime("updated")
}

// HasExpired checks if the mfa is expired, aka. whether it has been
// more than maxElapsed time since its creation.
func (m *MFA) HasExpired(maxElapsed time.Duration) bool {
	return time.Since(m.Created().Time()) > maxElapsed
}

func (app *BaseApp) registerMFAHooks() {
	recordRefHooks[*MFA](app, CollectionNameMFAs, CollectionTypeAuth)

	// run on every hour to cleanup expired mfa sessions
	app.Cron().Add("__pbMFACleanup__", "0 * * * *", func() {
		if err := app.DeleteExpiredMFAs(); err != nil {
			app.Logger().Warn("Failed to delete expired MFA sessions", "error", err)
		}
	})

	// delete existing mfas on password change
	app.OnRecordUpdate().Bind(&hook.Handler[*RecordEvent]{
		Func: func(e *RecordEvent) error {
			err := e.Next()
			if err != nil || !e.Record.Collection().IsAuth() {
				return err
			}

			old := e.Record.Original().GetString(FieldNamePassword + ":hash")
			new := e.Record.GetString(FieldNamePassword + ":hash")
			if old != new {
				err = e.App.DeleteAllMFAsByRecord(e.Record)
				if err != nil {
					e.App.Logger().Warn(
						"Failed to delete all previous mfas",
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
