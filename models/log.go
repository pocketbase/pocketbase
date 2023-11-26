package models

import (
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ Model = (*Log)(nil)

type Log struct {
	BaseModel

	Data    types.JsonMap `db:"data" json:"data"`
	Message string        `db:"message" json:"message"`
	Level   int           `db:"level" json:"level"`
}

func (m *Log) TableName() string {
	return "_logs"
}
