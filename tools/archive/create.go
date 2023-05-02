package archive

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
)

// Create creates a new zip archive from src dir content and saves it in dest path.
func Create(src, dest string) error {
	zf, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	if err := zipAddFS(zw, os.DirFS(src)); err != nil {
		// try to cleanup the created zip file
		os.Remove(dest)

		return err
	}

	return nil
}

// note remove after similar method is added in the std lib (https://github.com/golang/go/issues/54898)
func zipAddFS(w *zip.Writer, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
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
