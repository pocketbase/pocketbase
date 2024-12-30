package apis

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/osutils"
)

// DefaultInstallerFunc is the default PocketBase installer function.
//
// It will attempt to open a link in the browser (with a short-lived auth
// token for the systemSuperuser) to the installer UI so that users can
// create their own custom superuser record.
//
// See https://github.com/pocketbase/pocketbase/discussions/5814.
func DefaultInstallerFunc(app core.App, systemSuperuser *core.Record, baseURL string) error {
	token, err := systemSuperuser.NewStaticAuthToken(30 * time.Minute)
	if err != nil {
		return err
	}

	// launch url (ignore errors and always print a help text as fallback)
	url := fmt.Sprintf("%s/_/#/pbinstal/%s", strings.TrimRight(baseURL, "/"), token)
	_ = osutils.LaunchURL(url)
	color.Magenta("\n(!) Launch the URL below in the browser if it hasn't been open already to create your first superuser account:")
	color.New(color.Bold).Add(color.FgCyan).Println(url)
	color.New(color.FgHiBlack, color.Italic).Printf("(you can also create your first superuser by running: %s superuser upsert EMAIL PASS)\n\n", os.Args[0])

	return nil
}

func loadInstaller(
	app core.App,
	baseURL string,
	installerFunc func(app core.App, systemSuperuser *core.Record, baseURL string) error,
) error {
	if installerFunc == nil || !needInstallerSuperuser(app) {
		return nil
	}

	superuser, err := findOrCreateInstallerSuperuser(app)
	if err != nil {
		return err
	}

	return installerFunc(app, superuser, baseURL)
}

func needInstallerSuperuser(app core.App) bool {
	total, err := app.CountRecords(core.CollectionNameSuperusers, dbx.Not(dbx.HashExp{
		"email": core.DefaultInstallerEmail,
	}))

	return err == nil && total == 0
}

func findOrCreateInstallerSuperuser(app core.App) (*core.Record, error) {
	col, err := app.FindCachedCollectionByNameOrId(core.CollectionNameSuperusers)
	if err != nil {
		return nil, err
	}

	record, err := app.FindAuthRecordByEmail(col, core.DefaultInstallerEmail)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		record = core.NewRecord(col)
		record.SetEmail(core.DefaultInstallerEmail)
		record.SetRandomPassword()

		err = app.Save(record)
		if err != nil {
			return nil, err
		}
	}

	return record, nil
}
