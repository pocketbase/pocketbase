package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// CollectionNameWebAuthnCredentials is the name of the system collection
// that stores WebAuthn/Passkey credentials linked to auth records.
const CollectionNameWebAuthnCredentials = "_webauthnCredentials"

var (
	_ Model        = (*WebAuthnCredential)(nil)
	_ PreValidator = (*WebAuthnCredential)(nil)
	_ RecordProxy  = (*WebAuthnCredential)(nil)
)

// WebAuthnCredential defines a Record proxy for working with the
// _webauthnCredentials collection.
type WebAuthnCredential struct {
	*Record
}

// NewWebAuthnCredential instantiates and returns a new blank *WebAuthnCredential model.
//
// Example usage:
//
//	wc := core.NewWebAuthnCredential(app)
//	wc.SetRecordRef(user.Id)
//	wc.SetCollectionRef(user.Collection().Id)
//	wc.SetCredentialId("base64url-encoded-id")
//	wc.SetPublicKey("base64url-encoded-key")
//	app.Save(wc)
func NewWebAuthnCredential(app App) *WebAuthnCredential {
	m := &WebAuthnCredential{}

	c, err := app.FindCachedCollectionByNameOrId(CollectionNameWebAuthnCredentials)
	if err != nil {
		// this is just to make tests easier since it is a system collection and it is expected to be always accessible
		// (note: the loaded record is further checked on WebAuthnCredential.PreValidate())
		c = NewBaseCollection("@__invalid__")
	}

	m.Record = NewRecord(c)

	return m
}

// PreValidate implements the [PreValidator] interface and checks
// whether the proxy is properly loaded.
func (m *WebAuthnCredential) PreValidate(ctx context.Context, app App) error {
	if m.Record == nil || m.Record.Collection().Name != CollectionNameWebAuthnCredentials {
		return errors.New("missing or invalid WebAuthnCredential ProxyRecord")
	}

	return nil
}

// ProxyRecord returns the proxied Record model.
func (m *WebAuthnCredential) ProxyRecord() *Record {
	return m.Record
}

// SetProxyRecord loads the specified record model into the current proxy.
func (m *WebAuthnCredential) SetProxyRecord(record *Record) {
	m.Record = record
}

// CollectionRef returns the "collectionRef" field value.
func (m *WebAuthnCredential) CollectionRef() string {
	return m.GetString("collectionRef")
}

// SetCollectionRef updates the "collectionRef" record field value.
func (m *WebAuthnCredential) SetCollectionRef(collectionId string) {
	m.Set("collectionRef", collectionId)
}

// RecordRef returns the "recordRef" record field value.
func (m *WebAuthnCredential) RecordRef() string {
	return m.GetString("recordRef")
}

// SetRecordRef updates the "recordRef" record field value.
func (m *WebAuthnCredential) SetRecordRef(recordId string) {
	m.Set("recordRef", recordId)
}

// CredentialId returns the "credentialId" field value (base64url-encoded).
func (m *WebAuthnCredential) CredentialId() string {
	return m.GetString("credentialId")
}

// SetCredentialId updates the "credentialId" record field value.
func (m *WebAuthnCredential) SetCredentialId(id string) {
	m.Set("credentialId", id)
}

// PublicKey returns the "publicKey" field value (base64url-encoded).
func (m *WebAuthnCredential) PublicKey() string {
	return m.GetString("publicKey")
}

// SetPublicKey updates the "publicKey" record field value.
func (m *WebAuthnCredential) SetPublicKey(key string) {
	m.Set("publicKey", key)
}

// AttestationType returns the "attestationType" field value.
func (m *WebAuthnCredential) AttestationType() string {
	return m.GetString("attestationType")
}

// SetAttestationType updates the "attestationType" record field value.
func (m *WebAuthnCredential) SetAttestationType(at string) {
	m.Set("attestationType", at)
}

// Transport returns the "transport" field value as a slice of strings.
func (m *WebAuthnCredential) Transport() []string {
	raw := m.GetString("transport")
	if raw == "" {
		return nil
	}

	var result []string
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil
	}
	return result
}

// SetTransport updates the "transport" record field value.
func (m *WebAuthnCredential) SetTransport(transports []string) {
	m.Set("transport", transports)
}

// SignCount returns the "signCount" field value.
func (m *WebAuthnCredential) SignCount() uint32 {
	return uint32(m.GetInt("signCount"))
}

// SetSignCount updates the "signCount" record field value.
func (m *WebAuthnCredential) SetSignCount(count uint32) {
	m.Set("signCount", int(count))
}

// AAGUID returns the "aaguid" field value.
func (m *WebAuthnCredential) AAGUID() string {
	return m.GetString("aaguid")
}

// SetAAGUID updates the "aaguid" record field value.
func (m *WebAuthnCredential) SetAAGUID(val string) {
	m.Set("aaguid", val)
}

// Name returns the "name" field value (user-friendly credential label).
func (m *WebAuthnCredential) Name() string {
	return m.GetString("name")
}

// SetName updates the "name" record field value.
func (m *WebAuthnCredential) SetName(name string) {
	m.Set("name", name)
}

// Flags returns the "flags" field as a decoded CredentialFlags map.
func (m *WebAuthnCredential) Flags() webauthn.CredentialFlags {
	raw := m.GetString("flags")
	if raw == "" {
		return webauthn.CredentialFlags{}
	}

	var flags webauthn.CredentialFlags
	if err := json.Unmarshal([]byte(raw), &flags); err != nil {
		return webauthn.CredentialFlags{}
	}
	return flags
}

// SetFlags updates the "flags" record field value.
func (m *WebAuthnCredential) SetFlags(flags webauthn.CredentialFlags) {
	m.Set("flags", flags)
}

// ToWebAuthnCredential converts the stored base64url-encoded credential fields
// into a go-webauthn Credential struct for cryptographic operations.
func (m *WebAuthnCredential) ToWebAuthnCredential() (webauthn.Credential, error) {
	credId, err := base64.RawURLEncoding.DecodeString(m.CredentialId())
	if err != nil {
		return webauthn.Credential{}, errors.New("invalid credentialId encoding")
	}

	pubKey, err := base64.RawURLEncoding.DecodeString(m.PublicKey())
	if err != nil {
		return webauthn.Credential{}, errors.New("invalid publicKey encoding")
	}

	transports := m.Transport()
	protoTransports := make([]protocol.AuthenticatorTransport, len(transports))
	for i, t := range transports {
		protoTransports[i] = protocol.AuthenticatorTransport(t)
	}

	authenticator := webauthn.Authenticator{
		SignCount: m.SignCount(),
	}

	// Restore AAGUID if stored
	if aaguidStr := m.AAGUID(); aaguidStr != "" {
		if aaguidBytes, err := base64.RawURLEncoding.DecodeString(aaguidStr); err == nil {
			authenticator.AAGUID = aaguidBytes
		}
	}

	return webauthn.Credential{
		ID:              credId,
		PublicKey:       pubKey,
		AttestationType: m.AttestationType(),
		Transport:       protoTransports,
		Flags:           m.Flags(),
		Authenticator:   authenticator,
	}, nil
}

// FromWebAuthnCredential populates the record fields from a go-webauthn Credential struct.
func (m *WebAuthnCredential) FromWebAuthnCredential(cred webauthn.Credential) {
	m.SetCredentialId(base64.RawURLEncoding.EncodeToString(cred.ID))
	m.SetPublicKey(base64.RawURLEncoding.EncodeToString(cred.PublicKey))
	m.SetAttestationType(cred.AttestationType)
	m.SetSignCount(cred.Authenticator.SignCount)
	m.SetFlags(cred.Flags)

	if len(cred.Authenticator.AAGUID) > 0 {
		m.SetAAGUID(base64.RawURLEncoding.EncodeToString(cred.Authenticator.AAGUID))
	}

	transports := make([]string, len(cred.Transport))
	for i, t := range cred.Transport {
		transports[i] = string(t)
	}
	m.SetTransport(transports)
}

// registerWebAuthnCredentialHooks registers cascading delete hooks
// so that credentials are cleaned up when related collections or records are deleted.
func (app *BaseApp) registerWebAuthnCredentialHooks() {
	recordRefHooks[*WebAuthnCredential](app, CollectionNameWebAuthnCredentials, CollectionTypeAuth)
}
