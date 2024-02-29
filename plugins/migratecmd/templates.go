package migratecmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/models"
)

const (
	TemplateLangJS = "js"
	TemplateLangGo = "go"

	// note: this usually should be configurable similar to the jsvm plugin,
	// but for simplicity is static as users can easily change the
	// reference path if they use custom dirs structure
	jsTypesDirective = `/// <reference path="../pb_data/types.d.ts" />` + "\n"
)

var emptyTemplateErr = errors.New("empty template")

// -------------------------------------------------------------------
// JavaScript templates
// -------------------------------------------------------------------

func (p *plugin) jsBlankTemplate() (string, error) {
	const template = jsTypesDirective + `migrate((db) => {
  // add up queries...
}, (db) => {
  // add down queries...
})
`

	return template, nil
}

func (p *plugin) jsSnapshotTemplate(collections []*models.Collection) (string, error) {
	jsonData, err := marhshalWithoutEscape(collections, "  ", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = jsTypesDirective + `migrate((db) => {
  const snapshot = %s;

  const collections = snapshot.map((item) => new Collection(item));

  return Dao(db).importCollections(collections, true, null);
}, (db) => {
  return null;
})
`

	return fmt.Sprintf(template, string(jsonData)), nil
}

func (p *plugin) jsCreateTemplate(collection *models.Collection) (string, error) {
	jsonData, err := marhshalWithoutEscape(collection, "  ", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = jsTypesDirective + `migrate((db) => {
  const collection = new Collection(%s);

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId(%q);

  return dao.deleteCollection(collection);
})
`

	return fmt.Sprintf(template, string(jsonData), collection.Id), nil
}

func (p *plugin) jsDeleteTemplate(collection *models.Collection) (string, error) {
	jsonData, err := marhshalWithoutEscape(collection, "  ", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = jsTypesDirective + `migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId(%q);

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection(%s);

  return Dao(db).saveCollection(collection);
})
`

	return fmt.Sprintf(template, collection.Id, string(jsonData)), nil
}

func (p *plugin) jsDiffTemplate(new *models.Collection, old *models.Collection) (string, error) {
	if new == nil && old == nil {
		return "", errors.New("the diff template require at least one of the collection to be non-nil")
	}

	if new == nil {
		return p.jsDeleteTemplate(old)
	}

	if old == nil {
		return p.jsCreateTemplate(new)
	}

	upParts := []string{}
	downParts := []string{}
	varName := "collection"

	if old.Name != new.Name {
		upParts = append(upParts, fmt.Sprintf("%s.name = %q", varName, new.Name))
		downParts = append(downParts, fmt.Sprintf("%s.name = %q", varName, old.Name))
	}

	if old.Type != new.Type {
		upParts = append(upParts, fmt.Sprintf("%s.type = %q", varName, new.Type))
		downParts = append(downParts, fmt.Sprintf("%s.type = %q", varName, old.Type))
	}

	if old.System != new.System {
		upParts = append(upParts, fmt.Sprintf("%s.system = %t", varName, new.System))
		downParts = append(downParts, fmt.Sprintf("%s.system = %t", varName, old.System))
	}

	// ---
	// note: strconv.Quote is used because %q converts the rule operators in unicode char codes
	// ---

	formatRule := func(prop string, rule *string) string {
		if rule == nil {
			return fmt.Sprintf("%s.%s = null", varName, prop)
		}

		return fmt.Sprintf("%s.%s = %s", varName, prop, strconv.Quote(*rule))
	}

	if old.ListRule != new.ListRule {
		oldRule := formatRule("listRule", old.ListRule)
		newRule := formatRule("listRule", new.ListRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	if old.ViewRule != new.ViewRule {
		oldRule := formatRule("viewRule", old.ViewRule)
		newRule := formatRule("viewRule", new.ViewRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	if old.CreateRule != new.CreateRule {
		oldRule := formatRule("createRule", old.CreateRule)
		newRule := formatRule("createRule", new.CreateRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	if old.UpdateRule != new.UpdateRule {
		oldRule := formatRule("updateRule", old.UpdateRule)
		newRule := formatRule("updateRule", new.UpdateRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	if old.DeleteRule != new.DeleteRule {
		oldRule := formatRule("deleteRule", old.DeleteRule)
		newRule := formatRule("deleteRule", new.DeleteRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	// Options
	rawNewOptions, err := marhshalWithoutEscape(new.Options, "  ", "  ")
	if err != nil {
		return "", err
	}
	rawOldOptions, err := marhshalWithoutEscape(old.Options, "  ", "  ")
	if err != nil {
		return "", err
	}
	if !bytes.Equal(rawNewOptions, rawOldOptions) {
		upParts = append(upParts, fmt.Sprintf("%s.options = %s", varName, rawNewOptions))
		downParts = append(downParts, fmt.Sprintf("%s.options = %s", varName, rawOldOptions))
	}

	// Indexes
	rawNewIndexes, err := marhshalWithoutEscape(new.Indexes, "  ", "  ")
	if err != nil {
		return "", err
	}
	rawOldIndexes, err := marhshalWithoutEscape(old.Indexes, "  ", "  ")
	if err != nil {
		return "", err
	}
	if !bytes.Equal(rawNewIndexes, rawOldIndexes) {
		upParts = append(upParts, fmt.Sprintf("%s.indexes = %s", varName, rawNewIndexes))
		downParts = append(downParts, fmt.Sprintf("%s.indexes = %s", varName, rawOldIndexes))
	}

	// ensure new line between regular and collection fields
	if len(upParts) > 0 {
		upParts[len(upParts)-1] += "\n"
	}
	if len(downParts) > 0 {
		downParts[len(downParts)-1] += "\n"
	}

	// Schema
	// -----------------------------------------------------------------

	// deleted fields
	for _, oldField := range old.Schema.Fields() {
		if new.Schema.GetFieldById(oldField.Id) != nil {
			continue // exist
		}

		rawOldField, err := marhshalWithoutEscape(oldField, "  ", "  ")
		if err != nil {
			return "", err
		}

		upParts = append(upParts, "// remove")
		upParts = append(upParts, fmt.Sprintf("%s.schema.removeField(%q)\n", varName, oldField.Id))

		downParts = append(downParts, "// add")
		downParts = append(downParts, fmt.Sprintf("%s.schema.addField(new SchemaField(%s))\n", varName, rawOldField))
	}

	// created fields
	for _, newField := range new.Schema.Fields() {
		if old.Schema.GetFieldById(newField.Id) != nil {
			continue // exist
		}

		rawNewField, err := marhshalWithoutEscape(newField, "  ", "  ")
		if err != nil {
			return "", err
		}

		upParts = append(upParts, "// add")
		upParts = append(upParts, fmt.Sprintf("%s.schema.addField(new SchemaField(%s))\n", varName, rawNewField))

		downParts = append(downParts, "// remove")
		downParts = append(downParts, fmt.Sprintf("%s.schema.removeField(%q)\n", varName, newField.Id))
	}

	// modified fields
	for _, newField := range new.Schema.Fields() {
		oldField := old.Schema.GetFieldById(newField.Id)
		if oldField == nil {
			continue
		}

		rawNewField, err := marhshalWithoutEscape(newField, "  ", "  ")
		if err != nil {
			return "", err
		}

		rawOldField, err := marhshalWithoutEscape(oldField, "  ", "  ")
		if err != nil {
			return "", err
		}

		if bytes.Equal(rawNewField, rawOldField) {
			continue // no change
		}

		upParts = append(upParts, "// update")
		upParts = append(upParts, fmt.Sprintf("%s.schema.addField(new SchemaField(%s))\n", varName, rawNewField))

		downParts = append(downParts, "// update")
		downParts = append(downParts, fmt.Sprintf("%s.schema.addField(new SchemaField(%s))\n", varName, rawOldField))
	}

	// -----------------------------------------------------------------

	if len(upParts) == 0 && len(downParts) == 0 {
		return "", emptyTemplateErr
	}

	up := strings.Join(upParts, "\n  ")
	down := strings.Join(downParts, "\n  ")

	const template = jsTypesDirective + `migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId(%q)

  %s

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId(%q)

  %s

  return dao.saveCollection(collection)
})
`

	return fmt.Sprintf(
		template,
		old.Id, strings.TrimSpace(up),
		new.Id, strings.TrimSpace(down),
	), nil
}

// -------------------------------------------------------------------
// Go templates
// -------------------------------------------------------------------

func (p *plugin) goBlankTemplate() (string, error) {
	const template = `package %s

import (
	"github.com/pocketbase/dbx"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		// add up queries...

		return nil
	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}
`

	return fmt.Sprintf(template, filepath.Base(p.config.Dir)), nil
}

func (p *plugin) goSnapshotTemplate(collections []*models.Collection) (string, error) {
	jsonData, err := marhshalWithoutEscape(collections, "\t\t", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = `package %s

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := ` + "`%s`" + `

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
`
	return fmt.Sprintf(
		template,
		filepath.Base(p.config.Dir),
		escapeBacktick(string(jsonData)),
	), nil
}

func (p *plugin) goCreateTemplate(collection *models.Collection) (string, error) {
	jsonData, err := marhshalWithoutEscape(collection, "\t\t", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = `package %s

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := ` + "`%s`" + `

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId(%q)
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
`

	return fmt.Sprintf(
		template,
		filepath.Base(p.config.Dir),
		escapeBacktick(string(jsonData)),
		collection.Id,
	), nil
}

func (p *plugin) goDeleteTemplate(collection *models.Collection) (string, error) {
	jsonData, err := marhshalWithoutEscape(collection, "\t\t", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = `package %s

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId(%q)
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	}, func(db dbx.Builder) error {
		jsonData := ` + "`%s`" + `

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	})
}
`

	return fmt.Sprintf(
		template,
		filepath.Base(p.config.Dir),
		collection.Id,
		escapeBacktick(string(jsonData)),
	), nil
}

func (p *plugin) goDiffTemplate(new *models.Collection, old *models.Collection) (string, error) {
	if new == nil && old == nil {
		return "", errors.New("the diff template require at least one of the collection to be non-nil")
	}

	if new == nil {
		return p.goDeleteTemplate(old)
	}

	if old == nil {
		return p.goCreateTemplate(new)
	}

	upParts := []string{}
	downParts := []string{}
	varName := "collection"
	if old.Name != new.Name {
		upParts = append(upParts, fmt.Sprintf("%s.Name = %q\n", varName, new.Name))
		downParts = append(downParts, fmt.Sprintf("%s.Name = %q\n", varName, old.Name))
	}

	if old.Type != new.Type {
		upParts = append(upParts, fmt.Sprintf("%s.Type = %q\n", varName, new.Type))
		downParts = append(downParts, fmt.Sprintf("%s.Type = %q\n", varName, old.Type))
	}

	if old.System != new.System {
		upParts = append(upParts, fmt.Sprintf("%s.System = %t\n", varName, new.System))
		downParts = append(downParts, fmt.Sprintf("%s.System = %t\n", varName, old.System))
	}

	// ---
	// note: strconv.Quote is used because %q converts the rule operators in unicode char codes
	// ---

	formatRule := func(prop string, rule *string) string {
		if rule == nil {
			return fmt.Sprintf("%s.%s = nil\n", varName, prop)
		}

		return fmt.Sprintf("%s.%s = types.Pointer(%s)\n", varName, prop, strconv.Quote(*rule))
	}

	if old.ListRule != new.ListRule {
		oldRule := formatRule("ListRule", old.ListRule)
		newRule := formatRule("ListRule", new.ListRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	if old.ViewRule != new.ViewRule {
		oldRule := formatRule("ViewRule", old.ViewRule)
		newRule := formatRule("ViewRule", new.ViewRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	if old.CreateRule != new.CreateRule {
		oldRule := formatRule("CreateRule", old.CreateRule)
		newRule := formatRule("CreateRule", new.CreateRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	if old.UpdateRule != new.UpdateRule {
		oldRule := formatRule("UpdateRule", old.UpdateRule)
		newRule := formatRule("UpdateRule", new.UpdateRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	if old.DeleteRule != new.DeleteRule {
		oldRule := formatRule("DeleteRule", old.DeleteRule)
		newRule := formatRule("DeleteRule", new.DeleteRule)

		if oldRule != newRule {
			upParts = append(upParts, newRule)
			downParts = append(downParts, oldRule)
		}
	}

	// Options
	rawNewOptions, err := marhshalWithoutEscape(new.Options, "\t\t", "\t")
	if err != nil {
		return "", err
	}
	rawOldOptions, err := marhshalWithoutEscape(old.Options, "\t\t", "\t")
	if err != nil {
		return "", err
	}
	if !bytes.Equal(rawNewOptions, rawOldOptions) {
		upParts = append(upParts, "options := map[string]any{}")
		upParts = append(upParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), &options)", escapeBacktick(string(rawNewOptions)))))
		upParts = append(upParts, fmt.Sprintf("%s.SetOptions(options)\n", varName))
		// ---
		downParts = append(downParts, "options := map[string]any{}")
		downParts = append(downParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), &options)", escapeBacktick(string(rawOldOptions)))))
		downParts = append(downParts, fmt.Sprintf("%s.SetOptions(options)\n", varName))
	}

	// Indexes
	rawNewIndexes, err := marhshalWithoutEscape(new.Indexes, "\t\t", "\t")
	if err != nil {
		return "", err
	}
	rawOldIndexes, err := marhshalWithoutEscape(old.Indexes, "\t\t", "\t")
	if err != nil {
		return "", err
	}
	if !bytes.Equal(rawNewIndexes, rawOldIndexes) {
		upParts = append(upParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), &%s.Indexes)", escapeBacktick(string(rawNewIndexes)), varName))+"\n")
		// ---
		downParts = append(downParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), &%s.Indexes)", escapeBacktick(string(rawOldIndexes)), varName))+"\n")
	}

	// Schema
	// ---------------------------------------------------------------
	// deleted fields
	for _, oldField := range old.Schema.Fields() {
		if new.Schema.GetFieldById(oldField.Id) != nil {
			continue // exist
		}

		rawOldField, err := marhshalWithoutEscape(oldField, "\t\t", "\t")
		if err != nil {
			return "", err
		}

		fieldVar := fmt.Sprintf("del_%s", oldField.Name)

		upParts = append(upParts, "// remove")
		upParts = append(upParts, fmt.Sprintf("%s.Schema.RemoveField(%q)\n", varName, oldField.Id))

		downParts = append(downParts, "// add")
		downParts = append(downParts, fmt.Sprintf("%s := &schema.SchemaField{}", fieldVar))
		downParts = append(downParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), %s)", escapeBacktick(string(rawOldField)), fieldVar)))
		downParts = append(downParts, fmt.Sprintf("%s.Schema.AddField(%s)\n", varName, fieldVar))
	}

	// created fields
	for _, newField := range new.Schema.Fields() {
		if old.Schema.GetFieldById(newField.Id) != nil {
			continue // exist
		}

		rawNewField, err := marhshalWithoutEscape(newField, "\t\t", "\t")
		if err != nil {
			return "", err
		}

		fieldVar := fmt.Sprintf("new_%s", newField.Name)

		upParts = append(upParts, "// add")
		upParts = append(upParts, fmt.Sprintf("%s := &schema.SchemaField{}", fieldVar))
		upParts = append(upParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), %s)", escapeBacktick(string(rawNewField)), fieldVar)))
		upParts = append(upParts, fmt.Sprintf("%s.Schema.AddField(%s)\n", varName, fieldVar))

		downParts = append(downParts, "// remove")
		downParts = append(downParts, fmt.Sprintf("%s.Schema.RemoveField(%q)\n", varName, newField.Id))
	}

	// modified fields
	for _, newField := range new.Schema.Fields() {
		oldField := old.Schema.GetFieldById(newField.Id)
		if oldField == nil {
			continue
		}

		rawNewField, err := marhshalWithoutEscape(newField, "\t\t", "\t")
		if err != nil {
			return "", err
		}

		rawOldField, err := marhshalWithoutEscape(oldField, "\t\t", "\t")
		if err != nil {
			return "", err
		}

		if bytes.Equal(rawNewField, rawOldField) {
			continue // no change
		}

		fieldVar := fmt.Sprintf("edit_%s", newField.Name)

		upParts = append(upParts, "// update")
		upParts = append(upParts, fmt.Sprintf("%s := &schema.SchemaField{}", fieldVar))
		upParts = append(upParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), %s)", escapeBacktick(string(rawNewField)), fieldVar)))
		upParts = append(upParts, fmt.Sprintf("%s.Schema.AddField(%s)\n", varName, fieldVar))

		downParts = append(downParts, "// update")
		downParts = append(downParts, fmt.Sprintf("%s := &schema.SchemaField{}", fieldVar))
		downParts = append(downParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), %s)", escapeBacktick(string(rawOldField)), fieldVar)))
		downParts = append(downParts, fmt.Sprintf("%s.Schema.AddField(%s)\n", varName, fieldVar))
	}
	// ---------------------------------------------------------------

	if len(upParts) == 0 && len(downParts) == 0 {
		return "", emptyTemplateErr
	}

	up := strings.Join(upParts, "\n\t\t")
	down := strings.Join(downParts, "\n\t\t")
	combined := up + down

	// generate imports
	// ---
	var imports string

	if strings.Contains(combined, "json.Unmarshal(") ||
		strings.Contains(combined, "json.Marshal(") {
		imports += "\n\t\"encoding/json\"\n"
	}

	imports += "\n\t\"github.com/pocketbase/dbx\""
	imports += "\n\t\"github.com/pocketbase/pocketbase/daos\""
	imports += "\n\tm \"github.com/pocketbase/pocketbase/migrations\""

	if strings.Contains(combined, "schema.SchemaField{") {
		imports += "\n\t\"github.com/pocketbase/pocketbase/models/schema\""
	}

	if strings.Contains(combined, "types.Pointer(") {
		imports += "\n\t\"github.com/pocketbase/pocketbase/tools/types\""
	}
	// ---

	const template = `package %s

import (%s
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId(%q)
		if err != nil {
			return err
		}

		%s

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId(%q)
		if err != nil {
			return err
		}

		%s

		return dao.SaveCollection(collection)
	})
}
`

	return fmt.Sprintf(
		template,
		filepath.Base(p.config.Dir),
		imports,
		old.Id, strings.TrimSpace(up),
		new.Id, strings.TrimSpace(down),
	), nil
}

func marhshalWithoutEscape(v any, prefix string, indent string) ([]byte, error) {
	raw, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return nil, err
	}

	// unescape escaped unicode characters
	unescaped, err := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(raw)), `\\u`, `\u`))
	if err != nil {
		return nil, err
	}

	return []byte(unescaped), nil
}

func escapeBacktick(v string) string {
	return strings.ReplaceAll(v, "`", "` + \"`\" + `")
}

func goErrIf(v string) string {
	return "if err := " + v + "; err != nil {\n\t\t\treturn err\n\t\t}"
}
