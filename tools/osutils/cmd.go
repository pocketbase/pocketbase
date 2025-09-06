package osutils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

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

// YesNoPrompt performs a console prompt that asks the user for Yes/No answer.
//
// If the user just press Enter (aka. doesn't type anything) it returns the fallback value.
func YesNoPrompt(message string, fallback bool) bool {
	options := "Y/n"
	if !fallback {
		options = "y/N"
	}

	r := bufio.NewReader(os.Stdin)

	var s string
	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", message, options)

		s, _ = r.ReadString('\n')

		s = strings.ToLower(strings.TrimSpace(s))

		switch s {
		case "":
			return fallback
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
	}
}
