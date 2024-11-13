package migrations

import (
	"github.com/pocketbase/pocketbase/core"
)

// note: this migration will be deleted in future version

// add new OTP sentTo text field (if not already)
func init() {
	core.SystemMigrations.Register(func(txApp core.App) error {
		otpCollection, err := txApp.FindCollectionByNameOrId(core.CollectionNameOTPs)
		if err != nil {
			return err
		}

		field := otpCollection.Fields.GetByName("sentTo")
		if field != nil {
			return nil // already exists
		}

		otpCollection.Fields.Add(&core.TextField{
			Name:   "sentTo",
			System: true,
			Hidden: true,
		})

		return txApp.Save(otpCollection)
	}, nil)
}
