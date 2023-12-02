package daos

import (
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
)

// ExternalAuthQuery returns a new ExternalAuth select query.
func (dao *Dao) ExternalAuthQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&models.ExternalAuth{})
}

// FindAllExternalAuthsByRecord returns all ExternalAuth models
// linked to the provided auth record.
func (dao *Dao) FindAllExternalAuthsByRecord(authRecord *models.Record) ([]*models.ExternalAuth, error) {
	auths := []*models.ExternalAuth{}

	err := dao.ExternalAuthQuery().
		AndWhere(dbx.HashExp{
			"collectionId": authRecord.Collection().Id,
			"recordId":     authRecord.Id,
		}).
		OrderBy("created ASC").
		All(&auths)

	if err != nil {
		return nil, err
	}

	return auths, nil
}

// FindExternalAuthByRecordAndProvider returns the first available
// ExternalAuth model for the specified record data and provider.
func (dao *Dao) FindExternalAuthByRecordAndProvider(authRecord *models.Record, provider string) (*models.ExternalAuth, error) {
	model := &models.ExternalAuth{}

	err := dao.ExternalAuthQuery().
		AndWhere(dbx.HashExp{
			"collectionId": authRecord.Collection().Id,
			"recordId":     authRecord.Id,
			"provider":     provider,
		}).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

// FindFirstExternalAuthByExpr returns the first available
// ExternalAuth model that satisfies the non-nil expression.
func (dao *Dao) FindFirstExternalAuthByExpr(expr dbx.Expression) (*models.ExternalAuth, error) {
	model := &models.ExternalAuth{}

	err := dao.ExternalAuthQuery().
		AndWhere(dbx.Not(dbx.HashExp{"providerId": ""})). // exclude empty providerIds
		AndWhere(expr).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

// SaveExternalAuth upserts the provided ExternalAuth model.
func (dao *Dao) SaveExternalAuth(model *models.ExternalAuth) error {
	// extra check the model data in case the provider's API response
	// has changed and no longer returns the expected fields
	if model.CollectionId == "" || model.RecordId == "" || model.Provider == "" || model.ProviderId == "" {
		return errors.New("Missing required ExternalAuth fields.")
	}

	return dao.Save(model)
}

// DeleteExternalAuth deletes the provided ExternalAuth model.
func (dao *Dao) DeleteExternalAuth(model *models.ExternalAuth) error {
	return dao.Delete(model)
}
