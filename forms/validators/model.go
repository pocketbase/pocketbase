package validators

import (
	"database/sql"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

// UniqueId checks whether the provided model id already exists.
//
// Example:
//
//	validation.Field(&form.Id, validation.By(validators.UniqueId(form.dao, tableName)))
func UniqueId(dao *daos.Dao, tableName string) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(string)
		if v == "" {
			return nil // nothing to check
		}

		var foundId string

		err := dao.DB().
			Select("id").
			From(tableName).
			Where(dbx.HashExp{"id": v}).
			Limit(1).
			Row(&foundId)

		if (err != nil && !errors.Is(err, sql.ErrNoRows)) || foundId != "" {
			return validation.NewError("validation_invalid_id", "The model id is invalid or already exists.")
		}

		return nil
	}
}
