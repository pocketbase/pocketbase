package core

import "github.com/pocketbase/pocketbase/tools/types"

var (
	_ Model = (*Log)(nil)
)

const LogsTableName = "_logs"

type Log struct {
	BaseModel

	Created types.DateTime     `db:"created" json:"created"`
	Data    types.JSONMap[any] `db:"data" json:"data"`
	Message string             `db:"message" json:"message"`
	Level   int                `db:"level" json:"level"`
}

func (m *Log) TableName() string {
	return LogsTableName
}
