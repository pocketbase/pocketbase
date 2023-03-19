package models

import (
	"errors"

	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"golang.org/x/crypto/bcrypt"
)

var (
	_ Model = (*Admin)(nil)
)

type Admin struct {
	BaseModel

	Avatar          int            `db:"avatar" json:"avatar"`
	Email           string         `db:"email" json:"email"`
	TokenKey        string         `db:"tokenKey" json:"-"`
	PasswordHash    string         `db:"passwordHash" json:"-"`
	LastResetSentAt types.DateTime `db:"lastResetSentAt" json:"-"`
}

// TableName returns the Admin model SQL table name.
func (m *Admin) TableName() string {
	return "_admins"
}

// ValidatePassword validates a plain password against the model's password.
func (m *Admin) ValidatePassword(password string) bool {
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
func (m *Admin) SetPassword(password string) error {
	if password == "" {
		return errors.New("The provided plain password is empty")
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	m.PasswordHash = string(hashedPassword)
	m.LastResetSentAt = types.DateTime{} // reset

	// invalidate previously issued tokens
	return m.RefreshTokenKey()
}

// RefreshTokenKey generates and sets new random token key.
func (m *Admin) RefreshTokenKey() error {
	m.TokenKey = security.RandomString(50)
	return nil
}
