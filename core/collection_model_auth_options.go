package core

import (
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

func (m *Collection) unsetMissingOAuth2MappedFields() {
	if !m.IsAuth() {
		return
	}

	if m.OAuth2.MappedFields.Id != "" {
		if m.Fields.GetByName(m.OAuth2.MappedFields.Id) == nil {
			m.OAuth2.MappedFields.Id = ""
		}
	}

	if m.OAuth2.MappedFields.Name != "" {
		if m.Fields.GetByName(m.OAuth2.MappedFields.Name) == nil {
			m.OAuth2.MappedFields.Name = ""
		}
	}

	if m.OAuth2.MappedFields.Username != "" {
		if m.Fields.GetByName(m.OAuth2.MappedFields.Username) == nil {
			m.OAuth2.MappedFields.Username = ""
		}
	}

	if m.OAuth2.MappedFields.AvatarURL != "" {
		if m.Fields.GetByName(m.OAuth2.MappedFields.AvatarURL) == nil {
			m.OAuth2.MappedFields.AvatarURL = ""
		}
	}
}

func (m *Collection) setDefaultAuthOptions() {
	m.collectionAuthOptions = collectionAuthOptions{
		VerificationTemplate:       defaultVerificationTemplate,
		ResetPasswordTemplate:      defaultResetPasswordTemplate,
		ConfirmEmailChangeTemplate: defaultConfirmEmailChangeTemplate,
		AuthRule:                   types.Pointer(""),
		AuthAlert: AuthAlertConfig{
			Enabled:       true,
			EmailTemplate: defaultAuthAlertTemplate,
		},
		PasswordAuth: PasswordAuthConfig{
			Enabled:        true,
			IdentityFields: []string{FieldNameEmail},
		},
		MFA: MFAConfig{
			Enabled:  false,
			Duration: 1800, // 30min
		},
		OTP: OTPConfig{
			Enabled:       false,
			Duration:      180, // 3min
			Length:        8,
			EmailTemplate: defaultOTPTemplate,
		},
		AuthToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 604800, // 7 days
		},
		PasswordResetToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1800, // 30min
		},
		EmailChangeToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1800, // 30min
		},
		VerificationToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 259200, // 3days
		},
		FileToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 180, // 3min
		},
	}
}

var _ optionsValidator = (*collectionAuthOptions)(nil)

// collectionAuthOptions defines the options for the "auth" type collection.
type collectionAuthOptions struct {
	// AuthRule could be used to specify additional record constraints
	// applied after record authentication and right before returning the
	// auth token response to the client.
	//
	// For example, to allow only verified users you could set it to
	// "verified = true".
	//
	// Set it to empty string to allow any Auth collection record to authenticate.
	//
	// Set it to nil to disallow authentication altogether for the collection
	// (that includes password, OAuth2, etc.).
	AuthRule *string `form:"authRule" json:"authRule"`

	// ManageRule gives admin-like permissions to allow fully managing
	// the auth record(s), eg. changing the password without requiring
	// to enter the old one, directly updating the verified state and email, etc.
	//
	// This rule is executed in addition to the Create and Update API rules.
	ManageRule *string `form:"manageRule" json:"manageRule"`

	// AuthAlert defines options related to the auth alerts on new device login.
	AuthAlert AuthAlertConfig `form:"authAlert" json:"authAlert"`

	// OAuth2 specifies whether OAuth2 auth is enabled for the collection
	// and which OAuth2 providers are allowed.
	OAuth2 OAuth2Config `form:"oauth2" json:"oauth2"`

	// PasswordAuth defines options related to the collection password authentication.
	PasswordAuth PasswordAuthConfig `form:"passwordAuth" json:"passwordAuth"`

	// MFA defines options related to the Multi-factor authentication (MFA).
	MFA MFAConfig `form:"mfa" json:"mfa"`

	// OTP defines options related to the One-time password authentication (OTP).
	OTP OTPConfig `form:"otp" json:"otp"`

	// Various token configurations
	// ---
	AuthToken          TokenConfig `form:"authToken" json:"authToken"`
	PasswordResetToken TokenConfig `form:"passwordResetToken" json:"passwordResetToken"`
	EmailChangeToken   TokenConfig `form:"emailChangeToken" json:"emailChangeToken"`
	VerificationToken  TokenConfig `form:"verificationToken" json:"verificationToken"`
	FileToken          TokenConfig `form:"fileToken" json:"fileToken"`

	// Default email templates
	// ---
	VerificationTemplate       EmailTemplate `form:"verificationTemplate" json:"verificationTemplate"`
	ResetPasswordTemplate      EmailTemplate `form:"resetPasswordTemplate" json:"resetPasswordTemplate"`
	ConfirmEmailChangeTemplate EmailTemplate `form:"confirmEmailChangeTemplate" json:"confirmEmailChangeTemplate"`
}

func (o *collectionAuthOptions) validate(cv *collectionValidator) error {
	err := validation.ValidateStruct(o,
		validation.Field(
			&o.AuthRule,
			validation.By(cv.checkRule),
			validation.By(cv.ensureNoSystemRuleChange(cv.original.AuthRule)),
		),
		validation.Field(
			&o.ManageRule,
			validation.NilOrNotEmpty,
			validation.By(cv.checkRule),
			validation.By(cv.ensureNoSystemRuleChange(cv.original.ManageRule)),
		),
		validation.Field(&o.AuthAlert),
		validation.Field(&o.PasswordAuth),
		validation.Field(&o.OAuth2),
		validation.Field(&o.OTP),
		validation.Field(&o.MFA),
		validation.Field(&o.AuthToken),
		validation.Field(&o.PasswordResetToken),
		validation.Field(&o.EmailChangeToken),
		validation.Field(&o.VerificationToken),
		validation.Field(&o.FileToken),
		validation.Field(&o.VerificationTemplate, validation.Required),
		validation.Field(&o.ResetPasswordTemplate, validation.Required),
		validation.Field(&o.ConfirmEmailChangeTemplate, validation.Required),
	)
	if err != nil {
		return err
	}

	if o.MFA.Enabled {
		// if MFA is enabled require at least 2 auth methods
		//
		// @todo maybe consider disabling the check because if custom auth methods
		// are registered it may fail since we don't have mechanism to detect them at the moment
		authsEnabled := 0
		if o.PasswordAuth.Enabled {
			authsEnabled++
		}
		if o.OAuth2.Enabled {
			authsEnabled++
		}
		if o.OTP.Enabled {
			authsEnabled++
		}
		if authsEnabled < 2 {
			return validation.Errors{
				"mfa": validation.Errors{
					"enabled": validation.NewError("validation_mfa_not_enough_auths", "MFA requires at least 2 auth methods to be enabled."),
				},
			}
		}

		if o.MFA.Rule != "" {
			mfaRuleValidators := []validation.RuleFunc{
				cv.checkRule,
				cv.ensureNoSystemRuleChange(&cv.original.MFA.Rule),
			}

			for _, validator := range mfaRuleValidators {
				err := validator(&o.MFA.Rule)
				if err != nil {
					return validation.Errors{
						"mfa": validation.Errors{
							"rule": err,
						},
					}
				}
			}
		}
	}

	// extra check to ensure that only unique identity fields are used
	if o.PasswordAuth.Enabled {
		err = validation.Validate(o.PasswordAuth.IdentityFields, validation.By(cv.checkFieldsForUniqueIndex))
		if err != nil {
			return validation.Errors{
				"passwordAuth": validation.Errors{
					"identityFields": err,
				},
			}
		}
	}

	return nil
}

// -------------------------------------------------------------------

type EmailTemplate struct {
	Subject string `form:"subject" json:"subject"`
	Body    string `form:"body" json:"body"`
}

// Validate makes EmailTemplate validatable by implementing [validation.Validatable] interface.
func (t EmailTemplate) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Subject, validation.Required),
		validation.Field(&t.Body, validation.Required),
	)
}

// Resolve replaces the placeholder parameters in the current email
// template and returns its components as ready-to-use strings.
func (t EmailTemplate) Resolve(placeholders map[string]any) (subject, body string) {
	body = t.Body
	subject = t.Subject

	for k, v := range placeholders {
		vStr := cast.ToString(v)

		// replace subject placeholder params (if any)
		subject = strings.ReplaceAll(subject, k, vStr)

		// replace body placeholder params (if any)
		body = strings.ReplaceAll(body, k, vStr)
	}

	return subject, body
}

// -------------------------------------------------------------------

type AuthAlertConfig struct {
	Enabled       bool          `form:"enabled" json:"enabled"`
	EmailTemplate EmailTemplate `form:"emailTemplate" json:"emailTemplate"`
}

// Validate makes AuthAlertConfig validatable by implementing [validation.Validatable] interface.
func (c AuthAlertConfig) Validate() error {
	return validation.ValidateStruct(&c,
		// note: for now always run the email template validations even
		// if not enabled since it could be used separately
		validation.Field(&c.EmailTemplate),
	)
}

// -------------------------------------------------------------------

type TokenConfig struct {
	Secret string `form:"secret" json:"secret,omitempty"`

	// Duration specifies how long an issued token to be valid (in seconds)
	Duration int64 `form:"duration" json:"duration"`
}

// Validate makes TokenConfig validatable by implementing [validation.Validatable] interface.
func (c TokenConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Secret, validation.Required, validation.Length(30, 255)),
		validation.Field(&c.Duration, validation.Required, validation.Min(10), validation.Max(94670856)), // ~3y max
	)
}

// DurationTime returns the current Duration as [time.Duration].
func (c TokenConfig) DurationTime() time.Duration {
	return time.Duration(c.Duration) * time.Second
}

// -------------------------------------------------------------------

type OTPConfig struct {
	Enabled bool `form:"enabled" json:"enabled"`

	// Duration specifies how long the OTP to be valid (in seconds)
	Duration int64 `form:"duration" json:"duration"`

	// Length specifies the auto generated password length.
	Length int `form:"length" json:"length"`

	// EmailTemplate is the default OTP email template that will be send to the auth record.
	//
	// In addition to the system placeholders you can also make use of
	// [core.EmailPlaceholderOTPId] and [core.EmailPlaceholderOTP].
	EmailTemplate EmailTemplate `form:"emailTemplate" json:"emailTemplate"`
}

// Validate makes OTPConfig validatable by implementing [validation.Validatable] interface.
func (c OTPConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Duration, validation.When(c.Enabled, validation.Required, validation.Min(10), validation.Max(86400))),
		validation.Field(&c.Length, validation.When(c.Enabled, validation.Required, validation.Min(4))),
		// note: for now always run the email template validations even
		// if not enabled since it could be used separately
		validation.Field(&c.EmailTemplate),
	)
}

// DurationTime returns the current Duration as [time.Duration].
func (c OTPConfig) DurationTime() time.Duration {
	return time.Duration(c.Duration) * time.Second
}

// -------------------------------------------------------------------

type MFAConfig struct {
	Enabled bool `form:"enabled" json:"enabled"`

	// Duration specifies how long an issued MFA to be valid (in seconds)
	Duration int64 `form:"duration" json:"duration"`

	// Rule is an optional field to restrict MFA only for the records that satisfy the rule.
	//
	// Leave it empty to enable MFA for everyone.
	Rule string `form:"rule" json:"rule"`
}

// Validate makes MFAConfig validatable by implementing [validation.Validatable] interface.
func (c MFAConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Duration, validation.When(c.Enabled, validation.Required, validation.Min(10), validation.Max(86400))),
	)
}

// DurationTime returns the current Duration as [time.Duration].
func (c MFAConfig) DurationTime() time.Duration {
	return time.Duration(c.Duration) * time.Second
}

// -------------------------------------------------------------------

type PasswordAuthConfig struct {
	Enabled bool `form:"enabled" json:"enabled"`

	// IdentityFields is a list of field names that could be used as
	// identity during password authentication.
	//
	// Usually only fields that has single column UNIQUE index are accepted as values.
	IdentityFields []string `form:"identityFields" json:"identityFields"`
}

// Validate makes PasswordAuthConfig validatable by implementing [validation.Validatable] interface.
func (c PasswordAuthConfig) Validate() error {
	// strip duplicated values
	c.IdentityFields = list.ToUniqueStringSlice(c.IdentityFields)

	if !c.Enabled {
		return nil // no need to validate
	}

	return validation.ValidateStruct(&c,
		validation.Field(&c.IdentityFields, validation.Required),
	)
}

// -------------------------------------------------------------------

type OAuth2KnownFields struct {
	Id        string `form:"id" json:"id"`
	Name      string `form:"name" json:"name"`
	Username  string `form:"username" json:"username"`
	AvatarURL string `form:"avatarURL" json:"avatarURL"`
}

type OAuth2Config struct {
	Providers []OAuth2ProviderConfig `form:"providers" json:"providers"`

	MappedFields OAuth2KnownFields `form:"mappedFields" json:"mappedFields"`

	Enabled bool `form:"enabled" json:"enabled"`
}

// GetProviderConfig returns the first OAuth2ProviderConfig that matches the specified name.
//
// Returns false and zero config if no such provider is available in c.Providers.
func (c OAuth2Config) GetProviderConfig(name string) (config OAuth2ProviderConfig, exists bool) {
	for _, p := range c.Providers {
		if p.Name == name {
			return p, true
		}
	}
	return
}

// Validate makes OAuth2Config validatable by implementing [validation.Validatable] interface.
func (c OAuth2Config) Validate() error {
	if !c.Enabled {
		return nil // no need to validate
	}

	return validation.ValidateStruct(&c,
		// note: don't require providers for now as they could be externally registered/removed
		validation.Field(&c.Providers, validation.By(checkForDuplicatedProviders)),
	)
}

func checkForDuplicatedProviders(value any) error {
	configs, _ := value.([]OAuth2ProviderConfig)

	existing := map[string]struct{}{}

	for i, c := range configs {
		if c.Name == "" {
			continue // the name nonempty state is validated separately
		}
		if _, ok := existing[c.Name]; ok {
			return validation.Errors{
				strconv.Itoa(i): validation.Errors{
					"name": validation.NewError("validation_duplicated_provider", "The provider {{.name}} is already registered.").
						SetParams(map[string]any{"name": c.Name}),
				},
			}
		}
		existing[c.Name] = struct{}{}
	}

	return nil
}

type OAuth2ProviderConfig struct {
	// PKCE overwrites the default provider PKCE config option.
	//
	// This usually shouldn't be needed but some OAuth2 vendors, like the LinkedIn OIDC,
	// may require manual adjustment due to returning error if extra parameters are added to the request
	// (https://github.com/pocketbase/pocketbase/discussions/3799#discussioncomment-7640312)
	PKCE *bool `form:"pkce" json:"pkce"`

	Name         string         `form:"name" json:"name"`
	ClientId     string         `form:"clientId" json:"clientId"`
	ClientSecret string         `form:"clientSecret" json:"clientSecret,omitempty"`
	AuthURL      string         `form:"authURL" json:"authURL"`
	TokenURL     string         `form:"tokenURL" json:"tokenURL"`
	UserInfoURL  string         `form:"userInfoURL" json:"userInfoURL"`
	DisplayName  string         `form:"displayName" json:"displayName"`
	Extra        map[string]any `form:"extra" json:"extra"`
}

// Validate makes OAuth2ProviderConfig validatable by implementing [validation.Validatable] interface.
func (c OAuth2ProviderConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.By(checkProviderName)),
		validation.Field(&c.ClientId, validation.Required),
		validation.Field(&c.ClientSecret, validation.Required),
		validation.Field(&c.AuthURL, is.URL),
		validation.Field(&c.TokenURL, is.URL),
		validation.Field(&c.UserInfoURL, is.URL),
	)
}

func checkProviderName(value any) error {
	name, _ := value.(string)
	if name == "" {
		return nil // nothing to check
	}

	if _, err := auth.NewProviderByName(name); err != nil {
		return validation.NewError("validation_missing_provider", "Invalid or missing provider with name {{.name}}.").
			SetParams(map[string]any{"name": name})
	}

	return nil
}

// InitProvider returns a new auth.Provider instance loaded with the current OAuth2ProviderConfig options.
func (c OAuth2ProviderConfig) InitProvider() (auth.Provider, error) {
	provider, err := auth.NewProviderByName(c.Name)
	if err != nil {
		return nil, err
	}

	if c.ClientId != "" {
		provider.SetClientId(c.ClientId)
	}

	if c.ClientSecret != "" {
		provider.SetClientSecret(c.ClientSecret)
	}

	if c.AuthURL != "" {
		provider.SetAuthURL(c.AuthURL)
	}

	if c.UserInfoURL != "" {
		provider.SetUserInfoURL(c.UserInfoURL)
	}

	if c.TokenURL != "" {
		provider.SetTokenURL(c.TokenURL)
	}

	if c.DisplayName != "" {
		provider.SetDisplayName(c.DisplayName)
	}

	if c.PKCE != nil {
		provider.SetPKCE(*c.PKCE)
	}

	if c.Extra != nil {
		provider.SetExtra(c.Extra)
	}

	return provider, nil
}
