package core

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/pocketbase/tools/security"
)

// Supported record token types
const (
	TokenTypeAuth          = "auth"
	TokenTypeFile          = "file"
	TokenTypeVerification  = "verification"
	TokenTypePasswordReset = "passwordReset"
	TokenTypeEmailChange   = "emailChange"
)

// List with commonly used record token claims
const (
	TokenClaimId           = "id"
	TokenClaimType         = "type"
	TokenClaimCollectionId = "collectionId"
	TokenClaimEmail        = "email"
	TokenClaimNewEmail     = "newEmail"
	TokenClaimRefreshable  = "refreshable"
)

// Common token related errors
var (
	ErrNotAuthRecord     = errors.New("not an auth collection record")
	ErrMissingSigningKey = errors.New("missing or invalid signing key")
)

// NewStaticAuthToken generates and returns a new static record authentication token.
//
// Static auth tokens are similar to the regular auth tokens, but are
// non-refreshable and support custom duration.
//
// Zero or negative duration will fallback to the duration from the auth collection settings.
func (m *Record) NewStaticAuthToken(duration time.Duration) (string, error) {
	return m.newAuthToken(duration, false)
}

// NewAuthToken generates and returns a new record authentication token.
func (m *Record) NewAuthToken() (string, error) {
	return m.newAuthToken(0, true)
}

func (m *Record) newAuthToken(duration time.Duration, refreshable bool) (string, error) {
	if !m.Collection().IsAuth() {
		return "", ErrNotAuthRecord
	}

	key := (m.TokenKey() + m.Collection().AuthToken.Secret)
	if key == "" {
		return "", ErrMissingSigningKey
	}

	claims := jwt.MapClaims{
		TokenClaimType:         TokenTypeAuth,
		TokenClaimId:           m.Id,
		TokenClaimCollectionId: m.Collection().Id,
		TokenClaimRefreshable:  refreshable,
	}

	if duration <= 0 {
		duration = m.Collection().AuthToken.DurationTime()
	}

	return security.NewJWT(claims, key, duration)
}

// NewVerificationToken generates and returns a new record verification token.
func (m *Record) NewVerificationToken() (string, error) {
	if !m.Collection().IsAuth() {
		return "", ErrNotAuthRecord
	}

	key := (m.TokenKey() + m.Collection().VerificationToken.Secret)
	if key == "" {
		return "", ErrMissingSigningKey
	}

	return security.NewJWT(
		jwt.MapClaims{
			TokenClaimType:         TokenTypeVerification,
			TokenClaimId:           m.Id,
			TokenClaimCollectionId: m.Collection().Id,
			TokenClaimEmail:        m.Email(),
		},
		key,
		m.Collection().VerificationToken.DurationTime(),
	)
}

// NewPasswordResetToken generates and returns a new auth record password reset request token.
func (m *Record) NewPasswordResetToken() (string, error) {
	if !m.Collection().IsAuth() {
		return "", ErrNotAuthRecord
	}

	key := (m.TokenKey() + m.Collection().PasswordResetToken.Secret)
	if key == "" {
		return "", ErrMissingSigningKey
	}

	return security.NewJWT(
		jwt.MapClaims{
			TokenClaimType:         TokenTypePasswordReset,
			TokenClaimId:           m.Id,
			TokenClaimCollectionId: m.Collection().Id,
			TokenClaimEmail:        m.Email(),
		},
		key,
		m.Collection().PasswordResetToken.DurationTime(),
	)
}

// NewEmailChangeToken generates and returns a new auth record change email request token.
func (m *Record) NewEmailChangeToken(newEmail string) (string, error) {
	if !m.Collection().IsAuth() {
		return "", ErrNotAuthRecord
	}

	key := (m.TokenKey() + m.Collection().EmailChangeToken.Secret)
	if key == "" {
		return "", ErrMissingSigningKey
	}

	return security.NewJWT(
		jwt.MapClaims{
			TokenClaimType:         TokenTypeEmailChange,
			TokenClaimId:           m.Id,
			TokenClaimCollectionId: m.Collection().Id,
			TokenClaimEmail:        m.Email(),
			TokenClaimNewEmail:     newEmail,
		},
		key,
		m.Collection().EmailChangeToken.DurationTime(),
	)
}

// NewFileToken generates and returns a new record private file access token.
func (m *Record) NewFileToken() (string, error) {
	if !m.Collection().IsAuth() {
		return "", ErrNotAuthRecord
	}

	key := (m.TokenKey() + m.Collection().FileToken.Secret)
	if key == "" {
		return "", ErrMissingSigningKey
	}

	return security.NewJWT(
		jwt.MapClaims{
			TokenClaimType:         TokenTypeFile,
			TokenClaimId:           m.Id,
			TokenClaimCollectionId: m.Collection().Id,
		},
		key,
		m.Collection().FileToken.DurationTime(),
	)
}
