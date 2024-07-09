package tokens

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
)

// NewRecordAuthToken generates and returns a new auth record authentication token.
func NewRecordAuthToken(app core.App, record *models.Record) (string, error) {
	if !record.Collection().IsAuth() {
		return "", errors.New("the record is not from an auth collection")
	}

	return security.NewJWT(
		jwt.MapClaims{
			"id":           record.Id,
			"type":         TypeAuthRecord,
			"collectionId": record.Collection().Id,
		},
		(record.TokenKey() + app.Settings().RecordAuthToken.Secret),
		app.Settings().RecordAuthToken.Duration,
	)
}

// NewRecordVerifyToken generates and returns a new record verification token.
func NewRecordVerifyToken(app core.App, record *models.Record) (string, error) {
	if !record.Collection().IsAuth() {
		return "", errors.New("the record is not from an auth collection")
	}

	return security.NewJWT(
		jwt.MapClaims{
			"id":           record.Id,
			"type":         TypeAuthRecord,
			"collectionId": record.Collection().Id,
			"email":        record.Email(),
		},
		(record.TokenKey() + app.Settings().RecordVerificationToken.Secret),
		app.Settings().RecordVerificationToken.Duration,
	)
}

// NewRecordResetPasswordToken generates and returns a new auth record password reset request token.
func NewRecordResetPasswordToken(app core.App, record *models.Record) (string, error) {
	if !record.Collection().IsAuth() {
		return "", errors.New("the record is not from an auth collection")
	}

	return security.NewJWT(
		jwt.MapClaims{
			"id":           record.Id,
			"type":         TypeAuthRecord,
			"collectionId": record.Collection().Id,
			"email":        record.Email(),
		},
		(record.TokenKey() + app.Settings().RecordPasswordResetToken.Secret),
		app.Settings().RecordPasswordResetToken.Duration,
	)
}

// NewRecordChangeEmailToken generates and returns a new auth record change email request token.
func NewRecordChangeEmailToken(app core.App, record *models.Record, newEmail string) (string, error) {
	return security.NewJWT(
		jwt.MapClaims{
			"id":           record.Id,
			"type":         TypeAuthRecord,
			"collectionId": record.Collection().Id,
			"email":        record.Email(),
			"newEmail":     newEmail,
		},
		(record.TokenKey() + app.Settings().RecordEmailChangeToken.Secret),
		app.Settings().RecordEmailChangeToken.Duration,
	)
}

// NewRecordFileToken generates and returns a new record private file access token.
func NewRecordFileToken(app core.App, record *models.Record) (string, error) {
	if !record.Collection().IsAuth() {
		return "", errors.New("the record is not from an auth collection")
	}

	return security.NewJWT(
		jwt.MapClaims{
			"id":           record.Id,
			"type":         TypeAuthRecord,
			"collectionId": record.Collection().Id,
		},
		(record.TokenKey() + app.Settings().RecordFileToken.Secret),
		app.Settings().RecordFileToken.Duration,
	)
}
