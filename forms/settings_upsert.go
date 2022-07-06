package forms

import (
	"os"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

// SettingsUpsert defines app settings upsert form.
type SettingsUpsert struct {
	*core.Settings

	app core.App
}

// NewSettingsUpsert creates new settings upsert form from the provided app.
func NewSettingsUpsert(app core.App) *SettingsUpsert {
	form := &SettingsUpsert{app: app}

	// load the application settings into the form
	form.Settings, _ = app.Settings().Clone()

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *SettingsUpsert) Validate() error {
	return form.Settings.Validate()
}

// Submit validates the form and upserts the loaded settings.
//
// On success the app settings will be refreshed with the form ones.
func (form *SettingsUpsert) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	encryptionKey := os.Getenv(form.app.EncryptionEnv())

	saveErr := form.app.Dao().SaveParam(
		models.ParamAppSettings,
		form.Settings,
		encryptionKey,
	)
	if saveErr != nil {
		return saveErr
	}

	// explicitly trigger old logs deletion
	form.app.LogsDao().DeleteOldRequests(
		time.Now().AddDate(0, 0, -1*form.Settings.Logs.MaxDays),
	)

	// merge the application settings with the form ones
	return form.app.Settings().Merge(form.Settings)
}
