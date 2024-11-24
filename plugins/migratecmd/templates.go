package migratecmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

const (
	TemplateLangJS = "js"
	TemplateLangGo = "go"

	// note: this usually should be configurable similar to the jsvm plugin,
	// but for simplicity is static as users can easily change the
	// reference path if they use custom dirs structure
	jsTypesDirective = `/// <reference path="../pb_data/types.d.ts" />` + "\n"
)

var ErrEmptyTemplate = errors.New("empty template")

// -------------------------------------------------------------------
// JavaScript templates
// -------------------------------------------------------------------

func (p *plugin) jsBlankTemplate() (string, error) {
	const template = jsTypesDirective + `migrate((app) => {
  // add up queries...
}, (app) => {
  // add down queries...
})
`

	return template, nil
}

func (p *plugin) jsSnapshotTemplate(collections []*core.Collection) (string, error) {
	// unset timestamp fields
	var collectionsData = make([]map[string]any, len(collections))
	for i, c := range collections {
		data, err := toMap(c)
		if err != nil {
			return "", fmt.Errorf("failed to serialize %q into a map: %w", c.Name, err)
		}
		delete(data, "created")
		delete(data, "updated")
		deleteNestedMapKey(data, "oauth2", "providers")
		collectionsData[i] = data
	}

	jsonData, err := marhshalWithoutEscape(collectionsData, "  ", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = jsTypesDirective + `migrate((app) => {
  const snapshot = %s;

  return app.importCollections(snapshot, false);
}, (app) => {
  return null;
})
`

	return fmt.Sprintf(template, string(jsonData)), nil
}

func (p *plugin) jsCreateTemplate(collection *core.Collection) (string, error) {
	// unset timestamp fields
	collectionData, err := toMap(collection)
	if err != nil {
		return "", err
	}
	delete(collectionData, "created")
	delete(collectionData, "updated")
	deleteNestedMapKey(collectionData, "oauth2", "providers")

	jsonData, err := marhshalWithoutEscape(collectionData, "  ", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collection: %w", err)
	}

	const template = jsTypesDirective + `migrate((app) => {
  const collection = new Collection(%s);

  return app.save(collection);
}, (app) => {
  const collection = app.findCollectionByNameOrId(%q);

  return app.delete(collection);
})
`

	return fmt.Sprintf(template, string(jsonData), collection.Id), nil
}

func (p *plugin) jsDeleteTemplate(collection *core.Collection) (string, error) {
	// unset timestamp fields
	collectionData, err := toMap(collection)
	if err != nil {
		return "", err
	}
	delete(collectionData, "created")
	delete(collectionData, "updated")
	deleteNestedMapKey(collectionData, "oauth2", "providers")

	jsonData, err := marhshalWithoutEscape(collectionData, "  ", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = jsTypesDirective + `migrate((app) => {
  const collection = app.findCollectionByNameOrId(%q);

  return app.delete(collection);
}, (app) => {
  const collection = new Collection(%s);

  return app.save(collection);
})
`

	return fmt.Sprintf(template, collection.Id, string(jsonData)), nil
}

func (p *plugin) jsDiffTemplate(new *core.Collection, old *core.Collection) (string, error) {
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

	newMap, err := toMap(new)
	if err != nil {
		return "", err
	}

	oldMap, err := toMap(old)
	if err != nil {
		return "", err
	}

	// non-fields
	// -----------------------------------------------------------------

	upDiff := diffMaps(oldMap, newMap, "fields", "created", "updated")
	if len(upDiff) > 0 {
		downDiff := diffMaps(newMap, oldMap, "fields", "created", "updated")

		rawUpDiff, err := marhshalWithoutEscape(upDiff, "  ", "  ")
		if err != nil {
			return "", err
		}

		rawDownDiff, err := marhshalWithoutEscape(downDiff, "  ", "  ")
		if err != nil {
			return "", err
		}

		upParts = append(upParts, "// update collection data")
		upParts = append(upParts, fmt.Sprintf("unmarshal(%s, %s)", string(rawUpDiff), varName)+"\n")
		// ---
		downParts = append(downParts, "// update collection data")
		downParts = append(downParts, fmt.Sprintf("unmarshal(%s, %s)", string(rawDownDiff), varName)+"\n")
	}

	// fields
	// -----------------------------------------------------------------

	oldFieldsSlice, ok := oldMap["fields"].([]any)
	if !ok {
		return "", errors.New(`oldMap["fields"] is not []any`)
	}

	newFieldsSlice, ok := newMap["fields"].([]any)
	if !ok {
		return "", errors.New(`newMap["fields"] is not []any`)
	}

	// deleted fields
	for i, oldField := range old.Fields {
		if new.Fields.GetById(oldField.GetId()) != nil {
			continue // exist
		}

		rawOldField, err := marhshalWithoutEscape(oldFieldsSlice[i], "  ", "  ")
		if err != nil {
			return "", err
		}

		upParts = append(upParts, "// remove field")
		upParts = append(upParts, fmt.Sprintf("%s.fields.removeById(%q)\n", varName, oldField.GetId()))

		downParts = append(downParts, "// add field")
		downParts = append(downParts, fmt.Sprintf("%s.fields.addAt(%d, new Field(%s))\n", varName, i, rawOldField))
	}

	// created fields
	for i, newField := range new.Fields {
		if old.Fields.GetById(newField.GetId()) != nil {
			continue // exist
		}

		rawNewField, err := marhshalWithoutEscape(newFieldsSlice[i], "  ", "  ")
		if err != nil {
			return "", err
		}

		upParts = append(upParts, "// add field")
		upParts = append(upParts, fmt.Sprintf("%s.fields.addAt(%d, new Field(%s))\n", varName, i, rawNewField))

		downParts = append(downParts, "// remove field")
		downParts = append(downParts, fmt.Sprintf("%s.fields.removeById(%q)\n", varName, newField.GetId()))
	}

	// modified fields
	// (note currently ignoring order-only changes as it comes with too many edge-cases)
	for i, newField := range new.Fields {
		var rawNewField, rawOldField []byte

		rawNewField, err = marhshalWithoutEscape(newFieldsSlice[i], "  ", "  ")
		if err != nil {
			return "", err
		}

		var oldFieldIndex int

		for j, oldField := range old.Fields {
			if oldField.GetId() == newField.GetId() {
				rawOldField, err = marhshalWithoutEscape(oldFieldsSlice[j], "  ", "  ")
				if err != nil {
					return "", err
				}
				oldFieldIndex = j
				break
			}
		}

		if rawOldField == nil || bytes.Equal(rawNewField, rawOldField) {
			continue // new field or no change
		}

		upParts = append(upParts, "// update field")
		upParts = append(upParts, fmt.Sprintf("%s.fields.addAt(%d, new Field(%s))\n", varName, i, rawNewField))

		downParts = append(downParts, "// update field")
		downParts = append(downParts, fmt.Sprintf("%s.fields.addAt(%d, new Field(%s))\n", varName, oldFieldIndex, rawOldField))
	}

	// -----------------------------------------------------------------

	if len(upParts) == 0 && len(downParts) == 0 {
		return "", ErrEmptyTemplate
	}

	up := strings.Join(upParts, "\n  ")
	down := strings.Join(downParts, "\n  ")

	const template = jsTypesDirective + `migrate((app) => {
  const collection = app.findCollectionByNameOrId(%q)

  %s

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId(%q)

  %s

  return app.save(collection)
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
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// add up queries...

		return nil
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}
`

	return fmt.Sprintf(template, filepath.Base(p.config.Dir)), nil
}

func (p *plugin) goSnapshotTemplate(collections []*core.Collection) (string, error) {
	// unset timestamp fields
	var collectionsData = make([]map[string]any, len(collections))
	for i, c := range collections {
		data, err := toMap(c)
		if err != nil {
			return "", fmt.Errorf("failed to serialize %q into a map: %w", c.Name, err)
		}
		delete(data, "created")
		delete(data, "updated")
		deleteNestedMapKey(data, "oauth2", "providers")
		collectionsData[i] = data
	}

	jsonData, err := marhshalWithoutEscape(collectionsData, "\t\t", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = `package %s

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := ` + "`%s`" + `

		return app.ImportCollectionsByMarshaledJSON([]byte(jsonData), false)
	}, func(app core.App) error {
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

func (p *plugin) goCreateTemplate(collection *core.Collection) (string, error) {
	// unset timestamp fields
	collectionData, err := toMap(collection)
	if err != nil {
		return "", err
	}
	delete(collectionData, "created")
	delete(collectionData, "updated")
	deleteNestedMapKey(collectionData, "oauth2", "providers")

	jsonData, err := marhshalWithoutEscape(collectionData, "\t\t", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = `package %s

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := ` + "`%s`" + `

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(%q)
		if err != nil {
			return err
		}

		return app.Delete(collection)
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

func (p *plugin) goDeleteTemplate(collection *core.Collection) (string, error) {
	// unset timestamp fields
	collectionData, err := toMap(collection)
	if err != nil {
		return "", err
	}
	delete(collectionData, "created")
	delete(collectionData, "updated")
	deleteNestedMapKey(collectionData, "oauth2", "providers")

	jsonData, err := marhshalWithoutEscape(collectionData, "\t\t", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %w", err)
	}

	const template = `package %s

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(%q)
		if err != nil {
			return err
		}

		return app.Delete(collection)
	}, func(app core.App) error {
		jsonData := ` + "`%s`" + `

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
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

func (p *plugin) goDiffTemplate(new *core.Collection, old *core.Collection) (string, error) {
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

	newMap, err := toMap(new)
	if err != nil {
		return "", err
	}

	oldMap, err := toMap(old)
	if err != nil {
		return "", err
	}

	// non-fields
	// -----------------------------------------------------------------

	upDiff := diffMaps(oldMap, newMap, "fields", "created", "updated")
	if len(upDiff) > 0 {
		downDiff := diffMaps(newMap, oldMap, "fields", "created", "updated")

		rawUpDiff, err := marhshalWithoutEscape(upDiff, "\t\t", "\t")
		if err != nil {
			return "", err
		}

		rawDownDiff, err := marhshalWithoutEscape(downDiff, "\t\t", "\t")
		if err != nil {
			return "", err
		}

		upParts = append(upParts, "// update collection data")
		upParts = append(upParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), &%s)", escapeBacktick(string(rawUpDiff)), varName)))
		// ---
		downParts = append(downParts, "// update collection data")
		downParts = append(downParts, goErrIf(fmt.Sprintf("json.Unmarshal([]byte(`%s`), &%s)", escapeBacktick(string(rawDownDiff)), varName)))
	}

	// fields
	// -----------------------------------------------------------------

	oldFieldsSlice, ok := oldMap["fields"].([]any)
	if !ok {
		return "", errors.New(`oldMap["fields"] is not []any`)
	}

	newFieldsSlice, ok := newMap["fields"].([]any)
	if !ok {
		return "", errors.New(`newMap["fields"] is not []any`)
	}

	// deleted fields
	for i, oldField := range old.Fields {
		if new.Fields.GetById(oldField.GetId()) != nil {
			continue // exist
		}

		rawOldField, err := marhshalWithoutEscape(oldFieldsSlice[i], "\t\t", "\t")
		if err != nil {
			return "", err
		}

		upParts = append(upParts, "// remove field")
		upParts = append(upParts, fmt.Sprintf("%s.Fields.RemoveById(%q)\n", varName, oldField.GetId()))

		downParts = append(downParts, "// add field")
		downParts = append(downParts, goErrIf(fmt.Sprintf("%s.Fields.AddMarshaledJSONAt(%d, []byte(`%s`))", varName, i, escapeBacktick(string(rawOldField)))))
	}

	// created fields
	for i, newField := range new.Fields {
		if old.Fields.GetById(newField.GetId()) != nil {
			continue // exist
		}

		rawNewField, err := marhshalWithoutEscape(newFieldsSlice[i], "\t\t", "\t")
		if err != nil {
			return "", err
		}

		upParts = append(upParts, "// add field")
		upParts = append(upParts, goErrIf(fmt.Sprintf("%s.Fields.AddMarshaledJSONAt(%d, []byte(`%s`))", varName, i, escapeBacktick(string(rawNewField)))))

		downParts = append(downParts, "// remove field")
		downParts = append(downParts, fmt.Sprintf("%s.Fields.RemoveById(%q)\n", varName, newField.GetId()))
	}

	// modified fields
	// (note currently ignoring order-only changes as it comes with too many edge-cases)
	for i, newField := range new.Fields {
		var rawNewField, rawOldField []byte

		rawNewField, err = marhshalWithoutEscape(newFieldsSlice[i], "\t\t", "\t")
		if err != nil {
			return "", err
		}

		var oldFieldIndex int

		for j, oldField := range old.Fields {
			if oldField.GetId() == newField.GetId() {
				rawOldField, err = marhshalWithoutEscape(oldFieldsSlice[j], "\t\t", "\t")
				if err != nil {
					return "", err
				}
				oldFieldIndex = j
				break
			}
		}

		if rawOldField == nil || bytes.Equal(rawNewField, rawOldField) {
			continue // new field or no change
		}

		upParts = append(upParts, "// update field")
		upParts = append(upParts, goErrIf(fmt.Sprintf("%s.Fields.AddMarshaledJSONAt(%d, []byte(`%s`))", varName, i, escapeBacktick(string(rawNewField)))))

		downParts = append(downParts, "// update field")
		downParts = append(downParts, goErrIf(fmt.Sprintf("%s.Fields.AddMarshaledJSONAt(%d, []byte(`%s`))", varName, oldFieldIndex, escapeBacktick(string(rawOldField)))))
	}

	// ---------------------------------------------------------------

	if len(upParts) == 0 && len(downParts) == 0 {
		return "", ErrEmptyTemplate
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

	imports += "\n\t\"github.com/pocketbase/pocketbase/core\""
	imports += "\n\tm \"github.com/pocketbase/pocketbase/migrations\""
	// ---

	const template = `package %s

import (%s
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(%q)
		if err != nil {
			return err
		}

		%s

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId(%q)
		if err != nil {
			return err
		}

		%s

		return app.Save(collection)
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
	return "if err := " + v + "; err != nil {\n\t\t\treturn err\n\t\t}\n"
}

func toMap(v any) (map[string]any, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	result := map[string]any{}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func diffMaps(old, new map[string]any, excludeKeys ...string) map[string]any {
	diff := map[string]any{}

	for k, vNew := range new {
		if slices.Contains(excludeKeys, k) {
			continue
		}

		vOld, ok := old[k]
		if !ok {
			// new field
			diff[k] = vNew
			continue
		}

		// compare the serialized version of the values in case of slice or other custom type
		rawOld, _ := json.Marshal(vOld)
		rawNew, _ := json.Marshal(vNew)

		if !bytes.Equal(rawOld, rawNew) {
			// if both are maps add recursively only the changed fields
			vOldMap, ok1 := vOld.(map[string]any)
			vNewMap, ok2 := vNew.(map[string]any)
			if ok1 && ok2 {
				subDiff := diffMaps(vOldMap, vNewMap)
				if len(subDiff) > 0 {
					diff[k] = subDiff
				}
			} else {
				diff[k] = vNew
			}
		}
	}

	// unset missing fields
	for k := range old {
		if _, ok := diff[k]; ok || slices.Contains(excludeKeys, k) {
			continue // already added
		}

		if _, ok := new[k]; !ok {
			diff[k] = nil
		}
	}

	return diff
}

func deleteNestedMapKey(data map[string]any, parts ...string) {
	if len(parts) == 0 {
		return
	}

	if len(parts) == 1 {
		delete(data, parts[0])
		return
	}

	v, ok := data[parts[0]].(map[string]any)
	if ok {
		deleteNestedMapKey(v, parts[1:]...)
	}
}
