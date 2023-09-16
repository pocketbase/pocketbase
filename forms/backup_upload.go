package forms

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// BackupUpload is a request form for uploading a new app backup.
type BackupUpload struct {
	app core.App
	ctx context.Context

	File *filesystem.File `json:"file"`
}

// NewBackupUpload creates new BackupUpload request form.
func NewBackupUpload(app core.App) *BackupUpload {
	return &BackupUpload{
		app: app,
		ctx: context.Background(),
	}
}

// SetContext replaces the default form upload context with the provided one.
func (form *BackupUpload) SetContext(ctx context.Context) {
	form.ctx = ctx
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *BackupUpload) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.File,
			validation.Required,
			validation.By(validators.UploadedFileMimeType([]string{"application/zip"})),
			validation.By(form.checkUniqueName),
		),
	)
}

func (form *BackupUpload) checkUniqueName(value any) error {
	v, _ := value.(*filesystem.File)
	if v == nil {
		return nil // nothing to check
	}

	fsys, err := form.app.NewBackupsFilesystem()
	if err != nil {
		return err
	}
	defer fsys.Close()

	fsys.SetContext(form.ctx)

	if exists, err := fsys.Exists(v.OriginalName); err != nil || exists {
		return validation.NewError("validation_backup_name_exists", "Backup file with the specified name already exists.")
	}

	return nil
}

// Submit validates the form and upload the backup file.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before uploading the backup.
func (form *BackupUpload) Submit(interceptors ...InterceptorFunc[*filesystem.File]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	return runInterceptors(form.File, func(file *filesystem.File) error {
		fsys, err := form.app.NewBackupsFilesystem()
		if err != nil {
			return err
		}

		fsys.SetContext(form.ctx)

		return fsys.UploadFile(file, file.OriginalName)
	}, interceptors...)
}
