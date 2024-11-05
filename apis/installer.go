package apis

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/security"
)

const installerEmail = "__pbinstaller@example.com"
const installerHookId = "__pbinstallerHook"

func loadInstaller(app core.App, hostURL string) error {
	if !needInstallerSuperuser(app) {
		return nil
	}

	installerRecord, err := findOrCreateInstallerSuperuser(app)
	if err != nil {
		return err
	}

	token, err := installerRecord.NewStaticAuthToken(30 * time.Minute)
	if err != nil {
		return err
	}

	// prevent sending password reset emails to the installer address
	app.OnMailerRecordPasswordResetSend(core.CollectionNameSuperusers).Bind(&hook.Handler[*core.MailerRecordEvent]{
		Id: installerHookId,
		Func: func(e *core.MailerRecordEvent) error {
			if e.Record.Email() == installerEmail {
				return errors.New("cannot reset the password for the installer account")
			}

			return e.Next()
		},
	})

	// cleanup the installer account after the first superuser creation
	app.OnRecordCreate(core.CollectionNameSuperusers).Bind(&hook.Handler[*core.RecordEvent]{
		Id: installerHookId,
		Func: func(e *core.RecordEvent) error {
			if err := e.Next(); err != nil {
				return err
			}

			color.Green("Successfully created superuser %s! This message will no longer show on the next startup.\n\n", e.Record.Email())

			if err = e.App.Delete(installerRecord); err != nil {
				e.App.Logger().Error("Failed to remove installer superuser", "error", err)
			}

			app.OnRecordCreate().Unbind(installerHookId)
			app.OnMailerRecordPasswordResetSend().Unbind(installerHookId)

			return nil
		},
	})

	// launch url (ignore errors and always print a help text as fallback)
	url := fmt.Sprintf("%s/_/#/pbinstal/%s", hostURL, token)
	_ = launchURL(url)
	color.Magenta("\n(!) Launch the URL below in the browser if it hasn't been open already to create your first superuser account:")
	color.New(color.Bold).Add(color.FgCyan).Println(url)
	color.New(color.FgHiBlack, color.Italic).Printf("(you can also create your first superuser account by running '%s superuser upsert test@example.com yourpass' and restart the server)\n", os.Args[0])

	return nil
}

func needInstallerSuperuser(app core.App) bool {
	total, err := app.CountRecords(core.CollectionNameSuperusers, dbx.Not(dbx.HashExp{
		"email": installerEmail,
	}))
	return err == nil && total == 0
}

func findOrCreateInstallerSuperuser(app core.App) (*core.Record, error) {
	col, err := app.FindCachedCollectionByNameOrId(core.CollectionNameSuperusers)
	if err != nil {
		return nil, err
	}

	record, err := app.FindAuthRecordByEmail(col, installerEmail)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		record = core.NewRecord(col)
		record.SetEmail(installerEmail)
		record.SetPassword(security.RandomString(30))

		err = app.Save(record)
		if err != nil {
			return nil, err
		}
	}

	return record, nil
}

func launchURL(url string) error {
	if err := is.URL.Validate(url); err != nil {
		return err
	}

	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		// not sure if this is the best command but seems to be the most reliable based on the comments in
		// https://stackoverflow.com/questions/3739327/launching-a-website-via-the-windows-commandline#answer-49115945
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default: // linux, freebsd, etc.
		return exec.Command("xdg-open", url).Start()
	}
}
