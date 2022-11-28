package migratecmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
)

const collectionsCacheKey = "migratecmd_collections"

// onCollectionChange handles the automigration snapshot generation on
// collection change event (create/update/delete).
func (p *plugin) afterCollectionChange() func(*core.ModelEvent) error {
	return func(e *core.ModelEvent) error {
		if e.Model.TableName() != "_collections" {
			return nil // not a collection
		}

		// @todo replace with the OldModel when added to the ModelEvent
		oldCollections, err := p.getCachedCollections()
		if err != nil {
			return err
		}

		old, _ := oldCollections[e.Model.GetId()]

		new, err := p.app.Dao().FindCollectionByNameOrId(e.Model.GetId())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		var template string
		var templateErr error
		if p.options.TemplateLang == TemplateLangJS {
			template, templateErr = p.jsDiffTemplate(new, old)
		} else {
			template, templateErr = p.goDiffTemplate(new, old)
		}
		if templateErr != nil {
			return fmt.Errorf("failed to resolve template: %v", templateErr)
		}

		var action string
		switch {
		case new == nil:
			action = "deleted_" + old.Name
		case old == nil:
			action = "created_" + new.Name
		default:
			action = "updated_" + old.Name
		}

		appliedTime := time.Now().Unix()
		fileDest := filepath.Join(p.options.Dir, fmt.Sprintf("%d_%s.%s", appliedTime, action, p.options.TemplateLang))

		// ensure that the local migrations dir exist
		if err := os.MkdirAll(p.options.Dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create migration dir: %v", err)
		}

		if err := os.WriteFile(fileDest, []byte(template), 0644); err != nil {
			return fmt.Errorf("failed to save automigrate file: %v", err)
		}

		p.refreshCachedCollections()

		return nil
	}
}

func (p *plugin) refreshCachedCollections() error {
	if p.app.Dao() == nil {
		return errors.New("app is not initialized yet")
	}

	var collections []*models.Collection
	if err := p.app.Dao().CollectionQuery().All(&collections); err != nil {
		return err
	}

	mapped := map[string]*models.Collection{}
	for _, c := range collections {
		mapped[c.Id] = c
	}

	p.app.Cache().Set(collectionsCacheKey, mapped)

	return nil
}

func (p *plugin) getCachedCollections() (map[string]*models.Collection, error) {
	if !p.app.Cache().Has(collectionsCacheKey) {
		if err := p.refreshCachedCollections(); err != nil {
			return nil, err
		}
	}

	result, _ := p.app.Cache().Get(collectionsCacheKey).(map[string]*models.Collection)

	return result, nil
}

// getAllMigrationNames return sorted slice with both applied and new
// local migration file names.
func (p *plugin) getAllMigrationNames() ([]string, error) {
	names := []string{}

	for _, migration := range m.AppMigrations.Items() {
		names = append(names, migration.File)
	}

	localFiles, err := p.getLocalMigrationNames()
	if err != nil {
		return nil, err
	}
	for _, name := range localFiles {
		if !list.ExistInSlice(name, names) {
			names = append(names, name)
		}
	}

	sort.Slice(names, func(i int, j int) bool {
		return names[i] < names[j]
	})

	return names, nil
}

// getLocalMigrationNames returns a list with all local migration files
//
// Returns an empty slice if the migrations directory doesn't exist.
func (p *plugin) getLocalMigrationNames() ([]string, error) {
	files, err := os.ReadDir(p.options.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	result := make([]string, 0, len(files))

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		result = append(result, f.Name())
	}

	return result, nil
}
