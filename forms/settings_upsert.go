package forms

import (
	"os"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// SettingsUpsert is a [core.Settings] upsert (create/update) form.
type SettingsUpsert struct {
	*core.Settings

	app core.App
	dao *daos.Dao
}

// NewSettingsUpsert creates a new [SettingsUpsert] form with initializer
// config created from the provided [core.App] instance.
//
// If you want to submit the form as part of a transaction,
// you can change the default Dao via [SetDao()].
func NewSettingsUpsert(app core.App) *SettingsUpsert {
	form := &SettingsUpsert{
		app: app,
		dao: app.Dao(),
	}

	// load the application settings into the form
	form.Settings, _ = app.Settings().Clone()

	return form
}

// SetDao replaces the default form Dao instance with the provided one.
func (form *SettingsUpsert) SetDao(dao *daos.Dao) {
	form.dao = dao
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *SettingsUpsert) Validate() error {
	return form.Settings.Validate()
}

// Submit validates the form and upserts the loaded settings.
//
// On success the app settings will be refreshed with the form ones.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *SettingsUpsert) Submit(interceptors ...InterceptorFunc) error {
	if err := form.Validate(); err != nil {
		return err
	}

	encryptionKey := os.Getenv(form.app.EncryptionEnv())

	return runInterceptors(func() error {
		saveErr := form.dao.SaveParam(
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

		if form.Settings.Logs.MaxDays == 0 {
			// reclaim deleted logs disk space
			form.app.LogsDao().Vacuum()
		}

		// merge the application settings with the form ones
		return form.app.Settings().Merge(form.Settings)
	}, interceptors...)
}
