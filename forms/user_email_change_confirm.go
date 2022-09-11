package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
)

// UserEmailChangeConfirm specifies a user email change confirmation form.
type UserEmailChangeConfirm struct {
	config UserEmailChangeConfirmConfig

	Token    string `form:"token" json:"token"`
	Password string `form:"password" json:"password"`
}

// UserEmailChangeConfirmConfig is the [UserEmailChangeConfirm] factory initializer config.
//
// NB! App is required struct member.
type UserEmailChangeConfirmConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewUserEmailChangeConfirm creates a new [UserEmailChangeConfirm]
// form with initializer config created from the provided [core.App] instance.
//
// This factory method is used primarily for convenience (and backward compatibility).
// If you want to submit the form as part of another transaction, use
// [NewUserEmailChangeConfirmWithConfig] with explicitly set Dao.
func NewUserEmailChangeConfirm(app core.App) *UserEmailChangeConfirm {
	return NewUserEmailChangeConfirmWithConfig(UserEmailChangeConfirmConfig{
		App: app,
	})
}

// NewUserEmailChangeConfirmWithConfig creates a new [UserEmailChangeConfirm]
// form with the provided config or panics on invalid configuration.
func NewUserEmailChangeConfirmWithConfig(config UserEmailChangeConfirmConfig) *UserEmailChangeConfirm {
	form := &UserEmailChangeConfirm{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserEmailChangeConfirm) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Token,
			validation.Required,
			validation.By(form.checkToken),
		),
		validation.Field(
			&form.Password,
			validation.Required,
			validation.Length(1, 100),
			validation.By(form.checkPassword),
		),
	)
}

func (form *UserEmailChangeConfirm) checkToken(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	_, _, err := form.parseToken(v)

	return err
}

func (form *UserEmailChangeConfirm) checkPassword(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	user, _, _ := form.parseToken(form.Token)
	if user == nil || !user.ValidatePassword(v) {
		return validation.NewError("validation_invalid_password", "Missing or invalid user password.")
	}

	return nil
}

func (form *UserEmailChangeConfirm) parseToken(token string) (*models.User, string, error) {
	// check token payload
	claims, _ := security.ParseUnverifiedJWT(token)
	newEmail, _ := claims["newEmail"].(string)
	if newEmail == "" {
		return nil, "", validation.NewError("validation_invalid_token_payload", "Invalid token payload - newEmail must be set.")
	}

	// ensure that there aren't other users with the new email
	if !form.config.Dao.IsUserEmailUnique(newEmail, "") {
		return nil, "", validation.NewError("validation_existing_token_email", "The new email address is already registered: "+newEmail)
	}

	// verify that the token is not expired and its signature is valid
	user, err := form.config.Dao.FindUserByToken(
		token,
		form.config.App.Settings().UserEmailChangeToken.Secret,
	)
	if err != nil || user == nil {
		return nil, "", validation.NewError("validation_invalid_token", "Invalid or expired token.")
	}

	return user, newEmail, nil
}

// Submit validates and submits the user email change confirmation form.
// On success returns the updated user model associated to `form.Token`.
func (form *UserEmailChangeConfirm) Submit() (*models.User, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	user, newEmail, err := form.parseToken(form.Token)
	if err != nil {
		return nil, err
	}

	user.Email = newEmail
	user.Verified = true
	user.RefreshTokenKey() // invalidate old tokens

	if err := form.config.Dao.SaveUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
