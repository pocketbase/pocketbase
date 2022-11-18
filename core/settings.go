package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

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

	AdminAuthToken           TokenConfig `form:"adminAuthToken" json:"adminAuthToken"`
	AdminPasswordResetToken  TokenConfig `form:"adminPasswordResetToken" json:"adminPasswordResetToken"`
	RecordAuthToken          TokenConfig `form:"recordAuthToken" json:"recordAuthToken"`
	RecordPasswordResetToken TokenConfig `form:"recordPasswordResetToken" json:"recordPasswordResetToken"`
	RecordEmailChangeToken   TokenConfig `form:"recordEmailChangeToken" json:"recordEmailChangeToken"`
	RecordVerificationToken  TokenConfig `form:"recordVerificationToken" json:"recordVerificationToken"`

	// Deprecated: Will be removed in v0.9!
	EmailAuth EmailAuthConfig `form:"emailAuth" json:"emailAuth"`

	AppleAuth     AppleAuthProviderConfig `form:"appleAuth" json:"appleAuth"`
	GoogleAuth    BaseAuthProviderConfig  `form:"googleAuth" json:"googleAuth"`
	FacebookAuth  BaseAuthProviderConfig  `form:"facebookAuth" json:"facebookAuth"`
	GithubAuth    BaseAuthProviderConfig  `form:"githubAuth" json:"githubAuth"`
	GitlabAuth    BaseAuthProviderConfig  `form:"gitlabAuth" json:"gitlabAuth"`
	DiscordAuth   BaseAuthProviderConfig  `form:"discordAuth" json:"discordAuth"`
	TwitterAuth   BaseAuthProviderConfig  `form:"twitterAuth" json:"twitterAuth"`
	MicrosoftAuth BaseAuthProviderConfig  `form:"microsoftAuth" json:"microsoftAuth"`
	SpotifyAuth   BaseAuthProviderConfig  `form:"spotifyAuth" json:"spotifyAuth"`
	KakaoAuth     BaseAuthProviderConfig  `form:"kakaoAuth" json:"kakaoAuth"`
	TwitchAuth    BaseAuthProviderConfig  `form:"twitchAuth" json:"twitchAuth"`
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
			MaxDays: 5,
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
		RecordAuthToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1209600, // 14 days,
		},
		RecordPasswordResetToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1800, // 30 minutes,
		},
		RecordVerificationToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 604800, // 7 days,
		},
		RecordEmailChangeToken: TokenConfig{
			Secret:   security.RandomString(50),
			Duration: 1800, // 30 minutes,
		},
		AppleAuth: AppleAuthProviderConfig{
			BaseAuthProviderConfig: BaseAuthProviderConfig{
				Enabled: false,
			},
		},
		GoogleAuth: BaseAuthProviderConfig{
			Enabled: false,
		},
		FacebookAuth: BaseAuthProviderConfig{
			Enabled: false,
		},
		GithubAuth: BaseAuthProviderConfig{
			Enabled: false,
		},
		GitlabAuth: BaseAuthProviderConfig{
			Enabled: false,
		},
		DiscordAuth: BaseAuthProviderConfig{
			Enabled: false,
		},
		TwitterAuth: BaseAuthProviderConfig{
			Enabled: false,
		},
		MicrosoftAuth: BaseAuthProviderConfig{
			Enabled: false,
		},
		SpotifyAuth: BaseAuthProviderConfig{
			Enabled: false,
		},
		KakaoAuth: BaseAuthProviderConfig{
			Enabled: false,
		},
		TwitchAuth: BaseAuthProviderConfig{
			Enabled: false,
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
		validation.Field(&s.RecordAuthToken),
		validation.Field(&s.RecordPasswordResetToken),
		validation.Field(&s.RecordEmailChangeToken),
		validation.Field(&s.RecordVerificationToken),
		validation.Field(&s.Smtp),
		validation.Field(&s.S3),
		validation.Field(&s.AppleAuth),
		validation.Field(&s.GoogleAuth),
		validation.Field(&s.FacebookAuth),
		validation.Field(&s.GithubAuth),
		validation.Field(&s.GitlabAuth),
		validation.Field(&s.DiscordAuth),
		validation.Field(&s.TwitterAuth),
		validation.Field(&s.MicrosoftAuth),
		validation.Field(&s.SpotifyAuth),
		validation.Field(&s.KakaoAuth),
		validation.Field(&s.TwitchAuth),
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
		&clone.RecordAuthToken.Secret,
		&clone.RecordPasswordResetToken.Secret,
		&clone.RecordEmailChangeToken.Secret,
		&clone.RecordVerificationToken.Secret,
		&clone.AppleAuth.ClientSecret,
		&clone.AppleAuth.SigningKey,
		&clone.GoogleAuth.ClientSecret,
		&clone.FacebookAuth.ClientSecret,
		&clone.GithubAuth.ClientSecret,
		&clone.GitlabAuth.ClientSecret,
		&clone.DiscordAuth.ClientSecret,
		&clone.TwitterAuth.ClientSecret,
		&clone.MicrosoftAuth.ClientSecret,
		&clone.SpotifyAuth.ClientSecret,
		&clone.KakaoAuth.ClientSecret,
		&clone.TwitchAuth.ClientSecret,
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
		auth.NameApple:     s.AppleAuth,
		auth.NameGoogle:    s.GoogleAuth,
		auth.NameFacebook:  s.FacebookAuth,
		auth.NameGithub:    s.GithubAuth,
		auth.NameGitlab:    s.GitlabAuth,
		auth.NameDiscord:   s.DiscordAuth,
		auth.NameTwitter:   s.TwitterAuth,
		auth.NameMicrosoft: s.MicrosoftAuth,
		auth.NameSpotify:   s.SpotifyAuth,
		auth.NameKakao:     s.KakaoAuth,
		auth.NameTwitch:    s.TwitchAuth,
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

// AuthProviderConfig is a common interface for all provider configs.
type AuthProviderConfig interface {
	validation.Validatable

	// IsEnabled returns true if the provider is enabled in the config.
	IsEnabled() bool
	// SetupProvider loads the current AuthProviderConfig into the specified provider.
	SetupProvider(provider auth.Provider) error
}

// BaseAuthProviderConfig is the provider config used by all starndard providers
type BaseAuthProviderConfig struct {
	Enabled      bool   `form:"enabled" json:"enabled"`
	ClientId     string `form:"clientId" json:"clientId,omitempty"`
	ClientSecret string `form:"clientSecret" json:"clientSecret,omitempty"`
	AuthUrl      string `form:"authUrl" json:"authUrl,omitempty"`
	TokenUrl     string `form:"tokenUrl" json:"tokenUrl,omitempty"`
	UserApiUrl   string `form:"userApiUrl" json:"userApiUrl,omitempty"`
}

// IsEnabled implements [AuthProviderConfig] interface
func (c BaseAuthProviderConfig) IsEnabled() bool {
	return c.Enabled
}

// Validate implements [validation.Validatable] interface.
func (c BaseAuthProviderConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClientId, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.ClientSecret, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.AuthUrl, is.URL),
		validation.Field(&c.TokenUrl, is.URL),
		validation.Field(&c.UserApiUrl, is.URL),
	)
}

// SetupProvider implements [AuthProviderConfig] interface
func (c BaseAuthProviderConfig) SetupProvider(provider auth.Provider) error {
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

type AppleAuthProviderConfig struct {
	BaseAuthProviderConfig
	TeamId     string `form:"teamId" json:"teamId,omitempty"`
	KeyId      string `form:"keyId" json:"keyId,omitempty"`
	SigningKey string `form:"signingKey" json:"signingKey,omitempty"`
}

// Validate implements [validation.Validatable] interface.
func (c AppleAuthProviderConfig) Validate() error {
	requireSecret := c.Enabled && c.TeamId == "" && c.KeyId == "" && c.SigningKey == ""

	return validation.ValidateStruct(&c,
		validation.Field(&c.ClientId, validation.When(c.Enabled, validation.Required)),
		validation.Field(&c.ClientSecret, validation.When(requireSecret, validation.Required)),
		validation.Field(&c.AuthUrl, is.URL),
		validation.Field(&c.TokenUrl, is.URL),
		validation.Field(&c.UserApiUrl, is.URL),
		validation.Field(&c.TeamId, validation.When(c.Enabled && !requireSecret, validation.Required)),
		validation.Field(&c.KeyId, validation.When(c.Enabled && !requireSecret, validation.Required)),
		validation.Field(&c.SigningKey, validation.When(c.Enabled && !requireSecret, validation.Required)),
	)
}

// SetupProvider implements [AuthProviderConfig] interface
func (c AppleAuthProviderConfig) SetupProvider(provider auth.Provider) error {
	err := c.BaseAuthProviderConfig.SetupProvider(provider)
	if err != nil {
		return err
	}

	appleProvider, ok := provider.(*auth.Apple)
	if !ok {
		return nil
	}

	if c.TeamId != "" {
		appleProvider.SetTeamId(c.TeamId)
	}

	if c.KeyId != "" {
		appleProvider.SetKeyId(c.KeyId)
	}

	if c.SigningKey != "" {
		appleProvider.SetSigningKey(c.SigningKey)
	}

	if appleProvider.TeamId() != "" && appleProvider.KeyId() != "" && appleProvider.SigningKey() != "" {
		// TODO: maybe generate outside?
		validity := time.Hour*24*180 - time.Second
		clientSecret, err := auth.GenerateAppleClientSecret(
			appleProvider.SigningKey(),
			appleProvider.TeamId(),
			appleProvider.ClientId(),
			appleProvider.KeyId(),
			validity,
		)
		if err != nil {
			return err
		}
		appleProvider.SetClientSecret(clientSecret)
	}

	return nil
}

// -------------------------------------------------------------------

// Deprecated: Will be removed in v0.9!
type EmailAuthConfig struct {
	Enabled           bool     `form:"enabled" json:"enabled"`
	ExceptDomains     []string `form:"exceptDomains" json:"exceptDomains"`
	OnlyDomains       []string `form:"onlyDomains" json:"onlyDomains"`
	MinPasswordLength int      `form:"minPasswordLength" json:"minPasswordLength"`
}

// Deprecated: Will be removed in v0.9!
func (c EmailAuthConfig) Validate() error {
	return nil
}
