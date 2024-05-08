package logger

import (
	"log/slog"
	"time"

	"github.com/lilysnc/pocketbasepg/tools/types"
)

// Log is similar to [slog.Record] bit contains the log attributes as
// preformatted JSON map.
type Log struct {
	Time    time.Time
	Message string
	Level   slog.Level
	Data    types.JsonMap
}
