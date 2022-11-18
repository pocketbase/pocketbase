package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cobra"
)

// Temporary console command to update the pb_data structure to be compatible with the v0.8.0 changes.
//
// NB! It will be removed in v0.9.0!
func NewTempUpgradeCommand(app core.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades your existing pb_data to be compatible with the v0.8.x changes",
		Long: `
Upgrades your existing pb_data to be compatible with the v0.8.x changes
Prerequisites and caveats:
- already upgraded to v0.7.*
- no existing users collection
- existing profiles collection fields like email, username, verified, etc. will be renamed to username2, email2, etc.
`,
		Run: func(command *cobra.Command, args []string) {
			if err := upgrade(app); err != nil {
				color.Red("Error: %v", err)
			}
		},
	}

	return command
}

func upgrade(app core.App) error {
	if _, err := app.Dao().FindCollectionByNameOrId("users"); err == nil {
		return errors.New("It seems that you've already upgraded or have an existing 'users' collection.")
	}

	return app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		if err := migrateCollections(txDao); err != nil {
			return err
		}

		if err := migrateUsers(app, txDao); err != nil {
			return err
		}

		if err := resetMigrationsTable(txDao); err != nil {
			return err
		}

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Println("The pb_data upgrade completed successfully!")
		bold.Println("You can now start the application as usual with the 'serve' command.")
		bold.Println("Please review the migrated collection API rules and fields in the Admin UI and apply the necessary changes in your client-side code.")
		fmt.Println()

		return nil
	})
}

// -------------------------------------------------------------------

func migrateCollections(txDao *daos.Dao) error {
	// add new collection columns
	if _, err := txDao.DB().AddColumn("_collections", "type", "TEXT DEFAULT 'base' NOT NULL").Execute(); err != nil {
		return err
	}
	if _, err := txDao.DB().AddColumn("_collections", "options", "JSON DEFAULT '{}' NOT NULL").Execute(); err != nil {
		return err
	}

	ruleReplacements := []struct {
		old string
		new string
	}{
		{"expand", "expand2"},
		{"collecitonId", "collectionId2"},
		{"collecitonName", "collectionName2"},
		{"profile.userId", "profile.id"},

		// @collection.*
		{"@collection.profiles.userId", "@collection.users.id"},
		{"@collection.profiles.username", "@collection.users.username2"},
		{"@collection.profiles.email", "@collection.users.email2"},
		{"@collection.profiles.emailVisibility", "@collection.users.emailVisibility2"},
		{"@collection.profiles.verified", "@collection.users.verified2"},
		{"@collection.profiles.tokenKey", "@collection.users.tokenKey2"},
		{"@collection.profiles.passwordHash", "@collection.users.passwordHash2"},
		{"@collection.profiles.lastResetSentAt", "@collection.users.lastResetSentAt2"},
		{"@collection.profiles.lastVerificationSentAt", "@collection.users.lastVerificationSentAt2"},
		{"@collection.profiles.", "@collection.users."},

		// @request.*
		{"@request.user.profile.userId", "@request.auth.id"},
		{"@request.user.profile.username", "@request.auth.username2"},
		{"@request.user.profile.email", "@request.auth.email2"},
		{"@request.user.profile.emailVisibility", "@request.auth.emailVisibility2"},
		{"@request.user.profile.verified", "@request.auth.verified2"},
		{"@request.user.profile.tokenKey", "@request.auth.tokenKey2"},
		{"@request.user.profile.passwordHash", "@request.auth.passwordHash2"},
		{"@request.user.profile.lastResetSentAt", "@request.auth.lastResetSentAt2"},
		{"@request.user.profile.lastVerificationSentAt", "@request.auth.lastVerificationSentAt2"},
		{"@request.user.profile.", "@request.auth."},
		{"@request.user", "@request.auth"},
	}

	collections := []*models.Collection{}
	if err := txDao.CollectionQuery().All(&collections); err != nil {
		return err
	}

	for _, collection := range collections {
		collection.Type = models.CollectionTypeBase
		collection.NormalizeOptions()

		// rename profile fields
		// ---
		fieldsToRename := []string{
			"collectionId",
			"collectionName",
			"expand",
		}
		if collection.Name == "profiles" {
			fieldsToRename = append(fieldsToRename,
				"username",
				"email",
				"emailVisibility",
				"verified",
				"tokenKey",
				"passwordHash",
				"lastResetSentAt",
				"lastVerificationSentAt",
			)
		}
		for _, name := range fieldsToRename {
			f := collection.Schema.GetFieldByName(name)
			if f != nil {
				color.Blue("[%s - renamed field]", collection.Name)
				color.Yellow(" - old: %s", f.Name)
				color.Green(" - new: %s2", f.Name)
				fmt.Println()
				f.Name += "2"
			}
		}
		// ---

		// replace rule fields
		// ---
		rules := map[string]*string{
			"ListRule":   collection.ListRule,
			"ViewRule":   collection.ViewRule,
			"CreateRule": collection.CreateRule,
			"UpdateRule": collection.UpdateRule,
			"DeleteRule": collection.DeleteRule,
		}

		for ruleKey, rule := range rules {
			if rule == nil || *rule == "" {
				continue
			}

			originalRule := *rule

			for _, replacement := range ruleReplacements {
				re := regexp.MustCompile(regexp.QuoteMeta(replacement.old) + `\b`)
				*rule = re.ReplaceAllString(*rule, replacement.new)
			}

			*rule = replaceReversedLikes(*rule)

			if originalRule != *rule {
				color.Blue("[%s - replaced %s]:", collection.Name, ruleKey)
				color.Yellow(" - old: %s", strings.TrimSpace(originalRule))
				color.Green(" - new: %s", strings.TrimSpace(*rule))
				fmt.Println()
			}
		}
		// ---

		if err := txDao.SaveCollection(collection); err != nil {
			return err
		}
	}

	return nil
}

func migrateUsers(app core.App, txDao *daos.Dao) error {
	color.Blue(`[merging "_users" and "profiles"]:`)

	profilesCollection, err := txDao.FindCollectionByNameOrId("profiles")
	if err != nil {
		return err
	}

	originalProfilesCollectionId := profilesCollection.Id

	// change the profiles collection id to something else since we will be using
	// it for the new users collection in order to avoid renaming the storage dir
	_, idRenameErr := txDao.DB().NewQuery(fmt.Sprintf(
		`UPDATE {{_collections}}
			SET id = '%s'
			WHERE id = '%s';
		`,
		(originalProfilesCollectionId + "__old__"),
		originalProfilesCollectionId,
	)).Execute()
	if idRenameErr != nil {
		return idRenameErr
	}

	// refresh profiles collection
	profilesCollection, err = txDao.FindCollectionByNameOrId("profiles")
	if err != nil {
		return err
	}

	usersSchema, _ := profilesCollection.Schema.Clone()
	userIdField := usersSchema.GetFieldByName("userId")
	if userIdField != nil {
		usersSchema.RemoveField(userIdField.Id)
	}

	usersCollection := &models.Collection{}
	usersCollection.MarkAsNew()
	usersCollection.Id = originalProfilesCollectionId
	usersCollection.Name = "users"
	usersCollection.Type = models.CollectionTypeAuth
	usersCollection.Schema = *usersSchema
	usersCollection.CreateRule = types.Pointer("")
	if profilesCollection.ListRule != nil && *profilesCollection.ListRule != "" {
		*profilesCollection.ListRule = strings.ReplaceAll(*profilesCollection.ListRule, "userId", "id")
		usersCollection.ListRule = profilesCollection.ListRule
	}
	if profilesCollection.ViewRule != nil && *profilesCollection.ViewRule != "" {
		*profilesCollection.ViewRule = strings.ReplaceAll(*profilesCollection.ViewRule, "userId", "id")
		usersCollection.ViewRule = profilesCollection.ViewRule
	}
	if profilesCollection.UpdateRule != nil && *profilesCollection.UpdateRule != "" {
		*profilesCollection.UpdateRule = strings.ReplaceAll(*profilesCollection.UpdateRule, "userId", "id")
		usersCollection.UpdateRule = profilesCollection.UpdateRule
	}
	if profilesCollection.DeleteRule != nil && *profilesCollection.DeleteRule != "" {
		*profilesCollection.DeleteRule = strings.ReplaceAll(*profilesCollection.DeleteRule, "userId", "id")
		usersCollection.DeleteRule = profilesCollection.DeleteRule
	}

	// set auth options
	settings := app.Settings()
	authOptions := usersCollection.AuthOptions()
	authOptions.ManageRule = nil
	authOptions.AllowOAuth2Auth = true
	authOptions.AllowUsernameAuth = false
	authOptions.AllowEmailAuth = settings.EmailAuth.Enabled
	authOptions.MinPasswordLength = settings.EmailAuth.MinPasswordLength
	authOptions.OnlyEmailDomains = settings.EmailAuth.OnlyDomains
	authOptions.ExceptEmailDomains = settings.EmailAuth.ExceptDomains
	// twitter currently is the only provider that doesn't return an email
	authOptions.RequireEmail = !settings.TwitterAuth.Enabled

	usersCollection.SetOptions(authOptions)

	if err := txDao.SaveCollection(usersCollection); err != nil {
		return err
	}

	// copy the original users
	_, usersErr := txDao.DB().NewQuery(`
		INSERT INTO {{users}} (id, created, updated, username, email, emailVisibility, verified, tokenKey, passwordHash, lastResetSentAt, lastVerificationSentAt)
		SELECT id, created, updated, ("u_" || id), email, false, verified, tokenKey, passwordHash, lastResetSentAt, lastVerificationSentAt
		FROM {{_users}};
	`).Execute()
	if usersErr != nil {
		return usersErr
	}

	// generate the profile fields copy statements
	sets := []string{"id = p.id"}
	for _, f := range usersSchema.Fields() {
		sets = append(sets, fmt.Sprintf("%s = p.%s", f.Name, f.Name))
	}

	// copy profile fields
	_, copyProfileErr := txDao.DB().NewQuery(fmt.Sprintf(`
		UPDATE {{users}} as u
		SET %s
		FROM {{profiles}} as p
		WHERE u.id = p.userId;
	`, strings.Join(sets, ", "))).Execute()
	if copyProfileErr != nil {
		return copyProfileErr
	}

	profileRecords, err := txDao.FindRecordsByExpr("profiles")
	if err != nil {
		return err
	}

	// update all profiles and users fields to point to the new users collection
	collections := []*models.Collection{}
	if err := txDao.CollectionQuery().All(&collections); err != nil {
		return err
	}
	for _, collection := range collections {
		var hasChanges bool

		for _, f := range collection.Schema.Fields() {
			f.InitOptions()

			if f.Type == schema.FieldTypeUser {
				if collection.Name == "profiles" && f.Name == "userId" {
					continue
				}

				hasChanges = true

				// change the user field to a relation field
				options, _ := f.Options.(*schema.UserOptions)
				f.Type = schema.FieldTypeRelation
				f.Options = &schema.RelationOptions{
					CollectionId:  usersCollection.Id,
					MaxSelect:     &options.MaxSelect,
					CascadeDelete: options.CascadeDelete,
				}

				for _, p := range profileRecords {
					pId := p.Id
					pUserId := p.GetString("userId")
					// replace all user record id references with the profile id
					_, replaceErr := txDao.DB().NewQuery(fmt.Sprintf(`
						UPDATE %s
						SET [[%s]] = REPLACE([[%s]], '%s', '%s')
						WHERE [[%s]] LIKE ('%%%s%%');
					`, collection.Name, f.Name, f.Name, pUserId, pId, f.Name, pUserId)).Execute()
					if replaceErr != nil {
						return replaceErr
					}
				}
			}
		}

		if hasChanges {
			if err := txDao.Save(collection); err != nil {
				return err
			}
		}
	}

	if err := migrateExternalAuths(txDao, originalProfilesCollectionId); err != nil {
		return err
	}

	// drop _users table
	if _, err := txDao.DB().DropTable("_users").Execute(); err != nil {
		return err
	}

	// drop profiles table
	if _, err := txDao.DB().DropTable("profiles").Execute(); err != nil {
		return err
	}

	// delete profiles collection
	if err := txDao.Delete(profilesCollection); err != nil {
		return err
	}

	color.Green(` - Successfully merged "_users" and "profiles" into a new collection "users".`)
	fmt.Println()

	return nil
}

func migrateExternalAuths(txDao *daos.Dao, userCollectionId string) error {
	_, alterErr := txDao.DB().NewQuery(`
		-- crate new externalAuths table
		CREATE TABLE {{_newExternalAuths}} (
			[[id]]           TEXT PRIMARY KEY,
			[[collectionId]] TEXT NOT NULL,
			[[recordId]]     TEXT NOT NULL,
			[[provider]]     TEXT NOT NULL,
			[[providerId]]   TEXT NOT NULL,
			[[created]]      TEXT DEFAULT "" NOT NULL,
			[[updated]]      TEXT DEFAULT "" NOT NULL,
			---
			FOREIGN KEY ([[collectionId]]) REFERENCES {{_collections}} ([[id]]) ON UPDATE CASCADE ON DELETE CASCADE
		);

		-- copy all data from the old table to the new one
		INSERT INTO {{_newExternalAuths}}
		SELECT auth.id, "` + userCollectionId + `" as collectionId, [[profiles.id]] as recordId, auth.provider, auth.providerId, auth.created, auth.updated
		FROM {{_externalAuths}} auth
		INNER JOIN {{profiles}} on [[profiles.userId]] = [[auth.userId]];

		-- drop old table
		DROP TABLE {{_externalAuths}};

		-- rename new table
		ALTER TABLE {{_newExternalAuths}} RENAME TO {{_externalAuths}};

		-- create named indexes
		CREATE UNIQUE INDEX _externalAuths_record_provider_idx on {{_externalAuths}} ([[collectionId]], [[recordId]], [[provider]]);
		CREATE UNIQUE INDEX _externalAuths_provider_providerId_idx on {{_externalAuths}} ([[provider]], [[providerId]]);
	`).Execute()

	return alterErr
}

func resetMigrationsTable(txDao *daos.Dao) error {
	// reset the migration state to the new init
	_, err := txDao.DB().Delete("_migrations", dbx.HashExp{
		"file": "1661586591_add_externalAuths_table.go",
	}).Execute()

	return err
}

var reverseLikeRegex = regexp.MustCompile(`(['"]\w*['"])\s*(\~|!~)\s*([\w\@\.]*)`)

func replaceReversedLikes(rule string) string {
	parts := reverseLikeRegex.FindAllStringSubmatch(rule, -1)

	for _, p := range parts {
		if len(p) != 4 {
			continue
		}

		newPart := fmt.Sprintf("%s %s %s", p[3], p[2], p[1])

		rule = strings.ReplaceAll(rule, p[0], newPart)
	}

	return rule
}
