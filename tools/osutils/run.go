package osutils

import (
	"os"
	"strings"
)

var runDirs = []string{os.TempDir(), cacheDir()}

// IsProbablyGoRun loosely checks if the current program was started with "go run".
func IsProbablyGoRun() bool {
	for _, dir := range runDirs {
		if dir != "" && strings.HasPrefix(os.Args[0], dir) {
			return true
		}
	}

	return false
}

func cacheDir() string {
	dir := os.Getenv("GOCACHE")
	if dir == "off" {
		return ""
	}

	if dir == "" {
		dir, _ = os.UserCacheDir()
	}

	return dir
}
