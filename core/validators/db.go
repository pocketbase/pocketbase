package validators

import (
	"database/sql"
	"errors"
	"regexp"
	"slices"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pocketbase/dbx"
)

// UniqueId checks whether a field string id already exists in the specified table.
//
// Example:
//
//	validation.Field(&form.RelId, validation.By(validators.UniqueId(form.app.DB(), "tbl_example"))
func UniqueId(db dbx.Builder, tableName string) validation.RuleFunc {
	return func(value any) error {
		v, _ := value.(string)
		if v == "" {
			return nil // nothing to check
		}

		var foundId string

		err := db.
			Select("id").
			From(tableName).
			Where(dbx.HashExp{"id": v}).
			Limit(1).
			Row(&foundId)

		if (err != nil && !errors.Is(err, sql.ErrNoRows)) || foundId != "" {
			return validation.NewError("validation_invalid_or_existing_id", "The model id is invalid or already exists.")
		}

		return nil
	}
}

var regexViolateUniqueConstraint = regexp.MustCompile(`^Key \((.+)\)=.+ already exists\.$`)

// Input: Key (col1, "col 2", col3)=("value1", "value2", "value3") already exists.
// Output: col1, "col 2", col3
func extractColumnsFromPgErrorUniqueViolation(errorDetail string) []string {
	matches := regexViolateUniqueConstraint.FindStringSubmatch(errorDetail)
	if len(matches) < 2 {
		return nil // no match
	}

	var res []string
	columns := strings.Split(matches[1], ", ") // eg: col1, "col 2", col3
	for _, col := range columns {
		col = strings.Trim(col, `"`) // remove double quotes if any
		res = append(res, col)
	}
	return res
}

// NormalizeUniqueIndexError attempts to convert a
// "unique constraint failed" error into a validation.Errors.
//
// The provided err is returned as it is without changes if:
// - err is nil
// - err is already validation.Errors
// - err is not "unique constraint failed" error
func NormalizeUniqueIndexError(err error, tableOrAlias string, fieldNames []string) error {
	if err == nil {
		return err
	}

	if _, ok := err.(validation.Errors); ok {
		return err
	}

	// check for unique constraint failure
	/* SQLite:
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "unique constraint failed") {
		// note: extra space to unify multi-columns lookup
		msg = strings.ReplaceAll(strings.TrimSpace(msg), ",", " ") + " "

		normalizedErrs := validation.Errors{}

		for _, name := range fieldNames {
			// note: extra spaces to exclude table name with suffix matching the current one
			// 		 OR other fields starting with the current field name
			if strings.Contains(msg, strings.ToLower(" "+tableOrAlias+"."+name+" ")) {
				normalizedErrs[name] = validation.NewError("validation_not_unique", "Value must be unique")
			}
		}

		if len(normalizedErrs) > 0 {
			return normalizedErrs
		}
	}
	*/

	// check for unique constraint failure
	// PostgreSQL:
	if pgErr, ok := err.(*pgconn.PgError); ok {
		// The PostgreSQL 23505 UNIQUE VIOLATION error occurs when a unique constraint is violated.
		if pgErr.Code == "23505" {
			normalizedErrs := validation.Errors{}

			// extract the field name from pgError.Detail
			columns := extractColumnsFromPgErrorUniqueViolation(pgErr.Detail)

			for _, name := range fieldNames {
				if slices.Contains(columns, name) {
					normalizedErrs[name] = validation.NewError("validation_not_unique", "Value must be unique")
				}
			}

			if len(normalizedErrs) > 0 {
				return normalizedErrs
			}
		}
	}

	return err
}
