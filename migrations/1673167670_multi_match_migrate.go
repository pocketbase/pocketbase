package migrations

import (
	"regexp"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

// This migration replaces for backward compatibility the default operators
// (=, !=, >, etc.) with their any/opt equivalent (?=, ?=, ?>, etc.)
// in any muli-rel expression collection rule.
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		exprRegex := regexp.MustCompile(`([\@\'\"\w\.]+)\s*(=|!=|~|!~|>|>=|<|<=)\s*([\@\'\"\w\.]+)`)

		collections := []*models.Collection{}
		if err := dao.CollectionQuery().All(&collections); err != nil {
			return err
		}

		findCollection := func(nameOrId string) *models.Collection {
			for _, c := range collections {
				if c.Id == nameOrId || c.Name == nameOrId {
					return c
				}
			}

			return nil
		}

		var isMultiRelLiteral func(mainCollection *models.Collection, literal string) bool
		isMultiRelLiteral = func(mainCollection *models.Collection, literal string) bool {
			if strings.HasPrefix(literal, "@collection.") {
				return true
			}

			if strings.HasPrefix(literal, `"`) ||
				strings.HasPrefix(literal, `'`) ||
				strings.HasPrefix(literal, "@request.method") ||
				strings.HasPrefix(literal, "@request.data") ||
				strings.HasPrefix(literal, "@request.query") {
				return false
			}

			parts := strings.Split(literal, ".")
			if len(parts) <= 1 {
				return false
			}

			if strings.HasPrefix(literal, "@request.auth") && len(parts) >= 4 {
				// check each auth collection
				for _, c := range collections {
					if c.IsAuth() && isMultiRelLiteral(c, strings.Join(parts[2:], ".")) {
						return true
					}
				}

				return false
			}

			activeCollection := mainCollection

			for i, p := range parts {
				f := activeCollection.Schema.GetFieldByName(p)
				if f == nil || f.Type != schema.FieldTypeRelation {
					return false // not a  relation field
				}

				// is multi-relation and not the last prop
				opt, ok := f.Options.(*schema.RelationOptions)
				if ok && (opt.MaxSelect == nil || *opt.MaxSelect != 1) && i != len(parts)-1 {
					return true
				}

				activeCollection = findCollection(opt.CollectionId)
				if activeCollection == nil {
					return false
				}
			}

			return false
		}

		// replace all multi-match operators to their any/opt equivalent, eg. "=" => "?="
		migrateRule := func(collection *models.Collection, rule *string) (*string, error) {
			if rule == nil || *rule == "" {
				return rule, nil
			}

			newRule := *rule
			parts := exprRegex.FindAllStringSubmatch(newRule, -1)

			for _, p := range parts {
				if isMultiRelLiteral(collection, p[1]) || isMultiRelLiteral(collection, p[3]) {
					newRule = strings.ReplaceAll(newRule, p[0], p[1]+" ?"+p[2]+" "+p[3])
				}
			}

			return &newRule, nil
		}

		var ruleErr error
		for _, c := range collections {
			c.ListRule, ruleErr = migrateRule(c, c.ListRule)
			if ruleErr != nil {
				return ruleErr
			}

			c.ViewRule, ruleErr = migrateRule(c, c.ViewRule)
			if ruleErr != nil {
				return ruleErr
			}

			c.CreateRule, ruleErr = migrateRule(c, c.CreateRule)
			if ruleErr != nil {
				return ruleErr
			}

			c.UpdateRule, ruleErr = migrateRule(c, c.UpdateRule)
			if ruleErr != nil {
				return ruleErr
			}

			c.DeleteRule, ruleErr = migrateRule(c, c.DeleteRule)
			if ruleErr != nil {
				return ruleErr
			}

			if c.IsAuth() {
				opt := c.AuthOptions()
				opt.ManageRule, ruleErr = migrateRule(c, opt.ManageRule)
				if ruleErr != nil {
					return ruleErr
				}
				c.SetOptions(opt)
			}

			if err := dao.Save(c); err != nil {
				return err
			}
		}

		return nil
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collections := []*models.Collection{}
		if err := dao.CollectionQuery().All(&collections); err != nil {
			return err
		}

		anyOpRegex := regexp.MustCompile(`\?(=|!=|~|!~|>|>=|<|<=)`)

		// replace any/opt operators to their old versions, eg. "?=" => "="
		revertRule := func(rule *string) (*string, error) {
			if rule == nil || *rule == "" {
				return rule, nil
			}

			newRule := *rule
			newRule = anyOpRegex.ReplaceAllString(newRule, "${1}")

			return &newRule, nil
		}

		var ruleErr error
		for _, c := range collections {
			c.ListRule, ruleErr = revertRule(c.ListRule)
			if ruleErr != nil {
				return ruleErr
			}

			c.ViewRule, ruleErr = revertRule(c.ViewRule)
			if ruleErr != nil {
				return ruleErr
			}

			c.CreateRule, ruleErr = revertRule(c.CreateRule)
			if ruleErr != nil {
				return ruleErr
			}

			c.UpdateRule, ruleErr = revertRule(c.UpdateRule)
			if ruleErr != nil {
				return ruleErr
			}

			c.DeleteRule, ruleErr = revertRule(c.DeleteRule)
			if ruleErr != nil {
				return ruleErr
			}

			if c.IsAuth() {
				opt := c.AuthOptions()
				opt.ManageRule, ruleErr = revertRule(opt.ManageRule)
				if ruleErr != nil {
					return ruleErr
				}
				c.SetOptions(opt)
			}

			if err := dao.Save(c); err != nil {
				return err
			}
		}

		return nil
	})
}
