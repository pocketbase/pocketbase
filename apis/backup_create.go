package apis

import (
	"context"
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

func backupCreate(e *core.RequestEvent) error {
	if e.App.Store().Has(core.StoreKeyActiveBackup) {
		return e.BadRequestError("Try again later - another backup/restore process has already been started", nil)
	}

	form := new(backupCreateForm)
	form.app = e.App

	err := e.BindBody(form)
	if err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}

	err = form.validate()
	if err != nil {
		return e.BadRequestError("An error occurred while validating the submitted data.", err)
	}

	err = e.App.CreateBackup(context.Background(), form.Name)
	if err != nil {
		return e.BadRequestError("Failed to create backup.", err)
	}

	// we don't retrieve the generated backup file because it may not be
	// available yet due to the eventually consistent nature of some S3 providers
	return e.NoContent(http.StatusNoContent)
}

// -------------------------------------------------------------------

var backupNameRegex = regexp.MustCompile(`^[a-z0-9_-]+\.zip$`)

type backupCreateForm struct {
	app core.App

	Name string `form:"name" json:"name"`
}

func (form *backupCreateForm) validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Name,
			validation.Length(1, 150),
			validation.Match(backupNameRegex),
			validation.By(form.checkUniqueName),
		),
	)
}

func (form *backupCreateForm) checkUniqueName(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	fsys, err := form.app.NewBackupsFilesystem()
	if err != nil {
		return err
	}
	defer fsys.Close()

	if exists, err := fsys.Exists(v); err != nil || exists {
		return validation.NewError("validation_backup_name_exists", "The backup file name is invalid or already exists.")
	}

	return nil
}
