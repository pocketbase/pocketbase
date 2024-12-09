package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
)

const (
	paramsTable = "_params"

	paramsKeySettings = "settings"

	systemHookIdSettings = "__pbSettingsSystemHook__"
)

func (app *BaseApp) registerSettingsHooks() {
	saveFunc := func(me *ModelEvent) error {
		if err := me.Next(); err != nil {
			return err
		}

		if me.Model.PK() == paramsKeySettings {
			// auto reload the app settings because we don't know whether
			// the Settings model is the app one or a different one
			return errors.Join(
				me.App.Settings().PostScan(),
				me.App.ReloadSettings(),
			)
		}

		return nil
	}

	app.OnModelAfterCreateSuccess(paramsTable).Bind(&hook.Handler[*ModelEvent]{
		Id:       systemHookIdSettings,
		Func:     saveFunc,
		Priority: -999,
	})

	app.OnModelAfterUpdateSuccess(paramsTable).Bind(&hook.Handler[*ModelEvent]{
		Id:       systemHookIdSettings,
		Func:     saveFunc,
		Priority: -999,
	})

	app.OnModelDelete(paramsTable).Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdSettings,
		Func: func(me *ModelEvent) error {
			if me.Model.PK() == paramsKeySettings {
				return errors.New("the app params settings cannot be deleted")
			}

			return me.Next()
		},
		Priority: -999,
	})

	app.OnCollectionUpdate().Bind(&hook.Handler[*CollectionEvent]{
		Id: systemHookIdSettings,
		Func: func(e *CollectionEvent) error {
			oldCollection, err := e.App.FindCachedCollectionByNameOrId(e.Collection.Id)
			if err != nil {
				return fmt.Errorf("failed to retrieve old cached collection: %w", err)
			}

			err = e.Next()
			if err != nil {
				return err
			}

			// update existing rate limit rules on collection rename
			if oldCollection.Name != e.Collection.Name {
				var hasChange bool

				rules := e.App.Settings().RateLimits.Rules
				for i := 0; i < len(rules); i++ {
					if strings.HasPrefix(rules[i].Label, oldCollection.Name+":") {
						rules[i].Label = strings.Replace(rules[i].Label, oldCollection.Name+":", e.Collection.Name+":", 1)
						hasChange = true
					}
				}

				if hasChange {
					e.App.Settings().RateLimits.Rules = rules
					err = e.App.Save(e.App.Settings())
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
		Priority: 99,
	})
}

var (
	_ Model         = (*Settings)(nil)
	_ PostValidator = (*Settings)(nil)
	_ DBExporter    = (*Settings)(nil)
)

type settings struct {
	SMTP         SMTPConfig         `form:"smtp" json:"smtp"`
	Backups      BackupsConfig      `form:"backups" json:"backups"`
	S3           S3Config           `form:"s3" json:"s3"`
	Meta         MetaConfig         `form:"meta" json:"meta"`
	RateLimits   RateLimitsConfig   `form:"rateLimits" json:"rateLimits"`
	TrustedProxy TrustedProxyConfig `form:"trustedProxy" json:"trustedProxy"`
	Batch        BatchConfig        `form:"batch" json:"batch"`
	Logs         LogsConfig         `form:"logs" json:"logs"`
}

// Settings defines the PocketBase app settings.
type Settings struct {
	settings

	mu    sync.RWMutex
	isNew bool
}

func newDefaultSettings() *Settings {
	return &Settings{
		isNew: true,
		settings: settings{
			Meta: MetaConfig{
				AppName:       "Acme",
				AppURL:        "http://localhost:8090",
				HideControls:  false,
				SenderName:    "Support",
				SenderAddress: "support@example.com",
			},
			Logs: LogsConfig{
				MaxDays: 5,
				LogIP:   true,
			},
			SMTP: SMTPConfig{
				Enabled:  false,
				Host:     "smtp.example.com",
				Port:     587,
				Username: "",
				Password: "",
				TLS:      false,
			},
			Backups: BackupsConfig{
				CronMaxKeep: 3,
			},
			Batch: BatchConfig{
				Enabled:     false,
				MaxRequests: 50,
				Timeout:     3,
			},
			RateLimits: RateLimitsConfig{
				Enabled: false, // @todo once tested enough enable by default for new installations
				Rules: []RateLimitRule{
					{Label: "*:auth", MaxRequests: 2, Duration: 3},
					{Label: "*:create", MaxRequests: 20, Duration: 5},
					{Label: "/api/batch", MaxRequests: 3, Duration: 1},
					{Label: "/api/", MaxRequests: 300, Duration: 10},
				},
			},
		},
	}
}

// TableName implements [Model.TableName] interface method.
func (s *Settings) TableName() string {
	return paramsTable
}

// PK implements [Model.LastSavedPK] interface method.
func (s *Settings) LastSavedPK() any {
	return paramsKeySettings
}

// PK implements [Model.PK] interface method.
func (s *Settings) PK() any {
	return paramsKeySettings
}

// IsNew implements [Model.IsNew] interface method.
func (s *Settings) IsNew() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.isNew
}

// MarkAsNew implements [Model.MarkAsNew] interface method.
func (s *Settings) MarkAsNew() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isNew = true
}

// MarkAsNew implements [Model.MarkAsNotNew] interface method.
func (s *Settings) MarkAsNotNew() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isNew = false
}

// PostScan implements [Model.PostScan] interface method.
func (s *Settings) PostScan() error {
	s.MarkAsNotNew()
	return nil
}

// String returns a serialized string representation of the current settings.
func (s *Settings) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	raw, _ := json.Marshal(s)
	return string(raw)
}

// DBExport prepares and exports the current settings for db persistence.
func (s *Settings) DBExport(app App) (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := types.NowDateTime()

	result := map[string]any{
		"id": s.PK(),
	}

	if s.IsNew() {
		result["created"] = now
	}
	result["updated"] = now

	encoded, err := json.Marshal(s.settings)
	if err != nil {
		return nil, err
	}

	encryptionKey := os.Getenv(app.EncryptionEnv())
	if encryptionKey != "" {
		encryptVal, encryptErr := security.Encrypt(encoded, encryptionKey)
		if encryptErr != nil {
			return nil, encryptErr
		}

		result["value"] = encryptVal
	} else {
		result["value"] = encoded
	}

	return result, nil
}

// PostValidate implements the [PostValidator] interface and defines
// the Settings model validations.
func (s *Settings) PostValidate(ctx context.Context, app App) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return validation.ValidateStructWithContext(ctx, s,
		validation.Field(&s.Meta),
		validation.Field(&s.Logs),
		validation.Field(&s.SMTP),
		validation.Field(&s.S3),
		validation.Field(&s.Backups),
		validation.Field(&s.Batch),
		validation.Field(&s.RateLimits),
		validation.Field(&s.TrustedProxy),
	)
}

// Merge merges the "other" settings into the current one.
func (s *Settings) Merge(other *Settings) error {
	other.mu.RLock()
	defer other.mu.RUnlock()

	raw, err := json.Marshal(other.settings)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return json.Unmarshal(raw, &s)
}

// Clone creates a new deep copy of the current settings.
func (s *Settings) Clone() (*Settings, error) {
	clone := &Settings{
		isNew: s.isNew,
	}

	if err := clone.Merge(s); err != nil {
		return nil, err
	}

	return clone, nil
}

// MarshalJSON implements the [json.Marshaler] interface.
//
// Note that sensitive fields (S3 secret, SMTP password, etc.) are excluded.
func (s *Settings) MarshalJSON() ([]byte, error) {
	s.mu.RLock()
	copy := s.settings
	s.mu.RUnlock()

	sensitiveFields := []*string{
		&copy.SMTP.Password,
		&copy.S3.Secret,
		&copy.Backups.S3.Secret,
	}

	// mask all sensitive fields
	for _, v := range sensitiveFields {
		if v != nil && *v != "" {
			*v = ""
		}
	}

	return json.Marshal(copy)
}

// -------------------------------------------------------------------

type SMTPConfig struct {
	Enabled  bool   `form:"enabled" json:"enabled"`
	Port     int    `form:"port" json:"port"`
	Host     string `form:"host" json:"host"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password,omitempty"`

	// SMTP AUTH - PLAIN (default) or LOGIN
	AuthMethod string `form:"authMethod" json:"authMethod"`

	// Whether to enforce TLS encryption for the mail server connection.
	//
	// When set to false StartTLS command is send, leaving the server
	// to decide whether to upgrade the connection or not.
	TLS bool `form:"tls" json:"tls"`

	// LocalName is optional domain name or IP address used for the
	// EHLO/HELO exchange (if not explicitly set, defaults to "localhost").
	//
	// This is required only by some SMTP servers, such as Gmail SMTP-relay.
	LocalName string `form:"localName" json:"localName"`
}

// Validate makes SMTPConfig validatable by implementing [validation.Validatable] interface.
func (c SMTPConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Host,
			validation.When(c.Enabled, validation.Required),
			is.Host,
		),
		validation.Field(
			&c.Port,
			validation.When(c.Enabled, validation.Required),
			validation.Min(0),
		),
		validation.Field(
			&c.AuthMethod,
			// don't require it for backward compatibility
			// (fallback internally to PLAIN)
			// validation.When(c.Enabled, validation.Required),
			validation.In(mailer.SMTPAuthLogin, mailer.SMTPAuthPlain),
		),
		validation.Field(&c.LocalName, is.Host),
	)
}

// -------------------------------------------------------------------

type S3Config struct {
	Enabled        bool   `form:"enabled" json:"enabled"`
	Bucket         string `form:"bucket" json:"bucket"`
	Region         string `form:"region" json:"region"`
	Endpoint       string `form:"endpoint" json:"endpoint"`
	AccessKey      string `form:"accessKey" json:"accessKey"`
	Secret         string `form:"secret" json:"secret,omitempty"`
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

type BatchConfig struct {
	Enabled bool `form:"enabled" json:"enabled"`

	// MaxRequests is the maximum allowed batch request to execute.
	MaxRequests int `form:"maxRequests" json:"maxRequests"`

	// Timeout is the the max duration in seconds to wait before cancelling the batch transaction.
	Timeout int64 `form:"timeout" json:"timeout"`

	// MaxBodySize is the maximum allowed batch request body size in bytes.
	//
	// If not set, fallbacks to max ~128MB.
	MaxBodySize int64 `form:"maxBodySize" json:"maxBodySize"`
}

// Validate makes BatchConfig validatable by implementing [validation.Validatable] interface.
func (c BatchConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.MaxRequests, validation.When(c.Enabled, validation.Required), validation.Min(0)),
		validation.Field(&c.Timeout, validation.When(c.Enabled, validation.Required), validation.Min(0)),
		validation.Field(&c.MaxBodySize, validation.Min(0)),
	)
}

// -------------------------------------------------------------------

type BackupsConfig struct {
	// Cron is a cron expression to schedule auto backups, eg. "* * * * *".
	//
	// Leave it empty to disable the auto backups functionality.
	Cron string `form:"cron" json:"cron"`

	// CronMaxKeep is the the max number of cron generated backups to
	// keep before removing older entries.
	//
	// This field works only when the cron config has valid cron expression.
	CronMaxKeep int `form:"cronMaxKeep" json:"cronMaxKeep"`

	// S3 is an optional S3 storage config specifying where to store the app backups.
	S3 S3Config `form:"s3" json:"s3"`
}

// Validate makes BackupsConfig validatable by implementing [validation.Validatable] interface.
func (c BackupsConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.S3),
		validation.Field(&c.Cron, validation.By(checkCronExpression)),
		validation.Field(
			&c.CronMaxKeep,
			validation.When(c.Cron != "", validation.Required),
			validation.Min(1),
		),
	)
}

func checkCronExpression(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	_, err := cron.NewSchedule(v)
	if err != nil {
		return validation.NewError("validation_invalid_cron", err.Error())
	}

	return nil
}

// -------------------------------------------------------------------

type MetaConfig struct {
	AppName       string `form:"appName" json:"appName"`
	AppURL        string `form:"appURL" json:"appURL"`
	SenderName    string `form:"senderName" json:"senderName"`
	SenderAddress string `form:"senderAddress" json:"senderAddress"`
	HideControls  bool   `form:"hideControls" json:"hideControls"`
}

// Validate makes MetaConfig validatable by implementing [validation.Validatable] interface.
func (c MetaConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.AppName, validation.Required, validation.Length(1, 255)),
		validation.Field(&c.AppURL, validation.Required, is.URL),
		validation.Field(&c.SenderName, validation.Required, validation.Length(1, 255)),
		validation.Field(&c.SenderAddress, is.EmailFormat, validation.Required),
	)
}

// -------------------------------------------------------------------

type LogsConfig struct {
	MaxDays   int  `form:"maxDays" json:"maxDays"`
	MinLevel  int  `form:"minLevel" json:"minLevel"`
	LogIP     bool `form:"logIP" json:"logIP"`
	LogAuthId bool `form:"logAuthId" json:"logAuthId"`
}

// Validate makes LogsConfig validatable by implementing [validation.Validatable] interface.
func (c LogsConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.MaxDays, validation.Min(0)),
	)
}

// -------------------------------------------------------------------

type TrustedProxyConfig struct {
	// Headers is a list of explicit trusted header(s) to check.
	Headers []string `form:"headers" json:"headers"`

	// UseLeftmostIP specifies to use the left-mostish IP from the trusted headers.
	//
	// Note that this could be insecure when used with X-Forwarded-For header
	// because some proxies like AWS ELB allow users to prepend their own header value
	// before appending the trusted ones.
	UseLeftmostIP bool `form:"useLeftmostIP" json:"useLeftmostIP"`
}

// MarshalJSON implements the [json.Marshaler] interface.
func (c TrustedProxyConfig) MarshalJSON() ([]byte, error) {
	type alias TrustedProxyConfig

	// serialize as empty array
	if c.Headers == nil {
		c.Headers = []string{}
	}

	return json.Marshal(alias(c))
}

// Validate makes RateLimitRule validatable by implementing [validation.Validatable] interface.
func (c TrustedProxyConfig) Validate() error {
	return nil
}

// -------------------------------------------------------------------

type RateLimitsConfig struct {
	Rules   []RateLimitRule `form:"rules" json:"rules"`
	Enabled bool            `form:"enabled" json:"enabled"`
}

// FindRateLimitRule returns the first matching rule based on the provided labels.
//
// Optionally you can further specify a list of valid RateLimitRule.Audience values to further filter the matching rule
// (aka. the rule Audience will have to exist in one of the specified options).
func (c *RateLimitsConfig) FindRateLimitRule(searchLabels []string, optOnlyAudience ...string) (RateLimitRule, bool) {
	var prefixRules []int

	for i, label := range searchLabels {
		// check for direct match
		for j := range c.Rules {
			if label == c.Rules[j].Label &&
				(len(optOnlyAudience) == 0 || slices.Contains(optOnlyAudience, c.Rules[j].Audience)) {
				return c.Rules[j], true
			}

			if i == 0 && strings.HasSuffix(c.Rules[j].Label, "/") {
				prefixRules = append(prefixRules, j)
			}
		}

		// check for prefix match
		if len(prefixRules) > 0 {
			for j := range prefixRules {
				if strings.HasPrefix(label+"/", c.Rules[prefixRules[j]].Label) &&
					(len(optOnlyAudience) == 0 || slices.Contains(optOnlyAudience, c.Rules[prefixRules[j]].Audience)) {
					return c.Rules[prefixRules[j]], true
				}
			}
		}
	}

	return RateLimitRule{}, false
}

// MarshalJSON implements the [json.Marshaler] interface.
func (c RateLimitsConfig) MarshalJSON() ([]byte, error) {
	type alias RateLimitsConfig

	// serialize as empty array
	if c.Rules == nil {
		c.Rules = []RateLimitRule{}
	}

	return json.Marshal(alias(c))
}

// Validate makes RateLimitsConfig validatable by implementing [validation.Validatable] interface.
func (c RateLimitsConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Rules,
			validation.When(c.Enabled, validation.Required),
			validation.By(checkUniqueRuleLabel),
		),
	)
}

func checkUniqueRuleLabel(value any) error {
	rules, ok := value.([]RateLimitRule)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	existing := make([]string, 0, len(rules))

	for i, rule := range rules {
		fullKey := rule.Label + "@@" + rule.Audience

		var conflicts bool
		for _, key := range existing {
			if strings.HasPrefix(key, fullKey) || strings.HasPrefix(fullKey, key) {
				conflicts = true
				break
			}
		}

		if conflicts {
			return validation.Errors{
				strconv.Itoa(i): validation.Errors{
					"label": validation.NewError("validation_conflicting_rate_limit_rule", "Rate limit rule configuration with label {{.label}} already exists or conflicts with another rule.").
						SetParams(map[string]any{"label": rule.Label}),
				},
			}
		} else {
			existing = append(existing, fullKey)
		}
	}

	return nil
}

var rateLimitRuleLabelRegex = regexp.MustCompile(`^(\w+\ \/[\w\/-]*|\/[\w\/-]*|^\w+\:\w+|\*\:\w+|\w+)$`)

// The allowed RateLimitRule.Audience values
const (
	RateLimitRuleAudienceAll   = ""
	RateLimitRuleAudienceGuest = "@guest"
	RateLimitRuleAudienceAuth  = "@auth"
)

type RateLimitRule struct {
	// Label is the identifier of the current rule.
	//
	// It could be a tag, complete path or path prerefix (when ends with `/`).
	//
	// Example supported labels:
	//   - test_a (plain text "tag")
	//   - users:create
	//   - *:create
	//   - /
	//   - /api
	//   - POST /api/collections/
	Label string `form:"label" json:"label"`

	// Audience specifies the auth group the rule should apply for:
	//   - ""      - both guests and authenticated users (default)
	//   - "guest" - only for guests
	//   - "auth"  - only for authenticated users
	Audience string `form:"audience" json:"audience"`

	// Duration specifies the interval (in seconds) per which to reset
	// the counted/accumulated rate limiter tokens.
	Duration int64 `form:"duration" json:"duration"`

	// MaxRequests is the max allowed number of requests per Duration.
	MaxRequests int `form:"maxRequests" json:"maxRequests"`
}

// Validate makes RateLimitRule validatable by implementing [validation.Validatable] interface.
func (c RateLimitRule) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Label, validation.Required, validation.Match(rateLimitRuleLabelRegex)),
		validation.Field(&c.MaxRequests, validation.Required, validation.Min(1)),
		validation.Field(&c.Duration, validation.Required, validation.Min(1)),
		validation.Field(&c.Audience,
			validation.In(RateLimitRuleAudienceAll, RateLimitRuleAudienceGuest, RateLimitRuleAudienceAuth),
		),
	)
}

// DurationTime returns the tag's Duration as [time.Duration].
func (c RateLimitRule) DurationTime() time.Duration {
	return time.Duration(c.Duration) * time.Second
}
