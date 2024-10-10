package apis

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func backupUpload(e *core.RequestEvent) error {
	fsys, err := e.App.NewBackupsFilesystem()
	if err != nil {
		return err
	}
	defer fsys.Close()

	form := new(backupUploadForm)
	form.fsys = fsys
	files, _ := e.FindUploadedFiles("file")
	if len(files) > 0 {
		form.File = files[0]
	}

	err = form.validate()
	if err != nil {
		return e.BadRequestError("An error occurred while validating the submitted data.", err)
	}

	err = fsys.UploadFile(form.File, form.File.OriginalName)
	if err != nil {
		return e.BadRequestError("Failed to upload backup.", err)
	}

	// we don't retrieve the generated backup file because it may not be
	// available yet due to the eventually consistent nature of some S3 providers
	return e.NoContent(http.StatusNoContent)
}

// -------------------------------------------------------------------

type backupUploadForm struct {
	fsys *filesystem.System

	File *filesystem.File `json:"file"`
}

func (form *backupUploadForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.File,
			validation.Required,
			validation.By(validators.UploadedFileMimeType([]string{"application/zip"})),
			validation.By(form.checkUniqueName),
		),
	)
}

func (form *backupUploadForm) checkUniqueName(value any) error {
	v, _ := value.(*filesystem.File)
	if v == nil {
		return nil // nothing to check
	}

	// note: we use the original name because that is what we upload
	if exists, err := form.fsys.Exists(v.OriginalName); err != nil || exists {
		return validation.NewError("validation_backup_name_exists", "Backup file with the specified name already exists.")
	}

	return nil
}
