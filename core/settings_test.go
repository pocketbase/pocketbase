package core_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/auth"
)

func TestSettingsValidate(t *testing.T) {
	s := core.NewSettings()

	// set invalid settings data
	s.Meta.AppName = ""
	s.Logs.MaxDays = -10
	s.Smtp.Enabled = true
	s.Smtp.Host = ""
	s.S3.Enabled = true
	s.S3.Endpoint = "invalid"
	s.AdminAuthToken.Duration = -10
	s.AdminPasswordResetToken.Duration = -10
	s.UserAuthToken.Duration = -10
	s.UserPasswordResetToken.Duration = -10
	s.UserEmailChangeToken.Duration = -10
	s.UserVerificationToken.Duration = -10
	s.EmailAuth.Enabled = true
	s.EmailAuth.MinPasswordLength = -10
	s.GoogleAuth.Enabled = true
	s.GoogleAuth.ClientId = ""
	s.FacebookAuth.Enabled = true
	s.FacebookAuth.ClientId = ""
	s.GithubAuth.Enabled = true
	s.GithubAuth.ClientId = ""
	s.GitlabAuth.Enabled = true
	s.GitlabAuth.ClientId = ""
	s.DiscordAuth.Enabled = true
	s.DiscordAuth.ClientId = ""
	s.TwitterAuth.Enabled = true
	s.TwitterAuth.ClientId = ""

	// check if Validate() is triggering the members validate methods.
	err := s.Validate()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	expectations := []string{
		`"meta":{`,
		`"logs":{`,
		`"smtp":{`,
		`"s3":{`,
		`"adminAuthToken":{`,
		`"adminPasswordResetToken":{`,
		`"userAuthToken":{`,
		`"userPasswordResetToken":{`,
		`"userEmailChangeToken":{`,
		`"userVerificationToken":{`,
		`"emailAuth":{`,
		`"googleAuth":{`,
		`"facebookAuth":{`,
		`"githubAuth":{`,
		`"gitlabAuth":{`,
		`"discordAuth":{`,
	}

	errBytes, _ := json.Marshal(err)
	jsonErr := string(errBytes)
	for _, expected := range expectations {
		if !strings.Contains(jsonErr, expected) {
			t.Errorf("Expected error key %s in %v", expected, jsonErr)
		}
	}
}

func TestSettingsMerge(t *testing.T) {
	s1 := core.NewSettings()
	s1.Meta.AppUrl = "old_app_url"

	s2 := core.NewSettings()
	s2.Meta.AppName = "test"
	s2.Logs.MaxDays = 123
	s2.Smtp.Host = "test"
	s2.Smtp.Enabled = true
	s2.S3.Enabled = true
	s2.S3.Endpoint = "test"
	s2.AdminAuthToken.Duration = 1
	s2.AdminPasswordResetToken.Duration = 2
	s2.UserAuthToken.Duration = 3
	s2.UserPasswordResetToken.Duration = 4
	s2.UserEmailChangeToken.Duration = 5
	s2.UserVerificationToken.Duration = 6
	s2.EmailAuth.Enabled = false
	s2.EmailAuth.MinPasswordLength = 30
	s2.GoogleAuth.Enabled = true
	s2.GoogleAuth.ClientId = "google_test"
	s2.FacebookAuth.Enabled = true
	s2.FacebookAuth.ClientId = "facebook_test"
	s2.GithubAuth.Enabled = true
	s2.GithubAuth.ClientId = "github_test"
	s2.GitlabAuth.Enabled = true
	s2.GitlabAuth.ClientId = "gitlab_test"
	s2.DiscordAuth.Enabled = true
	s2.DiscordAuth.ClientId = "discord_test"
	s2.TwitterAuth.Enabled = true
	s2.TwitterAuth.ClientId = "twitter_test"

	if err := s1.Merge(s2); err != nil {
		t.Fatal(err)
	}

	s1Encoded, err := json.Marshal(s1)
	if err != nil {
		t.Fatal(err)
	}

	s2Encoded, err := json.Marshal(s2)
	if err != nil {
		t.Fatal(err)
	}

	if string(s1Encoded) != string(s2Encoded) {
		t.Fatalf("Expected the same serialization, got %v VS %v", string(s1Encoded), string(s2Encoded))
	}
}

func TestSettingsClone(t *testing.T) {
	s1 := core.NewSettings()

	s2, err := s1.Clone()
	if err != nil {
		t.Fatal(err)
	}

	s1Bytes, err := json.Marshal(s1)
	if err != nil {
		t.Fatal(err)
	}

	s2Bytes, err := json.Marshal(s2)
	if err != nil {
		t.Fatal(err)
	}

	if string(s1Bytes) != string(s2Bytes) {
		t.Fatalf("Expected equivalent serialization, got %v VS %v", string(s1Bytes), string(s2Bytes))
	}

	// verify that it is a deep copy
	s1.Meta.AppName = "new"
	if s1.Meta.AppName == s2.Meta.AppName {
		t.Fatalf("Expected s1 and s2 to have different Meta.AppName, got %s", s1.Meta.AppName)
	}
}

func TestSettingsRedactClone(t *testing.T) {
	s1 := core.NewSettings()
	s1.Meta.AppName = "test123" // control field
	s1.Smtp.Password = "test123"
	s1.Smtp.Tls = true
	s1.S3.Secret = "test123"
	s1.AdminAuthToken.Secret = "test123"
	s1.AdminPasswordResetToken.Secret = "test123"
	s1.UserAuthToken.Secret = "test123"
	s1.UserPasswordResetToken.Secret = "test123"
	s1.UserEmailChangeToken.Secret = "test123"
	s1.UserVerificationToken.Secret = "test123"
	s1.GoogleAuth.ClientSecret = "test123"
	s1.FacebookAuth.ClientSecret = "test123"
	s1.GithubAuth.ClientSecret = "test123"
	s1.GitlabAuth.ClientSecret = "test123"
	s1.DiscordAuth.ClientSecret = "test123"
	s1.TwitterAuth.ClientSecret = "test123"

	s2, err := s1.RedactClone()
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := json.Marshal(s2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"meta":{"appName":"test123","appUrl":"http://localhost:8090","hideControls":false,"senderName":"Support","senderAddress":"support@example.com","verificationTemplate":{"body":"\u003cp\u003eHello,\u003c/p\u003e\n\u003cp\u003eThank you for joining us at {APP_NAME}.\u003c/p\u003e\n\u003cp\u003eClick on the button below to verify your email address.\u003c/p\u003e\n\u003cp\u003e\n  \u003ca class=\"btn\" href=\"{ACTION_URL}\" target=\"_blank\" rel=\"noopener\"\u003eVerify\u003c/a\u003e\n\u003c/p\u003e\n\u003cp\u003e\n  Thanks,\u003cbr/\u003e\n  {APP_NAME} team\n\u003c/p\u003e","subject":"Verify your {APP_NAME} email","actionUrl":"{APP_URL}/_/#/users/confirm-verification/{TOKEN}"},"resetPasswordTemplate":{"body":"\u003cp\u003eHello,\u003c/p\u003e\n\u003cp\u003eClick on the button below to reset your password.\u003c/p\u003e\n\u003cp\u003e\n  \u003ca class=\"btn\" href=\"{ACTION_URL}\" target=\"_blank\" rel=\"noopener\"\u003eReset password\u003c/a\u003e\n\u003c/p\u003e\n\u003cp\u003e\u003ci\u003eIf you didn't ask to reset your password, you can ignore this email.\u003c/i\u003e\u003c/p\u003e\n\u003cp\u003e\n  Thanks,\u003cbr/\u003e\n  {APP_NAME} team\n\u003c/p\u003e","subject":"Reset your {APP_NAME} password","actionUrl":"{APP_URL}/_/#/users/confirm-password-reset/{TOKEN}"},"confirmEmailChangeTemplate":{"body":"\u003cp\u003eHello,\u003c/p\u003e\n\u003cp\u003eClick on the button below to confirm your new email address.\u003c/p\u003e\n\u003cp\u003e\n  \u003ca class=\"btn\" href=\"{ACTION_URL}\" target=\"_blank\" rel=\"noopener\"\u003eConfirm new email\u003c/a\u003e\n\u003c/p\u003e\n\u003cp\u003e\u003ci\u003eIf you didn't ask to change your email address, you can ignore this email.\u003c/i\u003e\u003c/p\u003e\n\u003cp\u003e\n  Thanks,\u003cbr/\u003e\n  {APP_NAME} team\n\u003c/p\u003e","subject":"Confirm your {APP_NAME} new email address","actionUrl":"{APP_URL}/_/#/users/confirm-email-change/{TOKEN}"}},"logs":{"maxDays":7},"smtp":{"enabled":false,"host":"smtp.example.com","port":587,"username":"","password":"******","tls":true},"s3":{"enabled":false,"bucket":"","region":"","endpoint":"","accessKey":"","secret":"******","forcePathStyle":false},"adminAuthToken":{"secret":"******","duration":1209600},"adminPasswordResetToken":{"secret":"******","duration":1800},"userAuthToken":{"secret":"******","duration":1209600},"userPasswordResetToken":{"secret":"******","duration":1800},"userEmailChangeToken":{"secret":"******","duration":1800},"userVerificationToken":{"secret":"******","duration":604800},"emailAuth":{"enabled":true,"exceptDomains":null,"onlyDomains":null,"minPasswordLength":8},"googleAuth":{"enabled":false,"allowRegistrations":true,"clientSecret":"******"},"facebookAuth":{"enabled":false,"allowRegistrations":true,"clientSecret":"******"},"githubAuth":{"enabled":false,"allowRegistrations":true,"clientSecret":"******"},"gitlabAuth":{"enabled":false,"allowRegistrations":true,"clientSecret":"******"},"discordAuth":{"enabled":false,"allowRegistrations":true,"clientSecret":"******"},"twitterAuth":{"enabled":false,"allowRegistrations":true,"clientSecret":"******"}}`

	if encodedStr := string(encoded); encodedStr != expected {
		t.Fatalf("Expected %v, got \n%v", expected, encodedStr)
	}
}

func TestNamedAuthProviderConfigs(t *testing.T) {
	s := core.NewSettings()

	s.GoogleAuth.ClientId = "google_test"
	s.FacebookAuth.ClientId = "facebook_test"
	s.GithubAuth.ClientId = "github_test"
	s.GitlabAuth.ClientId = "gitlab_test"
	s.GitlabAuth.Enabled = true
	s.DiscordAuth.ClientId = "discord_test"
	s.TwitterAuth.ClientId = "twitter_test"

	result := s.NamedAuthProviderConfigs()

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"discord":{"enabled":false,"allowRegistrations":true,"clientId":"discord_test"},"facebook":{"enabled":false,"allowRegistrations":true,"clientId":"facebook_test"},"github":{"enabled":false,"allowRegistrations":true,"clientId":"github_test"},"gitlab":{"enabled":true,"allowRegistrations":true,"clientId":"gitlab_test"},"google":{"enabled":false,"allowRegistrations":true,"clientId":"google_test"},"twitter":{"enabled":false,"allowRegistrations":true,"clientId":"twitter_test"}}`

	if encodedStr := string(encoded); encodedStr != expected {
		t.Fatalf("Expected the same serialization, got %v", encodedStr)
	}
}

func TestTokenConfigValidate(t *testing.T) {
	scenarios := []struct {
		config      core.TokenConfig
		expectError bool
	}{
		// zero values
		{
			core.TokenConfig{},
			true,
		},
		// invalid data
		{
			core.TokenConfig{
				Secret:   strings.Repeat("a", 5),
				Duration: 4,
			},
			true,
		},
		// valid secret but invalid duration
		{
			core.TokenConfig{
				Secret:   strings.Repeat("a", 30),
				Duration: 63072000 + 1,
			},
			true,
		},
		// valid data
		{
			core.TokenConfig{
				Secret:   strings.Repeat("a", 30),
				Duration: 100,
			},
			false,
		},
	}

	for i, scenario := range scenarios {
		result := scenario.config.Validate()

		if result != nil && !scenario.expectError {
			t.Errorf("(%d) Didn't expect error, got %v", i, result)
		}

		if result == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}

func TestSmtpConfigValidate(t *testing.T) {
	scenarios := []struct {
		config      core.SmtpConfig
		expectError bool
	}{
		// zero values (disabled)
		{
			core.SmtpConfig{},
			false,
		},
		// zero values (enabled)
		{
			core.SmtpConfig{Enabled: true},
			true,
		},
		// invalid data
		{
			core.SmtpConfig{
				Enabled: true,
				Host:    "test:test:test",
				Port:    -10,
			},
			true,
		},
		// valid data
		{
			core.SmtpConfig{
				Enabled: true,
				Host:    "example.com",
				Port:    100,
				Tls:     true,
			},
			false,
		},
	}

	for i, scenario := range scenarios {
		result := scenario.config.Validate()

		if result != nil && !scenario.expectError {
			t.Errorf("(%d) Didn't expect error, got %v", i, result)
		}

		if result == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}

func TestS3ConfigValidate(t *testing.T) {
	scenarios := []struct {
		config      core.S3Config
		expectError bool
	}{
		// zero values (disabled)
		{
			core.S3Config{},
			false,
		},
		// zero values (enabled)
		{
			core.S3Config{Enabled: true},
			true,
		},
		// invalid data
		{
			core.S3Config{
				Enabled:  true,
				Endpoint: "test:test:test",
			},
			true,
		},
		// valid data (url endpoint)
		{
			core.S3Config{
				Enabled:   true,
				Endpoint:  "https://localhost:8090",
				Bucket:    "test",
				Region:    "test",
				AccessKey: "test",
				Secret:    "test",
			},
			false,
		},
		// valid data (hostname endpoint)
		{
			core.S3Config{
				Enabled:   true,
				Endpoint:  "example.com",
				Bucket:    "test",
				Region:    "test",
				AccessKey: "test",
				Secret:    "test",
			},
			false,
		},
	}

	for i, scenario := range scenarios {
		result := scenario.config.Validate()

		if result != nil && !scenario.expectError {
			t.Errorf("(%d) Didn't expect error, got %v", i, result)
		}

		if result == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}

func TestMetaConfigValidate(t *testing.T) {
	invalidTemplate := core.EmailTemplate{
		Subject:   "test",
		ActionUrl: "test",
		Body:      "test",
	}

	noPlaceholdersTemplate := core.EmailTemplate{
		Subject:   "test",
		ActionUrl: "http://example.com",
		Body:      "test",
	}

	withPlaceholdersTemplate := core.EmailTemplate{
		Subject:   "test",
		ActionUrl: "http://example.com" + core.EmailPlaceholderToken,
		Body:      "test" + core.EmailPlaceholderActionUrl,
	}

	scenarios := []struct {
		config      core.MetaConfig
		expectError bool
	}{
		// zero values
		{
			core.MetaConfig{},
			true,
		},
		// invalid data
		{
			core.MetaConfig{
				AppName:                    strings.Repeat("a", 300),
				AppUrl:                     "test",
				SenderName:                 strings.Repeat("a", 300),
				SenderAddress:              "invalid_email",
				VerificationTemplate:       invalidTemplate,
				ResetPasswordTemplate:      invalidTemplate,
				ConfirmEmailChangeTemplate: invalidTemplate,
			},
			true,
		},
		// invalid data (missing required placeholders)
		{
			core.MetaConfig{
				AppName:                    "test",
				AppUrl:                     "https://example.com",
				SenderName:                 "test",
				SenderAddress:              "test@example.com",
				VerificationTemplate:       noPlaceholdersTemplate,
				ResetPasswordTemplate:      noPlaceholdersTemplate,
				ConfirmEmailChangeTemplate: noPlaceholdersTemplate,
			},
			true,
		},
		// valid data
		{
			core.MetaConfig{
				AppName:                    "test",
				AppUrl:                     "https://example.com",
				SenderName:                 "test",
				SenderAddress:              "test@example.com",
				VerificationTemplate:       withPlaceholdersTemplate,
				ResetPasswordTemplate:      withPlaceholdersTemplate,
				ConfirmEmailChangeTemplate: withPlaceholdersTemplate,
			},
			false,
		},
	}

	for i, scenario := range scenarios {
		result := scenario.config.Validate()

		if result != nil && !scenario.expectError {
			t.Errorf("(%d) Didn't expect error, got %v", i, result)
		}

		if result == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}

func TestEmailTemplateValidate(t *testing.T) {
	scenarios := []struct {
		emailTemplate  core.EmailTemplate
		expectedErrors []string
	}{
		// require values
		{
			core.EmailTemplate{},
			[]string{"subject", "actionUrl", "body"},
		},
		// missing placeholders
		{
			core.EmailTemplate{
				Subject:   "test",
				ActionUrl: "test",
				Body:      "test",
			},
			[]string{"actionUrl", "body"},
		},
		// valid data
		{
			core.EmailTemplate{
				Subject:   "test",
				ActionUrl: "test" + core.EmailPlaceholderToken,
				Body:      "test" + core.EmailPlaceholderActionUrl,
			},
			[]string{},
		},
	}

	for i, s := range scenarios {
		result := s.emailTemplate.Validate()

		// parse errors
		errs, ok := result.(validation.Errors)
		if !ok && result != nil {
			t.Errorf("(%d) Failed to parse errors %v", i, result)
			continue
		}

		// check errors
		if len(errs) > len(s.expectedErrors) {
			t.Errorf("(%d) Expected error keys %v, got %v", i, s.expectedErrors, errs)
		}
		for _, k := range s.expectedErrors {
			if _, ok := errs[k]; !ok {
				t.Errorf("(%d) Missing expected error key %q in %v", i, k, errs)
			}
		}
	}
}

func TestEmailTemplateResolve(t *testing.T) {
	allPlaceholders := core.EmailPlaceholderActionUrl + core.EmailPlaceholderToken + core.EmailPlaceholderAppName + core.EmailPlaceholderAppUrl

	scenarios := []struct {
		emailTemplate     core.EmailTemplate
		expectedSubject   string
		expectedBody      string
		expectedActionUrl string
	}{
		// no placeholders
		{
			emailTemplate: core.EmailTemplate{
				Subject:   "subject:",
				Body:      "body:",
				ActionUrl: "/actionUrl////",
			},
			expectedSubject:   "subject:",
			expectedActionUrl: "/actionUrl/",
			expectedBody:      "body:",
		},
		// with placeholders
		{
			emailTemplate: core.EmailTemplate{
				ActionUrl: "/actionUrl////" + allPlaceholders,
				Subject:   "subject:" + allPlaceholders,
				Body:      "body:" + allPlaceholders,
			},
			expectedActionUrl: fmt.Sprintf(
				"/actionUrl/%%7BACTION_URL%%7D%s%s%s",
				"token_test",
				"name_test",
				"url_test",
			),
			expectedSubject: fmt.Sprintf(
				"subject:%s%s%s%s",
				core.EmailPlaceholderActionUrl,
				core.EmailPlaceholderToken,
				"name_test",
				"url_test",
			),
			expectedBody: fmt.Sprintf(
				"body:%s%s%s%s",
				fmt.Sprintf(
					"/actionUrl/%%7BACTION_URL%%7D%s%s%s",
					"token_test",
					"name_test",
					"url_test",
				),
				"token_test",
				"name_test",
				"url_test",
			),
		},
	}

	for i, s := range scenarios {
		subject, body, actionUrl := s.emailTemplate.Resolve("name_test", "url_test", "token_test")

		if s.expectedSubject != subject {
			t.Errorf("(%d) Expected subject %q got %q", i, s.expectedSubject, subject)
		}

		if s.expectedBody != body {
			t.Errorf("(%d) Expected body \n%v got \n%v", i, s.expectedBody, body)
		}

		if s.expectedActionUrl != actionUrl {
			t.Errorf("(%d) Expected actionUrl \n%v got \n%v", i, s.expectedActionUrl, actionUrl)
		}
	}
}

func TestLogsConfigValidate(t *testing.T) {
	scenarios := []struct {
		config      core.LogsConfig
		expectError bool
	}{
		// zero values
		{
			core.LogsConfig{},
			false,
		},
		// invalid data
		{
			core.LogsConfig{MaxDays: -10},
			true,
		},
		// valid data
		{
			core.LogsConfig{MaxDays: 1},
			false,
		},
	}

	for i, scenario := range scenarios {
		result := scenario.config.Validate()

		if result != nil && !scenario.expectError {
			t.Errorf("(%d) Didn't expect error, got %v", i, result)
		}

		if result == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}

func TestAuthProviderConfigValidate(t *testing.T) {
	scenarios := []struct {
		config      core.AuthProviderConfig
		expectError bool
	}{
		// zero values (disabled)
		{
			core.AuthProviderConfig{},
			false,
		},
		// zero values (enabled)
		{
			core.AuthProviderConfig{Enabled: true},
			true,
		},
		// invalid data
		{
			core.AuthProviderConfig{
				Enabled:      true,
				ClientId:     "",
				ClientSecret: "",
				AuthUrl:      "test",
				TokenUrl:     "test",
				UserApiUrl:   "test",
			},
			true,
		},
		// valid data (only the required)
		{
			core.AuthProviderConfig{
				Enabled:      true,
				ClientId:     "test",
				ClientSecret: "test",
			},
			false,
		},
		// valid data (fill all fields)
		{
			core.AuthProviderConfig{
				Enabled:      true,
				ClientId:     "test",
				ClientSecret: "test",
				AuthUrl:      "https://example.com",
				TokenUrl:     "https://example.com",
				UserApiUrl:   "https://example.com",
			},
			false,
		},
	}

	for i, scenario := range scenarios {
		result := scenario.config.Validate()

		if result != nil && !scenario.expectError {
			t.Errorf("(%d) Didn't expect error, got %v", i, result)
		}

		if result == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}

func TestAuthProviderConfigSetupProvider(t *testing.T) {
	provider := auth.NewGithubProvider()

	// disabled config
	c1 := core.AuthProviderConfig{Enabled: false}
	if err := c1.SetupProvider(provider); err == nil {
		t.Errorf("Expected error, got nil")
	}

	c2 := core.AuthProviderConfig{
		Enabled:      true,
		ClientId:     "test_ClientId",
		ClientSecret: "test_ClientSecret",
		AuthUrl:      "test_AuthUrl",
		UserApiUrl:   "test_UserApiUrl",
		TokenUrl:     "test_TokenUrl",
	}
	if err := c2.SetupProvider(provider); err != nil {
		t.Error(err)
	}
	encoded, _ := json.Marshal(c2)
	expected := `{"enabled":true,"allowRegistrations":false,"clientId":"test_ClientId","clientSecret":"test_ClientSecret","authUrl":"test_AuthUrl","tokenUrl":"test_TokenUrl","userApiUrl":"test_UserApiUrl"}`
	if string(encoded) != expected {
		t.Errorf("Expected %s, got %s", expected, string(encoded))
	}
}

func TestEmailAuthConfigValidate(t *testing.T) {
	scenarios := []struct {
		config      core.EmailAuthConfig
		expectError bool
	}{
		// zero values (disabled)
		{
			core.EmailAuthConfig{},
			false,
		},
		// zero values (enabled)
		{
			core.EmailAuthConfig{Enabled: true},
			true,
		},
		// invalid data (only the required)
		{
			core.EmailAuthConfig{
				Enabled:           true,
				MinPasswordLength: 4,
			},
			true,
		},
		// valid data (only the required)
		{
			core.EmailAuthConfig{
				Enabled:           true,
				MinPasswordLength: 5,
			},
			false,
		},
		// invalid data (both OnlyDomains and ExceptDomains set)
		{
			core.EmailAuthConfig{
				Enabled:           true,
				MinPasswordLength: 5,
				OnlyDomains:       []string{"example.com", "test.com"},
				ExceptDomains:     []string{"example.com", "test.com"},
			},
			true,
		},
		// valid data (only onlyDomains set)
		{
			core.EmailAuthConfig{
				Enabled:           true,
				MinPasswordLength: 5,
				OnlyDomains:       []string{"example.com", "test.com"},
			},
			false,
		},
		// valid data (only exceptDomains set)
		{
			core.EmailAuthConfig{
				Enabled:           true,
				MinPasswordLength: 5,
				ExceptDomains:     []string{"example.com", "test.com"},
			},
			false,
		},
	}

	for i, scenario := range scenarios {
		result := scenario.config.Validate()

		if result != nil && !scenario.expectError {
			t.Errorf("(%d) Didn't expect error, got %v", i, result)
		}

		if result == nil && scenario.expectError {
			t.Errorf("(%d) Expected error, got nil", i)
		}
	}
}
