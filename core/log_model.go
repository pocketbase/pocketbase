package core

import (
	"encoding/json/v2"

	"github.com/pocketbase/pocketbase/tools/types"
)

var (
	_ Model      = (*Log)(nil)
	_ DBExporter = (*Log)(nil)
)

const LogsTableName = "_logs"

const (
	defaultMaxLogDataSize    = 16 << 10 // ~16kb
	defaultMaxLogMessageSize = 8000
)

type Log struct {
	BaseModel

	Created types.DateTime     `db:"created" json:"created"`
	Data    types.JSONMap[any] `db:"data" json:"data"`
	Message string             `db:"message" json:"message"`
	Level   int                `db:"level" json:"level"`
}

func (l *Log) TableName() string {
	return LogsTableName
}

// DBExport prepares and exports the current log model for db persistence.
//
// It also truncates the log's message and data to ensure that it is
// under app.Settings().Logs.MaxDataSize.
func (l *Log) DBExport(app App) (map[string]any, error) {
	result := map[string]any{
		"id":      l.Id,
		"created": l.Created,
		"level":   l.Level,
	}

	// truncate the raw message bytes
	// (this is expected to be very rare so it is ok even if multi-byte chars)
	if int64(len(l.Message)) > defaultMaxLogMessageSize {
		result["message"] = l.Message[:defaultMaxLogMessageSize]
	} else {
		result["message"] = l.Message
	}

	// @todo once added in the standard library consider replacing with
	// WithByteLimit and WithDepthLimit as suggested in https://github.com/golang/go/issues/56733
	if len(l.Data) == 0 {
		result["data"] = l.Data
	} else {
		maxDataSize := app.Settings().Logs.MaxDataSize
		if maxDataSize == 0 {
			maxDataSize = defaultMaxLogDataSize
		}

		rawData, err := l.Data.MarshalJSON()
		if err != nil {
			return nil, err
		}

		if int64(len(rawData)) > maxDataSize {
			truncatedData := types.JSONMap[any]{}

			// ignore syntax errors in case of truncated incomplete json
			//
			// jsonv2 stream decodes and all "valid" attrs read up to the
			// invalid part will be populated in truncatedData
			_ = json.Unmarshal(rawData[:maxDataSize], &truncatedData)

			truncatedData["__pb_truncated__"] = true

			rawData, err = truncatedData.MarshalJSON()
			if err != nil {
				return nil, err
			}
		}

		result["data"] = types.JSONRaw(rawData)
	}

	return result, nil
}
