package core_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/auth"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestCollectionAuthOptionsValidate(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name           string
		collection     func(app core.App) (*core.Collection, error)
		expectedErrors []string
	}{
		// authRule
		{
			name: "nil authRule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.AuthRule = nil
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "empty authRule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.AuthRule = types.Pointer("")
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "invalid authRule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.AuthRule = types.Pointer("missing != ''")
				return c, nil
			},
			expectedErrors: []string{"authRule"},
		},
		{
			name: "valid authRule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.AuthRule = types.Pointer("id != ''")
				return c, nil
			},
			expectedErrors: []string{},
		},

		// manageRule
		{
			name: "nil manageRule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.ManageRule = nil
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "empty manageRule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.ManageRule = types.Pointer("")
				return c, nil
			},
			expectedErrors: []string{"manageRule"},
		},
		{
			name: "invalid manageRule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.ManageRule = types.Pointer("missing != ''")
				return c, nil
			},
			expectedErrors: []string{"manageRule"},
		},
		{
			name: "valid manageRule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.ManageRule = types.Pointer("id != ''")
				return c, nil
			},
			expectedErrors: []string{},
		},

		// passwordAuth
		{
			name: "trigger passwordAuth validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.PasswordAuth = core.PasswordAuthConfig{
					Enabled: true,
				}
				return c, nil
			},
			expectedErrors: []string{"passwordAuth"},
		},
		{
			name: "passwordAuth with non-unique identity fields",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(&core.TextField{Name: "test"})
				c.PasswordAuth = core.PasswordAuthConfig{
					Enabled:        true,
					IdentityFields: []string{"email", "test"},
				}
				return c, nil
			},
			expectedErrors: []string{"passwordAuth"},
		},
		{
			name: "passwordAuth with non-unique identity fields",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.Fields.Add(&core.TextField{Name: "test"})
				c.AddIndex("auth_test_idx", true, "test", "")
				c.PasswordAuth = core.PasswordAuthConfig{
					Enabled:        true,
					IdentityFields: []string{"email", "test"},
				}
				return c, nil
			},
			expectedErrors: []string{},
		},

		// oauth2
		{
			name: "trigger oauth2 validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.OAuth2 = core.OAuth2Config{
					Enabled: true,
					Providers: []core.OAuth2ProviderConfig{
						{Name: "missing"},
					},
				}
				return c, nil
			},
			expectedErrors: []string{"oauth2"},
		},

		// otp
		{
			name: "trigger otp validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.OTP = core.OTPConfig{
					Enabled:  true,
					Duration: -10,
				}
				return c, nil
			},
			expectedErrors: []string{"otp"},
		},

		// mfa
		{
			name: "trigger mfa validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.MFA = core.MFAConfig{
					Enabled:  true,
					Duration: -10,
				}
				return c, nil
			},
			expectedErrors: []string{"mfa"},
		},
		{
			name: "mfa enabled with < 2 auth methods",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.MFA.Enabled = true
				c.PasswordAuth.Enabled = true
				c.OTP.Enabled = false
				c.OAuth2.Enabled = false
				return c, nil
			},
			expectedErrors: []string{"mfa"},
		},
		{
			name: "mfa enabled with >= 2 auth methods",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.MFA.Enabled = true
				c.PasswordAuth.Enabled = true
				c.OTP.Enabled = true
				c.OAuth2.Enabled = false
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "mfa disabled with invalid rule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.PasswordAuth.Enabled = true
				c.OTP.Enabled = true
				c.MFA.Enabled = false
				c.MFA.Rule = "invalid"
				return c, nil
			},
			expectedErrors: []string{},
		},
		{
			name: "mfa enabled with invalid rule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.PasswordAuth.Enabled = true
				c.OTP.Enabled = true
				c.MFA.Enabled = true
				c.MFA.Rule = "invalid"
				return c, nil
			},
			expectedErrors: []string{"mfa"},
		},
		{
			name: "mfa enabled with valid rule",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.PasswordAuth.Enabled = true
				c.OTP.Enabled = true
				c.MFA.Enabled = true
				c.MFA.Rule = "1=1"
				return c, nil
			},
			expectedErrors: []string{},
		},

		// tokens
		{
			name: "trigger authToken validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.AuthToken.Secret = ""
				return c, nil
			},
			expectedErrors: []string{"authToken"},
		},
		{
			name: "trigger passwordResetToken validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.PasswordResetToken.Secret = ""
				return c, nil
			},
			expectedErrors: []string{"passwordResetToken"},
		},
		{
			name: "trigger emailChangeToken validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.EmailChangeToken.Secret = ""
				return c, nil
			},
			expectedErrors: []string{"emailChangeToken"},
		},
		{
			name: "trigger verificationToken validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.VerificationToken.Secret = ""
				return c, nil
			},
			expectedErrors: []string{"verificationToken"},
		},
		{
			name: "trigger fileToken validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.FileToken.Secret = ""
				return c, nil
			},
			expectedErrors: []string{"fileToken"},
		},

		// templates
		{
			name: "trigger verificationTemplate validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.VerificationTemplate.Body = ""
				return c, nil
			},
			expectedErrors: []string{"verificationTemplate"},
		},
		{
			name: "trigger resetPasswordTemplate validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.ResetPasswordTemplate.Body = ""
				return c, nil
			},
			expectedErrors: []string{"resetPasswordTemplate"},
		},
		{
			name: "trigger confirmEmailChangeTemplate validations",
			collection: func(app core.App) (*core.Collection, error) {
				c := core.NewAuthCollection("new_auth")
				c.ConfirmEmailChangeTemplate.Body = ""
				return c, nil
			},
			expectedErrors: []string{"confirmEmailChangeTemplate"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			collection, err := s.collection(app)
			if err != nil {
				t.Fatalf("Failed to retrieve test collection: %v", err)
			}

			result := app.Validate(collection)

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestEmailTemplateValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		template       core.EmailTemplate
		expectedErrors []string
	}{
		{
			"zero value",
			core.EmailTemplate{},
			[]string{"subject", "body"},
		},
		{
			"non-empty data",
			core.EmailTemplate{
				Subject: "a",
				Body:    "b",
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.template.Validate()

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestEmailTemplateResolve(t *testing.T) {
	template := core.EmailTemplate{
		Subject: "test_subject {PARAM3} {PARAM1}-{PARAM2} repeat-{PARAM1}",
		Body:    "test_body {PARAM3} {PARAM2}-{PARAM1} repeat-{PARAM2}",
	}

	scenarios := []struct {
		name            string
		placeholders    map[string]any
		template        core.EmailTemplate
		expectedSubject string
		expectedBody    string
	}{
		{
			"no placeholders",
			nil,
			template,
			template.Subject,
			template.Body,
		},
		{
			"no matching placeholders",
			map[string]any{"{A}": "abc", "{B}": 456},
			template,
			template.Subject,
			template.Body,
		},
		{
			"at least one matching placeholder",
			map[string]any{"{PARAM1}": "abc", "{PARAM2}": 456},
			template,
			"test_subject {PARAM3} abc-456 repeat-abc",
			"test_body {PARAM3} 456-abc repeat-456",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			subject, body := s.template.Resolve(s.placeholders)

			if subject != s.expectedSubject {
				t.Fatalf("Expected subject\n%v\ngot\n%v", s.expectedSubject, subject)
			}

			if body != s.expectedBody {
				t.Fatalf("Expected body\n%v\ngot\n%v", s.expectedBody, body)
			}
		})
	}
}

func TestTokenConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.TokenConfig
		expectedErrors []string
	}{
		{
			"zero value",
			core.TokenConfig{},
			[]string{"secret", "duration"},
		},
		{
			"invalid data",
			core.TokenConfig{
				Secret:   strings.Repeat("a", 29),
				Duration: 9,
			},
			[]string{"secret", "duration"},
		},
		{
			"valid data",
			core.TokenConfig{
				Secret:   strings.Repeat("a", 30),
				Duration: 10,
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.config.Validate()

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestTokenConfigDurationTime(t *testing.T) {
	scenarios := []struct {
		config   core.TokenConfig
		expected time.Duration
	}{
		{core.TokenConfig{}, 0 * time.Second},
		{core.TokenConfig{Duration: 1234}, 1234 * time.Second},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%d", i, s.config.Duration), func(t *testing.T) {
			result := s.config.DurationTime()

			if result != s.expected {
				t.Fatalf("Expected duration %d, got %d", s.expected, result)
			}
		})
	}
}

func TestAuthAlertConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.AuthAlertConfig
		expectedErrors []string
	}{
		{
			"zero value (disabled)",
			core.AuthAlertConfig{},
			[]string{"emailTemplate"},
		},
		{
			"zero value (enabled)",
			core.AuthAlertConfig{Enabled: true},
			[]string{"emailTemplate"},
		},
		{
			"invalid template",
			core.AuthAlertConfig{
				EmailTemplate: core.EmailTemplate{Body: "", Subject: "b"},
			},
			[]string{"emailTemplate"},
		},
		{
			"valid data",
			core.AuthAlertConfig{
				EmailTemplate: core.EmailTemplate{Body: "a", Subject: "b"},
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.config.Validate()

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestOTPConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.OTPConfig
		expectedErrors []string
	}{
		{
			"zero value (disabled)",
			core.OTPConfig{},
			[]string{"emailTemplate"},
		},
		{
			"zero value (enabled)",
			core.OTPConfig{Enabled: true},
			[]string{"duration", "length", "emailTemplate"},
		},
		{
			"invalid length (< 3)",
			core.OTPConfig{
				Enabled:       true,
				EmailTemplate: core.EmailTemplate{Body: "a", Subject: "b"},
				Duration:      100,
				Length:        3,
			},
			[]string{"length"},
		},
		{
			"invalid duration (< 10)",
			core.OTPConfig{
				Enabled:       true,
				EmailTemplate: core.EmailTemplate{Body: "a", Subject: "b"},
				Duration:      9,
				Length:        100,
			},
			[]string{"duration"},
		},
		{
			"invalid duration (> 86400)",
			core.OTPConfig{
				Enabled:       true,
				EmailTemplate: core.EmailTemplate{Body: "a", Subject: "b"},
				Duration:      86401,
				Length:        100,
			},
			[]string{"duration"},
		},
		{
			"invalid template (triggering EmailTemplate validations)",
			core.OTPConfig{
				Enabled:       true,
				EmailTemplate: core.EmailTemplate{Body: "", Subject: "b"},
				Duration:      86400,
				Length:        4,
			},
			[]string{"emailTemplate"},
		},
		{
			"valid data",
			core.OTPConfig{
				Enabled:       true,
				EmailTemplate: core.EmailTemplate{Body: "a", Subject: "b"},
				Duration:      86400,
				Length:        4,
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.config.Validate()

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestOTPConfigDurationTime(t *testing.T) {
	scenarios := []struct {
		config   core.OTPConfig
		expected time.Duration
	}{
		{core.OTPConfig{}, 0 * time.Second},
		{core.OTPConfig{Duration: 1234}, 1234 * time.Second},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%d", i, s.config.Duration), func(t *testing.T) {
			result := s.config.DurationTime()

			if result != s.expected {
				t.Fatalf("Expected duration %d, got %d", s.expected, result)
			}
		})
	}
}

func TestMFAConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.MFAConfig
		expectedErrors []string
	}{
		{
			"zero value (disabled)",
			core.MFAConfig{},
			[]string{},
		},
		{
			"zero value (enabled)",
			core.MFAConfig{Enabled: true},
			[]string{"duration"},
		},
		{
			"invalid duration (< 10)",
			core.MFAConfig{Enabled: true, Duration: 9},
			[]string{"duration"},
		},
		{
			"invalid duration (> 86400)",
			core.MFAConfig{Enabled: true, Duration: 86401},
			[]string{"duration"},
		},
		{
			"valid data",
			core.MFAConfig{Enabled: true, Duration: 86400},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.config.Validate()

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestMFAConfigDurationTime(t *testing.T) {
	scenarios := []struct {
		config   core.MFAConfig
		expected time.Duration
	}{
		{core.MFAConfig{}, 0 * time.Second},
		{core.MFAConfig{Duration: 1234}, 1234 * time.Second},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%d", i, s.config.Duration), func(t *testing.T) {
			result := s.config.DurationTime()

			if result != s.expected {
				t.Fatalf("Expected duration %d, got %d", s.expected, result)
			}
		})
	}
}

func TestPasswordAuthConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.PasswordAuthConfig
		expectedErrors []string
	}{
		{
			"zero value (disabled)",
			core.PasswordAuthConfig{},
			[]string{},
		},
		{
			"zero value (enabled)",
			core.PasswordAuthConfig{Enabled: true},
			[]string{"identityFields"},
		},
		{
			"empty values",
			core.PasswordAuthConfig{Enabled: true, IdentityFields: []string{"", ""}},
			[]string{"identityFields"},
		},
		{
			"valid data",
			core.PasswordAuthConfig{Enabled: true, IdentityFields: []string{"abc"}},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.config.Validate()

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestOAuth2ConfigGetProviderConfig(t *testing.T) {
	scenarios := []struct {
		name           string
		providerName   string
		config         core.OAuth2Config
		expectedExists bool
	}{
		{
			"zero value",
			"gitlab",
			core.OAuth2Config{},
			false,
		},
		{
			"empty config with valid provider",
			"gitlab",
			core.OAuth2Config{},
			false,
		},
		{
			"non-empty config with missing provider",
			"gitlab",
			core.OAuth2Config{Providers: []core.OAuth2ProviderConfig{{Name: "google"}, {Name: "github"}}},
			false,
		},
		{
			"config with existing provider",
			"github",
			core.OAuth2Config{Providers: []core.OAuth2ProviderConfig{{Name: "google"}, {Name: "github"}}},
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			config, exists := s.config.GetProviderConfig(s.providerName)

			if exists != s.expectedExists {
				t.Fatalf("Expected exists %v, got %v", s.expectedExists, exists)
			}

			if exists {
				if config.Name != s.providerName {
					t.Fatalf("Expected config with name %q, got %q", s.providerName, config.Name)
				}
			} else {
				if config.Name != "" {
					t.Fatalf("Expected empty config, got %v", config)
				}
			}
		})
	}
}

func TestOAuth2ConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.OAuth2Config
		expectedErrors []string
	}{
		{
			"zero value (disabled)",
			core.OAuth2Config{},
			[]string{},
		},
		{
			"zero value (enabled)",
			core.OAuth2Config{Enabled: true},
			[]string{},
		},
		{
			"unknown provider",
			core.OAuth2Config{Enabled: true, Providers: []core.OAuth2ProviderConfig{
				{Name: "missing", ClientId: "abc", ClientSecret: "456"},
			}},
			[]string{"providers"},
		},
		{
			"known provider with invalid data",
			core.OAuth2Config{Enabled: true, Providers: []core.OAuth2ProviderConfig{
				{Name: "gitlab", ClientId: "abc", TokenURL: "!invalid!"},
			}},
			[]string{"providers"},
		},
		{
			"known provider with valid data",
			core.OAuth2Config{Enabled: true, Providers: []core.OAuth2ProviderConfig{
				{Name: "gitlab", ClientId: "abc", ClientSecret: "456", TokenURL: "https://example.com"},
			}},
			[]string{},
		},
		{
			"known provider with valid data (duplicated)",
			core.OAuth2Config{Enabled: true, Providers: []core.OAuth2ProviderConfig{
				{Name: "gitlab", ClientId: "abc1", ClientSecret: "1", TokenURL: "https://example1.com"},
				{Name: "google", ClientId: "abc2", ClientSecret: "2", TokenURL: "https://example2.com"},
				{Name: "gitlab", ClientId: "abc3", ClientSecret: "3", TokenURL: "https://example3.com"},
			}},
			[]string{"providers"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.config.Validate()

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestOAuth2ProviderConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.OAuth2ProviderConfig
		expectedErrors []string
	}{
		{
			"zero value",
			core.OAuth2ProviderConfig{},
			[]string{"name", "clientId", "clientSecret"},
		},
		{
			"minimum valid data",
			core.OAuth2ProviderConfig{Name: "gitlab", ClientId: "abc", ClientSecret: "456"},
			[]string{},
		},
		{
			"non-existing provider",
			core.OAuth2ProviderConfig{Name: "missing", ClientId: "abc", ClientSecret: "456"},
			[]string{"name"},
		},
		{
			"invalid urls",
			core.OAuth2ProviderConfig{
				Name:         "gitlab",
				ClientId:     "abc",
				ClientSecret: "456",
				AuthURL:      "!invalid!",
				TokenURL:     "!invalid!",
				UserInfoURL:  "!invalid!",
			},
			[]string{"authURL", "tokenURL", "userInfoURL"},
		},
		{
			"valid urls",
			core.OAuth2ProviderConfig{
				Name:         "gitlab",
				ClientId:     "abc",
				ClientSecret: "456",
				AuthURL:      "https://example.com/a",
				TokenURL:     "https://example.com/b",
				UserInfoURL:  "https://example.com/c",
			},
			[]string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := s.config.Validate()

			tests.TestValidationErrors(t, result, s.expectedErrors)
		})
	}
}

func TestOAuth2ProviderConfigInitProvider(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.OAuth2ProviderConfig
		expectedConfig core.OAuth2ProviderConfig
		expectedError  bool
	}{
		{
			"empty config",
			core.OAuth2ProviderConfig{},
			core.OAuth2ProviderConfig{},
			true,
		},
		{
			"missing provider",
			core.OAuth2ProviderConfig{
				Name:         "missing",
				ClientId:     "test_ClientId",
				ClientSecret: "test_ClientSecret",
				AuthURL:      "test_AuthURL",
				TokenURL:     "test_TokenURL",
				UserInfoURL:  "test_UserInfoURL",
				DisplayName:  "test_DisplayName",
				PKCE:         types.Pointer(true),
			},
			core.OAuth2ProviderConfig{
				Name:         "missing",
				ClientId:     "test_ClientId",
				ClientSecret: "test_ClientSecret",
				AuthURL:      "test_AuthURL",
				TokenURL:     "test_TokenURL",
				UserInfoURL:  "test_UserInfoURL",
				DisplayName:  "test_DisplayName",
				PKCE:         types.Pointer(true),
			},
			true,
		},
		{
			"existing provider minimal",
			core.OAuth2ProviderConfig{
				Name: "gitlab",
			},
			core.OAuth2ProviderConfig{
				Name:         "gitlab",
				ClientId:     "",
				ClientSecret: "",
				AuthURL:      "https://gitlab.com/oauth/authorize",
				TokenURL:     "https://gitlab.com/oauth/token",
				UserInfoURL:  "https://gitlab.com/api/v4/user",
				DisplayName:  "GitLab",
				PKCE:         types.Pointer(true),
			},
			false,
		},
		{
			"existing provider with all fields",
			core.OAuth2ProviderConfig{
				Name:         "gitlab",
				ClientId:     "test_ClientId",
				ClientSecret: "test_ClientSecret",
				AuthURL:      "test_AuthURL",
				TokenURL:     "test_TokenURL",
				UserInfoURL:  "test_UserInfoURL",
				DisplayName:  "test_DisplayName",
				PKCE:         types.Pointer(true),
				Extra:        map[string]any{"a": 1},
			},
			core.OAuth2ProviderConfig{
				Name:         "gitlab",
				ClientId:     "test_ClientId",
				ClientSecret: "test_ClientSecret",
				AuthURL:      "test_AuthURL",
				TokenURL:     "test_TokenURL",
				UserInfoURL:  "test_UserInfoURL",
				DisplayName:  "test_DisplayName",
				PKCE:         types.Pointer(true),
				Extra:        map[string]any{"a": 1},
			},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			provider, err := s.config.InitProvider()

			hasErr := err != nil
			if hasErr != s.expectedError {
				t.Fatalf("Expected hasErr %v, got %v", s.expectedError, hasErr)
			}

			if hasErr {
				if provider != nil {
					t.Fatalf("Expected nil provider, got %v", provider)
				}
				return
			}

			factory, ok := auth.Providers[s.expectedConfig.Name]
			if !ok {
				t.Fatalf("Missing factory for provider %q", s.expectedConfig.Name)
			}

			expectedType := fmt.Sprintf("%T", factory())
			providerType := fmt.Sprintf("%T", provider)
			if expectedType != providerType {
				t.Fatalf("Expected provider instanceof %q, got %q", expectedType, providerType)
			}

			if provider.ClientId() != s.expectedConfig.ClientId {
				t.Fatalf("Expected ClientId %q, got %q", s.expectedConfig.ClientId, provider.ClientId())
			}

			if provider.ClientSecret() != s.expectedConfig.ClientSecret {
				t.Fatalf("Expected ClientSecret %q, got %q", s.expectedConfig.ClientSecret, provider.ClientSecret())
			}

			if provider.AuthURL() != s.expectedConfig.AuthURL {
				t.Fatalf("Expected AuthURL %q, got %q", s.expectedConfig.AuthURL, provider.AuthURL())
			}

			if provider.UserInfoURL() != s.expectedConfig.UserInfoURL {
				t.Fatalf("Expected UserInfoURL %q, got %q", s.expectedConfig.UserInfoURL, provider.UserInfoURL())
			}

			if provider.TokenURL() != s.expectedConfig.TokenURL {
				t.Fatalf("Expected TokenURL %q, got %q", s.expectedConfig.TokenURL, provider.TokenURL())
			}

			if provider.DisplayName() != s.expectedConfig.DisplayName {
				t.Fatalf("Expected DisplayName %q, got %q", s.expectedConfig.DisplayName, provider.DisplayName())
			}

			if provider.PKCE() != *s.expectedConfig.PKCE {
				t.Fatalf("Expected PKCE %v, got %v", *s.expectedConfig.PKCE, provider.PKCE())
			}

			rawMeta, _ := json.Marshal(provider.Extra())
			expectedMeta, _ := json.Marshal(s.expectedConfig.Extra)
			if !bytes.Equal(rawMeta, expectedMeta) {
				t.Fatalf("Expected PKCE %v, got %v", *s.expectedConfig.PKCE, provider.PKCE())
			}
		})
	}
}
