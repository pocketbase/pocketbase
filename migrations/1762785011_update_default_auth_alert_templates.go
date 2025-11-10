package migrations

import (
	"github.com/pocketbase/pocketbase/core"
)

const oldAuthAlertTemplate = `<p>Hello,</p>
<p>We noticed a login to your ` + core.EmailPlaceholderAppName + ` account from a new location.</p>
<p>If this was you, you may disregard this email.</p>
<p><strong>If this wasn't you, you should immediately change your ` + core.EmailPlaceholderAppName + ` account password to revoke access from all other locations.</strong></p>
<p>
  Thanks,<br/>
  ` + core.EmailPlaceholderAppName + ` team
</p>`

func init() {
	core.SystemMigrations.Register(func(txApp core.App) error {
		collections, err := txApp.FindAllCollections(core.CollectionTypeAuth)
		if err != nil {
			return err
		}

		newTemplate := core.NewAuthCollection("up").AuthAlert.EmailTemplate.Body

		for _, c := range collections {
			if c.AuthAlert.EmailTemplate.Body != oldAuthAlertTemplate {
				continue
			}

			c.AuthAlert.EmailTemplate.Body = newTemplate

			err = txApp.Save(c)
			if err != nil {
				return err
			}
		}

		return nil
	}, func(txApp core.App) error {
		collections, err := txApp.FindAllCollections(core.CollectionTypeAuth)
		if err != nil {
			return err
		}

		newTemplate := core.NewAuthCollection("down").AuthAlert.EmailTemplate.Body

		for _, c := range collections {
			if c.AuthAlert.EmailTemplate.Body != newTemplate {
				continue
			}

			c.AuthAlert.EmailTemplate.Body = oldAuthAlertTemplate

			err = txApp.Save(c)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
