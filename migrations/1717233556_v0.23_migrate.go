package migrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

// note: this migration will be deleted in future version

func init() {
	core.SystemMigrations.Register(func(txApp core.App) error {
		// note: mfas and authOrigins tables are available only with v0.23
		hasUpgraded := txApp.HasTable(core.CollectionNameMFAs) && txApp.HasTable(core.CollectionNameAuthOrigins)
		if hasUpgraded {
			return nil
		}

		oldSettings, err := loadOldSettings(txApp)
		if err != nil {
			return fmt.Errorf("failed to fetch old settings: %w", err)
		}

		if err = migrateOldCollections(txApp, oldSettings); err != nil {
			return err
		}

		if err = migrateSuperusers(txApp, oldSettings); err != nil {
			return fmt.Errorf("failed to migrate admins->superusers: %w", err)
		}

		if err = migrateSettings(txApp, oldSettings); err != nil {
			return fmt.Errorf("failed to migrate settings: %w", err)
		}

		if err = migrateExternalAuths(txApp); err != nil {
			return fmt.Errorf("failed to migrate externalAuths: %w", err)
		}

		if err = createMFAsCollection(txApp); err != nil {
			return fmt.Errorf("failed to create mfas collection: %w", err)
		}

		if err = createOTPsCollection(txApp); err != nil {
			return fmt.Errorf("failed to create otps collection: %w", err)
		}

		if err = createAuthOriginsCollection(txApp); err != nil {
			return fmt.Errorf("failed to create authOrigins collection: %w", err)
		}

		err = os.Remove(filepath.Join(txApp.DataDir(), "logs.db"))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			txApp.Logger().Warn("Failed to delete old logs.db file", "error", err)
		}

		return nil
	}, nil)
}

// -------------------------------------------------------------------

func migrateSuperusers(txApp core.App, oldSettings *oldSettingsModel) error {
	// create new superusers collection and table
	err := createSuperusersCollection(txApp)
	if err != nil {
		return err
	}

	// update with the token options from the old settings
	superusersCollection, err := txApp.FindCollectionByNameOrId(core.CollectionNameSuperusers)
	if err != nil {
		return err
	}

	superusersCollection.AuthToken.Secret = zeroFallback(
		cast.ToString(getMapVal(oldSettings.Value, "adminAuthToken", "secret")),
		superusersCollection.AuthToken.Secret,
	)
	superusersCollection.AuthToken.Duration = zeroFallback(
		cast.ToInt64(getMapVal(oldSettings.Value, "adminAuthToken", "duration")),
		superusersCollection.AuthToken.Duration,
	)
	superusersCollection.PasswordResetToken.Secret = zeroFallback(
		cast.ToString(getMapVal(oldSettings.Value, "adminPasswordResetToken", "secret")),
		superusersCollection.PasswordResetToken.Secret,
	)
	superusersCollection.PasswordResetToken.Duration = zeroFallback(
		cast.ToInt64(getMapVal(oldSettings.Value, "adminPasswordResetToken", "duration")),
		superusersCollection.PasswordResetToken.Duration,
	)
	superusersCollection.FileToken.Secret = zeroFallback(
		cast.ToString(getMapVal(oldSettings.Value, "adminFileToken", "secret")),
		superusersCollection.FileToken.Secret,
	)
	superusersCollection.FileToken.Duration = zeroFallback(
		cast.ToInt64(getMapVal(oldSettings.Value, "adminFileToken", "duration")),
		superusersCollection.FileToken.Duration,
	)
	if err = txApp.Save(superusersCollection); err != nil {
		return fmt.Errorf("failed to migrate token configs: %w", err)
	}

	// copy old admins records into the new one
	_, err = txApp.DB().NewQuery(`
		INSERT INTO {{` + core.CollectionNameSuperusers + `}} ([[id]], [[verified]], [[email]], [[password]], [[tokenKey]],  [[created]], [[updated]])
		SELECT [[id]], true, [[email]], [[passwordHash]], [[tokenKey]], [[created]], [[updated]] FROM {{_admins}};
	`).Execute()
	if err != nil {
		return err
	}

	// remove old admins table
	_, err = txApp.DB().DropTable("_admins").Execute()
	if err != nil {
		return err
	}

	return nil
}

// -------------------------------------------------------------------

type oldSettingsModel struct {
	Id       string         `db:"id" json:"id"`
	Key      string         `db:"key" json:"key"`
	RawValue types.JSONRaw  `db:"value" json:"value"`
	Value    map[string]any `db:"-" json:"-"`
}

func loadOldSettings(txApp core.App) (*oldSettingsModel, error) {
	oldSettings := &oldSettingsModel{Value: map[string]any{}}
	err := txApp.DB().Select().From("_params").Where(dbx.HashExp{"key": "settings"}).One(oldSettings)
	if err != nil {
		return nil, err
	}

	// try without decrypt
	plainDecodeErr := json.Unmarshal(oldSettings.RawValue, &oldSettings.Value)

	// failed, try to decrypt
	if plainDecodeErr != nil {
		encryptionKey := os.Getenv(txApp.EncryptionEnv())

		// load without decryption has failed and there is no encryption key to use for decrypt
		if encryptionKey == "" {
			return nil, fmt.Errorf("invalid settings db data or missing encryption key %q", txApp.EncryptionEnv())
		}

		// decrypt
		decrypted, decryptErr := security.Decrypt(string(oldSettings.RawValue), encryptionKey)
		if decryptErr != nil {
			return nil, decryptErr
		}

		// decode again
		decryptedDecodeErr := json.Unmarshal(decrypted, &oldSettings.Value)
		if decryptedDecodeErr != nil {
			return nil, decryptedDecodeErr
		}
	}

	return oldSettings, nil
}

func migrateSettings(txApp core.App, oldSettings *oldSettingsModel) error {
	// renamed old params collection
	_, err := txApp.DB().RenameTable("_params", "_params_old").Execute()
	if err != nil {
		return err
	}

	// create new params table
	err = createParamsTable(txApp)
	if err != nil {
		return err
	}

	// migrate old settings
	newSettings := txApp.Settings()
	// ---
	newSettings.Meta.AppName = cast.ToString(getMapVal(oldSettings.Value, "meta", "appName"))
	newSettings.Meta.AppURL = strings.TrimSuffix(cast.ToString(getMapVal(oldSettings.Value, "meta", "appUrl")), "/")
	newSettings.Meta.HideControls = cast.ToBool(getMapVal(oldSettings.Value, "meta", "hideControls"))
	newSettings.Meta.SenderName = cast.ToString(getMapVal(oldSettings.Value, "meta", "senderName"))
	newSettings.Meta.SenderAddress = cast.ToString(getMapVal(oldSettings.Value, "meta", "senderAddress"))
	// ---
	newSettings.Logs.MaxDays = cast.ToInt(getMapVal(oldSettings.Value, "logs", "maxDays"))
	newSettings.Logs.MinLevel = cast.ToInt(getMapVal(oldSettings.Value, "logs", "minLevel"))
	newSettings.Logs.LogIP = cast.ToBool(getMapVal(oldSettings.Value, "logs", "logIp"))
	// ---
	newSettings.SMTP.Enabled = cast.ToBool(getMapVal(oldSettings.Value, "smtp", "enabled"))
	newSettings.SMTP.Port = cast.ToInt(getMapVal(oldSettings.Value, "smtp", "port"))
	newSettings.SMTP.Host = cast.ToString(getMapVal(oldSettings.Value, "smtp", "host"))
	newSettings.SMTP.Username = cast.ToString(getMapVal(oldSettings.Value, "smtp", "username"))
	newSettings.SMTP.Password = cast.ToString(getMapVal(oldSettings.Value, "smtp", "password"))
	newSettings.SMTP.AuthMethod = cast.ToString(getMapVal(oldSettings.Value, "smtp", "authMethod"))
	newSettings.SMTP.TLS = cast.ToBool(getMapVal(oldSettings.Value, "smtp", "tls"))
	newSettings.SMTP.LocalName = cast.ToString(getMapVal(oldSettings.Value, "smtp", "localName"))
	// ---
	newSettings.Backups.Cron = cast.ToString(getMapVal(oldSettings.Value, "backups", "cron"))
	newSettings.Backups.CronMaxKeep = cast.ToInt(getMapVal(oldSettings.Value, "backups", "cronMaxKeep"))
	newSettings.Backups.S3 = core.S3Config{
		Enabled:        cast.ToBool(getMapVal(oldSettings.Value, "backups", "s3", "enabled")),
		Bucket:         cast.ToString(getMapVal(oldSettings.Value, "backups", "s3", "bucket")),
		Region:         cast.ToString(getMapVal(oldSettings.Value, "backups", "s3", "region")),
		Endpoint:       cast.ToString(getMapVal(oldSettings.Value, "backups", "s3", "endpoint")),
		AccessKey:      cast.ToString(getMapVal(oldSettings.Value, "backups", "s3", "accessKey")),
		Secret:         cast.ToString(getMapVal(oldSettings.Value, "backups", "s3", "secret")),
		ForcePathStyle: cast.ToBool(getMapVal(oldSettings.Value, "backups", "s3", "forcePathStyle")),
	}
	// ---
	newSettings.S3 = core.S3Config{
		Enabled:        cast.ToBool(getMapVal(oldSettings.Value, "s3", "enabled")),
		Bucket:         cast.ToString(getMapVal(oldSettings.Value, "s3", "bucket")),
		Region:         cast.ToString(getMapVal(oldSettings.Value, "s3", "region")),
		Endpoint:       cast.ToString(getMapVal(oldSettings.Value, "s3", "endpoint")),
		AccessKey:      cast.ToString(getMapVal(oldSettings.Value, "s3", "accessKey")),
		Secret:         cast.ToString(getMapVal(oldSettings.Value, "s3", "secret")),
		ForcePathStyle: cast.ToBool(getMapVal(oldSettings.Value, "s3", "forcePathStyle")),
	}
	// ---
	err = txApp.Save(newSettings)
	if err != nil {
		return err
	}

	// remove old params table
	_, err = txApp.DB().DropTable("_params_old").Execute()
	if err != nil {
		return err
	}

	return nil
}

// -------------------------------------------------------------------

func migrateExternalAuths(txApp core.App) error {
	// renamed old externalAuths table
	_, err := txApp.DB().RenameTable("_externalAuths", "_externalAuths_old").Execute()
	if err != nil {
		return err
	}

	// create new externalAuths collection and table
	err = createExternalAuthsCollection(txApp)
	if err != nil {
		return err
	}

	// copy old externalAuths records into the new one
	_, err = txApp.DB().NewQuery(`
		INSERT INTO {{` + core.CollectionNameExternalAuths + `}} ([[id]], [[collectionRef]], [[recordRef]], [[provider]], [[providerId]], [[created]], [[updated]])
		SELECT [[id]], [[collectionId]], [[recordId]], [[provider]], [[providerId]], [[created]], [[updated]] FROM {{_externalAuths_old}};
	`).Execute()
	if err != nil {
		return err
	}

	// remove old externalAuths table
	_, err = txApp.DB().DropTable("_externalAuths_old").Execute()
	if err != nil {
		return err
	}

	return nil
}

// -------------------------------------------------------------------

func migrateOldCollections(txApp core.App, oldSettings *oldSettingsModel) error {
	oldCollections := []*OldCollectionModel{}
	err := txApp.DB().Select().From("_collections").All(&oldCollections)
	if err != nil {
		return err
	}

	for _, c := range oldCollections {
		dummyAuthCollection := core.NewAuthCollection("test")

		options := c.Options
		c.Options = types.JSONMap[any]{} // reset

		// update rules
		// ---
		c.ListRule = migrateRule(c.ListRule)
		c.ViewRule = migrateRule(c.ViewRule)
		c.CreateRule = migrateRule(c.CreateRule)
		c.UpdateRule = migrateRule(c.UpdateRule)
		c.DeleteRule = migrateRule(c.DeleteRule)

		// migrate fields
		// ---
		for i, field := range c.Schema {
			switch cast.ToString(field["type"]) {
			case "bool":
				field = toBoolField(field)
			case "number":
				field = toNumberField(field)
			case "text":
				field = toTextField(field)
			case "url":
				field = toURLField(field)
			case "email":
				field = toEmailField(field)
			case "editor":
				field = toEditorField(field)
			case "date":
				field = toDateField(field)
			case "select":
				field = toSelectField(field)
			case "json":
				field = toJSONField(field)
			case "relation":
				field = toRelationField(field)
			case "file":
				field = toFileField(field)
			}
			c.Schema[i] = field
		}

		// type specific changes
		switch c.Type {
		case "auth":
			// token configs
			// ---
			c.Options["authToken"] = map[string]any{
				"secret":   zeroFallback(cast.ToString(getMapVal(oldSettings.Value, "recordAuthToken", "secret")), dummyAuthCollection.AuthToken.Secret),
				"duration": zeroFallback(cast.ToInt64(getMapVal(oldSettings.Value, "recordAuthToken", "duration")), dummyAuthCollection.AuthToken.Duration),
			}
			c.Options["passwordResetToken"] = map[string]any{
				"secret":   zeroFallback(cast.ToString(getMapVal(oldSettings.Value, "recordPasswordResetToken", "secret")), dummyAuthCollection.PasswordResetToken.Secret),
				"duration": zeroFallback(cast.ToInt64(getMapVal(oldSettings.Value, "recordPasswordResetToken", "duration")), dummyAuthCollection.PasswordResetToken.Duration),
			}
			c.Options["emailChangeToken"] = map[string]any{
				"secret":   zeroFallback(cast.ToString(getMapVal(oldSettings.Value, "recordEmailChangeToken", "secret")), dummyAuthCollection.EmailChangeToken.Secret),
				"duration": zeroFallback(cast.ToInt64(getMapVal(oldSettings.Value, "recordEmailChangeToken", "duration")), dummyAuthCollection.EmailChangeToken.Duration),
			}
			c.Options["verificationToken"] = map[string]any{
				"secret":   zeroFallback(cast.ToString(getMapVal(oldSettings.Value, "recordVerificationToken", "secret")), dummyAuthCollection.VerificationToken.Secret),
				"duration": zeroFallback(cast.ToInt64(getMapVal(oldSettings.Value, "recordVerificationToken", "duration")), dummyAuthCollection.VerificationToken.Duration),
			}
			c.Options["fileToken"] = map[string]any{
				"secret":   zeroFallback(cast.ToString(getMapVal(oldSettings.Value, "recordFileToken", "secret")), dummyAuthCollection.FileToken.Secret),
				"duration": zeroFallback(cast.ToInt64(getMapVal(oldSettings.Value, "recordFileToken", "duration")), dummyAuthCollection.FileToken.Duration),
			}

			onlyVerified := cast.ToBool(options["onlyVerified"])
			if onlyVerified {
				c.Options["authRule"] = "verified=true"
			} else {
				c.Options["authRule"] = ""
			}

			c.Options["manageRule"] = nil
			if options["manageRule"] != nil {
				manageRule, err := cast.ToStringE(options["manageRule"])
				if err == nil && manageRule != "" {
					c.Options["manageRule"] = migrateRule(&manageRule)
				}
			}

			// passwordAuth
			identityFields := []string{}
			if cast.ToBool(options["allowEmailAuth"]) {
				identityFields = append(identityFields, "email")
			}
			if cast.ToBool(options["allowUsernameAuth"]) {
				identityFields = append(identityFields, "username")
			}
			c.Options["passwordAuth"] = map[string]any{
				"enabled":        len(identityFields) > 0,
				"identityFields": identityFields,
			}

			// oauth2
			// ---
			oauth2Providers := []map[string]any{}
			providerNames := []string{
				"googleAuth",
				"facebookAuth",
				"githubAuth",
				"gitlabAuth",
				"discordAuth",
				"twitterAuth",
				"microsoftAuth",
				"spotifyAuth",
				"kakaoAuth",
				"twitchAuth",
				"stravaAuth",
				"giteeAuth",
				"livechatAuth",
				"giteaAuth",
				"oidcAuth",
				"oidc2Auth",
				"oidc3Auth",
				"appleAuth",
				"instagramAuth",
				"vkAuth",
				"yandexAuth",
				"patreonAuth",
				"mailcowAuth",
				"bitbucketAuth",
				"planningcenterAuth",
			}
			for _, name := range providerNames {
				if !cast.ToBool(getMapVal(oldSettings.Value, name, "enabled")) {
					continue
				}
				oauth2Providers = append(oauth2Providers, map[string]any{
					"name":         strings.TrimSuffix(name, "Auth"),
					"clientId":     cast.ToString(getMapVal(oldSettings.Value, name, "clientId")),
					"clientSecret": cast.ToString(getMapVal(oldSettings.Value, name, "clientSecret")),
					"authURL":      cast.ToString(getMapVal(oldSettings.Value, name, "authUrl")),
					"tokenURL":     cast.ToString(getMapVal(oldSettings.Value, name, "tokenUrl")),
					"userInfoURL":  cast.ToString(getMapVal(oldSettings.Value, name, "userApiUrl")),
					"displayName":  cast.ToString(getMapVal(oldSettings.Value, name, "displayName")),
					"pkce":         getMapVal(oldSettings.Value, name, "pkce"),
				})
			}

			c.Options["oauth2"] = map[string]any{
				"enabled":   cast.ToBool(options["allowOAuth2Auth"]) && len(oauth2Providers) > 0,
				"providers": oauth2Providers,
				"mappedFields": map[string]string{
					"username": "username",
				},
			}

			// default email templates
			// ---
			emailTemplates := map[string]core.EmailTemplate{
				"verificationTemplate":       dummyAuthCollection.VerificationTemplate,
				"resetPasswordTemplate":      dummyAuthCollection.ResetPasswordTemplate,
				"confirmEmailChangeTemplate": dummyAuthCollection.ConfirmEmailChangeTemplate,
			}
			for name, fallback := range emailTemplates {
				c.Options[name] = map[string]any{
					"subject": zeroFallback(
						cast.ToString(getMapVal(oldSettings.Value, "meta", name, "subject")),
						fallback.Subject,
					),
					"body": zeroFallback(
						strings.ReplaceAll(
							cast.ToString(getMapVal(oldSettings.Value, "meta", name, "body")),
							"{ACTION_URL}",
							cast.ToString(getMapVal(oldSettings.Value, "meta", name, "actionUrl")),
						),
						fallback.Body,
					),
				}
			}

			// mfa
			// ---
			c.Options["mfa"] = map[string]any{
				"enabled":  dummyAuthCollection.MFA.Enabled,
				"duration": dummyAuthCollection.MFA.Duration,
				"rule":     dummyAuthCollection.MFA.Rule,
			}

			// otp
			// ---
			c.Options["otp"] = map[string]any{
				"enabled":  dummyAuthCollection.OTP.Enabled,
				"duration": dummyAuthCollection.OTP.Duration,
				"length":   dummyAuthCollection.OTP.Length,
				"emailTemplate": map[string]any{
					"subject": dummyAuthCollection.OTP.EmailTemplate.Subject,
					"body":    dummyAuthCollection.OTP.EmailTemplate.Body,
				},
			}

			// auth alerts
			// ---
			c.Options["authAlert"] = map[string]any{
				"enabled": dummyAuthCollection.AuthAlert.Enabled,
				"emailTemplate": map[string]any{
					"subject": dummyAuthCollection.AuthAlert.EmailTemplate.Subject,
					"body":    dummyAuthCollection.AuthAlert.EmailTemplate.Body,
				},
			}

			// add system field indexes
			// ---
			c.Indexes = append(types.JSONArray[string]{
				fmt.Sprintf("CREATE UNIQUE INDEX `_%s_username_idx` ON `%s` (username COLLATE NOCASE)", c.Id, c.Name),
				fmt.Sprintf("CREATE UNIQUE INDEX `_%s_email_idx` ON `%s` (`email`) WHERE `email` != ''", c.Id, c.Name),
				fmt.Sprintf("CREATE UNIQUE INDEX `_%s_tokenKey_idx` ON `%s` (`tokenKey`)", c.Id, c.Name),
			}, c.Indexes...)

			// prepend the auth system fields
			// ---
			tokenKeyField := map[string]any{
				"id":                  fieldIdChecksum("text", "tokenKey"),
				"type":                "text",
				"name":                "tokenKey",
				"system":              true,
				"hidden":              true,
				"required":            true,
				"presentable":         false,
				"primaryKey":          false,
				"min":                 30,
				"max":                 60,
				"pattern":             "",
				"autogeneratePattern": "[a-zA-Z0-9_]{50}",
			}
			passwordField := map[string]any{
				"id":          fieldIdChecksum("password", "password"),
				"type":        "password",
				"name":        "password",
				"presentable": false,
				"system":      true,
				"hidden":      true,
				"required":    true,
				"pattern":     "",
				"min":         cast.ToInt(options["minPasswordLength"]),
				"cost":        bcrypt.DefaultCost, // new default
			}
			emailField := map[string]any{
				"id":            fieldIdChecksum("email", "email"),
				"type":          "email",
				"name":          "email",
				"system":        true,
				"hidden":        false,
				"presentable":   false,
				"required":      cast.ToBool(options["requireEmail"]),
				"exceptDomains": cast.ToStringSlice(options["exceptEmailDomains"]),
				"onlyDomains":   cast.ToStringSlice(options["onlyEmailDomains"]),
			}
			emailVisibilityField := map[string]any{
				"id":          fieldIdChecksum("bool", "emailVisibility"),
				"type":        "bool",
				"name":        "emailVisibility",
				"system":      true,
				"hidden":      false,
				"presentable": false,
				"required":    false,
			}
			verifiedField := map[string]any{
				"id":          fieldIdChecksum("bool", "verified"),
				"type":        "bool",
				"name":        "verified",
				"system":      true,
				"hidden":      false,
				"presentable": false,
				"required":    false,
			}
			usernameField := map[string]any{
				"id":                  fieldIdChecksum("text", "username"),
				"type":                "text",
				"name":                "username",
				"system":              false,
				"hidden":              false,
				"required":            true,
				"presentable":         false,
				"primaryKey":          false,
				"min":                 3,
				"max":                 150,
				"pattern":             `^[\w][\w\.\-]*$`,
				"autogeneratePattern": "users[0-9]{6}",
			}
			c.Schema = append(types.JSONArray[types.JSONMap[any]]{
				passwordField,
				tokenKeyField,
				emailField,
				emailVisibilityField,
				verifiedField,
				usernameField,
			}, c.Schema...)

			// rename passwordHash records rable column to password
			// ---
			_, err = txApp.DB().RenameColumn(c.Name, "passwordHash", "password").Execute()
			if err != nil {
				return err
			}

			// delete unnecessary auth columns
			dropColumns := []string{"lastResetSentAt", "lastVerificationSentAt", "lastLoginAlertSentAt"}
			for _, drop := range dropColumns {
				// ignore errors in case the columns don't exist
				_, _ = txApp.DB().DropColumn(c.Name, drop).Execute()
			}
		case "view":
			c.Options["viewQuery"] = cast.ToString(options["query"])
		}

		// prepend the id field
		idField := map[string]any{
			"id":                  fieldIdChecksum("text", "id"),
			"type":                "text",
			"name":                "id",
			"system":              true,
			"required":            true,
			"presentable":         false,
			"hidden":              false,
			"primaryKey":          true,
			"min":                 15,
			"max":                 15,
			"pattern":             "^[a-z0-9]+$",
			"autogeneratePattern": "[a-z0-9]{15}",
		}
		c.Schema = append(types.JSONArray[types.JSONMap[any]]{idField}, c.Schema...)

		var addCreated, addUpdated bool

		if c.Type == "view" {
			// manually check if the view has created/updated columns
			columns, _ := txApp.TableColumns(c.Name)
			for _, c := range columns {
				if strings.EqualFold(c, "created") {
					addCreated = true
				} else if strings.EqualFold(c, "updated") {
					addUpdated = true
				}
			}
		} else {
			addCreated = true
			addUpdated = true
		}

		if addCreated {
			createdField := map[string]any{
				"id":          fieldIdChecksum("autodate", "created"),
				"type":        "autodate",
				"name":        "created",
				"system":      false,
				"presentable": false,
				"hidden":      false,
				"onCreate":    true,
				"onUpdate":    false,
			}
			c.Schema = append(c.Schema, createdField)
		}

		if addUpdated {
			updatedField := map[string]any{
				"id":          fieldIdChecksum("autodate", "updated"),
				"type":        "autodate",
				"name":        "updated",
				"system":      false,
				"presentable": false,
				"hidden":      false,
				"onCreate":    true,
				"onUpdate":    true,
			}
			c.Schema = append(c.Schema, updatedField)
		}

		if err = txApp.DB().Model(c).Update(); err != nil {
			return err
		}
	}

	_, err = txApp.DB().RenameColumn("_collections", "schema", "fields").Execute()
	if err != nil {
		return err
	}

	// run collection validations
	collections, err := txApp.FindAllCollections()
	if err != nil {
		return fmt.Errorf("failed to retrieve all collections: %w", err)
	}
	for _, c := range collections {
		err = txApp.Validate(c)
		if err != nil {
			return fmt.Errorf("migrated collection %q validation failure: %w", c.Name, err)
		}
	}

	return nil
}

type OldCollectionModel struct {
	Id         string                              `db:"id" json:"id"`
	Created    types.DateTime                      `db:"created" json:"created"`
	Updated    types.DateTime                      `db:"updated" json:"updated"`
	Name       string                              `db:"name" json:"name"`
	Type       string                              `db:"type" json:"type"`
	System     bool                                `db:"system" json:"system"`
	Schema     types.JSONArray[types.JSONMap[any]] `db:"schema" json:"schema"`
	Indexes    types.JSONArray[string]             `db:"indexes" json:"indexes"`
	ListRule   *string                             `db:"listRule" json:"listRule"`
	ViewRule   *string                             `db:"viewRule" json:"viewRule"`
	CreateRule *string                             `db:"createRule" json:"createRule"`
	UpdateRule *string                             `db:"updateRule" json:"updateRule"`
	DeleteRule *string                             `db:"deleteRule" json:"deleteRule"`
	Options    types.JSONMap[any]                  `db:"options" json:"options"`
}

func (c OldCollectionModel) TableName() string {
	return "_collections"
}

func migrateRule(rule *string) *string {
	if rule == nil {
		return nil
	}

	str := strings.ReplaceAll(*rule, "@request.data", "@request.body")

	return &str
}

func toBoolField(data map[string]any) map[string]any {
	return map[string]any{
		"type":        "bool",
		"id":          cast.ToString(data["id"]),
		"name":        cast.ToString(data["name"]),
		"system":      cast.ToBool(data["system"]),
		"required":    cast.ToBool(data["required"]),
		"presentable": cast.ToBool(data["presentable"]),
		"hidden":      false,
	}
}

func toNumberField(data map[string]any) map[string]any {
	return map[string]any{
		"type":        "number",
		"id":          cast.ToString(data["id"]),
		"name":        cast.ToString(data["name"]),
		"system":      cast.ToBool(data["system"]),
		"required":    cast.ToBool(data["required"]),
		"presentable": cast.ToBool(data["presentable"]),
		"hidden":      false,
		"onlyInt":     cast.ToBool(getMapVal(data, "options", "noDecimal")),
		"min":         getMapVal(data, "options", "min"),
		"max":         getMapVal(data, "options", "max"),
	}
}

func toTextField(data map[string]any) map[string]any {
	return map[string]any{
		"type":                "text",
		"id":                  cast.ToString(data["id"]),
		"name":                cast.ToString(data["name"]),
		"system":              cast.ToBool(data["system"]),
		"primaryKey":          cast.ToBool(data["primaryKey"]),
		"hidden":              cast.ToBool(data["hidden"]),
		"presentable":         cast.ToBool(data["presentable"]),
		"required":            cast.ToBool(data["required"]),
		"min":                 cast.ToInt(getMapVal(data, "options", "min")),
		"max":                 cast.ToInt(getMapVal(data, "options", "max")),
		"pattern":             cast.ToString(getMapVal(data, "options", "pattern")),
		"autogeneratePattern": cast.ToString(getMapVal(data, "options", "autogeneratePattern")),
	}
}

func toEmailField(data map[string]any) map[string]any {
	return map[string]any{
		"type":          "email",
		"id":            cast.ToString(data["id"]),
		"name":          cast.ToString(data["name"]),
		"system":        cast.ToBool(data["system"]),
		"required":      cast.ToBool(data["required"]),
		"presentable":   cast.ToBool(data["presentable"]),
		"hidden":        false,
		"exceptDomains": cast.ToStringSlice(getMapVal(data, "options", "exceptDomains")),
		"onlyDomains":   cast.ToStringSlice(getMapVal(data, "options", "onlyDomains")),
	}
}

func toURLField(data map[string]any) map[string]any {
	return map[string]any{
		"type":          "url",
		"id":            cast.ToString(data["id"]),
		"name":          cast.ToString(data["name"]),
		"system":        cast.ToBool(data["system"]),
		"required":      cast.ToBool(data["required"]),
		"presentable":   cast.ToBool(data["presentable"]),
		"hidden":        false,
		"exceptDomains": cast.ToStringSlice(getMapVal(data, "options", "exceptDomains")),
		"onlyDomains":   cast.ToStringSlice(getMapVal(data, "options", "onlyDomains")),
	}
}

func toEditorField(data map[string]any) map[string]any {
	return map[string]any{
		"type":        "editor",
		"id":          cast.ToString(data["id"]),
		"name":        cast.ToString(data["name"]),
		"system":      cast.ToBool(data["system"]),
		"required":    cast.ToBool(data["required"]),
		"presentable": cast.ToBool(data["presentable"]),
		"hidden":      false,
		"convertURLs": cast.ToBool(getMapVal(data, "options", "convertUrls")),
	}
}

func toDateField(data map[string]any) map[string]any {
	return map[string]any{
		"type":        "date",
		"id":          cast.ToString(data["id"]),
		"name":        cast.ToString(data["name"]),
		"system":      cast.ToBool(data["system"]),
		"required":    cast.ToBool(data["required"]),
		"presentable": cast.ToBool(data["presentable"]),
		"hidden":      false,
		"min":         cast.ToString(getMapVal(data, "options", "min")),
		"max":         cast.ToString(getMapVal(data, "options", "max")),
	}
}

func toJSONField(data map[string]any) map[string]any {
	return map[string]any{
		"type":        "json",
		"id":          cast.ToString(data["id"]),
		"name":        cast.ToString(data["name"]),
		"system":      cast.ToBool(data["system"]),
		"required":    cast.ToBool(data["required"]),
		"presentable": cast.ToBool(data["presentable"]),
		"hidden":      false,
		"maxSize":     cast.ToInt64(getMapVal(data, "options", "maxSize")),
	}
}

func toSelectField(data map[string]any) map[string]any {
	return map[string]any{
		"type":        "select",
		"id":          cast.ToString(data["id"]),
		"name":        cast.ToString(data["name"]),
		"system":      cast.ToBool(data["system"]),
		"required":    cast.ToBool(data["required"]),
		"presentable": cast.ToBool(data["presentable"]),
		"hidden":      false,
		"values":      cast.ToStringSlice(getMapVal(data, "options", "values")),
		"maxSelect":   cast.ToInt(getMapVal(data, "options", "maxSelect")),
	}
}

func toRelationField(data map[string]any) map[string]any {
	maxSelect := cast.ToInt(getMapVal(data, "options", "maxSelect"))
	if maxSelect <= 0 {
		maxSelect = 2147483647
	}

	return map[string]any{
		"type":          "relation",
		"id":            cast.ToString(data["id"]),
		"name":          cast.ToString(data["name"]),
		"system":        cast.ToBool(data["system"]),
		"required":      cast.ToBool(data["required"]),
		"presentable":   cast.ToBool(data["presentable"]),
		"hidden":        false,
		"collectionId":  cast.ToString(getMapVal(data, "options", "collectionId")),
		"cascadeDelete": cast.ToBool(getMapVal(data, "options", "cascadeDelete")),
		"minSelect":     cast.ToInt(getMapVal(data, "options", "minSelect")),
		"maxSelect":     maxSelect,
	}
}

func toFileField(data map[string]any) map[string]any {
	return map[string]any{
		"type":        "file",
		"id":          cast.ToString(data["id"]),
		"name":        cast.ToString(data["name"]),
		"system":      cast.ToBool(data["system"]),
		"required":    cast.ToBool(data["required"]),
		"presentable": cast.ToBool(data["presentable"]),
		"hidden":      false,
		"maxSelect":   cast.ToInt(getMapVal(data, "options", "maxSelect")),
		"maxSize":     cast.ToInt64(getMapVal(data, "options", "maxSize")),
		"thumbs":      cast.ToStringSlice(getMapVal(data, "options", "thumbs")),
		"mimeTypes":   cast.ToStringSlice(getMapVal(data, "options", "mimeTypes")),
		"protected":   cast.ToBool(getMapVal(data, "options", "protected")),
	}
}

func getMapVal(m map[string]any, keys ...string) any {
	if len(keys) == 0 {
		return nil
	}

	result, ok := m[keys[0]]
	if !ok {
		return nil
	}

	// end key reached
	if len(keys) == 1 {
		return result
	}

	if m, ok = result.(map[string]any); !ok {
		return nil
	}

	return getMapVal(m, keys[1:]...)
}

func zeroFallback[T comparable](v T, fallback T) T {
	var zero T

	if v == zero {
		return fallback
	}

	return v
}
