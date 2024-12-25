package core

import (
	"context"
	"errors"
	"time"

	"github.com/pocketbase/pocketbase/tools/types"
)

const CollectionNameOTPs = "_otps"

var (
	_ Model        = (*OTP)(nil)
	_ PreValidator = (*OTP)(nil)
	_ RecordProxy  = (*OTP)(nil)
)

// OTP defines a Record proxy for working with the otps collection.
type OTP struct {
	*Record
}

// NewOTP instantiates and returns a new blank *OTP model.
//
// Example usage:
//
//	otp := core.NewOTP(app)
//	otp.SetRecordRef(user.Id)
//	otp.SetCollectionRef(user.Collection().Id)
//	otp.SetPassword(security.RandomStringWithAlphabet(6, "1234567890"))
//	app.Save(otp)
func NewOTP(app App) *OTP {
	m := &OTP{}

	c, err := app.FindCachedCollectionByNameOrId(CollectionNameOTPs)
	if err != nil {
		// this is just to make tests easier since otp is a system collection and it is expected to be always accessible
		// (note: the loaded record is further checked on OTP.PreValidate())
		c = NewBaseCollection("__invalid__")
	}

	m.Record = NewRecord(c)

	return m
}

// PreValidate implements the [PreValidator] interface and checks
// whether the proxy is properly loaded.
func (m *OTP) PreValidate(ctx context.Context, app App) error {
	if m.Record == nil || m.Record.Collection().Name != CollectionNameOTPs {
		return errors.New("missing or invalid otp ProxyRecord")
	}

	return nil
}

// ProxyRecord returns the proxied Record model.
func (m *OTP) ProxyRecord() *Record {
	return m.Record
}

// SetProxyRecord loads the specified record model into the current proxy.
func (m *OTP) SetProxyRecord(record *Record) {
	m.Record = record
}

// CollectionRef returns the "collectionRef" field value.
func (m *OTP) CollectionRef() string {
	return m.GetString("collectionRef")
}

// SetCollectionRef updates the "collectionRef" record field value.
func (m *OTP) SetCollectionRef(collectionId string) {
	m.Set("collectionRef", collectionId)
}

// RecordRef returns the "recordRef" record field value.
func (m *OTP) RecordRef() string {
	return m.GetString("recordRef")
}

// SetRecordRef updates the "recordRef" record field value.
func (m *OTP) SetRecordRef(recordId string) {
	m.Set("recordRef", recordId)
}

// SentTo returns the "sentTo" record field value.
//
// It could be any string value (email, phone, message app id, etc.)
// and usually is used as part of the auth flow to update the verified
// user state in case for example the sentTo value matches with the user record email.
func (m *OTP) SentTo() string {
	return m.GetString("sentTo")
}

// SetSentTo updates the "sentTo" record field value.
func (m *OTP) SetSentTo(val string) {
	m.Set("sentTo", val)
}

// Created returns the "created" record field value.
func (m *OTP) Created() types.DateTime {
	return m.GetDateTime("created")
}

// Updated returns the "updated" record field value.
func (m *OTP) Updated() types.DateTime {
	return m.GetDateTime("updated")
}

// HasExpired checks if the otp is expired, aka. whether it has been
// more than maxElapsed time since its creation.
func (m *OTP) HasExpired(maxElapsed time.Duration) bool {
	return time.Since(m.Created().Time()) > maxElapsed
}

func (app *BaseApp) registerOTPHooks() {
	recordRefHooks[*OTP](app, CollectionNameOTPs, CollectionTypeAuth)

	// run on every hour to cleanup expired otp sessions
	app.Cron().Add("__pbOTPCleanup__", "0 * * * *", func() {
		if err := app.DeleteExpiredOTPs(); err != nil {
			app.Logger().Warn("Failed to delete expired OTP sessions", "error", err)
		}
	})
}
