package daos

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
)

// ParamQuery returns a new Param select query.
func (dao *Dao) ParamQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&models.Param{})
}

// FindParamByKey finds the first Param model with the provided key.
func (dao *Dao) FindParamByKey(key string) (*models.Param, error) {
	param := &models.Param{}

	err := dao.ParamQuery().
		AndWhere(dbx.HashExp{"key": key}).
		Limit(1).
		One(param)

	if err != nil {
		return nil, err
	}

	return param, nil
}

// SaveParam creates or updates a Param model by the provided key-value pair.
// The value argument will be encoded as json string.
//
// If `optEncryptionKey` is provided it will encrypt the value before storing it.
func (dao *Dao) SaveParam(key string, value any, optEncryptionKey ...string) error {
	param, _ := dao.FindParamByKey(key)
	if param == nil {
		param = &models.Param{Key: key}
	}

	normalizedValue := value

	// encrypt if optEncryptionKey is set
	if len(optEncryptionKey) > 0 && optEncryptionKey[0] != "" {
		encoded, encodingErr := json.Marshal(value)
		if encodingErr != nil {
			return encodingErr
		}

		encryptVal, encryptErr := security.Encrypt(encoded, optEncryptionKey[0])
		if encryptErr != nil {
			return encryptErr
		}

		normalizedValue = encryptVal
	}

	encodedValue := types.JsonRaw{}
	if err := encodedValue.Scan(normalizedValue); err != nil {
		return err
	}

	param.Value = encodedValue

	return dao.Save(param)
}

// DeleteParam deletes the provided Param model.
func (dao *Dao) DeleteParam(param *models.Param) error {
	return dao.Delete(param)
}
