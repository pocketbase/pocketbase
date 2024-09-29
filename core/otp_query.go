package core

import (
	"errors"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/types"
)

// FindAllOTPsByRecord returns all OTP models linked to the provided auth record.
func (app *BaseApp) FindAllOTPsByRecord(authRecord *Record) ([]*OTP, error) {
	result := []*OTP{}

	err := app.RecordQuery(CollectionNameOTPs).
		AndWhere(dbx.HashExp{
			"collectionRef": authRecord.Collection().Id,
			"recordRef":     authRecord.Id,
		}).
		OrderBy("created DESC").
		All(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindAllOTPsByCollection returns all OTP models linked to the provided collection.
func (app *BaseApp) FindAllOTPsByCollection(collection *Collection) ([]*OTP, error) {
	result := []*OTP{}

	err := app.RecordQuery(CollectionNameOTPs).
		AndWhere(dbx.HashExp{"collectionRef": collection.Id}).
		OrderBy("created DESC").
		All(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindOTPById returns a single OTP model by its id.
func (app *BaseApp) FindOTPById(id string) (*OTP, error) {
	result := &OTP{}

	err := app.RecordQuery(CollectionNameOTPs).
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteAllOTPsByRecord deletes all OTP models associated with the provided record.
//
// Returns a combined error with the failed deletes.
func (app *BaseApp) DeleteAllOTPsByRecord(authRecord *Record) error {
	models, err := app.FindAllOTPsByRecord(authRecord)
	if err != nil {
		return err
	}

	var errs []error
	for _, m := range models {
		if err := app.Delete(m); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// DeleteExpiredOTPs deletes the expired OTPs for all auth collections.
func (app *BaseApp) DeleteExpiredOTPs() error {
	authCollections, err := app.FindAllCollections(CollectionTypeAuth)
	if err != nil {
		return err
	}

	// note: perform even if OTP is disabled to ensure that there are no dangling old records
	for _, collection := range authCollections {
		minValidDate, err := types.ParseDateTime(time.Now().Add(-1 * collection.OTP.DurationTime()))
		if err != nil {
			return err
		}

		items := []*Record{}

		err = app.RecordQuery(CollectionNameOTPs).
			AndWhere(dbx.HashExp{"collectionRef": collection.Id}).
			AndWhere(dbx.NewExp("[[created]] < {:date}", dbx.Params{"date": minValidDate})).
			All(&items)
		if err != nil {
			return err
		}

		for _, item := range items {
			err = app.Delete(item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
