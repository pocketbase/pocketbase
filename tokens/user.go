package tokens

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
)

// NewUserAuthToken generates and returns a new user authentication token.
func NewUserAuthToken(app core.App, user *models.User) (string, error) {
	return security.NewToken(
		jwt.MapClaims{"id": user.Id, "type": "user"},
		(user.TokenKey + app.Settings().UserAuthToken.Secret),
		app.Settings().UserAuthToken.Duration,
	)
}

// NewUserVerifyToken generates and returns a new user verification token.
func NewUserVerifyToken(app core.App, user *models.User) (string, error) {
	return security.NewToken(
		jwt.MapClaims{"id": user.Id, "type": "user", "email": user.Email},
		(user.TokenKey + app.Settings().UserVerificationToken.Secret),
		app.Settings().UserVerificationToken.Duration,
	)
}

// NewUserResetPasswordToken generates and returns a new user password reset request token.
func NewUserResetPasswordToken(app core.App, user *models.User) (string, error) {
	return security.NewToken(
		jwt.MapClaims{"id": user.Id, "type": "user", "email": user.Email},
		(user.TokenKey + app.Settings().UserPasswordResetToken.Secret),
		app.Settings().UserPasswordResetToken.Duration,
	)
}

// NewUserChangeEmailToken generates and returns a new user change email request token.
func NewUserChangeEmailToken(app core.App, user *models.User, newEmail string) (string, error) {
	return security.NewToken(
		jwt.MapClaims{"id": user.Id, "type": "user", "email": user.Email, "newEmail": newEmail},
		(user.TokenKey + app.Settings().UserEmailChangeToken.Secret),
		app.Settings().UserEmailChangeToken.Duration,
	)
}
