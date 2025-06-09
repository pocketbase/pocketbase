package pocketbase

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase/core"
)

const (
	expectedDriverVersion = "v1.38.0"
	expectedLibcVersion   = "v1.65.10"

	// ModerncDepsCheckHookId is the id of the hook that performs the modernc.org/* deps checks.
	// It could be used for removing/unbinding the hook if you don't want the checks.
	ModerncDepsCheckHookId = "pbModerncDepsCheck"
)

// checkModerncDeps checks whether the current binary was buit with the
// expected and tested modernc driver dependencies.
//
// This is needed because modernc.org/libc doesn't follow semantic versioning
// and using a version different from the one in the go.mod of modernc.org/sqlite
// could have unintended side-effects and cause obscure build and runtime bugs
// (https://github.com/pocketbase/pocketbase/issues/6136).
func checkModerncDeps(app core.App) {
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

	var msg string
	if driverVersion != expectedDriverVersion {
		msg = fmt.Sprintf(
			"You are using modernc.org/sqlite %s which differs from the expected and tested %s.\n"+
				"Make sure to either manually update in your go.mod the dependency version to the expected one OR if you want to keep yours "+
				"ensure that its indirect modernc.org/libc dependency has the same version as in the https://gitlab.com/cznic/sqlite/-/blob/master/go.mod, "+
				"otherwise it could result in unexpected build or runtime errors.",
			driverVersion,
			expectedDriverVersion,
		)
		app.Logger().Warn(msg, slog.String("current", driverVersion), slog.String("expected", expectedDriverVersion))
	} else if libcVersion != expectedLibcVersion {
		msg = fmt.Sprintf(
			"You are using modernc.org/libc %s which differs from the expected and tested %s.\n"+
				"Please update your go.mod and manually set modernc.org/libc to %s, otherwise it could result in unexpected build or runtime errors "+
				"(you may have to also run 'go clean -modcache' to clear the cache if the warning persists).",
			libcVersion,
			expectedLibcVersion,
			expectedLibcVersion,
		)
		app.Logger().Warn(msg, slog.String("current", libcVersion), slog.String("expected", expectedLibcVersion))
	}

	// ensure that the message is printed to the default stderr too
	// (when in dev mode this is not needed because we print all logs)
	if msg != "" && !app.IsDev() {
		color.Yellow("\nWARN " + msg + "\n\n")
	}
}
