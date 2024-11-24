package migratecmd_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TestAutomigrateCollectionCreate(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		lang             string
		expectedTemplate string
	}{
		{
			migratecmd.TemplateLangJS,
			`
/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = new Collection({
    "authAlert": {
      "emailTemplate": {
        "body": "<p>Hello,</p>\n<p>We noticed a login to your {APP_NAME} account from a new location.</p>\n<p>If this was you, you may disregard this email.</p>\n<p><strong>If this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.</strong></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
        "subject": "Login from a new location"
      },
      "enabled": true
    },
    "authRule": "",
    "authToken": {
      "duration": 604800
    },
    "confirmEmailChangeTemplate": {
      "body": "<p>Hello,</p>\n<p>Click on the button below to confirm your new email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Confirm new email</a>\n</p>\n<p><i>If you didn't ask to change your email address, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
      "subject": "Confirm your {APP_NAME} new email address"
    },
    "createRule": null,
    "deleteRule": null,
    "emailChangeToken": {
      "duration": 1800
    },
    "fields": [
      {
        "autogeneratePattern": "[a-z0-9]{15}",
        "hidden": false,
        "id": "text@TEST_RANDOM",
        "max": 15,
        "min": 15,
        "name": "id",
        "pattern": "^[a-z0-9]+$",
        "presentable": false,
        "primaryKey": true,
        "required": true,
        "system": true,
        "type": "text"
      },
      {
        "cost": 0,
        "hidden": true,
        "id": "password@TEST_RANDOM",
        "max": 0,
        "min": 8,
        "name": "password",
        "pattern": "",
        "presentable": false,
        "required": true,
        "system": true,
        "type": "password"
      },
      {
        "autogeneratePattern": "[a-zA-Z0-9]{50}",
        "hidden": true,
        "id": "text@TEST_RANDOM",
        "max": 60,
        "min": 30,
        "name": "tokenKey",
        "pattern": "",
        "presentable": false,
        "primaryKey": false,
        "required": true,
        "system": true,
        "type": "text"
      },
      {
        "exceptDomains": null,
        "hidden": false,
        "id": "email@TEST_RANDOM",
        "name": "email",
        "onlyDomains": null,
        "presentable": false,
        "required": true,
        "system": true,
        "type": "email"
      },
      {
        "hidden": false,
        "id": "bool@TEST_RANDOM",
        "name": "emailVisibility",
        "presentable": false,
        "required": false,
        "system": true,
        "type": "bool"
      },
      {
        "hidden": false,
        "id": "bool@TEST_RANDOM",
        "name": "verified",
        "presentable": false,
        "required": false,
        "system": true,
        "type": "bool"
      }
    ],
    "fileToken": {
      "duration": 180
    },
    "id": "@TEST_RANDOM",
    "indexes": [
      "create index test on new_name (id)",
      "CREATE UNIQUE INDEX ` + "`" + `idx_tokenKey_@TEST_RANDOM` + "`" + ` ON ` + "`" + `new_name` + "`" + ` (` + "`" + `tokenKey` + "`" + `)",
      "CREATE UNIQUE INDEX ` + "`" + `idx_email_@TEST_RANDOM` + "`" + ` ON ` + "`" + `new_name` + "`" + ` (` + "`" + `email` + "`" + `) WHERE ` + "`" + `email` + "`" + ` != ''"
    ],
    "listRule": "@request.auth.id != '' && 1 > 0 || 'backtick` + "`" + `test' = 0",
    "manageRule": "1 != 2",
    "mfa": {
      "duration": 1800,
      "enabled": false,
      "rule": ""
    },
    "name": "new_name",
    "oauth2": {
      "enabled": false,
      "mappedFields": {
        "avatarURL": "",
        "id": "",
        "name": "",
        "username": ""
      }
    },
    "otp": {
      "duration": 180,
      "emailTemplate": {
        "body": "<p>Hello,</p>\n<p>Your one-time password is: <strong>{OTP}</strong></p>\n<p><i>If you didn't ask for the one-time password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
        "subject": "OTP for {APP_NAME}"
      },
      "enabled": false,
      "length": 8
    },
    "passwordAuth": {
      "enabled": true,
      "identityFields": [
        "email"
      ]
    },
    "passwordResetToken": {
      "duration": 1800
    },
    "resetPasswordTemplate": {
      "body": "<p>Hello,</p>\n<p>Click on the button below to reset your password.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Reset password</a>\n</p>\n<p><i>If you didn't ask to reset your password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
      "subject": "Reset your {APP_NAME} password"
    },
    "system": true,
    "type": "auth",
    "updateRule": null,
    "verificationTemplate": {
      "body": "<p>Hello,</p>\n<p>Thank you for joining us at {APP_NAME}.</p>\n<p>Click on the button below to verify your email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Verify</a>\n</p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
      "subject": "Verify your {APP_NAME} email"
    },
    "verificationToken": {
      "duration": 259200
    },
    "viewRule": "id = \"1\""
  });

  return app.save(collection);
}, (app) => {
  const collection = app.findCollectionByNameOrId("@TEST_RANDOM");

  return app.delete(collection);
})
`,
		},
		{
			migratecmd.TemplateLangGo,
			`
package _test_migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := ` + "`" + `{
			"authAlert": {
				"emailTemplate": {
					"body": "<p>Hello,</p>\n<p>We noticed a login to your {APP_NAME} account from a new location.</p>\n<p>If this was you, you may disregard this email.</p>\n<p><strong>If this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.</strong></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
					"subject": "Login from a new location"
				},
				"enabled": true
			},
			"authRule": "",
			"authToken": {
				"duration": 604800
			},
			"confirmEmailChangeTemplate": {
				"body": "<p>Hello,</p>\n<p>Click on the button below to confirm your new email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Confirm new email</a>\n</p>\n<p><i>If you didn't ask to change your email address, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
				"subject": "Confirm your {APP_NAME} new email address"
			},
			"createRule": null,
			"deleteRule": null,
			"emailChangeToken": {
				"duration": 1800
			},
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text@TEST_RANDOM",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"cost": 0,
					"hidden": true,
					"id": "password@TEST_RANDOM",
					"max": 0,
					"min": 8,
					"name": "password",
					"pattern": "",
					"presentable": false,
					"required": true,
					"system": true,
					"type": "password"
				},
				{
					"autogeneratePattern": "[a-zA-Z0-9]{50}",
					"hidden": true,
					"id": "text@TEST_RANDOM",
					"max": 60,
					"min": 30,
					"name": "tokenKey",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"exceptDomains": null,
					"hidden": false,
					"id": "email@TEST_RANDOM",
					"name": "email",
					"onlyDomains": null,
					"presentable": false,
					"required": true,
					"system": true,
					"type": "email"
				},
				{
					"hidden": false,
					"id": "bool@TEST_RANDOM",
					"name": "emailVisibility",
					"presentable": false,
					"required": false,
					"system": true,
					"type": "bool"
				},
				{
					"hidden": false,
					"id": "bool@TEST_RANDOM",
					"name": "verified",
					"presentable": false,
					"required": false,
					"system": true,
					"type": "bool"
				}
			],
			"fileToken": {
				"duration": 180
			},
			"id": "@TEST_RANDOM",
			"indexes": [
				"create index test on new_name (id)",
				"CREATE UNIQUE INDEX ` + "` + \"`\" + `" + `idx_tokenKey_@TEST_RANDOM` + "` + \"`\" + `" + ` ON ` + "` + \"`\" + `" + `new_name` + "` + \"`\" + `" + ` (` + "` + \"`\" + `" + `tokenKey` + "` + \"`\" + `" + `)",
				"CREATE UNIQUE INDEX ` + "` + \"`\" + `" + `idx_email_@TEST_RANDOM` + "` + \"`\" + `" + ` ON ` + "` + \"`\" + `" + `new_name` + "` + \"`\" + `" + ` (` + "` + \"`\" + `" + `email` + "` + \"`\" + `" + `) WHERE ` + "` + \"`\" + `" + `email` + "` + \"`\" + `" + ` != ''"
			],
			"listRule": "@request.auth.id != '' && 1 > 0 || 'backtick` + "` + \"`\" + `" + `test' = 0",
			"manageRule": "1 != 2",
			"mfa": {
				"duration": 1800,
				"enabled": false,
				"rule": ""
			},
			"name": "new_name",
			"oauth2": {
				"enabled": false,
				"mappedFields": {
					"avatarURL": "",
					"id": "",
					"name": "",
					"username": ""
				}
			},
			"otp": {
				"duration": 180,
				"emailTemplate": {
					"body": "<p>Hello,</p>\n<p>Your one-time password is: <strong>{OTP}</strong></p>\n<p><i>If you didn't ask for the one-time password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
					"subject": "OTP for {APP_NAME}"
				},
				"enabled": false,
				"length": 8
			},
			"passwordAuth": {
				"enabled": true,
				"identityFields": [
					"email"
				]
			},
			"passwordResetToken": {
				"duration": 1800
			},
			"resetPasswordTemplate": {
				"body": "<p>Hello,</p>\n<p>Click on the button below to reset your password.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Reset password</a>\n</p>\n<p><i>If you didn't ask to reset your password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
				"subject": "Reset your {APP_NAME} password"
			},
			"system": true,
			"type": "auth",
			"updateRule": null,
			"verificationTemplate": {
				"body": "<p>Hello,</p>\n<p>Thank you for joining us at {APP_NAME}.</p>\n<p>Click on the button below to verify your email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Verify</a>\n</p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
				"subject": "Verify your {APP_NAME} email"
			},
			"verificationToken": {
				"duration": 259200
			},
			"viewRule": "id = \"1\""
		}` + "`" + `

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("@TEST_RANDOM")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.lang, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			migrationsDir := filepath.Join(app.DataDir(), "_test_migrations")

			migratecmd.MustRegister(app, nil, migratecmd.Config{
				TemplateLang: s.lang,
				Automigrate:  true,
				Dir:          migrationsDir,
			})

			app.Bootstrap()

			collection := core.NewAuthCollection("new_name")
			collection.System = true
			collection.ListRule = types.Pointer("@request.auth.id != '' && 1 > 0 || 'backtick`test' = 0")
			collection.ViewRule = types.Pointer(`id = "1"`)
			collection.Indexes = types.JSONArray[string]{"create index test on new_name (id)"}
			collection.ManageRule = types.Pointer("1 != 2")
			//  should be ignored
			collection.OAuth2.Providers = []core.OAuth2ProviderConfig{{Name: "gitlab", ClientId: "abc", ClientSecret: "123"}}
			testSecret := strings.Repeat("a", 30)
			collection.AuthToken.Secret = testSecret
			collection.FileToken.Secret = testSecret
			collection.EmailChangeToken.Secret = testSecret
			collection.PasswordResetToken.Secret = testSecret
			collection.VerificationToken.Secret = testSecret

			// save the newly created dummy collection (with mock request event)
			event := new(core.CollectionRequestEvent)
			event.RequestEvent = &core.RequestEvent{}
			event.App = app
			event.Collection = collection
			err := app.OnCollectionCreateRequest().Trigger(event, func(e *core.CollectionRequestEvent) error {
				return e.App.Save(e.Collection)
			})
			if err != nil {
				t.Fatalf("Failed to save the created dummy collection, got: %v", err)
			}

			files, err := os.ReadDir(migrationsDir)
			if err != nil {
				t.Fatalf("Expected migrationsDir to be created, got %v", err)
			}

			if total := len(files); total != 1 {
				t.Fatalf("Expected 1 file to be generated, got %d: %v", total, files)
			}

			expectedName := "_created_new_name." + s.lang
			if !strings.Contains(files[0].Name(), expectedName) {
				t.Fatalf("Expected filename to contains %q, got %q", expectedName, files[0].Name())
			}

			fullPath := filepath.Join(migrationsDir, files[0].Name())
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("Failed to read the generated migration file: %v", err)
			}
			contentStr := strings.TrimSpace(string(content))

			// replace @TEST_RANDOM placeholder with a regex pattern
			expectedTemplate := strings.ReplaceAll(
				"^"+regexp.QuoteMeta(strings.TrimSpace(s.expectedTemplate))+"$",
				"@TEST_RANDOM",
				`\w+`,
			)
			if !list.ExistInSliceWithRegex(contentStr, []string{expectedTemplate}) {
				t.Fatalf("Expected template \n%v \ngot \n%v", s.expectedTemplate, contentStr)
			}
		})
	}
}

func TestAutomigrateCollectionDelete(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		lang             string
		expectedTemplate string
	}{
		{
			migratecmd.TemplateLangJS,
			`
/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("@TEST_RANDOM");

  return app.delete(collection);
}, (app) => {
  const collection = new Collection({
    "authAlert": {
      "emailTemplate": {
        "body": "<p>Hello,</p>\n<p>We noticed a login to your {APP_NAME} account from a new location.</p>\n<p>If this was you, you may disregard this email.</p>\n<p><strong>If this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.</strong></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
        "subject": "Login from a new location"
      },
      "enabled": true
    },
    "authRule": "",
    "authToken": {
      "duration": 604800
    },
    "confirmEmailChangeTemplate": {
      "body": "<p>Hello,</p>\n<p>Click on the button below to confirm your new email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Confirm new email</a>\n</p>\n<p><i>If you didn't ask to change your email address, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
      "subject": "Confirm your {APP_NAME} new email address"
    },
    "createRule": null,
    "deleteRule": null,
    "emailChangeToken": {
      "duration": 1800
    },
    "fields": [
      {
        "autogeneratePattern": "[a-z0-9]{15}",
        "hidden": false,
        "id": "text@TEST_RANDOM",
        "max": 15,
        "min": 15,
        "name": "id",
        "pattern": "^[a-z0-9]+$",
        "presentable": false,
        "primaryKey": true,
        "required": true,
        "system": true,
        "type": "text"
      },
      {
        "cost": 0,
        "hidden": true,
        "id": "password@TEST_RANDOM",
        "max": 0,
        "min": 8,
        "name": "password",
        "pattern": "",
        "presentable": false,
        "required": true,
        "system": true,
        "type": "password"
      },
      {
        "autogeneratePattern": "[a-zA-Z0-9]{50}",
        "hidden": true,
        "id": "text@TEST_RANDOM",
        "max": 60,
        "min": 30,
        "name": "tokenKey",
        "pattern": "",
        "presentable": false,
        "primaryKey": false,
        "required": true,
        "system": true,
        "type": "text"
      },
      {
        "exceptDomains": null,
        "hidden": false,
        "id": "email3885137012",
        "name": "email",
        "onlyDomains": null,
        "presentable": false,
        "required": true,
        "system": true,
        "type": "email"
      },
      {
        "hidden": false,
        "id": "bool@TEST_RANDOM",
        "name": "emailVisibility",
        "presentable": false,
        "required": false,
        "system": true,
        "type": "bool"
      },
      {
        "hidden": false,
        "id": "bool256245529",
        "name": "verified",
        "presentable": false,
        "required": false,
        "system": true,
        "type": "bool"
      }
    ],
    "fileToken": {
      "duration": 180
    },
    "id": "@TEST_RANDOM",
    "indexes": [
      "create index test on test123 (id)",
      "CREATE UNIQUE INDEX ` + "`" + `idx_tokenKey_@TEST_RANDOM` + "`" + ` ON ` + "`" + `test123` + "`" + ` (` + "`" + `tokenKey` + "`" + `)",
      "CREATE UNIQUE INDEX ` + "`" + `idx_email_@TEST_RANDOM` + "`" + ` ON ` + "`" + `test123` + "`" + ` (` + "`" + `email` + "`" + `) WHERE ` + "`" + `email` + "`" + ` != ''"
    ],
    "listRule": "@request.auth.id != '' && 1 > 0 || 'backtick` + "`" + `test' = 0",
    "manageRule": "1 != 2",
    "mfa": {
      "duration": 1800,
      "enabled": false,
      "rule": ""
    },
    "name": "test123",
    "oauth2": {
      "enabled": false,
      "mappedFields": {
        "avatarURL": "",
        "id": "",
        "name": "",
        "username": ""
      }
    },
    "otp": {
      "duration": 180,
      "emailTemplate": {
        "body": "<p>Hello,</p>\n<p>Your one-time password is: <strong>{OTP}</strong></p>\n<p><i>If you didn't ask for the one-time password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
        "subject": "OTP for {APP_NAME}"
      },
      "enabled": false,
      "length": 8
    },
    "passwordAuth": {
      "enabled": true,
      "identityFields": [
        "email"
      ]
    },
    "passwordResetToken": {
      "duration": 1800
    },
    "resetPasswordTemplate": {
      "body": "<p>Hello,</p>\n<p>Click on the button below to reset your password.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Reset password</a>\n</p>\n<p><i>If you didn't ask to reset your password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
      "subject": "Reset your {APP_NAME} password"
    },
    "system": false,
    "type": "auth",
    "updateRule": null,
    "verificationTemplate": {
      "body": "<p>Hello,</p>\n<p>Thank you for joining us at {APP_NAME}.</p>\n<p>Click on the button below to verify your email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Verify</a>\n</p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
      "subject": "Verify your {APP_NAME} email"
    },
    "verificationToken": {
      "duration": 259200
    },
    "viewRule": "id = \"1\""
  });

  return app.save(collection);
})
`,
		},
		{
			migratecmd.TemplateLangGo,
			`
package _test_migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("@TEST_RANDOM")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	}, func(app core.App) error {
		jsonData := ` + "`" + `{
			"authAlert": {
				"emailTemplate": {
					"body": "<p>Hello,</p>\n<p>We noticed a login to your {APP_NAME} account from a new location.</p>\n<p>If this was you, you may disregard this email.</p>\n<p><strong>If this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.</strong></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
					"subject": "Login from a new location"
				},
				"enabled": true
			},
			"authRule": "",
			"authToken": {
				"duration": 604800
			},
			"confirmEmailChangeTemplate": {
				"body": "<p>Hello,</p>\n<p>Click on the button below to confirm your new email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Confirm new email</a>\n</p>\n<p><i>If you didn't ask to change your email address, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
				"subject": "Confirm your {APP_NAME} new email address"
			},
			"createRule": null,
			"deleteRule": null,
			"emailChangeToken": {
				"duration": 1800
			},
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text@TEST_RANDOM",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"cost": 0,
					"hidden": true,
					"id": "password@TEST_RANDOM",
					"max": 0,
					"min": 8,
					"name": "password",
					"pattern": "",
					"presentable": false,
					"required": true,
					"system": true,
					"type": "password"
				},
				{
					"autogeneratePattern": "[a-zA-Z0-9]{50}",
					"hidden": true,
					"id": "text@TEST_RANDOM",
					"max": 60,
					"min": 30,
					"name": "tokenKey",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"exceptDomains": null,
					"hidden": false,
					"id": "email3885137012",
					"name": "email",
					"onlyDomains": null,
					"presentable": false,
					"required": true,
					"system": true,
					"type": "email"
				},
				{
					"hidden": false,
					"id": "bool@TEST_RANDOM",
					"name": "emailVisibility",
					"presentable": false,
					"required": false,
					"system": true,
					"type": "bool"
				},
				{
					"hidden": false,
					"id": "bool256245529",
					"name": "verified",
					"presentable": false,
					"required": false,
					"system": true,
					"type": "bool"
				}
			],
			"fileToken": {
				"duration": 180
			},
			"id": "@TEST_RANDOM",
			"indexes": [
				"create index test on test123 (id)",
				"CREATE UNIQUE INDEX ` + "` + \"`\" + `" + `idx_tokenKey_@TEST_RANDOM` + "` + \"`\" + `" + ` ON ` + "` + \"`\" + `" + `test123` + "` + \"`\" + `" + ` (` + "` + \"`\" + `" + `tokenKey` + "` + \"`\" + `" + `)",
				"CREATE UNIQUE INDEX ` + "` + \"`\" + `" + `idx_email_@TEST_RANDOM` + "` + \"`\" + `" + ` ON ` + "` + \"`\" + `" + `test123` + "` + \"`\" + `" + ` (` + "` + \"`\" + `" + `email` + "` + \"`\" + `" + `) WHERE ` + "` + \"`\" + `" + `email` + "` + \"`\" + `" + ` != ''"
			],
			"listRule": "@request.auth.id != '' && 1 > 0 || 'backtick` + "` + \"`\" + `" + `test' = 0",
			"manageRule": "1 != 2",
			"mfa": {
				"duration": 1800,
				"enabled": false,
				"rule": ""
			},
			"name": "test123",
			"oauth2": {
				"enabled": false,
				"mappedFields": {
					"avatarURL": "",
					"id": "",
					"name": "",
					"username": ""
				}
			},
			"otp": {
				"duration": 180,
				"emailTemplate": {
					"body": "<p>Hello,</p>\n<p>Your one-time password is: <strong>{OTP}</strong></p>\n<p><i>If you didn't ask for the one-time password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
					"subject": "OTP for {APP_NAME}"
				},
				"enabled": false,
				"length": 8
			},
			"passwordAuth": {
				"enabled": true,
				"identityFields": [
					"email"
				]
			},
			"passwordResetToken": {
				"duration": 1800
			},
			"resetPasswordTemplate": {
				"body": "<p>Hello,</p>\n<p>Click on the button below to reset your password.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Reset password</a>\n</p>\n<p><i>If you didn't ask to reset your password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
				"subject": "Reset your {APP_NAME} password"
			},
			"system": false,
			"type": "auth",
			"updateRule": null,
			"verificationTemplate": {
				"body": "<p>Hello,</p>\n<p>Thank you for joining us at {APP_NAME}.</p>\n<p>Click on the button below to verify your email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Verify</a>\n</p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
				"subject": "Verify your {APP_NAME} email"
			},
			"verificationToken": {
				"duration": 259200
			},
			"viewRule": "id = \"1\""
		}` + "`" + `

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.lang, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			migrationsDir := filepath.Join(app.DataDir(), "_test_migrations")

			// create dummy collection
			collection := core.NewAuthCollection("test123")
			collection.ListRule = types.Pointer("@request.auth.id != '' && 1 > 0 || 'backtick`test' = 0")
			collection.ViewRule = types.Pointer(`id = "1"`)
			collection.Indexes = types.JSONArray[string]{"create index test on test123 (id)"}
			collection.ManageRule = types.Pointer("1 != 2")
			if err := app.Save(collection); err != nil {
				t.Fatalf("Failed to save dummy collection, got: %v", err)
			}

			migratecmd.MustRegister(app, nil, migratecmd.Config{
				TemplateLang: s.lang,
				Automigrate:  true,
				Dir:          migrationsDir,
			})

			app.Bootstrap()

			// delete the newly created dummy collection (with mock request event)
			event := new(core.CollectionRequestEvent)
			event.RequestEvent = &core.RequestEvent{}
			event.App = app
			event.Collection = collection
			err := app.OnCollectionDeleteRequest().Trigger(event, func(e *core.CollectionRequestEvent) error {
				return e.App.Delete(e.Collection)
			})
			if err != nil {
				t.Fatalf("Failed to delete dummy collection, got: %v", err)
			}

			files, err := os.ReadDir(migrationsDir)
			if err != nil {
				t.Fatalf("Expected migrationsDir to be created, got: %v", err)
			}

			if total := len(files); total != 1 {
				t.Fatalf("Expected 1 file to be generated, got %d", total)
			}

			expectedName := "_deleted_test123." + s.lang
			if !strings.Contains(files[0].Name(), expectedName) {
				t.Fatalf("Expected filename to contains %q, got %q", expectedName, files[0].Name())
			}

			fullPath := filepath.Join(migrationsDir, files[0].Name())
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("Failed to read the generated migration file: %v", err)
			}
			contentStr := strings.TrimSpace(string(content))

			// replace @TEST_RANDOM placeholder with a regex pattern
			expectedTemplate := strings.ReplaceAll(
				"^"+regexp.QuoteMeta(strings.TrimSpace(s.expectedTemplate))+"$",
				"@TEST_RANDOM",
				`\w+`,
			)
			if !list.ExistInSliceWithRegex(contentStr, []string{expectedTemplate}) {
				t.Fatalf("Expected template \n%v \ngot \n%v", s.expectedTemplate, contentStr)
			}
		})
	}
}

func TestAutomigrateCollectionUpdate(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		lang             string
		expectedTemplate string
	}{
		{
			migratecmd.TemplateLangJS,
			`
/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("@TEST_RANDOM")

  // update collection data
  unmarshal({
    "createRule": "id = \"nil_update\"",
    "deleteRule": null,
    "fileToken": {
      "duration": 10
    },
    "indexes": [
      "create index test1 on test123_update (f1_name)",
      "CREATE UNIQUE INDEX ` + "`" + `idx_tokenKey_@TEST_RANDOM` + "`" + ` ON ` + "`" + `test123_update` + "`" + ` (` + "`" + `tokenKey` + "`" + `)",
      "CREATE UNIQUE INDEX ` + "`" + `idx_email_@TEST_RANDOM` + "`" + ` ON ` + "`" + `test123_update` + "`" + ` (` + "`" + `email` + "`" + `) WHERE ` + "`" + `email` + "`" + ` != ''"
    ],
    "listRule": "@request.auth.id != ''",
    "name": "test123_update",
    "oauth2": {
      "enabled": true
    },
    "updateRule": "id = \"2_update\""
  }, collection)

  // remove field
  collection.fields.removeById("f3_id")

  // add field
  collection.fields.addAt(8, new Field({
    "autogeneratePattern": "",
    "hidden": false,
    "id": "f4_id",
    "max": 0,
    "min": 0,
    "name": "f4_name",
    "pattern": "` + "`" + `test backtick` + "`" + `123",
    "presentable": false,
    "primaryKey": false,
    "required": false,
    "system": false,
    "type": "text"
  }))

  // update field
  collection.fields.addAt(7, new Field({
    "hidden": false,
    "id": "f2_id",
    "max": null,
    "min": 10,
    "name": "f2_name_new",
    "onlyInt": false,
    "presentable": false,
    "required": false,
    "system": false,
    "type": "number"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("@TEST_RANDOM")

  // update collection data
  unmarshal({
    "createRule": null,
    "deleteRule": "id = \"3\"",
    "fileToken": {
      "duration": 180
    },
    "indexes": [
      "create index test1 on test123 (f1_name)",
      "CREATE UNIQUE INDEX ` + "`" + `idx_tokenKey_@TEST_RANDOM` + "`" + ` ON ` + "`" + `test123` + "`" + ` (` + "`" + `tokenKey` + "`" + `)",
      "CREATE UNIQUE INDEX ` + "`" + `idx_email_@TEST_RANDOM` + "`" + ` ON ` + "`" + `test123` + "`" + ` (` + "`" + `email` + "`" + `) WHERE ` + "`" + `email` + "`" + ` != ''"
    ],
    "listRule": "@request.auth.id != '' && 1 != 2",
    "name": "test123",
    "oauth2": {
      "enabled": false
    },
    "updateRule": "id = \"2\""
  }, collection)

  // add field
  collection.fields.addAt(8, new Field({
    "hidden": false,
    "id": "f3_id",
    "name": "f3_name",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "bool"
  }))

  // remove field
  collection.fields.removeById("f4_id")

  // update field
  collection.fields.addAt(7, new Field({
    "hidden": false,
    "id": "f2_id",
    "max": null,
    "min": 10,
    "name": "f2_name",
    "onlyInt": false,
    "presentable": false,
    "required": false,
    "system": false,
    "type": "number"
  }))

  return app.save(collection)
})

`,
		},
		{
			migratecmd.TemplateLangGo,
			`
package _test_migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("@TEST_RANDOM")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(` + "`" + `{
			"createRule": "id = \"nil_update\"",
			"deleteRule": null,
			"fileToken": {
				"duration": 10
			},
			"indexes": [
				"create index test1 on test123_update (f1_name)",
				"CREATE UNIQUE INDEX ` + "` + \"`\" + `" + `idx_tokenKey_@TEST_RANDOM` + "` + \"`\" + `" + ` ON ` + "` + \"`\" + `" + `test123_update` + "` + \"`\" + `" + ` (` + "` + \"`\" + `" + `tokenKey` + "` + \"`\" + `" + `)",
				"CREATE UNIQUE INDEX ` + "` + \"`\" + `" + `idx_email_@TEST_RANDOM` + "` + \"`\" + `" + ` ON ` + "` + \"`\" + `" + `test123_update` + "` + \"`\" + `" + ` (` + "` + \"`\" + `" + `email` + "` + \"`\" + `" + `) WHERE ` + "` + \"`\" + `" + `email` + "` + \"`\" + `" + ` != ''"
			],
			"listRule": "@request.auth.id != ''",
			"name": "test123_update",
			"oauth2": {
				"enabled": true
			},
			"updateRule": "id = \"2_update\""
		}` + "`" + `), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("f3_id")

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(8, []byte(` + "`" + `{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "f4_id",
			"max": 0,
			"min": 0,
			"name": "f4_name",
			"pattern": "` + "` + \"`\" + `" + `test backtick` + "` + \"`\" + `" + `123",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}` + "`" + `)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(7, []byte(` + "`" + `{
			"hidden": false,
			"id": "f2_id",
			"max": null,
			"min": 10,
			"name": "f2_name_new",
			"onlyInt": false,
			"presentable": false,
			"required": false,
			"system": false,
			"type": "number"
		}` + "`" + `)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("@TEST_RANDOM")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(` + "`" + `{
			"createRule": null,
			"deleteRule": "id = \"3\"",
			"fileToken": {
				"duration": 180
			},
			"indexes": [
				"create index test1 on test123 (f1_name)",
				"CREATE UNIQUE INDEX ` + "` + \"`\" + `" + `idx_tokenKey_@TEST_RANDOM` + "` + \"`\" + `" + ` ON ` + "` + \"`\" + `" + `test123` + "` + \"`\" + `" + ` (` + "` + \"`\" + `" + `tokenKey` + "` + \"`\" + `" + `)",
				"CREATE UNIQUE INDEX ` + "` + \"`\" + `" + `idx_email_@TEST_RANDOM` + "` + \"`\" + `" + ` ON ` + "` + \"`\" + `" + `test123` + "` + \"`\" + `" + ` (` + "` + \"`\" + `" + `email` + "` + \"`\" + `" + `) WHERE ` + "` + \"`\" + `" + `email` + "` + \"`\" + `" + ` != ''"
			],
			"listRule": "@request.auth.id != '' && 1 != 2",
			"name": "test123",
			"oauth2": {
				"enabled": false
			},
			"updateRule": "id = \"2\""
		}` + "`" + `), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(8, []byte(` + "`" + `{
			"hidden": false,
			"id": "f3_id",
			"name": "f3_name",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "bool"
		}` + "`" + `)); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("f4_id")

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(7, []byte(` + "`" + `{
			"hidden": false,
			"id": "f2_id",
			"max": null,
			"min": 10,
			"name": "f2_name",
			"onlyInt": false,
			"presentable": false,
			"required": false,
			"system": false,
			"type": "number"
		}` + "`" + `)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.lang, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			migrationsDir := filepath.Join(app.DataDir(), "_test_migrations")

			// create dummy collection
			collection := core.NewAuthCollection("test123")
			collection.ListRule = types.Pointer("@request.auth.id != '' && 1 != 2")
			collection.ViewRule = types.Pointer(`id = "1"`)
			collection.UpdateRule = types.Pointer(`id = "2"`)
			collection.CreateRule = nil
			collection.DeleteRule = types.Pointer(`id = "3"`)
			collection.Indexes = types.JSONArray[string]{"create index test1 on test123 (f1_name)"}
			collection.ManageRule = types.Pointer("1 != 2")
			collection.Fields.Add(&core.TextField{
				Id:       "f1_id",
				Name:     "f1_name",
				Required: true,
			})
			collection.Fields.Add(&core.NumberField{
				Id:   "f2_id",
				Name: "f2_name",
				Min:  types.Pointer(10.0),
			})
			collection.Fields.Add(&core.BoolField{
				Id:   "f3_id",
				Name: "f3_name",
			})

			if err := app.Save(collection); err != nil {
				t.Fatalf("Failed to save dummy collection, got %v", err)
			}

			// init plugin
			migratecmd.MustRegister(app, nil, migratecmd.Config{
				TemplateLang: s.lang,
				Automigrate:  true,
				Dir:          migrationsDir,
			})
			app.Bootstrap()

			// update the dummy collection
			collection.Name = "test123_update"
			collection.ListRule = types.Pointer("@request.auth.id != ''")
			collection.ViewRule = types.Pointer(`id = "1"`) // no change
			collection.UpdateRule = types.Pointer(`id = "2_update"`)
			collection.CreateRule = types.Pointer(`id = "nil_update"`)
			collection.DeleteRule = nil
			collection.Indexes = types.JSONArray[string]{
				"create index test1 on test123_update (f1_name)",
			}
			collection.Fields.RemoveById("f3_id")
			collection.Fields.Add(&core.TextField{
				Id:      "f4_id",
				Name:    "f4_name",
				Pattern: "`test backtick`123",
			})
			f := collection.Fields.GetById("f2_id")
			f.SetName("f2_name_new")
			collection.OAuth2.Enabled = true
			collection.FileToken.Duration = 10
			//  should be ignored
			collection.OAuth2.Providers = []core.OAuth2ProviderConfig{{Name: "gitlab", ClientId: "abc", ClientSecret: "123"}}
			testSecret := strings.Repeat("b", 30)
			collection.AuthToken.Secret = testSecret
			collection.FileToken.Secret = testSecret
			collection.EmailChangeToken.Secret = testSecret
			collection.PasswordResetToken.Secret = testSecret
			collection.VerificationToken.Secret = testSecret

			// save the changes and trigger automigrate (with mock request event)
			event := new(core.CollectionRequestEvent)
			event.RequestEvent = &core.RequestEvent{}
			event.App = app
			event.Collection = collection
			err := app.OnCollectionUpdateRequest().Trigger(event, func(e *core.CollectionRequestEvent) error {
				return e.App.Save(e.Collection)
			})
			if err != nil {
				t.Fatalf("Failed to save dummy collection changes, got %v", err)
			}

			files, err := os.ReadDir(migrationsDir)
			if err != nil {
				t.Fatalf("Expected migrationsDir to be created, got: %v", err)
			}

			if total := len(files); total != 1 {
				t.Fatalf("Expected 1 file to be generated, got %d", total)
			}

			expectedName := "_updated_test123." + s.lang
			if !strings.Contains(files[0].Name(), expectedName) {
				t.Fatalf("Expected filename to contains %q, got %q", expectedName, files[0].Name())
			}

			fullPath := filepath.Join(migrationsDir, files[0].Name())
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("Failed to read the generated migration file: %v", err)
			}
			contentStr := strings.TrimSpace(string(content))

			// replace @TEST_RANDOM placeholder with a regex pattern
			expectedTemplate := strings.ReplaceAll(
				"^"+regexp.QuoteMeta(strings.TrimSpace(s.expectedTemplate))+"$",
				"@TEST_RANDOM",
				`\w+`,
			)
			if !list.ExistInSliceWithRegex(contentStr, []string{expectedTemplate}) {
				t.Fatalf("Expected template \n%v \ngot \n%v", s.expectedTemplate, contentStr)
			}
		})
	}
}

func TestAutomigrateCollectionNoChanges(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		lang string
	}{
		{
			migratecmd.TemplateLangJS,
		},
		{
			migratecmd.TemplateLangGo,
		},
	}

	for _, s := range scenarios {
		t.Run(s.lang, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			migrationsDir := filepath.Join(app.DataDir(), "_test_migrations")

			// create dummy collection
			collection := core.NewAuthCollection("test123")

			if err := app.Save(collection); err != nil {
				t.Fatalf("Failed to save dummy collection, got %v", err)
			}

			// init plugin
			migratecmd.MustRegister(app, nil, migratecmd.Config{
				TemplateLang: s.lang,
				Automigrate:  true,
				Dir:          migrationsDir,
			})
			app.Bootstrap()

			//  should be ignored
			collection.OAuth2.Providers = []core.OAuth2ProviderConfig{{Name: "gitlab", ClientId: "abc", ClientSecret: "123"}}
			testSecret := strings.Repeat("b", 30)
			collection.AuthToken.Secret = testSecret
			collection.FileToken.Secret = testSecret
			collection.EmailChangeToken.Secret = testSecret
			collection.PasswordResetToken.Secret = testSecret
			collection.VerificationToken.Secret = testSecret

			// resave without other changes and trigger automigrate (with mock request event)
			event := new(core.CollectionRequestEvent)
			event.RequestEvent = &core.RequestEvent{}
			event.App = app
			event.Collection = collection
			err := app.OnCollectionUpdateRequest().Trigger(event, func(e *core.CollectionRequestEvent) error {
				return e.App.Save(e.Collection)
			})
			if err != nil {
				t.Fatalf("Failed to save dummy collection update, got %v", err)
			}

			files, _ := os.ReadDir(migrationsDir)
			if total := len(files); total != 0 {
				t.Fatalf("Expected 0 files to be generated, got %d", total)
			}
		})
	}
}
