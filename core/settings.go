package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/rest"
	"github.com/pocketbase/pocketbase/tools/security"
)

// Settings defines common app configuration options.
type Settings struct {
	mux sync.RWMutex

	Meta MetaConfig `form:"meta" json:"meta"`
	Logs LogsConfig `form:"logs" json:"logs"`
	Smtp SmtpConfig `form:"smtp" json:"smtp"`
	S3   S3Config   `form:"s3" json:"s3"`

	AdminAuthToken          TokenConfig `form:"adminAuthToken" json:"adminAuthToken"`
	AdminPasswordResetToken TokenConfig `form:"adminPasswordResetToken" json:"adminPasswordResetToken"`
	UserAuthToken           TokenConfig `form:"userAuthToken" json:"userAuthToken"`
	UserPasswordResetToken  TokenConfig `form:"userPasswordResetToken" json:"userPasswordResetToken"`
	UserEmailChangeToken    TokenConfig `form:"userEmailChangeToken" json:"userEmailChangeToken"`
	UserVerificationToken   TokenConfig `form:"userVerificationToken" json:"userVerificationToken"`

	EmailAuth    EmailAuthConfig    `form:"emailAuth" json:"emailAuth"`
	GoogleAuth   AuthProviderConfig `form:"googleAuth" json:"googleAuth"`
	FacebookAuth AuthProviderConfig `form:"facebookAuth" json:"facebookAuth"`
	GithubAuth   AuthProviderConfig `form:"githubAuth" json:"githubAuth"`
	GitlabAuth   AuthProviderConfig `form:"gitlabAuth" json:"gitlabAuth"`
	DiscordAuth  AuthProviderConfig `form:"discordAuth" json:"discordAuth"`
	TwitterAuth  AuthProviderConfig `form:"twitterAuth" json:"twitterAuth"`
}

// NewSettings creates and returns a new default Settings instance.
func NewSettings() *Settings {
	return &Settings{
		Meta: MetaConfig{
			AppName:                    "Acme",
			AppUrl:                     "http://localhost:8090",
			HideControls:               false,
			SenderName:                 "Support",
			SenderAddress:              "support@example.com",
			VerificationTemplate:       defaultVerificationTemplate,
			ResetPasswordTemplate:      defaultResetPasswordTemplate,
			ConfirmEmailChangeTemplate: defaultConfirmEmailChangeTemplate,
		},

		Logs: LogsConfig{
			MaxDays: 7,
		},
		Smtp: SmtpConfig{
			Enabled:  false,
			Host:     "smtp.example.com",
			Port:     587,
			Username: "",
			Password: "",
			Tls:      false,
		},
		AdminAuthToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1209600, // 14 days,
		},
		AdminPasswordResetToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1800, // 30 minutes,
		},
		UserAuthToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1209600, // 14 days,
		},
		UserPasswordResetToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1800, // 30 minutes,
		},
		UserVerificationToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 604800, // 7 days,
		},
		UserEmailChangeToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1800, // 30 minutes,
		},
		EmailAuth: EmailAuthConfig{
			Enabled:           true,
			MinPasswordLength: 8,
		},
		GoogleAuth: AuthProviderConfig{
			Enabled:            false,
			AllowRegistrations: true,
		},
		FacebookAuth: AuthProviderConfig{
			Enabled:            false,
			AllowRegistrations: true,
		},
		GithubAuth: AuthProviderConfig{
			Enabled:            false,
			AllowRegistrations: true,
		},
		GitlabAuth: AuthProviderConfig{
			Enabled:            false,
			AllowRegistrations: true,
		},
		DiscordAuth: AuthProviderConfig{
			Enabled:            false,
			AllowRegistrations: true,
		},
		TwitterAuth: AuthProviderConfig{
			Enabled:            false,
			AllowRegistrations: true,
		},
	}
}

// Validate makes Settings validatable by implementing [validation.Validatable] interface.
func (s *Settings) Validate() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	return validation.ValidateStruct(s,
		validation.Field(&s.Meta),
		validation.Field(&s.Logs),
		validation.Field(&s.AdminAuthToken),
		validation.Field(&s.AdminPasswordResetToken),
		validation.Field(&s.UserAuthToken),
		validation.Field(&s.UserPasswordResetToken),
		validation.Field(&s.UserEmailChangeToken),
		validation.Field(&s.UserVerificationToken),
		validation.Field(&s.Smtp),
		validation.Field(&s.S3),
		validation.Field(&s.EmailAuth),
		validation.Field(&s.GoogleAuth),
		validation.Field(&s.FacebookAuth),
		validation.Field(&s.GithubAuth),
		validation.Field(&s.GitlabAuth),
		validation.Field(&s.DiscordAuth),
		validation.Field(&s.TwitterAuth),
	)
}

// Merge merges `other` settings into the current one.
func (s *Settings) Merge(other *Settings) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	bytes, err := json.Marshal(other)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

// Clone creates a new deep copy of the current settings.
func (s *Settings) Clone() (*Settings, error) {
	settings := &Settings{}
	if err := settings.Merge(s); err != nil {
		return nil, err
	}
	return settings, nil
}

// RedactClone creates a new deep copy of the current settings,
// while replacing the secret values with `******`.
func (s *Settings) RedactClone() (*Settings, error) {
	clone, err := s.Clone()
	if err != nil {
		return nil, err
	}

	mask := "******"

	sensitiveFields := []*string{
		&clone.Smtp.Password,
		&clone.S3.Secret,
		&clone.AdminAuthToken.Secret,
		&clone.AdminPasswordResetToken.Secret,
		&clone.UserAuthToken.Secret,
		&clone.UserPasswordResetToken.Secret,
		&clone.UserEmailChangeToken.Secret,
		&clone.UserVerificationToken.Secret,
		&clone.GoogleAuth.ClientSecret,
		&clone.FacebookAuth.ClientSecret,
		&clone.GithubAuth.ClientSecret,
		&clone.GitlabAuth.ClientSecret,
		&clone.DiscordAuth.ClientSecret,
		&clone.TwitterAuth.ClientSecret,
	}

	// mask all sensitive fields
	for _, v := range sensitiveFields {
		if v != nil && *v != "" {
			*v = mask
		}
	}

	return clone, nil
}

// NamedAuthProviderConfigs returns a map with all registered OAuth2
// provider configurations (indexed by their name identifier).
func (s *Settings) NamedAuthProviderConfigs() map[string]AuthProviderConfig {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return map[string]AuthProviderConfig{
		auth.NameGoogle:   s.GoogleAuth,
		auth.NameFacebook: s.FacebookAuth,
		auth.NameGithub:   s.GithubAuth,
		auth.NameGitlab:   s.GitlabAuth,
		auth.NameDiscord:  s.DiscordAuth,
		auth.NameTwitter:  s.TwitterAuth,
	}
}

// -------------------------------------------------------------------

type TokenConfig struct {
	Secret   string `form:"secret" json:"secret"`
	Duration int64  `form:"duration" json:"duration"`
}

// Validate makes TokenConfig validatable by implementing [validation.Validatable] interface.
func (c TokenConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Secret, validation.Required, validation.Length(30, 300)),
		validation.Field(&c.Duration, validation.Required, validation.Min(5), validation.Max(63072000)),
	)
}

// -------------------------------------------------------------------

type SmtpConfig struct {
	Enabled  bool   `form:"enabled" json:"enabled"`
	Host     string `form:"host" json:"host"`
	Port     int    `form:"port" json:"port"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`

	// Whether to enforce TLS encryption for the mail server connection.
	//
	// When set to false StartTLS command is send, leaving the server
	// to decide whether to upgrade the connection or not.
	Tls bool `form:"tls" json:"tls"`
}

// Validate makes SmtpConfig validatable by implementing [validation.Validatable] interface.
func (c SmtpConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, is.Host, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.Port, validation.When(c.Enabled, validation.Required), validation.Min(0)),
	)
}

// -------------------------------------------------------------------

type S3Config struct {
	Enabled        bool   `form:"enabled" json:"enabled"`
	Bucket         string `form:"bucket" json:"bucket"`
	Region         string `form:"region" json:"region"`
	Endpoint       string `form:"endpoint" json:"endpoint"`
	AccessKey      string `form:"accessKey" json:"accessKey"`
	Secret         string `form:"secret" json:"secret"`
	ForcePathStyle bool   `form:"forcePathStyle" json:"forcePathStyle"`
}

// Validate makes S3Config validatable by implementing [validation.Validatable] interface.
func (c S3Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Endpoint, is.URL, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.Bucket, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.Region, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.AccessKey, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.Secret, validation.When(c.Enabled, validation.Required)),
	)
}

// -------------------------------------------------------------------

type MetaConfig struct {
	AppName                    string        `form:"appName" json:"appName"`
	AppUrl                     string        `form:"appUrl" json:"appUrl"`
	HideControls               bool          `form:"hideControls" json:"hideControls"`
	SenderName                 string        `form:"senderName" json:"senderName"`
	SenderAddress              string        `form:"senderAddress" json:"senderAddress"`
	VerificationTemplate       EmailTemplate `form:"verificationTemplate" json:"verificationTemplate"`
	ResetPasswordTemplate      EmailTemplate `form:"resetPasswordTemplate" json:"resetPasswordTemplate"`
	ConfirmEmailChangeTemplate EmailTemplate `form:"confirmEmailChangeTemplate" json:"confirmEmailChangeTemplate"`
}

// Validate makes MetaConfig validatable by implementing [validation.Validatable] interface.
func (c MetaConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.AppName, validation.Required, validation.Length(1, 255)),
		validation.Field(&c.AppUrl, validation.Required, is.URL),
		validation.Field(&c.SenderName, validation.Required, validation.Length(1, 255)),
		validation.Field(&c.SenderAddress, is.EmailFormat, validation.Required),
		validation.Field(&c.VerificationTemplate, validation.Required),
		validation.Field(&c.ResetPasswordTemplate, validation.Required),
		validation.Field(&c.ConfirmEmailChangeTemplate, validation.Required),
	)
}

type EmailTemplate struct {
	Body      string `form:"body" json:"body"`
	Subject   string `form:"subject" json:"subject"`
	ActionUrl string `form:"actionUrl" json:"actionUrl"`
}

// Validate makes EmailTemplate validatable by implementing [validation.Validatable] interface.
func (t EmailTemplate) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Subject, validation.Required),
		validation.Field(
			&t.Body,
			validation.Required,
			validation.By(checkPlaceholderParams(EmailPlaceholderActionUrl)),
		),
		validation.Field(
			&t.ActionUrl,
			validation.Required,
			validation.By(checkPlaceholderParams(EmailPlaceholderToken)),
		),
	)
}

func checkPlaceholderParams(params ...string) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(string)

		for _, param := range params {
			if !strings.Contains(v, param) {
				return validation.NewError(
					"validation_missing_required_param",
					fmt.Sprintf("Missing required parameter %q", param),
				)
			}
		}

		return nil
	}
}

// Resolve replaces the placeholder parameters in the current email
// template and returns its components as ready-to-use strings.
func (t EmailTemplate) Resolve(
	appName string,
	appUrl,
	token string,
) (subject, body, actionUrl string) {
	// replace action url placeholder params (if any)
	actionUrlParams := map[string]string{
		EmailPlaceholderAppName: appName,
		EmailPlaceholderAppUrl:  appUrl,
		EmailPlaceholderToken:   token,
	}
	actionUrl = t.ActionUrl
	for k, v := range actionUrlParams {
		actionUrl = strings.ReplaceAll(actionUrl, k, v)
	}
	actionUrl, _ = rest.NormalizeUrl(actionUrl)

	// replace body placeholder params (if any)
	bodyParams := map[string]string{
		EmailPlaceholderAppName:   appName,
		EmailPlaceholderAppUrl:    appUrl,
		EmailPlaceholderToken:     token,
		EmailPlaceholderActionUrl: actionUrl,
	}
	body = t.Body
	for k, v := range bodyParams {
		body = strings.ReplaceAll(body, k, v)
	}

	// replace subject placeholder params (if any)
	subjectParams := map[string]string{
		EmailPlaceholderAppName: appName,
		EmailPlaceholderAppUrl:  appUrl,
	}
	subject = t.Subject
	for k, v := range subjectParams {
		subject = strings.ReplaceAll(subject, k, v)
	}

	return subject, body, actionUrl
}

// -------------------------------------------------------------------

type LogsConfig struct {
	MaxDays int `form:"maxDays" json:"maxDays"`
}

// Validate makes LogsConfig validatable by implementing [validation.Validatable] interface.
func (c LogsConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.MaxDays, validation.Min(0)),
	)
}

// -------------------------------------------------------------------

type EmailAuthConfig struct {
	Enabled           bool     `form:"enabled" json:"enabled"`
	ExceptDomains     []string `form:"exceptDomains" json:"exceptDomains"`
	OnlyDomains       []string `form:"onlyDomains" json:"onlyDomains"`
	MinPasswordLength int      `form:"minPasswordLength" json:"minPasswordLength"`
}

// Validate makes `EmailAuthConfig` validatable by implementing [validation.Validatable] interface.
func (c EmailAuthConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.ExceptDomains,
			validation.When(len(c.OnlyDomains) > 0, validation.Empty).Else(validation.Each(is.Domain)),
		),
		validation.Field(
			&c.OnlyDomains,
			validation.When(len(c.ExceptDomains) > 0, validation.Empty).Else(validation.Each(is.Domain)),
		),
		validation.Field(
			&c.MinPasswordLength,
			validation.When(c.Enabled, validation.Required),
			validation.Min(5),
			validation.Max(100),
		),
	)
}

// -------------------------------------------------------------------

type AuthProviderConfig struct {
	Enabled            bool   `form:"enabled" json:"enabled"`
	AllowRegistrations bool   `form:"allowRegistrations" json:"allowRegistrations"`
	ClientId           string `form:"clientId" json:"clientId,omitempty"`
	ClientSecret       string `form:"clientSecret" json:"clientSecret,omitempty"`
	AuthUrl            string `form:"authUrl" json:"authUrl,omitempty"`
	TokenUrl           string `form:"tokenUrl" json:"tokenUrl,omitempty"`
	UserApiUrl         string `form:"userApiUrl" json:"userApiUrl,omitempty"`
}

// Validate makes `ProviderConfig` validatable by implementing [validation.Validatable] interface.
func (c AuthProviderConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClientId, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.ClientSecret, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.AuthUrl, is.URL),
		validation.Field(&c.TokenUrl, is.URL),
		validation.Field(&c.UserApiUrl, is.URL),
	)
}

// SetupProvider loads the current AuthProviderConfig into the specified provider.
func (c AuthProviderConfig) SetupProvider(provider auth.Provider) error {
	if !c.Enabled {
		return errors.New("The provider is not enabled.")
	}

	if c.ClientId != "" {
		provider.SetClientId(c.ClientId)
	}

	if c.ClientSecret != "" {
		provider.SetClientSecret(c.ClientSecret)
	}

	if c.AuthUrl != "" {
		provider.SetAuthUrl(c.AuthUrl)
	}

	if c.UserApiUrl != "" {
		provider.SetUserApiUrl(c.UserApiUrl)
	}

	if c.TokenUrl != "" {
		provider.SetTokenUrl(c.TokenUrl)
	}

	return nil
}
