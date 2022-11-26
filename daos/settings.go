package daos

import (
	"encoding/json"
	"errors"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/settings"
	"github.com/pocketbase/pocketbase/tools/security"
)

// FindSettings returns and decode the serialized app settings param value.
//
// The method will first try to decode the param value without decryption.
// If it fails and optEncryptionKey is set, it will try again by first
// decrypting the value and then decode it again.
//
// Returns an error if it fails to decode the stored serialized param value.
func (dao *Dao) FindSettings(optEncryptionKey ...string) (*settings.Settings, error) {
	param, err := dao.FindParamByKey(models.ParamAppSettings)
	if err != nil {
		return nil, err
	}

	result := settings.New()

	// try first without decryption
	plainDecodeErr := json.Unmarshal(param.Value, result)

	// failed, try to decrypt
	if plainDecodeErr != nil {
		var encryptionKey string
		if len(optEncryptionKey) > 0 && optEncryptionKey[0] != "" {
			encryptionKey = optEncryptionKey[0]
		}

		// load without decrypt has failed and there is no encryption key to use for decrypt
		if encryptionKey == "" {
			return nil, errors.New("failed to load the stored app settings - missing or invalid encryption key")
		}

		// decrypt
		decrypted, decryptErr := security.Decrypt(string(param.Value), encryptionKey)
		if decryptErr != nil {
			return nil, decryptErr
		}

		// decode again
		decryptedDecodeErr := json.Unmarshal(decrypted, result)
		if decryptedDecodeErr != nil {
			return nil, decryptedDecodeErr
		}
	}

	return result, nil
}

// SaveSettings persists the specified settings configuration.
//
// If optEncryptionKey is set, then the stored serialized value will be encrypted with it.
func (dao *Dao) SaveSettings(newSettings *settings.Settings, optEncryptionKey ...string) error {
	return dao.SaveParam(models.ParamAppSettings, newSettings, optEncryptionKey...)
}
