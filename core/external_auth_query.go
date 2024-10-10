package core

import (
	"github.com/pocketbase/dbx"
)

// FindAllExternalAuthsByRecord returns all ExternalAuth models
// linked to the provided auth record.
func (app *BaseApp) FindAllExternalAuthsByRecord(authRecord *Record) ([]*ExternalAuth, error) {
	auths := []*ExternalAuth{}

	err := app.RecordQuery(CollectionNameExternalAuths).
		AndWhere(dbx.HashExp{
			"collectionRef": authRecord.Collection().Id,
			"recordRef":     authRecord.Id,
		}).
		OrderBy("created DESC").
		All(&auths)

	if err != nil {
		return nil, err
	}

	return auths, nil
}

// FindAllExternalAuthsByCollection returns all ExternalAuth models
// linked to the provided auth collection.
func (app *BaseApp) FindAllExternalAuthsByCollection(collection *Collection) ([]*ExternalAuth, error) {
	auths := []*ExternalAuth{}

	err := app.RecordQuery(CollectionNameExternalAuths).
		AndWhere(dbx.HashExp{"collectionRef": collection.Id}).
		OrderBy("created DESC").
		All(&auths)

	if err != nil {
		return nil, err
	}

	return auths, nil
}

// FindFirstExternalAuthByExpr returns the first available (the most recent created)
// ExternalAuth model that satisfies the non-nil expression.
func (app *BaseApp) FindFirstExternalAuthByExpr(expr dbx.Expression) (*ExternalAuth, error) {
	model := &ExternalAuth{}

	err := app.RecordQuery(CollectionNameExternalAuths).
		AndWhere(dbx.Not(dbx.HashExp{"providerId": ""})). // exclude empty providerIds
		AndWhere(expr).
		OrderBy("created DESC").
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}
