package core

import (
	"fmt"

	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/router"
)

const CollectionNameSuperusers = "_superusers"

func (app *BaseApp) registerSuperuserHooks() {
	app.OnRecordDelete(CollectionNameSuperusers).Bind(&hook.Handler[*RecordEvent]{
		Id: "pbSuperusersRecordDelete",
		Func: func(e *RecordEvent) error {
			originalApp := e.App
			txErr := e.App.RunInTransaction(func(txApp App) error {
				e.App = txApp

				total, err := e.App.CountRecords(CollectionNameSuperusers)
				if err != nil {
					return fmt.Errorf("failed to fetch total superusers count: %w", err)
				}

				if total == 1 {
					return router.NewBadRequestError("You can't delete the only existing superuser", nil)
				}

				return e.Next()
			})
			e.App = originalApp

			return txErr
		},
		Priority: -99,
	})

	recordSaveHandler := &hook.Handler[*RecordEvent]{
		Id: "pbSuperusersRecordSaveExec",
		Func: func(e *RecordEvent) error {
			e.Record.SetVerified(true) // always mark superusers as verified
			return e.Next()
		},
		Priority: -99,
	}
	app.OnRecordCreateExecute(CollectionNameSuperusers).Bind(recordSaveHandler)
	app.OnRecordUpdateExecute(CollectionNameSuperusers).Bind(recordSaveHandler)

	collectionSaveHandler := &hook.Handler[*CollectionEvent]{
		Id: "pbSuperusersCollectionSaveExec",
		Func: func(e *CollectionEvent) error {
			// don't allow name change even if executed with SaveNoValidate
			e.Collection.Name = CollectionNameSuperusers

			// for now don't allow superusers OAuth2 since we don't want
			// to accidentally create a new superuser by just OAuth2 signin
			e.Collection.OAuth2.Enabled = false
			e.Collection.OAuth2.Providers = nil

			// force password auth
			e.Collection.PasswordAuth.Enabled = true

			// for superusers we don't allow for now standalone OTP auth and always require to be combined with MFA
			if e.Collection.OTP.Enabled {
				e.Collection.MFA.Enabled = true
			}

			return e.Next()
		},
		Priority: 99,
	}
	app.OnCollectionCreateExecute(CollectionNameSuperusers).Bind(collectionSaveHandler)
	app.OnCollectionUpdateExecute(CollectionNameSuperusers).Bind(collectionSaveHandler)
}

// IsSuperuser returns whether the current record is a superuser, aka.
// whether the record is from the _superusers collection.
func (m *Record) IsSuperuser() bool {
	return m.Collection().Name == CollectionNameSuperusers
}
