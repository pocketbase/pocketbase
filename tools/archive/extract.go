package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract extracts the zip archive at "src" to "dest".
//
// Note that only dirs and regular files will be extracted.
// Symbolic links, named pipes, sockets, or any other irregular files
// are skipped because they come with too many edge cases and ambiguities.
func Extract(src, dest string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zr.Close()

	// normalize dest path to check later for Zip Slip
	dest = filepath.Clean(dest) + string(os.PathSeparator)

	for _, f := range zr.File {
		err := extractFile(f, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractFile extracts the provided zipFile into "basePath/zipFileName" path,
// creating all the necessary path directories.
func extractFile(zipFile *zip.File, basePath string) error {
	path := filepath.Join(basePath, zipFile.Name)

	// check for Zip Slip
	if !strings.HasPrefix(path, basePath) {
		return fmt.Errorf("invalid file path: %s", path)
	}

	r, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	// allow only dirs or regular files
	if zipFile.FileInfo().IsDir() {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	} else if zipFile.FileInfo().Mode().IsRegular() {
		// ensure that the file path directories are created
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, r)
		if err != nil {
			return err
		}
	}

	return nil
}
