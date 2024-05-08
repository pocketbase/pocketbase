package models

import "github.com/lilysnc/pocketbasepg/tools/types"

type BackupFileInfo struct {
	Key      string         `json:"key"`
	Size     int64          `json:"size"`
	Modified types.DateTime `json:"modified"`
}
