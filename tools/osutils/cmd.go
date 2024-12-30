package osutils

import (
	"os/exec"
	"runtime"

	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// LaunchURL attempts to open the provided url in the user's default browser.
//
// It is platform dependent and it uses:
//   - "open" on macOS
//   - "rundll32" on Windows
//   - "xdg-open" on everything else (Linux, FreeBSD, etc.)
func LaunchURL(url string) error {
	if err := is.URL.Validate(url); err != nil {
		return err
	}

	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		// not sure if this is the best command but seems to be the most reliable based on the comments in
		// https://stackoverflow.com/questions/3739327/launching-a-website-via-the-windows-commandline#answer-49115945
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}
