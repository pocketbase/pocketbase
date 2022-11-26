package settings_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tools/auth"
)

func TestSettingsValidate(t *testing.T) {
	s := settings.New()

	// set invalid settings data
	s.Meta.AppName = ""
	s.Logs.MaxDays = -10
	s.Smtp.Enabled = true
	s.Smtp.Host = ""
	s.S3.Enabled = true
	s.S3.Endpoint = "invalid"
	s.AdminAuthToken.Duration = -10
	s.AdminPasswordResetToken.Duration = -10
	s.RecordAuthToken.Duration = -10
	s.RecordPasswordResetToken.Duration = -10
	s.RecordEmailChangeToken.Duration = -10
	s.RecordVerificationToken.Duration = -10
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
	s.MicrosoftAuth.Enabled = true
	s.MicrosoftAuth.ClientId = ""
	s.SpotifyAuth.Enabled = true
	s.SpotifyAuth.ClientId = ""
	s.KakaoAuth.Enabled = true
	s.KakaoAuth.ClientId = ""
	s.TwitchAuth.Enabled = true
	s.TwitchAuth.ClientId = ""

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
		`"recordAuthToken":{`,
		`"recordPasswordResetToken":{`,
		`"recordEmailChangeToken":{`,
		`"recordVerificationToken":{`,
		`"googleAuth":{`,
		`"facebookAuth":{`,
		`"githubAuth":{`,
		`"gitlabAuth":{`,
		`"discordAuth":{`,
		`"twitterAuth":{`,
		`"microsoftAuth":{`,
		`"spotifyAuth":{`,
		`"kakaoAuth":{`,
		`"twitchAuth":{`,
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
	s1 := settings.New()
	s1.Meta.AppUrl = "old_app_url"

	s2 := settings.New()
	s2.Meta.AppName = "test"
	s2.Logs.MaxDays = 123
	s2.Smtp.Host = "test"
	s2.Smtp.Enabled = true
	s2.S3.Enabled = true
	s2.S3.Endpoint = "test"
	s2.AdminAuthToken.Duration = 1
	s2.AdminPasswordResetToken.Duration = 2
	s2.RecordAuthToken.Duration = 3
	s2.RecordPasswordResetToken.Duration = 4
	s2.RecordEmailChangeToken.Duration = 5
	s2.RecordVerificationToken.Duration = 6
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
	s2.MicrosoftAuth.Enabled = true
	s2.MicrosoftAuth.ClientId = "microsoft_test"
	s2.SpotifyAuth.Enabled = true
	s2.SpotifyAuth.ClientId = "spotify_test"
	s2.KakaoAuth.Enabled = true
	s2.KakaoAuth.ClientId = "kakao_test"
	s2.TwitchAuth.Enabled = true
	s2.TwitchAuth.ClientId = "twitch_test"

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
	s1 := settings.New()

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
	s1 := settings.New()
	s1.Meta.AppName = "test123" // control field
	s1.Smtp.Password = "test123"
	s1.Smtp.Tls = true
	s1.S3.Secret = "test123"
	s1.AdminAuthToken.Secret = "test123"
	s1.AdminPasswordResetToken.Secret = "test123"
	s1.RecordAuthToken.Secret = "test123"
	s1.RecordPasswordResetToken.Secret = "test123"
	s1.RecordEmailChangeToken.Secret = "test123"
	s1.RecordVerificationToken.Secret = "test123"
	s1.GoogleAuth.ClientSecret = "test123"
	s1.FacebookAuth.ClientSecret = "test123"
	s1.GithubAuth.ClientSecret = "test123"
	s1.GitlabAuth.ClientSecret = "test123"
	s1.DiscordAuth.ClientSecret = "test123"
	s1.TwitterAuth.ClientSecret = "test123"
	s1.MicrosoftAuth.ClientSecret = "test123"
	s1.SpotifyAuth.ClientSecret = "test123"
	s1.KakaoAuth.ClientSecret = "test123"
	s1.TwitchAuth.ClientSecret = "test123"

	s2, err := s1.RedactClone()
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := json.Marshal(s2)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"meta":{"appName":"test123","appUrl":"http://localhost:8090","hideControls":false,"senderName":"Support","senderAddress":"support@example.com","verificationTemplate":{"body":"\u003cp\u003eHello,\u003c/p\u003e\n\u003cp\u003eThank you for joining us at {APP_NAME}.\u003c/p\u003e\n\u003cp\u003eClick on the button below to verify your email address.\u003c/p\u003e\n\u003cp\u003e\n  \u003ca class=\"btn\" href=\"{ACTION_URL}\" target=\"_blank\" rel=\"noopener\"\u003eVerify\u003c/a\u003e\n\u003c/p\u003e\n\u003cp\u003e\n  Thanks,\u003cbr/\u003e\n  {APP_NAME} team\n\u003c/p\u003e","subject":"Verify your {APP_NAME} email","actionUrl":"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}"},"resetPasswordTemplate":{"body":"\u003cp\u003eHello,\u003c/p\u003e\n\u003cp\u003eClick on the button below to reset your password.\u003c/p\u003e\n\u003cp\u003e\n  \u003ca class=\"btn\" href=\"{ACTION_URL}\" target=\"_blank\" rel=\"noopener\"\u003eReset password\u003c/a\u003e\n\u003c/p\u003e\n\u003cp\u003e\u003ci\u003eIf you didn't ask to reset your password, you can ignore this email.\u003c/i\u003e\u003c/p\u003e\n\u003cp\u003e\n  Thanks,\u003cbr/\u003e\n  {APP_NAME} team\n\u003c/p\u003e","subject":"Reset your {APP_NAME} password","actionUrl":"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}"},"confirmEmailChangeTemplate":{"body":"\u003cp\u003eHello,\u003c/p\u003e\n\u003cp\u003eClick on the button below to confirm your new email address.\u003c/p\u003e\n\u003cp\u003e\n  \u003ca class=\"btn\" href=\"{ACTION_URL}\" target=\"_blank\" rel=\"noopener\"\u003eConfirm new email\u003c/a\u003e\n\u003c/p\u003e\n\u003cp\u003e\u003ci\u003eIf you didn't ask to change your email address, you can ignore this email.\u003c/i\u003e\u003c/p\u003e\n\u003cp\u003e\n  Thanks,\u003cbr/\u003e\n  {APP_NAME} team\n\u003c/p\u003e","subject":"Confirm your {APP_NAME} new email address","actionUrl":"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}"}},"logs":{"maxDays":5},"smtp":{"enabled":false,"host":"smtp.example.com","port":587,"username":"","password":"******","tls":true},"s3":{"enabled":false,"bucket":"","region":"","endpoint":"","accessKey":"","secret":"******","forcePathStyle":false},"adminAuthToken":{"secret":"******","duration":1209600},"adminPasswordResetToken":{"secret":"******","duration":1800},"recordAuthToken":{"secret":"******","duration":1209600},"recordPasswordResetToken":{"secret":"******","duration":1800},"recordEmailChangeToken":{"secret":"******","duration":1800},"recordVerificationToken":{"secret":"******","duration":604800},"emailAuth":{"enabled":false,"exceptDomains":null,"onlyDomains":null,"minPasswordLength":0},"googleAuth":{"enabled":false,"clientSecret":"******"},"facebookAuth":{"enabled":false,"clientSecret":"******"},"githubAuth":{"enabled":false,"clientSecret":"******"},"gitlabAuth":{"enabled":false,"clientSecret":"******"},"discordAuth":{"enabled":false,"clientSecret":"******"},"twitterAuth":{"enabled":false,"clientSecret":"******"},"microsoftAuth":{"enabled":false,"clientSecret":"******"},"spotifyAuth":{"enabled":false,"clientSecret":"******"},"kakaoAuth":{"enabled":false,"clientSecret":"******"},"twitchAuth":{"enabled":false,"clientSecret":"******"}}`

	if encodedStr := string(encoded); encodedStr != expected {
		t.Fatalf("Expected\n%v\ngot\n%v", expected, encodedStr)
	}
}

func TestNamedAuthProviderConfigs(t *testing.T) {
	s := settings.New()

	s.GoogleAuth.ClientId = "google_test"
	s.FacebookAuth.ClientId = "facebook_test"
	s.GithubAuth.ClientId = "github_test"
	s.GitlabAuth.ClientId = "gitlab_test"
	s.GitlabAuth.Enabled = true
	s.DiscordAuth.ClientId = "discord_test"
	s.TwitterAuth.ClientId = "twitter_test"
	s.MicrosoftAuth.ClientId = "microsoft_test"
	s.SpotifyAuth.ClientId = "spotify_test"
	s.KakaoAuth.ClientId = "kakao_test"
	s.TwitchAuth.ClientId = "twitch_test"

	result := s.NamedAuthProviderConfigs()

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	encodedStr := string(encoded)

	expectedParts := []string{
		`"discord":{"enabled":false,"clientId":"discord_test"}`,
		`"facebook":{"enabled":false,"clientId":"facebook_test"}`,
		`"github":{"enabled":false,"clientId":"github_test"}`,
		`"gitlab":{"enabled":true,"clientId":"gitlab_test"}`,
		`"google":{"enabled":false,"clientId":"google_test"}`,
		`"microsoft":{"enabled":false,"clientId":"microsoft_test"}`,
		`"spotify":{"enabled":false,"clientId":"spotify_test"}`,
		`"twitter":{"enabled":false,"clientId":"twitter_test"}`,
		`"kakao":{"enabled":false,"clientId":"kakao_test"}`,
		`"twitch":{"enabled":false,"clientId":"twitch_test"}`,
	}
	for _, p := range expectedParts {
		if !strings.Contains(encodedStr, p) {
			t.Fatalf("Expected \n%s \nin \n%s", p, encodedStr)
		}
	}
}

func TestTokenConfigValidate(t *testing.T) {
	scenarios := []struct {
		config      settings.TokenConfig
		expectError bool
	}{
		// zero values
		{
			settings.TokenConfig{},
			true,
		},
		// invalid data
		{
			settings.TokenConfig{
				Secret:   strings.Repeat("a", 5),
				Duration: 4,
			},
			true,
		},
		// valid secret but invalid duration
		{
			settings.TokenConfig{
				Secret:   strings.Repeat("a", 30),
				Duration: 63072000 + 1,
			},
			true,
		},
		// valid data
		{
			settings.TokenConfig{
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
		config      settings.SmtpConfig
		expectError bool
	}{
		// zero values (disabled)
		{
			settings.SmtpConfig{},
			false,
		},
		// zero values (enabled)
		{
			settings.SmtpConfig{Enabled: true},
			true,
		},
		// invalid data
		{
			settings.SmtpConfig{
				Enabled: true,
				Host:    "test:test:test",
				Port:    -10,
			},
			true,
		},
		// valid data
		{
			settings.SmtpConfig{
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
		config      settings.S3Config
		expectError bool
	}{
		// zero values (disabled)
		{
			settings.S3Config{},
			false,
		},
		// zero values (enabled)
		{
			settings.S3Config{Enabled: true},
			true,
		},
		// invalid data
		{
			settings.S3Config{
				Enabled:  true,
				Endpoint: "test:test:test",
			},
			true,
		},
		// valid data (url endpoint)
		{
			settings.S3Config{
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
			settings.S3Config{
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
	invalidTemplate := settings.EmailTemplate{
		Subject:   "test",
		ActionUrl: "test",
		Body:      "test",
	}

	noPlaceholdersTemplate := settings.EmailTemplate{
		Subject:   "test",
		ActionUrl: "http://example.com",
		Body:      "test",
	}

	withPlaceholdersTemplate := settings.EmailTemplate{
		Subject:   "test",
		ActionUrl: "http://example.com" + settings.EmailPlaceholderToken,
		Body:      "test" + settings.EmailPlaceholderActionUrl,
	}

	scenarios := []struct {
		config      settings.MetaConfig
		expectError bool
	}{
		// zero values
		{
			settings.MetaConfig{},
			true,
		},
		// invalid data
		{
			settings.MetaConfig{
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
			settings.MetaConfig{
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
			settings.MetaConfig{
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
		emailTemplate  settings.EmailTemplate
		expectedErrors []string
	}{
		// require values
		{
			settings.EmailTemplate{},
			[]string{"subject", "actionUrl", "body"},
		},
		// missing placeholders
		{
			settings.EmailTemplate{
				Subject:   "test",
				ActionUrl: "test",
				Body:      "test",
			},
			[]string{"actionUrl", "body"},
		},
		// valid data
		{
			settings.EmailTemplate{
				Subject:   "test",
				ActionUrl: "test" + settings.EmailPlaceholderToken,
				Body:      "test" + settings.EmailPlaceholderActionUrl,
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
	allPlaceholders := settings.EmailPlaceholderActionUrl + settings.EmailPlaceholderToken + settings.EmailPlaceholderAppName + settings.EmailPlaceholderAppUrl

	scenarios := []struct {
		emailTemplate     settings.EmailTemplate
		expectedSubject   string
		expectedBody      string
		expectedActionUrl string
	}{
		// no placeholders
		{
			emailTemplate: settings.EmailTemplate{
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
			emailTemplate: settings.EmailTemplate{
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
				settings.EmailPlaceholderActionUrl,
				settings.EmailPlaceholderToken,
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
		config      settings.LogsConfig
		expectError bool
	}{
		// zero values
		{
			settings.LogsConfig{},
			false,
		},
		// invalid data
		{
			settings.LogsConfig{MaxDays: -10},
			true,
		},
		// valid data
		{
			settings.LogsConfig{MaxDays: 1},
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
		config      settings.AuthProviderConfig
		expectError bool
	}{
		// zero values (disabled)
		{
			settings.AuthProviderConfig{},
			false,
		},
		// zero values (enabled)
		{
			settings.AuthProviderConfig{Enabled: true},
			true,
		},
		// invalid data
		{
			settings.AuthProviderConfig{
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
			settings.AuthProviderConfig{
				Enabled:      true,
				ClientId:     "test",
				ClientSecret: "test",
			},
			false,
		},
		// valid data (fill all fields)
		{
			settings.AuthProviderConfig{
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
	c1 := settings.AuthProviderConfig{Enabled: false}
	if err := c1.SetupProvider(provider); err == nil {
		t.Errorf("Expected error, got nil")
	}

	c2 := settings.AuthProviderConfig{
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

	if provider.ClientId() != c2.ClientId {
		t.Fatalf("Expected ClientId %s, got %s", c2.ClientId, provider.ClientId())
	}

	if provider.ClientSecret() != c2.ClientSecret {
		t.Fatalf("Expected ClientSecret %s, got %s", c2.ClientSecret, provider.ClientSecret())
	}

	if provider.AuthUrl() != c2.AuthUrl {
		t.Fatalf("Expected AuthUrl %s, got %s", c2.AuthUrl, provider.AuthUrl())
	}

	if provider.UserApiUrl() != c2.UserApiUrl {
		t.Fatalf("Expected UserApiUrl %s, got %s", c2.UserApiUrl, provider.UserApiUrl())
	}

	if provider.TokenUrl() != c2.TokenUrl {
		t.Fatalf("Expected TokenUrl %s, got %s", c2.TokenUrl, provider.TokenUrl())
	}
}
