package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract extracts the zip archive at src to dest.
func Extract(src, dest string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func(zr *zip.ReadCloser) {
		_ = zr.Close()
	}(zr)

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
	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(r)

	if zipFile.FileInfo().IsDir() {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	} else {
		// ensure that the file path directories are created
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		_, err = io.Copy(f, r)
		if err != nil {
			return err
		}
	}

	return nil
}
