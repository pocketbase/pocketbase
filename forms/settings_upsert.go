package forms

import (
	"os"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// SettingsUpsert specifies a [core.Settings] upsert (create/update) form.
type SettingsUpsert struct {
	*core.Settings

	config SettingsUpsertConfig
}

// SettingsUpsertConfig is the [SettingsUpsert] factory initializer config.
//
// NB! App is required struct member.
type SettingsUpsertConfig struct {
	App     core.App
	Dao     *daos.Dao
	LogsDao *daos.Dao
}

// NewSettingsUpsert creates a new [SettingsUpsert] form with initializer
// config created from the provided [core.App] instance.
//
// If you want to submit the form as part of another transaction, use
// [NewSettingsUpsertWithConfig] with explicitly set Dao.
func NewSettingsUpsert(app core.App) *SettingsUpsert {
	return NewSettingsUpsertWithConfig(SettingsUpsertConfig{
		App: app,
	})
}

// NewSettingsUpsertWithConfig creates a new [SettingsUpsert] form
// with the provided config or panics on invalid configuration.
func NewSettingsUpsertWithConfig(config SettingsUpsertConfig) *SettingsUpsert {
	form := &SettingsUpsert{config: config}

	if form.config.App == nil {
		panic("Missing required config.App instance.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	if form.config.LogsDao == nil {
		form.config.LogsDao = form.config.App.LogsDao()
	}

	// load the application settings into the form
	form.Settings, _ = config.App.Settings().Clone()

	return form
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

	encryptionKey := os.Getenv(form.config.App.EncryptionEnv())

	return runInterceptors(func() error {
		saveErr := form.config.Dao.SaveParam(
			models.ParamAppSettings,
			form.Settings,
			encryptionKey,
		)
		if saveErr != nil {
			return saveErr
		}

		// explicitly trigger old logs deletion
		form.config.LogsDao.DeleteOldRequests(
			time.Now().AddDate(0, 0, -1*form.Settings.Logs.MaxDays),
		)

		// merge the application settings with the form ones
		return form.config.App.Settings().Merge(form.Settings)
	}, interceptors...)
}
