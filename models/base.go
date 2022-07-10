// Package models implements all PocketBase DB models.
package models

import (
	"errors"

	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/crypto/bcrypt"
)

// ColumnValueMapper defines an interface for custom db model data serialization.
type ColumnValueMapper interface {
	// ColumnValueMap returns the data to be used when persisting the model.
	ColumnValueMap() map[string]any
}

// Model defines an interface with common methods that all db models should have.
type Model interface {
	TableName() string
	HasId() bool
	GetId() string
	GetCreated() types.DateTime
	GetUpdated() types.DateTime
	RefreshId()
	RefreshCreated()
	RefreshUpdated()
}

// -------------------------------------------------------------------
// BaseModel
// -------------------------------------------------------------------

// BaseModel defines common fields and methods used by all other models.
type BaseModel struct {
	Id      string         `db:"id" json:"id"`
	Created types.DateTime `db:"created" json:"created"`
	Updated types.DateTime `db:"updated" json:"updated"`
}

// HasId returns whether the model has a nonzero primary key (aka. id).
func (m *BaseModel) HasId() bool {
	return m.GetId() != ""
}

// GetId returns the model's id.
func (m *BaseModel) GetId() string {
	return m.Id
}

// GetCreated returns the model's Created datetime.
func (m *BaseModel) GetCreated() types.DateTime {
	return m.Created
}

// GetUpdated returns the model's Updated datetime.
func (m *BaseModel) GetUpdated() types.DateTime {
	return m.Updated
}

// RefreshId generates and sets a new model id.
//
// The generated id is a cryptographically random 15 characters length string
// (could change in the future).
func (m *BaseModel) RefreshId() {
	m.Id = security.RandomString(15)
}

// RefreshCreated updates the model's Created field with the current datetime.
func (m *BaseModel) RefreshCreated() {
	m.Created = types.NowDateTime()
}

// RefreshUpdated updates the model's Created field with the current datetime.
func (m *BaseModel) RefreshUpdated() {
	m.Updated = types.NowDateTime()
}

// -------------------------------------------------------------------
// BaseAccount
// -------------------------------------------------------------------

// BaseAccount defines common fields and methods used by auth models (aka. users and admins).
type BaseAccount struct {
	BaseModel

	Email           string         `db:"email" json:"email"`
	TokenKey        string         `db:"tokenKey" json:"-"`
	PasswordHash    string         `db:"passwordHash" json:"-"`
	LastResetSentAt types.DateTime `db:"lastResetSentAt" json:"lastResetSentAt"`
}

// ValidatePassword validates a plain password against the model's password.
func (m *BaseAccount) ValidatePassword(password string) bool {
	bytePassword := []byte(password)
	bytePasswordHash := []byte(m.PasswordHash)

	// comparing the password with the hash
	err := bcrypt.CompareHashAndPassword(bytePasswordHash, bytePassword)

	// nil means it is a match
	return err == nil
}

// SetPassword sets cryptographically secure string to `model.Password`.
//
// Additionally this method also resets the LastResetSentAt and the TokenKey fields.
func (m *BaseAccount) SetPassword(password string) error {
	if password == "" {
		return errors.New("The provided plain password is empty")
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return err
	}

	m.PasswordHash = string(hashedPassword)
	m.LastResetSentAt = types.DateTime{} // reset

	// invalidate previously issued tokens
	m.RefreshTokenKey()

	return nil
}

// RefreshTokenKey generates and sets new random token key.
func (m *BaseAccount) RefreshTokenKey() {
	m.TokenKey = security.RandomString(50)
}
