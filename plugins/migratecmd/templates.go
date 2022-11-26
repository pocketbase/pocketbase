package migratecmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/pocketbase/pocketbase/models"
)

const (
	TemplateLangJS = "js"
	TemplateLangGo = "go"
)

// -------------------------------------------------------------------
// JavaScript templates
// -------------------------------------------------------------------

func (p *plugin) jsCreateTemplate() (string, error) {
	const template = `migrate((db) => {
  // add up queries...
}, (db) => {
  // add down queries...
})
`

	return template, nil
}

func (p *plugin) jsSnapshotTemplate(collections []*models.Collection) (string, error) {
	jsonData, err := json.MarshalIndent(collections, "  ", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %v", err)
	}

	const template = `migrate((db) => {
  const snapshot = %s;

  const collections = snapshot.map((item) => unmarshal(item, new Collection()));

  return Dao(db).importCollections(collections, true, null);
}, (db) => {
  return null;
})
`

	return fmt.Sprintf(template, string(jsonData)), nil
}

// -------------------------------------------------------------------
// Go templates
// -------------------------------------------------------------------

func (p *plugin) goCreateTemplate() (string, error) {
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

	return fmt.Sprintf(template, filepath.Base(p.options.Dir)), nil
}

func (p *plugin) goSnapshotTemplate(collections []*models.Collection) (string, error) {
	jsonData, err := json.MarshalIndent(collections, "\t", "\t\t")
	if err != nil {
		return "", fmt.Errorf("failed to serialize collections list: %v", err)
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
	return fmt.Sprintf(template, filepath.Base(p.options.Dir), string(jsonData)), nil
}
