package forms

import (
	"context"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

var backupNameRegex = regexp.MustCompile(`^[a-z0-9_-]+\.zip$`)

// BackupCreate is a request form for creating a new app backup.
type BackupCreate struct {
	app core.App
	ctx context.Context

	Name string `form:"name" json:"name"`
}

// NewBackupCreate creates new BackupCreate request form.
func NewBackupCreate(app core.App) *BackupCreate {
	return &BackupCreate{
		app: app,
		ctx: context.Background(),
	}
}

// SetContext replaces the default form context with the provided one.
func (form *BackupCreate) SetContext(ctx context.Context) {
	form.ctx = ctx
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *BackupCreate) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Name,
			validation.Length(1, 100),
			validation.Match(backupNameRegex),
			validation.By(form.checkUniqueName),
		),
	)
}

func (form *BackupCreate) checkUniqueName(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	fsys, err := form.app.NewBackupsFilesystem()
	if err != nil {
		return err
	}
	defer fsys.Close()

	fsys.SetContext(form.ctx)

	if exists, err := fsys.Exists(v); err != nil || exists {
		return validation.NewError("validation_backup_name_exists", "The backup file name is invalid or already exists.")
	}

	return nil
}

// Submit validates the form and creates the app backup.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before creating the backup.
func (form *BackupCreate) Submit(interceptors ...InterceptorFunc[string]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	return runInterceptors(form.Name, func(name string) error {
		return form.app.CreateBackup(form.ctx, name)
	}, interceptors...)
}
