// Package ghupdate implements a new command to selfupdate the current
// PocketBase executable with the latest GitHub release.
package ghupdate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/archive"
	"github.com/spf13/cobra"
)

// HttpClient is a base HTTP client interface (usually used for test purposes).
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Options defines optional struct to customize the default plugin behavior.
//
// NB! This plugin is considered experimental and its options may change in the future.
type Options struct {
	// Owner specifies the account owner of the repository (default to "pocketbase").
	Owner string

	// Repo specifies the name of the repository (default to "pocketbase").
	Repo string

	// ArchiveExecutable specifies the name of the executable file in the release archive (default to "pocketbase").
	ArchiveExecutable string

	// Optional context to use when fetching and downloading the latest release.
	Context context.Context

	// The HTTP client to use when fetching and downloading the latest release.
	// Defaults to `http.DefaultClient`.
	HttpClient HttpClient
}

// MustRegister registers the ghupdate plugin to the provided app instance
// and panic if it fails.
func MustRegister(app core.App, rootCmd *cobra.Command, options *Options) {
	if err := Register(app, rootCmd, options); err != nil {
		panic(err)
	}
}

// Register registers the ghupdate plugin to the provided app instance.
func Register(app core.App, rootCmd *cobra.Command, options *Options) error {
	p := &plugin{
		app:            app,
		currentVersion: rootCmd.Version,
	}

	if options != nil {
		p.options = options
	} else {
		p.options = &Options{}
	}

	if p.options.Owner == "" {
		p.options.Owner = "pocketbase"
	}

	if p.options.Repo == "" {
		p.options.Repo = "pocketbase"
	}

	if p.options.ArchiveExecutable == "" {
		p.options.ArchiveExecutable = "pocketbase"
	}

	if p.options.HttpClient == nil {
		p.options.HttpClient = http.DefaultClient
	}

	if p.options.Context == nil {
		p.options.Context = context.Background()
	}

	rootCmd.AddCommand(p.updateCmd())

	return nil
}

type plugin struct {
	app            core.App
	currentVersion string
	options        *Options
}

func (p *plugin) updateCmd() *cobra.Command {
	var withBackup bool

	command := &cobra.Command{
		Use:   "update",
		Short: "Automatically updates the current PocketBase executable with the latest available version",
		// @todo remove after logs generalization
		// prevents printing the error log twice
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(command *cobra.Command, args []string) error {
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

	latest, err := fetchLatestRelease(
		p.options.Context,
		p.options.HttpClient,
		p.options.Owner,
		p.options.Repo,
	)
	if err != nil {
		return err
	}

	if compareVersions(strings.TrimPrefix(p.currentVersion, "v"), strings.TrimPrefix(latest.Tag, "v")) <= 0 {
		color.Green("You already have the latest PocketBase %s.", p.currentVersion)
		return nil
	}

	suffix := archiveSuffix(runtime.GOOS, runtime.GOARCH)
	if suffix == "" {
		return errors.New("unsupported platform")
	}

	asset, err := latest.findAssetBySuffix(suffix)
	if err != nil {
		return err
	}

	releaseDir := path.Join(p.app.DataDir(), core.LocalTempDirName)
	defer os.RemoveAll(releaseDir)

	color.Yellow("Downloading %s...", asset.Name)

	// download the release asset
	assetZip := path.Join(releaseDir, asset.Name)
	if err := downloadFile(p.options.Context, p.options.HttpClient, asset.DownloadUrl, assetZip); err != nil {
		return err
	}

	color.Yellow("Extracting %s...", asset.Name)

	extractDir := path.Join(releaseDir, "extracted_"+asset.Name)
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

	newExec := path.Join(extractDir, p.options.ArchiveExecutable)

	// rename the current executable
	if err := os.Rename(oldExec, renamedOldExec); err != nil {
		return fmt.Errorf("Failed to rename the current executable: %w", err)
	}

	tryToRevertExecChanges := func() {
		if revertErr := os.Rename(renamedOldExec, oldExec); revertErr != nil && p.app.IsDebug() {
			log.Println(revertErr)
		}
	}

	// replace with the extracted binary
	if err := os.Rename(newExec, oldExec); err != nil {
		tryToRevertExecChanges()
		return fmt.Errorf("Failed replacing the executable: %w", err)
	}

	if withBackup {
		color.Yellow("Creating pb_data backup...")

		backupName := fmt.Sprintf("@update_%s.zip", latest.Tag)
		if err := p.app.CreateBackup(p.options.Context, backupName); err != nil {
			tryToRevertExecChanges()
			return err
		}
	}

	color.HiBlack("---")
	color.Green("Update completed sucessfully! You can start the executable as usual.")

	return nil
}

func fetchLatestRelease(
	ctx context.Context,
	client HttpClient,
	owner string,
	repo string,
) (*release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
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

	result := &release{}
	if err := json.Unmarshal(rawBody, result); err != nil {
		return nil, err
	}

	return result, nil
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
