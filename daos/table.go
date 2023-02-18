package daos

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
)

// HasTable checks if a table with the provided name exists (case insensitive).
func (dao *Dao) HasTable(tableName string) bool {
	var exists bool

	err := dao.DB().Select("count(*)").
		From("sqlite_schema").
		AndWhere(dbx.HashExp{"type": "table"}).
		AndWhere(dbx.NewExp("LOWER([[name]])=LOWER({:tableName})", dbx.Params{"tableName": tableName})).
		Limit(1).
		Row(&exists)

	return err == nil && exists
}

// GetTableColumns returns all column names of a single table by its name.
func (dao *Dao) GetTableColumns(tableName string) ([]string, error) {
	columns := []string{}

	err := dao.DB().NewQuery("SELECT name FROM PRAGMA_TABLE_INFO({:tableName})").
		Bind(dbx.Params{"tableName": tableName}).
		Column(&columns)

	return columns, err
}

// GetTableInfo returns the `table_info` pragma result for the specified table.
func (dao *Dao) GetTableInfo(tableName string) ([]*models.TableInfoRow, error) {
	info := []*models.TableInfoRow{}

	err := dao.DB().NewQuery("SELECT * FROM PRAGMA_TABLE_INFO({:tableName})").
		Bind(dbx.Params{"tableName": tableName}).
		All(&info)

	return info, err
}

// DeleteTable drops the specified table.
//
// This method is a no-op if a table with the provided name doesn't exist.
//
// Be aware that this method is vulnerable to SQL injection and the
// "tableName" argument must come only from trusted input!
func (dao *Dao) DeleteTable(tableName string) error {
	_, err := dao.DB().NewQuery(fmt.Sprintf(
		"DROP TABLE IF EXISTS {{%s}}",
		tableName,
	)).Execute()

	return err
}

// Vacuum executes VACUUM on the current dao.DB() instance in order to
// reclaim unused db disk space.
func (dao *Dao) Vacuum() error {
	_, err := dao.DB().NewQuery("VACUUM").Execute()

	return err
}
