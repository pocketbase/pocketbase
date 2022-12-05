package migratecmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/migrate"
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

		old := oldCollections[e.Model.GetId()]

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
			if errors.Is(templateErr, emptyTemplateErr) {
				return nil // no changes
			}
			return fmt.Errorf("failed to resolve template: %w", templateErr)
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
		name := fmt.Sprintf("%d_%s.%s", appliedTime, action, p.options.TemplateLang)
		filePath := filepath.Join(p.options.Dir, name)

		return p.app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
			// insert the migration entry
			_, err := txDao.DB().Insert(migrate.DefaultMigrationsTable, dbx.Params{
				"file":    name,
				"applied": appliedTime,
			}).Execute()
			if err != nil {
				return err
			}

			// ensure that the local migrations dir exist
			if err := os.MkdirAll(p.options.Dir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create migration dir: %w", err)
			}

			if err := os.WriteFile(filePath, []byte(template), 0644); err != nil {
				return fmt.Errorf("failed to save automigrate file: %w", err)
			}

			p.updateSingleCachedCollection(new, old)

			return nil
		})
	}
}

func (p *plugin) updateSingleCachedCollection(new, old *models.Collection) {
	cached, _ := p.app.Cache().Get(collectionsCacheKey).(map[string]*models.Collection)

	switch {
	case new == nil:
		delete(cached, old.Id)
	default:
		cached[new.Id] = new
	}

	p.app.Cache().Set(collectionsCacheKey, cached)
}

func (p *plugin) refreshCachedCollections() error {
	if p.app.Dao() == nil {
		return errors.New("app is not initialized yet")
	}

	var collections []*models.Collection
	if err := p.app.Dao().CollectionQuery().All(&collections); err != nil {
		return err
	}

	cached := map[string]*models.Collection{}
	for _, c := range collections {
		cached[c.Id] = c
	}

	p.app.Cache().Set(collectionsCacheKey, cached)

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

func (p *plugin) hasCustomMigrations() bool {
	files, err := os.ReadDir(p.options.Dir)
	if err != nil {
		return false
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		return true
	}

	return false
}
