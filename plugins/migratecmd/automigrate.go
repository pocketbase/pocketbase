package migratecmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// automigrateOnCollectionChange handles the automigration snapshot
// generation on collection change request event (create/update/delete).
func (p *plugin) automigrateOnCollectionChange(e *core.CollectionRequestEvent) error {
	var err error
	var old *core.Collection
	if !e.Collection.IsNew() {
		old, err = e.App.FindCollectionByNameOrId(e.Collection.Id)
		if err != nil {
			return err
		}
	}

	err = e.Next()
	if err != nil {
		return err
	}

	new, err := p.app.FindCollectionByNameOrId(e.Collection.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// for now exclude OAuth2 configs from the migration
	if old != nil && old.IsAuth() {
		old.OAuth2.Providers = nil
	}
	if new != nil && new.IsAuth() {
		new.OAuth2.Providers = nil
	}

	var template string
	var templateErr error
	if p.config.TemplateLang == TemplateLangJS {
		template, templateErr = p.jsDiffTemplate(new, old)
	} else {
		template, templateErr = p.goDiffTemplate(new, old)
	}
	if templateErr != nil {
		if errors.Is(templateErr, ErrEmptyTemplate) {
			return nil // no changes
		}
		return fmt.Errorf("failed to resolve template: %w", templateErr)
	}

	var action string
	switch {
	case new == nil:
		action = "deleted_" + normalizeCollectionName(old.Name)
	case old == nil:
		action = "created_" + normalizeCollectionName(new.Name)
	default:
		action = "updated_" + normalizeCollectionName(old.Name)
	}

	name := fmt.Sprintf("%d_%s.%s", time.Now().Unix(), action, p.config.TemplateLang)
	filePath := filepath.Join(p.config.Dir, name)

	return p.app.RunInTransaction(func(txApp core.App) error {
		// insert the migration entry
		_, err := txApp.DB().Insert(core.DefaultMigrationsTable, dbx.Params{
			"file": name,
			// use microseconds for more granular applied time in case
			// multiple collection changes happens at the ~exact time
			"applied": time.Now().UnixMicro(),
		}).Execute()
		if err != nil {
			return err
		}

		// ensure that the local migrations dir exist
		if err := os.MkdirAll(p.config.Dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create migration dir: %w", err)
		}

		if err := os.WriteFile(filePath, []byte(template), 0644); err != nil {
			return fmt.Errorf("failed to save automigrate file: %w", err)
		}

		return nil
	})
}

func normalizeCollectionName(name string) string {
	// adds an extra "_" suffix to the name in case the collection ends
	// with "test" to prevent accidentally resulting in "_test.go"/"_test.js" files
	if strings.HasSuffix(strings.ToLower(name), "test") {
		name += "_"
	}

	return name
}
