package forms

import (
	"os"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models/settings"
)

// SettingsUpsert is a [settings.Settings] upsert (create/update) form.
type SettingsUpsert struct {
	*settings.Settings

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
func (form *SettingsUpsert) Submit(interceptors ...InterceptorFunc[*settings.Settings]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	return runInterceptors(form.Settings, func(s *settings.Settings) error {
		form.Settings = s

		oldSettings, err := form.app.Settings().Clone();
		if err != nil {
			return err
		}

		// eagerly merge the application settings with the form ones
		if err := form.app.Settings().Merge(form.Settings); err != nil {
			return err
		}

		// persists settings change
		encryptionKey := os.Getenv(form.app.EncryptionEnv())
		if err := form.dao.SaveSettings(form.Settings, encryptionKey); err != nil {
			// try to revert app settings
			form.app.Settings().Merge(oldSettings)

			return err
		}

		// explicitly trigger old logs deletion
		form.app.LogsDao().DeleteOldRequests(
			time.Now().AddDate(0, 0, -1*form.Settings.Logs.MaxDays),
		)

		if form.Settings.Logs.MaxDays == 0 {
			// no logs are allowed -> reclaim preserved disk space after the previous delete operation
			form.app.LogsDao().Vacuum()
		}

		return nil
	}, interceptors...)
}
