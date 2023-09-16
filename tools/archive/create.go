package archive

import (
	"archive/zip"
	"compress/flate"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Create creates a new zip archive from src dir content and saves it in dest path.
//
// You can specify skipPaths to skip/ignore certain directories and files (relative to src)
// preventing adding them in the final archive.
func Create(src string, dest string, skipPaths ...string) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	zf, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	// register a custom Deflate compressor
	zw.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestSpeed)
	})

	if err := zipAddFS(zw, os.DirFS(src), skipPaths...); err != nil {
		// try to cleanup at least the created zip file
		os.Remove(dest)

		return err
	}

	return nil
}

// note remove after similar method is added in the std lib (https://github.com/golang/go/issues/54898)
func zipAddFS(w *zip.Writer, fsys fs.FS, skipPaths ...string) error {
	return fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// skip
		for _, ignore := range skipPaths {
			if ignore == name ||
				strings.HasPrefix(filepath.Clean(name)+string(os.PathSeparator), filepath.Clean(ignore)+string(os.PathSeparator)) {
				return nil
			}
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		h.Name = name
		h.Method = zip.Deflate

		fw, err := w.CreateHeader(h)
		if err != nil {
			return err
		}

		f, err := fsys.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(fw, f)

		return err
	})
}
