package daos

import (
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

// LogQuery returns a new Log select query.
func (dao *Dao) LogQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&models.Log{})
}

// FindLogById finds a single Log entry by its id.
func (dao *Dao) FindLogById(id string) (*models.Log, error) {
	model := &models.Log{}

	err := dao.LogQuery().
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

type LogsStatsItem struct {
	Total int            `db:"total" json:"total"`
	Date  types.DateTime `db:"date" json:"date"`
}

// LogsStats returns hourly grouped requests logs statistics.
func (dao *Dao) LogsStats(expr dbx.Expression) ([]*LogsStatsItem, error) {
	result := []*LogsStatsItem{}

	query := dao.LogQuery().
		Select("count(id) as total", "strftime('%Y-%m-%d %H:00:00', created) as date").
		GroupBy("date")

	if expr != nil {
		query.AndWhere(expr)
	}

	err := query.All(&result)

	return result, err
}

// DeleteOldLogs delete all requests that are created before createdBefore.
func (dao *Dao) DeleteOldLogs(createdBefore time.Time) error {
	formattedDate := createdBefore.UTC().Format(types.DefaultDateLayout)
	expr := dbx.NewExp("[[created]] <= {:date}", dbx.Params{"date": formattedDate})

	_, err := dao.NonconcurrentDB().Delete((&models.Log{}).TableName(), expr).Execute()

	return err
}

// SaveLog upserts the provided Log model.
func (dao *Dao) SaveLog(log *models.Log) error {
	return dao.Save(log)
}
