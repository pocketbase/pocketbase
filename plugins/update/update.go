// Package update implements a new command to replace the current
// PocketBase executable with the one from the latest BugStage releases.
//
// Note that it expects the fetched release to have file(s) in the following formats:
//   - "*_linux_amd64.zip"
//   - "*_linux_arm64.zip"
//   - "*_linux_armv7.zip"
//   - "*_darwin_amd64.zip"
//   - "*_darwin_arm64.zip"
//   - "*_windows_amd64.zip"
//   - "*_windows_arm64.zip"
//
// The appropriate file is auto downloaded based on the user platform.
// If there is no file matching the user platform then an error is returned.
//
// Example plugin registration:
//
//	update.MustRegister(app, app.RootCmd, update.Config{
//		BaseURL:    "https://ganigeorgiev.com",
//		Project:    "pocketbase",
//		Repository: "pocketbase",
//	})
package update

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/archive"
	"github.com/pocketbase/pocketbase/tools/osutils"
	"github.com/spf13/cobra"
)

// HttpClient is a base HTTP client interface (usually used for test purposes).
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Config defines the config options of the update plugin.
type Config struct {
	// BaseURL is the base URL of the BugStage host (or similar compatible)
	// used to fetch the latest releases information.
	//
	// This field is required!
	BaseURL string

	// Project specifies the project name of the repository.
	//
	// This field is required!
	Project string

	// Repository specifies the name of the repository.
	//
	// This field is required!
	Repository string

	// Channel specifies the name of the release channel to check.
	//
	// If not set uses the system "default" channel.
	Channel string

	// ArchiveExecutable specifies the name of the executable file in the release archive
	// (default to the same value as in [Config.Repository]; an additional ".exe" check is also performed as a fallback).
	ArchiveExecutable string

	// Optional context to use when fetching and downloading the latest release.
	Context context.Context

	// The HTTP client to use when fetching and downloading the latest release.
	// Defaults to [http.DefaultClient].
	HttpClient HttpClient
}

// MustRegister registers the update plugin to the provided app instance
// and panic if it fails.
func MustRegister(app core.App, rootCmd *cobra.Command, config Config) {
	if err := Register(app, rootCmd, config); err != nil {
		panic(err)
	}
}

// Register registers the update plugin to the provided app instance.
func Register(app core.App, rootCmd *cobra.Command, config Config) error {
	if config.BaseURL == "" {
		return errors.New("update.Config.BaseURL is required and needs to point to a BugStage instance")
	}

	if config.Project == "" {
		return errors.New("update.Config.Project is required")
	}

	if config.Repository == "" {
		return errors.New("update.Config.Repository is required")
	}

	p := &plugin{
		app:            app,
		currentVersion: rootCmd.Version,
		config:         config,
	}

	p.config.BaseURL = strings.TrimRight(p.config.BaseURL, "/")

	if p.config.ArchiveExecutable == "" {
		p.config.ArchiveExecutable = p.config.Repository
	}

	if p.config.Channel == "" {
		p.config.Channel = "default"
	}

	if p.config.HttpClient == nil {
		p.config.HttpClient = http.DefaultClient
	}

	if p.config.Context == nil {
		p.config.Context = context.Background()
	}

	rootCmd.AddCommand(p.updateCmd())

	return nil
}

type plugin struct {
	app            core.App
	config         Config
	currentVersion string
}

func (p *plugin) updateCmd() *cobra.Command {
	var withBackup bool

	command := &cobra.Command{
		Use:          "update",
		Short:        "Automatically updates the current app executable with the latest available version",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			var needConfirm bool
			if isMaybeRunningInDocker() {
				needConfirm = true
				color.Yellow("NB! It seems that you are in a Docker container.")
				color.Yellow("The update command may not work as expected in this context because usually the version of the app is managed by the container image itself.")
			} else if isMaybeRunningInNixOS() {
				needConfirm = true
				color.Yellow("NB! It seems that you are in a NixOS.")
				color.Yellow("Due to the non-standard filesystem implementation of the environment, the update command may not work as expected.")
			}

			if needConfirm {
				confirm := osutils.YesNoPrompt("Do you want to proceed with the update?", false)
				if !confirm {
					fmt.Println("The command has been cancelled.")
					return nil
				}
			}

			return p.update(withBackup)
		},
	}

	command.PersistentFlags().BoolVar(
		&withBackup,
		"backup",
		true,
		"Creates a pb_data backup at the end of the update process",
	)

	return command
}

func (p *plugin) update(withBackup bool) error {
	color.Yellow("Fetching release information...")

	latest, err := p.fetchLatestRelease()
	if err != nil {
		return err
	}

	if compareVersions(strings.TrimPrefix(p.currentVersion, "v"), strings.TrimPrefix(latest.Tag, "v")) <= 0 {
		color.Green("You already have the latest version %s.", p.currentVersion)
		return nil
	}

	suffix := archiveSuffix(runtime.GOOS, runtime.GOARCH)
	if suffix == "" {
		return errors.New("unsupported platform")
	}

	filename, err := latest.findFileBySuffix(suffix)
	if err != nil {
		return err
	}

	fileURL := p.config.BaseURL + latest.buildFilePath(filename)

	releaseDir := filepath.Join(p.app.DataDir(), core.LocalTempDirName)
	defer os.RemoveAll(releaseDir)

	color.Yellow("Downloading %s...", filename)

	// download the release asset
	assetZip := filepath.Join(releaseDir, filename)
	if err := downloadFile(p.config.Context, p.config.HttpClient, fileURL, assetZip); err != nil {
		return err
	}

	color.Yellow("Extracting %s...", filename)

	extractDir := filepath.Join(releaseDir, "extracted_"+filename)
	defer os.RemoveAll(extractDir)

	if err := archive.Extract(assetZip, extractDir); err != nil {
		return err
	}

	color.Yellow("Replacing the executable...")

	oldExec, err := os.Executable()
	if err != nil {
		return err
	}
	renamedOldExec := oldExec + ".old"
	defer os.Remove(renamedOldExec)

	newExec := filepath.Join(extractDir, p.config.ArchiveExecutable)
	if _, err := os.Stat(newExec); err != nil {
		// try again with an .exe extension
		newExec = newExec + ".exe"
		if _, fallbackErr := os.Stat(newExec); fallbackErr != nil {
			return fmt.Errorf("the executable in the extracted path is missing or it is inaccessible: %v, %v", err, fallbackErr)
		}
	}

	// rename the current executable
	if err := os.Rename(oldExec, renamedOldExec); err != nil {
		return fmt.Errorf("failed to rename the current executable: %w", err)
	}

	tryToRevertExecChanges := func() {
		if revertErr := os.Rename(renamedOldExec, oldExec); revertErr != nil {
			p.app.Logger().Debug(
				"Failed to revert executable",
				slog.String("old", renamedOldExec),
				slog.String("new", oldExec),
				slog.String("error", revertErr.Error()),
			)
		}
	}

	// replace with the extracted binary
	if err := os.Rename(newExec, oldExec); err != nil {
		tryToRevertExecChanges()
		return fmt.Errorf("failed replacing the executable: %w", err)
	}

	// backup existing pb_data
	if withBackup {
		color.Yellow("Creating pb_data backup...")

		backupName := fmt.Sprintf("@update_%s.zip", latest.Tag)
		if err := p.app.CreateBackup(p.config.Context, backupName); err != nil {
			tryToRevertExecChanges()
			return err
		}
	}

	// print the release notes (if any)
	if latest.Content != "" {
		fmt.Print("\n")
		color.Cyan("%s release notes:", latest.Tag)
		color.Cyan(strings.TrimSpace(latest.Content))
		fmt.Print("\n")
	}

	color.HiBlack("---")
	fmt.Print("\n")
	color.Green("Update completed successfully! You can start the executable as usual.")
	fmt.Print("\n")

	return nil
}

func (p *plugin) fetchLatestRelease() (*release, error) {
	queryParams := url.Values{}
	queryParams.Set("perPage", "1")
	queryParams.Set("skipTotal", "1")
	queryParams.Set("sort", "-semver")
	queryParams.Set("filter", fmt.Sprintf(
		"channel.name=%q && channel.repository.name=%q && channel.repository.project.name=%q",
		p.config.Channel,
		p.config.Repository,
		p.config.Project,
	))

	releaseURL := fmt.Sprintf("%s/api/collections/releases/records?%s", p.config.BaseURL, queryParams.Encode())

	req, err := http.NewRequestWithContext(p.config.Context, "GET", releaseURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := p.config.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// http.Client doesn't treat non 2xx responses as error
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"(%d) failed to fetch latest releases:\n%s",
			res.StatusCode,
			string(rawBody),
		)
	}

	result := struct {
		Items []*release `json:"items"`
	}{}
	if err := json.Unmarshal(rawBody, &result); err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errors.New("no releases found")
	}

	return result.Items[0], nil
}

func downloadFile(
	ctx context.Context,
	client HttpClient,
	url string,
	destPath string,
) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// http.Client doesn't treat non 2xx responses as error
	if res.StatusCode >= 400 {
		return fmt.Errorf("(%d) failed to send download file request", res.StatusCode)
	}

	// ensure that the dest parent dir(s) exist
	if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
		return err
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, res.Body); err != nil {
		return err
	}

	return nil
}

func archiveSuffix(goos, goarch string) string {
	switch goos {
	case "linux":
		switch goarch {
		case "amd64":
			return "_linux_amd64.zip"
		case "arm64":
			return "_linux_arm64.zip"
		case "arm":
			return "_linux_armv7.zip"
		}
	case "darwin":
		switch goarch {
		case "amd64":
			return "_darwin_amd64.zip"
		case "arm64":
			return "_darwin_arm64.zip"
		}
	case "windows":
		switch goarch {
		case "amd64":
			return "_windows_amd64.zip"
		case "arm64":
			return "_windows_arm64.zip"
		}
	}

	return ""
}

func compareVersions(a, b string) int {
	aSplit := strings.Split(a, ".")
	aTotal := len(aSplit)

	bSplit := strings.Split(b, ".")
	bTotal := len(bSplit)

	limit := aTotal
	if bTotal > aTotal {
		limit = bTotal
	}

	for i := 0; i < limit; i++ {
		var x, y int

		if i < aTotal {
			x, _ = strconv.Atoi(aSplit[i])
		}

		if i < bTotal {
			y, _ = strconv.Atoi(bSplit[i])
		}

		if x < y {
			return 1 // b is newer
		}

		if x > y {
			return -1 // a is newer
		}
	}

	return 0 // equal
}

// note: not completely reliable as it may not work on all platforms
// but should at least provide a warning for the most common use cases
func isMaybeRunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

// note: untested
func isMaybeRunningInNixOS() bool {
	_, err := os.Stat("/etc/NIXOS")
	return err == nil
}
