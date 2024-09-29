package core

import (
	"errors"

	"github.com/pocketbase/dbx"
)

// FindAllAuthOriginsByRecord returns all AuthOrigin models linked to the provided auth record (in DESC order).
func (app *BaseApp) FindAllAuthOriginsByRecord(authRecord *Record) ([]*AuthOrigin, error) {
	result := []*AuthOrigin{}

	err := app.RecordQuery(CollectionNameAuthOrigins).
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

// FindAllAuthOriginsByCollection returns all AuthOrigin models linked to the provided collection (in DESC order).
func (app *BaseApp) FindAllAuthOriginsByCollection(collection *Collection) ([]*AuthOrigin, error) {
	result := []*AuthOrigin{}

	err := app.RecordQuery(CollectionNameAuthOrigins).
		AndWhere(dbx.HashExp{"collectionRef": collection.Id}).
		OrderBy("created DESC").
		All(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindAuthOriginById returns a single AuthOrigin model by its id.
func (app *BaseApp) FindAuthOriginById(id string) (*AuthOrigin, error) {
	result := &AuthOrigin{}

	err := app.RecordQuery(CollectionNameAuthOrigins).
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindAuthOriginByRecordAndFingerprint returns a single AuthOrigin model
// by its authRecord relation and fingerprint.
func (app *BaseApp) FindAuthOriginByRecordAndFingerprint(authRecord *Record, fingerprint string) (*AuthOrigin, error) {
	result := &AuthOrigin{}

	err := app.RecordQuery(CollectionNameAuthOrigins).
		AndWhere(dbx.HashExp{
			"collectionRef": authRecord.Collection().Id,
			"recordRef":     authRecord.Id,
			"fingerprint":   fingerprint,
		}).
		Limit(1).
		One(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteAllAuthOriginsByRecord deletes all AuthOrigin models associated with the provided record.
//
// Returns a combined error with the failed deletes.
func (app *BaseApp) DeleteAllAuthOriginsByRecord(authRecord *Record) error {
	models, err := app.FindAllAuthOriginsByRecord(authRecord)
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
