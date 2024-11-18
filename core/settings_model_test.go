package core_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

func TestSettingsDelete(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	err := app.Delete(app.Settings())
	if err == nil {
		t.Fatal("Exected settings delete to fail")
	}
}

func TestSettingsMerge(t *testing.T) {
	s1 := &core.Settings{}
	s1.Meta.AppURL = "app_url" // should be unset

	s2 := &core.Settings{}
	s2.Meta.AppName = "test"
	s2.Logs.MaxDays = 123
	s2.SMTP.Host = "test"
	s2.SMTP.Enabled = true
	s2.S3.Enabled = true
	s2.S3.Endpoint = "test"
	s2.Backups.Cron = "* * * * *"
	s2.Batch.Timeout = 15

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
		t.Fatalf("Expected the same serialization, got\n%v\nVS\n%v", string(s1Encoded), string(s2Encoded))
	}
}

func TestSettingsClone(t *testing.T) {
	s1 := &core.Settings{}
	s1.Meta.AppName = "test_name"

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
	s2.Meta.AppName = "new_test_name"
	if s1.Meta.AppName == s2.Meta.AppName {
		t.Fatalf("Expected s1 and s2 to have different Meta.AppName, got %s", s1.Meta.AppName)
	}
}

func TestSettingsMarshalJSON(t *testing.T) {
	settings := &core.Settings{}

	// control fields
	settings.Meta.AppName = "test123"
	settings.SMTP.Username = "abc"

	// secrets
	testSecret := "test_secret"
	settings.SMTP.Password = testSecret
	settings.S3.Secret = testSecret
	settings.Backups.S3.Secret = testSecret

	raw, err := json.Marshal(settings)
	if err != nil {
		t.Fatal(err)
	}
	rawStr := string(raw)

	expected := `{"smtp":{"enabled":false,"port":0,"host":"","username":"abc","authMethod":"","tls":false,"localName":""},"backups":{"cron":"","cronMaxKeep":0,"s3":{"enabled":false,"bucket":"","region":"","endpoint":"","accessKey":"","forcePathStyle":false}},"s3":{"enabled":false,"bucket":"","region":"","endpoint":"","accessKey":"","forcePathStyle":false},"meta":{"appName":"test123","appURL":"","senderName":"","senderAddress":"","hideControls":false},"rateLimits":{"rules":[],"enabled":false},"trustedProxy":{"headers":[],"useLeftmostIP":false},"batch":{"enabled":false,"maxRequests":0,"timeout":0,"maxBodySize":0},"logs":{"maxDays":0,"minLevel":0,"logIP":false,"logAuthId":false}}`

	if rawStr != expected {
		t.Fatalf("Expected\n%v\ngot\n%v", expected, rawStr)
	}
}

func TestSettingsValidate(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	s := app.Settings()

	// set invalid settings data
	s.Meta.AppName = ""
	s.Logs.MaxDays = -10
	s.SMTP.Enabled = true
	s.SMTP.Host = ""
	s.S3.Enabled = true
	s.S3.Endpoint = "invalid"
	s.Backups.Cron = "invalid"
	s.Backups.CronMaxKeep = -10
	s.Batch.Enabled = true
	s.Batch.MaxRequests = -1
	s.Batch.Timeout = -1
	s.RateLimits.Enabled = true
	s.RateLimits.Rules = nil

	// check if Validate() is triggering the members validate methods.
	err := app.Validate(s)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	expectations := []string{
		`"meta":{`,
		`"logs":{`,
		`"smtp":{`,
		`"s3":{`,
		`"backups":{`,
		`"batch":{`,
		`"rateLimits":{`,
	}

	errBytes, _ := json.Marshal(err)
	jsonErr := string(errBytes)
	for _, expected := range expectations {
		if !strings.Contains(jsonErr, expected) {
			t.Errorf("Expected error key %s in %v", expected, jsonErr)
		}
	}
}

func TestMetaConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.MetaConfig
		expectedErrors []string
	}{
		{
			"zero values",
			core.MetaConfig{},
			[]string{
				"appName",
				"appURL",
				"senderName",
				"senderAddress",
			},
		},
		{
			"invalid data",
			core.MetaConfig{
				AppName:       strings.Repeat("a", 300),
				AppURL:        "test",
				SenderName:    strings.Repeat("a", 300),
				SenderAddress: "invalid_email",
			},
			[]string{
				"appName",
				"appURL",
				"senderName",
				"senderAddress",
			},
		},
		{
			"valid data",
			core.MetaConfig{
				AppName:       "test",
				AppURL:        "https://example.com",
				SenderName:    "test",
				SenderAddress: "test@example.com",
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

func TestLogsConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.LogsConfig
		expectedErrors []string
	}{
		{
			"zero values",
			core.LogsConfig{},
			[]string{},
		},
		{
			"invalid data",
			core.LogsConfig{MaxDays: -1},
			[]string{"maxDays"},
		},
		{
			"valid data",
			core.LogsConfig{MaxDays: 2},
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

func TestSMTPConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.SMTPConfig
		expectedErrors []string
	}{
		{
			"zero values (disabled)",
			core.SMTPConfig{},
			[]string{},
		},
		{
			"zero values (enabled)",
			core.SMTPConfig{Enabled: true},
			[]string{"host", "port"},
		},
		{
			"invalid data",
			core.SMTPConfig{
				Enabled:    true,
				Host:       "test:test:test",
				Port:       -10,
				LocalName:  "invalid!",
				AuthMethod: "invalid",
			},
			[]string{"host", "port", "authMethod", "localName"},
		},
		{
			"valid data (no explicit auth method and localName)",
			core.SMTPConfig{
				Enabled: true,
				Host:    "example.com",
				Port:    100,
				TLS:     true,
			},
			[]string{},
		},
		{
			"valid data (explicit auth method and localName)",
			core.SMTPConfig{
				Enabled:    true,
				Host:       "example.com",
				Port:       100,
				AuthMethod: mailer.SMTPAuthLogin,
				LocalName:  "example.com",
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

func TestS3ConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.S3Config
		expectedErrors []string
	}{
		{
			"zero values (disabled)",
			core.S3Config{},
			[]string{},
		},
		{
			"zero values (enabled)",
			core.S3Config{Enabled: true},
			[]string{
				"bucket",
				"region",
				"endpoint",
				"accessKey",
				"secret",
			},
		},
		{
			"invalid data",
			core.S3Config{
				Enabled:  true,
				Endpoint: "test:test:test",
			},
			[]string{
				"bucket",
				"region",
				"endpoint",
				"accessKey",
				"secret",
			},
		},
		{
			"valid data (url endpoint)",
			core.S3Config{
				Enabled:   true,
				Endpoint:  "https://localhost:8090",
				Bucket:    "test",
				Region:    "test",
				AccessKey: "test",
				Secret:    "test",
			},
			[]string{},
		},
		{
			"valid data (hostname endpoint)",
			core.S3Config{
				Enabled:   true,
				Endpoint:  "example.com",
				Bucket:    "test",
				Region:    "test",
				AccessKey: "test",
				Secret:    "test",
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

func TestBackupsConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.BackupsConfig
		expectedErrors []string
	}{
		{
			"zero value",
			core.BackupsConfig{},
			[]string{},
		},
		{
			"invalid cron",
			core.BackupsConfig{
				Cron:        "invalid",
				CronMaxKeep: 0,
			},
			[]string{"cron", "cronMaxKeep"},
		},
		{
			"invalid enabled S3",
			core.BackupsConfig{
				S3: core.S3Config{
					Enabled: true,
				},
			},
			[]string{"s3"},
		},
		{
			"valid data",
			core.BackupsConfig{
				S3: core.S3Config{
					Enabled:   true,
					Endpoint:  "example.com",
					Bucket:    "test",
					Region:    "test",
					AccessKey: "test",
					Secret:    "test",
				},
				Cron:        "*/10 * * * *",
				CronMaxKeep: 1,
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

func TestBatchConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.BatchConfig
		expectedErrors []string
	}{
		{
			"zero value",
			core.BatchConfig{},
			[]string{},
		},
		{
			"zero value (enabled)",
			core.BatchConfig{Enabled: true},
			[]string{"maxRequests", "timeout"},
		},
		{
			"invalid data (negative values)",
			core.BatchConfig{
				MaxRequests: -1,
				Timeout:     -1,
				MaxBodySize: -1,
			},
			[]string{"maxRequests", "timeout", "maxBodySize"},
		},
		{
			"min fields valid data",
			core.BatchConfig{
				Enabled:     true,
				MaxRequests: 1,
				Timeout:     1,
			},
			[]string{},
		},
		{
			"all fields valid data",
			core.BatchConfig{
				Enabled:     true,
				MaxRequests: 10,
				Timeout:     1,
				MaxBodySize: 1,
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

func TestRateLimitsConfigValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.RateLimitsConfig
		expectedErrors []string
	}{
		{
			"zero value (disabled)",
			core.RateLimitsConfig{},
			[]string{},
		},
		{
			"zero value (enabled)",
			core.RateLimitsConfig{Enabled: true},
			[]string{"rules"},
		},
		{
			"invalid data",
			core.RateLimitsConfig{
				Enabled: true,
				Rules: []core.RateLimitRule{
					{
						Label:       "/123abc/",
						Duration:    1,
						MaxRequests: 2,
					},
					{
						Label:       "!abc",
						Duration:    -1,
						MaxRequests: -1,
					},
				},
			},
			[]string{"rules"},
		},
		{
			"valid data",
			core.RateLimitsConfig{
				Enabled: true,
				Rules: []core.RateLimitRule{
					{
						Label:       "123_abc",
						Duration:    1,
						MaxRequests: 2,
					},
					{
						Label:       "/456-abc",
						Duration:    1,
						MaxRequests: 2,
					},
				},
			},
			[]string{},
		},
		{
			"duplicated rules with the same audience",
			core.RateLimitsConfig{
				Enabled: true,
				Rules: []core.RateLimitRule{
					{
						Label:       "/a",
						Duration:    1,
						MaxRequests: 2,
					},
					{
						Label:       "/a",
						Duration:    2,
						MaxRequests: 3,
					},
				},
			},
			[]string{"rules"},
		},
		{
			"duplicated rule with conflicting audience (A)",
			core.RateLimitsConfig{
				Enabled: true,
				Rules: []core.RateLimitRule{
					{
						Label:       "/a",
						Duration:    1,
						MaxRequests: 2,
					},
					{
						Label:       "/a",
						Duration:    1,
						MaxRequests: 2,
						Audience:    core.RateLimitRuleAudienceGuest,
					},
				},
			},
			[]string{"rules"},
		},
		{
			"duplicated rule with conflicting audience (B)",
			core.RateLimitsConfig{
				Enabled: true,
				Rules: []core.RateLimitRule{
					{
						Label:       "/a",
						Duration:    1,
						MaxRequests: 2,
						Audience:    core.RateLimitRuleAudienceAuth,
					},
					{
						Label:       "/a",
						Duration:    1,
						MaxRequests: 2,
					},
				},
			},
			[]string{"rules"},
		},
		{
			"duplicated rule with non-conflicting audience",
			core.RateLimitsConfig{
				Enabled: true,
				Rules: []core.RateLimitRule{
					{
						Label:       "/a",
						Duration:    1,
						MaxRequests: 2,
						Audience:    core.RateLimitRuleAudienceAuth,
					},
					{
						Label:       "/a",
						Duration:    1,
						MaxRequests: 2,
						Audience:    core.RateLimitRuleAudienceGuest,
					},
				},
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

func TestRateLimitsFindRateLimitRule(t *testing.T) {
	limits := core.RateLimitsConfig{
		Rules: []core.RateLimitRule{
			{Label: "abc"},
			{Label: "def", Audience: core.RateLimitRuleAudienceGuest},
			{Label: "/test/a", Audience: core.RateLimitRuleAudienceGuest},
			{Label: "POST /test/a"},
			{Label: "/test/a/", Audience: core.RateLimitRuleAudienceAuth},
			{Label: "POST /test/a/"},
		},
	}

	scenarios := []struct {
		labels   []string
		audience []string
		expected string
	}{
		{[]string{}, []string{}, ""},
		{[]string{"missing"}, []string{}, ""},
		{[]string{"abc"}, []string{}, "abc"},
		{[]string{"abc"}, []string{core.RateLimitRuleAudienceGuest}, ""},
		{[]string{"abc"}, []string{core.RateLimitRuleAudienceAuth}, ""},
		{[]string{"def"}, []string{core.RateLimitRuleAudienceGuest}, "def"},
		{[]string{"def"}, []string{core.RateLimitRuleAudienceAuth}, ""},
		{[]string{"/test"}, []string{}, ""},
		{[]string{"/test/a"}, []string{}, "/test/a"},
		{[]string{"/test/a"}, []string{core.RateLimitRuleAudienceAuth}, "/test/a/"},
		{[]string{"/test/a"}, []string{core.RateLimitRuleAudienceGuest}, "/test/a"},
		{[]string{"GET /test/a"}, []string{}, ""},
		{[]string{"POST /test/a"}, []string{}, "POST /test/a"},
		{[]string{"/test/a/b/c"}, []string{}, "/test/a/"},
		{[]string{"/test/a/b/c"}, []string{core.RateLimitRuleAudienceAuth}, "/test/a/"},
		{[]string{"/test/a/b/c"}, []string{core.RateLimitRuleAudienceGuest}, ""},
		{[]string{"GET /test/a/b/c"}, []string{}, ""},
		{[]string{"POST /test/a/b/c"}, []string{}, "POST /test/a/"},
		{[]string{"/test/a", "abc"}, []string{}, "/test/a"}, // priority checks
	}

	for _, s := range scenarios {
		t.Run(strings.Join(s.labels, "_")+":"+strings.Join(s.audience, "_"), func(t *testing.T) {
			rule, ok := limits.FindRateLimitRule(s.labels, s.audience...)

			hasLabel := rule.Label != ""
			if hasLabel != ok {
				t.Fatalf("Expected hasLabel %v, got %v", hasLabel, ok)
			}

			if rule.Label != s.expected {
				t.Fatalf("Expected rule with label %q, got %q", s.expected, rule.Label)
			}
		})
	}
}

func TestRateLimitRuleValidate(t *testing.T) {
	scenarios := []struct {
		name           string
		config         core.RateLimitRule
		expectedErrors []string
	}{
		{
			"zero value",
			core.RateLimitRule{},
			[]string{"label", "duration", "maxRequests"},
		},
		{
			"invalid data",
			core.RateLimitRule{
				Label:       "@abc",
				Duration:    -1,
				MaxRequests: -1,
				Audience:    "invalid",
			},
			[]string{"label", "duration", "maxRequests", "audience"},
		},
		{
			"valid data (name)",
			core.RateLimitRule{
				Label:       "abc:123",
				Duration:    1,
				MaxRequests: 1,
			},
			[]string{},
		},
		{
			"valid data (name:action)",
			core.RateLimitRule{
				Label:       "abc:123",
				Duration:    1,
				MaxRequests: 1,
			},
			[]string{},
		},
		{
			"valid data (*:action)",
			core.RateLimitRule{
				Label:       "*:123",
				Duration:    1,
				MaxRequests: 1,
			},
			[]string{},
		},
		{
			"valid data (path /a/b)",
			core.RateLimitRule{
				Label:       "/a/b",
				Duration:    1,
				MaxRequests: 1,
			},
			[]string{},
		},
		{
			"valid data (path POST /a/b)",
			core.RateLimitRule{
				Label:       "POST /a/b/",
				Duration:    1,
				MaxRequests: 1,
			},
			[]string{},
		},
		{
			"invalid audience",
			core.RateLimitRule{
				Label:       "/a/b/",
				Duration:    1,
				MaxRequests: 1,
				Audience:    "invalid",
			},
			[]string{"audience"},
		},
		{
			"valid audience - " + core.RateLimitRuleAudienceGuest,
			core.RateLimitRule{
				Label:       "POST /a/b/",
				Duration:    1,
				MaxRequests: 1,
				Audience:    core.RateLimitRuleAudienceGuest,
			},
			[]string{},
		},
		{
			"valid audience - " + core.RateLimitRuleAudienceAuth,
			core.RateLimitRule{
				Label:       "POST /a/b/",
				Duration:    1,
				MaxRequests: 1,
				Audience:    core.RateLimitRuleAudienceAuth,
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

func TestRateLimitRuleDurationTime(t *testing.T) {
	scenarios := []struct {
		config   core.RateLimitRule
		expected time.Duration
	}{
		{core.RateLimitRule{}, 0 * time.Second},
		{core.RateLimitRule{Duration: 1234}, 1234 * time.Second},
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
