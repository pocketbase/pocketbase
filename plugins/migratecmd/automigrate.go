package migratecmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
)

const migrationsTable = "_migrations"
const automigrateSuffix = "_automigrate"

// onCollectionChange handles the automigration snapshot generation on
// collection change event (create/update/delete).
func (p *plugin) onCollectionChange() func(*core.ModelEvent) error {
	return func(e *core.ModelEvent) error {
		if e.Model.TableName() != "_collections" {
			return nil // not a collection
		}

		collections := []*models.Collection{}
		if err := p.app.Dao().CollectionQuery().OrderBy("created ASC").All(&collections); err != nil {
			return fmt.Errorf("failed to fetch collections list: %v", err)
		}
		if len(collections) == 0 {
			return errors.New("missing collections to automigrate")
		}

		oldFiles, err := p.getAllMigrationNames()
		if err != nil {
			return fmt.Errorf("failed to fetch migration files list: %v", err)
		}

		var template string
		var templateErr error
		if p.options.TemplateLang == TemplateLangJS {
			template, templateErr = p.jsSnapshotTemplate(collections)
		} else {
			template, templateErr = p.goSnapshotTemplate(collections)
		}
		if templateErr != nil {
			return fmt.Errorf("failed to resolve template: %v", templateErr)
		}

		appliedTime := time.Now().Unix()
		fileDest := filepath.Join(p.options.Dir, fmt.Sprintf("%d_automigrate.%s", appliedTime, p.options.TemplateLang))

		// ensure that the local migrations dir exist
		if err := os.MkdirAll(p.options.Dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create migration dir: %v", err)
		}

		if err := os.WriteFile(fileDest, []byte(template), 0644); err != nil {
			return fmt.Errorf("failed to save automigrate file: %v", err)
		}

		// remove the old untracked automigrate file
		// (only if the last one was automigrate!)
		if len(oldFiles) > 0 && strings.HasSuffix(oldFiles[len(oldFiles)-1], automigrateSuffix+"."+p.options.TemplateLang) {
			olfName := oldFiles[len(oldFiles)-1]
			oldPath := filepath.Join(p.options.Dir, olfName)

			isUntracked := exec.Command(p.options.GitPath, "ls-files", "--error-unmatch", oldPath).Run() != nil
			if isUntracked {
				// delete the old automigrate from the db if it was already applied
				_, err := p.app.Dao().DB().Delete(migrationsTable, dbx.HashExp{"file": olfName}).Execute()
				if err != nil {
					return fmt.Errorf("failed to delete last applied automigrate from the migration db: %v", err)
				}

				// delete the old automigrate file from the filesystem
				if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("failed to delete last automigrates from the filesystem: %v", err)
				}
			}
		}

		return nil
	}
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
