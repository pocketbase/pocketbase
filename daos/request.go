package daos

import (
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

// RequestQuery returns a new Request logs select query.
func (dao *Dao) RequestQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&models.Request{})
}

// FindRequestById finds a single Request log by its id.
func (dao *Dao) FindRequestById(id string) (*models.Request, error) {
	model := &models.Request{}

	err := dao.RequestQuery().
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(model)

	if err != nil {
		return nil, err
	}

	return model, nil
}

type RequestsStatsItem struct {
	Total int            `db:"total" json:"total"`
	Date  types.DateTime `db:"date" json:"date"`
}

// RequestsStats returns hourly grouped requests logs statistics.
func (dao *Dao) RequestsStats(expr dbx.Expression) ([]*RequestsStatsItem, error) {
	result := []*RequestsStatsItem{}

	query := dao.RequestQuery().
		Select("count(id) as total", "strftime('%Y-%m-%d %H:00:00', created) as date").
		GroupBy("date")

	if expr != nil {
		query.AndWhere(expr)
	}

	err := query.All(&result)

	return result, err
}

// DeleteOldRequests delete all requests that are created before createdBefore.
func (dao *Dao) DeleteOldRequests(createdBefore time.Time) error {
	m := models.Request{}
	tableName := m.TableName()

	formattedDate := createdBefore.UTC().Format(types.DefaultDateLayout)
	expr := dbx.NewExp("[[created]] <= {:date}", dbx.Params{"date": formattedDate})

	_, err := dao.NonconcurrentDB().Delete(tableName, expr).Execute()

	return err
}

// SaveRequest upserts the provided Request model.
func (dao *Dao) SaveRequest(request *models.Request) error {
	return dao.Save(request)
}
