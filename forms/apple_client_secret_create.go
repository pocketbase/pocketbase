package forms

import (
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/pocketbase/core"
)

var privateKeyRegex = regexp.MustCompile(`(?m)-----BEGIN PRIVATE KEY----[\s\S]+-----END PRIVATE KEY-----`)

// AppleClientSecretCreate is a form struct to generate a new Apple Client Secret.
//
// Reference: https://developer.apple.com/documentation/sign_in_with_apple/generate_and_validate_tokens
type AppleClientSecretCreate struct {
	app core.App

	// ClientId is the identifier of your app (aka. Service ID).
	ClientId string `form:"clientId" json:"clientId"`

	// TeamId is a 10-character string associated with your developer account
	// (usually could be found next to your name in the Apple Developer site).
	TeamId string `form:"teamId" json:"teamId"`

	// KeyId is a 10-character key identifier generated for the "Sign in with Apple"
	// private key associated with your developer account.
	KeyId string `form:"keyId" json:"keyId"`

	// PrivateKey is the private key associated to your app.
	// Usually wrapped within -----BEGIN PRIVATE KEY----- X -----END PRIVATE KEY-----.
	PrivateKey string `form:"privateKey" json:"privateKey"`

	// Duration specifies how long the generated JWT should be considered valid.
	// The specified value must be in seconds and max 15777000 (~6months).
	Duration int `form:"duration" json:"duration"`
}

// NewAppleClientSecretCreate creates a new [AppleClientSecretCreate] form with initializer
// config created from the provided [core.App] instances.
func NewAppleClientSecretCreate(app core.App) *AppleClientSecretCreate {
	form := &AppleClientSecretCreate{
		app: app,
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *AppleClientSecretCreate) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.ClientId, validation.Required),
		validation.Field(&form.TeamId, validation.Required, validation.Length(10, 10)),
		validation.Field(&form.KeyId, validation.Required, validation.Length(10, 10)),
		validation.Field(&form.PrivateKey, validation.Required, validation.Match(privateKeyRegex)),
		validation.Field(&form.Duration, validation.Required, validation.Min(1), validation.Max(15777000)),
	)
}

// Submit validates the form and returns a new Apple Client Secret JWT.
func (form *AppleClientSecretCreate) Submit() (string, error) {
	if err := form.Validate(); err != nil {
		return "", err
	}

	signKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(strings.TrimSpace(form.PrivateKey)))
	if err != nil {
		return "", err
	}

	now := time.Now()

	claims := &jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"https://appleid.apple.com"},
		Subject:   form.ClientId,
		Issuer:    form.TeamId,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(form.Duration) * time.Second)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = form.KeyId

	return token.SignedString(signKey)
}
