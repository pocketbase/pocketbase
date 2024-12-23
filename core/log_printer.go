package core

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase/tools/logger"
	"github.com/pocketbase/pocketbase/tools/store"
	"github.com/spf13/cast"
)

var cachedColors = store.New[string, *color.Color](nil)

// getColor returns [color.Color] object and cache it (if not already).
func getColor(attrs ...color.Attribute) (c *color.Color) {
	cacheKey := fmt.Sprint(attrs)
	if c = cachedColors.Get(cacheKey); c == nil {
		c = color.New(attrs...)
		cachedColors.Set(cacheKey, c)
	}
	return
}

// printLog prints the provided log to the stderr.
// (note: defined as variable to overwriting in the tests)
var printLog = func(log *logger.Log) {
	var str strings.Builder

	switch log.Level {
	case slog.LevelDebug:
		str.WriteString(getColor(color.Bold, color.FgHiBlack).Sprint("DEBUG "))
		str.WriteString(getColor(color.FgWhite).Sprint(log.Message))
	case slog.LevelInfo:
		str.WriteString(getColor(color.Bold, color.FgWhite).Sprint("INFO "))
		str.WriteString(getColor(color.FgWhite).Sprint(log.Message))
	case slog.LevelWarn:
		str.WriteString(getColor(color.Bold, color.FgYellow).Sprint("WARN "))
		str.WriteString(getColor(color.FgYellow).Sprint(log.Message))
	case slog.LevelError:
		str.WriteString(getColor(color.Bold, color.FgRed).Sprint("ERROR "))
		str.WriteString(getColor(color.FgRed).Sprint(log.Message))
	default:
		str.WriteString(getColor(color.Bold, color.FgCyan).Sprintf("[%d] ", log.Level))
		str.WriteString(getColor(color.FgCyan).Sprint(log.Message))
	}

	str.WriteString("\n")

	if v, ok := log.Data["type"]; ok && cast.ToString(v) == "request" {
		padding := 0
		keys := []string{"error", "details"}
		for _, k := range keys {
			if v := log.Data[k]; v != nil {
				str.WriteString(getColor(color.FgHiRed).Sprintf("%s└─ %v", strings.Repeat(" ", padding), v))
				str.WriteString("\n")
				padding += 3
			}
		}
	} else if len(log.Data) > 0 {
		str.WriteString(getColor(color.FgHiBlack).Sprintf("└─ %v", log.Data))
		str.WriteString("\n")
	}

	fmt.Print(str.String())
}
