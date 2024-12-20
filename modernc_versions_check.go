package pocketbase

import (
	"fmt"
	"log/slog"
	"runtime/debug"
)

const (
	expectedDriverVersion = "v1.34.3"
	expectedLibcVersion   = "v1.55.3"

	// ModerncDepsCheckHookId is the id of the hook that performs the modernc.org/* deps checks.
	// It could be used for removing/unbinding the hook if you don't want the checks.
	ModerncDepsCheckHookId = "pbModerncDepsCheck"
)

func checkModerncDeps(logger *slog.Logger) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return // no build info (probably compiled without module support)
	}

	var driverVersion, libcVersion string

	for _, dep := range info.Deps {
		switch dep.Path {
		case "modernc.org/libc":
			libcVersion = dep.Version
		case "modernc.org/sqlite":
			driverVersion = dep.Version
		}

		// no need to further search if both deps are located
		if driverVersion != "" && libcVersion != "" {
			break
		}
	}

	// not using the default driver
	if driverVersion == "" {
		return
	}

	if driverVersion != expectedDriverVersion {
		logger.Warn(fmt.Sprintf(
			"You are using  modernc.org/sqlite %s which differs from the tested %s.\n"+
				"Make sure to either manually update in your go.mod the dependency version to the expected one OR if you want to keep yours "+
				"ensure that its indirect modernc.org/libc dependency has the same version as in the https://gitlab.com/cznic/sqlite/-/blob/master/go.mod, "+
				"otherwise it could result in unexpected build or runtime errors.",
			driverVersion,
			expectedDriverVersion,
		), slog.String("current", driverVersion), slog.String("expected", expectedDriverVersion))
	} else if libcVersion != expectedLibcVersion {
		logger.Warn(fmt.Sprintf(
			"You are using a modernc.org/libc %s which differs from the tested %s.\n"+
				"Please update your go.mod and manually set modernc.org/libc to %s, otherwise it could result in unexpected build or runtime errors "+
				"(you may have to also run 'go clean -modcache' to clear the cache if the warning persists).",
			libcVersion,
			expectedLibcVersion,
			expectedLibcVersion,
		), slog.String("current", libcVersion), slog.String("expected", expectedLibcVersion))
	}
}
