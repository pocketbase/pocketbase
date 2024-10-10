package core

import (
	"errors"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/types"
)

// FindAllMFAsByRecord returns all MFA models linked to the provided auth record.
func (app *BaseApp) FindAllMFAsByRecord(authRecord *Record) ([]*MFA, error) {
	result := []*MFA{}

	err := app.RecordQuery(CollectionNameMFAs).
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

// FindAllMFAsByCollection returns all MFA models linked to the provided collection.
func (app *BaseApp) FindAllMFAsByCollection(collection *Collection) ([]*MFA, error) {
	result := []*MFA{}

	err := app.RecordQuery(CollectionNameMFAs).
		AndWhere(dbx.HashExp{"collectionRef": collection.Id}).
		OrderBy("created DESC").
		All(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindMFAById returns a single MFA model by its id.
func (app *BaseApp) FindMFAById(id string) (*MFA, error) {
	result := &MFA{}

	err := app.RecordQuery(CollectionNameMFAs).
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteAllMFAsByRecord deletes all MFA models associated with the provided record.
//
// Returns a combined error with the failed deletes.
func (app *BaseApp) DeleteAllMFAsByRecord(authRecord *Record) error {
	models, err := app.FindAllMFAsByRecord(authRecord)
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

// DeleteExpiredMFAs deletes the expired MFAs for all auth collections.
func (app *BaseApp) DeleteExpiredMFAs() error {
	authCollections, err := app.FindAllCollections(CollectionTypeAuth)
	if err != nil {
		return err
	}

	// note: perform even if MFA is disabled to ensure that there are no dangling old records
	for _, collection := range authCollections {
		minValidDate, err := types.ParseDateTime(time.Now().Add(-1 * collection.MFA.DurationTime()))
		if err != nil {
			return err
		}

		items := []*Record{}

		err = app.RecordQuery(CollectionNameMFAs).
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
